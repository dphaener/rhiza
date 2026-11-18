package api

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// TestWatchdogStalled drives the stalled-step recovery with a fake clock,
// fake sessions and a fake run driver.
func TestWatchdogStalled(t *testing.T) {
	t.Parallel()
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	sessions := session.NewStore(t.TempDir())
	report := func(id string, auto bool, at time.Time) {
		rec := session.Record{ID: id, Autonomous: auto}
		rec.Status.Hook = &session.Report{State: session.StateWorking, At: at}
		if err := sessions.Write(rec); err != nil {
			t.Fatal(err)
		}
	}
	journal := func(id string, auto bool, sid string) run.Journal {
		return run.Journal{ID: id, Entries: []run.Entry{
			{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: auto}},
			{Type: run.TypeEntered, At: t0, Step: "build", Visit: 1},
			{Type: run.TypeSession, At: t0, Step: "build", Session: sid},
		}}
	}
	var log, sent, stopped, resumed []string
	d := watchdogDeps{
		sessions: sessions,
		step:     func(run.State) workflow.Step { return workflow.Step{Kind: "agent"} },
		send:     func(sid, _ string) error { sent = append(sent, sid); return nil },
		stop:     func(st run.State, _ []string) error { stopped = append(stopped, st.ID); return nil },
		classify: fakeClassify,
		resume:   func(id string) error { resumed = append(resumed, id); return nil },
		log:      func(f string, a ...any) { log = append(log, strings.TrimSpace(fmt.Sprintf(f, a...))) },
	}
	pol := watchdogPolicyOf(nil)
	_ = pol
	pol = WatchdogPolicy{Enabled: true, AgentNudge: 20 * time.Minute, AgentRetry: 40 * time.Minute, CommandLimit: 30 * time.Minute, WaitLimit: 90 * time.Minute, RetryBudget: 1, AutonomousOnly: true}

	report("auto", true, t0)
	report("busy", true, t0)
	report("owner", false, t0)
	js := []run.Journal{journal("A", true, "auto"), journal("B", true, "busy"), journal("O", false, "owner")}
	var mem []daemon.WatchdogRun
	var holds []daemon.WatchdogHold
	poll := func(now time.Time) {
		t.Helper()
		var err error
		if mem, holds, err = watchdogStalled(pol, js, now, mem, d); err != nil {
			t.Fatal(err)
		}
	}

	poll(t0.Add(19 * time.Minute))
	if len(sent)+len(stopped) != 0 {
		t.Fatalf("before the first limit: sent %v stopped %v", sent, stopped)
	}
	poll(t0.Add(21 * time.Minute))
	if strings.Join(sent, ",") != "auto,busy" || len(stopped) != 0 {
		t.Fatalf("first limit: sent %v stopped %v", sent, stopped)
	}
	// busy makes progress in between; auto and the owner's run stay silent.
	report("busy", true, t0.Add(30*time.Minute))
	poll(t0.Add(41 * time.Minute))
	if strings.Join(stopped, ",") != "A" || strings.Join(resumed, ",") != "A" {
		t.Fatalf("second limit: stopped %v resumed %v", stopped, resumed)
	}
	if len(holds) != 1 || holds[0].Run != "O" {
		t.Fatalf("holds %+v, want the attended run's", holds)
	}
	var lines int
	for _, l := range log {
		if strings.Contains(l, "run A") && strings.Contains(l, "stopped and resumed") {
			lines++
		}
	}
	if lines != 1 {
		t.Errorf("want one stop line for A, log %q", log)
	}
	for _, m := range mem {
		if m.Run == "A" && m.Spent != 1 {
			t.Errorf("A spent %d, want 1", m.Spent)
		}
		if m.Run == "B" && m.Spent != 0 {
			t.Errorf("B spent %d, want 0", m.Spent)
		}
	}
	// The budget is spent: the next stall of A is held, not retried. B, silent
	// again since its progress, is retried once, within its own budget.
	poll(t0.Add(90 * time.Minute))
	if strings.Join(stopped, ",") != "A,B" {
		t.Fatalf("over budget: stopped %v", stopped)
	}
	if !slicesHoldFor(holds, "A") || !slicesHoldFor(holds, "O") {
		t.Fatalf("holds %+v, want A and O", holds)
	}
	for _, id := range []string{"O"} {
		for _, s := range append(append([]string{}, stopped...), resumed...) {
			if s == id {
				t.Errorf("attended run %s was stopped or resumed", id)
			}
		}
	}
	if items := watchdogHoldItems(holds); len(items) != len(holds) || items[0].Kind != InboxStalled {
		t.Errorf("hold items %+v", items)
	}
}

