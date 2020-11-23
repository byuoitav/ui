package client

import (
	"context"
	"fmt"

	"github.com/byuoitav/ui"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Builder struct {
	DataService  ui.DataService
	AVController ui.AVController
	Log          *zap.Logger
}

func (b *Builder) New(ctx context.Context, room, controlGroup string) (ui.Client, error) {
	client := &client{
		// TODO validate roomID?
		roomID:         room,
		controlGroupID: controlGroup,
		dataService:    b.DataService,
		avController:   b.AVController,
		log:            b.Log,
		outgoing:       make(chan []byte, 1),
	}

	id := ui.RequestID(ctx)
	if id != "" {
		client.log = b.Log.Named(id)
	}

	client.handlers = map[string]messageHandler{
		"setPower":    client.setPower,
		"setVolume":   client.setVolume,
		"setMute":     client.setMute,
		"setInput":    client.setInput,
		"tiltUp":      client.tiltUp,
		"tiltDown":    client.tiltDown,
		"panLeft":     client.panLeft,
		"panRight":    client.panRight,
		"panTiltStop": client.panTiltStop,
		"zoomIn":      client.zoomIn,
		"zoomOut":     client.zoomOut,
		"zoomStop":    client.zoomStop,
		"setPreset":   client.setPreset,
	}

	// get initial state
	errg, gctx := errgroup.WithContext(ctx)

	errg.Go(func() error {
		return client.updateConfig(gctx)
	})

	errg.Go(func() error {
		return client.updateRoomState(gctx)
	})

	if err := errg.Wait(); err != nil {
		return nil, fmt.Errorf("unable to get data for client: %w", err)
	}

	// TODO start update routines

	client.sendJSONMsg("room", client.Room())
	return client, nil
}
