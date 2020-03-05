package bff

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/byuoitav/common/structs"
	"go.uber.org/zap"
)

type shareState int

const (
	stateCantShare shareState = iota
	stateCanShare
	stateIsMaster
	stateIsActiveMinion
	stateIsInactiveMinion
)

const (
	lazSharing = "-sharing-"
)

// the types of objects that are stored in lazarette for each display group
type lazShareData struct {
	State shareState `json:"state"`

	// if they are in a group, this is who that leader is
	Master ID `json:"master,omitempty"`
}

func (c *Client) getActiveAndInactiveForDisplayGroup(group ID) ([]ID, []ID) {
	var active []ID
	var inactive []ID

	c.lazs.Range(func(key, value interface{}) bool {
		skey, ok := key.(string)
		if !ok || !strings.HasPrefix(skey, lazSharing) {
			return true
		}

		shareData, ok := value.(lazShareData)
		if !ok {
			return true
		}

		// make sure it's in this group
		if shareData.Master != group {
			return true
		}

		// strip prefix off of id
		id := strings.TrimPrefix(skey, lazSharing)

		// add it to active/inactive list
		switch shareData.State {
		case stateIsActiveMinion:
			active = append(active, ID(id))
		case stateIsInactiveMinion:
			inactive = append(inactive, ID(id))
		}

		return true
	})

	return active, inactive
}

func (c *Client) getShareData(group ID) (lazShareData, error) {
	var data lazShareData

	idata, ok := c.lazs.Load(lazSharing + string(group))
	if !ok {
		return data, errors.New("no share data found")
	}

	data, ok = idata.(lazShareData)
	if !ok {
		return data, errors.New("unexpected type")
	}

	return data, nil
}

// SetSharing .
type SetSharing struct {
}

// SetSharingMessage is a message enabling or disabling sharing
type SetSharingMessage struct {
	Group   ID       `json:"group,omitempty"`
	Options []string `json:"opts,omitempty"`
}

// Do .
func (ss SetSharing) Do(c *Client, data []byte) {
	var msg SetSharingMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setSharing: %s", err))
		return
	}

	// TODO we may need some kind of lock in lazarette to handle the case of
	// two people starting a share at the exact same time,
	// right now, whoever's request gets put into lazarette second will win

	// get the group we are talking about
	room := c.GetRoom()
	cg := room.ControlGroups[c.selectedControlGroupID]
	// TODO make sure cg is not nil

	group, err := cg.DisplayGroups.GetDisplayGroup(msg.Group)
	if err != nil {
		// handle err
	}

	switch group.ShareInfo.State {
	case stateCanShare:
		c.Info("Starting share", zap.String("master", string(msg.Group)), zap.Strings("minions", msg.Options))
		ss.Share(c, msg)
	case stateIsMaster:
		c.Info("Stopping share", zap.String("master", string(msg.Group)))
		ss.Unshare(c, msg)
	case stateIsActiveMinion:
		c.Info("Leaving share", zap.String("group", string(msg.Group)))
		// ss.LeaveShare()
	}
}

