package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/byuoitav/ui/handlers"
	"github.com/byuoitav/ui/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func main() {
	var port int
	var logLevel int

	pflag.IntVarP(&port, "port", "p", 8080, "port to run the server on")
	pflag.IntVarP(&logLevel, "log-level", "l", 2, "level of logging wanted. 1=DEBUG, 2=INFO, 3=WARN, 4=ERROR, 5=PANIC")
	pflag.Parse()

	setLog := func(level int) error {
		switch level {
		case 1:
			fmt.Printf("\nSetting log level to *debug*\n\n")
			log.Config.Level.SetLevel(zap.DebugLevel)
		case 2:
			fmt.Printf("\nSetting log level to *info*\n\n")
			log.Config.Level.SetLevel(zap.InfoLevel)
		case 3:
			fmt.Printf("\nSetting log level to *warn*\n\n")
			log.Config.Level.SetLevel(zap.WarnLevel)
		case 4:
			fmt.Printf("\nSetting log level to *error*\n\n")
			log.Config.Level.SetLevel(zap.ErrorLevel)
		case 5:
			fmt.Printf("\nSetting log level to *panic*\n\n")
			log.Config.Level.SetLevel(zap.PanicLevel)
		default:
			return errors.New("invalid log level: must be [1-4]")
		}

		return nil
	}

	// set the initial log level
	if err := setLog(logLevel); err != nil {
		log.P.Fatal("unable to set log level", zap.Error(err), zap.Int("got", logLevel))
	}

	// build echo server
	e := echo.New()

	// register new clients
	e.GET("/ws", handlers.NewClient)
	e.GET("/ws/:key", handlers.NewClient)

	// handle load balancer status check
	e.GET("/status", func(c echo.Context) error {
		return c.String(http.StatusOK, "healthy")
	})

	// set the log level
	e.GET("/log/:level", func(c echo.Context) error {
		level, err := strconv.Atoi(c.Param("level"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if err := setLog(level); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.String(http.StatusOK, fmt.Sprintf("Set log level to %v", level))
	})

	// serve uis
	e.Group("/", middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "dragonfruit",
		Index:  "index.html",
		HTML5:  true,
		Browse: true,
	}))

	// TODO add an if check on these ui's, only serve them if on a pi
	e.Group("/blueberry", middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "blueberry",
		Index:  "index.html",
		HTML5:  true,
		Browse: true,
	}))

	e.Group("/cherry", middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "cherry",
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
