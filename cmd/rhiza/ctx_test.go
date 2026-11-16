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
	env := []string{"HOME=" + home, "GIT_CONFIG_GLOBAL=" + filepath.Join(home, "gitconfig"), "PATH=" + os.Getenv("PATH"), "RHIZA_TEST_NOW=2026-10-24T12:00:00Z", "RHIZA_TEST_SEED=1"}
	rhiza := func(wd string, args ...string) (int, string, string) {
		t.Helper()
		var stdout, stderr bytes.Buffer
		code := run(args, wd, "", strings.NewReader(""), &stdout, &stderr, env)
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
	out, err := shellWithStdin("sleep 20; echo real: slow", "", os.Environ())
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
	out, err := shellWithStdin("cat", "flake: network", os.Environ())
	if err != nil || out != "flake: network" {
		t.Fatalf("got %q, %v", out, err)
	}
}

// TestChildEnv checks that run() hands the children it starts the
// environment it was given, not the process's: git reads the given HOME's
// GIT_CONFIG_GLOBAL, the index database and the run's temp folder go under
// the given XDG_CACHE_HOME, a workflow step sees the given variables, and a
// variable only the process has stays unset in it.
func TestChildEnv(t *testing.T) {
	proc := t.TempDir()
	t.Setenv("HOME", proc)
	t.Setenv("XDG_CACHE_HOME", filepath.Join(proc, "cache"))
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(proc, "gitconfig"))
	t.Setenv("RHIZA_CHILD_ENV_PROCESS_ONLY", "process")
	writeFile(t, filepath.Join(proc, "gitconfig"), "[user]\n\temail = process@example.test\n\tname = Process\n[rhiza \"ship\"]\n\tlocalGate = process-gate\n")

	home := t.TempDir()
	cache := filepath.Join(home, "given-cache")
	writeFile(t, filepath.Join(home, "gitconfig"), "[user]\n\temail = given@example.test\n\tname = Given\n[rhiza \"ship\"]\n\tlocalGate = given-gate\n")
	env := []string{
		"HOME=" + home,
		"XDG_CACHE_HOME=" + cache,
		"GIT_CONFIG_GLOBAL=" + filepath.Join(home, "gitconfig"),
		"PATH=" + os.Getenv("PATH"),
		"RHIZA_CHILD_ENV_MARKER=given",
		"RHIZA_TEST_NOW=2026-10-24T12:00:00Z",
		"RHIZA_TEST_SEED=1",
	}
	// The repo is project repo of its own Rhiza root, as in the scripts.
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "rhiza.toml"), "[projects.repo]\npath = \"repo\"\n\n[daemon]\n\"rhiza.daemon.failureClassifier\" = \"\"\n")
	repo := filepath.Join(root, "repo")
	if err := os.Mkdir(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
		cmd.Env = env
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	rhiza := func(args ...string) string {
		t.Helper()
		var stdout, stderr bytes.Buffer
		if code := run(args, repo, "", strings.NewReader(""), &stdout, &stderr, env); code != 0 {
			t.Fatalf("rhiza %v: exit %d\n%s%s", args, code, stdout.String(), stderr.String())
		}
		return stdout.String()
	}

	git("init", "-q")
	rhiza("init")
	git("add", "-A")
	git("commit", "-qm", "init")

	// git config, read by a git child, is the given GIT_CONFIG_GLOBAL's.
	settings := rhiza("daemon", "settings", "--json")
	if !strings.Contains(settings, `"given-gate"`) || strings.Contains(settings, "process-gate") {
		t.Errorf("daemon settings read the process's git config, not the given one:\n%s", settings)
	}

	// The index database is under the given XDG_CACHE_HOME.
	rhiza("status", "--json")
	if dbs, _ := filepath.Glob(filepath.Join(cache, "rhiza", "index", "*.db")); len(dbs) == 0 {
		t.Errorf("no index database under the given %s/rhiza/index", cache)
	}
	if _, err := os.Stat(filepath.Join(proc, "cache", "rhiza", "index")); err == nil {
		t.Errorf("an index database was written under the process's cache dir")
	}

	// A workflow step sees the given environment, its run's temp folder
	// under the given XDG_CACHE_HOME, and not the process's variable.
	rhiza("change", "new", "Env", "--title", "Env", "--json")
	out := filepath.Join(t.TempDir(), "step.env")
	writeFile(t, filepath.Join(repo, ".rhiza", "workflows", "childenv.toml"), `name = "childenv"

[env]
FOO = "{run_tmp}"

[step.a]
kind = "command"
run = "echo \"marker=${RHIZA_CHILD_ENV_MARKER-unset} process=${RHIZA_CHILD_ENV_PROCESS_ONLY-unset} home=$HOME foo=$FOO\" > `+out+`"
on.pass = "done"
`)
	if resp := rhiza("run", "start", "--foreground", "--workflow", "childenv", "--json"); !strings.Contains(resp, `"state": "succeeded"`) {
		t.Fatalf("run start: %s", resp)
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(string(b))
	want := "marker=given process=unset home=" + home + " foo=" + filepath.Join(cache, "rhiza", "run-tmp") + string(filepath.Separator)
	if !strings.HasPrefix(got, want) {
		t.Errorf("step env = %q, want it to start %q", got, want)
	}
	if _, err := os.Stat(filepath.Join(proc, "cache", "rhiza", "run-tmp")); err == nil {
		t.Errorf("a run temp folder was made under the process's cache dir")
	}
}

func writeFile(t *testing.T, path, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}
