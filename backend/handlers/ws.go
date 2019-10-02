package handlers

import (
	"net/http"
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
	// client.close?

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.P.Warn("unable to upgrade connection", zap.Error(err))
		return c.String(http.StatusBadRequest, "unable to upgrade connection "+err.Error())
	}
	defer ws.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		for msg := range client.Out {
			client.Info("writing message to client", zap.ByteString("toClient", msg))

			err := ws.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				// log.Printf("failed to write message: %s\n", err)
				return
			}
		}
	}()

	wg.Wait()
	return nil
}
