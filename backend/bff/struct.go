package bff

import "strings"

// ShareState is one of 7 possible share states

// Room .
type Room struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`

	ControlGroups        map[string]ControlGroup `json:"controlGroups"`
	SelectedControlGroup ID                      `json:"selectedControlGroup"`
}

type DisplayGroups []DisplayGroup

// ControlGroup .
type ControlGroup struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`

	PoweredOn bool `json:"poweredOn"`

	fullDisplayGroups DisplayGroups
	DisplayGroups     DisplayGroups  `json:"displayGroups,omitempty"`
	Inputs            []Input        `json:"inputs"`
	AudioGroups       []AudioGroup   `json:"audioGroups,omitempty"`
	PresentGroups     []PresentGroup `json:"presentGroups,omitempty"`

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
	State   shareState `json:"state"`
	Options []string   `json:"opts,omitempty"`
	Master  ID         `json:"master,omitempty"`
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

func IDsToStrings(ids []ID) []string {
	var strs []string

	for i := range ids {
		strs = append(strs, string(ids[i]))
	}

	return strs
}

func StringsToIDs(strings []string) []ID {
	var ids []ID

	for i := range strings {
		ids = append(ids, ID(strings[i]))
	}

	return ids
}

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

func IntP(i int) *int {
	return &i
}
