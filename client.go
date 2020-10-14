package ui

import "context"

type Client interface {
	HandleMessage([]byte)
	OutgoingMessages() chan []byte
}

type ClientConfig interface {
	New(ctx context.Context, room, controlGroup string) (Client, error)
}
