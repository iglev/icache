package icache

import "fmt"

var (
	// ErrNotFound not found err
	ErrNotFound = fmt.Errorf("err not found")
	// ErrCacheIf CacheIf err
	ErrCacheIf = fmt.Errorf("err CacheIf")
	// ErrGetterIf GetterIf err
	ErrGetterIf = fmt.Errorf("err GetterIf")
	// ErrRateLimit rate limit err
	ErrRateLimit = fmt.Errorf("err ratelimit")
)

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
