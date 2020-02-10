package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/lazarette/lazarette"
	"github.com/byuoitav/ui/log"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func remove(l []string, index int) []string {
	l[index] = l[len(l)-1]
	return l[:len(l)-1]
}

// Client represents a client of the bff
type Client struct {
	buildingID             string
	roomID                 string
	selectedControlGroupID string

	room     structs.Room
	state    structs.PublicRoom
	lazState LazState

	shareMutex *sync.RWMutex
	sharing    Sharing
	// TODO get shareable
	shareable Shareable

	uiConfig UIConfig

	// if this channel is closed, then all goroutines
	// spawned by the client should exit
	kill  chan struct{}
	close sync.Once

	ws         *websocket.Conn
	httpClient *http.Client

	// messages going out to the client
	Out chan Message

	// events put in this channel get sent to the hub
	SendEvent chan events.Event

	lazContext context.Context
	lazCancel  context.CancelFunc

	*zap.Logger
}

// RegisterClient registers a new client
func RegisterClient(ctx context.Context, ws *websocket.Conn, roomID, controlGroupID string) (*Client, error) {
	log.P.Info("Registering client", zap.String("roomID", roomID), zap.String("controlGroupID", controlGroupID), zap.String("name", ws.RemoteAddr().String()))

	split := strings.Split(roomID, "-")
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid roomID %q - must match format BLDG-ROOM", roomID)
	}
	lazAddr := os.Getenv("LAZARETTE_ADDR")
	if len(lazAddr) == 0 {
		return nil, fmt.Errorf("LAZARETTE_ADDR not set")
	}

	lazContext, lazCancel := context.WithCancel(ctx)

	conn, err := grpc.DialContext(lazContext, lazAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		lazCancel()
		return nil, fmt.Errorf("unable to connect with grpc to lazarette %v", err)
	}

	remote := lazarette.NewLazaretteClient(conn)

	lazState := LazState{
		Client: remote,
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
		lazState:   lazState,
		Logger:     log.P.Named(ws.RemoteAddr().String()),
		shareMutex: new(sync.Mutex),
	}

	c.lazContext = lazContext
	c.lazCancel = lazCancel

	// setup shoudn't take longer than 10 seconds
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	errCh := make(chan error, 3)
	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer close(doneCh)
		wg.Wait()
	}()
	// 1
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

	// 2
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

	// 3
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

	// 4 - LazState setup
	go func() {
		c.Info("LazState: entered")
		defer wg.Done()
		var err error
		setSharing := true

		kSharing := &lazarette.Key{
			Key: fmt.Sprintf("%v-_sharing_displays", roomID),
		}
		jSharingDisplays, err := c.lazState.Client.Get(ctx, kSharing)
		if err != nil {
			c.Info("unable to get sharing displays:", zap.Error(err))
			setSharing = false
		}

		if setSharing {
			c.Debug("LazState: got sharing")
			var sharingDisplays Sharing
			err = json.Unmarshal(jSharingDisplays.Data, &sharingDisplays)
			if err != nil {
				c.Warn("unable to unmarshal sharing displays: ", zap.Error(err))
				errCh <- fmt.Errorf("unable to unmarshal volume: %v", err)
			}
			c.Info("LazState: unmarhalled sharing")
			c.shareMutex.Lock()
			c.sharing = sharingDisplays
			c.shareMutex.Unlock()
		}

		c.Debug("LazState: setup finished")

	}()
	go func() {
		c.lazState.Subscription, err = c.lazState.Client.Subscribe(c.lazContext, &lazarette.Key{Key: roomID})
		if err != nil {
			c.Warn("unable to subscribe to lazarette", zap.Error(err))
			return
		}
		for {
			select {
			case <-ctx.Done():
				return
			case <-c.kill:
				return
			default:
				kv, err := c.lazState.Subscription.Recv()
				switch {
				case err == io.EOF:
					return
				case err != nil:
					c.Warn("something went wrong receiving change from remote", zap.Error(err))
					continue
				case kv == nil:
					continue
				}
				if strings.Contains(kv.Key, "_sharing_displays") {
					var sharingDisplays Sharing
					err = json.Unmarshal(kv.Data, &sharingDisplays)
					if err != nil {
						c.Warn("unable to unmarshal sharing", zap.Error(err))
					} else {
						c.shareMutex.Lock()
						c.sharing = sharingDisplays
						c.shareMutex.Unlock()
					}
				}

			}
		}
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

	// Set up who can share to who
	var s Shareable
	for _, p := range c.uiConfig.Presets {
		for i, d := range p.ShareableDisplays {
			shar := remove(p.ShareableDisplays, i)
			shareable := make([]ID, len(shar))
			for j := range shar {
				shareable[j] = ID(shar[j])
			}
			s[ID(d)] = shareable
		}
	}
	c.shareable = s

	//check if controlgroup is empty, if not turn on displays in controlgroup
	if len(c.selectedControlGroupID) > 0 {
		var displays []ID
		for _, display := range room.ControlGroups[c.selectedControlGroupID].DisplayBlocks {
			displays = append(displays, display.ID)
		}

		setPowerMessage := SetPowerMessage{
			DisplayBlocks: displays,
			Status:        "on",
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

// Wait waits until a client is dead
func (c *Client) Wait() {
	<-c.kill
}

// Close closes the client
func (c *Client) Close() {
	c.close.Do(func() {
		c.Info("Closing client. Bye!")

		c.lazCancel()
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
