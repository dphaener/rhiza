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
	"time"
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
	// AskReply reports whether one hook event and its stdin payload is the
	// agent asking the owner a question, and if so returns what the hook
	// prints on stdout to deny that call, telling the agent reason. It is
	// pure: it builds bytes only.
	AskReply(event string, payload []byte, reason string) (reply []byte, ok bool, err error)
	// HookMessage reads the message a session sent another out of one hook
	// event and its stdin payload. The bool is false when the event carries
	// no message. It errors on a payload it cannot read, or on a message
	// call with no recipient or no message.
	HookMessage(event string, payload []byte) (Message, bool, error)
	// ContextReply returns what a hook prints on stdout for the agent to
	// add text to its context at its next step, when event is a tool-call
	// boundary that can carry it; the bool is false when it cannot. It is
	// pure: it builds bytes only.
	ContextReply(event, context string) ([]byte, bool, error)
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
	// TokenUse returns the session's token use read from transcript. Each
	// of the agent's messages counts once, however many transcript lines
	// repeat it; lines its helper agents left in the same transcript count
	// toward Helpers. Unreadable lines are skipped, not fatal. It reads
	// nothing itself.
	TokenUse(transcript []byte) (TokenUse, error)
	// UsageLimit reports whether transcript ends with the agent stopped by
	// its usage limit, and when the limit resets if the agent said. It
	// reads nothing itself.
	UsageLimit(transcript []byte) (UsageLimit, bool)
}

// UsageLimit is a session stopped by the agent's usage limit.
type UsageLimit struct {
	// ResetAt is when the limit resets; zero when the agent did not say.
	ResetAt time.Time
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
	// StatusLineCommand is the argv the agent runs to render its status
	// line, with the status line JSON on stdin; empty installs none.
	StatusLineCommand []string
	// Prompt is the agent's first prompt; empty starts it with none.
	Prompt string
	// Resume continues the session SessionID instead of starting it; Prompt
	// is then what the resumed session is told.
	Resume bool
	// Model is the model the agent starts on; "" means the agent's
	// default.
	Model string
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

// Message is one message a session sent another through the agent.
type Message struct {
	// SessionID is the agent's own session ID from the payload.
	SessionID string
	// To is the recipient, verbatim, as the agent named it.
	To string
	// Text is the message, verbatim; its raw JSON when Structured.
	Text string
	// Summary is the agent's short summary of the message; "" when none.
	Summary string
	// Structured reports whether the message was not a string, so Text
	// holds its raw JSON.
	Structured bool
}

// TokenCounts are the tokens of some of a session's model calls.
type TokenCounts struct {
	// Input is the tokens sent that were neither written to nor read from
	// the cache.
	Input int64
	// Output is the tokens the model generated.
	Output int64
	// CacheCreation is the input tokens written to the cache.
	CacheCreation int64
	// CacheRead is the input tokens read from the cache.
	CacheRead int64
}

// TokenUse is a session's token use: its own calls, and those its helper
// agents logged in its transcript.
type TokenUse struct {
	// Session is the tokens of the session's own model calls.
	Session TokenCounts
	// Helpers is the tokens of the calls its helper agents logged in the
	// session's transcript, kept apart from Session.
	Helpers TokenCounts
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

// PaneLimiter is the optional part of an adapter that reads the text of the
// agent's terminal pane: whether it shows the agent parked on a usage-limit
// prompt that a key press dismisses. It is pure.
type PaneLimiter interface {
	PaneUsageLimit(pane string) bool
}

// Truster is implemented by an adapter whose agent stops at a prompt in a
// folder its owner has not trusted yet. Callers that start the agent mark
// the folders trusted first. Like the rest of an adapter it is pure: the
// caller reads and writes the file.
type Truster interface {
	// TrustFile returns the path of the file that records trusted folders,
	// or "" when getenv gives no place to look for it.
	TrustFile(getenv func(string) string) string
	// Trust returns current (the file's bytes, nil when it is missing) with
	// every folder in dirs trusted and every other setting unchanged.
	Trust(current []byte, dirs []string) ([]byte, error)
}

// UsageRecord is one agent message's counted tokens and when it was logged.
type UsageRecord struct {
	// At is the message's timestamp.
	At time.Time
	// Tokens is the message's input, output and cache-creation tokens;
	// cache reads are never counted.
	Tokens int64
}

// UsageReader is the optional part of an adapter that reports the agent's
// token use on the whole machine, the owner's own sessions included. Unlike
// the rest of an adapter it reads files, and only those the agent keeps under
// its own config directory.
type UsageReader interface {
	// UsageRecords returns the token records of every transcript the agent
	// keeps on the machine, helpers' included, that are timestamped in
	// [from, to), in time order. Each message counts once, by its id.
	// getenv looks up an environment variable. Unreadable files and lines
	// are skipped, not fatal.
	UsageRecords(getenv func(string) string, from, to time.Time) ([]UsageRecord, error)
}

// LiveWindow is one rate-limit window an agent reports about the account's
// live use.
type LiveWindow struct {
	// Kind is the window kind: "session" or "weekly".
	Kind string
	// UsedPercentage is the share of the window's limit used, 0 to 100.
	UsedPercentage float64
	// ResetsAt is when the window resets.
	ResetsAt time.Time
}

// StatusLineReader is the optional part of an adapter whose agent renders a
// status line from a JSON payload on stdin that carries the account's live
// rate limits.
type StatusLineReader interface {
	// LiveUsage returns the rate-limit windows the payload carries; none
	// when it carries none. A payload that does not parse is an error.
	LiveUsage(payload []byte) ([]LiveWindow, error)
	// UserStatusLine returns the status line command the owner's own agent
	// settings configure, "" when there is none. getenv looks up an
	// environment variable.
	UserStatusLine(getenv func(string) string) string
}

// LiveSnapshot is the latest live figure recorded for one rate-limit window.
type LiveSnapshot struct {
	// Kind is the window kind: "session" or "weekly".
	Kind string
	// UsedPercentage is the share of the window's limit used, 0 to 100.
	UsedPercentage float64
	// ResetsAt is when the window resets.
	ResetsAt time.Time
	// CapturedAt is when the figure was recorded.
	CapturedAt time.Time
}

// LiveUsageReader is the optional part of an adapter that reports the live
// rate-limit figures its agent's status lines recorded. An adapter without
// it reports none.
type LiveUsageReader interface {
	// LiveSnapshots returns the latest recorded snapshot of each window kind
	// in the usage state folder dir, ordered by kind. A missing or corrupt
	// record gives none, never an error.
	LiveSnapshots(dir string) []LiveSnapshot
}
