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
)

const (
	// maxMessageSize is the max message size allowed from the websocket
	maxMessageSize = 512 // bytes

	// pongWait is max time we'll wait for the websocket between pings
	pongWait = 60 * time.Second

	// pingPeriod is how often we'll send a ping to the client
	//	3/4 of the pongWait time
	pingPeriod = (pongWait * 3) / 4

	// time allowed to send a message to the client
	writeWait = 10 * time.Second
)

func (h *handlers) Websocket(c *gin.Context) {
	// upgrade the connection
	ws, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Warn("unable to upgrade connection", zap.Error(err))
		c.String(http.StatusBadRequest, "unable to upgrade connection %s", err)
		return
	}
	// TODO attempt graceful shutdown of the connection?
	defer ws.Close() // nolint:errcheck

	closeWith := func(msg string) {
		h.log.Warn(msg)

		// max control frame size is 125 bytes (https://tools.ietf.org/html/rfc6455#section-5.5)
		cmsg := websocket.FormatCloseMessage(4000, msg)
		if len(cmsg) > 125 {
			cmsg = cmsg[:125]
		}

		_ = ws.WriteMessage(websocket.CloseMessage, cmsg)
	}

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	room, controlGroup, err := h.roomAndControlGroup(ctx, c.Param("key"))
	if err != nil {
		closeWith(err.Error())
		return
	}

	client, err := h.builder.New(ctx, room, controlGroup)
	if err != nil {
		closeWith(fmt.Sprintf("unable to build client: %s", err))
		return
	}
	defer client.Close()

	// add to clients so we can get stats/refresh it
	h.clientsMu.Lock()
	h.clients[c.Request.RemoteAddr] = client
	h.clientsMu.Unlock()

	// remove it from the map when we're done
	defer func() {
		h.clientsMu.Lock()
		delete(h.clients, c.Request.RemoteAddr)
		h.clientsMu.Unlock()
	}()

	// start the read/write pumps
	errCh := make(chan error)
	go func() {
		err := h.writePump(ctx, client, ws)
		select {
		case errCh <- err:
		default:
		}
	}()
	go func() {
		err := h.readPump(ctx, client, ws)
		select {
		case errCh <- err:
		default:
		}
	}()

	// wait for something to be done
	select {
	case err := <-errCh:
		if err != nil {
			h.log.Warn("something went wrong in the pumps", zap.Error(err))
			return
		}

		h.log.Info("Closing client due to websocket close")
	case <-ctx.Done():
		h.log.Info("Closing client & websocket due to context", zap.String("reason", ctx.Err().Error()))
	case <-client.Done():
		h.log.Info("Closing websocket due to client close")
	}
}

func (h *handlers) roomAndControlGroup(ctx context.Context, key string) (room, controlGroup string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if h.deviceID == "" {
		// we get the control group using the key
		h.log.Info("Getting room/preset from control key", zap.String("key", key))

		room, controlGroup, err = h.dataService.RoomAndControlGroup(ctx, key)
		// TODO switch on different kinds of errors
		if err != nil {
			return "", "", fmt.Errorf("unable to get room/controlGroup using key: %w", err)
		}

		return room, controlGroup, nil
	}

	// we get the control group based on how this service was configured
	controlGroup, err = h.dataService.ControlGroup(ctx, h.roomID, h.deviceID)
	if err != nil {
		return "", "", fmt.Errorf("unable to get controlGroup: %w", err)
	}

	return h.roomID, controlGroup, nil
}

// readPump receives messages on a websocket and passes them to the associated client
func (h *handlers) readPump(ctx context.Context, client ui.Client, ws *websocket.Conn) error {
	// set max message size
	ws.SetReadLimit(maxMessageSize)

	// define what to do when we get a pong
	_ = ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		return ws.SetReadDeadline(time.Now().Add(pongWait))
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
