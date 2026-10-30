package run

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/workflow"
)

func (f *fakeRunner) StartFix(spec FixSpec) (string, error) {
	f.fixSpecs = append(f.fixSpecs, spec)
	return fixAttemptSID(spec.Cycle, spec.Attempt), nil
}

func (f *fakeRunner) ReviewCommit(since string) (string, string, error) {
	f.reviewSince = append(f.reviewSince, since)
	if len(f.reviewCommits) > 0 {
		c := f.reviewCommits[0]
		f.reviewCommits = f.reviewCommits[1:]
		return c[0], c[1], nil
	}
	return "C-review", "", nil
}

func (f *fakeRunner) CommitCheckWrites(session string) (string, error) {
	f.checkWrites = append(f.checkWrites, session)
	if err := f.checkWritesErr[session]; err != nil {
		return "", err
	}
	return "W-" + session, nil
}

// fixSID is the fake runner's session ID for the change set's fixer in
// cycle's first attempt.
func fixSID(cycle int) string { return fixAttemptSID(cycle, 1) }

// fixAttemptSID is the session ID for the change set's fixer in cycle on
// attempt: a retry (attempt above 1) gets its own.
func fixAttemptSID(cycle, attempt int) string {
	id := "S-fixer-" + strconv.Itoa(cycle)
	if attempt > 1 {
		id += "-r" + strconv.Itoa(attempt)
	}
	return id
}

// fixFindingIDs are each finding's "check id severity".
func fixFindingIDs(fs []Finding) []string {
	var out []string
	for _, f := range fs {
		out = append(out, f.Check+" "+f.ID+" "+f.Severity)
	}
	return out
}

// entryAt is the index of s in entries, failing the test without it.
func entryAt(t *testing.T, entries []string, s string) int {
	t.Helper()
	i := slices.Index(entries, s)
	if i < 0 {
		t.Fatalf("no %q in %v", s, entries)
	}
	return i
}

// The acceptance of leaf review-fix-change-set: two confirmed findings
// naming no leaf, from two checks, go to one fixer session for the change
// set, in the run's worktree, told both findings; it passes on a new
// commit whose Review: trailer names the run (ReviewCommit since its
// start), and only after it is judged do the failed checks, and only
// they, run again.
func TestFixChangeSetOneFixer(t *testing.T) {
	e, r, _, _ := setup(t, checksWF)
	passChecks(r, 1, threeChecks...)
	passChecks(r, 2, "alpha", "gamma")
	r.findings[checkSID("alpha", "check", 1)] = []Finding{sev("f1", SeverityMajor), sev("f2", SeverityMinor)}
	r.findings[checkSID("gamma", "check", 1)] = []Finding{sev("f1", SeverityBlocker)}
	for _, c := range []string{"alpha", "gamma"} {
		r.outcome[checkSID(c, "critic", 1)] = workflow.OutcomePass
	}
	r.outcome[fixSID(2)] = workflow.OutcomePass
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if len(r.leafSpecs) != 0 {
		t.Errorf("leaf sessions started: %v", leafSpecRoles(r))
	}
	if len(r.fixSpecs) != 1 {
		t.Fatalf("change-set fixer sessions = %+v, want one", r.fixSpecs)
	}
	fx := r.fixSpecs[0]
	if fx.Cycle != 2 || fx.Step.Name != "review" {
		t.Errorf("fixer spec = %+v", fx)
	}
	if want := []string{"alpha f1 major", "gamma f1 blocker"}; !reflect.DeepEqual(fixFindingIDs(fx.Findings), want) {
		t.Errorf("fixer findings = %v, want %v", fixFindingIDs(fx.Findings), want)
	}
	for _, v := range []string{EnvCycle + "=2", EnvRole + "=fixer"} {
		if !slices.Contains(fx.Env, v) {
			t.Errorf("fixer env %v lacks %s", fx.Env, v)
		}
	}
	for _, v := range fx.Env {
		if strings.HasPrefix(v, EnvLeaf+"=") || strings.HasPrefix(v, EnvCheck+"=") || strings.HasPrefix(v, EnvAttempt+"=") {
			t.Errorf("fixer env has %s", v)
		}
	}
	if want := []string{fakeHead}; !reflect.DeepEqual(r.reviewSince, want) {
		t.Errorf("ReviewCommit calls = %v, want %v", r.reviewSince, want)
	}
	want := []string{"alpha check 1", "beta check 1", "gamma check 1", "alpha critic 1", "gamma critic 1", "alpha check 2", "gamma check 2"}
	if !reflect.DeepEqual(roles(r), want) {
		t.Errorf("check sessions %v, want %v", roles(r), want)
	}
	entries := fixEntries(t, e)
	if entryAt(t, entries, "check-cycle 2 alpha gamma") > entryAt(t, entries, "review-fix fixer 2") {
		t.Errorf("the fixer started before cycle 2 began: %v", entries)
	}
	if entryAt(t, entries, "review-fix-judged fixer 2 pass") > entryAt(t, entries, "check-session alpha check 2") {
		t.Errorf("a check of cycle 2 started before the fixer was judged: %v", entries)
	}
	fix, ok := st.ReviewFix("review", 2)
	if !ok || !fix.Judged || fix.Outcome != workflow.OutcomePass || fix.Commit != "C-review" || fix.Since != fakeHead || fix.Session != fixSID(2) {
		t.Errorf("fixer state = %+v, %v", fix, ok)
	}
	if !slices.Contains(r.closed, fixSID(2)) {
		t.Errorf("the fixer's pane was not closed at the step's end: %v", r.closed)
	}
	if !slices.Contains(r.ended, "review:"+fixSID(2)) {
		t.Errorf("AgentEnded not called for the fixer: %v", r.ended)
	}
}

