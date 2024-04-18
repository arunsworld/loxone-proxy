package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

const pushoverURL = "https://api.pushover.net/1/messages.json"

type pushoverMessage struct {
	Token   string `json:"token"`
	User    string `json:"user"`
	Message string `json:"message"`
	// Optional
	Priority string `json:"priority"` // high priority = 1; emergency priority = 2
	Title    string `json:"title"`
	TTL      string `json:"ttl"`
}

type pushoverController struct {
	pushoverUser        string
	pushoverAppToken    string
	msgQueue            chan pushoverMessage
	stopProcessingQueue context.CancelFunc
}

func newPushoverController(pushoverUser, pushoverAppToken string) (*pushoverController, error) {
	if pushoverUser == "" {
		return nil, fmt.Errorf("PUSHOVER_USER not found")
	}
	if pushoverAppToken == "" {
		return nil, fmt.Errorf("PUSHOVER_LOXONE_APP_TOKEN not found")
	}
	result := &pushoverController{
		pushoverUser:     pushoverUser,
		pushoverAppToken: pushoverAppToken,
		msgQueue:         make(chan pushoverMessage),
	}
	ctx, stopProcessingQueue := context.WithCancel(context.Background())
	result.stopProcessingQueue = stopProcessingQueue
	go result.processQueue(ctx)

	return result, nil
}

func (pc *pushoverController) send(msg string) {
	pc.msgQueue <- pushoverMessage{
		Message: msg,
	}
}

func (pc *pushoverController) close() {
	pc.stopProcessingQueue()
}

func (pc *pushoverController) processQueue(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("pushoverController processQueue terminating")
			return
		case msg := <-pc.msgQueue:
			if err := pc.deliverMessage(msg); err != nil {
				log.Error().Err(err).Str("content", msg.Message).Msg("unable to deliver push notification")
			}
		}
	}
}

func (pc *pushoverController) deliverMessage(msg pushoverMessage) error {
	msg.Token = pc.pushoverAppToken
	msg.User = pc.pushoverUser

	contentBytes := new(bytes.Buffer)
	if err := json.NewEncoder(contentBytes).Encode(msg); err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	resp, err := http.Post(pushoverURL, "application/json", contentBytes)
	if err != nil {
		return fmt.Errorf("error doing POST to pushoverURL: %w", err)
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received bad status code: %d - %s", resp.StatusCode, string(body))
	}

	log.Debug().Str("content", msg.Message).Msg("pushover msg sent successfully")

	return nil
}
