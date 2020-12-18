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

func (c *client) doStateTransition(ctx context.Context, modify func(ui.State) ui.State, stateControls ...ui.StateControlConfig) error {
	states := c.curStates(true)

	var transitions []ui.StateTransition
	for i, stateControl := range stateControls {
		transition, ok := c.matchStateTransition(states, stateControl.StateTransitions)
		if !ok {
			return fmt.Errorf("no transition matched the current state for stateControl %d", i)
		}

		transitions = append(transitions, transition)
	}

	var setStates []string
	for _, transition := range transitions {
		setStates = append(setStates, transition.Action.SetStates...)
	}

	state := c.mergeStates(setStates...)
	if modify != nil {
		state = modify(state.Copy())
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
		go c.publishState(context.Background())
	}

	for _, transition := range transitions {
		for _, req := range transition.Action.Requests {
			if err := c.doGenericRequest(ctx, req); err != nil {
				c.log.Error("unable to do generic request", zap.Error(err), zap.String("url", req.URL.String()))
			}
		}
	}

	return nil
}

func (c *client) matchStateTransition(states map[string]bool, transitions []ui.StateTransition) (ui.StateTransition, bool) {
	for _, transition := range transitions {
		if c.matchStates(states, transition.MatchStates) {
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
