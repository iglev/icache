package icache

import "context"

// GetterIf getter interface
type GetterIf interface {
	Get(context.Context, string, SinkIf) error
}

// GetterIfFunc func
type GetterIfFunc func(context.Context, string, SinkIf) error

// Get get
func (f GetterIfFunc) Get(ctx context.Context, strKey string, ifSink SinkIf) error {
	return f(ctx, strKey, ifSink)
}
