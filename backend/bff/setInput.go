package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/byuoitav/common/structs"
	"go.uber.org/zap"
)

const (
	inputBecomeActivePrefix = "becomeActive_"
)

// HTTPRequest .
type HTTPRequest struct {
	Method string          `json:"method"`
	URL    string          `json:"url"`
	Body   json.RawMessage `json:"body"`
}

// SetInput .
type SetInput struct {
	OnSameInput HTTPRequest `json:"onSameInput"`
}

// SetInputMessage .
type SetInputMessage struct {
	DisplayGroup ID `json:"displayGroup"`
	Input        ID `json:"input"`
}

// Do .
func (si SetInput) Do(c *Client, data []byte) {
	var msg SetInputMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		c.Warn("invalid value for setInput", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setInput: %s", err))
		return
	}

	room := c.GetRoom()
	cg := room.ControlGroups[c.selectedControlGroupID]

	// find the display group by ID
	group, err := GetDisplayGroupByID(cg.DisplayGroups, msg.DisplayGroup)
	if err != nil {
		c.Warn("failed to change input", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to change input: %s", err))
		return
	}

	// build the state object
	var state structs.PublicRoom
	input := msg.Input.GetName()

	// validate input is valid for myself
	// validate input is valid for all minions

	switch group.ShareInfo.State {
	case stateIsActiveMinion:
		// TODO this is illegal, return an error
	case stateIsMaster:
		// get every display group in this room
		allGroups := room.GetAllDisplayGroups()

		// get all of my active minions
		active, _ := c.getActiveAndInactiveForDisplayGroup(msg.DisplayGroup)

		c.Info("Setting input as sharing master", zap.String("displayGroup", string(msg.DisplayGroup)), zap.String("input", string(msg.Input)), zap.Strings("activeMinions", IDsToStrings(active)))

		// set each active minion's input
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
						Name:  disp.ID.GetName(),
						Input: input,
					},
					Blanked: BoolP(false),
				})
			}
		}
	case stateIsInactiveMinion:
		// TODO if input == my masters id, then i should switch into their active group
		// c.Info("Rejoining share group", zap.String("displayGroup", string(msg.DisplayGroup)))
		// and change input to whatever the master's input is
		// and update lazarette!
	default:
		c.Info("Setting input", zap.String("displayGroup", string(msg.DisplayGroup)), zap.String("input", string(msg.Input)))
	}

	// figure out share stuff
	//if shareMap := c.getShareMap(); shareMap != nil {
	//	if data, ok := shareMap[msg.DisplayGroup]; ok {
	//		switch data.State {
	//		case MinionInactive:
	//			swap := data.Master == msg.Input
	//			if swap {
	//				// we are switching into the active list
	//				masterGroup, err := GetDisplayGroupByID(cg.DisplayGroups, data.Master)
	//				if err != nil {
	//					// TODO
	//				}
	//				inputName = masterGroup.Input.GetName()
	//				master := shareMap[data.Master]
	//				if index, ok := contain(master.Inactive, msg.DisplayGroup); ok {
	//					// and add me to the active list
	//					master.Active = append(master.Active, msg.DisplayGroup)
	//					// remove me from the inactive list
	//					master.Inactive = removeID(master.Inactive, index)
	//					shareMap[data.Master] = master
	//					c.lazUpdates <- lazMessage{ Key:  lazSharingDisplays,
	//						Data: shareMap,
	//					}
	//				} else {
	//					// TODO  this is an error too
	//					// Since you are an inactive minion, you should be on the inactive list
	//				}
	//			}
	//		}
	//	}
	//}

	// change input for all of my displays
	for _, disp := range group.Displays {
		state.Displays = append(state.Displays, structs.Display{
			PublicDevice: structs.PublicDevice{
				Name:  disp.ID.GetName(),
				Input: input,
			},
			Blanked: BoolP(false),
		})
	}

	// this shouldn't take longer than 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// make the state changes
	if err := c.SendAPIRequest(ctx, state); err != nil {
		c.Warn("failed to change input", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to change input: %s", err))
	}

	c.Info("Finished setInput", zap.String("displayGroup", string(msg.DisplayGroup)), zap.String("input", string(msg.Input)))
}
