package api

import (
	"slices"
	"strings"
	"testing"
)

// With RHIZA_CPUPROFILE set in the daemon's environment, each driver it
// detaches profiles into the daemon's path with the run's ID added; unset,
// the driver's environment has none.
func TestDaemonCPUProfilePassedToDrivers(t *testing.T) {
	run := func(prof string) Program {
		t.Helper()
		ctx := testCtx(t.TempDir())
		common := t.TempDir()
		ctx.Git = func(args ...string) (string, error) { return common, nil }
		ctx.Self = "rhiza"
		ctx.repoSession = true
		ctx.Getenv = func(k string) string {
			if k == "RHIZA_CPUPROFILE" {
				return prof
			}
			return ""
		}
		var got Program
		ctx.Detach = func(p Program, _ string) error { got = p; return nil }
		if err := daemonResume(ctx, "01RUNID"); err != nil {
			t.Fatal(err)
		}
		return got
	}
	p := run("/dir/daemon.pprof")
	if !slices.Contains(p.Env, "RHIZA_CPUPROFILE=/dir/daemon.pprof.01RUNID") {
		t.Fatalf("env %v", p.Env)
	}
	for _, e := range run("").Env {
		if strings.HasPrefix(e, "RHIZA_CPUPROFILE=") {
			t.Fatalf("unset profile still passed: %v", e)
		}
	}
}
