package client

import (
	"encoding/json"
)

func (c *client) setPower(data []byte) {
	var msg struct {
		ControlGroup string `json:"controlGroup"`
		On           bool   `json:"on"`
		All          bool   `json:"all"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		// TODO log/send error
		return
	}

	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	//cgID := c.controlGroupID
	//if msg.ControlGroup != "" {
	//	cgID = msg.ControlGroup
	//}

	//// make sure control group exists
	//cg, ok := c.config.ControlGroups[cgID]
	//if !ok {
	//	// TODO log/send invalid control group error
	//	return
	//}

	if !msg.On {
		// c.doControlSet(cg.PowerOff)
		return
	}
}