// TestWatchdogStalledSpawnedChild replays the 2026-11-16 stall: a conduct
// run whose conductor spawned a fixer that went idle without saying done.
// The fixer is nudged at the first limit, not the conductor blocked in its
// spawn command; past the second the run is stopped with both panes closed
// and resumed, charged once. A fixer that said done, or whose spawn ended,
// is neither nudged nor counted, and a conductor waiting on a fixer that
// reports is left alone however old its own report.
func TestWatchdogStalledSpawnedChild(t *testing.T) {
	t.Parallel()
	t0 := time.Date(2026, 11, 16, 0, 0, 0, 0, time.UTC)
	pol := WatchdogPolicy{Enabled: true, AgentNudge: 20 * time.Minute, AgentRetry: 40 * time.Minute, CommandLimit: 30 * time.Minute, WaitLimit: 90 * time.Minute, RetryBudget: 1, AutonomousOnly: true}
	type fixture struct {
		sessions               *session.Store
		sent, stopped, resumed []string
		d                      watchdogDeps
	}
	setup := func(t *testing.T) *fixture {
		f := &fixture{sessions: session.NewStore(t.TempDir())}
		f.d = watchdogDeps{
			sessions: f.sessions,
			step:     func(run.State) workflow.Step { return workflow.Step{Kind: "agent"} },
			send:     func(sid, _ string) error { f.sent = append(f.sent, sid); return nil },
			stop: func(st run.State, sids []string) error {
				f.stopped = append(f.stopped, st.ID+":"+strings.Join(sids, ","))
				return nil
			},
			resume: func(id string) error { f.resumed = append(f.resumed, id); return nil },
			log:    func(string, ...any) {},
		}
		return f
	}
	// The conductor works in its spawn command, its last report two hours
	// old; the fixer is autonomous and idle by its hook.
	report := func(t *testing.T, s *session.Store, id, state string, at time.Time, done bool) {
		t.Helper()
		rec := session.Record{ID: id, Autonomous: true}
		rec.Status.Hook = &session.Report{State: state, At: at}
		if done {
			rec.Status.Done = &session.Report{State: session.StateDone, At: at}
		}
		if err := s.Write(rec); err != nil {
			t.Fatal(err)
		}
	}
	conduct := func(ended bool) []run.Journal {
		es := []run.Entry{
			{Type: run.TypeStarted, At: t0.Add(-3 * time.Hour), Started: &run.Started{Autonomous: true}},
			{Type: run.TypeEntered, At: t0.Add(-3 * time.Hour), Step: "conduct", Visit: 1},
			{Type: run.TypeSession, At: t0.Add(-3 * time.Hour), Step: "conduct", Session: "conductor"},
			{Type: run.TypeSpawn, At: t0, Step: "conduct", Parent: "conductor", Session: "fixer", Role: "fixer"},
		}
		if ended {
			es = append(es, run.Entry{Type: run.TypeSpawnEnded, At: t0, Step: "conduct", Parent: "conductor", Session: "fixer", Outcome: workflow.OutcomePass})
		}
		return []run.Journal{{ID: "R", Entries: es}}
	}
	poll := func(t *testing.T, f *fixture, js []run.Journal, mem []daemon.WatchdogRun, now time.Time) []daemon.WatchdogRun {
		t.Helper()
		mem, _, err := watchdogStalled(pol, js, now, mem, f.d)
		if err != nil {
			t.Fatal(err)
		}
		return mem
	}

	t.Run("idle child", func(t *testing.T) {
		f := setup(t)
		report(t, f.sessions, "conductor", session.StateWorking, t0.Add(-2*time.Hour), false)
		report(t, f.sessions, "fixer", session.StateIdle, t0, false)
		js := conduct(false)
		var mem []daemon.WatchdogRun
		mem = poll(t, f, js, mem, t0.Add(19*time.Minute))
		if len(f.sent)+len(f.stopped) != 0 {
			t.Fatalf("before the first limit: sent %v stopped %v", f.sent, f.stopped)
		}
		mem = poll(t, f, js, mem, t0.Add(21*time.Minute))
		if !slices.Equal(f.sent, []string{"fixer"}) || len(f.stopped) != 0 {
			t.Fatalf("first limit: sent %v stopped %v, want the fixer nudged", f.sent, f.stopped)
		}
		mem = poll(t, f, js, mem, t0.Add(41*time.Minute))
		if !slices.Equal(f.stopped, []string{"R:fixer,conductor"}) || !slices.Equal(f.resumed, []string{"R"}) {
			t.Fatalf("second limit: stopped %v resumed %v, want both panes closed and one resume", f.stopped, f.resumed)
		}
		if len(mem) != 1 || mem[0].Spent != 1 {
			t.Errorf("memory %+v, want the budget charged once", mem)
		}
		if !slices.Equal(f.sent, []string{"fixer"}) {
			t.Errorf("sent %v, want the fixer only", f.sent)
		}
	})
	t.Run("child said done", func(t *testing.T) {
		f := setup(t)
		report(t, f.sessions, "conductor", session.StateWorking, t0.Add(-2*time.Hour), false)
		report(t, f.sessions, "fixer", session.StateIdle, t0, true)
		js := conduct(false)
		var mem []daemon.WatchdogRun
		for _, m := range []time.Duration{21, 41, 90} {
			mem = poll(t, f, js, mem, t0.Add(m*time.Minute))
		}
		if len(f.sent)+len(f.stopped)+len(f.resumed) != 0 {
			t.Fatalf("done child: sent %v stopped %v resumed %v", f.sent, f.stopped, f.resumed)
		}
	})
	t.Run("spawn ended", func(t *testing.T) {
		f := setup(t)
		report(t, f.sessions, "conductor", session.StateWorking, t0, false)
		report(t, f.sessions, "fixer", session.StateIdle, t0.Add(-2*time.Hour), false)
		js := conduct(true)
		poll(t, f, js, nil, t0.Add(21*time.Minute))
		if slices.Contains(f.sent, "fixer") {
			t.Fatalf("ended spawn's child was nudged: sent %v", f.sent)
		}
		// With its spawn ended the conductor waits on nothing: its own
		// silence counts again.
		if !slices.Equal(f.sent, []string{"conductor"}) {
			t.Errorf("sent %v, want the conductor", f.sent)
		}
	})
	t.Run("child reports", func(t *testing.T) {
		f := setup(t)
		report(t, f.sessions, "conductor", session.StateWorking, t0.Add(-5*time.Hour), false)
		js := conduct(false)
		var mem []daemon.WatchdogRun
		for _, m := range []time.Duration{21, 41, 61, 81} {
			now := t0.Add(m * time.Minute)
			report(t, f.sessions, "fixer", session.StateWorking, now.Add(-time.Minute), false)
			mem = poll(t, f, js, mem, now)
		}
		if len(f.sent)+len(f.stopped)+len(f.resumed) != 0 {
			t.Fatalf("waiting conductor: sent %v stopped %v resumed %v", f.sent, f.stopped, f.resumed)
		}
	})
}

