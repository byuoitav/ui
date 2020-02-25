package bff

// ShareState is one of 7 possible share states
type ShareState int

const (
	// Nothing means that you can't share at all
	Nothing ShareState = iota + 1

	// Share means that you can share right now
	Share

	// Unshare means that you are currently sharing, and that you could unshare
	Unshare

	// Link means that you can link
	Link

	// Unlink means that you currently are linked, and you could unlink
	Unlink

	// MinionActive means that you are being shared to and are participating in that share
	MinionActive

	// MinionInactive means that you are being shared to but you are NOT participating in that share
	MinionInactive
)

//// Shareable .
//type Shareable map[ID][]ID
//
//// Sharing .
//type Sharing map[ID]ShareGroups

// ShareGroups .
//type ShareGroups struct {
//	Input    ID   `json:"input"`
//	Active   []ID `json:"active"`
//	Inactive []ID `json:"inactive"`
//	Linked   []ID `json:"linked"`
//}

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

	//TODO switch power to be a boolean?
	Power string `json:"power"`

	MediaAudio struct {
		Level int  `json:"level"`
		Muted bool `json:"muted"`
	} `json:"mediaAudio"`

	DisplayBlocks []DisplayBlock `json:"displayBlocks"`
	Inputs        []Input        `json:"inputs"`
	AudioGroups   []AudioGroup   `json:"audioGroups"`
	PresentGroups []PresentGroup `json:"presentGroups"`

	Support Support `json:"support"`
}

// Support .
type Support struct {
	HelpRequested bool `json:"helpRequested"`

	HelpMessage string `json:"helpMessage"`
	HelpEnabled bool   `json:"helpEnabled"`
}

// DisplayBlock .
type DisplayBlock struct {
	ID ID `json:"id"`

	Blanked bool `json:"blanked"`

	Outputs []IconPair `json:"outputs"`
	Input   ID         `json:"input"`
	Share   ShareInfo  `json:"share"`
}

// ShareInfo .
type ShareInfo struct {
	State   ShareState `json:"state"`
	Master  ID         `json:"master"`
	Options []string   `json:"options"`
}

// Input .
type Input struct {
	ID ID `json:"id"`
	IconPair

	SubInputs []Input `json:"subInputs"`
	Disabled  bool    `json:"disabled"`
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

// Icon .
type Icon struct {
	Icon string `json:"icon"`
}

// IconPair .
type IconPair struct {
	ID   ID     `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Icon
}

// ID .
type ID string

// BoolP .
func BoolP(b bool) *bool {
	return &b
}
