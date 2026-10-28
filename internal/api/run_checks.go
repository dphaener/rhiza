package api

// This file holds the checks step's side of the stepRunner: a check
// session started through session.start's path, as a leaf session is,
// with the check and cycle in its record and environment.

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// stepRunner runs a checks step: it is a run.CheckRunner.
var _ run.CheckRunner = (*stepRunner)(nil)

// StartCheck starts the check session spec names as StartLeaf starts a
// leaf session: a session on the run's change set, with the check's
// adapter and mode, else the run's, its prompt the check's rendered with
// {check} and {cycle} too, so it is told only the change set and its
// check. The session record carries the check's name, the cycle and the
// role, which its agent gets as RHIZA_CHECK, RHIZA_CYCLE and RHIZA_ROLE,
// and the model spec names. It works in the run's worktree unless spec
// names a Dir, where it then runs with RHIZA_ROOT set to the worktree's
// Rhiza root. A spec with an Anchor splits that session's pane, in the
// direction the step's split gives; an anchor with no record opens a new
// tab.
func (r *stepRunner) StartCheck(spec run.CheckSpec) (string, error) {
	role := spec.Check.Role
	if spec.Role != workflow.RoleCheck {
		rp := spec.Step.Role(spec.Role)
		if rp == nil {
			return "", fmt.Errorf("step %q has no %s", spec.Step.Name, spec.Role)
		}
		role = *rp
	}
	req := agentStepRequest(workflow.Step{Name: spec.Step.Name, Agent: role.Agent, PermissionMode: role.PermissionMode}, r.mode, r.agent)
	req.Change, req.Run, req.autonomous = r.change, r.run, r.autonomous
	req.check, req.cycle, req.role, req.model = spec.Check.Name, spec.Cycle, spec.Role, spec.Model
	t, err := resolveChangeWorktree(r.ctx, r.change, "no agent can work on it")
	if err != nil {
		return "", err
	}
	if spec.Dir != "" {
		req.dir, req.root = spec.Dir, filepath.Join(t.worktree, t.prefix)
	}
	req.split = spec.Step.Split
	if spec.Anchor != "" {
		if a, err := r.sessions.Read(spec.Anchor); err == nil {
			req.anchor = a.Tab
		}
	}
	req.Prompt = role.RenderCheckPrompt(workflow.CheckVars{
		PromptVars: workflow.PromptVars{
			Change: t.cs.String(), Slug: session.Slug(t.title), Title: t.title, Run: r.run, Step: spec.Step.Name,
		},
		Check: spec.Check.Name, Cycle: spec.Cycle,
	})
	resp, err := runSessionStart(r.ctx, req)
	if err != nil {
		return "", err
	}
	return resp.(SessionStartResponse).Session.ID, nil
}

// StartFix starts the change set's fixer spec names as StartLeaf starts a
// leaf's: a session on the run's change set in the run's worktree, with
// the fixer's adapter and mode, else the run's, its prompt
// changeSetFixerPrompt's. Its record carries the cycle and role fixer,
// which its agent gets as RHIZA_CYCLE and RHIZA_ROLE, and no leaf or
// check. A spec with an Anchor splits that session's pane.
func (r *stepRunner) StartFix(spec run.FixSpec) (string, error) {
	role := spec.Step.Role(workflow.RoleFixer)
	if role == nil {
		return "", fmt.Errorf("step %q has no fixer", spec.Step.Name)
	}
	req := agentStepRequest(workflow.Step{Name: spec.Step.Name, Agent: role.Agent, PermissionMode: role.PermissionMode}, r.mode, r.agent)
	req.Change, req.Run, req.autonomous = r.change, r.run, r.autonomous
	req.cycle, req.role, req.model = spec.Cycle, workflow.RoleFixer, spec.Model
	t, err := resolveChangeWorktree(r.ctx, r.change, "no agent can work on it")
	if err != nil {
		return "", err
	}
	req.split = spec.Step.Split
	if spec.Anchor != "" {
		if a, err := r.sessions.Read(spec.Anchor); err == nil {
			req.anchor = a.Tab
		}
	}
	req.Prompt = changeSetFixerPrompt(*role, spec, workflow.PromptVars{
		Change: t.cs.String(), Slug: session.Slug(t.title), Title: t.title, Run: r.run, Step: spec.Step.Name,
	})
	resp, err := runSessionStart(r.ctx, req)
	if err != nil {
		return "", err
	}
	return resp.(SessionStartResponse).Session.ID, nil
}

// changeSetFixerPrompt is the change set's fixer's prompt: role's fixer
// prompt rendered with vars, {leaf} empty and the cycle, then the
// findings naming no leaf, one per line as a leaf fixer's are, and the
// instruction to fix them in one new commit with the trailer
// `Review: <run>`.
func changeSetFixerPrompt(role workflow.Role, spec run.FixSpec, vars workflow.PromptVars) string {
	p := role.RenderFixerPrompt(workflow.FixerVars{PromptVars: vars, Cycle: spec.Cycle})
	return strings.TrimRight(p, "\n") + "\n\n" + fmt.Sprintf("Review cycle %d confirmed these findings, which name no leaf of the change set:\n\n%s\n\n"+
		"Fix them, nothing else, in one new commit with the trailer `Review: %s`. `rhiza run status --json` has them too.",
		spec.Cycle, findingLines(spec.Findings), vars.Run)
}

// OpenCheckCheckout adds a throwaway detached worktree at commit where
// check's session in role on cycle works, as OpenCheckout does a
// verifier's: under the git common dir's rhiza/checkouts, in a folder of
// the run's own (<run>/<check>-<role>-<cycle>), so it never meets a leaf
// verifier's <run>-<leaf>-<attempt>. The run's worktree is never touched.
func (r *stepRunner) OpenCheckCheckout(check, role string, cycle int, commit string) (string, error) {
	return r.openCheckout(filepath.Join(r.run, fmt.Sprintf("%s-%s-%d", check, role, cycle)), commit)
}
