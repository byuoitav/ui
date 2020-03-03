package bff

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/byuoitav/lazarette/lazarette"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	lazSharingDisplays = "-sharingDisplays"
)

// LazaretteState represents the current lazarette state in a room
type LazaretteState struct {
	*sync.Map
}

// ShareDataMap is contained in the LazaretteState map and contains information about sharing in a room
type ShareDataMap map[ID]ShareData

// ShareData is all the information a sharing display needs to know
type ShareData struct {
	State    ShareState
	Active   []ID
	Inactive []ID
	Master   ID
}

type lazMessage struct {
	Key  string
	Data interface{}
}

func (c *Client) getShareMap() ShareDataMap {
	if ishareMap, ok := c.lazs.Load(lazSharingDisplays); ok {
		if shareMap, ok := ishareMap.(ShareDataMap); ok {
			return shareMap
		}
	}

	return nil
}

func (c *Client) setShareMap(l ShareDataMap) {
	c.lazs.Store(lazSharingDisplays, l)
}

// ConnectToLazarette dials lazarette and returns a new client. The connection will be killed
// when ctx expires.
func ConnectToLazarette(ctx context.Context, addr string) (lazarette.LazaretteClient, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("unable to open grpc connection: %s", err)
	}

	// TODO reconnect

	return lazarette.NewLazaretteClient(conn), nil
}

func (c *Client) updateLazaretteState(laz lazarette.LazaretteClient) {
	for {
		select {
		case <-c.kill:
			return
		case message := <-c.lazUpdates:
			c.stats.Lazarette.UpdatesSent++

			key := fmt.Sprintf("%v%v", c.roomID, message.Key)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			// cancel needs to be called before this block ends!!!

			data, err := json.Marshal(message.Data)
			if err != nil {
				c.Warn("unable to marshal lazarette message", zap.String("key", key), zap.Error(err))
				cancel()
				continue
			}

			_, err = laz.Set(ctx, &lazarette.KeyValue{
				Timestamp: ptypes.TimestampNow(),
				Key:       key,
				Data:      data,
			})
			if err != nil {
				c.Warn("unable to set updated key to lazarette", zap.String("key", key), zap.Error(err))
				cancel()
			}

			// store it in our local map
			c.lazs.Store(key, data)
			cancel()
		}
	}
}

func (c *Client) subLazaretteState(sub lazarette.Lazarette_SubscribeClient) {
	for {
		select {
		case <-c.kill:
			return
		default:
			kv, err := sub.Recv()
			switch {
			case errors.Is(err, io.EOF):
				c.Warn("lazarette stream ended", zap.Error(err))
				return
			case err != nil:
				s := status.Convert(err)
				if s.Code() == codes.Canceled || s.Code() == codes.DeadlineExceeded {
					c.Warn("ending lazarette stream", zap.Error(s.Err()))
					return
				}

				c.Warn("lazarette stream error", zap.Error(err))
				continue
			case kv == nil:
				continue
			}

			c.stats.Lazarette.UpdatesRecieved++

			// strip off beginning roomID so that we only have the actual key
			key := strings.TrimPrefix(kv.GetKey(), c.roomID)

			// stick the value into our map
			switch key {
			case lazSharingDisplays:
				var data ShareDataMap
				if err := json.Unmarshal(kv.GetData(), &data); err != nil {
					c.Warn("unable to parse share data from lazarette", zap.Error(err))
					continue
				}

				c.Debug("Got lazarette update", zap.String("key", lazSharingDisplays), zap.ByteString("data", kv.GetData()))
				c.lazs.Store(lazSharingDisplays, data)
			default:
			}

			// TODO get a new room and send it?
		}
	}
}
