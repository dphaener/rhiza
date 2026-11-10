package run

import (
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/workflow"
)

// signoffWF is a review: a checks step of two checks, then a signoff step
// whose changes go back to it.
const signoffWF = `[step.review]
kind = "checks"
interval = "5s"
max_visits = 3
on.pass = "signoff"
[[step.review.check]]
name = "alpha"
prompt = "check {check}"
[[step.review.check]]
name = "beta"
prompt = "check {check}"
[step.review.critic]
prompt = "criticise {check}"
[step.review.fixer]
prompt = "fix {leaf}"
[step.review.verifier]
prompt = "verify {leaf}"
[step.signoff]
kind = "signoff"
max_visits = 3
on.pass = "done"
on.changes = "review"
`

// appendSignoff journals a sign-off of the signoff step's visit, as rhiza
// run approve and request-changes do.
func appendSignoff(t *testing.T, e *Engine, visit int, outcome string, findings ...Finding) {
	t.Helper()
	if err := e.Store.Append(runID, Entry{
		Type: TypeSignoff, Step: "signoff", Visit: visit, Actor: "ada", Rule: RuleOwner, Outcome: outcome, Findings: findings, At: at(1),
	}); err != nil {
		t.Fatal(err)
	}
}

// signoffEntries are the journal's signoff entries.
func signoffEntries(t *testing.T, e *Engine) []Entry {
	t.Helper()
	j, err := e.Store.Read(runID)
	if err != nil {
		t.Fatal(err)
	}
	var out []Entry
	for _, en := range j.Entries {
		if en.Type == TypeSignoff {
			out = append(out, en)
		}
	}
	return out
}

