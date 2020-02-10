package bff

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type GetKeyConfig struct {
}

type GetKeyMessage struct {
	ControlGroupID string `json:"controlGroupID"`
}

type KeyResponse struct {
	ControlKey string `json:"controlKey"`
	ControlURL string `json:"controlURL"`
}

func (ck GetKeyConfig) Do(c *Client, data []byte) {
	if strings.HasPrefix(c.ws.RemoteAddr().String(), "127.0.0.1") || strings.HasPrefix(c.ws.RemoteAddr().String(), "localhost") {
		var msg GetKeyMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			c.Out <- ErrorMessage(fmt.Errorf("invalid control group id: %s", err))
			return
		}

		hostname := os.Getenv("SYSTEM_ID")
		if len(hostname) == 0 {
			c.Out <- ErrorMessage(fmt.Errorf("cannot get hostname to find control key"))
			return
		}

		hostArray := strings.Split(hostname, "-")

		url := fmt.Sprintf("%s/%s-%s %s/getControlKey", os.Getenv("CODE_SERVICE_URL"), hostArray[0], hostArray[1], msg.ControlGroupID)
		resp, err := c.httpClient.Get(url)
		if err != nil {
			c.Out <- ErrorMessage(fmt.Errorf("failed to get control key for the group %s: %s", msg.ControlGroupID, err))
			return
		}

		var k KeyResponse
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.Out <- ErrorMessage(fmt.Errorf("failed to read response body: %s", err))
			return
		}

		if err := json.Unmarshal(body, &k); err != nil {
			c.Out <- ErrorMessage(fmt.Errorf("invalid control key: %s", err))
			return
		}

		if c.room.Designation != "production" {
			k.ControlURL = "rooms.stg.byu.edu"
		} else {
			k.ControlURL = "rooms.byu.edu"
		}

		j, err := JSONMessage("mobileControl", k)
		if err != nil {
			c.Out <- ErrorMessage(fmt.Errorf("failed to format control key: %s", err))
			return
		}

		c.Out <- j
	}
}
