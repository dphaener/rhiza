package api

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/proctree"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// runawayFixture fakes everything the runaway pass touches: the process
// table is a slice the test sets, never the host's ps.
type runawayFixture struct {
	sessions                     *session.Store
	table                        []proctree.Process
	reads                        int
	sent, stopped, resumed, bugs []string
	nudges                       []string
	journalled                   []string
	log                          []string
	d                            watchdogDeps
}

func newRunawayFixture(t *testing.T) *runawayFixture {
	t.Helper()
	f := &runawayFixture{sessions: session.NewStore(t.TempDir())}
	f.d = watchdogDeps{
		sessions: f.sessions,
		step:     func(run.State) workflow.Step { return workflow.Step{Kind: "agent"} },
		send: func(sid, text string) error {
			f.sent = append(f.sent, sid)
			f.nudges = append(f.nudges, text)
			return nil
		},
		stop: func(st run.State, sids []string) error {
			f.stopped = append(f.stopped, st.ID+":"+strings.Join(sids, ","))
			return nil
		},
		resume: func(id string) error { f.resumed = append(f.resumed, id); return nil },
		procs: func() ([]proctree.Process, error) {
			f.reads++
			return f.table, nil
		},
		stopped: func(id, step, detail string) error {
			f.journalled = append(f.journalled, id+"|"+step+"|"+detail)
			return nil
		},
		requests: func() ([]daemon.Request, error) { return nil, nil },
		file:     func(text string) error { f.bugs = append(f.bugs, text); return nil },
		log:      func(fm string, a ...any) { f.log = append(f.log, fmt.Sprintf(fm, a...)) },
	}
	return f
}

// report writes a session record working since at.
func (f *runawayFixture) report(t *testing.T, sid string, auto bool, at time.Time) {
	t.Helper()
	rec := session.Record{ID: sid, Autonomous: auto}
	rec.Status.Hook = &session.Report{State: session.StateWorking, At: at}
	if err := f.sessions.Write(rec); err != nil {
		t.Fatal(err)
	}
}

// tree adds sid's tree to the table: its `rhiza session run` root, under
// an unrelated shell, and n-1 descendants, each of rss bytes, starting at
// pid base.
func (f *runawayFixture) tree(sid string, base, n int, rss int64) {
	f.table = append(f.table, proctree.Process{PID: base, PPID: 1, RSS: rss, Args: "/usr/bin/rhiza session run " + sid})
	for k := 1; k < n; k++ {
		parent := base
		if k > 1 {
			parent = base + 1
		}
		f.table = append(f.table, proctree.Process{PID: base + k, PPID: parent, RSS: rss, Args: "node worker"})
	}
}

var runawayPolicy = WatchdogPolicy{
	Enabled: true, AgentNudge: 20 * time.Minute, AgentRetry: 40 * time.Minute, CommandLimit: 30 * time.Minute,
	WaitLimit: 90 * time.Minute, RetryBudget: 1, FileBugs: true, AutonomousOnly: true,
	RunawayProcs: 400, RunawayMemory: 8 << 30, RunawayGrace: 5 * time.Minute,
}

func runawayJournal(id string, auto bool, t0 time.Time, sid string) run.Journal {
	return run.Journal{ID: id, Entries: []run.Entry{
		{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: auto}},
		{Type: run.TypeEntered, At: t0, Step: "build", Visit: 1},
		{Type: run.TypeSession, At: t0, Step: "build", Session: sid},
	}}
}

// runawayPoll runs one poll of every watchdog pass over js, all live.
func runawayPoll(t *testing.T, pol WatchdogPolicy, f *runawayFixture, js []run.Journal, mem []daemon.WatchdogRun, now time.Time) []daemon.WatchdogRun {
	t.Helper()
	mem, _, err := watchdogPasses(pol, js, func(string) bool { return true }, now, mem, f.d)
	if err != nil {
		t.Fatal(err)
	}
	return mem
}

func runawayNudges(mem []daemon.WatchdogRun, id string) []daemon.RunawayNudge {
	for _, m := range mem {
		if m.Run == id {
			return m.Runaway
		}
	}
	return nil
}

