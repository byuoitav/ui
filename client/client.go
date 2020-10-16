package client

import (
	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/ui"
)

type client struct {
	// static info about the room
	roomID         string
	controlGroupID string

	// structs to ~do stuff~ with
	dataService  ui.DataService
	avController ui.AVController

	// info that we update occasionally
	state  avcontrol.StateResponse
	config ui.Config

	// TODO controlKey/url
}

// TODO
// Things we need:
// A list of controlGroups they are allowed to switch to
func (c *client) Room() Room {
	room := Room{
		ID:                   c.roomID,
		ControlGroups:        make(map[string]ControlGroup),
		SelectedControlGroup: c.controlGroupID,
	}

	// build all of the controlGroups
	for id, cGroup := range c.config.ControlGroups {
		group := ControlGroup{
			ID: id,
			Support: Support{ // TODO this should be pulled from cache
				HelpRequested: false,
				HelpMessage:   "Request Help",
				HelpEnabled:   true,
			},
			PoweredOn: true, // TODO need to add to config, how to decide if poweredOn - or is it just !PoweredOff?
		}

		// TODO sharing
		for _, cDisplay := range cGroup.Displays {
			display := DisplayGroup{
				ID: cDisplay.Name,
				Displays: []IconPair{
					{
						Name: cDisplay.Name,
						Icon: cDisplay.Icon,
					},
				},
			}

			// build each of the sources
			// TODO subinputs
			for _, cSource := range cDisplay.Sources {
				input := Input{
					IconPair: IconPair{
						Name: cSource.Name,
						Icon: cSource.Icon,
					},
				}

				var curInput bool
				if display.Input == "" {
					curInput = c.stateMatches(cSource.APIRequest)
					if curInput {
						display.Input = cSource.Name
					}
				}

				if cSource.Visible || curInput {
					display.Inputs = append(display.Inputs, input)
				}
			}
		}

		room.ControlGroups[id] = group
	}

	return room
}
