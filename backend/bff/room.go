package bff

import (
	"errors"
	"fmt"
	"sort"
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

			// Add share info
			shareData, err := c.getShareData(group.ID)
			switch {
			//If you are blueberry and have no shareable displays
			case len(preset.ShareableDisplays) == 0 && len(preset.Displays) == 1:
				group.ShareInfo.State = stateCantShare

			// No shareable data found
			case err != nil:
				// if there is no share data (yet), but there are sharable displays
				// then allow them to share to those options
				group.ShareInfo.State = stateCanShare

				if len(preset.Displays) == 1 {
					// blueberry case
					group.ShareInfo.Options = convertNamesToIDStrings(c.roomID, preset.ShareableDisplays)
				} else {
					// cherry case
					var options []string
					for _, option := range preset.Displays {
						if option != name {
							options = append(options, option)
						}
					}

					group.ShareInfo.Options = convertNamesToIDStrings(c.roomID, options)
				}

			case shareData.State == stateIsMaster:
				group.ShareInfo.State = shareData.State
				outputIcon = "dynamic_feed"

			case shareData.State == stateIsActiveMinion || shareData.State == stateIsInactiveMinion:
				group.ShareInfo.State = shareData.State
				group.ShareInfo.Master = shareData.Master

			default:
				group.ShareInfo.State = shareData.State
				group.ShareInfo.Master = shareData.Master

				// blueberry case
				if len(preset.Displays) == 1 {
					group.ShareInfo.Options = convertNamesToIDStrings(c.roomID, preset.ShareableDisplays)
				} else { // cherry case
					var options []string
					for _, option := range preset.Displays {
						if option != name {
							options = append(options, option)
						}
					}
					group.ShareInfo.Options = convertNamesToIDStrings(c.roomID, options)
				}
			}

			group.Displays = append(group.Displays, IconPair{
				ID:   ID(config.ID),
				Name: config.DisplayName,
				Icon: outputIcon,
			})

			cg.DisplayGroups = append(cg.DisplayGroups, group)
		}

		cg.fullDisplayGroups = append(cg.fullDisplayGroups, cg.DisplayGroups...)

		// check displays groups that i need to get rid of (cherry)
		if len(cg.DisplayGroups) > 1 {
			var keep DisplayGroups

			for i := range cg.DisplayGroups {
				if cg.DisplayGroups[i].ShareInfo.State == stateCanShare || cg.DisplayGroups[i].ShareInfo.State == stateIsMaster {
					keep = append(keep, cg.DisplayGroups[i])
				}
			}

			// add minions to their masters
			for i := range cg.DisplayGroups {
				if cg.DisplayGroups[i].ShareInfo.State == stateIsActiveMinion {
					// find the master
					for j := range keep {
						if keep[j].ID == cg.DisplayGroups[i].ShareInfo.Master {
							keep[j].Displays = append(keep[j].Displays, cg.DisplayGroups[i].Displays...)
						}
					}
				}
			}

			// recreate display groups for this controlgroup
			cg.DisplayGroups = nil
			cg.DisplayGroups = append(cg.DisplayGroups, keep...)
		}

		// Now sort the display groups by size
		sort.SliceStable(cg.DisplayGroups, func(i, j int) bool {
			return len(cg.DisplayGroups[i].Displays) > len(cg.DisplayGroups[j].Displays)
		})

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
		//fmt.Printf("DisplayGroup Count: %v || State: %v\n", len(cg.DisplayGroups), cg.DisplayGroups[0].ShareInfo.State)
		if len(cg.DisplayGroups) == 1 && cg.DisplayGroups[0].ShareInfo.State == stateIsInactiveMinion {
			cg.Inputs = append(cg.Inputs, Input{
				ID: ID(inputBecomeActive),
				IconPair: IconPair{
					Name: string(cg.DisplayGroups[0].ShareInfo.Master),
					Icon: "share",
				},
			})
		}

		// create this control-groups's media audio info
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
			c.Out <- ErrorMessage(errors.New("caleb was actually right and caught a divide-by-zero error"))
			cg.MediaAudio.Level = 69
		} else {
			cg.MediaAudio.Level /= len(preset.AudioDevices)
		}

		// create the control-groups's audio groups.
		// if audioGroups are present in the config, then use those.
		// if not, create a mics audio group with all of the independentAudioDevices in it
		// if there are no audioGroups or independentAudioDevices, don't create any groups
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

		// order the audiogroups alphabetically
		sort.Slice(cg.AudioGroups, func(i, j int) bool {
			return len(cg.AudioGroups[i].ID) < len(cg.AudioGroups[j].ID)
		})

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
