package ui

import avcontrol "github.com/byuoitav/av-control-api"

// Config represents the program for a room that can be used
// to control a room
type Config struct {
	ID            string
	ControlPanels map[string]string
	ControlGroups map[string]ControlGroup
}

// ControlGroup represents a group of Devices and inputs. These groups
// are used for logical grouping and displaying on different UIs
type ControlGroup struct {
	PowerOff ControlSet
	PowerOn  ControlSet

	Displays map[string]DisplayConfig
	Audio    AudioConfig
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

	// map of group name -> audio name -> controlSet
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
