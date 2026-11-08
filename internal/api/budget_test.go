package api

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/gen"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/quiettest"
	"github.com/dphaener/rhiza/internal/rootguard"
	"github.com/dphaener/rhiza/internal/tmpsweep"
)

// readBudget is the most a read's median may take, fresh workspace and
// index open to answer: milestone 1's target for every read on the
// 50,000-node model (D14, as revised).
const readBudget = 50 * time.Millisecond

// readBudgetRuns is how many fresh runs TestReadBudget times per read.
const readBudgetRuns = 20

// readBudgetRetries is how many more times a read whose median is over
// readBudget is timed again. Load from outside the test run comes in short
// spikes; the read fails only if every attempt is over.
const readBudgetRetries = 2

// readTiming is one attempt's timings of a read.
type readTiming struct {
	median, min, max time.Duration
	size             int
}

// timeWithinBudget calls measure up to 1+readBudgetRetries times, stopping
// at the first attempt whose median is under readBudget, and logs every
// attempt through log. It reports whether some attempt was under budget.
func timeWithinBudget(label string, measure func() (readTiming, error), log func(string)) (bool, error) {
	attempts := 1 + readBudgetRetries
	for n := 1; n <= attempts; n++ {
		rt, err := measure()
		if err != nil {
			return false, err
		}
		log(fmt.Sprintf("%s attempt %d/%d: median %v (min %v, max %v), size %d", label, n, attempts, rt.median, rt.min, rt.max, rt.size))
		if rt.median < readBudget {
			return true, nil
		}
	}
	return false, nil
}

// generated is gen.Default() written once, with its index built, shared by
// TestReadBudget and the benchmarks. It lives under one os.MkdirTemp folder
// that TestMain removes.
type generated struct {
	root      string
	indexPath string
	summary   gen.Summary
	// subtree is the deep subtree tree reads (D14): the deep sample's
	// grandparent, two levels above the deepest.
	subtree id.ID
	err     error
}

var (
	genOnce sync.Once
	genDir  string
	genned  generated
)

func TestMain(m *testing.M) {
	// A run killed by a go test timeout never reached the RemoveAll below.
	tmpsweep.StaleTemp("rhiza-api-budget-")
	diskFreeBytes = func(dir string) (uint64, error) {
		panic("test read the host's real free disk space for " + dir + ": pin it with pinFreeBytes(ctx) or set RHIZA_TEST_FREE_BYTES beside RHIZA_TEST_NOW")
	}
	code := rootguard.Main(m.Run)
	if genDir != "" {
		os.RemoveAll(genDir)
	}
	os.Exit(code)
}

// defaultModel returns the shared generated model, writing and indexing it
// on first use.
func defaultModel(t testing.TB) generated {
	t.Helper()
	genOnce.Do(func() {
		dir, err := os.MkdirTemp("", "rhiza-api-budget-")
		if err != nil {
			genned.err = err
			return
		}
		genDir = dir
		root := filepath.Join(dir, "repo")
		if err := os.Mkdir(root, 0o755); err != nil {
			genned.err = err
			return
		}
		s, err := gen.Write(root, gen.Default())
		if err != nil {
			genned.err = err
			return
		}
		indexPath := filepath.Join(dir, "index.db")
		x, err := index.Open(root, indexPath)
		if err != nil {
			genned.err = err
			return
		}
		as, err := x.Ancestors(index.AsIs, s.Deep.ID)
		if err == nil && len(as) < 2 {
			err = fmt.Errorf("deep node %s has %d ancestors, want 2 or more", s.Deep.ID, len(as))
		}
		if cerr := x.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			genned.err = err
			return
		}
		genned = generated{root: root, indexPath: indexPath, summary: s, subtree: as[len(as)-2].ID}
	})
	if genned.err != nil {
		t.Fatalf("generate default model: %v", genned.err)
	}
	return genned
}

// budgetLock holds an exclusive lock on a file shared by every test that
// builds the 50,000-node model (the budget tests of internal/api,
// internal/index and internal/slice), until t ends. Each takes it before
// it writes and indexes its model: go test runs packages at once, and a
// budget timed while another package builds or reads its own model would
// measure that contention, not the read.
func budgetLock(t *testing.T) {
	t.Helper()
	f, err := os.OpenFile(filepath.Join(os.TempDir(), "rhiza-budget.lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		f.Close()
	})
}

// budgetRead is one read the budget times: an operation and its request.
type budgetRead struct {
	name string
	op   string
	req  any
}

// answerSize is how much a read's response holds: the nodes of show, ls
// and tree, the change sets of why, log and ready, the groups of roadmap.
// The budget wants it positive, so that no read is timed answering
// nothing.
func answerSize(resp any) int {
	switch r := resp.(type) {
	case ShowResponse:
		if r.Node != nil {
			return 1
		}
	case LsResponse:
		return len(r.Nodes)
	case TreeResponse:
		return len(r.Nodes)
	case WhyResponse:
		return len(r.ChangeSets)
	case LogResponse:
		return len(r.ChangeSets)
	case RoadmapResponse:
		return len(r.Groups)
	case ReadyResponse:
		return len(r.ChangeSets)
	}
	return 0
}

