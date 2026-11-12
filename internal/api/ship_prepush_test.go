package api

import (
	"github.com/dphaener/rhiza/internal/id"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// prePushRepo is a real repo with a Rhiza model, an origin it has pushed
// branch cs/x to, and a ctx running real git and programs in it.
type prePushRepo struct {
	t      *testing.T
	root   string
	origin string
	ctx    *Ctx
}

func (r *prePushRepo) git(args ...string) string {
	r.t.Helper()
	return gitRun(r.t, r.root, args...)
}

func gitRun(t *testing.T, dir string, args ...string) string {
	t.Helper()
	c := exec.Command("git", args...)
	c.Dir = dir
	c.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	out, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func newPrePushRepo(t *testing.T) *prePushRepo {
	t.Helper()
	prePushCheck = shipPrePushCheck
	t.Cleanup(func() { prePushCheck = func(*Ctx, string) error { return nil } })
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(home, "gitconfig"))
	t.Setenv("GIT_CONFIG_SYSTEM", os.DevNull)
	base := t.TempDir()
	r := &prePushRepo{t: t, root: filepath.Join(base, "repo"), origin: filepath.Join(base, "origin.git")}
	gitRun(t, base, "init", "-q", "--bare", "-b", "main", r.origin)
	gitRun(t, base, "init", "-q", "-b", "main", r.root)
	r.git("remote", "add", "origin", r.origin)
	ctx := testCtx(r.root)
	ctx.Git = func(args ...string) (string, error) { return gitAt(ctx, r.root, args...) }
	ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		c := exec.Command(name, args...)
		c.Dir = dir
		c.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		var out strings.Builder
		c.Stdout, c.Stderr = &out, &out
		err := c.Run()
		return []byte(out.String()), err
	}
	r.ctx = ctx
	if _, err := callOp(t, ctx, "init", `{}`); err != nil {
		t.Fatal(err)
	}
	if _, err := callOp(t, ctx, "render", `{}`); err != nil {
		t.Fatal(err)
	}
	r.git("add", "-A")
	r.git("commit", "-q", "-m", "base")
	r.git("push", "-q", "origin", "HEAD:main")
	r.git("checkout", "-q", "-b", "cs/x")
	r.git("push", "-q", "-u", "origin", "cs/x")
	return r
}

func (r *prePushRepo) remote() string {
	return gitRun(r.t, r.root, "ls-remote", "origin", "refs/heads/cs/x")
}

// strayLog is a log entry the store does not read as canonical.
func (r *prePushRepo) strayLog() string {
	r.t.Helper()
	dir := filepath.Join(r.root, ".rhiza", "log", "01M6WHRFPM0C0X3H6VBFT4R19Y")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		r.t.Fatal(err)
	}
	p := filepath.Join(dir, "01M6WHRFPM0C0X3H6VBFT4R19Z.json")
	if err := os.WriteFile(p, []byte("not json\n"), 0o644); err != nil {
		r.t.Fatal(err)
	}
	return p
}

func (r *prePushRepo) commitFile(name, text string) {
	r.t.Helper()
	p := filepath.Join(r.root, name)
	os.MkdirAll(filepath.Dir(p), 0o755)
	if err := os.WriteFile(p, []byte(text), 0o644); err != nil {
		r.t.Fatal(err)
	}
	r.git("add", "-A")
	r.git("commit", "-q", "-m", "work "+name)
}

func (r *prePushRepo) worktrees() int {
	return strings.Count(r.git("worktree", "list", "--porcelain"), "worktree ")
}

func (r *prePushRepo) setCheck(cmd string) {
	r.t.Helper()
	r.git("config", PrePushCheckKey, cmd)
}

