package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/byuoitav/common/structs"
	"go.uber.org/zap"
)

type SetMuted struct {
}

type SetMutedMessage struct {
	AudioDevice ID   `json:"audioDevice"`
	Muted       bool `json:"muted"`
}

func (sm SetMuted) Do(c *Client, data []byte) {
	var msg SetMutedMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.Warn("invalid value for setMuted", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setMuted: %s", err))
		return
	}

	// this shouldn't take longer than 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get the current control group
	cg := c.GetRoom().ControlGroups[c.selectedControlGroupID]
	c.Info("Setting muted", zap.String("on", string(msg.AudioDevice)), zap.Bool("to", msg.Muted), zap.String("controlGroup", string(cg.ID)))

	// build request to send av api
	// if audioDevice isn't set, then they want to change the media mute
	// if it is, just change then given audio device
	var state structs.PublicRoom
	if len(msg.AudioDevice) == 0 {
		// to change media mute, we set the mute on _all_ of the matching presets' audioDevices
		preset, err := c.GetPresetByName(string(cg.ID))
		if err != nil {
			c.Warn("failed to set muted on media audio", zap.Error(err))
			c.Out <- ErrorMessage(fmt.Errorf("failed to set muted on media audio: %w", err))
		}

		// add each device to the av api request
		for _, dev := range preset.AudioDevices {
			state.AudioDevices = append(state.AudioDevices, structs.AudioDevice{
				PublicDevice: structs.PublicDevice{
					Name: dev,
				},
				Muted: &msg.Muted,
			})
		}
	} else {
		state.AudioDevices = append(state.AudioDevices, structs.AudioDevice{
			PublicDevice: structs.PublicDevice{
				Name: msg.AudioDevice.GetName(),
			},
			Muted: &msg.Muted,
		})
	}

	if err := c.SendAPIRequest(ctx, state); err != nil {
		c.Warn("failed to set muted", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to set muted: %w", err))
	}

	c.Info("Finished setting muted", zap.String("on", string(msg.AudioDevice)), zap.Bool("to", msg.Muted), zap.String("controlGroup", string(cg.ID)))
}
