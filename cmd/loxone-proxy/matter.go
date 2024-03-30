package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/tom-code/gomat"
	"github.com/tom-code/gomat/symbols"
)

const defaultMatterPort = 5540

type device struct {
	name     string
	deviceID uint64
	deviceIP string
	port     int // typically 5540 (udp port)
}

type job struct {
	device device
	state  uint32
}

type matterController struct {
	mu                  sync.RWMutex
	devicesReady        chan struct{}
	fabric              *gomat.Fabric
	secureChannels      map[string]*gomat.SecureChannel
	jobQueue            chan job
	stopProcessingQueue context.CancelFunc
}

const fabricID uint64 = 0x110
const adminUser uint64 = 100

func newMatterController(devices []device) (*matterController, error) {
	cm := gomat.NewFileCertManager(fabricID)
	if err := cm.Load(); err != nil {
		return nil, err
	}
	fabric := gomat.NewFabric(fabricID, cm)
	mc := &matterController{
		fabric:         fabric,
		jobQueue:       make(chan job, 500),
		devicesReady:   make(chan struct{}),
		secureChannels: make(map[string]*gomat.SecureChannel),
	}

	go mc.setupDevices(devices)

	ctx, stopProcessingQueue := context.WithCancel(context.Background())
	mc.stopProcessingQueue = stopProcessingQueue
	go mc.processQueue(ctx)

	return mc, nil
}

func (mc *matterController) setupDevices(devices []device) {
	for _, d := range devices {
		if err := mc.setupDevice(d); err != nil {
			log.Error().Err(err).Str("device", d.name).Msg("unable to setup device")
		}
	}
	close(mc.devicesReady)
	log.Info().Msg("all devices are now ready")
}

func (mc *matterController) processQueue(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("matterController processQueue terminating")
		case job := <-mc.jobQueue:
			select {
			case <-mc.devicesReady:
				if err := mc.switchLights(job.device, job.state); err != nil {
					log.Error().Err(err).Str("device", job.device.name).Msg("unable to switch device")
				}
			case <-ctx.Done():
				log.Debug().Msg("matterController processQueue terminating with a job in hand but devices not yet ready")
			}
		}
	}
}

func (mc *matterController) setupDevice(d device) error {
	mc.mu.RLock()
	_, ok := mc.secureChannels[d.name]
	if ok {
		mc.mu.RUnlock()
		log.Info().Str("device", d.name).Msg("device already setup")
		return nil
	}
	mc.mu.RUnlock()

	port := d.port
	if port == 0 {
		port = defaultMatterPort
	}
	secureChannel, err := gomat.StartSecureChannel(net.ParseIP(d.deviceIP), port, 0)
	if err != nil {
		return fmt.Errorf("unable to start secure channel for %s to %s: %w", d.name, d.deviceIP, err)
	}
	secureChannel, err = gomat.SigmaExchange(mc.fabric, adminUser, d.deviceID, secureChannel)
	if err != nil {
		secureChannel.Close()
		return fmt.Errorf("unable to sigmaExchange on channel for %s to %s: %w", d.name, d.deviceIP, err)
	}

	mc.mu.Lock()
	defer mc.mu.Unlock()

	if _, ok := mc.secureChannels[d.name]; ok {
		log.Info().Str("device", d.name).Msg("device already setup")
		return nil
	}
	mc.secureChannels[d.name] = &secureChannel
	log.Info().Str("device", d.name).Msg("device is setup")
	return nil
}

func (mc *matterController) removeDevice(d device) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	sc, ok := mc.secureChannels[d.name]
	if !ok {
		log.Info().Str("device", d.name).Msg("device asked to be removed, but it doesn't exist")
		return
	}
	sc.Close()
	delete(mc.secureChannels, d.name)
	log.Info().Str("device", d.name).Msg("device is removed")
}

func (mc *matterController) switchLightsOn(device device) {
	mc.jobQueue <- job{
		device: device,
		state:  symbols.COMMAND_ID_OnOff_On,
	}
}

func (mc *matterController) switchLightsOff(device device) {
	mc.jobQueue <- job{
		device: device,
		state:  symbols.COMMAND_ID_OnOff_Off,
	}
}

func (mc *matterController) switchLights(device device, command uint32) error {
	mc.mu.RLock()
	secureChannel, ok := mc.secureChannels[device.name]
	mc.mu.RUnlock()

	if !ok {
		if err := mc.setupDevice(device); err != nil {
			return err
		}
		mc.mu.RLock()
		secureChannel, ok = mc.secureChannels[device.name]
		mc.mu.RUnlock()
		if !ok {
			return fmt.Errorf("BUG???")
		}
	}

	if err := mc.trySwitchingLights(secureChannel, command); err == nil {
		return nil
	} else {
		log.Error().Err(err).Str("device", device.name).Msg("error switching lights - will attempt again")
	}
	mc.removeDevice(device)
	if err := mc.setupDevice(device); err != nil {
		log.Error().Err(err).Str("device", device.name).Msg("error setting up device")
		return err
	}

	mc.mu.RLock()
	secureChannel, ok = mc.secureChannels[device.name]
	mc.mu.RUnlock()
	if !ok {
		return fmt.Errorf("BUG 2???")
	}

	return mc.trySwitchingLights(secureChannel, command)
}

func (mc *matterController) trySwitchingLights(secureChannel *gomat.SecureChannel, command uint32) error {
	to_send := gomat.EncodeIMInvokeRequest(1,
		symbols.CLUSTER_ID_OnOff,
		command,
		[]byte{}, false, uint16(rand.Intn(0xffff)))
	secureChannel.Send(to_send)

	resp, err := secureChannel.Receive()
	if err != nil {
		return fmt.Errorf("error receiving from channel: %w", err)
	}
	status, err := resp.Tlv.GetIntRec([]int{1, 0, 1, 1, 0})
	if err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}
	if status != 0 {
		return fmt.Errorf("status received was %d instead of 0", status)
	}
	return nil
}

func (mc *matterController) close() {
	mc.devicesReady = make(chan struct{})
	mc.stopProcessingQueue()

	mc.mu.RLock()
	defer mc.mu.RUnlock()

	for _, sc := range mc.secureChannels {
		sc.Close()
	}
	log.Info().Msg("matter controller closed")
}
