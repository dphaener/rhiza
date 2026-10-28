package api

// This file holds the checks step's side of the stepRunner: a check
// session started through session.start's path, as a leaf session is,
// with the check and cycle in its record and environment.

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
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
// findings naming no leaf (the cycle's, or in an owner cycle the
// owner's), one per line as a leaf fixer's are, and the
// instruction to fix them in one new commit with the trailer
// `Review: <run>`.
func changeSetFixerPrompt(role workflow.Role, spec run.FixSpec, vars workflow.PromptVars) string {
	p := role.RenderFixerPrompt(workflow.FixerVars{PromptVars: vars, Cycle: spec.Cycle})
	head := fmt.Sprintf("Review cycle %d confirmed these findings", spec.Cycle)
	if spec.Owner {
		head = "The owner requested these changes"
	}
	return strings.TrimRight(p, "\n") + "\n\n" + fmt.Sprintf("%s, which name no leaf of the change set:\n\n%s\n\n"+
		"Fix them, nothing else, in one new commit with the trailer `Review: %s`. `rhiza run status --json` has them too.",
		head, findingLines(spec.Findings), vars.Run) + answerLines(spec.Answers)
}

// OpenCheckCheckout adds a throwaway detached worktree at commit where
// check's session in role on cycle works, as OpenCheckout does a
// verifier's: under the git common dir's rhiza/checkouts, in a folder of
// the run's own (<run>/<check>-<role>-<cycle>), so it never meets a leaf
// verifier's <run>-<leaf>-<attempt>. The run's worktree is never touched.
func (r *stepRunner) OpenCheckCheckout(check, role string, cycle int, commit string) (string, error) {
	return r.openCheckout(filepath.Join(r.run, fmt.Sprintf("%s-%s-%d", check, role, cycle)), commit)
}

// checkQuestionPrefix starts the handle of every question node a capped
// review adds on a finding naming no leaf, or on a check's failure:
// review-failed-<run8>-<n>.
const checkQuestionPrefix = "review-failed-"

// reviewQuestionIntent is the intent of the record change set a capped
// review's question goes in once the run's change set has landed.
const reviewQuestionIntent = "A question review run %s raised"

// AddReviewQuestion writes question node review-failed-<run8>-<n> (n the
// next number free in the change set's view), with statement text and no
// link, into the run's change set in the model as the run's worktree
// holds it, as the run's caller's actor, in one entry, as AddLeafQuestion
// does. Once the run's change set has landed, the question goes in a
// record change set of its own, opened and landed with it. It refuses a
// change set otherwise finished.
func (r *stepRunner) AddReviewQuestion(text string) (run.Question, error) {
	cs, err := id.Parse(r.change)
	if err != nil {
		return run.Question{}, fmt.Errorf("run %s names change set %q: %v", r.run, r.change, err)
	}
	ww, err := r.worktreeWorkspace(true)
	if err != nil {
		return run.Question{}, err
	}
	c, ok := ww.changeSet(cs)
	if !ok {
		return run.Question{}, NotFoundf("change set %s is not in the log", cs)
	}
	landed := c.State == changeset.StateLanded
	if c.State.Terminal() && !landed {
		return run.Question{}, Conflictf("change set %s is %s; no question can be added to it", cs, c.State).
			WithNext("rhiza change show " + cs.String())
	}
	target, ops := cs, c.Ops
	if landed {
		target, ops = id.ID{}, nil
	}
	m := ww.model(target)
	short := strings.ToLower(r.run)
	if len(short) > 8 {
		short = short[:8]
	}
	prefix := checkQuestionPrefix + short + "-"
	handle := ""
	for n := 1; handle == ""; n++ {
		if h := prefix + strconv.Itoa(n); !answerProposed(m, ops, h) {
			handle = h
		}
	}
	actor, forID, err := ww.resolveActor(target, r.ctx)
	if err != nil {
		return run.Question{}, err
	}
	add := AddRequest{Kind: string(op.KindQuestion), Handle: handle, Fields: []string{"statement=" + text}}
	items := []applyItem{{op: "add", build: func(b writeBase) (edit, error) { return buildAdd(b, add) }}}
	e, err := buildApply(ww, target, r.ctx.IDs, items)
	if err == nil {
		err = e.check("add")
	}
	if err != nil {
		return run.Question{}, err
	}
	if len(e.nodes) != 1 {
		return run.Question{}, Errorf("adding question %s touched %d nodes, not 1", handle, len(e.nodes))
	}
	if landed {
		_, _, err = ww.writeRecord(fmt.Sprintf(reviewQuestionIntent, r.run), e.ops, actor, forID)
	} else {
		_, err = ww.appendEntry(cs, changeset.Ops{Ops: e.ops, Drop: e.drop}, actor, forID)
	}
	if err != nil {
		return run.Question{}, err
	}
	return run.Question{ID: e.nodes[0].String(), Handle: handle}, nil
}
