package couch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/ui"
)

type config struct {
	ID string `json:"_id"`

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

		Cameras []struct {
			Name        string     `json:"name"`
			TiltUp      controlSet `json:"tiltUp"`
			TiltDown    controlSet `json:"tiltDown"`
			PanLeft     controlSet `json:"panLeft"`
			PanRight    controlSet `json:"panRight"`
			PanTiltStop controlSet `json:"panTiltStop"`
			ZoomIn      controlSet `json:"zoomIn"`
			ZoomOut     controlSet `json:"zoomOut"`
			ZoomStop    controlSet `json:"zoomStop"`
			// Stream      controlSet `json:"stream"`
			// Reboot      string         `json:"reboot"` // this service doesn't care about this, savePreset, and stream
			Presets []struct {
				Name      string     `json:"name"`
				SetPreset controlSet `json:"setPreset"`
				// SavePreset string `json:"savePreset"`
			} `json:"presets"`
		} `json:"cameras"`
	} `json:"controlGroups"`
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

	return config.convert()
}

func (c config) convert() (ui.Config, error) {
	var err error
	config := ui.Config{
		ID:            c.ID,
		ControlGroups: make(map[string]ui.ControlGroup, len(c.ControlGroups)),
	}

	for k, v := range c.ControlGroups {
		var controlGroup ui.ControlGroup

		controlGroup.PowerOff, err = v.PowerOff.convert()
		if err != nil {
			return config, fmt.Errorf("invalid control set 'powerOff': %w", err)
		}

		controlGroup.PowerOn, err = v.PowerOff.convert()
		if err != nil {
			return config, fmt.Errorf("invalid control set 'powerOn': %w", err)
		}

		controlGroup.Audio.Media.Volume, err = v.Audio.Media.Volume.convert()
		if err != nil {
			return config, fmt.Errorf("invalid control set 'volume': %w", err)
		}

		controlGroup.Audio.Media.Mute, err = v.Audio.Media.Mute.convert()
		if err != nil {
			return config, fmt.Errorf("invalid control set 'mute': %w", err)
		}

		controlGroup.Audio.Media.Unmute, err = v.Audio.Media.Unmute.convert()
		if err != nil {
			return config, fmt.Errorf("invalid control set 'unmute': %w", err)
		}

		for _, disp := range v.Displays {
			uiDisp := ui.DisplayConfig{
				Name: disp.Name,
				Icon: disp.Icon,
			}

			for _, source := range disp.Sources {
				sourceConfig := ui.SourceConfig{
					Name:    source.Name,
					Icon:    source.Icon,
					Visible: source.Visible,
				}

				sourceConfig.ControlSet, err = source.controlSet.convert()
				if err != nil {
					return config, fmt.Errorf("invalid control set '%s.%s': %w", disp.Name, source.Name, err)
				}

				for _, subSource := range source.Sources {
					subSourceConfig := ui.SourceConfig{
						Name:    subSource.Name,
						Icon:    subSource.Icon,
						Visible: subSource.Visible,
					}

					subSourceConfig.ControlSet, err = subSource.controlSet.convert()
					if err != nil {
						return config, fmt.Errorf("invalid control set '%s.%s.%s': %w", disp.Name, source.Name, subSource.Name, err)
					}

					sourceConfig.Sources = append(sourceConfig.Sources, subSourceConfig)
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
				audioDeviceConfig := ui.AudioDeviceConfig{
					Name: ad.Name,
				}

				audioDeviceConfig.Volume, err = ad.Volume.convert()
				if err != nil {
					return config, fmt.Errorf("invalid control set '%s.%s.volume': %w", ag.Name, ad.Name, err)
				}

				audioDeviceConfig.Mute, err = ad.Mute.convert()
				if err != nil {
					return config, fmt.Errorf("invalid control set '%s.%s.mute': %w", ag.Name, ad.Name, err)
				}

				audioDeviceConfig.Unmute, err = ad.Unmute.convert()
				if err != nil {
					return config, fmt.Errorf("invalid control set '%s.%s.unmute': %w", ag.Name, ad.Name, err)
				}

				audioGroup.AudioDevices = append(audioGroup.AudioDevices, audioDeviceConfig)
			}

			controlGroup.Audio.Groups = append(controlGroup.Audio.Groups, audioGroup)
		}

		for _, cam := range v.Cameras {
			camera := ui.CameraConfig{
				Name: cam.Name,
			}

			camera.TiltUp, err = cam.TiltUp.convert()
			if err != nil {
				return config, fmt.Errorf("invalid control set '%s.tiltUp': %w", cam.Name, err)
			}

			camera.TiltDown, err = cam.TiltDown.convert()
			if err != nil {
				return config, fmt.Errorf("invalid control set '%s.tiltDown': %w", cam.Name, err)
			}

			camera.PanLeft, err = cam.PanLeft.convert()
			if err != nil {
				return config, fmt.Errorf("invalid control set '%s.panLeft': %w", cam.Name, err)
			}

			camera.PanRight, err = cam.PanRight.convert()
			if err != nil {
				return config, fmt.Errorf("invalid control set '%s.panRight': %w", cam.Name, err)
			}

			camera.PanTiltStop, err = cam.PanTiltStop.convert()
			if err != nil {
				return config, fmt.Errorf("invalid control set '%s.panTiltStop': %w", cam.Name, err)
			}

			camera.ZoomIn, err = cam.ZoomIn.convert()
			if err != nil {
				return config, fmt.Errorf("invalid control set '%s.zoomIn': %w", cam.Name, err)
			}

			camera.ZoomOut, err = cam.ZoomOut.convert()
			if err != nil {
				return config, fmt.Errorf("invalid control set '%s.zoomOut': %w", cam.Name, err)
			}

			camera.ZoomStop, err = cam.ZoomStop.convert()
			if err != nil {
				return config, fmt.Errorf("invalid control set '%s.zoomStop': %w", cam.Name, err)
			}

			for _, pre := range cam.Presets {
				preset := ui.CameraPresetConfig{
					Name: pre.Name,
				}

				preset.SetPreset, err = pre.SetPreset.convert()
				if err != nil {
					return config, fmt.Errorf("invalid control set '%s.%s.zoomStop': %w", cam.Name, pre.Name, err)
				}

				camera.Presets = append(camera.Presets, preset)
			}

			controlGroup.Cameras = append(controlGroup.Cameras, camera)
		}

		config.ControlGroups[k] = controlGroup
	}

	return config, nil
}

func (c controlSet) convert() (ui.ControlSet, error) {
	cs := ui.ControlSet{
		APIRequest: c.APIRequest,
		Requests:   make([]ui.GenericControlRequest, len(c.Requests)),
	}

	for i := range c.Requests {
		cs.Requests[i].Body = c.Requests[i].Body
		cs.Requests[i].Method = c.Requests[i].Method

		u, err := url.Parse(c.Requests[i].URL)
		if err != nil {
			return cs, err
		}

		cs.Requests[i].URL = u
	}

	return cs, nil
}
