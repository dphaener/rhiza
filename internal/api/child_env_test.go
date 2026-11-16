package api

import (
	"os"
	"strings"
	"testing"
)

// TestChildEnv checks that a hook or gate child started through a Ctx gets
// the Ctx's Environ, not the process's: a marker only Environ holds is set
// in the child, and a variable the process has but Environ lacks is unset
// there.
func TestChildEnv(t *testing.T) {
	t.Setenv("RHIZA_CHILD_ENV_PROCESS_ONLY", "process")
	if _, ok := os.LookupEnv("RHIZA_CHILD_ENV_MARKER"); ok {
		t.Fatal("the process already has RHIZA_CHILD_ENV_MARKER")
	}
	env := []string{"PATH=" + os.Getenv("PATH"), "RHIZA_CHILD_ENV_MARKER=given"}
	ctx := &Ctx{Root: t.TempDir(), Environ: env, Getenv: func(string) string { return "" }}
	probe := `echo "marker=${RHIZA_CHILD_ENV_MARKER-unset} process=${RHIZA_CHILD_ENV_PROCESS_ONLY-unset} added=${RHIZA_CHILD_ENV_ADDED-unset}"`

	out, err := runHookCommandOutput(ctx, probe, nil, []string{"RHIZA_CHILD_ENV_ADDED=yes"})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := strings.TrimSpace(out), "marker=given process=unset added=yes"; got != want {
		t.Errorf("hook child: %q, want %q", got, want)
	}

	// The gate's environment, as runShipGate builds it.
	gate := childEnv(ctx, "RHIZA_GATE_HEAD=abc")
	for _, kv := range gate {
		if strings.HasPrefix(kv, "RHIZA_CHILD_ENV_PROCESS_ONLY=") {
			t.Errorf("gate env carries the process's %s", kv)
		}
	}
	if !contains(gate, "RHIZA_CHILD_ENV_MARKER=given") || !contains(gate, "RHIZA_GATE_HEAD=abc") {
		t.Errorf("gate env = %q, want the marker and RHIZA_GATE_HEAD", gate)
	}
	// The head build's environment drops SHIP_HEAD_RHIZA from Environ.
	ctx.Environ = append(env, envShipRhiza+"=/x/rhiza")
	for _, kv := range headChildEnv(ctx) {
		if strings.HasPrefix(kv, envShipRhiza+"=") {
			t.Errorf("head env keeps %s", kv)
		}
	}
	// childEnv never aliases Environ: an addition does not write into it.
	base := make([]string, 1, 4)
	base[0] = "A=1"
	ctx.Environ = base
	_ = childEnv(ctx, "B=2")
	if got := base[:2][1]; got != "" {
		t.Errorf("childEnv wrote %q into Environ's backing array", got)
	}
	// A Ctx with no Environ, a programmatic caller's, gives no environment,
	// never the process's, and never a nil one a child would inherit.
	if got := childEnv(&Ctx{}); got == nil || len(got) != 0 {
		t.Errorf("nil Environ: child env = %q, want empty", got)
	}
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
