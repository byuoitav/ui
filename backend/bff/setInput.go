package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/byuoitav/common/structs"
	"go.uber.org/zap"
)

type SetInput struct {
	OnSameInput HttpRequest `json:"onSameInput"`
}

type SetInputMessage struct {
	DisplayID string `json:"display"`
	InputID   string `json:"input"`
}

func (si SetInput) Do(c *Client, data []byte) {
	var msg SetInputMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setInput: %s", err))
		return
	}

	cg := c.GetRoom().ControlGroups[c.selectedControlGroupID]
	if len(cg.ID) == 0 {
		// error
	}

	c.Info("setting input", zap.String("on", msg.DisplayID), zap.String("to", msg.InputID), zap.String("controlGroup", string(cg.ID)))

	// find the display by ID
	var disp Display
	for i := range cg.Displays {
		if cg.Displays[i].ID == ID(msg.DisplayID) {
			disp = cg.Displays[i]
			break
		}
	}

	if len(disp.ID) <= 0 {
		// error
		fmt.Printf("no!!!\n")
		return
	}

	var state structs.PublicRoom
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

	if len(si.OnSameInput.URL) > 0 {
		// send mute request
	}

	err = c.SendAPIRequest(context.TODO(), state)
	if err != nil {
		c.Warn("failed to change input", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to change input: %s", err))
	}
}
