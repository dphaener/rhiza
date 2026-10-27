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

// fakeHead is the fake runner's worktree HEAD.
const fakeHead = "H0"

func (f *fakeRunner) Head() (string, error) {
	if len(f.heads) > 0 {
		h := f.heads[0]
		f.heads = f.heads[1:]
		return h, nil
	}
	return fakeHead, nil
}

func (f *fakeRunner) LeafCommit(leaf Leaf, since string) (string, string, error) {
	f.leafSince = append(f.leafSince, leaf.Handle+" "+since)
	if c, ok := f.leafCommit[leaf.Handle]; ok {
		return c[0], c[1], nil
	}
	return "C-" + leaf.Handle, "", nil
}

func (f *fakeRunner) Evidence(leaf Leaf, session string) (bool, string, error) {
	f.evidenceCalls = append(f.evidenceCalls, leaf.Handle+" "+session)
	if ev, ok := f.evidence[leaf.Handle]; ok {
		return ev.ok, ev.detail, ev.err
	}
	return true, "", nil
}

func (f *fakeRunner) CommitEvidence(leaf Leaf, session string) (string, error) {
	f.commitCalls = append(f.commitCalls, leaf.Handle+" "+session)
	if err := f.evidenceCommitErr[leaf.Handle]; err != nil {
		return "", err
	}
	return "E-" + leaf.Handle, nil
}

func (f *fakeRunner) AddLeafQuestion(leaf Leaf, text string) (Question, error) {
	if f.leafQuestionErr != nil {
		return Question{}, f.leafQuestionErr
	}
	f.leafQuestions = append(f.leafQuestions, leaf.Handle+" "+text)
	n := strconv.Itoa(len(f.leafQuestions))
	return Question{ID: "QL" + n, Handle: "leaf-failed-" + n}, nil
}

func (f *fakeRunner) LeafAnswer(question string) (string, string, bool, error) {
	f.leafAnswerN[question]++
	if a, ok := f.leafAnswers[question]; ok {
		return a[0], a[1], true, nil
	}
	return "", "", false, nil
}

// checkoutDir is the fake runner's checkout for a leaf's verifier.
func checkoutDir(handle string, attempt int) string {
	return "/checkouts/" + handle + "-" + strconv.Itoa(attempt)
}

func (f *fakeRunner) OpenCheckout(leaf Leaf, attempt int, commit string) (string, error) {
	if f.checkoutErr != nil {
		return "", f.checkoutErr
	}
	dir := checkoutDir(leaf.Handle, attempt)
	f.checkouts = append(f.checkouts, "open "+dir+" "+commit)
	return dir, nil
}

func (f *fakeRunner) RemoveCheckout(dir string) error {
	f.checkouts = append(f.checkouts, "remove "+dir)
	return f.removeErr
}

