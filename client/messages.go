package client

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type message map[string]json.RawMessage

func (c *client) HandleMessage(b []byte) {
	var msg message
	if err := json.Unmarshal(b, &msg); err != nil {
		// TODO log error/return error?
		return
	}

	for k, v := range msg {
		fmt.Printf("%v\n", v)
		switch k {
		case "setInput":
		case "setMuted":
		case "setVolume":
		case "setPower":
			c.setPower(v)
		case "setBlanked":
		case "helpRequest":
		case "setSharing":
		case "selectControlGroup":
		default:
			// c.Warn("received message with unknown key", zap.String("key", k), zap.ByteString("val", v))
			// c.Out <- ErrorMessage(fmt.Errorf("unknown key %q", k))
		}
	}
}

func (c *client) OutgoingMessages() chan []byte {
	// TODO this should probably return a 'copy' of this channel...
	return c.outgoing
}

func (c *client) sendMessage(msg message) {
	b, err := json.Marshal(msg)
	if err != nil {
		// TODO log error
		return
	}

	select {
	case c.outgoing <- b:
	default:
	}
}

func (c *client) sendStringMessage(key string, format string, a ...interface{}) {
	m := make(map[string]json.RawMessage)
	m[key] = []byte(strconv.Quote(fmt.Sprintf(format, a...)))

	c.sendMessage(m)
}

func (c *client) sendJSONMsg(k string, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		// TODO log error
		return
	}

	m := make(map[string]json.RawMessage)
	m[k] = b

	c.sendMessage(m)
}
