package couch

import "context"

func (d *DataService) UIForDevice(ctx context.Context, room, id string) (string, error) {
	return "", nil
}

func (d *DataService) ControlGroup(ctx context.Context, room, id string) (string, error) {
	return "ITB 1006", nil
}

func (d *DataService) RoomAndControlGroup(ctx context.Context, key string) (string, string, error) {
	return "", "", nil
}
