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
	All       bool `json:"all"`
	PoweredOn bool `json:"poweredOn"`
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
	var state structs.PublicRoom
	if msg.All {
		for _, cg := range room.ControlGroups {
			c.Info("Setting power", zap.String("to", status), zap.String("controlGroup", string(cg.ID)))

			// go through all of the display groups and turn on all of their displays
			for _, group := range cg.DisplayGroups {
				for _, disp := range group.Displays {
					state.Displays = append(state.Displays, structs.Display{
						PublicDevice: structs.PublicDevice{
							Name:  disp.ID.GetName(),
							Power: status,
						},
					})
				}
			}
		}
	} else {
		cg := room.ControlGroups[c.selectedControlGroupID]
		c.Info("Setting power", zap.String("to", status), zap.String("controlGroup", string(cg.ID)))

		// go through all of the display groups and turn on all of their displays
		for _, group := range cg.DisplayGroups {
			for _, disp := range group.Displays {
				state.Displays = append(state.Displays, structs.Display{
					PublicDevice: structs.PublicDevice{
						Name:  disp.ID.GetName(),
						Power: status,
					},
				})
			}
		}
	}

	if err := c.SendAPIRequest(ctx, state); err != nil {
		c.Warn("failed to set power", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to set power: %s", err))
	}

	c.Info("Finished setting power", zap.String("to", status), zap.String("room", string(room.ID)))
}

// PowerOffAll .
//func (sp SetPower) PowerOffAll(c *Client) error {
//	controlGroups := c.GetRoom().ControlGroups
//	if controlGroups == nil {
//		// error
//		return fmt.Errorf("Control Groups not found %q", c.selectedControlGroupID)
//	}
//
//	c.Info("Powering off all devices in the room.")
//	var disp []DisplayBlock
//	for _, cg := range controlGroups {
//		c.Info("Powering off all devices in the room.")
//
//		for _, d := range cg.DisplayBlocks {
//			if !contains(disp, d) {
//				disp = append(disp, d)
//			}
//		}
//	}
//
//	if len(disp) <= 0 {
//		// error
//		fmt.Printf("no!!!\n")
//		return fmt.Errorf("the display(s) are less than or equal to zero")
//	}
//
//	var state structs.PublicRoom
//	for _, display := range disp {
//		for _, out := range display.Outputs {
//			// TODO write a getnamefromid func
//			dSplit := strings.Split(string(out.ID), "-")
//			display := structs.Display{
//				PublicDevice: structs.PublicDevice{
//					Name:  dSplit[2],
//					Power: "standby",
//				},
//			}
//
//			state.Displays = append(state.Displays, display)
//		}
//	}
//
//	err := c.SendAPIRequest(context.Background(), state)
//	if err != nil {
//		c.Warn("failed to set power", zap.Error(err))
//		c.Out <- ErrorMessage(fmt.Errorf("failed to set power: %s", err))
//	}
//
//	return nil
//}

//func contains(s []DisplayBlock, e DisplayBlock) bool {
//	for _, a := range s {
//		if a.ID == e.ID {
//			return true
//		}
//	}
//	return false
//}
