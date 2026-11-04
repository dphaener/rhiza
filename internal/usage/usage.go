// Package usage works out the usage windows of an agent's subscription from
// timestamped token records. It is pure: records and limit resets come in,
// windows come out, and nothing reads a clock, a file or an adapter.
package usage

import (
	"sort"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
)

// Kind names a kind of usage window.
type Kind string

// The two window kinds the governor tracks.
const (
	Session Kind = "session"
	Weekly  Kind = "weekly"
)

// Kinds lists every window kind, in a fixed order.
var Kinds = []Kind{Session, Weekly}

// Length is how long a window of this kind lasts. An unknown kind has none.
func (k Kind) Length() time.Duration {
	switch k {
	case Session:
		return 5 * time.Hour
	case Weekly:
		return 7 * 24 * time.Hour
	}
	return 0
}

// Limit is one usage-limit hit: when it was seen and the reset time its
// message gave.
type Limit struct {
	HitAt   time.Time
	ResetAt time.Time
}

// Window is one usage window: it runs from Start to Reset (exclusive), and
// Used counts the tokens of the records inside it.
type Window struct {
	Kind  Kind
	Start time.Time
	Reset time.Time
	Used  int64
}

// Contains reports whether t falls inside the window.
func (w Window) Contains(t time.Time) bool {
	return !t.Before(w.Start) && t.Before(w.Reset)
}

// Windows returns every window of the kind that the records and limits imply,
// oldest first. A window starts at the first record at or after the previous
// window's reset. A limit hit inside a window moves that window to the reset
// time the limit gave: it ends there and starts one window length earlier,
// but never before the previous window's reset. Records need not be sorted.
func Windows(kind Kind, records []adapter.UsageRecord, limits []Limit) []Window {
	length := kind.Length()
	if length <= 0 || len(records) == 0 {
		return nil
	}
	recs := append([]adapter.UsageRecord(nil), records...)
	sort.SliceStable(recs, func(i, j int) bool { return recs[i].At.Before(recs[j].At) })
	lims := append([]Limit(nil), limits...)
	sort.SliceStable(lims, func(i, j int) bool { return lims[i].HitAt.Before(lims[j].HitAt) })

	var out []Window
	var prevReset time.Time
	for _, r := range recs {
		if n := len(out); n > 0 && r.At.Before(out[n-1].Reset) {
			if !r.At.Before(out[n-1].Start) {
				out[n-1].Used += r.Tokens
			}
			continue
		}
		w := Window{Kind: kind, Start: r.At, Reset: r.At.Add(length)}
		for _, l := range lims {
			if l.HitAt.Before(w.Start) || !l.HitAt.Before(w.Reset) || !l.ResetAt.After(l.HitAt) {
				continue
			}
			w.Reset = l.ResetAt
			w.Start = l.ResetAt.Add(-length)
			if w.Start.Before(prevReset) {
				w.Start = prevReset
			}
			break
		}
		if r.At.Before(w.Start) {
			// A limit pulled the start past this record; it belongs to no window.
			continue
		}
		w.Used = r.Tokens
		out = append(out, w)
		prevReset = w.Reset
	}
	return out
}

// Current returns the window of the kind that holds now, and false when none
// does (the last window has reset, or there are no records).
func Current(kind Kind, records []adapter.UsageRecord, limits []Limit, now time.Time) (Window, bool) {
	ws := Windows(kind, records, limits)
	for i := len(ws) - 1; i >= 0; i-- {
		if ws[i].Contains(now) {
			return ws[i], true
		}
	}
	return Window{}, false
}

// KindOf says which window a limit hit belongs to, from its reset time: a
// reset 5 hours or less after the hit is the session window, a later one the
// weekly window. A hit with no reset time, or one that resets before it was
// seen, teaches nothing.
func KindOf(l Limit) (Kind, bool) {
	d := l.ResetAt.Sub(l.HitAt)
	if l.ResetAt.IsZero() || d <= 0 {
		return "", false
	}
	if d <= Session.Length() {
		return Session, true
	}
	return Weekly, true
}

// UsedAt is the tokens the kind's window held at the hit: those of its
// records inside the window the hit falls in, up to the hit. Records should
// reach back at least two window lengths before the hit. It is false when no
// window holds the hit.
func UsedAt(kind Kind, records []adapter.UsageRecord, l Limit) (int64, bool) {
	var upto []adapter.UsageRecord
	for _, r := range records {
		if !r.At.After(l.HitAt) {
			upto = append(upto, r)
		}
	}
	w, ok := Current(kind, upto, []Limit{l}, l.HitAt)
	return w.Used, ok
}
