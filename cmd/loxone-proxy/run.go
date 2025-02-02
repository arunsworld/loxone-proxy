package main

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func run(ctx *cli.Context) error {
	webApp := fiber.New(fiber.Config{
		Prefork: false,
	})
	go appShutdownOnCtxCancel(ctx.Context, webApp)
	webApp.Use(logger.New(), etag.New(), compress.New())

	setupRoutes(webApp)

	port := ctx.Int("port")
	hostname := ctx.String("hostname")
	log.Debug().Int("port", port).Str("hostname", hostname).Msg("serving webserver...")
	return webApp.Listen(fmt.Sprintf("%s:%d", hostname, port))
}

func appShutdownOnCtxCancel(ctx context.Context, app *fiber.App) {
	<-ctx.Done()
	if err := app.Shutdown(); err != nil {
		log.Printf("error while shutting down app: %v", err)
	}
}
