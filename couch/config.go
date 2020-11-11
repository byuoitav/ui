package couch

import (
	"context"
	"encoding/json"
	"fmt"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/ui"
)

type config struct {
	ID            string `json:"_id"`
	ControlPanels map[string]struct {
		UIType string `json:"uiType"`

		// TODO divider sensors
		ControlGroup string `json:"controlGroup"`
	} `json:"controlPanels"`
	ControlGroups map[string]struct {
		PowerOff controlSet `json:"powerOff"`
		PowerOn  controlSet `json:"powerOn"`
		Displays map[string]struct {
			Icon    string `json:"icon"`
			Visible bool   `json:"visible"`
			controlSet

			Sources map[string]struct {
				Icon    string `json:"icon"`
				Visible bool   `json:"visible"`
				controlSet

				Sources map[string]struct {
					Icon    string `json:"icon"`
					Visible bool   `json:"visible"`
					controlSet
				} `json:"sources"`
			} `json:"sources"`
		} `json:"displays"`
		Audio struct {
			Media  controlSet                       `json:"media"`
			Groups map[string]map[string]controlSet `json:"groups"`
		} `json:"audio"`
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

	return config.convert(), nil
}

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
			Displays: make(map[string]ui.DisplayConfig, len(v.Displays)),
		}

		controlGroup.Audio.Media = v.Audio.Media.convert()
		controlGroup.Audio.Groups = make(map[string]map[string]ui.ControlSet, len(v.Audio.Groups))

		for id, disp := range v.Displays {
			uiDisp := ui.DisplayConfig{
				Icon:    disp.Icon,
				Sources: make(map[string]ui.SourceConfig, len(disp.Sources)),
			}

			for name, source := range disp.Sources {
				uiDisp.Sources[name] = ui.SourceConfig{
					Icon:       source.Icon,
					Visible:    source.Visible,
					ControlSet: source.controlSet.convert(),
					Sources:    make(map[string]ui.SourceConfig, len(source.Sources)),
				}

				for subName, subSource := range source.Sources {
					uiDisp.Sources[name].Sources[subName] = ui.SourceConfig{
						Icon:       subSource.Icon,
						Visible:    subSource.Visible,
						ControlSet: subSource.controlSet.convert(),
					}
				}
			}

			controlGroup.Displays[id] = uiDisp
		}

		for gName, g := range v.Audio.Groups {
			group := make(map[string]ui.ControlSet)
			for name, set := range g {
				group[name] = set.convert()
			}

			controlGroup.Audio.Groups[gName] = group
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
