package api

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
)

// A start pending at the end of a poll delays the self-update to the end of
// each later poll with none pending.
func TestDaemonLoopSelfUpdateWaitsForPending(t *testing.T) {
	common := t.TempDir()
	lock, err := daemon.TryLock(common, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	polls := 0
	var updatedAt []int
	// Polls 1 and 2 leave a start pending, polls 3 and 4 do not.
	pendingAt := map[int]int{1: 1, 2: 1}
	ctx := testCtx(t.TempDir())
	ticks, stop := make(chan time.Time), make(chan struct{})
	done := make(chan error, 1)
	go func() {
		l, err := daemonLoop(ctx, common, lock,
			func() error { polls++; return nil },
			func() int { return pendingAt[polls] },
			nil, func() { updatedAt = append(updatedAt, polls) }, ticks, stop)
		if l != nil {
			l.Unlock()
		}
		done <- err
	}()
	for range 3 {
		ticks <- time.Now()
	}
	close(stop)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if polls != 4 || len(updatedAt) != 2 || updatedAt[0] != 3 || updatedAt[1] != 4 {
		t.Fatalf("polls %d, self-update at polls %v; want 4 polls, self-update at [3 4]", polls, updatedAt)
	}
}

// selfUpdateRepo is daemonRepo whose git config sets
// rhiza.daemon.selfUpdate to value ("" for unset), whose main checkout is
// /main, and whose origin/main has moved past it. It records every git and
// make run in the main checkout.
func selfUpdateRepo(t *testing.T, value string, named bool) (*Ctx, func() []string) {
	t.Helper()
	ctx, common := daemonRepo(t)
	if named {
		toml := "self_update_project = \"repo\"\n\n[projects.repo]\npath = \".\"\n"
		if err := os.WriteFile(filepath.Join(ctx.Root, "rhiza.toml"), []byte(toml), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var mu sync.Mutex
	var ran []string
	record := func(s string) {
		mu.Lock()
		defer mu.Unlock()
		ran = append(ran, s)
	}
	ctx.Git = func(args ...string) (string, error) {
		switch {
		case len(args) == 2 && args[0] == "rev-parse" && args[1] == "--git-common-dir":
			return common + "\n", nil
		case len(args) == 3 && args[0] == "config" && args[2] == "rhiza.daemon.selfUpdate" && value != "":
			return value + "\n", nil
		}
		return "", errors.New("unexpected git " + strings.Join(args, " "))
	}
	ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		if dir != ctx.Root {
			return nil, errors.New("unexpected " + name)
		}
		record(name + " " + strings.Join(args, " "))
		if name == "make" {
			return nil, nil
		}
		switch args[0] {
		case "rev-parse":
			if args[2] == "HEAD" {
				return []byte("aaa\n"), nil
			}
			return []byte("bbb\n"), nil
		case "merge-base":
			return nil, errors.New("exit status 1")
		case "symbolic-ref":
			return []byte("main\n"), nil
		}
		return nil, nil
	}
	return ctx, func() []string {
		mu.Lock()
		defer mu.Unlock()
		return append([]string(nil), ran...)
	}
}

// With rhiza.daemon.selfUpdate off, unset or false, the daemon neither
// fetches nor runs make update; turned on, the same daemon fetches,
// fast-forwards and runs make update once over several polls.
func TestDaemonSelfUpdateSetting(t *testing.T) {
	start := daemon.Image{Size: 10, Mod: time.Unix(1700000000, 0)}
	old := daemonReexec
	t.Cleanup(func() { daemonReexec = old })
	daemonReexec = func() *daemon.Reexec { return (&reexecFake{img: start}).reexec(start) }

	for _, value := range []string{"", "false"} {
		ctx, ran := selfUpdateRepo(t, value, true)
		if err := runDaemonTicks(t, ctx, 2, nil); err != nil {
			t.Fatal(err)
		}
		if got := ran(); len(got) != 0 {
			t.Fatalf("selfUpdate %q: ran %q, want nothing", value, got)
		}
	}

	// selfUpdate on but no project named: nothing runs.
	ctx, ran := selfUpdateRepo(t, "true", false)
	if err := runDaemonTicks(t, ctx, 2, nil); err != nil {
		t.Fatal(err)
	}
	if got := ran(); len(got) != 0 {
		t.Fatalf("no self_update_project: ran %q, want nothing", got)
	}

	ctx, ran = selfUpdateRepo(t, "true", true)
	if err := runDaemonTicks(t, ctx, 2, nil); err != nil {
		t.Fatal(err)
	}
	got := strings.Join(ran(), "\n")
	for _, want := range []string{"git fetch origin main", "git merge --ff-only bbb", "make update"} {
		if !strings.Contains(got, want) {
			t.Fatalf("selfUpdate true: ran\n%s\nwant %q", got, want)
		}
	}
	if n := strings.Count(got, "make update"); n != 1 {
		t.Fatalf("selfUpdate true: make update ran %d times over 3 polls, want 1", n)
	}
}

// hangProcess is a program whose output never ends until it is killed.
type hangProcess struct {
	r      *io.PipeReader
	w      *io.PipeWriter
	killed chan struct{}
	once   sync.Once
}

func (h *hangProcess) Output() io.Reader { return h.r }
func (h *hangProcess) Wait() (int, error) {
	<-h.killed
	return 137, nil
}
func (h *hangProcess) Kill() {
	h.once.Do(func() { close(h.killed); h.w.Close() })
}

// A self-update fetch that hangs is killed after its timeout, without a
// terminal prompt, and reported as failed.
func TestFetchMainWithinTimesOut(t *testing.T) {
	ctx := testCtx(t.TempDir())
	var got Program
	h := &hangProcess{killed: make(chan struct{})}
	h.r, h.w = io.Pipe()
	ctx.Start = func(p Program) (Process, error) { got = p; return h, nil }
	err := fetchMainWithin(ctx, "/main", 10*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("err %v", err)
	}
	if got.Name != "git" || strings.Join(got.Args, " ") != "fetch origin main" || got.Dir != "/main" ||
		strings.Join(got.Env, " ") != "GIT_TERMINAL_PROMPT=0" {
		t.Fatalf("program %+v", got)
	}
}

func TestFetchMainWithinExitCode(t *testing.T) {
	ctx := testCtx(t.TempDir())
	ctx.Start = (&shStart{codes: map[string]int{}}).start
	if err := fetchMainWithin(ctx, "/main", time.Second); err != nil {
		t.Fatalf("err %v", err)
	}
	ctx.Start = func(Program) (Process, error) { return &fakeProcess{out: strings.NewReader(""), code: 128}, nil }
	if err := fetchMainWithin(ctx, "/main", time.Second); err == nil {
		t.Fatal("exit 128 not an error")
	}
}
