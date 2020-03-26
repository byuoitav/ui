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

func (sp SetPower) Do(c *Client, data []byte) {
	var msg SetPowerMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		c.Warn("invalid value for setPower", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setPower: %s", err))
		return
	}

	if err := sp.DoWithMessage(c, msg); err != nil {
		c.Warn("failed to setPower", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to setPower: %w", err))
	}
}

func (sp SetPower) DoWithMessage(c *Client, msg SetPowerMessage) error {
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

	addDisplayGroup := func(group DisplayGroup) error {
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

		// set the media volume to 30 if we are turning it on
		if msg.PoweredOn {
			preset, err := c.GetPresetByName(string(group.ID))
			if err != nil {
				return fmt.Errorf("unable to set media audio level: %w", err)
			}

			// set media audio
			for i := range preset.AudioDevices {
				state.AudioDevices = append(state.AudioDevices, structs.AudioDevice{
					PublicDevice: structs.PublicDevice{
						Name: preset.AudioDevices[i],
					},
					Volume: IntP(30),
				})
			}
		}

		return nil
	}

	if msg.All {
		c.Info("Setting power for the entire room", zap.String("to", status))

		// set power on everything in the room
		for _, group := range room.GetAllDisplayGroups() {
			err := addDisplayGroup(group)
			if err != nil {
				return err
			}
		}
	} else {
		// set power on everything in my controlGroup
		cg := room.ControlGroups[c.selectedControlGroupID]
		c.Info("Setting power", zap.String("to", status), zap.String("controlGroup", string(cg.ID)))

		for _, group := range cg.fullDisplayGroups {
			err := addDisplayGroup(group)
			if err != nil {
				return err
			}
		}
	}

	if err := c.SendAPIRequest(ctx, state); err != nil {
		return err
	}

	c.Info("Finished setting power", zap.String("to", status), zap.Bool("entireRoom?", msg.All))
	return nil
}
