package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/byuoitav/ui"
)

type Publisher struct {
	URL    string
	System string
}

type event struct {
	GeneratingSystem string      `json:"generating-system"`
	Timestamp        time.Time   `json:"timestamp"`
	Tags             []string    `json:"event-tags"`
	TargetDevice     deviceInfo  `json:"target-device"`
	AffectedRoom     roomInfo    `json:"affected-room"`
	Key              string      `json:"key"`
	Value            string      `json:"value"`
	User             string      `json:"user"`
	Data             interface{} `json:"data,omitempty"`
}

type roomInfo struct {
	BuildingID string `json:"buildingID,omitempty"`
	RoomID     string `json:"roomID,omitempty"`
}

type deviceInfo struct {
	roomInfo
	DeviceID string `json:"deviceID,omitempty"`
}

func (p *Publisher) Publish(ctx context.Context, uiEvent ui.Event) error {
	e := event{
		GeneratingSystem: p.System,
		Timestamp:        time.Now(),
		Key:              uiEvent.Key,
		Value:            uiEvent.Value,
		Data:             uiEvent.Data,
		Tags:             uiEvent.Tags,
		AffectedRoom:     newRoomInfo(uiEvent.Room),
		TargetDevice:     newDeviceInfo(uiEvent.Device),
	}

	reqBody, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("unable to marshal event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.URL, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("unable to build request: %w", err)
	}

	req.Header.Add("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("got a %v response from API", resp.StatusCode)
	}

	return nil
}

func newRoomInfo(room string) roomInfo {
	split := strings.Split(room, "-")
	if len(split) != 2 {
		return roomInfo{
			RoomID: room,
		}
	}

	return roomInfo{
		BuildingID: split[0],
		RoomID:     split[0] + "-" + split[1],
	}
}

func newDeviceInfo(device string) deviceInfo {
	split := strings.Split(device, "-")
	if len(split) != 3 {
		return deviceInfo{
			DeviceID: device,
		}
	}

	return deviceInfo{
		roomInfo: roomInfo{
			BuildingID: split[0],
			RoomID:     split[0] + "-" + split[1],
		},
		DeviceID: split[0] + "-" + split[1] + "-" + split[2],
	}
}
