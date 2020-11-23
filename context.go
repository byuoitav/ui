package ui

import "context"

type contextKey int

const (
	_keyRequestID contextKey = iota
)

// RequestID pulls a request ID from a context.Context.
func RequestID(ctx context.Context) string {
	id, ok := ctx.Value(_keyRequestID).(string)
	if !ok {
		return ""
	}

	return id
}

// WithRequestID returns a new context.Context, based on ctx, with the request id set.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, _keyRequestID, id)
}
