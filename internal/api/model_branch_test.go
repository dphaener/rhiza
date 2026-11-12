package api

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// modelGH is a fake gh with the pull requests pr create made, by branch.
type modelGH struct {
	exec    func(string, []string, string, []string) ([]byte, error)
	prs     map[string]bool
	creates int
	calls   int
	fail    bool
}

func (g *modelGH) run(name string, args []string, dir string, env []string) ([]byte, error) {
	if name != "gh" {
		return g.exec(name, args, dir, env)
	}
	g.calls++
	if g.fail {
		return nil, errors.New("gh: not authenticated")
	}
	k := strings.Join(args, " ")
	switch {
	case strings.HasPrefix(k, "pr view "):
		b := args[2]
		if !g.prs[b] {
			return nil, errors.New("no pull requests found for branch " + b)
		}
		return []byte(`{"number":7,"url":"https://github.com/o/r/pull/7","state":"OPEN","body":"","headRefOid":"x","statusCheckRollup":[]}`), nil
	case strings.HasPrefix(k, "pr create "):
		for i, a := range args {
			if a == "--head" {
				g.prs[args[i+1]] = true
			}
		}
		g.creates++
		return []byte("https://github.com/o/r/pull/7\n"), nil
	}
	return nil, errors.New("unexpected gh " + k)
}

type modelFixture struct {
	ctx    *Ctx
	gh     *modelGH
	root   string
	origin string
	git    func(dir string, args ...string) string
	cs     string
	stderr *bytes.Buffer
}

// modelBranchFixture is a repo with a bare origin holding main, a change set
// opened with the step off, and a fake gh.
func modelBranchFixture(t *testing.T, withOrigin bool) *modelFixture {
	t.Helper()
	tmp := t.TempDir()
	root := filepath.Join(tmp, "repo")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	ctx := testCtx(root)
	exec := hermeticExec(t, tmp)
	git := func(dir string, args ...string) string {
		t.Helper()
		out, err := exec("git", args, dir, nil)
		if err != nil {
			t.Fatal(err)
		}
		return strings.TrimSpace(string(out))
	}
	origin := filepath.Join(tmp, "origin.git")
	git(tmp, "init", "-q", "--bare", origin)
	git(root, "init", "-q")
	if withOrigin {
		git(root, "remote", "add", "origin", origin)
	}
	if _, err := callOp(t, ctx, "init", `{}`); err != nil {
		t.Fatal(err)
	}
	if _, err := callOp(t, ctx, "render", `{}`); err != nil {
		t.Fatal(err)
	}
	git(root, "add", "-A")
	git(root, "commit", "-q", "-m", "init")
	if withOrigin {
		git(root, "push", "-q", "origin", "main")
		git(root, "fetch", "-q", "origin")
	}
	f := &modelFixture{ctx: ctx, root: root, origin: origin, git: git, stderr: &bytes.Buffer{}}
	// The change set is opened before the step is wired, so its open entry
	// stays a plain untracked file.
	f.cs = newChange(t, ctx, `{"intent":"Shape"}`).ChangeSet.String()
	f.gh = &modelGH{exec: exec, prs: map[string]bool{}}
	ctx.Exec = f.gh.run
	ctx.Stderr = f.stderr
	return f
}

func (f *modelFixture) branches(t *testing.T) []string {
	t.Helper()
	out := f.git(f.root, "ls-remote", "--heads", f.origin, "refs/heads/model/*")
	var bs []string
	for _, l := range strings.Split(out, "\n") {
		if _, ref, ok := strings.Cut(l, "\t"); ok {
			bs = append(bs, strings.TrimPrefix(ref, "refs/heads/"))
		}
	}
	return bs
}

