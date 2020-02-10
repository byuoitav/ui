package bff

import (
	"encoding/json"

	"github.com/byuoitav/lazarette/lazarette"
)

// ShareState is one of 6 possible share states
type ShareState int

/*
	Nothing		   - cannot share
	Share		   - can share
	Unshare        - can stop sharing (is currently sharing)
	Link           - can link
	Unlink         - can unlink (is currently linked)
	MinionActive   - is shared to and is displaying the share
	MinionInactive - is shared to and is not displaying the share
*/
const (
	Nothing = iota
	Share
	Unshare
	Link
	Unlink
	MinionActive
	MinionInactive
)

// LazState .
type LazState struct {
	Client       lazarette.LazaretteClient
	Subscription lazarette.Lazarette_SubscribeClient
}

// Shareable .
type Shareable map[ID][]ID

// Sharing .
type Sharing map[ID]ShareGroups

// ShareGroups .
type ShareGroups struct {
	Input    ID   `json:"input"`
	Active   []ID `json:"active"`
	Inactive []ID `json:"inactive"`
	Linked   []ID `json:"linked"`
}

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
	//TODO am right?
	Power string `json:"power"`

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

	Outputs []IconPair `json:"outputs"`
	Input   ID         `json:"input"`
	Share   ShareInfo  `json:"share"`
}

// ShareInfo .
type ShareInfo struct {
	Options []string   `json:"shareOptions"`
	State   ShareState `json:"shareState"`
	Master  ID         `json:"shareMaster"`
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

// HTTPRequest .
type HTTPRequest struct {
	Method string          `json:"method"`
	URL    string          `json:"url"`
	Body   json.RawMessage `json:"body"`
}

// BoolP .
func BoolP(b bool) *bool {
	return &b
}
