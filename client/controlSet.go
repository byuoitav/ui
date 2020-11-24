package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/byuoitav/ui"
	"go.uber.org/zap"
)

// TODO need to make these things happen in parallel (?)
// TODO make sure APIRequest actually exists
func (c *client) doControlSet(ctx context.Context, cs ui.ControlSet) error {
	if len(cs.APIRequest.Devices) > 0 {
		c.log.Debug("Setting room state", zap.Any("req", cs.APIRequest))

		state, err := c.avController.SetRoomState(ctx, c.roomID, cs.APIRequest)
		if err != nil {
			return fmt.Errorf("unable to set room state: %w", err)
		}

		for i := range state.Errors {
			c.log.Warn("error in set room response", zap.Any("stateError", state.Errors[i]))
		}

		// update room state, send update room to frontend
		c.updateRoomStateFromState(state)
		c.sendJSONMsg("room", c.Room())
	}

	for _, req := range cs.Requests {
		if err := c.doGenericRequest(ctx, req); err != nil {
			c.log.Error("unable to do generic request", zap.Error(err), zap.String("url", req.URL.String()))
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