func (f *fakeRunner) ClosePane(session string) error {
	if f.closeErr != nil {
		return f.closeErr
	}
	f.closed = append(f.closed, session)
	return nil
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

// builderFails are the leaf entries of a leaf whose builder fails each of
// its MaxAttempts attempts.
func builderFails(handle string) []string {
	var out []string
	for n := 1; n <= MaxAttempts; n++ {
		out = append(out, "leaf-session "+handle+" builder "+strconv.Itoa(n), "leaf-attempt "+handle+" builder "+strconv.Itoa(n)+" fail")
	}
	return append(out, "leaf "+handle+" failed")
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

// A builder that fails each attempt fails its leaf with no verifier, its
// retries going before the next leaf; the next leaf, independent of it,
// still runs, and the step fails naming the leaf that did not pass.
func TestLeavesStepBuilderFails(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = []Leaf{twoLeaves[0], {ID: "L2", Handle: "leaf-b"}}
	passAll(r, twoLeaves)
	r.outcome[leafSID("leaf-a", "builder", 1)] = workflow.OutcomeFail
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:fail"}) {
		t.Fatalf("outcomes = %v", got)
	}
	want := append(builderFails("leaf-a"),
		"leaf-session leaf-b builder 1", "leaf-attempt leaf-b builder 1 pass",
		"leaf-session leaf-b verifier 1", "leaf-attempt leaf-b verifier 1 pass", "leaf leaf-b passed",
	)
	if got := leafEntries(t, e); !reflect.DeepEqual(got, want) {
		t.Errorf("leaf entries =\n%v\nwant\n%v", got, want)
	}
	if d := st.Steps[0].Detail; d != "1 of 2 leaves passed; not passed: leaf-a" {
		t.Errorf("detail = %q", d)
	}
}

// A verifier's fail on the last attempt fails the leaf.
func TestLeavesStepVerifierFails(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = twoLeaves[:1]
	for n := 1; n <= MaxAttempts; n++ {
		r.outcome[leafSID("leaf-a", "builder", n)] = workflow.OutcomePass
		r.outcome[leafSID("leaf-a", "verifier", n)] = workflow.OutcomeFail
	}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:fail"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if l := st.Leaves[0]; l.Outcome != LeafFailed || !strings.HasPrefix(l.Detail, "attempt 3: verifier failed") {
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

// A leaf whose dependency failed is skipped, naming the dependency, and
// never started.
func TestLeavesStepDependencyFailedSkips(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = twoLeaves
	passAll(r, twoLeaves)
	r.outcome[leafSID("leaf-a", "builder", 1)] = workflow.OutcomeFail
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	want := append(builderFails("leaf-a"), "leaf leaf-b skipped")
	if got := leafEntries(t, e); !reflect.DeepEqual(got, want) {
		t.Errorf("leaf entries =\n%v\nwant\n%v", got, want)
	}
	if l := st.Leaves[1]; l.Outcome != LeafSkipped || !strings.Contains(l.Detail, "leaf-a") {
		t.Errorf("leaf-b = %+v", l)
	}
	if d := st.Steps[0].Detail; d != "0 of 2 leaves passed; not passed: leaf-a, leaf-b" {
		t.Errorf("detail = %q", d)
	}
}

// abcLeaves are leaves A, B and C in leaf order, C depending on A.
var abcLeaves = []Leaf{{ID: "LA", Handle: "a"}, {ID: "LB", Handle: "b"}, {ID: "LC", Handle: "c", DependsOn: []string{"LA"}}}

// Builders run one at a time in leaf order, while a built leaf's verifier
// runs alongside the next builder: B's builder starts while A's verifier
// is still live, and C, which depends on A, starts only once A passed.
func TestLeafOrderOverlapAndGating(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = abcLeaves
	passAll(r, abcLeaves)
	av := leafSID("a", "verifier", 1)
	// A's verifier runs for five polls; B's builder finishes on its third.
	r.agentOut[av] = []bool{false, false, false, false, false}
	r.agentOut[leafSID("b", "builder", 1)] = []bool{false, false}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	got := leafEntries(t, e)
	bStart, avDone := slices.Index(got, "leaf-session b builder 1"), slices.Index(got, "leaf-attempt a verifier 1 pass")
	avStart := slices.Index(got, "leaf-session a verifier 1")
	if bStart < 0 || avStart < 0 || !(avStart < bStart && bStart < avDone) {
		t.Errorf("B's builder did not start while A's verifier was live:\n%v", got)
	}
	cStart, aPassed := slices.Index(got, "leaf-session c builder 1"), slices.Index(got, "leaf a passed")
	if cStart < 0 || aPassed < 0 || cStart < aPassed {
		t.Errorf("C's builder started before A passed:\n%v", got)
	}
	// One builder live at a time: each builder starts only after the one
	// before it was judged.
	var builders []string
	for _, en := range got {
		if strings.HasPrefix(en, "leaf-session") && strings.Contains(en, "builder") {
			builders = append(builders, en)
		}
	}
	if want := []string{"leaf-session a builder 1", "leaf-session b builder 1", "leaf-session c builder 1"}; !reflect.DeepEqual(builders, want) {
		t.Errorf("builders started %v, want %v", builders, want)
	}
	if slices.Index(got, "leaf-attempt a builder 1 pass") > bStart || slices.Index(got, "leaf-attempt b builder 1 pass") > cStart {
		t.Errorf("a builder started while another was live:\n%v", got)
	}
	seen := map[string]bool{}
	for _, s := range r.leafSpecs {
		if k := s.Leaf.Handle + " " + s.Role; seen[k] {
			t.Errorf("%s started twice", k)
		} else {
			seen[k] = true
		}
	}
	if len(r.leafSpecs) != 6 {
		t.Errorf("started %d sessions, want 6", len(r.leafSpecs))
	}
}

// nextLeafActions polls a live verifier and starts the next ready builder
// in the same turn, holds back a leaf whose dependency is in flight, and
// starts no builder while one is live.
func TestLeafOrderNextActions(t *testing.T) {
	built := LeafState{Leaf: "LA",
		Sessions: []LeafSession{{Role: "builder", Attempt: 1, Session: "SAb"}, {Role: "verifier", Attempt: 1, Session: "SAv"}},
		Attempts: []LeafAttempt{{Role: "builder", Attempt: 1, Outcome: "pass"}}}
	got := nextLeafActions(abcLeaves, map[string]LeafState{"LA": built})
	if len(got) != 2 || got[0].kind != leafPoll || got[0].session != "SAv" ||
		got[1].kind != leafStart || got[1].leaf.ID != "LB" || got[1].role != workflow.RoleBuilder {
		t.Fatalf("verifier live: %+v", got)
	}
	states := map[string]LeafState{"LA": built, "LB": {Leaf: "LB", Sessions: []LeafSession{{Role: "builder", Attempt: 1, Session: "SBb"}}}}
	got = nextLeafActions(abcLeaves, states)
	if len(got) != 2 || got[0].session != "SAv" || got[1].session != "SBb" {
		t.Fatalf("builder live: %+v", got)
	}
	// B passed, A's verifier still live: C waits on A, nothing starts.
	states["LB"] = LeafState{Leaf: "LB", Outcome: LeafPassed}
	got = nextLeafActions(abcLeaves, states)
	if len(got) != 1 || got[0].kind != leafPoll || got[0].session != "SAv" {
		t.Fatalf("C waits: %+v", got)
	}
	states["LA"] = LeafState{Leaf: "LA", Outcome: LeafPassed}
	got = nextLeafActions(abcLeaves, states)
	if len(got) != 1 || got[0].kind != leafStart || got[0].leaf.ID != "LC" {
		t.Fatalf("C ready: %+v", got)
	}
}

// A journal cut while A's verifier and B's builder are live rejoins both
// sessions, starting neither again.
func TestLeafOrderResume(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = abcLeaves
	passAll(r, abcLeaves)
	av, bb := leafSID("a", "verifier", 1), leafSID("b", "builder", 1)
	for _, en := range []Entry{
		{Type: TypeEntered, Step: "build", Visit: 1},
		{Type: TypeLeafSession, Step: "build", Leaf: "LA", Handle: "a", Role: "builder", Attempt: 1, Session: leafSID("a", "builder", 1)},
		{Type: TypeLeafAttempt, Step: "build", Leaf: "LA", Handle: "a", Role: "builder", Attempt: 1, Outcome: "pass"},
		{Type: TypeLeafSession, Step: "build", Leaf: "LA", Handle: "a", Role: "verifier", Attempt: 1, Session: av},
		{Type: TypeLeafSession, Step: "build", Leaf: "LB", Handle: "b", Role: "builder", Attempt: 1, Session: bb},
	} {
		en.At = at(1)
		if err := e.Store.Append(runID, en); err != nil {
			t.Fatal(err)
		}
	}
	r.agentOut[av] = []bool{false}
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
	if want := []string{"b verifier", "c builder", "c verifier"}; !reflect.DeepEqual(started, want) {
		t.Errorf("started %v, want %v", started, want)
	}
	if r.polls[av] != 2 || r.polls[bb] != 1 {
		t.Errorf("polls: A verifier %d, B builder %d; want both rejoined (2, 1)", r.polls[av], r.polls[bb])
	}
	for _, l := range st.Leaves {
		if l.Outcome != LeafPassed || len(l.Sessions) != 2 {
			t.Errorf("leaf %s = %+v", l.Handle, l)
		}
	}
}

// Two leaf sessions blocked together each get one blocked and one
// unblocked notice, however often they are polled.
func TestLeafOrderBlocksPerSession(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = abcLeaves[:2]
	passAll(r, abcLeaves)
	av, bb := leafSID("a", "verifier", 1), leafSID("b", "builder", 1)
	blk := AgentStatus{Blocked: &Block{Kind: "permission"}}
	r.agentSeq[av] = []AgentStatus{blk, blk, blk, blk}
	r.agentSeq[bb] = []AgentStatus{blk, blk, {}, blk, blk}
	var notices []string
	e.Notify = func(n Notice) { notices = append(notices, n.Session+" "+strconv.FormatBool(n.Blocked)) }
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	count := map[string]int{}
	for _, n := range notices {
		count[n]++
	}
	want := map[string]int{av + " true": 1, av + " false": 1, bb + " true": 2, bb + " false": 2}
	if !reflect.DeepEqual(count, want) {
		t.Errorf("notices = %v, want counts %v", notices, want)
	}
}

// A builder that exits 0 passes only with a commit naming its leaf: with
// none its attempt fails with LeafCommit's detail and no verifier starts;
// with one its attempt records the commit, and the verifier starts on it,
// its leaf-session journalling it and its env carrying RHIZA_COMMIT. The
// builder's leaf-session journals HEAD as it started, the since of
// LeafCommit.
func TestLeavesStepLeafCommit(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = []Leaf{{ID: "L1", Handle: "leaf-a"}, {ID: "L2", Handle: "leaf-b"}}
	passAll(r, twoLeaves)
	r.leafCommit = map[string][2]string{"leaf-a": {"", "no new commit since H0"}, "leaf-b": {"abc123", ""}}
	for n := 2; n <= MaxAttempts; n++ {
		r.outcome[leafSID("leaf-a", "builder", n)] = workflow.OutcomePass
	}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:fail"}) {
		t.Fatalf("outcomes = %v", got)
	}
	a, b := st.Leaves[0], st.Leaves[1]
	if j, _ := a.Judged(workflow.RoleBuilder, 1); j.Outcome != workflow.OutcomeFail || j.Detail != "no new commit since H0" || j.Commit != "" {
		t.Errorf("leaf-a builder attempt = %+v", j)
	}
	if _, ok := a.Session(workflow.RoleVerifier, 1); ok || a.Outcome != LeafFailed {
		t.Errorf("leaf-a = %+v, want failed with no verifier", a)
	}
	if !strings.Contains(a.Detail, "no new commit since H0") {
		t.Errorf("leaf-a detail = %q", a.Detail)
	}
	if s, _ := a.Session(workflow.RoleBuilder, 1); s.Commit != fakeHead {
		t.Errorf("leaf-a builder session commit = %q, want %q", s.Commit, fakeHead)
	}
	if j, _ := b.Judged(workflow.RoleBuilder, 1); j.Outcome != workflow.OutcomePass || j.Commit != "abc123" {
		t.Errorf("leaf-b builder attempt = %+v", j)
	}
	if s, _ := b.Session(workflow.RoleVerifier, 1); s.Commit != "abc123" {
		t.Errorf("leaf-b verifier session commit = %q", s.Commit)
	}
	if b.Outcome != LeafPassed {
		t.Errorf("leaf-b = %+v", b)
	}
	var verifier LeafSpec
	for _, s := range r.leafSpecs {
		if s.Leaf.Handle == "leaf-b" && s.Role == workflow.RoleVerifier {
			verifier = s
		}
		if s.Role == workflow.RoleBuilder && (s.Commit != "" || !slices.Contains(s.Env, "RHIZA_ROLE=builder")) {
			t.Errorf("builder spec = %+v", s)
		}
	}
	if verifier.Commit != "abc123" || !slices.Contains(verifier.Env, "RHIZA_COMMIT=abc123") || !slices.Contains(verifier.Env, "RHIZA_ROLE=verifier") {
		t.Errorf("verifier spec = %+v", verifier)
	}
	if want := []string{"leaf-a H0", "leaf-a H0", "leaf-a H0", "leaf-b H0"}; !reflect.DeepEqual(r.leafSince, want) {
		t.Errorf("LeafCommit calls = %v, want %v", r.leafSince, want)
	}
}

// A builder that exits non-zero fails on its exit alone: LeafCommit is not
// asked.
func TestLeavesStepBuilderExitFailsWithoutCommitCheck(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = []Leaf{{ID: "L1", Handle: "leaf-a"}}
	r.outcome[leafSID("leaf-a", workflow.RoleBuilder, 1)] = workflow.OutcomeFail
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	if len(r.leafSince) != 0 {
		t.Errorf("LeafCommit calls = %v, want none", r.leafSince)
	}
}

// verifierAttempt is the judgement of leaf-a's verifier on attempt 1.
func verifierAttempt(t *testing.T, st State) LeafAttempt {
	t.Helper()
	j, ok := st.Leaves[0].Judged(workflow.RoleVerifier, 1)
	if !ok {
		t.Fatalf("verifier not judged: %+v", st.Leaves[0])
	}
	return j
}

// A verifier that ends well passes the attempt only with Evidence its own
// actor wrote; ending with neither Evidence nor findings, or with Evidence
// another actor (the builder's) wrote, fails it.
func TestLeavesStepVerifierEvidence(t *testing.T) {
	for _, tc := range []struct {
		name       string
		evidence   *fakeEvidence
		outcome    string
		detail     string
		leafResult string
	}{
		{name: "own Evidence", outcome: workflow.OutcomePass, detail: "Evidence written by the verifier", leafResult: LeafPassed},
		{name: "neither", evidence: &fakeEvidence{}, outcome: workflow.OutcomeFail,
			detail: "verifier ended with neither Evidence nor findings", leafResult: LeafFailed},
		{name: "builder's Evidence", evidence: &fakeEvidence{detail: "Evidence ev-a was written by claude-code-builder, not the verifier's actor claude-code-verifier"},
			outcome:    workflow.OutcomeFail,
			detail:     "verifier ended with no Evidence of its own and no findings: Evidence ev-a was written by claude-code-builder, not the verifier's actor claude-code-verifier",
			leafResult: LeafFailed},
		{name: "error", evidence: &fakeEvidence{err: errors.New("boom")}, outcome: workflow.OutcomeFail,
			detail: "read its Evidence: boom", leafResult: LeafFailed},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e, r, _, _ := setup(t, leavesWF)
			r.leaves = twoLeaves[:1]
			passAll(r, twoLeaves)
			if tc.evidence != nil {
				r.evidence = map[string]fakeEvidence{"leaf-a": *tc.evidence}
			}
			st, err := e.Run(runID)
			if err != nil {
				t.Fatal(err)
			}
			if j := verifierAttempt(t, st); j.Outcome != tc.outcome || j.Detail != tc.detail {
				t.Errorf("verifier judged %s (%s), want %s (%s)", j.Outcome, j.Detail, tc.outcome, tc.detail)
			}
			if o := st.Leaves[0].Outcome; o != tc.leafResult {
				t.Errorf("leaf outcome = %s, want %s", o, tc.leafResult)
			}
			if want := []string{"leaf-a " + leafSID("leaf-a", "verifier", 1)}; !reflect.DeepEqual(r.evidenceCalls, want) {
				t.Errorf("Evidence calls = %v, want %v", r.evidenceCalls, want)
			}
		})
	}
}

