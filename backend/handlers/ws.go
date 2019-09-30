package handlers

import (
	"net/http"

	"github.com/byuoitav/ui/bff"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
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
	// check that the room ID is valid, or do that in middleware

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.String(http.StatusBadRequest, "unable to upgrade connection "+err.Error())
	}
	defer ws.Close()

	bff.RegisterClient("ITB-1101", "The Cube")
	return nil
}
