package run

import (
	"reflect"
	"strings"
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

const includedMid = `
name = "mid"
[step.inner]
kind = "workflow"
workflow = "inner"
on.pass = "done"
`

const includingWrap = `
name = "wrap"
[step.mid]
kind = "workflow"
workflow = "mid"
on.pass = "done"
`

// A workflow two includes deep runs its steps under paths with two '/': a
// failure stops the run there, its detail names the whole path, and a retry
// resumes at that step, not at the first segment.
func TestNestedWorkflowStepResumesAtFullPath(t *testing.T) {
	read := func(name string) ([]byte, error) {
		if name == "mid" {
			return []byte(includedMid), nil
		}
		return []byte(includedInner), nil
	}
	def, err := workflow.ParseWith("wrap", []byte(includingWrap), read)
	if err != nil {
		t.Fatal(err)
	}
	store := newStore(t)
	s := Started{Workflow: "wrap", Definition: includingWrap, Includes: def.Includes, Change: "CS1", Title: "Wrapping"}
	if err := Start(store, runID, s, t0, def); err != nil {
		t.Fatal(err)
	}
	r := newFakeRunner()
	r.store = store
	c := &fakeClock{now: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
	e := &Engine{Store: store, Runner: r, Clock: c, Now: func() time.Time { return at(1) }}

	const nested = "mid/inner/two"
	if _, ok := def.Steps[nested]; !ok {
		t.Fatalf("flat steps have no %s: %v", nested, def.Order)
	}
	r.codes[nested] = []int{1, 0}
	st, err := e.Run(runID)
	if err != nil || st.Result != Failed {
		t.Fatalf("first run: %v, %q", err, st.Result)
	}
	if got, want := outcomes(st), []string{"mid/inner/one:pass", nested + ":fail"}; !reflect.DeepEqual(got, want) {
		t.Errorf("failed run outcomes = %v, want %v", got, want)
	}
	if !strings.Contains(st.Detail, nested) {
		t.Errorf("failure detail %q does not name %s", st.Detail, nested)
	}
	st, err = e.Retry(runID)
	if err != nil || st.Result != Succeeded {
		t.Fatalf("retry: %v, %q (%s)", err, st.Result, st.Detail)
	}
	if r.commands["mid/inner/one"] != 1 || r.commands[nested] != 2 || r.commands["mid/inner/three"] != 1 {
		t.Errorf("commands = %v", r.commands)
	}
	want := []string{"mid/inner/one:pass", nested + ":fail", nested + ":pass", "mid/inner/three:pass"}
	if got := outcomes(st); !reflect.DeepEqual(got, want) {
		t.Errorf("outcomes = %v, want %v", got, want)
	}
}

const includedHandback = `
name = "inner"
[step.conduct]
kind = "agent"
prompt = "conduct"
on.pass = "gate"
[step.gate]
kind = "command"
run = "gate"
on.fail = "handback:conduct"
`

const includingHandback = `
name = "wf"
[step.w]
kind = "workflow"
workflow = "inner"
`

// An included step failing hands back to the included agent step and the
// run continues from it.
func TestHandbackIncluded(t *testing.T) {
	def, err := workflow.ParseWith("wf", []byte(includingHandback), func(string) ([]byte, error) { return []byte(includedHandback), nil })
	if err != nil {
		t.Fatal(err)
	}
	store := newStore(t)
	s := Started{Workflow: "wf", Definition: includingHandback, Includes: def.Includes, Change: "CS1", Title: "Handing back"}
	if err := Start(store, runID, s, t0, def); err != nil {
		t.Fatal(err)
	}
	r := newFakeRunner()
	r.store = store
	r.outcome["S-inner/conduct"] = workflow.OutcomePass
	r.codes["inner/gate"] = []int{1, 0}
	c := &fakeClock{now: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
	e := &Engine{Store: store, Runner: r, Clock: c, Now: func() time.Time { return at(1) }}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"inner/conduct:pass", "inner/gate:fail", "inner/conduct:pass", "inner/gate:pass"}
	if st.Result != Succeeded || !reflect.DeepEqual(outcomes(st), want) {
		t.Fatalf("result %q steps %v", st.Result, outcomes(st))
	}
	if len(st.Handbacks) != 1 || st.Handbacks[0].Target != "inner/conduct" || st.Handbacks[0].Step != "inner/gate" {
		t.Errorf("handbacks %+v", st.Handbacks)
	}
}
