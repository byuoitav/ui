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
	Displays []ID   `json:"displays"`
	Status   string `json:"status"`
}

func (sp SetPower) DoWithMessage(c *Client, msg SetPowerMessage) error {
	cg := c.GetRoom().ControlGroups[c.selectedControlGroupID]
	if len(cg.ID) == 0 {
		// error
		return fmt.Errorf("len(display.ID) is equal to zero")
	}

	for i := range msg.Displays {
		c.Info("setting Power", zap.String("on", fmt.Sprintf("%v", msg.Displays[i])), zap.String("to", msg.Status), zap.String("controlGroup", string(cg.ID)))
	}

	// find the display by ID
	var disp []Display
	for i := range cg.Displays {
		for j := range msg.Displays {
			if cg.Displays[i].ID == ID(msg.Displays[j]) {
				disp = append(disp, cg.Displays[i])
				break
			}
		}
	}

	if len(disp) <= 0 {
		// error
		fmt.Printf("no!!!\n")
		return fmt.Errorf("the display(s) are less than or equal to zero")
	}

	var state structs.PublicRoom
	for _, display := range disp {
		for _, out := range display.Outputs {
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
	}

	err := c.SendAPIRequest(state)
	if err != nil {
		c.Warn("failed to change input", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to change input: %s", err))
	}

	return nil
}

func (sp SetPower) Do(c *Client, data []byte) {
	var msg SetPowerMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setPower: %s", err))
		return
	}

	err = sp.DoWithMessage(c, msg)
	if err != nil {
		c.Out <- ErrorMessage(fmt.Errorf("error occured when calling DoWithMessage: %s", err))
		return
	}

}
