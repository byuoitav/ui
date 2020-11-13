package bff

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/byuoitav/av-control-api/client"
	"github.com/byuoitav/common/structs"
	"go.uber.org/zap"
)

const (
	inputBecomeActive = "_becomeActiveMinion"
)

// HTTPRequest .
type HTTPRequest struct {
	Method string          `json:"method"`
	URL    string          `json:"url"`
	Body   json.RawMessage `json:"body"`
}

// SetInput .
type SetInput struct {
	OnSameInput HTTPRequest `json:"onSameInput"`
}

// SetInputMessage .
type SetInputMessage struct {
	DisplayGroup client.ID `json:"displayGroup"`
	Input        client.ID `json:"input"`
}

// Do .
func (si SetInput) Do(c *Client, data []byte) {
	var msg SetInputMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		c.Warn("invalid value for setInput", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setInput: %s", err))
		return
	}

	room := c.GetRoom()
	cg := room.ControlGroups[c.selectedControlGroupID]

	// find the display group by ID
	group, err := GetDisplayGroupByID(cg.DisplayGroups, msg.DisplayGroup)
	if err != nil {
		c.Warn("failed to change input", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to change input: %s", err))
		return
	}

	// build the state object
	var state structs.PublicRoom
	input := msg.Input.GetName()

	// validate input is valid for myself
	// validate input is valid for all minions

	switch group.ShareInfo.State {
	case stateIsActiveMinion:
		// TODO this is illegal, return an error
	case stateIsMaster:
		// get every display group in this room
		allGroups := room.GetAllDisplayGroups()

		// get all of my active minions
		active, _ := c.getActiveAndInactiveForDisplayGroup(msg.DisplayGroup)

		c.Info("Setting input as sharing master", zap.String("displayGroup", string(msg.DisplayGroup)), zap.String("input", string(msg.Input)), zap.Strings("activeMinions", IDsToStrings(active)))

		// set each active minion's input
		for i := range active {
			mgroup, err := allGroups.GetDisplayGroup(active[i])
			if err != nil {
				// invalid display group id in active list
				// TODO validate active list?
				continue
			}

			for _, disp := range mgroup.Displays {
				state.Displays = append(state.Displays, structs.Display{
					PublicDevice: structs.PublicDevice{
						Name:  disp.ID.GetName(),
						Input: input,
					},
					Blanked: client.BoolP(false),
				})
			}
		}
	case stateIsInactiveMinion:
		if string(msg.Input) != inputBecomeActive {
			// this is just a normal change input request
			c.Info("Setting input", zap.String("displayGroup", string(msg.DisplayGroup)), zap.String("input", string(msg.Input)))
			break
		}

		// find the masters input
		allGroups := room.GetAllDisplayGroups()

		for i := range allGroups {
			if allGroups[i].ID == group.ShareInfo.Master {
				input = allGroups[i].Input.GetName()
				break
			}
		}

		if input == inputBecomeActive {
			err := errors.New("cannot change input, invalid master")
			c.Warn("setInput failed", zap.Error(err))
			c.Out <- ErrorMessage(err)
			return
		}

		c.Info("Becoming an active minion", zap.String("displayGroup", string(msg.DisplayGroup)), zap.String("master", string(group.ShareInfo.Master)))

		// mute my audio devices
		audioDevices, err := cg.GetMediaAudioDeviceIDs(c.uiConfig.Presets)
		if err != nil {
			err := errors.New("cannot change input, preset not found")
			c.Warn("setInput failed", zap.Error(err))
			c.Out <- ErrorMessage(err)
			return
		}
		for i := range audioDevices {
			state.AudioDevices = append(state.AudioDevices, structs.AudioDevice{
				PublicDevice: structs.PublicDevice{
					Name: audioDevices[i].GetName(),
				},
				Muted: client.BoolP(true),
			})
		}

		// update my state in lazarette
		c.lazUpdates <- lazMessage{
			Key: lazSharing + string(group.ID),
			Data: lazShareData{
				State:  stateIsActiveMinion,
				Master: group.ShareInfo.Master,
			},
		}
	default:
		c.Info("Setting input", zap.String("displayGroup", string(msg.DisplayGroup)), zap.String("input", string(msg.Input)))
	}

	// change input for all of my displays
	for _, disp := range group.Displays {
		state.Displays = append(state.Displays, structs.Display{
			PublicDevice: structs.PublicDevice{
				Name:  disp.ID.GetName(),
				Input: input,
			},
			Blanked: client.BoolP(false),
		})
	}

	// this shouldn't take longer than 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// make the state changes
	if err := c.SendAPIRequest(ctx, state); err != nil {
		c.Warn("failed to change input", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to change input: %s", err))
	}

	c.Info("Finished setInput", zap.String("displayGroup", string(msg.DisplayGroup)), zap.String("input", string(msg.Input)))
}
