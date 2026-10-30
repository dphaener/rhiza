package api

import (
	"slices"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// run.open on a change set with a live run whose process recorded a Rex
// block checks the block is still there, then focuses it in its session.
// A change set with no live run, or a run that is not live, calls nothing
// and says so; a live run whose block is gone or unrecorded says its tab
// is gone and focuses nothing.
func TestRunOpenFocusesLiveRunTab(t *testing.T) {
	f := newSessionFixture(t, "x")
	rs := run.NewStore(run.Dir(f.common))
	at := time.Unix(1700000000, 0).UTC()
	journal := func(id, change string, p *run.Process) {
		t.Helper()
		if err := rs.Create(id, run.Entry{Type: run.TypeStarted, At: at, Step: "s", Started: &run.Started{Workflow: "w", Change: change}}); err != nil {
			t.Fatal(err)
		}
		if p != nil {
			if err := rs.Append(id, run.Entry{Type: run.TypeProcess, At: at, Process: p}); err != nil {
				t.Fatal(err)
			}
		}
	}
	hold := func(id string) {
		t.Helper()
		unlock, ok, err := rs.TryLock(id)
		if err != nil || !ok {
			t.Fatalf("lock %s: %v %v", id, ok, err)
		}
		t.Cleanup(unlock)
	}
	journal("LIVE", "CSLIVE", &run.Process{RexSession: "session:s", RexWindow: "window:w", RexBlock: "block:7"})
	hold("LIVE")
	journal("ENDED", "CSIDLE", &run.Process{RexSession: "session:s", RexBlock: "block:9"})
	journal("NOBLOCK", "CSNOBLOCK", &run.Process{})
	hold("NOBLOCK")

	var calls [][]string
	inspectFails := false
	f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		calls = append(calls, append([]string{name}, args...))
		if inspectFails && args[0] == "block" {
			return nil, Errorf("no such block")
		}
		return []byte(`{"block":{"id":"block:7"}}`), nil
	}
	open := func(change string) RunOpenResponse {
		t.Helper()
		calls = nil
		resp, err := callOp(t, f.ctx, "run.open", `{"change":"`+change+`"}`)
		if err != nil {
			t.Fatal(err)
		}
		return resp.(RunOpenResponse)
	}

	got := open("CSLIVE")
	if want := (RunOpenResponse{Change: "CSLIVE", Live: true, Run: "LIVE", Opened: true, Tab: "block:7"}); got != want {
		t.Errorf("live run: %+v, want %+v", got, want)
	}
	wantCalls := [][]string{
		{"rex", "block", "inspect", "-s", "session:s", "--json", "block:7"},
		{"rex", "focus", "-s", "session:s", "block:7"},
	}
	if !slices.EqualFunc(calls, wantCalls, slices.Equal) {
		t.Errorf("live run rex calls = %q, want %q", calls, wantCalls)
	}

	for _, change := range []string{"CSIDLE", "CSNONE"} {
		if got := open(change); got != (RunOpenResponse{Change: change}) || len(calls) != 0 {
			t.Errorf("%s: %+v, rex calls %q; want no live run and no calls", change, got, calls)
		}
	}

	if got := open("CSNOBLOCK"); !got.Live || got.Opened || !got.TabGone || len(calls) != 0 {
		t.Errorf("unrecorded block: %+v, rex calls %q; want tab gone and no calls", got, calls)
	}

	inspectFails = true
	if got := open("CSLIVE"); !got.Live || got.Opened || !got.TabGone || len(calls) != 1 {
		t.Errorf("gone block: %+v, rex calls %q; want tab gone after one inspect", got, calls)
	}
}
