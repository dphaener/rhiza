package api

import (
	"github.com/dphaener/rhiza/internal/run"
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
	t.Run("unidentified flake only skips with the reason", func(t *testing.T) {
		f := gateFixture(t, "echo gate-ran; exit 1", false)
		if err := writeAtomic(ciStabilityPath(f.common), []byte(`{"runs":[{"head":"a","result":"green","time":"x","flakes":["ci / UNKNOWN STEP"]}]}`)); err != nil {
			t.Fatal(err)
		}
		r := shipGate(t, f)
		if !r.Skipped || r.Output != "" || r.ExitCode() != 0 || !strings.HasPrefix(r.Reason, "CI is stable") || !strings.Contains(r.Reason, "ignored unidentified flake") || !strings.Contains(r.Reason, "ci / UNKNOWN STEP") {
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

// ship rebase runs the gate on the head it merged, before the push: the
// gate step before it saw only the head without main's new commits. A
// failing gate pushes nothing and leaves main merged into the branch,
// unpushed and recorded, for the gate's fixer: the head before the merge
// already passed the gate, so a fixer there would see no failure.
func TestShipRebaseRunsLocalGate(t *testing.T) {
	setup := func(t *testing.T, gate string) (*sessionFixture, string) {
		f, state := shipRebaseFixture(t)
		f.git.out["config --get "+LocalGateKey] = gate + "\n"
		if err := writeAtomic(ciStabilityPath(f.common), []byte(`{"runs":[{"head":"a","result":"green","time":"x","flakes":["TestPond"]}]}`)); err != nil {
			t.Fatal(err)
		}
		return f, state
	}
	t.Run("failing", func(t *testing.T) {
		f, state := setup(t, "echo FAIL TestMerged; exit 1")
		f.git.out[rebaseMerge] = ""
		r := shipRebase(t, f)
		if r.ExitCode() != 1 || r.Gate == nil || r.Gate.Passed || r.Pushed != "" || r.Kept != "beef000111" || !strings.Contains(r.Detail, "nothing was pushed") {
			t.Fatalf("rebase = %+v", r)
		}
		for _, c := range f.git.calls {
			if strings.HasPrefix(c, "push") {
				t.Errorf("pushed after a failing gate: %q", c)
			}
			if strings.HasPrefix(c, "reset") || strings.HasPrefix(c, "merge --abort") {
				t.Errorf("the merge the gate failed on was not kept: %q", c)
			}
		}
		if !slices.Contains(f.git.calls, rebaseMerge) {
			t.Errorf("main was not merged into the branch; calls %q", f.git.calls)
		}
		if st, _ := readShipState(state); st.GateFailed != "beef000111" || st.Pushed != "" {
			t.Errorf("state = %+v; want the kept head recorded, nothing pushed", st)
		}
		if got := human(t, r); !strings.Contains(got, "FAIL TestMerged") || !strings.Contains(got, "for the gate's fixer") || !strings.Contains(got, "Full log: "+r.Gate.Log) {
			t.Errorf("human = %q", got)
		}
		// The rebase's gate writes the same log ship gate does.
		want := filepath.Join(f.common, "rhiza", "gate", "cs", "ship-it.log")
		if r.Gate.Log != want {
			t.Errorf("gate log = %q, want %q", r.Gate.Log, want)
		}
		if b, err := os.ReadFile(want); err != nil || string(b) != "FAIL TestMerged\n" {
			t.Errorf("log = %q, %v; want the whole output", b, err)
		}
	})
	t.Run("a conflict runs no gate", func(t *testing.T) {
		f, state := setup(t, "touch "+filepath.Join(t.TempDir(), "ran"))
		f.git.out[rebaseConflict] = "feature.txt\n"
		r := shipRebase(t, f)
		if r.ExitCode() != 1 || r.Gate != nil || r.Kept != "" || !slices.Equal(r.Conflicts, []string{"feature.txt"}) {
			t.Fatalf("rebase = %+v", r)
		}
		if st, _ := readShipState(state); st.GateFailed != "" || st.MergeConflict == nil {
			t.Errorf("state = %+v; want the conflict recorded, no kept head", st)
		}
	})
	t.Run("passing", func(t *testing.T) {
		f, state := setup(t, "echo gate-ok")
		f.git.out[rebaseMerge] = ""
		f.git.out[rebasePush] = ""
		r := shipRebase(t, f)
		if r.ExitCode() != 0 || r.Gate == nil || !r.Gate.Passed || r.Pushed != "beef000111" {
			t.Fatalf("rebase = %+v", r)
		}
		if !slices.Contains(f.git.calls, rebasePush) {
			t.Errorf("not pushed; calls %q", f.git.calls)
		}
		if st, _ := readShipState(state); st.GateFailed != "" || st.Pushed != "beef000111" {
			t.Errorf("state = %+v", st)
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

// While the branch is at the merge ship rebase left after the gate failed
// on it, ship check fails, whatever the PR's CI on the pushed head says:
// the run routes that to the gate's fixer. A wait (--settled) is not held
// by it, and a new head on the branch ends it.
func TestShipCheckFailsOnGateFailedHead(t *testing.T) {
	f, g, state := shipFixture(t)
	g.answers[shipPRView] = openPR(greenRollup)
	g.answers["git rev-parse HEAD"] = "beef000111\n"
	if err := writeShipState(state, shipState{GateFailed: "beef000111"}); err != nil {
		t.Fatal(err)
	}
	r := shipCheck(t, f, `{}`)
	if r.ExitCode() != 1 || r.Passed || !strings.Contains(r.Detail, "local gate failed on beef000") {
		t.Errorf("check on the kept head = %+v", r)
	}
	if r := shipCheck(t, f, `{"settled":true}`); r.ExitCode() != 0 {
		t.Errorf("settled check on the kept head = %+v", r)
	}
	g.answers["git rev-parse HEAD"] = "cafe000222\n"
	if r := shipCheck(t, f, `{}`); r.ExitCode() != 0 || !r.Passed {
		t.Errorf("check after a fix = %+v", r)
	}
}

// The same holds for the head the pre-push check failed on.
func TestShipCheckFailsOnPrePushFailedHead(t *testing.T) {
	f, g, state := shipFixture(t)
	g.answers[shipPRView] = openPR(greenRollup)
	g.answers["git rev-parse HEAD"] = "beef000111\n"
	if err := writeShipState(state, shipState{PrePushFailed: "beef000111"}); err != nil {
		t.Fatal(err)
	}
	r := shipCheck(t, f, `{}`)
	if r.ExitCode() != 1 || r.Passed || !strings.Contains(r.Detail, "pre-push check failed on beef000") {
		t.Errorf("check on the kept head = %+v", r)
	}
	if r := shipCheck(t, f, `{"settled":true}`); r.ExitCode() != 0 {
		t.Errorf("settled check on the kept head = %+v", r)
	}
	g.answers["git rev-parse HEAD"] = "cafe000222\n"
	if r := shipCheck(t, f, `{}`); r.ExitCode() != 0 || !r.Passed {
		t.Errorf("check after a fix = %+v", r)
	}
}

// ship gate --failed also passes on a head the pre-push check failed on in
// ship open-pr or land (recorded in the ship state), so the run's triage
// step sends it to the gate's fixer; a failure that was not the pre-push
// check (a gh or push error) records nothing and the triage fails.
func TestShipGateFailedOnPrePushFailure(t *testing.T) {
	f, g, state := shipFixture(t)
	g.answers["git rev-parse HEAD"] = "beef000111\n"
	ask := func() ShipGateResponse {
		t.Helper()
		r, err := callOp(t, f.ctx, "ship.gate", `{"failed":true}`)
		if err != nil {
			t.Fatal(err)
		}
		return r.(ShipGateResponse)
	}
	if r := ask(); r.ExitCode() != 1 {
		t.Errorf("after a gh or push error: %+v", r)
	}
	tgt, err := findShipTarget(f.ctx)
	if err != nil {
		t.Fatal(err)
	}
	shipRecordPrePushFailed(tgt, &ShipPrePushCheck{Head: "beef000111", Lines: []string{"x"}})
	if st, _ := readShipState(state); st.PrePushFailed != "beef000111" {
		t.Fatalf("state = %+v", st)
	}
	if r := ask(); r.ExitCode() != 0 || r.Failed != "beef000111" {
		t.Errorf("on the pre-push-failed head: %+v", r)
	}
	g.answers["git rev-parse HEAD"] = "cafe000222\n"
	if r := ask(); r.ExitCode() != 1 {
		t.Errorf("a fixed head still counts as the failure: %+v", r)
	}
}

// The gate's whole output goes to a log file the response names, and a
// failing gate's human output ends with its failing lines and the path.
func TestShipGateFullLog(t *testing.T) {
	script := "echo '--- FAIL: TestX'; seq 1 20000 | sed 's/^/filler line /'; echo FAIL; echo 'make: *** [Makefile:64: test] Error 1'; exit 2"
	f := gateFixture(t, script, true)
	r, err := callOp(t, f.ctx, "ship.gate", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	g := r.(ShipGateResponse)
	if g.Passed || g.Log == "" {
		t.Fatalf("gate = %+v", g)
	}
	b, err := os.ReadFile(g.Log)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) <= 64<<10 || !strings.HasPrefix(string(b), "--- FAIL: TestX\n") || !strings.Contains(string(b), "Error 1") {
		t.Errorf("log holds %d bytes, want the whole output", len(b))
	}
	lines := strings.Split(strings.TrimRight(human(t, g), "\n"), "\n")
	if len(lines) > 20 {
		lines = lines[len(lines)-20:]
	}
	tail := strings.Join(lines, "\n")
	if !strings.Contains(tail, "--- FAIL: TestX") || !strings.Contains(tail, g.Log) {
		t.Errorf("last 20 lines = %q", tail)
	}
	t.Run("a passing gate writes it too", func(t *testing.T) {
		f := gateFixture(t, "echo fine", true)
		g := shipGate(t, f)
		if !g.Passed || g.Log == "" {
			t.Fatalf("gate = %+v", g)
		}
		if b, _ := os.ReadFile(g.Log); !strings.Contains(string(b), "fine") {
			t.Errorf("log = %q", b)
		}
	})
}

// In a run, each gate run writes its own log in the run's directory and
// never overwrites an earlier one; outside a run the log is per branch.
func TestShipGateRunLog(t *testing.T) {
	script := "echo '--- FAIL: TestX'; seq 1 20000 | sed 's/^/filler line /'; echo FAIL; echo 'make: *** [Makefile:64: test] Error 1'; exit 2"
	f := gateFixture(t, script, true)
	f.ctx.Getenv = func(k string) string {
		if k == envRun {
			return "01RUNLOG"
		}
		return ""
	}
	g1 := shipGate(t, f)
	g2 := shipGate(t, f)
	wantDir := filepath.Join(run.Dir(f.common), "01RUNLOG")
	if filepath.Dir(g1.Log) != wantDir || filepath.Dir(g2.Log) != wantDir || g1.Log == g2.Log {
		t.Fatalf("logs = %q, %q; want two files in %s", g1.Log, g2.Log, wantDir)
	}
	for _, p := range []string{g1.Log, g2.Log} {
		b, err := os.ReadFile(p)
		if err != nil || len(b) <= 64<<10 || !strings.Contains(string(b), "Error 1") {
			t.Errorf("%s holds %d bytes, %v; want the whole output", p, len(b), err)
		}
	}
	lines := strings.Split(strings.TrimRight(human(t, g1), "\n"), "\n")
	tail := strings.Join(lines[max(0, len(lines)-20):], "\n")
	if !strings.Contains(tail, "--- FAIL: TestX") || !strings.Contains(tail, "Full log: "+g1.Log) {
		t.Errorf("last 20 lines = %q", tail)
	}
	if len(tail) > stepTailBytes {
		t.Errorf("tail is %d bytes, over the journal's %d: the log path would be cut", len(tail), stepTailBytes)
	}
	pass := gateFixture(t, "echo fine", true)
	pass.ctx.Getenv = f.ctx.Getenv
	if g := shipGate(t, pass); !g.Passed || !strings.HasPrefix(g.Log, filepath.Join(run.Dir(pass.common), "01RUNLOG")) {
		t.Errorf("passing gate = %+v", g)
	}
}
