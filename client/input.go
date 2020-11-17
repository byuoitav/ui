package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// TODO displayGroups
func (c *client) setInput(data []byte) {
	var msg struct {
		Display   string `json:"display"`
		Source    string `json:"source"`
		SubSource string `json:"subSource"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cg, ok := c.config.ControlGroups[c.controlGroupID]
	if !ok {
		// TODO log/send invalid control group error
		return
	}

	for _, disp := range cg.Displays {
		if msg.Display != disp.Name {
			continue
		}

		for _, src := range disp.Sources {
			if msg.Source != src.Name {
				continue
			}

			if msg.SubSource == "" {
				c.doControlSet(ctx, src.ControlSet)
				return
			}

			for _, subSrc := range src.Sources {
				if msg.SubSource == subSrc.Name {
					c.doControlSet(ctx, subSrc.ControlSet)
					return
				}
			}

			// error about subsource
		}

		// error about source
	}

	// error about display
}
