package run

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/workflow"
)

// onLeaf is f with its finding naming leaf lf.
func onLeaf(f Finding, lf Leaf) Finding {
	f.Leaf, f.LeafHandle = lf.ID, lf.Handle
	return f
}

// fixLeaves are two independent leaves for the fix phase's tests.
var fixLeaves = []Leaf{{ID: "L1", Handle: "leaf-a"}, {ID: "L2", Handle: "leaf-b"}}

// passFixes makes the fixer and verifier of each leaf exit 0 in cycle.
func passFixes(r *fakeRunner, cycle int, leaves ...Leaf) {
	for _, l := range leaves {
		for _, role := range []string{workflow.RoleFixer, workflow.RoleVerifier} {
			r.outcome[leafSID(l.Handle, role, cycle)] = workflow.OutcomePass
		}
	}
}

// leafSpecRoles are each started leaf session's "handle role attempt".
func leafSpecRoles(r *fakeRunner) []string {
	var out []string
	for _, s := range r.leafSpecs {
		out = append(out, s.Leaf.Handle+" "+s.Role+" "+strconv.Itoa(s.Attempt))
	}
	return out
}

// fixEntries are the journal's checks and leaf entries as "type
// check-or-handle role cycle-or-attempt outcome", in order.
func fixEntries(t *testing.T, e *Engine) []string {
	t.Helper()
	j, err := e.Store.Read(runID)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, en := range j.Entries {
		switch en.Type {
		case TypeCheckCycle:
			out = append(out, strings.Join(append([]string{en.Type, attemptStr(en.Cycle)}, en.Checks...), " "))
		case TypeCheckSession, TypeCheckJudged, TypeCheck:
			out = append(out, strings.Join(strings.Fields(strings.Join([]string{en.Type, en.Check, en.Role, attemptStr(en.Cycle), en.Outcome}, " ")), " "))
		case TypeReviewFix, TypeReviewFixJudged:
			out = append(out, strings.Join(strings.Fields(strings.Join([]string{en.Type, en.Role, attemptStr(en.Cycle), en.Outcome}, " ")), " "))
		case TypeLeafSession, TypeLeafAttempt, TypeLeaf:
			out = append(out, strings.Join(strings.Fields(strings.Join([]string{en.Type, en.Handle, en.Role, attemptStr(en.Attempt), en.Outcome}, " ")), " "))
		}
	}
	return out
}

