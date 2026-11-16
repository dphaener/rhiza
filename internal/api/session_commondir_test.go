package api

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
)

// countingGit runs git in dir through exec and records every call's args.
type countingGit struct {
	mu    sync.Mutex
	calls [][]string
	run   func(args ...string) (string, error)
}

func newCountingGit(exec func(string, []string, string, []string) ([]byte, error), dir string) *countingGit {
	g := &countingGit{}
	g.run = func(args ...string) (string, error) {
		g.mu.Lock()
		g.calls = append(g.calls, append([]string(nil), args...))
		g.mu.Unlock()
		out, err := exec("git", args, dir, nil)
		return string(out), err
	}
	return g
}

// commonDirCalls counts the calls that asked for the git common dir.
func (g *countingGit) commonDirCalls() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	n := 0
	for _, c := range g.calls {
		if slices.Contains(c, "--git-common-dir") {
			n++
		}
	}
	return n
}

// TestGitCommonDirOnce: two gitCommonDir calls on one Ctx make one git
// call, rev-parse --path-format=absolute --git-common-dir, and return the
// same absolute, cleaned path, from the main checkout and from a worktree;
// orphanLogEntries on that Ctx asks git for the common dir no more.
func TestGitCommonDirOnce(t *testing.T) {
	t.Parallel()
	f := modelBranchFixture(t, true)
	wt := filepath.Join(filepath.Dir(f.root), "repo.wt")
	f.git(f.root, "worktree", "add", "-q", "-b", "cs/wt", wt)
	main, err := filepath.EvalSymlinks(filepath.Join(f.root, ".git"))
	if err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{f.root, wt} {
		t.Run(filepath.Base(dir), func(t *testing.T) {
			g := newCountingGit(f.gh.exec, dir)
			ctx := testCtx(t, dir)
			ctx.Git = g.run
			first, err := gitCommonDir(ctx)
			if err != nil {
				t.Fatal(err)
			}
			second, err := gitCommonDir(ctx)
			if err != nil {
				t.Fatal(err)
			}
			want := []string{"rev-parse", "--path-format=absolute", "--git-common-dir"}
			if len(g.calls) != 1 || !slices.Equal(g.calls[0], want) {
				t.Fatalf("git calls %q, want one %q", g.calls, want)
			}
			if first != second || first != main || !filepath.IsAbs(first) || filepath.Clean(first) != first {
				t.Fatalf("common dirs %q and %q, want %q twice", first, second, main)
			}
			orphanLogEntries(ctx)
			if n := g.commonDirCalls(); n != 1 {
				t.Fatalf("orphanLogEntries asked for the common dir again: %d calls in %q", n, g.calls)
			}
			// A copy moved to another checkout asks for its own.
			c := *ctx
			c.Root = t.TempDir()
			c.Git = func(args ...string) (string, error) { return "/elsewhere/.git\n", nil }
			if d, err := gitCommonDir(&c); err != nil || d != "/elsewhere/.git" {
				t.Fatalf("copy at another root: %q, %v", d, err)
			}
			if d, _ := gitCommonDir(ctx); d != main {
				t.Fatalf("after the copy: %q, want %q", d, main)
			}
		})
	}
}

// TestOrphanLogSymlinkRoot: a main checkout reached through a symlinked
// directory keeps git's resolved common dir, equal to rev-parse
// --absolute-git-dir, so orphanLogEntries still lists an untracked orphan
// log entry.
func TestOrphanLogSymlinkRoot(t *testing.T) {
	t.Parallel()
	f := modelBranchFixture(t, true)
	link := filepath.Join(filepath.Dir(f.root), "link")
	if err := os.Symlink(f.root, link); err != nil {
		t.Fatal(err)
	}
	ctx := testCtx(t, link)
	ctx.Git = newCountingGit(f.gh.exec, link).run
	common, err := gitCommonDir(ctx)
	if err != nil {
		t.Fatal(err)
	}
	gitDir, err := ctx.Git("rev-parse", "--absolute-git-dir")
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Clean(strings.TrimSpace(gitDir)); common != want {
		t.Fatalf("common dir %q, want --absolute-git-dir's %q", common, want)
	}
	if strings.HasPrefix(common, link+string(filepath.Separator)) {
		t.Fatalf("common dir %q keeps the symlinked root", common)
	}
	orphans := orphanLogEntries(ctx)
	if !slices.ContainsFunc(orphans, func(p string) bool {
		return strings.HasPrefix(p, ".rhiza/log/"+f.cs+"/")
	}) {
		t.Fatalf("orphans %q, want an entry of %s", orphans, f.cs)
	}
}

// TestShareCommonDir: a worktree added from ctx.Root's repository, as the
// pre-push check's clean checkout is, takes ctx.Root's kept common dir, so
// a copy of ctx moved there asks git nothing; a from other than ctx.Root
// shares nothing.
func TestShareCommonDir(t *testing.T) {
	t.Parallel()
	f := modelBranchFixture(t, true)
	g := newCountingGit(f.gh.exec, f.root)
	ctx := testCtx(t, f.root)
	ctx.Git = g.run
	tmp := t.TempDir()
	shareCommonDir(ctx, f.root, tmp)
	want, err := gitCommonDir(ctx)
	if err != nil {
		t.Fatal(err)
	}
	c := *ctx
	c.Root = tmp
	c.Git = func(args ...string) (string, error) { return "", errors.New("asked git in tmp") }
	if d, err := gitCommonDir(&c); err != nil || d != want {
		t.Fatalf("shared common dir %q, %v; want %q", d, err, want)
	}
	if n := g.commonDirCalls(); n != 1 {
		t.Fatalf("%d common dir calls, want 1: %q", n, g.calls)
	}
	other := t.TempDir()
	shareCommonDir(ctx, other, tmp+"-x")
	c.Root = tmp + "-x"
	if _, err := gitCommonDir(&c); err == nil {
		t.Fatal("a from other than ctx.Root shared its common dir")
	}
}
