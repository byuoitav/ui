package client

import (
	"context"
	"fmt"
)

func (c *client) updateRoomState(ctx context.Context) error {
	state, err := c.avController.RoomState(ctx, c.roomID)
	if err != nil {
		return fmt.Errorf("unable to get state: %w", err)
	}

	// TODO something with the errors...?

	state.Errors = nil
	c.state = state
	return nil
}
