package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type routes struct {
	webApp *fiber.App
	mc     *matterController
	pc     *pushoverController
	nc     nodeRedClient
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
	case "livingRoom":
	case "sittingRoom":
		r.mc.switchLightsOff(sittingRoomLight01)
		if err := r.nc.turnOff("sitting-room-lady"); err != nil {
			log.Printf("err turning off sitting-room-lady: %v", err)
			return err
		} else {
			log.Printf("turned off sitting-room-lady")
		}
	case "diningRoom":
		if err := r.nc.turnOff("dining-room-left"); err != nil {
			log.Printf("err turning off dining-room-left: %v", err)
			return err
		} else {
			log.Printf("turned dining-room-lef on")
		}
	case "theDen":
	case "doorbellChime":
		r.mc.switchLightsOff(doorbellChime)
	}
	return c.Status(http.StatusOK).Send([]byte("OK"))
}

func (r *routes) lightsOn(c *fiber.Ctx) error {
	room := string(c.Body())
	log.Info().Str("room", room).Msg("lights on...")
	switch room {
	case "livingRoom":
	case "sittingRoom":
		r.mc.switchLightsOn(sittingRoomLight01)
		if err := r.nc.turnOn("sitting-room-lady"); err != nil {
			log.Printf("err turning on sitting-room-lady: %v", err)
			return err
		} else {
			log.Printf("turned on sitting-room-lady")
		}
	case "diningRoom":
		if err := r.nc.turnOn("dining-room-left"); err != nil {
			log.Printf("err turning on dining-room-left: %v", err)
			return err
		} else {
			log.Printf("turned dining-room-lef on")
		}
	case "theDen":
	case "doorbellChime":
		r.mc.switchLightsOn(doorbellChime)
		r.pc.send("Someone rang the doorbell")
	}
	return c.Status(http.StatusOK).Send([]byte("OK"))
}

func (r *routes) boilerOn(c *fiber.Ctx) error {
	log.Info().Msg("boiler is triggered on...")
	r.pc.send("Boiler has turned on")
	return c.Status(http.StatusOK).Send([]byte("OK"))
}

func (r *routes) boilerOff(c *fiber.Ctx) error {
	log.Info().Msg("boiler is triggered off...")
	r.pc.send("Boiler has turned off")
	return c.Status(http.StatusOK).Send([]byte("OK"))
}

func (r *routes) burglarAlarmArmed(c *fiber.Ctx) error {
	log.Info().Msg("burglar alarm is armed...")
	r.pc.send("Burglar alarm armed")
	return c.Status(http.StatusOK).Send([]byte("OK"))
}

func (r *routes) burglarAlarmDisarmed(c *fiber.Ctx) error {
	log.Info().Msg("burglar alarm is disarmed...")
	r.pc.send("Burglar alarm disarmed")
	return c.Status(http.StatusOK).Send([]byte("OK"))
}

func (r *routes) burglarAlarmTriggered(c *fiber.Ctx) error {
	log.Info().Msg("burglar alarm has triggered...")
	r.pc.send("Burglar alarm triggered")
	return c.Status(http.StatusOK).Send([]byte("OK"))
}

func (r *routes) burglarAlarmDiffused(c *fiber.Ctx) error {
	log.Info().Msg("burglar alarm diffused...")
	r.pc.send("Burglar alarm diffused")
	return c.Status(http.StatusOK).Send([]byte("OK"))
}