// TestWatchdogStalledReason checks that each stall retry and give-up names
// its reason (decision stall-reason-on-run): a stalled spawned fixer's retry
// journals one warning naming the fixer, its role, its parent and its last
// report, the charged attempt names it too, and once the budget is spent the
// hold and its inbox item do. A nudge journals nothing, nor does a step that
// moved since; a step's own session is named as the step's session.
func TestWatchdogStalledReason(t *testing.T) {
	t.Parallel()
	t0 := time.Date(2026, 11, 16, 0, 0, 0, 0, time.UTC)
	pol := WatchdogPolicy{Enabled: true, AgentNudge: 20 * time.Minute, AgentRetry: 40 * time.Minute, CommandLimit: 30 * time.Minute, WaitLimit: 90 * time.Minute, RetryBudget: 1, AutonomousOnly: true}
	type warning struct{ run, detail string }
	setup := func(t *testing.T) (*session.Store, *[]warning, watchdogDeps) {
		sessions := session.NewStore(t.TempDir())
		var warns []warning
		d := watchdogDeps{
			sessions: sessions,
			step:     func(run.State) workflow.Step { return workflow.Step{Kind: "agent"} },
			send:     func(string, string) error { return nil },
			stop:     func(run.State, []string) error { return nil },
			resume:   func(string) error { return nil },
			warn:     func(id, detail string) error { warns = append(warns, warning{id, detail}); return nil },
			now:      func() time.Time { return t0.Add(time.Hour) },
			log:      func(string, ...any) {},
		}
		return sessions, &warns, d
	}
	report := func(t *testing.T, s *session.Store, id string, at time.Time) {
		t.Helper()
		rec := session.Record{ID: id, Autonomous: true}
		rec.Status.Hook = &session.Report{State: session.StateIdle, At: at}
		if err := s.Write(rec); err != nil {
			t.Fatal(err)
		}
	}
	poll := func(t *testing.T, d watchdogDeps, js []run.Journal, mem []daemon.WatchdogRun, now time.Time) ([]daemon.WatchdogRun, []daemon.WatchdogHold) {
		t.Helper()
		mem, holds, err := watchdogStalled(pol, js, now, mem, d)
		if err != nil {
			t.Fatal(err)
		}
		return mem, holds
	}
	fixerAt := t0.Add(-5 * time.Minute)
	conduct := []run.Journal{{ID: "R", Entries: []run.Entry{
		{Type: run.TypeStarted, At: t0.Add(-3 * time.Hour), Started: &run.Started{Autonomous: true}},
		{Type: run.TypeEntered, At: t0.Add(-3 * time.Hour), Step: "conduct", Visit: 1},
		{Type: run.TypeSession, At: t0.Add(-3 * time.Hour), Step: "conduct", Session: "conductor-1"},
		{Type: run.TypeSpawn, At: t0, Step: "conduct", Parent: "conductor-1", Session: "fixer-1", Role: "fixer"},
	}}}
	fixerAtS := fixerAt.UTC().Format(time.RFC3339)

	t.Run("spawned fixer", func(t *testing.T) {
		sessions, warns, d := setup(t)
		report(t, sessions, "conductor-1", t0.Add(-2*time.Hour))
		report(t, sessions, "fixer-1", fixerAt)
		mem, _ := poll(t, d, conduct, nil, t0.Add(21*time.Minute))
		if len(*warns) != 0 {
			t.Fatalf("a nudge journalled %v, want nothing", *warns)
		}
		mem, holds := poll(t, d, conduct, mem, t0.Add(41*time.Minute))
		if len(*warns) != 1 || (*warns)[0].run != "R" {
			t.Fatalf("warnings %v, want one on run R", *warns)
		}
		w := (*warns)[0].detail
		for _, want := range []string{"fixer-1", "role fixer", "conductor-1", fixerAtS} {
			if !strings.Contains(w, want) {
				t.Errorf("warning %q does not name %q", w, want)
			}
		}
		if len(holds) != 0 || len(mem) != 1 || len(mem[0].Attempts) != 1 {
			t.Fatalf("holds %+v memory %+v, want one charged attempt", holds, mem)
		}
		if a := mem[0].Attempts[0].Action; !strings.Contains(a, "fixer-1") || !strings.Contains(a, "role fixer") {
			t.Errorf("attempt %q does not name the fixer and its role", a)
		}
		// Silent again past the second limit: the budget is spent, so the
		// run is held, and the hold and its item name the fixer.
		mem, holds = poll(t, d, conduct, mem, t0.Add(90*time.Minute))
		if len(*warns) != 1 {
			t.Errorf("warnings %v after the give-up, want the retry's only", *warns)
		}
		if len(holds) != 1 {
			t.Fatalf("holds %+v, want one", holds)
		}
		items := watchdogHoldItems(holds)
		for _, s := range []string{holds[0].Detail, items[0].Detail} {
			for _, want := range []string{"fixer-1", "role fixer", "conductor-1", fixerAtS} {
				if !strings.Contains(s, want) {
					t.Errorf("hold detail %q does not name %q", s, want)
				}
			}
		}
		// The budget hold that replaces it lists the charged attempt, which
		// names the fixer.
		bh := watchdogBudgetHolds(pol, conduct, mem, holds)
		if len(bh) != 1 || !strings.Contains(bh[0].Detail, "fixer-1") || !strings.Contains(watchdogHoldItems(bh)[0].Detail, "role fixer") {
			t.Errorf("budget holds %+v, want one naming the fixer and its role", bh)
		}
	})
	t.Run("moved since", func(t *testing.T) {
		sessions, warns, d := setup(t)
		report(t, sessions, "conductor-1", t0.Add(-2*time.Hour))
		report(t, sessions, "fixer-1", fixerAt)
		mem, _ := poll(t, d, conduct, nil, t0.Add(21*time.Minute))
		report(t, sessions, "fixer-1", t0.Add(30*time.Minute))
		mem, _ = poll(t, d, conduct, mem, t0.Add(41*time.Minute))
		if len(*warns) != 0 || len(mem[0].Attempts) != 0 {
			t.Errorf("a step that moved journalled %v, attempts %+v", *warns, mem[0].Attempts)
		}
	})
	t.Run("step's own session", func(t *testing.T) {
		sessions, warns, d := setup(t)
		report(t, sessions, "builder-1", t0)
		js := []run.Journal{{ID: "B", Entries: []run.Entry{
			{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true}},
			{Type: run.TypeEntered, At: t0, Step: "build", Visit: 1},
			{Type: run.TypeSession, At: t0, Step: "build", Session: "builder-1"},
		}}}
		mem, _ := poll(t, d, js, nil, t0.Add(21*time.Minute))
		poll(t, d, js, mem, t0.Add(41*time.Minute))
		if len(*warns) != 1 {
			t.Fatalf("warnings %v, want one", *warns)
		}
		w := (*warns)[0].detail
		if !strings.Contains(w, "session builder-1 (the step's session) last reported at "+t0.Format(time.RFC3339)) {
			t.Errorf("warning %q does not name the step's session", w)
		}
		if strings.Contains(w, "spawned by") {
			t.Errorf("warning %q names a parent", w)
		}
	})
}

