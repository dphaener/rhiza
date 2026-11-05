package daemon

import "time"

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
	last time.Time
}

// Due reports whether a sweep is due at now, given the free space and the
// minFree guard. A minFree of 0 (guard off) or unknown free space never
// makes a sweep urgent. Call Done after the sweep ran.
func (s *SweepSchedule) Due(now time.Time, free, minFree uint64, freeKnown bool) bool {
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
func (s *SweepSchedule) Done(now time.Time) { s.last = now }
