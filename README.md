# icache
开箱即用的缓存库，插件式实现和支持singleflight回源

## CacheIf 缓存器接口
```golang
type CacheIf interface {
	Get(context.Context, string) (interface{}, error)
	Set(context.Context, string, interface{}, int32) error
	Del(context.Context, string) error
	IsErrNotFound(err error) bool
}
```

## GetterIf 回源接口
```golang
type GetterIf interface {
	Get(context.Context, string, SinkIf) error
}
```

## FlightGroupIf flight group
```golang
type FlightGroupIf interface {
	Do(string, func() (interface{}, error)) (interface{}, error)
}
// 两个协程并发回调回源
// 1: Get("key")
// 2: Get("key")
// 1: loadCache("key")
// 2: loadCache("key")
// 1: load("key")
// 2: load("key")
// 1: flightGroup.Do("key", fn)
// 1: fn()
// 2: flightGroup.Do("key", fn)
// 2: fn()
```

## Stats 操作统计
```golang
type Stats struct {
	GetCnt  int64 // cache get op cnt
	DelCnt  int64 // cache del op cnt
	HitCnt  int64 // cache hit cnt
	ErrCnt  int64 // cache errors cnt
	MissCnt int64 // cache miss cnt

	SourceCnt    int64 // get from source cnt
	SourceHitCnt int64 // get from source hit cnt
	SourceErrCnt int64 // get from source err cnt
}
```

## example
see icache_test.go
