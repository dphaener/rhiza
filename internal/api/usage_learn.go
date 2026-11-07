package api

// This file teaches the daemon its usage limits. A run driver that sees its
// session stop at a usage limit records the hit; the daemon reads the hits at
// each poll and saves each window's use at the latest hit as that window's
// learned limit, which stands in for a limit setting that is unset.

import (
	"strconv"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/rootcfg"
	"github.com/dphaener/rhiza/internal/usage"
)

// recordLimitHit records that a run driver's session stopped at a usage
// limit now. It is best effort: a driver never fails over it. A hit with no
// reset time is not recorded, for it teaches nothing.
func recordLimitHit(ctx *Ctx, a adapter.Adapter, l adapter.UsageLimit) {
	if l.ResetAt.IsZero() {
		return
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return
	}
	daemon.RecordLimitHitIn(limitsDir(ctx, common), daemon.LimitHit{HitAt: ctx.Now().UTC(), ResetAt: l.ResetAt.UTC(), Agent: a.Name()})
}

// limitsDir is where usage-limit hits are kept: the root's state folder
// when the working directory is under a root, so every project's drivers
// feed the one governor, else the repo's own folder under its common dir.
func limitsDir(ctx *Ctx, common string) string {
	if root, ok := rootcfg.Find(ctx.Wd); ok {
		return daemon.RootStateDir(root)
	}
	return daemon.Dir(common)
}

// learnLimits folds the recorded hits into state's learned limits: for each
// window kind the latest hit wins, with the window's use at that hit. A hit
// whose use cannot be read teaches nothing yet and is tried again next poll.
func learnLimits(ctx *Ctx, common string, state *daemon.State) {
	hits := daemon.ReadLimitHitsIn(limitsDir(ctx, common))
	if len(hits) == 0 {
		return
	}
	known := map[string]daemon.LearnedLimit{}
	for _, l := range state.Learned {
		known[l.Kind] = l
	}
	for _, h := range hits {
		lim := usage.Limit{HitAt: h.HitAt, ResetAt: h.ResetAt}
		kind, ok := usage.KindOf(lim)
		if !ok || !h.HitAt.After(known[string(kind)].HitAt) {
			continue
		}
		used, ok := usageAtHit(ctx, h.Agent, kind, lim)
		if !ok || used < 1 {
			continue
		}
		known[string(kind)] = daemon.LearnedLimit{Kind: string(kind), Tokens: used, HitAt: h.HitAt}
	}
	state.Learned = nil
	for _, k := range usage.Kinds {
		if l, ok := known[string(k)]; ok {
			state.Learned = append(state.Learned, l)
		}
	}
}

// ratchetLimits raises a learned limit with use between hits: when a current
// window's use passes its kind's learned limit and no recorded hit of that
// kind falls inside the window, the use becomes the learned limit, with the
// hit time it had. A window holding a hit is left to learnLimits, which
// teaches the use at that hit.
func ratchetLimits(ctx *Ctx, common string, state *daemon.State, windows map[usage.Kind]usage.Window) {
	if len(state.Learned) == 0 {
		return
	}
	raiseLearned(state, windows, daemon.ReadLimitHitsIn(limitsDir(ctx, common)))
}

// raiseLearned is ratchetLimits over the hits already read.
func raiseLearned(state *daemon.State, windows map[usage.Kind]usage.Window, hits []daemon.LimitHit) {
	for i, l := range state.Learned {
		w, ok := windows[usage.Kind(l.Kind)]
		if !ok || w.Used <= l.Tokens {
			continue
		}
		held := false
		for _, h := range hits {
			if k, ok := usage.KindOf(usage.Limit{HitAt: h.HitAt, ResetAt: h.ResetAt}); ok && string(k) == l.Kind && w.Contains(h.HitAt) {
				held = true
				break
			}
		}
		if !held {
			state.Learned[i].Tokens = w.Used
		}
	}
}

// usageAtHit is the kind's window use at the hit, read from the token
// records of the agent that saw it (else of every agent that reports them).
func usageAtHit(ctx *Ctx, agent string, kind usage.Kind, l usage.Limit) (int64, bool) {
	names := adapter.Names()
	if _, ok := adapter.Lookup(agent); ok {
		names = []string{agent}
	}
	from := l.HitAt.Add(-2 * kind.Length())
	var recs []adapter.UsageRecord
	for _, n := range names {
		a, _ := adapter.Lookup(n)
		r, ok := a.(adapter.UsageReader)
		if !ok {
			continue
		}
		if ctx.Getenv == nil {
			continue
		}
		got, err := r.UsageRecords(ctx.Getenv, from, l.HitAt.Add(time.Nanosecond))
		if err == nil {
			recs = append(recs, got...)
		}
	}
	return usage.UsedAt(kind, recs, l)
}

// usageLimit is the limit of a window kind: the live one when the kind has
// a live snapshot, else the setting when it is set, else the learned one. A
// window with none has no known limit and never holds. Every caller that
// needs a kind's limit goes through here.
func usageLimit(rows []Setting, learned []daemon.LearnedLimit, live map[usage.Kind]int64, kind usage.Kind) (int64, bool) {
	if n, ok := live[kind]; ok && n > 0 {
		return n, true
	}
	if v := settingValue(rows, "rhiza.daemon.usage."+string(kind)); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			return n, true
		}
	}
	for _, l := range learned {
		if l.Kind == string(kind) && l.Tokens > 0 {
			return l.Tokens, true
		}
	}
	return 0, false
}
