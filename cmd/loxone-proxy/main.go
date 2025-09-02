package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var appName = "loxone proxy"

func main() {
	app := &cli.App{
		Name:  appName,
		Flags: createFlags(),
		Action: func(ctx *cli.Context) error {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			if ctx.Bool("debug") {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
				log.Debug().Msg("debug mode turned on")
			}

			return run(ctx)
		},
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02T15:04:05.999Z07:00",
	})
	log.Info().Msgf("%s starting...", appName)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	if err := app.RunContext(ctx, os.Args); err != nil {
		cancel()
		log.Fatal().Msg(err.Error())
	}
	cancel()
	log.Info().Msgf("%s terminated...", appName)
}

func createFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "hostname", Usage: "hostname (blank OK, or localhost, or IP address)", Value: "", EnvVars: []string{"HOST"}},
		&cli.IntFlag{Name: "port", Value: 6160, EnvVars: []string{"PORT"}},
		&cli.BoolFlag{Name: "debug", Usage: "turn on debug mode", Aliases: []string{"d"}, EnvVars: []string{"DEBUG"}},
		&cli.StringFlag{Name: "rpi", Usage: "rpi hostname", Value: "", EnvVars: []string{"RPI"}},
	}
}
