package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var (
		port     int
		logLevel string

		avAPIAddr      string
		keyServiceAddr string
		host           string
		lazaretteAddr  string
		lazaretteSSL   bool
	)

	pflag.IntVarP(&port, "port", "P", 8080, "port to run the server on")
	pflag.StringVarP(&logLevel, "log-level", "L", "", "level to log at. refer to https://godoc.org/go.uber.org/zap/zapcore#Level for options")
	pflag.StringVarP(&host, "host", "", "rooms.av.byu.edu", "host of this server to display")
	pflag.StringVarP(&avAPIAddr, "control-api", "", "localhost:8000", "address of the av-control-api to use")
	pflag.StringVarP(&keyServiceAddr, "key-service", "", "control-keys.avs.byu.edu", "address of the code service to use")
	pflag.StringVarP(&lazaretteAddr, "lazarette", "l", "localhost:7777", "address of the lazarette cache to use")
	pflag.BoolVar(&lazaretteSSL, "lazarette-use-ssl", false, "include to enable lazarette tls/ssl")
	pflag.Parse()

	config, log := logger(logLevel)
	defer log.Sync() // nolint:errcheck

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