func slicesHoldFor(holds []daemon.WatchdogHold, id string) bool {
	for _, h := range holds {
		if h.Run == id {
			return true
		}
	}
	return false
}

// TestWatchdogTransient drives the transient-failure recovery: each known
// signature leads to one resume --retry, an unknown failure is left alone,
// and so is a run that is live.
func TestWatchdogTransient(t *testing.T) {
	t.Parallel()
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	failed := func(id, out string, auto bool) run.Journal {
		return run.Journal{ID: id, Entries: []run.Entry{
			{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: auto}},
			{Type: run.TypeEntered, At: t0, Step: "build", Visit: 1},
			{Type: run.TypeEnded, At: t0, Step: "build", Visit: 1, Outcome: "fail", Detail: "exit 1\n" + out},
			{Type: run.TypeFinished, At: t0, Step: "build", End: "failed", Detail: "step build failed"},
		}}
	}
	js := []run.Journal{
		failed("work", "open /tmp/go-build123/b001/importcfg: no such file or directory", true),
		failed("lock", "fatal: Unable to create '/r/.git/index.lock': File exists.", true),
		failed("gh", "gh: HTTP 502 Bad Gateway", true),
		failed("real", "FAIL: TestThing expected 1 got 2", true),
		failed("live", "gh: HTTP 503 Service Unavailable", true),
		failed("owner", "gh: HTTP 502 Bad Gateway", false),
	}
	var log, resumed, classified []string
	d := watchdogDeps{
		classify: func(out string) FailureVerdict { classified = append(classified, out); return fakeClassify(out) },
		resume:   func(id string) error { resumed = append(resumed, id); return nil },
		log:      func(f string, a ...any) { log = append(log, fmt.Sprintf(f, a...)) },
	}
	pol := WatchdogPolicy{Enabled: true, RetryBudget: 1, AutonomousOnly: true}
	live := func(id string) bool { return id == "live" }
	mem, err := watchdogTransient(pol, js, live, nil, d)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(resumed, ","); got != "work,lock,gh" {
		t.Fatalf("resumed %s, want work,lock,gh", got)
	}
	if len(log) != 4 || !strings.Contains(log[3], "classified real") || !strings.Contains(log[3], "the test is failing") {
		t.Errorf("want one line per resume and one for the real failure, log %q", log)
	}
	if len(classified) != 4 {
		t.Errorf("classified %d failures, want 4", len(classified))
	}
	// The budget is one: a second poll over the same failures acts on none,
	// and the real failure, already judged, is neither classified nor logged
	// again.
	if mem, err = watchdogTransient(pol, js, live, mem, d); err != nil {
		t.Fatal(err)
	}
	if len(resumed) != 3 || len(classified) != 4 || len(log) != 4 {
		t.Errorf("second poll: resumed %v, classified %d, log %q", resumed, len(classified), log)
	}
	// A new failure of the same run (the journal moved on) is judged anew.
	real := &js[3]
	real.Entries = append(real.Entries,
		run.Entry{Type: run.TypeEntered, At: t0.Add(time.Hour), Step: "build", Visit: 2},
		run.Entry{Type: run.TypeEnded, At: t0.Add(time.Hour), Step: "build", Visit: 2, Outcome: "fail", Detail: "exit 1\ngh: HTTP 503"},
		run.Entry{Type: run.TypeFinished, At: t0.Add(time.Hour), Step: "build", End: "failed", Detail: "step build failed"})
	if _, err := watchdogTransient(pol, js, live, mem, d); err != nil {
		t.Fatal(err)
	}
	if len(classified) != 5 || strings.Join(resumed, ",") != "work,lock,gh,real" {
		t.Errorf("new failure: classified %d, resumed %v", len(classified), resumed)
	}
}

