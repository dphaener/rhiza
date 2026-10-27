package run

import (
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/dphaener/rhiza/internal/workflow"
)

// paneStates builds a step's leaf states from sessions given as
// "leaf role attempt" with their judgement and pane, in start order.
type paneSession struct {
	leaf, role string
	attempt    int
	outcome    string // "" while live
	closed     bool
}

func paneStates(ss []paneSession) map[string]LeafState {
	out := map[string]LeafState{}
	for i, s := range ss {
		l := out[s.leaf]
		l.Leaf, l.Handle, l.Step = s.leaf, s.leaf, "build"
		l.Sessions = append(l.Sessions, LeafSession{Role: s.role, Attempt: s.attempt,
			Session: leafSID(s.leaf, s.role, s.attempt), Seq: i, PaneClosed: s.closed})
		if s.outcome != "" {
			l.Attempts = append(l.Attempts, LeafAttempt{Role: s.role, Attempt: s.attempt, Outcome: s.outcome})
		}
		out[s.leaf] = l
	}
	return out
}

// leafPane anchors a leaf session's pane from the journal alone: the
// first builder opens the tab, the first verifier splits its builder, a
// later verifier the most recent live verifier, and the next builder
// takes a passed builder's pane but only splits a failed one's.
func TestLeafPaneAnchor(t *testing.T) {
	pass, fail := workflow.OutcomePass, workflow.OutcomeFail
	b, v := workflow.RoleBuilder, workflow.RoleVerifier
	for _, c := range []struct {
		name             string
		ss               []paneSession
		act              leafAction
		anchor, replaces string
	}{
		{"first builder opens the tab", nil,
			leafAction{leaf: Leaf{ID: "a"}, role: b, attempt: 1}, "", ""},
		{"first verifier splits its builder", []paneSession{{"a", b, 1, pass, false}},
			leafAction{leaf: Leaf{ID: "a"}, role: v, attempt: 1}, leafSID("a", b, 1), ""},
		{"next builder takes the passed builder's pane", []paneSession{{"a", b, 1, pass, false}, {"a", v, 1, "", false}},
			leafAction{leaf: Leaf{ID: "b"}, role: b, attempt: 1}, leafSID("a", b, 1), leafSID("a", b, 1)},
		{"later verifier splits the live verifier", []paneSession{
			{"a", b, 1, pass, true}, {"a", v, 1, "", false}, {"b", b, 1, pass, false}},
			leafAction{leaf: Leaf{ID: "b"}, role: v, attempt: 1}, leafSID("a", v, 1), ""},
		{"third verifier splits the most recent verifier", []paneSession{
			{"a", b, 1, pass, true}, {"a", v, 1, "", false}, {"b", b, 1, pass, true}, {"b", v, 1, "", false}, {"c", b, 1, pass, false}},
			leafAction{leaf: Leaf{ID: "c"}, role: v, attempt: 1}, leafSID("b", v, 1), ""},
		{"a passed verifier's closed pane is no anchor", []paneSession{
			{"a", b, 1, pass, true}, {"a", v, 1, pass, true}, {"b", b, 1, pass, false}},
			leafAction{leaf: Leaf{ID: "b"}, role: v, attempt: 1}, leafSID("b", b, 1), ""},
		{"a retry builder splits a failed builder and keeps it", []paneSession{{"a", b, 1, fail, false}},
			leafAction{leaf: Leaf{ID: "a"}, role: b, attempt: 2}, leafSID("a", b, 1), ""},
		{"a builder with only a verifier pane left splits it", []paneSession{{"a", b, 1, pass, true}, {"a", v, 1, fail, false}},
			leafAction{leaf: Leaf{ID: "a"}, role: b, attempt: 2}, leafSID("a", v, 1), ""},
		{"every pane closed opens the tab anew", []paneSession{{"a", b, 1, pass, true}, {"a", v, 1, pass, true}},
			leafAction{leaf: Leaf{ID: "b"}, role: b, attempt: 1}, "", ""},
	} {
		anchor, replaces := leafPane(paneStates(c.ss), c.act)
		if anchor != c.anchor || replaces != c.replaces {
			t.Errorf("%s: leafPane = %q, %q, want %q, %q", c.name, anchor, replaces, c.anchor, c.replaces)
		}
	}
}

