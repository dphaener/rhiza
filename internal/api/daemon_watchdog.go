package api

// This file is the daemon watchdog's stalled-step recovery (decision
// watchdog-scope-and-defaults): a live run whose step has had no journal
// entry or session report past its limit is nudged at the first limit and,
// if nothing moves by the second, stopped and resumed with --retry against
// the run's own watchdog budget. An owner-attended run only gets a log line
// and an inbox item.

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// watchdogDeps is everything the stalled-step check touches outside the
// journals: tests replace each.
type watchdogDeps struct {
	sessions *session.Store
	// step is the kind of st's current step and its own stall limit.
	step   func(st run.State) (kind string, stall time.Duration)
	send   func(sid, text string) error
	stop   func(st run.State, sessions []string) error
	resume func(runID string) error
	log    func(format string, args ...any)
}

// watchdogProgress is the latest sign of life of st's current step: the last
// journal entry or the last report of one of the step's sessions. ok is
// false when the step is not stalled-able: one of its sessions waits on an
// answer, said it is done, paused at its usage limit or is gone.
func watchdogProgress(j run.Journal, st run.State, d watchdogDeps) (at time.Time, sids []string, ok bool) {
	if n := len(j.Entries); n > 0 {
		at = j.Entries[n-1].At
	}
	if len(st.Blocks) > 0 {
		return at, nil, false
	}
	for _, sid := range stepSessions(st) {
		rec, err := d.sessions.Read(sid)
		if err != nil {
			continue
		}
		if !session.Stallable(rec) {
			return at, nil, false
		}
		rep, _ := session.Resolve(rec)
		if rep.At.After(at) {
			at = rep.At
		}
		sids = append(sids, sid)
	}
	return at, sids, true
}

// watchdogStalled checks each live run of journals at now and returns the
// watchdog's memory and holds after it. A failed nudge, stop or resume is
// not recorded, so the next poll tries again; the first error is returned.
func watchdogStalled(pol WatchdogPolicy, journals []run.Journal, now time.Time, prev []daemon.WatchdogRun, d watchdogDeps) ([]daemon.WatchdogRun, []daemon.WatchdogHold, error) {
	var out []daemon.WatchdogRun
	var holds []daemon.WatchdogHold
	var first error
	fail := func(err error) {
		if first == nil {
			first = err
		}
	}
	for _, j := range journals {
		st := j.State()
		mem := daemon.WatchdogRun{Run: j.ID}
		if i := slices.IndexFunc(prev, func(w daemon.WatchdogRun) bool { return w.Run == j.ID }); i >= 0 {
			mem = prev[i]
		}
		out = append(out, mem)
		at, sids, ok := watchdogProgress(j, st, d)
		cur := &out[len(out)-1]
		if !ok || at.IsZero() {
			continue
		}
		if !at.Equal(cur.Progress) {
			cur.Progress, cur.Nudged = at, false
		}
		kind, stall := d.step(st)
		second := pol.StallLimit(kind, stall)
		first1 := second
		if len(sids) > 0 {
			first1 = pol.AgentNudge
			if stall > 0 {
				first1 = stall / 2
			}
		}
		silent := now.Sub(at)
		if silent <= first1 {
			continue
		}
		if !st.Autonomous && pol.AutonomousOnly {
			if silent > second {
				holds = append(holds, daemon.WatchdogHold{
					Run: j.ID, Change: st.Change, Step: st.Current, Since: at,
					Detail: fmt.Sprintf("step %s has made no progress since %s: the run is yours, the watchdog does not stop or retry it", st.Current, at.UTC().Format(time.RFC3339)),
				})
				d.log("watchdog: run %s step %s stalled since %s: attended, left to the owner", j.ID, st.Current, at.UTC().Format(time.RFC3339))
			}
			continue
		}
		if silent <= second {
			if len(sids) > 0 && !cur.Nudged {
				sent := true
				for _, sid := range sids {
					if err := d.send(sid, daemonNudgeText); err != nil {
						fail(fmt.Errorf("nudge session %s: %w", sid, err))
						sent = false
					}
				}
				if sent {
					cur.Nudged = true
					d.log("watchdog: run %s step %s silent since %s: nudged", j.ID, st.Current, at.UTC().Format(time.RFC3339))
				}
			}
			continue
		}
		if cur.Spent >= pol.RetryBudget {
			holds = append(holds, daemon.WatchdogHold{
				Run: j.ID, Change: st.Change, Step: st.Current, Since: at,
				Detail: fmt.Sprintf("step %s stalled and the watchdog's retry budget of %d is spent: it stopped and retried %d times", st.Current, pol.RetryBudget, cur.Spent),
			})
			continue
		}
		if err := d.stop(st, sids); err != nil {
			fail(fmt.Errorf("stop run %s: %w", j.ID, err))
			continue
		}
		if err := d.resume(j.ID); err != nil {
			fail(fmt.Errorf("resume run %s: %w", j.ID, err))
			continue
		}
		cur.Spent++
		cur.Nudged = false
		d.log("watchdog: run %s step %s silent since %s: stopped and resumed with --retry (budget %d of %d)", j.ID, st.Current, at.UTC().Format(time.RFC3339), cur.Spent, pol.RetryBudget)
	}
	return out, holds, first
}

