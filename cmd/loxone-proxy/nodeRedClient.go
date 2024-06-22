package main

import (
	"fmt"
	"net/http"
	"net/url"
)

const nodeEndPointURL = "http://192.168.1.62:1880"

type nodeRedClient struct{}

func (client nodeRedClient) turnOn(room string) error {
	roomEndPoint, err := url.JoinPath(nodeEndPointURL, room)
	if err != nil {
		return err
	}
	roomEndPoint = fmt.Sprintf("%s?on=true", roomEndPoint)
	_, err = http.Get(roomEndPoint)
	return err
}

func (client nodeRedClient) turnOff(room string) error {
	roomEndPoint, err := url.JoinPath(nodeEndPointURL, room)
	if err != nil {
		return err
	}
	roomEndPoint = fmt.Sprintf("%s?on=false", roomEndPoint)
	_, err = http.Get(roomEndPoint)
	return err
}
