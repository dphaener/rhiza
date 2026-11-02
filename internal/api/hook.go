package api

// This file holds rhiza hook ls: the daemon's notify and event commands,
// each with where it is set. A live daemon's recorded values win over git
// config, because they are what it runs.

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/daemon"
)

// The hooks, and the sources of a hook's command.
const (
	HookNotify = "notify"
	HookEvent  = "event"

	HookFlag   = "daemon flag"
	HookConfig = "git config"
	HookUnset  = "unset"
)

// hookKeys maps a hook to its git config key.
var hookKeys = map[string]string{
	HookNotify: "rhiza.notifyCommand",
	HookEvent:  "rhiza.eventCommand",
}

// hookNames is the hooks in listing order.
var hookNames = []string{HookNotify, HookEvent}

func init() {
	register(Operation{
		Name:     "hook.ls",
		Summary:  "List the notify and event hook commands, each with its command and source: daemon flag (from a live daemon), git config or unset",
		Request:  reflect.TypeFor[HookLsRequest](),
		Response: reflect.TypeFor[HookLsResponse](),
		Run:      runHookLs,
	})
}

// HookLsRequest takes nothing.
type HookLsRequest struct{}

// HookRow is one hook command and where it comes from.
type HookRow struct {
	Name    string `json:"name" desc:"The hook: notify or event"`
	Command string `json:"command" desc:"The command, run through sh -c; empty when unset"`
	Source  string `json:"source" desc:"Where it is set: daemon flag, git config or unset"`
}

func (h HookRow) text() string {
	if h.Source == HookUnset {
		return HookUnset
	}
	return fmt.Sprintf("%s (%s)", h.Command, h.Source)
}

// HookLsResponse lists the hooks.
type HookLsResponse struct {
	Hooks []HookRow `json:"hooks" desc:"The notify hook, then the event hook"`
}

// Human writes the hooks as text.
func (r HookLsResponse) Human(w io.Writer) {
	for _, h := range r.Hooks {
		fmt.Fprintf(w, "%s: %s\n", h.Name, h.text())
	}
}

// resolveHooks takes each hook's command from its daemon flag, else git
// config, else unset.
func resolveHooks(ctx *Ctx, notifyFlag, eventFlag string) []daemon.HookValue {
	flags := map[string]string{HookNotify: notifyFlag, HookEvent: eventFlag}
	var out []daemon.HookValue
	for _, name := range hookNames {
		h := daemon.HookValue{Name: name, Source: HookUnset}
		if c := strings.TrimSpace(flags[name]); c != "" {
			h.Command, h.Source = c, HookFlag
		} else if ctx.Git != nil {
			if v, err := ctx.Git("config", "--get", hookKeys[name]); err == nil {
				if c := strings.TrimSpace(v); c != "" {
					h.Command, h.Source = c, HookConfig
				}
			}
		}
		out = append(out, h)
	}
	return out
}

// hookCommand is the command of the hook called name, "" when unset.
func hookCommand(hooks []daemon.HookValue, name string) string {
	for _, h := range hooks {
		if h.Name == name {
			return h.Command
		}
	}
	return ""
}

func runHookLs(ctx *Ctx, req any) (any, error) {
	hooks := resolveHooks(ctx, "", "")
	if common, err := gitCommonDir(ctx); err == nil {
		if _, running := daemon.Running(common); running {
			if st, ok := daemon.ReadState(common); ok && st.Process == daemon.ReadHolder(common) && len(st.Hooks) > 0 {
				hooks = st.Hooks
			}
		}
	}
	resp := HookLsResponse{Hooks: []HookRow{}}
	for _, h := range hooks {
		resp.Hooks = append(resp.Hooks, HookRow{Name: h.Name, Command: h.Command, Source: h.Source})
	}
	return resp, nil
}