// daemonWatchdog is the poll's stalled-step pass over the live runs.
func daemonWatchdog(ctx *Ctx, rs *run.Store, state *daemon.State, pol WatchdogPolicy) error {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil
	}
	all, err := rs.List()
	if err != nil {
		return Errorf("read runs: %v", err)
	}
	var live []run.Journal
	for _, j := range all {
		if st := j.State(); !st.Ended && rs.Live(j.ID) {
			live = append(live, j)
		}
	}
	d := watchdogDeps{
		sessions: session.NewStore(session.Dir(common)),
		step: func(st run.State) (string, time.Duration) {
			def, err := st.ParseDefinition()
			if err != nil {
				return "agent", 0
			}
			s := def.Steps[st.Current]
			return s.Kind, s.Stall
		},
		send: func(sid, text string) error { return daemonSend(ctx, sid, text) },
		stop: func(st run.State, sids []string) error {
			for _, sid := range sids {
				if rec, err := session.NewStore(session.Dir(common)).Read(sid); err == nil {
					if _, err := closeSessionPane(ctx, rec, ctx.Root); err != nil {
						return err
					}
				}
			}
			return terminateRun(st)
		},
		resume: func(id string) error {
			_, err := queueContinue(ctx, common, id, true, false)
			return err
		},
		log: func(format string, args ...any) {
			if ctx.Stderr != nil {
				fmt.Fprintf(ctx.Stderr, "rhiza: daemon: "+format+"\n", args...)
			}
		},
	}
	mem, holds, err := watchdogStalled(pol, live, ctx.Now(), state.Watchdog, d)
	state.Watchdog, state.WatchdogHolds = mem, holds
	return err
}

// terminateRun signals the process that drives st on this host to stop; a
// run with no process here is already stopped.
func terminateRun(st run.State) error {
	host, _ := os.Hostname()
	p := st.Process
	if p == nil || p.Pid <= 0 || p.Host != host {
		return nil
	}
	return run.Terminate(p.Pid)
}

// watchdogHoldItems are the stalled items for holds.
func watchdogHoldItems(holds []daemon.WatchdogHold) []InboxItem {
	var out []InboxItem
	for _, h := range holds {
		at := h.Since.UTC().Format(time.RFC3339)
		sum := sha256.Sum256([]byte(strings.Join([]string{InboxStalled, "watchdog", h.Run, at}, "\x00")))
		out = append(out, InboxItem{
			ID: hex.EncodeToString(sum[:])[:7], Kind: InboxStalled, At: at,
			Run: h.Run, Change: h.Change, Step: h.Step, Detail: h.Detail,
			Next: "rhiza run status " + h.Run,
		})
	}
	return out
}

// inboxWatchdog is the watchdog's held items from the daemon's saved state.
func inboxWatchdog(ctx *Ctx) []InboxItem {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil
	}
	st, ok := daemon.ReadState(common)
	if !ok {
		return nil
	}
	return watchdogHoldItems(st.WatchdogHolds)
}
