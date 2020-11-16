package ui

import (
	"context"

	avcontrol "github.com/byuoitav/av-control-api"
)

type AVController interface {
	RoomState(context.Context, string) (avcontrol.StateResponse, error)
	SetRoomState(context.Context, string, avcontrol.StateRequest) (avcontrol.StateResponse, error)
}
