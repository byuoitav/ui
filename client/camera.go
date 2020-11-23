package client

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

type cameraControlMsg struct {
	Camera string `json:"camera"`
}

func (c *client) tiltUp(data []byte) {
	var msg cameraControlMsg
	if err := json.Unmarshal(data, &msg); err != nil {
		c.log.Warn("unable to unmarshal message", zap.Error(err), zap.ByteString("data", data))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log := c.log.With(zap.String("controlGroup", c.controlGroupID))

	// make sure control group exists
	cg, ok := c.config.ControlGroups[c.controlGroupID]
	if !ok {
		log.Warn("invalid control group")
		return
	}

	for _, cam := range cg.Cameras {
		if cam.Name == msg.Camera {
			log.Info("Tilting up", zap.String("camera", msg.Camera))
			c.doControlSet(ctx, cam.TiltUp)
			return
		}
	}

	c.log.Warn("invalid camera", zap.String("camera", msg.Camera))
}

func (c *client) tiltDown(data []byte) {
	var msg cameraControlMsg
	if err := json.Unmarshal(data, &msg); err != nil {
		c.log.Warn("unable to unmarshal message", zap.Error(err), zap.ByteString("data", data))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log := c.log.With(zap.String("controlGroup", c.controlGroupID))

	// make sure control group exists
	cg, ok := c.config.ControlGroups[c.controlGroupID]
	if !ok {
		log.Warn("invalid control group")
		return
	}

	for _, cam := range cg.Cameras {
		if cam.Name == msg.Camera {
			log.Info("Tilting down", zap.String("camera", msg.Camera))
			c.doControlSet(ctx, cam.TiltDown)
			return
		}
	}

	c.log.Warn("invalid camera", zap.String("camera", msg.Camera))
}

func (c *client) panLeft(data []byte) {
	var msg cameraControlMsg
	if err := json.Unmarshal(data, &msg); err != nil {
		c.log.Warn("unable to unmarshal message", zap.Error(err), zap.ByteString("data", data))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log := c.log.With(zap.String("controlGroup", c.controlGroupID))

	// make sure control group exists
	cg, ok := c.config.ControlGroups[c.controlGroupID]
	if !ok {
		log.Warn("invalid control group")
		return
	}

	for _, cam := range cg.Cameras {
		if cam.Name == msg.Camera {
			log.Info("Panning left", zap.String("camera", msg.Camera))
			c.doControlSet(ctx, cam.PanLeft)
			return
		}
	}

	c.log.Warn("invalid camera", zap.String("camera", msg.Camera))
}

func (c *client) panRight(data []byte) {
	var msg cameraControlMsg
	if err := json.Unmarshal(data, &msg); err != nil {
		c.log.Warn("unable to unmarshal message", zap.Error(err), zap.ByteString("data", data))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log := c.log.With(zap.String("controlGroup", c.controlGroupID))

	// make sure control group exists
	cg, ok := c.config.ControlGroups[c.controlGroupID]
	if !ok {
		log.Warn("invalid control group")
		return
	}

	for _, cam := range cg.Cameras {
		if cam.Name == msg.Camera {
			log.Info("Panning right", zap.String("camera", msg.Camera))
			c.doControlSet(ctx, cam.PanRight)
			return
		}
	}

	c.log.Warn("invalid camera", zap.String("camera", msg.Camera))
}

func (c *client) panTiltStop(data []byte) {
	var msg cameraControlMsg
	if err := json.Unmarshal(data, &msg); err != nil {
		c.log.Warn("unable to unmarshal message", zap.Error(err), zap.ByteString("data", data))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log := c.log.With(zap.String("controlGroup", c.controlGroupID))

	// make sure control group exists
	cg, ok := c.config.ControlGroups[c.controlGroupID]
	if !ok {
		log.Warn("invalid control group")
		return
	}

	for _, cam := range cg.Cameras {
		if cam.Name == msg.Camera {
			log.Info("Stopping pan/tilt", zap.String("camera", msg.Camera))
			c.doControlSet(ctx, cam.PanTiltStop)
			return
		}
	}

	c.log.Warn("invalid camera", zap.String("camera", msg.Camera))
}

func (c *client) zoomIn(data []byte) {
	var msg cameraControlMsg
	if err := json.Unmarshal(data, &msg); err != nil {
		c.log.Warn("unable to unmarshal message", zap.Error(err), zap.ByteString("data", data))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log := c.log.With(zap.String("controlGroup", c.controlGroupID))

	// make sure control group exists
	cg, ok := c.config.ControlGroups[c.controlGroupID]
	if !ok {
		log.Warn("invalid control group")
		return
	}

	for _, cam := range cg.Cameras {
		if cam.Name == msg.Camera {
			log.Info("Zooming in", zap.String("camera", msg.Camera))
			c.doControlSet(ctx, cam.ZoomIn)
			return
		}
	}

	c.log.Warn("invalid camera", zap.String("camera", msg.Camera))
}

func (c *client) zoomOut(data []byte) {
	var msg cameraControlMsg
	if err := json.Unmarshal(data, &msg); err != nil {
		c.log.Warn("unable to unmarshal message", zap.Error(err), zap.ByteString("data", data))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log := c.log.With(zap.String("controlGroup", c.controlGroupID))

	// make sure control group exists
	cg, ok := c.config.ControlGroups[c.controlGroupID]
	if !ok {
		log.Warn("invalid control group")
		return
	}

	for _, cam := range cg.Cameras {
		if cam.Name == msg.Camera {
			log.Info("Zooming out", zap.String("camera", msg.Camera))
			c.doControlSet(ctx, cam.ZoomOut)
			return
		}
	}

	c.log.Warn("invalid camera", zap.String("camera", msg.Camera))
}

func (c *client) zoomStop(data []byte) {
	var msg cameraControlMsg
	if err := json.Unmarshal(data, &msg); err != nil {
		c.log.Warn("unable to unmarshal message", zap.Error(err), zap.ByteString("data", data))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log := c.log.With(zap.String("controlGroup", c.controlGroupID))

	// make sure control group exists
	cg, ok := c.config.ControlGroups[c.controlGroupID]
	if !ok {
		log.Warn("invalid control group")
		return
	}

	for _, cam := range cg.Cameras {
		if cam.Name == msg.Camera {
			log.Info("Stopping zoom", zap.String("camera", msg.Camera))
			c.doControlSet(ctx, cam.ZoomStop)
			return
		}
	}

	c.log.Warn("invalid camera", zap.String("camera", msg.Camera))
}

func (c *client) setPreset(data []byte) {
	var msg struct {
		cameraControlMsg
		Preset string `json:"preset"`
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		c.log.Warn("unable to unmarshal message", zap.Error(err), zap.ByteString("data", data))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log := c.log.With(zap.String("controlGroup", c.controlGroupID))

	// make sure control group exists
	cg, ok := c.config.ControlGroups[c.controlGroupID]
	if !ok {
		log.Warn("invalid control group")
		return
	}

	for _, cam := range cg.Cameras {
		if cam.Name == msg.Camera {
			for _, preset := range cam.Presets {
				if preset.Name == msg.Preset {
					log.Info("Setting preset", zap.String("camera", msg.Camera), zap.String("preset", msg.Preset))
					c.doControlSet(ctx, preset.SetPreset)
					return
				}
			}

			c.log.Warn("invalid preset", zap.String("camera", msg.Camera), zap.String("preset", msg.Preset))
			return
		}
	}

	c.log.Warn("invalid camera", zap.String("camera", msg.Camera))
}
