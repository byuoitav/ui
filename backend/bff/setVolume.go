package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/byuoitav/common/structs"
	"go.uber.org/zap"
)

type SetVolume struct {
}

type SetVolumeMessage struct {
	AudioDevice ID  `json:"audioDevice"`
	Level       int `json:"level"`
}

func (sv SetVolume) Do(c *Client, data []byte) {
	var msg SetVolumeMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.Warn("invalid value for setVolume", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setVolume: %s", err))
		return
	}

	// this shouldn't take longer than 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get the current control group
	cg := c.GetRoom().ControlGroups[c.selectedControlGroupID]
	c.Info("setting volume", zap.String("on", string(msg.AudioDevice)), zap.Int("to", msg.Level), zap.String("controlGroup", string(cg.ID)))

	// build request to send av api
	// if audioDevice isn't set, then they want to change the media level
	// if it is, just change the given audio device
	var state structs.PublicRoom
	if len(msg.AudioDevice) > 0 {
		// to change media volume, we set the volume on _all_ of the matching presets' audioDevices
		preset, err := c.GetPresetByName(string(cg.ID))
		if err != nil {
			c.Warn("failed to set volume on media audio", zap.Error(err))
			c.Out <- ErrorMessage(fmt.Errorf("failed to set volume media audio: %w", err))
		}

		// add each device to the av api request
		for _, dev := range preset.AudioDevices {
			state.AudioDevices = append(state.AudioDevices, structs.AudioDevice{
				PublicDevice: structs.PublicDevice{
					Name: dev,
				},
				Volume: &msg.Level,
			})
		}
	} else {
		state.AudioDevices = append(state.AudioDevices, structs.AudioDevice{
			PublicDevice: structs.PublicDevice{
				Name: msg.AudioDevice.GetName(),
			},
			Volume: &msg.Level,
		})
	}

	if err := c.SendAPIRequest(ctx, state); err != nil {
		c.Warn("failed to set volume", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to set volume: %w", err))
	}
}
