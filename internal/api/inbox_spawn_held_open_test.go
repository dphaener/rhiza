package api

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// inbox.open finds a spawn-held item by the ID rhiza inbox lists. A held
// spawn or step session has no session yet, so it opens the run's own
// journalled tab when the run has one, and is refused, like the
// machine-wide items, when it has none.
func TestInboxOpenSpawnHeld(t *testing.T) {
	t.Parallel()
	f := newSessionFixture(t, "x")
	rs := run.NewStore(run.Dir(f.common))
	at := time.Unix(1700000000, 0).UTC()
	const disk = "machine hold: disk-low: 1.0 GiB free, under the minimum of 2.0 GiB"
	journal := func(id string, p *run.Process) {
		t.Helper()
		es := []run.Entry{
			{Type: run.TypeStarted, At: at, Step: "build", Started: &run.Started{Workflow: "w", Change: "CS-" + id}},
			{Type: run.TypeEntered, At: at, Step: "build", Visit: 1},
			{Type: run.TypeSpawnHeld, At: at.Add(time.Second), Step: "build", Detail: disk},
		}
		if p != nil {
			es = append(es, run.Entry{Type: run.TypeProcess, At: at, Process: p})
		}
		if err := rs.Create(id, es[0]); err != nil {
			t.Fatal(err)
		}
		for _, e := range es[1:] {
			if err := rs.Append(id, e); err != nil {
				t.Fatal(err)
			}
		}
		unlock, ok, err := rs.TryLock(id)
		if err != nil || !ok {
			t.Fatalf("lock %s: %v %v", id, ok, err)
		}
		t.Cleanup(unlock)
	}
	journal("NOTAB", nil)
	journal("TAB", &run.Process{RexSession: "session:s", RexBlock: "block:7"})

	var calls [][]string
	f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		calls = append(calls, append([]string{name}, args...))
		return []byte(`{"block":{"id":"block:7"}}`), nil
	}
	list, err := callOp(t, f.ctx, "inbox", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]string{}
	for _, it := range list.(InboxResponse).Items {
		if it.Kind == InboxSpawnHeld {
			ids[it.Run] = it.ID
		}
	}
	if ids["NOTAB"] == "" || ids["TAB"] == "" {
		t.Fatalf("spawn-held items %v, want one per run", ids)
	}

	_, err = callOp(t, f.ctx, "inbox.open", `{"item":"`+ids["NOTAB"]+`"}`)
	var e *Error
	if !errors.As(err, &e) || e.Code != CodeUsage || !strings.Contains(e.Message, "has no tab to open") {
		t.Fatalf("open a held step with no run tab: %v, want usage: no tab to open", err)
	}
	if len(calls) != 0 {
		t.Errorf("refused open called %v", calls)
	}

	resp, err := callOp(t, f.ctx, "inbox.open", `{"item":"`+ids["TAB"]+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	want := InboxOpenResponse{Item: ids["TAB"], Kind: InboxSpawnHeld, Run: "TAB", Opened: true, Tab: "block:7"}
	if got := resp.(InboxOpenResponse); got != want {
		t.Errorf("open a held step in a run with a tab: %+v, want %+v", got, want)
	}
	wantCalls := [][]string{
		{"rex", "block", "inspect", "-s", "session:s", "--json", "block:7"},
		{"rex", "focus", "-s", "session:s", "block:7"},
	}
	if !slices.EqualFunc(calls, wantCalls, slices.Equal) {
		t.Errorf("calls %v, want %v", calls, wantCalls)
	}
}
