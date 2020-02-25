package bff

import (
	"fmt"

	"github.com/byuoitav/common/structs"
)

func GetDeviceConfigByID(devices []structs.Device, id ID) structs.Device {
	for i := range devices {
		if id == ID(devices[i].ID) {
			return devices[i]
		}
	}

	return structs.Device{}
}

func GetDeviceConfigByName(devices []structs.Device, name string) structs.Device {
	for i := range devices {
		if name == devices[i].Name {
			return devices[i]
		}
	}

	return structs.Device{}
}

func GetDisplayStateByName(displays []structs.Display, name string) structs.Display {
	for i := range displays {
		if name == displays[i].Name {
			return displays[i]
		}
	}

	return structs.Display{}
}

func GetAudioDeviceStateByName(audioDevices []structs.AudioDevice, name string) structs.AudioDevice {
	for i := range audioDevices {
		if name == audioDevices[i].Name {
			return audioDevices[i]
		}
	}

	return structs.AudioDevice{}
}

func GetAudioDeviceByID(audioGroups []AudioGroup, id ID) (AudioDevice, error) {
	for i := range audioGroups {
		for j := range audioGroups[i].AudioDevices {
			if id == audioGroups[i].AudioDevices[j].ID {
				return audioGroups[i].AudioDevices[j], nil
			}
		}
	}

	return AudioDevice{}, fmt.Errorf("audioDevice %q not found", id)
}

func GetDisplayGroupByID(groups []DisplayGroup, id ID) (DisplayGroup, error) {
	for i := range groups {
		if groups[i].ID == id {
			return groups[i], nil
		}
	}

	return DisplayGroup{}, fmt.Errorf("displayGroup %q not found", id)
}

func (c *Client) GetPresetByName(name string) (Preset, error) {
	for i := range c.uiConfig.Presets {
		if name == c.uiConfig.Presets[i].Name {
			return c.uiConfig.Presets[i], nil
		}
	}

	return Preset{}, fmt.Errorf("preset %q not found", name)
}
