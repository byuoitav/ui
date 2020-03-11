package bff

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

// SelectControlGroup .
type SelectControlGroup struct {
}

// SelectControlGroupMessage .
type SelectControlGroupMessage struct {
	ID ID `json:"id"`
}

// Do .
func (s SelectControlGroup) Do(c *Client, data []byte) {
	var msg SelectControlGroupMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.Warn("invalid value for selectControlGroup", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for selectControlGroup: %s", err))
		return
	}

	var id string

	if len(msg.ID) > 0 {
		// Validate that the control group is real
		for _, cg := range c.GetRoom().ControlGroups {
			if cg.ID == msg.ID {
				id = string(cg.ID)
				break
			}
		}

		if len(id) == 0 {
			c.Warn("invalid control group", zap.String("id", string(msg.ID)))
			c.Out <- ErrorMessage(fmt.Errorf("invalid control group: %s", msg.ID))
			return
		}
	}

	// Otherwise set the control group on the client
	c.selectedControlGroupID = id

	// And send the updated room to the front end
	roomMsg, err := JSONMessage("room", c.GetRoom())
	if err != nil {
		c.Warn("unable to make new room message", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("unable to make new room message: %w", err))
		return
	}

	c.Out <- roomMsg
}
