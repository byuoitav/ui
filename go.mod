module github.com/byuoitav/ui

go 1.15

require (
	github.com/byuoitav/av-control-api v0.1.0
	github.com/gin-gonic/gin v1.6.3
	github.com/gorilla/websocket v1.4.2
	github.com/spf13/pflag v1.0.5
	go.uber.org/zap v1.16.0
	golang.org/x/sync v0.0.0-20201008141435-b3e1573b7520
)

replace github.com/byuoitav/av-control-api => ../av-control-api
