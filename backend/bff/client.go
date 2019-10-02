package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/byuoitav/common/structs"
)

type Client struct {
	buildingID string
	roomID     string

	room     structs.Room
	state    structs.PublicRoom
	uiConfig structs.UIConfig

	httpClient *http.Client

	Out chan []byte
	// probably add a send function
	// 	In  chan []byte
}

func RegisterClient(ctx context.Context, roomID, controlGroupID string) (*Client, error) {
	split := strings.Split(roomID, "-")
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid roomID - must match format BLDG-ROOM")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client := &Client{
		buildingID: split[0],
		roomID:     roomID,
		httpClient: &http.Client{},
		Out:        make(chan []byte, 8),
	}

	errCh := make(chan error, 3)
	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer close(doneCh)
		wg.Wait()
	}()

	go func() {
		var err error
		defer wg.Done()

		client.state, err = GetRoomState(ctx, client.httpClient, client.roomID)
		if err != nil {
			errCh <- fmt.Errorf("unable to get ui config: %v", err)
		}
	}()

	go func() {
		var err error
		defer wg.Done()

		client.room, err = GetRoomConfig(ctx, client.httpClient, client.roomID)
		if err != nil {
			errCh <- fmt.Errorf("unable to get room config: %v", err)
		}
	}()

	go func() {
		var err error
		defer wg.Done()

		client.uiConfig, err = GetUIConfig(ctx, client.httpClient, client.roomID)
		if err != nil {
			errCh <- fmt.Errorf("unable to get ui config: %v", err)
		}
	}()

	select {
	case err := <-errCh:
		return nil, fmt.Errorf("unable to get room information: %v", err)
	case <-ctx.Done():
		return nil, fmt.Errorf("unable to get room information: all requests timed out")
	case <-doneCh:
	}

	// write the inital room info
	room := client.GetRoom()
	buf, err := json.Marshal(room)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal room: %s", err)
	}

	client.Out <- buf
	return client, nil
}

func (c *Client) GetRoom() Room {
	var room Room
	room.ID = ID(c.roomID)
	room.Name = c.room.Name

	for _, preset := range c.uiConfig.Presets {
		var cg ControlGroup
		cg.ID = ID(preset.Name)
		cg.Name = preset.Name

		for _, name := range preset.Displays {
			var d Display
			config := GetDeviceConfigByName(c.room.Devices, name)
			state := GetDisplayStateByName(c.state.Displays, name)

			d.ID = ID(config.ID)
			d.Name = config.DisplayName
			d.Input = ID(state.Input)

			if state.Blanked != nil {
				d.Blanked = *state.Blanked
			}

			cg.Displays = append(cg.Displays, d)
		}

		for _, name := range preset.Inputs {
			var i Input
			config := GetDeviceConfigByName(c.room.Devices, name)

			i.ID = ID(config.ID)
			i.Name = config.DisplayName
			// i.Icon = preset.
			// TODO subinputs

			cg.Inputs = append(cg.Inputs, i)
		}

		for group, audioDevices := range preset.AudioGroups {
			var ag AudioGroup
			ag.ID = ID(group)
			ag.Name = group

			for _, name := range audioDevices {
				var ad AudioDevice
				config := GetDeviceConfigByName(c.room.Devices, name)
				state := GetAudioDeviceStateByName(c.state.AudioDevices, name)

				ad.ID = ID(config.ID)
				ad.Name = config.DisplayName

				if state.Volume != nil {
					ad.Level = *state.Volume
				}

				if state.Muted != nil {
					ad.Muted = *state.Muted
				}

				ag.AudioDevices = append(ag.AudioDevices, ad)
			}

			cg.AudioGroups = append(cg.AudioGroups, ag)
		}

		// TODO PresentGroups
		room.ControlGroups = append(room.ControlGroups, cg)
	}

	return room
}
