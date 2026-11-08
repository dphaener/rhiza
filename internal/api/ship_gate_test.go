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

// The gate can be committed in .rhiza/settings.json, which every checkout
// shares; git config overrides it. A key in it that is no project setting
// is refused.
func TestShipLocalGateSettingsFile(t *testing.T) {
	f := gateFixture(t, "", true)
	file := filepath.Join(f.ctx.Root, ".rhiza", "settings.json")
	if err := os.WriteFile(file, []byte(`{"rhiza.ship.localGate": "echo from-file"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	r := shipGate(t, f)
	if r.Command != "echo from-file" || !r.Passed || !strings.Contains(r.Output, "from-file") {
		t.Errorf("gate = %+v", r)
	}
	rows, err := settingsOf(f.ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range rows {
		if s.Key == LocalGateKey && (s.Value != "echo from-file" || s.Source != ".rhiza/settings.json" || s.Scope != ScopeProject) {
			t.Errorf("setting = %+v", s)
		}
	}
	f.git.out["config --get "+LocalGateKey] = "echo from-config\n"
	if r = shipGate(t, f); r.Command != "echo from-config" {
		t.Errorf("git config did not win: %+v", r)
	}
	for body, want := range map[string]string{
		`{"rhiza.nope": "x"}`:            "not a setting",
		`{"rhiza.daemon.maxRuns": "2"}`:  "machine setting",
		`{"rhiza.ship.localGate": 3}`:    "settings.json",
		`{"rhiza.ship.localGate": "  "}`: "must not be empty",
	} {
		if err := os.WriteFile(file, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := callOp(t, f.ctx, "ship.gate", `{}`)
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: err = %v, want %q", body, err, want)
		}
	}
}

// ship rebase runs the gate on the head it rebased, before the push: the
// gate step before it saw only the head without main's new commits. A
// failing gate puts the branch back and pushes nothing.
func TestShipRebaseRunsLocalGate(t *testing.T) {
	setup := func(t *testing.T, gate string) *sessionFixture {
		f, _ := shipRebaseFixture(t)
		f.git.out["rebase "+rebaseOnto] = ""
		f.git.out["config --get "+LocalGateKey] = gate + "\n"
		if err := writeAtomic(ciStabilityPath(f.common), []byte(`{"runs":[{"head":"a","result":"green","time":"x","flakes":["TestPond"]}]}`)); err != nil {
			t.Fatal(err)
		}
		return f
	}
	t.Run("failing", func(t *testing.T) {
		f := setup(t, "echo FAIL TestMerged; exit 1")
		f.git.out["reset -q --keep beef000111"] = ""
		r := shipRebase(t, f)
		if r.ExitCode() != 1 || r.Gate == nil || r.Gate.Passed || r.Pushed != "" || !strings.Contains(r.Detail, "nothing was pushed") {
			t.Fatalf("rebase = %+v", r)
		}
		for _, c := range f.git.calls {
			if strings.HasPrefix(c, "push") {
				t.Errorf("pushed after a failing gate: %q", c)
			}
		}
		if !slices.Contains(f.git.calls, "reset -q --keep beef000111") {
			t.Errorf("the branch was not put back; calls %q", f.git.calls)
		}
		if got := human(t, r); !strings.Contains(got, "FAIL TestMerged") || !strings.Contains(got, "nothing was pushed") {
			t.Errorf("human = %q", got)
		}
	})
	t.Run("passing", func(t *testing.T) {
		f := setup(t, "echo gate-ok")
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
