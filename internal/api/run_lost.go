package api

// This file holds how a run finds a session that ended while no run
// watched it (a reboot killed its wrapper, so no exit was ever recorded):
// on a rejoin (Rejoin), the session's record and its Rex block; on every
// poll (AgentStatus), the record alone. Such a session gets an exit of
// state "lost" written into its record, so exit-is-final holds and every
// reader agrees.

import (
	"fmt"
	"strings"

	"github.com/dphaener/rhiza/internal/session"
)

// Reasons a session is lost, the detail of its lost exit.
const (
	lostHookEnded = "SessionEnd without an exit report"
	lostNoBlock   = "Rex block not found"
)

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
	rec, err = r.markLostIfEnded(rec)
	if err != nil || rec.Status.Exited() {
		return err
	}
	r.setBlockPending(rec.ID, true)
	return r.checkBlock(rec)
}

// markLostIfEnded writes a lost exit into rec's record when it has the
// SessionEnd mark and no exit, and returns the record as it stands.
func (r *stepRunner) markLostIfEnded(rec session.Record) (session.Record, error) {
	if rec.Status.Exited() || !rec.Status.HookEnded() {
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
