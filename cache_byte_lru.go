package icache

import (
	"context"
	"fmt"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

/*
	Get(context.Context, string) (interface{}, error)
	Set(context.Context, string, interface{}, int32) error
	Del(context.Context, string) error
	IsErrNotFound(err error) bool
*/

// LRUByteCache
type LRUByteCache struct {
	lru *lru.Cache
}

// NewLRUByteCache new lru cache
func NewLRUByteCache(iSize int) CacheIf {
	return NewLRUByteCacheWithEvict(iSize, nil)
}

// NewLRUByteCacheWithEvict new lru cache
func NewLRUByteCacheWithEvict(iSize int, onEvicted func(key interface{}, value interface{})) CacheIf {
	cache, err := lru.NewWithEvict(iSize, onEvicted)
	if err != nil {
		panic(err)
	}
	return &LRUByteCache{lru: cache}
}

type lruByteItem struct {
	val      []byte
	expireTs int64
}

// Get get
func (c *LRUByteCache) Get(ctx context.Context, strKey string) (interface{}, error) {
	valIf, ok := c.lru.Get(strKey)
	if !ok {
		return nil, ErrNotFound
	}
	item := valIf.(*lruByteItem)
	// check ttl
	if item.expireTs > 0 && time.Now().Unix() > item.expireTs {
		c.lru.Remove(strKey)
		return nil, ErrNotFound
	}
	return item.val, nil
}

// Set set
func (c *LRUByteCache) Set(ctx context.Context, strKey string, valIf interface{}, iTTL int32) error {
	item := &lruByteItem{
		// val: valIf,
	}
	switch valIf.(type) {
	case []byte:
		item.val = valIf.([]byte)
	case string:
		item.val = []byte(valIf.(string))
	default:
		return fmt.Errorf("LRUByteCache only support []byte and string type")
	}
	if iTTL > 0 {
		item.expireTs = time.Now().Unix() + int64(iTTL)
	}
	c.lru.Add(strKey, item)
	return nil
}

// Del del
func (c *LRUByteCache) Del(ctx context.Context, strKey string) error {
	c.lru.Remove(strKey)
	return nil
}

// IsErrNotFound is not found err
func (c *LRUByteCache) IsErrNotFound(err error) bool {
	return err == ErrNotFound
}
