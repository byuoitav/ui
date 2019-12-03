package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/byuoitav/ui/bff"
	"github.com/byuoitav/ui/log"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type preset struct {
	RoomID     string `json:"RoomID"`
	PresetName string `json:"PresetName"`
}

var (
	upgrader = websocket.Upgrader{
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func NewClient(c echo.Context) (string, error) {
	// TODO check that the room ID is valid, or do that in middleware
	controlKey := c.Param("key")
	var resp preset
	url := fmt.Sprintf("https://control-keys.avs.byu.edu/%s/getPreset", controlKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("an error occured while making the call: %w", err)
	}

	res, gerr := http.DefaultClient.Do(req)
	if gerr != nil {
		return "", fmt.Errorf("error when making call: %w", gerr)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error when unmarshalling the response: %w", err)
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Printf("%s/n", body)
		return "", fmt.Errorf("error when unmarshalling the response: %w", err)
	}

	client, err := bff.RegisterClient(c.Request().Context(), resp.RoomID, resp.PresetName, c.Request().RemoteAddr)
	if err != nil {
		log.P.Warn("unable to register client", zap.Error(err))
		return "", c.String(http.StatusInternalServerError, err.Error())
	}
	// TODO client.close?

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.P.Warn("unable to upgrade connection", zap.Error(err))
		return "", c.String(http.StatusBadRequest, "unable to upgrade connection "+err.Error())
	}
	defer ws.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	// send messages out
	go func() {
		defer wg.Done()

		for msg := range client.Out {
			var data []byte
			for _, v := range msg {
				data = v
				break
			}
			/*
				data, err := json.Marshal(msg)
				if err != nil {
					client.Warn("unable to marshal message to send to client", zap.Error(err))
					fmt.Printf("\nmap: %s\n", msg)
					continue
				}
			*/

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

	// recv messages
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
	return "", nil
}