// A change-set fixer that ends without a new commit whose Review: trailer
// names the run is retried within its cycle; one that fails every attempt
// (1 + MaxFixCycles) fails the cycle as a failed leaf fixer does: each
// check of the cycle fails with no session. The next cycle tries the fix
// again; with no cycle left the step fails.
func TestFixChangeSetNoCommitFailsCycle(t *testing.T) {
	e, r, _, _ := setup(t, checksWF)
	passChecks(r, 1, threeChecks...)
	r.findings[checkSID("alpha", "check", 1)] = []Finding{sev("f1", SeverityMajor)}
	r.outcome[checkSID("alpha", "critic", 1)] = workflow.OutcomePass
	for n := 1; n <= 1+MaxFixCycles; n++ {
		r.outcome[fixAttemptSID(2, n)] = workflow.OutcomePass
		r.outcome[fixAttemptSID(3, n)] = workflow.OutcomePass
		r.reviewCommits = append(r.reviewCommits, [2]string{"", "no new commit since H0"}, [2]string{"", "no new commit since H0"})
	}
	r.reviewCommits = r.reviewCommits[:2*(1+MaxFixCycles)]
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:fail"}) {
		t.Fatalf("outcomes = %v", got)
	}
	var got [][2]int
	for _, s := range r.fixSpecs {
		got = append(got, [2]int{s.Cycle, s.Attempt})
	}
	if want := [][2]int{{2, 1}, {2, 2}, {2, 3}, {3, 1}, {3, 2}, {3, 3}}; !reflect.DeepEqual(got, want) {
		t.Errorf("fixer sessions (cycle, attempt) = %v, want %v", got, want)
	}
	if want := []string{"alpha check 1", "beta check 1", "gamma check 1", "alpha critic 1"}; !reflect.DeepEqual(roles(r), want) {
		t.Errorf("check sessions %v, want %v", roles(r), want)
	}
	c, _ := st.Check("review", "alpha")
	o, _ := c.Outcome(2)
	if want := "cycle 2: fix failed: change set: fixer failed: no new commit since H0"; o.Outcome != CheckFailed || o.Detail != want {
		t.Errorf("alpha's cycle 2 outcome = %+v, want detail %q", o, want)
	}
	if fix, _ := st.ReviewFix("review", 2); fix.Attempt != 1+MaxFixCycles || fix.Outcome != workflow.OutcomeFail || fix.Detail != "no new commit since H0" {
		t.Errorf("fixer judgement = %+v", fix)
	}
	if s := st.Steps[0]; s.Detail != "2 of 3 checks passed after 3 cycles; failed: alpha; unanswered questions: review-failed-1" {
		t.Errorf("step detail = %q", s.Detail)
	}
}

