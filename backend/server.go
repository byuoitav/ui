package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/byuoitav/ui/bff"
	"github.com/byuoitav/ui/handlers"
	"github.com/byuoitav/ui/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

func main() {
	var (
		port     int
		logLevel int

		avApiAddr         string
		codeServiceAddr   string
		remoteControlAddr string
		lazaretteAddr     string
	)

	pflag.IntVarP(&port, "port", "P", 8080, "port to run the server on")
	pflag.IntVarP(&logLevel, "log-level", "L", 2, "level of logging wanted. 1=DEBUG, 2=INFO, 3=WARN, 4=ERROR, 5=PANIC")
	pflag.StringVarP(&avApiAddr, "av-api", "a", "localhost:8000", "address of the av-control-api to use")
	pflag.StringVarP(&codeServiceAddr, "code-service", "c", "control-keys.avs.byu.edu", "address of the code service to use")
	pflag.StringVarP(&remoteControlAddr, "remote-control", "r", "rooms.av.byu.edu", "address of the remote control to show")
	pflag.StringVarP(&lazaretteAddr, "lazarette", "l", "localhost:7777", "address of the lazarette cache to use")
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

	// handle load balancer status check
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "healthy")
	})

	// set the log level
	e.GET("/logz/:level", func(c echo.Context) error {
		level, err := strconv.Atoi(c.Param("level"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if err := setLog(level); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.String(http.StatusOK, fmt.Sprintf("Set log level to %v", level))
	})

	bffhandlers := handlers.BFF{
		AvApiAddr:         avApiAddr,
		CodeServiceAddr:   codeServiceAddr,
		RemoteControlAddr: remoteControlAddr,
		LazaretteAddr:     lazaretteAddr,
	}

	bffg := e.Group("", bffhandlers.SetupMiddleware)

	// register new clients
	bffg.GET("/ws", bffhandlers.NewClient)
	bffg.GET("/ws/:key", bffhandlers.NewClient)

	// uicontrol endpoints
	bffg.GET("/uicontrol/refresh", bffhandlers.RefreshClients)

	// stats
	statsg := bffg.Group("/statsz")
	statsg.GET("", bffhandlers.Stats)
	statsg.GET("/clients", bffhandlers.ClientStats)

	// group to get ui config
	single := singleflight.Group{}

	// serve the correct ui for this room
	e.Group("/", func(next echo.HandlerFunc) echo.HandlerFunc {
		root := "dragonfruit"
		id := os.Getenv("SYSTEM_ID")
		idsp := strings.Split(id, "-")

		if len(idsp) == 3 {
			roomID := fmt.Sprintf("%v-%v", idsp[0], idsp[1])

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// cherry/blueberry will return slowish because this gets called every time.
			// i dont think that really matters to much, but just a thought. singleflight should mitigate some of that
			iconfig, err, _ := single.Do(roomID, func() (interface{}, error) {
				config, err := bff.GetUIConfig(ctx, http.DefaultClient, roomID)
				if err != nil {
					return nil, err
				}

				return config, nil
			})
			if err != nil {
				return errHandler(http.StatusInternalServerError, err)
			}

			config, ok := iconfig.(bff.UIConfig)
			if !ok {
				return errHandler(http.StatusInternalServerError, fmt.Errorf("unexpected response from getting UI config: got %T", iconfig))
			}

			// find my config
			for _, panel := range config.Panels {
				if strings.EqualFold(panel.Hostname, id) {
					switch {
					case strings.Contains(panel.UIPath, "blueberry"):
						root = "blueberry"
					case strings.Contains(panel.UIPath, "cherry"):
						root = "cherry"
					}
				}
			}
		}

		return middleware.StaticWithConfig(middleware.StaticConfig{
			Root:   root,
			Index:  "index.html",
			HTML5:  true,
			Browse: true,
		})(next)
	})

	addr := fmt.Sprintf(":%d", port)
	log.P.Info("Starting server", zap.String("addr", addr))

	err := e.Start(addr)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.P.Fatal("failed to start server", zap.Error(err))
	}
}

func errHandler(code int, err error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.HTML(code, fmt.Sprintf(`
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
					<h1>%s</h1>
					<span>%s</span>
					<br /> <br /> <br />
					<span>This page will refresh in 10 seconds.</span>
				</body>
			</html>
		`, http.StatusText(code), err))
	}
}