func TestShipPrePushCheck(t *testing.T) {
	t.Run("a clean head passes and the temporary checkout is removed", func(t *testing.T) {
		r := newPrePushRepo(t)
		r.commitFile("a.txt", "a\n")
		if err := shipPrePushCheck(r.ctx, r.root); err != nil {
			t.Fatal(err)
		}
		if n := r.worktrees(); n != 1 {
			t.Errorf("%d worktrees left, want only the repo", n)
		}
	})
	t.Run("a committed non-canonical log entry fails", func(t *testing.T) {
		r := newPrePushRepo(t)
		r.strayLog()
		r.git("add", "-A")
		r.git("commit", "-q", "-m", "stray")
		err := shipPrePushCheck(r.ctx, r.root)
		if err == nil || !strings.Contains(err.Error(), "rhiza check: log") || !strings.Contains(err.Error(), shortSHA(r.git("rev-parse", "HEAD"))) {
			t.Fatalf("err = %v", err)
		}
		if n := r.worktrees(); n != 1 {
			t.Errorf("%d worktrees left after a failure", n)
		}
	})
	t.Run("a committed stale render fails", func(t *testing.T) {
		r := newPrePushRepo(t)
		files, _ := filepath.Glob(filepath.Join(r.root, ".rhiza", "rendered", "*.md"))
		if len(files) == 0 {
			t.Skip("init renders nothing here")
		}
		if err := os.WriteFile(files[0], []byte("stale\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		r.git("add", "-A")
		r.git("commit", "-q", "-m", "stale")
		err := shipPrePushCheck(r.ctx, r.root)
		if err == nil || !strings.Contains(err.Error(), "rhiza check: render") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("a failing command fails with its lines", func(t *testing.T) {
		r := newPrePushRepo(t)
		r.setCheck("echo FAIL generated file; exit 3")
		err := shipPrePushCheck(r.ctx, r.root)
		if err == nil || !strings.Contains(err.Error(), "FAIL generated file") || !strings.Contains(err.Error(), "exit status 3") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("the command runs in a checkout of the head, not the worktree", func(t *testing.T) {
		r := newPrePushRepo(t)
		r.commitFile("a.txt", "committed\n")
		os.WriteFile(filepath.Join(r.root, "a.txt"), []byte("dirty\n"), 0o644)
		r.setCheck(`test "$(cat a.txt)" = committed && test "$PWD" != ` + r.root)
		if err := shipPrePushCheck(r.ctx, r.root); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("a fix only in the working tree still fails", func(t *testing.T) {
		r := newPrePushRepo(t)
		r.commitFile("flag.txt", "stale\n")
		r.setCheck("grep -q fixed flag.txt")
		os.WriteFile(filepath.Join(r.root, "flag.txt"), []byte("fixed\n"), 0o644)
		if err := shipPrePushCheck(r.ctx, r.root); err == nil {
			t.Fatal("a working-tree-only fix passed")
		}
		// An uncommitted stray log entry beside a bad committed head fails too.
		r.git("checkout", "--", "flag.txt")
		r.strayLog()
		r.git("add", "-A")
		r.git("commit", "-q", "-m", "stray")
		r.git("rm", "-q", "-r", "--cached", ".rhiza/log/01M6WHRFPM0C0X3H6VBFT4R19Y")
		os.RemoveAll(filepath.Join(r.root, ".rhiza", "log", "01M6WHRFPM0C0X3H6VBFT4R19Y"))
		if err := shipPrePushCheck(r.ctx, r.root); err == nil {
			t.Fatal("an uncommitted removal of a committed stray entry passed")
		}
	})
	t.Run("an untracked stray log entry with a clean head passes", func(t *testing.T) {
		r := newPrePushRepo(t)
		r.commitFile("a.txt", "a\n")
		r.strayLog()
		if err := shipPrePushCheck(r.ctx, r.root); err != nil {
			t.Fatal(err)
		}
	})
}

func TestShipPushRunsPrePushCheck(t *testing.T) {
	t.Run("failing pushes nothing", func(t *testing.T) {
		r := newPrePushRepo(t)
		before := r.remote()
		r.commitFile("a.txt", "a\n")
		r.setCheck("echo BAD; exit 1")
		pushed, err := shipPush(r.ctx, "cs/x")
		if err == nil || pushed || !strings.Contains(err.Error(), "BAD") || r.remote() != before {
			t.Fatalf("pushed %v, err %v, remote %q was %q", pushed, err, r.remote(), before)
		}
	})
	t.Run("passing pushes", func(t *testing.T) {
		r := newPrePushRepo(t)
		r.commitFile("a.txt", "a\n")
		r.setCheck("true")
		pushed, err := shipPush(r.ctx, "cs/x")
		if err != nil || !pushed || !strings.HasPrefix(r.remote(), r.git("rev-parse", "HEAD")) {
			t.Fatalf("pushed %v, err %v", pushed, err)
		}
	})
}

func TestShipPrePushCheckWatchdog(t *testing.T) {
	r := newPrePushRepo(t)
	// main moves on with a file; the branch has its own work.
	other := filepath.Join(t.TempDir(), "other")
	gitRun(t, filepath.Dir(other), "clone", "-q", r.origin, other)
	os.WriteFile(filepath.Join(other, "m.txt"), []byte("m\n"), 0o644)
	gitRun(t, other, "add", ".")
	gitRun(t, other, "commit", "-q", "-m", "main moves")
	gitRun(t, other, "push", "-q", "origin", "HEAD:main")
	r.commitFile("a.txt", "a\n")
	r.git("push", "-q", "origin", "cs/x")
	r.setCheck("echo BAD; exit 1")
	cs, _ := id.Parse("01M6WHRFPM0C0X3H6VBFT4R19Y")
	tg := shipTarget{cs: cs, branch: "cs/x", worktree: r.root, main: r.root, common: filepath.Join(r.root, ".git")}
	before, remote := r.git("rev-parse", "HEAD"), r.remote()
	res, err := watchdogMergeMain(r.ctx, tg, "")
	if err == nil || !strings.Contains(err.Error(), "BAD") || res.Pushed != "" {
		t.Fatalf("res %+v, err %v", res, err)
	}
	if got := r.git("rev-parse", "HEAD"); got != before {
		t.Errorf("branch left at %s, want %s", got, before)
	}
	if r.remote() != remote {
		t.Error("origin's branch moved")
	}
	r.setCheck("true")
	if res, err = watchdogMergeMain(r.ctx, tg, ""); err != nil || res.Pushed == "" {
		t.Fatalf("passing: %+v, %v", res, err)
	}
}

// ship rebase pushes nothing when the pre-push check fails on the merged
// head, puts the branch back where it was, and says why in its error.
func TestShipRebasePrePushCheckFails(t *testing.T) {
	f, _ := shipRebaseFixture(t)
	f.git.out[rebaseMerge] = ""
	f.git.out["reset -q --keep beef000111"] = ""
	prePushCheck = func(*Ctx, string) error {
		return Errorf("the pre-push check failed on beef000, so nothing was pushed:\nBAD")
	}
	t.Cleanup(func() { prePushCheck = func(*Ctx, string) error { return nil } })
	_, err := callOp(t, f.ctx, "ship.rebase", `{}`)
	if err == nil || !strings.Contains(err.Error(), "BAD") {
		t.Fatalf("err = %v", err)
	}
	for _, c := range f.git.calls {
		if strings.HasPrefix(c, "push") {
			t.Errorf("pushed: %q", c)
		}
	}
	if !slices.Contains(f.git.calls, "reset -q --keep beef000111") {
		t.Errorf("branch not put back; calls %q", f.git.calls)
	}
}
