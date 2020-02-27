package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

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

type NewClientConfig struct {
	AvApiAddr         string
	CodeServiceAddr   string
	RemoteControlAddr string
}

func NewClientHandler(config NewClientConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		// open the websocket
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			log.P.Warn("unable to upgrade connection", zap.Error(err))
			return c.String(http.StatusBadRequest, "unable to upgrade connection "+err.Error())
		}
		defer ws.Close()

		closeWithReason := func(msg string) error {
			// max control frame size is 125 bytes (https://tools.ietf.org/html/rfc6455#section-5.5)
			cmsg := websocket.FormatCloseMessage(4000, msg)
			if len(cmsg) > 125 {
				cmsg = cmsg[:125]
			}

			if err := ws.WriteMessage(websocket.CloseMessage, cmsg); err != nil {
				log.P.Warn("unable to write close message", zap.Error(err))
			}

			return err
		}

		bffconfig := bff.ClientConfig{
			AvApiAddr:         config.AvApiAddr,
			CodeServiceAddr:   config.CodeServiceAddr,
			RemoteControlAddr: config.RemoteControlAddr,
		}

		// if it is coming from localhost then don't worry about a key
		hostname := os.Getenv("SYSTEM_ID")
		if len(hostname) > 0 {
			log.P.Info("using hostname for localhost")

			hostnameArray := strings.Split(hostname, "-")
			bffconfig.RoomID = fmt.Sprintf("%s-%s", hostnameArray[0], hostnameArray[1])

			uiConfig, err := bff.GetUIConfig(c.Request().Context(), http.DefaultClient, bffconfig.RoomID)
			if err != nil {
				return closeWithReason(fmt.Sprintf("unable to get ui config: %s", err))
			}

			for _, p := range uiConfig.Panels {
				if p.Hostname == hostname {
					bffconfig.ControlGroupID = p.Preset
					break
				}
			}
		} else {
			// if not localhost then use the code service to get the info
			log.P.Info("Getting room/preset from control key", zap.String("key", c.Param("key")))

			room, cgID, err := bff.GetRoomAndControlGroup(c.Request().Context(), config.CodeServiceAddr, c.Param("key"))
			switch {
			case errors.Is(err, bff.ErrInvalidControlKey):
				return closeWithReason("Invalid control key")
			case err != nil:
				return closeWithReason(fmt.Sprintf("unable to get room/control group: %s", err))
			}

			bffconfig.RoomID = room
			bffconfig.ControlGroupID = cgID
		}

		client, err := bff.RegisterClient(c.Request().Context(), ws, bffconfig)
		if err != nil {
			log.P.Warn("unable to register client", zap.Error(err))
			return closeWithReason(fmt.Sprintf("unable to register client: %s", err))
		}

		log.P.Info("Successfully registered client", zap.String("client", c.Request().RemoteAddr))

		// if this function exits, the websocket connection is closed
		// so we need to wait for the client to be finished
		client.Wait()
		return nil
	}
}
