package ui

import (
	"net/url"

	avcontrol "github.com/byuoitav/av-control-api"
)

// Config represents the program for a room that can be used
// to control a room
type Config struct {
	ID            string
	ControlGroups map[string]ControlGroup
}

// ControlGroup represents a group of Devices and inputs. These groups
// are used for logical grouping and displaying on different UIs
type ControlGroup struct {
	// PowerOff represents the state the room needs to be in to be considered "off"
	PowerOff ControlSet

	// PowerOn is the state we set to turn on the room
	PowerOn ControlSet

	Displays []DisplayConfig
	Audio    AudioConfig
	Cameras  []CameraConfig
}

// CameraConfig represents a Camera and its associated control endpoints
type CameraConfig struct {
	Name        string
	TiltUp      ControlSet
	TiltDown    ControlSet
	PanLeft     ControlSet
	PanRight    ControlSet
	PanTiltStop ControlSet
	ZoomIn      ControlSet
	ZoomOut     ControlSet
	ZoomStop    ControlSet
	Presets     []CameraPresetConfig
}

// CameraPresetConfig represents a preset on a camera
type CameraPresetConfig struct {
	Name      string
	SetPreset ControlSet
}

// DisplayConfig represents a Display and its associated controls
// for a given group
type DisplayConfig struct {
	Name    string
	Icon    string
	Sources []SourceConfig
}

// SourceConfig represents a source and its associated controls
type SourceConfig struct {
	Name    string
	Icon    string
	Visible bool
	ControlSet

	// Sources represent sub-sources of the parent source
	Sources []SourceConfig
}

// AudioConfig contains information about audio controls in the room
type AudioConfig struct {
	Media AudioDeviceConfig

	Groups []AudioGroupConfig
}

type AudioGroupConfig struct {
	Name         string
	AudioDevices []AudioDeviceConfig
}

type AudioDeviceConfig struct {
	Name   string
	Volume ControlSet `json:"volume"`
	Mute   ControlSet `json:"mute"`
	Unmute ControlSet `json:"unmute"`
}

// ControlSet represents the request to be made (both to the
// AV Control API and other arbitrary locations) in order to set the room
// to a given state
type ControlSet struct {
	APIRequest avcontrol.StateRequest
	Requests   []GenericControlRequest
}

// GenericControlRequest contains the information necessary to make a generic
// HTTP request in association with making state changes in a room
type GenericControlRequest struct {
	URL    *url.URL
	Method string
	Body   []byte
}

func (cs *ControlSet) Copy() *ControlSet {
	if cs == nil {
		return nil
	}

	n := &ControlSet{}

	if cs.APIRequest.Devices != nil {
		n.APIRequest.Devices = make(map[avcontrol.DeviceID]avcontrol.DeviceState, len(cs.APIRequest.Devices))

		for id, state := range cs.APIRequest.Devices {
			nState := avcontrol.DeviceState{}

			if state.PoweredOn != nil {
				b := *state.PoweredOn
				nState.PoweredOn = &b
			}

			if state.Blanked != nil {
				b := *state.Blanked
				nState.Blanked = &b
			}

			if state.Inputs != nil {
				nState.Inputs = make(map[string]avcontrol.Input, len(state.Inputs))

				for out, in := range state.Inputs {
					input := avcontrol.Input{}

					if in.Audio != nil {
						c := *in.Audio
						input.Audio = &c
					}

					if in.Video != nil {
						c := *in.Video
						input.Video = &c
					}

					if in.AudioVideo != nil {
						c := *in.AudioVideo
						input.AudioVideo = &c
					}

					nState.Inputs[out] = input
				}
			}

			if state.Volumes != nil {
				nState.Volumes = make(map[string]int, len(state.Volumes))

				for block, vol := range state.Volumes {
					nState.Volumes[block] = vol
				}
			}

			if state.Mutes != nil {
				nState.Mutes = make(map[string]bool, len(state.Mutes))

				for block, m := range state.Mutes {
					nState.Mutes[block] = m
				}
			}

			n.APIRequest.Devices[id] = nState
		}
	}

	if cs.Requests != nil {
		n.Requests = make([]GenericControlRequest, len(cs.Requests))

		for i, req := range cs.Requests {
			u := *req.URL
			n.Requests[i].URL = &u
			n.Requests[i].Method = req.Method

			n.Requests[i].Body = make([]byte, len(req.Body))
			copy(n.Requests[i].Body, req.Body)
		}
	}

	return n
}
