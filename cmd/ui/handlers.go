package main

import (
	"context"
	"net/http"
	"path"
	"path/filepath"
	"time"

	"github.com/byuoitav/ui"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

type handlers struct {
	roomID   string
	deviceID string

	log         *zap.Logger
	single      singleflight.Group
	dataService ui.DataService
	config      ui.ClientConfig
	upgrader    websocket.Upgrader
}

func (h *handlers) ServeUI(c *gin.Context) {
	root := "dragonfruit/"

	if h.deviceID != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// cherry/blueberry will return slowish because this gets called every time.
		// i don't think that really matters too much, but just a thought.
		// singleflight should help mitigate some of that ¯\_(ツ)_/¯
		// if it's unbearable, we could always cache what ui it should use

		ui, err, _ := h.single.Do(h.deviceID, func() (interface{}, error) {
			ui, err := h.dataService.UIForDevice(ctx, h.roomID, h.deviceID)
			if err != nil {
				return nil, err
			}

			return ui, nil
		})

		if err != nil {
			c.Writer.Header().Add("Content-Type", "text/html")
			c.String(http.StatusInternalServerError, `
			<!DOCTYPE html>
			<html>
				<head>
					<script>
						function refresh() {
							setTimeout(() => {window.location.reload()}, 10000)
						}
					</script>
				</head>
				<body onload="refresh()">
					<h1>Internal Server Error</h1>
					<span>%s</span>
					<br /> <br /> <br />
					<span>This page will refresh in 10 seconds.</span>
				</body>
			</html>
			`, err)
			return
		}

		root = ui.(string) + "/"
	}

	dir, file := path.Split(c.Request.RequestURI)
	if file == "" || filepath.Ext(file) == "" {
		c.File(root + "index.html")
	} else {
		c.File(root + path.Join(dir, file))
	}
}
