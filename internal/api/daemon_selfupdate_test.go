package api

import (
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

// A root that still sets the removed self_update_project or
// rhiza.daemon.selfUpdate gets a message naming rhiza.daemon.autoUpdate.
func TestDaemonSettingsReportRemovedSelfUpdate(t *testing.T) {
	for name, toml := range map[string]string{
		"project": "self_update_project = \"repo\"\n\n[projects.repo]\npath = \".\"\n",
		"key":     "[projects.repo]\npath = \".\"\n\n[daemon]\n\"rhiza.daemon.selfUpdate\" = \"true\"\n",
	} {
		ctx, _ := daemonRepo(t)
		if err := os.WriteFile(filepath.Join(ctx.Root, "rhiza.toml"), []byte(toml), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := settingsOf(ctx, nil)
		if err == nil || !strings.Contains(err.Error(), "rhiza.daemon.autoUpdate") {
			t.Fatalf("%s: settingsOf error = %v, want one naming rhiza.daemon.autoUpdate", name, err)
		}
	}
	ctx, _ := daemonRepo(t)
	_, err := runDaemonSet(ctx, DaemonSetRequest{Key: "rhiza.daemon.selfUpdate", Value: "true"})
	if err == nil || !strings.Contains(err.Error(), "rhiza.daemon.autoUpdate") {
		t.Fatalf("daemon set selfUpdate error = %v", err)
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
