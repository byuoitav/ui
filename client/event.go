package client

import (
	"encoding/json"
	"fmt"
)

func (c *client) event(data []byte) {
	var msg struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Printf("error: %s\n", err)
		// TODO log/send error
		return
	}

	fmt.Printf("event: %v\n", msg)

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
}
