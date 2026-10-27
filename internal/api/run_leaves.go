package api

// This file holds the leaves step's side of the stepRunner: the run change
// set's leaves, read from the model where the run's sessions write, and a
// leaf session started through session.start's path, as an agent step's
// is, with the leaf and attempt in its record and environment, and the
// commits a builder's work is judged by, read with git in the run's
// worktree, and a passed verifier's Evidence entries committed there alone.

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/check"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/store"
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
// model spec names. A spec with a Dir (a verifier's checkout) runs its
// agent there, with RHIZA_ROOT set to the worktree's Rhiza root so its
// rhiza writes reach the build worktree's log. Its tab still opens in the
// worktree: `rhiza session run` writes there too, and the checkout is
// removed once the session ends while the tab may stay.
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
	if spec.Dir != "" {
		req.dir, req.root = spec.Dir, filepath.Join(t.worktree, t.prefix)
	}
	req.Prompt = leafPrompt(*role, spec, workflow.PromptVars{
		Change: t.cs.String(), Slug: session.Slug(t.title), Title: t.title, Run: r.run, Step: spec.Step.Name,
	})
	resp, err := runSessionStart(r.ctx, req)
	if err != nil {
		return "", err
	}
	return resp.(SessionStartResponse).Session.ID, nil
}