// TestWatchdogRunaway drives an autonomous run whose step session's tree
// holds 450 processes: nudged once with the reading, then a grace later
// stopped, its step journalled stopped with the reading, never resumed,
// and a runaway bug filed.
func TestWatchdogRunaway(t *testing.T) {
	t.Parallel()
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	f := newRunawayFixture(t)
	f.report(t, "auto", true, t0)
	f.tree("auto", 1000, 450, 2<<20)
	f.tree("other", 5000, 900, 2<<20) // a session of no run: never read as one
	js := []run.Journal{runawayJournal("R", true, t0, "auto")}

	mem := runawayPoll(t, runawayPolicy, f, js, nil, t0.Add(time.Minute))
	if !slices.Equal(f.sent, []string{"auto"}) || len(f.stopped) != 0 {
		t.Fatalf("first poll: sent %v stopped %v, want auto nudged", f.sent, f.stopped)
	}
	if n := f.nudges[0]; !strings.Contains(n, "450 processes") || !strings.Contains(n, "runawayProcs 400") {
		t.Errorf("nudge %q does not name the count and the threshold", n)
	}
	if ns := runawayNudges(mem, "R"); len(ns) != 1 || ns[0].Session != "auto" || !ns[0].At.Equal(t0.Add(time.Minute)) {
		t.Fatalf("nudge memory %+v", ns)
	}

	mem = runawayPoll(t, runawayPolicy, f, js, mem, t0.Add(4*time.Minute))
	if len(f.sent) != 1 || len(f.stopped) != 0 {
		t.Fatalf("inside the grace: sent %v stopped %v, want one nudge and no stop", f.sent, f.stopped)
	}

	mem = runawayPoll(t, runawayPolicy, f, js, mem, t0.Add(6*time.Minute))
	if !slices.Equal(f.stopped, []string{"R:auto"}) {
		t.Fatalf("after the grace: stopped %v, want R with auto's pane", f.stopped)
	}
	if len(f.journalled) != 1 {
		t.Fatalf("journalled %v, want one stopped entry", f.journalled)
	}
	j := f.journalled[0]
	for _, want := range []string{"R|build|", "session auto", "450 processes", "900.0 MiB resident memory", "runawayProcs 400"} {
		if !strings.Contains(j, want) {
			t.Errorf("stopped entry %q does not name %q", j, want)
		}
	}
	if len(f.resumed) != 0 {
		t.Errorf("resumed %v: a runaway is never resumed or retried", f.resumed)
	}
	if len(f.bugs) != 1 || !strings.HasPrefix(f.bugs[0], "watchdog runaway session at step build: run R") {
		t.Errorf("bugs %q, want one runaway bug", f.bugs)
	}
	if ns := runawayNudges(mem, "R"); len(ns) != 0 {
		t.Errorf("nudge memory after the stop %+v", ns)
	}
	for _, m := range mem {
		if m.Run == "R" && (m.Spent != 0 || len(m.Attempts) != 0) {
			t.Errorf("the stop was charged as a retry: %+v", m)
		}
	}
}

// TestWatchdogRunawaySpawnedChild catches a spawned fixer whose tree is
// over runawayMemory: the fixer is nudged, not its waiting conductor, and
// the stop closes both panes and names the fixer with its role.
func TestWatchdogRunawaySpawnedChild(t *testing.T) {
	t.Parallel()
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	f := newRunawayFixture(t)
	f.report(t, "conductor", true, t0)
	f.report(t, "fixer", true, t0)
	f.tree("conductor", 100, 3, 100<<20)
	f.tree("fixer", 2000, 10, 1<<30)
	js := []run.Journal{{ID: "C", Entries: []run.Entry{
		{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true}},
		{Type: run.TypeEntered, At: t0, Step: "conduct", Visit: 1},
		{Type: run.TypeSession, At: t0, Step: "conduct", Session: "conductor"},
		{Type: run.TypeSpawn, At: t0, Step: "conduct", Parent: "conductor", Session: "fixer", Role: "fixer"},
	}}}
	mem := runawayPoll(t, runawayPolicy, f, js, nil, t0.Add(time.Minute))
	if !slices.Equal(f.sent, []string{"fixer"}) || !strings.Contains(f.nudges[0], "runawayMemory 8.0 GiB") {
		t.Fatalf("sent %v nudges %q, want the fixer nudged over runawayMemory", f.sent, f.nudges)
	}
	runawayPoll(t, runawayPolicy, f, js, mem, t0.Add(7*time.Minute))
	if !slices.Equal(f.stopped, []string{"C:fixer,conductor"}) {
		t.Fatalf("stopped %v, want the fixer's and the conductor's panes", f.stopped)
	}
	if len(f.journalled) != 1 || !strings.Contains(f.journalled[0], "session fixer (role fixer, spawned by session conductor): 10 processes, 10.0 GiB resident memory, over runawayMemory 8.0 GiB") {
		t.Fatalf("journalled %q", f.journalled)
	}
	if len(f.resumed) != 0 {
		t.Errorf("resumed %v", f.resumed)
	}
}