// TestWatchdogTransientBudgetAcrossPolls runs whole polls, the stalled-step
// pass included, over a run that keeps failing on a transient error: the
// budget of one holds across them, and a run whose continue is already
// queued is not charged.
func TestWatchdogTransientBudgetAcrossPolls(t *testing.T) {
	t.Parallel()
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	js := []run.Journal{{ID: "gh", Entries: []run.Entry{
		{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true}},
		{Type: run.TypeEntered, At: t0, Step: "build", Visit: 1},
		{Type: run.TypeEnded, At: t0, Step: "build", Visit: 1, Outcome: "fail", Detail: "exit 1\ngh: HTTP 502 Bad Gateway"},
		{Type: run.TypeFinished, At: t0, Step: "build", End: "failed", Detail: "step build failed"},
	}}}
	var resumed []string
	queued := true
	d := watchdogDeps{
		classify: fakeClassify,
		resume:   func(id string) error { resumed = append(resumed, id); return nil },
		queued:   func(string) bool { return queued },
		log:      func(string, ...any) {},
	}
	pol := WatchdogPolicy{Enabled: true, RetryBudget: 1, AutonomousOnly: true}
	nolive := func(string) bool { return false }
	var mem []daemon.WatchdogRun
	poll := func() {
		t.Helper()
		var err error
		if mem, _, err = watchdogPasses(pol, js, nolive, t0.Add(time.Hour), mem, d); err != nil {
			t.Fatal(err)
		}
	}
	poll()
	if len(resumed) != 0 {
		t.Fatalf("continue already queued: resumed %v", resumed)
	}
	queued = false
	for range 5 {
		poll()
	}
	if strings.Join(resumed, ",") != "gh" {
		t.Fatalf("five polls resumed %v, want gh once", resumed)
	}
	if len(mem) != 1 || mem[0].Spent != 1 {
		t.Errorf("memory %+v, want gh spent 1", mem)
	}
	// A run gone from the store is forgotten.
	js = nil
	poll()
	if len(mem) != 0 {
		t.Errorf("memory %+v after the run is gone", mem)
	}
}

