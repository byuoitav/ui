package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/byuoitav/ui"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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
		return
	}

	// could we combine these requests into a single request somehow?
	errg := &errgroup.Group{}
	for i := range controlSets {
		errg.Go(func() error {
			return c.doControlSet(ctx, controlSets[i])
		})
	}

	if err := errg.Wait(); err != nil {
		c.log.Error("unable to set input", zap.Error(err))
	}

	/*
		// ugh. what should the device be for this
		// i guess we need to have a general discussion about this - what should be the devices for our events
		event := ui.Event{
			Room:   c.roomID,
			Device: msg.DisplayGroup,
			Tags:   []string{"core-state"},
			Key:    "input",
			Value:  msg.Source,
			// get ip from client
		}

		if msg.SubSource != "" {
			event.Value = fmt.Sprintf("%s.%s", msg.Source, msg.SubSource)
		}

		c.publisher.Publish(context.Background(), event)
	*/
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
