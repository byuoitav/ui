package bff

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/ui/log"
	"go.uber.org/zap"
)

type Client struct {
	id                     string
	buildingID             string
	roomID                 string
	selectedControlGroupID string

	room     structs.Room
	state    structs.PublicRoom
	uiConfig UIConfig

	httpClient *http.Client

	// messages going out to the client
	Out chan Message

	// events put in this channel get sent to the hub
	SendEvent chan events.Event

	*zap.Logger
}

func RegisterClient(ctx context.Context, roomID, controlGroupID, name string) (*Client, error) {
	log.P.Info("Registering client", zap.String("roomID", roomID), zap.String("controlGroupID", controlGroupID), zap.String("name", name))

	split := strings.Split(roomID, "-")
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid roomID %q - must match format BLDG-ROOM", roomID)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	c := &Client{
		id:         name,
		buildingID: split[0],
		roomID:     roomID,
		httpClient: &http.Client{},
		Out:        make(chan Message, 1),
		SendEvent:  make(chan events.Event),
		Logger:     log.P.Named(name),
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

		c.state, err = GetRoomState(ctx, c.httpClient, c.roomID)
		if err != nil {
			c.Warn("unable to get room state", zap.Error(err))
			errCh <- fmt.Errorf("unable to get ui config: %v", err)
		}

		c.Debug("Successfully got room state")
	}()

	go func() {
		var err error
		defer wg.Done()

		c.room, err = GetRoomConfig(ctx, c.httpClient, c.roomID)
		if err != nil {
			c.Warn("unable to get room config", zap.Error(err))
			errCh <- fmt.Errorf("unable to get room config: %v", err)
		}

		c.Debug("Successfully got room config")
	}()

	go func() {
		var err error
		defer wg.Done()

		c.uiConfig, err = GetUIConfig(ctx, c.httpClient, c.roomID)
		if err != nil {
			c.Warn("unable to get ui config", zap.Error(err))
			errCh <- fmt.Errorf("unable to get ui config: %v", err)
		}

		c.Debug("Successfully got ui config")
	}()

	select {
	case err := <-errCh:
		return nil, fmt.Errorf("unable to get room information: %v", err)
	case <-ctx.Done():
		return nil, fmt.Errorf("unable to get room information: all requests timed out")
	case <-doneCh:
	}

	room := c.GetRoom()
	if _, ok := room.ControlGroups[controlGroupID]; ok {
		c.selectedControlGroupID = controlGroupID
	}

	//check if controlgroup is empty, if not turn on displays in controlgroup
	if c.selectedControlGroupID != "" {
		var displays []ID
		for _, display := range room.ControlGroups[c.selectedControlGroupID].Displays {
			displays = append(displays, display.ID)
		}

		setPowerMessage := SetPowerMessage{
			Displays: displays,
			Status:   "on",
		}

		// turn the control group on - this will send the room to the client
		err := c.CurrentPreset().Actions.SetPower.DoWithMessage(ctx, c, setPowerMessage)
		if err != nil {
			return nil, fmt.Errorf("unable to power on the following displays %s: %s", displays, err)
		}

		go c.HandleEvents()
		return c, nil
	}

	c.Info("Got all initial information, sending room to client")

	// write the inital room info
	msg, err := JSONMessage("room", c.GetRoom())
	if err != nil {
		return nil, fmt.Errorf("unable to marshal room: %s", err)
	}

	c.Out <- msg

	go c.HandleEvents()
	return c, nil
}

