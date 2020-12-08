package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (c *client) setPower(data []byte) {
	var msg struct {
		On  bool `json:"on"`
		All bool `json:"all"` // TODO...still something we want?
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Printf("error: %s\n", err)
		// TODO log/send error
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

	if !msg.On {
		c.doStateTransition(ctx, nil, cg.PowerOff)
		return
	}

	c.doStateTransition(ctx, nil, cg.PowerOn)
}
