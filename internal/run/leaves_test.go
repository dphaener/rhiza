package run

import (
	"errors"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/workflow"
)

func (f *fakeRunner) Leaves(change string) ([]Leaf, error) {
	if f.leavesErr != nil {
		return nil, f.leavesErr
	}
	return f.leaves, nil
}

func (f *fakeRunner) StartLeaf(spec LeafSpec) (string, error) {
	if err := f.leafErr[spec.Leaf.Handle+" "+spec.Role]; err != nil {
		return "", err
	}
	f.leafSpecs = append(f.leafSpecs, spec)
	return leafSID(spec.Leaf.Handle, spec.Role, spec.Attempt), nil
}

// leafSID is the fake runner's session ID for a leaf session.
func leafSID(handle, role string, attempt int) string {
	return "S-" + handle + "-" + role + "-" + strconv.Itoa(attempt)
}

const leavesWF = `[step.build]
kind = "leaves"
interval = "5s"
[step.build.builder]
prompt = "build {leaf}"
[step.build.verifier]
prompt = "verify {leaf}"
`

var twoLeaves = []Leaf{{ID: "L1", Handle: "leaf-a"}, {ID: "L2", Handle: "leaf-b", DependsOn: []string{"L1"}}}

// passAll makes every leaf session of leaves exit 0.
func passAll(r *fakeRunner, leaves []Leaf) {
	for _, l := range leaves {
		for _, role := range []string{workflow.RoleBuilder, workflow.RoleVerifier} {
			r.outcome[leafSID(l.Handle, role, 1)] = workflow.OutcomePass
		}
	}
}

// leafEntries are the journal's leaf entries as "type leaf role attempt
// outcome".
func leafEntries(t *testing.T, e *Engine) []string {
	t.Helper()
	j, err := e.Store.Read(runID)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, en := range j.Entries {
		switch en.Type {
		case TypeLeafSession, TypeLeafAttempt, TypeLeaf:
			out = append(out, strings.Join(strings.Fields(strings.Join([]string{en.Type, en.Handle, en.Role, attemptStr(en.Attempt), en.Outcome}, " ")), " "))
		}
	}
	return out
}

func attemptStr(n int) string {
	if n == 0 {
		return ""
	}
	return strconv.Itoa(n)
}

// A leaves step over two leaves journals a session, an attempt and an
// outcome per leaf and role, builder before verifier and leaf by leaf, and
// passes when both leaves pass.
func TestLeavesStepTwoLeaves(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = twoLeaves
	passAll(r, twoLeaves)
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	want := []string{
		"leaf-session leaf-a builder 1", "leaf-attempt leaf-a builder 1 pass",
		"leaf-session leaf-a verifier 1", "leaf-attempt leaf-a verifier 1 pass",
		"leaf leaf-a passed",
		"leaf-session leaf-b builder 1", "leaf-attempt leaf-b builder 1 pass",
		"leaf-session leaf-b verifier 1", "leaf-attempt leaf-b verifier 1 pass",
		"leaf leaf-b passed",
	}
	if got := leafEntries(t, e); !reflect.DeepEqual(got, want) {
		t.Errorf("leaf entries =\n%v\nwant\n%v", got, want)
	}
	if len(st.Leaves) != 2 || st.Leaves[0].Handle != "leaf-a" || st.Leaves[1].Handle != "leaf-b" {
		t.Fatalf("Leaves = %+v", st.Leaves)
	}
	for _, l := range st.Leaves {
		if l.Outcome != LeafPassed || len(l.Sessions) != 2 || len(l.Attempts) != 2 || l.Step != "build" {
			t.Errorf("leaf %s = %+v", l.Handle, l)
		}
	}
	if st.Steps[0].Detail != "2 of 2 leaves passed" {
		t.Errorf("detail = %q", st.Steps[0].Detail)
	}
	spec := r.leafSpecs[0]
	if spec.Step.Name != "build" || spec.Role != workflow.RoleBuilder || spec.Attempt != 1 || spec.Leaf.ID != "L1" {
		t.Errorf("first spec = %+v", spec)
	}
	for _, kv := range []string{"RHIZA_STEP=build", "RHIZA_CHANGE=CS1", "RHIZA_LEAF=leaf-a", "RHIZA_ATTEMPT=1"} {
		if !slices.Contains(spec.Env, kv) {
			t.Errorf("env %v lacks %s", spec.Env, kv)
		}
	}
	if want := []string{"build:" + leafSID("leaf-a", "builder", 1), "build:" + leafSID("leaf-a", "verifier", 1),
		"build:" + leafSID("leaf-b", "builder", 1), "build:" + leafSID("leaf-b", "verifier", 1)}; !reflect.DeepEqual(r.ended, want) {
		t.Errorf("AgentEnded = %v, want %v", r.ended, want)
	}
	if len(st.Sessions) != 4 {
		t.Errorf("Sessions = %v, want the 4 leaf sessions", st.Sessions)
	}
}

