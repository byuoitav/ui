package bff

import (
	"os"

	"github.com/byuoitav/av-api/base"
	"github.com/byuoitav/device-monitoring/messenger"
)

func (c *Client) HandleEvents() {
	mess, err := messenger.BuildMessenger(os.Getenv("HUB_ADDRESS"), base.Messenger, 1)
	if err != nil {
	}

	mess.SubscribeToRooms(c.roomID)

	for {
		event := mess.ReceiveEvent()
		// TODO do something with event
	}

	// TODO close messenger when it's done!
}
