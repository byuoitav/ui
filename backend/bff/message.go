package bff

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Message map[string]json.RawMessage

type SetInputMessage struct {
	DisplayID string `json:"display"`
	InputID   string `json:"input"`
}

func ErrorMessage(format string, a ...interface{}) Message {
	return StringMessage("error", format, a...)
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
			c.CurrentPreset().Actions.SetInput.Do(c, v)
		case "setMuted":
		case "setVolume":
		default:
			// c.Warn("received message with unknown key", zap.String("key", k), zap.ByteString("val", v))
			fmt.Printf("v: %s", v)
			c.Out <- ErrorMessage("unknown key %q", k)
		}
	}
}
