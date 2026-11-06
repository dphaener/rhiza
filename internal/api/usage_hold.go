package api

// This file is the usage governor's hold. At each poll the daemon works out,
// for each usage window it knows a limit for, whether the use so far plus
// what the runs live and the next one are expected to burn would eat into
// the reserve. If so it holds everything it would launch, like the disk-low
// guard, until the window resets. Live runs and the owner's own starts are
// untouched, and so is a milestone request's plan run.

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/usage"
)

// usageMeanRuns is how many of the latest finished runs set the expected
// tokens of one run.
const usageMeanRuns = 10

// usageHold is a window that holds the daemon's launches.
type usageHold struct {
	Kind      usage.Kind
	Used      int64
	Projected int64
	Limit     int64
	Reserve   int
	Reset     time.Time
}

// Detail is the hold in words, for a skip.
func (h usageHold) Detail() string {
	return usageHoldDetail(daemon.UsageHold{Window: string(h.Kind), Used: h.Used, Projected: h.Projected, Limit: h.Limit, Reserve: h.Reserve, Reset: h.Reset})
}

// usageHoldDetail is a saved hold in words.
func usageHoldDetail(h daemon.UsageHold) string {
	return fmt.Sprintf("%s window: %d tokens used, %d projected, over the limit of %d less a %d%% reserve; resets %s",
		h.Window, h.Used, h.Projected, h.Limit, h.Reserve, h.Reset.UTC().Format(time.RFC3339))
}

// usageHoldState is the hold to save after a poll: nil when none holds, else
// prev with the poll's numbers, keeping its Since while the hold goes on.
func usageHoldState(h *usageHold, prev *daemon.UsageHold, now time.Time) *daemon.UsageHold {
	if h == nil {
		return nil
	}
	s := &daemon.UsageHold{Since: now.UTC(), Window: string(h.Kind), Used: h.Used, Projected: h.Projected, Limit: h.Limit, Reserve: h.Reserve, Reset: h.Reset}
	if prev != nil && prev.Window == s.Window && prev.Reset.Equal(s.Reset) {
		s.Since = prev.Since
	}
	return s
}

// usageReserve is the reserve setting, a percent of a limit.
func usageReserve(rows []Setting) int {
	n, err := strconv.Atoi(settingValue(rows, "rhiza.daemon.usage.reserve"))
	if err != nil || n < 0 || n > 100 {
		return 20
	}
	return n
}

// usageRunMean is the mean counted tokens of the latest finished runs that
// recorded any, and false when there are none.
func usageRunMean(runs []run.Journal) (int64, bool) {
	type fin struct {
		at     time.Time
		tokens int64
	}
	var done []fin
	for _, j := range runs {
		st := j.State()
		if !st.Ended {
			continue
		}
		var n int64
		for _, u := range st.Tokens {
			for _, c := range []run.TokenCounts{u.Session, u.Helpers} {
				n += c.Input + c.Output + c.CacheCreation
			}
		}
		if n > 0 {
			done = append(done, fin{st.EndedAt, n})
		}
	}
	if len(done) == 0 {
		return 0, false
	}
	slices.SortStableFunc(done, func(a, b fin) int { return b.at.Compare(a.at) })
	if len(done) > usageMeanRuns {
		done = done[:usageMeanRuns]
	}
	var sum int64
	for _, d := range done {
		sum += d.tokens
	}
	return sum / int64(len(done)), true
}

// usageRecordsFor reads the token records of every agent that reports them,
// for [from, to).
func usageRecordsFor(ctx *Ctx, from, to time.Time) []adapter.UsageRecord {
	if ctx.Getenv == nil {
		return nil
	}
	var recs []adapter.UsageRecord
	for _, n := range adapter.Names() {
		a, _ := adapter.Lookup(n)
		r, ok := a.(adapter.UsageReader)
		if !ok {
			continue
		}
		if got, err := r.UsageRecords(ctx.Getenv, from, to); err == nil {
			recs = append(recs, got...)
		}
	}
	return recs
}

