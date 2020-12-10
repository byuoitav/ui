package client

import (
	"sync"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/ui"
	"go.uber.org/zap"
)

var _ ui.Client = &client{}

type client struct {
	// static info about the room
	roomID         string
	controlGroupID string

	// structs to ~do stuff~ with
	dataService  ui.DataService
	avController ui.AVController
	publisher    ui.EventPublisher
	log          *zap.Logger

	state   avcontrol.StateResponse
	stateMu sync.RWMutex

	config   ui.Config
	configMu sync.RWMutex

	handlers map[string]messageHandler

	outgoing chan []byte
	kill     chan struct{}
	killOnce sync.Once

	// TODO controlKey/url
}

func (c *client) Close() {
	c.killOnce.Do(func() {
		c.log.Info("Closing client")

		// close the kill chan to clean up all resources
		close(c.kill)

		// close outgoing chan to make sure no more messages are sent
		close(c.outgoing)
	})
}

func (c *client) Done() <-chan struct{} {
	return c.kill
}

// TODO
// Things we need:
// A list of controlGroups they are allowed to switch to
func (c *client) Room() Room {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()

	c.configMu.RLock()
	defer c.configMu.RUnlock()

	room := Room{
		Name:                 c.roomID,
		ControlGroups:        make(map[string]ControlGroup),
		SelectedControlGroup: c.controlGroupID,
	}

	for cgName, cg := range c.config.ControlGroups {
		group := ControlGroup{
			Name: cgName,
			Support: Support{ // TODO this should be pulled from cache
				HelpRequested: false,
				HelpMessage:   "Request Help",
				HelpEnabled:   true,
			},
			PoweredOn: !c.doesStateMatch(cg.PowerOff.MatchStates...),
		}

		// build each display group
		// TODO sharing
		for _, disp := range cg.Displays {
			display := DisplayGroup{
				Name:    disp.Name,
				Blanked: c.doesStateMatch(disp.Blank.MatchStates...),
				Displays: []IconPair{
					{
						Name: disp.Name,
						Icon: disp.Icon,
					},
				},
			}

			// build each of the sources
			for _, source := range disp.Sources {
				input := Input{
					IconPair: IconPair{
						Name: source.Name,
						Icon: source.Icon,
					},
				}

				var curInput bool
				if display.Input == "" {
					curInput = c.doesStateMatch(source.MatchStates...)
					if curInput {
						display.Input = source.Name
					}
				}

				if source.Visible || curInput {
					display.Inputs = append(display.Inputs, input)
				}
			}

			group.DisplayGroups = append(group.DisplayGroups, display)
		}

		// build media audio info
		group.MediaAudio.Level = c.getVolume(cg.MediaAudio.Volume.MatchStates...)
		group.MediaAudio.Muted = c.doesStateMatch(cg.MediaAudio.Mute.MatchStates...)

		// build audio groups
		for _, ag := range cg.AudioGroups {
			audioGroup := AudioGroup{
				Name: ag.Name,
				// Muted is true if all of the audio devices in this group
				// are muted
				Muted: true,
			}

			for _, ad := range ag.AudioDevices {
				audioDevice := AudioDevice{
					IconPair: IconPair{
						Name: ad.Name,
					},
					Level: c.getVolume(ad.Volume.MatchStates...),
					Muted: c.doesStateMatch(ad.Mute.MatchStates...),
				}

				if !audioDevice.Muted {
					audioGroup.Muted = false
				}

				audioGroup.AudioDevices = append(audioGroup.AudioDevices, audioDevice)
			}

			group.AudioGroups = append(group.AudioGroups, audioGroup)
		}

		for _, cam := range cg.Cameras {
			camera := Camera{
				Name: cam.Name,
			}

			for _, pre := range cam.Presets {
				camera.Presets = append(camera.Presets, pre.Name)
			}

			group.Cameras = append(group.Cameras, camera)
		}

		room.ControlGroups[cgName] = group
	}

	return room
}
