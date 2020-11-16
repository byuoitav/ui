package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (c *client) setMute(data []byte) {
	var msg struct {
		ControlGroup string `json:"controlGroup"`
		Mute         bool   `json:"mute"`
		AudioGroup   string `json:"audioGroup"`
		AudioDevice  string `json:"audioDevice"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cgID := c.controlGroupID
	if msg.ControlGroup != "" {
		cgID = msg.ControlGroup
	}

	// make sure control group exists
	cg, ok := c.config.ControlGroups[cgID]
	if !ok {
		// TODO log/send invalid control group error
		return
	}

	if msg.AudioGroup == "" && msg.AudioDevice == "" {
		cs := cg.Audio.Media.Copy()
		cs.APIRequest = fillMuteRequest(cs.APIRequest, msg.Mute)
		c.doControlSet(ctx, *cs)
		return
	}

	for _, ag := range cg.Audio.Groups {
		if ag.Name != msg.AudioGroup {
			continue
		}

		for _, ad := range ag.AudioDevices {
			if ad.Name != msg.AudioDevice {
				continue
			}

			cs := ad.Copy()
			cs.APIRequest = fillMuteRequest(cs.APIRequest, msg.Mute)
			c.doControlSet(ctx, *cs)
			return
		}
	}

	fmt.Printf("invalid!!!\n")
	// TODO some kind of invalid ag/ad error
}
