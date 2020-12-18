package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/byuoitav/ui"
)

func (c *client) helpRequest(data []byte) {
	var msg struct {
		RequestType string `json:"requestType"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Printf("error unmarshaling: %s\n", err)
		return
	}

	event := ui.Event{
		Room:  c.roomID,
		Key:   "help-request",
		Value: msg.RequestType,
	}

	if err := c.publisher.Publish(context.TODO(), event); err != nil {
		fmt.Printf("error sending event: %s\n", err)
		return
	}
}
