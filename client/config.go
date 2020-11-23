package client

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

func (c *client) updateConfig(ctx context.Context) error {
	c.log.Debug("Updating room state")

	config, err := c.dataService.Config(ctx, c.roomID)
	if err != nil {
		c.log.Error("unable to update room config", zap.Error(err))
		return fmt.Errorf("unable to get config: %w", err)
	}

	c.log.Debug("Successfully updated room config")

	c.configMu.Lock()
	defer c.configMu.Unlock()

	c.config = config
	return nil
}