// TestWatchdogRunawayLeftAlone holds the failure cases: under both
// thresholds, back under after the nudge, a threshold of 0, an attended
// run under autonomousOnly and a session with no tree are never stopped.
func TestWatchdogRunawayLeftAlone(t *testing.T) {
	t.Parallel()
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	quiet := func(t *testing.T, f *runawayFixture) {
		t.Helper()
		if len(f.sent)+len(f.stopped)+len(f.journalled)+len(f.bugs)+len(f.resumed) != 0 {
			t.Fatalf("sent %v stopped %v journalled %v bugs %v resumed %v, want nothing", f.sent, f.stopped, f.journalled, f.bugs, f.resumed)
		}
	}

	t.Run("under both", func(t *testing.T) {
		t.Parallel()
		f := newRunawayFixture(t)
		f.report(t, "s", true, t0)
		f.tree("s", 10, 400, 10<<20)
		js := []run.Journal{runawayJournal("R", true, t0, "s")}
		mem := runawayPoll(t, runawayPolicy, f, js, nil, t0.Add(time.Minute))
		runawayPoll(t, runawayPolicy, f, js, mem, t0.Add(10*time.Minute))
		quiet(t, f)
	})

	t.Run("back under", func(t *testing.T) {
		t.Parallel()
		f := newRunawayFixture(t)
		f.report(t, "s", true, t0)
		f.tree("s", 10, 450, 1<<20)
		js := []run.Journal{runawayJournal("R", true, t0, "s")}
		mem := runawayPoll(t, runawayPolicy, f, js, nil, t0.Add(time.Minute))
		if len(f.sent) != 1 {
			t.Fatalf("sent %v, want a nudge", f.sent)
		}
		f.table = nil
		f.tree("s", 10, 50, 1<<20)
		mem = runawayPoll(t, runawayPolicy, f, js, mem, t0.Add(7*time.Minute))
		if len(f.stopped) != 0 {
			t.Fatalf("stopped %v a tree back under", f.stopped)
		}
		if ns := runawayNudges(mem, "R"); len(ns) != 0 {
			t.Fatalf("nudge not cleared: %+v", ns)
		}
		// Over again: a fresh nudge, not a stop on the old one.
		f.table = nil
		f.tree("s", 10, 450, 1<<20)
		runawayPoll(t, runawayPolicy, f, js, mem, t0.Add(8*time.Minute))
		if len(f.sent) != 2 || len(f.stopped) != 0 {
			t.Fatalf("over again: sent %v stopped %v, want a second nudge only", f.sent, f.stopped)
		}
	})

	t.Run("threshold off", func(t *testing.T) {
		t.Parallel()
		f := newRunawayFixture(t)
		f.report(t, "s", true, t0)
		f.tree("s", 10, 5000, 1<<30)
		js := []run.Journal{runawayJournal("R", true, t0, "s")}
		pol := runawayPolicy
		pol.RunawayProcs, pol.RunawayMemory = 0, 0
		mem := runawayPoll(t, pol, f, js, nil, t0.Add(time.Minute))
		runawayPoll(t, pol, f, js, mem, t0.Add(10*time.Minute))
		quiet(t, f)
		if f.reads != 0 {
			t.Errorf("read the table %d times with both thresholds off", f.reads)
		}
		// One threshold off: only the other is passed.
		pol = runawayPolicy
		pol.RunawayProcs = 0
		f.table = nil
		f.tree("s", 10, 5000, 1<<20)
		runawayPoll(t, pol, f, js, nil, t0.Add(time.Minute))
		quiet(t, f)
	})

	t.Run("attended", func(t *testing.T) {
		t.Parallel()
		f := newRunawayFixture(t)
		f.report(t, "owner", false, t0)
		f.tree("owner", 10, 450, 1<<20)
		js := []run.Journal{runawayJournal("O", false, t0, "owner")}
		mem := runawayPoll(t, runawayPolicy, f, js, nil, t0.Add(time.Minute))
		runawayPoll(t, runawayPolicy, f, js, mem, t0.Add(10*time.Minute))
		quiet(t, f)
		if !slices.ContainsFunc(f.log, func(l string) bool {
			return strings.Contains(l, "run O") && strings.Contains(l, "450 processes") && strings.Contains(l, "left to the owner")
		}) {
			t.Errorf("log %q, want the attended runaway logged", f.log)
		}
	})

	t.Run("no tree", func(t *testing.T) {
		t.Parallel()
		f := newRunawayFixture(t)
		f.report(t, "s", true, t0)
		f.tree("someone-else", 10, 450, 1<<20)
		f.table = append(f.table, proctree.Process{PID: 9000, PPID: 1, Args: "rex tab -- rhiza session run s"})
		js := []run.Journal{runawayJournal("R", true, t0, "s")}
		mem := runawayPoll(t, runawayPolicy, f, js, nil, t0.Add(time.Minute))
		runawayPoll(t, runawayPolicy, f, js, mem, t0.Add(10*time.Minute))
		quiet(t, f)
	})
}

