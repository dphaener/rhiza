package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// run milestone-check passes when the plan's change set serves a milestone
// directly or through an intent under one, and fails, naming the change
// set and the fix, when it serves none: no intent at all, or an intent
// under no milestone.
func TestRunMilestoneCheck(t *testing.T) {
	ctx, _, _, _ := openPRFixture(t, false)
	ctx.IndexPath = filepath.Join(t.TempDir(), "index.db")
	writePolicy(t, ctx.Root, "[schedule]\nbacklog = \"skip\"\n")
	cs := newChange(t, ctx, `{"intent":"Plan it","title":"Plan it"}`).ChangeSet.String()
	store := run.NewStore(run.Dir(filepath.Join(ctx.Root, ".git")))
	err := store.Create("R2", run.Entry{
		Type: run.TypeStarted, Step: "milestone-check", At: time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC),
		Started: &run.Started{Workflow: "plan", Change: cs, Title: "Plan it"},
	})
	if err != nil {
		t.Fatal(err)
	}
	runID := "R2"
	ctx.Getenv = func(k string) string {
		if k == envRun {
			return runID
		}
		return ""
	}
	call := func(name, req string) {
		t.Helper()
		if _, err := callOp(t, ctx, name, req); err != nil {
			t.Fatalf("%s %s: %v", name, req, err)
		}
	}
	wantFail := func(what string) {
		t.Helper()
		_, err := callOp(t, ctx, "run.milestone-check", `{}`)
		if err == nil {
			t.Fatalf("%s: passed; want a failure", what)
		}
		if !strings.Contains(err.Error(), cs) || !strings.Contains(err.Error(), "no milestone") || !strings.Contains(err.Error(), "Fix:") {
			t.Errorf("%s: error %q lacks the change set, \"no milestone\" or the fix", what, err)
		}
	}
	wantPass := func(what string) {
		t.Helper()
		resp, err := callOp(t, ctx, "run.milestone-check", `{}`)
		if err != nil {
			t.Fatalf("%s: %v", what, err)
		}
		if r := resp.(RunMilestoneCheckResponse); len(r.ChangeSets) != 1 || r.ChangeSets[0] != cs {
			t.Errorf("%s: response %+v", what, r)
		}
	}

	for _, req := range []string{
		`{"kind":"intent","handle":"loose","fields":["statement=No milestone above it"]}`,
		`{"kind":"intent","handle":"m1","fields":["statement=A milestone","milestone=1"]}`,
		`{"kind":"intent","handle":"point","parent":"m1","fields":["statement=A done point"]}`,
	} {
		recordWrite(t, ctx, "add", req)
	}

	wantFail("serving nothing")

	call("link", `{"from":"`+cs+`","type":"serves","to":"loose"}`)
	wantFail("serving an intent under no milestone")

	call("link", `{"from":"`+cs+`","type":"serves","to":"point"}`)
	wantPass("serving an intent under a milestone (and one under none)")

	// A second plan run, its change set serving the milestone directly.
	cs = newChange(t, ctx, `{"intent":"Plan more","title":"Plan more"}`).ChangeSet.String()
	if err := store.Create("R3", run.Entry{
		Type: run.TypeStarted, Step: "milestone-check", At: time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC),
		Started: &run.Started{Workflow: "plan", Change: cs, Title: "Plan more"},
	}); err != nil {
		t.Fatal(err)
	}
	runID = "R3"
	wantFail("a second change set serving nothing")
	call("link", `{"from":"`+cs+`","type":"serves","to":"m1"}`)
	wantPass("serving a milestone directly")

	// A third, serving an intent it adds itself under the milestone: the
	// intent is not landed, so only the change set's what-if view has it.
	cs = newChange(t, ctx, `{"intent":"Plan sub","title":"Plan sub"}`).ChangeSet.String()
	if err := store.Create("R4", run.Entry{
		Type: run.TypeStarted, Step: "milestone-check", At: time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC),
		Started: &run.Started{Workflow: "plan", Change: cs, Title: "Plan sub"},
	}); err != nil {
		t.Fatal(err)
	}
	runID = "R4"
	call("add", `{"kind":"intent","handle":"sub","parent":"m1","fields":["statement=A new done point"]}`)
	call("link", `{"from":"`+cs+`","type":"serves","to":"sub"}`)
	wantPass("serving an unlanded intent of its own under a milestone")
}

func writePolicy(t *testing.T, root, body string) {
	t.Helper()
	p := filepath.Join(root, ".rhiza", "policy.toml")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// Without a policy file (or with backlog = "after") a plan serving no
// milestone passes; started with a milestone, the change set must serve
// that one, directly or through an intent under it, and one serving
// another milestone or none fails naming it.
func TestRunMilestoneCheckBacklogAndNamedMilestone(t *testing.T) {
	ctx, _, _, _ := openPRFixture(t, false)
	ctx.IndexPath = filepath.Join(t.TempDir(), "index.db")
	store := run.NewStore(run.Dir(filepath.Join(ctx.Root, ".git")))
	m1 := recordWrite(t, ctx, "add", `{"kind":"intent","handle":"m1","fields":["statement=First milestone","milestone=1"]}`).Nodes[0].String()
	m2 := recordWrite(t, ctx, "add", `{"kind":"intent","handle":"m2","fields":["statement=Second milestone","milestone=2"]}`).Nodes[0].String()
	recordWrite(t, ctx, "add", `{"kind":"intent","handle":"point","parent":"m1","fields":["statement=A done point"]}`)
	n := 0
	start := func(milestone string, serves ...string) (string, string) {
		t.Helper()
		n++
		cs := newChange(t, ctx, `{"intent":"Plan `+string(rune('a'+n))+`","title":"Plan"}`).ChangeSet.String()
		runID := "RM" + string(rune('A'+n))
		if err := store.Create(runID, run.Entry{
			Type: run.TypeStarted, Step: "milestone-check", At: time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC),
			Started: &run.Started{Workflow: "plan", Change: cs, Title: "Plan", Milestone: milestone},
		}); err != nil {
			t.Fatal(err)
		}
		for _, to := range serves {
			if _, err := callOp(t, ctx, "link", `{"from":"`+cs+`","type":"serves","to":"`+to+`"}`); err != nil {
				t.Fatal(err)
			}
		}
		ctx.Getenv = func(k string) string {
			if k == envRun {
				return runID
			}
			return ""
		}
		return cs, runID
	}
	check := func() error {
		_, err := callOp(t, ctx, "run.milestone-check", `{}`)
		return err
	}

	start("")
	if err := check(); err != nil {
		t.Errorf("no milestone served, no policy file: %v", err)
	}
	writePolicy(t, ctx.Root, "[schedule]\nbacklog = \"before\"\n")
	if err := check(); err != nil {
		t.Errorf("no milestone served, backlog before: %v", err)
	}

	start(m1, "m1")
	if err := check(); err != nil {
		t.Errorf("serving the named milestone directly: %v", err)
	}
	start(m1, "point")
	if err := check(); err != nil {
		t.Errorf("serving the named milestone through an intent: %v", err)
	}
	cs, _ := start(m1, "m2")
	if err := check(); err == nil || !strings.Contains(err.Error(), cs) || !strings.Contains(err.Error(), m1) || !strings.Contains(err.Error(), "Fix:") {
		t.Errorf("serving another milestone: error %v lacks the change set, the milestone or the fix", err)
	}
	cs, _ = start(m2)
	if err := check(); err == nil || !strings.Contains(err.Error(), cs) || !strings.Contains(err.Error(), m2) {
		t.Errorf("serving none under a named milestone: error %v", err)
	}
}
