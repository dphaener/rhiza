package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/quiettest"
)

// changeShowBudget times rhiza change show of the in-flight change sets of
// g, the 50,000-node model, under the read budget, with the index fresh and
// after one entry is appended to an unrelated change set, logging each
// attempt through log. TestReadBudget runs it last, since it appends to
// the shared model.
func changeShowBudget(t *testing.T, g generated, log func(string)) {
	t.Helper()
	s := g.summary
	timeShow := func(label string, cs id.ID) {
		t.Helper()
		ok, err := timeWithinBudget(label, func() (readTiming, error) {
			var laps []quiettest.Lap
			for range readBudgetRuns {
				ctx := testCtx(g.root)
				ctx.IndexPath = g.indexPath
				ctx.ChangeRef = cs.String()
				reads := countingCtx(ctx)
				o, _ := Lookup("change.show")
				sw := quiettest.Start()
				resp, err := o.Run(ctx, ChangeShowRequest{})
				laps = append(laps, sw.Stop())
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
			return readTiming{quiettest.Summarize(laps), 1}, nil
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
