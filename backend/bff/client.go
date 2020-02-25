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

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Client represents a client of the bff
type Client struct {
	buildingID             string
	roomID                 string
	selectedControlGroupID string

	room     structs.Room
	state    structs.PublicRoom
	uiConfig UIConfig

	lazs       LazaretteState
	lazUpdates chan lazarette.KeyValue

	//shareMutex sync.RWMutex
	//sharing    Sharing
	//// TODO get shareable
	//shareable Shareable

	// if this channel is closed, then all goroutines
	// spawned by the client should exit
	kill      chan struct{}
	closeOnce sync.Once

	ws         *websocket.Conn
	httpClient *http.Client

	// messages going out to the client
	Out chan Message

	// events put in this channel get sent to the hub
	SendEvent chan events.Event

	*zap.Logger
}

// RegisterClient registers a new client
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
		// shareMutex: sync.RWMutex{},
		lazUpdates: make(chan lazarette.KeyValue),
		lazs: LazaretteState{
			Map: &sync.Map{},
		},
	}

	// setup shoudn't take longer than 10 seconds
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// create the errgroup for all these setup functions
	g, ctx := errgroup.WithContext(ctx)

	// get the room state
	g.Go(func() error {
		var err error
		c.state, err = GetRoomState(ctx, c.httpClient, c.roomID)
		if err != nil {
			return fmt.Errorf("unable to get ui config: %w", err)
		}

		c.Debug("Successfully got room state")
		return nil
	})

	// get the room config
	g.Go(func() error {
		var err error

		c.room, err = GetRoomConfig(ctx, c.httpClient, c.roomID)
		if err != nil {
			return fmt.Errorf("unable to get room config: %w", err)
		}

		c.Debug("Successfully got room config")
		return nil
	})

	// get the ui config
	g.Go(func() error {
		var err error

		c.uiConfig, err = GetUIConfig(ctx, c.httpClient, c.roomID)
		if err != nil {
			return fmt.Errorf("unable to get ui config: %w", err)
		}

		c.Debug("Successfully got ui config")
		return nil
	})

	// connect to lazarette
	g.Go(func() error {
		var err error

		laz, err := ConnectToLazarette(ctx)
		if err != nil {
			return fmt.Errorf("unable to connect to lazarette: %w", err)
		}

		// create the subscription
		sub, err := laz.Subscribe(ctx, &lazarette.Key{Key: roomID})
		if err != nil {
			return fmt.Errorf("unable to subscribe to lazarette: %w", err)
		}

		go c.syncLazaretteState(sub)
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("unable to get setup client: %w", err)
	}

	room := c.GetRoom()
	if _, ok := room.ControlGroups[controlGroupID]; ok {
		c.selectedControlGroupID = controlGroupID
	}

	// TODO this should happen in get room
	// Set up who can share to who
	//s := make(map[ID][]ID)
	//for _, p := range c.uiConfig.Presets {
	//	for i, d := range p.ShareableDisplays {
	//		shar := remove(p.ShareableDisplays, i)
	//		shareable := make([]ID, len(shar))
	//		for j := range shar {
	//			shareable[j] = ID(shar[j])
	//		}
	//		s[ID(d)] = shareable
	//	}
	//}
	//c.shareable = s

	// send the inital room info
	msg, err := JSONMessage("room", c.GetRoom())
	if err != nil {
		return nil, fmt.Errorf("unable to marshal room: %s", err)
	}

	c.Out <- msg

	c.Info("Got all initial information, sent room to client. Starting ws/event goroutines")

	go c.handleEvents()
	go c.readPump()
	go c.writePump()
	return c, nil
}

func remove(l []string, index int) []string {
	l[index] = l[len(l)-1]
	return l[:len(l)-1]
}

// Wait waits until a client is dead
func (c *Client) Wait() {
	<-c.kill
}

// Close closes the client
func (c *Client) Close() {
	c.closeOnce.Do(func() {
		c.Info("Closing client. Bye!")

		// close the kill chan to clean up all resources
		close(c.kill)

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
