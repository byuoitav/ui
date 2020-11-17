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

func (c *client) getVolume(req avcontrol.StateRequest) int {
	vols := []int{}

	for id, rDev := range req.Devices {
		sDev, ok := c.state.Devices[id]
		if !ok {
			continue
		}

		for block, rDevVol := range rDev.Volumes {
			// only count devices who's request is configured with -1
			if rDevVol != -1 {
				continue
			}

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

func fillVolumeRequest(req avcontrol.StateRequest, vol int) avcontrol.StateRequest {
	for _, rDev := range req.Devices {
		for block, rDevVol := range rDev.Volumes {
			if rDevVol != -1 {
				continue
			}

			rDev.Volumes[block] = vol
		}
	}

	return req
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
