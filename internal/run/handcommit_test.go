package run

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// handRunner is a fakeRunner whose HandCommits returns, on each scan, the
// next of scans (past the end, the last), or err.
type handRunner struct {
	*fakeRunner
	scans [][]Intervention
	err   error
	calls int
}

func (h *handRunner) HandCommits(st State) ([]Intervention, error) {
	h.calls++
	if h.err != nil {
		return nil, h.err
	}
	if len(h.scans) == 0 {
		return nil, nil
	}
	i := min(h.calls-1, len(h.scans)-1)
	return h.scans[i], nil
}

// TestHandCommitInterventions: the engine scans on entering each step and
// before the run finishes, journalling each hand commit once, at the step
// it scanned on; a scan that finds a journalled commit again adds none.
func TestHandCommitInterventions(t *testing.T) {
	e, r, _, _ := setup(t, cmdWF)
	r.codes["test"] = []int{0}
	a := Intervention{Commit: "aaa", Actor: "sam", Detail: "hand commit aaa: fix typo"}
	b := Intervention{Commit: "bbb", Actor: "kim", Detail: "hand commit bbb: tweak"}
	h := &handRunner{fakeRunner: r, scans: [][]Intervention{nil, {a}, {a, b}}}
	e.Runner = h
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded {
		t.Fatalf("result %q detail %q", st.Result, st.Detail)
	}
	// Enter test, enter ship, before the finish.
	if h.calls != 3 {
		t.Errorf("scans = %d, want 3", h.calls)
	}
	j, err := e.Store.Read(runID)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, iv := range j.State().Interventions {
		got = append(got, strings.Join([]string{iv.Kind, iv.Commit, iv.Actor, iv.Step}, " "))
	}
	want := []string{"hand-commit aaa sam ship", "hand-commit bbb kim ship"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("interventions = %q, want %q", got, want)
	}
	// A later scan of the same commits (a resume) adds none.
	h.scans = [][]Intervention{{a, b}}
	h.calls = 0
	r2 := &running{e: e, j: j}
	r2.scanHandCommits("ship")
	if j2, _ := e.Store.Read(runID); len(j2.State().Interventions) != 2 {
		t.Errorf("rescan journalled %d interventions, want 2", len(j2.State().Interventions))
	}
}

// TestHandCommitScanErrorChangesNothing (failure case): a scan that fails,
// or finds a commit with no actor, journals nothing and the run goes on.
func TestHandCommitScanErrorChangesNothing(t *testing.T) {
	for _, h := range []*handRunner{
		{err: errors.New("git log failed")},
		{scans: [][]Intervention{{{Commit: "ccc"}}}},
	} {
		e, r, _, _ := setup(t, cmdWF)
		r.codes["test"] = []int{0}
		h.fakeRunner = r
		e.Runner = h
		var logs []string
		e.Log = func(s string) { logs = append(logs, s) }
		st, err := e.Run(runID)
		if err != nil {
			t.Fatal(err)
		}
		if st.Result != Succeeded || len(st.Interventions) != 0 {
			t.Fatalf("result %q interventions %v", st.Result, st.Interventions)
		}
		if !strings.Contains(strings.Join(logs, "\n"), "hand commit") {
			t.Errorf("no hand commit line logged: %q", logs)
		}
	}
}
