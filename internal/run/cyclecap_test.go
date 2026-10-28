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

func (f *fakeRunner) AddReviewQuestion(text string) (Question, error) {
	if f.reviewQuestionErr != nil {
		return Question{}, f.reviewQuestionErr
	}
	f.reviewQuestions = append(f.reviewQuestions, text)
	n := strconv.Itoa(len(f.reviewQuestions))
	return Question{ID: "QR" + n, Handle: "review-failed-" + n}, nil
}

// capWF is checksWF followed, on pass, by a command step standing in for
// sign-off: a capped review must never reach it.
var capWF = strings.Replace(checksWF, "interval = \"5s\"\n", "interval = \"5s\"\non.pass = \"signoff\"\n", 1) + `
[step.signoff]
kind = "command"
run = "signoff"
`

// capSetup is setup over capWF, autonomous when asked.
func capSetup(t *testing.T, autonomous bool) (*Engine, *fakeRunner) {
	t.Helper()
	def, err := workflow.Parse("wf", []byte(capWF))
	if err != nil {
		t.Fatal(err)
	}
	store := newStore(t)
	if err := Start(store, runID, Started{Workflow: "wf", Definition: capWF, Change: "CS1", Title: "Capped", Autonomous: autonomous}, t0, def); err != nil {
		t.Fatal(err)
	}
	r := newFakeRunner()
	r.store = store
	c := &fakeClock{now: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
	return &Engine{Store: store, Runner: r, Clock: c, Now: func() time.Time { return at(1) }}, r
}

// failAlphaOnLeaf makes check alpha find confirmed major f1 on leaf-a in
// each of cycles, every other check pass in cycle 1, and every fix of
// leaf-a in cycles pass.
func failAlphaOnLeaf(r *fakeRunner, cycles ...int) {
	a := fixLeaves[0]
	r.leaves = fixLeaves
	passChecks(r, 1, threeChecks...)
	for _, n := range cycles {
		passChecks(r, n, "alpha")
		r.findings[checkSID("alpha", "check", n)] = []Finding{onLeaf(sev("f1", SeverityMajor), a)}
		r.outcome[checkSID("alpha", "critic", n)] = workflow.OutcomePass
		if n > 1 {
			passFixes(r, n, a)
		}
	}
}

// reviewQuestions are the journal's review-question entries.
func reviewQuestionEntries(t *testing.T, e *Engine) []Entry {
	t.Helper()
	var out []Entry
	for _, en := range mustEntries(t, e) {
		if en.Type == TypeReviewQuestion {
			out = append(out, en)
		}
	}
	return out
}

// The first part of leaf review-cycle-cap's acceptance: a check failing
// three times gets two fix cycles (cycles 2 and 3), then its surviving
// confirmed finding becomes a question linked blocks to the leaf it
// names, the step fails, the next step (sign-off) is never entered and
// the run fails.
func TestCycleCapRaisesQuestions(t *testing.T) {
	e, r := capSetup(t, false)
	failAlphaOnLeaf(r, 1, 2, 3)
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:fail"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if st.Result != Failed {
		t.Errorf("run result = %q, want failed", st.Result)
	}
	if r.commands["signoff"] != 0 || slices.Contains(entryTypes(t, e), "signoff") {
		t.Errorf("the sign-off step ran: %v", r.commands)
	}
	for _, en := range mustEntries(t, e) {
		if en.Type == TypeEntered && en.Step == "signoff" {
			t.Errorf("the sign-off step was entered")
		}
	}
	if want := []string{"leaf-a fixer 2", "leaf-a verifier 2", "leaf-a fixer 3", "leaf-a verifier 3"}; !reflect.DeepEqual(leafSpecRoles(r), want) {
		t.Errorf("fix cycles' sessions %v, want %v", leafSpecRoles(r), want)
	}
	if len(st.Cycles) != 3 {
		t.Errorf("cycles = %+v, want 3", st.Cycles)
	}
	if len(r.leafQuestions) != 1 || len(r.reviewQuestions) != 0 {
		t.Fatalf("leaf questions %v, review questions %v", r.leafQuestions, r.reviewQuestions)
	}
	q := r.leafQuestions[0]
	for _, s := range []string{"leaf-a ", "check alpha", "f1", "major", "a.go:1", "bad f1", "3 cycles", "leaf-a change"} {
		if !strings.Contains(q, s) {
			t.Errorf("question %q lacks %q", q, s)
		}
	}
	qs := reviewQuestionEntries(t, e)
	if len(qs) != 1 {
		t.Fatalf("review-question entries = %+v", qs)
	}
	if en := qs[0]; en.Step != "review" || en.Check != "alpha" || en.Cycle != 3 || en.Finding != "f1" ||
		en.Leaf != "L1" || en.Handle != "leaf-a" || en.Question != "QL1" || en.QuestionHandle != "leaf-failed-1" {
		t.Errorf("review-question entry = %+v", en)
	}
	if got := st.ReviewQuestionsFor("review", 3); len(got) != 1 || got[0].Question != "QL1" {
		t.Errorf("ReviewQuestionsFor = %+v", got)
	}
	if d := st.Steps[0].Detail; d != "2 of 3 checks passed after 3 cycles; failed: alpha; unanswered questions: leaf-failed-1" {
		t.Errorf("detail = %q", d)
	}
}

