package api

import (
	"testing"

	"github.com/dphaener/rhiza/internal/run"
)

// The Overview lists unfinished runs by their change set's title, workflow
// and step, and leaves finished runs out.
func TestWebRunsSourceListsLiveRuns(t *testing.T) {
	ctx, common := daemonRepo(t)
	rs := run.NewStore(run.Dir(common))
	mk := func(id, wf, title string, es ...run.Entry) {
		t.Helper()
		first := run.Entry{Type: run.TypeStarted, Started: &run.Started{Workflow: wf, Change: "C-" + id, Title: title}}
		if err := rs.Create(id, first); err != nil {
			t.Fatal(err)
		}
		for _, e := range es {
			if err := rs.Append(id, e); err != nil {
				t.Fatal(err)
			}
		}
	}
	mk("R1", "build", "Retry webhook delivery", run.Entry{Type: run.TypeEntered, Step: "leaves", Visit: 1})
	mk("R2", "build", "Already done", run.Entry{Type: run.TypeFinished, Detail: "ok"})
	got := WebRunsSource(ctx)()
	if len(got) != 1 {
		t.Fatalf("runs = %+v", got)
	}
	r := got[0]
	if r.ID != "R1" || r.Title != "Retry webhook delivery" || r.Workflow != "build" || r.Step != "leaves" {
		t.Errorf("run = %+v", r)
	}
}
