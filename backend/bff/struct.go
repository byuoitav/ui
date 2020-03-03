package bff

import "strings"

// ShareState is one of 7 possible share states
type ShareState int

const (
	// CantShare means that you can't share at all
	CantShare ShareState = iota

	// Share means that you can share right now
	Share

	// Unshare means that you are currently sharing, and that you could unshare
	Unshare

	// MinionActive means that you are being shared to and are participating in that share
	MinionActive

	// MinionInactive means that you are being shared to but you are NOT participating in that share
	MinionInactive
)

// Room .
type Room struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`

	ControlGroups        map[string]ControlGroup `json:"controlGroups"`
	SelectedControlGroup ID                      `json:"selectedControlGroup"`
}

// ControlGroup .
type ControlGroup struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`

	PoweredOn bool `json:"poweredOn"`

	DisplayGroups []DisplayGroup `json:"displayGroups,omitempty"`
	Inputs        []Input        `json:"inputs"`
	AudioGroups   []AudioGroup   `json:"audioGroups,omitempty"`
	PresentGroups []PresentGroup `json:"presentGroups,omitempty"`

	MediaAudio struct {
		Level int  `json:"level"`
		Muted bool `json:"muted"`
	} `json:"mediaAudio"`

	ControlInfo struct {
		Key string `json:"key,omitempty"`
		URL string `json:"url,omitempty"`
	} `json:"controlInfo,omitempty"`

	Support Support `json:"support"`
}

// Support .
type Support struct {
	HelpRequested bool   `json:"helpRequested"`
	HelpMessage   string `json:"helpMessage"`
	HelpEnabled   bool   `json:"helpEnabled"`
}

// DisplayGroup .
type DisplayGroup struct {
	ID ID `json:"id"`

	Displays []IconPair `json:"displays"`
	Blanked  bool       `json:"blanked"`
	Input    ID         `json:"input"`

	ShareInfo ShareInfo `json:"shareInfo,omitempty"`
}

// ShareInfo .
type ShareInfo struct {
	State   ShareState `json:"state"`
	Options []string   `json:"options,omitempty"`
}

// Input .
type Input struct {
	ID ID `json:"id"`
	IconPair

	SubInputs []Input `json:"subInputs,omitempty"`
}

// AudioGroup .
type AudioGroup struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`

	AudioDevices []AudioDevice `json:"audioDevices"`
	Muted        bool          `json:"muted"`
}

// AudioDevice .
type AudioDevice struct {
	ID ID `json:"id"`
	IconPair

	Level int  `json:"level"`
	Muted bool `json:"muted"`
}

// PresentGroup .
type PresentGroup struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`

	Items []PresentItem `json:"items"`
}

// PresentItem .
type PresentItem struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
}

// IconPair .
type IconPair struct {
	ID   ID     `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Icon string `json:"icon"`
}

// ID .
type ID string

// GetName gets the name of an ID
func (i ID) GetName() string {
	split := strings.Split(string(i), "-")
	if len(split) != 3 {
		return string(i)
	}

	return split[2]
}

// BoolP .
func BoolP(b bool) *bool {
	return &b
}
