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
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for setSharing: %w", err))
		return
	}

	// TODO we may need some kind of lock in lazarette to handle the case of
	// two people starting a share at the exact same time,
	// right now, whoever's request gets put into lazarette second will win

	// get the group we are talking about
	room := c.GetRoom()
	cg, ok := room.ControlGroups[c.selectedControlGroupID]
	if !ok {
		c.Out <- ErrorMessage(errors.New("sharing: selectedControlGroupID not in room.ControlGroups"))
		return
	}

	group, err := cg.DisplayGroups.GetDisplayGroup(msg.Group)
	if err != nil {
		c.Out <- ErrorMessage(errors.New("sharing: msg.Group not in cg.DisplayGroups"))
		return
	}

	switch group.ShareInfo.State {
	case stateCanShare:
		c.Info("Starting share", zap.String("master", string(msg.Group)), zap.Strings("minions", msg.Options))
		ss.Share(c, msg)
	case stateIsMaster:
		c.Info("Stopping share", zap.String("master", string(msg.Group)))
		ss.Unshare(c, msg)
	case stateIsActiveMinion:
		c.Info("Becoming an inactive minion", zap.String("group", string(msg.Group)))
		ss.becomeInactiveMinion(c, msg)
	}
}

// Share starts sharing
func (ss SetSharing) Share(c *Client, msg SetSharingMessage) {
	room := c.GetRoom()
	cg := room.ControlGroups[c.selectedControlGroupID]

	dgroups := room.GetAllDisplayGroups()

	// Find all of the display group names
	disps := make(map[ID]int)
	for i, name := range dgroups {
		fmt.Printf("\n%s\n", name.ID)
		disps[name.ID] = i
	}

	// get the current input that the master is on
	// and validate that master group id is valid
	var input string
	for i := range cg.DisplayGroups {
		if cg.DisplayGroups[i].ID == msg.Group {
			input = cg.DisplayGroups[i].Input.GetName()
			break
		}
	}

	if len(input) == 0 {
		err := errors.New("sharing: msg.Group not found in cg.DisplayGroups")
		c.Warn("failed to start share", zap.Error(err))
		c.Out <- ErrorMessage(err)
		return
	}

	/*
		// TODO do this validation only on blueberry
			preset, err := c.GetPresetByName(msg.Group.GetName())
			if err != nil {
				// If we can't find
				c.Warn("failed to start share", zap.Error(err))
				c.Out <- ErrorMessage(errors.New("sharing: no preset found for msg.Group"))
				return
			}
			// FOR BLUEBERRY
			// Find the set of shareable displays and see if it is a super set of msg.Options
			if len(room.ControlGroups) == 1 {
			}
			shareable := make(map[string]bool)

			for _, name := range preset.ShareableDisplays {
				shareable[name] = true
			}
	*/
	for _, id := range msg.Options {
		// validate that options exist in the room's display groups
		/*
			c.Out <- StringMessage(id, "")
			c.Out <- StringMessage(room.ID.GetName()+"-", "")
			c.Out <- StringMessage(strings.TrimPrefix(id, room.ID.GetName()+"-"), "")
		*/
		if _, ok := disps[ID(id)]; !ok {
			c.Warn("failed to start share", zap.Error(errors.New("sharing: option "+id+" not found in cg.DisplayGroups")))
			c.Out <- ErrorMessage(errors.New("sharing: option " + id + " not found in cg.DisplayGroups"))
			return
		}
		/*
			// validate that inputs are valid for minions by
			// validating that options exist in the preset's shareable displays
			if _, ok := shareable[strings.TrimPrefix(id, room.ID.GetName()+"-")]; !ok {
				c.Warn("failed to start share", zap.Error(err))
				c.Out <- ErrorMessage(errors.New("sharing: option " + id + " not found in preset.ShareableDisplays"))
				return
			}
		*/
	}

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
			c.Warn("failed to start share", zap.Error(err))
			c.Out <- ErrorMessage(fmt.Errorf("sharing: could not get display group by ID: %w", err))
			return
		}

		// Update all the minion displays to be the master input
		for _, disp := range mgroup.Displays {
			state.Displays = append(state.Displays, structs.Display{
				PublicDevice: structs.PublicDevice{
					Name:  disp.ID.GetName(),
					Input: input,
					Power: "on",
				},
			})
		}

		// Mute the minions
		mcg, err := GetControlGroupByDisplayGroupID(room.ControlGroups, mgroup.ID)
		if err != nil {
			c.Warn("failed to start share", zap.Error(err))
			c.Out <- ErrorMessage(fmt.Errorf("sharing: could not get control group by display group id: %w", err))
			return
		}
		preset, err := c.GetPresetByName(string(mcg.ID))
		if err != nil {
			c.Warn("failed to start share", zap.Error(err))
			c.Out <- ErrorMessage(fmt.Errorf("sharing: could not get preset by name: %w", err))
			return
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
		c.Warn("failed to start share", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to start share: %w", err))
		return
	}

	// let the frontend know that sharing is complete
	c.Out <- StringMessage("shareStarted", "")
}

