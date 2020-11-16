package client

// Room .
type Room struct {
	Name string `json:"name"`

	ControlGroups        map[string]ControlGroup `json:"controlGroups"`
	SelectedControlGroup string                  `json:"selectedControlGroup"`
}

type DisplayGroups []DisplayGroup

// ControlGroup .
type ControlGroup struct {
	Name string `json:"name"`

	DisplayGroups DisplayGroups  `json:"displayGroups,omitempty"`
	AudioGroups   []AudioGroup   `json:"audioGroups,omitempty"`
	PresentGroups []PresentGroup `json:"presentGroups,omitempty"`

	PoweredOn  bool `json:"poweredOn"`
	MediaAudio struct {
		Level int  `json:"level"`
		Muted bool `json:"muted"`
	} `json:"mediaAudio"`

	ControlInfo struct {
		Key string `json:"key,omitempty"`
		URL string `json:"url,omitempty"`
	} `json:"controlInfo,omitempty"`

	Cameras []Camera `json:"cameras,omitempty"`

	Support Support `json:"support"`
}

type Camera struct {
	DisplayName string         `json:"displayName"`
	TiltUp      string         `json:"tiltUp"`
	TiltDown    string         `json:"tiltDown"`
	PanLeft     string         `json:"panLeft"`
	PanRight    string         `json:"panRight"`
	PanTiltStop string         `json:"panTiltStop"`
	ZoomIn      string         `json:"zoomIn"`
	ZoomOut     string         `json:"zoomOut"`
	ZoomStop    string         `json:"zoomStop"`
	Stream      string         `json:"stream"`
	Reboot      string         `json:"reboot,omitempty"`
	Presets     []CameraPreset `json:"presets"`
}

type CameraPreset struct {
	DisplayName string `json:"displayName"`
	SetPreset   string `json:"setPreset"`
	SavePreset  string `json:"savePreset,omitempty"`
}

// Support .
type Support struct {
	HelpRequested bool   `json:"helpRequested"`
	HelpMessage   string `json:"helpMessage"`
	HelpEnabled   bool   `json:"helpEnabled"`
}

// DisplayGroup .
type DisplayGroup struct {
	Name string `json:"name"`

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
	IconPair

	SubInputs []Input `json:"subInputs,omitempty"`
}

// AudioGroup .
type AudioGroup struct {
	Name string `json:"name"`

	AudioDevices []AudioDevice `json:"audioDevices"`
	Muted        bool          `json:"muted"`
}

// AudioDevice .
type AudioDevice struct {
	IconPair

	Level int  `json:"level"`
	Muted bool `json:"muted"`
}

// PresentGroup .
type PresentGroup struct {
	Name string `json:"name"`

	Items []PresentItem `json:"items"`
}

// PresentItem .
type PresentItem struct {
	Name string `json:"name"`
}

// IconPair .
type IconPair struct {
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
