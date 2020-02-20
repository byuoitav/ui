package bff

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// TODO all messages should take a ctx

type Message map[string]json.RawMessage

func ErrorMessage(err error) Message {
	return StringMessage("error", err.Error())
}

func StringMessage(key string, format string, a ...interface{}) Message {
	m := make(map[string]json.RawMessage)
	m[key] = []byte(strconv.Quote(fmt.Sprintf(format, a...)))
	return m
}

func JSONMessage(key string, val interface{}) (Message, error) {
	data, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	m := make(map[string]json.RawMessage)
	m[key] = data
	return m, nil
}

func (c *Client) HandleMessage(msg Message) {
	for k, v := range msg {
		switch k {
		case "setInput":
			//This will also blank and unblank the device
			c.CurrentPreset().Actions.SetInput.Do(c, v)
		case "setMuted":
			c.CurrentPreset().Actions.SetMuted.Do(c, v)
		case "setVolume":
			c.CurrentPreset().Actions.SetVolume.Do(c, v)
		case "setPower":
			c.CurrentPreset().Actions.SetPower.Do(c, v)
		case "turnOffRoom":
			_ = c.CurrentPreset().Actions.SetPower.PowerOffAll(c)
		case "helpRequest":
			c.CurrentPreset().Actions.HelpRequest.Do(c, v)
		case "getControlKey":
			c.CurrentPreset().Actions.GetControlKey.Do(c, v)
		case "setSharing":
			c.CurrentPreset().Actions.SetSharing.Do(c, v)
		default:
			// c.Warn("received message with unknown key", zap.String("key", k), zap.ByteString("val", v))
			fmt.Printf("v: %s", v)
			c.Out <- ErrorMessage(fmt.Errorf("unknown key %q", k))
		}
	}
}
