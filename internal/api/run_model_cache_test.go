package api

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"syscall"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/gen"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/quiettest"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// driverFixture is a stepRunner over a repository with a large log and its
// index, whose polls the tests count.
type driverFixture struct {
	r        *stepRunner
	ctx      *Ctx
	cs       id.ID
	question id.ID
	sessions *session.Store
	reads    *int
}

// poll does what a driver's tick asks of the model: the leaves, the answer
// to a question as the owner's Decision and as the arbiter's, a verifier's
// Evidence, and the session's status.
func (f *driverFixture) poll(t testing.TB, leaf string) int {
	t.Helper()
	leaves, err := f.r.Leaves(f.cs.String())
	if err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := f.r.LeafAnswer(f.question.String()); err != nil {
		t.Fatal(err)
	}
	if _, err := f.r.Answer(f.question.String(), ""); err != nil {
		t.Fatal(err)
	}
	if leaf != "" {
		if _, _, err := f.r.Evidence(leafOf(leaves, leaf), "S-v"); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := f.r.AgentStatus("S-v"); err != nil {
		t.Fatal(err)
	}
	return len(leaves)
}

func leafOf(leaves []run.Leaf, handle string) run.Leaf {
	for _, l := range leaves {
		if l.Handle == handle {
			return l
		}
	}
	return run.Leaf{}
}

func newDriverFixture(t testing.TB, root, indexPath string, cs, actor id.ID, debounce time.Duration) *driverFixture {
	t.Helper()
	ctx := testCtx(root)
	ctx.IndexPath = indexPath
	ctx.Git = func(args ...string) (string, error) { return "\n", nil }
	reads := countingCtx(ctx)
	sessions := session.NewStore(filepath.Join(t.TempDir(), "sessions"))
	if err := sessions.Write(session.Record{ID: "S-v", Actor: actor.String(), Change: cs.String()}); err != nil {
		t.Fatal(err)
	}
	r := &stepRunner{ctx: ctx, change: cs.String(), run: "R1", worktree: root, sessions: sessions, debounce: debounce}
	t.Cleanup(r.closeModel)
	q, err := ctx.IDs.New()
	if err != nil {
		t.Fatal(err)
	}
	return &driverFixture{r: r, ctx: ctx, cs: cs, question: q, sessions: sessions, reads: reads}
}

// An idle poll does no model work; a poll after an entry in an unrelated
// change set catches up that change set alone, computes no answer again
// and reads the log never; a
// poll after an entry in the run's change set, or a session done mark,
// sees it.
func TestDriverPollIsCheap(t *testing.T) {
	tmp := t.TempDir()
	root := filepath.Join(tmp, "repo")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := gen.Write(root, gen.Options{Seed: 7, Nodes: 800, Links: 300, Landed: 120, InFlight: 6}); err != nil {
		t.Fatal(err)
	}
	setup := testCtx(root)
	setup.IndexPath = filepath.Join(tmp, "index.db")
	cs := newChange(t, setup, `{"intent":"Build"}`).ChangeSet
	mustWrite(t, setup, "add", `{"kind":"capability","handle":"leaf-a"}`)
	mustWrite(t, setup, "add", `{"kind":"capability","handle":"leaf-b"}`)
	ada := recordAda(t, setup, "ada@example.com")

	f := newDriverFixture(t, root, setup.IndexPath, cs, ada, -1)
	if n := f.poll(t, ""); n != 2 {
		t.Fatalf("leaves = %d, want 2", n)
	}
	f.poll(t, "leaf-a")
	rm := f.r.rm
	if *f.reads != 0 {
		t.Fatalf("the first polls read the log %d times", *f.reads)
	}

	// Idle: nothing changed.
	computed := rm.computed
	for range 5 {
		rm.caught = index.CatchUp{Fresh: true}
		f.poll(t, "leaf-a")
		if c := rm.caught; !c.Fresh || c.Read+c.Decoded+c.Folded != 0 {
			t.Fatalf("an idle poll caught up: %+v", c)
		}
	}
	if rm.computed != computed {
		t.Errorf("idle polls computed %d answers, want none", rm.computed-computed)
	}
	if *f.reads != 0 {
		t.Errorf("idle polls read the log %d times", *f.reads)
	}

	// An entry in an unrelated change set.
	other := *setup
	other.ChangeSet = id.ID{}
	other.ChangeRef = ""
	newChange(t, &other, `{"intent":"Unrelated"}`)
	rm.caught = index.CatchUp{Fresh: true}
	f.poll(t, "leaf-a")
	if c := rm.caught; c.Read != 2 && c.Read != 1 || c.Folded != 1 || c.Decoded > 3 {
		t.Errorf("a poll after an unrelated entry caught up %+v, want its change set alone", c)
	}
	// No answer, the Evidence's included, is computed again.
	if n := rm.computed - computed; n != 0 {
		t.Errorf("a poll after an unrelated entry computed %d answers, want none", n)
	}
	rm.caught = index.CatchUp{Fresh: true}
	f.poll(t, "leaf-a")
	if c := rm.caught; c.Read+c.Decoded+c.Folded != 0 {
		t.Errorf("the next poll caught up %+v", c)
	}
	if *f.reads != 0 {
		t.Errorf("polls read the log %d times", *f.reads)
	}

	// An entry in the run's change set and a session done mark.
	setup.ChangeSet = cs
	mustWrite(t, setup, "add", `{"kind":"capability","handle":"leaf-c"}`)
	if n := f.poll(t, ""); n != 3 {
		t.Errorf("leaves after a write = %d, want 3", n)
	}
	rec, err := f.sessions.Read("S-v")
	if err != nil {
		t.Fatal(err)
	}
	rec.Status.Done = &session.Report{State: "done", At: time.Now()}
	if err := f.sessions.Write(rec); err != nil {
		t.Fatal(err)
	}
	st, err := f.r.AgentStatus("S-v")
	if err != nil || !st.SaidDone {
		t.Errorf("AgentStatus = %+v, %v; want SaidDone", st, err)
	}
	if *f.reads != 0 {
		t.Errorf("polls read the log %d times", *f.reads)
	}
}

// cpuTime is this process's CPU time so far.
func cpuTime() time.Duration {
	var ru syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &ru); err != nil {
		return 0
	}
	return time.Duration(ru.Utime.Nano() + ru.Stime.Nano())
}

