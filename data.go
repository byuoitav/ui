package ui

import "context"

type DataService interface {
	// UIForDevice returns the ui that should be used for the given room/id.
	UIForDevice(ctx context.Context, room, id string) (string, error)

	// ControlGroup takes the room and id of a device, and returns the control group that that devices should use.
	// TODO this should also take into account the divider sensors?
	// or maybe that should be some combo interface that takes this and key service?
	ControlGroup(ctx context.Context, room, id string) (string, error)

	// RoomAndControlGroup figures out which room and control group is associated with the given key.
	RoomAndControlGroup(ctx context.Context, key string) (string, string, error)

	// Config returns the config for the given room.
	Config(ctx context.Context, room string) (Config, error)
}
