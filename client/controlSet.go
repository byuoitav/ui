package client

import (
	"context"
	"fmt"

	"github.com/byuoitav/ui"
)

// TODO need to reconcile state coming back and update our current state
// TODO need to make these things happen in parallel (?)
// TODO make sure APIRequest actually exists
func (c *client) doControlSet(ctx context.Context, cs ui.ControlSet) error {
	fmt.Printf("doing control set\n")
	if err := c.avController.SetRoomState(ctx, c.roomID, cs.APIRequest); err != nil {
		fmt.Printf("error: %s\n", err)
	}

	for _, req := range cs.Requests {
		c.doGenericRequest(ctx, req)
	}

	return nil
}

func (c *client) doGenericRequest(ctx context.Context, req ui.GenericControlRequest) error {
	return nil
}
