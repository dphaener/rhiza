package api

// This file holds ship.gate: while CI is unstable, the project's local gate
// command (rhiza.ship.localGate) runs in the run's worktree before a push,
// and the push waits for it to exit 0. While CI is stable, or with no gate
// set, it passes at once and the run keeps its scoped tests.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

func init() {
	register(Operation{
		Name:     "ship.gate",
		Summary:  "Before a push, while CI is unstable, run the project's local gate command (rhiza.ship.localGate) in the change set's worktree and fail if it does; with CI stable or no gate set, do nothing. The command gets RHIZA_GATE_HEAD (the commit it checks) and RHIZA_GATE_BASE (the merge-base of that head with origin/main, empty if none is found) in its environment, so it can check only what the branch changes",
		Request:  reflect.TypeFor[ShipGateRequest](),
		Response: reflect.TypeFor[ShipGateResponse](),
		Run:      runShipGate,
	})
}

// LocalGateKey is the setting holding the project's local gate command.
const LocalGateKey = "rhiza.ship.localGate"

// PrePushCheckKey is the setting holding the project's pre-push check
// command, run after rhiza check before every push of a run.
const PrePushCheckKey = "rhiza.ship.prePushCheck"

// shipGateTimeout bounds one gate run; shipGateTail is how much of its
// output a response keeps (the end, where a failure is named).
const (
	shipGateTimeout = time.Hour
	shipGateTail    = 16 << 10
)

// ShipGateRequest takes nothing to run the gate, which is a project
// setting; Failed asks a routing question instead.
type ShipGateRequest struct {
	Failed bool `json:"failed,omitempty" flag:"failed" desc:"Run no gate: exit 0 only if the branch is still at an unpushed head the local gate failed on (ship rebase left it) or the pre-push check failed on (ship open-pr or land), else exit 1 (a conflict or another failure); reads the state file only (for a step that routes a gate failing on a rebased head to a fixer)"`
}