// A check failing with no confirmed finding to name (its session failed)
// and a finding naming no leaf of the change set each raise a question in
// the run's change set (AddReviewQuestion), one per finding or failure.
// One that cannot be added is named in the step's detail, and the run
// still fails; a resume adds it then, and never raises one twice.
func TestCycleCapQuestionsWithoutLeaf(t *testing.T) {
	e, r := capSetup(t, false)
	passChecks(r, 1, "alpha", "gamma")
	for _, n := range []int{1, 2, 3} {
		passChecks(r, n, "alpha")
		r.findings[checkSID("alpha", "check", n)] = []Finding{sev("f1", SeverityBlocker), sev("f2", SeverityMinor)}
		r.outcome[checkSID("alpha", "critic", n)] = workflow.OutcomePass
		if n > 1 {
			r.outcome[fixSID(n)] = workflow.OutcomePass
		}
	}
	r.reviewQuestionErr = errors.New("change set CS1 is landed")
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed {
		t.Fatalf("run result = %q", st.Result)
	}
	d := st.Steps[0].Detail
	for _, s := range []string{"failed: alpha, beta", "check alpha finding f1: its question could not be added: change set CS1 is landed",
		"check beta: its question could not be added"} {
		if !strings.Contains(d, s) {
			t.Errorf("detail %q lacks %q", d, s)
		}
	}
	if n := len(reviewQuestionEntries(t, e)); n != 0 {
		t.Errorf("%d review-question entries for questions never added", n)
	}

	r.reviewQuestionErr = nil
	if st, err = e.Retry(runID); err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed {
		t.Fatalf("retried run result = %q", st.Result)
	}
	if len(r.reviewQuestions) != 2 {
		t.Fatalf("review questions = %v", r.reviewQuestions)
	}
	if q := r.reviewQuestions[0]; !strings.Contains(q, "check alpha") || !strings.Contains(q, "f1") || !strings.Contains(q, "blocker") {
		t.Errorf("alpha's question = %q", q)
	}
	if q := r.reviewQuestions[1]; !strings.Contains(q, "check beta") || !strings.Contains(q, "check failed") {
		t.Errorf("beta's question = %q", q)
	}
	var got []string
	for _, en := range reviewQuestionEntries(t, e) {
		got = append(got, en.Check+" "+en.Finding+" "+en.Question)
	}
	if want := []string{"alpha f1 QR1", "beta  QR2"}; !reflect.DeepEqual(got, want) {
		t.Errorf("review-question entries = %v, want %v", got, want)
	}

	if _, err = e.Retry(runID); err != nil {
		t.Fatal(err)
	}
	if len(r.reviewQuestions) != 2 || len(reviewQuestionEntries(t, e)) != 2 {
		t.Errorf("a resume raised the questions again: %v", r.reviewQuestions)
	}
}

