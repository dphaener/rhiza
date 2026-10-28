package daemon

import (
	"testing"

	"github.com/dphaener/rhiza/internal/run"
)

func TestResumable(t *testing.T) {
	dead := func(int) bool { return false }
	proc := run.Entry{Type: run.TypeProcess, Process: &run.Process{Host: "h", Pid: 1}}
	journal := func(id string, es ...run.Entry) run.Journal {
		return run.Journal{ID: id, Entries: append([]run.Entry{{Type: run.TypeStarted, Started: &run.Started{Change: id}}}, es...)}
	}
	runs := []run.Journal{
		journal("interrupted", proc),
		journal("again", proc, run.Entry{Type: run.TypeDaemonResumed}),
		journal("failed", proc, run.Entry{Type: run.TypeFinished, End: run.Failed}),
	}
	got := Resumable(runs, "h", dead)
	if len(got) != 1 || got[0] != "interrupted" {
		t.Fatalf("Resumable = %v, want [interrupted]", got)
	}
}
