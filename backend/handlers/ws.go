package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

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

	preset := struct {
		RoomID     string `json:"RoomID"`
		PresetName string `json:"PresetName"`
	}{}

	if strings.HasPrefix(c.Request().Host, "localhost") || strings.HasPrefix(c.Request().RemoteAddr, "127.0.0.1") {
		// if it is coming from localhost then don't worry about a key
		hostname := os.Getenv("SYSTEM_ID")
		if len(hostname) > 0 {
			log.P.Info("using hostname for localhost")

			hostnameArray := strings.Split(hostname, "-")
			preset.RoomID = fmt.Sprintf("%s-%s", hostnameArray[0], hostnameArray[1])

			uiConfig, err := bff.GetUIConfig(c.Request().Context(), http.DefaultClient, preset.RoomID)
			if err != nil {
				return closeWithReason(fmt.Sprintf("unable to get ui config: %s", err))
			}

			for _, p := range uiConfig.Panels {
				if p.Hostname == hostname {
					preset.PresetName = p.Preset
					break
				}
			}
		} else {
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

			if err = json.Unmarshal(body, &preset); err != nil {
				return closeWithReason(fmt.Sprintf("unable to parse response from code service: %s. response: %s", err, body))
			}
		}

	} else {
		// if not localhost then use the code service to get the info
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

		if err = json.Unmarshal(body, &preset); err != nil {
			return closeWithReason(fmt.Sprintf("unable to parse response from code service: %s. response: %s", err, body))
		}
	}

	client, err := bff.RegisterClient(c.Request().Context(), ws, preset.RoomID, preset.PresetName)
	if err != nil {
		log.P.Warn("unable to register client", zap.Error(err))
		return closeWithReason(fmt.Sprintf("unable to register client: %s", err))
	}

	log.P.Info("Successfully registered client", zap.String("client", c.Request().RemoteAddr))

	// if this function exists, the websocket connection is closed
	// so we need to wait for the client to be finished
	client.Wait()

	return nil
}
