package client

import (
	"context"
	"fmt"
)

func (c *client) updateConfig(ctx context.Context) error {
	config, err := c.dataService.Config(ctx, c.roomID)
	if err != nil {
		return fmt.Errorf("unable to get config: %w", err)
	}

	c.config = config
	return nil
}
