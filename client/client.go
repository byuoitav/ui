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
	log          *zap.Logger

	state   avcontrol.StateResponse
	stateMu sync.RWMutex

	config   ui.Config
	configMu sync.RWMutex

	handlers map[string]messageHandler

	outgoing chan []byte

	// TODO controlKey/url
}

// TODO
// Things we need:
// A list of controlGroups they are allowed to switch to
// switch back to lists instead of maps for order
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
			PoweredOn: !c.stateMatches(cg.PowerOff.APIRequest),
		}

		// build each display group
		// TODO sharing
		for _, disp := range cg.Displays {
			display := DisplayGroup{
				Name: disp.Name,
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
					curInput = c.stateMatches(source.APIRequest)
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
		group.MediaAudio.Level = c.getVolume(cg.Audio.Media.Volume.APIRequest)
		group.MediaAudio.Muted = c.stateMatches(cg.Audio.Media.Mute.APIRequest)

		// build audio groups
		for _, ag := range cg.Audio.Groups {
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
					Level: c.getVolume(ad.Volume.APIRequest),
					Muted: c.stateMatches(ad.Mute.APIRequest),
				}

				if !audioDevice.Muted {
					audioGroup.Muted = false
				}

				audioGroup.AudioDevices = append(audioGroup.AudioDevices, audioDevice)
			}

			group.AudioGroups = append(group.AudioGroups, audioGroup)
		}

		/*
			for _, cam := range cGroup.Cameras {
				camera := Camera{
					DisplayName: cam.DisplayName,
					TiltUp:      cam.TiltUp,
					TiltDown:    cam.TiltDown,
					PanLeft:     cam.PanLeft,
					PanRight:    cam.PanRight,
					PanTiltStop: cam.PanTiltStop,
					ZoomIn:      cam.ZoomIn,
					ZoomOut:     cam.ZoomOut,
					ZoomStop:    cam.ZoomStop,
					Stream:      cam.Stream,
					Reboot:      cam.Reboot,
				}

				for _, preset := range cam.Presets {
					pre := CameraPreset{
						DisplayName: preset.DisplayName,
						SetPreset:   preset.SetPreset,
						SavePreset:  preset.SavePreset,
					}

					camera.Presets = append(camera.Presets, pre)
				}

				group.Cameras = append(group.Cameras, camera)
			}
		*/

		room.ControlGroups[cgName] = group
	}

	return room
}