// A verifier judged pass has its Evidence committed once, the commit
// journalled on its attempt and on the passed leaf; a commit that fails
// fails the attempt, and so the leaf, with that detail. A verifier with
// no Evidence of its own commits nothing.
func TestLeavesStepEvidenceCommit(t *testing.T) {
	t.Run("pass", func(t *testing.T) {
		e, r, _, _ := setup(t, leavesWF)
		r.leaves = twoLeaves[:1]
		passAll(r, twoLeaves)
		st, err := e.Run(runID)
		if err != nil {
			t.Fatal(err)
		}
		if want := []string{"leaf-a " + leafSID("leaf-a", "verifier", 1)}; !reflect.DeepEqual(r.commitCalls, want) {
			t.Errorf("CommitEvidence calls = %v, want %v", r.commitCalls, want)
		}
		if j := verifierAttempt(t, st); j.Outcome != workflow.OutcomePass || j.Commit != "E-leaf-a" {
			t.Errorf("verifier judged %+v, want pass with E-leaf-a", j)
		}
		if l := st.Leaves[0]; l.Outcome != LeafPassed || l.Commit != "E-leaf-a" {
			t.Errorf("leaf = %+v, want passed with E-leaf-a", l)
		}
	})
	t.Run("commit fails", func(t *testing.T) {
		e, r, _, _ := setup(t, leavesWF)
		r.leaves = twoLeaves[:1]
		passAll(r, twoLeaves)
		r.evidenceCommitErr = map[string]error{"leaf-a": errors.New("index.lock exists")}
		st, err := e.Run(runID)
		if err != nil {
			t.Fatal(err)
		}
		if j := verifierAttempt(t, st); j.Outcome != workflow.OutcomeFail || j.Detail != "commit its Evidence: index.lock exists" || j.Commit != "" {
			t.Errorf("verifier judged %+v", j)
		}
		if l := st.Leaves[0]; l.Outcome != LeafFailed || l.Commit != "" {
			t.Errorf("leaf = %+v, want failed with no commit", l)
		}
	})
	t.Run("no Evidence", func(t *testing.T) {
		e, r, _, _ := setup(t, leavesWF)
		r.leaves = twoLeaves[:1]
		passAll(r, twoLeaves)
		r.evidence = map[string]fakeEvidence{"leaf-a": {}}
		if _, err := e.Run(runID); err != nil {
			t.Fatal(err)
		}
		if len(r.commitCalls) != 0 {
			t.Errorf("CommitEvidence called without Evidence: %v", r.commitCalls)
		}
	})
}

