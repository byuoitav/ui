package client

import (
	"context"
	"fmt"

	"github.com/byuoitav/ui"
)

// TODO need to make these things happen in parallel (?)
// TODO make sure APIRequest actually exists
func (c *client) doControlSet(ctx context.Context, cs ui.ControlSet) error {
	fmt.Printf("doing control set\n")
	state, err := c.avController.SetRoomState(ctx, c.roomID, cs.APIRequest)
	if err != nil {
		return fmt.Errorf("unable to set room state: %w", err)
	}

	for range state.Errors {
		// send these errors to the frontend?
		// or just log them? idk
	}

	// update room state, send update room to frontend
	c.updateRoomStateFromState(state)
	c.sendJSONMsg("room", c.Room())

	for _, req := range cs.Requests {
		c.doGenericRequest(ctx, req)
	}

	return nil
}

func (c *client) doGenericRequest(ctx context.Context, req ui.GenericControlRequest) error {
	return nil
}
