package bff

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/byuoitav/lazarette/lazarette"
	"google.golang.org/grpc"
)

type LazaretteState struct {
	*sync.Map
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

func (c *Client) syncLazaretteState(sub lazarette.Lazarette_SubscribeClient) {
	for {
		select {
		case <-c.kill:
			return
		case kv := <-c.lazUpdates:
		default:
			kv, err := sub.Recv()
			switch {
			case err == io.EOF:
				// TODO
			case err != nil:
				// TODO
			case kv == nil:
				continue
			}

			// strip off beginning roomID so that we only have the actual key
			key := strings.TrimPrefix(kv.GetKey(), c.roomID)

			// stick the value into our map
			switch key {
			case "-sharingDisplays":
				var sharingDisplays Sharing
				if err := json.Unmarshal(kv.GetData(), &sharingDisplays); err != nil {
					// TODO
				}

				c.lazs.Store(key, sharingDisplays)
			default:
			}

			// TODO get a new room and send it?
		}
	}
}
