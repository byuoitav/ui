package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/byuoitav/common/structs"
	"go.uber.org/zap"
)

// HTTPRequest .
type HTTPRequest struct {
	Method string          `json:"method"`
	URL    string          `json:"url"`
	Body   json.RawMessage `json:"body"`
}

// SetInput .
type SetInput struct {
	OnSameInput HTTPRequest `json:"onSameInput"`
}

// SetInputMessage .
type SetInputMessage struct {
	DisplayGroup ID `json:"displayGroup"`
	Input        ID `json:"input"`
}

// Do .
func (si SetInput) Do(c *Client, data []byte) {
	var msg SetInputMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		c.Warn("invalid value for setInput", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setInput: %s", err))
		return
	}

	// this shouldn't take longer than 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cg := c.GetRoom().ControlGroups[c.selectedControlGroupID]
	c.Info("Setting input", zap.String("on", string(msg.DisplayGroup)), zap.String("to", string(msg.Input)), zap.String("controlGroup", string(cg.ID)))

	// find the display group by ID
	group, err := GetDisplayGroupByID(cg.DisplayGroups, msg.DisplayGroup)
	if err != nil {
		c.Warn("failed to change input", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to change input: %s", err))
		return
	}

	// build the state object
	var state structs.PublicRoom

	inputName := msg.Input.GetName()

	// figure out share stuff
	if shareMap := c.getShareMap(); shareMap != nil {
		if data, ok := shareMap[msg.DisplayGroup]; ok {
			switch data.State {
			case MinionActive:
				return
			case Unshare:
				for _, active := range data.Active {
					minionGroup, err := GetDisplayGroupByID(cg.DisplayGroups, active)
					if err != nil {
						// TODO
					}

					for _, disp := range minionGroup.Displays {
						state.Displays = append(state.Displays, structs.Display{
							PublicDevice: structs.PublicDevice{
								Name:  disp.ID.GetName(),
								Input: msg.Input.GetName(),
							},
							Blanked: BoolP(false),
						})
					}
				}
			case MinionInactive:
				swap := data.Master == msg.Input
				if swap {
					// we are switching into the active list
					masterGroup, err := GetDisplayGroupByID(cg.DisplayGroups, data.Master)
					if err != nil {
						// TODO
					}

					inputName = masterGroup.Input.GetName()

					// TODO update lazarette state
				}
			}
		}
	}

	for _, disp := range group.Displays {
		display := structs.Display{
			PublicDevice: structs.PublicDevice{
				Name:  disp.ID.GetName(),
				Input: inputName,
			},
			Blanked: BoolP(false),
		}

		// Add each display to the list of displays to change on the new state
		state.Displays = append(state.Displays, display)
	}

	// make the state changes
	if err := c.SendAPIRequest(ctx, state); err != nil {
		c.Warn("failed to change input", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to change input: %s", err))
	}

	c.Info("Finished setting input", zap.String("on", string(msg.DisplayGroup)), zap.String("to", string(msg.Input)), zap.String("controlGroup", string(cg.ID)))
}

func removeID()
