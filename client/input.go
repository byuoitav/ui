package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/byuoitav/ui"
)

func (c *client) setInput(data []byte) {
	var msg struct {
		DisplayGroup string `json:"displayGroup"`
		Source       string `json:"source"`
		SubSource    string `json:"subSource"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	room := c.Room()

	// find all of the controlSets for the affected displays
	var controlSets []ui.ControlSet
	for _, dg := range room.ControlGroups[c.controlGroupID].DisplayGroups {
		if msg.DisplayGroup != dg.Name {
			continue
		}

		for _, d := range dg.Displays {
			// find the matching controlSet in the config
			cfg, ok := c.getSourceConfig(d.Name, msg.Source, msg.SubSource)
			if !ok {
				// error
				return
			}

			controlSets = append(controlSets, cfg.ControlSet)
		}
	}

	if len(controlSets) == 0 {
		// some error
	}

	// TODO parallelize this?
	for _, set := range controlSets {
		c.doControlSet(ctx, set)
	}
}

func (c *client) getSourceConfig(disp, src, subSrc string) (ui.SourceConfig, bool) {
	c.configMu.RLock()
	defer c.configMu.RUnlock()

	for _, cDisp := range c.config.ControlGroups[c.controlGroupID].Displays {
		if cDisp.Name == disp {
			for _, cSrc := range cDisp.Sources {
				if cSrc.Name == src {
					if subSrc == "" {
						return cSrc, true
					}

					for _, cSubSrc := range cSrc.Sources {
						if cSubSrc.Name == subSrc {
							return cSubSrc, true
						}
					}
				}
			}
		}
	}

	return ui.SourceConfig{}, false
}
