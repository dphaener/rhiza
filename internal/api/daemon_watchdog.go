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
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/proctree"
	"github.com/dphaener/rhiza/internal/rootcfg"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// watchdogDeps is everything the stalled-step check touches outside the
// journals: tests replace each.
type watchdogDeps struct {
	sessions *session.Store
	// step is st's current step: its kind, stall limit and timeout.
	step func(st run.State) workflow.Step
	send func(sid, text string) error
	// stop closes the panes of sessions, the step's stall sessions and the
	// parents waiting on their spawns, and stops st's process.
	stop   func(st run.State, sessions []string) error
	resume func(runID string) error
	// queued is whether a continue of the run is already queued: the
	// failure passes leave such a run be, uncharged. nil is never.
	queued func(runID string) bool
	// ci reads the PR of st's change set: its head, its checks and the
	// attempt of each Actions run on the head. rerun reruns the failed jobs
	// of one Actions run of st's PR.
	ci    func(st run.State) (watchdogPR, error)
	rerun func(st run.State, actionsRun string) error
	// merge merges origin/main into the branch of st's PR in the run's
	// worktree and pushes it; on a conflict it aborts and names the files.
	merge func(st run.State) (watchdogMerge, error)
	// requests reads the plan request queue; nil is none.
	requests func() ([]daemon.Request, error)
	// file queues a bug request with the text; nil files none. all is every
	// run's journal of the poll, for reading which requests are still open.
	file func(text string) error
	all  []run.Journal
	// flakes reads the CI stability record's runs; nil is none.
	flakes func() ([]CIRun, error)
	// changeState gives a change set's state and whether the log has it; nil
	// knows none.
	changeState func(change string) (string, bool)
	// warn appends a warning entry with detail to the run's journal; nil
	// journals none.
	warn func(runID, detail string) error
	// now is the clock stamping journalled attempts; nil stamps none.
	now func() time.Time
	// classify reads a failed step's output and says whether a rerun would
	// likely clear it; nil classifies nothing, so every failure is real.
	classify func(output string) FailureVerdict
	// procs reads the host's process table for the runaway pass; nil reads
	// none, so the pass does nothing.
	procs func() ([]proctree.Process, error)
	// stopped journals a stopped entry on runID's step with detail; nil
	// journals none.
	stopped func(runID, step, detail string) error
	log     func(format string, args ...any)
}

// charge spends one retry of m's budget and journals the recovery tried.
func (d watchdogDeps) charge(m *daemon.WatchdogRun, action string) {
	m.Spent++
	d.attempt(m, action)
}

// attempt journals a recovery tried on m's run, charged to no budget.
func (d watchdogDeps) attempt(m *daemon.WatchdogRun, action string) {
	var at time.Time
	if d.now != nil {
		at = d.now()
	}
	m.Attempts = append(m.Attempts, daemon.WatchdogAttempt{At: at, Action: action})
}

