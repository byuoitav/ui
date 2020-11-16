package av

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/ui"
)

var _ ui.AVController = &Controller{}

type Controller struct {
	// BaseURL is the base url of the av-control-api server
	BaseURL string
}

func (a *Controller) RoomState(ctx context.Context, room string) (avcontrol.StateResponse, error) {
	var state avcontrol.StateResponse

	url := fmt.Sprintf("%s/api/v1/room/%s/state", a.BaseURL, room)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return state, fmt.Errorf("unable to build request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return state, fmt.Errorf("unable to do request: %w", err)
	}
	defer resp.Body.Close()

	// TODO check response status code?

	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return state, fmt.Errorf("unable to decode response: %w", err)
	}

	return state, err
}

func (a *Controller) SetRoomState(ctx context.Context, room string, state avcontrol.StateRequest) (avcontrol.StateResponse, error) {
	var sresp avcontrol.StateResponse

	reqBody, err := json.Marshal(state)
	if err != nil {
		return sresp, fmt.Errorf("unable to marshal state: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/room/%s/state", a.BaseURL, room)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return sresp, fmt.Errorf("unable to build request: %w", err)
	}

	req.Header.Add("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return sresp, fmt.Errorf("unable to do request: %w", err)
	}
	defer resp.Body.Close()

	// TODO check response status code?

	if err := json.NewDecoder(resp.Body).Decode(&sresp); err != nil {
		return sresp, fmt.Errorf("unable to decode response: %w", err)
	}

	return sresp, nil
}
