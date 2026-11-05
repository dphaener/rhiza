package api

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/web"
)

type webFake struct {
	addr   string
	closed bool
}

func (f *webFake) Close() error { f.closed = true; return nil }

// The daemon's web listener follows rhiza.web.port poll by poll: default,
// moved, off, back on, unchanged. A fake factory stands in for the socket.
// Each read of the port setting takes its value from the test over a
// channel, and the test checks the listeners only once the daemon waits for
// its next tick, so every check sees exactly one finished poll.
func TestDaemonWebListenerFollowsPort(t *testing.T) {
	var fakes []*webFake
	old := webFactory
	webFactory = func(addr string, _ http.Handler) (web.Listener, error) {
		f := &webFake{addr: addr}
		fakes = append(fakes, f)
		return f, nil
	}
	t.Cleanup(func() { webFactory = old })

	ctx, common := daemonRepo(t)
	ports := make(chan string)
	ctx.Git = func(args ...string) (string, error) {
		if len(args) == 2 && args[0] == "rev-parse" {
			return common + "\n", nil
		}
		if len(args) == 3 && args[0] == "config" && args[2] == web.PortKey {
			if p := <-ports; p != "" {
				return p + "\n", nil
			}
		}
		return "", errors.New("unset")
	}
	ctx.Getenv = func(string) string { return "" }
	ticks := make(chan time.Time)
	oldTicker := daemonTicker
	daemonTicker = func() (<-chan time.Time, func()) { return ticks, func() {} }
	t.Cleanup(func() { daemonTicker = oldTicker })
	stop := make(chan struct{})
	ctx.Interrupts = func() (<-chan struct{}, func()) { return stop, func() {} }
	done := make(chan error, 1)
	go func() {
		_, err := callOp(t, ctx, "daemon", `{"permission_mode":"manual"}`)
		done <- err
	}()
	feed := func(port string) {
		t.Helper()
		select {
		case ports <- port:
		case err := <-done:
			t.Fatalf("daemon stopped: %v", err)
		case <-time.After(10 * time.Second):
			t.Fatal("daemon never read rhiza.web.port")
		}
	}
	tick := func() {
		t.Helper()
		select {
		case ticks <- time.Now():
		case err := <-done:
			t.Fatalf("daemon stopped: %v", err)
		case <-time.After(10 * time.Second):
			t.Fatal("daemon never waited for a tick")
		}
	}
	open := func() []string {
		var out []string
		for _, f := range fakes {
			if !f.closed {
				out = append(out, f.addr)
			}
		}
		return out
	}

	feed("") // the settings the daemon reads before its first poll
	steps := []struct {
		port   string
		starts int
		open   []string
	}{
		{"", 1, []string{"127.0.0.1:7787"}},     // default: started
		{"8080", 2, []string{"127.0.0.1:8080"}}, // moved
		{"off", 2, nil},                         // stopped
		{"7787", 3, []string{"127.0.0.1:7787"}}, // restarted
		{"7787", 3, []string{"127.0.0.1:7787"}}, // unchanged: left alone
	}
	for i, s := range steps {
		feed(s.port)
		// The daemon takes the next tick only once this poll is over.
		tick()
		if len(fakes) != s.starts || !slices.Equal(open(), s.open) {
			var all []string
			for _, f := range fakes {
				all = append(all, fmt.Sprintf("%s closed=%v", f.addr, f.closed))
			}
			t.Fatalf("poll %d (port %q): listeners %v, want %d started and %v open", i, s.port, all, s.starts, s.open)
		}
		if i+1 < len(steps) {
			continue
		}
		// That tick started one more poll: let it read the same port, then stop.
		feed(s.port)
	}
	close(stop)
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("daemon did not stop on interrupt")
	}
	if got := open(); len(got) != 0 {
		t.Fatalf("listeners %v still open after the daemon stopped", got)
	}
}

// --once and RHIZA_TEST_NOW open no listener.
func TestDaemonWebNoListenerOnceOrTestNow(t *testing.T) {
	opened := 0
	old := webFactory
	webFactory = func(string, http.Handler) (web.Listener, error) { opened++; return &webFake{}, nil }
	t.Cleanup(func() { webFactory = old })

	ctx, _ := daemonRepo(t)
	ctx.IndexPath = filepath.Join(t.TempDir(), "index.db")
	ctx.Getenv = func(string) string { return "" }
	if _, err := callOp(t, ctx, "daemon", `{"permission_mode":"manual","once":true}`); err != nil {
		t.Fatal(err)
	}
	ctx, _ = daemonRepo(t)
	ctx.IndexPath = filepath.Join(t.TempDir(), "index.db")
	ctx.Getenv = func(k string) string {
		if k == "RHIZA_TEST_NOW" {
			return "2026-01-01T00:00:00Z"
		}
		return ""
	}
	if err := runDaemonTicks(t, ctx, 1, nil); err != nil {
		t.Fatal(err)
	}
	if opened != 0 {
		t.Fatalf("%d listeners opened", opened)
	}
}
