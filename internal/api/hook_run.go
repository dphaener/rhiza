package api

// This file holds rhiza hook run: one hook command run by hand with a
// sample payload, through the daemon's own runner. It records nothing: the
// daemon's handled items are untouched.

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/daemon"
)

func init() {
	register(Operation{
		Name:     "hook.run",
		Summary:  "Run the notify or event hook command once with a sample payload, showing its output and exit status; writes nothing",
		Request:  reflect.TypeFor[HookRunRequest](),
		Response: reflect.TypeFor[HookRunResponse](),
		Run:      runHookRun,
	})
}

// HookRunRequest names the hook.
type HookRunRequest struct {
	Hook    string `json:"hook" arg:"0" desc:"The hook to run: notify or event"`
	Command string `json:"command,omitempty" desc:"Run this command instead of the hook's own, through sh -c in the repo root"`
}

// HookRunResponse is how the command ended.
type HookRunResponse struct {
	Hook    string `json:"hook" desc:"The hook that ran"`
	Command string `json:"command" desc:"The command that ran"`
	Source  string `json:"source" desc:"Where the command is set: command flag, daemon flag or git config"`
	Payload string `json:"payload" desc:"The sample JSON sent on stdin"`
	Output  string `json:"output" desc:"What the command wrote to stdout and stderr"`
	Code    int    `json:"code" desc:"The command's exit code: 0 when it succeeded, else non-zero (1 when it timed out or did not start)"`
	Error   string `json:"error,omitempty" desc:"Why it failed; absent on success"`
}

// ExitCode is 0 when the command succeeded, else non-zero.
func (r HookRunResponse) ExitCode() int { return r.Code }

// Human writes the output, then the status.
func (r HookRunResponse) Human(w io.Writer) {
	if r.Output != "" {
		fmt.Fprint(w, r.Output)
		if !strings.HasSuffix(r.Output, "\n") {
			fmt.Fprintln(w)
		}
	}
	if r.Code == 0 {
		fmt.Fprintf(w, "%s hook: exit status 0\n", r.Hook)
		return
	}
	fmt.Fprintf(w, "%s hook failed: %s (exit status %d)\n", r.Hook, r.Error, r.Code)
}

func runHookRun(ctx *Ctx, req any) (any, error) {
	r := req.(HookRunRequest)
	var in []byte
	var env []string
	switch r.Hook {
	case HookNotify:
		it := InboxItem{ID: "sample", Kind: "question", At: "2026-01-01T00:00:00Z", Run: "sample", Change: "sample", ChangeTitle: "Sample change set", Question: "A sample question"}
		in, _ = json.Marshal(it)
		env = []string{
			"RHIZA_INBOX_KIND=" + it.Kind,
			"RHIZA_INBOX_ID=" + it.ID,
			"RHIZA_INBOX_CHANGE=" + it.Change,
			"RHIZA_INBOX_TITLE=" + inboxNotice(it.ChangeTitle, it),
		}
	case HookEvent:
		ev := lifecycleEvent{Event: "accepted", Change: "sample", Entry: "sample", Title: "Sample change set", At: "2026-01-01T00:00:00Z", Actor: "sample"}
		in, _ = json.Marshal(ev)
		env = []string{
			"RHIZA_EVENT=" + ev.Event,
			"RHIZA_EVENT_CHANGE=" + ev.Change,
			"RHIZA_EVENT_ENTRY=" + ev.Entry,
			"RHIZA_EVENT_TITLE=" + ev.Title,
		}
	default:
		return nil, Usagef("unknown hook %q: use notify or event", r.Hook).WithNext("rhiza hook ls")
	}
	command, source := strings.TrimSpace(r.Command), "command flag"
	if command == "" {
		hooks := resolveHooks(ctx, "", "")
		if common, err := gitCommonDir(ctx); err == nil {
			if _, running := daemon.Running(common); running {
				if st, ok := daemon.ReadState(common); ok && st.Process == daemon.ReadHolder(common) && len(st.Hooks) > 0 {
					hooks = st.Hooks
				}
			}
		}
		for _, h := range hooks {
			if h.Name == r.Hook {
				command, source = h.Command, h.Source
			}
		}
		if command == "" {
			return nil, Usagef("the %s hook is unset", r.Hook).WithNext(fmt.Sprintf("git config %s '<command>'", hookKeys[r.Hook]))
		}
	}
	resp := HookRunResponse{Hook: r.Hook, Command: command, Source: source, Payload: string(in)}
	out, err := runHookCommandOutput(ctx, command, in, env)
	resp.Output = out
	if err != nil {
		resp.Error = err.Error()
		resp.Code = 1
		var ee *exec.ExitError
		if errors.As(err, &ee) && ee.ExitCode() > 0 {
			resp.Code = ee.ExitCode()
		}
	}
	return resp, nil
}
