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
	"go.uber.org/zap"
)

func (c *Client) HandleEvents() {
	mess, err := messenger.BuildMessenger(os.Getenv("HUB_ADDRESS"), base.Messenger, 1)
	if err != nil {
		fmt.Printf("%s", err)
	}

	mess.SubscribeToRooms(c.roomID)

	for {
		event := mess.ReceiveEvent()
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
			newstate, changed = handleVolume(state, event.Value, event.TargetDevice.DeviceID)
		case "muted":
			newstate, changed = handleMuted(state, event.Value, event.TargetDevice.DeviceID)
		case "power":
			newstate, changed = handlePower(state, event.Value, event.TargetDevice.DeviceID)
		case "input":
			newstate, changed = handleInput(state, event.Value, event.TargetDevice.DeviceID)
		case "blanked":
			newstate, changed = handleBlanked(state, event.Value, event.TargetDevice.DeviceID)
		default:
			continue
		}

		if !changed {
			continue
		}

		fmt.Println("\n A new room state is being sent!")
		c.state = newstate
		msg, err := JSONMessage("room", c.state)
		if err != nil {
			c.Warn("failed to create JSON message with new room state", zap.Error(err))
		}

		fmt.Printf("Sending new room down the pipeline: %v", c.state)
		c.Out <- msg

	}

	// TODO close messenger when it's done!
}

func handleVolume(state structs.PublicRoom, volume, targetDevice string) (structs.PublicRoom, bool) {
	intVolume, _ := strconv.Atoi(volume)
	deviceId := getDeviceId(targetDevice)

	for i, device := range state.AudioDevices {
		if device.Name != deviceId {
			continue
		}

		if device.Volume == &intVolume {
			continue
		}

		newState := state
		newState.AudioDevices[i].Volume = &intVolume
		return newState, true
	}

	return state, false
}

func handleMuted(state structs.PublicRoom, muted, targetDevice string) (structs.PublicRoom, bool) {
	isMuted, _ := strconv.ParseBool(muted)
	deviceId := getDeviceId(targetDevice)

	for i, device := range state.AudioDevices {
		if device.Name != deviceId {
			continue
		}

		if device.Muted == &isMuted {
			continue
		}

		newState := state
		newState.AudioDevices[i].Muted = &isMuted
		return newState, true
	}

	return state, false
}
func handlePower(state structs.PublicRoom, power, targetDevice string) (structs.PublicRoom, bool) {
	deviceId := getDeviceId(targetDevice)

	for i, device := range state.Displays {
		if device.Name != deviceId {
			continue
		}

		if device.Power == power {
			continue
		}

		newState := state
		newState.Displays[i].Power = power
		return newState, true
	}

	return state, false
}
func handleInput(state structs.PublicRoom, input, targetDevice string) (structs.PublicRoom, bool) {
	deviceId := getDeviceId(targetDevice)

	for i, device := range state.Displays {
		if device.Name != deviceId {
			continue
		}

		if device.Input == input {
			continue
		}

		newState := state
		newState.Displays[i].Input = input
		return newState, true
	}

	return state, false
}
func handleBlanked(state structs.PublicRoom, blanked, targetDevice string) (structs.PublicRoom, bool) {
	isBlanked, _ := strconv.ParseBool(blanked)
	deviceId := getDeviceId(targetDevice)

	for i, device := range state.Displays {
		if device.Name != deviceId {
			continue
		}

		if device.Blanked == &isBlanked {
			continue
		}

		newState := state
		newState.Displays[i].Blanked = &isBlanked
		return newState, true
	}

	return state, false
}

func getDeviceId(targetDevice string) string {
	splitDevice := strings.Split(targetDevice, "-")
	return splitDevice[2]
}
