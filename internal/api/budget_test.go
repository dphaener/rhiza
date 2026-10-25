package api

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/dphaener/genesis/internal/gen"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/index"
)

// readBudget is the most a read's median may take, fresh workspace and
// index open to answer: milestone 1's target for every read on the
// 50,000-node model (D14, as revised).
const readBudget = 50 * time.Millisecond

// readBudgetRuns is how many fresh runs TestReadBudget times per read.
const readBudgetRuns = 20

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
	code := m.Run()
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
		dir, err := os.MkdirTemp("", "genesis-api-budget-")
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
	f, err := os.OpenFile(filepath.Join(os.TempDir(), "genesis-budget.lock"), os.O_CREATE|os.O_RDWR, 0o600)
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

func TestReadBudget(t *testing.T) {
	if testing.Short() {
		t.Skip("the read budget generates the 50,000-node model; skipped under -short")
	}
	budgetLock(t)
	g := defaultModel(t)
	for _, v := range []string{viewAsIs, viewToBe} {
		for _, rd := range budgetReads(g, v) {
			var times []time.Duration
			var size int
			for range readBudgetRuns {
				start := time.Now()
				resp, err := runRead(g, rd)
				times = append(times, time.Since(start))
				if err != nil {
					t.Fatalf("%s %s: %v", rd.name, v, err)
				}
				if size = answerSize(resp); size == 0 {
					t.Fatalf("%s %s answered nothing: %+v", rd.name, v, resp)
				}
			}
			slices.Sort(times)
			median := (times[readBudgetRuns/2-1] + times[readBudgetRuns/2]) / 2
			t.Logf("%-7s %-5s median %v (min %v, max %v), size %d", rd.name, v, median, times[0], times[len(times)-1], size)
			if median >= readBudget {
				t.Errorf("%s %s median %v, want < %v", rd.name, v, median, readBudget)
			}
		}
	}
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
// set: not budgeted (D14).
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
