package api

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/quiettest"
)

// applyBudget is the most the median of one rhiza apply write may take on
// the 50,000-node model.
const applyBudget = 200 * time.Millisecond

// TestReadBudgetApply times one rhiza apply write (a set of the deep node)
// into a draft change set of the 50,000-node model, index open to answer,
// and requires that it reads the log not at all. It runs after
// TestReadBudget (file order) because it appends to the shared model.
func TestReadBudgetApply(t *testing.T) {
	if testing.Short() {
		t.Skip("the read budget generates the 50,000-node model; skipped under -short")
	}
	q := quiettest.Wait(t)
	budgetLock(t)
	timingLog := openTimingLog(t, q.String())
	g := defaultModel(t)
	file := filepath.Join(t.TempDir(), "writes.json")
	apply := func(i int) (time.Duration, error) {
		body := fmt.Sprintf(`{"writes":[{"op":"set","ref":%q,"fields":["statement=Budget %d"]}]}`, g.summary.Deep.ID.String(), i)
		if err := os.WriteFile(file, []byte(body), 0o644); err != nil {
			return 0, err
		}
		ctx := testCtx(g.root)
		ctx.IndexPath = g.indexPath
		ctx.ChangeRef = g.summary.Draft.String()
		seededIDs(ctx, int64(1000+i))
		reads := countingCtx(ctx)
		start := time.Now()
		_, err := callOp(t, ctx, "apply", fmt.Sprintf(`{"file":%q}`, file))
		d := time.Since(start)
		if err == nil && *reads != 0 {
			err = fmt.Errorf("apply read the log %d times", *reads)
		}
		return d, err
	}
	// The first write may record an actor; two more, a racy window apart,
	// leave the index's stats fresh.
	for i := range 3 {
		if _, err := apply(i); err != nil {
			t.Fatal(err)
		}
		time.Sleep(2200 * time.Millisecond)
	}
	// Load from outside the test run comes in spikes: the write fails only
	// if every attempt's median is over the budget.
	for attempt := 1; attempt <= 1+readBudgetRetries; attempt++ {
		var times []time.Duration
		for i := range readBudgetRuns {
			d, err := apply(10 + attempt*100 + i)
			if err != nil {
				t.Fatal(err)
			}
			times = append(times, d)
			time.Sleep(2200 * time.Millisecond)
		}
		slices.Sort(times)
		median := (times[readBudgetRuns/2-1] + times[readBudgetRuns/2]) / 2
		line := fmt.Sprintf("apply one write attempt %d: median %v (min %v, max %v)", attempt, median, times[0], times[len(times)-1])
		t.Log(line)
		timingLog(line)
		if median < applyBudget {
			return
		}
	}
	t.Errorf("apply median over %v on every attempt", applyBudget)
}
