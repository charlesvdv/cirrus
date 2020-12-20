package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/charlesvdv/cirrus/backend/api/rest"
	"github.com/charlesvdv/cirrus/backend/db"
	"github.com/charlesvdv/cirrus/backend/pkg/user"
)

func main() {
	runMainApp(os.Args, mainApp)
}

type appConfig struct {
	port             int
	debug            bool
	hostname         string
	databaseType     string
	postgresHost     string
	postgresPort     uint16
	postgresUser     string
	postgresPassword string
	postgresDatabase string
}

func (c appConfig) validate() error {
	if c.databaseType != "postgres" {
		return fmt.Errorf("Invalid database type '%s'. Expected 'postgres'", c.databaseType)
	}

	return nil
}

func (c appConfig) postgresConfig() db.PostgresConfig {
	return db.PostgresConfig{
		Port:     c.postgresPort,
		Host:     c.postgresHost,
		User:     c.postgresUser,
		Password: c.postgresPassword,
		Database: c.postgresDatabase,
	}
}

func mainApp(conf appConfig) {
	if err := conf.validate(); err != nil {
		log.Fatal().Err(err).Msg("Invalid configuration")
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	database, err := db.NewPostgresDatabase(conf.postgresConfig())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	err = database.UpdateSchemas()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to update schema")
	}

	userService := user.NewUserService(&database, &user.PostgresRepository{})
	rootHandler := rest.NewRootHandler()
	rootHandler.Register(
		rest.NewUserHandler(&userService),
	)

	address := fmt.Sprintf("%v:%v", conf.hostname, conf.port)
	log.Info().Msgf("Listening on %v", address)
	if err := http.ListenAndServe(address, rootHandler.Get()); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
