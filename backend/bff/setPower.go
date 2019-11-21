package bff

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/byuoitav/common/structs"
	"go.uber.org/zap"
)

type SetPower struct {
}

type SetPowerMessage struct {
	DisplayID string `json:"display"`
	Status    string `json:"status"`
}

func (sp SetPower) Do(c *Client, data []byte) {
	var msg SetPowerMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setPower: %s", err))
		return
	}

	cg := c.GetRoom().ControlGroups[c.selectedControlGroupID]
	if len(cg.ID) == 0 {
		// error
	}

	c.Info("setting Power", zap.String("on", msg.DisplayID), zap.String("to", msg.Status), zap.String("controlGroup", string(cg.ID)))

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
		display.Power = msg.Status

		state.Displays = append(state.Displays, display)
	}

	err = c.SendAPIRequest(state)
	if err != nil {
		c.Warn("failed to change input", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to change input: %s", err))
	}
}
