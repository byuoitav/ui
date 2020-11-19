package couch

import (
	"context"
	"fmt"
)

type panelConfig struct {
	ID string `json:"_id"`

	ControlPanels map[string]struct {
		UIType string `json:"uiType"`

		// TODO divider sensors
		ControlGroup string `json:"controlGroup"`
	} `json:"controlPanels"`
}

func (d *DataService) panelConfig(ctx context.Context, room string) (panelConfig, error) {
	var config panelConfig

	db := d.client.DB(ctx, d.database)
	if err := db.Get(ctx, room).ScanDoc(&config); err != nil {
		return config, fmt.Errorf("unable to get/scan room: %w", err)
	}

	return config, nil
}

func (d *DataService) UIForDevice(ctx context.Context, room, id string) (string, error) {
	config, err := d.panelConfig(ctx, room)
	if err != nil {
		return "", fmt.Errorf("unable to get panel config: %w", err)
	}

	v, ok := config.ControlPanels[id]
	if !ok {
		return "", fmt.Errorf("%s is not configured in %s", id, room)
	}

	return v.UIType, nil
}

func (d *DataService) ControlGroup(ctx context.Context, room, id string) (string, error) {
	config, err := d.panelConfig(ctx, room)
	if err != nil {
		return "", fmt.Errorf("unable to get panel config: %w", err)
	}

	v, ok := config.ControlPanels[id]
	if !ok {
		return "", fmt.Errorf("%s is not configured in %s", id, room)
	}

	return v.ControlGroup, nil
}

func (d *DataService) RoomAndControlGroup(ctx context.Context, key string) (string, string, error) {
	return "", "", nil
}
