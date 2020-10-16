package ui

import "context"

type Client interface {
	HandleMessage([]byte)
	OutgoingMessages() chan []byte
}

type ClientBuilder interface {
	New(ctx context.Context, room, controlGroup string) (Client, error)
}
