package api

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/adapter/claudecode"
	"github.com/dphaener/rhiza/internal/quiettest"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

// idleBudget is the most an unchanged root poll may cost in process CPU:
// 0.05 core at the 5 s poll is 250 ms; the budget is well inside that.
const idleBudget = 100 * time.Millisecond

// idleCounts is what a poll did: log reads (both projects), plan reads, the
// change sets folded, the picks computed and the transcript files read.
type idleCounts struct{ logReads, planReads, folds, picks, transcripts int64 }

func (f *idleFix) counts() idleCounts {
	_, reads := claudecode.UsageScanStats()
	return idleCounts{
		logReads:    int64(f.home.logReads + f.other.logReads),
		planReads:   planReadCount.Load(),
		folds:       planFoldCount.Load(),
		picks:       daemonPickCount.Load(),
		transcripts: reads,
	}
}

func (c idleCounts) minus(o idleCounts) idleCounts {
	return idleCounts{c.logReads - o.logReads, c.planReads - o.planReads, c.folds - o.folds, c.picks - o.picks, c.transcripts - o.transcripts}
}

// settle makes the log's folders old news again, so their stamp is trusted.
func (f *idleFix) settle(t testing.TB) {
	t.Helper()
	old := time.Now().Add(-time.Hour)
	f.home.backdate(t, old)
	f.other.backdate(t, old)
}

// TestRootIdlePoll runs the root daemon's poll over two projects with
// nothing changed: the second poll folds nothing and computes no pick. Each
// stamp then changes in turn, and the next poll does the work again.
func TestRootIdlePoll(t *testing.T) {
	if testing.Short() {
		t.Skip("builds the production-sized idle fixture; skipped under -short")
	}
	f := newIdleFix(t)
	f.rootPoll(t)
	first := f.counts()
	if first.logReads == 0 || first.planReads == 0 || first.picks == 0 {
		t.Fatalf("first poll did no work: %+v", first)
	}
	f.rootPoll(t)
	if d := f.counts().minus(first); d != (idleCounts{}) {
		t.Fatalf("the second poll did work with nothing changed: %+v", d)
	}
	f.rootPoll(t)
	if d := f.counts().minus(first); d != (idleCounts{}) {
		t.Fatalf("the third poll did work with nothing changed: %+v", d)
	}

	step := 0
	changes := []struct {
		name string
		do   func()
		want func(d idleCounts) string
	}{
		{"log entry", func() {
			f.home.call("change.new", fmt.Sprintf(`{"intent":"Fresh %d"}`, step))
			f.settle(t)
		}, func(d idleCounts) string {
			if d.logReads == 0 || d.planReads == 0 || d.picks == 0 {
				return "want a log read, a plan read and a pick"
			}
			return ""
		}},
		{"ref move", func() { f.other.refs = fmt.Sprintf("%040d refs/heads/main\n", step) }, func(d idleCounts) string {
			if d.logReads == 0 || d.planReads == 0 || d.picks == 0 {
				return "want a log read, a plan read and a pick"
			}
			return ""
		}},
		{"setting", func() { f.home.config["rhiza.maxRuns"] = fmt.Sprint(step) }, func(d idleCounts) string {
			if d.logReads == 0 || d.planReads == 0 || d.picks == 0 {
				return "want a log read, a plan read and a pick"
			}
			return ""
		}},
		{"run record", func() {
			src := "[step.s]\nkind = \"agent\"\nprompt = \"work\"\n"
			def, err := workflow.Parse("wf", []byte(src))
			if err != nil {
				t.Fatal(err)
			}
			if err := run.Start(f.projs[0].rs, fmt.Sprintf("N%d", step), run.Started{Workflow: "wf", Definition: src, Change: f.home.target}, f.home.ctx.Now(), def); err != nil {
				t.Fatal(err)
			}
		}, func(d idleCounts) string {
			if d.picks == 0 {
				return "want a pick"
			}
			return ""
		}},
		{"transcript", func() {
			p := firstTranscript(t, f)
			fh, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, 0o644)
			if err != nil {
				t.Fatal(err)
			}
			defer fh.Close()
			ts := f.home.ctx.Now().UTC().Format(time.RFC3339)
			fmt.Fprintf(fh, `{"type":"assistant","timestamp":%q,"requestId":"req-new","message":{"id":"msg-new","model":"claude-sonnet-5-5","usage":{"input_tokens":5,"output_tokens":5}}}`+"\n", ts)
		}, func(d idleCounts) string {
			if d.transcripts == 0 {
				return "want a transcript read"
			}
			return ""
		}},
	}
	for _, c := range changes {
		step++
		c.do()
		before := f.counts()
		f.rootPoll(t)
		if why := c.want(f.counts().minus(before)); why != "" {
			t.Errorf("after a %s: %s, got %+v", c.name, why, f.counts().minus(before))
		}
		// And the poll after that is idle again.
		f.settle(t)
		f.rootPoll(t) // may refold once more: a stamp taken in a racy window
		before = f.counts()
		f.rootPoll(t)
		if d := f.counts().minus(before); d != (idleCounts{}) {
			t.Errorf("a poll after the %s did work with nothing changed: %+v", c.name, d)
		}
	}
}

func firstTranscript(t testing.TB, f *idleFix) string {
	t.Helper()
	config := f.home.ctx.Getenv("CLAUDE_CONFIG_DIR")
	var found string
	_ = filepath.WalkDir(filepath.Join(config, "projects"), func(p string, d os.DirEntry, err error) error {
		if err == nil && found == "" && filepath.Ext(p) == ".jsonl" {
			found = p
		}
		return nil
	})
	if found == "" {
		t.Fatal("no transcript")
	}
	return found
}

// TestRootIdlePollBudget requires that an unchanged root poll over the
// idle fixture costs under idleBudget of process CPU (median).
func TestRootIdlePollBudget(t *testing.T) {
	if testing.Short() {
		t.Skip("builds the production-sized idle fixture; skipped under -short")
	}
	quiettest.Wait(t)
	budgetLock(t)
	f := newIdleFix(t)
	f.rootPoll(t)
	f.rootPoll(t)
	attempts := 1 + readBudgetRetries
	inconclusive := ""
	for n := 1; n <= attempts; n++ {
		var laps []quiettest.Lap
		for range readBudgetRuns {
			sw := quiettest.Start()
			f.rootPoll(t)
			laps = append(laps, sw.Stop())
		}
		tm := quiettest.Summarize(laps)
		t.Logf("idle root poll attempt %d/%d: %v", n, attempts, tm)
		over, why := tm.Over(idleBudget)
		if !over && why == "" {
			return
		}
		if why != "" {
			inconclusive = why
		}
	}
	if inconclusive != "" {
		t.Log("idle root poll " + inconclusive)
		return
	}
	t.Fatalf("an unchanged root poll took over %v of own time on every attempt", idleBudget)
}
