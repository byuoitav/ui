package couch

import (
	"context"
	"encoding/json"
	"fmt"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/ui"
)

type config struct {
	ID string `json:"_id"`

	ControlPanels map[string]struct {
		UIType string `json:"uiType"`

		// TODO divider sensors
		ControlGroup string `json:"controlGroup"`
	} `json:"controlPanels"`

	ControlGroups map[string]struct {
		PowerOff controlSet `json:"powerOff"`
		PowerOn  controlSet `json:"powerOn"`

		Displays []struct {
			Name string `json:"name"`
			Icon string `json:"icon"`

			Sources []struct {
				Name    string `json:"name"`
				Icon    string `json:"icon"`
				Visible bool   `json:"visible"`
				controlSet

				Sources []struct {
					Name    string `json:"name"`
					Icon    string `json:"icon"`
					Visible bool   `json:"visible"`
					controlSet
				} `json:"sources"`
			} `json:"sources"`
		} `json:"displays"`

		Audio struct {
			Media struct {
				Volume controlSet `json:"volume"`
				Mute   controlSet `json:"mute"`
				Unmute controlSet `json:"unmute"`
			} `json:"media"`
			Groups []struct {
				Name         string `json:"name"`
				AudioDevices []struct {
					Name   string     `json:"name"`
					Volume controlSet `json:"volume"`
					Mute   controlSet `json:"mute"`
					Unmute controlSet `json:"unmute"`
				} `json:"audioDevices"`
			} `json:"groups"`
		} `json:"audio"`

		Cameras struct {
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
			Reboot      string         `json:"reboot"`
			Presets     []CameraPreset `json:"presets"`
		} `json:"cameras"`
	} `json:"controlGroups"`
}

type CameraPreset struct {
	DisplayName string `json:"displayName"`
	SetPreset   string `json:"setPreset"`
	SavePreset  string `json:"savePreset"`
}

type controlSet struct {
	APIRequest avcontrol.StateRequest
	Requests   []struct {
		URL    string          `json:"url"`
		Method string          `json:"method"`
		Body   json.RawMessage `json:"body"`
	} `json:"requests"`
}

func (d *DataService) config(ctx context.Context, room string) (config, error) {
	var config config

	db := d.client.DB(ctx, d.database)
	if err := db.Get(ctx, room).ScanDoc(&config); err != nil {
		return config, fmt.Errorf("unable to get/scan room: %w", err)
	}

	return config, nil
}

func (d *DataService) Config(ctx context.Context, room string) (ui.Config, error) {
	config, err := d.config(ctx, room)
	if err != nil {
		return ui.Config{}, err
	}

	return config.convert(), nil
}

// TODO camera stuff
func (c config) convert() ui.Config {
	config := ui.Config{
		ID:            c.ID,
		ControlPanels: make(map[string]ui.ControlPanelConfig),
		ControlGroups: make(map[string]ui.ControlGroup, len(c.ControlGroups)),
	}

	for k, v := range c.ControlPanels {
		config.ControlPanels[k] = ui.ControlPanelConfig{
			UIType:       v.UIType,
			ControlGroup: v.ControlGroup,
		}
	}

	for k, v := range c.ControlGroups {
		controlGroup := ui.ControlGroup{
			PowerOff: v.PowerOff.convert(),
			PowerOn:  v.PowerOn.convert(),
		}

		controlGroup.Audio.Media.Volume = v.Audio.Media.Volume.convert()
		controlGroup.Audio.Media.Mute = v.Audio.Media.Mute.convert()
		controlGroup.Audio.Media.Unmute = v.Audio.Media.Unmute.convert()

		for _, disp := range v.Displays {
			uiDisp := ui.DisplayConfig{
				Name: disp.Name,
				Icon: disp.Icon,
			}

			for _, source := range disp.Sources {
				sourceConfig := ui.SourceConfig{
					Name:       source.Name,
					Icon:       source.Icon,
					Visible:    source.Visible,
					ControlSet: source.controlSet.convert(),
				}

				for _, subSource := range source.Sources {
					sourceConfig.Sources = append(sourceConfig.Sources, ui.SourceConfig{
						Name:       subSource.Name,
						Icon:       subSource.Icon,
						Visible:    subSource.Visible,
						ControlSet: subSource.controlSet.convert(),
					})
				}

				uiDisp.Sources = append(uiDisp.Sources, sourceConfig)
			}

			controlGroup.Displays = append(controlGroup.Displays, uiDisp)
		}

		for _, ag := range v.Audio.Groups {
			audioGroup := ui.AudioGroupConfig{
				Name: ag.Name,
			}

			for _, ad := range ag.AudioDevices {
				audioGroup.AudioDevices = append(audioGroup.AudioDevices, ui.AudioDeviceConfig{
					Name:   ad.Name,
					Volume: ad.Volume.convert(),
					Mute:   ad.Mute.convert(),
					Unmute: ad.Unmute.convert(),
				})
			}

			controlGroup.Audio.Groups = append(controlGroup.Audio.Groups, audioGroup)
		}

		config.ControlGroups[k] = controlGroup
	}

	return config
}

func (c controlSet) convert() ui.ControlSet {
	cs := ui.ControlSet{
		APIRequest: c.APIRequest,
		Requests:   make([]ui.GenericControlRequest, len(c.Requests)),
	}

	for i := range c.Requests {
		cs.Requests[i].Body = c.Requests[i].Body
		cs.Requests[i].Method = c.Requests[i].Method
		cs.Requests[i].URL = c.Requests[i].URL
	}

	return cs
}
