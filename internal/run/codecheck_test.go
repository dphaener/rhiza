package run

import "testing"

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
