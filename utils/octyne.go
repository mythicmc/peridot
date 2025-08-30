package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

var httpc = http.Client{
	Transport: &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", "/tmp/octyne.sock.42069")
		},
	},
}

type OctyneServerStatus int

const (
	OctyneServerStatusStopped OctyneServerStatus = iota
	OctyneServerStatusRunning
	OctyneServerStatusCrashed
)

func OctyneGetServers() (map[string]OctyneServerStatus, error) {
	resp, err := httpc.Get("http://unix/servers?extrainfo=true")
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response error %s", resp.Status)
	}
	var responseBody struct {
		Servers map[string]struct {
			Status OctyneServerStatus `json:"status"`
		} `json:"servers"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, fmt.Errorf("failed to decode server statuses: %v", err)
	}

	statuses := make(map[string]OctyneServerStatus, len(responseBody.Servers))
	for name, info := range responseBody.Servers {
		statuses[name] = info.Status
	}
	return statuses, nil
}

func OctyneGetServerStatus(server string) (OctyneServerStatus, error) {
	resp, err := httpc.Get("http://unix/server/" + server)
	if err != nil {
		return 0, err
	} else if resp.StatusCode != 200 {
		return 0, fmt.Errorf("response error %s", resp.Status)
	}
	var responseBody struct {
		Status OctyneServerStatus `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return 0, fmt.Errorf("failed to decode server statuses: %v", err)
	}
	return responseBody.Status, nil
}

func OctyneStartServer(server string) error {
	return octynePostServer(server, "START")
}

func OctyneTerminateServer(server string) error {
	return octynePostServer(server, "TERM")
}

func octynePostServer(server, action string) error {
	if resp, err := httpc.Post(
		"http://unix/server/"+server, "text/plain", bytes.NewReader([]byte(action)),
	); err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("response error %s", resp.Status)
	}
	return nil
}
