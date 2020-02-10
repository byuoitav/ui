package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/byuoitav/common/structs"
	"go.uber.org/zap"
)

// SetSharing .
type SetSharing struct {
}

// SetSharingMessage .
type SetSharingMessage struct {
	Master  ID   `json:"master"`
	Minions []ID `json:"minions"`
}

func contain(l []ID, id ID) (int, bool) {
	for i, e := range l {
		if id == e {
			return i, true
		}
	}
	return -1, false
}

func subArray(big []ID, small []ID) bool {
	for _, v := range small {
		if _, ok := contain(big, v); !ok {
			return false
		}
	}
	return true
}

func removeID(l []ID, index int) []ID {
	l[index] = l[len(l)-1]
	return l[:len(l)-1]
}

func getShareable(presets []Preset, id ID) ([]string, error) {
	for _, p := range presets {
		for _, d := range p.Displays {
			if d == string(id) {
				return p.ShareableDisplays, nil
			}
		}
	}
	return nil, fmt.Errorf("display not found")
}

// On Legacy
func (ss SetSharing) On(c *Client, data []byte) {
	var msg SetSharingMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setSharing: %s", err))
		return
	}

	// Validate that everything can actually be shared to

	if !subArray(c.shareable[msg.Master], msg.Minions) {
		c.Warn("cannot share to all displays in minion list")
		return
	}

	// Remove all minions from other (in)active lists and from sharing

	for _, min := range msg.Minions {
		for master, lists := range c.sharing {
			if min == master { // Absorbing another master
				for _, m := range lists.Active {
					msg.Minions = append(msg.Minions, m)
				}
				for _, m := range lists.Inactive {
					msg.Minions = append(msg.Minions, m)
				}
				delete(c.sharing, master)
			} else if i, exists := contain(lists.Active, min); exists { //Active
				removeID(lists.Active, i)
			} else if i, exists := contain(lists.Inactive, min); exists { //Inactive
				removeID(lists.Inactive, i)
			}
		}
	}

	cg := c.GetRoom().ControlGroups[c.selectedControlGroupID]

	// find the display by ID
	disp, err := getDisplay(cg, ID(msg.Master))
	if err != nil {
		fmt.Printf("no!!!\n")
		return
	}

	c.sharing[msg.Master] = ShareGroups{
		Active: msg.Minions,
		Input:  disp.Input,
	}
	// Change all the shared peeps inputs

	// create public room with new input info, mute all minions
	var state structs.PublicRoom
	for _, m := range msg.Minions {

		d, err := getDisplay(cg, m)
		if err != nil {
			fmt.Printf("no!!!\n")
			return
		}

		for _, out := range d.Outputs {

			// TODO write a getnamefromid func
			dSplit := strings.Split(string(out.ID), "-")
			display := structs.Display{
				PublicDevice: structs.PublicDevice{
					Name: dSplit[2],
				},
			}

			if disp.Input == "blank" {
				display.Blanked = BoolP(true)
			} else {
				iSplit := strings.Split(string(disp.Input), "-")
				display.Input = iSplit[2]
				display.Blanked = BoolP(false)
			}

			state.Displays = append(state.Displays, display)
		}
	}

	err = c.SendAPIRequest(context.TODO(), state)
	if err != nil {
		c.Warn("failed to change input", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to change input: %s", err))
	}
}

// Off Legacy
func (ss SetSharing) Off(c *Client, data []byte) {
	var msg SetSharingMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setSharing: %s", err))
		return
	}

	cg := c.GetRoom().ControlGroups[c.selectedControlGroupID]

	var state structs.PublicRoom

	for _, m := range msg.Minions {
		// find the display by ID
		disp, err := getDisplay(cg, m)
		if err != nil {
			fmt.Printf("no!!!\n")
			return
		}
		var input string
		done := false
		for _, p := range c.uiConfig.Presets {
			for _, d := range p.Displays {
				if ID(d) == m {
					input = p.Inputs[0]
					done = true
					break
				}
			}
			if done {
				break
			}
		}

		for _, out := range disp.Outputs {

			dSplit := strings.Split(string(out.ID), "-")
			display := structs.Display{
				PublicDevice: structs.PublicDevice{
					Name: dSplit[2],
				},
			}

			if input == "blank" {
				display.Blanked = BoolP(true)
			} else {
				iSplit := strings.Split(string(input), "-")
				display.Input = iSplit[2]
				display.Blanked = BoolP(false)
			}

			state.Displays = append(state.Displays, display)
		}
	}

	err := c.SendAPIRequest(context.TODO(), state)
	if err != nil {
		c.Warn("failed to change input", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to change input: %s", err))
	}

	delete(c.sharing, msg.Master)

}