// A change-set fixer that fails its first attempt and passes its second
// is retried inside the same cycle: no new cycle begins, and the cycle's
// checks start only after the retry passes.
func TestFixChangeSetRetriedInCycle(t *testing.T) {
	e, r, _, _ := setup(t, checksWF)
	passChecks(r, 1, threeChecks...)
	passChecks(r, 2, "alpha")
	r.findings[checkSID("alpha", "check", 1)] = []Finding{sev("f1", SeverityMajor)}
	r.outcome[checkSID("alpha", "critic", 1)] = workflow.OutcomePass
	r.outcome[fixAttemptSID(2, 1)] = workflow.OutcomeFail
	r.outcome[fixAttemptSID(2, 2)] = workflow.OutcomePass
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if len(st.Cycles) != 2 {
		t.Errorf("%d cycles, want 2: the retry runs in cycle 2", len(st.Cycles))
	}
	if len(r.fixSpecs) != 2 || r.fixSpecs[0].Attempt != 1 || r.fixSpecs[1].Attempt != 2 || r.fixSpecs[1].Cycle != 2 {
		t.Errorf("fixer sessions = %+v", r.fixSpecs)
	}
	entries := fixEntries(t, e)
	if entryAt(t, entries, "review-fix-judged fixer 2 pass") > entryAt(t, entries, "check-session alpha check 2") {
		t.Errorf("a check of cycle 2 started before the retry passed: %v", entries)
	}
	if fix, _ := st.ReviewFix("review", 2); fix.Attempt != 2 || fix.Outcome != workflow.OutcomePass {
		t.Errorf("latest fixer = %+v", fix)
	}
}

// A change-set fixer whose session fails fails the cycle too, without
// asking for its commit.
func TestFixChangeSetSessionFails(t *testing.T) {
	e, r, _, _ := setup(t, checksWF)
	passChecks(r, 1, threeChecks...)
	passChecks(r, 3, "alpha")
	r.findings[checkSID("alpha", "check", 1)] = []Finding{sev("f1", SeverityMajor)}
	r.outcome[checkSID("alpha", "critic", 1)] = workflow.OutcomePass
	r.outcome[fixSID(3)] = workflow.OutcomePass
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	c, _ := st.Check("review", "alpha")
	if o, _ := c.Outcome(2); o.Detail != "cycle 2: fix failed: change set: fixer failed: exit 0" {
		t.Errorf("alpha's cycle 2 outcome = %+v", o)
	}
	if want := []string{fakeHead}; !reflect.DeepEqual(r.reviewSince, want) {
		t.Errorf("ReviewCommit calls = %v, want only cycle 3's", r.reviewSince)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:pass"}) {
		t.Errorf("outcomes = %v", got)
	}
}

// The change set's fixer runs after every leaf's fixer, never beside
// one, though it may run beside the last leaf's verifier; a finding
// naming a leaf the change set does not have is the change set's too.
// The checks run again only once all are judged.
func TestFixChangeSetAfterLeafFixers(t *testing.T) {
	e, r, _, _ := setup(t, checksWF)
	r.leaves = fixLeaves
	a, b := fixLeaves[0], fixLeaves[1]
	passChecks(r, 1, threeChecks...)
	passChecks(r, 2, "alpha")
	stranger := Finding{ID: "f3", Severity: SeverityMajor, Leaf: "LX", LeafHandle: "gone", Reason: "bad f3"}
	r.findings[checkSID("alpha", "check", 1)] = []Finding{onLeaf(sev("f1", SeverityMajor), b), sev("f2", SeverityMajor), stranger, onLeaf(sev("f4", SeverityMajor), a)}
	r.outcome[checkSID("alpha", "critic", 1)] = workflow.OutcomePass
	passFixes(r, 2, a, b)
	r.outcome[fixSID(2)] = workflow.OutcomePass
	r.onPoll = func(session string, n int) {
		if session == fixSID(2) && n == 0 && len(r.leafSpecs) != 4 {
			t.Errorf("the change set's fixer polled with %v started", leafSpecRoles(r))
		}
	}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if want := []string{"leaf-a fixer 2", "leaf-a verifier 2", "leaf-b fixer 2", "leaf-b verifier 2"}; !reflect.DeepEqual(leafSpecRoles(r), want) {
		t.Fatalf("leaf sessions %v, want %v", leafSpecRoles(r), want)
	}
	if len(r.fixSpecs) != 1 {
		t.Fatalf("change-set fixers = %+v", r.fixSpecs)
	}
	if want := []string{"alpha f2 major", "alpha f3 major"}; !reflect.DeepEqual(fixFindingIDs(r.fixSpecs[0].Findings), want) {
		t.Errorf("change-set fixer findings = %v, want %v", fixFindingIDs(r.fixSpecs[0].Findings), want)
	}
	for _, s := range r.leafSpecs {
		for _, f := range s.Findings {
			if f.Leaf == "" || f.Leaf == "LX" {
				t.Errorf("leaf %s's %s got finding %+v", s.Leaf.Handle, s.Role, f)
			}
		}
	}
	entries := fixEntries(t, e)
	if entryAt(t, entries, "leaf-attempt leaf-b fixer 2 pass") > entryAt(t, entries, "review-fix fixer 2") {
		t.Errorf("the change set's fixer started before leaf-b's fixer was judged: %v", entries)
	}
	if entryAt(t, entries, "review-fix-judged fixer 2 pass") > entryAt(t, entries, "check-session alpha check 2") ||
		entryAt(t, entries, "leaf-attempt leaf-b verifier 2 pass") > entryAt(t, entries, "check-session alpha check 2") {
		t.Errorf("cycle 2's check started before the fix phase ended: %v", entries)
	}
}

