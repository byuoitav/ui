package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/byuoitav/ui/bff"
	"github.com/byuoitav/ui/log"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

var (
	upgrader = websocket.Upgrader{
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func NewClient(c echo.Context) error {
	// TODO check that the room ID is valid, or do that in middleware
	client, err := bff.RegisterClient(c.Request().Context(), c.Param("key"), "", c.Request().RemoteAddr)
	if err != nil {
		log.P.Warn("unable to register client", zap.Error(err))
		return c.String(http.StatusInternalServerError, err.Error())
	}
	// TODO client.close?

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.P.Warn("unable to upgrade connection", zap.Error(err))
		return c.String(http.StatusBadRequest, "unable to upgrade connection "+err.Error())
	}
	defer ws.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	// send messages out
	go func() {
		defer wg.Done()

		for msg := range client.Out {
			data, err := json.Marshal(msg)
			if err != nil {
				client.Warn("unable to marshal message to send to client", zap.Error(err))
				continue
			}

			// log that we are sending a message
			if _, ok := msg["error"]; ok {
				client.Warn("sending error to client", zap.ByteString("msg", data))
			} else {
				client.Debug("Sending message to client", zap.ByteString("msg", data))
			}

			err = ws.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				client.Error("failed to write message", zap.Error(err))
				return // ?
			}
		}
	}()

	// recv messages
	go func() {
		defer wg.Done()

		for {
			msgType, msg, err := ws.ReadMessage()
			switch {
			case err != nil:
				client.Error("failed to read messsage", zap.Error(err))
				switch {
				case errors.Is(err, io.ErrUnexpectedEOF) || strings.Contains(err.Error(), io.ErrUnexpectedEOF.Error()):
					ws.Close()
					return
				}
			case msgType == websocket.PingMessage:
				// send a pong
			default:
				var m bff.Message
				err = json.Unmarshal(msg, &m)
				if err != nil {
					client.Warn("unable to unmarshal message", zap.Error(err))
					client.Out <- bff.ErrorMessage("unable to parse message: %w", err)
					continue
				}

				resps := client.HandleMessage(m)
				for resp := range resps {
					client.Out <- resp
				}
			}
		}
	}()

	wg.Wait()
	return nil
}
