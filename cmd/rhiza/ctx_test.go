package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestRunWorkingDir checks that run() works in the working directory it is
// given, not the process's: from a process outside any repo it initialises
// and reads a scratch repo, its errors land in the stderr it was given, and
// a working directory that does not exist fails with exit code 1.
func TestRunWorkingDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(home, "gitconfig"))
	outside := t.TempDir()
	t.Chdir(outside)
	repo := t.TempDir()
	if out, err := exec.Command("git", "-C", repo, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	env := map[string]string{"HOME": home, "RHIZA_TEST_NOW": "2026-10-24T12:00:00Z", "RHIZA_TEST_SEED": "1"}
	getenv := func(k string) string { return env[k] }
	rhiza := func(wd string, args ...string) (int, string, string) {
		t.Helper()
		var stdout, stderr bytes.Buffer
		code := run(args, wd, "", strings.NewReader(""), &stdout, &stderr, getenv)
		return code, stdout.String(), stderr.String()
	}

	code, out, errs := rhiza(repo, "init")
	if code != 0 || !strings.Contains(out, "Initialised a Rhiza repository in "+repo) {
		t.Fatalf("init: exit %d\n%s%s", code, out, errs)
	}
	if _, err := os.Stat(filepath.Join(repo, ".rhiza")); err != nil {
		t.Errorf("init wrote no .rhiza in the given folder: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outside, ".rhiza")); err == nil {
		t.Errorf("init wrote .rhiza in the process's working directory")
	}
	if code, out, errs := rhiza(repo, "status", "--json"); code != 0 || !strings.Contains(out, "{") || errs != "" {
		t.Errorf("status in the repo: exit %d\n%s%s", code, out, errs)
	}
	if code, out, errs := rhiza(outside, "status"); code != 4 || out != "" || !strings.Contains(errs, "no Rhiza repository in "+outside) {
		t.Errorf("status outside the repo: exit %d\nstdout %q\nstderr %q", code, out, errs)
	}
	if code, out, errs := rhiza(repo, "show", "nosuch"); code != 4 || out != "" || !strings.Contains(errs, `no node or link matches "nosuch"`) {
		t.Errorf("show nosuch: exit %d\nstdout %q\nstderr %q", code, out, errs)
	}

	missing := filepath.Join(outside, "missing")
	if code, out, errs := rhiza(missing, "status"); code != 1 || out != "" || !strings.Contains(errs, "working directory") {
		t.Errorf("missing working directory: exit %d\nstdout %q\nstderr %q", code, out, errs)
	}
}

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
		ctx, err := newCtx(globals{wd: t.TempDir()}, nil, env(dir))
		if err != nil {
			t.Fatal(err)
		}
		if ctx.Home != dir {
			t.Errorf("Home = %q, want %q", ctx.Home, dir)
		}
	})

	t.Run("unset", func(t *testing.T) {
		ctx, err := newCtx(globals{wd: t.TempDir()}, nil, env(""))
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
