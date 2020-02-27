package bff

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/byuoitav/code-service/codemap"
	"go.uber.org/zap"
)

var (
	ErrInvalidControlKey = errors.New("invalid room control key")
)

func GetControlKey(ctx context.Context, codeServiceURL, roomID, controlGroupID string) (string, error) {
	url := fmt.Sprintf("http://%s/%s %s/getControlKey", codeServiceURL, roomID, controlGroupID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("unable to build request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("unable to make request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return "", fmt.Errorf("no control key found")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response: %w", err)
	}

	var key codemap.ControlKey
	if err := json.Unmarshal(body, &key); err != nil {
		return "", fmt.Errorf("unable to parse response: %w", err)
	}

	return key.ControlKey, nil
}

func GetRoomAndControlGroup(ctx context.Context, codeServiceURL, key string) (string, string, error) {
	url := fmt.Sprintf("http://%s/%s/getPreset", codeServiceURL, key)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", fmt.Errorf("unable to build request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("unable to make request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return "", "", fmt.Errorf("invalid room control key")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("unable to read response: %w", err)
	}

	preset := struct {
		RoomID string `json:"RoomID"`
		Name   string `json:"PresetName"`
	}{}

	if err := json.Unmarshal(body, &preset); err != nil {
		return "", "", fmt.Errorf("unable to parse response: %w", err)
	}

	return preset.RoomID, preset.Name, nil
}

func (c *Client) updateControlKey() {
	updateOne := func(id string) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		key, err := GetControlKey(ctx, c.config.CodeServiceAddr, c.roomID, id)
		if err != nil {
			c.Warn("unable to update control key", zap.String("room", c.roomID), zap.String("controlGroup", id), zap.Error(err))
			return
		}

		c.controlKeysMu.Lock()
		defer c.controlKeysMu.Unlock()

		c.controlKeys[id] = key
	}

	updateAll := func() {
		room := c.GetRoom()

		// TODO delete ids that no longer exist? not a big deal
		for id, _ := range room.ControlGroups {
			updateOne(id)
		}

		// send the updated room
		msg, err := JSONMessage("room", c.GetRoom())
		if err != nil {
			c.Warn("unable to build updated room message", zap.Error(err))
			return
		}

		c.Out <- msg
	}

	// set the codes initially
	updateAll()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// then update the codes every minute
	for {
		select {
		case <-c.kill:
			return
		case <-ticker.C:
			updateAll()
		}
	}
}
