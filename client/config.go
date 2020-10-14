package client

import (
	"context"
	"fmt"

	"github.com/byuoitav/ui"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	DataService  ui.DataService
	AVController ui.AVController
}

func (c *Config) New(ctx context.Context, room, controlGroup string) (ui.Client, error) {
	client := &client{
		// TODO validate roomID?
		roomID:         room,
		controlGroupID: controlGroup,
	}

	// get initial state
	errg, gctx := errgroup.WithContext(ctx)

	// get the config
	errg.Go(func() error {
		return nil
	})

	// get the room state
	errg.Go(func() error {
		return nil
	})

	if err := errg.Wait(); err != nil {
		return nil, fmt.Errorf("unable to get data for client: %w", err)
	}

	return client, nil
}
