package api

import (
	"os"
	"path/filepath"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/session"
)

// removeRunScratch removes the scratch folders of the agent sessions that
// run rid started. Call it only once the run has ended: a live run's
// sessions still use theirs. Errors are skipped: this is housekeeping.
func removeRunScratch(ctx *Ctx, sessions *session.Store, rid string) {
	if ctx.DryRun || rid == "" {
		return
	}
	recs, err := sessions.List()
	if err != nil {
		return
	}
	getenv := ctxGetenv(ctx)
	for _, rec := range recs {
		if rec.Run != rid {
			continue
		}
		a, ok := adapter.Lookup(rec.Agent)
		if !ok {
			continue
		}
		for _, g := range a.ScratchGlobs(rec.ID, getenv) {
			matches, err := filepath.Glob(g)
			if err != nil {
				continue
			}
			for _, p := range matches {
				_ = os.RemoveAll(p)
			}
		}
	}
}

// removeWorktreeIndex removes the index database of the worktree at path.
// Call it before the worktree itself goes: the database is named by the
// resolved path. Errors are skipped: this is housekeeping.
func removeWorktreeIndex(ctx *Ctx, path string) {
	if ctx.DryRun || path == "" {
		return
	}
	_ = index.Remove(path)
}
