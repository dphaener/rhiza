package main

import (
	"strings"
	"testing"
	"time"
)

// TestNewCtxHome checks that newCtx passes the home directory from HOME,
// and leaves it empty when there is none.
func TestNewCtxHome(t *testing.T) {
	env := func(home string) func(string) string {
		return func(k string) string {
			if k == "HOME" {
				return home
			}
			return ""
		}
	}

	t.Run("set", func(t *testing.T) {
		dir := t.TempDir()
		ctx, err := newCtx(globals{}, nil, env(dir))
		if err != nil {
			t.Fatal(err)
		}
		if ctx.Home != dir {
			t.Errorf("Home = %q, want %q", ctx.Home, dir)
		}
	})

	t.Run("unset", func(t *testing.T) {
		ctx, err := newCtx(globals{}, nil, env(""))
		if err != nil {
			t.Fatal(err)
		}
		if ctx.Home != "" {
			t.Errorf("Home = %q, want empty", ctx.Home)
		}
	})
}

// TestShellWithStdinTimesOut checks that a classifier that hangs is killed,
// with the processes it started, after classifierTimeout.
func TestShellWithStdinTimesOut(t *testing.T) {
	old := classifierTimeout
	classifierTimeout = 200 * time.Millisecond
	t.Cleanup(func() { classifierTimeout = old })

	start := time.Now()
	out, err := shellWithStdin("sleep 20; echo real: slow", "")
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("err = %v, want a timeout", err)
	}
	if out != "" {
		t.Errorf("out = %q, want none", out)
	}
	if d := time.Since(start); d > 5*time.Second {
		t.Errorf("took %s, want it killed near the timeout", d)
	}
}

// TestShellWithStdin checks that a classifier within the timeout reads its
// stdin and returns its output.
func TestShellWithStdin(t *testing.T) {
	out, err := shellWithStdin("cat", "flake: network")
	if err != nil || out != "flake: network" {
		t.Fatalf("got %q, %v", out, err)
	}
}
