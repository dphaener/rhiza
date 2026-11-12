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
