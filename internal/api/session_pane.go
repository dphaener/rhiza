package api

// This file holds session.pane and session.key, which the inbox app's
// nudge uses: pane reads the text of a session's terminal and asks its
// adapter whether it shows a usage-limit prompt; key presses Escape or
// Enter in it. Nothing here names an agent.

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/adapter"
)

func init() {
	register(Operation{
		Name:     "session.pane",
		Summary:  "Read the text on a session's terminal, and say whether it shows the agent parked at a usage-limit prompt",
		Request:  reflect.TypeFor[SessionPaneRequest](),
		Response: reflect.TypeFor[SessionPaneResponse](),
		Run:      runSessionPane,
	})
	register(Operation{
		Name:     "session.key",
		Summary:  "Press a key (escape or enter) in a running agent session's terminal",
		Write:    true,
		Request:  reflect.TypeFor[SessionKeyRequest](),
		Response: reflect.TypeFor[SessionKeyResponse](),
		Run:      runSessionKey,
	})
}

// SessionPaneRequest names the session.
type SessionPaneRequest struct {
	Session string `json:"session" arg:"0" desc:"The session: its full ID or a unique prefix of 4+ characters"`
}

// SessionPaneResponse is the session's pane text.
type SessionPaneResponse struct {
	Session    string `json:"session" desc:"The session's full ID"`
	Text       string `json:"text" desc:"The text on the session's terminal, screen and scrollback"`
	UsageLimit bool   `json:"usage_limit" desc:"True when the session's adapter reads the text as the agent parked at a usage-limit prompt"`
}

// SessionKeyRequest names the session and the key.
type SessionKeyRequest struct {
	Session string `json:"session" arg:"0" desc:"The session: its full ID or a unique prefix of 4+ characters"`
	Key     string `json:"key" arg:"1" desc:"The key: escape or enter"`
}

// SessionKeyResponse is the key pressed.
type SessionKeyResponse struct {
	Session string `json:"session" desc:"The session's full ID"`
	Key     string `json:"key" desc:"The key pressed"`
	DryRun  bool   `json:"dry_run,omitempty" desc:"True if nothing was sent"`
}

func (r SessionPaneResponse) Human(w io.Writer) {
	fmt.Fprint(w, r.Text)
	if r.UsageLimit {
		fmt.Fprintf(w, "\nThe session is parked at a usage-limit prompt.\n")
	}
}

func (r SessionKeyResponse) Human(w io.Writer) {
	if r.DryRun {
		fmt.Fprintf(w, "Would press %s in session %s\n", r.Key, r.Session)
		return
	}
	fmt.Fprintf(w, "Pressed %s in session %s\n", r.Key, r.Session)
}

func runSessionPane(ctx *Ctx, req any) (any, error) {
	r := req.(SessionPaneRequest)
	st, err := sessionStore(ctx)
	if err != nil {
		return nil, err
	}
	rec, err := readSession(st, strings.TrimSpace(r.Session))
	if err != nil {
		return nil, err
	}
	if err := messageable(rec); err != nil {
		return nil, err
	}
	if ctx.Exec == nil {
		return nil, Errorf("cannot read a pane: this caller runs no programs")
	}
	out, err := ctx.Exec("rex", []string{"capture", "-b", rec.Tab.Block}, rec.Worktree, nil)
	if err != nil {
		return nil, Errorf("rex could not read session %s's tab: %v", rec.ID, err).WithNext("rhiza session show " + rec.ID)
	}
	resp := SessionPaneResponse{Session: rec.ID, Text: string(out)}
	if a, ok := adapter.Lookup(rec.Agent); ok {
		if pl, ok := a.(adapter.PaneLimiter); ok {
			resp.UsageLimit = pl.PaneUsageLimit(resp.Text)
		}
	}
	return resp, nil
}

func runSessionKey(ctx *Ctx, req any) (any, error) {
	r := req.(SessionKeyRequest)
	key := strings.ToLower(strings.TrimSpace(r.Key))
	if key != "escape" && key != "enter" {
		return nil, Usagef("key %q: want escape or enter", r.Key).WithNext("rhiza describe session.key")
	}
	st, err := sessionStore(ctx)
	if err != nil {
		return nil, err
	}
	rec, err := readSession(st, strings.TrimSpace(r.Session))
	if err != nil {
		return nil, err
	}
	if err := messageable(rec); err != nil {
		return nil, err
	}
	resp := SessionKeyResponse{Session: rec.ID, Key: key, DryRun: ctx.DryRun}
	if ctx.DryRun {
		return resp, nil
	}
	if ctx.Exec == nil {
		return nil, Errorf("cannot press a key: this caller runs no programs")
	}
	if _, err := ctx.Exec("rex", []string{"send-key", "-b", rec.Tab.Block, key}, rec.Worktree, nil); err != nil {
		return nil, Errorf("rex could not press %s in session %s's tab: %v", key, rec.ID, err).WithNext("rhiza session show " + rec.ID)
	}
	return resp, nil
}