// Unshare stops sharing
func (ss SetSharing) Unshare(c *Client, msg SetSharingMessage) {
	// reset everyone in my active group to their default inputs
	active, inactive := c.getActiveAndInactiveForDisplayGroup(msg.Group)
	if len(active) == 0 && len(inactive) == 0 {
		c.Warn(msg.Group.GetName() + " was sharing to nobody... (active and inactive lists both empty)")
	}

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
			c.Warn("failed to stop share", zap.Error(err))
			c.Out <- ErrorMessage(fmt.Errorf("sharing: could not get control group by ID for %s: %w", active[i], err))
			return
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
			c.Warn("failed to stop share", zap.Error(err))
			c.Out <- ErrorMessage(fmt.Errorf("sharing: could not get display group by ID for %s: %w", active[i], err))
			return
		}

		mcg, err := GetControlGroupByDisplayGroupID(room.ControlGroups, mgroup.ID)
		if err != nil {
			c.Warn("failed to stop share", zap.Error(err))
			c.Out <- ErrorMessage(fmt.Errorf("sharing: could not get control group by display group id for %s: %w", mgroup.ID, err))
			return
		}

		preset, err := c.GetPresetByName(string(mcg.ID))
		if err != nil {
			c.Warn("failed to stop share", zap.Error(err))
			c.Out <- ErrorMessage(fmt.Errorf("sharing: could not get preset by name for %s: %w", mgroup.ID, err))
			return
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
	// make sure we even need to do this
	if len(state.Displays) == 0 && len(state.AudioDevices) == 0 {
		c.Out <- StringMessage("nothing to unshare", "")
		c.Out <- StringMessage("shareEnded", "")
		return

	}
	if err := c.SendAPIRequest(ctx, state); err != nil {
		c.Warn("failed to end share", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to end share: %w", err))
		return
	}

	c.Out <- StringMessage("shareEnded", "")
}

// becomeInactiveMinion causes a minion in a share group to become inactive. Becoming inactive does a few things on the UI:
// 1. The modal saying that you are being shared to disappears
// 2. A new input shows up to rejoin the share group (see (*client).GetRoom())
// 3. The state of the display group is set to muted=false, blanked=false, and input=default
//
// This function is only ever called from blueberry.
// TODO: we should probably validate that it is blueberry.
func (ss SetSharing) becomeInactiveMinion(c *Client, msg SetSharingMessage) {
	room := c.GetRoom()
	cg := room.ControlGroups[c.selectedControlGroupID]

	// get the group
	group, err := cg.DisplayGroups.GetDisplayGroup(msg.Group)
	if err != nil {
		c.Warn("failed to become an inactive minion", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("sharing: could not get display group for %s: %w", msg.Group, err))
		return
	}

	// build the av-api state
	var state structs.PublicRoom

	// change the displays
	for i := range group.Displays {
		state.Displays = append(state.Displays, structs.Display{
			PublicDevice: structs.PublicDevice{
				Name:  group.Displays[i].ID.GetName(),
				Input: cg.Inputs[0].ID.GetName(),
			},
			Blanked: BoolP(false),
		})
	}

	// change the audio devices
	audioDevices, err := cg.GetMediaAudioDeviceIDs(c.uiConfig.Presets)
	if err != nil {
		c.Warn("failed to become an inactive minion", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("sharing: could not get media audio device ids: %w", err))
		return
	}

	for i := range audioDevices {
		state.AudioDevices = append(state.AudioDevices, structs.AudioDevice{
			PublicDevice: structs.PublicDevice{
				Name: group.Displays[i].ID.GetName(),
			},
			Muted: BoolP(false),
		})
	}

	// update lazarette
	c.lazUpdates <- lazMessage{
		Key: lazSharing + string(msg.Group),
		Data: lazShareData{
			State:  stateIsInactiveMinion,
			Master: group.ShareInfo.Master,
		},
	}

	// don't take longer than 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// send the api request
	if err := c.SendAPIRequest(ctx, state); err != nil {
		c.Warn("failed to become an inactive minion", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("failed to become an inactive minion: %w", err))
		return
	}

	c.Out <- StringMessage("becameInactive", "")
}
