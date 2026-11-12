package api

// A model write made in a main checkout for a change set with no live run is
// also committed, in the same operation, on the short-lived branch
// model/<change set>-<entry> cut from origin/main in a temporary worktree,
// pushed, and given a model-only pull request. A later write to the same
// change set while that pull request is open goes onto the same branch. The
// step never fails the write: the entry is already in the main checkout's
// log, and a failed step warns and leaves it for rhiza model commit-untracked.
//
// The step runs only where there is something to ship to: a repo with no
// origin/main stays local. Under RHIZA_TEST_NOW it runs only when a test
// opts in with RHIZA_TEST_MODEL_BRANCH, so scripts with an origin and no
// fake gh never reach gh.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/store"
)

// modelBranchPrefix starts every model branch's name.
const modelBranchPrefix = "model/"

func branchShort(s string) string {
	if len(s) > 8 {
		return s[:8]
	}
	return s
}

// modelBranchEnabled reports whether the model-branch step may run.
func modelBranchEnabled(ctx *Ctx) bool {
	if ctx.DryRun || ctx.Exec == nil {
		return false
	}
	if ctx.Getenv != nil && ctx.Getenv("RHIZA_TEST_NOW") != "" && ctx.Getenv("RHIZA_TEST_MODEL_BRANCH") == "" {
		return false
	}
	return true
}

// modelBranchStep returns the function that ships a written entry on a model
// branch, or nil when the step cannot apply to this context.
func modelBranchStep(ctx *Ctx) func(entry.Entry) error {
	if !modelBranchEnabled(ctx) {
		return nil
	}
	return func(e entry.Entry) error {
		top, common, ok := gitDirsOf(ctx.Root)
		if !ok || !samePath(filepath.Join(top, ".git"), common) {
			return nil // not a git repo, or a linked worktree
		}
		if _, live, err := run.NewStore(run.Dir(common)).LiveOf(e.ChangeSet.String(), ""); err != nil || live {
			return nil
		}
		if _, err := ctx.Exec("git", []string{"rev-parse", "--verify", "--quiet", "refs/remotes/origin/main"}, top, nil); err != nil {
			return nil // no origin/main: nowhere to ship
		}
		if err := shipOnModelBranch(ctx, top, e); err != nil && ctx.Stderr != nil {
			fmt.Fprintf(ctx.Stderr, "warning: entry %s of change set %s is written but not on a model branch (%v); rhiza health reports it and rhiza model commit-untracked retries\n", e.ID, e.ChangeSet, err)
		}
		return nil
	}
}

func shipOnModelBranch(ctx *Ctx, top string, e entry.Entry) (err error) {
	return shipModelWrite(ctx, top, e.ChangeSet.String(), e.ID.String(), func(root string) error {
		st, err := store.Open(root)
		if err != nil {
			return err
		}
		if err := st.Write(e); err != nil && !isEntryExists(err) {
			return err
		}
		return nil
	})
}

// shipModelWrite puts a write on change set cs's model branch. tag names a
// new branch (model/<cs>-<tag>); an open model PR of cs takes the write on
// its branch instead. write makes the write in the temporary worktree's
// project root.
func shipModelWrite(ctx *Ctx, top, cs, tag string, write func(root string) error) (err error) {
	git := func(dir string, args ...string) (string, error) {
		out, err := ctx.Exec("git", args, dir, nil)
		return strings.TrimSpace(string(out)), err
	}
	branch := modelBranchPrefix + cs + "-" + tag
	base := "origin/main"
	reuse := false
	if _, err := git(top, "fetch", "-q", "origin", "main"); err != nil {
		return err
	}
	// An open model PR for this change set takes the write on its branch.
	if out, err := git(top, "ls-remote", "--heads", "origin", "refs/heads/"+modelBranchPrefix+cs+"-*"); err == nil {
		for _, line := range strings.Split(out, "\n") {
			_, ref, ok := strings.Cut(line, "\t")
			name := strings.TrimPrefix(ref, "refs/heads/")
			if !ok || name == ref {
				continue
			}
			v, err := gh(ctx, top, "pr", "view", name, "--json", "number,url,state,headRefOid,body,statusCheckRollup")
			if err != nil {
				continue
			}
			if pv, err := parsePRView(v); err == nil && strings.EqualFold(pv.State, "OPEN") {
				branch, reuse = name, true
				break
			}
		}
	}
	if reuse {
		if _, err := git(top, "fetch", "-q", "origin", branch); err != nil {
			return err
		}
		base = "FETCH_HEAD"
	}

	tmp, err := os.MkdirTemp("", "rhiza-model-*")
	if err != nil {
		return err
	}
	wt := filepath.Join(tmp, "wt")
	defer func() {
		_, _ = git(top, "worktree", "remove", "--force", wt)
		_ = os.RemoveAll(tmp)
		_, _ = git(top, "worktree", "prune")
	}()
	if _, err := git(top, "worktree", "add", "-q", "--detach", wt, base); err != nil {
		return err
	}
	root := wt
	if rel, err := filepath.Rel(top, ctx.Root); err == nil && rel != "." {
		root = filepath.Join(wt, rel)
	}
	if err := write(root); err != nil {
		return err
	}
	// Render's output rides the same commit when the write changed the
	// landed model, so the PR passes the self-check.
	rc := *ctx
	rc.Root, rc.DryRun = root, false
	if _, err := runRender(&rc, nil); err != nil {
		return err
	}
	if _, err := git(wt, "add", "-A", "--", store.Dir); err != nil {
		return err
	}
	if out, _ := git(wt, "status", "--porcelain", "--", store.Dir); out == "" {
		return nil // the branch already holds this entry
	}
	msg := fmt.Sprintf("Record model write %s on change set %s\n\nChange: %s\n", branchShort(tag), branchShort(cs), cs)
	if _, err := git(wt, "commit", "-q", "-m", msg); err != nil {
		return err
	}
	if _, err := git(wt, "push", "-q", "origin", "HEAD:refs/heads/"+branch); err != nil {
		return err
	}
	if reuse {
		return nil
	}
	body, err := os.CreateTemp(tmp, "body-*.md")
	if err != nil {
		return err
	}
	fmt.Fprintf(body, "Change set %s: a model-only write.\n", cs)
	if err := body.Close(); err != nil {
		return err
	}
	title := fmt.Sprintf("Model write on change set %s", branchShort(cs))
	_, err = gh(ctx, top, "pr", "create", "--base", shipMain, "--head", branch, "--title", title, "--body-file", body.Name())
	return err
}

func isEntryExists(err error) bool {
	return err != nil && strings.Contains(err.Error(), store.ErrEntryExists.Error())
}
