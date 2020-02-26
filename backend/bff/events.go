package bff

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/byuoitav/central-event-system/hub/base"
	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/common/v2/events"
	"go.uber.org/zap"
)

func (c *Client) handleEvents() {
	mess, err := messenger.BuildMessenger(os.Getenv("HUB_ADDRESS"), base.Messenger, 1)
	if err != nil {
		c.Error("unable to build messenger", zap.Error(err))
	}

	mess.SubscribeToRooms(c.roomID)

	// receive events
	eventCh := make(chan base.EventWrapper, 1)
	mess.SetReceiveChannel(eventCh)

	defer func() {
		c.Info("Closing event messenger")
		mess.UnsubscribeFromRooms(c.roomID)
		mess.Kill()
		close(eventCh)
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	// send events
	go func() {
		defer wg.Done()

		for {
			select {
			case event := <-c.SendEvent:
				mess.SendEvent(event)
			case <-c.kill:
				return
			}
		}
	}()

	go func() {
		defer wg.Done()

		for {
			select {
			case eventWrap := <-eventCh:
				var event events.Event
				if err := json.Unmarshal(eventWrap.Event, &event); err != nil {
					c.Warn("received an invalid event", zap.Error(err))
					continue
				}

				isCoreState := false

				for _, tag := range event.EventTags {
					if tag == "core-state" {
						isCoreState = true
					}
				}

				if !isCoreState {
					continue
				}

				c.Debug("Received core state event", zap.String("key", event.Key), zap.String("value", event.Value), zap.String("on", event.TargetDevice.DeviceID))

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
					c.Debug("Ignoring no change event")
					continue
				}

				c.Info("Updating room from event", zap.String("key", event.Key), zap.String("value", event.Value), zap.String("on", event.TargetDevice.DeviceID))
				c.state = newstate

				// send an updated room to the client
				msg, err := JSONMessage("room", c.GetRoom())
				if err != nil {
					c.Warn("failed to create JSON message with new room state", zap.Error(err))
					continue
				}

				c.Out <- msg
			case <-c.kill:
				return
			}
		}
	}()

	wg.Wait()
}

func handleVolume(state structs.PublicRoom, volume, targetDevice string) (structs.PublicRoom, bool) {
	intVolume, _ := strconv.Atoi(volume)
	deviceID := getDeviceID(targetDevice)
	changed := false

	for i := range state.AudioDevices {
		if state.AudioDevices[i].Name != deviceID {
			continue
		}

		if state.AudioDevices[i].Volume == &intVolume {
			continue
		}

		state.AudioDevices[i].Volume = &intVolume
		changed = true
	}

	return state, changed
}

func handleMuted(state structs.PublicRoom, muted, targetDevice string) (structs.PublicRoom, bool) {
	isMuted, _ := strconv.ParseBool(muted)
	deviceID := getDeviceID(targetDevice)
	changed := false

	for i := range state.AudioDevices {
		if state.AudioDevices[i].Name != deviceID {
			continue
		}

		if state.AudioDevices[i].Muted == &isMuted {
			continue
		}

		state.AudioDevices[i].Muted = &isMuted
		changed = true
	}

	return state, changed
}
func handlePower(state structs.PublicRoom, power, targetDevice string) (structs.PublicRoom, bool) {
	deviceID := getDeviceID(targetDevice)
	changed := false

	for i := range state.Displays {
		if state.Displays[i].Name != deviceID {
			continue
		}

		if state.Displays[i].Power == power {
			continue
		}

		state.Displays[i].Power = power
		changed = true
	}

	return state, changed
}
func handleInput(state structs.PublicRoom, input, targetDevice string) (structs.PublicRoom, bool) {
	deviceID := getDeviceID(targetDevice)
	changed := false

	for i := range state.Displays {
		if state.Displays[i].Name != deviceID {
			continue
		}

		if state.Displays[i].Input == input {
			continue
		}

		state.Displays[i].Input = input
		changed = true
	}

	return state, changed
}
func handleBlanked(state structs.PublicRoom, blanked, targetDevice string) (structs.PublicRoom, bool) {
	isBlanked, _ := strconv.ParseBool(blanked)
	deviceID := getDeviceID(targetDevice)
	changed := false

	for i := range state.Displays {
		if state.Displays[i].Name != deviceID {
			continue
		}

		if state.Displays[i].Blanked == &isBlanked {
			continue
		}

		state.Displays[i].Blanked = &isBlanked
		changed = true
	}

	return state, changed
}

func getDeviceID(targetDevice string) string {
	splitDevice := strings.Split(targetDevice, "-")
	return splitDevice[2]
}