// Confirmed findings naming leaves go, after the cycle, to one fixer per
// leaf in leaf order, one at a time, in the run's worktree, each told its
// leaf's findings; each fixed leaf's verifier checks the fix commit
// alongside the next fixer, and only once every fixer and verifier is
// judged do the failed checks, and only they, run again. No leaf outcome
// is written.
func TestFixLeafOrder(t *testing.T) {
	e, r, _, _ := setup(t, checksWF)
	r.leaves = fixLeaves
	a, b := fixLeaves[0], fixLeaves[1]
	passChecks(r, 1, threeChecks...)
	passChecks(r, 2, "alpha", "gamma")
	r.findings[checkSID("alpha", "check", 1)] = []Finding{onLeaf(sev("f1", SeverityMajor), b), onLeaf(sev("f2", SeverityMajor), a), onLeaf(sev("f3", SeverityMinor), a)}
	r.findings[checkSID("gamma", "check", 1)] = []Finding{onLeaf(sev("f1", SeverityBlocker), a)}
	for _, c := range []string{"alpha", "gamma"} {
		r.outcome[checkSID(c, "critic", 1)] = workflow.OutcomePass
	}
	passFixes(r, 2, a, b)
	r.onPoll = func(session string, n int) {
		if session == leafSID("leaf-a", "verifier", 2) && n == 0 && !slices.Contains(leafSpecRoles(r), "leaf-b fixer 2") {
			t.Errorf("leaf-a's verifier polled before leaf-b's fixer started: %v", leafSpecRoles(r))
		}
		if session == leafSID("leaf-b", "fixer", 2) && n == 0 && len(r.leafSpecs) != 3 {
			t.Errorf("leaf-b's fixer polled with %v started", leafSpecRoles(r))
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
	want := []string{"alpha check 1", "beta check 1", "gamma check 1", "alpha critic 1", "gamma critic 1", "alpha check 2", "gamma check 2"}
	if !reflect.DeepEqual(roles(r), want) {
		t.Errorf("check sessions %v, want %v", roles(r), want)
	}

	fa := r.leafSpecs[0]
	if fa.Dir != "" || fa.Commit != "" || fa.Cycle != 2 || fa.Step.Name != "review" {
		t.Errorf("leaf-a's fixer spec = %+v", fa)
	}
	var got []string
	for _, f := range fa.Findings {
		got = append(got, f.Check+" "+f.ID+" "+f.Severity)
	}
	if want := []string{"alpha f2 major", "gamma f1 blocker"}; !reflect.DeepEqual(got, want) {
		t.Errorf("leaf-a's fixer findings = %v, want %v", got, want)
	}
	for _, v := range []string{EnvLeaf + "=leaf-a", EnvAttempt + "=2", EnvCycle + "=2", EnvRole + "=fixer"} {
		if !slices.Contains(fa.Env, v) {
			t.Errorf("leaf-a's fixer env %v lacks %s", fa.Env, v)
		}
	}
	if fb := r.leafSpecs[2]; len(fb.Findings) != 1 || fb.Findings[0].ID != "f1" || fb.Findings[0].Check != "alpha" {
		t.Errorf("leaf-b's fixer findings = %+v", fb.Findings)
	}
	va := r.leafSpecs[1]
	if va.Role != workflow.RoleVerifier || va.Commit != "C-leaf-a" || va.Dir != checkoutDir("leaf-a", 2) || len(va.Findings) != 0 {
		t.Errorf("leaf-a's verifier spec = %+v", va)
	}
	if !slices.Contains(va.Env, EnvCommit+"=C-leaf-a") || !slices.Contains(va.Env, EnvCycle+"=2") {
		t.Errorf("leaf-a's verifier env = %v", va.Env)
	}
	if want := []string{"leaf-a " + fakeHead, "leaf-b " + fakeHead}; !reflect.DeepEqual(r.leafSince, want) {
		t.Errorf("LeafCommit calls = %v, want %v", r.leafSince, want)
	}
	if want := []string{"leaf-a " + leafSID("leaf-a", "verifier", 2), "leaf-b " + leafSID("leaf-b", "verifier", 2)}; !reflect.DeepEqual(r.commitCalls, want) {
		t.Errorf("CommitEvidence calls = %v, want %v", r.commitCalls, want)
	}

	entries := fixEntries(t, e)
	at := func(s string) int {
		t.Helper()
		i := slices.Index(entries, s)
		if i < 0 {
			t.Fatalf("no %q in %v", s, entries)
		}
		return i
	}
	if at("check-cycle 2 alpha gamma") > at("leaf-session leaf-a fixer 2") {
		t.Errorf("the fixer started before cycle 2 began: %v", entries)
	}
	for _, s := range []string{"leaf-attempt leaf-a verifier 2 pass", "leaf-attempt leaf-b verifier 2 pass"} {
		if at(s) > at("check-session alpha check 2") {
			t.Errorf("a check of cycle 2 started before %q: %v", s, entries)
		}
	}
	for _, en := range entries {
		if strings.HasPrefix(en, TypeLeaf+" ") {
			t.Errorf("a leaf outcome was journalled: %s", en)
		}
	}
	l, _ := st.Leaf("review", "L1")
	if j, _ := l.Judged(workflow.RoleFixer, 2); j.Outcome != workflow.OutcomePass || j.Commit != "C-leaf-a" {
		t.Errorf("leaf-a's fixer judgement = %+v", j)
	}
	if j, _ := l.Judged(workflow.RoleVerifier, 2); j.Outcome != workflow.OutcomePass || j.Commit != "E-leaf-a" {
		t.Errorf("leaf-a's verifier judgement = %+v", j)
	}
	if !slices.Contains(r.checkouts, "remove "+checkoutDir("leaf-a", 2)) {
		t.Errorf("leaf-a's verifier checkout not removed: %v", r.checkouts)
	}
}

// A fixer that ends with no new commit naming its leaf fails the cycle:
// no later fixer starts, no verifier runs, and each check of the cycle
// fails without a session. The next cycle tries the fix again; with no
// cycle left the step fails.
func TestFixLeafNoCommitFailsCycle(t *testing.T) {
	e, r, _, _ := setup(t, checksWF)
	r.leaves = fixLeaves
	a, b := fixLeaves[0], fixLeaves[1]
	passChecks(r, 1, threeChecks...)
	r.findings[checkSID("alpha", "check", 1)] = []Finding{onLeaf(sev("f1", SeverityMajor), a), onLeaf(sev("f2", SeverityMajor), b)}
	r.outcome[checkSID("alpha", "critic", 1)] = workflow.OutcomePass
	passFixes(r, 2, a, b)
	passFixes(r, 3, a, b)
	r.leafCommit = map[string][2]string{"leaf-a": {"", "no new commit since H0"}}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:fail"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if want := []string{"leaf-a fixer 2", "leaf-a fixer 3"}; !reflect.DeepEqual(leafSpecRoles(r), want) {
		t.Errorf("leaf sessions %v, want %v", leafSpecRoles(r), want)
	}
	if want := []string{"alpha check 1", "beta check 1", "gamma check 1", "alpha critic 1"}; !reflect.DeepEqual(roles(r), want) {
		t.Errorf("check sessions %v, want %v", roles(r), want)
	}
	c, _ := st.Check("review", "alpha")
	o, _ := c.Outcome(2)
	if want := "cycle 2: fix failed: leaf leaf-a: fixer failed: no new commit since H0"; o.Outcome != CheckFailed || o.Detail != want {
		t.Errorf("alpha's cycle 2 outcome = %+v, want detail %q", o, want)
	}
	if o, _ := c.Outcome(3); o.Outcome != CheckFailed || !strings.HasPrefix(o.Detail, "cycle 3: fix failed: ") {
		t.Errorf("alpha's cycle 3 outcome = %+v", o)
	}
	if s := st.Steps[0]; s.Detail != "2 of 3 checks passed after 3 cycles; failed: alpha" {
		t.Errorf("step detail = %q", s.Detail)
	}
}

// A fixed leaf's verifier that finds problems, or writes no Evidence of
// its own, fails the cycle too.
func TestFixLeafVerifierFailsCycle(t *testing.T) {
	e, r, _, _ := setup(t, checksWF)
	r.leaves = fixLeaves
	a := fixLeaves[0]
	passChecks(r, 1, threeChecks...)
	passChecks(r, 3, "beta")
	r.findings[checkSID("beta", "check", 1)] = []Finding{onLeaf(sev("f1", SeverityMajor), a)}
	r.outcome[checkSID("beta", "critic", 1)] = workflow.OutcomePass
	passFixes(r, 2, a)
	passFixes(r, 3, a)
	r.evidence = map[string]fakeEvidence{"leaf-a": {}}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	c, _ := st.Check("review", "beta")
	o, _ := c.Outcome(2)
	if want := "cycle 2: fix failed: leaf leaf-a: verifier failed: verifier ended with neither Evidence nor findings"; o.Outcome != CheckFailed || o.Detail != want {
		t.Errorf("beta's cycle 2 outcome = %+v, want %q", o, want)
	}
	if want := []string{"leaf-a fixer 2", "leaf-a verifier 2", "leaf-a fixer 3", "leaf-a verifier 3"}; !reflect.DeepEqual(leafSpecRoles(r), want) {
		t.Errorf("leaf sessions %v, want %v", leafSpecRoles(r), want)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:fail"}) {
		t.Errorf("outcomes = %v", got)
	}
}

// Findings that name no leaf, and those the critic dropped or that are
// only minor, start no leaf fixer: one naming no leaf goes to the change
// set's fixer instead, and the failed check runs again.
func TestFixLeafNoneWithoutLeafFindings(t *testing.T) {
	e, r, _, _ := setup(t, checksWF)
	r.leaves = fixLeaves
	a := fixLeaves[0]
	passChecks(r, 1, threeChecks...)
	passChecks(r, 2, "alpha")
	r.findings[checkSID("alpha", "check", 1)] = []Finding{sev("f1", SeverityMajor), onLeaf(sev("f2", SeverityMajor), a), onLeaf(sev("f3", SeverityNit), a)}
	crit := checkSID("alpha", "critic", 1)
	r.outcome[crit] = workflow.OutcomePass
	r.outcome[fixSID(2)] = workflow.OutcomePass
	r.verdicts[crit] = []Verdict{{Finding: "f1", Verdict: Confirmed, Reason: "real"}, {Finding: "f2", Verdict: Dropped, Reason: "no"}}
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
	if len(r.fixSpecs) != 1 || !reflect.DeepEqual(fixFindingIDs(r.fixSpecs[0].Findings), []string{"alpha f1 major"}) {
		t.Errorf("change-set fixers = %+v", r.fixSpecs)
	}
}

// fixActions, the fix phase's pure decision: fixers one at a time in
// leaf order; a fixed leaf's verifier starts alongside the next fixer;
// nothing more starts after a failure; done once all are judged.
func TestFixActions(t *testing.T) {
	a, b := fixLeaves[0], fixLeaves[1]
	found := map[string][]Finding{a.ID: {sev("f1", SeverityMajor)}, b.ID: {sev("f2", SeverityMajor)}}
	sess := func(l LeafState, role string, judged string, commit string) LeafState {
		l.Sessions = append(l.Sessions, LeafSession{Role: role, Attempt: 2, Session: role + "-" + l.Handle})
		if judged != "" {
			l.Attempts = append(l.Attempts, LeafAttempt{Role: role, Attempt: 2, Outcome: judged, Commit: commit})
		}
		return l
	}
	la, lb := LeafState{Leaf: a.ID, Handle: a.Handle}, LeafState{Leaf: b.ID, Handle: b.Handle}
	show := func(acts []checkAction) []string {
		var out []string
		for _, x := range acts {
			s := strconv.Itoa(int(x.kind)) + ":" + x.fix.leaf.Handle + ":" + x.fix.role
			if x.fix.commit != "" {
				s += "@" + x.fix.commit
			}
			out = append(out, s)
		}
		return out
	}
	for _, tt := range []struct {
		name    string
		fixes   map[string]LeafState
		want    []string
		failure string
	}{
		{"first fixer in leaf order", nil, []string{"6:leaf-a:fixer"}, ""},
		{"one fixer at a time", map[string]LeafState{a.ID: sess(la, "fixer", "", "")}, []string{"7:leaf-a:fixer"}, ""},
		{"verifier alongside the next fixer", map[string]LeafState{a.ID: sess(la, "fixer", "pass", "C1")},
			[]string{"6:leaf-a:verifier@C1", "6:leaf-b:fixer"}, ""},
		{"polls", map[string]LeafState{a.ID: sess(sess(la, "fixer", "pass", "C1"), "verifier", "", ""), b.ID: sess(lb, "fixer", "", "")},
			[]string{"7:leaf-a:verifier", "7:leaf-b:fixer"}, ""},
		{"fixer failed: nothing more starts", map[string]LeafState{a.ID: sess(la, "fixer", "fail", "")}, nil, "leaf leaf-a: fixer failed: "},
		{"failure waits for live sessions", map[string]LeafState{
			a.ID: sess(sess(la, "fixer", "pass", "C1"), "verifier", "", ""), b.ID: sess(lb, "fixer", "fail", "")},
			[]string{"7:leaf-a:verifier"}, "leaf leaf-b: fixer failed: "},
		{"all judged", map[string]LeafState{
			a.ID: sess(sess(la, "fixer", "pass", "C1"), "verifier", "pass", "E1"),
			b.ID: sess(sess(lb, "fixer", "pass", "C2"), "verifier", "pass", "E2")}, nil, ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			acts, failure := fixActions(2, fixLeaves, found, tt.fixes)
			if got := show(acts); !reflect.DeepEqual(got, tt.want) || failure != tt.failure {
				t.Errorf("fixActions = %v, %q; want %v, %q", got, failure, tt.want, tt.failure)
			}
		})
	}
}

