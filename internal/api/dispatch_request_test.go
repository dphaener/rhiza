package api

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
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
	// The workflows the rules name must resolve, or the daemon starts no run.
	plan, err := os.ReadFile(filepath.Join("..", "workflow", "builtin", "plan.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(f.ctx.Root, ".rhiza", "workflows"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"bugfix", "roadmap"} {
		body := strings.Replace(string(plan), `name = "plan"`, `name = "`+name+`"`, 1)
		if err := os.WriteFile(filepath.Join(f.ctx.Root, ".rhiza", "workflows", name+".toml"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
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
		if _, err := daemonPlanRequests(f.ctx, rs, nil, 0, 5, "auto", onlyKind, nil); err != nil {
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

// The ready workflow is the change set's recorded workflow, else the [ready]
// rule for its kind, else ship-change.
func TestReadyWorkflowOrder(t *testing.T) {
	root := t.TempDir()
	bug := changeset.Open{Kind: changeset.KindIntention, RequestKind: changeset.RequestBug}
	got := func(o changeset.Open) string {
		t.Helper()
		name, err := readyWorkflow(root, o)
		if err != nil {
			t.Fatal(err)
		}
		return name
	}
	if g := got(bug); g != "ship-change" {
		t.Errorf("no file: %q, want ship-change", g)
	}
	if err := os.MkdirAll(filepath.Join(root, ".rhiza"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".rhiza", "dispatch.toml"), []byte("[ready]\nbug = \"bugfix\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if g := got(bug); g != "bugfix" {
		t.Errorf("rule: %q, want bugfix", g)
	}
	if g := got(changeset.Open{Kind: changeset.KindIntention}); g != "ship-change" {
		t.Errorf("feature: %q, want ship-change", g)
	}
	bug.Workflow = "hotfix"
	if g := got(bug); g != "hotfix" {
		t.Errorf("recorded: %q, want hotfix over the rule", g)
	}
}

// An explicit run start --workflow wins over the change set's recorded
// workflow, the [ready] rule for its kind and the ship-change default.
func TestStartWorkflowExplicitWins(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".rhiza"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".rhiza", "dispatch.toml"), []byte("[ready]\nbug = \"bugfix\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name string
		open changeset.Open
	}{
		{"recorded", changeset.Open{Kind: changeset.KindIntention, RequestKind: changeset.RequestBug, Workflow: "hotfix"}},
		{"rule", changeset.Open{Kind: changeset.KindIntention, RequestKind: changeset.RequestBug}},
		{"default", changeset.Open{Kind: changeset.KindIntention}},
	} {
		got, err := startWorkflow(root, "other", c.open)
		if err != nil || got != "other" {
			t.Errorf("%s: %q, %v; want the explicit other", c.name, got, err)
		}
		want, _ := readyWorkflow(root, c.open)
		if got, err := startWorkflow(root, " ", c.open); err != nil || got != want {
			t.Errorf("%s with no --workflow: %q, %v; want the ready workflow %q", c.name, got, err, want)
		}
	}
}

// run start on an existing change set runs the [ready] rule's workflow with
// no --workflow, and the explicit --workflow over it.
func TestRunStartExplicitWorkflowOverRule(t *testing.T) {
	f := newSessionFixture(t, "explicit")
	newChange(t, f.ctx, `{"intent":"Plain work","title":"Plain work"}`)
	writeWorkflow(t, f.ctx.Root, "bugfix", "name = \"bugfix\"\n"+twoCommands)
	writeWorkflow(t, f.ctx.Root, "other", "name = \"other\"\n"+twoCommands)
	if err := os.WriteFile(filepath.Join(f.ctx.Root, ".rhiza", "dispatch.toml"), []byte("[ready]\nfeature = \"bugfix\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	f.ctx.DryRun = true
	for _, c := range []struct{ req, want string }{
		{`{"foreground":true}`, "bugfix"},
		{`{"foreground":true,"workflow":"other"}`, "other"},
	} {
		resp, err := callOp(t, f.ctx, "run.start", c.req)
		if err != nil {
			t.Fatalf("%s: %v", c.req, err)
		}
		if got := resp.(RunResponse).Run.Workflow; got != c.want {
			t.Errorf("%s: workflow %q, want %q", c.req, got, c.want)
		}
	}
}
