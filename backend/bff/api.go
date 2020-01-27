package bff

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/byuoitav/common/structs"
	"go.uber.org/zap"
)

// SendAPIRequest .
func (c *Client) SendAPIRequest(ctx context.Context, room structs.PublicRoom) error {
	body, err := json.Marshal(room)
	if err != nil {
		return err
	}

	roomSplit := strings.Split(c.roomID, "-")

	url := fmt.Sprintf("http://itb-1006-cp1.byu.edu:8000/buildings/%s/rooms/%s", c.buildingID, roomSplit[1])
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	c.Debug("sending API request", zap.String("url", url), zap.ByteString("body", body))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var newState structs.PublicRoom
	err = json.Unmarshal(data, &newState)
	if err != nil {
		return err
	}

	c.updateRoom(newState)
	c.Info("Updated room, sending to client")

	roomMsg, err := JSONMessage("room", c.GetRoom())
	if err != nil {
		return err
	}

	c.Out <- roomMsg
	return nil
}

// return if there are changes (true/false)
func (c *Client) updateRoom(newRoom structs.PublicRoom) {
	for _, disp := range newRoom.Displays {
		// go find the current one
		curIdx := -1
		for i := range c.state.Displays {
			if disp.Name == c.state.Displays[i].Name {
				curIdx = i
				break
			}
		}

		if curIdx < 0 {
			c.state.Displays = append(c.state.Displays, disp)
		} else {
			// merge the display
			if len(disp.Power) > 0 {
				c.state.Displays[curIdx].Power = disp.Power
			}

			if len(disp.Input) > 0 {
				c.state.Displays[curIdx].Input = disp.Input
			}

			if disp.Blanked != nil {
				c.state.Displays[curIdx].Blanked = disp.Blanked
			}
		}
	}

	for _, ad := range newRoom.AudioDevices {
		// go find the current one
		curIdx := -1
		for i := range c.state.AudioDevices {
			if ad.Name == c.state.AudioDevices[i].Name {
				curIdx = i
				break
			}
		}

		if curIdx < 0 {
			c.state.AudioDevices = append(c.state.AudioDevices, ad)
		} else {
			// merge the display
			if ad.Volume != nil {
				c.state.AudioDevices[curIdx].Volume = ad.Volume
			}

			if ad.Muted != nil {
				c.state.AudioDevices[curIdx].Muted = ad.Muted
			}
		}
	}
}