// A builder that fails fails its leaf with no verifier; the next leaf
// still runs, and the step fails naming the leaf that did not pass.
func TestLeavesStepBuilderFails(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = twoLeaves
	passAll(r, twoLeaves)
	r.outcome[leafSID("leaf-a", "builder", 1)] = workflow.OutcomeFail
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:fail"}) {
		t.Fatalf("outcomes = %v", got)
	}
	want := []string{
		"leaf-session leaf-a builder 1", "leaf-attempt leaf-a builder 1 fail", "leaf leaf-a failed",
		"leaf-session leaf-b builder 1", "leaf-attempt leaf-b builder 1 pass",
		"leaf-session leaf-b verifier 1", "leaf-attempt leaf-b verifier 1 pass", "leaf leaf-b passed",
	}
	if got := leafEntries(t, e); !reflect.DeepEqual(got, want) {
		t.Errorf("leaf entries =\n%v\nwant\n%v", got, want)
	}
	if d := st.Steps[0].Detail; d != "1 of 2 leaves passed; not passed: leaf-a" {
		t.Errorf("detail = %q", d)
	}
}

// A verifier's fail fails the leaf.
func TestLeavesStepVerifierFails(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = twoLeaves[:1]
	passAll(r, twoLeaves)
	r.outcome[leafSID("leaf-a", "verifier", 1)] = workflow.OutcomeFail
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:fail"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if l := st.Leaves[0]; l.Outcome != LeafFailed || !strings.Contains(l.Detail, "verifier failed") {
		t.Errorf("leaf = %+v", l)
	}
}

// An attempt is judged when the agent says done, never on idle alone; the
// conductor sleeps the step's interval between polls that judge nothing.
func TestLeavesStepSaidDoneNotIdle(t *testing.T) {
	e, r, c, _ := setup(t, leavesWF)
	r.leaves = twoLeaves[:1]
	passAll(r, twoLeaves)
	b := leafSID("leaf-a", "builder", 1)
	r.agentSeq[b] = []AgentStatus{{Idle: true}, {Idle: true}, {SaidDone: true}}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if r.polls[b] != 3 {
		t.Errorf("builder polled %d times, want 3", r.polls[b])
	}
	if a, _ := st.Leaves[0].Judged(workflow.RoleBuilder, 1); a.Outcome != workflow.OutcomePass || a.Detail != "said done" {
		t.Errorf("builder judged %+v", a)
	}
	if !reflect.DeepEqual(c.sleeps, []time.Duration{5 * time.Second, 5 * time.Second}) {
		t.Errorf("sleeps = %v", c.sleeps)
	}
}

// A resume rejoins the leaf session the journal holds and starts nothing
// the journal already has.
func TestLeavesStepResumeRejoins(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = twoLeaves
	passAll(r, twoLeaves)
	for _, en := range []Entry{
		{Type: TypeEntered, Step: "build", Visit: 1},
		{Type: TypeLeafSession, Step: "build", Leaf: "L1", Handle: "leaf-a", Role: "builder", Attempt: 1, Session: leafSID("leaf-a", "builder", 1)},
		{Type: TypeLeafAttempt, Step: "build", Leaf: "L1", Handle: "leaf-a", Role: "builder", Attempt: 1, Outcome: "pass"},
		{Type: TypeLeafSession, Step: "build", Leaf: "L1", Handle: "leaf-a", Role: "verifier", Attempt: 1, Session: leafSID("leaf-a", "verifier", 1)},
	} {
		en.At = at(1)
		if err := e.Store.Append(runID, en); err != nil {
			t.Fatal(err)
		}
	}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	var started []string
	for _, s := range r.leafSpecs {
		started = append(started, s.Leaf.Handle+" "+s.Role)
	}
	if want := []string{"leaf-b builder", "leaf-b verifier"}; !reflect.DeepEqual(started, want) {
		t.Errorf("started %v, want %v", started, want)
	}
	if r.polls[leafSID("leaf-a", "verifier", 1)] != 1 {
		t.Errorf("leaf-a's verifier was not rejoined")
	}
}

// A blocked leaf session is journalled and notified once, and its block
// ends when the session ends.
func TestLeavesStepBlocked(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = twoLeaves[:1]
	passAll(r, twoLeaves)
	b := leafSID("leaf-a", "builder", 1)
	r.agentSeq[b] = []AgentStatus{{Blocked: &Block{Kind: "permission"}}, {Blocked: &Block{Kind: "permission"}}}
	var notices []string
	e.Notify = func(n Notice) { notices = append(notices, n.Step+" "+n.Session+" "+strconv.FormatBool(n.Blocked)) }
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if want := []string{"build " + b + " true", "build " + b + " false"}; !reflect.DeepEqual(notices, want) {
		t.Errorf("notices = %v, want %v", notices, want)
	}
	types := entryTypes(t, e)
	if !slices.Contains(types, "blocked:build") || !slices.Contains(types, "unblocked:build") {
		t.Errorf("entries %v lack the block", types)
	}
}

