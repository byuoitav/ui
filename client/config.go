package client

import (
	"context"
	"fmt"

	"github.com/byuoitav/ui"
	"golang.org/x/sync/errgroup"
)

type Builder struct {
	DataService  ui.DataService
	AVController ui.AVController
}

func (b *Builder) New(ctx context.Context, room, controlGroup string) (ui.Client, error) {
	client := &client{
		// TODO validate roomID?
		roomID:         room,
		controlGroupID: controlGroup,
		dataService:    b.DataService,
		avController:   b.AVController,
		outgoing:       make(chan []byte, 1),
	}

	client.handlers = map[string]messageHandler{
		"setPower": client.setPower,
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
