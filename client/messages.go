package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

type message map[string]json.RawMessage

type messageHandler func(b []byte)

func (c *client) HandleMessage(b []byte) {
	var msg message
	if err := json.Unmarshal(b, &msg); err != nil {
		c.log.Warn("unable to parse message", zap.Error(err), zap.ByteString("msg", b))
		return
	}

	for k, v := range msg {
		if handler, ok := c.handlers[k]; ok {
			c.log.Debug("Calling handler for message", zap.String("key", k), zap.ByteString("val", v))
			handler(v)
		} else {
			c.log.Warn("no handler registered", zap.String("key", k), zap.ByteString("val", v))
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
		c.log.Warn("unable to send message", zap.Error(err))
		return
	}

	select {
	case c.outgoing <- b:
	default:
		c.log.Warn("outgoing channel was full; not sending message to client")
	}
}

func (c *client) sendStringMessage(key string, format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	c.log.Debug("Sending string message", zap.String("key", key), zap.String("val", str))

	m := make(map[string]json.RawMessage)
	m[key] = []byte(strconv.Quote(str))

	c.sendMessage(m)
}

func (c *client) sendJSONMsg(k string, v interface{}) {
	c.log.Debug("Sending json message", zap.String("key", k), zap.Any("val", v))

	b, err := json.Marshal(v)
	if err != nil {
		c.log.Warn("unable to send json message", zap.Error(err))
		return
	}

	m := make(map[string]json.RawMessage)
	m[k] = b

	c.sendMessage(m)
}
