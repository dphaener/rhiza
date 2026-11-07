package api

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// The daemon starts a queued request's plan run with the workflow the
// dispatch file's [request] table names for its kind: plan by default, and
// a milestone's never with --autonomous.
func TestDaemonStartsRequestWithDispatchedWorkflow(t *testing.T) {
	f := newSessionFixture(t, "disp")
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	f.ctx.Now = func() time.Time { return now }
	f.ctx.repoSession = true
	f.ctx.Self = "/bin/rhiza"
	var detached [][]string
	f.ctx.Detach = func(p Program, logPath string) error {
		detached = append(detached, p.Args)
		return nil
	}
	file := filepath.Join(f.ctx.Root, ".rhiza", "dispatch.toml")
	if err := os.WriteFile(file, []byte("[request]\nbug = \"bugfix\"\nmilestone = \"roadmap\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rs := run.NewStore(run.Dir(f.common))
	started := func(kind string, onlyKind string) []string {
		t.Helper()
		if _, err := callOp(t, f.ctx, "request", `{"text":"Thing","kind":"`+kind+`"}`); err != nil {
			t.Fatal(err)
		}
		detached = nil
		if _, err := daemonPlanRequests(f.ctx, rs, nil, 0, 5, "auto", onlyKind); err != nil {
			t.Fatal(err)
		}
		if len(detached) != 1 {
			t.Fatalf("%s: %d starts, want one", kind, len(detached))
		}
		return detached[0]
	}
	flag := func(args []string, name string) string {
		i := slices.Index(args, name)
		if i < 0 || i+1 >= len(args) {
			return ""
		}
		return args[i+1]
	}
	bug := started("bug", "")
	if flag(bug, "--workflow") != "bugfix" || !slices.Contains(bug, "--new") {
		t.Errorf("bug start %q, want --workflow bugfix --new", bug)
	}
	feature := started("feature", "")
	if flag(feature, "--workflow") != "plan" {
		t.Errorf("feature start %q, want --workflow plan", feature)
	}
	ms := started("milestone", "milestone")
	if flag(ms, "--workflow") != "roadmap" || !slices.Contains(ms, "--milestone") || slices.Contains(ms, "--autonomous") {
		t.Errorf("milestone start %q, want --workflow roadmap --milestone without --autonomous", ms)
	}
}

// A run under a workflow not named plan, started with --new, is a plan run:
// run resume's milestone owner gate stops it when autonomous.
func TestRunResumeGatesMilestonePlanUnderOtherWorkflow(t *testing.T) {
	f := newSessionFixture(t, "gate")
	cs := newChange(t, f.ctx, `{"intent":"Roadmap","title":"Roadmap","kind":"record"}`).ChangeSet
	rs := run.NewStore(run.Dir(f.common))
	create := func(id string, plan bool) {
		t.Helper()
		s := run.Started{Workflow: "roadmap", Change: cs.String(), Definition: twoCommands, Autonomous: true, Plan: plan}
		err := rs.Create(id, run.Entry{Type: run.TypeStarted, At: time.Now(), Step: "one", Started: &s})
		if err != nil {
			t.Fatal(err)
		}
	}
	create("RPLAN", true)
	_, err := callOp(t, f.ctx, "run.resume", `{"foreground":true,"run":"RPLAN"}`)
	if err == nil || !strings.Contains(err.Error(), "milestone planning needs the owner") {
		t.Fatalf("resume of a plan run under workflow roadmap: %v, want the owner gate", err)
	}
	create("ROTHER", false)
	_, err = callOp(t, f.ctx, "run.resume", `{"foreground":true,"run":"ROTHER"}`)
	if err != nil && strings.Contains(err.Error(), "milestone planning needs the owner") {
		t.Fatalf("resume of a non-plan run: %v, want no owner gate", err)
	}
}

// The inbox does not list a landed change set's interrupted plan run as
// skipped, whatever workflow it runs (the daemon resumes it); a non-plan run
// that never reached land is listed as skipped.
func TestInboxKeepsPlanRunUnderOtherWorkflow(t *testing.T) {
	landed := func(string) (string, bool) { return "landed", true }
	items := func(plan bool) int {
		j := run.Journal{ID: "R1", Entries: []run.Entry{
			{Type: run.TypeStarted, At: time.Now(), Step: "one", Started: &run.Started{Workflow: "bugfix", Change: "CS1", Plan: plan}},
			{Type: run.TypeProcess, Process: &run.Process{Host: "h", Pid: 1}},
		}}
		return len(inboxItems([]run.Journal{j}, "h", func(int) bool { return false }, func(string) bool { return false }, landed, nil, nil))
	}
	if got := items(true); got != 0 {
		t.Errorf("plan run under bugfix: %d items, want 0 (resumable, not skipped)", got)
	}
	if got := items(false); got != 1 {
		t.Errorf("non-plan run under bugfix that never landed: %d items, want 1 skipped", got)
	}
}

// run resume lets a request's plan run through when its change set landed
// early (the interview lands its record), whatever workflow it runs; a run
// not started for a request, which never reached land, is still refused.
func TestRunResumeLetsLandedPlanRunUnderOtherWorkflow(t *testing.T) {
	f := newSessionFixture(t, "landed")
	cs := newChange(t, f.ctx, `{"intent":"Fix the bug","title":"Fix the bug"}`).ChangeSet
	for _, s := range []string{"propose", "accept", "land"} {
		move(t, f.ctx, s, `{}`)
	}
	rs := run.NewStore(run.Dir(f.common))
	create := func(id string, plan bool) {
		t.Helper()
		s := run.Started{Workflow: "bugfix", Change: cs.String(), Definition: twoCommands, Plan: plan}
		if err := rs.Create(id, run.Entry{Type: run.TypeStarted, At: time.Now(), Step: "one", Started: &s}); err != nil {
			t.Fatal(err)
		}
	}
	f.ctx.DryRun = true
	create("RPLAN", true)
	resp, err := callOp(t, f.ctx, "run.resume", `{"foreground":true,"run":"RPLAN"}`)
	if err != nil {
		t.Fatalf("resume of a plan run under bugfix on a landed change set: %v, want it resumed", err)
	}
	if rr := resp.(RunResponse); rr.Run.ID != "RPLAN" || !rr.DryRun {
		t.Errorf("resume response %+v", rr)
	}
	create("ROTHER", false)
	_, err = callOp(t, f.ctx, "run.resume", `{"foreground":true,"run":"ROTHER"}`)
	if err == nil || !strings.Contains(err.Error(), "no longer accepted") {
		t.Fatalf("resume of a non-plan run on a landed change set: %v, want it refused", err)
	}
}
