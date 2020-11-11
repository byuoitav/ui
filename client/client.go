package client

import (
	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/ui"
)

var _ ui.Client = &client{}

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
			PoweredOn: !c.stateMatches(cGroup.PowerOff.APIRequest),
		}

		// build each display group
		// TODO sharing
		for dispName, cDisplay := range cGroup.Displays {
			display := DisplayGroup{
				ID: dispName,
				Displays: []IconPair{
					{
						Name: dispName,
						Icon: cDisplay.Icon,
					},
				},
			}

			// build each of the sources
			// TODO subinputs
			for sourceName, cSource := range cDisplay.Sources {
				input := Input{
					IconPair: IconPair{
						Name: sourceName,
						Icon: cSource.Icon,
					},
				}

				var curInput bool
				if display.Input == "" {
					curInput = c.stateMatches(cSource.APIRequest)
					if curInput {
						display.Input = sourceName
					}
				}

				if cSource.Visible || curInput {
					display.Inputs = append(display.Inputs, input)
				}
			}
		}

		// build media audio info
		group.MediaAudio.Level = c.getVolume(cGroup.Audio.Media.APIRequest, c.state)
		group.MediaAudio.Muted = c.getMuted(cGroup.Audio.Media.APIRequest, c.state)

		// build audio groups
		for gID, cAudioGroup := range cGroup.Audio.Groups {
			audioGroup := AudioGroup{
				ID:    gID,
				Name:  gID,
				Muted: true,
			}

			for aID, cAudio := range cAudioGroup {
				audio := AudioDevice{
					// TODO need to get icon (need to change config)
					ID:    aID,
					Level: c.getVolume(cAudio.APIRequest, c.state),
					Muted: c.getMuted(cAudio.APIRequest, c.state),
				}

				if !audio.Muted {
					audioGroup.Muted = false
				}

				audioGroup.AudioDevices = append(audioGroup.AudioDevices, audio)
			}

			group.AudioGroups = append(group.AudioGroups, audioGroup)
		}

		room.ControlGroups[id] = group

		// TODO controlInfo
		// TODO presentGroups
	}

	return room
}
