package api

import (
	"fmt"
	"os"
	"path/filepath"
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
	apply := func(i int) (quiettest.Lap, error) {
		body := fmt.Sprintf(`{"writes":[{"op":"set","ref":%q,"fields":["statement=Budget %d"]}]}`, g.summary.Deep.ID.String(), i)
		if err := os.WriteFile(file, []byte(body), 0o644); err != nil {
			return quiettest.Lap{}, err
		}
		ctx := testCtx(g.root)
		ctx.IndexPath = g.indexPath
		ctx.ChangeRef = g.summary.Draft.String()
		seededIDs(ctx, int64(1000+i))
		reads := countingCtx(ctx)
		sw := quiettest.Start()
		_, err := callOp(t, ctx, "apply", fmt.Sprintf(`{"file":%q}`, file))
		d := sw.Stop()
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
	// if every attempt's median own time is over the budget, and is
	// inconclusive if one was over only on a saturated machine
	// (quiettest.Timing.Over).
	inconclusive := ""
	for attempt := 1; attempt <= 1+readBudgetRetries; attempt++ {
		var laps []quiettest.Lap
		for i := range readBudgetRuns {
			d, err := apply(10 + attempt*100 + i)
			if err != nil {
				t.Fatal(err)
			}
			laps = append(laps, d)
			time.Sleep(2200 * time.Millisecond)
		}
		tm := quiettest.Summarize(laps)
		line := fmt.Sprintf("apply one write attempt %d: %v", attempt, tm)
		t.Log(line)
		timingLog(line)
		over, why := tm.Over(applyBudget)
		if !over && why == "" {
			return
		}
		if why != "" {
			inconclusive = why
		}
	}
	if inconclusive != "" {
		t.Log("apply one write " + inconclusive)
		timingLog("apply one write " + inconclusive)
		return
	}
	t.Errorf("apply median own time over %v on every attempt", applyBudget)
}
