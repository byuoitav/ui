package main

import (
	"github.com/byuoitav/ui/handlers"
	"github.com/labstack/echo"
)

func main() {
	// build echo server
	e := echo.New()

	bffhandlers := handlers.BFF{
		AvAPIAddr:         avAPIAddr,
		CodeServiceAddr:   codeServiceAddr,
		RemoteControlAddr: remoteControlAddr,
		LazaretteAddr:     lazaretteAddr,
		LazaretteSSL:      lazaretteSSL,
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
}