// Findings a verifier's process journals on its attempt, outside the
// engine, fail the attempt even when its own Evidence is there too; the
// engine reads the journal again before judging.
func TestLeavesStepVerifierFindings(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = twoLeaves[:1]
	passAll(r, twoLeaves)
	v := leafSID("leaf-a", "verifier", 1)
	r.agentSeq[v] = []AgentStatus{{}}
	r.onPoll = func(session string, n int) {
		if session != v || n != 0 {
			return
		}
		if err := e.Store.Append(runID, Entry{
			Type: TypeFindings, At: at(2), Step: "build", Leaf: "L1", Handle: "leaf-a", Attempt: 1,
			Findings: []Finding{{ID: "f1", Reason: "no test"}, {ID: "f2", Reason: "wrong exit code"}},
		}); err != nil {
			t.Fatal(err)
		}
	}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if j := verifierAttempt(t, st); j.Outcome != workflow.OutcomeFail || j.Detail != "verifier found 2 problem(s)" {
		t.Errorf("verifier judged %s (%s)", j.Outcome, j.Detail)
	}
	l := st.Leaves[0]
	if l.Outcome != LeafFailed || len(l.AttemptFindings(1)) != 2 || len(l.AttemptFindings(2)) != 0 {
		t.Errorf("leaf = %+v", l)
	}
	if len(st.Rounds) != 0 {
		t.Errorf("a leaf's findings opened a plan round: %+v", st.Rounds)
	}
	if len(r.evidenceCalls) != 0 {
		t.Errorf("Evidence called with findings journalled: %v", r.evidenceCalls)
	}
}

// A verifier whose process exits non-zero fails its attempt without its
// Evidence being looked at.
func TestLeavesStepVerifierExitFails(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = twoLeaves[:1]
	passAll(r, twoLeaves)
	r.outcome[leafSID("leaf-a", "verifier", 1)] = workflow.OutcomeFail
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if j := verifierAttempt(t, st); j.Outcome != workflow.OutcomeFail || len(r.evidenceCalls) != 0 {
		t.Errorf("verifier judged %+v; Evidence calls %v", j, r.evidenceCalls)
	}
}

// A verifier works in a checkout of its builder's commit, opened before
// it starts, named in its LeafSpec and leaf-session Detail, and removed
// once it is judged, whether it passed or failed. A removal that fails is
// only logged.
func TestLeavesStepVerifierCheckout(t *testing.T) {
	for _, tc := range []struct {
		name, outcome, leaf string
		removeErr           error
	}{
		{"pass", workflow.OutcomePass, LeafPassed, nil},
		{"fail", workflow.OutcomeFail, LeafFailed, nil},
		{"remove fails", workflow.OutcomePass, LeafPassed, errors.New("locked")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e, r, _, _ := setup(t, leavesWF)
			r.leaves = twoLeaves[:1]
			passAll(r, twoLeaves)
			r.outcome[leafSID("leaf-a", "verifier", 1)] = tc.outcome
			r.removeErr = tc.removeErr
			st, err := e.Run(runID)
			if err != nil {
				t.Fatal(err)
			}
			if st.Leaves[0].Outcome != tc.leaf {
				t.Fatalf("leaf = %+v", st.Leaves[0])
			}
			dir := checkoutDir("leaf-a", 1)
			if want := []string{"open " + dir + " C-leaf-a", "remove " + dir}; !reflect.DeepEqual(r.checkouts, want) {
				t.Errorf("checkouts = %v, want %v", r.checkouts, want)
			}
			if len(r.leafSpecs) < 2 || r.leafSpecs[0].Dir != "" || r.leafSpecs[1].Dir != dir {
				t.Errorf("specs = %+v", r.leafSpecs)
			}
			if s, _ := st.Leaves[0].Session(workflow.RoleVerifier, 1); s.Detail != dir {
				t.Errorf("verifier session = %+v", s)
			}
		})
	}
}

// A checkout that cannot be opened ends the step with an error, and no
// verifier starts.
func TestLeavesStepVerifierCheckoutFails(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = twoLeaves[:1]
	passAll(r, twoLeaves)
	r.checkoutErr = errors.New("no such commit")
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:error"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if len(r.leafSpecs) != 1 {
		t.Errorf("specs = %+v", r.leafSpecs)
	}
}

