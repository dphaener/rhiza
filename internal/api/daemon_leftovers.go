package api

import (
	"fmt"
	"slices"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/leftovers"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// daemonSweepLeftovers prunes what ended runs left on disk when the
// schedule says so (at start, then hourly): index databases whose worktree
// is gone, and old scratch folders of sessions that belong to no live run,
// in each project whose git common dir is in commons (every project the
// daemon serves). It writes one stderr line with what was freed, only when
// something was. The sweep runs on a goroutine of its own, so a slow disk
// never holds up the poll; none starts while one runs.
func daemonSweepLeftovers(ctx *Ctx, sched *daemon.SweepSchedule, commons []string) {
	if !sched.Start(ctx.Now(), 0, 0, false) {
		return
	}
	c := *ctx
	commons = slices.Clone(commons)
	go func() {
		defer sched.Finish()
		sweepLeftoversNow(&c, commons)
	}()
}

// sweepLeftoversWork does the sweep's disk work and returns what it freed;
// a test replaces it to hold the sweep.
var sweepLeftoversWork = func(ctx *Ctx, commons []string) (index.Pruned, leftovers.Freed) {
	var idx index.Pruned
	if dir, err := index.IndexDir(); err == nil {
		idx = index.PruneOrphans(dir)
	}
	var scr leftovers.Freed
	seen := map[string]bool{}
	for _, common := range commons {
		if seen[common] {
			continue
		}
		seen[common] = true
		f := sweepSessionScratch(ctx, common)
		scr.Folders += f.Folders
		scr.Bytes += f.Bytes
	}
	return idx, scr
}

// sweepLeftoversNow runs one sweep and writes its stderr line.
func sweepLeftoversNow(ctx *Ctx, commons []string) {
	idx, scr := sweepLeftoversWork(ctx, commons)
	if (idx.Files > 0 || scr.Folders > 0) && ctx.Stderr != nil {
		fmt.Fprintf(ctx.Stderr, "rhiza: daemon: leftovers sweep freed %s: %d index databases, %d session scratch folders\n",
			daemon.FormatSize(uint64(idx.Bytes+scr.Bytes)), idx.Files, scr.Folders)
	}
}

// sweepSessionScratch removes old scratch folders of the sessions recorded
// under common whose run has ended or is unknown.
func sweepSessionScratch(ctx *Ctx, common string) leftovers.Freed {
	recs, err := session.NewStore(session.Dir(common)).List()
	if err != nil {
		return leftovers.Freed{}
	}
	runs, err := run.NewStore(run.Dir(common)).List()
	if err != nil {
		return leftovers.Freed{}
	}
	open := map[string]bool{}
	for _, j := range runs {
		if !j.State().Ended {
			open[j.ID] = true
		}
	}
	ss := make([]leftovers.Session, 0, len(recs))
	for _, r := range recs {
		ss = append(ss, leftovers.Session{ID: r.ID, Agent: r.Agent, Run: r.Run})
	}
	getenv := ctxGetenv(ctx)
	return leftovers.Scratch(ss, func(id string) bool { return open[id] },
		func(agent string) []string {
			a, ok := adapter.Lookup(agent)
			if !ok {
				return nil
			}
			return a.ScratchParents(getenv)
		}, nil, ctx.Now())
}
