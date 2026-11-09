package api

// This file holds run.code-check: a review's code check. It runs the
// check op in-process and build-check (buildCheckRun) on the run's change
// set, both in the run's worktree (the caller's checkout once it is
// gone), and journals the outcome as a code-check entry. It passes when
// rhiza check finds no failure and every leaf has Evidence by an actor
// that wrote none of its commits; otherwise it names each cause. It
// writes nothing but that journal entry.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/check"
	"github.com/dphaener/rhiza/internal/run"
)

func init() {
	register(Operation{
		Name:     "run.code-check",
		Summary:  "Run rhiza check and build-check on the change set of the run RHIZA_RUN names and journal the outcome as its code check: pass when rhiza check finds no failure and every leaf has Evidence by a verifier, else fail naming each cause",
		Write:    true,
		Request:  reflect.TypeFor[RunCodeCheckRequest](),
		Response: reflect.TypeFor[RunCodeCheckResponse](),
		Run:      runRunCodeCheck,
	})
}

const runCodeCheckNext = "rhiza describe run.code-check"

// RunCodeCheckRequest takes nothing: the run is RHIZA_RUN's.
type RunCodeCheckRequest struct{}

// RunCodeCheckResponse is the code check's outcome. It is printed whether
// or not it passed; a failing check exits 1.
type RunCodeCheckResponse struct {
	Run    string         `json:"run" desc:"The run's full ID"`
	Change string         `json:"change" desc:"The run's change set"`
	Step   string         `json:"step,omitempty" desc:"The step the run is at, which the code-check entry names"`
	Passed bool           `json:"passed" desc:"True if rhiza check found no failure and every leaf is OK"`
	Causes []string       `json:"causes" desc:"Why it failed, one cause per line: each failure rhiza check found, then each leaf lacking Evidence by a verifier; empty when it passed"`
	Check  CheckResponse  `json:"check" desc:"rhiza check's report on the run's worktree"`
	Leaves []RunBuildLeaf `json:"leaves" desc:"build-check's leaves of the change set, in leaf order"`
	DryRun bool           `json:"dry_run,omitempty" desc:"True if nothing was journalled"`
}

// ExitCode is 0 when the check passed, else 1.
func (r RunCodeCheckResponse) ExitCode() int {
	if r.Passed {
		return 0
	}
	return 1
}

// runRunCodeCheck checks the unfinished run RHIZA_RUN names: rhiza check
// in its worktree, then build-check of its change set, and journals the
// outcome with the causes as its detail.
func runRunCodeCheck(ctx *Ctx, _ any) (any, error) {
	rr, err := openReviewRun(ctx, "run code-check", runCodeCheckNext)
	if err != nil {
		return nil, err
	}
	if ctx.Exec == nil || ctx.Git == nil {
		return nil, Errorf("run code-check runs rhiza check and reads git history, and this caller runs no programs")
	}
	cr, err := checkIn(ctx, rr.state.Worktree)
	if err != nil {
		return nil, err
	}
	br, err := buildCheckRun(ctx, rr.id, rr.state.Change, rr.state.Worktree)
	if err != nil {
		return nil, err
	}
	resp := codeCheck(cr, br)
	resp.Step, resp.DryRun = rr.state.Current, ctx.DryRun
	e := run.Entry{Type: run.TypeCodeCheck, At: ctx.Now(), Step: rr.state.Current,
		Outcome: run.CodeCheckPass, Detail: strings.Join(resp.Causes, "; ")}
	if !resp.Passed {
		e.Outcome = run.CodeCheckFail
	}
	if err := rr.journal(ctx, "code check", e); err != nil {
		return nil, err
	}
	return resp, nil
}

// codeCheck is the code check's outcome from rhiza check's report cr and
// build-check's br: the failures cr found, then the leaves br found not
// OK, each a cause.
func codeCheck(cr CheckResponse, br RunBuildCheckResponse) RunCodeCheckResponse {
	resp := RunCodeCheckResponse{Run: br.Run, Change: br.Change, Causes: []string{}, Check: cr, Leaves: br.Leaves}
	for _, f := range cr.Findings {
		if f.Severity == check.SeverityFail {
			resp.Causes = append(resp.Causes, "rhiza check: "+f.line())
		}
	}
	if !cr.OK && len(resp.Causes) == 0 {
		resp.Causes = append(resp.Causes, "rhiza check failed")
	}
	for _, l := range br.Leaves {
		if !l.OK {
			resp.Causes = append(resp.Causes, fmt.Sprintf("leaf %s: %s", l.Handle, l.Why))
		}
	}
	resp.Passed = len(resp.Causes) == 0
	return resp
}

// checkIn runs the check op in worktree, at the folder matching the
// caller's place below its git top, as buildCheckRun reads it; in the
// caller's checkout when worktree does not exist.
func checkIn(ctx *Ctx, worktree string) (CheckResponse, error) {
	c := *ctx
	c.DryRun = false
	if _, serr := os.Stat(worktree); worktree != "" && serr == nil {
		prefix, err := ctx.Git("rev-parse", "--show-prefix")
		if err != nil {
			return CheckResponse{}, Errorf("%v", err)
		}
		dir := filepath.Join(worktree, strings.TrimSpace(prefix))
		c.Root = dir
		c.Git = func(args ...string) (string, error) {
			out, err := ctx.Exec("git", append([]string{"-C", dir}, args...), dir, nil)
			return string(out), err
		}
	}
	r, err := runCheck(&c, CheckRequest{SkipReference: true})
	if err != nil {
		return CheckResponse{}, err
	}
	return r.(CheckResponse), nil
}

// Human reports the check: one line when it passed, else each cause.
func (r RunCodeCheckResponse) Human(w io.Writer) {
	if r.Passed {
		fmt.Fprintf(w, "code-check passed: rhiza check is ok and %d %s Evidence by a verifier\n",
			len(r.Leaves), shipPlural(len(r.Leaves), "leaf has", "leaves have"))
		return
	}
	fmt.Fprintf(w, "code-check failed: %d %s\n\n", len(r.Causes), shipPlural(len(r.Causes), "cause", "causes"))
	for _, c := range r.Causes {
		fmt.Fprintf(w, "- %s\n", c)
	}
	fmt.Fprintf(w, "\nnext: rhiza run status %s\n", r.Run)
}
