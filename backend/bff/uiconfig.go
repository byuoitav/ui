package bff

import (
	"encoding/json"
)

// UIConfig - a representation of a ui config
type UIConfig struct {
	ID                  string               `json:"_id,omitempty"`
	Panels              []Panel              `json:"panels"`
	Presets             []Preset             `json:"presets"`
	InputConfiguration  []IOConfiguration    `json:"inputConfiguration"`
	OutputConfiguration []IOConfiguration    `json:"outputConfiguration"`
	AudioConfiguration  []AudioConfiguration `json:"audioConfiguration"`
	PseudoInputs        []PseudoInput        `json:"pseudoInputs,omitempty"`

	ActionsConfig json.RawMessage `json:"actions,omitempty"`
}

// Preset - a representation of what is controlled by this preset.
type Preset struct {
	Name                    string              `json:"name"`
	Icon                    string              `json:"icon"`
	Displays                []string            `json:"displays"`
	ShareableDisplays       []string            `json:"shareableDisplays"`
	AudioDevices            []string            `json:"audioDevices"`
	Inputs                  []string            `json:"inputs"`
	IndependentAudioDevices []string            `json:"independentAudioDevices,omitempty"`
	AudioGroups             map[string][]string `json:"audioGroups,omitempty"`
	VolumeMatches           []string            `json:"volumeMatches,omitempty"`

	ActionsConfig json.RawMessage `json:"actions,omitempty"`
	Actions       struct {
		SetInput    SetInput    `json:"setInput,omitempty"`
		SetVolume   SetVolume   `json:"setVolume,omitempty"`
		SetMuted    SetMuted    `json:"setMuted,omitempty"`
		SetPower    SetPower    `json:"setPower,omitempty"`
		HelpRequest HelpRequest `json:"helpRequest,omitempty"`
		SetSharing  SetSharing  `json:"setSharing,omitempty"`
	} `json:"-"`
}

// Panel - a representation of a touchpanel and which preset it has.
type Panel struct {
	Hostname string   `json:"hostname"`
	UIPath   string   `json:"uipath"`
	Preset   string   `json:"preset"`
	Features []string `json:"features"`
}

// AudioConfiguration - a representation of how the audio is configured when using multiple displays.
type AudioConfiguration struct {
	Display      string   `json:"display"`
	AudioDevices []string `json:"audioDevices"`
	RoomWide     bool     `json:"roomWide"`
}

// IOConfiguration - a representation of an input or output device.
type IOConfiguration struct {
	Name        string            `json:"name"`
	Icon        string            `json:"icon"`
	Displayname *string           `json:"displayname,omitempty"`
	SubInputs   []IOConfiguration `json:"subInputs,omitempty"`
}

// PseudoInput - a fake input I guess
type PseudoInput struct {
	Displayname string `json:"displayname"`
	Config      []struct {
		Input   string   `json:"input"`
		Outputs []string `json:"outputs"`
	} `json:"config"`
}

// UnmarshalJSON .
func (u *UIConfig) UnmarshalJSON(b []byte) error {
	type Alias UIConfig
	config := &struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(b, &config); err != nil {
		return err
	}

	for i := range u.Presets {
		// first unmarshal the room wide action config into this preset's actions
		if len(u.ActionsConfig) > 0 {
			if err := json.Unmarshal(u.ActionsConfig, &u.Presets[i].Actions); err != nil {
				return err
			}
		}

		// then do the preset specific config
		if len(u.Presets[i].ActionsConfig) > 0 {
			if err := json.Unmarshal(u.Presets[i].ActionsConfig, &u.Presets[i].Actions); err != nil {
				return err
			}
		}
	}

	return nil
}
