package api

// This file is the watchdog's recovery of a PR with no CI because it
// conflicts with main (GitHub starts no workflows on a conflicting PR): the
// daemon merges origin/main into the branch in the run's worktree and pushes
// it normally, never rebasing or forcing, and records the push in the run's
// ship state so the next ship check waits for CI on it.

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

// watchdogMerge is what merging main into a PR branch did: Pushed is the
// head pushed, Conflicts the files that stopped the merge (it was aborted
// and the branch is unchanged). Both are empty when main was already in.
type watchdogMerge struct {
	Pushed    string
	Conflicts []string
}

// watchdogMergeMain merges origin/main into t's branch, checked out in t's
// worktree, pushes it without force and records the pushed head in the ship
// state. A conflicting merge is aborted, leaving the branch as it was.
func watchdogMergeMain(ctx *Ctx, t shipTarget) (watchdogMerge, error) {
	var res watchdogMerge
	if fi, err := os.Stat(t.worktree); err != nil || !fi.IsDir() {
		return res, fmt.Errorf("the worktree %s is gone", t.worktree)
	}
	git := func(args ...string) (string, error) { return gitAt(ctx, t.worktree, args...) }
	if out, err := git("symbolic-ref", "--short", "HEAD"); err != nil || strings.TrimSpace(out) != t.branch {
		return res, fmt.Errorf("%s is not on %s", t.worktree, t.branch)
	}
	if out, err := git("status", "--porcelain", "--untracked-files=no"); err != nil || len(gitLines(out)) > 0 {
		return res, fmt.Errorf("%s has uncommitted changes", t.worktree)
	}
	if _, err := git("fetch", "-q", "origin"); err != nil {
		return res, err
	}
	// Follow a push from elsewhere; diverged work is not ours to merge.
	if out, err := git("rev-parse", "--verify", "--quiet", "refs/remotes/origin/"+t.branch); err == nil {
		remote := strings.TrimSpace(out)
		if _, err := git("merge-base", "--is-ancestor", remote, "HEAD"); err != nil {
			if _, err := git("merge-base", "--is-ancestor", "HEAD", remote); err != nil {
				return res, fmt.Errorf("origin's %s and the worktree have diverged", t.branch)
			}
			if _, err := git("merge", "--ff-only", remote); err != nil {
				return res, err
			}
		}
	}
	onto, err := git("rev-parse", "--verify", "refs/remotes/origin/"+shipMain)
	if err != nil {
		return res, fmt.Errorf("origin has no %s branch: %w", shipMain, err)
	}
	onto = strings.TrimSpace(onto)
	if _, err := git("merge-base", "--is-ancestor", onto, "HEAD"); err == nil {
		return res, nil
	}
	orig, err := git("rev-parse", "HEAD")
	if err != nil {
		return res, err
	}
	orig = strings.TrimSpace(orig)
	if _, err := git("merge", "--no-edit", "-m", "Merge origin/"+shipMain+" into "+t.branch, onto); err != nil {
		out, derr := git("diff", "--name-only", "--diff-filter=U")
		git("merge", "--abort")
		if conflicts := gitLines(out); derr == nil && len(conflicts) > 0 {
			res.Conflicts = conflicts
			return res, nil
		}
		return res, err
	}
	if _, err := git("push", "origin", "HEAD:refs/heads/"+t.branch); err != nil {
		git("reset", "-q", "--keep", orig)
		return res, err
	}
	head, err := git("rev-parse", "HEAD")
	if err != nil {
		return res, err
	}
	res.Pushed = strings.TrimSpace(head)
	st, err := readShipState(t.statePath())
	if err != nil {
		return res, err
	}
	st.Pushed, st.Head, st.Checks = res.Pushed, "", nil
	return res, writeShipState(t.statePath(), st)
}

// hasConflictHandling reports whether st's workflow has a step that settles
// a conflicting PR itself: the ship check with --settled rebases it, and
// ship rebase does.
func hasConflictHandling(st run.State) bool {
	def, err := st.ParseDefinition()
	if err != nil {
		return false
	}
	for _, s := range def.Steps {
		for _, cmd := range []string{s.Run, s.Until} {
			if strings.Contains(cmd, "rhiza ship rebase") || strings.Contains(cmd, "rhiza ship check") && strings.Contains(cmd, "--settled") {
				return true
			}
		}
	}
	return false
}

