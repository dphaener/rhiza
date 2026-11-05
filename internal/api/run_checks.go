package api

// This file holds the checks step's side of the stepRunner: a check
// session started through session.start's path, as a leaf session is,
// with the check and cycle in its record and environment.

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/store"
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
// Rhiza root. The session opens its own tab.
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
// check. The session opens its own tab.
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
	head := fmt.Sprintf("Confirmed findings to fix in review cycle %d", spec.Cycle)
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

// CommitCheckWrites commits, on HEAD of the run's worktree, the log entry
// files there that check or critic session session's actor wrote and no
// commit holds yet (a QA check's Evidence, a model check's discovered-from
// draft change set), with the entries recording that actor (actorEntries),
// and nothing else, as one chore(rhiza) commit "record <check> check
// writes" ("critic writes" for a critic) with Change: and Actor: (the
// session's) trailers and no Leaf: trailer. It builds the commit as
// CommitEvidence does (commitOnly), so other uncommitted edits stay as
// they were and no later fixer sweeps the check's writes into its own
// commit. A session whose actor wrote no entry still has the uncommitted
// entries recording that actor committed, the same way. With nothing
// uncommitted it makes no commit and returns "".
func (r *stepRunner) CommitCheckWrites(session string) (string, error) {
	cs, err := id.Parse(r.change)
	if err != nil {
		return "", fmt.Errorf("run %s names change set %q: %v", r.run, r.change, err)
	}
	rec, err := r.sessions.Read(session)
	if err != nil {
		return "", fmt.Errorf("read session %s: %v", session, err)
	}
	w, err := r.worktreeWorkspace(false)
	if err != nil {
		return "", err
	}
	m := w.model(cs)
	actor, err := liveActor(m, rec.Actor)
	if err != nil {
		return "", fmt.Errorf("session %s's actor %s: %s", session, rec.Actor, AsError(err).Message)
	}
	prefix, err := r.ctx.Git("rev-parse", "--show-prefix")
	if err != nil {
		return "", err
	}
	dir := path.Join(filepath.ToSlash(strings.TrimSpace(prefix)), store.Dir)
	out, err := r.worktreeGit("ls-files", "--others", "--modified", "--exclude-standard", "--full-name", "-z", "--", path.Join(dir, store.LogDir))
	if err != nil {
		return "", err
	}
	var rels []string
	for _, p := range strings.Split(out, "\x00") {
		if rel, ok := strings.CutPrefix(p, dir+"/"); ok {
			rels = append(rels, rel)
		}
	}
	var paths []string
	add := func(rel string) {
		if p := path.Join(dir, rel); !slices.Contains(paths, p) {
			paths = append(paths, p)
		}
	}
	for _, e := range store.ReadEntries(os.DirFS(filepath.Join(r.worktree, filepath.FromSlash(dir))), rels).Entries {
		if e.Actor == actor {
			add(path.Join(store.LogDir, e.ChangeSet.String(), e.ID.String()+".json"))
		}
	}
	wrote := len(paths) > 0
	for _, e := range actorEntries(w.changeSets(), actor, cs) {
		rel := path.Join(store.LogDir, e.cs.String(), e.entry.String()+".json")
		if wrote || slices.Contains(rels, rel) {
			add(rel)
		}
	}
	if len(paths) == 0 {
		return "", nil
	}
	role := rec.Role
	if role == "" {
		role = workflow.RoleCheck
	}
	msg := fmt.Sprintf("chore(rhiza): record %s %s writes\n\nChange: %s\nActor: %s\n", rec.Check, role, cs, nodeName(m, actor))
	return r.commitOnly(paths, msg)
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
	add := AddRequest{Kind: string(op.KindQuestion), Handle: handle, Fields: []string{"statement=" + text}, lenient: true}
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
