package api

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// run.open on a change set with a live run whose process recorded a Rex
// block, named by its ID or a unique prefix, checks the block is still
// there, then focuses it in its session.
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
	live := newChange(t, f.ctx, `{"intent":"Live","title":"Live"}`).ChangeSet.String()
	idle := newChange(t, f.ctx, `{"intent":"Idle","title":"Idle"}`).ChangeSet.String()
	none := newChange(t, f.ctx, `{"intent":"None","title":"None"}`).ChangeSet.String()
	noBlock := newChange(t, f.ctx, `{"intent":"No block","title":"No block"}`).ChangeSet.String()
	journal("LIVE", live, &run.Process{RexSession: "session:s", RexWindow: "window:w", RexBlock: "block:7"})
	hold("LIVE")
	journal("ENDED", idle, &run.Process{RexSession: "session:s", RexBlock: "block:9"})
	journal("NOBLOCK", noBlock, &run.Process{})
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

	wantCalls := [][]string{
		{"rex", "block", "inspect", "-s", "session:s", "--json", "block:7"},
		{"rex", "focus", "-s", "session:s", "block:7"},
	}
	for _, ref := range []string{live, strings.ToLower(live[:len(live)-1])} {
		got := open(ref)
		if want := (RunOpenResponse{Change: live, Live: true, Run: "LIVE", Opened: true, Tab: "block:7"}); got != want {
			t.Errorf("live run by %s: %+v, want %+v", ref, got, want)
		}
		if !slices.EqualFunc(calls, wantCalls, slices.Equal) {
			t.Errorf("live run by %s: rex calls = %q, want %q", ref, calls, wantCalls)
		}
	}

	for _, change := range []string{idle, none} {
		if got := open(change); got != (RunOpenResponse{Change: change}) || len(calls) != 0 {
			t.Errorf("%s: %+v, rex calls %q; want no live run and no calls", change, got, calls)
		}
	}

	if got := open(noBlock); !got.Live || got.Opened || !got.TabGone || len(calls) != 0 {
		t.Errorf("unrecorded block: %+v, rex calls %q; want tab gone and no calls", got, calls)
	}

	inspectFails = true
	if got := open(live); !got.Live || got.Opened || !got.TabGone || len(calls) != 1 {
		t.Errorf("gone block: %+v, rex calls %q; want tab gone after one inspect", got, calls)
	}
}
