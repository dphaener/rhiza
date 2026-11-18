package api

// This file holds how a run finds a session that ended while no run
// watched it (a reboot killed its wrapper, so no exit was ever recorded):
// on a rejoin (Rejoin), the session's record and its Rex block; on every
// poll (AgentStatus), the record alone, once its SessionEnd mark is older
// than hookEndGrace. Such a session gets an exit of
// state "lost" written into its record, so exit-is-final holds and every
// reader agrees.

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

var _ run.SpawnWaker = (*stepRunner)(nil)

// Reasons a session is lost, the detail of its lost exit.
const (
	lostHookEnded = "SessionEnd without an exit report"
	lostNoBlock   = "Rex block not found"
)

// hookEndGrace is how long a poll leaves a SessionEnd mark with no exit
// unjudged. Claude's SessionEnd hook runs before its process exits, and the
// wrapper records the exit only after it has waited on the process and
// stopped following it, so a live session that ends normally has the mark
// without an exit for a moment; a poll in that window must not call it
// lost. A rejoin judges the mark at once: it follows a resume, when no
// wrapper the run started is expected to be alive.
const hookEndGrace = time.Minute

// Rejoin is run.Rejoiner: the run rejoins session sid as a resume, before
// its first poll. It reads the record, then asks Rex whether the session's
// block still exists. A Rex error that is not a "not found" answer judges
// nothing: it comes back for the engine to log, and the check is retried on
// the session's later polls until Rex answers.
func (r *stepRunner) Rejoin(sid string) error {
	rec, err := readSession(r.sessions, sid)
	if err != nil {
		return err
	}
	rec, err = r.markLostIfEnded(rec, 0)
	if err != nil || rec.Status.Exited() {
		return err
	}
	r.setBlockPending(rec.ID, true)
	return r.checkBlock(rec)
}

// markLostIfEnded writes a lost exit into rec's record when it has the
// SessionEnd mark, at least grace old, and no exit, and returns the record
// as it stands.
func (r *stepRunner) markLostIfEnded(rec session.Record, grace time.Duration) (session.Record, error) {
	if rec.Status.Exited() || !rec.Status.HookEnded() {
		return rec, nil
	}
	if grace > 0 && r.ctx.Now().Sub(rec.Status.Hook.At) < grace {
		return rec, nil
	}
	return r.markLost(rec.ID, lostHookEnded)
}

// markLost writes an exit of state lost, naming why, into session sid's
// record, unless the session has exited meanwhile: an exit is final.
func (r *stepRunner) markLost(sid, why string) (session.Record, error) {
	rec, err := r.sessions.Update(sid, func(rec *session.Record) {
		if rec.Status.Exited() {
			return
		}
		rec.Status.Exit = &session.Report{State: session.StateLost, Detail: why, At: r.ctx.Now().UTC()}
	})
	if err != nil {
		return rec, err
	}
	r.setBlockPending(rec.ID, false)
	return rec, nil
}

func (r *stepRunner) setBlockPending(sid string, on bool) {
	r.lostMu.Lock()
	defer r.lostMu.Unlock()
	if r.blockPending == nil {
		r.blockPending = map[string]bool{}
	}
	if on {
		r.blockPending[sid] = true
	} else {
		delete(r.blockPending, sid)
	}
}

func (r *stepRunner) isBlockPending(sid string) bool {
	r.lostMu.Lock()
	defer r.lostMu.Unlock()
	return r.blockPending[sid]
}

// checkBlock asks Rex whether rec's block still exists (`rex block
// inspect`). A "not found" answer marks the session lost; an answer that
// the block exists ends the pending check; any other error is returned and
// the check stays pending. A record with no block, or a Runner that cannot
// run rex, has nothing to ask.
func (r *stepRunner) checkBlock(rec session.Record) error {
	if rec.Tab.Block == "" || r.ctx.Exec == nil {
		r.setBlockPending(rec.ID, false)
		return nil
	}
	args := []string{"block", "inspect"}
	if rec.Tab.Session != "" {
		args = append(args, "-s", rec.Tab.Session)
	}
	args = append(args, "--json", rec.Tab.Block)
	if _, err := r.ctx.Exec("rex", args, r.dir(), nil); err != nil {
		if strings.Contains(err.Error(), "not found") {
			_, err := r.markLost(rec.ID, lostNoBlock)
			return err
		}
		return fmt.Errorf("ask Rex whether block %s of session %s exists: %v", rec.Tab.Block, rec.ID, err)
	}
	r.setBlockPending(rec.ID, false)
	return nil
}

