package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type routes struct {
	webApp *fiber.App
	mc     *matterController
}

func (r *routes) setup() {
	r.webApp.Post("/ON", r.lightsOn)
	r.webApp.Post("/OFF", r.lightsOff)
}

func (r *routes) lightsOff(c *fiber.Ctx) error {
	room := string(c.Body())
	log.Info().Str("room", room).Msg("lights off...")
	switch room {
	case "sittingRoom":
		r.mc.switchLightsOff(sittingRoomLight01)
	}
	return c.Status(http.StatusOK).Send([]byte("OK"))
}

func (r *routes) lightsOn(c *fiber.Ctx) error {
	room := string(c.Body())
	log.Info().Str("room", room).Msg("lights on...")
	switch room {
	case "sittingRoom":
		r.mc.switchLightsOn(sittingRoomLight01)
	}
	return c.Status(http.StatusOK).Send([]byte("OK"))
}
