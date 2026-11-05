package api

import (
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
