package api

// This file holds run.build-check: whether every leaf of a run's change
// set was verified by someone other than its builder. A leaf passes when
// live Evidence verifies it whose creating actor is named in no Actor:
// trailer of a commit whose Leaf: trailer names the leaf, over the run
// worktree's history. It reads the model and git only; it writes nothing.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/check"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

func init() {
	register(Operation{
		Name:     "run.build-check",
		Summary:  "Check the change set of the run RHIZA_RUN names: pass when every leaf has live Evidence by an actor that wrote none of its commits, else fail naming each leaf that lacks it",
		Write:    false,
		Request:  reflect.TypeFor[RunBuildCheckRequest](),
		Response: reflect.TypeFor[RunBuildCheckResponse](),
		Run:      runRunBuildCheck,
	})
}

const runBuildCheckNext = "rhiza describe run.build-check"

// RunBuildCheckRequest takes nothing: the run is RHIZA_RUN's.
type RunBuildCheckRequest struct{}

// RunBuildLeaf is one leaf of the run's change set and whether a
// verifier's Evidence covers it.
type RunBuildLeaf struct {
	Leaf     string   `json:"leaf" desc:"The leaf capability's ID"`
	Handle   string   `json:"handle" desc:"Its handle"`
	OK       bool     `json:"ok" desc:"True if live Evidence verifies it written by an actor that wrote none of its commits"`
	Evidence []string `json:"evidence" desc:"The live Evidence that verifies it, by handle (ID when it has none), in link order; empty when none does"`
	Why      string   `json:"why,omitempty" desc:"Why it is not OK; absent when it is"`
}

// RunBuildCheckResponse is the check's outcome over every leaf. It is
// printed whether or not it passed; a failing check exits 1.
type RunBuildCheckResponse struct {
	Run    string         `json:"run" desc:"The run's full ID"`
	Change string         `json:"change" desc:"The run's change set"`
	Passed bool           `json:"passed" desc:"True if every leaf is OK"`
	Leaves []RunBuildLeaf `json:"leaves" desc:"The change set's leaves in leaf order"`
}

// ExitCode is 0 when every leaf is OK, else 1.
func (r RunBuildCheckResponse) ExitCode() int {
	if r.Passed {
		return 0
	}
	return 1
}

// runRunBuildCheck reads the run RHIZA_RUN names, finished or not, and
// judges its change set's leaves (buildCheck) in the model and history of
// its worktree, or of the caller's checkout once the worktree is gone.
func runRunBuildCheck(ctx *Ctx, _ any) (any, error) {
	rr, err := openRun(ctx, "run build-check", runBuildCheckNext, true)
	if err != nil {
		return nil, err
	}
	return buildCheckRun(ctx, rr.id, rr.state.Change, rr.state.Worktree)
}

// buildCheckRun is build-check over change set change of run rid, read
// in worktree (the caller's checkout when it does not exist). Other
// checks, such as a review's, may call it with their run.
func buildCheckRun(ctx *Ctx, rid, change, worktree string) (RunBuildCheckResponse, error) {
	cs, err := id.Parse(change)
	if err != nil {
		return RunBuildCheckResponse{}, Errorf("run %s names change set %q: %v", rid, change, err)
	}
	if ctx.Exec == nil {
		return RunBuildCheckResponse{}, Errorf("run build-check reads git history, and this caller runs no programs")
	}
	var w *workspace
	dir := ctx.Root
	if _, serr := os.Stat(worktree); worktree != "" && serr == nil {
		prefix, err := ctx.Git("rev-parse", "--show-prefix")
		if err != nil {
			return RunBuildCheckResponse{}, Errorf("%v", err)
		}
		dir = filepath.Join(worktree, strings.TrimSpace(prefix))
		w, err = openWorkspace(dir, ctx.Now, ctx.IDs, true)
		if err != nil {
			return RunBuildCheckResponse{}, err
		}
	} else if w, err = openCtxWorkspace(ctx); err != nil {
		return RunBuildCheckResponse{}, err
	}
	c, ok := w.changeSet(cs)
	if !ok {
		return RunBuildCheckResponse{}, Errorf("run %s names change set %s, which is not in the log", rid, cs).WithNext("rhiza run status " + rid)
	}
	git := func(args ...string) (string, error) {
		out, err := ctx.Exec("git", append([]string{"-C", dir}, args...), dir, nil)
		return string(out), err
	}
	m := w.model(cs)
	leaves := buildCheck(m, w.changeSets(), changeLeaves(m, c.Ops), gitCommits(git))
	resp := RunBuildCheckResponse{Run: rid, Change: cs.String(), Passed: true, Leaves: leaves}
	for _, l := range leaves {
		if !l.OK {
			resp.Passed = false
		}
	}
	return resp, nil
}

