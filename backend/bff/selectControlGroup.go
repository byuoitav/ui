package bff

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type SelectControlGroup struct {
}

type SelectControlGroupMessage struct {
	ID ID `json:"id"`
}

func (s SelectControlGroup) Do(c *Client, data []byte) {
	var msg SelectControlGroupMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.Warn("invalid value for selectControlGroup", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for selectControlGroup: %s", err))
		return
	}

	if err := s.DoWithMessage(c, msg); err != nil {
		c.Warn("failed to selectControlGroup", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to selectControlGroup: %w", err))
	}
}

func (s SelectControlGroup) DoWithMessage(c *Client, msg SelectControlGroupMessage) error {
	var id string

	if len(msg.ID) > 0 {
		room := c.GetRoom()

		// Validate that the control group is real
		for _, cg := range room.ControlGroups {
			if cg.ID == msg.ID {
				id = string(cg.ID)
				break
			}
		}

		if len(id) == 0 {
			return fmt.Errorf("invalid control group %q", msg.ID)
		}

		// turn on the control group if it's not already on
		if !room.ControlGroups[id].PoweredOn {
			preset, err := c.GetPresetByName(id)
			if err != nil {
				return fmt.Errorf("no matching preset found: %w", err)
			}

			err = preset.Actions.SetPower.DoWithMessage(c, SetPowerMessage{
				PoweredOn: true,
			})
			if err != nil {
				return fmt.Errorf("failed to turn on controlGroup: %w", err)
			}
		}
	}

	// Otherwise set the control group on the client
	c.selectedControlGroupID = id

	// And send the updated room to the front end
	roomMsg, err := JSONMessage("room", c.GetRoom())
	if err != nil {
		return fmt.Errorf("unable to create new room message: %w", err)
	}

	c.Out <- roomMsg
	return nil
}
