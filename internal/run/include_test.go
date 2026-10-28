package run

import (
	"reflect"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/workflow"
)

const includedInner = `
name = "inner"
[step.one]
kind = "command"
run = "true"
on.pass = "two"
[step.two]
kind = "command"
run = "true"
on.pass = "three"
[step.three]
kind = "command"
run = "true"
`

const includingOuter = `
name = "wf"
[step.w]
kind = "workflow"
workflow = "inner"
skip = ["two"]
on.pass = "after"
[step.after]
kind = "command"
run = "true"
`

func setupIncluding(t *testing.T) (*Engine, *fakeRunner) {
	t.Helper()
	def, err := workflow.ParseWith("wf", []byte(includingOuter), func(string) ([]byte, error) { return []byte(includedInner), nil })
	if err != nil {
		t.Fatal(err)
	}
	store := newStore(t)
	s := Started{Workflow: "wf", Definition: includingOuter, Includes: def.Includes, Change: "CS1", Title: "Including"}
	if err := Start(store, runID, s, t0, def); err != nil {
		t.Fatal(err)
	}
	r := newFakeRunner()
	r.store = store
	c := &fakeClock{now: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
	return &Engine{Store: store, Runner: r, Clock: c, Now: func() time.Time { return at(1) }}, r
}

// A workflow step's inner steps run in the same journal under the inner
// workflow's name; the skipped one does not run, and the outer on.pass is
// taken.
func TestWorkflowStepRunsInnerSteps(t *testing.T) {
	e, r := setupIncluding(t)
	st, err := e.Run(runID)
	if err != nil || st.Result != Succeeded {
		t.Fatalf("run: %v, %q (%s)", err, st.Result, st.Detail)
	}
	want := []string{"inner/one:pass", "inner/three:pass", "after:pass"}
	if got := outcomes(st); !reflect.DeepEqual(got, want) {
		t.Errorf("outcomes = %v, want %v", got, want)
	}
	if r.commands["inner/two"] != 0 {
		t.Error("skipped step ran")
	}
}

// An inner step failing fails the outer step, and a retry resumes at that
// inner step.
func TestWorkflowStepInnerFailureResumesThere(t *testing.T) {
	e, r := setupIncluding(t)
	r.codes["inner/three"] = []int{1, 0}
	st, err := e.Run(runID)
	if err != nil || st.Result != Failed {
		t.Fatalf("first run: %v, %q", err, st.Result)
	}
	if st.Detail == "" {
		t.Error("failed run has no detail")
	}
	st, err = e.Retry(runID)
	if err != nil || st.Result != Succeeded {
		t.Fatalf("retry: %v, %q (%s)", err, st.Result, st.Detail)
	}
	if r.commands["inner/one"] != 1 || r.commands["inner/three"] != 2 || r.commands["after"] != 1 {
		t.Errorf("commands = %v", r.commands)
	}
}
