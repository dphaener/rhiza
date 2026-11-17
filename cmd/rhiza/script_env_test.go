package main

import (
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

// leakVars are the calling session's variables a run sets in the test
// process's own environment, which no script may see.
var leakVars = []string{
	"RHIZA_CHANGE", "RHIZA_ACTOR", "RHIZA_SESSION", "RHIZA_RUN",
	"SHIP_HEAD_BUILD", "SHIP_HEAD_RHIZA", "REX_SESSION", "REX_BLOCK",
}

// leakMarker is the value TestScriptEnvIsolated gives name.
func leakMarker(name string) string { return "leak-marker-" + name }

// checkEnvIsolated fails ts when env (NAME=value lines) holds a leakVars
// name or marker value, naming each one, or lacks RHIZA_TEST_NOW or
// RHIZA_TEST_SEED; where names the environment in the message.
func checkEnvIsolated(ts *testscript.TestScript, where string, env []string) {
	have := map[string]bool{}
	var leaked []string
	for _, kv := range env {
		n, v, _ := strings.Cut(kv, "=")
		have[n] = true
		for _, leak := range leakVars {
			if n == leak || strings.Contains(v, leakMarker(leak)) {
				leaked = append(leaked, leak+" (in "+kv+")")
			}
		}
	}
	if len(leaked) > 0 {
		ts.Fatalf("%s: leaked from the test process: %s", where, strings.Join(leaked, ", "))
	}
	for _, n := range []string{"RHIZA_TEST_NOW", "RHIZA_TEST_SEED"} {
		if !have[n] {
			ts.Fatalf("%s: %s is missing", where, n)
		}
	}
}

// TestScriptEnvIsolated shows a script isolated from the test process's
// own RHIZA_*, SHIP_* and REX_* variables: neither an exec'd command's
// environment nor scriptEnv (what the in-process rhiza is handed) holds
// them, while setup's RHIZA_TEST_NOW and RHIZA_TEST_SEED are there. The
// script lives in t.TempDir(), so TestScripts never runs it.
func TestScriptEnvIsolated(t *testing.T) {
	for _, n := range leakVars {
		t.Setenv(n, leakMarker(n))
	}
	f := writeFixture(t, "env-isolated-fixture", `
exec sh -c env
env-isolated
`)
	p := scriptParams(t, []string{f})
	p.Cmds["env-isolated"] = func(ts *testscript.TestScript, neg bool, args []string) {
		if neg || len(args) != 0 {
			ts.Fatalf("usage: env-isolated")
		}
		exec := strings.Split(strings.TrimSpace(ts.ReadFile("stdout")), "\n")
		checkEnvIsolated(ts, "exec sh -c env", exec)
		checkEnvIsolated(ts, "scriptEnv", scriptEnv(ts))
	}
	testscript.Run(t, p)
}
