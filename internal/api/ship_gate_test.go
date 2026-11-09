package api

import (
	"os"
	"path/filepath"
	"slices"
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

// A head the gate passed on, with a clean worktree, passes again without
// running the gate: the gate before rebase or land re-checks the head the
// gate before open-pr passed. A new head runs it again.
func TestShipLocalGateOncePerHead(t *testing.T) {
	f, g, _ := shipFixture(t)
	log := filepath.Join(t.TempDir(), "gate.log")
	f.git.out["config --get "+LocalGateKey] = "echo ran >> " + log + "\n"
	if err := writeAtomic(ciStabilityPath(f.common), []byte(`{"runs":[{"head":"a","result":"green","time":"x","flakes":["TestPond"]}]}`)); err != nil {
		t.Fatal(err)
	}
	g.answers["git rev-parse HEAD"] = "beef000111\n"
	g.answers["git status --porcelain --untracked-files=no"] = ""
	runs := func() int {
		b, _ := os.ReadFile(log)
		return strings.Count(string(b), "ran\n")
	}
	r := shipGate(t, f)
	if !r.Passed || r.Head != "beef000111" || runs() != 1 {
		t.Fatalf("first gate = %+v, runs %d", r, runs())
	}
	r = shipGate(t, f)
	if !r.Passed || r.Skipped || !strings.Contains(r.Reason, "already passed on beef000") || runs() != 1 {
		t.Errorf("second gate on the same head = %+v, runs %d", r, runs())
	}
	g.answers["git status --porcelain --untracked-files=no"] = " M x.go\n"
	if r = shipGate(t, f); !r.Passed || runs() != 2 {
		t.Errorf("gate with a dirty tree = %+v, runs %d", r, runs())
	}
	g.answers["git status --porcelain --untracked-files=no"] = ""
	g.answers["git rev-parse HEAD"] = "cafe000222\n"
	if r = shipGate(t, f); !r.Passed || r.Head != "cafe000222" || runs() != 3 {
		t.Errorf("gate on a new head = %+v, runs %d", r, runs())
	}
}

// ship rebase runs the gate on the head it rebased, before the push: the
// gate step before it saw only the head without main's new commits. A
// failing gate pushes nothing and leaves main merged into the branch as it
// was, unpushed and recorded, for the gate's fixer: the head before the
// rebase already passed the gate, so a fixer there would see no failure.
func TestShipRebaseRunsLocalGate(t *testing.T) {
	setup := func(t *testing.T, gate string) (*sessionFixture, string) {
		f, state := shipRebaseFixture(t)
		f.git.out["rebase "+rebaseOnto] = ""
		f.git.out["config --get "+LocalGateKey] = gate + "\n"
		if err := writeAtomic(ciStabilityPath(f.common), []byte(`{"runs":[{"head":"a","result":"green","time":"x","flakes":["TestPond"]}]}`)); err != nil {
			t.Fatal(err)
		}
		return f, state
	}
	const merge = "merge --no-edit -m Merge origin/main into cs/ship-it " + rebaseOnto
	t.Run("failing", func(t *testing.T) {
		f, state := setup(t, "echo FAIL TestMerged; exit 1")
		f.git.out["reset -q --keep beef000111"] = ""
		f.git.out[merge] = ""
		r := shipRebase(t, f)
		if r.ExitCode() != 1 || r.Gate == nil || r.Gate.Passed || r.Pushed != "" || r.Kept != "beef000111" || !strings.Contains(r.Detail, "nothing was pushed") {
			t.Fatalf("rebase = %+v", r)
		}
		for _, c := range f.git.calls {
			if strings.HasPrefix(c, "push") {
				t.Errorf("pushed after a failing gate: %q", c)
			}
		}
		i := slices.Index(f.git.calls, "reset -q --keep beef000111")
		if i < 0 || !slices.Contains(f.git.calls[i:], merge) {
			t.Errorf("main was not merged into the branch as it was; calls %q", f.git.calls)
		}
		if st, _ := readShipState(state); st.GateFailed != "beef000111" || st.Pushed != "" {
			t.Errorf("state = %+v; want the kept head recorded, nothing pushed", st)
		}
		if got := human(t, r); !strings.Contains(got, "FAIL TestMerged") || !strings.Contains(got, "for the gate's fixer") {
			t.Errorf("human = %q", got)
		}
	})
	t.Run("failing, and the merge conflicts", func(t *testing.T) {
		f, state := setup(t, "exit 1")
		f.git.out["reset -q --keep beef000111"] = ""
		f.git.out["merge --abort"] = ""
		f.git.out[rebaseConflict] = "feature.txt\n"
		r := shipRebase(t, f)
		if r.ExitCode() != 1 || r.Kept != "" || !strings.Contains(r.Detail, "the branch is as it was") {
			t.Fatalf("rebase = %+v", r)
		}
		if n := strings.Count(strings.Join(f.git.calls, "\n"), "reset -q --keep beef000111"); n != 2 {
			t.Errorf("put back %d times; want 2; calls %q", n, f.git.calls)
		}
		if st, _ := readShipState(state); st.GateFailed != "" {
			t.Errorf("state = %+v; want no kept head", st)
		}
	})
	t.Run("passing", func(t *testing.T) {
		f, _ := setup(t, "echo gate-ok")
		f.git.out[rebasePush] = ""
		r := shipRebase(t, f)
		if r.ExitCode() != 0 || r.Gate == nil || !r.Gate.Passed || r.Pushed != "beef000111" {
			t.Fatalf("rebase = %+v", r)
		}
		if !slices.Contains(f.git.calls, rebasePush) {
			t.Errorf("not pushed; calls %q", f.git.calls)
		}
	})
}

// ship gate --failed runs no gate: it passes only while the branch is at
// the head ship rebase left after the gate failed on it, so a rebase
// conflict ends the run and a failing gate reaches the fixer.
func TestShipGateFailed(t *testing.T) {
	f, g, state := shipFixture(t)
	f.git.out["config --get "+LocalGateKey] = "touch " + filepath.Join(t.TempDir(), "ran") + "\n"
	g.answers["git rev-parse HEAD"] = "beef000111\n"
	ask := func() ShipGateResponse {
		t.Helper()
		r, err := callOp(t, f.ctx, "ship.gate", `{"failed":true}`)
		if err != nil {
			t.Fatal(err)
		}
		return r.(ShipGateResponse)
	}
	if r := ask(); r.ExitCode() != 1 || r.Failed != "" {
		t.Errorf("no kept head: %+v", r)
	}
	st, _ := readShipState(state)
	st.GateFailed = "beef000111"
	if err := writeShipState(state, st); err != nil {
		t.Fatal(err)
	}
	r := ask()
	if r.ExitCode() != 0 || r.Failed != "beef000111" || r.Command != "" {
		t.Errorf("kept head: %+v", r)
	}
	if got := human(t, r); !strings.Contains(got, "rhiza ship gate") {
		t.Errorf("human = %q", got)
	}
	g.answers["git rev-parse HEAD"] = "cafe000222\n"
	if r := ask(); r.ExitCode() != 1 {
		t.Errorf("a fixed head still counts as the failure: %+v", r)
	}
}
