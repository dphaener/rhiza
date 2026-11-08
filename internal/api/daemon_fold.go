package api

// This file holds the daemon's fold cache: what a poll derives from the
// folded log, kept while the log, the git refs and the current change set
// are unchanged. The clock- and liveness-dependent steps of a poll never
// live here; they run every poll against what the cache holds.

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"slices"
	"sync/atomic"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/policy"
)

// planReadCount and planFoldCount count the plan reads and the in-flight
// change sets folded by them, process-wide; tests read them to see what a
// daemon poll skipped.
var planReadCount, planFoldCount atomic.Int64

// daemonMove is one lifecycle move the events pass may run its command for.
type daemonMove struct {
	ev lifecycleEvent
	id id.ID
}

// daemonFolded is what one fold of the log gave the daemon.
type daemonFolded struct {
	res       planResult
	cands     []string
	planSkips []daemon.Skip
	sets      map[id.ID]changeset.ChangeSet
	overlaps  []InboxItem
	moves     []daemonMove
}

// changeSet is the folded change set cs, false if the log has none.
func (f *daemonFolded) changeSet(cs id.ID) (changeset.ChangeSet, bool) {
	c, ok := f.sets[cs]
	return c, ok
}

// stateOf gives the state of a change set by ID, as changeStateOf does.
func (f *daemonFolded) stateOf(change string) (string, bool) {
	cs, err := id.Parse(change)
	if err != nil {
		return "", false
	}
	c, ok := f.sets[cs]
	return string(c.State), ok
}

// daemonFoldCache keeps the last fold with the stamp it was taken under.
type daemonFoldCache struct {
	stamp  string
	folded *daemonFolded
}

// invalidate drops the cached fold, as a write of the daemon's own does.
func (c *daemonFoldCache) invalidate() { c.folded, c.stamp = nil, "" }

// daemonStamp keys the cache: the log's stat fingerprint, the git refs, the
// daemon's settings and the change set the reads default to. It is "" when the log's stamp cannot
// be trusted (a stat failed, or a folder changed in the last clock tick),
// which never matches a cached one.
func daemonStamp(ctx *Ctx) string {
	if ctx.DryRun {
		return ""
	}
	fp, racy := index.LogStamp(ctx.Root)
	if fp == "" || racy {
		return ""
	}
	h := sha256.New()
	h.Write([]byte(fp))
	h.Write([]byte{0})
	if ctx.Git != nil {
		if out, err := ctx.Git("for-each-ref", "--format=%(objectname) %(refname)"); err == nil {
			h.Write([]byte(out))
		}
	}
	h.Write([]byte{0})
	if ctx.Git != nil {
		// The daemon's settings live in git config.
		if out, err := ctx.Git("config", "--get-regexp", `^rhiza\.`); err == nil {
			h.Write([]byte(out))
		}
	}
	h.Write([]byte{0})
	if data, err := os.ReadFile(filepath.Join(ctx.Root, filepath.FromSlash(policy.Path))); err == nil {
		h.Write(data)
	}
	h.Write([]byte{0})
	h.Write([]byte(ctx.ChangeSet.String() + "|" + ctx.ChangeRef))
	if w, err := openCtxWorkspace(ctx); err == nil {
		if data, err := os.ReadFile(w.currentPath()); err == nil {
			h.Write(data)
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}

// daemonFoldOf returns the fold of the log for this poll: the cached one
// when the stamp is unchanged, else a new one, which it keeps.
func daemonFoldOf(ctx *Ctx) (*daemonFolded, error) {
	if ctx.daemonFold == nil {
		ctx.daemonFold = &daemonFoldCache{}
	}
	c := ctx.daemonFold
	stamp := daemonStamp(ctx)
	if stamp != "" && c.folded != nil && c.stamp == stamp {
		return c.folded, nil
	}
	c.invalidate()
	f, err := daemonFold(ctx)
	if err != nil {
		return nil, err
	}
	if stamp != "" {
		c.folded, c.stamp = f, stamp
	}
	return f, nil
}

// daemonFold folds the log and derives what the poll passes need of it.
func daemonFold(ctx *Ctx) (*daemonFolded, error) {
	res, _, _, err := planRead(ctx, "daemon", "")
	if err != nil {
		return nil, err
	}
	ws, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	l := ws.log()
	list, _ := changeset.FoldLog(l)
	f := &daemonFolded{res: res, sets: make(map[id.ID]changeset.ChangeSet, len(list))}
	for _, c := range list {
		f.sets[c.ID] = c
	}
	sched := daemonSchedule{Since: func(cs id.ID) time.Time { return acceptedSince(f.sets[cs]) }}
	if pol, err := policy.Read(os.DirFS(ctx.Root)); err != nil {
		sched.Detail = []string{err.Error()}
	} else {
		sched.Backlog = pol.Backlog
	}
	f.cands, f.planSkips = daemonCandidates(res, func(cs id.ID) bool {
		c, ok := f.sets[cs]
		return ok && len(changeLeaves(newModelSets(list, cs), c.Ops)) > 0
	}, sched)
	m0 := newModelSets(list, id.ID{})
	if ovs := gmodel.Overlaps(list); len(ovs) > 0 {
		f.overlaps = overlapItems(ovs, m0)
	}
	for _, c := range list {
		for _, mv := range c.Lifecycle {
			name, ok := eventNames[mv.Type]
			if !ok {
				continue
			}
			f.moves = append(f.moves, daemonMove{id: mv.Entry, ev: lifecycleEvent{
				Event: name, Change: c.ID.String(), Entry: mv.Entry.String(),
				Title: c.Open.DisplayTitle(), Ref: c.Open.Ref, At: rfc3339(mv.At), Actor: feedActor(m0, mv.Actor),
			}})
		}
	}
	slices.SortFunc(f.moves, func(a, b daemonMove) int { return id.Compare(a.id, b.id) })
	return f, nil
}

// acceptedSince is when change set c was last accepted or reopened.
func acceptedSince(c changeset.ChangeSet) (at time.Time) {
	for _, m := range c.Lifecycle {
		if m.Type == entry.TypeAccept || m.Type == entry.TypeReopen {
			at = m.At
		}
	}
	return at
}
