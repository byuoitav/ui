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
	"github.com/byuoitav/lazarette/lazarette"
	"github.com/byuoitav/ui/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// ClientConfig represents some configuration options for the client
type ClientConfig struct {
	RoomID         string
	ControlGroupID string

	AvApiAddr         string
	CodeServiceAddr   string
	RemoteControlAddr string
	LazaretteAddr     string
	LazaretteSSL      bool
}

// Client represents a client of the bff
type Client struct {
	config ClientConfig
	stats  ClientStats

	buildingID             string
	roomID                 string
	selectedControlGroupID string

	room     structs.Room
	state    structs.PublicRoom
	uiConfig UIConfig

	lazs       LazaretteState
	lazConn    *grpc.ClientConn
	lazUpdates chan lazMessage

	// controlKeys is periodically updated
	controlKeysMu sync.RWMutex
	controlKeys   map[string]string

	// if this channel is closed, then all goroutines
	// spawned by the client should exit
	kill      chan struct{}
	closeOnce sync.Once
	killed    bool

	ws         *websocket.Conn
	httpClient *http.Client

	// messages going out to the client
	Out chan Message

	// events put in this channel get sent to the hub
	SendEvent chan events.Event

	*zap.Logger
}

// RegisterClient registers a new client
func RegisterClient(ctx context.Context, ws *websocket.Conn, config ClientConfig) (*Client, error) {
	split := strings.Split(config.RoomID, "-")
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid roomID %q - must match format BLDG-ROOM", config.RoomID)
	}

	// build our client
	c := &Client{
		config:     config,
		buildingID: split[0],
		roomID:     config.RoomID,
		kill:       make(chan struct{}),
		ws:         ws,
		httpClient: &http.Client{},
		Out:        make(chan Message, 1),
		SendEvent:  make(chan events.Event),
		Logger:     log.P.Named(ws.RemoteAddr().String()),
		lazUpdates: make(chan lazMessage),
		lazs: LazaretteState{
			Map: &sync.Map{},
		},
		controlKeys: make(map[string]string),
	}

	c.Info("Registering client", zap.String("roomID", config.RoomID), zap.String("controlGroupID", config.ControlGroupID))

	// init stats
	c.stats.AvControlApi.ResponseCodes = make(map[int]uint)
	now := time.Now()
	c.stats.CreatedAt = &now

	// setup shoudn't take longer than 10 seconds
	sctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// create the errgroup for all these setup functions
	g, gctx := errgroup.WithContext(sctx)

	// get the room state
	g.Go(func() error {
		c.stats.Routines++
		defer c.stats.decRoutines()

		var err error
		c.state, err = GetRoomState(gctx, c.httpClient, c.config.AvApiAddr, c.roomID)
		if err != nil {
			return fmt.Errorf("unable to get room state: %w", err)
		}

		c.Debug("Successfully got room state")
		return nil
	})

	// get the room config
	g.Go(func() error {
		c.stats.Routines++
		defer c.stats.decRoutines()

		var err error
		c.room, err = GetRoomConfig(gctx, c.httpClient, c.config.AvApiAddr, c.roomID)
		if err != nil {
			return fmt.Errorf("unable to get room config: %w", err)
		}

		c.Debug("Successfully got room config")
		return nil
	})

	// get the ui config
	g.Go(func() error {
		c.stats.Routines++
		defer c.stats.decRoutines()

		var err error
		c.uiConfig, err = GetUIConfig(gctx, c.httpClient, c.roomID)
		if err != nil {
			return fmt.Errorf("unable to get ui config: %w", err)
		}

		c.Debug("Successfully got ui config")
		return nil
	})

	// connect to lazarette
	g.Go(func() error {
		c.stats.Routines++
		defer c.stats.decRoutines()

		var err error
		c.lazConn, err = createGrpcConn(gctx, c.config.LazaretteAddr, c.config.LazaretteSSL)
		if err != nil {
			return fmt.Errorf("unable to connect to lazarette: %w", err)
		}

		c.Debug("Successfully connected to lazarette")
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("unable to setup client: %w", err)
	}

	// start lazarette subscription
	laz := lazarette.NewLazaretteClient(c.lazConn)
	sub, err := laz.Subscribe(ctx, &lazarette.Key{Key: c.roomID})
	if err != nil {
		return nil, fmt.Errorf("unable to subscribe to lazarette: %w", err)
	}

	// build the initial room
	room := c.GetRoom()
	if _, ok := room.ControlGroups[config.ControlGroupID]; ok {
		c.selectedControlGroupID = config.ControlGroupID
	}

	// send the inital room info
	msg, err := JSONMessage("room", c.GetRoom())
	if err != nil {
		return nil, fmt.Errorf("unable to marshal room: %s", err)
	}

	c.Out <- msg

	c.Info("Got all initial information, sent room to client. Starting ws/event goroutines")

	// start data update routines
	go func() {
		c.stats.Routines++
		defer c.stats.decRoutines()

		c.subLazaretteState(sub)
	}()

	go func() {
		c.stats.Routines++
		defer c.stats.decRoutines()

		c.updateLazaretteState(laz)
	}()

	go func() {
		c.stats.Routines++
		defer c.stats.decRoutines()

		c.updateControlKey()
	}()

	go func() {
		c.stats.Routines++
		defer c.stats.decRoutines()

		c.handleEvents()
	}()

	go func() {
		c.stats.Routines++
		defer c.stats.decRoutines()

		c.readPump()
	}()

	go func() {
		c.stats.Routines++
		defer c.stats.decRoutines()

		c.writePump()
	}()

	return c, nil
}

// Wait waits until a client is dead
func (c *Client) Wait() {
	<-c.kill
}

func (c *Client) Killed() bool {
	return c.killed
}

// Close closes the client
func (c *Client) Close() {
	c.closeOnce.Do(func() {
		c.Info("Closing client. Bye!")
		c.killed = true

		// close the kill chan to clean up all resources
		close(c.kill)

		// close the lazarette connection
		c.lazConn.Close()

		// close our websocket with the frontend
		c.ws.Close()
	})
}

// CurrentPreset returns the current preset of the room
func (c *Client) CurrentPreset() Preset {
	for _, p := range c.uiConfig.Presets {
		if p.Name == c.selectedControlGroupID {
			return p
		}
	}

	return Preset{}
}

// Refresh refreshes the client
func (c *Client) Refresh() {
	c.Info("Sending refresh message")
	c.Out <- StringMessage("refresh", "")
}