// usageWindowsNow returns the current window of each kind that has one, and
// the windows to save for the next poll. A kind with a saved window builds
// its chain from that window's start while it holds, else from its reset, so
// the window stays where it was found; one with none looks back two window
// lengths. Only the limit hits of a kind move its windows. The records are
// read once, for the widest span any kind needs.
func usageWindowsNow(ctx *Ctx, common string, saved []daemon.UsageWindow) (map[usage.Kind]usage.Window, []daemon.UsageWindow) {
	now := ctx.Now()
	var hits []usage.Limit
	for _, h := range daemon.ReadLimitHitsIn(limitsDir(ctx, common)) {
		hits = append(hits, usage.Limit{HitAt: h.HitAt, ResetAt: h.ResetAt})
	}
	prev := map[usage.Kind]*daemon.UsageWindow{}
	for i := range saved {
		prev[usage.Kind(saved[i].Kind)] = &saved[i]
	}
	from := map[usage.Kind]time.Time{}
	widest := now
	for _, kind := range usage.Kinds {
		f := now.Add(-2 * kind.Length())
		if p := prev[kind]; p != nil {
			f = p.Reset
			if now.Before(p.Reset) {
				f = p.Start
			}
		}
		from[kind] = f
		if f.Before(widest) {
			widest = f
		}
	}
	recs := usageRecordsFor(ctx, widest, now.Add(time.Nanosecond))
	windows := map[usage.Kind]usage.Window{}
	var next []daemon.UsageWindow
	for _, kind := range usage.Kinds {
		var kr []adapter.UsageRecord
		for _, r := range recs {
			if !r.At.Before(from[kind]) {
				kr = append(kr, r)
			}
		}
		var kh []usage.Limit
		for _, h := range hits {
			if k, ok := usage.KindOf(h); ok && k == kind {
				kh = append(kh, h)
			}
		}
		if w, ok := usage.Current(kind, kr, kh, now); ok {
			windows[kind] = w
			next = append(next, daemon.UsageWindow{Kind: string(kind), Start: w.Start.UTC(), Reset: w.Reset.UTC(), Used: w.Used})
		} else if p := prev[kind]; p != nil {
			// No window holds now; the next one starts at or after this reset.
			next = append(next, *p)
		}
	}
	return windows, next
}

// usageHoldOf returns the first window whose projected use passes its limit
// less the reserve, nil when none does. A window with no known limit, or no
// current window (the last one has reset), never holds.
func usageHoldOf(rows []Setting, learned []daemon.LearnedLimit, runs []run.Journal, host string, windows map[usage.Kind]usage.Window) *usageHold {
	reserve := usageReserve(rows)
	live := int64(0)
	for _, j := range runs {
		if j.State().Liveness(host, nil) == run.LiveRunning {
			live++
		}
	}
	mean, _ := usageRunMean(runs)
	for _, kind := range usage.Kinds {
		limit, ok := usageLimit(rows, learned, kind)
		if !ok {
			continue
		}
		w, ok := windows[kind]
		if !ok {
			continue
		}
		projected := w.Used + mean*(live+1)
		if projected > limit-limit*int64(reserve)/100 {
			return &usageHold{Kind: kind, Used: w.Used, Projected: projected, Limit: limit, Reserve: reserve, Reset: w.Reset}
		}
	}
	return nil
}

// governorHold is the machine-wide hold from a limit hit any project's
// driver recorded: it stops new launches until the hit's reset time. nil
// when no hit holds.
func governorHold(ctx *Ctx, common string) *usageHold {
	h, ok := daemon.ActiveHold(limitsDir(ctx, common), ctx.Now())
	if !ok {
		return nil
	}
	kind, _ := usage.KindOf(usage.Limit{HitAt: h.HitAt, ResetAt: h.ResetAt})
	return &usageHold{Kind: kind, Reset: h.ResetAt}
}