// retryBlock is the poll's retry of a rejoin's block check that Rex did not
// answer. Rex still failing is logged and judges nothing; the record is
// returned as it stands.
func (r *stepRunner) retryBlock(rec session.Record) session.Record {
	if rec.Status.Exited() || !r.isBlockPending(rec.ID) {
		return rec
	}
	if err := r.checkBlock(rec); err != nil {
		if r.ctx.Stderr != nil {
			fmt.Fprintf(r.ctx.Stderr, "rejoin: %v\n", err)
		}
		return rec
	}
	if again, err := readSession(r.sessions, rec.ID); err == nil {
		return again
	}
	return rec
}

// Child is run.ChildWatcher: where spawned child session sid stands. A
// record that is gone is reported as such; a child whose SessionEnd mark
// has no exit past the grace is marked lost first, as a poll of a step's
// own session does.
func (r *stepRunner) Child(sid string) (run.ChildState, error) {
	rec, err := r.sessions.Read(sid)
	var nf *session.NotFoundError
	if errors.As(err, &nf) {
		return run.ChildState{Gone: true, Detail: "session record is gone"}, nil
	}
	if err != nil {
		return run.ChildState{}, err
	}
	if rec, err = r.markLostIfEnded(rec, hookEndGrace); err != nil {
		return run.ChildState{}, err
	}
	outcome, detail, done := childOutcome(rec)
	return run.ChildState{Ended: done, Outcome: outcome, Detail: detail}, nil
}

// spawnWakeTimeout is how long WakeParent waits for an idle parent to
// reach its prompt, short so one parent cannot hold up the engine's poll;
// a working parent gets the pointer queued at once.
const spawnWakeTimeout = 30 * time.Second

// WakeParent is run.SpawnWaker: the end of no-wait child sp, just
// recorded, wakes its parent session. A child that said done waits at its
// prompt, so its pane is closed first, as a waiting spawn does; a failed
// close is reported and still wakes the parent. The parent gets a message
// naming the child, its role, leaf, outcome and detail and pointing at
// rhiza run status for the rest, queued while it works.
func (r *stepRunner) WakeParent(sp run.Spawn) error {
	var errs []error
	if rec, err := readSession(r.sessions, sp.Session); err != nil {
		errs = append(errs, err)
	} else if rec.Status.SaidDone() {
		if _, err := closeSessionPane(r.ctx, rec, r.ctx.Root); err != nil {
			errs = append(errs, fmt.Errorf("child session %s said done, but its pane stays open: %v", sp.Session, err))
		}
	}
	if err := r.messageParent(sp.Parent, spawnWakeText(r.run, sp)); err != nil {
		errs = append(errs, fmt.Errorf("message parent session %s: %v", sp.Parent, err))
	}
	return errors.Join(errs...)
}

// messageParent sends live session sid text as a message from the run's
// caller, as session.message does: queued while the session works, else
// typed once it is idle, waiting at most spawnWakeTimeout.
func (r *stepRunner) messageParent(sid, text string) error {
	rec, err := readSession(r.sessions, sid)
	if err != nil {
		return err
	}
	if err := messageable(rec); err != nil {
		return err
	}
	from, err := messageSender(r.ctx)
	if err != nil {
		return err
	}
	dir, rel, err := messageFolder(r.ctx, rec)
	if err != nil {
		return err
	}
	name, err := writeMessage(dir, messageContent(from, rec, rfc3339(r.ctx.Now()), text))
	if err != nil {
		return err
	}
	pointer := messagePointer(from.Handle, rel+"/"+name)
	if queued, err := queuePointer(r.sessions, rec.ID, pointer); err != nil || queued {
		return err
	}
	if _, err := typePointer(r.ctx, r.sessions, rec.ID, pointer, spawnWakeTimeout); err != nil {
		return withNote(err, "the message is in "+filepath.Join(dir, name))
	}
	return nil
}

// spawnWakeText is the message that wakes the parent of no-wait child sp.
func spawnWakeText(runID string, sp run.Spawn) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Your spawned child session %s ended.\n\n", sp.Session)
	fmt.Fprintf(&b, "- Role: %s\n", sp.Role)
	if sp.Leaf != "" {
		fmt.Fprintf(&b, "- Leaf: %s\n", sp.Leaf)
	}
	fmt.Fprintf(&b, "- Outcome: %s\n", sp.Outcome)
	if sp.Detail != "" {
		fmt.Fprintf(&b, "- Detail: %s\n", sp.Detail)
	}
	fmt.Fprintf(&b, "\nRun `rhiza run status %s --json` for the rest.\n", runID)
	return b.String()
}