// Share starts sharing
func (ss SetSharing) Share(c *Client, msg SetSharingMessage) {
	room := c.GetRoom()
	cg := room.ControlGroups[c.selectedControlGroupID]

	dgroups := room.GetAllDisplayGroups()

	// get the current input that the master is on
	var input string
	for i := range cg.DisplayGroups {
		if cg.DisplayGroups[i].ID == msg.Group {
			input = cg.DisplayGroups[i].Input.GetName()
		}
	}

	if len(input) == 0 {
		// handle
	}

	// blanked?
	// validate that options are valid
	// validate that master group id is valid
	// validate that inputs are valid for minions

	var state structs.PublicRoom
	toMute := make(map[string]structs.AudioDevice)
	for _, minion := range msg.Options {
		c.lazUpdates <- lazMessage{
			Key: lazSharing + minion,
			Data: lazShareData{
				State:  stateIsActiveMinion,
				Master: msg.Group,
			},
		}

		mgroup, err := GetDisplayGroupByID(dgroups, ID(minion))
		if err != nil {
			// handle
		}

		// Update all the minion displays to be the master input
		for _, disp := range mgroup.Displays {
			state.Displays = append(state.Displays, structs.Display{
				PublicDevice: structs.PublicDevice{
					Name:  disp.ID.GetName(),
					Input: input,
				},
			})
		}

		// Mute the minions
		mcg, err := GetControlGroupByDisplayGroupID(room.ControlGroups, mgroup.ID)
		preset, err := c.GetPresetByName(string(mcg.ID))
		if err != nil {
			// handle
		}
		for _, audio := range preset.AudioDevices {
			if ID(audio) == msg.Group {
				continue
			}
			toMute[audio] = structs.AudioDevice{
				PublicDevice: structs.PublicDevice{
					Name: audio,
				},
				Muted: BoolP(true),
			}

		}
	}

	// Actually mute the audio devices
	for _, dev := range toMute {
		state.AudioDevices = append(state.AudioDevices, dev)
	}

	// update the masters lazarette data
	c.lazUpdates <- lazMessage{
		Key: lazSharing + string(msg.Group),
		Data: lazShareData{
			State: stateIsMaster,
		},
	}

	// don't take longer than 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// send the api request
	if err := c.SendAPIRequest(ctx, state); err != nil {
		c.Warn("failed to set sharing", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to set sharing: %w", err))
	}

	// let the frontend know that sharing is complete
	c.Out <- StringMessage("shareStarted", "")
}

// Unshare stops sharing
func (ss SetSharing) Unshare(c *Client, msg SetSharingMessage) {
	// reset everyone in my active group to their default inputs
	active, inactive := c.getActiveAndInactiveForDisplayGroup(msg.Group)

	var state structs.PublicRoom
	room := c.GetRoom()
	dgroups := room.GetAllDisplayGroups()

	toUnmute := make(map[string]structs.AudioDevice)
	// process the active devices
	for i := range active {
		// reset its state
		c.lazUpdates <- lazMessage{
			Key: lazSharing + string(active[i]),
			Data: lazShareData{
				State: stateCanShare,
			},
		}

		// get the default input for this group
		cg, err := GetControlGroupByDisplayGroupID(room.ControlGroups, active[i])
		if err != nil {
			// handle
		}

		state.Displays = append(state.Displays, structs.Display{
			PublicDevice: structs.PublicDevice{
				Name:  active[i].GetName(),
				Input: cg.Inputs[0].ID.GetName(),
			},
			Blanked: BoolP(false),
		})

		mgroup, err := GetDisplayGroupByID(dgroups, active[i])
		if err != nil {
			// handle
		}

		mcg, err := GetControlGroupByDisplayGroupID(room.ControlGroups, mgroup.ID)
		preset, err := c.GetPresetByName(string(mcg.ID))
		if err != nil {
			// handle
		}
		for _, audio := range preset.AudioDevices {
			if ID(audio) == msg.Group {
				continue
			}
			toUnmute[audio] = structs.AudioDevice{
				PublicDevice: structs.PublicDevice{
					Name: audio,
				},
				Muted: BoolP(false),
			}

		}
	}

	// Actually unmute the audio devices
	for _, dev := range toUnmute {
		state.AudioDevices = append(state.AudioDevices, dev)
	}

	// process the inactive devices
	for i := range inactive {
		// reset its state
		c.lazUpdates <- lazMessage{
			Key: lazSharing + string(inactive[i]),
			Data: lazShareData{
				State: stateCanShare,
			},
		}

		// don't need to change it's state at all!
		// just remove it from the group.
	}

	// reset the masters state
	c.lazUpdates <- lazMessage{
		Key: lazSharing + string(msg.Group),
		Data: lazShareData{
			State: stateCanShare,
		},
	}

	// don't take longer than 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// send the api request
	// TODO make sure we even need to do this
	if err := c.SendAPIRequest(ctx, state); err != nil {
		c.Warn("failed to set sharing", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to set sharing: %w", err))
	}

	c.Out <- StringMessage("shareEnded", "")
}
