package api

import (
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// A project's needs-you count is its open questions, sign-offs, failed runs
// and stuck runs; a hold is in no project's count.
func TestProjectNeedsYou(t *testing.T) {
	now := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)
	at := now.Add(-3 * time.Hour)
	journal := func(id string, es ...run.Entry) run.Journal {
		started := run.Entry{Type: run.TypeStarted, At: at, Started: &run.Started{Workflow: "build", Change: "C-" + id}}
		return run.Journal{ID: id, Entries: append([]run.Entry{started}, es...)}
	}
	enter := func(step string) run.Entry { return run.Entry{Type: run.TypeEntered, At: at, Step: step, Visit: 1} }
	agent := func(step string, last time.Time) []run.Entry {
		return []run.Entry{enter(step), {Type: run.TypeSession, At: last, Step: step, Session: "S"}}
	}
	base := func() []run.Journal {
		return []run.Journal{
			journal("ROPEN", enter("gate"), run.Entry{Type: run.TypeStopped, At: at, Step: "gate", Detail: run.OpenStop + "q1; answer it"}),
			journal("RSIGN", enter("signoff"), run.Entry{Type: run.TypeStopped, At: at, Step: "signoff", Detail: run.SignoffWait}),
			journal("RFAIL", enter("open-pr"),
				run.Entry{Type: run.TypeEnded, At: at, Step: "open-pr", Visit: 1, Outcome: "fail", Detail: "exit 3"},
				run.Entry{Type: run.TypeFinished, At: at, Step: "open-pr", Detail: "failed", End: run.Failed}),
		}
	}
	hold := InboxItem{Kind: InboxDiskLow, At: at.Format(time.RFC3339)}
	count := func(extra run.Journal, withExtra bool) int {
		runs := base()
		live := map[string]bool{}
		if withExtra {
			runs = append(runs, extra)
			live[extra.ID] = true
		}
		items := inboxItems(runs, "h", nil, func(id string) bool { return live[id] }, func(string) (string, bool) { return "", false }, nil, nil, nil, nil)
		items = append(items, hold)
		return len(projectNeedsYou(items, runs, func(id string) bool { return live[id] }, now))
	}
	if got := count(run.Journal{}, false); got != 3 {
		t.Fatalf("question, sign-off, failed run and a hold: count %d, want 3", got)
	}
	if got := count(journal("RSTUCK", agent("leaves", now.Add(-21*time.Minute))...), true); got != 4 {
		t.Errorf("a run 21 minutes quiet is stuck: count %d, want 4", got)
	}
	if got := count(journal("RQUIET", agent("leaves", now.Add(-19*time.Minute))...), true); got != 3 {
		t.Errorf("a run 19 minutes quiet is not stuck: count %d, want 3", got)
	}
	if got := count(journal("RFRESH", agent("leaves", now.Add(-5*time.Minute))...), true); got != 3 {
		t.Errorf("a run 5 minutes quiet adds nothing: count %d, want 3", got)
	}
	if got := len(projectNeedsYou(nil, nil, nil, now)); got != 0 {
		t.Errorf("an empty project: count %d, want 0", got)
	}
}
