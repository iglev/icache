package icache

import "github.com/juju/ratelimit"

// Option option
type Option struct {
	f func(ic *ICache)
}

// SetGetter set getter
func SetGetter(getter GetterIf) Option {
	return Option{func(ic *ICache) {
		ic.getter = getter
	}}
}

// SetFlightGroup set flight group
func SetFlightGroup(flightGroup FlightGroupIf) Option {
	return Option{func(ic *ICache) {
		ic.flightGroup = flightGroup
	}}
}

// SetCache set cache
func SetCache(cache CacheIf) Option {
	return Option{func(ic *ICache) {
		ic.cache = cache
	}}
}

// SetRateLimit set rate limit
func SetRateLimit(iPerSecLimit int64) Option {
	return Option{func(ic *ICache) {
		ic.rateLimiter = ratelimit.NewBucketWithRate(1.0, iPerSecLimit)
	}}
}
