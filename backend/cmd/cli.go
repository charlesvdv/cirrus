package main

import (
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func runMainApp(args []string) {
	conf := appConfig{}

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
		},
		Action: func(c *cli.Context) error {
			mainApp(conf)
			return nil
		},
	}

	err := app.Run(args)
	if err != nil {
		log.Fatal().Err(err).Msg("Application failed unexpectedly")
	}
}
