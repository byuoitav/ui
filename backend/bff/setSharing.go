package bff

//
//import (
//	"context"
//	"fmt"
//	"strings"
//
//	"github.com/byuoitav/common/structs"
//	"go.uber.org/zap"
//)
//
////
//func contain(l []ID, id ID) (int, bool) {
//	for i, e := range l {
//		if id == e {
//			return i, true
//		}
//	}
//	return -1, false
//}
//
//func getDisplayGroup(cg ControlGroup, id ID) (DisplayGroup, error) {
//	var d DisplayGroup
//	for _, dGroup := range cg.DisplayGroups {
//		if dGroup.ID == id {
//			return dGroup, nil
//		}
//	}
//
//	return d, fmt.Errorf("Could not find display")
//}
//
////
////func subArray(big []ID, small []ID) bool {
////	for _, v := range small {
////		if _, ok := contain(big, v); !ok {
////			return false
////		}
////	}
////	return true
////}
////
//func removeID(l []ID, index int) []ID {
//	l[index] = l[len(l)-1]
//	return l[:len(l)-1]
//}
//
//// Do sets the sharing state
//// Off Legacy
//func (ss SetSharing) Off(c *Client, msg SetSharingMessage) {
//	if shareMap := c.getShareMap(); shareMap != nil {
//
//		var state structs.PublicRoom
//
//		cg := c.GetRoom().ControlGroups[c.selectedControlGroupID]
//
//		mState := shareMap[msg.Master].State
//		switch mState {
//		case Unshare:
//			// Each minion should be able to share now
//			// And their input should be their default again
//			for _, minion := range shareMap[msg.Master].Active {
//				shareMap[minion] = ShareData{
//					State: Share,
//				}
//
//				dg, err := getDisplayGroup(cg, msg.Master)
//				if err != nil {
//					fmt.Printf("no!!!\n")
//					return
//				}
//				var input string
//				done := false
//				for _, p := range c.uiConfig.Presets {
//					for _, d := range p.Displays {
//						if ID(d) == msg.Master {
//							input = p.Inputs[0]
//							done = true
//							break
//						}
//					}
//					if done {
//						break
//					}
//				}
//
//				for _, out := range dg.Displays {
//
//					dSplit := strings.Split(string(out.ID), "-")
//					display := structs.Display{
//						PublicDevice: structs.PublicDevice{
//							Name: dSplit[2],
//						},
//					}
//
//					if input == "blank" {
//						display.Blanked = BoolP(true)
//					} else {
//						iSplit := strings.Split(string(input), "-")
//						display.Input = iSplit[2]
//						display.Blanked = BoolP(false)
//					}
//
//					state.Displays = append(state.Displays, display)
//				}
//
//			}
//
//			for _, minion := range shareMap[msg.Master].Inactive {
//				shareMap[minion] = ShareData{
//					State: Share,
//				}
//
//				dg, err := getDisplayGroup(cg, msg.Master)
//				if err != nil {
//					fmt.Printf("no!!!\n")
//					return
//				}
//				var input string
//				done := false
//				for _, p := range c.uiConfig.Presets {
//					for _, d := range p.Displays {
//						if ID(d) == msg.Master {
//							input = p.Inputs[0]
//							done = true
//							break
//						}
//					}
//					if done {
//						break
//					}
//				}
//
//				for _, out := range dg.Displays {
//
//					dSplit := strings.Split(string(out.ID), "-")
//					display := structs.Display{
//						PublicDevice: structs.PublicDevice{
//							Name: dSplit[2],
//						},
//					}
//
//					if input == "blank" {
//						display.Blanked = BoolP(true)
//					} else {
//						iSplit := strings.Split(string(input), "-")
//						display.Input = iSplit[2]
//						display.Blanked = BoolP(false)
//					}
//
//					state.Displays = append(state.Displays, display)
//				}
//
//			}
//			shareMap[msg.Master] = ShareData{
//				State: Share,
//			}
//		case MinionActive:
//			master := shareMap[msg.Master].Master
//			if index, ok := contain(shareMap[master].Active, msg.Master); ok {
//				inactive := append(shareMap[master].Inactive, shareMap[master].Active[index])
//				active := removeID(shareMap[master].Active, index)
//				shareMap[master] = ShareData{
//					State:    mState,
//					Active:   active,
//					Inactive: inactive,
//				}
//			}
//			dg, err := getDisplayGroup(cg, msg.Master)
//			if err != nil {
//				fmt.Printf("no!!!\n")
//				return
//			}
//			var input string
//			done := false
//			for _, p := range c.uiConfig.Presets {
//				for _, d := range p.Displays {
//					if ID(d) == msg.Master {
//						input = p.Inputs[0]
//						done = true
//						break
//					}
//				}
//				if done {
//					break
//				}
//			}
//
//			for _, out := range dg.Displays {
//
//				dSplit := strings.Split(string(out.ID), "-")
//				display := structs.Display{
//					PublicDevice: structs.PublicDevice{
//						Name: dSplit[2],
//					},
//				}
//
//				if input == "blank" {
//					display.Blanked = BoolP(true)
//				} else {
//					iSplit := strings.Split(string(input), "-")
//					display.Input = iSplit[2]
//					display.Blanked = BoolP(false)
//				}
//
//				state.Displays = append(state.Displays, display)
//			}
//		}
//
//		c.lazUpdates <- lazMessage{
//			Key:  lazSharingDisplays,
//			Data: shareMap,
//		}
//
//		// Send the updated room to the AV API
//		err := c.SendAPIRequest(context.TODO(), state)
//		if err != nil {
//			c.Warn("failed to change input", zap.Error(err))
//			c.Out <- ErrorMessage(fmt.Errorf("failed to change input: %s", err))
//		}
//	}
//
//}
//
////
////// Off Legacy
////func (ss SetSharing) Off(c *Client, msg SetSharingMessage) {
////
////	//TODO remove inactive minions from their list
////
////	cg := c.GetRoom().ControlGroups[c.selectedControlGroupID]
////
////	var state structs.PublicRoom
////	minions := c.sharing[msg.Master]
////	for _, m := range minions.Active {
////		// find the display by ID
////		disp, err := getDisplay(cg, m)
////		if err != nil {
////			fmt.Printf("no!!!\n")
////			return
////		}
////		var input string
////		done := false
////		for _, p := range c.uiConfig.Presets {
////			for _, d := range p.Displays {
////				if ID(d) == m {
////					input = p.Inputs[0]
////					done = true
////					break
////				}
////			}
////			if done {
////				break
////			}
////		}
////
////		for _, out := range disp.Outputs {
////
////			dSplit := strings.Split(string(out.ID), "-")
////			display := structs.Display{
////				PublicDevice: structs.PublicDevice{
////					Name: dSplit[2],
////				},
////			}
////
////			if input == "blank" {
////				display.Blanked = BoolP(true)
////			} else {
////				iSplit := strings.Split(string(input), "-")
////				display.Input = iSplit[2]
////				display.Blanked = BoolP(false)
////			}
////
////			state.Displays = append(state.Displays, display)
////		}
////	}
////
////	go updateLazSharing(context.TODO(), c)
////
////	err := c.SendAPIRequest(context.TODO(), state)
////	if err != nil {
////		c.Warn("failed to change input", zap.Error(err))
////		c.Out <- ErrorMessage(fmt.Errorf("failed to change input: %s", err))
////	}
////	c.shareMutex.Lock()
////	delete(c.sharing, msg.Master)
////	c.shareMutex.Unlock()
////
////}
