package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/charlesvdv/cirrus/backend/api/rest"
	"github.com/charlesvdv/cirrus/backend/pkg/user"
)

func main() {
	runMainApp(os.Args)
}

type appConfig struct {
	port     int
	hostname string
}

func mainApp(conf appConfig) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	rootHandler := rest.NewRootHandler()
	rootHandler.Register(
		rest.NewUserHandler(&user.UserManager{}),
	)

	address := fmt.Sprintf("%v:%v", conf.hostname, conf.port)
	log.Info().Msgf("Listening on %v", address)
	if err := http.ListenAndServe(address, rootHandler.Get()); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
		os.Exit(1)
	}
}
