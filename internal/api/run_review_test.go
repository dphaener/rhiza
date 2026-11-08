package api

import (
	"testing"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/store"
)

// A check or critic session's RHIZA_CYCLE names its round's cycle when it
// is one of the step's cycles that runs the check; anything else is 0, so
// the cycle is inferred.
func TestEnvCheckCycle(t *testing.T) {
	st := run.State{Cycles: []run.CheckCycle{
		{Step: "review", Cycle: 1, Checks: []string{"alpha", "beta"}},
		{Step: "review", Cycle: 2, Checks: []string{"beta"}},
	}}
	for _, tt := range []struct {
		env, check string
		want       int
	}{
		{"1", "alpha", 1},
		{" 2 ", "beta", 2},
		{"2", "alpha", 0},
		{"3", "beta", 0},
		{"0", "alpha", 0},
		{"x", "alpha", 0},
		{"", "alpha", 0},
	} {
		ctx := &Ctx{Getenv: func(k string) string {
			if k == envCycle {
				return tt.env
			}
			return ""
		}}
		if got := envCheckCycle(ctx, st, "review", tt.check); got != tt.want {
			t.Errorf("RHIZA_CYCLE=%q for %s: %d, want %d", tt.env, tt.check, got, tt.want)
		}
	}
}

// Two run IDs that share their first 8 characters (the same start time)
// get different question handles carrying the full IDs.
func TestReviewQuestionHandleFullID(t *testing.T) {
	a := "01M6VY7E0AAAAAAAAAAAAAAAAA"
	b := "01M6VY7E0BBBBBBBBBBBBBBBBB"
	ha, hb := reviewQuestionHandle(a, 3, "f1"), reviewQuestionHandle(b, 3, "f1")
	if ha == hb {
		t.Fatalf("handles collide: %s", ha)
	}
	if want := "review-01m6vy7e0aaaaaaaaaaaaaaaaa-r3-f1"; ha != want {
		t.Errorf("handle = %s, want %s", ha, want)
	}
}

// The asked entry for a finding is found by its journalled handle, the base
// or base-N, and the autonomous gate takes it for the question left.
func TestReviewAskedFor(t *testing.T) {
	base := reviewQuestionHandle("01M6VY7E0AAAAAAAAAAAAAAAAA", 3, "f1")
	other := reviewQuestionHandle("01M6VY7E0BBBBBBBBBBBBBBBBB", 3, "f1")
	asked := []run.Asked{{Handle: other, Question: "x"}, {Handle: base + "-2", Question: "q2"}}
	a, ok := reviewAskedFor(asked, base)
	if !ok || a.Question != "q2" {
		t.Fatalf("reviewAskedFor = %v, %v; want the -2 entry", a, ok)
	}
	if _, ok := reviewAskedFor(asked[:1], base); ok {
		t.Error("another run's entry matched")
	}
	if _, ok := reviewAskedFor([]run.Asked{{Handle: base + "-x"}}, base); ok {
		t.Error("a non-numeric suffix matched")
	}
	st := run.State{ID: "01M6VY7E0AAAAAAAAAAAAAAAAA", Asked: []run.Asked{{Handle: base + "-2", Question: "01M6W7QGXS2VNCJ3KMFRMW78XE"}},
		Rounds: []run.Round{{Number: 3, Findings: []run.Finding{{ID: "f1", Node: "n"}}, Verdicts: []run.Verdict{{Finding: "f1", Verdict: run.Confirmed}}}}}
	m := newModel(store.Log{}, id.ID{})
	open, err := gateReview(nil, m, reviewRun{id: st.ID, state: st})
	if err != nil {
		t.Fatalf("gateReview with the -2 question: %v", err)
	}
	if !open {
		t.Error("gateReview: the unanswered -2 question is not open")
	}
	st.Asked = []run.Asked{{Handle: other, Question: "01M6W7QGXS2VNCJ3KMFRMW78XE"}}
	if _, err := gateReview(nil, m, reviewRun{id: st.ID, state: st}); err == nil {
		t.Error("gateReview took another run's question for the finding")
	}
}
