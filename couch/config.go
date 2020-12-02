package couch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/ui"
)

func (d *DataService) config(ctx context.Context, room string) (Config, error) {
	var config Config

	db := d.client.DB(ctx, d.database)
	if err := db.Get(ctx, room).ScanDoc(&config); err != nil {
		return Config{}, fmt.Errorf("unable to get/scan room: %w", err)
	}

	return config, nil
}

func (d *DataService) Config(ctx context.Context, room string) (ui.Config, error) {
	config, err := d.config(ctx, room)
	if err != nil {
		return ui.Config{}, err
	}

	return config.convert()
}

type Config struct {
	ID            string                  `json:"_id"`
	ControlGroups map[string]ControlGroup `json:"controlGroups"`
	States        map[string]State        `json:"states"`
}

type State avcontrol.StateRequest

func (c Config) convert() (ui.Config, error) {
	config := ui.Config{
		ID:            c.ID,
		ControlGroups: make(map[string]ui.ControlGroup, len(c.ControlGroups)),
		States:        make(map[string]ui.State, len(c.States)),
	}

	for k, state := range c.States {
		config.States[k] = ui.State(state)
	}

	for k, cg := range c.ControlGroups {
		group, err := cg.convert()
		if err != nil {
			return config, fmt.Errorf("unable to convert controlGroup %q: %w", k, err)
		}

		config.ControlGroups[k] = group
	}

	return config, nil
}

type ControlGroup struct {
	PowerOff StateControlConfig `json:"powerOff"`
	PowerOn  StateControlConfig `json:"powerOn"`

	Displays []Display      `json:"displays"`
	Audio    AudioConfig    `json:"audio"`
	Cameras  []CameraConfig `json:"cameras"`
}

func (cg ControlGroup) convert() (ui.ControlGroup, error) {
	var err error
	res := ui.ControlGroup{
		Displays: make([]ui.DisplayConfig, len(cg.Displays)),
		Audio: ui.AudioConfig{
			Groups: make([]ui.AudioGroupConfig, len(cg.Audio.Groups)),
		},
		Cameras: make([]ui.CameraConfig, len(cg.Cameras)),
	}

	res.PowerOff, err = cg.PowerOff.convert()
	if err != nil {
		return res, fmt.Errorf("unable to convert powerOff: %w", err)
	}

	res.PowerOn, err = cg.PowerOn.convert()
	if err != nil {
		return res, fmt.Errorf("unable to convert powerOn: %w", err)
	}

	for i := range cg.Displays {
		res.Displays[i], err = cg.Displays[i].convert()
		if err != nil {
			return res, fmt.Errorf("unable to convert display %d: %w", i, err)
		}
	}

	res.Audio.Media, err = cg.Audio.Media.convert()
	if err != nil {
		return res, fmt.Errorf("unable to convert audio.media: %w", err)
	}

	for i := range cg.Audio.Groups {
		res.Audio.Groups[i], err = cg.Audio.Groups[i].convert()
		if err != nil {
			return res, fmt.Errorf("unable to convert audio group %d: %w", i, err)
		}
	}

	for i := range cg.Cameras {
		res.Cameras[i], err = cg.Cameras[i].convert()
		if err != nil {
			return res, fmt.Errorf("unable to convert camera %d: %w", i, err)
		}
	}

	return res, nil
}

type Display struct {
	Name string `json:"name"`
	Icon string `json:"icon"`

	Sources []Source `json:"sources"`
}

func (d Display) convert() (ui.DisplayConfig, error) {
	var err error
	res := ui.DisplayConfig{
		Name:    d.Name,
		Icon:    d.Icon,
		Sources: make([]ui.SourceConfig, len(d.Sources)),
	}

	for i := range d.Sources {
		res.Sources[i], err = d.Sources[i].convert()
		if err != nil {
			return res, fmt.Errorf("unable to convert source %d: %w", i, err)
		}
	}

	return res, err
}

type Source struct {
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Visible bool   `json:"visible"`
	StateControlConfig

	Sources []Source `json:"sources"`
}

func (s Source) convert() (ui.SourceConfig, error) {
	var err error
	res := ui.SourceConfig{
		Name:    s.Name,
		Icon:    s.Icon,
		Visible: s.Visible,
		Sources: make([]ui.SourceConfig, len(s.Sources)),
	}

	res.StateControlConfig, err = s.StateControlConfig.convert()
	if err != nil {
		return res, err
	}

	for i := range s.Sources {
		res.Sources[i], err = s.Sources[i].convert()
		if err != nil {
			return res, fmt.Errorf("unable to convert subsource %d: %w", i, err)
		}
	}

	return res, nil
}

type AudioConfig struct {
	Media  AudioDevice  `json:"media"`
	Groups []AudioGroup `json:"groups"`
}

type AudioGroup struct {
	Name         string        `json:"name"`
	AudioDevices []AudioDevice `json:"audioDevices"`
}

func (a AudioGroup) convert() (ui.AudioGroupConfig, error) {
	var err error
	res := ui.AudioGroupConfig{
		Name:         a.Name,
		AudioDevices: make([]ui.AudioDeviceConfig, len(a.AudioDevices)),
	}

	for i := range a.AudioDevices {
		res.AudioDevices[i], err = a.AudioDevices[i].convert()
		if err != nil {
			return res, fmt.Errorf("unable to convert audio device %d: %w", i, err)
		}
	}

	return res, nil
}

type AudioDevice struct {
	Name   string             `json:"name"`
	Volume StateControlConfig `json:"volume"`
	Mute   StateControlConfig `json:"mute"`
	Unmute StateControlConfig `json:"unmute"`
}

