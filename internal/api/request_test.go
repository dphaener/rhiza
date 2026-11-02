package api

import (
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

// A plan request whose run never journals (run start failed in its tab
// before journaling) goes back in the queue once daemonStartGrace passed,
// and is started again, with a new run, daemonUnjournaledHold after its
// first start.
func TestDaemonPlanRequestsRequeuesUnjournaled(t *testing.T) {
	f := newSessionFixture(t, "req")
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	f.ctx.Now = func() time.Time { return now }
	if _, err := callOp(t, f.ctx, "request", `{"text":"Dark mode"}`); err != nil {
		t.Fatal(err)
	}
	rs := run.NewStore(run.Dir(f.common))
	poll := func() []string {
		t.Helper()
		f.rex.calls = nil
		got, err := daemonPlanRequests(f.ctx, rs, nil, 0, 2, "auto")
		if err != nil {
			t.Fatal(err)
		}
		return got
	}
	request := func() daemon.Request {
		t.Helper()
		reqs, err := daemon.ReadRequests(f.common)
		if err != nil || len(reqs) != 1 {
			t.Fatalf("requests %v, %v", reqs, err)
		}
		return reqs[0]
	}
	if got := poll(); len(got) != 1 || len(f.rex.calls) != 1 {
		t.Fatalf("first poll started %v with rex calls %q, want one start", got, f.rex.calls)
	}
	first := request()
	if first.Run == "" {
		t.Fatal("started request has no run")
	}
	now = now.Add(daemonStartGrace - time.Second)
	if got := poll(); len(got) != 0 || request().Run != first.Run {
		t.Fatalf("within the grace: started %v, run %q, want none and %q", got, request().Run, first.Run)
	}
	now = now.Add(time.Second)
	if got := poll(); len(got) != 0 || len(f.rex.calls) != 0 {
		t.Fatalf("after the grace: started %v, want none until the hold passed", got)
	}
	q := request()
	if q.Run != "" || !q.RetryAt.Equal(first.StartedAt.Add(daemonUnjournaledHold)) {
		t.Fatalf("after the grace: run %q, retry at %v, want back in the queue until %v", q.Run, q.RetryAt, first.StartedAt.Add(daemonUnjournaledHold))
	}
	if state, _ := daemon.RequestState(q, nil); state != daemon.RequestQueued {
		t.Errorf("after the grace: state %s, want queued", state)
	}
	now = first.StartedAt.Add(daemonUnjournaledHold)
	if got := poll(); len(got) != 1 || len(f.rex.calls) != 1 {
		t.Fatalf("after the hold: started %v, want the request again", got)
	}
	if q := request(); q.Run == "" || q.Run == first.Run || !q.RetryAt.IsZero() {
		t.Errorf("after the hold: run %q (first %q), retry at %v, want a new run", q.Run, first.Run, q.RetryAt)
	}
}
