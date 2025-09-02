package main

import (
	"fmt"
	"net/http"
	"net/url"
)

type nodeRedClient struct {
	nodeEndPointURL string
}

func (client nodeRedClient) turnOn(room string) error {
	roomEndPoint, err := url.JoinPath(client.nodeEndPointURL, room)
	if err != nil {
		return err
	}
	roomEndPoint = fmt.Sprintf("%s?on=true", roomEndPoint)
	_, err = http.Get(roomEndPoint)
	return err
}

func (client nodeRedClient) turnOff(room string) error {
	roomEndPoint, err := url.JoinPath(client.nodeEndPointURL, room)
	if err != nil {
		return err
	}
	roomEndPoint = fmt.Sprintf("%s?on=false", roomEndPoint)
	_, err = http.Get(roomEndPoint)
	return err
}
