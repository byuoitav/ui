package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (c *client) setMute(data []byte) {
	var msg struct {
		Mute        bool   `json:"mute"`
		AudioGroup  string `json:"audioGroup"`
		AudioDevice string `json:"audioDevice"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cg, ok := c.config.ControlGroups[c.controlGroupID]
	if !ok {
		// TODO log/send invalid control group error
		return
	}

	if msg.AudioGroup == "" && msg.AudioDevice == "" {
		if msg.Mute {
			c.doStateTransition(ctx, cg.Audio.Media.Mute, nil)
		} else {
			c.doStateTransition(ctx, cg.Audio.Media.Unmute, nil)
		}

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

			if msg.Mute {
				c.doStateTransition(ctx, ad.Mute, nil)
			} else {
				c.doStateTransition(ctx, ad.Unmute, nil)
			}

			return
		}
	}

	fmt.Printf("invalid!!!\n")
	// TODO some kind of invalid ag/ad error
}