// watchdogPRNoCI handles each failed run that ended waiting on its PR's CI
// when the PR has no checks and GitHub reports it CONFLICTING: it merges
// main into the branch and resumes the run with --retry to wait for CI on
// the pushed head. If the merge conflicts it resumes the run so the run's
// own conflict handling takes it, or holds an inbox item when the workflow
// has none or the budget is spent. It returns the memory and holds after,
// and the first error.
func watchdogPRNoCI(pol WatchdogPolicy, journals []run.Journal, live func(id string) bool, mem []daemon.WatchdogRun, d watchdogDeps) ([]daemon.WatchdogRun, []daemon.WatchdogHold, error) {
	var holds []daemon.WatchdogHold
	var first error
	fail := func(err error) {
		if first == nil {
			first = err
		}
	}
	if d.ci == nil || d.merge == nil {
		return mem, nil, nil
	}
	for _, j := range journals {
		st := j.State()
		if live(j.ID) || !st.Ended || st.Result != run.Failed || st.Abandoned() {
			continue
		}
		if !st.Autonomous && pol.AutonomousOnly {
			continue
		}
		step, err := st.FailedStep()
		if err != nil || !waitingOnCI(st, step) {
			continue
		}
		pr, err := d.ci(st)
		if err != nil {
			fail(fmt.Errorf("read the PR of run %s: %w", j.ID, err))
			continue
		}
		if pr.Head == "" || len(pr.Checks) > 0 || !strings.EqualFold(pr.Mergeable, "CONFLICTING") {
			continue
		}
		i := slices.IndexFunc(mem, func(w daemon.WatchdogRun) bool { return w.Run == j.ID })
		if i < 0 {
			mem = append(mem, daemon.WatchdogRun{Run: j.ID})
			i = len(mem) - 1
		}
		m := &mem[i]
		res, err := d.merge(st)
		if err != nil {
			fail(fmt.Errorf("merge main into the PR of run %s: %w", j.ID, err))
			continue
		}
		resume := func() bool {
			if m.Spent >= pol.RetryBudget || d.queued != nil && d.queued(j.ID) {
				return false
			}
			if err := d.resume(j.ID); err != nil {
				fail(fmt.Errorf("resume run %s: %w", j.ID, err))
				return false
			}
			d.charge(m, fmt.Sprintf("step %s: PR has no CI: resumed with --retry", step))
			return true
		}
		switch {
		case len(res.Conflicts) > 0:
			files := strings.Join(res.Conflicts, ", ")
			if hasConflictHandling(st) && resume() {
				d.log("watchdog: run %s step %s: PR has no CI and merging main conflicts in %s: left the branch unchanged, resumed with --retry for the run's conflict handling (budget %d of %d)", j.ID, step, files, m.Spent, pol.RetryBudget)
				continue
			}
			holds = append(holds, daemon.WatchdogHold{
				Run: j.ID, Change: st.Change, Step: step, Since: j.Entries[len(j.Entries)-1].At,
				Detail: fmt.Sprintf("the PR has no CI because it conflicts with main in %s, and merging main in by hand is needed: the workflow has no conflict handling, or the watchdog's retry budget of %d is spent", files, pol.RetryBudget),
			})
			d.log("watchdog: run %s step %s: PR has no CI and merging main conflicts in %s: left to the owner", j.ID, step, files)
		case res.Pushed != "":
			d.log("watchdog: run %s step %s: PR has no CI because it conflicts with main: merged origin/main and pushed %s", j.ID, step, shortSHA(res.Pushed))
			if resume() {
				d.log("watchdog: run %s resumed with --retry to wait for CI on %s (budget %d of %d)", j.ID, shortSHA(res.Pushed), m.Spent, pol.RetryBudget)
			}
		}
	}
	return mem, holds, first
}
