package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/byuoitav/ui/bff"
	"github.com/byuoitav/ui/log"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

var (
	upgrader = websocket.Upgrader{
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func NewClient(c echo.Context) error {
	// TODO check if it's coming from localhost and accept that, figure out which preset it's supposed to be

	// open the websocket
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.P.Warn("unable to upgrade connection", zap.Error(err))
		return c.String(http.StatusBadRequest, "unable to upgrade connection "+err.Error())
	}
	defer ws.Close()

	closeWithReason := func(msg string) error {
		err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(4000, msg))
		if err != nil {
			log.P.Warn("unable to write close message", zap.Error(err))
		}

		return err
	}

	url := fmt.Sprintf("%s/%s/getPreset", os.Getenv("CODE_SERVICE_URL"), c.Param("key"))

	req, err := http.NewRequestWithContext(c.Request().Context(), "GET", url, nil)
	if err != nil {
		return closeWithReason(fmt.Sprintf("unable to build request to check room code: %s", err))
	}

	log.P.Info("Getting room/preset from control key", zap.String("key", c.Param("key")))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return closeWithReason(fmt.Sprintf("unable to make request to check room code: %s", err))
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return closeWithReason("invalid room control key")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return closeWithReason(fmt.Sprintf("unable to read response from code service: %s", err))
	}

	preset := struct {
		RoomID     string `json:"RoomID"`
		PresetName string `json:"PresetName"`
	}{}

	if err = json.Unmarshal(body, &preset); err != nil {
		return closeWithReason(fmt.Sprintf("unable to parse response from code service: %s. response: %s", err, body))
	}

	client, err := bff.RegisterClient(c.Request().Context(), preset.RoomID, preset.PresetName, c.Request().RemoteAddr)
	if err != nil {
		log.P.Warn("unable to register client", zap.Error(err))
		return closeWithReason(fmt.Sprintf("unable to register client: %s", err))
	}

	log.P.Info("Successfully registered client", zap.String("client", c.Request().RemoteAddr))

	wg := sync.WaitGroup{}
	wg.Add(2)

	// send messages out from backend
	go func() {
		defer wg.Done()

		for msg := range client.Out {
			data, err := json.Marshal(msg)
			if err != nil {
				client.Warn("unable to marshal message to send to client", zap.Error(err))
				continue
			}

			// log that we are sending a message
			if _, ok := msg["error"]; ok {
				client.Warn("sending error to client", zap.ByteString("message", data))
			} else {
				client.Debug("Sending message to client", zap.ByteString("message", data))
			}

			err = ws.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				client.Error("failed to write message", zap.Error(err))
				return // ?
			}
		}
	}()

	// recv messages from frontend
	go func() {
		defer wg.Done()

		for {
			msgType, msg, err := ws.ReadMessage()
			switch {
			case err != nil:
				client.Error("failed to read messsage", zap.Error(err))

				switch {
				case errors.Is(err, io.ErrUnexpectedEOF) || strings.Contains(err.Error(), io.ErrUnexpectedEOF.Error()):
					ws.Close()
					return
					// TODO what other errors are we getting?
				default:
					ws.Close()
					return
				}
			case msgType == websocket.PingMessage:
				// send a pong
			default:
				var m bff.Message
				err = json.Unmarshal(msg, &m)
				if err != nil {
					client.Warn("unable to unmarshal message", zap.Error(err))
					client.Out <- bff.ErrorMessage(fmt.Errorf("unable to parse message: %s", err))
					continue
				}

				go client.HandleMessage(m)
			}
		}
	}()

	wg.Wait()
	return nil
}
