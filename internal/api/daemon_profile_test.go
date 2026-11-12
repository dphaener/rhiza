package api

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/rootcfg"
)

// The daemon captures a CPU profile of itself on request: the file is a
// non-empty gzip pprof at --out, polling goes on during and after the
// capture, and a second request while one runs is refused. The daemon loop
// runs in-process with its ticks and the capture's clock under the test's
// control; nothing binds a socket.
func TestDaemonProfile(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, rootcfg.File), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	rootLock, err := daemon.TryLockRoot(root, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	defer rootLock.Unlock()
	common := t.TempDir()
	lock, err := daemon.TryLock(common, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	fire := make(chan time.Time)
	var asked time.Duration
	prof := newDaemonProfiler(daemon.RootStateDir(root))
	prof.After = func(d time.Duration) <-chan time.Time { asked = d; return fire }
	prof.Reset()
	profileWaitEvery = time.Millisecond
	t.Cleanup(func() { profileWaitEvery = 200 * time.Millisecond })

	polls := make(chan int, 100)
	n := 0
	poll := func() error { prof.Poll(); n++; polls <- n; return nil }
	ctx := testCtx(root)
	ctx.Wd = root
	ticks, stop := make(chan time.Time), make(chan struct{})
	done := make(chan error, 1)
	go func() {
		l, err := daemonLoop(ctx, common, lock, poll, func() int { return 0 }, nil, nil, ticks, stop)
		if l != nil {
			l.Unlock()
		}
		done <- err
	}()
	<-polls

	out := filepath.Join(t.TempDir(), "daemon.pprof")
	type result struct {
		resp any
		err  error
	}
	first := make(chan result, 1)
	go func() {
		resp, err := callOp(t, ctx, "daemon.profile", `{"seconds":1,"out":"`+out+`"}`)
		first <- result{resp, err}
	}()
	// The daemon picks the request up at its next poll.
	reqPath := filepath.Join(daemon.RootStateDir(root), profileRequestFile)
	for i := 0; ; i++ {
		if _, err := os.Stat(reqPath); err == nil {
			break
		}
		if i > 2000 {
			t.Fatal("the request file never appeared")
		}
		time.Sleep(time.Millisecond)
	}
	ticks <- time.Now()
	<-polls
	if asked != time.Second {
		t.Fatalf("capture timer asked for %v, want 1s", asked)
	}

	// A second request, while the capture runs, is refused.
	if _, err := callOp(t, ctx, "daemon.profile", `{"seconds":1}`); err == nil || !strings.Contains(err.Error(), "capture is already in progress") {
		t.Fatalf("second request: %v, want a capture in progress error", err)
	}
	// Polling goes on during the capture.
	ticks <- time.Now()
	<-polls

	fire <- time.Now()
	r := <-first
	if r.err != nil {
		t.Fatal(r.err)
	}
	if got := r.resp.(DaemonProfileResponse).Path; got != out {
		t.Fatalf("path %q, want %q", got, out)
	}
	f, err := os.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("the profile is not gzip: %v", err)
	}
	if b, err := io.ReadAll(gz); err != nil || len(b) == 0 {
		t.Fatalf("profile body: %d bytes, %v", len(b), err)
	}

	// Polling goes on afterwards.
	before := n
	ticks <- time.Now()
	<-polls
	ticks <- time.Now()
	<-polls
	if n <= before {
		t.Fatal("the daemon stopped polling after the capture")
	}
	close(stop)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

// With no --out the profile lands in the root's state folder.
func TestDaemonProfileDefaultPath(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, rootcfg.File), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	rootLock, err := daemon.TryLockRoot(root, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	defer rootLock.Unlock()
	prof := newDaemonProfiler(daemon.RootStateDir(root))
	prof.After = func(time.Duration) <-chan time.Time { c := make(chan time.Time, 1); c <- time.Now(); return c }
	profileWaitEvery = time.Millisecond
	t.Cleanup(func() { profileWaitEvery = 200 * time.Millisecond })
	ctx := testCtx(root)
	ctx.Wd = root
	go func() {
		for i := 0; i < 5000; i++ {
			prof.Poll()
			time.Sleep(time.Millisecond)
		}
	}()
	resp, err := callOp(t, ctx, "daemon.profile", `{"seconds":1}`)
	if err != nil {
		t.Fatal(err)
	}
	p := resp.(DaemonProfileResponse).Path
	if filepath.Dir(p) != daemon.RootStateDir(root) {
		t.Fatalf("default path %q is not in the state folder", p)
	}
	if fi, err := os.Stat(p); err != nil || fi.Size() == 0 {
		t.Fatalf("profile %q: %v", p, err)
	}
}

// Without a daemon the operation says so.
func TestDaemonProfileNoDaemon(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, rootcfg.File), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	ctx := testCtx(root)
	ctx.Wd = root
	if _, err := callOp(t, ctx, "daemon.profile", `{}`); err == nil || !strings.Contains(err.Error(), "no daemon") {
		t.Fatalf("got %v, want a no daemon error", err)
	}
}
