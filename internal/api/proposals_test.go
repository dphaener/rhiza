package api

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

var proposalsAt = time.Date(2026, 11, 13, 9, 0, 0, 0, time.UTC)

// proposalRun is a journal of a run of workflow on change set cs: started,
// then gated naming gated (when any), then ended with end ("" leaves it
// live) and detail.
func proposalRun(id, workflow, cs string, gated []string, end, detail string) run.Journal {
	es := []run.Entry{{Type: run.TypeStarted, At: proposalsAt, Started: &run.Started{Workflow: workflow, Change: cs}}}
	if gated != nil {
		es = append(es, run.Entry{Type: run.TypeGated, At: proposalsAt.Add(time.Minute), Changes: gated})
	}
	if end != "" {
		es = append(es, run.Entry{Type: run.TypeFinished, At: proposalsAt.Add(2 * time.Minute), End: end, Detail: detail})
	}
	return run.Journal{ID: id, Entries: es}
}

func TestUncoveredProposals(t *testing.T) {
	const a, b, rec = "A", "B", "REC"
	tests := []struct {
		name string
		runs []run.Journal
		held map[string]bool
		want []string
	}{
		{"no runs", nil, nil, []string{a, b}},
		{"live run covers", []run.Journal{proposalRun("r1", "plan-review", a, nil, "", "")}, nil, []string{b}},
		{"pending start covers", nil, map[string]bool{a: true}, []string{b}},
		{"gated by a plan run", []run.Journal{proposalRun("r1", "plan", a, []string{a, b}, run.Succeeded, "")}, nil, nil},
		{"gated by a plan-review run on another change set", []run.Journal{proposalRun("r1", "plan-review", "X", []string{a}, run.Succeeded, "")}, nil, []string{b}},
		{"plan-milestone on the record change set gates the sets it names", []run.Journal{proposalRun("r1", "plan-milestone", rec, []string{a}, run.Succeeded, "")}, nil, []string{b}},
		{"a gated entry of another workflow covers nothing", []run.Journal{proposalRun("r1", "ship-change", "X", []string{a}, run.Succeeded, "")}, nil, []string{a, b}},
		{"failed plan-review covers", []run.Journal{proposalRun("r1", "plan-review", a, nil, run.Failed, "boom")}, nil, []string{b}},
		{"failed plan covers", []run.Journal{proposalRun("r1", "plan", a, nil, run.Failed, "boom")}, nil, []string{b}},
		{"abandoned run covers nothing", []run.Journal{proposalRun("r1", "plan-review", a, nil, run.Failed, run.AbandonedPrefix+"the owner")}, nil, []string{a, b}},
		{"finished plan-review without a gate covers nothing", []run.Journal{proposalRun("r1", "plan-review", a, nil, run.Succeeded, "")}, nil, []string{a, b}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := uncoveredProposals([]string{a, b}, tc.runs, tc.held)
			if !slices.Equal(got, tc.want) {
				t.Errorf("uncoveredProposals = %v, want %v", got, tc.want)
			}
		})
	}
}

// A plan-review start not journaled yet holds its change set: the daemon's
// next poll in the same process starts nothing more for it, and a pause
// starts none at all. --max-runs bounds the starts, the ship pick counting.
func TestDaemonPlansProposalsPending(t *testing.T) {
	poll := func(f *daemonFoldFix, max int, paused bool) {
		t.Helper()
		var skips []daemon.Skip
		if _, _, _, err := daemonPoll1(f.ctx, f.rs, f.pending, f.unj, 0, max, 0, false, "manual", paused, nil, nil, &skips); err != nil {
			t.Fatal(err)
		}
	}

	f := newDaemonFoldFix(t)
	// The fixture covers its proposals; drop the runs that do.
	runs, err := f.rs.List()
	if err != nil {
		t.Fatal(err)
	}
	for _, j := range runs {
		if err := os.Remove(filepath.Join(f.runsDir, j.ID+".jsonl")); err != nil {
			t.Fatal(err)
		}
	}
	poll(f, 2, true)
	if len(f.picked) != 0 {
		t.Fatalf("paused poll started %v", f.picked)
	}
	// Three proposed change sets with leaves and one ready: two slots.
	poll(f, 2, false)
	if len(f.picked) != 2 {
		t.Fatalf("first poll with 2 slots started %v, want 2", f.picked)
	}
	reviews := len(f.fold(t).proposed)
	if reviews != 3 {
		t.Fatalf("fixture has %d proposed change sets with leaves, want 3", reviews)
	}
	poll(f, 2, false)
	if len(f.picked) != 2 {
		t.Fatalf("a second poll with the starts pending started more: %v", f.picked)
	}
	// Slots free up: the rest follow, once each.
	poll(f, 10, false)
	poll(f, 10, false)
	seen := map[string]bool{}
	for _, cs := range f.picked {
		if seen[cs] {
			t.Fatalf("%s started twice: %v", cs, f.picked)
		}
		seen[cs] = true
	}
	if len(f.picked) != 4 {
		t.Fatalf("started %v, want the 3 proposals and the ready change set", f.picked)
	}
	// Past the grace the starts never journaled: they are held as not
	// journaled, not started again, until daemonUnjournaledHold passes.
	at := f.clock
	f.ctx.Now = func() time.Time { return at.Add(daemonStartGrace + time.Second) }
	poll(f, 10, false)
	f.ctx.Now = func() time.Time { return at.Add(daemonUnjournaledHold - time.Second) }
	poll(f, 10, false)
	if len(f.picked) != 4 {
		t.Fatalf("an unjournaled plan-review start was started again within the hold: %v", f.picked)
	}
	f.ctx.Now = func() time.Time { return at.Add(daemonUnjournaledHold + time.Minute) }
	poll(f, 10, false)
	// The hold lapsed for the ship pick too: all four start again, once.
	if !slices.Equal(f.picked[4:], f.picked[:4]) {
		t.Fatalf("after the hold started %v, want the 3 proposals and the ready change set again", f.picked)
	}
}

func (f *daemonFoldFix) fold(t *testing.T) *daemonFolded {
	t.Helper()
	fold, err := daemonFoldOf(f.ctx)
	if err != nil {
		t.Fatal(err)
	}
	return fold
}
