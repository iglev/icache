package icache

import "sync/atomic"

// Stats stat
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

// AddGet add get
func (s *Stats) AddGet(n int64) {
	atomic.AddInt64(&s.GetCnt, n)
}

// AddDel add del
func (s *Stats) AddDel(n int64) {
	atomic.AddInt64(&s.DelCnt, n)
}

// AddHit add hit
func (s *Stats) AddHit(n int64) {
	atomic.AddInt64(&s.HitCnt, n)
}

// AddErr add err
func (s *Stats) AddErr(n int64) {
	atomic.AddInt64(&s.ErrCnt, n)
}

// AddMiss add miss
func (s *Stats) AddMiss(n int64) {
	atomic.AddInt64(&s.MissCnt, n)
}

// AddSource add source
func (s *Stats) AddSource(n int64) {
	atomic.AddInt64(&s.SourceCnt, n)
}

// AddSourceHit add source hit
func (s *Stats) AddSourceHit(n int64) {
	atomic.AddInt64(&s.SourceHitCnt, n)
}

// AddSourceErr add source err
func (s *Stats) AddSourceErr(n int64) {
	atomic.AddInt64(&s.SourceErrCnt, n)
}
