package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/byuoitav/common/structs"
	"go.uber.org/zap"
)

// SetInput .
type SetInput struct {
	OnSameInput HTTPRequest `json:"onSameInput"`
}

// SetInputMessage .
type SetInputMessage struct {
	DisplayID string `json:"display"`
	InputID   string `json:"input"`
}

// GetDisplay .
func getDisplay(cg ControlGroup, id ID) (DisplayBlock, error) {
	var disp DisplayBlock
	for i := range cg.DisplayBlocks {
		if cg.DisplayBlocks[i].ID == id {
			return cg.DisplayBlocks[i], nil
		}
	}
	return disp, fmt.Errorf("error display not found")
}

// Do .
func (si SetInput) Do(c *Client, data []byte) {
	var msg SetInputMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setInput: %s", err))
		return
	}

	cg := c.GetRoom().ControlGroups[c.selectedControlGroupID]

	/* TODO
	if len(cg.ID) == 0 {
		// error
	}
	*/

	c.Info("setting input", zap.String("on", msg.DisplayID), zap.String("to", msg.InputID), zap.String("controlGroup", string(cg.ID)))

	var state structs.PublicRoom

	// Go through all sharing groups
	for master, list := range c.sharing {
		// If the master is changing input
		if master == ID(msg.DisplayID) {
			// All active
			for _, m := range list.Active {
				disp, err := getDisplay(cg, m)
				if err != nil {
					fmt.Printf("no!!!\n")
					return
				}
				for _, out := range disp.Outputs {
					// TODO write a getnamefromid func
					dSplit := strings.Split(string(out.ID), "-")
					display := structs.Display{
						PublicDevice: structs.PublicDevice{
							Name: dSplit[2],
						},
					}

					if msg.InputID == "blank" {
						display.Blanked = BoolP(true)
					} else {
						iSplit := strings.Split(string(msg.InputID), "-")
						display.Input = iSplit[2]
						display.Blanked = BoolP(false)
					}

					state.Displays = append(state.Displays, display)
				}
			}
			return
		}
		done := false
		for i, a := range list.Active {
			if a == ID(msg.DisplayID) {
				NewActive := removeID(list.Active, i)
				Inactive := append(list.Inactive, a)
				input := list.Input
				c.sharing[master] = ShareGroups{
					Input:    input,
					Active:   NewActive,
					Inactive: Inactive,
				}
				done = true
				break
			}
		}
		if done {
			break
		}
	}

	// find the display by ID
	disp, err := getDisplay(cg, ID(msg.DisplayID))
	if err != nil {
		fmt.Printf("no!!!\n")
		return
	}

	for _, out := range disp.Outputs {
		// TODO write a getnamefromid func
		dSplit := strings.Split(string(out.ID), "-")
		display := structs.Display{
			PublicDevice: structs.PublicDevice{
				Name: dSplit[2],
			},
		}

		if msg.InputID == "blank" {
			display.Blanked = BoolP(true)
		} else {
			iSplit := strings.Split(string(msg.InputID), "-")
			display.Input = iSplit[2]
			display.Blanked = BoolP(false)
		}

		state.Displays = append(state.Displays, display)
	}

	/* TODO
	if len(si.OnSameInput.URL) > 0 {
		// send mute request
	}
	*/

	err = c.SendAPIRequest(context.TODO(), state)
	if err != nil {
		c.Warn("failed to change input", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to change input: %s", err))
	}
}