// budgetReads is every budgeted read of g in view v (D14): show of a deep
// node by handle, the first page of ls, a deep subtree, why and log of a
// deep node, roadmap and ready.
func budgetReads(g generated, v string) []budgetRead {
	s := g.summary
	return []budgetRead{
		{"show", "show", ShowRequest{Ref: s.Deep.Handle, View: v}},
		{"ls", "ls", LsRequest{View: v}},
		{"tree", "tree", TreeRequest{Ref: g.subtree.String(), View: v}},
		{"why", "why", WhyRequest{Ref: s.Deep.ID.String(), View: v}},
		{"log", "log", LogRequest{Ref: s.Deep.ID.String(), View: v}},
		{"roadmap", "roadmap", RoadmapRequest{View: v}},
		{"ready", "ready", ReadyRequest{View: v}},
	}
}

// runRead runs rd once as the CLI does: a fresh Ctx, so a fresh workspace
// and index open.
func runRead(g generated, rd budgetRead) (any, error) {
	ctx := testCtx(g.root)
	ctx.IndexPath = g.indexPath
	o, _ := Lookup(rd.op)
	return o.Run(ctx, rd.req)
}

// openTimingLog returns a function that appends a line to the file named by
// RHIZA_BUDGET_LOG, so that CI can print the timings after the run, passing
// or failing. The file is opened in append mode, starts its block with a
// header naming GOMAXPROCS, the CPU count and quiet (what quiettest.Wait
// reported), and is closed when t ends. An open or write error fails t.
// With the variable unset it returns a no-op.
func openTimingLog(t *testing.T, quiet string) func(line string) {
	t.Helper()
	path := os.Getenv("RHIZA_BUDGET_LOG")
	if path == "" {
		return func(string) {}
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("RHIZA_BUDGET_LOG: %v", err)
	}
	t.Cleanup(func() {
		if err := f.Close(); err != nil {
			t.Errorf("RHIZA_BUDGET_LOG %s: %v", path, err)
		}
	})
	write := func(line string) {
		t.Helper()
		if _, err := fmt.Fprintln(f, line); err != nil {
			t.Fatalf("RHIZA_BUDGET_LOG %s: %v", path, err)
		}
	}
	write(fmt.Sprintf("read budget: GOMAXPROCS=%d NumCPU=%d, %s", runtime.GOMAXPROCS(0), runtime.NumCPU(), quiet))
	return write
}

func TestReadBudget(t *testing.T) {
	if testing.Short() {
		t.Skip("the read budget generates the 50,000-node model; skipped under -short")
	}
	q := quiettest.Wait(t)
	t.Logf("other packages: %s", q)
	budgetLock(t)
	timingLog := openTimingLog(t, q.String())
	g := defaultModel(t)
	var late string
	for _, v := range []string{viewAsIs, viewToBe, viewWhatIfPref + g.summary.Draft.String()} {
		for _, rd := range budgetReads(g, v) {
			label := fmt.Sprintf("%-7s %-5s", rd.name, strings.TrimSuffix(v, g.summary.Draft.String()))
			log := func(line string) {
				t.Log(line)
				timingLog(line)
			}
			ok, err := timeWithinBudget(label, func() (readTiming, error) {
				var times []time.Duration
				var size int
				for range readBudgetRuns {
					start := time.Now()
					resp, err := runRead(g, rd)
					times = append(times, time.Since(start))
					if err != nil {
						return readTiming{}, fmt.Errorf("%s %s: %w", rd.name, v, err)
					}
					if size = answerSize(resp); size == 0 {
						return readTiming{}, fmt.Errorf("%s %s answered nothing: %+v", rd.name, v, resp)
					}
				}
				slices.Sort(times)
				median := (times[readBudgetRuns/2-1] + times[readBudgetRuns/2]) / 2
				return readTiming{median, times[0], times[len(times)-1], size}, nil
			}, log)
			if err != nil {
				t.Fatal(err)
			}
			if !ok {
				t.Errorf("%s %s median over %v on every attempt", rd.name, v, readBudget)
			}
			if late == "" {
				late = quiettest.Late(rd.name, v)
			}
		}
	}
	// The default view with a landed current change set: no view, so the
	// change set's state comes from the index, not a fold of the log.
	landed := landedChangeSet(t, g)
	for _, rd := range budgetReads(g, "") {
		label := fmt.Sprintf("%-7s %-5s", rd.name, "default")
		log := func(line string) {
			t.Log(line)
			timingLog(line)
		}
		ok, err := timeWithinBudget(label, func() (readTiming, error) {
			var times []time.Duration
			var size int
			for range readBudgetRuns {
				ctx := testCtx(g.root)
				ctx.IndexPath = g.indexPath
				ctx.ChangeSet = landed
				o, _ := Lookup(rd.op)
				start := time.Now()
				resp, err := o.Run(ctx, rd.req)
				times = append(times, time.Since(start))
				if err != nil {
					return readTiming{}, fmt.Errorf("%s default: %w", rd.name, err)
				}
				if size = answerSize(resp); size == 0 {
					return readTiming{}, fmt.Errorf("%s default answered nothing: %+v", rd.name, resp)
				}
			}
			slices.Sort(times)
			median := (times[readBudgetRuns/2-1] + times[readBudgetRuns/2]) / 2
			return readTiming{median, times[0], times[len(times)-1], size}, nil
		}, log)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("%s default median over %v on every attempt", rd.name, readBudget)
		}
		if late == "" {
			late = quiettest.Late(rd.name, "default")
		}
	}
	if late != "" {
		t.Log(late)
		timingLog(late)
	}
	// Last: it appends an entry to the shared model.
	changeShowBudget(t, g, func(line string) {
		t.Log(line)
		timingLog(line)
	})
	driverPollBudget(t, g, func(line string) {
		t.Log(line)
		timingLog(line)
	})
}