// watchdogProgress is the latest sign of life of st's current step: the last
// journal entry or the last report of one of the step's stall sessions
// (stallSessions: its open spawned children counted, a parent waiting on
// one not, so the parent's progress is its child's). ok is false when the
// step is not stalled-able: one of those sessions waits on an answer, said
// it is done, paused at its usage limit or is gone, or the step's fresh
// session waits on a pre-start hold (run.State.StepHeld).
func watchdogProgress(j run.Journal, st run.State, d watchdogDeps) (at time.Time, sids []string, ok bool) {
	if n := len(j.Entries); n > 0 {
		at = j.Entries[n-1].At
	}
	if len(st.Blocks) > 0 {
		return at, nil, false
	}
	if _, held := st.StepHeld(st.Current); held {
		return at, nil, false
	}
	for _, sid := range stallSessions(st) {
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

// stallReason names why st's step counts as stalled, for the run's journal,
// the charged attempt and the hold (decision stall-reason-on-run): each of
// sids, the stall sessions, with its role (a spawned child's role and
// parent, or the step's session) and the time it last reported; with no
// session, the step's last journal entry at at.
func stallReason(st run.State, sids []string, at time.Time, d watchdogDeps) string {
	if len(sids) == 0 {
		return fmt.Sprintf("no session, last journal entry at %s", at.UTC().Format(time.RFC3339))
	}
	var parts []string
	for _, sid := range sids {
		who := sessionRole(st, sid)
		last := "never reported"
		if rec, err := d.sessions.Read(sid); err == nil {
			if rep, _ := session.Resolve(rec); !rep.At.IsZero() {
				last = "last reported at " + rep.At.UTC().Format(time.RFC3339)
			}
		}
		parts = append(parts, fmt.Sprintf("session %s (%s) %s", sid, who, last))
	}
	return strings.Join(parts, "; ")
}

// sessionRole names sid's part in st's current step: a spawned child's
// role and parent, or the step's session.
func sessionRole(st run.State, sid string) string {
	who := "the step's session"
	for _, sp := range st.Spawns {
		if sp.Step == st.Current && !sp.Ended && sp.Session == sid {
			role := sp.Role
			if role == "" {
				role = "spawned"
			}
			who = fmt.Sprintf("role %s, spawned by session %s", role, sp.Parent)
		}
	}
	return who
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
		step := d.step(st)
		stall := step.Stall
		second := pol.StallLimit(step.Kind, stall, step.Timeout)
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
		reason := stallReason(st, sids, at, d)
		if cur.Spent >= pol.RetryBudget {
			holds = append(holds, daemon.WatchdogHold{
				Run: j.ID, Change: st.Change, Step: st.Current, Since: at,
				Detail: fmt.Sprintf("step %s stalled (%s) and the watchdog's retry budget of %d is spent: it stopped and retried %d times", st.Current, reason, pol.RetryBudget, cur.Spent),
			})
			continue
		}
		closing := slices.Clone(sids)
		for _, sid := range spawnWaiting(st) {
			if !slices.Contains(closing, sid) {
				closing = append(closing, sid)
			}
		}
		if err := d.stop(st, closing); err != nil {
			fail(fmt.Errorf("stop run %s: %w", j.ID, err))
			continue
		}
		if err := d.resume(j.ID); err != nil {
			fail(fmt.Errorf("resume run %s: %w", j.ID, err))
			continue
		}
		d.charge(cur, fmt.Sprintf("step %s stalled (%s): stopped and resumed with --retry", st.Current, reason))
		cur.Nudged = false
		if d.warn != nil {
			if err := d.warn(j.ID, fmt.Sprintf("watchdog: step %s stalled: %s; stopped and resumed with --retry (budget %d of %d)", st.Current, reason, cur.Spent, pol.RetryBudget)); err != nil {
				d.log("watchdog: run %s: warning not journalled: %v", j.ID, err)
			}
		}
		d.fileBug(pol, j.ID, bugStall, st.Current, fmt.Sprintf("no journal or session progress since %s: %s", at.UTC().Format(time.RFC3339), reason))
		d.log("watchdog: run %s step %s silent since %s: stopped and resumed with --retry (budget %d of %d)", j.ID, st.Current, at.UTC().Format(time.RFC3339), cur.Spent, pol.RetryBudget)
	}
	return out, holds, first
}

// daemonWatchdog is the poll's watchdog over this repo's runs.
func daemonWatchdog(ctx *Ctx, rs *run.Store, state *daemon.State, pol WatchdogPolicy) error {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil
	}
	all, err := rs.List()
	if err != nil {
		return Errorf("read runs: %v", err)
	}
	d := watchdogDeps{
		sessions: session.NewStore(session.Dir(common)),
		step: func(st run.State) workflow.Step {
			def, err := st.ParseDefinition()
			if err != nil {
				return workflow.Step{Kind: "agent"}
			}
			return def.Steps[st.Current]
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
		classify: ctx.Verdict,
		resume: func(id string) error {
			_, err := queueContinue(ctx, common, id, true, false)
			return err
		},
		warn: func(id, detail string) error {
			return rs.Append(id, run.Entry{Type: run.TypeWarning, At: ctx.Now(), Detail: detail})
		},
		stopped: func(id, step, detail string) error {
			return rs.Append(id, run.Entry{Type: run.TypeStopped, At: ctx.Now(), Step: step, Detail: detail})
		},
		procs: watchdogProcTable(ctx),
		now:   ctx.Now,
		log: func(format string, args ...any) {
			if ctx.Stderr != nil {
				fmt.Fprintf(ctx.Stderr, "rhiza: daemon: "+format+"\n", args...)
			}
		},
	}
	fold := newInboxFold(ctx)
	defer fold.close()
	d.changeState = fold.stateOf
	d.queued = func(id string) bool {
		q, err := daemon.ReadContinues(common)
		return err == nil && slices.ContainsFunc(q, func(c daemon.Continue) bool { return c.Run == id })
	}
	d.requests = func() ([]daemon.Request, error) { return daemon.ReadRequests(common) }
	d.flakes = func() ([]CIRun, error) {
		rec, err := readCIStability(ciStabilityPath(common))
		return rec.Runs, err
	}
	d.file = func(text string) error {
		return daemon.UpdateRequests(common, func(reqs []daemon.Request) ([]daemon.Request, error) {
			rid, err := newRequestID(ctx, reqs)
			if err != nil {
				return nil, err
			}
			return append(reqs, daemon.Request{ID: rid, Kind: daemon.KindBug, Text: text, At: ctx.Now().UTC()}), nil
		})
	}
	d.ci = func(st run.State) (watchdogPR, error) { return watchdogReadPR(ctx, st) }
	d.rerun = func(st run.State, id string) error {
		t, err := watchdogShipTarget(ctx, st)
		if err != nil {
			return err
		}
		_, err = gh(ctx, t.dir(), "run", "rerun", id, "--failed")
		return err
	}
	d.merge = func(st run.State) (watchdogMerge, error) {
		t, err := watchdogShipTarget(ctx, st)
		if err != nil {
			return watchdogMerge{}, err
		}
		return watchdogMergeMain(shipRunEnv(ctx, st, rs.TmpDir(st.ID)), t, st.ID)
	}
	d.log = watchdogJournalLog(d.log, state, ctx.Now)
	mem, holds, err := watchdogPasses(pol, all, rs.Live, ctx.Now(), state.Watchdog, d)
	state.Watchdog, state.WatchdogHolds = mem, holds
	return err
}

// watchdogPasses is one poll of the watchdog over all runs: the runaway
// pass, then the stalled-step pass over the live ones it did not stop, then the transient-failure and CI
// passes over the failed ones. The memory of a run that is not live is
// carried over from prev, so its budget holds across polls; a run gone from
// all is forgotten. It returns the memory and holds after, and the first
// error.
func watchdogPasses(pol WatchdogPolicy, all []run.Journal, live func(id string) bool, now time.Time, prev []daemon.WatchdogRun, d watchdogDeps) ([]daemon.WatchdogRun, []daemon.WatchdogHold, error) {
	d.all = all
	d.ci = watchdogReadOnce(d.ci)
	var running []run.Journal
	for _, j := range all {
		if st := j.State(); !st.Ended && live(j.ID) {
			running = append(running, j)
		}
	}
	prev, stopped, err := watchdogRunaway(pol, running, now, prev, d)
	running = slices.DeleteFunc(running, func(j run.Journal) bool { return slices.Contains(stopped, j.ID) })
	mem, holds, serr := watchdogStalled(pol, running, now, prev, d)
	if err == nil {
		err = serr
	}
	for _, j := range all {
		if slices.ContainsFunc(mem, func(w daemon.WatchdogRun) bool { return w.Run == j.ID }) {
			continue
		}
		if i := slices.IndexFunc(prev, func(w daemon.WatchdogRun) bool { return w.Run == j.ID }); i >= 0 {
			mem = append(mem, prev[i])
		}
	}
	mem, terr := watchdogTransient(pol, all, live, mem, d)
	if err == nil {
		err = terr
	}
	mem, terr = watchdogCIRerun(pol, all, live, mem, d)
	if err == nil {
		err = terr
	}
	if d.requests != nil {
		reqs, rerr := d.requests()
		if rerr == nil {
			mem, rerr = watchdogRequests(pol, reqs, all, live, mem, d)
		}
		if err == nil {
			err = rerr
		}
	}
	if ferr := watchdogFlakyTests(pol, d); err == nil {
		err = ferr
	}
	mem, nholds, terr := watchdogPRNoCI(pol, all, live, mem, d)
	holds = append(holds, nholds...)
	if err == nil {
		err = terr
	}
	mem, nholds, terr = watchdogGreenUnmerged(pol, all, live, mem, d)
	holds = append(holds, nholds...)
	if err == nil {
		err = terr
	}
	holds = watchdogBudgetHolds(pol, all, mem, holds)
	return mem, watchdogLiveHolds(all, holds, d), err
}

// watchdogLiveHolds drops the holds on a run that is abandoned or whose
// change set is withdrawn, rejected or superseded: the owner has settled it,
// and its hold would stay in the inbox with nothing to dismiss it.
func watchdogLiveHolds(all []run.Journal, holds []daemon.WatchdogHold, d watchdogDeps) []daemon.WatchdogHold {
	return slices.DeleteFunc(holds, func(h daemon.WatchdogHold) bool {
		if i := slices.IndexFunc(all, func(j run.Journal) bool { return j.ID == h.Run }); i >= 0 && all[i].State().Abandoned() {
			return true
		}
		if d.changeState == nil || h.Change == "" {
			return false
		}
		state, ok := d.changeState(h.Change)
		return ok && changeset.Dropped(changeset.State(state))
	})
}

// watchdogReadOnce wraps ci so each run's PR is read at most once per poll:
// the CI rerun, no-CI and green-but-unmerged passes share the read. A
// rerun or a push by an earlier pass leaves the later ones the read from
// before it, which they act on no differently. nil stays nil.
func watchdogReadOnce(ci func(st run.State) (watchdogPR, error)) func(st run.State) (watchdogPR, error) {
	if ci == nil {
		return nil
	}
	type read struct {
		pr  watchdogPR
		err error
	}
	seen := map[string]read{}
	return func(st run.State) (watchdogPR, error) {
		if r, ok := seen[st.ID]; ok {
			return r.pr, r.err
		}
		pr, err := ci(st)
		seen[st.ID] = read{pr, err}
		return pr, err
	}
}

// watchdogPRIdle reports whether the PR passes leave j's PR unread: its run
// failed with the watchdog's budget spent on attempts (the budget hold owns
// it and nothing more is resumed), or its PR was found gone, or its head red
// again after the rerun, while the journal still ends where it did then. A
// resume moves the journal on, and the PR is read again.
func watchdogPRIdle(pol WatchdogPolicy, mem []daemon.WatchdogRun, j run.Journal, st run.State) bool {
	i := slices.IndexFunc(mem, func(w daemon.WatchdogRun) bool { return w.Run == j.ID })
	if i < 0 || len(j.Entries) == 0 {
		return false
	}
	m := mem[i]
	if st.Result == run.Failed && len(m.Attempts) > 0 && m.Spent > 0 && m.Spent >= pol.RetryBudget {
		return true
	}
	last := j.Entries[len(j.Entries)-1].At
	return !m.PRGone.IsZero() && m.PRGone.Equal(last) || !m.CIDone.IsZero() && m.CIDone.Equal(last)
}

// markPRGone remembers that j's change set has no open PR, so the PR passes
// skip it until its journal moves on: an ended run's PR is not reopened by
// itself.
func markPRGone(mem []daemon.WatchdogRun, j run.Journal) []daemon.WatchdogRun {
	if len(j.Entries) == 0 {
		return mem
	}
	i := slices.IndexFunc(mem, func(w daemon.WatchdogRun) bool { return w.Run == j.ID })
	if i < 0 {
		mem = append(mem, daemon.WatchdogRun{Run: j.ID})
		i = len(mem) - 1
	}
	mem[i].PRGone = j.Entries[len(j.Entries)-1].At
	return mem
}

// watchdogBudgetHolds adds, for each run that is not done well and whose
// budget is spent, one hold listing every attempt; it replaces any other
// hold on that run, so the owner gets one item. The hold's time is the last
// attempt's, so its inbox item stays the same across polls.
func watchdogBudgetHolds(pol WatchdogPolicy, all []run.Journal, mem []daemon.WatchdogRun, holds []daemon.WatchdogHold) []daemon.WatchdogHold {
	for _, m := range mem {
		if len(m.Attempts) == 0 || m.Spent == 0 || m.Spent < pol.RetryBudget {
			continue
		}
		i := slices.IndexFunc(all, func(j run.Journal) bool { return j.ID == m.Run })
		if i < 0 {
			continue
		}
		st := all[i].State()
		if st.Ended && st.Result != run.Failed {
			continue
		}
		step := st.Current
		if st.Ended {
			if s, err := st.FailedStep(); err == nil {
				step = s
			}
		}
		var tried []string
		for _, a := range m.Attempts {
			tried = append(tried, fmt.Sprintf("%s at %s", a.Action, a.At.UTC().Format(time.RFC3339)))
		}
		holds = slices.DeleteFunc(holds, func(h daemon.WatchdogHold) bool { return h.Run == m.Run })
		holds = append(holds, daemon.WatchdogHold{
			Run: m.Run, Change: st.Change, Step: step, Since: m.Attempts[len(m.Attempts)-1].At,
			Detail: fmt.Sprintf("the watchdog's retry budget of %d is spent and it acts on this run no more; it tried: %s", pol.RetryBudget, strings.Join(tried, "; ")),
		})
	}
	return holds
}

// classifyOnce asks d's classifier about the failed step of j, unless m
// marks the failure Judged at the journal's last entry: one the watchdog
// left to the inbox is not asked about again while the journal still ends
// there. ok is false for such a failure.
func classifyOnce(m *daemon.WatchdogRun, j run.Journal, st run.State, step string, d watchdogDeps) (v FailureVerdict, ok bool) {
	if len(j.Entries) == 0 {
		return FailureVerdict{}, false
	}
	last := j.Entries[len(j.Entries)-1].At
	if !m.Judged.IsZero() && m.Judged.Equal(last) {
		return FailureVerdict{}, false
	}
	v = FailureVerdict{Reason: "no failure classifier"}
	if d.classify != nil {
		v = d.classify(failedDetail(st, step))
	}
	return v, true
}

// classifyTransientOnce is classifyOnce for a pass that retries transient
// failures only: a real verdict, or none, is logged and marked Judged at
// the journal's last entry, and ok is false for it.
func classifyTransientOnce(m *daemon.WatchdogRun, j run.Journal, st run.State, step string, d watchdogDeps) (FailureVerdict, bool) {
	v, ok := classifyOnce(m, j, st, step, d)
	if !ok || v.Transient {
		return v, ok
	}
	m.Judged = j.Entries[len(j.Entries)-1].At
	d.log("watchdog: run %s step %s failure classified real, left to the inbox: %s", j.ID, step, v.Reason)
	return v, false
}

// failedRetryDue reports whether the watchdog under pol will still retry
// j's failure with failedRetries, m being its memory of the run: the run
// failed and is not abandoned, the watchdog watches it, a retry is left,
// and the watchdog has not left this failure to the inbox (m.Judged at the
// journal's last entry). A transient failure within retryBudget is retried
// by that budget instead; this says nothing of it.
func failedRetryDue(pol WatchdogPolicy, m daemon.WatchdogRun, j run.Journal, st run.State) bool {
	if !pol.Enabled || m.FailedRetries >= pol.FailedRetries {
		return false
	}
	if !st.Ended || st.Result != run.Failed || st.Abandoned() || !st.Autonomous && pol.AutonomousOnly {
		return false
	}
	if n := len(j.Entries); n == 0 || m.Judged.Equal(j.Entries[n-1].At) {
		return false
	}
	return true
}

// watchdogTransient resumes with --retry each failed run of journals: one
// whose failing step's output the classifier calls transient within the
// run's retryBudget, one it calls real, or cannot judge, within
// failedRetries (each charged as an attempt in FailedRetries). A failure
// neither budget covers is left to the inbox, each verdict and reason
// logged, as is an attended run when the policy is autonomous-only. A run
// with a live process is never touched, and a real failure whose resume
// errors is left to the inbox. It returns the memory after, and
// the first error.
func watchdogTransient(pol WatchdogPolicy, journals []run.Journal, live func(id string) bool, mem []daemon.WatchdogRun, d watchdogDeps) ([]daemon.WatchdogRun, error) {
	var first error
	for _, j := range journals {
		st := j.State()
		if live(j.ID) || !st.Ended || st.Result != run.Failed || st.Abandoned() {
			continue
		}
		if !st.Autonomous && pol.AutonomousOnly {
			continue
		}
		step, err := st.FailedStep()
		if err != nil {
			continue
		}
		i := slices.IndexFunc(mem, func(w daemon.WatchdogRun) bool { return w.Run == j.ID })
		if i < 0 {
			mem = append(mem, daemon.WatchdogRun{Run: j.ID})
			i = len(mem) - 1
		}
		if mem[i].Spent >= pol.RetryBudget && mem[i].FailedRetries >= pol.FailedRetries || d.queued != nil && d.queued(j.ID) {
			continue
		}
		v, ok := classifyOnce(&mem[i], j, st, step, d)
		if !ok {
			continue
		}
		transient := v.Transient && mem[i].Spent < pol.RetryBudget
		if !transient && (v.Transient || mem[i].FailedRetries >= pol.FailedRetries) {
			mem[i].Judged = j.Entries[len(j.Entries)-1].At
			if v.Transient {
				d.log("watchdog: run %s step %s failed on a transient error (%s) with the retry budget of %d spent, left to the inbox", j.ID, step, v.Reason, pol.RetryBudget)
			} else {
				d.log("watchdog: run %s step %s failure classified real, left to the inbox: %s", j.ID, step, v.Reason)
			}
			continue
		}
		if err := d.resume(j.ID); err != nil {
			if first == nil {
				first = fmt.Errorf("resume run %s: %w", j.ID, err)
			}
			// A real failure the watchdog cannot resume is left to the
			// inbox, or failedRetryDue would hide it every poll.
			if !transient {
				mem[i].Judged = j.Entries[len(j.Entries)-1].At
				d.log("watchdog: run %s step %s failure classified real (%s), its resume failed (%v), left to the inbox", j.ID, step, v.Reason, err)
			}
			continue
		}
		if !transient {
			mem[i].FailedRetries++
			d.attempt(&mem[i], fmt.Sprintf("step %s failed (%s): resumed with --retry", step, v.Reason))
			d.log("watchdog: run %s step %s failure classified real (%s): resumed with --retry (failed retry %d of %d)", j.ID, step, v.Reason, mem[i].FailedRetries, pol.FailedRetries)
			continue
		}
		d.charge(&mem[i], fmt.Sprintf("step %s failed on a transient error (%s): resumed with --retry", step, v.Reason))
		d.log("watchdog: run %s step %s failed on a transient error (%s): resumed with --retry (budget %d of %d)", j.ID, step, v.Reason, mem[i].Spent, pol.RetryBudget)
		d.fileBug(pol, j.ID, bugTransient, step, failedDetail(st, step))
	}
	return mem, first
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

// failedRetryItems drops from items each failed run the watchdog will
// still retry with failedRetries (failedRetryDue) while up says a daemon
// works the project, and names, on a failed run's item, the retries it
// already tried. mem is the watchdog's memory of runs.
func failedRetryItems(items []InboxItem, runs []run.Journal, pol WatchdogPolicy, mem []daemon.WatchdogRun, up bool) []InboxItem {
	out := items[:0:0]
	for _, it := range items {
		if it.Kind != InboxFailed {
			out = append(out, it)
			continue
		}
		mi := slices.IndexFunc(mem, func(w daemon.WatchdogRun) bool { return w.Run == it.Run })
		ri := slices.IndexFunc(runs, func(j run.Journal) bool { return j.ID == it.Run })
		var m daemon.WatchdogRun
		if mi >= 0 {
			m = mem[mi]
		}
		if up && ri >= 0 && failedRetryDue(pol, m, runs[ri], runs[ri].State()) {
			continue
		}
		if m.FailedRetries > 0 {
			var tried []string
			for _, a := range m.Attempts {
				tried = append(tried, fmt.Sprintf("%s at %s", a.Action, a.At.UTC().Format(time.RFC3339)))
			}
			it.Detail = strings.TrimRight(it.Detail, "\n") + fmt.Sprintf("\nthe watchdog retried it %d of %d times (failedRetries) and leaves it to you; it tried: %s", m.FailedRetries, pol.FailedRetries, strings.Join(tried, "; "))
		}
		out = append(out, it)
	}
	return out
}

// inboxFailedRetries is failedRetryItems for the project ctx works in, from
// the daemon's saved state and the project's watchdog policy; items stay
// as they are when no daemon ever worked it.
func inboxFailedRetries(ctx *Ctx, common string, items []InboxItem, runs []run.Journal) []InboxItem {
	st, ok := daemon.ReadState(common)
	if !ok {
		return items
	}
	pol, err := WatchdogPolicyOf(ctx)
	if err != nil {
		pol = defaultWatchdogPolicy()
	}
	_, up := daemon.Running(common)
	if !up {
		if root, ok := rootcfg.Find(ctx.Wd); ok {
			_, up = daemon.RunningRoot(root)
		}
	}
	return failedRetryItems(items, runs, pol, st.Watchdog, up)
}

// watchdogPR is a run's PR as the CI rerun sees it: Attempts is the attempt
// of each Actions run on Head, by run ID (1 for a run never rerun).
type watchdogPR struct {
	Head     string
	Checks   []ShipCheck
	Attempts map[string]int
	// Mergeable is GitHub's verdict: MERGEABLE, CONFLICTING or UNKNOWN.
	Mergeable string
}

// watchdogShipTarget resolves the PR target of st's change set.
func watchdogShipTarget(ctx *Ctx, st run.State) (shipTarget, error) {
	cs, err := id.Parse(st.Change)
	if err != nil {
		return shipTarget{}, err
	}
	c := *ctx
	c.ChangeSet = cs
	return findShipTarget(&c)
}

// watchdogReadPR reads st's PR and the attempts of the Actions runs on its
// head with gh. A change set with no open PR gives an empty head, and so
// does a PR the change set's ship record does not name st's run as the owner
// of (decision watchdog-pr-owned-by-ship-record-run): a record naming
// another run, or none, matches no run.
func watchdogReadPR(ctx *Ctx, st run.State) (watchdogPR, error) {
	t, err := watchdogShipTarget(ctx, st)
	if err != nil {
		return watchdogPR{}, err
	}
	out, err := gh(ctx, t.dir(), "pr", "view", t.branch, "--json", "number,url,state,headRefOid,body,mergeable,statusCheckRollup")
	if noPR(err) {
		return watchdogPR{}, nil
	}
	if err != nil {
		return watchdogPR{}, err
	}
	v, err := parsePRView(out)
	if err != nil {
		return watchdogPR{}, err
	}
	if !strings.EqualFold(v.State, "OPEN") || !prOwnedBy(v.Body, t.cs) {
		return watchdogPR{}, nil
	}
	rec, err := readShipState(t.statePath())
	if err != nil {
		return watchdogPR{}, err
	}
	if rec.Run == "" || rec.Run != st.ID {
		return watchdogPR{}, nil
	}
	pr := watchdogPR{Head: v.Head, Checks: v.Checks, Attempts: map[string]int{}, Mergeable: v.Mergeable}
	out, err = gh(ctx, t.dir(), "run", "list", "--commit", v.Head, "--json", "databaseId,attempt")
	if err != nil {
		return watchdogPR{}, err
	}
	var runs []struct {
		DatabaseID int64 `json:"databaseId"`
		Attempt    int   `json:"attempt"`
	}
	if err := json.Unmarshal(out, &runs); err != nil {
		return watchdogPR{}, fmt.Errorf("gh run list: %v", err)
	}
	for _, r := range runs {
		pr.Attempts[strconv.FormatInt(r.DatabaseID, 10)] = r.Attempt
	}
	return pr, nil
}

// waitingOnCI reports whether the step of st that failed is a ship check on
// the PR's CI: the step that waited for or required green checks, or the
// triage that found no failing test to fix (an infra red). The open-PR probe
// is not.
func waitingOnCI(st run.State, step string) bool {
	def, err := st.ParseDefinition()
	if err != nil {
		return false
	}
	s := def.Steps[step]
	for _, cmd := range []string{s.Run, s.Until} {
		if strings.HasPrefix(cmd, "rhiza ship check") && !strings.Contains(cmd, "--pr-open") && !strings.Contains(cmd, "--main") {
			return true
		}
	}
	return false
}

// rerunSpent reports whether a failing Actions run of pr's head is already
// at an attempt above pol's reruns: its checks settled red after a rerun.
func rerunSpent(pr watchdogPR, pol WatchdogPolicy) bool {
	for _, c := range pr.Checks {
		if passingConclusions[strings.ToUpper(c.Conclusion)] || c.RunID == "" {
			continue
		}
		if pr.Attempts[c.RunID] > pol.CIReruns {
			return true
		}
	}
	return false
}

// watchdogCIRerun handles each failed run that ended waiting on its PR's CI
// with red checks: it reruns the failed jobs once per PR head and, on a
// later poll, resumes the run with --retry once the checks are green. A
// head with an Actions run already at attempt above the policy's reruns
// (the plan step's own rerun leaves one) counts as rerun, and a second red
// on the head is left to the owner. A run whose budget is spent, whose PR
// is gone, or whose head is red again after its rerun is not read again
// while its journal stays as it was (watchdogPRIdle). It returns the memory
// after, and the first error.
func watchdogCIRerun(pol WatchdogPolicy, journals []run.Journal, live func(id string) bool, mem []daemon.WatchdogRun, d watchdogDeps) ([]daemon.WatchdogRun, error) {
	var first error
	fail := func(err error) {
		if first == nil {
			first = err
		}
	}
	if pol.CIReruns <= 0 || d.ci == nil {
		return mem, nil
	}
	for _, j := range journals {
		st := j.State()
		if live(j.ID) || !st.Ended || st.Result != run.Failed || st.Abandoned() {
			continue
		}
		if !st.Autonomous && pol.AutonomousOnly {
			continue
		}
		step, err := st.FailedStep()
		if err != nil || !waitingOnCI(st, step) || watchdogPRIdle(pol, mem, j, st) {
			continue
		}
		i := slices.IndexFunc(mem, func(w daemon.WatchdogRun) bool { return w.Run == j.ID })
		if i >= 0 && mem[i].Spent >= pol.RetryBudget {
			// A rerun could never be followed by the resume.
			continue
		}
		pr, err := d.ci(st)
		if err != nil {
			fail(fmt.Errorf("read the PR of run %s: %w", j.ID, err))
			continue
		}
		if pr.Head == "" {
			mem = markPRGone(mem, j)
			continue
		}
		if i < 0 {
			mem = append(mem, daemon.WatchdogRun{Run: j.ID})
			i = len(mem) - 1
		}
		settled, green, failed := checksVerdict(pr.Checks)
		m := &mem[i]
		if m.CIHead == pr.Head {
			if !m.CIResumed && settled && !green && rerunSpent(pr, pol) {
				// Red again after the rerun: left to the owner, and not
				// read again until the journal moves on.
				m.CIDone = j.Entries[len(j.Entries)-1].At
				continue
			}
			// Rerun by the watchdog: resume once, when the rerun passed.
			if m.CIResumed || !green || m.Spent >= pol.RetryBudget || d.queued != nil && d.queued(j.ID) {
				continue
			}
			if err := d.resume(j.ID); err != nil {
				fail(fmt.Errorf("resume run %s: %w", j.ID, err))
				continue
			}
			d.charge(m, fmt.Sprintf("step %s: CI passed after the rerun on %s: resumed with --retry", step, shortSHA(pr.Head)))
			m.CIResumed = true
			d.log("watchdog: run %s step %s: CI passed after the rerun on %s: resumed with --retry (budget %d of %d)", j.ID, step, shortSHA(pr.Head), m.Spent, pol.RetryBudget)
			continue
		}
		if !settled || green {
			continue
		}
		var runs []string
		seen := map[string]bool{}
		skip := false
		for _, c := range pr.Checks {
			if passingConclusions[strings.ToUpper(c.Conclusion)] {
				continue
			}
			if c.RunID == "" {
				skip = true
				break
			}
			if !seen[c.RunID] {
				seen[c.RunID] = true
				runs = append(runs, c.RunID)
			}
		}
		if rerunSpent(pr, pol) {
			// Already rerun (the plan step's own rerun): a second red on
			// the head is the owner's, and is not read again until the
			// journal moves on.
			m.CIDone = j.Entries[len(j.Entries)-1].At
			continue
		}
		if skip || len(runs) == 0 {
			continue
		}
		// Record before rerunning: a crash between the two costs a missed
		// rerun, never a second one.
		m.CIHead, m.CIResumed = pr.Head, false
		var errs int
		for _, id := range runs {
			if err := d.rerun(st, id); err != nil {
				fail(fmt.Errorf("rerun CI run %s of run %s: %w", id, j.ID, err))
				errs++
			}
		}
		if errs == len(runs) {
			m.CIHead = ""
			continue
		}
		d.log("watchdog: run %s step %s failed on red CI (%s) at %s: reran the failed jobs of %s", j.ID, step, strings.Join(failed, ", "), shortSHA(pr.Head), shipPlural(len(runs), "run", "runs"))
	}
	return mem, first
}
