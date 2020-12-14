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
	// TODO this function will have to get the room, device -> call ControlGroup
	RoomAndControlGroup(ctx context.Context, key string) (string, string, error)

	// Config returns the config for the given room.
	Config(ctx context.Context, room string) (Config, error)
}

// Event - i'm holding off on events for a few days
type Event struct {
	Tags   []string
	Room   string
	Device string

	IP string

	Key   string
	Value string
}

type EventPublisher interface {
	Publish(context.Context, Event) error
}
