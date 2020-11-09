package client

import (
	"context"
	"fmt"

	avcontrol "github.com/byuoitav/av-control-api"
)

func (c *client) updateRoomState(ctx context.Context) error {
	state, err := c.avController.RoomState(ctx, c.roomID)
	if err != nil {
		return fmt.Errorf("unable to get state: %w", err)
	}

	// TODO something with the errors...?
	// TODO probably need to mutex the state...

	state.Errors = nil
	c.state = state
	return nil
}

func (c *client) stateMatches(req avcontrol.StateRequest) bool {
	for dID, d := range req.Devices {
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

func (c *client) getVolume(req avcontrol.StateRequest, state avcontrol.StateResponse) int {
	vols := []int{}

	for id, rDev := range req.Devices {
		sDev, ok := c.state.Devices[id]
		if !ok {
			continue
		}

		for block, _ := range rDev.Volumes {
			sVol, ok := sDev.Volumes[block]
			if !ok {
				continue
			}

			vols = append(vols, sVol)
		}
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

func (c *client) getMuted(req avcontrol.StateRequest, state avcontrol.StateResponse) bool {
	mutes := []bool{}

	for id, rDev := range req.Devices {
		sDev, ok := c.state.Devices[id]
		if !ok {
			continue
		}

		for block, _ := range rDev.Mutes {
			sMute, ok := sDev.Mutes[block]
			if !ok {
				continue
			}

			mutes = append(mutes, sMute)
		}
	}

	if len(mutes) == 0 {
		return false
	}

	for _, muted := range mutes {
		if !muted {
			return false
		}
	}

	return true
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