// ShipGateResponse is how the gate went.
type ShipGateResponse struct {
	Command string `json:"command,omitempty" desc:"The gate command (rhiza.ship.localGate); empty if none is set"`
	Skipped bool   `json:"skipped" desc:"True if the gate did not run: CI is stable, or no gate is set"`
	Reason  string `json:"reason,omitempty" desc:"Why it ran (why CI is unstable) or why it was skipped"`
	Passed  bool   `json:"passed" desc:"True if the gate ran and exited 0"`
	Output  string `json:"output,omitempty" desc:"The end of the gate's combined output, when it ran"`
	Error   string `json:"error,omitempty" desc:"How a failing gate failed: its exit status or a timeout"`
	Log     string `json:"log,omitempty" desc:"The file holding the gate's whole combined output (in the run's directory, one file per gate run; outside a run rhiza/gate/<branch>.log under the git common dir, replaced on each run), when it ran"`
	Head    string `json:"head,omitempty" desc:"The commit the gate passed on"`
	Failed  string `json:"failed,omitempty" desc:"With --failed: the unpushed head ship rebase left after the gate failed on it, when the branch is still at it"`
	// asked is true for --failed, whose exit code says whether Failed is set.
	asked bool
	// failLines are the failing lines of the whole output, for Human.
	failLines []string
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
		fmt.Fprintf(w, "The local gate failed on %s, which was left unpushed (a rebase merged origin/%s in, or the pre-push check failed): run rhiza ship gate there, or see the pre-push check in rhiza run status, for the failure.\n", shortSHA(r.Failed), shipMain)
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
		for _, l := range r.failLines {
			fmt.Fprintln(w, l)
		}
		if r.Log != "" {
			fmt.Fprintf(w, "Full log: %s\n", r.Log)
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

// shipRecordPrePushFailed records the head the pre-push check failed on in
// the ship state, for ship gate --failed. Best effort: the failure is
// reported either way.
func shipRecordPrePushFailed(t shipTarget, c *ShipPrePushCheck) {
	st, err := readShipState(t.statePath())
	if err != nil {
		return
	}
	st.PrePushFailed = c.Head
	writeShipState(t.statePath(), st)
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
	if h := strings.TrimSpace(head); h != "" && (h == st.GateFailed || h == st.PrePushFailed) {
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
		reason := "CI is stable"
		if v.Reason != "" {
			reason += " (" + v.Reason + ")"
		}
		return ShipGateResponse{Command: gate, Skipped: true, Reason: reason}, nil
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
	clean := head != "" && shipTreeClean(ctx, t.dir())
	tree := ""
	if clean {
		tree = shipCodeTree(ctx, t.dir())
	}
	if tree != "" && st.GatedTree == tree && st.GatedCommand == gate && st.Gated != "" {
		resp.Passed, resp.Head = true, head
		resp.Reason = v.Reason + "; it already passed on " + shortSHA(st.Gated) + " with the same code outside .rhiza/, so the gate was not run again"
		return resp, nil
	}
	c, cancel := context.WithTimeout(context.Background(), shipGateTimeout)
	defer cancel()
	cmd := exec.CommandContext(c, "sh", "-c", gate)
	cmd.Dir = t.dir()
	base, _ := gitAt(ctx, t.dir(), "merge-base", "HEAD", "origin/"+shipMain)
	cmd.Env = append(os.Environ(), "RHIZA_GATE_HEAD="+head, "RHIZA_GATE_BASE="+strings.TrimSpace(base))
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	cmd.WaitDelay = time.Second
	runErr := cmd.Run()
	resp.Output = out.String()
	resp.failLines = shipGateFailLines(resp.Output)
	resp.Log = shipWriteGateLog(ctx, t, out.Bytes())
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
		if clean {
			st, err := readShipState(t.statePath())
			if err != nil {
				return ShipGateResponse{}, err
			}
			st.Gated, st.GatedTree, st.GatedCommand = head, tree, gate
			if err := writeShipState(t.statePath(), st); err != nil {
				return ShipGateResponse{}, err
			}
		}
	}
	return resp, nil
}

// shipGateMaxFailLines caps the failing lines a failed gate's output ends with.
const shipGateMaxFailLines = 10

// shipGateFailLines picks a gate output's failing go test lines ("--- FAIL:",
// "FAIL", "panic:"), the first shipGateMaxFailLines of them.
func shipGateFailLines(out string) []string {
	var lines []string
	for _, l := range strings.Split(out, "\n") {
		l = strings.TrimRight(l, "\r")
		t := strings.TrimLeft(l, " \t")
		if strings.HasPrefix(t, "--- FAIL:") || strings.HasPrefix(l, "FAIL") || strings.HasPrefix(l, "panic:") {
			lines = append(lines, l)
			if len(lines) == shipGateMaxFailLines {
				break
			}
		}
	}
	return lines
}

// shipWriteGateLog writes the gate's whole output to its own file and
// returns the path; empty if it could not be written (the response still has
// the tail). In a run (RHIZA_RUN) the file is <runs>/<run id>/gate-<n>.log,
// n counting up so a later gate of the run never overwrites an earlier one;
// outside a run it is <common>/rhiza/gate/<branch>.log, replaced each time.
func shipWriteGateLog(ctx *Ctx, t *shipTarget, out []byte) string {
	if t.common == "" {
		return ""
	}
	if id := strings.TrimSpace(ctxGetenv(ctx)(envRun)); id != "" {
		if j, err := run.NewStore(run.Dir(t.common)).Read(id); err == nil {
			id = j.ID
		}
		if strings.ContainsAny(id, `/\`) || strings.HasPrefix(id, ".") {
			return ""
		}
		dir := filepath.Join(run.Dir(t.common), id)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return ""
		}
		for n := 1; n < 10000; n++ {
			path := filepath.Join(dir, fmt.Sprintf("gate-%d.log", n))
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
			if errors.Is(err, fs.ErrExist) {
				continue
			}
			if err != nil {
				return ""
			}
			_, err = f.Write(out)
			if cerr := f.Close(); err == nil {
				err = cerr
			}
			if err != nil {
				return ""
			}
			return path
		}
		return ""
	}
	if t.branch == "" {
		return ""
	}
	path := filepath.Join(t.common, "rhiza", "gate", filepath.FromSlash(t.branch)+".log")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return ""
	}
	if err := writeAtomic(path, out); err != nil {
		return ""
	}
	return path
}

// shipTreeClean reports whether dir has no uncommitted changes to tracked
// files, so a gate run there tested its HEAD.
func shipTreeClean(ctx *Ctx, dir string) bool {
	out, err := gitAt(ctx, dir, "status", "--porcelain", "--untracked-files=no")
	return err == nil && strings.TrimSpace(out) == ""
}

// shipCodeTree hashes HEAD's tree with the top-level .rhiza/ entry left out,
// so commits that change only the model, its rendered files or the log keep
// the hash. Empty if it cannot be worked out: nothing is then reused.
func shipCodeTree(ctx *Ctx, dir string) string {
	out, err := gitAt(ctx, dir, "ls-tree", "HEAD")
	if err != nil || strings.TrimSpace(out) == "" {
		return ""
	}
	h := sha256.New()
	for _, l := range strings.Split(out, "\n") {
		if l == "" {
			continue
		}
		if _, name, ok := strings.Cut(l, "\t"); ok && name == ".rhiza" {
			continue
		}
		io.WriteString(h, l+"\n")
	}
	return hex.EncodeToString(h.Sum(nil))
}
