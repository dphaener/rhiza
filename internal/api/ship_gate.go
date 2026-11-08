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

// ShipGateRequest takes nothing: the gate is a project setting.
type ShipGateRequest struct{}

// ShipGateResponse is how the gate went.
type ShipGateResponse struct {
	Command string `json:"command,omitempty" desc:"The gate command (rhiza.ship.localGate); empty if none is set"`
	Skipped bool   `json:"skipped" desc:"True if the gate did not run: CI is stable, or no gate is set"`
	Reason  string `json:"reason,omitempty" desc:"Why it ran (why CI is unstable) or why it was skipped"`
	Passed  bool   `json:"passed" desc:"True if the gate ran and exited 0"`
	Output  string `json:"output,omitempty" desc:"The end of the gate's combined output, when it ran"`
	Error   string `json:"error,omitempty" desc:"How a failing gate failed: its exit status or a timeout"`
}

// ExitCode is 0 when the gate passed or was skipped, else 1.
func (r ShipGateResponse) ExitCode() int {
	if r.Skipped || r.Passed {
		return 0
	}
	return 1
}

// Human says how the gate stands, with a failing gate's output.
func (r ShipGateResponse) Human(w io.Writer) {
	switch {
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

func runShipGate(ctx *Ctx, _ any) (any, error) {
	if err := refuseShipDryRun(ctx, "ship.gate"); err != nil {
		return nil, err
	}
	rows, err := settingsOf(ctx, nil)
	if err != nil {
		return nil, err
	}
	gate := strings.TrimSpace(settingValue(rows, LocalGateKey))
	if gate == "" {
		return ShipGateResponse{Skipped: true, Reason: "no gate is set (" + LocalGateKey + ")"}, nil
	}
	v, err := ciStabilityVerdict(ctx)
	if err != nil {
		return nil, err
	}
	if !v.Unstable {
		return ShipGateResponse{Command: gate, Skipped: true, Reason: "CI is stable"}, nil
	}
	t, err := findShipTarget(ctx)
	if err != nil {
		return nil, err
	}
	resp := ShipGateResponse{Command: gate, Reason: v.Reason}
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
		resp.Passed = true
	}
	return resp, nil
}