// changeSetFindings, the change set's fixer's findings: the standing
// serious ones of the confirmed rounds that name no leaf of leaves.
func TestFixChangeSetFindings(t *testing.T) {
	a := fixLeaves[0]
	x := CheckState{Check: "x",
		Sessions: []CheckSession{{Role: "check", Cycle: 1}, {Role: "critic", Cycle: 1}},
		Judged:   []CheckJudged{{Role: "check", Cycle: 1, Outcome: "pass"}, {Role: "critic", Cycle: 1, Outcome: "pass"}},
		Rounds: []Round{{Number: 1, Cycle: 1, Findings: []Finding{
			sev("f1", SeverityMajor), onLeaf(sev("f2", SeverityMajor), a), sev("f3", SeverityNit), sev("f4", SeverityBlocker),
			{ID: "f5", Severity: SeverityMajor, Leaf: "LX"},
		}, Verdicts: []Verdict{{Finding: "f4", Verdict: Dropped}}}},
	}
	got := fixFindingIDs(changeSetFindings(2, []string{"x"}, map[string]CheckState{"x": x}, fixLeaves))
	if want := []string{"x f1 major", "x f5 major"}; !reflect.DeepEqual(got, want) {
		t.Errorf("changeSetFindings = %v, want %v", got, want)
	}
	if got := changeSetFindings(1, []string{"x"}, map[string]CheckState{"x": x}, fixLeaves); len(got) != 0 {
		t.Errorf("cycle 1 has findings to fix: %v", got)
	}
}

