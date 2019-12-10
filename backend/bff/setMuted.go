package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/byuoitav/common/structs"
	"go.uber.org/zap"
)

type SetMuted struct {
}

type SetMutedMessage struct {
	AudioDeviceID string `json:"audioDevice"`
	Muted         bool   `json:"muted"`
}

func (sm SetMuted) Do(c *Client, data []byte) {
	var msg SetMutedMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setMuted: %s", err))
		return
	}

	// get the current control group
	cg := c.GetRoom().ControlGroups[c.selectedControlGroupID]
	/* TODO
	if len(cg.ID) == 0 {
		// error
	}
	*/

	// find the audio device
	ad, err := GetAudioDeviceByID(cg.AudioGroups, ID(msg.AudioDeviceID))
	if err != nil {
		c.Out <- ErrorMessage(err)
		return
	}

	c.Info("setting muted", zap.String("on", msg.AudioDeviceID), zap.Bool("muted", msg.Muted), zap.String("controlGroup", string(cg.ID)))

	// build request to send av api
	var state structs.PublicRoom
	// TODO check length of split
	idSplit := strings.Split(string(ad.ID), "-")
	state.AudioDevices = append(state.AudioDevices, structs.AudioDevice{
		PublicDevice: structs.PublicDevice{
			Name: idSplit[2],
		},
		Muted: &msg.Muted,
	})

	if err := c.SendAPIRequest(context.TODO(), state); err != nil {
		c.Warn("failed to set muted", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to set muted: %s", err))
	}
}