func TestModelBranch(t *testing.T) {
	t.Run("AcceptShipsBranch", func(t *testing.T) {
		f := modelBranchFixture(t, true)
		exec := f.ctx.Exec
		f.ctx.Exec = nil // propose stays a plain write
		if _, err := callOp(t, f.ctx, "propose", `{}`); err != nil {
			t.Fatal(err)
		}
		f.ctx.Exec = exec
		if _, err := callOp(t, f.ctx, "accept", `{}`); err != nil {
			t.Fatal(err)
		}
		bs := f.branches(t)
		if len(bs) != 1 || !strings.HasPrefix(bs[0], "model/"+f.cs+"-") {
			t.Fatalf("branches = %v", bs)
		}
		f.git(f.root, "fetch", "-q", "origin", bs[0])
		files := strings.Fields(f.git(f.root, "diff", "--name-only", "origin/main", "FETCH_HEAD"))
		if len(files) != 1 || !strings.HasPrefix(files[0], ".rhiza/log/") {
			t.Fatalf("branch holds %v, want exactly one log entry", files)
		}
		if f.gh.creates != 1 {
			t.Fatalf("created %d PRs", f.gh.creates)
		}
		// The entry is in the main checkout too, and the temp worktree is gone.
		if wts := f.git(f.root, "worktree", "list"); strings.Count(wts, "\n") != 0 {
			t.Errorf("worktrees left: %s", wts)
		}
		if _, err := os.Stat(filepath.Join(f.root, files[0])); err != nil {
			t.Errorf("main checkout lacks the entry: %v", err)
		}

		// A second write goes onto the same branch and PR.
		before := f.git(f.root, "ls-remote", f.origin, "refs/heads/"+bs[0])
		mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"x"}`)
		if bs2 := f.branches(t); len(bs2) != 1 {
			t.Fatalf("branches after second write = %v", bs2)
		}
		if f.gh.creates != 1 {
			t.Errorf("second write opened a PR: %d", f.gh.creates)
		}
		if after := f.git(f.root, "ls-remote", f.origin, "refs/heads/"+bs[0]); after == before {
			t.Error("second write did not extend the branch")
		}
	})

	t.Run("RecordApplyCarriesRender", func(t *testing.T) {
		f := modelBranchFixture(t, true)
		c := *f.ctx
		c.Record = true
		if _, err := applyFile(t, &c, treeWrites); err != nil {
			t.Fatal(err)
		}
		bs := f.branches(t)
		if len(bs) != 1 || f.gh.creates != 1 {
			t.Fatalf("branches %v, PRs %d", bs, f.gh.creates)
		}
		f.git(f.root, "fetch", "-q", "origin", bs[0])
		files := f.git(f.root, "diff", "--name-only", "origin/main", "FETCH_HEAD")
		if !strings.Contains(files, ".rhiza/rendered/") {
			t.Fatalf("branch lacks the render: %s", files)
		}
		wt := filepath.Join(t.TempDir(), "wt")
		f.git(f.root, "worktree", "add", "-q", "--detach", wt, "FETCH_HEAD")
		cc := testCtx(wt)
		cc.Git = func(args ...string) (string, error) {
			out, err := f.gh.exec("git", args, wt, nil)
			return string(out), err
		}
		resp, err := callOp(t, cc, "check", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		if h := human(t, resp); strings.Contains(h, "stale") {
			t.Errorf("check reports stale render: %s", h)
		}
	})

	t.Run("NoOriginStaysLocal", func(t *testing.T) {
		f := modelBranchFixture(t, false)
		if _, err := callOp(t, f.ctx, "propose", `{}`); err != nil {
			t.Fatal(err)
		}
		if f.gh.calls != 0 || f.stderr.Len() != 0 {
			t.Errorf("gh calls %d, stderr %q", f.gh.calls, f.stderr)
		}
		if out := f.git(f.root, "branch", "--list", "model/*"); out != "" {
			t.Errorf("branch cut: %s", out)
		}
	})

	t.Run("FailingGHWarns", func(t *testing.T) {
		f := modelBranchFixture(t, true)
		f.gh.fail = true
		if _, err := callOp(t, f.ctx, "propose", `{}`); err != nil {
			t.Fatalf("write failed: %v", err)
		}
		if !strings.Contains(f.stderr.String(), "warning") {
			t.Errorf("no warning: %q", f.stderr)
		}
		if wts := f.git(f.root, "worktree", "list"); strings.Count(wts, "\n") != 0 {
			t.Errorf("worktrees left: %s", wts)
		}
	})

	t.Run("DryRunPushesNothing", func(t *testing.T) {
		f := modelBranchFixture(t, true)
		c := *f.ctx
		c.DryRun = true
		if _, err := callOp(t, &c, "propose", `{}`); err != nil {
			t.Fatal(err)
		}
		if bs := f.branches(t); len(bs) != 0 || f.gh.calls != 0 {
			t.Errorf("branches %v, gh calls %d", bs, f.gh.calls)
		}
	})
}

// TestDaemonCommitsUntracked: a daemon tick commits an orphan entry of the
// main checkout onto a model branch with a PR, and the next tick has
// nothing left to do.
func TestDaemonCommitsUntracked(t *testing.T) {
	f := modelBranchFixture(t, true)
	f.ctx.Git = func(args ...string) (string, error) {
		out, err := f.gh.exec("git", args, f.root, nil)
		return string(out), err
	}
	daemonCommitUntracked(f.ctx)
	bs := f.branches(t)
	if len(bs) == 0 || f.gh.creates != len(bs) {
		t.Fatalf("branches %v, PRs %d", bs, f.gh.creates)
	}
	if !strings.HasPrefix(strings.Join(bs, " "), "model/") || !strings.Contains(strings.Join(bs, " "), "model/"+f.cs+"-") {
		t.Errorf("no branch for change set %s: %v", f.cs, bs)
	}
	if f.stderr.Len() != 0 {
		t.Errorf("stderr: %s", f.stderr)
	}
	n := f.gh.creates
	daemonCommitUntracked(f.ctx)
	if f.gh.creates != n || len(f.branches(t)) != len(bs) {
		t.Errorf("second tick shipped again: %d PRs", f.gh.creates)
	}
}
