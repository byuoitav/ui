package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
	AvAPIAddr         string
	CodeServiceAddr   string
	RemoteControlAddr string
	LazaretteAddr     string
	LazaretteSSL      bool

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
	cconfig := bff.ClientConfig{
		AvAPIAddr:         b.AvAPIAddr,
		CodeServiceAddr:   b.CodeServiceAddr,
		RemoteControlAddr: b.RemoteControlAddr,
		LazaretteAddr:     b.LazaretteAddr,
		LazaretteSSL:      b.LazaretteSSL,
	}

	// if it is coming from localhost then don't worry about a key
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

		ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
		defer cancel()

		room, cgID, err := bff.GetRoomAndControlGroup(ctx, b.CodeServiceAddr, c.Param("key"))
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
