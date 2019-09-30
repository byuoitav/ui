package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/byuoitav/ui/handlers"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.GET("ws/:key", handlers.NewClient)

	err := e.Start(":88")
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("failed to start server: %s", err)
	}
}
