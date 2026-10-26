// Package adapter is the seam between Rhiza and the agent programs it
// launches. Callers see only the Adapter interface; the one package per agent
// that implements it is the only place that names the agent.
//
// An adapter builds argv, interprets hook payloads, says where the agent keeps
// a session's transcript and reads the owner's answers out of one. It never
// runs anything and never reads files or the environment itself.
package adapter

import (
	"fmt"
	"slices"
	"sort"
	"sync"
)

// Adapter describes how to launch one kind of agent and how to read its hooks.
type Adapter interface {
	// Name is the registry key, for example the value of a --agent flag.
	Name() string
	// PermissionModes lists every permission mode the agent accepts.
	PermissionModes() []string
	// Bypass reports whether mode switches the agent's permission checks off.
	Bypass(mode string) bool
	// Command returns the argv that starts the agent for l. It errors when
	// l.PermissionMode is empty or not one of PermissionModes.
	Command(l Launch) ([]string, error)
	// HookStatus maps one hook event and its stdin payload to a status
	// report. The bool is false when the event leaves the status unchanged.
	HookStatus(event string, payload []byte) (HookReport, bool, error)
	// TranscriptGlobs returns file glob patterns, most specific first, where
	// the agent keeps the session's own transcript (not those of any helper
	// agents it started). getenv looks up an environment variable. It is
	// pure: it reads nothing, and returns no patterns for a session ID that
	// could escape or widen a pattern.
	TranscriptGlobs(sessionID string, getenv func(string) string) []string
	// Answers returns every question in transcript that the agent asked and
	// the owner answered, in transcript order; the questions of one call keep
	// their order. Unreadable entries are skipped, not fatal.
	Answers(transcript []byte) ([]Answer, error)
}

// Answer is one question the agent asked the owner and the owner's answer,
// split by who wrote each part: option labels are the agent's wording that
// the owner picked; Typed and Notes are the owner's own words.
type Answer struct {
	// ToolUseID identifies the agent's call that asked the question, so one
	// answer can be told apart from another with the same text.
	ToolUseID string
	// Question is the question text, verbatim.
	Question string
	// Header is the question's short label, if the agent gave one.
	Header string
	// Options are the labels of the options the agent offered, in order.
	Options []string
	// MultiSelect reports whether the owner could pick several options.
	MultiSelect bool
	// Picked are the option labels the owner chose, verbatim, in the order
	// they appear in the answer. Empty when the owner picked none.
	Picked []string
	// Typed is the text the owner typed that is not an option label,
	// verbatim; "" when the owner typed nothing.
	Typed string
	// Notes is the annotation the owner typed alongside the answer, verbatim;
	// "" when there is none.
	Notes string
}

// Launch is what an adapter needs to build a command.
type Launch struct {
	// SessionID is the Rhiza-generated ID the agent must use for its session.
	SessionID string
	// PermissionMode is mandatory and must be one of PermissionModes.
	PermissionMode string
	// HookCommand is the argv prefix the agent runs for each hook; the adapter
	// appends the event name.
	HookCommand []string
}

// HookReport is the status a hook event reports.
type HookReport struct {
	// SessionID is the agent's own session ID from the payload, so callers
	// can verify it is the session they launched.
	SessionID string
	// State is one of idle, working, blocked or done.
	State string
	// Kind qualifies a blocked state, for example permission.
	Kind string
	// Detail is free text from the agent, if any.
	Detail string
}

// CheckMode returns an error unless mode is a non-empty member of
// a.PermissionModes. Adapters call it from Command.
func CheckMode(a Adapter, mode string) error {
	if mode == "" {
		return fmt.Errorf("a permission mode is required; one of %v", a.PermissionModes())
	}
	if !slices.Contains(a.PermissionModes(), mode) {
		return fmt.Errorf("unknown permission mode %q; one of %v", mode, a.PermissionModes())
	}
	return nil
}

var (
	mu       sync.RWMutex
	adapters = map[string]Adapter{}
	def      string
)

// Register adds a to the registry. It panics on an empty or duplicate name,
// like database/sql.Register.
func Register(a Adapter) {
	mu.Lock()
	defer mu.Unlock()
	register(a)
}

// RegisterDefault adds a to the registry and makes it the default adapter.
func RegisterDefault(a Adapter) {
	mu.Lock()
	defer mu.Unlock()
	register(a)
	def = a.Name()
}

func register(a Adapter) {
	name := a.Name()
	if name == "" {
		panic("adapter: Register with an empty name")
	}
	if _, dup := adapters[name]; dup {
		panic("adapter: Register called twice for " + name)
	}
	adapters[name] = a
}

// Lookup returns the adapter registered under name.
func Lookup(name string) (Adapter, bool) {
	mu.RLock()
	defer mu.RUnlock()
	a, ok := adapters[name]
	return a, ok
}

// Default returns the default adapter, or nil when none is registered as one.
func Default() Adapter {
	mu.RLock()
	defer mu.RUnlock()
	return adapters[def]
}

// Names returns the registered adapter names, sorted.
func Names() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(adapters))
	for n := range adapters {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
