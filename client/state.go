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
	state.Errors = nil

	c.stateMu.Lock()
	defer c.stateMu.Unlock()

	c.state = state
	return nil
}

// TODO this should (i think?) send events for the things that have changed
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

			cur.Inputs[iID] = curInput
		}

		for block, v := range d.Volumes {
			cur.Volumes[block] = v
		}

		for block, m := range d.Mutes {
			cur.Mutes[block] = m
		}

		c.state.Devices[dID] = cur
	}
}

// stateMatches assumes that (*client).stateMu has already been at least read locked.
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
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()

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
