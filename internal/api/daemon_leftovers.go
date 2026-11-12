package api

import (
	"fmt"

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
// something was.
func daemonSweepLeftovers(ctx *Ctx, sched *daemon.SweepSchedule, commons []string) {
	now := ctx.Now()
	if !sched.Due(now, 0, 0, false) {
		return
	}
	sched.Done(now)
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
		func(s leftovers.Session) []string {
			a, ok := adapter.Lookup(s.Agent)
			if !ok {
				return nil
			}
			return a.ScratchGlobs(s.ID, getenv)
		}, ctx.Now())
}
