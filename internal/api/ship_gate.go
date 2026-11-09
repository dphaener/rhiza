package api

// This file holds ship.gate: while CI is unstable, the project's local gate
// command (rhiza.ship.localGate) runs in the run's worktree before a push,
// and the push waits for it to exit 0. While CI is stable, or with no gate
// set, it passes at once and the run keeps its scoped tests.

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

func init() {
	register(Operation{
		Name:     "ship.gate",
		Summary:  "Before a push, while CI is unstable, run the project's local gate command (rhiza.ship.localGate) in the change set's worktree and fail if it does; with CI stable or no gate set, do nothing",
		Request:  reflect.TypeFor[ShipGateRequest](),
		Response: reflect.TypeFor[ShipGateResponse](),
		Run:      runShipGate,
	})
}

// LocalGateKey is the setting holding the project's local gate command.
const LocalGateKey = "rhiza.ship.localGate"

// shipGateTimeout bounds one gate run; shipGateTail is how much of its
// output a response keeps (the end, where a failure is named).
const (
	shipGateTimeout = time.Hour
	shipGateTail    = 16 << 10
)

// ShipGateRequest takes nothing to run the gate, which is a project
// setting; Failed asks a routing question instead.
type ShipGateRequest struct {
	Failed bool `json:"failed,omitempty" flag:"failed" desc:"Run no gate: exit 0 only if the last ship rebase left an unpushed head the local gate failed on and the branch is still at it, else exit 1 (a conflict or another failure); reads the state file only (for a step that routes a gate failing on a rebased head to a fixer)"`
}

// ShipGateResponse is how the gate went.
type ShipGateResponse struct {
	Command string `json:"command,omitempty" desc:"The gate command (rhiza.ship.localGate); empty if none is set"`
	Skipped bool   `json:"skipped" desc:"True if the gate did not run: CI is stable, or no gate is set"`
	Reason  string `json:"reason,omitempty" desc:"Why it ran (why CI is unstable) or why it was skipped"`
	Passed  bool   `json:"passed" desc:"True if the gate ran and exited 0"`
	Output  string `json:"output,omitempty" desc:"The end of the gate's combined output, when it ran"`
	Error   string `json:"error,omitempty" desc:"How a failing gate failed: its exit status or a timeout"`
	Head    string `json:"head,omitempty" desc:"The commit the gate passed on"`
	Failed  string `json:"failed,omitempty" desc:"With --failed: the unpushed head ship rebase left after the gate failed on it, when the branch is still at it"`
	// asked is true for --failed, whose exit code says whether Failed is set.
	asked bool
}

// ExitCode is 0 when the gate passed or was skipped, else 1.
func (r ShipGateResponse) ExitCode() int {
	if r.asked {
		if r.Failed != "" {
			return 0
		}
		return 1
	}
	if r.Skipped || r.Passed {
		return 0
	}
	return 1
}

// Human says how the gate stands, with a failing gate's output.
func (r ShipGateResponse) Human(w io.Writer) {
	switch {
	case r.asked && r.Failed != "":
		fmt.Fprintf(w, "The local gate failed on %s, which ship rebase left unpushed with origin/%s merged in: run rhiza ship gate there to see the failure.\n", shortSHA(r.Failed), shipMain)
	case r.asked:
		fmt.Fprintln(w, "No head the local gate failed on after a rebase is on the branch.")
	case r.Skipped:
		fmt.Fprintf(w, "Local gate skipped: %s.\n", r.Reason)
	case r.Passed:
		fmt.Fprintf(w, "Local gate passed (%s): %s.\n", r.Command, r.Reason)
	default:
		if r.Output != "" {
			fmt.Fprintln(w, strings.TrimRight(r.Output, "\n"))
		}
		fmt.Fprintf(w, "Local gate failed (%s): %s. CI is unstable (%s): fix it here, never push to see.\n", r.Command, r.Error, r.Reason)
	}
}

func runShipGate(ctx *Ctx, req any) (any, error) {
	if req.(ShipGateRequest).Failed {
		return shipGateFailed(ctx)
	}
	if err := refuseShipDryRun(ctx, "ship.gate"); err != nil {
		return nil, err
	}
	return shipLocalGate(ctx, nil)
}

// shipGateFailed answers ship gate --failed: whether the branch is at the
// head ship rebase left after the gate failed on it, so the step after a
// failed rebase sends that to the gate's fixer and a conflict ends the run.
func shipGateFailed(ctx *Ctx) (any, error) {
	t, err := findShipTarget(ctx)
	if err != nil {
		return nil, err
	}
	st, err := readShipState(t.statePath())
	if err != nil {
		return nil, err
	}
	resp := ShipGateResponse{asked: true}
	head, _ := gitAt(ctx, t.dir(), "rev-parse", "HEAD")
	if h := strings.TrimSpace(head); h != "" && h == st.GateFailed {
		resp.Failed = h
	}
	return resp, nil
}

// shipLocalGate runs the local gate in the change set's worktree (t, or
// the one findShipTarget finds when t is nil) while CI is unstable. A head
// the gate already passed on, recorded in the ship state as Gated, passes
// again without running it: the gate before rebase or land, say, re-checks
// the head the gate before open-pr already passed.
func shipLocalGate(ctx *Ctx, t *shipTarget) (ShipGateResponse, error) {
	rows, err := settingsOf(ctx, nil)
	if err != nil {
		return ShipGateResponse{}, err
	}
	gate := strings.TrimSpace(settingValue(rows, LocalGateKey))
	if gate == "" {
		return ShipGateResponse{Skipped: true, Reason: "no gate is set (" + LocalGateKey + ")"}, nil
	}
	v, err := ciStabilityVerdict(ctx)
	if err != nil {
		return ShipGateResponse{}, err
	}
	if !v.Unstable {
		return ShipGateResponse{Command: gate, Skipped: true, Reason: "CI is stable"}, nil
	}
	if t == nil {
		found, err := findShipTarget(ctx)
		if err != nil {
			return ShipGateResponse{}, err
		}
		t = &found
	}
	resp := ShipGateResponse{Command: gate, Reason: v.Reason}
	head, _ := gitAt(ctx, t.dir(), "rev-parse", "HEAD")
	head = strings.TrimSpace(head)
	st, err := readShipState(t.statePath())
	if err != nil {
		return ShipGateResponse{}, err
	}
	if head != "" && st.Gated == head && shipTreeClean(ctx, t.dir()) {
		resp.Passed, resp.Head, resp.Reason = true, head, v.Reason+"; it already passed on "+shortSHA(head)
		return resp, nil
	}
	c, cancel := context.WithTimeout(context.Background(), shipGateTimeout)
	defer cancel()
	cmd := exec.CommandContext(c, "sh", "-c", gate)
	cmd.Dir = t.dir()
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
		resp.Error = "timed out after " + shipGateTimeout.String()
	case runErr != nil:
		resp.Error = runErr.Error()
	default:
		resp.Passed, resp.Head = true, head
		if head != "" && shipTreeClean(ctx, t.dir()) {
			st, err := readShipState(t.statePath())
			if err != nil {
				return ShipGateResponse{}, err
			}
			st.Gated = head
			if err := writeShipState(t.statePath(), st); err != nil {
				return ShipGateResponse{}, err
			}
		}
	}
	return resp, nil
}

// shipTreeClean reports whether dir has no uncommitted changes to tracked
// files, so a gate run there tested its HEAD.
func shipTreeClean(ctx *Ctx, dir string) bool {
	out, err := gitAt(ctx, dir, "status", "--porcelain", "--untracked-files=no")
	return err == nil && strings.TrimSpace(out) == ""
}