func (a AudioDevice) convert() (ui.AudioDeviceConfig, error) {
	var err error
	res := ui.AudioDeviceConfig{
		Name: a.Name,
	}

	res.Volume, err = a.Volume.convert()
	if err != nil {
		return res, fmt.Errorf("unable to convert volume: %w", err)
	}

	res.Mute, err = a.Mute.convert()
	if err != nil {
		return res, fmt.Errorf("unable to convert mute: %w", err)
	}

	res.Unmute, err = a.Unmute.convert()
	if err != nil {
		return res, fmt.Errorf("unable to convert unmute: %w", err)
	}

	return res, nil
}

type CameraConfig struct {
	Name        string             `json:"name"`
	TiltUp      StateControlConfig `json:"tiltUp"`
	TiltDown    StateControlConfig `json:"tiltDown"`
	PanLeft     StateControlConfig `json:"panLeft"`
	PanRight    StateControlConfig `json:"panRight"`
	PanTiltStop StateControlConfig `json:"panTiltStop"`
	ZoomIn      StateControlConfig `json:"zoomIn"`
	ZoomOut     StateControlConfig `json:"zoomOut"`
	ZoomStop    StateControlConfig `json:"zoomStop"`
	// Reboot      string         `json:"reboot"` // this service doesn't care about this, savePreset, and stream
	// Stream      controlSet `json:"stream"` // left in for shipyard memory

	Presets []CameraPresetConfig `json:"presets"`
}

func (c CameraConfig) convert() (ui.CameraConfig, error) {
	var err error
	res := ui.CameraConfig{
		Name:    c.Name,
		Presets: make([]ui.CameraPresetConfig, len(c.Presets)),
	}

	res.TiltUp, err = c.TiltUp.convert()
	if err != nil {
		return res, fmt.Errorf("unable to convert tiltUp: %w", err)
	}

	res.TiltDown, err = c.TiltDown.convert()
	if err != nil {
		return res, fmt.Errorf("unable to convert tiltDown: %w", err)
	}

	res.PanLeft, err = c.PanLeft.convert()
	if err != nil {
		return res, fmt.Errorf("unable to convert panLeft: %w", err)
	}

	res.PanRight, err = c.PanRight.convert()
	if err != nil {
		return res, fmt.Errorf("unable to convert panRight: %w", err)
	}

	res.PanTiltStop, err = c.PanTiltStop.convert()
	if err != nil {
		return res, fmt.Errorf("unable to convert panTiltStop: %w", err)
	}

	res.ZoomIn, err = c.ZoomIn.convert()
	if err != nil {
		return res, fmt.Errorf("unable to convert zoomIn: %w", err)
	}

	res.ZoomOut, err = c.ZoomOut.convert()
	if err != nil {
		return res, fmt.Errorf("unable to convert zoomOut: %w", err)
	}

	res.ZoomStop, err = c.ZoomStop.convert()
	if err != nil {
		return res, fmt.Errorf("unable to convert zoomStop: %w", err)
	}

	for i := range c.Presets {
		res.Presets[i], err = c.Presets[i].convert()
		if err != nil {
			return res, fmt.Errorf("unable to convert preset %d: %w", i, err)
		}
	}

	return res, nil
}

type CameraPresetConfig struct {
	Name      string             `json:"name"`
	SetPreset StateControlConfig `json:"setPreset"`
	// SavePreset string `json:"savePreset"`
}

func (c CameraPresetConfig) convert() (ui.CameraPresetConfig, error) {
	var err error
	res := ui.CameraPresetConfig{
		Name: c.Name,
	}

	res.SetPreset, err = c.SetPreset.convert()
	if err != nil {
		return res, fmt.Errorf("unable to convert setPreset: %w", err)
	}

	return res, nil
}

type StateControlConfig struct {
	MatchStates      []string          `json:"matchStates"`
	StateTransitions []StateTransition `json:"stateTransitions"`
}

func (c StateControlConfig) convert() (ui.StateControlConfig, error) {
	var err error
	res := ui.StateControlConfig{
		StateTransitions: make([]ui.StateTransition, len(c.StateTransitions)),
		MatchStates:      c.MatchStates,
	}

	for i := range c.StateTransitions {
		res.StateTransitions[i].MatchStates = c.StateTransitions[i].MatchStates

		res.StateTransitions[i].Action, err = c.StateTransitions[i].Action.convert()
		if err != nil {
			return res, fmt.Errorf("unable to convert action: %w", err)
		}
	}

	return res, nil
}

type StateTransition struct {
	MatchStates []string              `json:"matchStates"`
	Action      StateTransitionAction `json:"action"`
}

type StateTransitionAction struct {
	SetStates []string         `json:"setStates"`
	Requests  []GenericRequest `json:"requests"`
}

func (a StateTransitionAction) convert() (ui.StateTransitionAction, error) {
	res := ui.StateTransitionAction{
		Requests:  make([]ui.GenericRequest, len(a.Requests)),
		SetStates: a.SetStates,
	}

	for i := range a.Requests {
		res.Requests[i].Body = a.Requests[i].Body
		res.Requests[i].Method = a.Requests[i].Method

		if res.Requests[i].Method == "" {
			res.Requests[i].Method = http.MethodGet
		}

		u, err := url.Parse(a.Requests[i].URL)
		if err != nil {
			return res, err
		}

		res.Requests[i].URL = u
	}

	return res, nil
}

type GenericRequest struct {
	URL    string          `json:"url"`
	Method string          `json:"method"`
	Body   json.RawMessage `json:"body"`
}