// buildCheck judges each of leaves over view m, the change sets sets
// (for who created a node) and commits (newest first, with trailers): a
// leaf is OK when some live Evidence with a live verifies link to it was
// created by an actor named, by handle or ID, in no Actor: trailer of a
// commit whose Leaf: trailer names the leaf's handle.
func buildCheck(m Model, sets []changeset.ChangeSet, leaves []ChangeLeaf, commits []check.Commit) []RunBuildLeaf {
	name := func(n id.ID) string {
		if a, ok := m.Node(n); ok && a.Handle != "" {
			return a.Handle
		}
		return n.String()
	}
	out := make([]RunBuildLeaf, 0, len(leaves))
	for _, l := range leaves {
		r := RunBuildLeaf{Leaf: l.ID.String(), Handle: l.Handle, Evidence: []string{}}
		// authorOf is the first commit naming the leaf that actor wrote.
		authorOf := func(actor id.ID) (string, bool) {
			for _, c := range commits {
				if !slices.Contains(c.Leaves, l.Handle) {
					continue
				}
				if slices.Contains(c.Actors, name(actor)) || slices.Contains(c.Actors, actor.String()) {
					return c.Hash, true
				}
			}
			return "", false
		}
		var whys []string
		for _, ln := range m.LinksTo(l.ID, op.LinkVerifies) {
			ev, ok := m.Node(ln.From)
			if !ok || ev.Retired || ev.Kind != op.KindEvidence {
				continue
			}
			r.Evidence = append(r.Evidence, name(ev.ID))
			who, ok := nodeCreator(sets, ev.ID)
			if !ok {
				whys = append(whys, fmt.Sprintf("its Evidence %s was written by an unknown actor", name(ev.ID)))
				continue
			}
			if commit, wrote := authorOf(who); wrote {
				whys = append(whys, fmt.Sprintf("its Evidence %s was written by %s, who wrote its commit %s", name(ev.ID), name(who), shortRev(commit)))
				continue
			}
			r.OK = true
		}
		switch {
		case r.OK:
		case len(whys) == 0:
			r.Why = "no Evidence verifies it"
		default:
			r.Why = strings.Join(whys, "; ")
		}
		out = append(out, r)
	}
	return out
}

// Human reports the check: one line when it passed, else each leaf that
// is not OK and why.
func (r RunBuildCheckResponse) Human(w io.Writer) {
	if r.Passed {
		fmt.Fprintf(w, "build-check passed: %d leaves have Evidence by a verifier\n", len(r.Leaves))
		return
	}
	failed := 0
	for _, l := range r.Leaves {
		if !l.OK {
			failed++
		}
	}
	fmt.Fprintf(w, "build-check failed: %d of %d leaves lack Evidence by a verifier\n\n", failed, len(r.Leaves))
	for _, l := range r.Leaves {
		if !l.OK {
			fmt.Fprintf(w, "- %s: %s\n", l.Handle, l.Why)
		}
	}
	fmt.Fprintf(w, "\nnext: rhiza run status %s\n", r.Run)
}
