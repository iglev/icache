package icache

import (
	"context"
)

// CacheIf cache interface
type CacheIf interface {
	Get(context.Context, string) (interface{}, error)
	Set(context.Context, string, interface{}, int32) error
	Del(context.Context, string) error
	IsErrNotFound(err error) bool
}