// A runner that cannot run leaves, or fails to list or start them, ends
// the step with an error outcome.
func TestLeavesStepRunnerErrors(t *testing.T) {
	e, _, _, _ := setup(t, leavesWF)
	e.Runner = struct{ Runner }{newFakeRunner()}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:error"}) || !strings.Contains(st.Steps[0].Detail, "cannot run leaves") {
		t.Errorf("no LeafRunner: %v %q", got, st.Steps[0].Detail)
	}

	e, r, _, _ := setup(t, leavesWF)
	r.leavesErr = errors.New("no change set")
	if st, _ = e.Run(runID); outcomes(st)[0] != "build:error" || st.Steps[0].Detail != "no change set" {
		t.Errorf("Leaves error: %v %q", outcomes(st), st.Steps[0].Detail)
	}

	e, r, _, _ = setup(t, leavesWF)
	r.leaves = twoLeaves
	r.leafErr = map[string]error{"leaf-a builder": errors.New("rex is gone")}
	if st, _ = e.Run(runID); outcomes(st)[0] != "build:error" || st.Steps[0].Detail != "leaf leaf-a: start its builder: rex is gone" {
		t.Errorf("StartLeaf error: %v %q", outcomes(st), st.Steps[0].Detail)
	}
}

// A new visit of a leaves step starts its leaves anew.
func TestLeavesFoldNewVisit(t *testing.T) {
	j := Journal{ID: runID, Entries: []Entry{
		{Type: TypeEntered, Step: "build", Visit: 1},
		{Type: TypeLeafSession, Step: "build", Leaf: "L1", Handle: "leaf-a", Role: "builder", Attempt: 1, Session: "S1"},
		{Type: TypeLeaf, Step: "build", Leaf: "L1", Outcome: LeafFailed, Detail: "no"},
		{Type: TypeEnded, Step: "build", Outcome: "fail", Next: "build"},
	}}
	st := j.State()
	if len(st.Leaves) != 1 || st.Leaves[0].Outcome != LeafFailed || st.Leaves[0].Handle != "leaf-a" {
		t.Fatalf("Leaves = %+v", st.Leaves)
	}
	j.Entries = append(j.Entries, Entry{Type: TypeEntered, Step: "build", Visit: 2})
	if st := j.State(); len(st.Leaves) != 0 {
		t.Errorf("after a new visit Leaves = %+v, want none", st.Leaves)
	}
}

// With no leaves open, nextLeafActions ends the step; with one open, it
// moves only that leaf.
func TestNextLeafActions(t *testing.T) {
	got := nextLeafActions(twoLeaves, map[string]LeafState{})
	if len(got) != 1 || got[0].kind != leafStart || got[0].leaf.ID != "L1" || got[0].role != workflow.RoleBuilder || got[0].attempt != 1 {
		t.Errorf("fresh: %+v", got)
	}
	states := map[string]LeafState{
		"L1": {Leaf: "L1", Outcome: LeafPassed},
		"L2": {Leaf: "L2", Sessions: []LeafSession{{Role: "builder", Attempt: 1, Session: "S2"}}},
	}
	if got := nextLeafActions(twoLeaves, states); len(got) != 1 || got[0].kind != leafPoll || got[0].session != "S2" {
		t.Errorf("polling: %+v", got)
	}
	states["L2"] = LeafState{Leaf: "L2", Outcome: LeafPassed}
	if got := nextLeafActions(twoLeaves, states); len(got) != 1 || got[0].kind != leavesDone || got[0].outcome != workflow.OutcomePass {
		t.Errorf("done: %+v", got)
	}
}

// A leaf session starts on its role's model for the attempt, journaled on
// its leaf-session entry; a role with no models names none.
func TestLeavesModel(t *testing.T) {
	src := strings.Replace(leavesWF, "prompt = \"build {leaf}\"\n", "prompt = \"build {leaf}\"\nmodels = [\"sonnet\", \"opus\"]\n", 1)
	e, r, _, _ := setup(t, src)
	r.leaves = twoLeaves[:1]
	passAll(r, r.leaves)
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, s := range r.leafSpecs {
		got = append(got, s.Role+":"+s.Model)
	}
	if want := []string{"builder:sonnet", "verifier:"}; !reflect.DeepEqual(got, want) {
		t.Errorf("specs = %q, want %q", got, want)
	}
	j, err := e.Store.Read(runID)
	if err != nil {
		t.Fatal(err)
	}
	for _, en := range j.Entries {
		if en.Type == TypeLeafSession && en.Role == workflow.RoleBuilder && en.Model != "sonnet" {
			t.Errorf("builder leaf-session model = %q, want sonnet", en.Model)
		}
	}
}
