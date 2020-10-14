package client

import "github.com/byuoitav/ui"

type client struct {
	// static info about the room
	roomID         string
	controlGroupID string

	// structs to ~do stuff~ with
	dataService  ui.DataService
	avController ui.AVController

	// info that we update occasionally
	state    ui.RoomState
	uiConfig ui.UIConfig
}