// A leaves step's panes over a run: each session's spec names its anchor,
// a passed verifier's pane closes as it passes, a passed builder's when
// the next builder takes its place, and a failed session's stays until
// the step ends.
func TestLeavesStepPanes(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = []Leaf{{ID: "LA", Handle: "a"}, {ID: "LB", Handle: "b"}}
	av, bv := leafSID("a", "verifier", 1), leafSID("b", "verifier", 1)
	ab, bb := leafSID("a", "builder", 1), leafSID("b", "builder", 1)
	r.outcome[ab], r.outcome[av], r.outcome[bb] = workflow.OutcomePass, workflow.OutcomePass, workflow.OutcomePass
	r.outcome[bv] = workflow.OutcomeFail
	// A's verifier is still live when B's verifier starts.
	r.agentOut[av] = []bool{false, false, false, false, false, false}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:fail"}) {
		t.Fatalf("outcomes = %v", got)
	}
	anchors := map[string]string{}
	for _, s := range r.leafSpecs {
		anchors[leafSID(s.Leaf.Handle, s.Role, s.Attempt)] = s.Anchor
	}
	b2, b3 := leafSID("b", "builder", 2), leafSID("b", "builder", 3)
	want := map[string]string{ab: "", av: ab, bb: ab, bv: av, b2: bb, b3: b2}
	if !reflect.DeepEqual(anchors, want) {
		t.Errorf("anchors = %v\nwant %v", anchors, want)
	}

	j, err := e.Store.Read(runID)
	if err != nil {
		t.Fatal(err)
	}
	var seq []string
	for _, en := range j.Entries {
		switch en.Type {
		case TypePaneClosed:
			seq = append(seq, "closed "+en.Session)
		case TypeLeafSession:
			seq = append(seq, "start "+en.Session)
		case TypeLeafAttempt:
			seq = append(seq, "judged "+en.Role+" "+en.Handle+" "+attemptStr(en.Attempt)+" "+en.Outcome)
		case TypeLeaf:
			seq = append(seq, "leaf "+en.Handle+" "+en.Outcome)
		case TypeEnded:
			seq = append(seq, "ended")
		}
	}
	at := func(s string) int {
		i := slices.Index(seq, s)
		if i < 0 {
			t.Fatalf("%q not in\n%v", s, seq)
		}
		return i
	}
	if at("closed "+av) != at("judged verifier a 1 pass")+1 {
		t.Errorf("A's verifier's pane did not close as it passed:\n%v", seq)
	}
	if at("closed "+ab) != at("start "+bb)+1 {
		t.Errorf("A's builder's pane did not close as B's builder took it:\n%v", seq)
	}
	if at("closed "+bb) != at("start "+b2)+1 {
		t.Errorf("B's passed builder's pane did not close as its retry took it:\n%v", seq)
	}
	failedAt := at("leaf b failed")
	for _, s := range []string{bv, b2, b3} {
		if i := at("closed " + s); i < failedAt || i > at("ended") {
			t.Errorf("failed session %s's pane closed before the step ended:\n%v", s, seq)
		}
	}
	if !reflect.DeepEqual(r.closed, []string{ab, bb, av, bv, b2, b3}) && !reflect.DeepEqual(r.closed, []string{ab, av, bb, bv, b2, b3}) {
		t.Errorf("closed = %v", r.closed)
	}
	for _, l := range st.Leaves {
		for _, s := range l.Sessions {
			if !s.PaneClosed {
				t.Errorf("session %s's pane is open in the fold after the step", s.Session)
			}
		}
	}
}

// A pane that cannot be closed is only logged, and stays open in the
// journal, so the step's end tries it again; the run carries on.
func TestLeavesStepPaneCloseFails(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = []Leaf{{ID: "LA", Handle: "a"}}
	passAll(r, r.leaves)
	r.closeErr = errors.New("rex is down")
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"build:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	for _, s := range st.Leaves[0].Sessions {
		if s.PaneClosed {
			t.Errorf("session %s's pane closed in the journal though ClosePane failed", s.Session)
		}
	}
}

// A live session's pane stays when the step ends on an error: a resume
// rejoins it.
func TestLeavesStepLivePaneStays(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = []Leaf{{ID: "LA", Handle: "a"}}
	ab := leafSID("a", "builder", 1)
	r.agentOut[ab] = []bool{false}
	r.agentErr[ab] = errors.New("rex is down")
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	if len(r.closed) != 0 {
		t.Errorf("closed = %v, want none", r.closed)
	}
}
