package api

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/run"
)

// openCycleJournal is a check whose cycle 1 failed and whose cycle 2 is
// still open.
func openCycleJournal() run.Journal {
	return run.Journal{ID: "R1", Entries: []run.Entry{
		{Type: run.TypeStarted, Started: &run.Started{Workflow: "review", Change: "C1", Title: "Tidy up"}},
		{Type: run.TypeCheckCycle, Step: "review", Cycle: 1, Checks: []string{"alpha"}},
		{Type: run.TypeCheck, Step: "review", Check: "alpha", Cycle: 1, Outcome: run.CheckFailed, Detail: "cycle 1: confirmed f1"},
		{Type: run.TypeCheckCycle, Step: "review", Cycle: 2, Checks: []string{"alpha"}},
	}}
}

func TestOpenCycleReadsOpen(t *testing.T) {
	st := openCycleJournal().State()
	c, ok := st.Check("review", "alpha")
	if !ok {
		t.Fatal("no check alpha")
	}
	rc := runCheckInfo(st, c)
	if rc.Outcome != "" {
		t.Errorf("outcome = %q, want empty while cycle 2 is open", rc.Outcome)
	}
	var b bytes.Buffer
	RunInfo{ID: "R1", State: "running", Checks: []RunCheck{rc}}.Human(&b)
	if want := "- alpha (step review): open in cycle 2 (cycle 1 failed)\n"; !strings.Contains(b.String(), want) {
		t.Errorf("status text lacks %q:\n%s", want, b.String())
	}

	sum := reviewSummary(st)
	if !strings.Contains(sum, "### alpha: open\n") {
		t.Errorf("summary does not head alpha open:\n%s", sum)
	}

	// Once cycle 2 is judged, its outcome is the check's.
	j := openCycleJournal()
	j.Entries = append(j.Entries, run.Entry{Type: run.TypeCheck, Step: "review", Check: "alpha", Cycle: 2, Outcome: run.CheckPassed})
	st = j.State()
	c, _ = st.Check("review", "alpha")
	if rc := runCheckInfo(st, c); rc.Outcome != "passed" {
		t.Errorf("outcome = %q, want passed", rc.Outcome)
	}
}
