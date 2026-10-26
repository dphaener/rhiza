// Package claudecode is the adapter for Claude Code, the first built-in
// adapter and the default. It is the only package that names the agent.
package claudecode

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dphaener/rhiza/internal/adapter"
)

// Name is the registry name of this adapter.
const Name = "claude-code"

// hookEvents are the events whose hooks the settings JSON installs.
var hookEvents = []string{"SessionStart", "UserPromptSubmit", "PreToolUse", "Notification", "Stop", "SessionEnd"}

var modes = []string{"acceptEdits", "auto", "bypassPermissions", "manual", "dontAsk", "plan"}

func init() { adapter.RegisterDefault(Adapter{}) }

// Adapter launches Claude Code.
type Adapter struct{}

// Name implements adapter.Adapter.
func (Adapter) Name() string { return Name }

// PermissionModes implements adapter.Adapter.
func (Adapter) PermissionModes() []string { return append([]string(nil), modes...) }

// Bypass implements adapter.Adapter.
func (Adapter) Bypass(mode string) bool { return mode == "bypassPermissions" }

// Command implements adapter.Adapter. The permission mode is always an
// argument and never part of the settings JSON.
func (a Adapter) Command(l adapter.Launch) ([]string, error) {
	if err := adapter.CheckMode(a, l.PermissionMode); err != nil {
		return nil, err
	}
	if l.SessionID == "" {
		return nil, fmt.Errorf("a session ID is required")
	}
	settings, err := Settings(l.HookCommand)
	if err != nil {
		return nil, err
	}
	return []string{"claude", "--session-id", l.SessionID, "--permission-mode", l.PermissionMode, "--settings", settings}, nil
}

type hookCmd struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

type hookGroup struct {
	Hooks []hookCmd `json:"hooks"`
}

type settings struct {
	Hooks map[string][]hookGroup `json:"hooks"`
}

// Settings returns the settings JSON: command hooks only, never permissions.
func Settings(hookCommand []string) (string, error) {
	if len(hookCommand) == 0 {
		return "", fmt.Errorf("a hook command is required")
	}
	s := settings{Hooks: map[string][]hookGroup{}}
	for _, ev := range hookEvents {
		cmd := shellJoin(append(append([]string(nil), hookCommand...), ev))
		s.Hooks[ev] = []hookGroup{{Hooks: []hookCmd{{Type: "command", Command: cmd}}}}
	}
	b, err := json.Marshal(s)
	return string(b), err
}

func shellJoin(args []string) string {
	out := make([]string, len(args))
	for i, a := range args {
		out[i] = shellQuote(a)
	}
	return strings.Join(out, " ")
}

func shellQuote(s string) string {
	if s != "" && strings.IndexFunc(s, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("-_./:=@%+,", r))
	}) < 0 {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

type payload struct {
	SessionID        string `json:"session_id"`
	NotificationType string `json:"notification_type"`
	Message          string `json:"message"`
}

// HookStatus implements adapter.Adapter.
func (Adapter) HookStatus(event string, raw []byte) (adapter.HookReport, bool, error) {
	var p payload
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &p); err != nil {
			return adapter.HookReport{}, false, fmt.Errorf("%s hook payload: %w", event, err)
		}
	}
	r := adapter.HookReport{SessionID: p.SessionID}
	switch event {
	case "SessionStart", "Stop":
		r.State = "idle"
	case "UserPromptSubmit", "PreToolUse":
		r.State = "working"
	case "SessionEnd":
		r.State = "done"
	case "Notification":
		switch p.NotificationType {
		case "permission_prompt":
			r.State, r.Kind, r.Detail = "blocked", "permission", p.Message
		case "idle_prompt":
			r.State = "idle"
		default:
			return r, false, nil
		}
	default:
		return adapter.HookReport{}, false, fmt.Errorf("unknown hook event %q", event)
	}
	return r, true, nil
}
