// Package session holds Rhiza's record of an agent session: what was
// launched, where, and the latest status each source has reported. It is pure
// data and file handling: it runs no program, talks to no terminal and names
// no agent.
package session

import "time"

// Tab identifies the terminal tab a session runs in.
type Tab struct {
	Session string `json:"session"`
	Window  string `json:"window"`
	Block   string `json:"block"`
	Title   string `json:"title"`
}

// Record is one session, stored as <id>.json.
type Record struct {
	ID     string `json:"id"`
	Agent  string `json:"agent"`
	Change string `json:"change"`
	// Writes is the change set the agent writes into (its RHIZA_CHANGE)
	// when that is not Change: an arbiter answers a question in the change
	// set the question lives in. WritesRecord marks Writes as a record
	// change set opened for the agent, landed once its work passes.
	Writes       string `json:"writes,omitempty"`
	WritesRecord bool   `json:"writes_record,omitempty"`
	// Advises is the session an advisor recommends an answer for: the
	// session of a normal run that is blocked on a question. Empty for
	// every other session.
	Advises string `json:"advises,omitempty"`
	Title   string `json:"title"`
	Step    string `json:"step"`
	Run     string `json:"run,omitempty"`
	// Autonomous is true when an autonomous run launched the session: the
	// owner is not there to answer it.
	Autonomous bool `json:"autonomous,omitempty"`
	// Leaf is the handle of the leaf a leaves step launched the session
	// on, Attempt the leaf's attempt and Role the session's role (builder
	// or verifier); empty for every other session. Commit is the commit a
	// verifier checks.
	Leaf    string `json:"leaf,omitempty"`
	Attempt int    `json:"attempt,omitempty"`
	Role    string `json:"role,omitempty"`
	Commit  string `json:"commit,omitempty"`
	// Check is the name of the check a checks step launched the session
	// for, and Cycle the check's cycle; Role is then check. Empty for
	// every other session.
	Check string `json:"check,omitempty"`
	Cycle int    `json:"cycle,omitempty"`
	// Model is the model the agent starts on; "" means the agent's
	// default.
	Model string `json:"model,omitempty"`
	// Dir is where the agent runs when that is not Worktree: a verifier's
	// throwaway checkout of the commit it checks. Root is then the Rhiza
	// root its rhiza writes go to (its RHIZA_ROOT), in Worktree.
	Dir            string    `json:"dir,omitempty"`
	Root           string    `json:"root,omitempty"`
	Actor          string    `json:"actor"`
	For            string    `json:"for"`
	PermissionMode string    `json:"permission_mode"`
	Prompt         string    `json:"prompt,omitempty"`
	Worktree       string    `json:"worktree"`
	Branch         string    `json:"branch"`
	Tab            Tab       `json:"tab"`
	StartedAt      time.Time `json:"started_at"`
	Status         Status    `json:"status"`
	// Tokens is the latest reading of the session's token use, nil until
	// one is recorded.
	Tokens *TokenUse `json:"tokens,omitempty"`
}

// TokenCounts are the tokens a set of model calls used, by kind.
type TokenCounts struct {
	Input         int64 `json:"input"`
	Output        int64 `json:"output"`
	CacheCreation int64 `json:"cache_creation"`
	CacheRead     int64 `json:"cache_read"`
}

// TokenUse is a reading of a session's token use: its own calls, and
// apart from them the calls its helper agents logged in its transcript.
type TokenUse struct {
	Session TokenCounts `json:"session"`
	Helpers TokenCounts `json:"helpers"`
}