// nextCheckActions opens a cycle past the first with its fix phase: the
// first fixer, before any check of the cycle starts; the checks start
// only once every fixer and verifier is judged; a failed fix ends each
// check of the cycle failed, and then the next cycle begins.
func TestNextCheckActionsFixPhase(t *testing.T) {
	a := fixLeaves[0]
	step := workflow.Step{Name: "review", Checks: []workflow.Check{{Name: "x"}, {Name: "y"}}}
	cycles := []CheckCycle{{Step: "review", Cycle: 1, Checks: []string{"x", "y"}}, {Step: "review", Cycle: 2, Checks: []string{"x"}}}
	x := CheckState{Check: "x",
		Sessions: []CheckSession{{Role: "check", Cycle: 1, Session: "SX"}, {Role: "critic", Cycle: 1, Session: "CX"}},
		Judged:   []CheckJudged{{Role: "check", Cycle: 1, Outcome: "pass"}, {Role: "critic", Cycle: 1, Outcome: "pass"}},
		Outcomes: []CheckOutcome{{Cycle: 1, Outcome: CheckFailed}},
		Rounds:   []Round{{Number: 1, Cycle: 1, Findings: []Finding{onLeaf(sev("f1", SeverityMajor), a)}}},
	}
	y := CheckState{Check: "y", Outcomes: []CheckOutcome{{Cycle: 1, Outcome: CheckPassed}}}
	states := map[string]CheckState{"x": x, "y": y}
	fixed := func(fixer, verifier string) map[string]LeafState {
		l := LeafState{Leaf: a.ID, Handle: a.Handle, Sessions: []LeafSession{{Role: "fixer", Attempt: 2, Session: "F"}}}
		l.Attempts = append(l.Attempts, LeafAttempt{Role: "fixer", Attempt: 2, Outcome: fixer, Commit: "C"})
		if verifier != "" {
			l.Sessions = append(l.Sessions, LeafSession{Role: "verifier", Attempt: 2, Session: "V"})
			l.Attempts = append(l.Attempts, LeafAttempt{Role: "verifier", Attempt: 2, Outcome: verifier})
		}
		return map[string]LeafState{a.ID: l}
	}
	show := func(acts []checkAction) []string {
		var out []string
		for _, x := range acts {
			out = append(out, strconv.Itoa(int(x.kind))+":"+x.check+x.fix.role+":"+x.outcome+x.detail)
		}
		return out
	}
	for _, tt := range []struct {
		name  string
		fixes map[string]LeafState
		want  []string
	}{
		{"fixer first", nil, []string{"6:fixer:"}},
		{"checks after the fix phase", fixed("pass", "pass"), []string{"2:x:"}},
		{"verifier before the checks", fixed("pass", ""), []string{"6:verifier:"}},
		{"failed fix ends the cycle's checks", fixed("fail", ""), []string{"4:x:failedcycle 2: fix failed: leaf leaf-a: fixer failed: "}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := show(nextCheckActions(step, cycles, states, fixLeaves, tt.fixes, nil)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("actions = %v, want %v", got, tt.want)
			}
		})
	}
	ended := x
	ended.Outcomes = append(ended.Outcomes, CheckOutcome{Cycle: 2, Outcome: CheckFailed})
	got := show(nextCheckActions(step, cycles, map[string]CheckState{"x": ended, "y": y}, fixLeaves, fixed("fail", ""), nil))
	if want := []string{"1::"}; !reflect.DeepEqual(got, want) {
		t.Errorf("after the failed fix = %v, want the next cycle %v", got, want)
	}
}
