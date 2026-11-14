package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// spawnFixture is run R1 in a step "build" with a calling session, and a
// spawner whose start writes the child's record, exited with code.
type spawnFixture struct {
	f     *sessionFixture
	store *run.Store
	sp    spawner
	brief string
	code  int
	reqs  []SessionStartRequest
	slept int
	// onSleep runs at each sleep, to lift a hold.
	onSleep func()
}

func newSpawnFixture(t *testing.T) *spawnFixture {
	t.Helper()
	f, _, _ := runFixture(t, `{"permission_mode":"manual"}`)
	t0 := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	s := &spawnFixture{f: f, store: run.NewStore(run.Dir(f.common))}
	if err := s.store.Create("R1", run.Entry{
		Type: run.TypeStarted, Step: "build", At: t0,
		Started: &run.Started{Workflow: "w", Change: "01M76F4G0JQ649RR29M9QPFX68", PermissionMode: "manual", Agent: "claude-code"},
	}); err != nil {
		t.Fatal(err)
	}
	sessions := session.NewStore(session.Dir(f.common))
	if err := sessions.Write(session.Record{ID: "parent-1", Run: "R1", Step: "build"}); err != nil {
		t.Fatal(err)
	}
	s.brief = filepath.Join(t.TempDir(), "child-brief.md")
	if err := os.WriteFile(s.brief, []byte("Check leaf run-spawn.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	f.ctx.Getenv = func(k string) string {
		switch k {
		case envRun:
			return "R1"
		case envSession:
			return "parent-1"
		}
		return ""
	}
	f.ctx.Sleep = func(time.Duration) {
		s.slept++
		if s.onSleep != nil {
			s.onSleep()
		}
	}
	s.sp = spawner{ctx: f.ctx, start: func(r SessionStartRequest) (string, error) {
		s.reqs = append(s.reqs, r)
		if _, err := r.briefFor("child-1"); err != nil {
			return "", err
		}
		code := s.code
		return "child-1", sessions.Write(session.Record{ID: "child-1", Run: "R1", Step: r.Step, Parent: r.parent,
			Status: session.Status{Exit: &session.Report{State: session.StateExited, Code: &code}}})
	}}
	return s
}

func (s *spawnFixture) req() RunSpawnRequest {
	return RunSpawnRequest{Role: "verifier", Leaf: "run-spawn", Brief: s.brief}
}

func (s *spawnFixture) state(t *testing.T) run.State {
	t.Helper()
	j, err := s.store.Read("R1")
	if err != nil {
		t.Fatal(err)
	}
	return j.State()
}

func TestRunSpawn(t *testing.T) {
	t.Run("Pass", func(t *testing.T) {
		s := newSpawnFixture(t)
		resp, err := s.sp.spawn(s.req())
		if err != nil {
			t.Fatal(err)
		}
		if resp.Outcome != "pass" || resp.ExitCode() != 0 || resp.Held || resp.Session != "child-1" || resp.Step != "build" || resp.Parent != "parent-1" {
			t.Fatalf("response %+v", resp)
		}
		st := s.state(t)
		if len(st.Spawns) != 1 {
			t.Fatalf("spawns %+v", st.Spawns)
		}
		sp := st.Spawns[0]
		if sp.Step != "build" || sp.Parent != "parent-1" || sp.Session != "child-1" || sp.Role != "verifier" || sp.Leaf != "run-spawn" || !sp.Ended || sp.Outcome != "pass" {
			t.Errorf("spawn %+v", sp)
		}
		var seen bool
		for _, rs := range st.Sessions {
			seen = seen || rs.Session == "child-1" && rs.Step == "build"
		}
		if !seen {
			t.Errorf("child not among the run's sessions under its calling step: %+v", st.Sessions)
		}
		r := s.reqs[0]
		if r.Run != "R1" || r.PermissionMode != "manual" || r.Agent != "claude-code" || r.parent != "parent-1" || r.role != "verifier" || r.leaf != "run-spawn" {
			t.Errorf("start request %+v", r)
		}
		b, err := os.ReadFile(resp.Brief)
		if err != nil {
			t.Fatal(err)
		}
		want := run.RenderSpawnBrief("Check leaf run-spawn.\n")
		if string(b) != want || !strings.HasPrefix(want, "# Brief\n\n## Contract\n\n") || !strings.Contains(want, run.ContractLastLine+"\n\n## Brief\n\nCheck leaf run-spawn.\n") {
			t.Errorf("brief.md\n%s", b)
		}
		if _, err := os.Stat(filepath.Join(filepath.Dir(resp.Brief), sliceFile)); err == nil {
			t.Error("a spawned child's directory holds a slice.json")
		}
	})

	t.Run("Fail", func(t *testing.T) {
		s := newSpawnFixture(t)
		s.code = 1
		resp, err := s.sp.spawn(s.req())
		if err != nil {
			t.Fatal(err)
		}
		if resp.Outcome != "fail" || resp.Detail != "exit 1" || resp.ExitCode() != 1 {
			t.Fatalf("response %+v", resp)
		}
		if sp := s.state(t).Spawns[0]; sp.Outcome != "fail" || sp.Detail != "exit 1" {
			t.Errorf("spawn %+v", sp)
		}
	})

	t.Run("OutsideRun", func(t *testing.T) {
		s := newSpawnFixture(t)
		s.f.ctx.Getenv = func(string) string { return "" }
		_, err := s.sp.spawn(s.req())
		if e := AsError(err); e == nil || e.Code != CodeUsage || !strings.Contains(e.Message, "RHIZA_RUN") {
			t.Fatalf("err %v", err)
		}
		if len(s.reqs) != 0 {
			t.Error("a child started")
		}
	})

	t.Run("MissingBrief", func(t *testing.T) {
		s := newSpawnFixture(t)
		r := s.req()
		r.Brief = filepath.Join(filepath.Dir(s.brief), "nope.md")
		_, err := s.sp.spawn(r)
		if e := AsError(err); e == nil || e.Code != CodeUsage || !strings.Contains(e.Message, r.Brief) {
			t.Fatalf("err %v", err)
		}
		if len(s.reqs) != 0 || len(s.state(t).Spawns) != 0 {
			t.Error("a child started")
		}
	})

	t.Run("UsageHold", func(t *testing.T) {
		s := newSpawnFixture(t)
		hold := daemon.UsageHold{Since: s.f.ctx.Now(), Window: "5h", Used: 90, Projected: 95, Limit: 100, Reserve: 20, Reset: s.f.ctx.Now().Add(time.Hour)}
		write := func(st daemon.State) {
			if err := daemon.WriteState(s.f.common, st); err != nil {
				t.Fatal(err)
			}
		}
		write(daemon.State{UsageHold: &hold})
		s.onSleep = func() { write(daemon.State{}) }
		resp, err := s.sp.spawn(s.req())
		if err != nil {
			t.Fatal(err)
		}
		if !resp.Held || resp.Outcome != "pass" || s.slept != 1 {
			t.Fatalf("response %+v, slept %d", resp, s.slept)
		}
		holds := s.state(t).SpawnHolds
		if len(holds) != 1 || holds[0].Held() || holds[0].Step != "build" || !strings.Contains(holds[0].Detail, "usage hold") {
			t.Errorf("holds %+v", holds)
		}
	})

	t.Run("MachineHold", func(t *testing.T) {
		s := newSpawnFixture(t)
		write := func(st daemon.State) {
			if err := daemon.WriteState(s.f.common, st); err != nil {
				t.Fatal(err)
			}
		}
		write(daemon.State{MachineHold: &daemon.MachineHold{Since: s.f.ctx.Now(), Reason: daemon.SkipLoadHigh, Detail: "load 9.0 over 8.0", Load: true}})
		s.onSleep = func() {
			if len(s.state(t).SpawnHolds) != 1 || !s.state(t).SpawnHolds[0].Held() {
				t.Error("the wait was not journaled as held while it waited")
			}
			write(daemon.State{})
		}
		resp, err := s.sp.spawn(s.req())
		if err != nil {
			t.Fatal(err)
		}
		if !resp.Held || s.slept != 1 {
			t.Fatalf("response %+v, slept %d", resp, s.slept)
		}
		holds := s.state(t).SpawnHolds
		if len(holds) != 1 || holds[0].Held() || !strings.Contains(holds[0].Detail, "load 9.0") {
			t.Errorf("holds %+v", holds)
		}
	})

	t.Run("MaxRunsHoldsNothing", func(t *testing.T) {
		s := newSpawnFixture(t)
		if err := daemon.WriteState(s.f.common, daemon.State{MaxRuns: 1, Runs: []daemon.LiveRun{{}}}); err != nil {
			t.Fatal(err)
		}
		resp, err := s.sp.spawn(s.req())
		if err != nil {
			t.Fatal(err)
		}
		if resp.Held || s.slept != 0 || len(s.state(t).SpawnHolds) != 0 {
			t.Fatalf("response %+v, slept %d", resp, s.slept)
		}
	})
}

// In a conductor run's conduct step (declaring checks plan, code and qa) a
// check-plan child gets check plan and a cycle, so its RHIZA_CHECK and
// RHIZA_CYCLE are set; a critic names its check with --check; a --check
// the step does not declare is refused.
func TestRunSpawnCheck(t *testing.T) {
	conductor := func(t *testing.T) *spawnFixture {
		s := newSpawnFixture(t)
		def := "[step.conduct]\nkind = \"agent\"\nprompt = \"go\"\n" +
			"[[step.conduct.check]]\nname = \"plan\"\n[[step.conduct.check]]\nname = \"code\"\n[[step.conduct.check]]\nname = \"qa\"\n"
		if err := s.store.Create("R2", run.Entry{
			Type: run.TypeStarted, Step: "conduct", At: s.f.ctx.Now(),
			Started: &run.Started{Workflow: "conduct", Definition: def, Change: "01M76F4G0JQ649RR29M9QPFX68", PermissionMode: "manual", Agent: "claude-code"},
		}); err != nil {
			t.Fatal(err)
		}
		if err := s.store.Append("R2", run.Entry{Type: run.TypeEntered, Step: "conduct", Visit: 1, At: s.f.ctx.Now()}); err != nil {
			t.Fatal(err)
		}
		get := s.f.ctx.Getenv
		s.f.ctx.Getenv = func(k string) string {
			if k == envRun {
				return "R2"
			}
			return get(k)
		}
		return s
	}
	t.Run("CheckRole", func(t *testing.T) {
		s := conductor(t)
		resp, err := s.sp.spawn(RunSpawnRequest{Role: "check-plan", Brief: s.brief})
		if err != nil {
			t.Fatal(err)
		}
		if r := s.reqs[0]; r.check != "plan" || r.cycle != 1 {
			t.Errorf("start request check %q cycle %d", r.check, r.cycle)
		}
		if resp.Check != "plan" || resp.Cycle != 1 {
			t.Errorf("response %+v", resp)
		}
	})
	t.Run("Critic", func(t *testing.T) {
		s := conductor(t)
		if err := s.store.Append("R2", run.Entry{Type: run.TypeFindings, Step: "conduct", Check: "qa", At: s.f.ctx.Now(),
			Findings: []run.Finding{{ID: "f1", Check: "qa", Severity: "blocker", Place: "a.go:1", Reason: "r"}}}); err != nil {
			t.Fatal(err)
		}
		if _, err := s.sp.spawn(RunSpawnRequest{Role: "critic", Check: "qa", Brief: s.brief}); err != nil {
			t.Fatal(err)
		}
		if r := s.reqs[0]; r.check != "qa" || r.cycle != 1 {
			t.Errorf("start request check %q cycle %d", r.check, r.cycle)
		}
	})
	t.Run("UnknownCheck", func(t *testing.T) {
		s := conductor(t)
		_, err := s.sp.spawn(RunSpawnRequest{Role: "critic", Check: "scope", Brief: s.brief})
		if e := AsError(err); e == nil || e.Code != CodeUsage || !strings.Contains(e.Message, "plan, code, qa") {
			t.Fatalf("err %v", err)
		}
		if len(s.reqs) != 0 {
			t.Error("a child started")
		}
	})
	t.Run("NoChecksStep", func(t *testing.T) {
		s := newSpawnFixture(t)
		if _, err := s.sp.spawn(RunSpawnRequest{Role: "verifier", Leaf: "run-spawn", Brief: s.brief}); err != nil {
			t.Fatal(err)
		}
		if r := s.reqs[0]; r.check != "" || r.cycle != 0 {
			t.Errorf("start request check %q cycle %d", r.check, r.cycle)
		}
	})
}
