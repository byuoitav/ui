package ui

import (
	"net/url"

	avcontrol "github.com/byuoitav/av-control-api"
)

type State avcontrol.StateRequest

// Config represents the program for a room that can be used
// to control a room
type Config struct {
	ID            string
	ControlGroups map[string]ControlGroup
	States        map[string]State
}

// ControlGroup represents a group of Devices and inputs. These groups
// are used for logical grouping and displaying on different UIs
type ControlGroup struct {
	// PowerOff represents the state the room needs to be in to be considered "off"
	PowerOff StateControlConfig

	// PowerOn is the state we set to turn on the room
	PowerOn StateControlConfig

	Displays []DisplayConfig
	Audio    AudioConfig
	Cameras  []CameraConfig
}

// CameraConfig represents a Camera and its associated control endpoints
type CameraConfig struct {
	Name        string
	TiltUp      StateControlConfig
	TiltDown    StateControlConfig
	PanLeft     StateControlConfig
	PanRight    StateControlConfig
	PanTiltStop StateControlConfig
	ZoomIn      StateControlConfig
	ZoomOut     StateControlConfig
	ZoomStop    StateControlConfig
	Presets     []CameraPresetConfig
}

// CameraPresetConfig represents a preset on a camera
type CameraPresetConfig struct {
	Name      string
	SetPreset StateControlConfig
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
	StateControlConfig

	// Sources represent sub-sources of the parent source
	Sources []SourceConfig
}

// AudioConfig contains information about audio controls in the room
type AudioConfig struct {
	Media  AudioDeviceConfig
	Groups []AudioGroupConfig
}

type AudioGroupConfig struct {
	Name         string
	AudioDevices []AudioDeviceConfig
}

type AudioDeviceConfig struct {
	Name   string
	Volume StateControlConfig `json:"volume"`
	Mute   StateControlConfig `json:"mute"`
	Unmute StateControlConfig `json:"unmute"`
}

type StateControlConfig struct {
	MatchStates      []string
	StateTransitions []StateTransition
}

type StateTransition struct {
	MatchStates []string
	Action      StateTransitionAction
}

type StateTransitionAction struct {
	SetStates []string
	Requests  []GenericRequest
}

// GenericRequest contains the information necessary to make a generic
// HTTP request in association with making state changes in a room
type GenericRequest struct {
	URL    *url.URL
	Method string
	Body   []byte
}

func (s State) Copy() State {
	if s.Devices == nil {
		return State{}
	}

	res := State{
		Devices: make(map[avcontrol.DeviceID]avcontrol.DeviceState, len(s.Devices)),
	}

	for id, dev := range s.Devices {
		newDev := avcontrol.DeviceState{}

		if dev.PoweredOn != nil {
			t := *dev.PoweredOn
			newDev.PoweredOn = &t
		}

		if dev.Blanked != nil {
			t := *dev.Blanked
			newDev.Blanked = &t
		}

		if dev.Inputs != nil {
			newDev.Inputs = make(map[string]avcontrol.Input, len(dev.Inputs))

			for out, in := range dev.Inputs {
				newInput := avcontrol.Input{}

				if in.Audio != nil {
					t := *in.Audio
					newInput.Audio = &t
				}

				if in.Video != nil {
					t := *in.Video
					newInput.Video = &t
				}

				if in.AudioVideo != nil {
					t := *in.AudioVideo
					newInput.AudioVideo = &t
				}

				newDev.Inputs[out] = newInput
			}
		}

		if dev.Volumes != nil {
			newDev.Volumes = make(map[string]int, len(dev.Volumes))

			for block, vol := range dev.Volumes {
				newDev.Volumes[block] = vol
			}
		}

		if dev.Mutes != nil {
			newDev.Mutes = make(map[string]bool, len(dev.Mutes))

			for block, muted := range dev.Mutes {
				newDev.Mutes[block] = muted
			}
		}

		res.Devices[id] = newDev
	}

	return res
}