// TestWatchdogCIRerun drives the CI rerun: red checks on a head are rerun
// once, the run is resumed once they pass, and a head already at attempt 2,
// or already rerun, gets no second rerun, while a new head gets its own.
func TestWatchdogCIRerun(t *testing.T) {
	t.Parallel()
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	const def = "name = \"w\"\n[step.wait]\nkind = \"command\"\nrun = \"rhiza ship check --settled\"\n"
	failed := func(id string) run.Journal {
		return run.Journal{ID: id, Entries: []run.Entry{
			{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true, Definition: def}},
			{Type: run.TypeEntered, At: t0, Step: "wait", Visit: 1},
			{Type: run.TypeEnded, At: t0, Step: "wait", Visit: 1, Outcome: "fail", Detail: "failed: test"},
			{Type: run.TypeFinished, At: t0, Step: "wait", End: "failed", Detail: "step wait failed"},
		}}
	}
	red := watchdogPR{Head: "aaaa111", Checks: []ShipCheck{{Name: "test", Status: "COMPLETED", Conclusion: "FAILURE", RunID: "7"}}, Attempts: map[string]int{"7": 1}}
	green := watchdogPR{Head: "aaaa111", Checks: []ShipCheck{{Name: "test", Status: "COMPLETED", Conclusion: "SUCCESS", RunID: "7"}}, Attempts: map[string]int{"7": 2}}
	pr := red
	var reran, resumed, log []string
	d := watchdogDeps{
		ci:       func(run.State) (watchdogPR, error) { return pr, nil },
		rerun:    func(_ run.State, id string) error { reran = append(reran, id); return nil },
		classify: fakeClassify,
		resume:   func(id string) error { resumed = append(resumed, id); return nil },
		log:      func(f string, a ...any) { log = append(log, fmt.Sprintf(f, a...)) },
	}
	pol := WatchdogPolicy{Enabled: true, RetryBudget: 2, CIReruns: 1, AutonomousOnly: true}
	nolive := func(string) bool { return false }
	js := []run.Journal{failed("r")}
	poll := func(mem []daemon.WatchdogRun) []daemon.WatchdogRun {
		t.Helper()
		mem, err := watchdogCIRerun(pol, js, nolive, mem, d)
		if err != nil {
			t.Fatal(err)
		}
		return mem
	}

	mem := poll(nil)
	if strings.Join(reran, ",") != "7" || len(resumed) != 0 {
		t.Fatalf("red head: reran %v resumed %v", reran, resumed)
	}
	rerun := slices.Clone(mem)
	// Still running after the rerun: nothing yet.
	pr = red
	pr.Checks = []ShipCheck{{Name: "test", Status: "IN_PROGRESS", RunID: "7"}}
	pr.Attempts = map[string]int{"7": 2}
	mem = poll(mem)
	// The rerun passed: the run is resumed once.
	pr = green
	mem = poll(mem)
	mem = poll(mem)
	if strings.Join(resumed, ",") != "r" {
		t.Fatalf("green: resumed %v, want one", resumed)
	}
	if len(log) != 2 {
		t.Errorf("want one line per action, log %q", log)
	}
	// Red again after the rerun (attempt 2): no second rerun, no resume, and
	// the head is left to the owner until the journal moves on.
	resumed = nil
	pr = red
	pr.Attempts = map[string]int{"7": 2}
	red2 := poll(slices.Clone(rerun))
	if len(reran) != 1 || len(resumed) != 0 || red2[0].CIDone.IsZero() {
		t.Fatalf("red again: reran %v resumed %v mem %+v", reran, resumed, red2)
	}
	pr = green
	poll(red2)
	if len(resumed) != 0 {
		t.Fatalf("left to the owner: resumed %v", resumed)
	}
	// A head whose run is already at attempt 2 (the plan step's rerun) is not
	// rerun again.
	reran = nil
	pr = red
	pr.Attempts = map[string]int{"7": 2}
	poll(nil)
	if len(reran) != 0 {
		t.Errorf("attempt above 1 reran %v", reran)
	}
	// A new head is allowed its own rerun.
	pr = red
	pr.Head = "bbbb222"
	poll(mem)
	if strings.Join(reran, ",") != "7" {
		t.Errorf("new head reran %v, want 7", reran)
	}
	// A model-only PR's head: the Go jobs SKIPPED, the model jobs green. Nothing
	// is red, so nothing is rerun and no rerun is spent.
	reran = nil
	pr = watchdogPR{Head: "cccc333", Attempts: map[string]int{"7": 1}}
	for _, n := range []string{"lint", "scripts", "api", "rest", "budgets"} {
		pr.Checks = append(pr.Checks, ShipCheck{Name: n, Status: "COMPLETED", Conclusion: "SKIPPED", RunID: "7"})
	}
	for _, n := range []string{"changes", "model", "check"} {
		pr.Checks = append(pr.Checks, ShipCheck{Name: n, Status: "COMPLETED", Conclusion: "SUCCESS", RunID: "7"})
	}
	if rerunSpent(pr, pol) {
		t.Error("rerunSpent counts skipped jobs")
	}
	poll(nil)
	if len(reran) != 0 {
		t.Errorf("skipped jobs reran %v", reran)
	}
}

// With ciReruns N the watchdog reruns a head's failed jobs while an attempt
// is at most N, and never once an attempt is above N.
func TestWatchdogCIRerunBudgetN(t *testing.T) {
	t.Parallel()
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	const def = "name = \"w\"\n[step.wait]\nkind = \"command\"\nrun = \"rhiza ship check --settled\"\n"
	j := run.Journal{ID: "r", Entries: []run.Entry{
		{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true, Definition: def}},
		{Type: run.TypeEntered, At: t0, Step: "wait", Visit: 1},
		{Type: run.TypeEnded, At: t0, Step: "wait", Visit: 1, Outcome: "fail", Detail: "failed: test"},
		{Type: run.TypeFinished, At: t0, Step: "wait", End: "failed", Detail: "step wait failed"},
	}}
	for _, n := range []int{1, 2, 3} {
		for attempt := 1; attempt <= n+1; attempt++ {
			pr := watchdogPR{Head: "aaaa111", Checks: []ShipCheck{{Name: "test", Status: "COMPLETED", Conclusion: "FAILURE", RunID: "7"}}, Attempts: map[string]int{"7": attempt}}
			var reran []string
			d := watchdogDeps{
				ci:       func(run.State) (watchdogPR, error) { return pr, nil },
				rerun:    func(_ run.State, id string) error { reran = append(reran, id); return nil },
				classify: fakeClassify,
				resume:   func(string) error { return nil },
				log:      func(string, ...any) {},
			}
			pol := WatchdogPolicy{Enabled: true, RetryBudget: 2, CIReruns: n, AutonomousOnly: true}
			if _, err := watchdogCIRerun(pol, []run.Journal{j}, func(string) bool { return false }, nil, d); err != nil {
				t.Fatal(err)
			}
			want := 0
			if attempt <= n {
				want = 1
			}
			if len(reran) != want {
				t.Errorf("ciReruns %d, attempt %d: reran %v, want %d reruns", n, attempt, reran, want)
			}
		}
	}
}

