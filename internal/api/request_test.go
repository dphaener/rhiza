package api

import (
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

// A plan request whose run never journals (run start failed in its tab
// before journaling) goes back in the queue once daemonStartGrace passed,
// and is started again, with a new run, daemonUnjournaledHold after its
// first start.
func TestDaemonPlanRequestsRequeuesUnjournaled(t *testing.T) {
	f := newSessionFixture(t, "req")
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	f.ctx.Now = func() time.Time { return now }
	f.ctx.repoSession = true // the rhiza-<repo> session is daemon-repo-session's concern
	f.ctx.Self = "/bin/rhiza"
	var detached [][]string
	f.ctx.Detach = func(p Program, logPath string) error {
		detached = append(detached, p.Args)
		return nil
	}
	if _, err := callOp(t, f.ctx, "request", `{"text":"Dark mode"}`); err != nil {
		t.Fatal(err)
	}
	rs := run.NewStore(run.Dir(f.common))
	poll := func() []string {
		t.Helper()
		detached = nil
		got, err := daemonPlanRequests(f.ctx, rs, nil, 0, 2, "auto", "", nil)
		if err != nil {
			t.Fatal(err)
		}
		return got
	}
	request := func() daemon.Request {
		t.Helper()
		reqs, err := daemon.ReadRequests(f.common)
		if err != nil || len(reqs) != 1 {
			t.Fatalf("requests %v, %v", reqs, err)
		}
		return reqs[0]
	}
	if got := poll(); len(got) != 1 || len(detached) != 1 {
		t.Fatalf("first poll started %v with rex calls %q, want one start", got, detached)
	}
	first := request()
	if first.Run == "" {
		t.Fatal("started request has no run")
	}
	now = now.Add(daemonStartGrace - time.Second)
	if got := poll(); len(got) != 0 || request().Run != first.Run {
		t.Fatalf("within the grace: started %v, run %q, want none and %q", got, request().Run, first.Run)
	}
	now = now.Add(time.Second)
	if got := poll(); len(got) != 0 || len(detached) != 0 {
		t.Fatalf("after the grace: started %v, want none until the hold passed", got)
	}
	q := request()
	if q.Run != "" || !q.RetryAt.Equal(first.StartedAt.Add(daemonUnjournaledHold)) {
		t.Fatalf("after the grace: run %q, retry at %v, want back in the queue until %v", q.Run, q.RetryAt, first.StartedAt.Add(daemonUnjournaledHold))
	}
	if state, _ := daemon.RequestState(q, nil); state != daemon.RequestQueued {
		t.Errorf("after the grace: state %s, want queued", state)
	}
	now = first.StartedAt.Add(daemonUnjournaledHold)
	if got := poll(); len(got) != 1 || len(detached) != 1 {
		t.Fatalf("after the hold: started %v, want the request again", got)
	}
	if q := request(); q.Run == "" || q.Run == first.Run || !q.RetryAt.IsZero() {
		t.Errorf("after the hold: run %q (first %q), retry at %v, want a new run", q.Run, first.Run, q.RetryAt)
	}
}

func TestRequestPlanRunArgsPassRef(t *testing.T) {
	for _, ref := range []string{"https://example.com/issues/7", "--dry-run", "--change=01M6VY7JMM5NPB3QVG1ECZH1YW", "-h"} {
		with := planRunArgs(daemon.Request{ID: "r1", Ref: ref}, "plan", "default", false, "manual")
		if !slices.Contains(with, "--ref="+ref) || slices.Contains(with, ref) {
			t.Fatalf("with ref %q: want the one word --ref=%s in %q", ref, ref, with)
		}
	}
	without := planRunArgs(daemon.Request{ID: "r2"}, "plan", "default", false, "manual")
	if slices.ContainsFunc(without, func(a string) bool { return strings.HasPrefix(a, "--ref") }) {
		t.Fatalf("without a ref: %q", without)
	}
}

// A queued milestone request is demand while no other milestone request's
// run is live or unjournaled: a started milestone request stays in the
// queue after its run ends, and blocks the next one only while its run
// runs, or while its start is within daemonStartGrace and not journaled.
func TestDaemonMilestoneWaitsAfterEndedRun(t *testing.T) {
	f := newSessionFixture(t, "req")
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	f.ctx.Now = func() time.Time { return now }
	host, _ := os.Hostname()
	rs := run.NewStore(run.Dir(f.common))
	const first = "01M6YRUNFIRST0000000000000"
	write := func(firstRun string) {
		t.Helper()
		if err := daemon.WriteRequests(f.common, []daemon.Request{
			{ID: "r1", Kind: daemon.KindMilestone, Text: "First", At: now.Add(-time.Hour), Run: firstRun, StartedAt: now.Add(-time.Second)},
			{ID: "r2", Kind: daemon.KindMilestone, Text: "Second", At: now.Add(-time.Minute)},
		}); err != nil {
			t.Fatal(err)
		}
	}
	waits := func() bool {
		t.Helper()
		runs, err := rs.List()
		if err != nil {
			t.Fatal(err)
		}
		return daemonMilestoneWaits(f.ctx, rs, runs, host)
	}

	// The first's start has not journaled, within the grace: it blocks.
	write(first)
	if waits() {
		t.Fatal("first start unjournaled within the grace: second waits, want blocked")
	}
	// Its run journals and runs (a process on another host is running).
	if err := rs.Create(first, run.Entry{Type: run.TypeStarted, At: now, Started: &run.Started{Workflow: "plan-milestone"}}); err != nil {
		t.Fatal(err)
	}
	if err := rs.Append(first, run.Entry{Type: run.TypeProcess, At: now, Process: &run.Process{Host: host + "-elsewhere", Pid: 1}}); err != nil {
		t.Fatal(err)
	}
	if waits() {
		t.Fatal("first run live: second waits, want blocked")
	}
	// Its run ends: the request stays in the queue with its run set, and
	// the second is demand again.
	if err := rs.Append(first, run.Entry{Type: run.TypeFinished, At: now, End: run.Succeeded}); err != nil {
		t.Fatal(err)
	}
	if !waits() {
		t.Fatal("first run finished: second blocked, want it to wait for a slot")
	}
}

// A queued milestone request waits only while another milestone request's
// plan run holds the owner: live at the interview step, stopped at a
// sign-off (or gate) wait, or its start unjournaled within daemonStartGrace.
// Live at any later step it blocks nothing.
func TestDaemonMilestoneWaitsOwnerSteps(t *testing.T) {
	f := newSessionFixture(t, "req")
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	f.ctx.Now = func() time.Time { return now }
	host, _ := os.Hostname()
	rs := run.NewStore(run.Dir(f.common))
	const first = "01M6YRUNOWNER00000000000000"
	write := func(startedAt time.Time) {
		t.Helper()
		if err := daemon.WriteRequests(f.common, []daemon.Request{
			{ID: "r1", Kind: daemon.KindMilestone, Text: "First", At: now.Add(-time.Hour), Run: first, StartedAt: startedAt},
			{ID: "r2", Kind: daemon.KindMilestone, Text: "Second", At: now.Add(-time.Minute)},
		}); err != nil {
			t.Fatal(err)
		}
	}
	waits := func() bool {
		t.Helper()
		runs, err := rs.List()
		if err != nil {
			t.Fatal(err)
		}
		return daemonMilestoneWaits(f.ctx, rs, runs, host)
	}
	append1 := func(e run.Entry) {
		t.Helper()
		e.At = now
		if err := rs.Append(first, e); err != nil {
			t.Fatal(err)
		}
	}

	// Unjournaled within the grace: the second waits.
	write(now.Add(-time.Second))
	if waits() {
		t.Fatal("unjournaled within the grace: second can start, want it to wait")
	}
	// Journaled, live (another host) at the interview: still waits.
	if err := rs.Create(first, run.Entry{Type: run.TypeStarted, At: now, Started: &run.Started{Workflow: "plan-milestone"}}); err != nil {
		t.Fatal(err)
	}
	append1(run.Entry{Type: run.TypeProcess, Process: &run.Process{Host: host + "-elsewhere", Pid: 1}})
	append1(run.Entry{Type: run.TypeEntered, Step: "interview"})
	if waits() {
		t.Fatal("live at the interview: second can start, want it to wait")
	}
	// Live at a later step: the second is demand.
	for _, step := range []string{"milestone-check", "review"} {
		append1(run.Entry{Type: run.TypeEntered, Step: step})
		if !waits() {
			t.Fatalf("live at %s: second blocked, want it to start", step)
		}
	}
	// Stopped at the gate, waiting on the owner: it holds the owner again.
	append1(run.Entry{Type: run.TypeEntered, Step: "gate"})
	append1(run.Entry{Type: run.TypeStopped, Detail: "waiting for the owner"})
	if waits() {
		t.Fatal("stopped at the gate: second can start, want it to wait")
	}
}
