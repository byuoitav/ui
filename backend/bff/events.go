package bff

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	//"github.com/byuoitav/av-api/base"
	"github.com/byuoitav/central-event-system/hub/base"
	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/device-monitoring/messenger"
)

func (c *Client) HandleEvents() {
	mess, err := messenger.BuildMessenger(os.Getenv("HUB_ADDRESS"), base.Messenger, 1)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("IN THE HANDLE EVENTS FUNC")
	fmt.Printf("\n%s\n", c.roomID)
	mess.SubscribeToRooms(c.roomID)
	fmt.Printf("\n%s\n", os.Getenv("HUB_ADDRESS"))

	for {
		event := mess.ReceiveEvent()
		fmt.Printf("\nEVENT RECIEVED: %s", event)
		isCoreState := false

		for _, tag := range event.EventTags {
			if tag == "core-state" {
				isCoreState = true
			}
		}

		if !isCoreState {
			continue
		}

		state := c.state
		var changed bool
		var newstate structs.PublicRoom
		switch event.Key {
		case "volume":
			fmt.Printf("RECIEVED VOLUME EVENT")
			newstate, changed = handleVolume(state, event.Value, event.TargetDevice.DeviceID)
		case "muted":
			fmt.Printf("RECIEVED MUTED EVENT")
			newstate, changed = handleMuted(state, event.Value, event.TargetDevice.DeviceID)
		case "power":
			fmt.Printf("RECIEVED POWER EVENT")
			newstate, changed = handlePower(state, event.Value, event.TargetDevice.DeviceID)
		case "input":
			fmt.Printf("RECIEVED INPUT EVENT")
			newstate, changed = handleInput(state, event.Value, event.TargetDevice.DeviceID)
		case "blanked":
			fmt.Printf("RECIEVED BLANKED EVENT")
			newstate, changed = handleBlanked(state, event.Value, event.TargetDevice.DeviceID)
		default:
			continue
		}

		if changed {
			c.state = newstate
			msg, err := JSONMessage("room", c.state)
			if err != nil {
				// error
			}
			fmt.Println(msg)

			c.Out <- msg
		}

	}

	// TODO close messenger when it's done!
}

func handleVolume(state structs.PublicRoom, volume, targetDevice string) (structs.PublicRoom, bool) {
	intVolume, _ := strconv.Atoi(volume)
	deviceId := getDeviceId(targetDevice)

	for i, device := range state.AudioDevices {
		if device.Name == deviceId {
			if device.Volume != &intVolume {
				fmt.Println("CHANGING VOLUME STATE")
				newState := state
				newState.AudioDevices[i].Volume = &intVolume
				return newState, true
			}

			continue
		}
	}

	return state, false
}
func handleMuted(state structs.PublicRoom, muted, targetDevice string) (structs.PublicRoom, bool) {
	isMuted, _ := strconv.ParseBool(muted)
	deviceId := getDeviceId(targetDevice)

	for i, device := range state.AudioDevices {
		if device.Name == deviceId {
			if device.Muted != &isMuted {
				fmt.Println("CHANGING MUTED STATE")
				newState := state
				newState.AudioDevices[i].Muted = &isMuted
				return newState, true
			}

			continue
		}
	}

	return state, false
}
func handlePower(state structs.PublicRoom, power, targetDevice string) (structs.PublicRoom, bool) {
	deviceId := getDeviceId(targetDevice)

	for i, device := range state.Displays {
		if device.Name == deviceId {
			if device.Power != power {
				fmt.Println("CHANGING POWER STATE")
				newState := state
				newState.Displays[i].Power = power
				return newState, true
			}

			continue
		}
	}

	return state, false
}
func handleInput(state structs.PublicRoom, input, targetDevice string) (structs.PublicRoom, bool) {
	deviceId := getDeviceId(targetDevice)

	for i, device := range state.Displays {
		if device.Name == deviceId {
			if device.Input != input {
				fmt.Println("CHANGING INPUT STATE")
				newState := state
				newState.Displays[i].Input = input
				return newState, true
			}

			continue
		}
	}
	return state, false
}
func handleBlanked(state structs.PublicRoom, blanked, targetDevice string) (structs.PublicRoom, bool) {
	isBlanked, _ := strconv.ParseBool(blanked)
	deviceId := getDeviceId(targetDevice)

	for i, device := range state.Displays {
		if device.Name == deviceId {
			if device.Blanked != &isBlanked {
				fmt.Println("CHANGING BLANKED STATE")
				newState := state
				newState.Displays[i].Blanked = &isBlanked
				return newState, true
			}

			continue
		}
	}
	return state, false
}

func getDeviceId(targetDevice string) string {
	splitDevice := strings.Split(targetDevice, "-")
	return splitDevice[2]
}