// TestWatchdogCIReadsBounded checks the PR reads of the CI passes through
// watchdogPasses: a run's PR is read at most once per poll by the three
// passes together, and not at all while its journal stays where the PR was
// found gone, its head red again after the rerun, or its budget spent.
func TestWatchdogCIReadsBounded(t *testing.T) {
	t.Parallel()
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	const def = "name = \"w\"\n[step.wait]\nkind = \"command\"\nrun = \"rhiza ship check --settled\"\non.pass = \"merge\"\n[step.merge]\nkind = \"command\"\nrun = \"rhiza ship merge\"\n"
	failed := func(id string) run.Journal {
		return run.Journal{ID: id, Entries: []run.Entry{
			{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true, Definition: def}},
			{Type: run.TypeEntered, At: t0, Step: "wait", Visit: 1},
			{Type: run.TypeEnded, At: t0, Step: "wait", Visit: 1, Outcome: "fail", Detail: "failed: test"},
			{Type: run.TypeFinished, At: t0, Step: "wait", End: "failed", Detail: "step wait failed"},
		}}
	}
	red := watchdogPR{Head: "aaaa111", Mergeable: "MERGEABLE", Checks: []ShipCheck{{Name: "test", Status: "COMPLETED", Conclusion: "FAILURE", RunID: "7"}}, Attempts: map[string]int{"7": 2}}
	prs := map[string]watchdogPR{"gone": {}, "red": red, "spent": red}
	reads := map[string]int{}
	d := watchdogDeps{
		ci:       func(st run.State) (watchdogPR, error) { reads[st.ID]++; return prs[st.ID], nil },
		rerun:    func(run.State, string) error { t.Fatal("reran"); return nil },
		classify: fakeClassify,
		resume:   func(string) error { t.Fatal("resumed"); return nil },
		merge:    func(run.State) (watchdogMerge, error) { t.Fatal("merged"); return watchdogMerge{}, nil },
		log:      func(string, ...any) {},
	}
	pol := WatchdogPolicy{Enabled: true, RetryBudget: 1, CIReruns: 1, AutonomousOnly: true}
	none := func(string) bool { return false }
	js := []run.Journal{failed("gone"), failed("red"), failed("spent")}
	mem := []daemon.WatchdogRun{{Run: "spent", Spent: 1, Attempts: []daemon.WatchdogAttempt{{At: t0, Action: "x"}}}}
	for range 3 {
		var err error
		mem, _, err = watchdogPasses(pol, js, none, t0, mem, d)
		if err != nil {
			t.Fatal(err)
		}
	}
	if reads["gone"] != 1 || reads["red"] != 1 || reads["spent"] != 0 {
		t.Fatalf("reads over three polls %v, want gone 1, red 1, spent 0", reads)
	}
	// A journal that moved on is read again, once for the poll.
	moved := failed("red")
	moved.Entries[len(moved.Entries)-1].At = t0.Add(time.Minute)
	js[1] = moved
	watchdogPasses(pol, js, none, t0, mem, d)
	if reads["red"] != 2 {
		t.Fatalf("red read %d times after the journal moved, want 2", reads["red"])
	}
}

// TestWatchdogNestedStep: a stalled step three levels deep is named whole in
// the watchdog's log and in the inbox hold it leaves.
func TestWatchdogNestedStep(t *testing.T) {
	t.Parallel()
	const nested = "ship-change/build/build"
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	sessions := session.NewStore(t.TempDir())
	rec := session.Record{ID: "s1", Autonomous: false}
	rec.Status.Hook = &session.Report{State: session.StateWorking, At: t0}
	if err := sessions.Write(rec); err != nil {
		t.Fatal(err)
	}
	j := run.Journal{ID: "N", Entries: []run.Entry{
		{Type: run.TypeStarted, At: t0, Started: &run.Started{}},
		{Type: run.TypeEntered, At: t0, Step: nested, Visit: 1},
		{Type: run.TypeSession, At: t0, Step: nested, Session: "s1"},
	}}
	var log []string
	d := watchdogDeps{
		sessions: sessions,
		step:     func(run.State) workflow.Step { return workflow.Step{Kind: "agent"} },
		send:     func(string, string) error { return nil },
		stop:     func(run.State, []string) error { return nil },
		classify: fakeClassify,
		resume:   func(string) error { return nil },
		log:      func(f string, a ...any) { log = append(log, fmt.Sprintf(f, a...)) },
	}
	pol := WatchdogPolicy{Enabled: true, AgentNudge: 20 * time.Minute, AgentRetry: 40 * time.Minute, CommandLimit: 30 * time.Minute, WaitLimit: 90 * time.Minute, RetryBudget: 1, AutonomousOnly: true}
	_, holds, err := watchdogStalled(pol, []run.Journal{j}, t0.Add(50*time.Minute), nil, d)
	if err != nil {
		t.Fatal(err)
	}
	if len(holds) != 1 || holds[0].Step != nested {
		t.Fatalf("holds %+v, want one at %s", holds, nested)
	}
	if !strings.Contains(holds[0].Detail, "step "+nested+" ") {
		t.Errorf("hold detail %q does not name %s", holds[0].Detail, nested)
	}
	items := watchdogHoldItems(holds)
	if len(items) != 1 || items[0].Step != nested {
		t.Errorf("inbox items %+v, want step %s", items, nested)
	}
}

