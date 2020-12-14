package ui

import "context"

type Client interface {
	HandleMessage([]byte)
	OutgoingMessages() <-chan []byte
	Done() <-chan struct{}
	Close()
	Refresh()
}

type ClientBuilder interface {
	New(ctx context.Context, room, controlGroup string) (Client, error)
}
