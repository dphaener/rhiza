package api

import (
	"errors"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/dphaener/rhiza/internal/web"
)

type webFake struct{ closed bool }

func (f *webFake) Close() error { f.closed = true; return nil }

// The daemon's web listener follows rhiza.web.port poll by poll: default,
// moved, off, back on. A fake factory stands in for the socket.
func TestDaemonWebListenerFollowsPort(t *testing.T) {
	var events []string
	old := webFactory
	webFactory = func(addr string, _ http.Handler) (web.Listener, error) {
		events = append(events, "start "+addr)
		return &webFake{}, nil
	}
	t.Cleanup(func() { webFactory = old })

	ctx, common := daemonRepo(t)
	port := ""
	ctx.Git = func(args ...string) (string, error) {
		if len(args) == 2 && args[0] == "rev-parse" {
			return common + "\n", nil
		}
		if len(args) == 3 && args[0] == "config" && args[2] == web.PortKey && port != "" {
			return port + "\n", nil
		}
		return "", errors.New("unset")
	}
	ctx.Getenv = func(string) string { return "" }
	ports := []string{"8080", "off", "7787", "7787"}
	if err := runDaemonTicks(t, ctx, 4, func(i int) { port = ports[i-1] }); err != nil {
		t.Fatal(err)
	}
	want := []string{"start 127.0.0.1:7787", "start 127.0.0.1:8080", "start 127.0.0.1:7787"}
	if len(events) != len(want) {
		t.Fatalf("events %v, want %v", events, want)
	}
	for i := range want {
		if events[i] != want[i] {
			t.Fatalf("events %v, want %v", events, want)
		}
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
