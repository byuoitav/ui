package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/lazarette/lazarette"
	"github.com/golang/protobuf/ptypes/timestamp"
	"go.uber.org/zap"
)

type SetVolume struct {
}

type SetVolumeMessage struct {
	AudioDeviceID string `json:"audioDevice"`
	Level         int    `json:"level"`
}

func (sv SetVolume) Do(c *Client, data []byte) {
	var msg SetVolumeMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setVolume: %s", err))
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

	c.Info("setting volume", zap.String("on", msg.AudioDeviceID), zap.Int("level", msg.Level), zap.String("controlGroup", string(cg.ID)))

	// build request to send av api
	var state structs.PublicRoom
	// TODO check length of split
	idSplit := strings.Split(string(ad.ID), "-")
	state.AudioDevices = append(state.AudioDevices, structs.AudioDevice{
		PublicDevice: structs.PublicDevice{
			Name: idSplit[2],
		},
		Volume: &msg.Level,
	})

	if err := c.SendAPIRequest(context.TODO(), state); err != nil {
		c.Warn("failed to set volume", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to set volume: %s", err))
	} else {
		jVol, err := json.Marshal(msg.Level)
		if err != nil {
			c.Warn("failed to create JSON for volume", zap.Error(err))
			c.Out <- ErrorMessage(fmt.Errorf("failed to create JSON for volume: %s", err))
			return
		}
		c.Info("Setting volume in lazarette")
		_, err = c.lazState.Client.Set(context.TODO(), &lazarette.KeyValue{
			Key:  fmt.Sprintf("%v-_master_volume", c.roomID),
			Data: jVol,
			Timestamp: &timestamp.Timestamp{
				Seconds: time.Now().Unix(),
			},
		})
		if err != nil {
			c.Warn("failed to set volume", zap.Error(err))
			c.Out <- ErrorMessage(fmt.Errorf("failed to set volume in %s: %v", c.roomID, err))
		}

	}

}
