package api

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/quiettest"
)

// TestChangeShowBudget times rhiza change show of the in-flight change sets
// of the 50,000-node model under the read budget, with the index fresh and
// after one entry is appended to an unrelated change set. It runs after
// TestReadBudget (file order) because it appends to the shared model.
func TestChangeShowBudget(t *testing.T) {
	if testing.Short() {
		t.Skip("the read budget generates the 50,000-node model; skipped under -short")
	}
	q := quiettest.Wait(t)
	budgetLock(t)
	timingLog := openTimingLog(t, q.String())
	g := defaultModel(t)
	s := g.summary
	log := func(line string) {
		t.Log(line)
		timingLog(line)
	}
	timeShow := func(label string, cs id.ID) {
		t.Helper()
		ok, err := timeWithinBudget(label, func() (readTiming, error) {
			var times []time.Duration
			for range readBudgetRuns {
				ctx := testCtx(g.root)
				ctx.IndexPath = g.indexPath
				ctx.ChangeRef = cs.String()
				reads := countingCtx(ctx)
				o, _ := Lookup("change.show")
				start := time.Now()
				resp, err := o.Run(ctx, ChangeShowRequest{})
				times = append(times, time.Since(start))
				if err != nil {
					return readTiming{}, fmt.Errorf("%s: %w", label, err)
				}
				if r := resp.(ChangeShowResponse); r.ID != cs || len(r.Ops) == 0 {
					return readTiming{}, fmt.Errorf("%s answered %+v", label, r)
				}
				if *reads != 0 {
					return readTiming{}, fmt.Errorf("%s read the log %d times", label, *reads)
				}
			}
			slices.Sort(times)
			median := (times[readBudgetRuns/2-1] + times[readBudgetRuns/2]) / 2
			return readTiming{median, times[0], times[len(times)-1], 1}, nil
		}, log)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("%s median over %v on every attempt", label, readBudget)
		}
	}
	for _, c := range []struct {
		name string
		cs   id.ID
	}{{"draft", s.Draft}, {"accepted", s.Accepted}} {
		timeShow("change show "+c.name+" fresh", c.cs)
	}
	// Append an entry to an unrelated change set: a new draft.
	ctx := testCtx(g.root)
	ctx.IndexPath = g.indexPath
	if _, err := callOp(t, ctx, "change.new", `{"intent":"Unrelated to the others"}`); err != nil {
		t.Fatal(err)
	}
	// The first read after the append catches the index up. Until the
	// folder's change is two seconds old (the index's racy window) every
	// open lists the log again, so one more read after that stores the
	// fresh stats; the timed reads are the ones that follow.
	warm := testCtx(g.root)
	warm.IndexPath = g.indexPath
	warm.ChangeRef = s.Draft.String()
	for range 2 {
		if _, err := callOp(t, warm, "change.show", `{}`); err != nil {
			t.Fatal(err)
		}
		time.Sleep(2200 * time.Millisecond)
	}
	for _, c := range []struct {
		name string
		cs   id.ID
	}{{"draft", s.Draft}, {"accepted", s.Accepted}} {
		timeShow("change show "+c.name+" after an unrelated entry", c.cs)
	}
}
