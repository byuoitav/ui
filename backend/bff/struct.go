package bff

import "encoding/json"

type Room struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`

	ControlGroups        map[string]ControlGroup `json:"controlGroups"`
	SelectedControlGroup ID                      `json:"selectedControlGroup"`
}

type ControlGroup struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`

	Displays      []Display      `json:"displays"`
	Inputs        []Input        `json:"inputs"`
	AudioGroups   []AudioGroup   `json:"audioGroups"`
	PresentGroups []PresentGroup `json:"presentGroups"`

	Support Support `json:"support"`
}

type Support struct {
	HelpRequested bool `json:"helpRequested"`

	HelpMessage string `json:"helpMessage"`
	HelpEnabled bool   `json:"helpEnabled"`
}

type Display struct {
	ID ID `json:"id"`

	Outputs []IconPair `json:"outputs"`

	Input   ID   `json:"input"`
	Blanked bool `json:"blanked"`
}

type Input struct {
	ID ID `json:"id"`
	IconPair

	SubInputs []Input `json:"subInputs"`
	Disabled  bool    `json:"disabled"`
}

type AudioGroup struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`

	AudioDevices []AudioDevice `json:"audioDevices"`
	Muted        bool          `json:"muted"`
}

type AudioDevice struct {
	ID ID `json:"id"`
	IconPair

	Level int  `json:"level"`
	Muted bool `json:"muted"`
}

type PresentGroup struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`

	Items []PresentItem `json:"items"`
}

type PresentItem struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
}

type Icon struct {
	Icon string `json:"icon"`
}

type IconPair struct {
	ID   ID     `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Icon
}

type ID string

type HttpRequest struct {
	Method string          `json:"method"`
	URL    string          `json:"url"`
	Body   json.RawMessage `json:"body"`
}
