package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/ui"
	"go.uber.org/zap"
)

func (c *client) doStateTransition(ctx context.Context, cs ui.StateControlConfig, modify func(ui.State) ui.State) error {
	transition, ok := c.matchStateTransition(cs.StateTransitions)
	if !ok {
		return fmt.Errorf("no transition matched the current state")
	}

	state := c.mergeStates(transition.Action.SetStates...)
	if modify != nil {
		state = modify(state)
	}

	if len(state.Devices) > 0 {
		c.log.Debug("Setting room state", zap.Any("req", state))

		newState, err := c.avController.SetRoomState(ctx, c.roomID, avcontrol.StateRequest(state))
		if err != nil {
			return fmt.Errorf("unable to set room state: %w", err)
		}

		for i := range newState.Errors {
			c.log.Warn("error in set room response", zap.Any("stateError", newState.Errors[i]))
		}

		// update room state, send update room to frontend
		c.updateRoomStateFromState(newState)
		c.sendJSONMsg("room", c.Room())
	}

	for _, req := range transition.Action.Requests {
		if err := c.doGenericRequest(ctx, req); err != nil {
			c.log.Error("unable to do generic request", zap.Error(err), zap.String("url", req.URL.String()))
		}
	}

	return nil
}

func (c *client) matchStateTransition(transitions []ui.StateTransition) (ui.StateTransition, bool) {
	// could optimize this by remembering which states we've already matched
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()

	c.configMu.RLock()
	defer c.configMu.RUnlock()

	for _, transition := range transitions {
		if c.doesStateMatch(transition.MatchStates...) {
			return transition, true
		}
	}

	return ui.StateTransition{}, false
}

func (c *client) doGenericRequest(ctx context.Context, gcr ui.GenericRequest) error {
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
