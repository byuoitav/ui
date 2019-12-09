package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/byuoitav/ui/handlers"
	"github.com/byuoitav/ui/log"
	"github.com/labstack/echo"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func main() {
	var port int
	var logLevel int

	pflag.IntVarP(&port, "port", "p", 8080, "port to run the server on")
	pflag.IntVarP(&logLevel, "log-level", "l", 3, "level of logging wanted. 1=DEBUG, 2=INFO, 3=WARN, 4=ERROR, 5=PANIC")
	pflag.Parse()

	switch logLevel {
	case 1:
		log.P.Info("Setting log level to debug")
		log.Config.Level.SetLevel(zap.DebugLevel)
	case 2:
		log.P.Info("Setting log level to info")
		log.Config.Level.SetLevel(zap.InfoLevel)
	case 3:
		log.P.Info("Setting log level to warn")
		log.Config.Level.SetLevel(zap.WarnLevel)
	case 4:
		log.P.Info("Setting log level to error")
		log.Config.Level.SetLevel(zap.ErrorLevel)
	case 5:
		log.P.Info("Setting log level to panic")
		log.Config.Level.SetLevel(zap.PanicLevel)
	default:
		log.P.Fatal("invalid log level. must be [1-4]", zap.Int("got", logLevel))
	}

	e := echo.New()

	e.GET("ws/:key", handlers.NewClient)
	e.Group("/", middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "dragonfruit",
		Index:  "index.html",
		HTML5:  true,
		Browse: true,
	}))

	addr := fmt.Sprintf(":%d", port)
	err := e.Start(addr)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.P.Fatal("failed to start server", zap.Error(err))
	}
}
