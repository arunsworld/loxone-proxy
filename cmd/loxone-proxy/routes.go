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
	// lights
	r.webApp.Post("/ON", r.lightsOn)
	r.webApp.Post("/OFF", r.lightsOff)
	// boiler
	r.webApp.Post("/BOILER_ON", r.boilerOn)
	r.webApp.Post("/BOILER_OFF", r.boilerOff)
	// burglar alarm arming
	r.webApp.Post("/BURGLAR_ALARM_ARMED", r.burglarAlarmArmed)
	r.webApp.Post("/BURGLAR_ALARM_DISARMED", r.burglarAlarmDisarmed)
	// burglar alarm trigger
	r.webApp.Post("/BURGLAR_ALARM_TRIGGERED", r.burglarAlarmTriggered)
	r.webApp.Post("/BURGLAR_ALARM_DIFFUSED", r.burglarAlarmDiffused)
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

func (r *routes) boilerOn(c *fiber.Ctx) error {
	log.Info().Msg("boiler is triggered on...")
	return c.Status(http.StatusOK).Send([]byte("OK"))
}

func (r *routes) boilerOff(c *fiber.Ctx) error {
	log.Info().Msg("boiler is triggered off...")
	return c.Status(http.StatusOK).Send([]byte("OK"))
}

func (r *routes) burglarAlarmArmed(c *fiber.Ctx) error {
	log.Info().Msg("burglar alarm is armed...")
	return c.Status(http.StatusOK).Send([]byte("OK"))
}

func (r *routes) burglarAlarmDisarmed(c *fiber.Ctx) error {
	log.Info().Msg("burglar alarm is disarmed...")
	return c.Status(http.StatusOK).Send([]byte("OK"))
}

func (r *routes) burglarAlarmTriggered(c *fiber.Ctx) error {
	log.Info().Msg("burglar alarm has triggered...")
	return c.Status(http.StatusOK).Send([]byte("OK"))
}

func (r *routes) burglarAlarmDiffused(c *fiber.Ctx) error {
	log.Info().Msg("burglar alarm diffused...")
	return c.Status(http.StatusOK).Send([]byte("OK"))
}
