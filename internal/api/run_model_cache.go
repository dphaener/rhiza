package api

// This file holds the run driver's reads of the model through the index. A
// driver polls its sessions every few seconds for as long as a run lives, so
// what a poll asks of the model (the leaves, an answer to a question, a
// verifier's Evidence) must cost nothing while the log has not changed: the
// runner keeps the run worktree's index open, catches it up with a stat of
// the log's folders, and keeps each answer until the index's stamp moves.
// Without an index (ctx.IndexPath empty) the answers are read from the log,
// as before.

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
)

// kept is an answer kept, with the change set it was read for (zero for an
// answer that reads the whole log).
type kept struct {
	v  any
	cs id.ID
}

// runModel is the runner's open index of the run's worktree, and the
// answers read from it at stamp.
type runModel struct {
	root  string
	x     *index.Index
	stamp string
	vals  map[string]kept
	// sigs is, for each change set the answers read, the index's
	// signature of what it reads when they were computed.
	sigs map[id.ID]string
	// caught sums what the catch-ups did since a test last cleared it.
	caught index.CatchUp
	// ix is the write view over x the answers are computed with.
	ix *writeIndex
	// checked is when the index was last caught up.
	checked time.Time
	// computed counts the answers computed rather than kept.
	computed int
}

// worktreeRoot is the Rhiza root of the run's worktree, as worktreeWorkspace
// opens it. The path prefix is asked of git once.
func (r *stepRunner) worktreeRoot() (string, error) {
	if r.wtRoot != "" {
		return r.wtRoot, nil
	}
	prefix, err := r.ctx.Git("rev-parse", "--show-prefix")
	if err != nil {
		return "", err
	}
	root := filepath.Join(r.worktree, strings.TrimSpace(prefix))
	r.wtRoot = root
	return root, nil
}

// closeModel closes the runner's open index.
func (r *stepRunner) closeModel() {
	r.modelMu.Lock()
	defer r.modelMu.Unlock()
	if r.rm != nil && r.rm.x != nil {
		r.rm.x.Close()
	}
	r.rm = nil
}

// readModel returns the answer for key over the run worktree's model:
// the one kept at the index's stamp, else compute's, which is given a
// read-only workspace over the index. Nothing is computed, decoded or
// folded while the log is as it was. An error is not kept. With no index
// to read, compute runs over the log every time.
func (r *stepRunner) readModel(key string, cs id.ID, compute func(w *workspace) (any, error)) (any, error) {
	root, err := r.worktreeRoot()
	if err != nil {
		return nil, err
	}
	w, err := openWorkspace(root, r.ctx.Now, r.ctx.IDs, true)
	if err != nil {
		return nil, err
	}
	w.readLog = r.ctx.readLog
	if r.ctx.IndexPath == "" {
		return compute(w)
	}
	r.modelMu.Lock()
	defer r.modelMu.Unlock()
	rm, err := r.catchUpModel(root)
	if err != nil {
		return nil, err
	}
	if k, ok := rm.vals[key]; ok {
		return k.v, nil
	}
	w.ix = rm.ix
	v, err := compute(w)
	if err != nil {
		return nil, err
	}
	if ferr := w.ix.failure(); ferr != nil {
		rm.x.Close()
		r.rm = nil
		return nil, ferr
	}
	rm.vals[key] = kept{v, cs}
	rm.note(cs)
	rm.computed++
	return v, nil
}

// catchUpModel opens the index of the worktree at root, or catches the one
// open up, and drops the answers kept when its stamp has moved.
func (r *stepRunner) catchUpModel(root string) (*runModel, error) {
	rm := r.rm
	if rm != nil && rm.root != root {
		rm.x.Close()
		rm = nil
	}
	if rm == nil {
		path := r.ctx.IndexPath
		if r.ctx.Root != root {
			p, err := index.DefaultPath(root)
			if err != nil {
				return nil, Errorf("%s", err.Error())
			}
			path = p
		}
		x, err := index.Open(root, path)
		if err != nil {
			return nil, Errorf("%s", err.Error())
		}
		rm = &runModel{root: root, x: x, vals: map[string]kept{}, sigs: map[id.ID]string{}}
		rm.ix = &writeIndex{ctx: r.ctx, x: x, models: map[id.ID]Model{}, relevant: map[id.ID][]id.ID{}}
		r.rm = rm
	} else if r.refreshDue(rm) {
		if err := rm.x.Refresh(); err != nil {
			rm.x.Close()
			r.rm = nil
			return nil, Errorf("%s", err.Error())
		}
	} else {
		return rm, nil
	}
	rm.checked = time.Now()
	c := rm.x.LastCatchUp()
	rm.caught.Fresh = rm.caught.Fresh && c.Fresh
	rm.caught.Rebuilt = rm.caught.Rebuilt || c.Rebuilt
	rm.caught.Read += c.Read
	rm.caught.Decoded += c.Decoded
	rm.caught.Folded += c.Folded
	stamp, err := rm.x.Stamp()
	if err != nil {
		return nil, Errorf("%s", err.Error())
	}
	if stamp != rm.stamp {
		rm.revalidate()
		rm.stamp = stamp
	}
	return rm, nil
}

// note records the signature of what the answers for change set cs read, if
// it is not known yet.
func (rm *runModel) note(cs id.ID) {
	if cs.IsZero() {
		return
	}
	if _, ok := rm.sigs[cs]; ok {
		return
	}
	rel := rm.ix.relevant[cs]
	if rel == nil {
		return
	}
	if sig, err := rm.x.Signature(rel); err == nil {
		rm.sigs[cs] = sig
	}
}

// revalidate runs when the log has changed: it keeps the answers whose
// change sets read the same things as when they were computed (an entry in
// another change set that is in flight changes nothing they read) and drops
// the rest.
func (rm *runModel) revalidate() {
	stale := map[id.ID]bool{}
	for cs, was := range rm.sigs {
		rel := rm.ix.relevant[cs]
		now, err := rm.x.Signature(rel)
		if rel == nil || err != nil || now != was {
			stale[cs] = true
			delete(rm.sigs, cs)
		}
	}
	for key, k := range rm.vals {
		if _, ok := rm.sigs[k.cs]; !ok || k.cs.IsZero() || stale[k.cs] {
			delete(rm.vals, key)
		}
	}
	for cs := range rm.ix.models {
		if _, ok := rm.sigs[cs]; !ok {
			delete(rm.ix.models, cs)
			delete(rm.ix.relevant, cs)
		}
	}
}

// pollDebounce is how long after a catch-up the next calls of a poll take
// the index as it is: a driver's tick asks several things of the model in a
// row, and stats the log's folders once. A write of the runner's own
// (worktreeWorkspace with write) lifts it for writeSettle, so what it
// wrote is read back at once.
const (
	pollDebounce = 50 * time.Millisecond
	writeSettle  = 2 * time.Second
)

// refreshDue reports whether the index must be caught up before it
// answers. A negative r.debounce (tests) always catches up.
func (r *stepRunner) refreshDue(rm *runModel) bool {
	d := r.debounce
	if d == 0 {
		d = pollDebounce
	}
	now := time.Now()
	return d < 0 || rm.checked.IsZero() || now.Sub(rm.checked) >= d || now.Before(r.writeUntil)
}

// wrote notes that the runner is about to write the model.
func (r *stepRunner) wrote() {
	r.modelMu.Lock()
	r.writeUntil = time.Now().Add(writeSettle)
	r.modelMu.Unlock()
}
