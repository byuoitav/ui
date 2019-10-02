package main

import (
	"errors"
	"net/http"

	"github.com/byuoitav/ui/handlers"
	"github.com/byuoitav/ui/log"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func main() {
	log.Config.Level.SetLevel(zap.DebugLevel)
	e := echo.New()

	e.GET("ws/:key", handlers.NewClient)

	err := e.Start(":88")
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.P.Fatal("failed to start server", zap.Error(err))
	}
}
