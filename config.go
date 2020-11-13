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
	PowerOff ControlSet
	PowerOn  ControlSet

	Displays map[string]DisplayConfig
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
	Icon    string
	Sources map[string]SourceConfig
}

// SourceConfig represents a source and its associated controls
type SourceConfig struct {
	Icon    string
	Visible bool
	ControlSet

	// Sources represent sub-sources of the parent source
	Sources map[string]SourceConfig
}

// AudioConfig contains information about audio controls in the room
type AudioConfig struct {
	Media ControlSet

	// TODO do we need icons at any level for these things?
	// Groups is a map of group name -> audio name -> controlSet
	Groups map[string]map[string]ControlSet
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
