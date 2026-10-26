package session

import "time"

// States a report can carry.
const (
	StateStarting = "starting"
	StateWorking  = "working"
	StateIdle     = "idle"
	StateBlocked  = "blocked"
	StateDone     = "done"
	StateError    = "error"
	StateExited   = "exited"
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
}

// Exited reports whether the session's process has exited. A live source's
// report that arrives after that must not be recorded: it would not revive
// the session.
func (s Status) Exited() bool { return s.Exit != nil }

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
