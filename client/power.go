package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (c *client) setPower(data []byte) {
	var msg struct {
		ControlGroup string `json:"controlGroup"`
		On           bool   `json:"on"`
		All          bool   `json:"all"` // TODO
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Printf("error: %s\n", err)
		// TODO log/send error
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

	if !msg.On {
		c.doControlSet(ctx, cg.PowerOff)
		return
	}

	c.doControlSet(ctx, cg.PowerOn)
}
