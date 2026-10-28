package run

import (
	"os"
	"os/exec"
	"testing"
)

// deadPid is the pid of a process that has exited.
func deadPid(t *testing.T) int {
	t.Helper()
	cmd := exec.Command("true")
	if err := cmd.Run(); err != nil {
		t.Skipf("no true to run: %v", err)
	}
	return cmd.Process.Pid
}

func processJournal(p *Process, finished bool) Journal {
	j := Journal{ID: "01RUNPROC", Entries: []Entry{startedEntry("a", 0)}}
	if p != nil {
		j.Entries = append(j.Entries, Entry{Type: TypeProcess, At: at(1), Process: p})
	}
	if finished {
		j.Entries = append(j.Entries, Entry{Type: TypeFinished, At: at(2), End: Succeeded})
	}
	return j
}

func TestLiveness(t *testing.T) {
	host, _ := os.Hostname()
	mine := &Process{Host: host, Pid: os.Getpid(), RexWindow: "window:w", RexBlock: "block:b"}
	gone := &Process{Host: host, Pid: deadPid(t)}
	cases := []struct {
		name string
		j    Journal
		want Liveness
	}{
		{"running", processJournal(mine, false), LiveRunning},
		{"interrupted", processJournal(gone, false), LiveInterrupted},
		{"finished", processJournal(mine, true), NotLive},
		{"finished, process gone", processJournal(gone, true), NotLive},
		{"no record", processJournal(nil, false), LiveInterrupted},
		{"other host", processJournal(&Process{Host: host + "-x", Pid: gone.Pid}, false), LiveRunning},
	}
	for _, c := range cases {
		st := c.j.State()
		if got := st.Liveness(host, nil); got != c.want {
			t.Errorf("%s: %s, want %s", c.name, got, c.want)
		}
		if st.Live() != (c.want != NotLive) {
			t.Errorf("%s: Live = %v", c.name, st.Live())
		}
	}
}

// The latest process entry wins, and the Rex window and block fold.
func TestStateProcessLatestWins(t *testing.T) {
	host, _ := os.Hostname()
	j := processJournal(&Process{Host: host, Pid: deadPid(t), RexWindow: "window:1", RexBlock: "block:1"}, false)
	now := &Process{Host: host, Pid: os.Getpid(), RexWindow: "window:2", RexBlock: "block:2"}
	j.Entries = append(j.Entries, Entry{Type: TypeProcess, At: at(3), Process: now})
	st := j.State()
	if st.Process == nil || *st.Process != *now {
		t.Fatalf("Process = %+v", st.Process)
	}
	if st.Liveness(host, nil) != LiveRunning {
		t.Error("not running")
	}
}