// A review whose checks pass stops at its signoff step, waiting for the
// owner, and stays there on a resume, writing no sign-off; an approved
// sign-off of the visit then passes it and the run ends done.
func TestSignoffWaitsThenApproved(t *testing.T) {
	e, r, _, _ := setup(t, signoffWF)
	passChecks(r, 1, "alpha", "beta")
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Ended || st.Current != "signoff" || st.Stopped != SignoffWait {
		t.Fatalf("state: ended %v, current %q, stopped %q", st.Ended, st.Current, st.Stopped)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if st, err = e.Run(runID); err != nil || st.Stopped != SignoffWait || st.Ended {
		t.Fatalf("resume: %v, stopped %q", err, st.Stopped)
	}
	if got := signoffEntries(t, e); len(got) != 0 {
		t.Fatalf("signoff entries before the owner's: %+v", got)
	}
	types := entryTypes(t, e)
	if got := types[len(types)-1]; got != "stopped:signoff" {
		t.Errorf("last entry %s", got)
	}

	appendSignoff(t, e, 1, SignoffApproved)
	if st, err = e.Run(runID); err != nil {
		t.Fatal(err)
	}
	if !st.Ended || st.Result != Succeeded {
		t.Fatalf("after approve: ended %v result %q", st.Ended, st.Result)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:pass", "signoff:pass"}) {
		t.Errorf("outcomes = %v", got)
	}
	if d := st.Steps[1].Detail; d != "approved by ada" {
		t.Errorf("signoff detail %q", d)
	}
	if len(r.checkSpecs) != 2 {
		t.Errorf("check sessions %v: the approved review ran checks again", roles(r))
	}
}

// An autonomous run passes the signoff step at once, journalling a
// sign-off of rule autonomous.
func TestSignoffAutonomous(t *testing.T) {
	e, r := autonomousSetup(t, signoffWF)
	r.store = e.Store
	passChecks(r, 1, "alpha", "beta")
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if !st.Ended || st.Result != Succeeded {
		t.Fatalf("ended %v result %q stopped %q", st.Ended, st.Result, st.Stopped)
	}
	got := signoffEntries(t, e)
	if len(got) != 1 || got[0].Rule != RuleAutonomous || got[0].Outcome != SignoffApproved || got[0].Visit != 1 || got[0].Actor != "" {
		t.Fatalf("signoff entries %+v", got)
	}
	if d := st.Steps[1].Detail; d != "approved by rule autonomous" {
		t.Errorf("signoff detail %q", d)
	}
}

// A request for changes after the cap was used ends the signoff step
// changes, back to the checks step, which opens one owner cycle: the
// named leaf's fixer, then its verifier, get the owner's findings, then
// only the named check runs again; it passes, and the run waits at the
// sign-off again. A resume restarts none of the fix sessions.
func TestSignoffChangesOwnerCycle(t *testing.T) {
	e, r, _, _ := setup(t, signoffWF)
	r.leaves = fixLeaves
	a := fixLeaves[0]
	passChecks(r, 1, "beta")
	passChecks(r, 3, "alpha")
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Stopped != SignoffWait {
		t.Fatalf("stopped %q, outcomes %v", st.Stopped, outcomes(st))
	}
	if cy, _ := st.LatestCycle("review"); cy.Cycle != 1+MaxFixCycles {
		t.Fatalf("latest cycle %d; want the cap used", cy.Cycle)
	}

	f := onLeaf(Finding{ID: "f1", Check: "alpha", Severity: SeverityMajor, Place: "a.go:1", Reason: "the owner wants it"}, a)
	appendSignoff(t, e, 1, SignoffChanges, f)
	passFixes(r, 4, a)
	passChecks(r, 4, "alpha")
	if st, err = e.Run(runID); err != nil {
		t.Fatal(err)
	}
	if st.Ended || st.Stopped != SignoffWait || st.Current != "signoff" || st.CurrentVisit != 2 {
		t.Fatalf("after changes: ended %v stopped %q current %s visit %d, outcomes %v", st.Ended, st.Stopped, st.Current, st.CurrentVisit, outcomes(st))
	}
	if got, want := outcomes(st), []string{"review:pass", "signoff:changes", "review:pass"}; !reflect.DeepEqual(got, want) {
		t.Errorf("outcomes = %v, want %v", got, want)
	}
	if want := []string{"leaf-a fixer 8", "leaf-a verifier 8"}; !reflect.DeepEqual(leafSpecRoles(r), want) {
		t.Fatalf("leaf sessions %v, want %v", leafSpecRoles(r), want)
	}
	if fs := r.leafSpecs[0].Findings; len(fs) != 1 || fs[0].ID != "f1" || fs[0].Check != "alpha" || fs[0].Reason != "the owner wants it" {
		t.Errorf("fixer findings %+v", fs)
	}
	if got := roles(r); got[len(got)-1] != "alpha check 4" || slices.Contains(got, "beta check 4") {
		t.Errorf("check sessions %v", got)
	}
	cy, _ := st.LatestCycle("review")
	if !cy.Owner() || cy.Cycle != 4 || !reflect.DeepEqual(cy.Checks, []string{"alpha"}) || len(cy.Findings) != 1 {
		t.Errorf("owner cycle %+v", cy)
	}
	if _, ok := st.PendingChanges(); ok {
		t.Errorf("the request for changes is still pending")
	}
	if l, ok := st.Leaf("review", a.ID); !ok || len(l.Sessions) != 2 {
		t.Errorf("leaf-a's fix history across the re-entry: %+v, %v", l, ok)
	}

	if st, err = e.Run(runID); err != nil || st.Stopped != SignoffWait {
		t.Fatalf("resume: %v, stopped %q", err, st.Stopped)
	}
	if len(r.leafSpecs) != 2 {
		t.Errorf("a resume restarted fix sessions: %v", leafSpecRoles(r))
	}
}

// The owner cycle is one cycle outside the cap: a check failing it gets
// no further fix cycle, even with cycles of the cap left; its failure
// raises questions and the checks step fails.
func TestSignoffOwnerCycleFailsCapped(t *testing.T) {
	e, r, _, _ := setup(t, signoffWF)
	passChecks(r, 1, "alpha", "beta")
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	appendSignoff(t, e, 1, SignoffChanges, Finding{ID: "f1", Check: "beta", Severity: SeverityMinor, Place: "b.go:2", Reason: "tidy"})
	r.outcome["S-fixer-2"] = workflow.OutcomePass
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if !st.Ended || st.Result != Failed {
		t.Fatalf("ended %v result %q, outcomes %v", st.Ended, st.Result, outcomes(st))
	}
	if cy, _ := st.LatestCycle("review"); !cy.Owner() || cy.Cycle != 2 {
		t.Errorf("latest cycle %+v; want the owner cycle 2 and no cycle after it", cy)
	}
	if fx, ok := st.ReviewFix("review", 2); !ok || fx.Outcome != workflow.OutcomePass {
		t.Errorf("change-set fixer of the owner cycle: %+v, %v", fx, ok)
	}
	if len(r.reviewQuestions) != 1 {
		t.Errorf("review questions %v", r.reviewQuestions)
	}
}

// Re-entering the checks step for the owner's changes keeps its fix
// history: the fixer and verifier a cycle of the first visit ran stay in
// the fold, beside the owner cycle's, and none of them runs again.
func TestSignoffReentryKeepsFixHistory(t *testing.T) {
	e, r, _, _ := setup(t, signoffWF)
	r.leaves = fixLeaves
	a := fixLeaves[0]
	passChecks(r, 1, "alpha", "beta")
	r.findings[checkSID("alpha", "check", 1)] = []Finding{onLeaf(sev("f1", SeverityMajor), a)}
	r.outcome[checkSID("alpha", "critic", 1)] = workflow.OutcomePass
	passFixes(r, 2, a)
	passChecks(r, 2, "alpha")
	if st, err := e.Run(runID); err != nil || st.Stopped != SignoffWait {
		t.Fatalf("first visit: %v, stopped %q, outcomes %v", err, st.Stopped, outcomes(st))
	}
	appendSignoff(t, e, 1, SignoffChanges, onLeaf(Finding{ID: "f1", Check: "alpha", Severity: SeverityMajor, Place: "a.go:2", Reason: "again"}, a))
	passFixes(r, 3, a)
	passChecks(r, 3, "alpha")
	st, err := e.Run(runID)
	if err != nil || st.Stopped != SignoffWait {
		t.Fatalf("after changes: %v, stopped %q, outcomes %v", err, st.Stopped, outcomes(st))
	}
	want := []string{"leaf-a fixer 2", "leaf-a verifier 2", "leaf-a fixer 5", "leaf-a verifier 5"}
	if !reflect.DeepEqual(leafSpecRoles(r), want) {
		t.Fatalf("leaf sessions %v, want %v", leafSpecRoles(r), want)
	}
	l, ok := st.Leaf("review", a.ID)
	if !ok || len(l.Sessions) != 4 || len(l.Attempts) != 4 {
		t.Errorf("leaf-a across the re-entry: %d sessions, %d attempts (%v); want 4 and 4", len(l.Sessions), len(l.Attempts), ok)
	}
}

// signoffRun starts an autonomous run of src with the given recorded
// stage autonomy, over fakes.
func signoffRun(t *testing.T, src string, aut []StageAutonomy) (*Engine, *fakeRunner) {
	t.Helper()
	def, err := workflow.Parse("wf", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	store := newStore(t)
	if err := Start(store, runID, Started{Workflow: "wf", Definition: src, Change: "CS1", Title: "Sign", Autonomous: true, Autonomy: aut}, t0, def); err != nil {
		t.Fatal(err)
	}
	r := newFakeRunner()
	r.store = store
	c := &fakeClock{now: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
	return &Engine{Store: store, Runner: r, Clock: c, Now: func() time.Time { return at(1) }}, r
}

// An autonomous run whose workflow declared signoff owner waits for the
// owner at the signoff step, journalling no sign-off of rule autonomous;
// the owner's approval then passes it.
func TestSignoffWorkflowOwnerBeatsAutonomous(t *testing.T) {
	e, r := signoffRun(t, signoffWF, []StageAutonomy{{Stage: "signoff", Mode: "owner", Source: SourceWorkflow, Workflow: "docs"}})
	passChecks(r, 1, "alpha", "beta")
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Ended || st.Stopped != SignoffWait {
		t.Fatalf("ended %v stopped %q", st.Ended, st.Stopped)
	}
	if got := signoffEntries(t, e); len(got) != 0 {
		t.Fatalf("signoff entries %+v", got)
	}
	appendSignoff(t, e, 1, SignoffApproved)
	if st, err = e.Run(runID); err != nil || !st.Ended || st.Result != Succeeded {
		t.Fatalf("after approve: %+v, %v", st, err)
	}
}

// A workflow declaring signoff auto, or the setting (even owner), still
// signs off under rule autonomous in an autonomous run.
func TestSignoffAutoAndSettingSignOffAutonomously(t *testing.T) {
	for name, aut := range map[string][]StageAutonomy{
		"workflow auto":  {{Stage: "signoff", Mode: "auto", Source: SourceWorkflow, Workflow: "docs"}},
		"setting auto":   {{Stage: "signoff", Mode: "auto", Source: SourceSetting}},
		"setting owner":  {{Stage: "signoff", Mode: "owner", Source: SourceSetting}},
		"other workflow": {{Stage: "ship", Mode: "owner", Source: SourceWorkflow, Workflow: "docs"}},
	} {
		t.Run(name, func(t *testing.T) {
			e, r := signoffRun(t, signoffWF, aut)
			passChecks(r, 1, "alpha", "beta")
			st, err := e.Run(runID)
			if err != nil {
				t.Fatal(err)
			}
			if !st.Ended || st.Result != Succeeded {
				t.Fatalf("ended %v stopped %q", st.Ended, st.Stopped)
			}
			if got := signoffEntries(t, e); len(got) != 1 || got[0].Rule != RuleAutonomous {
				t.Fatalf("signoff entries %+v", got)
			}
		})
	}
}

// A journal without the autonomy record falls back to the started
// definition's declaration.
func TestSignoffOwnerFallsBackToDefinition(t *testing.T) {
	e, r := signoffRun(t, "[autonomy]\nsignoff = \"owner\"\n\n"+signoffWF, nil)
	passChecks(r, 1, "alpha", "beta")
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Ended || st.Stopped != SignoffWait {
		t.Fatalf("ended %v stopped %q", st.Ended, st.Stopped)
	}
	if got := signoffEntries(t, e); len(got) != 0 {
		t.Fatalf("signoff entries %+v", got)
	}
}

// Under a workflow-declared owner, the owner's request for changes sends
// an autonomous run back to the checks step, and its next signoff visit
// waits for the owner again rather than signing off under rule
// autonomous.
func TestSignoffWorkflowOwnerChangesGoBack(t *testing.T) {
	e, r := signoffRun(t, signoffWF, []StageAutonomy{{Stage: "signoff", Mode: "owner", Source: SourceWorkflow, Workflow: "docs"}})
	r.leaves = fixLeaves
	a := fixLeaves[0]
	passChecks(r, 1, "alpha", "beta")
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Stopped != SignoffWait {
		t.Fatalf("stopped %q, outcomes %v", st.Stopped, outcomes(st))
	}
	appendSignoff(t, e, 1, SignoffChanges, onLeaf(Finding{ID: "f1", Check: "alpha", Severity: SeverityMajor, Place: "a.go:1", Reason: "the owner wants it"}, a))
	passFixes(r, 2, a)
	passChecks(r, 2, "alpha")
	if st, err = e.Run(runID); err != nil {
		t.Fatal(err)
	}
	if st.Ended || st.Stopped != SignoffWait || st.Current != "signoff" || st.CurrentVisit != 2 {
		t.Fatalf("after changes: ended %v stopped %q current %s visit %d, outcomes %v", st.Ended, st.Stopped, st.Current, st.CurrentVisit, outcomes(st))
	}
	if got, want := outcomes(st), []string{"review:pass", "signoff:changes", "review:pass"}; !reflect.DeepEqual(got, want) {
		t.Errorf("outcomes = %v, want %v", got, want)
	}
	for _, s := range signoffEntries(t, e) {
		if s.Rule == RuleAutonomous {
			t.Errorf("autonomous sign-off %+v", s)
		}
	}
}
