package main

import (
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func runMainApp(args []string, main func(appConfig)) {
	conf := appConfig{}

	postgresPort := uint(0)

	app := &cli.App{
		Name:  "cirrus",
		Usage: "A modern file storage service",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Usage:       "HTTP server port",
				Value:       8000,
				Aliases:     []string{"p"},
				EnvVars:     []string{"CIRRUS_PORT"},
				Destination: &conf.port,
			},
			&cli.StringFlag{
				Name:        "hostname",
				Usage:       "HTTP server hostname",
				Value:       "",
				EnvVars:     []string{"CIRRUS_HOSTNAME"},
				Destination: &conf.hostname,
			},
			&cli.StringFlag{
				Name:        "database-type",
				Usage:       "Database type",
				Value:       "postgres",
				EnvVars:     []string{"CIRRUS_DATABASE"},
				Destination: &conf.databaseType,
			},
			&cli.StringFlag{
				Name:        "postgres-host",
				Usage:       "Host for postgres database",
				Value:       "localhost",
				EnvVars:     []string{"POSTGRES_HOST"},
				Destination: &conf.postgresHost,
			},
			&cli.UintFlag{
				Name:        "postgres-port",
				Usage:       "Port for postgres database",
				Value:       5432,
				EnvVars:     []string{"POSTGRES_PORT"},
				Destination: &postgresPort,
			},
			&cli.StringFlag{
				Name:        "postgres-user",
				Usage:       "User for postgres database",
				Value:       "postgres",
				EnvVars:     []string{"POSTGRES_USER"},
				Destination: &conf.postgresUser,
			},
			&cli.StringFlag{
				Name:        "postgres-password",
				Usage:       "Password for postgres database",
				EnvVars:     []string{"POSTGRES_PASSWORD"},
				Destination: &conf.postgresPassword,
			},
			&cli.StringFlag{
				Name:        "postgres-database",
				Usage:       "Database name for postgres database",
				Value:       "cirrus",
				EnvVars:     []string{"POSTGRES_DB"},
				Destination: &conf.postgresDatabase,
			},
			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "Debug mode",
				Value:       false,
				EnvVars:     []string{"CIRRUS_DEBUG"},
				Destination: &conf.debug,
			},
		},
		Action: func(c *cli.Context) error {
			conf.postgresPort = uint16(postgresPort)

			main(conf)
			return nil
		},
	}

	err := app.Run(args)
	if err != nil {
		log.Fatal().Err(err).Msg("Application failed unexpectedly")
	}
}
