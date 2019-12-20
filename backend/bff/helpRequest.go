package bff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/byuoitav/common/v2/events"
	"go.uber.org/zap"
)

type HelpRequest struct {
}

type HelpRequestMessage struct {
	Message string `json:"msg"`
}

func (hr HelpRequest) Do(c *Client, data []byte) {
	var msg HelpRequestMessage

	if err := json.Unmarshal(data, &msg); err != nil {
		c.Out <- ErrorMessage(fmt.Errorf("invalid value for helpRequest: %s", err))
		return
	}

	c.Info("Sending help request", zap.String("helpMsg", msg.Message))

	c.SendEvent <- events.Event{
		GeneratingSystem: os.Getenv("SYSTEM_ID"),
		User:             c.id,
		Key:              "help-request",
		Value:            "confirm",
		Data:             msg.Message,
		Timestamp:        time.Now(),
		EventTags: []string{
			events.Alert,
		},
		AffectedRoom: events.BasicRoomInfo{
			BuildingID: c.buildingID,
			RoomID:     c.roomID,
		},
		TargetDevice: events.BasicDeviceInfo{
			BasicRoomInfo: events.BasicRoomInfo{
				BuildingID: c.buildingID,
				RoomID:     c.roomID,
			},
			DeviceID: c.roomID + "-" + c.selectedControlGroupID,
		},
	}
}
