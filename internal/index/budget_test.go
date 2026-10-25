package index

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/dphaener/genesis/internal/gen"
	"github.com/dphaener/genesis/internal/model"
)

// budget is the most the fresh-read path's median may take (D13).
const budget = 25 * time.Millisecond

// budgetRuns is how many fresh reads TestBudget times.
const budgetRuns = 20

// generated is gen.Default() written once, with its index built, shared by
// TestBudget and the benchmarks. It lives under one os.MkdirTemp folder
// that TestMain removes.
type generated struct {
	root    string
	summary gen.Summary
	// dbPath is the index built from scratch; build is how long that took.
	dbPath string
	build  time.Duration
	err    error
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
		dir, err := os.MkdirTemp("", "genesis-index-budget-")
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
		dbPath := filepath.Join(dir, "index.db")
		start := time.Now()
		x, err := Open(root, dbPath)
		build := time.Since(start)
		if err != nil {
			genned.err = err
			return
		}
		if err := x.Close(); err != nil {
			genned.err = err
			return
		}
		genned = generated{root: root, summary: s, dbPath: dbPath, build: build}
	})
	if genned.err != nil {
		t.Fatalf("generate default model: %v", genned.err)
	}
	return genned
}

// freshRead is the budget's unit of work: open the index fresh, answer the
// queries a command typically asks, close it.
func freshRead(g generated) (resolved model.Node, ancestors, children []model.Node, links []model.Link, err error) {
	s := g.summary
	x, err := Open(g.root, g.dbPath)
	if err != nil {
		return
	}
	defer func() {
		if cerr := x.Close(); err == nil {
			err = cerr
		}
	}()
	if resolved, err = x.Resolve(ToBe, s.Deep.Handle); err != nil {
		return
	}
	if ancestors, err = x.Ancestors(ToBe, s.Deep.ID); err != nil {
		return
	}
	if children, err = x.Children(ToBe, s.Parent.ID); err != nil {
		return
	}
	links, err = x.Links(ToBe, s.Linked.ID, Both)
	return
}

func TestBudget(t *testing.T) {
	if testing.Short() {
		t.Skip("budget generates the 50,000-node model; skipped under -short")
	}
	g := defaultModel(t)
	s := g.summary
	var times []time.Duration
	for range budgetRuns {
		start := time.Now()
		resolved, ancestors, children, links, err := freshRead(g)
		times = append(times, time.Since(start))
		if err != nil {
			t.Fatal(err)
		}
		if resolved.ID != s.Deep.ID {
			t.Fatalf("Resolve(%q) = %s, want %s", s.Deep.Handle, resolved.ID, s.Deep.ID)
		}
		if len(ancestors) != s.Depth-1 {
			t.Fatalf("Ancestors(deep) = %d nodes, want %d", len(ancestors), s.Depth-1)
		}
		if len(children) == 0 {
			t.Fatalf("Children(parent) is empty, want about %d", s.LiveChildren)
		}
		if live := liveLinks(links); live < s.LinksOut+s.LinksIn {
			t.Fatalf("Links(linked) has %d live links, want at least %d", live, s.LinksOut+s.LinksIn)
		}
	}
	slices.Sort(times)
	median := (times[budgetRuns/2-1] + times[budgetRuns/2]) / 2
	t.Logf("model: %d entries, %d change sets, %d nodes, %d links", s.Entries, s.ChangeSets, s.Nodes, s.Links)
	t.Logf("scratch build: %v", g.build)
	t.Logf("fresh read over %d runs: min %v, median %v, max %v", budgetRuns, times[0], median, times[len(times)-1])
	if median >= budget {
		t.Fatalf("fresh read median %v, want < %v", median, budget)
	}
}

// liveLinks counts the links in ls that are not retired.
func liveLinks(ls []model.Link) int {
	n := 0
	for _, l := range ls {
		if !l.Retired {
			n++
		}
	}
	return n
}

// copyFile copies src over dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return fmt.Errorf("copy %s: %w", src, err)
	}
	return nil
}
