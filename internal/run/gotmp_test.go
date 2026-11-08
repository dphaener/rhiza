package run

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const gotmpWF = `
name = "wf"
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

// A command, a wait and an agent step each get GOTMPDIR as the run's own
// folder, which exists as the step starts and is gone once the run ends.
func TestStepsGetRunPrivateGoTmp(t *testing.T) {
	e, r, _, _ := setup(t, gotmpWF)
	want := e.Store.TmpDir(runID)
	if filepath.Base(want) != "gotmp" || filepath.Base(filepath.Dir(want)) != runID {
		t.Fatalf("TmpDir = %q, want <runs>/%s/gotmp", want, runID)
	}
	check := func(step string) {
		if got := envOf(r.envs[step], EnvGoTmp); got != want {
			t.Errorf("%s: GOTMPDIR = %q, want %q", step, got, want)
		}
	}
	r.onCommand = func(step string) {
		check(step)
		if fi, err := os.Stat(want); err != nil || !fi.IsDir() {
			t.Errorf("%s: GOTMPDIR folder missing when the step ran: %v", step, err)
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
	if _, err := os.Stat(filepath.Dir(want)); !os.IsNotExist(err) {
		t.Errorf("run folder still there after the run finished: %v", err)
	}
}

// A failed run leaves no folder behind either.
func TestGoTmpGoneAfterFailedRun(t *testing.T) {
	e, r, _, _ := setup(t, gotmpWF)
	r.codes["build"] = []int{1}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result == Succeeded {
		t.Fatal("run succeeded, want failed")
	}
	if got := envOf(r.envs["build"], EnvGoTmp); got != e.Store.TmpDir(runID) {
		t.Errorf("GOTMPDIR = %q", got)
	}
	if _, err := os.Stat(filepath.Dir(e.Store.TmpDir(runID))); !os.IsNotExist(err) {
		t.Errorf("run folder still there after the run failed: %v", err)
	}
}
