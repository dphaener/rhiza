package api

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/rootcfg"
)

func TestDaemonHomePrefersWorkingProject(t *testing.T) {
	root := t.TempDir()
	for _, n := range []string{"a", "b"} {
		if err := os.MkdirAll(filepath.Join(root, n, ".rhiza"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, n, ".rhiza", "repo.json"), []byte("{}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(root, "bare"), 0o755); err != nil {
		t.Fatal(err)
	}
	reg := rootcfg.Config{Projects: map[string]rootcfg.Project{
		"a": {Path: "a"}, "b": {Path: "b"}, "bare": {Path: "bare"},
	}}
	for wd, want := range map[string]string{
		filepath.Join(root, "b"):    "b",
		filepath.Join(root, "bare"): "a", // not initialised: the first initialised one
		root:                        "a",
	} {
		name, _, ok := daemonHome(&Ctx{Wd: wd, Root: wd}, root, reg)
		if !ok || name != want {
			t.Errorf("daemonHome from %s = %q, %v; want %q", wd, name, ok, want)
		}
	}
	if projectInitialised(filepath.Join(root, "bare")) {
		t.Error("a bare checkout counts as initialised")
	}
}

// pollGitStarts records each program a poll git call starts, each its own
// git that never exits until killed.
type pollGitStarts struct {
	progs []Program
	procs []*hangProcess
}

func pollGitCtx(t *testing.T, dir string, st *pollGitStarts) *Ctx {
	t.Helper()
	ctx := testCtx(dir)
	ctx.Git = func(...string) (string, error) { t.Fatal("unbounded Git used"); return "", nil }
	ctx.GitConfig = func(string) (string, error) { t.Fatal("unbounded GitConfig used"); return "", nil }
	ctx.Exec = func(string, []string, string, []string) ([]byte, error) {
		t.Fatal("untimed Exec used")
		return nil, nil
	}
	ctx.Start = func(p Program) (Process, error) {
		h := &hangProcess{killed: make(chan struct{})}
		h.r, h.w = io.Pipe()
		st.progs = append(st.progs, p)
		st.procs = append(st.procs, h)
		return h, nil
	}
	return ctx
}

// checkPollGitBounded asserts c's Git and GitConfig each fail after the
// limit with an error naming the command and the limit, each killing its
// own git, started without a terminal prompt.
func checkPollGitBounded(t *testing.T, c *Ctx, dir string, st *pollGitStarts) {
	t.Helper()
	if _, err := c.Git("status"); err == nil || !strings.Contains(err.Error(), "git -C "+dir+" status: timed out after 10ms") {
		t.Fatalf("Git err %v", err)
	}
	if _, err := c.GitConfig("user.name"); err == nil || !strings.Contains(err.Error(), "git -C "+dir+" config --get user.name: timed out after 10ms") {
		t.Fatalf("GitConfig err %v", err)
	}
	if len(st.progs) != 2 {
		t.Fatalf("programs %+v", st.progs)
	}
	for i, h := range st.procs {
		select {
		case <-h.killed:
		default:
			t.Fatalf("git %d (%v) was not killed", i, st.progs[i].Args)
		}
		if env := strings.Join(st.progs[i].Env, " "); env != "GIT_TERMINAL_PROMPT=0" {
			t.Fatalf("git %d env %q", i, env)
		}
	}
}

func TestPollGitBounded(t *testing.T) {
	old := pollGitTimeout
	pollGitTimeout = 10 * time.Millisecond
	t.Cleanup(func() { pollGitTimeout = old })
	dir := t.TempDir()
	var st pollGitStarts
	checkPollGitBounded(t, projectCtx(pollGitCtx(t, dir, &st), dir), dir, &st)

	// A git that exits in time returns its output unchanged.
	ctx := testCtx(dir)
	ctx.Git = func(...string) (string, error) { return "", nil }
	ctx.Exec = func(string, []string, string, []string) ([]byte, error) { return nil, nil }
	ctx.Start = func(Program) (Process, error) { return &fakeProcess{out: strings.NewReader("main\n")}, nil }
	pc := projectCtx(ctx, dir)
	if out, err := pc.Git("branch"); err != nil || out != "main\n" {
		t.Fatalf("out %q err %v", out, err)
	}
	if v, err := pc.GitConfig("user.name"); err != nil || v != "main" {
		t.Fatalf("config %q err %v", v, err)
	}
}

func TestPollGitBoundedHome(t *testing.T) {
	old := pollGitTimeout
	pollGitTimeout = 10 * time.Millisecond
	t.Cleanup(func() { pollGitTimeout = old })
	root := t.TempDir()
	home := filepath.Join(root, "a")
	if err := os.MkdirAll(filepath.Join(home, ".rhiza"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".rhiza", "repo.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	reg := rootcfg.Config{Projects: map[string]rootcfg.Project{"a": {Path: "a"}}}
	var st pollGitStarts
	ctx := pollGitCtx(t, home, &st)
	ctx.Wd = home
	_, hc, ok := daemonHome(ctx, root, reg)
	if !ok {
		t.Fatal("no home")
	}
	checkPollGitBounded(t, hc, home, &st)
}
