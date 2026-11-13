package api

// This file holds run.code-check: a review's code check. It runs the
// check op in-process and build-check (buildCheckRun) on the run's change
// set, both in the run's worktree (the caller's checkout once it is
// gone), and journals the outcome as a code-check entry. It passes when
// rhiza check finds no failure and every leaf has Evidence by an actor
// that wrote none of its commits; otherwise it names each cause. It
// writes nothing but that journal entry.
//
// It also runs the model lints (internal/modellint) and lists their hits
// as lints. A hit is advice for the model check: never a cause, and it
// never changes passed or the exit code.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/check"
	"github.com/dphaener/rhiza/internal/id"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/modellint"
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
	Lints  []RunLintHit   `json:"lints" desc:"The model lints' hits (stale references, unapplied Decisions, capabilities whose code the branch removed); a hit is never a cause and does not change passed"`
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
	resp.Lints, err = codeCheckLints(ctx, rr.state.Change, rr.state.Worktree)
	if err != nil {
		return nil, err
	}
	resp.Step, resp.DryRun = rr.state.Current, ctx.DryRun
	e := run.Entry{Type: run.TypeCodeCheck, At: ctx.Now(), Step: rr.state.Current,
		Outcome: run.CodeCheckPass, Detail: strings.Join(resp.Causes, "; "), Lints: runLintEntries(resp.Lints)}
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
	resp := RunCodeCheckResponse{Run: br.Run, Change: br.Change, Causes: []string{}, Check: cr, Leaves: br.Leaves, Lints: []RunLintHit{}}
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
	r, err := runCheck(&c, CheckRequest{SkipReference: true, BranchNewerFormat: true})
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
	} else {
		fmt.Fprintf(w, "code-check failed: %d %s\n\n", len(r.Causes), shipPlural(len(r.Causes), "cause", "causes"))
		for _, c := range r.Causes {
			fmt.Fprintf(w, "- %s\n", c)
		}
	}
	if len(r.Lints) > 0 {
		fmt.Fprintf(w, "\nlints: %d %s\n\n", len(r.Lints), shipPlural(len(r.Lints), "hit", "hits"))
		for _, h := range r.Lints {
			fmt.Fprintf(w, "- %s\n", h.line())
		}
	}
	if !r.Passed {
		fmt.Fprintf(w, "\nnext: rhiza run status %s\n", r.Run)
	}
}

// RunLintHit is one model lint hit: a node and what is wrong with it.
type RunLintHit struct {
	Node   string `json:"node" desc:"The node's full ID"`
	Handle string `json:"handle,omitempty" desc:"The node's handle"`
	Field  string `json:"field,omitempty" desc:"The field the hit is in: statement or acceptance"`
	Kind   string `json:"kind" enum:"path,test,script,unretired-capability,unapplied-decision,removed-code-capability" desc:"What is wrong: a path, Go test or testscript the head tree lacks, a capability all of whose paths are gone, a Decision never applied, a capability whose code the branch removed"`
	Ref    string `json:"ref,omitempty" desc:"The missing or removed reference, as written"`
}

func (h RunLintHit) line() string {
	s := h.Kind + " " + cmpName(h.Handle, h.Node)
	if h.Field != "" {
		s += " " + h.Field
	}
	if h.Ref != "" {
		s += ": " + h.Ref
	}
	return s
}

func cmpName(handle, node string) string {
	if handle != "" {
		return handle
	}
	return node
}

// runLintEntries is the hits as the journal keeps them.
func runLintEntries(hs []RunLintHit) []run.LintHit {
	var out []run.LintHit
	for _, h := range hs {
		out = append(out, run.LintHit{Node: h.Node, Handle: h.Handle, Field: h.Field, Kind: h.Kind, Ref: h.Ref})
	}
	return out
}

// codeCheckLints runs the model lints on change set change in worktree
// (the caller's checkout when it is gone): the head tree, the Leaf:
// trailers of the commits from the merge-base to HEAD and the paths
// deleted in merge-base..HEAD. The result is never nil.
func codeCheckLints(ctx *Ctx, change, worktree string) ([]RunLintHit, error) {
	out := []RunLintHit{}
	cs, err := id.Parse(change)
	if err != nil {
		return out, Errorf("run names change set %q: %v", change, err)
	}
	dir := ctx.Root
	var w *workspace
	if _, serr := os.Stat(worktree); worktree != "" && serr == nil {
		prefix, err := ctx.Git("rev-parse", "--show-prefix")
		if err != nil {
			return nil, Errorf("%v", err)
		}
		dir = filepath.Join(worktree, strings.TrimSpace(prefix))
		if w, err = openWorkspace(dir, ctx.Now, ctx.IDs, true); err != nil {
			return nil, err
		}
	} else if w, err = openCtxWorkspace(ctx); err != nil {
		return nil, err
	}
	c, ok := w.changeSet(cs)
	if !ok {
		return out, nil
	}
	git := func(args ...string) (string, error) {
		o, err := ctx.Exec("git", append([]string{"-C", dir}, args...), dir, nil)
		return string(o), err
	}
	sets := w.changeSets()
	m, _ := gmodel.WhatIf(sets, cs)
	if c.State == changeset.StateLanded {
		m, _ = gmodel.AsIs(sets)
	}

	// top is the git top, whose paths the lints' references name.
	top := dir
	if o, err := git("rev-parse", "--show-toplevel"); err == nil && strings.TrimSpace(o) != "" {
		top = strings.TrimSpace(o)
	}
	head := os.DirFS(top)

	var committed, deleted []string
	if base, _, _, err := checkBase(git, ""); err == nil && base != "" {
		for _, cm := range gitUnclaimed(git, base) {
			committed = append(committed, cm.Leaves...)
		}
		if o, err := git("diff", "--name-only", "--diff-filter=D", "--no-renames", base+"..HEAD"); err == nil {
			for _, l := range strings.Split(o, "\n") {
				if l = strings.TrimSpace(l); l != "" {
					deleted = append(deleted, l)
				}
			}
		}
	}

	var hits []modellint.Hit
	hits = append(hits, modellint.StaleRefs(m, c, head)...)
	hits = append(hits, modellint.UnappliedDecisions(m, c, committed)...)
	hits = append(hits, modellint.RemovedCode(m, c, head, deleted)...)
	for _, h := range hits {
		out = append(out, RunLintHit{Node: h.Node.String(), Handle: h.Handle, Field: h.Field, Kind: h.Kind, Ref: h.Ref})
	}
	return out, nil
}