// A verifier whose launch fails has its checkout removed.
func TestLeavesStepVerifierStartFailsRemovesCheckout(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = twoLeaves[:1]
	passAll(r, twoLeaves)
	r.leafErr = map[string]error{"leaf-a verifier": errors.New("no rex")}
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	dir := checkoutDir("leaf-a", 1)
	if want := []string{"open " + dir + " C-leaf-a", "remove " + dir}; !reflect.DeepEqual(r.checkouts, want) {
		t.Errorf("checkouts = %v, want %v", r.checkouts, want)
	}
}

// On a resume, the checkout of a verifier already judged is removed (a
// killed run never reached it), while a live verifier's stays until it is
// judged.
func TestLeavesStepResumeRemovesEndedCheckouts(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = []Leaf{{ID: "L1", Handle: "leaf-a"}, {ID: "L2", Handle: "leaf-b"}}
	passAll(r, r.leaves)
	for _, en := range []Entry{
		{Type: TypeEntered, Step: "build", Visit: 1},
		{Type: TypeLeafSession, Step: "build", Leaf: "L1", Handle: "leaf-a", Role: "builder", Attempt: 1, Session: leafSID("leaf-a", "builder", 1)},
		{Type: TypeLeafAttempt, Step: "build", Leaf: "L1", Handle: "leaf-a", Role: "builder", Attempt: 1, Outcome: "pass", Commit: "C-leaf-a"},
		{Type: TypeLeafSession, Step: "build", Leaf: "L1", Handle: "leaf-a", Role: "verifier", Attempt: 1, Session: leafSID("leaf-a", "verifier", 1), Commit: "C-leaf-a", Detail: "/old/a"},
		{Type: TypeLeafAttempt, Step: "build", Leaf: "L1", Handle: "leaf-a", Role: "verifier", Attempt: 1, Outcome: "pass"},
		{Type: TypeLeafSession, Step: "build", Leaf: "L2", Handle: "leaf-b", Role: "builder", Attempt: 1, Session: leafSID("leaf-b", "builder", 1)},
		{Type: TypeLeafAttempt, Step: "build", Leaf: "L2", Handle: "leaf-b", Role: "builder", Attempt: 1, Outcome: "pass", Commit: "C-leaf-b"},
		{Type: TypeLeafSession, Step: "build", Leaf: "L2", Handle: "leaf-b", Role: "verifier", Attempt: 1, Session: leafSID("leaf-b", "verifier", 1), Commit: "C-leaf-b", Detail: "/old/b"},
	} {
		en.At = at(1)
		if err := e.Store.Append(runID, en); err != nil {
			t.Fatal(err)
		}
	}
	r.agentSeq[leafSID("leaf-b", "verifier", 1)] = []AgentStatus{{Idle: true}}
	r.onPoll = func(session string, n int) {
		if session == leafSID("leaf-b", "verifier", 1) && n == 1 && len(r.checkouts) != 1 {
			t.Errorf("checkouts while leaf-b's verifier is live = %v", r.checkouts)
		}
	}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if want := []string{"remove /old/a", "remove /old/b"}; !reflect.DeepEqual(r.checkouts, want) {
		t.Errorf("checkouts = %v, want %v", r.checkouts, want)
	}
}

// A leaf whose verifier's pass was journalled with its Evidence commit, but
// whose leaf entry was not (a run killed in between), ends passed with
// that commit, without asking for another.
func TestNextForLeafEvidenceCommitResume(t *testing.T) {
	l := LeafState{
		Sessions: []LeafSession{
			{Role: workflow.RoleBuilder, Attempt: 1, Session: "b"},
			{Role: workflow.RoleVerifier, Attempt: 1, Session: "v"},
		},
		Attempts: []LeafAttempt{
			{Role: workflow.RoleBuilder, Attempt: 1, Outcome: workflow.OutcomePass, Commit: "C1"},
			{Role: workflow.RoleVerifier, Attempt: 1, Outcome: workflow.OutcomePass, Commit: "E1"},
		},
	}
	a := nextForLeaf(Leaf{ID: "L1", Handle: "leaf-a"}, l)
	if a.kind != leafEnd || a.outcome != LeafPassed || a.commit != "E1" {
		t.Errorf("nextForLeaf = %+v, want leafEnd passed with E1", a)
	}
}

// retriesWF is leavesWF with a model ladder on the builder.
var retriesWF = strings.Replace(leavesWF, "prompt = \"build {leaf}\"\n", "prompt = \"build {leaf}\"\nmodels = [\"sonnet\", \"opus\", \"opus\"]\n", 1)

// builderSpecs are r's builder LeafSpecs as "<handle> <attempt> <model>
// <RHIZA_ATTEMPT>", in start order.
func builderSpecs(r *fakeRunner) []string {
	var out []string
	for _, s := range r.leafSpecs {
		if s.Role != workflow.RoleBuilder {
			continue
		}
		env := ""
		for _, kv := range s.Env {
			if v, ok := strings.CutPrefix(kv, EnvAttempt+"="); ok {
				env = v
			}
		}
		out = append(out, s.Leaf.Handle+" "+strconv.Itoa(s.Attempt)+" "+s.Model+" "+env)
	}
	return out
}

