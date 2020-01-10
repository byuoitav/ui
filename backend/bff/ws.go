package bff

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// maxMessageSize is the max message size allowed from the websocket
	maxMessageSize = 512 // bytes

	// pongWait is max time we'll wait for the websocket between pings
	pongWait = 60 * time.Second

	// pingPeriod is how often we'll send a ping to the client
	//	3/4 of the pongWait time
	pingPeriod = (pongWait * 3) / 4

	// time allowed to send a message to the client
	writeWait = 10 * time.Second
)

// readPump receives messages from the frontend and passes them to
//	c.HandleMessage()
func (c *Client) readPump() {
	defer c.Close()

	// set max message size
	c.ws.SetReadLimit(maxMessageSize)

	// define what to do on receving a pong
	_ = c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		_ = c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// read messages from websocket
	for {
		select {
		case <-c.kill:
			return
		default:
			_, msg, err := c.ws.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					c.Error("websocket closed unexpectedly", zap.Error(err))
				} else {
					c.Error("unable to read message from websocket", zap.Error(err))
				}

				return
			}

			// parse message
			var m Message
			if err := json.Unmarshal(msg, &m); err != nil {
				c.Warn("unable to parse message", zap.Error(err))
				c.Out <- ErrorMessage(fmt.Errorf("unable to parse message: %s", err))
				continue
			}

			// handle message
			go c.HandleMessage(m)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer c.Close()
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.Out:
			if !ok {
				// my channel got closed, must be time to stop
				_ = c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// marshal the message
			data, err := json.Marshal(msg)
			if err != nil {
				c.Warn("unable to marshal message to send to client", zap.Error(err))
				continue
			}

			// log that we are sending a message
			if _, ok := msg["error"]; ok {
				c.Warn("sending error to client", zap.ByteString("message", data))
			} else {
				c.Debug("Sending message to client", zap.ByteString("message", data))
			}

			// set our write deadline
			_ = c.ws.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.ws.WriteMessage(websocket.TextMessage, data); err != nil {
				c.Error("unable to write message to client", zap.Error(err))
				return // bye
			}
		case <-ticker.C:
			_ = c.ws.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case <-c.kill:
			return
		}
	}
}