// The second part: once a Decision answers the question, a resume journals
// an answered check-cycle and runs one fresh cycle, whose fixer is given
// the answering Decision; the check passing then, the review passes and
// sign-off follows. Resumed before any answer, the step fails again at
// once, naming the open question, and runs nothing.
func TestCycleCapAnsweredResume(t *testing.T) {
	e, r := capSetup(t, false)
	failAlphaOnLeaf(r, 1, 2, 3)
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	before := len(r.checkSpecs) + len(r.leafSpecs)
	st, err := e.Retry(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed || len(r.checkSpecs)+len(r.leafSpecs) != before {
		t.Fatalf("an unanswered resume ran sessions or passed: %s", st.Result)
	}
	if d := st.Steps[len(st.Steps)-1].Detail; !strings.HasSuffix(d, "unanswered questions: leaf-failed-1") {
		t.Errorf("detail = %q", d)
	}
	if len(r.leafQuestions) != 1 {
		t.Errorf("questions raised again: %v", r.leafQuestions)
	}

	r.leafAnswers["QL1"] = [2]string{"rename-it", "Rename the flag instead."}
	passFixes(r, 4, fixLeaves[0])
	passChecks(r, 4, "alpha")
	if st, err = e.Retry(runID); err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || r.commands["signoff"] != 1 {
		t.Fatalf("result %q, signoff ran %d times", st.Result, r.commands["signoff"])
	}
	cy, _ := st.LatestCycle("review")
	if cy.Cycle != 4 || !cy.Answered() || !reflect.DeepEqual(cy.Checks, []string{"alpha"}) ||
		!reflect.DeepEqual(cy.Questions, []string{"QL1"}) || !reflect.DeepEqual(cy.Decisions, []string{"rename-it"}) {
		t.Errorf("answered cycle = %+v", cy)
	}
	var fixer LeafSpec
	for _, s := range r.leafSpecs {
		if s.Role == workflow.RoleFixer && s.Cycle == 4 {
			fixer = s
		}
	}
	if want := []ReviewAnswer{{Question: "QL1", QuestionHandle: "leaf-failed-1", Decision: "rename-it", Statement: "Rename the flag instead."}}; !reflect.DeepEqual(fixer.Answers, want) {
		t.Errorf("cycle 4 fixer's answers = %+v, want %+v", fixer.Answers, want)
	}
	if len(fixer.Findings) != 1 || fixer.Findings[0].ID != "f1" {
		t.Errorf("cycle 4 fixer's findings = %+v", fixer.Findings)
	}
	for _, en := range mustEntries(t, e) {
		if strings.Contains(en.Detail, "Rename the flag") {
			t.Errorf("the Decision's statement is in the journal: %+v", en)
		}
	}
}

// The third part: in an autonomous run the arbiter is started on the
// question at once; its answer gives the review one fresh cycle in the
// same visit, the fixer told the arbiter's Decision. The check failing
// again ends the run failed, with no second arbitration, then or on a
// resume.
func TestCycleCapAutonomous(t *testing.T) {
	e, r := capSetup(t, true)
	failAlphaOnLeaf(r, 1, 2, 3, 4)
	r.answers["QL1"] = Answer{OK: true, Decision: "arbiter-ruling", Statement: "Keep the flag."}
	r.leafAnswers["QL1"] = [2]string{"arbiter-ruling", "Keep the flag."}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed || r.commands["signoff"] != 0 {
		t.Fatalf("result %q, signoff ran %d times", st.Result, r.commands["signoff"])
	}
	if want := []string{"arbiter QL1"}; !reflect.DeepEqual(r.arbiters, want) {
		t.Errorf("arbiters = %v, want %v", r.arbiters, want)
	}
	if want := []string{"leaf-a fixer 2", "leaf-a verifier 2", "leaf-a fixer 3", "leaf-a verifier 3", "leaf-a fixer 4", "leaf-a verifier 4"}; !reflect.DeepEqual(leafSpecRoles(r), want) {
		t.Errorf("leaf sessions %v, want %v", leafSpecRoles(r), want)
	}
	if len(st.Cycles) != 4 || !st.Cycles[3].Answered() || !reflect.DeepEqual(st.Cycles[3].Decisions, []string{"arbiter-ruling"}) {
		t.Errorf("cycles = %+v", st.Cycles)
	}
	if fx := r.leafSpecs[4]; len(fx.Answers) != 1 || fx.Answers[0].Decision != "arbiter-ruling" || fx.Answers[0].Statement != "Keep the flag." {
		t.Errorf("cycle 4 fixer's answers = %+v", fx.Answers)
	}
	if n := countType(t, e, TypeEntered); n != 1 {
		t.Errorf("%d entered entries, want the one visit", n)
	}
	if len(r.leafQuestions) != 2 {
		t.Errorf("leaf questions = %v, want cycle 3's and cycle 4's", r.leafQuestions)
	}

	if _, err := e.Retry(runID); err != nil {
		t.Fatal(err)
	}
	if len(r.arbiters) != 1 || len(st.Cycles) != 4 {
		t.Errorf("a resume arbitrated again: %v", r.arbiters)
	}
	if st, _ := e.Store.Read(runID); len(st.State().Cycles) != 4 {
		t.Errorf("a resume ran another cycle: %+v", st.State().Cycles)
	}
}
