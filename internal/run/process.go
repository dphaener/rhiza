package run

import (
	"os"
	"strconv"
)

// Process is the process that runs a run: where it is and, when it runs in
// Rex, the session, window and block of its tab. Host and Pid are always
// set; the Rex fields only for a run in Rex. RexWindow may be empty when
// Rex could not name the block's window, and RexSession when the process
// had no REX_SESSION. Rex finds a block only within its session, so a
// caller in another session needs RexSession to reach the tab.
type Process struct {
	Host       string `json:"host"`
	Pid        int    `json:"pid"`
	RexSession string `json:"rex_session,omitempty"`
	RexWindow  string `json:"rex_window,omitempty"`
	RexBlock   string `json:"rex_block,omitempty"`
}

// Liveness is where a run stands as a process.
type Liveness string

const (
	// NotLive: the journal finished (done or failed), so the run counts
	// for nothing, as run-one-live has it.
	NotLive Liveness = "not-live"
	// LiveRunning: the journal has not finished and its process is running.
	LiveRunning Liveness = "running"
	// LiveInterrupted: the journal has not finished and its process is
	// gone, or was never recorded; the run can be resumed.
	LiveInterrupted Liveness = "interrupted"
)

// Live reports whether the run counts as live: its journal has not
// finished. It is the test LiveOf uses.
func (st State) Live() bool { return !st.Ended }

// Liveness reports whether st is not live, live and running, or live and
// interrupted. host is the machine asking: a process on another host
// cannot be looked up, so it counts as running. alive reports whether a
// pid runs on this host; nil means PidAlive. A run with no process record
// is interrupted: nothing says it runs (resuming it is still refused while
// another process holds its lock).
func (st State) Liveness(host string, alive func(pid int) bool) Liveness {
	if st.Ended {
		return NotLive
	}
	p := st.Process
	if p == nil || p.Pid <= 0 {
		return LiveInterrupted
	}
	if p.Host != host {
		return LiveRunning
	}
	if alive == nil {
		alive = PidAlive
	}
	if alive(p.Pid) {
		return LiveRunning
	}
	return LiveInterrupted
}

// ThisProcess is the calling process: its host name and pid, with the Rex
// session, window and block it runs in as given.
func ThisProcess(session, window, block string) Process {
	host, _ := os.Hostname()
	return Process{Host: host, Pid: os.Getpid(), RexSession: session, RexWindow: window, RexBlock: block}
}

// String is the process as "host:pid".
func (p Process) String() string { return p.Host + ":" + strconv.Itoa(p.Pid) }
