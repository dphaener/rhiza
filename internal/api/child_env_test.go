package api

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
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

// TestChildEnvPath checks that a hook child is found on the PATH of the
// Ctx's Environ, never the process's: a fake sh first on that PATH runs
// the hook, the same sh through a relative entry is refused (exec.ErrDot),
// and with a PATH that holds no sh the hook fails as not found though the
// process's PATH has one.
func TestChildEnvPath(t *testing.T) {
	realSh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("no sh")
	}
	wd := t.TempDir()
	marker := filepath.Join(t.TempDir(), "fake-sh")
	if err := os.MkdirAll(filepath.Join(wd, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	fake := "#!" + realSh + "\necho fake >> " + marker + "\nexec " + realSh + " \"$@\"\n"
	if err := os.WriteFile(filepath.Join(wd, "bin", "sh"), []byte(fake), 0o755); err != nil {
		t.Fatal(err)
	}
	ctx := &Ctx{Root: t.TempDir(), Wd: wd, Environ: []string{"PATH=" + filepath.Join(wd, "bin")}, Getenv: func(string) string { return "" }}
	if out, err := runHookCommandOutput(ctx, "echo ran", nil, nil); err != nil || strings.TrimSpace(out) != "ran" {
		t.Fatalf("hook: %q, %v", out, err)
	}
	if b, err := os.ReadFile(marker); err != nil || strings.TrimSpace(string(b)) != "fake" {
		t.Errorf("the Ctx's PATH's sh did not run the hook (marker %q, %v)", b, err)
	}

	// A relative entry's sh is refused, as os/exec refuses it.
	ctx.Environ = []string{"PATH=bin"}
	if _, err := runHookCommandOutput(ctx, "echo ran", nil, nil); !errors.Is(err, exec.ErrDot) {
		t.Errorf("hook with sh through a relative PATH entry: %v, want exec.ErrDot", err)
	}

	// Failure case: sh is on the process's PATH only.
	ctx.Environ = []string{"PATH=" + t.TempDir()}
	if _, err := runHookCommandOutput(ctx, "echo ran", nil, nil); !errors.Is(err, exec.ErrNotFound) {
		t.Errorf("hook with no sh on the Ctx's PATH: %v, want not found", err)
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
