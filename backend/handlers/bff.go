package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/byuoitav/ui/bff"
	"github.com/byuoitav/ui/log"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type BFF struct {
	AvApiAddr         string
	CodeServiceAddr   string
	RemoteControlAddr string
	LazaretteAddr     string

	init     sync.Once
	upgrader websocket.Upgrader

	clientsMu sync.Mutex
	clients   map[string]*bff.Client
}

func (b *BFF) SetupMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		b.init.Do(func() {
			b.clients = make(map[string]*bff.Client)
			b.upgrader = websocket.Upgrader{
				EnableCompression: true,
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}
		})

		return next(c)
	}
}

func (b *BFF) NewClient(c echo.Context) error {
	// open the websocket
	ws, err := b.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.P.Warn("unable to upgrade connection", zap.Error(err))
		return c.String(http.StatusBadRequest, "unable to upgrade connection "+err.Error())
	}
	defer ws.Close()

	closeWithReason := func(msg string) error {
		log.P.Warn("unable to create new client", zap.String("error", msg))

		// max control frame size is 125 bytes (https://tools.ietf.org/html/rfc6455#section-5.5)
		cmsg := websocket.FormatCloseMessage(4000, msg)
		if len(cmsg) > 125 {
			cmsg = cmsg[:125]
		}

		if err := ws.WriteMessage(websocket.CloseMessage, cmsg); err != nil {
			log.P.Warn("unable to write close message", zap.Error(err))
		}

		return nil
	}

	cconfig := bff.ClientConfig{
		AvApiAddr:         b.AvApiAddr,
		CodeServiceAddr:   b.CodeServiceAddr,
		RemoteControlAddr: b.RemoteControlAddr,
		LazaretteAddr:     b.LazaretteAddr,
	}

	// if it is coming from localhost then don't worry about a key
	hostname := os.Getenv("SYSTEM_ID")
	if len(hostname) > 0 {
		log.P.Info("using hostname for localhost")

		hostnameArray := strings.Split(hostname, "-")
		cconfig.RoomID = fmt.Sprintf("%s-%s", hostnameArray[0], hostnameArray[1])

		uiConfig, err := bff.GetUIConfig(c.Request().Context(), http.DefaultClient, cconfig.RoomID)
		if err != nil {
			return closeWithReason(fmt.Sprintf("unable to get ui config: %s", err))
		}

		for _, p := range uiConfig.Panels {
			if p.Hostname == hostname {
				cconfig.ControlGroupID = p.Preset
				break
			}
		}
	} else {
		// if not localhost then use the code service to get the info
		log.P.Info("Getting room/preset from control key", zap.String("key", c.Param("key")))

		room, cgID, err := bff.GetRoomAndControlGroup(c.Request().Context(), b.CodeServiceAddr, c.Param("key"))
		switch {
		case errors.Is(err, bff.ErrInvalidControlKey):
			return closeWithReason("Invalid control key")
		case err != nil:
			return closeWithReason(fmt.Sprintf("unable to get room/control group: %s", err))
		}

		cconfig.RoomID = room
		cconfig.ControlGroupID = cgID
	}

	client, err := bff.RegisterClient(c.Request().Context(), ws, cconfig)
	if err != nil {
		log.P.Warn("unable to register client", zap.Error(err))
		return closeWithReason(fmt.Sprintf("unable to register client: %s", err))
	}

	// add this client to the map
	b.clientsMu.Lock()
	b.clients[c.Request().RemoteAddr] = client
	b.clientsMu.Unlock()

	// defer deleting it from the map
	defer func() {
		go func() {
			// delay so that we can get stats off of it for a bit after it dies
			time.Sleep(30 * time.Second)

			b.clientsMu.Lock()
			defer b.clientsMu.Unlock()

			delete(b.clients, c.Request().RemoteAddr)
		}()
	}()

	log.P.Info("Successfully registered client", zap.String("client", c.Request().RemoteAddr))

	// if this function exits, the websocket connection is closed
	// so we need to wait for the client to be finished
	client.Wait()
	return nil
}

func (b *BFF) RefreshClients(c echo.Context) error {
	count := 0

	b.clientsMu.Lock()
	defer b.clientsMu.Unlock()

	for _, client := range b.clients {
		if !client.Killed() {
			count++
			go client.Refresh()
		}
	}

	log.P.Info("Refreshed clients.", zap.Int("count", count))
	return c.String(http.StatusOK, fmt.Sprintf("Successfully refreshed %d clients.", count))
}

func (b *BFF) Stats(c echo.Context) error {
	var stats []bff.ClientStats

	b.clientsMu.Lock()
	defer b.clientsMu.Unlock()

	for _, v := range b.clients {
		// don't include killed ones in agg
		if !v.Killed() {
			stats = append(stats, v.Stats())
		}
	}

	return c.JSON(http.StatusOK, bff.AggregateStats(stats))
}

func (b *BFF) ClientStats(c echo.Context) error {
	stats := make(map[string]bff.ClientStats)

	b.clientsMu.Lock()
	defer b.clientsMu.Unlock()

	for k, v := range b.clients {
		stats[k] = v.Stats()
	}

	return c.JSON(http.StatusOK, stats)
}
