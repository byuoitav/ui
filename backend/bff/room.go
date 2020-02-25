package bff

import "fmt"

//func containsID(ids []ID, id ID) int {
//	for index, i := range ids {
//		if i == id {
//			return index
//		}
//	}
//
//	return -1
//}

// GetRoom .
func (c *Client) GetRoom() Room {
	room := Room{
		ID:                   ID(c.roomID),
		Name:                 c.room.Name,
		ControlGroups:        make(map[string]ControlGroup),
		SelectedControlGroup: ID(c.selectedControlGroupID),
	}

	//var masters []ID
	//active := make(map[ID]ID)
	//inactive := make(map[ID]ID)

	//c.shareMutex.RLock()
	//for master, mins := range c.sharing {
	//	masters = append(masters, master)
	//	for _, a := range mins.Active {
	//		active[a] = master
	//	}
	//	for _, i := range mins.Inactive {
	//		inactive[i] = master
	//	}
	//}
	//c.shareMutex.RUnlock()

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

			//s := ShareInfo{
			//	Options: preset.ShareableDisplays,
			//}

			//// Set the different possible share states of a room
			//if m := containsID(masters, ID(name)); m >= 0 {
			//	s.State = Unshare
			//} else if master, ok := active[ID(name)]; ok {
			//	s.State = MinionActive
			//	s.Master = master
			//} else if master, ok := inactive[ID(name)]; ok {
			//	s.State = MinionInactive
			//	s.Master = master
			//} else if _, ok := c.shareable[ID(name)]; ok {
			//	s.State = Share
			//} else /*else if linkable?!?!*/ {
			//	s.State = Nothing
			//}
			//if s.State == MinionActive {
			//	c.shareMutex.RLock()
			//	curInput = string(c.sharing[s.Master].Input)
			//	c.shareMutex.RUnlock()
			//} else if s.State == MinionInactive {
			//	cg.Inputs = append(cg.Inputs, Input{
			//		ID: ID("Mirror ") + s.Master,
			//		IconPair: IconPair{
			//			Name: "Mirror " + string(s.Master),
			//			Icon: Icon{"settings_input_hdmi"},
			//		},
			//		Disabled: false,
			//	})
			//}

			group := DisplayGroup{
				ID:      ID(config.ID),
				Blanked: blanked,
				Input:   ID(curInput),
				// Share:   s,
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

			i := Input{
				ID: ID(config.ID),
				IconPair: IconPair{
					Name: config.DisplayName,
					Icon: icon,
				},
			}

			cg.Inputs = append(cg.Inputs, i)
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
		cg.MediaAudio.Level /= len(preset.AudioDevices)

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
