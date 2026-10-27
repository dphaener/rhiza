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
	ID             string    `json:"id"`
	Agent          string    `json:"agent"`
	Change         string    `json:"change"`
	Title          string    `json:"title"`
	Step           string    `json:"step"`
	Run            string    `json:"run,omitempty"`
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
