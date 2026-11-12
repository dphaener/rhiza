package run

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const gotmpWF = `
name = "wf"
[env]
FOO = "{run_tmp}"
PLAIN = "x"
[step.build]
kind = "command"
run = "make"
on.pass = "ci"
on.fail = "failed"

[step.ci]
kind = "wait"
until = "gh pr checks"
interval = "1s"
timeout = "5s"
on.pass = "agent"
on.timeout = "failed"

[step.agent]
kind = "agent"
prompt = "go {change}"
interval = "2s"
on.pass = "done"
`

func envOf(env []string, key string) string {
	for _, kv := range env {
		if v, ok := strings.CutPrefix(kv, key+"="); ok {
			return v
		}
	}
	return ""
}

// A command, a wait and an agent step each get FOO as the run's own
// folder, which exists as the step starts and is gone once the run ends.
func TestStepsGetEnvRunTmp(t *testing.T) {
	e, r, _, _ := setup(t, gotmpWF)
	want := e.Store.TmpDir(runID)
	if want != filepath.Join(os.Getenv("XDG_CACHE_HOME"), "rhiza", "run-tmp", runID) {
		t.Fatalf("TmpDir = %q, want <user cache>/rhiza/run-tmp/%s", want, runID)
	}
	check := func(step string) {
		if got := envOf(r.envs[step], "FOO"); got != want {
			t.Errorf("%s: FOO = %q, want %q", step, got, want)
		}
	}
	r.onCommand = func(step string) {
		check(step)
		if fi, err := os.Stat(want); err != nil || !fi.IsDir() {
			t.Errorf("%s: FOO folder missing when the step ran: %v", step, err)
		}
	}
	r.holds["ci"] = []bool{true}
	r.outcome["S-agent"] = "pass"
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded {
		t.Fatalf("result %q steps %v", st.Result, outcomes(st))
	}
	check("build")
	check("ci")
	check("agent")
	if _, err := os.Stat(want); !os.IsNotExist(err) {
		t.Errorf("run folder still there after the run finished: %v", err)
	}
}

// A failed run leaves no folder behind either.
func TestRunTmpGoneAfterFailedRun(t *testing.T) {
	e, r, _, _ := setup(t, gotmpWF)
	r.codes["build"] = []int{1}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result == Succeeded {
		t.Fatal("run succeeded, want failed")
	}
	if got := envOf(r.envs["build"], "FOO"); got != e.Store.TmpDir(runID) {
		t.Errorf("FOO = %q", got)
	}
	if _, err := os.Stat(e.Store.TmpDir(runID)); !os.IsNotExist(err) {
		t.Errorf("run folder still there after the run failed: %v", err)
	}
}

// A workflow without [env] gets no variable of Rhiza's setting and no folder.
func TestNoEnvNoRunTmp(t *testing.T) {
	e, r, _, _ := setup(t, strings.Replace(gotmpWF, "[env]\nFOO = \"{run_tmp}\"\nPLAIN = \"x\"\n", "", 1))
	r.codes["build"] = []int{1}
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	for _, kv := range r.envs["build"] {
		if strings.HasPrefix(kv, "GOTMPDIR=") || strings.HasPrefix(kv, "FOO=") {
			t.Errorf("unexpected %s", kv)
		}
	}
}

// Without a user cache dir the folder is under the OS temp dir, per user.
func TestTmpDirFallsBackToOSTemp(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", "")
	t.Setenv("HOME", "")
	got := TmpDir("r1")
	want := filepath.Join(os.TempDir(), "rhiza-run-tmp-"+strconv.Itoa(os.Getuid()), "r1")
	if got != want {
		t.Errorf("TmpDir = %q, want %q", got, want)
	}
}

// Ending a run removes its own folder only: the shared run-tmp parent and
// another run's folder stay.
func TestRunTmpRemovalSparesOthers(t *testing.T) {
	e, r, _, _ := setup(t, gotmpWF)
	other := TmpDir("other-run")
	if err := os.MkdirAll(other, 0o755); err != nil {
		t.Fatal(err)
	}
	r.holds["ci"] = []bool{true}
	r.outcome["S-agent"] = "pass"
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(e.Store.TmpDir(runID)); !os.IsNotExist(err) {
		t.Errorf("own folder still there: %v", err)
	}
	if _, err := os.Stat(other); err != nil {
		t.Errorf("another run's folder gone: %v", err)
	}
}
