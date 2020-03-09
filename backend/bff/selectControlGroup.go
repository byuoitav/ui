package bff

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type SelectControlGroup struct {
}

type SelectControlGroupMessage struct {
	ID ID `json:"id"`
}

func (s SelectControlGroupMessage) Do(c *Client, data []byte) {
	var msg SelectControlGroupMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.Warn("invalid value for selectControlGroup", zap.Error(err))
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for selectControlGroup: %s", err))
		return
	}
}
