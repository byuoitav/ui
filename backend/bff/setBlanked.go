package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/byuoitav/common/structs"
	"go.uber.org/zap"
)

// SetBlanked .
type SetBlanked struct {
}

// SetBlankedMessage is a message on who to (un)blank
type SetBlankedMessage struct {
	DisplayGroup ID   `json:"displayGroup"`
	Blanked      bool `json:"blanked"`
}

// Do .
func (sb SetBlanked) Do(c *Client, data []byte) {
	var msg SetBlankedMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		c.Warn("invalid value for setBlanked", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setBlanked: %s", err))
		return
	}

	// this shouldn't take longer than 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	room := c.GetRoom()
	cg := room.ControlGroups[c.selectedControlGroupID]
	c.Info("Setting blanked", zap.String("on", string(msg.DisplayGroup)), zap.Bool("to", msg.Blanked), zap.String("controlGroup", string(cg.ID)))

	// find the display group
	group, err := GetDisplayGroupByID(cg.DisplayGroups, msg.DisplayGroup)
	if err != nil {
		c.Warn("failed to set blanked", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to set blanked: %s", err))
		return
	}

	var state structs.PublicRoom

	switch group.ShareInfo.State {
	case stateIsActiveMinion:
		// TODO this is illegal, return an error
	case stateIsMaster:
		// get every display group in this room
		allGroups := room.GetAllDisplayGroups()

		// get all of my active minions
		active, _ := c.getActiveAndInactiveForDisplayGroup(msg.DisplayGroup)

		c.Info("Setting blanked as sharing master", zap.String("displayGroup", string(msg.DisplayGroup)), zap.Bool("blanked", msg.Blanked), zap.Strings("activeMinions", IDsToStrings(active)))

		// go through this display group and set all of it's displays blanked status
		for i := range active {
			mgroup, err := allGroups.GetDisplayGroup(active[i])
			if err != nil {
				// invalid display group id in active list
				// TODO validate active list?
				continue
			}
			for _, disp := range mgroup.Displays {
				state.Displays = append(state.Displays, structs.Display{
					PublicDevice: structs.PublicDevice{
						Name: disp.ID.GetName(),
					},
					Blanked: BoolP(msg.Blanked),
				})

			}
		}
	default:
		c.Info("Setting blanked", zap.String("on", string(msg.DisplayGroup)), zap.Bool("to", msg.Blanked), zap.String("controlGroup", string(cg.ID)))
	}

	// go through this display group and set all of it's displays blanked status
	for _, disp := range group.Displays {
		display := structs.Display{
			PublicDevice: structs.PublicDevice{
				Name: disp.ID.GetName(),
			},
			Blanked: BoolP(msg.Blanked),
		}

		// Add each display to the list of displays to change on the new state
		state.Displays = append(state.Displays, display)
	}

	if err := c.SendAPIRequest(ctx, state); err != nil {
		c.Warn("failed to set blanked", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to set blanked: %s", err))
	}

	c.Info("Finished setting blanked", zap.String("on", string(msg.DisplayGroup)), zap.Bool("to", msg.Blanked), zap.String("controlGroup", string(cg.ID)))
}
