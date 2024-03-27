package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func setupRoutes(webApp *fiber.App) {
	webApp.Post("/ON", lightsOn)
	webApp.Post("/OFF", lightsOff)
}

func lightsOff(c *fiber.Ctx) error {
	room := string(c.Body())
	log.Info().Str("room", room).Msg("lights off...")
	return c.Status(http.StatusOK).Send([]byte("OK"))
}

func lightsOn(c *fiber.Ctx) error {
	room := string(c.Body())
	log.Info().Str("room", room).Msg("lights on...")
	return c.Status(http.StatusOK).Send([]byte("OK"))
}