// landedChangeSet returns a landed change set of g's index.
func landedChangeSet(t *testing.T, g generated) id.ID {
	t.Helper()
	x, err := index.Open(g.root, g.indexPath)
	if err != nil {
		t.Fatal(err)
	}
	defer x.Close()
	infos, err := x.ChangeSets()
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range infos {
		if c.State == changeset.StateLanded {
			return c.ID
		}
	}
	t.Fatal("the generated model has no landed change set")
	return id.ID{}
}

// benchRead runs the read named name in to-be b.N times.
func benchRead(b *testing.B, name string) {
	g := defaultModel(b)
	i := slices.IndexFunc(budgetReads(g, viewToBe), func(rd budgetRead) bool { return rd.name == name })
	rd := budgetReads(g, viewToBe)[i]
	for b.Loop() {
		if _, err := runRead(g, rd); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadShow(b *testing.B)    { benchRead(b, "show") }
func BenchmarkReadLs(b *testing.B)      { benchRead(b, "ls") }
func BenchmarkReadTree(b *testing.B)    { benchRead(b, "tree") }
func BenchmarkReadWhy(b *testing.B)     { benchRead(b, "why") }
func BenchmarkReadLog(b *testing.B)     { benchRead(b, "log") }
func BenchmarkReadRoadmap(b *testing.B) { benchRead(b, "roadmap") }
func BenchmarkReadReady(b *testing.B)   { benchRead(b, "ready") }

// BenchmarkTreeWhole shows the whole capability tree: not budgeted (D14).
func BenchmarkTreeWhole(b *testing.B) {
	g := defaultModel(b)
	rd := budgetRead{"tree", "tree", TreeRequest{View: viewToBe}}
	for b.Loop() {
		if _, err := runRead(g, rd); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkShowWhatIf shows a deep node in the what-if of the draft change
// set.
func BenchmarkShowWhatIf(b *testing.B) {
	g := defaultModel(b)
	s := g.summary
	rd := budgetRead{"show", "show", ShowRequest{Ref: s.Deep.Handle, View: viewWhatIfPref + s.Draft.String()}}
	for b.Loop() {
		if _, err := runRead(g, rd); err != nil {
			b.Fatal(err)
		}
	}
}

func TestTimeWithinBudget(t *testing.T) {
	over := readTiming{median: readBudget, size: 1}
	under := readTiming{median: readBudget - 1, size: 1}
	run := func(seq ...readTiming) (bool, int, []string) {
		var lines []string
		calls := 0
		ok, err := timeWithinBudget("x", func() (readTiming, error) {
			rt := seq[min(calls, len(seq)-1)]
			calls++
			return rt, nil
		}, func(l string) { lines = append(lines, l) })
		if err != nil {
			t.Fatal(err)
		}
		return ok, calls, lines
	}
	if ok, calls, lines := run(over, under); !ok || calls != 2 || len(lines) != 2 {
		t.Errorf("over then under: ok %v, %d calls, log %q", ok, calls, lines)
	}
	if ok, calls, lines := run(over, over, over); ok || calls != 3 || len(lines) != 3 {
		t.Errorf("over three times: ok %v, %d calls, log %q", ok, calls, lines)
	}
	if ok, calls, _ := run(under); !ok || calls != 1 {
		t.Errorf("under at once: ok %v, %d calls", ok, calls)
	}
	if readBudget != 50*time.Millisecond {
		t.Errorf("readBudget = %v, want 50ms", readBudget)
	}
}