// leafPrompt is role's prompt for spec's session, rendered with vars and
// the leaf and attempt. A retry builder's (attempt past 1) ends with a
// line naming the leaf and the attempt and telling it to read the earlier
// attempts' findings from the run and add a fix commit for the leaf.
func leafPrompt(role workflow.Role, spec run.LeafSpec, vars workflow.PromptVars) string {
	p := role.RenderPrompt(workflow.LeafVars{PromptVars: vars, Leaf: spec.Leaf.Handle, Attempt: spec.Attempt})
	if spec.Role != workflow.RoleBuilder || spec.Attempt <= 1 {
		return p
	}
	return strings.TrimRight(p, "\n") + fmt.Sprintf("\n\nThis is attempt %d of %d on leaf %s: earlier attempts failed. "+
		"Read their findings with `rhiza run status --json`, then add a fix commit with the trailer `Leaf: %s`.",
		spec.Attempt, run.MaxAttempts, spec.Leaf.Handle, spec.Leaf.Handle)
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

// Evidence reports whether leaf has live Evidence verifying it that the
// actor of verifier session session created, in the run change set's view
// of the model as the run's worktree holds it, where its sessions write.
// Without, detail is "" when no live Evidence verifies the leaf, else
// names the other actors that wrote what does.
func (r *stepRunner) Evidence(leaf run.Leaf, session string) (bool, string, error) {
	cs, err := id.Parse(r.change)
	if err != nil {
		return false, "", fmt.Errorf("run %s names change set %q: %v", r.run, r.change, err)
	}
	lid, err := id.Parse(leaf.ID)
	if err != nil {
		return false, "", fmt.Errorf("leaf %s has ID %q: %v", leaf.Handle, leaf.ID, err)
	}
	rec, err := r.sessions.Read(session)
	if err != nil {
		return false, "", fmt.Errorf("read session %s: %v", session, err)
	}
	w, err := r.worktreeWorkspace(false)
	if err != nil {
		return false, "", err
	}
	m := w.model(cs)
	actor, err := liveActor(m, rec.Actor)
	if err != nil {
		return false, "", fmt.Errorf("session %s's actor %s: %s", session, rec.Actor, AsError(err).Message)
	}
	ok, detail := leafEvidence(m, w.changeSets(), lid, actor)
	return ok, detail, nil
}

// leafEvidence is Evidence's judgement over view m and the change sets
// sets: whether a live Evidence node with a live verifies link to leaf was
// created by actor (the actor stamped on the entry that first wrote its
// create), and, when none was, which actors created the ones that verify
// it ("" when none does).
func leafEvidence(m Model, sets []changeset.ChangeSet, leaf, actor id.ID) (bool, string) {
	name := func(n id.ID) string {
		if a, ok := m.Node(n); ok && a.Handle != "" {
			return a.Handle
		}
		return n.String()
	}
	var others []string
	for _, l := range m.LinksTo(leaf, op.LinkVerifies) {
		ev, ok := m.Node(l.From)
		if !ok || ev.Retired || ev.Kind != op.KindEvidence {
			continue
		}
		who, ok := nodeCreator(sets, ev.ID)
		if ok && who == actor {
			return true, ""
		}
		by := "an unknown actor"
		if ok {
			by = name(who)
		}
		others = append(others, fmt.Sprintf("Evidence %s was written by %s", name(ev.ID), by))
	}
	if len(others) == 0 {
		return false, ""
	}
	return false, strings.Join(others, "; ") + ", not by the verifier's actor " + name(actor)
}

// nodeCreator is the actor stamped on the entry that first wrote node n's
// create, in whichever of sets holds it (model.CreatedBy over the change
// sets' ops rather than a node's history).
func nodeCreator(sets []changeset.ChangeSet, n id.ID) (id.ID, bool) {
	for _, c := range sets {
		for _, o := range c.Ops {
			if o.Type != op.TypeNodeCreate || o.Node != n {
				continue
			}
			if rs := c.History[o.ID]; len(rs) > 0 && !rs[0].Actor.IsZero() {
				return rs[0].Actor, true
			}
			return id.ID{}, false
		}
	}
	return id.ID{}, false
}

// CommitEvidence commits the log entry files holding the Evidence of
// leaf that verifier session session's actor wrote (the entries that
// wrote its create and its verifies link to the leaf), and nothing else,
// on HEAD of the run's worktree, as its own chore(rhiza) commit naming the
// leaf, with Change: and Actor: (the verifier's) trailers and no Leaf:
// trailer. The entries of the verifier actor's own record go in the same
// commit (actorEntries), so the actor the Evidence names is in its tree. It builds the commit in a temporary index, so the builder's
// staged and unstaged edits stay as they were; the real index is only
// brought up to the new HEAD for the committed paths. When those files
// are already committed and unchanged (a resume after the commit), it
// makes no new commit and returns the last one that touched them.
func (r *stepRunner) CommitEvidence(leaf run.Leaf, session string) (string, error) {
	cs, err := id.Parse(r.change)
	if err != nil {
		return "", fmt.Errorf("run %s names change set %q: %v", r.run, r.change, err)
	}
	lid, err := id.Parse(leaf.ID)
	if err != nil {
		return "", fmt.Errorf("leaf %s has ID %q: %v", leaf.Handle, leaf.ID, err)
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
	entries := evidenceEntries(m, w.changeSets(), lid, actor)
	if len(entries) == 0 {
		return "", fmt.Errorf("no Evidence by %s verifies %s", nodeName(m, actor), leaf.Handle)
	}
	for _, e := range actorEntries(w.changeSets(), actor, cs) {
		if !slices.Contains(entries, e) {
			entries = append(entries, e)
		}
	}
	prefix, err := r.worktreeGit("rev-parse", "--show-prefix")
	if err != nil {
		return "", err
	}
	paths := make([]string, 0, len(entries))
	for _, e := range entries {
		paths = append(paths, filepath.ToSlash(filepath.Join(strings.TrimSpace(prefix), store.Dir, store.LogDir, e.cs.String(), e.entry.String()+".json")))
	}
	status, err := r.worktreeGit(append([]string{"status", "--porcelain", "--untracked-files=all", "--"}, paths...)...)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(status) == "" {
		out, err := r.worktreeGit(append([]string{"log", "-1", "--format=%H", "--"}, paths...)...)
		if err != nil {
			return "", err
		}
		if c := strings.TrimSpace(out); c != "" {
			return c, nil
		}
		return "", fmt.Errorf("the Evidence entries of %s are neither committed nor to be committed: %s", leaf.Handle, strings.Join(paths, ", "))
	}
	name := nodeName(m, actor)
	msg := fmt.Sprintf("chore(rhiza): evidence for %s\n\nChange: %s\nActor: %s\n", leaf.Handle, cs, name)
	return r.commitOnly(paths, msg)
}

// commitOnly commits paths, as the worktree holds them, on HEAD with
// message msg, through a temporary index read from HEAD, and moves HEAD
// (its branch) to the new commit only if it has not moved meanwhile. The
// real index is then updated for paths alone, so they show as committed
// while every other staged or unstaged change stays as it was.
func (r *stepRunner) commitOnly(paths []string, msg string) (string, error) {
	old, err := r.Head()
	if err != nil {
		return "", err
	}
	tmp, err := os.MkdirTemp("", "rhiza-index-")
	if err != nil {
		return "", fmt.Errorf("temporary index: %v", err)
	}
	defer os.RemoveAll(tmp)
	env := append(os.Environ(), "GIT_INDEX_FILE="+filepath.Join(tmp, "index"))
	git := func(args ...string) (string, error) {
		out, err := r.ctx.Exec("git", append([]string{"-C", r.worktree}, args...), r.worktree, env)
		return strings.TrimSpace(string(out)), err
	}
	if _, err := git("read-tree", old); err != nil {
		return "", err
	}
	if _, err := git(append([]string{"update-index", "--add", "--"}, paths...)...); err != nil {
		return "", err
	}
	tree, err := git("write-tree")
	if err != nil {
		return "", err
	}
	commit, err := git("commit-tree", tree, "-p", old, "-m", msg)
	if err != nil {
		return "", err
	}
	if _, err := r.worktreeGit("update-ref", "-m", "commit: "+strings.SplitN(msg, "\n", 2)[0], "HEAD", commit, old); err != nil {
		return "", err
	}
	if _, err := r.worktreeGit(append([]string{"update-index", "--add", "--"}, paths...)...); err != nil {
		return "", err
	}
	return commit, nil
}

// actorEntries are the entries that record actor: every entry of the
// record change set whose ops create it, or, when the run's change set
// run creates it, only the entries that wrote that create. Committed
// ones come back too; committing them again changes nothing.
func actorEntries(sets []changeset.ChangeSet, actor, run id.ID) []evidenceEntry {
	for _, c := range sets {
		for _, o := range c.Ops {
			if o.Type != op.TypeNodeCreate || o.Node != actor {
				continue
			}
			var out []evidenceEntry
			if c.ID != run {
				for _, e := range c.Applied {
					out = append(out, evidenceEntry{c.ID, e})
				}
				return out
			}
			for _, rv := range c.History[o.ID] {
				out = append(out, evidenceEntry{c.ID, rv.Entry})
			}
			return out
		}
	}
	return nil
}

// evidenceEntry is one log entry, in change set cs, that wrote an
// Evidence node's create or its verifies link, or recorded its actor.
type evidenceEntry struct {
	cs, entry id.ID
}

// evidenceEntries are the entries that wrote the create of the first
// live Evidence node actor created with a live verifies link to leaf, and
// that link's create (every version of each), in the order first seen and
// without repeats: one entry when an apply wrote both, two when an add and
// a link did. They are none when actor wrote no such Evidence.
func evidenceEntries(m Model, sets []changeset.ChangeSet, leaf, actor id.ID) []evidenceEntry {
	for _, l := range m.LinksTo(leaf, op.LinkVerifies) {
		ev, ok := m.Node(l.From)
		if !ok || ev.Retired || ev.Kind != op.KindEvidence {
			continue
		}
		if who, ok := nodeCreator(sets, ev.ID); !ok || who != actor {
			continue
		}
		var out []evidenceEntry
		for _, c := range sets {
			for _, o := range c.Ops {
				if !(o.Type == op.TypeNodeCreate && o.Node == ev.ID) && !(o.Type == op.TypeLinkCreate && o.Link == l.ID) {
					continue
				}
				for _, rv := range c.History[o.ID] {
					if e := (evidenceEntry{c.ID, rv.Entry}); !slices.Contains(out, e) {
						out = append(out, e)
					}
				}
			}
		}
		return out
	}
	return nil
}

// OpenCheckout adds a throwaway detached worktree at commit for leaf's
// verifier on attempt, under the git common dir's rhiza/checkouts, and
// returns its path. A leftover at that path (a launch that never got
// journalled) is removed first. The build worktree's index and files are
// never touched: a new worktree has its own.
func (r *stepRunner) OpenCheckout(leaf run.Leaf, attempt int, commit string) (string, error) {
	out, err := r.worktreeGit("rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return "", err
	}
	common := strings.TrimSpace(out)
	dir := filepath.Join(common, "rhiza", "checkouts", fmt.Sprintf("%s-%s-%d", r.run, leaf.Handle, attempt))
	if err := os.MkdirAll(filepath.Dir(dir), 0o755); err != nil {
		return "", fmt.Errorf("checkout folder: %v", err)
	}
	if _, err := os.Stat(dir); err == nil {
		if err := r.RemoveCheckout(dir); err != nil {
			return "", err
		}
	}
	if _, err := r.worktreeGit("worktree", "add", "--detach", dir, commit); err != nil {
		return "", err
	}
	return dir, nil
}

// RemoveCheckout removes a checkout OpenCheckout added, uncommitted
// changes and all. One already gone is only pruned from git's list.
func (r *stepRunner) RemoveCheckout(dir string) error {
	if _, err := os.Stat(dir); errors.Is(err, fs.ErrNotExist) {
		_, err := r.worktreeGit("worktree", "prune")
		return err
	}
	_, err := r.worktreeGit("worktree", "remove", "--force", dir)
	return err
}

// worktreeGit runs git in the run's worktree and returns its output.
func (r *stepRunner) worktreeGit(args ...string) (string, error) {
	if r.ctx.Exec == nil {
		return "", fmt.Errorf("this caller runs no programs")
	}
	out, err := r.ctx.Exec("git", append([]string{"-C", r.worktree}, args...), r.worktree, nil)
	return string(out), err
}
