package bff

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/byuoitav/lazarette/lazarette"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	lazSharingDisplays = "-sharingDisplays"
)

type LazaretteState struct {
	*sync.Map
}

type ShareDataMap map[ID]ShareData

type ShareData struct {
	State    ShareState
	Active   []ID
	Inactive []ID
	Master   ID
}

func (c *Client) getShareMap() ShareDataMap {
	if ishareMap, ok := c.lazs.Load(lazSharingDisplays); ok {
		if shareMap, ok := ishareMap.(ShareDataMap); ok {
			return shareMap
		}
	}

	return nil
}

func ConnectToLazarette(ctx context.Context) (lazarette.LazaretteClient, error) {
	lazAddr := os.Getenv("LAZARETTE_ADDR")
	if len(lazAddr) == 0 {
		return nil, fmt.Errorf("LAZARETTE_ADDR not set")
	}

	conn, err := grpc.DialContext(ctx, lazAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("unable to open grpc connection: %s", err)
	}

	return lazarette.NewLazaretteClient(conn), nil
}

func (c *Client) syncLazaretteState(laz lazarette.LazaretteClient, sub lazarette.Lazarette_SubscribeClient) {
	for {
		select {
		case <-c.kill:
			return
		case kv := <-c.lazUpdates:
			c.stats.Lazarette.UpdatesSent++
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			_, err := laz.Set(ctx, &kv)
			if err != nil {
				cancel()
				c.Warn("unable to set updated key to lazarette", zap.String("key", kv.GetKey()), zap.Error(err))
			}

			cancel()
		default:
			kv, err := sub.Recv()
			switch {
			case errors.Is(err, io.EOF):
				c.Warn("lazarette stream ended", zap.Error(err))
				return
			case err != nil:
				// c.Warn("lazarette stream error", zap.Error(err))
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

				c.lazs.Store(lazSharingDisplays, data)
			default:
			}

			// TODO get a new room and send it?
		}
	}
}