// TestWatchdogRunawayNoHostTable: under RHIZA_TEST_NOW the daemon has no
// process-table reader, so no test daemon runs ps.
func TestWatchdogRunawayNoHostTable(t *testing.T) {
	t.Parallel()
	env := map[string]string{"RHIZA_TEST_NOW": "2030-01-01T00:00:00Z"}
	ctx := &Ctx{Getenv: func(k string) string { return env[k] }}
	if watchdogProcTable(ctx) != nil {
		t.Fatal("the daemon reads the host's process table under RHIZA_TEST_NOW")
	}
	delete(env, "RHIZA_TEST_NOW")
	if watchdogProcTable(ctx) == nil {
		t.Fatal("a real daemon has no process-table reader")
	}
}

// TestWatchdogRunawayBYO catches a runaway session of the bring-your-own
// example workflow's agent step work (myteam-flow): its tree, read from a
// fake process table, holds 450 processes. It is nudged with the reading,
// then a grace later the run is stopped and step work journalled stopped
// with the runaway reason, never resumed.
func TestWatchdogRunawayBYO(t *testing.T) {
	t.Parallel()
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	def := byoWorkflow(t)
	f := newRunawayFixture(t)
	f.d.step = func(st run.State) workflow.Step { return def.Steps[st.Current] }
	f.report(t, "worker", true, t0)
	f.tree("worker", 1000, 450, 2<<20)
	js := []run.Journal{{ID: "W", Entries: []run.Entry{
		{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true, Workflow: "myteam-flow"}},
		{Type: run.TypeEntered, At: t0, Step: "services", Visit: 1},
		{Type: run.TypeEnded, At: t0, Step: "services", Visit: 1, Outcome: "pass"},
		{Type: run.TypeEntered, At: t0, Step: "work", Visit: 1},
		{Type: run.TypeSession, At: t0, Step: "work", Session: "worker"},
	}}}
	mem := runawayPoll(t, runawayPolicy, f, js, nil, t0.Add(time.Minute))
	if !slices.Equal(f.sent, []string{"worker"}) || !strings.Contains(f.nudges[0], "450 processes") {
		t.Fatalf("sent %v nudges %q, want worker nudged with its reading", f.sent, f.nudges)
	}
	runawayPoll(t, runawayPolicy, f, js, mem, t0.Add(7*time.Minute))
	if !slices.Equal(f.stopped, []string{"W:worker"}) {
		t.Fatalf("stopped %v, want W with worker's pane", f.stopped)
	}
	if len(f.journalled) != 1 {
		t.Fatalf("journalled %q, want one stopped entry", f.journalled)
	}
	for _, want := range []string{"W|work|", "session worker", "450 processes", "over runawayProcs 400"} {
		if !strings.Contains(f.journalled[0], want) {
			t.Errorf("stopped entry %q does not name %q", f.journalled[0], want)
		}
	}
	if len(f.resumed) != 0 {
		t.Errorf("resumed %v: a runaway is never resumed", f.resumed)
	}
	if len(f.bugs) != 1 || !strings.HasPrefix(f.bugs[0], "watchdog runaway session at step work: run W") {
		t.Errorf("bugs %q, want one runaway bug at step work", f.bugs)
	}
}
