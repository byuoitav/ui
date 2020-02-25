package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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
	sharingChanged := false

	// Go through all sharing groups
	c.shareMutex.Lock()
	for master, list := range c.sharing {
		done := false
		// If the master is changing input
		if master == ID(msg.DisplayID) {
			// Each active gets their outputs added to the public room with the input being the input of the master
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
			done = true
		} else {
			// Otherwise go through each active member of the list
			for i, a := range list.Active {
				// If the active member is the changed input
				if a == ID(msg.DisplayID) {
					//Remove it from the active list and add it to the inactive list
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
		}
		if done {
			sharingChanged = true
			break
		}
	}
	c.shareMutex.Unlock()

	// find the display by ID
	disp, err := getDisplay(cg, ID(msg.DisplayID))
	if err != nil {
		fmt.Printf("no!!!\n")
		return
	}

	// For each of the displays outputs
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
		// Add each display to the list of displays to change on the new state
		state.Displays = append(state.Displays, display)
	}

	/* TODO
	if len(si.OnSameInput.URL) > 0 {
		// send mute request
	}
	*/
	if sharingChanged {
		go updateLazSharing(context.TODO(), c)
	}
	// Make the state changes
	err = c.SendAPIRequest(context.TODO(), state)
	if err != nil {
		c.Warn("failed to change input", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to change input: %s", err))
	}
}
