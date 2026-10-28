package api

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

// AddReviewQuestion writes a question node, linked to nothing, into the
// run's change set; once that change set has landed, into a landed record
// change set of its own, each handle the next free.
func TestAddReviewQuestion(t *testing.T) {
	_, r, wctx := arbiterFixture(t)
	cs := mustParseID(r.change)
	q, err := r.AddReviewQuestion("Review check alpha still fails after 3 cycles. How?")
	if err != nil {
		t.Fatal(err)
	}
	if q.Handle != "review-failed-01runarb-1" {
		t.Errorf("handle %q", q.Handle)
	}
	ww, err := r.worktreeWorkspace(false)
	if err != nil {
		t.Fatal(err)
	}
	m := ww.model(cs)
	n, ok := m.Node(mustParseID(q.ID))
	if !ok || n.Kind != "question" || fieldText(n.Fields["statement"]) != "Review check alpha still fails after 3 cycles. How?" {
		t.Fatalf("question node %+v", n)
	}
	c, _ := ww.changeSet(cs)
	for _, o := range c.Ops {
		if strings.HasPrefix(string(o.Type), "link") {
			t.Errorf("the change set links something: %+v", o)
		}
	}
	if !questionCreatedIn(c, q.Handle) {
		t.Errorf("the question is not created in the run's change set %s", cs)
	}

	for _, s := range []string{"propose", "accept", "land"} {
		move(t, wctx, s, `{}`)
	}
	q2, err := r.AddReviewQuestion("Again, landed?")
	if err != nil {
		t.Fatal(err)
	}
	if q2.Handle != "review-failed-01runarb-2" {
		t.Errorf("landed handle %q", q2.Handle)
	}
	if ww, err = r.worktreeWorkspace(false); err != nil {
		t.Fatal(err)
	}
	if c, _ := ww.changeSet(cs); questionCreatedIn(c, q2.Handle) {
		t.Errorf("the question went into the landed change set")
	}
	home, ok := questionHome(ww, mustParseID(q2.ID))
	if !ok || home.Open.Kind != changeset.KindRecord || home.State != changeset.StateLanded {
		t.Fatalf("the landed run's question lives in %+v", home)
	}
	if n, ok := ww.model(id.ID{}).Node(mustParseID(q2.ID)); !ok || fieldText(n.Fields["statement"]) != "Again, landed?" {
		t.Errorf("landed question node %+v", n)
	}
}

// questionCreatedIn says whether c's ops create a node with handle h.
func questionCreatedIn(c changeset.ChangeSet, h string) bool {
	for _, o := range c.Ops {
		if o.Type == "node.create" && o.Handle == h {
			return true
		}
	}
	return false
}

// A fixer of an answered cycle is told each answering Decision, a leaf's
// and the change set's alike; any other fixer is not.
func TestCycleCapFixerPrompts(t *testing.T) {
	role := workflow.Role{Prompt: "fix {leaf}"}
	vars := workflow.PromptVars{Run: "R1"}
	answers := []run.ReviewAnswer{
		{Question: "Q1", QuestionHandle: "leaf-failed-r1-1", Decision: "rename-it", Statement: "Rename the\nflag."},
		{Question: "Q2", Decision: "keep-it"},
	}
	fs := []run.Finding{{ID: "f1", Check: "alpha", Severity: "major", Reason: "bad"}}
	leaf := fixerPrompt(role, run.LeafSpec{Role: workflow.RoleFixer, Leaf: run.Leaf{Handle: "leaf-a"}, Cycle: 4, Findings: fs, Answers: answers}, vars)
	cs := changeSetFixerPrompt(role, run.FixSpec{Cycle: 4, Findings: fs, Answers: answers}, vars)
	for _, p := range []string{leaf, cs} {
		for _, s := range []string{"- Decision rename-it answers question leaf-failed-r1-1: Rename the flag.", "- Decision keep-it answers question Q2"} {
			if !strings.Contains(p, s) {
				t.Errorf("prompt lacks %q:\n%s", s, p)
			}
		}
	}
	plain := fixerPrompt(role, run.LeafSpec{Role: workflow.RoleFixer, Leaf: run.Leaf{Handle: "leaf-a"}, Cycle: 2, Findings: fs}, vars)
	if strings.Contains(plain, "Decision") {
		t.Errorf("a fixer within the cap is told of Decisions:\n%s", plain)
	}
}

// run status lists the questions a capped check raised and the Decisions
// that gave it a fresh cycle.
func TestCycleCapStatus(t *testing.T) {
	j := run.Journal{ID: "R1", Entries: []run.Entry{
		{Type: run.TypeCheckCycle, Step: "review", Cycle: 3, Checks: []string{"alpha"}},
		{Type: run.TypeCheck, Step: "review", Check: "alpha", Cycle: 3, Outcome: run.CheckFailed},
		{Type: run.TypeReviewQuestion, Step: "review", Check: "alpha", Cycle: 3, Finding: "f1", Leaf: "L1", Handle: "leaf-a", Question: "Q1", QuestionHandle: "leaf-failed-r1-1"},
		{Type: run.TypeCheckCycle, Step: "review", Cycle: 4, Checks: []string{"alpha"}, Kind: run.CycleAnswered, Questions: []string{"Q1"}, Decisions: []string{"rename-it"}},
	}}
	st := j.State()
	c, _ := st.Check("review", "alpha")
	rc := runCheckInfo(st, c)
	want := []RunCheckQuestion{{Cycle: 3, Finding: "f1", Leaf: "L1", LeafHandle: "leaf-a", Question: "Q1", QuestionHandle: "leaf-failed-r1-1"}}
	if !reflect.DeepEqual(rc.Questions, want) {
		t.Errorf("questions = %+v, want %+v", rc.Questions, want)
	}
	if len(rc.Cycles) != 2 || rc.Cycles[0].Decisions != nil || !reflect.DeepEqual(rc.Cycles[1].Decisions, []string{"rename-it"}) {
		t.Errorf("cycles = %+v", rc.Cycles)
	}
	var b bytes.Buffer
	RunInfo{ID: "R1", Checks: []RunCheck{rc}}.Human(&b)
	for _, s := range []string{"  - cycle 4: open, answered by rename-it\n", "  - question leaf-failed-r1-1 after cycle 3 on f1, blocks leaf leaf-a\n"} {
		if !strings.Contains(b.String(), s) {
			t.Errorf("status lacks %q:\n%s", s, b.String())
		}
	}
}
