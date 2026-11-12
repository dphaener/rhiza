package api

// This file holds ship.generate: the project's generate command
// (rhiza.ship.generate) regenerates its checked-in generated files in the
// current worktree, and whatever it changed is staged, so the render step
// of a build or review run commits it with the rendered folder. With no
// command set it does nothing.

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"strings"
	"time"
)

func init() {
	register(Operation{
		Name:     "ship.generate",
		Summary:  "Run the project's generate command (rhiza.ship.generate) in the current worktree and stage the files it changed; with no command set, do nothing",
		Request:  reflect.TypeFor[ShipGenerateRequest](),
		Response: reflect.TypeFor[ShipGenerateResponse](),
		Run:      runShipGenerate,
	})
}

// GenerateKey is the setting holding the project's generate command.
const GenerateKey = "rhiza.ship.generate"

// shipGenerateTimeout bounds one generate run.
const shipGenerateTimeout = 30 * time.Minute

// ShipGenerateRequest takes nothing: the command is a project setting.
type ShipGenerateRequest struct{}

// ShipGenerateResponse is how the generate command went.
type ShipGenerateResponse struct {
	Command string   `json:"command,omitempty" desc:"The generate command (rhiza.ship.generate); empty if none is set, which skips it"`
	Skipped bool     `json:"skipped" desc:"True if no generate command is set"`
	Passed  bool     `json:"passed" desc:"True if the command ran and exited 0"`
	Staged  []string `json:"staged,omitempty" desc:"The files the command changed, now staged"`
	Output  string   `json:"output,omitempty" desc:"The end of the command's combined output, when it ran"`
	Error   string   `json:"error,omitempty" desc:"How a failing command failed: its exit status or a timeout"`
}

// ExitCode is 0 when the command passed or was skipped, else 1.
func (r ShipGenerateResponse) ExitCode() int {
	if r.Skipped || r.Passed {
		return 0
	}
	return 1
}

// Human says how the command went, with a failing command's output.
func (r ShipGenerateResponse) Human(w io.Writer) {
	switch {
	case r.Skipped:
		fmt.Fprintf(w, "No generate command set (%s): nothing to generate.\n", GenerateKey)
	case r.Passed:
		fmt.Fprintf(w, "Generate passed (%s): %d changed file(s) staged.\n", r.Command, len(r.Staged))
	default:
		if r.Output != "" {
			fmt.Fprintln(w, strings.TrimRight(r.Output, "\n"))
		}
		fmt.Fprintf(w, "Generate failed (%s): %s.\n", r.Command, r.Error)
	}
}

func runShipGenerate(ctx *Ctx, req any) (any, error) {
	if err := refuseShipDryRun(ctx, "ship.generate"); err != nil {
		return nil, err
	}
	rows, err := settingsOf(ctx, nil)
	if err != nil {
		return nil, err
	}
	gen := strings.TrimSpace(settingValue(rows, GenerateKey))
	if gen == "" {
		return ShipGenerateResponse{Skipped: true}, nil
	}
	dir := ctx.Wd
	if dir == "" {
		dir = ctx.Root
	}
	top, err := gitAt(ctx, dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, err
	}
	dir = strings.TrimSpace(top)
	before, err := shipGenerateDirty(ctx, dir)
	if err != nil {
		return nil, err
	}
	resp := ShipGenerateResponse{Command: gen}
	c, cancel := context.WithTimeout(context.Background(), shipGenerateTimeout)
	defer cancel()
	cmd := exec.CommandContext(c, "sh", "-c", gen)
	cmd.Dir = dir
	cmd.Env = os.Environ()
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	cmd.WaitDelay = time.Second
	runErr := cmd.Run()
	resp.Output = out.String()
	if len(resp.Output) > shipGateTail {
		resp.Output = "...\n" + resp.Output[len(resp.Output)-shipGateTail:]
	}
	switch {
	case c.Err() != nil:
		resp.Error = "timed out after " + shipGenerateTimeout.String()
		return resp, nil
	case runErr != nil:
		resp.Error = runErr.Error()
		return resp, nil
	}
	after, err := shipGenerateDirty(ctx, dir)
	if err != nil {
		return nil, err
	}
	// Stage what the command changed: a path whose status is new or
	// different since before, never the run's own uncommitted records.
	for path, st := range after {
		if before[path] != st {
			resp.Staged = append(resp.Staged, path)
		}
	}
	sort.Strings(resp.Staged)
	if len(resp.Staged) > 0 {
		args := append([]string{"add", "-A", "--"}, resp.Staged...)
		if _, err := gitAt(ctx, dir, args...); err != nil {
			return nil, err
		}
	}
	resp.Passed = true
	return resp, nil
}

// shipGenerateDirty maps each path git status lists, untracked files
// included, to its status line plus a hash of its content, so a path the
// command rewrote again differs from before.
func shipGenerateDirty(ctx *Ctx, dir string) (map[string]string, error) {
	out, err := gitAt(ctx, dir, "status", "--porcelain", "-z", "--untracked-files=all")
	if err != nil {
		return nil, err
	}
	m := map[string]string{}
	fields := strings.Split(out, "\x00")
	for i := 0; i < len(fields); i++ {
		f := fields[i]
		if len(f) < 4 {
			continue
		}
		code, path := f[:2], f[3:]
		if code[0] == 'R' || code[0] == 'C' {
			i++ // the rename's source is the next field
		}
		sum, _ := gitAt(ctx, dir, "hash-object", "--", path)
		m[path] = code + strings.TrimSpace(sum)
	}
	return m, nil
}
