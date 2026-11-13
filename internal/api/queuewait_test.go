package api

import (
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

// queueWaitPoll runs one poll with the given slots, holds and pause.
func queueWaitPoll(t *testing.T, f *daemonFoldFix, max int, hold *usageHold, machine *daemon.MachineHold, paused bool) []string {
	t.Helper()
	var skips []daemon.Skip
	started, _, _, err := daemonPoll1(f.ctx, f.rs, f.pending, f.unj, 0, max, 0, false, "", paused, hold, machine, &skips)
	if err != nil {
		t.Fatal(err)
	}
	return started
}

func queueWaitLines(t *testing.T, f *daemonFoldFix) []daemon.QueueWait {
	t.Helper()
	common, err := gitCommonDir(f.ctx)
	if err != nil {
		t.Fatal(err)
	}
	return daemon.ReadQueueWaits(common)
}

func wantQueueWaits(t *testing.T, f *daemonFoldFix, want ...string) []daemon.QueueWait {
	t.Helper()
	got := queueWaitLines(t, f)
	var reasons []string
	for _, w := range got {
		if w.Change != f.target || w.Title == "" {
			t.Fatalf("line %+v: want change %s and a title", w, f.target)
		}
		reasons = append(reasons, w.Reason)
	}
	if len(reasons) != len(want) {
		t.Fatalf("reasons %v, want %v", reasons, want)
	}
	for i := range want {
		if reasons[i] != want[i] {
			t.Fatalf("reasons %v, want %v", reasons, want)
		}
	}
	return got
}

func TestQueueWaitJournal(t *testing.T) {
	f := usageFix(t)
	hold := &usageHold{Kind: "session", Used: 900, Projected: 900, Limit: 1000, Reserve: 20, Reset: f.clock.Add(time.Hour)}
	machine := &daemon.MachineHold{Reason: daemon.SkipLoadHigh, Detail: "load 4.0 over 3.0"}

	queueWaitPoll(t, f, 0, nil, nil, false)
	queueWaitPoll(t, f, 0, nil, nil, false)
	queueWaitPoll(t, f, 4, hold, nil, false)
	queueWaitPoll(t, f, 4, hold, nil, false)
	queueWaitPoll(t, f, 4, nil, machine, false)
	got := wantQueueWaits(t, f, "no-slot", "usage-reserve", "load-high")
	if got[1].Detail == "" || got[2].Detail != "load 4.0 over 3.0" {
		t.Fatalf("details %+v", got)
	}
	if started := queueWaitPoll(t, f, 4, nil, nil, false); len(started) != 1 || started[0] != f.target {
		t.Fatalf("started %v", started)
	}
	wantQueueWaits(t, f, "no-slot", "usage-reserve", "load-high", "started")
	// Held, then live, with every hold on during the run: nothing more.
	queueWaitPoll(t, f, 4, nil, nil, false)
	queueWaitPoll(t, f, 4, hold, nil, false)
	queueWaitPoll(t, f, 4, nil, machine, false)
	queueWaitPoll(t, f, 4, nil, nil, true)
	j := run.Journal{ID: "live", Entries: []run.Entry{{Type: run.TypeStarted, At: time.Now(), Started: &run.Started{Change: f.target}}}}
	if err := f.rs.Create(j.ID, j.Entries[0]); err != nil {
		t.Fatal(err)
	}
	queueWaitPoll(t, f, 4, nil, nil, false)
	queueWaitPoll(t, f, 4, hold, nil, false)
	queueWaitPoll(t, f, 4, nil, machine, false)
	queueWaitPoll(t, f, 4, nil, nil, true)
	wantQueueWaits(t, f, "no-slot", "usage-reserve", "load-high", "started")
}

func TestQueueWaitFirstPollStarts(t *testing.T) {
	f := usageFix(t)
	queueWaitPoll(t, f, 4, nil, nil, false)
	wantQueueWaits(t, f, "started")
}

func TestQueueWaitRestartWritesNoDuplicate(t *testing.T) {
	f := usageFix(t)
	queueWaitPoll(t, f, 0, nil, nil, false)
	// A fresh daemon: no memory of the earlier polls but the journal.
	f.pending, f.unj = map[string]daemonPending{}, map[string]daemonPending{}
	queueWaitPoll(t, f, 0, nil, nil, false)
	wantQueueWaits(t, f, "no-slot")
}
