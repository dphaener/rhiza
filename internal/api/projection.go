package api

// This file holds the projection: the one interface the model is read
// through when a call has an index (ctx.IndexPath set). It answers the
// change sets with their states and folds, a change set's entries, the
// stored problems and the model of a view (the current change set's while
// it is in flight, else as-is), each with the dry run's pending entries
// laid over it. Its one implementation is the index with the pending
// entries overlaid (writeIndex, write_indexed.go). With ctx.IndexPath empty
// there is no projection and each read keeps its log path.

import (
	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
)

// projection is the read side of the model over the index. Every read
// takes the dry run's pending entries (nil for none) and lays them over
// what is stored. A read that meets an index error answers nothing and
// keeps the error for failure.
type projection interface {
	// changeSetIDs is the IDs of the stored and pending change sets,
	// sorted.
	changeSetIDs(pending []entry.Entry) []id.ID
	// changeSets folds every change set: the stored folds, with those the
	// pending entries touch folded again.
	changeSets(pending []entry.Entry) []changeset.ChangeSet
	// entries is change set cs's stored and pending entries, sorted by ID.
	entries(pending []entry.Entry, cs id.ID) []entry.Entry
	// problems is the problems the index stores (log, fold, as-is,
	// to-be); pending entries are not folded into them.
	problems() []index.Problem
	// model is the view of change set current (zero for as-is): its
	// what-if while it is in flight, else as-is.
	model(pending []entry.Entry, current id.ID) Model
	// failure is the first index error a read met, nil for none.
	failure() error
	// wrote notes that an entry was written to disk: the next read
	// catches up with it.
	wrote()
	// close releases the index.
	close()
}

// projectionSlot is a workspace's projection, opened on first use. Copies
// of a workspace (a dry-run copy in apply) share it, so the index is
// opened at most once per call and closed once.
type projectionSlot struct {
	open func() (projection, error)
	p    projection
}

// get opens the projection on the first call and returns it after. An
// index that cannot be opened gives a projection that answers nothing and
// fails with the open's error.
func (s *projectionSlot) get() projection {
	if s.p == nil && s.open != nil {
		open := s.open
		s.open = nil
		p, err := open()
		if err != nil {
			p = failedProjection{err: err}
		}
		s.p = p
	}
	return s.p
}

func (s *projectionSlot) close() {
	if s != nil && s.p != nil {
		s.p.close()
		s.p = nil
	}
	if s != nil {
		s.open = nil
	}
}

// openProjection opens the projection of ctx's index. memo keeps each
// view's model for the life of the projection (no pending entries asked);
// a caller that holds the projection across writes does not ask for it.
func openProjection(ctx *Ctx, memo bool) (projection, error) {
	ix, err := openWriteIndex(ctx)
	if err != nil {
		return nil, err
	}
	if memo {
		ix.models = map[id.ID]Model{}
	}
	return ix, nil
}

// useProjection gives w a projection over ctx's index, opened when a read
// first asks for it; with ctx.IndexPath empty it gives none. The caller
// closes it with closeProjection.
func (w *workspace) useProjection(ctx *Ctx) {
	if ctx.IndexPath == "" {
		w.proj = nil
		return
	}
	w.proj = &projectionSlot{open: func() (projection, error) { return openProjection(ctx, false) }}
}

// setProjection gives w the open projection p.
func (w *workspace) setProjection(p projection) {
	w.proj = &projectionSlot{p: p}
}

// projection is w's projection, opened on first use; nil when w reads the
// log.
func (w *workspace) projection() projection {
	if w.proj == nil {
		return nil
	}
	return w.proj.get()
}

// closeProjection closes w's projection if it was opened.
func (w *workspace) closeProjection() {
	w.proj.close()
}

// projectionFailure is the error w's projection met, nil for none or for a
// projection not opened.
func (w *workspace) projectionFailure() error {
	if w.proj == nil || w.proj.p == nil {
		return nil
	}
	return w.proj.p.failure()
}

// failedProjection is a projection whose index did not open: it answers
// nothing and fails with err.
type failedProjection struct{ err error }

func (f failedProjection) changeSetIDs([]entry.Entry) []id.ID             { return nil }
func (f failedProjection) changeSets([]entry.Entry) []changeset.ChangeSet { return nil }
func (f failedProjection) entries([]entry.Entry, id.ID) []entry.Entry     { return nil }
func (f failedProjection) problems() []index.Problem                      { return nil }
func (f failedProjection) model([]entry.Entry, id.ID) Model {
	return &writeModel{ix: &writeIndex{err: f.err}}
}
func (f failedProjection) failure() error { return f.err }
func (f failedProjection) wrote()         {}
func (f failedProjection) close()         {}
