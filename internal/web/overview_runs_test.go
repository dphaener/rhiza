package web_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/web"
)

// The Overview at / reads live runs from the repository's run journals: a
// running run shows its change set's title, workflow, current step and a
// link to its page, and a finished run is left out.
func TestOverviewListsLiveRunsFromJournals(t *testing.T) {
	ctx := fixtureCtx(t)
	rs := run.NewStore(run.Dir(filepath.Join(ctx.Root, ".git")))
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
	mk("01LIVE", "build", "Retry webhook delivery", run.Entry{Type: run.TypeEntered, Step: "leaves", Visit: 1})
	mk("01DONE", "plan", "Already shipped", run.Entry{Type: run.TypeFinished, Detail: "ok"})

	body := serve(web.HandlerWithRuns(nil, api.WebRunsSource(ctx), nil), "GET", "/").Body.String()
	for _, w := range []string{
		`href="/runs/01LIVE"`, "Retry webhook delivery", "<small>build</small>", ">leaves</span>",
	} {
		if !strings.Contains(body, w) {
			t.Errorf("/ lacks %q:\n%s", w, body)
		}
	}
	for _, w := range []string{"01DONE", "Already shipped"} {
		if strings.Contains(body, w) {
			t.Errorf("/ lists the finished run (%q):\n%s", w, body)
		}
	}
}
