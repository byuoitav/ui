package bff

import (
	"encoding/json"

	"go.uber.org/zap"
)

type SetVolume struct {
	AudioDeviceID string `json:"audioDevice"`
	Level         int    `json:"level"`
}

func (c *Client) HandleMessage(msg Message) chan Message {
	resps := make(chan Message)

	go func() {
		for k, v := range msg {
			switch k {
			case "setInput":
			case "setMuted":
			case "setVolume":
				var val SetVolume
				err := json.Unmarshal(v, &val)
				if err != nil {
					resps <- ErrorMessage("invalid value for key %q: %s", k, err)
					return
				}

				c.Info("setting volume", zap.String("id", val.AudioDeviceID), zap.Int("level", val.Level))
			default:
				// c.Warn("received message with unknown key", zap.String("key", k), zap.ByteString("val", v))
				resps <- ErrorMessage("unknown key %q", k)
			}
		}
	}()

	return resps
}
