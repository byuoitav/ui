package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/byuoitav/ui"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	// maxMessageSize is the max message size allowed from the websocket
	maxMessageSize = 512 // bytes

	// pongWait is max time we'll wait for the websocket between pings
	pongWait = 60 * time.Second

	// pingPeriod is how often we'll send a ping to the client
	//	3/4 of the pongWait time
	pingPeriod = pongWait * (3 / 4)

	// time allowed to send a message to the client
	writeWait = 10 * time.Second

	// duration after getting an initial room message to wait for more
	roomDebounceDuration = 400 * time.Millisecond
)

func (h *handlers) Websocket(c *gin.Context) {
	// upgrade the connection
	ws, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Warn("unable to upgrade connection", zap.Error(err))
		c.String(http.StatusBadRequest, "unable to upgrade connection %s", err)
		return
	}
	defer ws.Close() // nolint:errcheck

	closeWith := func(msg string) {
		h.log.Warn("unable to create new client", zap.String("error", msg))

		// max control frame size is 125 bytes (https://tools.ietf.org/html/rfc6455#section-5.5)
		cmsg := websocket.FormatCloseMessage(4000, msg)
		if len(cmsg) > 125 {
			cmsg = cmsg[:125]
		}

		_ = ws.WriteMessage(websocket.CloseMessage, cmsg)
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var room, controlGroup string
	if h.deviceID == "" {
		// we get the control group using the key
		key := c.Param("key")
		h.log.Info("Getting room/preset from control key", zap.String("key", key))

		var err error

		room, controlGroup, err = h.dataService.RoomAndControlGroup(ctx, key)
		// TODO switch on different kinds of errors
		if err != nil {
			closeWith(fmt.Sprintf("unable to get room and control group: %s", err))
			return
		}
	} else {
		// we get the control group based on how this service was configured
		var err error

		controlGroup, err = h.dataService.ControlGroup(ctx, h.roomID, h.deviceID)
		if err != nil {
			closeWith(fmt.Sprintf("unable to get control group: %s", err))
			return
		}
	}

	// TODO save client in some sort of cache so we can send refresh messages/get stats
	client := h.config.New(room, controlGroup)

	errg, gctx := errgroup.WithContext(c.Request.Context())

	errg.Go(func() error {
		return h.writePump(gctx, client, ws)
	})

	errg.Go(func() error {
		return h.readPump(gctx, client, ws)
	})

	if err := errg.Wait(); err != nil {
		h.log.Warn("something went wrong in the pumps", zap.Error(err))
	}
}

// readPump receives messages on a websocket and passes them to the associated client
func (h *handlers) readPump(ctx context.Context, client ui.Client, ws *websocket.Conn) error {
	// set max message size
	ws.SetReadLimit(maxMessageSize)

	// define what to do when we get a pong
	_ = ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		// TODO why not just return this?
		_ = ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// read messages from websocket
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, msg, err := ws.ReadMessage()
			switch {
			case websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway):
				// not really an error? the websocket was just closed by the client
				return nil
			case err != nil:
				return fmt.Errorf("unable to read message: %w", err)
			}

			go client.HandleMessage(msg)
		}
	}
}

// writePump gets outgoing messages from the client and writes them to the websocket
func (h *handlers) writePump(ctx context.Context, client ui.Client, ws *websocket.Conn) error {
	// how frequently to write a ping to the websocket
	ping := time.NewTicker(pingPeriod)
	defer ping.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ping.C:
			_ = ws.SetWriteDeadline(time.Now().Add(writeWait))

			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return fmt.Errorf("unable to write ping: %w", err)
			}
		case msg, ok := <-client.OutgoingMessages():
			if !ok {
				return ws.WriteMessage(websocket.CloseMessage, []byte{})
			}

			_ = ws.SetWriteDeadline(time.Now().Add(writeWait))

			if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				return fmt.Errorf("unable to write message: %w", err)
			}
		}
	}
}
