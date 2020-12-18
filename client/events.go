package client

import (
	"context"
	"time"

	"github.com/byuoitav/ui"
	"go.uber.org/zap"
)

func (c *client) publishState(ctx context.Context) {
	states := c.curStates(true)

	var arr []string
	for state, v := range states {
		if v {
			arr = append(arr, state)
		}
	}

	event := ui.Event{
		Key:   "states",
		Value: "",
		Data:  arr,
		Room:  c.roomID,
	}

	c.log.Debug("Publishing states", zap.Any("event", event))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := c.publisher.Publish(ctx, event); err != nil {
		c.log.Warn("unable to publish states", zap.Error(err))
	}
}
