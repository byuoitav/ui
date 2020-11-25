package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/byuoitav/ui/av"
	"github.com/byuoitav/ui/client"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/singleflight"
)

type dataServiceConfig struct {
	Addr     string
	Username string
	Password string
}

func main() {
	var (
		port     int
		logLevel string

		avAPIURL          string
		keyServiceAddr    string
		host              string
		roomID            string
		deviceID          string
		lazaretteAddr     string
		lazaretteSSL      bool
		dataServiceConfig dataServiceConfig
	)

	pflag.IntVarP(&port, "port", "P", 8080, "port to run the server on")
	pflag.StringVarP(&logLevel, "log-level", "L", "info", "level to log at. refer to https://godoc.org/go.uber.org/zap/zapcore#Level for options")
	pflag.StringVarP(&host, "host", "", "rooms.av.byu.edu", "host of this server to display")
	pflag.StringVarP(&roomID, "room", "", "", "room this device is in")
	pflag.StringVarP(&deviceID, "device", "", "", "id of this device")
	pflag.StringVarP(&avAPIURL, "control-api", "", "http://localhost:8000", "base url of the av-control-api server to use")
	pflag.StringVarP(&keyServiceAddr, "key-service", "", "control-keys.avs.byu.edu", "address of the code service to use")
	pflag.StringVarP(&lazaretteAddr, "lazarette", "l", "localhost:7777", "address of the lazarette cache to use")
	pflag.BoolVar(&lazaretteSSL, "lazarette-use-ssl", false, "include to enable lazarette tls/ssl")
	pflag.StringVar(&dataServiceConfig.Addr, "db-address", "", "database address")
	pflag.StringVar(&dataServiceConfig.Username, "db-username", "", "database username")
	pflag.StringVar(&dataServiceConfig.Password, "db-password", "", "database password")
	pflag.Parse()

	// ctx for setup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, log := logger(logLevel)
	defer log.Sync() // nolint:errcheck

	ds := dataService(ctx, dataServiceConfig)

	handlers := &handlers{
		roomID:      roomID,
		deviceID:    deviceID,
		log:         log,
		single:      singleflight.Group{},
		dataService: ds,
		config: &client.Builder{
			DataService: ds,
			AVController: &av.Controller{
				BaseURL: avAPIURL,
			},
			Log: log,
		},
		upgrader: websocket.Upgrader{
			EnableCompression: true,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	r := gin.New()
	r.Use(gin.Recovery())

	debug := r.Group("/debug")
	debug.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "healthy")
	})
	debug.GET("/logz", func(c *gin.Context) {
		c.String(http.StatusOK, config.Level.String())
	})
	debug.GET("/logz/:level", func(c *gin.Context) {
		var level zapcore.Level
		if err := level.Set(c.Param("level")); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		fmt.Printf("***\n\tSetting log level to %s\n***\n", level.String())
		config.Level.SetLevel(level)
		c.String(http.StatusOK, config.Level.String())
	})

	api := r.Group("/api/v1/")
	api.GET("/ws", handlers.Websocket)
	api.GET("/ws/:key", handlers.Websocket)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("unable to bind listener", zap.Error(err))
	}

	log.Info("Starting server", zap.String("on", lis.Addr().String()))
	err = r.RunListener(lis)
	switch {
	case errors.Is(err, http.ErrServerClosed):
	case err != nil:
		log.Fatal("failed to serve", zap.Error(err))
	}
}
