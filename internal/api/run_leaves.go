package api

// This file holds the leaves step's side of the stepRunner: the run change
// set's leaves, read from the model where the run's sessions write, and a
// leaf session started through session.start's path, as an agent step's
// is, with the leaf and attempt in its record and environment, and the
// commits a builder's work is judged by, read with git in the run's
// worktree.

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/check"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// Leaves are change set ref's leaves in leaf order (changeLeaves), read
// from the model as the run's worktree holds it, or the caller's checkout
// before the worktree exists. Each leaf's DependsOn are the IDs of the
// leaves among them it depends on.
func (r *stepRunner) Leaves(ref string) ([]run.Leaf, error) {
	cs, err := id.Parse(ref)
	if err != nil {
		return nil, fmt.Errorf("run %s names change set %q: %v", r.run, ref, err)
	}
	var w *workspace
	if _, serr := os.Stat(r.worktree); serr == nil {
		w, err = r.worktreeWorkspace(false)
	} else {
		w, err = openCtxWorkspace(r.ctx)
	}
	if err != nil {
		return nil, err
	}
	c, ok := w.changeSet(cs)
	if !ok {
		return nil, fmt.Errorf("change set %s is not in the log", cs)
	}
	m := w.model(cs)
	found := changeLeaves(m, c.Ops)
	in := map[id.ID]bool{}
	for _, l := range found {
		in[l.ID] = true
	}
	out := make([]run.Leaf, 0, len(found))
	for _, l := range found {
		deps := []string{}
		for _, d := range m.DependsOn(l.ID) {
			if in[d] && !slices.Contains(deps, d.String()) {
				deps = append(deps, d.String())
			}
		}
		out = append(out, run.Leaf{ID: l.ID.String(), Handle: l.Handle, DependsOn: deps})
	}
	return out, nil
}

// StartLeaf starts the leaf session spec names as StartAgent starts an
// agent step's: a session on the run's change set, in its own tab, with
// the role's adapter and mode, else the run's, and its prompt rendered
// with {leaf} and {attempt} too. The session record carries the leaf's
// handle, the attempt, the role and a verifier's commit, which its agent
// gets as RHIZA_LEAF, RHIZA_ATTEMPT, RHIZA_ROLE and RHIZA_COMMIT, and the
// model spec names.
func (r *stepRunner) StartLeaf(spec run.LeafSpec) (string, error) {
	role := spec.Step.Role(spec.Role)
	if role == nil {
		return "", fmt.Errorf("step %q has no %s", spec.Step.Name, spec.Role)
	}
	req := agentStepRequest(workflow.Step{Name: spec.Step.Name, Agent: role.Agent, PermissionMode: role.PermissionMode}, r.mode, r.agent)
	req.Change, req.Run, req.autonomous = r.change, r.run, r.autonomous
	req.leaf, req.attempt, req.model = spec.Leaf.Handle, spec.Attempt, spec.Model
	req.role, req.commit = spec.Role, spec.Commit
	t, err := resolveChangeWorktree(r.ctx, r.change, "no agent can work on it")
	if err != nil {
		return "", err
	}
	req.Prompt = role.RenderPrompt(workflow.LeafVars{
		PromptVars: workflow.PromptVars{
			Change: t.cs.String(), Slug: session.Slug(t.title), Title: t.title, Run: r.run, Step: spec.Step.Name,
		},
		Leaf: spec.Leaf.Handle, Attempt: spec.Attempt,
	})
	resp, err := runSessionStart(r.ctx, req)
	if err != nil {
		return "", err
	}
	return resp.(SessionStartResponse).Session.ID, nil
}

// Head is the commit at HEAD of the run's worktree, where builders work.
func (r *stepRunner) Head() (string, error) {
	out, err := r.worktreeGit("rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// LeafCommit is the newest commit in since..HEAD of the run's worktree
// whose Leaf: trailer names leaf's handle, its trailers read as check
// reads them (parseTrailerLog). With none, commit is "" and detail says
// why: no new commit since since, or the new commits name other leaves,
// or none.
func (r *stepRunner) LeafCommit(leaf run.Leaf, since string) (commit, detail string, err error) {
	out, err := r.worktreeGit("log", "--format="+trailerLogFormat, since+"..HEAD")
	if err != nil {
		return "", "", err
	}
	return leafCommit(parseTrailerLog(out), leaf.Handle, since)
}

// leafCommit picks the newest of commits (newest first) whose Leaf:
// trailers name handle, else says why none does.
func leafCommit(commits []check.Commit, handle, since string) (string, string, error) {
	if len(commits) == 0 {
		return "", "no new commit since " + shortRev(since), nil
	}
	var others []string
	for _, c := range commits {
		if slices.Contains(c.Leaves, handle) {
			return c.Hash, "", nil
		}
		for _, l := range c.Leaves {
			if !slices.Contains(others, l) {
				others = append(others, l)
			}
		}
	}
	if len(others) == 0 {
		return "", fmt.Sprintf("%d new commit(s) since %s, none with a Leaf: trailer", len(commits), shortRev(since)), nil
	}
	return "", "new commits name another leaf: " + strings.Join(others, ", "), nil
}

// worktreeGit runs git in the run's worktree and returns its output.
func (r *stepRunner) worktreeGit(args ...string) (string, error) {
	if r.ctx.Exec == nil {
		return "", fmt.Errorf("this caller runs no programs")
	}
	out, err := r.ctx.Exec("git", append([]string{"-C", r.worktree}, args...), r.worktree, nil)
	return string(out), err
}