// nextCheckActions slots the change set's fixer after the leaves'
// fixers: it starts only once every leaf fixer is judged pass, alongside
// a live verifier; it is polled while live even after a leaf's failure;
// its failure ends the cycle's checks failed; the checks start once it
// passed.
func TestFixChangeSetActions(t *testing.T) {
	a := fixLeaves[0]
	step := workflow.Step{Name: "review", Checks: []workflow.Check{{Name: "x"}, {Name: "y"}}}
	cycles := []CheckCycle{{Step: "review", Cycle: 1, Checks: []string{"x", "y"}}, {Step: "review", Cycle: 2, Checks: []string{"x"}}}
	x := CheckState{Check: "x",
		Sessions: []CheckSession{{Role: "check", Cycle: 1, Session: "SX"}, {Role: "critic", Cycle: 1, Session: "CX"}},
		Judged:   []CheckJudged{{Role: "check", Cycle: 1, Outcome: "pass"}, {Role: "critic", Cycle: 1, Outcome: "pass"}},
		Outcomes: []CheckOutcome{{Cycle: 1, Outcome: CheckFailed}},
		Rounds:   []Round{{Number: 1, Cycle: 1, Findings: []Finding{onLeaf(sev("f1", SeverityMajor), a), sev("f2", SeverityMajor)}}},
	}
	y := CheckState{Check: "y", Outcomes: []CheckOutcome{{Cycle: 1, Outcome: CheckPassed}}}
	states := map[string]CheckState{"x": x, "y": y}
	leaf := func(fixer, verifier string) map[string]LeafState {
		l := LeafState{Leaf: a.ID, Handle: a.Handle, Sessions: []LeafSession{{Role: "fixer", Attempt: 2, Session: "F"}}}
		if fixer != "" {
			l.Attempts = append(l.Attempts, LeafAttempt{Role: "fixer", Attempt: 2, Outcome: fixer, Commit: "C"})
		}
		if verifier != "-" {
			l.Sessions = append(l.Sessions, LeafSession{Role: "verifier", Attempt: 2, Session: "V"})
			if verifier != "" {
				l.Attempts = append(l.Attempts, LeafAttempt{Role: "verifier", Attempt: 2, Outcome: verifier})
			}
		}
		return map[string]LeafState{a.ID: l}
	}
	fix := func(outcome string) []ReviewFix {
		f := ReviewFix{Step: "review", Cycle: 2, Attempt: 1, Session: "CF"}
		if outcome != "" {
			f.Judged, f.Outcome, f.Detail = true, outcome, "why"
		}
		return []ReviewFix{f}
	}
	exhausted := fix("fail")
	exhausted[0].Attempt = 1 + MaxFixCycles
	show := func(acts []checkAction) []string {
		var out []string
		for _, x := range acts {
			s := strconv.Itoa(int(x.kind)) + ":" + x.check + x.fix.role + ":" + x.outcome + x.detail
			if x.kind == reviewFixStart {
				s += strings.Join(fixFindingIDs(x.findings), ",")
			}
			out = append(out, s)
		}
		return out
	}
	for _, tt := range []struct {
		name  string
		fixes map[string]LeafState
		cs    []ReviewFix
		want  []string
	}{
		{"leaf fixer first", nil, nil, []string{"6:fixer:"}},
		{"not beside a live leaf fixer", leaf("", "-"), nil, []string{"7:fixer:"}},
		{"beside the leaf's verifier", leaf("pass", ""), nil, []string{"8::x f2 major", "7:verifier:"}},
		{"polled", leaf("pass", "pass"), fix(""), []string{"9::"}},
		{"polled after a leaf's failure", leaf("pass", "fail"), fix(""), []string{"9::"}},
		{"not started after a leaf's failure", leaf("fail", "-"), nil, []string{"4:x:failedcycle 2: fix failed: leaf leaf-a: fixer failed: "}},
		{"its failure is retried in the cycle", leaf("pass", "pass"), fix("fail"), []string{"8::x f2 major"}},
		{"its last failure ends the cycle", leaf("pass", "pass"), exhausted, []string{"4:x:failedcycle 2: fix failed: change set: fixer failed: why"}},
		{"checks after it passed", leaf("pass", "pass"), fix("pass"), []string{"2:x:"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := show(nextCheckActions(step, cycles, states, fixLeaves, tt.fixes, tt.cs)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("actions = %v, want %v", got, tt.want)
			}
		})
	}
	other := []ReviewFix{{Step: "review", Cycle: 3, Judged: true, Outcome: "pass"}}
	if got := show(nextCheckActions(step, cycles, states, fixLeaves, leaf("pass", "pass"), other)); !reflect.DeepEqual(got, []string{"8::x f2 major"}) {
		t.Errorf("another cycle's fixer counted: %v", got)
	}
}

// The fold keeps a change-set fixer's session and judgement per cycle,
// and a pane-closed entry naming its session marks its pane closed, with
// no leaf or check state made for it.
func TestReviewFixFold(t *testing.T) {
	j := Journal{ID: runID, Entries: []Entry{
		{Type: TypeEntered, Step: "review", At: t0},
		{Type: TypeReviewFix, Step: "review", Role: workflow.RoleFixer, Cycle: 2, Session: "S1", Model: "opus", Commit: "H0", At: t0},
		{Type: TypeReviewFixJudged, Step: "review", Role: workflow.RoleFixer, Cycle: 2, Outcome: workflow.OutcomePass, Commit: "C1", Detail: "said done", At: t0},
		{Type: TypePaneClosed, Step: "review", Role: workflow.RoleFixer, Cycle: 2, Session: "S1", At: t0},
	}}
	st := j.State()
	f, ok := st.ReviewFix("review", 2)
	if !ok || f.Session != "S1" || f.Model != "opus" || f.Since != "H0" || !f.Judged || f.Outcome != workflow.OutcomePass ||
		f.Commit != "C1" || f.Detail != "said done" || !f.PaneClosed || f.Seq != 1 {
		t.Errorf("ReviewFix = %+v, %v", f, ok)
	}
	if len(st.Leaves) != 0 || len(st.Checks) != 0 {
		t.Errorf("leaves %+v, checks %+v", st.Leaves, st.Checks)
	}
	if _, ok := st.ReviewFix("review", 3); ok {
		t.Error("a fixer in cycle 3")
	}
}
