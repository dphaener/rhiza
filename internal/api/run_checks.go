package api

// This file holds the checks step's side of the stepRunner: a check
// session started through session.start's path, as a leaf session is,
// with the check and cycle in its record and environment.

import (
	"fmt"
	"path/filepath"

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
