package bff

type Room struct {
	IDInfo
	Icon

	ControlGroups        []ControlGroup `json:"controlGroups"`
	SelectedControlGroup ID             `json:"selectedControlGroup"`

	// SharingGroups []SharingGroup `json:"sharingGroups"`
}

type ControlGroup struct {
	IDInfo
	Icon

	Displays      []Display      `json:"displays"`
	Inputs        []Input        `json:"inputs"`
	AudioGroups   []AudioGroup   `json:"audioGroups"`
	PresentGroups []PresentGroup `json:"presentGroups"`

	HelpRequested bool `json:"helpRequested"`
}

type Display struct {
	IDInfo
	Icon

	Input   ID   `json:"input"`
	Blanked bool `json:"blanked"`
	// allowedInputs ?
}

type Input struct {
	IDInfo
	Icon

	SubInputs []Input `json:"subInputs"`
}

type AudioGroup struct {
	IDInfo

	AudioDevices []AudioDevice `json:"audioDevices"`
}

type AudioDevice struct {
	IDInfo
	Icon // ?

	// should these be pointers?
	Level int  `json:"level"`
	Muted bool `json:"muted"`
}

type PresentGroup struct {
	IDInfo

	Items []PresentItem `json:"presentItems"`
}

type PresentItem struct {
	IDInfo
}

type Icon struct {
	Icon string `json:"icon"`
}

type IDInfo struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
}

type ID string