const (
	idlePollCPUBudget      = 5 * time.Millisecond
	unrelatedPollCPUBudget = 50 * time.Millisecond
)

// driverPollBudget measures the CPU an idle driver poll costs on g, the
// 50,000-node model, and a poll after one unrelated entry. Like
// changeShowBudget it appends to the shared model, so TestReadBudget runs
// it last.
func driverPollBudget(t *testing.T, g generated, log func(string)) {
	t.Helper()
	// The run's change set: a few leaves, as a planned change set has.
	mk := testCtx(g.root)
	mk.IndexPath = g.indexPath
	mk.IDs = id.NewGenerator(mk.Now, rand.New(rand.NewSource(771)))
	cs := newChange(t, mk, `{"intent":"A run under the driver budget"}`).ChangeSet
	for _, h := range []string{"driver-leaf-a", "driver-leaf-b", "driver-leaf-c"} {
		mustWrite(t, mk, "add", fmt.Sprintf(`{"kind":"capability","handle":%q}`, h))
	}
	// The verifier session's actor, whose Evidence each poll judges.
	verifier := recordAda(t, mk, "driver-verifier@example.com")
	f := newDriverFixture(t, g.root, g.indexPath, cs, verifier, 0)
	measure := func(label string, budget time.Duration, runs int, between func()) {
		t.Helper()
		var per []time.Duration
		var laps []quiettest.Lap
		for range runs {
			if between != nil {
				between()
			}
			// A driver polls seconds apart: each poll's first call
			// catches up, the rest of its calls take the index as it is.
			time.Sleep(2 * pollDebounce)
			sw := quiettest.Start()
			start := cpuTime()
			f.poll(t, "driver-leaf-a")
			per = append(per, cpuTime()-start)
			laps = append(laps, sw.Stop())
		}
		slices.Sort(per)
		median := per[len(per)/2]
		// CPU time leaves out waiting for a CPU, but a saturated machine
		// still spends more of it on the same poll: over the budget only
		// then (outside CI), the poll is inconclusive, not failed. Over is
		// at or over its budget; the poll's has always been more than.
		tm := quiettest.Timing{Own: median, Pressure: quiettest.Summarize(laps).Pressure}
		log(fmt.Sprintf("driver poll %-28s median %v CPU (min %v, max %v), budget %v, cpu pressure %.0f%%", label, median, per[0], per[len(per)-1], budget, 100*tm.Pressure))
		if over, why := tm.Over(budget + 1); over {
			t.Errorf("driver poll %s: median %v of CPU over %v", label, median, budget)
		} else if why != "" {
			log("driver poll " + label + " " + why)
		}
	}
	// Warm the answers and let the racy window pass, so the index holds
	// its stat fingerprint. Each poll judges a verifier's Evidence too.
	f.poll(t, "driver-leaf-a")
	time.Sleep(2200 * time.Millisecond)
	f.poll(t, "driver-leaf-a")
	measure("idle", idlePollCPUBudget, readBudgetRuns, nil)

	// An entry in another small change set that is in flight. (A new change
	// set would have the index compute its views again, and the change sets
	// of the generated model are large to fold: no poll could pay for
	// either.)
	other := testCtx(g.root)
	other.IndexPath = g.indexPath
	other.IDs = id.NewGenerator(other.Now, rand.New(rand.NewSource(991)))
	other.ChangeSet = newChange(t, other, `{"intent":"Unrelated to the driver"}`).ChangeSet
	f.poll(t, "driver-leaf-a")
	n := 0
	measure("after an unrelated entry", unrelatedPollCPUBudget, 9, func() {
		n++
		req := fmt.Sprintf(`{"kind":"capability","handle":"unrelated-%d"}`, n)
		if _, err := callOp(t, other, "add", req); err != nil {
			t.Fatal(err)
		}
		// The catch-up of an entry lists the whole log, as any reader
		// opening the index does (the index's own cost, timed by its
		// budget): another reader does it here. Until the folder's
		// change is two seconds old (the racy window) every open lists
		// the log again, so a second open after that stores the stats.
		for range 2 {
			x, err := index.Open(g.root, g.indexPath)
			if err != nil {
				t.Fatal(err)
			}
			x.Close()
			time.Sleep(2200 * time.Millisecond)
		}
	})
	if *f.reads != 0 {
		t.Errorf("driver polls read the log %d times", *f.reads)
	}
}
