package daemon

import (
	"sync"
	"time"
)

const (
	// SweepEvery is the longest time between two Go build cache sweeps.
	SweepEvery = time.Hour
	// SweepUrgentEvery is the least time between two sweeps forced by low
	// free space.
	SweepUrgentEvery = 10 * time.Minute
)

// SweepSchedule decides when the daemon sweeps the Go build cache: at once
// when it starts, then hourly, and at once when free space is under twice
// minFree, but not again within SweepUrgentEvery of the last sweep.
type SweepSchedule struct {
	mu      sync.Mutex
	last    time.Time
	running bool
	wg      sync.WaitGroup
}

// Due reports whether a sweep is due at now, given the free space and the
// minFree guard. A minFree of 0 (guard off) or unknown free space never
// makes a sweep urgent. Call Done after the sweep ran.
func (s *SweepSchedule) Due(now time.Time, free, minFree uint64, freeKnown bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.due(now, free, minFree, freeKnown)
}

func (s *SweepSchedule) due(now time.Time, free, minFree uint64, freeKnown bool) bool {
	if s.last.IsZero() {
		return true
	}
	since := now.Sub(s.last)
	if since >= SweepEvery {
		return true
	}
	return freeKnown && minFree > 0 && free < 2*minFree && since >= SweepUrgentEvery
}

// Done records that a sweep ran at now.
func (s *SweepSchedule) Done(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last = now
}

// Start claims a sweep for a goroutine of its own: it reports true, and
// marks the sweep running and started at now, only when none is running and
// one is due. The caller must call Finish when the sweep ends.
func (s *SweepSchedule) Start(now time.Time, free, minFree uint64, freeKnown bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running || !s.due(now, free, minFree, freeKnown) {
		return false
	}
	s.running = true
	s.wg.Add(1)
	s.last = now
	return true
}

// Finish marks the sweep Start claimed as ended.
func (s *SweepSchedule) Finish() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		s.running = false
		s.wg.Done()
	}
}

// Wait blocks until the running sweep, if any, has finished.
func (s *SweepSchedule) Wait() { s.wg.Wait() }