// A verifier that fails twice, once with findings and once on its exit,
// sends each failed attempt to a fresh builder session on the next model
// of the ladder, its spec and RHIZA_ATTEMPT naming the attempt and its
// commit check measured from its own start HEAD; the third attempt's
// verifier passes the leaf. The earlier attempts' findings stay in the
// run for the retry builder to read.
func TestLeavesStepRetries(t *testing.T) {
	e, r, _, _ := setup(t, retriesWF)
	r.leaves = twoLeaves[:1]
	r.heads = []string{"H1", "H2", "H3"}
	for n := 1; n <= 3; n++ {
		r.outcome[leafSID("leaf-a", "builder", n)] = workflow.OutcomePass
	}
	v1 := leafSID("leaf-a", "verifier", 1)
	r.outcome[v1] = workflow.OutcomePass
	r.agentSeq[v1] = []AgentStatus{{}}
	r.onPoll = func(session string, n int) {
		if session != v1 || n != 0 {
			return
		}
		if err := e.Store.Append(runID, Entry{
			Type: TypeFindings, At: at(2), Step: "build", Leaf: "L1", Handle: "leaf-a", Attempt: 1,
			Findings: []Finding{{ID: "f1", Reason: "no test"}},
		}); err != nil {
			t.Fatal(err)
		}
	}
	r.outcome[leafSID("leaf-a", "verifier", 2)] = workflow.OutcomeFail
	r.outcome[leafSID("leaf-a", "verifier", 3)] = workflow.OutcomePass
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if want := []string{"leaf-a 1 sonnet 1", "leaf-a 2 opus 2", "leaf-a 3 opus 3"}; !reflect.DeepEqual(builderSpecs(r), want) {
		t.Errorf("builders = %q, want %q", builderSpecs(r), want)
	}
	l := st.Leaves[0]
	if l.Outcome != LeafPassed || l.Detail != "attempt 3: verifier passed" || l.Commit != "E-leaf-a" {
		t.Errorf("leaf = %+v", l)
	}
	seen := map[string]bool{}
	for n := 1; n <= 3; n++ {
		s, ok := l.Session(workflow.RoleBuilder, n)
		if !ok || seen[s.Session] || s.Commit != "H"+strconv.Itoa(n) {
			t.Errorf("builder session %d = %+v (ok %v)", n, s, ok)
		}
		seen[s.Session] = true
	}
	if want := []string{"leaf-a H1", "leaf-a H2", "leaf-a H3"}; !reflect.DeepEqual(r.leafSince, want) {
		t.Errorf("LeafCommit calls = %v, want %v", r.leafSince, want)
	}
	if f := l.AttemptFindings(1); len(f) != 1 || f[0].Reason != "no test" {
		t.Errorf("attempt 1 findings = %+v", f)
	}
}

// A verifier that fails all MaxAttempts attempts fails the leaf after the
// third builder: no fourth builder starts.
func TestLeavesStepRetriesCap(t *testing.T) {
	e, r, _, _ := setup(t, retriesWF)
	r.leaves = twoLeaves[:1]
	for n := 1; n <= MaxAttempts+1; n++ {
		r.outcome[leafSID("leaf-a", "builder", n)] = workflow.OutcomePass
		r.outcome[leafSID("leaf-a", "verifier", n)] = workflow.OutcomeFail
	}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:fail"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if want := []string{"leaf-a 1 sonnet 1", "leaf-a 2 opus 2", "leaf-a 3 opus 3"}; !reflect.DeepEqual(builderSpecs(r), want) {
		t.Errorf("builders = %q, want %q", builderSpecs(r), want)
	}
	l := st.Leaves[0]
	if l.Outcome != LeafFailed || l.Detail != "attempt 3: verifier failed: exit 0" {
		t.Errorf("leaf = %+v", l)
	}
	if _, ok := l.Session(workflow.RoleBuilder, 4); ok || l.Attempt() != 3 {
		t.Errorf("a fourth attempt started: %+v", l.Sessions)
	}
}

// A builder that fails, on its exit or with no leaf commit, uses up an
// attempt as a verifier's fail does: the next attempt's builder starts,
// with no verifier on the failed one.
func TestLeavesStepRetriesBuilderFailCounts(t *testing.T) {
	e, r, _, _ := setup(t, retriesWF)
	r.leaves = twoLeaves[:1]
	r.outcome[leafSID("leaf-a", "builder", 1)] = workflow.OutcomeFail
	r.outcome[leafSID("leaf-a", "builder", 2)] = workflow.OutcomePass
	r.outcome[leafSID("leaf-a", "verifier", 2)] = workflow.OutcomePass
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"leaf-session leaf-a builder 1", "leaf-attempt leaf-a builder 1 fail",
		"leaf-session leaf-a builder 2", "leaf-attempt leaf-a builder 2 pass",
		"leaf-session leaf-a verifier 2", "leaf-attempt leaf-a verifier 2 pass",
		"leaf leaf-a passed",
	}
	if got := leafEntries(t, e); !reflect.DeepEqual(got, want) {
		t.Errorf("leaf entries =\n%v\nwant\n%v", got, want)
	}
	if l := st.Leaves[0]; l.Outcome != LeafPassed || l.Detail != "attempt 2: verifier passed" {
		t.Errorf("leaf = %+v", l)
	}
}

// A retrying leaf's next builder goes before a new leaf's first one, and
// waits while another builder is live.
func TestLeavesRetryOrder(t *testing.T) {
	leaves := []Leaf{{ID: "LA", Handle: "a"}, {ID: "LB", Handle: "b"}}
	failedA := LeafState{Leaf: "LA",
		Sessions: []LeafSession{{Role: "builder", Attempt: 1, Session: "SAb"}, {Role: "verifier", Attempt: 1, Session: "SAv"}},
		Attempts: []LeafAttempt{{Role: "builder", Attempt: 1, Outcome: "pass"}, {Role: "verifier", Attempt: 1, Outcome: "fail"}}}
	got := nextLeafActions(leaves, map[string]LeafState{"LA": failedA})
	if len(got) != 1 || got[0].kind != leafStart || got[0].leaf.ID != "LA" || got[0].role != workflow.RoleBuilder || got[0].attempt != 2 {
		t.Fatalf("retry first: %+v", got)
	}
	liveB := LeafState{Leaf: "LB", Sessions: []LeafSession{{Role: "builder", Attempt: 1, Session: "SBb"}}}
	got = nextLeafActions(leaves, map[string]LeafState{"LA": failedA, "LB": liveB})
	if len(got) != 1 || got[0].kind != leafPoll || got[0].session != "SBb" {
		t.Fatalf("retry waits on a live builder: %+v", got)
	}
}