// TestWatchdogTransientClassifierError leaves a failure to the inbox when the
// classifier cannot run or is not set, even if the output looks transient.
func TestWatchdogTransientClassifierError(t *testing.T) {
	t.Parallel()
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	js := []run.Journal{{ID: "gh", Entries: []run.Entry{
		{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true}},
		{Type: run.TypeEntered, At: t0, Step: "build", Visit: 1},
		{Type: run.TypeEnded, At: t0, Step: "build", Visit: 1, Outcome: "fail", Detail: "exit 1\ngh: HTTP 502 Bad Gateway"},
		{Type: run.TypeFinished, At: t0, Step: "build", End: "failed", Detail: "step build failed"},
	}}}
	pol := WatchdogPolicy{Enabled: true, RetryBudget: 1, AutonomousOnly: true}
	nolive := func(string) bool { return false }
	for name, classify := range map[string]func(string) FailureVerdict{
		"error": func(string) FailureVerdict {
			return ClassifyFailure("c", "x", func(string, string) (string, error) { return "", errors.New("exit 1") })
		},
		"none": nil,
	} {
		var log []string
		d := watchdogDeps{
			classify: classify,
			resume:   func(string) error { t.Fatalf("%s: resumed", name); return nil },
			log:      func(f string, a ...any) { log = append(log, fmt.Sprintf(f, a...)) },
		}
		mem, err := watchdogTransient(pol, js, nolive, nil, d)
		if err != nil {
			t.Fatal(err)
		}
		// A second poll leaves the judged failure alone, logging nothing.
		if _, err := watchdogTransient(pol, js, nolive, mem, d); err != nil {
			t.Fatal(err)
		}
		want := map[string]string{"error": "classifier failed", "none": "no failure classifier"}[name]
		if len(log) != 1 || !strings.Contains(log[0], "left to the inbox") || !strings.Contains(log[0], want) {
			t.Errorf("%s: log %q", name, log)
		}
	}
}

// A step whose fresh session waits on a pre-start hold has no session yet;
// however long the hold lasts the watchdog neither nudges, stops nor
// charges it. Once the session entry releases the hold, the step is
// stall-able again.
func TestWatchdogStepHoldNotStalled(t *testing.T) {
	t.Parallel()
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	pol := WatchdogPolicy{Enabled: true, AgentNudge: 20 * time.Minute, AgentRetry: 40 * time.Minute, CommandLimit: 30 * time.Minute, WaitLimit: 90 * time.Minute, RetryBudget: 1, AutonomousOnly: true}
	var sent, stopped []string
	d := watchdogDeps{
		sessions: session.NewStore(t.TempDir()),
		step:     func(run.State) workflow.Step { return workflow.Step{Kind: "agent"} },
		send:     func(sid, _ string) error { sent = append(sent, sid); return nil },
		stop:     func(st run.State, _ []string) error { stopped = append(stopped, st.ID); return nil },
		classify: fakeClassify,
		resume:   func(string) error { return nil },
		log:      func(string, ...any) {},
	}
	j := run.Journal{ID: "H", Entries: []run.Entry{
		{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true}},
		{Type: run.TypeEntered, At: t0, Step: "build", Visit: 1},
		{Type: run.TypeSpawnHeld, At: t0, Step: "build", Detail: "machine hold: memory-low: 512.0 MiB available, under the minimum of 1.0 GiB"},
	}}
	mem, holds, err := watchdogStalled(pol, []run.Journal{j}, t0.Add(5*time.Hour), nil, d)
	if err != nil {
		t.Fatal(err)
	}
	if len(sent)+len(stopped)+len(holds) != 0 || len(mem) != 1 || mem[0].Spent != 0 {
		t.Fatalf("held step: sent %v stopped %v holds %+v mem %+v", sent, stopped, holds, mem)
	}
	// The session entry releases the hold: a silent step with no live
	// session counts again.
	j.Entries = append(j.Entries, run.Entry{Type: run.TypeSession, At: t0, Step: "build", Session: "gone"})
	if _, _, err := watchdogStalled(pol, []run.Journal{j}, t0.Add(5*time.Hour), nil, d); err != nil {
		t.Fatal(err)
	}
	if strings.Join(stopped, ",") != "H" {
		t.Fatalf("released step: stopped %v", stopped)
	}
}
