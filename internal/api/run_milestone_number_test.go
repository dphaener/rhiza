package api

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// A plan run for a new milestone draws one more than the highest landed
// milestone number and every number any run journal holds, so two runs
// started back to back get 36 then 37, and a number held only by an ended
// run is still skipped.
func TestReserveMilestoneNumber(t *testing.T) {
	ctx, _, _, _ := openPRFixture(t, false)
	for _, req := range []string{
		`{"kind":"intent","handle":"milestone-12","fields":["statement=Twelve","milestone=12"]}`,
		`{"kind":"intent","handle":"milestone-35","fields":["statement=Thirty-five","milestone=35"]}`,
		`{"kind":"intent","handle":"point","parent":"milestone-12","fields":["statement=Not a milestone"]}`,
	} {
		recordWrite(t, ctx, "add", req)
	}
	rs := run.NewStore(run.Dir(filepath.Join(ctx.Root, ".git")))
	start := func(rid string, n int) {
		t.Helper()
		err := rs.Create(rid, run.Entry{
			Type: run.TypeStarted, Step: "plan", At: time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC),
			Started: &run.Started{Workflow: "plan-milestone", Change: "C" + rid, MilestoneNumber: n},
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	draw := func(want int) {
		t.Helper()
		got, err := reserveMilestoneNumber(ctx, rs)
		if err != nil || got != want {
			t.Fatalf("reserveMilestoneNumber = %d, %v; want %d", got, err, want)
		}
	}
	draw(36)
	start("M1", 36)
	draw(37)
	start("M2", 37)
	// A run that planned an existing milestone reserved nothing.
	start("M3", 0)
	draw(38)
	// A failed run keeps its number.
	if err := rs.Append("M2", run.Entry{Type: run.TypeEnded, Step: "plan", Outcome: "fail", At: time.Date(2026, 11, 1, 13, 0, 0, 0, time.UTC)}); err != nil {
		t.Fatal(err)
	}
	draw(38)
	if got := mustState(t, rs, "M2").MilestoneNumber; got != 37 {
		t.Errorf("journal state MilestoneNumber = %d, want 37", got)
	}
	info := runInfo(ctx, mustState(t, rs, "M2"), false, nil)
	if info.MilestoneNumber != 37 {
		t.Errorf("run status milestone_number = %d, want 37", info.MilestoneNumber)
	}
}

func mustState(t *testing.T, rs *run.Store, rid string) run.State {
	t.Helper()
	j, err := rs.Read(rid)
	if err != nil {
		t.Fatal(err)
	}
	return j.State()
}

func TestLaunchEnvMilestoneNumber(t *testing.T) {
	has := func(l []string, s string) bool { return slices.Contains(l, s) }
	set, unset := launchEnv(session.Record{ID: "S1", Change: "CS1", Run: "R1", MilestoneNumber: 36})
	if !has(set, "RHIZA_MILESTONE_NUMBER=36") || has(unset, "RHIZA_MILESTONE_NUMBER") {
		t.Errorf("new-milestone plan: set %v unset %v", set, unset)
	}
	set, unset = launchEnv(session.Record{ID: "S1", Change: "CS1", Run: "R1"})
	if !has(unset, "RHIZA_MILESTONE_NUMBER") {
		t.Errorf("no number: unset %v lacks RHIZA_MILESTONE_NUMBER", unset)
	}
	for _, s := range set {
		if strings.HasPrefix(s, "RHIZA_MILESTONE_NUMBER=") {
			t.Errorf("no number: set %v", set)
		}
	}
}

// run milestone-check on a plan run that reserved a milestone number passes
// a new milestone intent using it (handle milestone-N, field milestone=N,
// done points prefixed mN-) and fails one with the wrong handle, the wrong
// field or an unprefixed done point, naming the intent and the number.
func TestRunMilestoneCheckEnforcesNumber(t *testing.T) {
	ctx, _, _, _ := openPRFixture(t, false)
	ctx.IndexPath = filepath.Join(t.TempDir(), "index.db")
	store := run.NewStore(run.Dir(filepath.Join(ctx.Root, ".git")))
	n := 0
	start := func(adds ...string) {
		t.Helper()
		n++
		cs := newChange(t, ctx, `{"intent":"Plan `+string(rune('a'+n))+`","title":"Plan"}`).ChangeSet.String()
		runID := "RN" + string(rune('A'+n))
		if err := store.Create(runID, run.Entry{
			Type: run.TypeStarted, Step: "milestone-check", At: time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC),
			Started: &run.Started{Workflow: "plan-milestone", Change: cs, Title: "Plan", MilestoneNumber: 7},
		}); err != nil {
			t.Fatal(err)
		}
		ctx.Getenv = func(k string) string {
			if k == envRun {
				return runID
			}
			return ""
		}
		for _, a := range adds {
			if _, err := callOp(t, ctx, "add", a); err != nil {
				t.Fatal(err)
			}
		}
	}
	check := func() error {
		_, err := callOp(t, ctx, "run.milestone-check", `{}`)
		return err
	}
	wantFail := func(what, intent string) {
		t.Helper()
		err := check()
		if err == nil {
			t.Fatalf("%s: passed; want a failure", what)
		}
		for _, s := range []string{intent, "milestone number 7", "Fix:"} {
			if !strings.Contains(err.Error(), s) {
				t.Errorf("%s: error %q lacks %q", what, err, s)
			}
		}
	}

	start(
		`{"kind":"intent","handle":"milestone-7","fields":["statement=Seven","milestone=7"]}`,
		`{"kind":"intent","handle":"m7-first","parent":"milestone-7","fields":["statement=A done point"]}`,
	)
	if err := check(); err != nil {
		t.Errorf("reserved number used: %v", err)
	}
	start(`{"kind":"intent","handle":"milestone-8","fields":["statement=Eight","milestone=7"]}`)
	wantFail("wrong handle", "milestone-8")
	start(`{"kind":"intent","handle":"milestone-7","fields":["statement=Seven","milestone=8"]}`)
	wantFail("wrong milestone field", "milestone-7")
	start(
		`{"kind":"intent","handle":"milestone-7","fields":["statement=Seven","milestone=7"]}`,
		`{"kind":"intent","handle":"first","parent":"milestone-7","fields":["statement=A done point"]}`,
	)
	wantFail("done point without the prefix", `"first"`)
	// The planner lands the record during the interview, before this step:
	// its new milestone intent is still checked.
	start(`{"kind":"intent","handle":"milestone-8","fields":["statement=Eight","milestone=8"]}`)
	wantFail("landed record, before landing", "milestone-8")
	for _, s := range []string{"propose", "accept", "land"} {
		if _, err := callOp(t, ctx, s, `{}`); err != nil {
			t.Fatalf("%s: %v", s, err)
		}
	}
	wantFail("landed record", "milestone-8")
}
