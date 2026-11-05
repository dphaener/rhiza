package api

import (
	"bytes"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

func interviewJournal(id string, started run.Started, ended bool) run.Journal {
	at := time.Date(2026, 11, 5, 6, 0, 0, 0, time.UTC)
	j := run.Journal{ID: id, Entries: []run.Entry{
		{Type: run.TypeStarted, At: at, Step: "interview", Started: &started},
		{Type: run.TypeSession, At: at.Add(time.Second), Step: "interview", Session: "S1"},
	}}
	if ended {
		j.Entries = append(j.Entries, run.Entry{Type: run.TypeEnded, At: at.Add(time.Minute), Step: "interview", Outcome: "pass"}, run.Entry{Type: run.TypeFinished, At: at.Add(time.Minute), End: run.Succeeded})
	}
	return j
}

// A live milestone request's plan run is an interview item from the
// interview's launch, with its tab, before any question block; it is gone
// once the run ended, and a feature's plan run has none.
func TestInboxMilestoneInterview(t *testing.T) {
	tab := func(sid string) session.Tab {
		return session.Tab{Session: "rs", Block: "blk", Title: "interview " + sid}
	}
	items := func(j run.Journal, live bool) []InboxItem {
		return inboxItems([]run.Journal{j}, "h", nil, func(string) bool { return live }, func(string) (string, bool) { return "", false }, nil, tab)
	}
	ms := run.Started{Change: "CS1", MilestonePlan: true}
	got := items(interviewJournal("R1", ms, false), true)
	if len(got) != 1 {
		t.Fatalf("items %+v, want one interview", got)
	}
	it := got[0]
	if it.Kind != InboxInterview || it.Change != "CS1" || it.Run != "R1" || it.Session != "S1" || it.Tab != "interview S1" || it.TabBlock != "blk" || it.TabSession != "rs" {
		t.Errorf("item %+v", it)
	}
	if rs, blk := itemTab(it, run.State{}); rs != "rs" || blk != "blk" {
		t.Errorf("itemTab = %q %q, want the interview's tab", rs, blk)
	}
	var out bytes.Buffer
	InboxResponse{Items: got}.Human(&out)
	if out.Len() == 0 {
		t.Error("no human output")
	}
	if got := items(interviewJournal("R1", ms, true), true); len(got) != 0 {
		t.Errorf("ended run: %+v, want none", got)
	}
	if got := items(interviewJournal("R1", run.Started{Change: "CS1"}, false), true); len(got) != 0 {
		t.Errorf("feature plan run: %+v, want none", got)
	}
	if got := items(interviewJournal("R1", run.Started{Change: "CS1", MilestonePlan: true, Autonomous: true}, false), true); len(got) != 0 {
		t.Errorf("autonomous run: %+v, want none", got)
	}
}

// The daemon notifies a milestone interview once, with the interview text.
func TestDaemonNotifyMilestoneInterview(t *testing.T) {
	ctx := testCtx(t.TempDir())
	var rex bytes.Buffer
	ctx.Rex = &rex
	rs := run.NewStore(t.TempDir())
	j := interviewJournal("R1", run.Started{Change: "CS1", Title: "Plan it", MilestonePlan: true}, false)
	if err := rs.Create("R1", j.Entries[0]); err != nil {
		t.Fatal(err)
	}
	if err := rs.Append("R1", j.Entries[1]); err != nil {
		t.Fatal(err)
	}
	unlock, ok, err := rs.TryLock("R1")
	if err != nil || !ok {
		t.Fatalf("lock: %v %v", ok, err)
	}
	defer unlock()
	var state daemon.State
	for range 2 {
		if err := daemonNotify(ctx, rs, &state, ""); err != nil {
			t.Fatal(err)
		}
	}
	var blocked []map[string]string
	for _, r := range rexReports(t, rex.String()) {
		if r["state"] == "blocked" {
			blocked = append(blocked, r)
		}
	}
	if len(blocked) != 1 || blocked[0]["msg"] != "Plan it: a milestone interview waits on you" {
		t.Fatalf("blocked reports %v, want one interview notice", blocked)
	}
}
