package api

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
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
	return writeStrayLog(r.t, r.root)
}

func writeStrayLog(t *testing.T, root string) string {
	t.Helper()
	dir := filepath.Join(root, ".rhiza", "log", "01M6WHRFPM0C0X3H6VBFT4R19Y")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(dir, "01M6WHRFPM0C0X3H6VBFT4R19Z.json")
	if err := os.WriteFile(p, []byte("not json\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

// prePushBreaks are the heads the pre-push check fails on, each made in the
// checkout at root (committed, the working tree left clean) with a line
// its failure carries.
var prePushBreaks = []struct {
	name, want string
	apply      func(t *testing.T, root string, git func(...string) string)
}{
	{"a committed non-canonical log entry", "rhiza check: log", func(t *testing.T, root string, git func(...string) string) {
		writeStrayLog(t, root)
		git("add", "-A")
		git("commit", "-q", "-m", "stray log entry")
	}},
	{"a committed stale render", "rhiza check: render", func(t *testing.T, root string, git func(...string) string) {
		files, _ := filepath.Glob(filepath.Join(root, ".rhiza", "rendered", "*.md"))
		if len(files) == 0 {
			t.Fatal("nothing rendered to make stale")
		}
		if err := os.WriteFile(files[0], []byte("stale\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		git("add", "-A")
		git("commit", "-q", "-m", "stale render")
	}},
	{"a failing prePushCheck command", "FAIL generated file", func(t *testing.T, root string, git func(...string) string) {
		git("config", PrePushCheckKey, "echo FAIL generated file; exit 3")
	}},
}

// realPrePushCheck runs the real pre-push check on r whatever checkout a
// fake-git op asks about, counting its runs.
func realPrePushCheck(t *testing.T, r *prePushRepo) *int {
	n := new(int)
	prePushCheck = func(*Ctx, string) error {
		*n++
		return shipPrePushCheck(r.ctx, r.root)
	}
	return n
}

func setCIStability(t *testing.T, common string, unstable bool) {
	t.Helper()
	rec := `{"runs":[{"head":"a","result":"green","time":"x"}]}`
	if unstable {
		rec = `{"runs":[{"head":"a","result":"green","time":"x","flakes":["TestPond"]}]}`
	}
	if err := writeAtomic(ciStabilityPath(common), []byte(rec)); err != nil {
		t.Fatal(err)
	}
}

func noPush(t *testing.T, calls []string) {
	t.Helper()
	for _, c := range calls {
		if strings.HasPrefix(c, "push") {
			t.Errorf("pushed: %q", c)
		}
	}
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
	t.Run("a render fixed only in the working tree still fails", func(t *testing.T) {
		r := newPrePushRepo(t)
		prePushBreaks[1].apply(t, r.root, r.git)
		if _, err := callOp(t, r.ctx, "render", `{}`); err != nil {
			t.Fatal(err)
		}
		if r.git("status", "--porcelain") == "" {
			t.Fatal("rendering again changed nothing in the working tree")
		}
		err := shipPrePushCheck(r.ctx, r.root)
		if err == nil || !strings.Contains(err.Error(), "rhiza check: render") {
			t.Fatalf("err = %v", err)
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

func TestShipPrePushCheckShipPush(t *testing.T) {
	for _, b := range prePushBreaks {
		t.Run(b.name+" pushes nothing", func(t *testing.T) {
			r := newPrePushRepo(t)
			before := r.remote()
			r.commitFile("a.txt", "a\n")
			b.apply(t, r.root, r.git)
			pushed, err := shipPush(r.ctx, "cs/x")
			if err == nil || pushed || !strings.Contains(err.Error(), b.want) || r.remote() != before {
				t.Fatalf("pushed %v, err %v, remote %q was %q", pushed, err, r.remote(), before)
			}
			if c, ok := prePushFailed(err); !ok || c.Head != r.git("rev-parse", "HEAD") || !strings.Contains(strings.Join(c.Lines, "\n"), b.want) {
				t.Errorf("failure = %+v", c)
			}
			if n := r.worktrees(); n != 1 {
				t.Errorf("%d worktrees left", n)
			}
		})
	}
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

// ship open-pr, against a real origin: a failing check pushes nothing,
// opens or edits no pull request, and exits 1 with the failing lines in
// its response; a passing one pushes as before, whatever CI's state.
func TestShipPrePushCheckOpenPR(t *testing.T) {
	setup := func(t *testing.T) (*Ctx, *openPRGH, string, func(...string) string) {
		ctx, g, origin, git := openPRFixture(t, true)
		prePushCheck = shipPrePushCheck
		t.Cleanup(func() { prePushCheck = func(*Ctx, string) error { return nil } })
		if _, err := callOp(t, ctx, "render", `{}`); err != nil {
			t.Fatal(err)
		}
		git("add", "-A")
		git("commit", "-q", "-m", "render")
		return ctx, g, origin, git
	}
	remote := func(g *openPRGH, origin string) string {
		out, _ := g.exec("git", []string{"--git-dir", origin, "for-each-ref", "refs/heads"}, "", nil)
		return string(out)
	}
	for _, b := range prePushBreaks {
		t.Run(b.name, func(t *testing.T) {
			ctx, g, origin, git := setup(t)
			b.apply(t, ctx.Root, git)
			before := remote(g, origin)
			resp, err := callOp(t, ctx, "ship.open-pr", `{}`)
			if err != nil {
				t.Fatal(err)
			}
			r := resp.(ShipOpenPRResponse)
			if r.ExitCode() != 1 || r.Pushed || r.Created || r.PrePushCheck == nil || r.PrePushCheck.Head != git("rev-parse", "HEAD") ||
				!strings.Contains(strings.Join(r.PrePushCheck.Lines, "\n"), b.want) {
				t.Fatalf("open-pr = %+v", r)
			}
			if got := remote(g, origin); got != before {
				t.Errorf("origin moved: %q, was %q", got, before)
			}
			for _, c := range g.calls {
				if strings.HasPrefix(c, "pr create") || strings.HasPrefix(c, "pr edit") {
					t.Errorf("gh %s after a failing check", c)
				}
			}
			if got := human(t, r); !strings.Contains(got, b.want) || !strings.Contains(got, "nothing was pushed") {
				t.Errorf("human = %q", got)
			}
			if n := strings.Count(git("worktree", "list", "--porcelain"), "worktree "); n != 1 {
				t.Errorf("%d worktrees left", n)
			}
		})
	}
	for _, unstable := range []bool{false, true} {
		t.Run("passing pushes, CI unstable "+map[bool]string{false: "no", true: "yes"}[unstable], func(t *testing.T) {
			ctx, g, origin, git := setup(t)
			setCIStability(t, filepath.Join(ctx.Root, ".git"), unstable)
			git("config", PrePushCheckKey, "true")
			resp, err := callOp(t, ctx, "ship.open-pr", `{}`)
			if err != nil {
				t.Fatal(err)
			}
			r := resp.(ShipOpenPRResponse)
			if r.ExitCode() != 0 || !r.Pushed || !r.Created || r.PrePushCheck != nil || !strings.Contains(remote(g, origin), git("rev-parse", "HEAD")) {
				t.Fatalf("open-pr = %+v", r)
			}
		})
	}
}

// ship land: a failing check on the land commit pushes nothing and exits 1
// with the failing lines in its response.
func TestShipPrePushCheckLand(t *testing.T) {
	for _, b := range prePushBreaks {
		t.Run(b.name, func(t *testing.T) {
			f, state, _ := shipLandFixture(t)
			r := newPrePushRepo(t)
			b.apply(t, r.root, r.git)
			runs := realPrePushCheck(t, r)
			resp, err := callOp(t, f.ctx, "ship.land", `{}`)
			if err != nil {
				t.Fatal(err)
			}
			l := resp.(ShipLandResponse)
			if *runs != 1 || l.ExitCode() != 1 || l.Pushed != "" || l.PrePushCheck == nil || !strings.Contains(strings.Join(l.PrePushCheck.Lines, "\n"), b.want) {
				t.Fatalf("land = %+v, runs %d", l, *runs)
			}
			noPush(t, f.git.calls)
			if st, _ := readShipState(state); st.Pushed != "" {
				t.Errorf("state = %+v", st)
			}
			if got := human(t, l); !strings.Contains(got, b.want) || !strings.Contains(got, "nothing was pushed") {
				t.Errorf("human = %q", got)
			}
		})
	}
	t.Run("passing pushes", func(t *testing.T) {
		f, _, _ := shipLandFixture(t)
		runs := realPrePushCheck(t, newPrePushRepo(t))
		resp, err := callOp(t, f.ctx, "ship.land", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		if l := resp.(ShipLandResponse); *runs != 1 || l.ExitCode() != 0 || l.Pushed != "beef000111" || !slices.Contains(f.git.calls, "push") {
			t.Fatalf("land = %+v, runs %d", l, *runs)
		}
	})
}

// ship rebase: a failing check on the merged head pushes nothing, keeps
// the merged head on the branch and records it, and exits 1 with the
// failing lines in its response; ship gate --failed then passes on it, so
// the step goes on to its fixer.
func TestShipPrePushCheckRebase(t *testing.T) {
	for _, b := range prePushBreaks {
		t.Run(b.name, func(t *testing.T) {
			f, state := shipRebaseFixture(t)
			f.git.out[rebaseMerge] = ""
			r := newPrePushRepo(t)
			b.apply(t, r.root, r.git)
			realPrePushCheck(t, r)
			resp, err := callOp(t, f.ctx, "ship.rebase", `{}`)
			if err != nil {
				t.Fatal(err)
			}
			rb := resp.(ShipRebaseResponse)
			if rb.ExitCode() != 1 || rb.Pushed != "" || rb.PrePushCheck == nil || !strings.Contains(strings.Join(rb.PrePushCheck.Lines, "\n"), b.want) ||
				!strings.Contains(rb.Detail, "nothing was pushed") {
				t.Fatalf("rebase = %+v", rb)
			}
			noPush(t, f.git.calls)
			for _, c := range f.git.calls {
				if strings.HasPrefix(c, "reset ") {
					t.Errorf("the merged head was not kept: %q", c)
				}
			}
			if rb.Kept != rb.PrePushCheck.Head {
				t.Errorf("kept %q, the check failed on %q", rb.Kept, rb.PrePushCheck.Head)
			}
			if st, _ := readShipState(state); st.Pushed != "" || st.GateFailed != "" || st.PrePushFailed != rb.PrePushCheck.Head {
				t.Errorf("state = %+v", st)
			}
			if got := human(t, rb); !strings.Contains(got, b.want) || !strings.Contains(got, "nothing was pushed") || !strings.Contains(got, "for the step's fixer") {
				t.Errorf("human = %q", got)
			}
			// The triage after the failed rebase passes on the kept head,
			// so the step goes on to its fixer.
			exec := f.ctx.Exec
			f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
				if name == "git" && strings.Join(args, " ") == "rev-parse HEAD" {
					return []byte(rb.PrePushCheck.Head + "\n"), nil
				}
				return exec(name, args, dir, env)
			}
			g, err := callOp(t, f.ctx, "ship.gate", `{"failed":true}`)
			if err != nil {
				t.Fatal(err)
			}
			if gr := g.(ShipGateResponse); gr.ExitCode() != 0 || gr.Failed != rb.PrePushCheck.Head {
				t.Errorf("ship gate --failed after the rebase = %+v", gr)
			}
		})
	}
}

// ship check --settled on a pull request GitHub reports CONFLICTING merges
// main as ship rebase does; when the pre-push check fails on the merge, the
// wait ends (settled, exit 0) with the merge kept unpushed, rather than
// polling and checking again until it times out, and ship gate --failed
// then passes on the kept head.
func TestShipPrePushCheckSettleConflict(t *testing.T) {
	b := prePushBreaks[0]
	f, state := shipRebaseFixture(t)
	f.git.out[rebaseMerge] = ""
	r := newPrePushRepo(t)
	b.apply(t, r.root, r.git)
	realPrePushCheck(t, r)
	body, _ := json.Marshal("Change set " + currentOf(t, f.ctx).String() + ": Ship it.")
	view := `{"number":7,"url":"https://github.com/o/r/pull/7","state":"OPEN","headRefOid":"abc1234def","mergeable":"CONFLICTING","body":` + string(body) + `,"statusCheckRollup":[]}`
	exec := f.ctx.Exec
	kept := ""
	f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		switch a := strings.Join(args, " "); {
		case name == "gh" && strings.HasPrefix(a, "pr view cs/ship-it"):
			return []byte(view), nil
		case name == "git" && a == "rev-parse HEAD" && kept != "":
			return []byte(kept + "\n"), nil
		}
		return exec(name, args, dir, env)
	}
	resp, err := callOp(t, f.ctx, "ship.check", `{"settled":true}`)
	if err != nil {
		t.Fatal(err)
	}
	c := resp.(ShipCheckResponse)
	if c.ExitCode() != 0 || !c.Settled || !c.Passed || !strings.Contains(c.Detail, "pre-push check failed") || !strings.Contains(c.Detail, "for the step's fixer") {
		t.Fatalf("settled check = %+v", c)
	}
	noPush(t, f.git.calls)
	st, _ := readShipState(state)
	if st.PrePushFailed == "" || st.Pushed != "" {
		t.Fatalf("state = %+v", st)
	}
	kept = st.PrePushFailed
	g, err := callOp(t, f.ctx, "ship.gate", `{"failed":true}`)
	if err != nil {
		t.Fatal(err)
	}
	if gr := g.(ShipGateResponse); gr.ExitCode() != 0 || gr.Failed != kept {
		t.Errorf("ship gate --failed after the settled check = %+v", gr)
	}
}

// With the pre-push check passing, ship rebase pushes whatever CI's state,
// and the local gate runs only while CI is unstable; the check runs both
// times.
func TestShipPrePushCheckWithLocalGate(t *testing.T) {
	for _, unstable := range []bool{false, true} {
		t.Run(map[bool]string{false: "CI stable", true: "CI unstable"}[unstable], func(t *testing.T) {
			f, _ := shipRebaseFixture(t)
			gateLog := filepath.Join(t.TempDir(), "gate.log")
			f.git.out["config --get "+LocalGateKey] = "echo ran >> " + gateLog + "\n"
			setCIStability(t, f.common, unstable)
			f.git.out[rebaseMerge] = ""
			f.git.out[rebasePush] = ""
			r := newPrePushRepo(t)
			r.setCheck("true")
			runs := realPrePushCheck(t, r)
			rb := shipRebase(t, f)
			_, gerr := os.Stat(gateLog)
			gateRan := gerr == nil
			if rb.ExitCode() != 0 || rb.Pushed != "beef000111" || rb.PrePushCheck != nil || *runs != 1 || !slices.Contains(f.git.calls, rebasePush) {
				t.Fatalf("rebase = %+v, check runs %d", rb, *runs)
			}
			if gateRan != unstable || (rb.Gate != nil) != unstable {
				t.Errorf("gate ran %v (resp %+v) with CI unstable %v", gateRan, rb.Gate, unstable)
			}
		})
	}
}

// The watchdog's merge of origin/main, against a real origin: a failing
// check on the merged head pushes nothing and leaves the branch at its
// pre-merge head, with the failing lines in its error; a passing one
// pushes the merge.
func TestShipPrePushCheckWatchdog(t *testing.T) {
	setup := func(t *testing.T) (*prePushRepo, shipTarget) {
		r := newPrePushRepo(t)
		// main moves on with a file; the branch has its own work.
		other := filepath.Join(t.TempDir(), "other")
		gitRun(t, filepath.Dir(other), "clone", "-q", r.origin, other)
		os.WriteFile(filepath.Join(other, "m.txt"), []byte("m\n"), 0o644)
		gitRun(t, other, "add", ".")
		gitRun(t, other, "commit", "-q", "-m", "main moves")
		gitRun(t, other, "push", "-q", "origin", "HEAD:main")
		r.commitFile("a.txt", "a\n")
		cs, _ := id.Parse("01M6WHRFPM0C0X3H6VBFT4R19Y")
		return r, shipTarget{cs: cs, branch: "cs/x", worktree: r.root, main: r.root, common: filepath.Join(r.root, ".git")}
	}
	for _, b := range prePushBreaks {
		t.Run(b.name, func(t *testing.T) {
			r, tg := setup(t)
			b.apply(t, r.root, r.git)
			r.git("push", "-q", "origin", "cs/x")
			before, remote := r.git("rev-parse", "HEAD"), r.remote()
			res, err := watchdogMergeMain(r.ctx, tg, "")
			if err == nil || !strings.Contains(err.Error(), b.want) || res.Pushed != "" {
				t.Fatalf("res %+v, err %v", res, err)
			}
			if got := r.git("rev-parse", "HEAD"); got != before {
				t.Errorf("branch left at %s, want %s", got, before)
			}
			if r.remote() != remote {
				t.Error("origin's branch moved")
			}
			if n := r.worktrees(); n != 1 {
				t.Errorf("%d worktrees left", n)
			}
		})
	}
	t.Run("passing pushes the merge", func(t *testing.T) {
		r, tg := setup(t)
		r.git("push", "-q", "origin", "cs/x")
		r.setCheck("true")
		res, err := watchdogMergeMain(r.ctx, tg, "")
		if err != nil || res.Pushed == "" || !strings.HasPrefix(r.remote(), res.Pushed) {
			t.Fatalf("passing: %+v, %v", res, err)
		}
	})
}

// When the run's workflow names a head build (SHIP_HEAD_BUILD) the check
// runs a build of the head being pushed, as ship land renders with one, not
// this binary's in-process check: a fake build command "builds" a rhiza
// that prints report as its check --json.
func TestShipPrePushCheckHeadBuild(t *testing.T) {
	const build = `go build -o "$SHIP_HEAD_RHIZA" ./cmd/rhiza`
	setup := func(t *testing.T, report string, buildFails bool) (*prePushRepo, *[]string) {
		r := newPrePushRepo(t)
		r.ctx.Getenv = func(key string) string {
			if key == envShipBuild {
				return build
			}
			return ""
		}
		var calls []string
		real := r.ctx.Exec
		r.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
			if name == "sh" && len(args) == 2 && args[1] == build {
				var bin string
				for _, kv := range env {
					if v, ok := strings.CutPrefix(kv, envShipRhiza+"="); ok {
						bin = v
					}
				}
				calls = append(calls, "build -o "+bin)
				if buildFails {
					return []byte("cmd/rhiza/main.go:1: syntax error\n"), errors.New("exit status 1")
				}
				script := "#!/bin/sh\necho '" + report + "'\n"
				if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
					t.Fatal(err)
				}
				return nil, nil
			}
			if filepath.Base(name) == "rhiza" {
				calls = append(calls, "rhiza "+strings.Join(args, " ")+" in "+dir)
			}
			return real(name, args, dir, env)
		}
		return r, &calls
	}
	t.Run("the head's check governs, not this binary's", func(t *testing.T) {
		// This binary would call the render stale; the head's build is the
		// rendering ship land writes, and its check passes.
		r, calls := setup(t, `{"ok":true,"base":"","findings":[]}`, false)
		prePushBreaks[1].apply(t, r.root, r.git)
		if err := shipPrePushCheck(r.ctx, r.root); err != nil {
			t.Fatal(err)
		}
		if len(*calls) != 2 || !strings.HasPrefix((*calls)[0], "build -o /") || !strings.HasPrefix((*calls)[1], "rhiza check --json --skip-reference --branch-newer-format in ") ||
			strings.HasSuffix((*calls)[1], " in "+r.root) {
			t.Errorf("calls %q", *calls)
		}
	})
	t.Run("the head's failing findings fail it", func(t *testing.T) {
		r, _ := setup(t, `{"ok":false,"base":"","findings":[{"severity":"fail","rule":"render","path":".rhiza/rendered/index.md","message":"differs from what render writes"}]}`, false)
		err := shipPrePushCheck(r.ctx, r.root)
		if err == nil || !strings.Contains(err.Error(), "rhiza check: render: .rhiza/rendered/index.md: differs from what render writes") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("a head that does not build fails it", func(t *testing.T) {
		r, _ := setup(t, "", true)
		err := shipPrePushCheck(r.ctx, r.root)
		if err == nil || !strings.Contains(err.Error(), "could not build rhiza from the head") || !strings.Contains(err.Error(), "syntax error") {
			t.Fatalf("err = %v", err)
		}
		if n := r.worktrees(); n != 1 {
			t.Errorf("%d worktrees left", n)
		}
	})
}

// prePushJournal is the pre-push entries the run whose worktree is r's
// checkout holds.
func (r *prePushRepo) prePushJournal(rs *run.Store, id string) []run.PrePushCheck {
	r.t.Helper()
	j, err := rs.Read(id)
	if err != nil {
		r.t.Fatal(err)
	}
	return j.State().PrePushChecks
}

func TestPrePushJournal(t *testing.T) {
	r := newPrePushRepo(t)
	rs := run.NewStore(run.Dir(filepath.Join(r.root, ".git")))
	const id = "01M77PREPUSH0000000000000A"
	first := run.Entry{Type: run.TypeStarted, At: time.Now(), Step: "ship",
		Started: &run.Started{Workflow: "ship-change", Change: "01M5PWGEF85Y185JZ2ZGCNH40N", Worktree: r.root}}
	if err := rs.Create(id, first); err != nil {
		t.Fatal(err)
	}

	r.commitFile("a.txt", "a\n")
	if err := shipPrePushCheck(r.ctx, r.root); err != nil {
		t.Fatal(err)
	}
	passedHead := strings.TrimSpace(r.git("rev-parse", "HEAD"))
	got := r.prePushJournal(rs, id)
	if len(got) != 1 || got[0].Outcome != run.CheckPassed || got[0].Commit != passedHead || got[0].Command != "rhiza check" || got[0].Output != "" {
		t.Fatalf("after a pass: %+v", got)
	}

	r.setCheck("echo FAIL generated file; exit 3")
	if err := shipPrePushCheck(r.ctx, r.root); err == nil {
		t.Fatal("check passed")
	}
	got = r.prePushJournal(rs, id)
	if len(got) != 2 {
		t.Fatalf("after a failure: %+v", got)
	}
	f := got[1]
	if f.Outcome != run.CheckFailed || f.Commit != passedHead || !strings.Contains(f.Command, "exit 3") || !strings.Contains(f.Output, "FAIL generated file") {
		t.Fatalf("failed entry: %+v", f)
	}
	st, _ := rs.Read(id)
	if l, ok := st.State().LatestPrePush(); !ok || l.Outcome != run.CheckFailed {
		t.Fatalf("latest: %+v", l)
	}
}

// The watchdog acts on a run from outside its steps: shipRunEnv lets its
// pre-push check see the workflow's [env] (SHIP_HEAD_BUILD among it) over
// the daemon's own environment.
func TestShipRunEnv(t *testing.T) {
	const def = "name = \"w\"\n[env]\nSHIP_HEAD_RHIZA = \"{run_tmp}/rhiza-head\"\nSHIP_HEAD_BUILD = 'go build -o \"$SHIP_HEAD_RHIZA\" ./cmd/rhiza'\n[step.wait]\nkind = \"command\"\nrun = \"true\"\n"
	j := run.Journal{ID: "r", Entries: []run.Entry{{Type: run.TypeStarted, At: time.Now(), Started: &run.Started{Definition: def}}}}
	ctx := testCtx(t.TempDir())
	ctx.Getenv = func(key string) string { return "own-" + key }
	got := ctxGetenv(shipRunEnv(ctx, j.State(), "/tmp/run"))
	if v := got(envShipBuild); v != `go build -o "$SHIP_HEAD_RHIZA" ./cmd/rhiza` {
		t.Errorf("%s = %q", envShipBuild, v)
	}
	if v := got(envShipRhiza); v != "/tmp/run/rhiza-head" {
		t.Errorf("%s = %q", envShipRhiza, v)
	}
	if v := got("PATH"); v != "own-PATH" {
		t.Errorf("PATH = %q", v)
	}
	plain := run.Journal{ID: "p", Entries: []run.Entry{{Type: run.TypeStarted, At: time.Now(), Started: &run.Started{Definition: "name = \"w\"\n[step.wait]\nkind = \"command\"\nrun = \"true\"\n"}}}}
	if c := shipRunEnv(ctx, plain.State(), "/tmp/run"); c != ctx {
		t.Error("a workflow with no [env] changed the context")
	}
}
