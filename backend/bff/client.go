package bff

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/ui/log"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	buildingID             string
	roomID                 string
	selectedControlGroupID string

	room     structs.Room
	state    structs.PublicRoom
	uiConfig UIConfig

	// if this channel is closed, then all goroutines
	// spawned by the client should exit
	kill chan struct{}

	ws         *websocket.Conn
	httpClient *http.Client

	// messages going out to the client
	Out chan Message

	// events put in this channel get sent to the hub
	SendEvent chan events.Event

	*zap.Logger
}

func RegisterClient(ctx context.Context, ws *websocket.Conn, roomID, controlGroupID string) (*Client, error) {
	log.P.Info("Registering client", zap.String("roomID", roomID), zap.String("controlGroupID", controlGroupID), zap.String("name", ws.RemoteAddr().String()))

	split := strings.Split(roomID, "-")
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid roomID %q - must match format BLDG-ROOM", roomID)
	}

	// build our client
	c := &Client{
		buildingID: split[0],
		roomID:     roomID,
		kill:       make(chan struct{}),
		ws:         ws,
		httpClient: &http.Client{},
		Out:        make(chan Message, 1),
		SendEvent:  make(chan events.Event),
		Logger:     log.P.Named(ws.RemoteAddr().String()),
	}

	// setup shoudn't take longer than 10 seconds
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	errCh := make(chan error, 3)
	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer close(doneCh)
		wg.Wait()
	}()

	go func() {
		var err error
		defer wg.Done()

		c.state, err = GetRoomState(ctx, c.httpClient, c.roomID)
		if err != nil {
			c.Warn("unable to get room state", zap.Error(err))
			errCh <- fmt.Errorf("unable to get ui config: %v", err)
		}

		c.Debug("Successfully got room state")
	}()

	go func() {
		var err error
		defer wg.Done()

		c.room, err = GetRoomConfig(ctx, c.httpClient, c.roomID)
		if err != nil {
			c.Warn("unable to get room config", zap.Error(err))
			errCh <- fmt.Errorf("unable to get room config: %v", err)
		}

		c.Debug("Successfully got room config")
	}()

	go func() {
		var err error
		defer wg.Done()

		c.uiConfig, err = GetUIConfig(ctx, c.httpClient, c.roomID)
		if err != nil {
			c.Warn("unable to get ui config", zap.Error(err))
			errCh <- fmt.Errorf("unable to get ui config: %v", err)
		}

		c.Debug("Successfully got ui config")
	}()

	select {
	case err := <-errCh:
		return nil, fmt.Errorf("unable to get room information: %v", err)
	case <-ctx.Done():
		return nil, fmt.Errorf("unable to get room information: all requests timed out")
	case <-doneCh:
	}

	room := c.GetRoom()
	if _, ok := room.ControlGroups[controlGroupID]; ok {
		c.selectedControlGroupID = controlGroupID
	}

	//check if controlgroup is empty, if not turn on displays in controlgroup
	if len(c.selectedControlGroupID) > 0 {
		var displays []ID
		for _, display := range room.ControlGroups[c.selectedControlGroupID].Displays {
			displays = append(displays, display.ID)
		}

		setPowerMessage := SetPowerMessage{
			Displays: displays,
			Status:   "on",
		}

		// turn the control group on - this will send the room to the client
		err := c.CurrentPreset().Actions.SetPower.DoWithMessage(ctx, c, setPowerMessage)
		if err != nil {
			return nil, fmt.Errorf("unable to power on the following displays %s: %s", displays, err)
		}
	} else {
		// write the inital room info
		msg, err := JSONMessage("room", c.GetRoom())
		if err != nil {
			return nil, fmt.Errorf("unable to marshal room: %s", err)
		}

		c.Out <- msg
	}

	c.Info("Got all initial information, sent room to client. Starting ws/event goroutines")

	go c.handleEvents()
	go c.readPump()
	go c.writePump()

	return c, nil
}

func (c *Client) CurrentPreset() Preset {
	for _, p := range c.uiConfig.Presets {
		if p.Name == c.selectedControlGroupID {
			return p
		}
	}

	return Preset{}
}
