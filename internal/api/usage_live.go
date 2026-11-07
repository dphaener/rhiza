package api

// This file records the account's live usage. Every agent session Rhiza
// launches runs this command as its status line; it keeps the latest
// rate-limit snapshot per window kind where the usage governor reads it,
// then hands the payload to the owner's own status line command, which the
// session's settings replaced.

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"reflect"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/daemon"
)

func init() {
	register(Operation{
		Name:     "session.status-line",
		Summary:  "Record the live rate limits an agent's status line reports, read from stdin, then print the owner's own status line; the agent's status line runs this",
		Request:  reflect.TypeFor[SessionStatusLineRequest](),
		Response: reflect.TypeFor[SessionStatusLineResponse](),
		Run:      runSessionStatusLine,
	})
}

// userStatusLineTimeout bounds the owner's own status line command.
const userStatusLineTimeout = 5 * time.Second

// SessionStatusLineRequest has no arguments: the payload is stdin.
type SessionStatusLineRequest struct{}

// SessionStatusLineResponse says what the status line recorded. Its text is
// the owner's own status line, if they have one.
type SessionStatusLineResponse struct {
	Recorded int `json:"recorded" desc:"How many rate-limit windows were recorded"`
	// text is the owner's status line, printed on stdout.
	text string
}

// Human writes the owner's own status line, and nothing else.
func (r SessionStatusLineResponse) Human(w io.Writer) { io.WriteString(w, r.text) }

// runSessionStatusLine never fails: a failing status line would get in the
// agent's way, so problems are dropped and the call succeeds.
func runSessionStatusLine(ctx *Ctx, req any) (any, error) {
	var payload []byte
	if ctx.Stdin != nil {
		payload, _ = io.ReadAll(ctx.Stdin)
	}
	var resp SessionStatusLineResponse
	user := ""
	for _, name := range adapter.Names() {
		a, _ := adapter.Lookup(name)
		r, ok := a.(adapter.StatusLineReader)
		if !ok {
			continue
		}
		if user == "" && ctx.Getenv != nil {
			user = r.UserStatusLine(ctx.Getenv)
		}
		windows, err := r.LiveUsage(payload)
		if err != nil || len(windows) == 0 {
			continue
		}
		resp.Recorded += recordLive(ctx, name, windows)
	}
	if user != "" && !ctx.DryRun {
		resp.text = runUserStatusLine(ctx, user, payload)
	}
	return resp, nil
}

// recordLive saves the windows as live snapshots and returns how many it
// saved; best effort.
func recordLive(ctx *Ctx, agent string, windows []adapter.LiveWindow) int {
	if ctx.DryRun {
		return 0
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return 0
	}
	now := ctx.Now().UTC()
	snaps := make([]daemon.LiveSnapshot, 0, len(windows))
	for _, w := range windows {
		snaps = append(snaps, daemon.LiveSnapshot{Kind: w.Kind, UsedPercentage: w.UsedPercentage, ResetsAt: w.ResetsAt, CapturedAt: now, Agent: agent})
	}
	if daemon.RecordLive(limitsDir(ctx, common), snaps) != nil {
		return 0
	}
	return len(snaps)
}

// runUserStatusLine runs the owner's status line command with the payload on
// stdin and returns what it printed; "" when it fails or times out.
func runUserStatusLine(ctx *Ctx, command string, payload []byte) string {
	c, cancel := context.WithTimeout(context.Background(), userStatusLineTimeout)
	defer cancel()
	cmd := exec.CommandContext(c, "sh", "-c", command)
	cmd.Env = os.Environ()
	cmd.Stdin = bytes.NewReader(payload)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.WaitDelay = time.Second
	if cmd.Run() != nil {
		return ""
	}
	return out.String()
}
