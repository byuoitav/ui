package client

import (
	"encoding/json"
)

func (c *client) HandleMessage(msg []byte) {
}

func (c *client) OutgoingMessages() chan []byte {
	// TODO this should probably return a 'copy' of this channel...
	return c.outgoing
}

func (c *client) sendJSONMsg(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		// TODO log error
	}

	select {
	case c.outgoing <- b:
	default:
	}
}
