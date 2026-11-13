package api

import (
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// A change set a plan run gated and then journaled accepted raises no accept
// item, even while main's model still says proposed; one the run did not
// accept keeps its item.
func TestInboxSkipsRunAcceptedGated(t *testing.T) {
	at := time.Date(2026, 11, 13, 14, 0, 0, 0, time.UTC)
	journal := func(extra ...run.Entry) run.Journal {
		es := []run.Entry{
			{Type: run.TypeStarted, At: at, Started: &run.Started{Workflow: "plan-bug", Change: "C1", Worktree: t.TempDir(), Plan: true}},
			{Type: run.TypeGated, At: at.Add(time.Second), Changes: []string{"C1", "C2"}},
		}
		return run.Journal{ID: "R1", Entries: append(es, extra...)}
	}
	proposed := func(string) (string, bool) { return "proposed", true }
	accepts := func(j run.Journal) []string {
		var out []string
		for _, it := range inboxItems([]run.Journal{j}, "h", nil, func(string) bool { return false }, proposed, nil, nil, nil, nil) {
			if it.Kind == InboxAccept {
				if it.Next != ownerAcceptNext(it.Change) {
					t.Errorf("next = %q", it.Next)
				}
				out = append(out, it.Change)
			}
		}
		return out
	}
	if got := accepts(journal()); len(got) != 2 {
		t.Fatalf("no accepted entry lists %v, want both", got)
	}
	acc := run.Entry{Type: run.TypeAccepted, At: at.Add(2 * time.Second), Change: "C1", LogEntry: "E1"}
	if got := accepts(journal(acc)); len(got) != 1 || got[0] != "C2" {
		t.Fatalf("one accepted lists %v, want [C2]", got)
	}
	rej := run.Entry{Type: run.TypeRejected, At: at.Add(2 * time.Second), Change: "C1", LogEntry: "E1"}
	if got := accepts(journal(rej)); len(got) != 2 {
		t.Fatalf("a rejected entry lists %v, want both", got)
	}
}
