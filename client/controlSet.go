package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/byuoitav/ui"
)

// TODO need to make these things happen in parallel (?)
// TODO make sure APIRequest actually exists
func (c *client) doControlSet(ctx context.Context, cs ui.ControlSet) error {
	fmt.Printf("doing control set\n")
	if len(cs.APIRequest.Devices) > 0 {
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
	}

	for _, req := range cs.Requests {
		if err := c.doGenericRequest(ctx, req); err != nil {
			fmt.Printf("err: %s\n", err)
		}
	}

	return nil
}

func (c *client) doGenericRequest(ctx context.Context, gcr ui.GenericControlRequest) error {
	req, err := http.NewRequestWithContext(ctx, gcr.Method, gcr.URL.String(), bytes.NewReader(gcr.Body))
	if err != nil {
		return fmt.Errorf("unable to build request: %w", err)
	}

	// should we add support for headers?

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to do request: %w", err)
	}
	defer resp.Body.Close()

	// we don't really care about the body...or the response code...
	// so i guess we're just done now
	return nil
}
