package client

import (
	"encoding/json"
	"fmt"
)

func (c *client) HandleMessage(msg []byte) {
	b, _ := json.Marshal(c.state)
	fmt.Printf("state: %s\n", b)
	fmt.Printf("room: %v\n", c.Room())
}

func (c *client) OutgoingMessages() chan []byte {
	return make(chan []byte)
}
