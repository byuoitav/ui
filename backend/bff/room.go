package bff

import (
	"errors"
	"fmt"
)

// GetRoom .
func (c *Client) GetRoom() Room {
	c.controlKeysMu.RLock()
	defer c.controlKeysMu.RUnlock()

	room := Room{
		ID:                   ID(c.roomID),
		Name:                 c.room.Name,
		ControlGroups:        make(map[string]ControlGroup),
		SelectedControlGroup: ID(c.selectedControlGroupID),
	}

	// create all of the presets for this room
	for _, preset := range c.uiConfig.Presets {
		cg := ControlGroup{
			ID:   ID(preset.Name),
			Name: preset.Name,
			Support: Support{
				HelpRequested: false, // TODO this info should be pulled from lazarette
				HelpMessage:   "Request Help",
				HelpEnabled:   true,
			},
		}

		cg.ControlInfo.Key = c.controlKeys[preset.Name]
		cg.ControlInfo.URL = c.config.RemoteControlAddr

		// poweredOn should be true if all of the displays are powered on
		poweredOn := true

		// create a displaygroup for each of the preset's displays
		// right now, each display in the room gets its own group
		// so for now, a display group will only have 1 display in it
		for _, name := range preset.Displays {
			config := GetDeviceConfigByName(c.room.Devices, name)
			state := GetDisplayStateByName(c.state.Displays, name)
			outputIcon := "tv" // default icon

			// find this icons display
			for _, IOconfig := range c.uiConfig.OutputConfiguration {
				if config.Name != IOconfig.Name {
					continue
				}

				outputIcon = IOconfig.Icon
			}

			// If any displays has its power off then the room is not entirely on
			if state.Power != "on" {
				poweredOn = false
			}

			// figure out what the current input for this display is
			// we are assuming that input is 'roomID-inputname'
			// if the input isn't set, then we default to the first input
			curInput := fmt.Sprintf("%s-%s", c.roomID, state.Input)
			if len(state.Input) == 0 && len(preset.Inputs) > 0 {
				curInput = fmt.Sprintf("%s-%s", c.roomID, preset.Inputs[0])
			}

			blanked := false
			if state.Blanked != nil && *state.Blanked {
				blanked = true
			}

			group := DisplayGroup{
				ID:      ID(config.ID),
				Blanked: blanked,
				Input:   ID(curInput),
				ShareInfo: ShareInfo{
					State: stateCantShare,
				},
			}

			shareData, err := c.getShareData(group.ID)
			switch {
			case len(preset.ShareableDisplays) == 0:
				group.ShareInfo.State = stateCantShare
			case err != nil:
				// if there is no share data (yet), but there are sharable displays
				// then allow them to share to those options
				group.ShareInfo.State = stateCanShare
				group.ShareInfo.Options = convertNamesToIDStrings(c.roomID, preset.ShareableDisplays)
			case shareData.State == stateIsMaster:
				group.ShareInfo.State = shareData.State
			case shareData.State == stateIsActiveMinion || shareData.State == stateIsInactiveMinion:
				group.ShareInfo.State = shareData.State
				group.ShareInfo.Master = shareData.Master
			default:
				group.ShareInfo.State = shareData.State
				group.ShareInfo.Options = convertNamesToIDStrings(c.roomID, preset.ShareableDisplays)
				group.ShareInfo.Master = shareData.Master
			}

			group.Displays = append(group.Displays, IconPair{
				ID:   ID(config.ID),
				Name: config.DisplayName,
				Icon: outputIcon,
			})

			cg.DisplayGroups = append(cg.DisplayGroups, group)
		}

		cg.PoweredOn = poweredOn

		// create the list of inputs available in this control group
		// TODO subinputs
		for _, name := range preset.Inputs {
			config := GetDeviceConfigByName(c.room.Devices, name)
			icon := "settings_input_hdmi" // default icon

			// figure out which icon to use
			for _, IOconfig := range c.uiConfig.InputConfiguration {
				if config.Name != IOconfig.Name {
					continue
				}

				icon = IOconfig.Icon
			}

			cg.Inputs = append(cg.Inputs, Input{
				ID: ID(config.ID),
				IconPair: IconPair{
					Name: config.DisplayName,
					Icon: icon,
				},
			})
		}

		// create an extra input if our ONLY display group is an inactive minion
		// the input will let them become an active minion again
		if len(cg.DisplayGroups) == 1 && cg.DisplayGroups[0].ShareInfo.State == stateIsInactiveMinion {
			cg.Inputs = append(cg.Inputs, Input{
				ID: ID(inputBecomeActive),
				IconPair: IconPair{
					Name: string(cg.DisplayGroups[0].ShareInfo.Master),
					Icon: "share",
				},
			})
		}

		// create this cg's media audio info
		// MediaAudio information is tied to the audioDevices array from the preset
		// MediaAudio.Muted is true if ALL of the devices are muted
		// MediaAudio.Level is the average level of the devices
		cg.MediaAudio.Muted = true
		for _, name := range preset.AudioDevices {
			state := GetAudioDeviceStateByName(c.state.AudioDevices, name)
			if state.Volume != nil {
				cg.MediaAudio.Level += (*state.Volume)
			}

			if state.Muted != nil && !(*state.Muted) {
				cg.MediaAudio.Muted = false
			}
		}
		if len(preset.AudioDevices) == 0 {
			c.Out <- ErrorMessage(errors.New("Caleb was Actually right and caught a divide-by-zero error"))
			cg.MediaAudio.Level = 69
		} else {
			cg.MediaAudio.Level /= len(preset.AudioDevices)
		}

		// create the cg's audio groups.
		// if audioGroups are present in the config, then use those.
		// if not, create a mics audio group with all of the independentAudioDevices in it
		// if there are no audioGroups or independentAudioDevices, dont't create any groups
		if len(preset.AudioGroups) > 0 {
			// create a group for each audioGroup in the preset
			for id, audioDevices := range preset.AudioGroups {
				group := AudioGroup{
					ID:    ID(id),
					Name:  id,
					Muted: true,
				}

				for _, name := range audioDevices {
					config := GetDeviceConfigByName(c.room.Devices, name)
					state := GetAudioDeviceStateByName(c.state.AudioDevices, name)

					// figure out which icon to use - default to 'mic'
					icon := "mic"
					for _, IOconfig := range c.uiConfig.OutputConfiguration {
						if config.Name != IOconfig.Name {
							continue
						}

						icon = IOconfig.Icon
					}

					dev := AudioDevice{
						ID: ID(config.ID),
						IconPair: IconPair{
							Name: config.DisplayName,
							Icon: icon,
						},
					}

					if state.Volume != nil {
						dev.Level = *state.Volume
					}

					if state.Muted != nil {
						dev.Muted = *state.Muted
					}

					if !dev.Muted {
						group.Muted = false
					}

					group.AudioDevices = append(group.AudioDevices, dev)
				}

				cg.AudioGroups = append(cg.AudioGroups, group)
			}
		} else if len(preset.IndependentAudioDevices) > 0 {
			// create an audio group for all of the independentAudioDevices
			group := AudioGroup{
				ID:    "micsAG",
				Name:  "Microphones",
				Muted: true,
			}

			for _, name := range preset.IndependentAudioDevices {
				config := GetDeviceConfigByName(c.room.Devices, name)
				state := GetAudioDeviceStateByName(c.state.AudioDevices, name)
				icon := "mic"

				// figure out which icon to use
				for _, IOconfig := range c.uiConfig.OutputConfiguration {
					if config.Name != IOconfig.Name {
						continue
					}

					icon = IOconfig.Icon
				}

				dev := AudioDevice{
					ID: ID(config.ID),
					IconPair: IconPair{
						Name: config.DisplayName,
						Icon: icon,
					},
				}

				if state.Volume != nil {
					dev.Level = *state.Volume
				}

				if state.Muted != nil {
					dev.Muted = *state.Muted
				}

				if !dev.Muted {
					group.Muted = false
				}

				group.AudioDevices = append(group.AudioDevices, dev)
			}

			cg.AudioGroups = append(cg.AudioGroups, group)
		}

		// set this cg in the controlgroups map
		room.ControlGroups[string(cg.ID)] = cg
	}

	return room
}

func convertNamesToIDStrings(roomID string, names []string) []string {
	var ids []string
	for i := range names {
		ids = append(ids, fmt.Sprintf("%s-%s", roomID, names[i]))
	}

	return ids
}