// A journal cut after attempt 1's verifier failed resumes with attempt 2's
// builder, starting nothing of attempt 1 again.
func TestLeavesStepRetriesResume(t *testing.T) {
	e, r, _, _ := setup(t, retriesWF)
	r.leaves = twoLeaves[:1]
	r.outcome[leafSID("leaf-a", "builder", 2)] = workflow.OutcomePass
	r.outcome[leafSID("leaf-a", "verifier", 2)] = workflow.OutcomePass
	for _, en := range []Entry{
		{Type: TypeEntered, Step: "build", Visit: 1},
		{Type: TypeLeafSession, Step: "build", Leaf: "L1", Handle: "leaf-a", Role: "builder", Attempt: 1, Session: leafSID("leaf-a", "builder", 1), Commit: "H0"},
		{Type: TypeLeafAttempt, Step: "build", Leaf: "L1", Handle: "leaf-a", Role: "builder", Attempt: 1, Outcome: "pass", Commit: "C-leaf-a"},
		{Type: TypeLeafSession, Step: "build", Leaf: "L1", Handle: "leaf-a", Role: "verifier", Attempt: 1, Session: leafSID("leaf-a", "verifier", 1), Commit: "C-leaf-a"},
		{Type: TypeLeafAttempt, Step: "build", Leaf: "L1", Handle: "leaf-a", Role: "verifier", Attempt: 1, Outcome: "fail", Detail: "verifier found 1 problem(s)"},
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
	var started []string
	for _, s := range r.leafSpecs {
		started = append(started, s.Leaf.Handle+" "+s.Role+" "+strconv.Itoa(s.Attempt)+" "+s.Model)
	}
	if want := []string{"leaf-a builder 2 opus", "leaf-a verifier 2 "}; !reflect.DeepEqual(started, want) {
		t.Errorf("started %q, want %q", started, want)
	}
	if l := st.Leaves[0]; l.Outcome != LeafPassed || l.Detail != "attempt 2: verifier passed" {
		t.Errorf("leaf = %+v", l)
	}
}

const leafFailsWF = `[step.build]
kind = "leaves"
interval = "5s"
[step.build.builder]
prompt = "build {leaf}"
models = ["sonnet", "opus", "opus"]
[step.build.verifier]
prompt = "verify {leaf}"
`

// abc are leaves A, B and C, C depending on A.
var abc = []Leaf{{ID: "LA", Handle: "leaf-a"}, {ID: "LB", Handle: "leaf-b"}, {ID: "LC", Handle: "leaf-c", DependsOn: []string{"LA"}}}

// leafFailsRun runs abc with leaf-a's builder failing every attempt and
// the other leaves passing on attempt 1; leaf-a passes from attempt from
// on, when from > 0.
func leafFailsRun(t *testing.T, e *Engine, r *fakeRunner, from int) State {
	t.Helper()
	r.leaves = abc
	passAll(r, abc)
	r.outcome[leafSID("leaf-a", "builder", 1)] = workflow.OutcomeFail
	if from > 0 {
		for _, role := range []string{workflow.RoleBuilder, workflow.RoleVerifier} {
			r.outcome[leafSID("leaf-a", role, from)] = workflow.OutcomePass
		}
	}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	return st
}

// leafOutcomes are each leaf's handle and outcome, in State.Leaves order.
func leafOutcomes(st State) []string {
	var out []string
	for _, l := range st.Leaves {
		out = append(out, l.Handle+":"+l.Outcome)
	}
	return out
}

// A leaf that fails all its attempts raises one question linked to it,
// journalled on its leaf entry; its dependent is skipped, the independent
// leaf is built and verified, and the run ends failed.
func TestLeafFails(t *testing.T) {
	e, r, _, _ := setup(t, leafFailsWF)
	st := leafFailsRun(t, e, r, 0)
	if st.Result != Failed || !reflect.DeepEqual(outcomes(st), []string{"build:fail"}) {
		t.Fatalf("result %q steps %v", st.Result, outcomes(st))
	}
	if got, want := leafOutcomes(st), []string{"leaf-a:failed", "leaf-b:passed", "leaf-c:skipped"}; !reflect.DeepEqual(got, want) {
		t.Errorf("leaves %v, want %v", got, want)
	}
	if len(r.leafQuestions) != 1 || !strings.HasPrefix(r.leafQuestions[0], "leaf-a Leaf leaf-a failed all 3 attempts of round 1: attempt 1 builder failed: exit 0; attempt 2 builder failed") {
		t.Errorf("questions %q", r.leafQuestions)
	}
	a, _ := st.Leaf("build", "LA")
	if a.Question != "QL1" || a.QuestionHandle != "leaf-failed-1" || a.Round() != 1 {
		t.Errorf("leaf-a = %+v", a)
	}
	b, _ := st.Leaf("build", "LB")
	if _, ok := b.Judged(workflow.RoleVerifier, 1); !ok || b.Commit != "E-leaf-b" {
		t.Errorf("leaf-b was not verified: %+v", b)
	}
	if c, _ := st.Leaf("build", "LC"); c.Detail != "its dependency leaf-a did not pass" || len(c.Sessions) != 0 {
		t.Errorf("leaf-c = %+v", c)
	}
	if len(r.arbiters) != 0 || len(r.leafAnswerN) != 0 {
		t.Errorf("a normal run arbitrated %v or read answers %v", r.arbiters, r.leafAnswerN)
	}
}

// A failed verifier's findings are in its leaf's question.
func TestLeafFailsQuestionNamesFindings(t *testing.T) {
	l := LeafState{Handle: "leaf-a"}
	for n := 1; n <= MaxAttempts; n++ {
		l.Sessions = append(l.Sessions, LeafSession{Role: workflow.RoleBuilder, Attempt: n})
		l.Attempts = append(l.Attempts,
			LeafAttempt{Role: workflow.RoleBuilder, Attempt: n, Outcome: workflow.OutcomePass},
			LeafAttempt{Role: workflow.RoleVerifier, Attempt: n, Outcome: workflow.OutcomeFail, Detail: "verifier found 1 problem(s)"})
		l.Findings = append(l.Findings, LeafFinding{Attempt: n, Finding: Finding{ID: "f1", Reason: "No test.\nAt all."}})
	}
	got := leafQuestionText(l)
	want := "Leaf leaf-a failed all 3 attempts of round 1: attempt 1 verifier failed: verifier found 1 problem(s) (f1 No test. At all); " +
		"attempt 2 verifier failed: verifier found 1 problem(s) (f1 No test. At all); attempt 3 verifier failed: verifier found 1 problem(s) (f1 No test. At all). " +
		"How should it be built, or should it change?"
	if got != want {
		t.Errorf("text\n%q\nwant\n%q", got, want)
	}
}

// A question that cannot be added leaves the leaf failed without one,
// named in its detail, and the run carries on to fail.
func TestLeafFailsQuestionNotAdded(t *testing.T) {
	e, r, _, _ := setup(t, leafFailsWF)
	r.leafQuestionErr = errors.New("change set landed")
	st := leafFailsRun(t, e, r, 0)
	a, _ := st.Leaf("build", "LA")
	if st.Result != Failed || a.Question != "" || !strings.HasSuffix(a.Detail, "; its question could not be added: change set landed") {
		t.Errorf("result %q leaf-a %+v", st.Result, a)
	}
	if got, want := leafOutcomes(st), []string{"leaf-a:failed", "leaf-b:passed", "leaf-c:skipped"}; !reflect.DeepEqual(got, want) {
		t.Errorf("leaves %v, want %v", got, want)
	}
}

// After a Decision answers the failed leaf's question, a retry keeps the
// passed leaf, gives the failed one a fresh round of 3 attempts on the
// ladder's first model with the Decision in its builder's spec, then
// builds its dependent, and the run succeeds.
func TestLeafFailsResume(t *testing.T) {
	e, r, _, _ := setup(t, leafFailsWF)
	if st := leafFailsRun(t, e, r, 0); st.Result != Failed {
		t.Fatalf("first run %q", st.Result)
	}
	before := len(r.leafSpecs)
	for _, role := range []string{workflow.RoleBuilder, workflow.RoleVerifier} {
		r.outcome[leafSID("leaf-a", role, 4)] = workflow.OutcomePass
	}
	r.leafAnswers["QL1"] = [2]string{"use-a-map", "Use a map, not a list."}
	st, err := e.Retry(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded {
		t.Fatalf("retry %q: %s", st.Result, st.Detail)
	}
	if got, want := leafOutcomes(st), []string{"leaf-a:passed", "leaf-b:passed", "leaf-c:passed"}; !reflect.DeepEqual(got, want) {
		t.Errorf("leaves %v, want %v", got, want)
	}
	var started []string
	for _, s := range r.leafSpecs[before:] {
		started = append(started, s.Leaf.Handle+" "+s.Role+" "+strconv.Itoa(s.Attempt))
	}
	if want := []string{"leaf-a builder 4", "leaf-a verifier 4", "leaf-c builder 1", "leaf-c verifier 1"}; !reflect.DeepEqual(started, want) {
		t.Errorf("started %v, want %v", started, want)
	}
	spec := r.leafSpecs[before]
	if spec.Round != 2 || spec.RoundAttempt != 1 || spec.Model != "sonnet" || spec.Question != "QL1" ||
		spec.Decision != "use-a-map" || spec.Answer != "Use a map, not a list." || !slices.Contains(spec.Env, "RHIZA_ATTEMPT=4") {
		t.Errorf("fresh round spec = %+v", spec)
	}
	if r.leafSpecs[2].Model != "opus" || r.leafSpecs[2].Attempt != 3 {
		t.Errorf("round 1 attempt 3 spec = %+v", r.leafSpecs[2])
	}
	a, _ := st.Leaf("build", "LA")
	if want := []LeafRound{{Round: 2, After: 3, Question: "QL1", QuestionHandle: "leaf-failed-1", Decision: "use-a-map", At: at(1)}}; !reflect.DeepEqual(a.Rounds, want) {
		t.Errorf("rounds %+v, want %+v", a.Rounds, want)
	}
	types := entryTypes(t, e)
	if !slices.Contains(types, "leaf-round:build") {
		t.Errorf("journal %v lacks leaf-round", types)
	}
}

// A retry while the failed leaf's question is unanswered starts nothing:
// the leaf stays failed, its dependent skipped, and the run fails again.
func TestLeafFailsResumeUnanswered(t *testing.T) {
	e, r, _, _ := setup(t, leafFailsWF)
	leafFailsRun(t, e, r, 0)
	before := len(r.leafSpecs)
	st, err := e.Retry(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed || len(r.leafSpecs) != before || r.leafAnswerN["QL1"] != 1 {
		t.Errorf("result %q, %d new sessions, %d answer reads", st.Result, len(r.leafSpecs)-before, r.leafAnswerN["QL1"])
	}
	if got, want := leafOutcomes(st), []string{"leaf-a:failed", "leaf-b:passed", "leaf-c:skipped"}; !reflect.DeepEqual(got, want) {
		t.Errorf("leaves %v, want %v", got, want)
	}
}

// In an autonomous run the arbiter is started on the failed leaf's
// question at once; its pass gives the leaf one fresh round, whose
// failure raises a second question that is never arbitrated, and the run
// ends failed.
func TestLeafFailsAutonomous(t *testing.T) {
	e, r := autonomousSetup(t, leafFailsWF)
	r.answers["QL1"] = Answer{OK: true, Decision: "assume-a-map", Statement: "Use a map."}
	r.leafAnswers["QL1"] = [2]string{"assume-a-map", "Use a map."}
	st := leafFailsRun(t, e, r, 0)
	if st.Result != Failed {
		t.Fatalf("result %q", st.Result)
	}
	if !reflect.DeepEqual(r.arbiters, []string{"arbiter QL1"}) {
		t.Errorf("arbiters %v, want one on QL1", r.arbiters)
	}
	if len(r.leafQuestions) != 2 || !strings.HasPrefix(r.leafQuestions[1], "leaf-a Leaf leaf-a failed all 3 attempts of round 2: attempt 4 builder failed") {
		t.Errorf("questions %q", r.leafQuestions)
	}
	a, _ := st.Leaf("build", "LA")
	if a.Outcome != LeafFailed || a.Question != "QL2" || len(a.Rounds) != 1 || a.Attempt() != 6 {
		t.Errorf("leaf-a = %+v", a)
	}
	var round2 []string
	for _, s := range r.leafSpecs {
		if s.Leaf.Handle == "leaf-a" && s.Round == 2 {
			round2 = append(round2, s.Model+" "+s.Decision+" "+s.Answer)
		}
	}
	if want := []string{"sonnet assume-a-map Use a map.", "opus assume-a-map Use a map.", "opus assume-a-map Use a map."}; !reflect.DeepEqual(round2, want) {
		t.Errorf("round 2 builders %q, want %q", round2, want)
	}
	if got, want := leafOutcomes(st), []string{"leaf-a:failed", "leaf-b:passed", "leaf-c:skipped"}; !reflect.DeepEqual(got, want) {
		t.Errorf("leaves %v, want %v", got, want)
	}
	if len(r.leafAnswerN) != 1 {
		t.Errorf("answer reads %v: the second question was read", r.leafAnswerN)
	}
}

// In an autonomous run an arbiter that fails leaves the leaf failed with
// no fresh round, and the run ends failed.
func TestLeafFailsAutonomousArbiterFails(t *testing.T) {
	e, r := autonomousSetup(t, leafFailsWF)
	st := leafFailsRun(t, e, r, 0)
	a, _ := st.Leaf("build", "LA")
	if st.Result != Failed || a.Outcome != LeafFailed || len(a.Rounds) != 0 || a.Attempt() != 3 {
		t.Errorf("result %q leaf-a %+v", st.Result, a)
	}
	if arb, ok := st.Arbitration("QL1"); !ok || arb.Outcome != workflow.OutcomeFail || !reflect.DeepEqual(r.arbiters, []string{"arbiter QL1"}) {
		t.Errorf("arbitration %+v arbiters %v", arb, r.arbiters)
	}
}
