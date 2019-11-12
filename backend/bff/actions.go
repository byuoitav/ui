package bff

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/byuoitav/common/structs"
	"go.uber.org/zap"
)

type SetInput struct {
	OnSameInput HttpRequest `json:"onSameInput"`
}

func (si SetInput) Do(c *Client, data []byte) {
	var msg SetInputMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		c.Out <- ErrorMessage("invalid value for setInput: %s", err)
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
		iSplit := strings.Split(string(msg.InputID), "-")

		state.Displays = append(state.Displays, structs.Display{
			PublicDevice: structs.PublicDevice{
				Name:  dSplit[2],
				Input: iSplit[2], // TODO do i need to get the name?
			},
		})
	}

	if len(si.OnSameInput.URL) > 0 {
		// send mute request
	}

	err = c.SendAPIRequest(state)
	if err != nil {
		c.Warn("failed to change input", zap.Error(err))
		c.Out <- ErrorMessage("failed to change input: %s", err)
	}
}
