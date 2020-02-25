package bff

func containsID(ids []ID, id ID) int {
	for index, i := range ids {
		if i == id {
			return index
		}
	}
	return -1
}

// GetRoom .
func (c *Client) GetRoom() Room {
	room := Room{
		ID:                   ID(c.roomID),
		Name:                 c.room.Name,
		ControlGroups:        make(map[string]ControlGroup),
		SelectedControlGroup: ID(c.selectedControlGroupID),
	}

	var masters []ID
	active := make(map[ID]ID)
	inactive := make(map[ID]ID)

	c.shareMutex.RLock()
	for master, mins := range c.sharing {
		masters = append(masters, master)
		for _, a := range mins.Active {
			active[a] = master
		}
		for _, i := range mins.Inactive {
			inactive[i] = master
		}
	}
	c.shareMutex.RUnlock()

	for _, preset := range c.uiConfig.Presets {
		cg := ControlGroup{
			ID:   ID(preset.Name),
			Name: preset.Name,
			Support: Support{
				HelpRequested: false, // This info also needs to be saved...
				HelpMessage:   "Request Help",
				HelpEnabled:   true,
			},
		}

		power := true

		for _, name := range preset.Displays {
			config := GetDeviceConfigByName(c.room.Devices, name)
			state := GetDisplayStateByName(c.state.Displays, name)
			outputIcon := "tv"

			for _, IOconfig := range c.uiConfig.OutputConfiguration {
				if config.Name != IOconfig.Name {
					continue
				}

				outputIcon = IOconfig.Icon
			}

			// If any displays has its power off then the room is not entirely on
			if state.Power != "on" {
				power = false
			}

			// figure out what the current input for this display is
			// we are assuming that input is roomID - input name
			// unless it's blanked, then the "input" is blank
			curInput := c.roomID + "-" + state.Input
			if state.Input == "" {
				curInput = c.roomID + "-" + preset.Inputs[0]
			}
			blanked := false
			if state.Blanked != nil && *state.Blanked {
				blanked = true
			}

			s := ShareInfo{
				Options: preset.ShareableDisplays,
			}

			// Set the different possible share states of a room
			if m := containsID(masters, ID(name)); m >= 0 {
				s.State = Unshare
			} else if master, ok := active[ID(name)]; ok {
				s.State = MinionActive
				s.Master = master
			} else if master, ok := inactive[ID(name)]; ok {
				s.State = MinionInactive
				s.Master = master
			} else if _, ok := c.shareable[ID(name)]; ok {
				s.State = Share
			} else /*else if linkable?!?!*/ {
				s.State = Nothing
			}
			if s.State == MinionActive {
				c.shareMutex.RLock()
				curInput = string(c.sharing[s.Master].Input)
				c.shareMutex.RUnlock()
			} else if s.State == MinionInactive {
				cg.Inputs = append(cg.Inputs, Input{
					ID: ID("Mirror ") + s.Master,
					IconPair: IconPair{
						Name: "Mirror " + string(s.Master),
						Icon: Icon{"settings_input_hdmi"},
					},
					Disabled: false,
				})
			}
			d := DisplayBlock{
				ID:      ID(config.ID),
				Blanked: blanked,
				Input:   ID(curInput),
				Share:   s,
			}

			// TODO outputs when we do sharing
			d.Outputs = append(d.Outputs, IconPair{
				ID:   ID(config.ID),
				Name: config.DisplayName,
				Icon: Icon{outputIcon},
			})

			cg.DisplayBlocks = append(cg.DisplayBlocks, d)
		}

		if power {
			cg.Power = "on"
		} else {
			cg.Power = "standby"
		}

		// add a blank input as the first input if we aren't on blueberry
		cg.Inputs = append(cg.Inputs, Input{
			ID: ID("blank"),
			IconPair: IconPair{
				Name: "Blank",
				Icon: Icon{"crop_landscape"},
			},
			Disabled: false,
		})

		for _, name := range preset.Inputs {
			config := GetDeviceConfigByName(c.room.Devices, name)
			inputIcon := "settings_input_hdmi"

			for _, IOconfig := range c.uiConfig.InputConfiguration {
				if config.Name != IOconfig.Name {
					continue
				}

				inputIcon = IOconfig.Icon
			}

			i := Input{
				ID: ID(config.ID),
				IconPair: IconPair{
					Name: config.DisplayName,
					Icon: Icon{inputIcon},
				},
				Disabled: false, // TODO look at the current displays reachable inputs to determine
			}

			// TODO subinputs

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
					audioIcon := "mic"
					for _, IOconfig := range c.uiConfig.OutputConfiguration {
						if config.Name != IOconfig.Name {
							continue
						}

						audioIcon = IOconfig.Icon
					}

					dev := AudioDevice{
						ID: ID(config.ID),
						IconPair: IconPair{
							Name: config.DisplayName,
							Icon: Icon{audioIcon},
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

				dev := AudioDevice{
					ID: ID(config.ID),
					IconPair: IconPair{
						Name: config.DisplayName,
						Icon: Icon{"mic"}, // TODO get mic icon from outputconfig
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
