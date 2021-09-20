package icache

import (
	"context"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

/*
	Get(context.Context, string) (interface{}, error)
	Set(context.Context, string, interface{}, int32) error
	Del(context.Context, string) error
	IsErrNotFound(err error) bool
*/

// LRUObjCache lru obj cache
type LRUObjCache struct {
	lru *lru.Cache
}

// NewLRUObjCache new lru cache
func NewLRUObjCache(iSize int) CacheIf {
	return NewLRUObjCacheWithEvict(iSize, nil)
}

// NewLRUObjCacheWithEvict new lru cache
func NewLRUObjCacheWithEvict(iSize int, onEvicted func(key interface{}, value interface{})) CacheIf {
	cache, err := lru.NewWithEvict(iSize, onEvicted)
	if err != nil {
		panic(err)
	}
	return &LRUObjCache{lru: cache}
}

type lruItem struct {
	val      interface{}
	expireTs int64
}

func (c *LRUObjCache) Get(ctx context.Context, strKey string) (interface{}, error) {
	valIf, ok := c.lru.Get(strKey)
	if !ok {
		return nil, ErrNotFound
	}
	item := valIf.(*lruItem)
	// check ttl
	if item.expireTs > 0 && time.Now().Unix() > item.expireTs {
		c.lru.Remove(strKey)
		return nil, ErrNotFound
	}
	return item.val, nil
}

func (c *LRUObjCache) Set(ctx context.Context, strKey string, valIf interface{}, iTTL int32) error {
	item := &lruItem{
		val: valIf,
	}
	if iTTL > 0 {
		item.expireTs = time.Now().Unix() + int64(iTTL)
	}
	c.lru.Add(strKey, item)
	return nil
}

func (c *LRUObjCache) Del(ctx context.Context, strKey string) error {
	c.lru.Remove(strKey)
	return nil
}

func (c *LRUObjCache) IsErrNotFound(err error) bool {
	return err == ErrNotFound
}
