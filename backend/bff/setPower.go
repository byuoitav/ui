package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/byuoitav/common/structs"
	"go.uber.org/zap"
)

type SetPower struct {
}

type SetPowerMessage struct {
	PoweredOn bool `json:"poweredOn"`
	All       bool `json:"all"`
}

// Do .
func (sp SetPower) Do(c *Client, data []byte) {
	var msg SetPowerMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		c.Warn("invalid value for setPower", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setPower: %s", err))
		return
	}

	// this shouldn't take longer than 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// convert power status to the string the av-api wants
	status := "standby"
	if msg.PoweredOn {
		status = "on"
	}

	room := c.GetRoom()

	// build the state body
	var state structs.PublicRoom

	addDisplayGroup := func(group DisplayGroup) {
		// Dissolve share group it's a master
		if group.ShareInfo.State == stateIsMaster {
			// TODO might be better to inline the logic here, but this is easier for now
			var ss SetSharing
			ss.Unshare(c, SetSharingMessage{
				Group: group.ID,
			})
		}

		for _, disp := range group.Displays {
			state.Displays = append(state.Displays, structs.Display{
				PublicDevice: structs.PublicDevice{
					Name:  disp.ID.GetName(),
					Power: status,
				},
			})
		}
	}

	if msg.All {
		c.Info("Setting power for the entire room", zap.String("to", status))

		// set power on everything in the room
		for _, group := range room.GetAllDisplayGroups() {
			addDisplayGroup(group)
		}
	} else {
		// set power on everything in my controlGroup
		cg := room.ControlGroups[c.selectedControlGroupID]
		c.Info("Setting power", zap.String("to", status), zap.String("controlGroup", string(cg.ID)))

		for _, group := range cg.fullDisplayGroups {
			addDisplayGroup(group)
		}
	}

	if err := c.SendAPIRequest(ctx, state); err != nil {
		c.Warn("failed to set power", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to set power: %s", err))
	}

	c.Info("Finished setting power", zap.String("to", status), zap.Bool("entireRoom?", msg.All))
}
