package ui

import avcontrol "github.com/byuoitav/av-control-api"

// Config represents the program for a room that can be used
// to control a room
type Config struct {
	ID            string
	ControlPanels map[string]ControlPanelConfig
	ControlGroups map[string]ControlGroup
}

type ControlPanelConfig struct {
	UIType       string `json:"uiType"`
	ControlGroup string `json:"controlGroup"`
}

// ControlGroup represents a group of Devices and inputs. These groups
// are used for logical grouping and displaying on different UIs
type ControlGroup struct {
	// PowerOff represents the state the room needs to be in to be considered "off"
	PowerOff ControlSet

	// PowerOn is the state we set to turn on the room
	PowerOn ControlSet

	Displays []DisplayConfig
	Audio    AudioConfig
	Cameras  []CameraConfig
}

// CameraConfig represents a Camera and its associated control endpoints
type CameraConfig struct {
	DisplayName string
	TiltUp      string
	TiltDown    string
	PanLeft     string
	PanRight    string
	PanTiltStop string
	ZoomIn      string
	ZoomOut     string
	ZoomStop    string
	Stream      string
	Reboot      string
	Presets     []CameraPresetConfig
}

// CameraPresetConfig represents a preset on a camera
type CameraPresetConfig struct {
	DisplayName string
	SetPreset   string
	SavePreset  string
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
	ControlSet

	// Sources represent sub-sources of the parent source
	Sources []SourceConfig
}

// AudioConfig contains information about audio controls in the room
type AudioConfig struct {
	Media ControlSet

	Groups []AudioGroupConfig
}

type AudioGroupConfig struct {
	Name         string
	AudioDevices []AudioDeviceConfig
}

type AudioDeviceConfig struct {
	Name string
	ControlSet
}

// ControlSet represents the request to be made (both to the
// AV Control API and other arbitrary locations) in order to set the room
// to a given state
type ControlSet struct {
	APIRequest avcontrol.StateRequest
	Requests   []GenericControlRequest
}

// GenericControlRequest contains the information necessary to make a generic
// HTTP request in association with making state changes in a room
type GenericControlRequest struct {
	URL    string
	Method string
	Body   []byte
}
