package run

import (
	"reflect"
	"testing"
)

// A code-check entry folds into State.CodeChecks, in journal order, and
// survives a new visit of its step.
func TestCodeCheckFold(t *testing.T) {
	j := Journal{ID: runID, Entries: []Entry{
		{Type: TypeEntered, Step: "code", At: t0},
		{Type: TypeCodeCheck, Step: "code", Outcome: CodeCheckFail, Detail: "leaf alpha: no Evidence verifies it", At: t0},
		{Type: TypeEntered, Step: "code", At: t0},
		{Type: TypeCodeCheck, Step: "code", Outcome: CodeCheckPass, At: t0},
	}}
	st := j.State()
	if len(st.CodeChecks) != 2 || st.CodeChecks[0].Outcome != CodeCheckFail || st.CodeChecks[0].Detail != "leaf alpha: no Evidence verifies it" {
		t.Fatalf("CodeChecks = %+v", st.CodeChecks)
	}
	c, ok := st.LatestCodeCheck()
	if !ok || c.Step != "code" || c.Outcome != CodeCheckPass {
		t.Errorf("LatestCodeCheck = %+v, %v", c, ok)
	}
	if _, ok := (Journal{ID: runID}).State().LatestCodeCheck(); ok {
		t.Error("a code check in an empty journal")
	}
}

// A code-check entry's lint hits survive the journal: appended, read back
// and folded, they come out of State().CodeChecks as written.
func TestCodeCheckLintsRoundTrip(t *testing.T) {
	s := newStore(t)
	if err := s.Create(runID, startedEntry("code", 0)); err != nil {
		t.Fatal(err)
	}
	lints := []LintHit{
		{Node: "01NODEA", Handle: "alpha", Field: "acceptance", Kind: "path", Ref: "internal/nope.go"},
		{Node: "01NODEB", Handle: "use-x", Kind: "unapplied-decision"},
		{Node: "01NODEC", Handle: "legacy", Field: "statement", Kind: "removed-code-capability", Ref: "legacy/old.txt"},
	}
	if err := s.Append(runID, Entry{Type: TypeCodeCheck, Step: "code", Outcome: CodeCheckPass, Lints: lints, At: at(1)}); err != nil {
		t.Fatal(err)
	}
	j, err := s.Read(runID)
	if err != nil {
		t.Fatal(err)
	}
	st := j.State()
	if len(st.CodeChecks) != 1 || st.CodeChecks[0].Outcome != CodeCheckPass || !reflect.DeepEqual(st.CodeChecks[0].Lints, lints) {
		t.Fatalf("CodeChecks = %+v", st.CodeChecks)
	}
}
