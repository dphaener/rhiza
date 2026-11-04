package api

// This file is the usage governor's hold. At each poll the daemon works out,
// for each usage window it knows a limit for, whether the use so far plus
// what the runs live and the next one are expected to burn would eat into
// the reserve. If so it holds everything it would launch, like the disk-low
// guard, until the window resets. Live runs and the owner's own starts are
// untouched.

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
	return fmt.Sprintf("%s window: %d tokens used, %d projected, over the limit of %d less a %d%% reserve; resets %s",
		h.Kind, h.Used, h.Projected, h.Limit, h.Reserve, h.Reset.UTC().Format(time.RFC3339))
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

// usageHoldOf returns the first window whose projected use passes its limit
// less the reserve, nil when none does. A window with no known limit, or no
// current window (the last one has reset), never holds.
func usageHoldOf(ctx *Ctx, common string, rows []Setting, learned []daemon.LearnedLimit, runs []run.Journal, host string) *usageHold {
	now := ctx.Now()
	var hits []usage.Limit
	for _, h := range daemon.ReadLimitHits(common) {
		hits = append(hits, usage.Limit{HitAt: h.HitAt, ResetAt: h.ResetAt})
	}
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
		recs := usageRecordsFor(ctx, now.Add(-2*kind.Length()), now.Add(time.Nanosecond))
		w, ok := usage.Current(kind, recs, hits, now)
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
