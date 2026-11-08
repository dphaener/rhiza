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
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// watchdogDeps is everything the stalled-step check touches outside the
// journals: tests replace each.
type watchdogDeps struct {
	sessions *session.Store
	// step is st's current step: its kind, stall limit and timeout.
	step   func(st run.State) workflow.Step
	send   func(sid, text string) error
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
	// now is the clock stamping journalled attempts; nil stamps none.
	now func() time.Time
	log func(format string, args ...any)
}

// charge spends one retry of m's budget and journals the recovery tried.
func (d watchdogDeps) charge(m *daemon.WatchdogRun, action string) {
	var at time.Time
	if d.now != nil {
		at = d.now()
	}
	m.Spent++
	m.Attempts = append(m.Attempts, daemon.WatchdogAttempt{At: at, Action: action})
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
		d.charge(cur, fmt.Sprintf("step %s stalled: stopped and resumed with --retry", st.Current))
		cur.Nudged = false
		d.fileBug(pol, j.ID, bugStall, st.Current, fmt.Sprintf("no journal or session progress since %s", at.UTC().Format(time.RFC3339)))
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
		resume: func(id string) error {
			_, err := queueContinue(ctx, common, id, true, false)
			return err
		},
		now: ctx.Now,
		log: func(format string, args ...any) {
			if ctx.Stderr != nil {
				fmt.Fprintf(ctx.Stderr, "rhiza: daemon: "+format+"\n", args...)
			}
		},
	}
	d.queued = func(id string) bool {
		q, err := daemon.ReadContinues(common)
		return err == nil && slices.ContainsFunc(q, func(c daemon.Continue) bool { return c.Run == id })
	}
	d.requests = func() ([]daemon.Request, error) { return daemon.ReadRequests(common) }
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
		return watchdogMergeMain(ctx, t, st.ID)
	}
	d.log = watchdogJournalLog(d.log, state, ctx.Now)
	mem, holds, err := watchdogPasses(pol, all, rs.Live, ctx.Now(), state.Watchdog, d)
	state.Watchdog, state.WatchdogHolds = mem, holds
	return err
}

// watchdogPasses is one poll of the watchdog over all runs: the
// stalled-step pass over the live ones, then the transient-failure and CI
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
	mem, holds, err := watchdogStalled(pol, running, now, prev, d)
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
	return mem, watchdogBudgetHolds(pol, all, mem, holds), err
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
	if st.Result == run.Failed && len(m.Attempts) > 0 && m.Spent >= pol.RetryBudget {
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
		if len(m.Attempts) == 0 || m.Spent < pol.RetryBudget {
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

// transientSignatures are the failure outputs a rerun is expected to clear:
// a go build's $WORK directory removed under it, a git index.lock held by
// another process, and a network or GitHub 5xx error.
var transientSignatures = []*regexp.Regexp{
	regexp.MustCompile(`(?i)go-build\d*/[^\n]*no such file or directory`),
	regexp.MustCompile(`(?i)index\.lock['"]?:? (file exists|exists)|unable to create '[^\n]*index\.lock'`),
	regexp.MustCompile(`(?i)(http|status|error|returned|received)[ :]*(50[0234])\b|\b50[0234] (bad gateway|service unavailable|gateway time-?out|internal server error)`),
	regexp.MustCompile(`(?i)(connection reset by peer|i/o timeout|tls handshake timeout|temporary failure in name resolution)`),
}

// transientSignature names the signature text matches, or "".
func transientSignature(text string) string {
	for _, re := range transientSignatures {
		if m := re.FindString(text); m != "" {
			return strings.TrimSpace(m)
		}
	}
	return ""
}

// watchdogTransient resumes with --retry each failed run of journals whose
// failing step's output matches a transient signature, within the run's
// budget. A failure matching none is left to the owner, as is an attended
// run when the policy is autonomous-only. A run with a live process is never touched. It returns the memory after, and the first error.
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
		sig := transientSignature(failedDetail(st, step))
		if sig == "" {
			continue
		}
		i := slices.IndexFunc(mem, func(w daemon.WatchdogRun) bool { return w.Run == j.ID })
		if i < 0 {
			mem = append(mem, daemon.WatchdogRun{Run: j.ID})
			i = len(mem) - 1
		}
		if mem[i].Spent >= pol.RetryBudget || d.queued != nil && d.queued(j.ID) {
			continue
		}
		if err := d.resume(j.ID); err != nil {
			if first == nil {
				first = fmt.Errorf("resume run %s: %w", j.ID, err)
			}
			continue
		}
		d.charge(&mem[i], fmt.Sprintf("step %s failed on a transient error (%q): resumed with --retry", step, sig))
		d.log("watchdog: run %s step %s failed on a transient error (%q): resumed with --retry (budget %d of %d)", j.ID, step, sig, mem[i].Spent, pol.RetryBudget)
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
// head with gh. A change set with no open PR gives an empty head.
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
