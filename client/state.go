package client

import (
	"context"
	"fmt"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/ui"
	"go.uber.org/zap"
)

// TODO add timeout
func (c *client) updateRoomState(ctx context.Context) error {
	c.log.Debug("Updating room state")

	state, err := c.avController.RoomState(ctx, c.roomID)
	if err != nil {
		c.log.Error("unable to update room state", zap.Error(err))
		return fmt.Errorf("unable to get state: %w", err)
	}

	for i := range state.Errors {
		c.log.Warn("error in API response", zap.Any("stateError", state.Errors[i]))
	}

	c.log.Debug("Successfully updated room state")
	state.Errors = nil

	c.stateMu.Lock()
	defer c.stateMu.Unlock()

	c.state = state
	return nil
}

// TODO remove devices from c.state? or do we just not care?
func (c *client) updateRoomStateFromState(state avcontrol.StateResponse) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()

	for dID, d := range state.Devices {
		cur := c.state.Devices[dID]

		if d.PoweredOn != nil {
			cur.PoweredOn = d.PoweredOn
		}

		if d.Blanked != nil {
			cur.Blanked = d.Blanked
		}

		for iID, i := range d.Inputs {
			curInput := cur.Inputs[iID]

			if i.Audio != nil {
				curInput.Audio = i.Audio
			}

			if i.Video != nil {
				curInput.Video = i.Video
			}

			if i.AudioVideo != nil {
				curInput.AudioVideo = i.AudioVideo
			}

			if cur.Inputs == nil {
				cur.Inputs = make(map[string]avcontrol.Input)
			}
			cur.Inputs[iID] = curInput
		}

		for block, v := range d.Volumes {
			if cur.Volumes == nil {
				cur.Volumes = make(map[string]int)
			}
			cur.Volumes[block] = v
		}

		for block, m := range d.Mutes {
			if cur.Mutes == nil {
				cur.Mutes = make(map[string]bool)
			}
			cur.Mutes[block] = m
		}

		if c.state.Devices == nil {
			c.state.Devices = make(map[avcontrol.DeviceID]avcontrol.DeviceState)
		}
		c.state.Devices[dID] = cur
	}
}

// curStates returns a map of all the current states to whether or not the room is in that state.
// If lock is true, curStates() will read lock (*client).configMu and (*client.stateMu) before
// matching states. If it's false, the caller must have a read or write lock on those.
func (c *client) curStates(lock bool) map[string]bool {
	if lock {
		c.stateMu.RLock()
		defer c.stateMu.RUnlock()

		c.configMu.RLock()
		defer c.configMu.RUnlock()
	}

	check := func(state ui.State) bool {
		for dID, d := range state.Devices {
			dd, ok := c.state.Devices[dID]
			if !ok {
				return false
			}

			if !boolMatches(d.PoweredOn, dd.PoweredOn) {
				return false
			}

			if !boolMatches(d.Blanked, dd.Blanked) {
				return false
			}

			for iID, i := range d.Inputs {
				ii, ok := dd.Inputs[iID]
				if !ok {
					return false
				}

				if !stringMatches(i.Audio, ii.Audio) {
					return false
				}

				if !stringMatches(i.Video, ii.Video) {
					return false
				}

				if !stringMatches(i.AudioVideo, ii.AudioVideo) {
					return false
				}
			}

			for block, v := range d.Volumes {
				vv, ok := dd.Volumes[block]
				if !ok {
					return false
				}

				if v != vv {
					return false
				}
			}

			for block, m := range d.Mutes {
				mm, ok := dd.Mutes[block]
				if !ok {
					return false
				}

				if m != mm {
					return false
				}
			}
		}

		return true
	}

	states := make(map[string]bool, len(c.config.States))
	for name, state := range c.config.States {
		states[name] = check(state)
	}

	return states
}

// matchStates returns true if all of matchStates is true in states.
func (c *client) matchStates(states map[string]bool, matchStates []string) bool {
	for _, match := range matchStates {
		if !states[match] {
			return false
		}
	}

	return true
}

func (c *client) mergeStates(states ...string) ui.State {
	if len(states) == 0 {
		return ui.State{}
	}

	c.stateMu.RLock()
	defer c.stateMu.RUnlock()

	c.configMu.RLock()
	defer c.configMu.RUnlock()

	base := ui.State{
		Devices: make(map[avcontrol.DeviceID]avcontrol.DeviceState),
	}

	for _, s := range states {
		state, ok := c.config.States[s]
		if !ok {
			continue
		}

		for dID, d := range state.Devices {
			cur := base.Devices[dID]

			if d.PoweredOn != nil {
				cur.PoweredOn = d.PoweredOn
			}

			if d.Blanked != nil {
				cur.Blanked = d.Blanked
			}

			if d.Inputs != nil && cur.Inputs == nil {
				cur.Inputs = make(map[string]avcontrol.Input)
			}

			if d.Volumes != nil && cur.Volumes == nil {
				cur.Volumes = make(map[string]int)
			}

			if d.Mutes != nil && cur.Mutes == nil {
				cur.Mutes = make(map[string]bool)
			}

			for iID, i := range d.Inputs {
				curInput := cur.Inputs[iID]

				if i.Audio != nil {
					curInput.Audio = i.Audio
				}

				if i.Video != nil {
					curInput.Video = i.Video
				}

				if i.AudioVideo != nil {
					curInput.AudioVideo = i.AudioVideo
				}

				cur.Inputs[iID] = curInput
			}

			for block, v := range d.Volumes {
				cur.Volumes[block] = v
			}

			for block, m := range d.Mutes {
				cur.Mutes[block] = m
			}

			base.Devices[dID] = cur
		}
	}

	return base
}

// getVolumes assumes that (*client).stateMu and (*client).configMu have already been read locked.
func (c *client) getVolume(matchStates []string) int {
	getVols := func(s string) []int {
		state, ok := c.config.States[s]
		if !ok {
			return []int{}
		}

		var vols []int

		for id, dev := range state.Devices {
			sDev, ok := c.state.Devices[id]
			if !ok {
				continue
			}

			for block, tmpl := range dev.Volumes {
				// only count devices who's request is configured with -1
				if tmpl != -1 {
					continue
				}

				vol, ok := sDev.Volumes[block]
				if !ok {
					continue
				}

				vols = append(vols, vol)
			}
		}

		return vols
	}

	var vols []int
	for _, state := range matchStates {
		vols = append(vols, getVols(state)...)
	}

	if len(vols) == 0 {
		return -1
	}

	avg := vols[0]
	for i := 1; i < len(vols); i++ {
		avg += vols[i]
		avg /= 2
	}

	return avg
}

// stringMatches returns true if:
// - a is nil
// - a is not nil, b is not nil, and their values are the same
func stringMatches(a, b *string) bool {
	switch {
	case a == nil:
		return true
	case b != nil && *a == *b:
		return true
	}

	return false
}

// boolMatches returns true if:
// - a is nil
// - a is not nil, b is not nil, and their values are the same
func boolMatches(a, b *bool) bool {
	switch {
	case a == nil:
		return true
	case b != nil && *a == *b:
		return true
	}

	return false
}
