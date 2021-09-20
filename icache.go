package icache

import (
	"context"
	"fmt"

	"github.com/iglev/icache/singleflight"
	"github.com/juju/ratelimit"
)

// ICache ICache struct
type ICache struct {
	cache       CacheIf
	getter      GetterIf
	flightGroup FlightGroupIf

	stats       Stats
	rateLimiter *ratelimit.Bucket
}

// NewICache new ICache
func NewICache(opts ...Option) (*ICache, error) {
	ic := &ICache{
		stats: Stats{},
	}

	// do opt
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.f(ic)
		}
	}

	// check
	if ic.cache == nil {
		return nil, ErrCacheIf
	}
	if ic.getter == nil {
		return nil, ErrGetterIf
	}
	if ic.flightGroup == nil {
		ic.flightGroup = &singleflight.Group{}
	}

	return ic, nil
}

// Get get key
func (ic *ICache) Get(ctx context.Context, strKey string, dest SinkIf) error {
	ic.stats.AddGet(1)
	if dest == nil {
		ic.stats.AddErr(1)
		return fmt.Errorf("nil dest")
	}
	view, err := ic.loadCache(ctx, strKey)
	if err != nil {
		if !ic.cache.IsErrNotFound(err) {
			ic.stats.AddErr(1)
			// loadCache fail, whatever
			// go on ic.load
		}
	} else {
		// hit cache
		ic.stats.AddHit(1)
		return dest.SetView(view)
	}

	// miss
	bDestSetView := false
	view, bDestSetView, err = ic.load(ctx, strKey, dest)
	if err != nil {
		return err
	}
	if bDestSetView {
		return nil
	}
	return dest.SetView(view)
}

// Delete del key
func (ic *ICache) Delete(ctx context.Context, strKey string) error {
	ic.stats.AddDel(1)
	err := ic.cache.Del(ctx, strKey)
	if err != nil {
		ic.stats.AddErr(1)
		return err
	}
	return nil
}

// GetStat get stat
func (ic *ICache) GetStat() Stats {
	return ic.stats
}

func (ic *ICache) setCache(ctx context.Context, strKey string, view View) error {
	return ic.cache.Set(ctx, strKey, view.v, view.ttl)
}

// load cache
func (ic *ICache) loadCache(ctx context.Context, strKey string) (View, error) {
	var view View
	valIf, err := ic.cache.Get(ctx, strKey)
	if err != nil {
		return view, err
	}
	view.v = valIf
	return view, nil
}

// load load
func (ic *ICache) load(ctx context.Context, strKey string, dest SinkIf) (View, bool, error) {
	ic.stats.AddMiss(1)
	bDestSetView := false
	viewIf, err := ic.flightGroup.Do(strKey, func() (interface{}, error) {
		if view, err := ic.loadCache(ctx, strKey); err == nil {
			// hit
			ic.stats.AddHit(1)
			return view, nil
		} else if !ic.cache.IsErrNotFound(err) {
			// loadCache fail, go on
			ic.stats.AddErr(1)
		}
		// miss
		ic.stats.AddSource(1)
		view, err := ic.loadSource(ctx, strKey, dest)
		if err != nil {
			ic.stats.AddSourceErr(1)
			return nil, err
		}
		ic.stats.AddSourceHit(1)
		bDestSetView = true
		ic.setCache(ctx, strKey, view)
		return view, nil
	})
	if err != nil {
		return View{}, false, err
	}
	view := viewIf.(View)
	return view, bDestSetView, nil
}

// loadSource load source
func (ic *ICache) loadSource(ctx context.Context, strKey string, dest SinkIf) (View, error) {
	// rate limit
	if ic.rateLimiter != nil {
		if cnt := ic.rateLimiter.TakeAvailable(1); cnt <= 0 {
			return View{}, ErrRateLimit
		}
	}
	err := ic.getter.Get(ctx, strKey, dest)
	if err != nil {
		return View{}, err
	}
	return dest.GetView()
}