func (c *Client) GetRoom() Room {
	room := Room{
		ID:                   ID(c.roomID),
		Name:                 c.room.Name,
		ControlGroups:        make(map[string]ControlGroup),
		SelectedControlGroup: "", // TODO where is this saved? c?
	}

	for _, preset := range c.uiConfig.Presets {
		cg := ControlGroup{
			ID:   ID(preset.Name),
			Name: preset.Name,
			Support: Support{
				HelpRequested: false, // This info also needs to be saved...
				HelpMessage:   "Request Help",
				HelpEnabled:   true,
			},
		}

		for _, name := range preset.Displays {
			config := GetDeviceConfigByName(c.room.Devices, name)
			state := GetDisplayStateByName(c.state.Displays, name)
			outputIcon := "tv"

			for _, IOconfig := range c.uiConfig.OutputConfiguration {
				if config.Name != IOconfig.Name {
					continue
				}

				outputIcon = IOconfig.Icon
			}

			// figure out what the current input for this display is
			// we are assuming that input is roomid - input name
			// unless it's blanked, then the "input" is blank
			curInput := c.roomID + "-" + state.Input
			if state.Blanked != nil && *state.Blanked {
				curInput = "blank"
			}

			d := Display{
				ID:    ID(config.ID),
				Input: ID(curInput),
			}

			// TODO outputs when we do sharing
			d.Outputs = append(d.Outputs, IconPair{
				ID:   ID(config.ID),
				Name: config.DisplayName,
				Icon: Icon{outputIcon},
			})

			cg.Displays = append(cg.Displays, d)
		}

		// add a blank input as the first input
		cg.Inputs = append(cg.Inputs, Input{
			ID: ID("blank"),
			IconPair: IconPair{
				Name: "Blank",
				Icon: Icon{"crop_landscape"},
			},
			Disabled: false,
		})

		for _, name := range preset.Inputs {
			config := GetDeviceConfigByName(c.room.Devices, name)
			inputIcon := "settings_input_hdmi"

			for _, IOconfig := range c.uiConfig.InputConfiguration {
				if config.Name != IOconfig.Name {
					continue
				}

				inputIcon = IOconfig.Icon
			}

			i := Input{
				ID: ID(config.ID),
				IconPair: IconPair{
					Name: config.DisplayName,
					Icon: Icon{inputIcon},
				},
				Disabled: false, // TODO look at the current displays reachable inputs to determine
			}

			// TODO subinputs

			cg.Inputs = append(cg.Inputs, i)
		}

		if len(preset.AudioGroups) > 0 {
			for group, audioDevices := range preset.AudioGroups {
				ag := AudioGroup{
					ID:    ID(group),
					Name:  group,
					Muted: true,
				}

				for _, name := range audioDevices {
					config := GetDeviceConfigByName(c.room.Devices, name)
					state := GetAudioDeviceStateByName(c.state.AudioDevices, name)
					audioIcon := "mic"

					for _, IOconfig := range c.uiConfig.OutputConfiguration {
						if config.Name != IOconfig.Name {
							continue
						}

						audioIcon = IOconfig.Icon
					}

					ad := AudioDevice{
						ID: ID(config.ID),
						IconPair: IconPair{
							Name: config.DisplayName,
							Icon: Icon{audioIcon},
						},
					}

					if state.Volume != nil {
						ad.Level = *state.Volume
					}

					if state.Muted != nil {
						ad.Muted = *state.Muted
					}

					if !ad.Muted {
						ag.Muted = false
					}

					ag.AudioDevices = append(ag.AudioDevices, ad)
				}

				cg.AudioGroups = append(cg.AudioGroups, ag)
			}
		} else {
			// create the displaysAG
			if len(preset.AudioDevices) >= 1 {
				ag := AudioGroup{
					ID:    "displaysAG",
					Name:  "Display Volume Mixing",
					Muted: true,
				}

				for _, name := range preset.AudioDevices {
					config := GetDeviceConfigByName(c.room.Devices, name)
					state := GetAudioDeviceStateByName(c.state.AudioDevices, name)

					ad := AudioDevice{
						ID: ID(config.ID),
						IconPair: IconPair{
							Name: config.DisplayName,
							Icon: Icon{"tv"},
						},
					}

					if state.Volume != nil {
						ad.Level = *state.Volume
					}

					if state.Muted != nil {
						ad.Muted = *state.Muted
					}

					if !ad.Muted {
						ag.Muted = false
					}

					ag.AudioDevices = append(ag.AudioDevices, ad)
				}

				cg.AudioGroups = append(cg.AudioGroups, ag)
			}

			// create the micsAG
			if len(preset.IndependentAudioDevices) >= 1 {

				ag := AudioGroup{
					ID:    "micsAG",
					Name:  "Microphones",
					Muted: true,
				}

				for _, name := range preset.IndependentAudioDevices {
					config := GetDeviceConfigByName(c.room.Devices, name)
					state := GetAudioDeviceStateByName(c.state.AudioDevices, name)

					ad := AudioDevice{
						ID: ID(config.ID),
						IconPair: IconPair{
							Name: config.DisplayName,
							Icon: Icon{"mic"},
						},
					}

					if state.Volume != nil {
						ad.Level = *state.Volume
					}

					if state.Muted != nil {
						ad.Muted = *state.Muted
					}

					if !ad.Muted {
						ag.Muted = false
					}

					ag.AudioDevices = append(ag.AudioDevices, ad)
				}

				cg.AudioGroups = append(cg.AudioGroups, ag)
			}

		}

		room.ControlGroups[string(cg.ID)] = cg
		// TODO PresentGroups
	}

	room.SelectedControlGroup = ID(c.selectedControlGroupID)

	return room
}

func (c *Client) CurrentPreset() Preset {
	for _, p := range c.uiConfig.Presets {
		if p.Name == c.selectedControlGroupID {
			return p
		}
	}

	return Preset{}
}
