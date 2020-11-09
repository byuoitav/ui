package client

// Room .
type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"` // TODO remove? what does the ui use this for

	ControlGroups        map[string]ControlGroup `json:"controlGroups"`
	SelectedControlGroup string                  `json:"selectedControlGroup"`
}

type DisplayGroups []DisplayGroup

// ControlGroup .
type ControlGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"` // TODO do we need both?

	PoweredOn bool `json:"poweredOn"`

	// fullDisplayGroups DisplayGroups
	DisplayGroups DisplayGroups  `json:"displayGroups,omitempty"`
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
	ID string `json:"id"`

	Displays []IconPair `json:"displays"`
	Inputs   []Input    `json:"inputs"`
	Input    string     `json:"input"`

	ShareInfo ShareInfo `json:"shareInfo,omitempty"`
}

// ShareInfo .
type ShareInfo struct {
	// State   shareState `json:"state"`
	Options []string `json:"opts,omitempty"`
	Master  string   `json:"master,omitempty"`
}

// Input .
type Input struct {
	ID string `json:"id"`
	IconPair

	SubInputs []Input `json:"subInputs,omitempty"`
}

// AudioGroup .
type AudioGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	AudioDevices []AudioDevice `json:"audioDevices"`
	Muted        bool          `json:"muted"`
}

// AudioDevice .
type AudioDevice struct {
	ID string `json:"id"`
	IconPair

	Level int  `json:"level"`
	Muted bool `json:"muted"`
}

// PresentGroup .
type PresentGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Items []PresentItem `json:"items"`
}

// PresentItem .
type PresentItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// IconPair .
type IconPair struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Icon string `json:"icon"`
}

// BoolP .
func BoolP(b bool) *bool {
	return &b
}

func IntP(i int) *int {
	return &i
}
