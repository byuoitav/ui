package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/byuoitav/ui"
)

// ok so the question for both volume/mute:
// how do we handle controlSets where we aren't actually supposed to set the "template" muted/volume values?
// i'm thinking we may have to separate the controlSets into a `setVolume` controlSet *and* and `setMute` controlSet
func (c *client) setVolume(data []byte) {
	var msg struct {
		Volume      int    `json:"volume"`
		AudioGroup  string `json:"audioGroup"`
		AudioDevice string `json:"audioDevice"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// make sure control group exists
	cg, ok := c.config.ControlGroups[c.controlGroupID]
	if !ok {
		// TODO log/send invalid control group error
		return
	}

	modify := func(state ui.State) ui.State {
		for _, dev := range state.Devices {
			for block, vol := range dev.Volumes {
				if vol != -1 {
					continue
				}

				dev.Volumes[block] = msg.Volume
			}
		}

		return state
	}

	if msg.AudioGroup == "" && msg.AudioDevice == "" {
		c.doStateTransition(ctx, modify, cg.Audio.Media.Volume)
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

			c.doStateTransition(ctx, modify, ad.Volume)
			return
		}
	}

	fmt.Printf("invalid!!!\n")
	// TODO some kind of invalid ag/ad error
}
