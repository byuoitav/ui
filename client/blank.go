package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/byuoitav/ui"
)

func (c *client) setBlank(data []byte) {
	var msg struct {
		DisplayGroup string `json:"displayGroup"`
		Blanked      bool   `json:"blanked"`
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Printf("error: %s\n", err)
		// TODO log/send error
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	room := c.Room()

	// find all of the controlSets for the affected displays
	var stateControls []ui.StateControlConfig
	for _, dg := range room.ControlGroups[c.controlGroupID].DisplayGroups {
		if msg.DisplayGroup != dg.Name {
			continue
		}

		for _, d := range dg.Displays {
			// find the matching state control config
			cfg, ok := c.getDisplayConfig(d.Name)
			if !ok {
				// error
				return
			}

			if msg.Blanked {
				stateControls = append(stateControls, cfg.Blank)
			} else {
				stateControls = append(stateControls, cfg.Unblank)
			}
		}

		break
	}

	if len(stateControls) == 0 {
		// some error
		return
	}

	_ = c.doStateTransition(ctx, nil, stateControls...)
}

func (c *client) getDisplayConfig(disp string) (ui.DisplayConfig, bool) {
	c.configMu.RLock()
	defer c.configMu.RUnlock()

	for _, cDisp := range c.config.ControlGroups[c.controlGroupID].Displays {
		if cDisp.Name == disp {
			return cDisp, true
		}
	}

	return ui.DisplayConfig{}, false
}
