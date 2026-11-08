package api

import (
	"strings"
	"testing"
)

// gateFixture is a ship fixture whose worktree exists, with gate set (or
// not) and the stability record holding a flake (or not).
func gateFixture(t *testing.T, gate string, unstable bool) *sessionFixture {
	t.Helper()
	f, _, _ := shipFixture(t)
	if gate != "" {
		f.git.out["config --get "+LocalGateKey] = gate + "\n"
	}
	rec := `{"runs":[{"head":"a","result":"green","time":"x"}]}`
	if unstable {
		rec = `{"runs":[{"head":"a","result":"green","time":"x","flakes":["TestPond"]}]}`
	}
	if err := writeAtomic(ciStabilityPath(f.common), []byte(rec)); err != nil {
		t.Fatal(err)
	}
	return f
}

func shipGate(t *testing.T, f *sessionFixture) ShipGateResponse {
	t.Helper()
	r, err := callOp(t, f.ctx, "ship.gate", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	return r.(ShipGateResponse)
}

func TestShipLocalGate(t *testing.T) {
	t.Run("unstable runs the gate in the worktree", func(t *testing.T) {
		f := gateFixture(t, "pwd; echo gate-ran", true)
		r := shipGate(t, f)
		if r.Skipped || !r.Passed || r.ExitCode() != 0 || !strings.Contains(r.Output, f.worktree) || !strings.Contains(r.Output, "gate-ran") {
			t.Errorf("resp = %+v, worktree %s", r, f.worktree)
		}
	})
	t.Run("stable skips", func(t *testing.T) {
		f := gateFixture(t, "echo gate-ran >&2; exit 1", false)
		r := shipGate(t, f)
		if !r.Skipped || r.Output != "" || r.ExitCode() != 0 || r.Reason != "CI is stable" {
			t.Errorf("resp = %+v", r)
		}
	})
	t.Run("no gate skips", func(t *testing.T) {
		f := gateFixture(t, "", true)
		r := shipGate(t, f)
		if !r.Skipped || r.Command != "" || r.ExitCode() != 0 {
			t.Errorf("resp = %+v", r)
		}
	})
	t.Run("failing gate fails with its output", func(t *testing.T) {
		f := gateFixture(t, "echo FAIL TestGarden; exit 3", true)
		r := shipGate(t, f)
		if r.Skipped || r.Passed || r.ExitCode() != 1 || !strings.Contains(r.Output, "FAIL TestGarden") || !strings.Contains(r.Error, "exit status 3") {
			t.Fatalf("resp = %+v", r)
		}
		var b strings.Builder
		r.Human(&b)
		if !strings.Contains(b.String(), "FAIL TestGarden") || !strings.Contains(b.String(), "Local gate failed") {
			t.Errorf("human = %q", b.String())
		}
	})
	t.Run("dry run refused", func(t *testing.T) {
		f := gateFixture(t, "true", true)
		f.ctx.DryRun = true
		if _, err := callOp(t, f.ctx, "ship.gate", `{}`); err == nil {
			t.Error("dry run ran the gate")
		}
	})
}
