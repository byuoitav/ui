package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type cameraControlMsg struct {
	Camera string `json:"camera"`
}

func (c *client) tiltUp(data []byte) {
	var msg cameraControlMsg
	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Printf("error: %s\n", err)
		// TODO log/send error
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// make sure control group exists
	cg, ok := c.config.ControlGroups[c.controlGroupID]
	if !ok {
		// TODO log/send invalid control group error
		return
	}

	for _, cam := range cg.Cameras {
		if cam.Name == msg.Camera {
			c.doControlSet(ctx, cam.TiltUp)
			return
		}
	}

	// TODO camera not found
}
