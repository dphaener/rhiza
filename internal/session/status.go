package session

import (
	"strings"
	"time"
)

// States a report can carry.
const (
	StateStarting = "starting"
	StateWorking  = "working"
	StateIdle     = "idle"
	StateBlocked  = "blocked"
	StateDone     = "done"
	StateError    = "error"
	StateExited   = "exited"
	// StateLost is the state of an exit a run wrote for a session it found
	// ended while no run watched it: the process's own exit never came.
	StateLost = "lost"
)

// Kinds of a blocked report.
const (
	KindPermission = "permission"
	KindQuestion   = "question"
	KindAuth       = "auth"
)

// Source names who made a report.
type Source string

const (
	SourceHook    Source = "hook"
	SourceRex     Source = "rex"
	SourceExit    Source = "exit"
	SourceDefault Source = "default"
)

// Report is the latest word from one source.
type Report struct {
	State  string    `json:"state"`
	Kind   string    `json:"kind,omitempty"`
	Detail string    `json:"detail,omitempty"`
	At     time.Time `json:"at"`
	// Code is the process exit code; only the exit source sets it.
	Code *int `json:"code,omitempty"`
}

// Status holds the latest report per source. A nil report means that source
// has not spoken. Each writer touches only its own field.
type Status struct {
	Hook *Report `json:"hook,omitempty"`
	Rex  *Report `json:"rex,omitempty"`
	Exit *Report `json:"exit,omitempty"`
	// Done is the agent's own word that its work is finished, written only
	// by `rhiza session done`. It is not a source Resolve weighs: it says
	// nothing of whether the session works, waits or has exited, and is
	// read on its own (SaidDone).
	Done *Report `json:"done,omitempty"`
	// Asked is a question the agent of an autonomous session put to the
	// owner, which its hook denied: the session waits for the arbiter's
	// answer, which arrives as a message. Later hooks do not touch it; it is
	// cleared once the answer is delivered. Like Done, Resolve does not
	// weigh it.
	Asked *Report `json:"asked,omitempty"`
}

// SaidDone reports whether the agent has said its work is done.
func (s Status) SaidDone() bool { return s.Done != nil }

// Exited reports whether the session's process has exited. A live source's
// report that arrives after that must not be recorded: it would not revive
// the session.
func (s Status) Exited() bool { return s.Exit != nil }

// HookEnded reports whether the session's hook said it ended (SessionEnd
// reports done), whether or not an exit was recorded.
func (s Status) HookEnded() bool { return s.Hook != nil && s.Hook.State == StateDone }

// ValidState reports whether state is one of the known states.
func ValidState(state string) bool {
	switch state {
	case StateStarting, StateWorking, StateIdle, StateBlocked, StateDone, StateError, StateExited:
		return true
	}
	return false
}

// Resolve returns the report that wins, and its source. Once the process
// has exited, the exit is the status: no hook or Rex report outranks it.
// While the session is live, hooks come first, then Rex events, then a bare
// starting report when no source has spoken. Precedence is strict: a present
// lower source never overrides a higher one, whatever the times say.
func Resolve(r Record) (Report, Source) {
	switch {
	case r.Status.Exit != nil:
		return *r.Status.Exit, SourceExit
	case r.Status.Hook != nil:
		return *r.Status.Hook, SourceHook
	case r.Status.Rex != nil:
		return *r.Status.Rex, SourceRex
	}
	return Report{State: StateStarting, At: r.StartedAt}, SourceDefault
}

// StallAfter is how long a working session may go without a report before
// it counts as stalled.
const StallAfter = 20 * time.Minute

// LimitPausedDetail is the detail of the working report a session stopped
// at its usage limit gets in place of idle.
const LimitPausedDetail = "paused at its usage limit"

// Stalled reports whether the session has been silent for more than
// StallAfter at now. Only a live session whose hook status resolves to
// working, or idle in a session an autonomous run launched, which has not
// said it is done and is not waiting on an answer, can stall: an idle
// prompt counts even when a background shell is open. An idle session of
// an owner-attended run never stalls, as it may be waiting on the owner.
// Blocked, exited and usage-paused sessions never do.
func Stalled(r Record, now time.Time) bool {
	if !Stallable(r) {
		return false
	}
	rep, _ := Resolve(r)
	return now.Sub(rep.At) > StallAfter
}

// Stallable reports whether the session could be stalled given enough
// silence: Stalled's conditions but for the time.
func Stallable(r Record) bool {
	rep, src := Resolve(r)
	if src != SourceHook {
		return false
	}
	if rep.State != StateWorking && (rep.State != StateIdle || !r.Autonomous) {
		return false
	}
	if r.Status.SaidDone() || r.Status.Asked != nil {
		return false
	}
	if strings.HasPrefix(rep.Detail, LimitPausedDetail) {
		return false
	}
	return true
}
