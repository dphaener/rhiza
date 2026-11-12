package api

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

func TestRemoveRunScratchTakesOnlyTheEndedRunsSessions(t *testing.T) {
	home := t.TempDir()
	ctx := &Ctx{Home: home}
	st := session.NewStore(filepath.Join(t.TempDir(), "sessions"))
	scratch := map[string]string{}
	for id, rn := range map[string]string{"s-ended": "run-a", "s-live": "run-b", "s-none": ""} {
		if err := st.Write(session.Record{ID: id, Agent: interviewAgent, Run: rn}); err != nil {
			t.Fatal(err)
		}
		d := filepath.Join(home, "scratch", id)
		if err := os.MkdirAll(filepath.Join(d, "scratchpad"), 0o755); err != nil {
			t.Fatal(err)
		}
		scratch[id] = d
	}
	removeRunScratch(ctx, st, "run-a")
	if _, err := os.Stat(scratch["s-ended"]); !os.IsNotExist(err) {
		t.Fatalf("the ended run's scratch folder stays: %v", err)
	}
	for _, id := range []string{"s-live", "s-none"} {
		if _, err := os.Stat(scratch[id]); err != nil {
			t.Fatalf("%s's scratch folder went: %v", id, err)
		}
	}
}

func TestRemoveWorktreeIndexKeepsOtherIndexes(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CACHE_HOME", filepath.Join(home, "cache"))
	ctx := &Ctx{Home: home}
	base := t.TempDir()
	dbs := map[string]string{}
	for _, n := range []string{"gone", "live"} {
		d := filepath.Join(base, n)
		if err := os.Mkdir(d, 0o755); err != nil {
			t.Fatal(err)
		}
		p, err := index.DefaultPath(d)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		dbs[n] = p
	}
	removeWorktreeIndex(ctx, filepath.Join(base, "gone"))
	if _, err := os.Stat(dbs["gone"]); !os.IsNotExist(err) {
		t.Fatalf("the removed worktree's index stays: %v", err)
	}
	if _, err := os.Stat(dbs["live"]); err != nil {
		t.Fatalf("a live worktree's index went: %v", err)
	}
}

func TestRemoveCheckoutRemovesItsIndex(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CACHE_HOME", filepath.Join(home, "cache"))
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(home, ".gitconfig"))
	repo := filepath.Join(t.TempDir(), "repo")
	gitIn := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	if err := os.Mkdir(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	gitIn(repo, "init", "-q")
	gitIn(repo, "-c", "user.name=T", "-c", "user.email=t@t.invalid", "commit", "-q", "--allow-empty", "-m", "init")
	r := &stepRunner{ctx: &Ctx{Home: home, Exec: func(name string, args []string, dir string, env []string) ([]byte, error) {
		cmd := exec.Command(name, args...)
		cmd.Dir = dir
		return cmd.Output()
	}}, worktree: repo, run: "r1"}
	dir, err := r.OpenCheckout(run.Leaf{Handle: "leaf"}, 1, "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	dbs := map[string]string{}
	for n, d := range map[string]string{"checkout": dir, "worktree": repo} {
		p, err := index.DefaultPath(d)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		dbs[n] = p
	}
	if err := r.RemoveCheckout(dir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("the checkout stays: %v", err)
	}
	if _, err := os.Stat(dbs["checkout"]); !os.IsNotExist(err) {
		t.Fatalf("the removed checkout's index stays: %v", err)
	}
	if _, err := os.Stat(dbs["worktree"]); err != nil {
		t.Fatalf("the run's worktree index went: %v", err)
	}
}
