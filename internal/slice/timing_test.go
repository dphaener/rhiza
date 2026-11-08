package slice

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/gen"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/quiettest"
	"github.com/dphaener/rhiza/internal/tmpsweep"
)

// sliceBudget is the most a slice's median may take, open to printed.
const sliceBudget = 50 * time.Millisecond

// sliceRuns is how many times TestBudget times each slice.
const sliceRuns = 20

// sliceTokens is the token budget each timed slice is fitted to.
const sliceTokens = 4000

// generated is gen.Default() written once, with its index built, shared by
// TestBudget. It lives under one os.MkdirTemp folder that TestMain removes.
type generated struct {
	root    string
	dbPath  string
	summary gen.Summary
	err     error
}

var (
	genOnce sync.Once
	genDir  string
	genned  generated
)

func TestMain(m *testing.M) {
	// A run killed by a go test timeout never reached the RemoveAll below.
	tmpsweep.StaleTemp("rhiza-slice-budget-")
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
		dir, err := os.MkdirTemp("", "rhiza-slice-budget-")
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
		x, err := index.Open(root, dbPath)
		if err != nil {
			genned.err = err
			return
		}
		if err := x.Close(); err != nil {
			genned.err = err
			return
		}
		genned = generated{root: root, dbPath: dbPath, summary: s}
	})
	if genned.err != nil {
		t.Fatalf("generate default model: %v", genned.err)
	}
	return genned
}

// busyFocus finds, breadth first from the root above the summary's deep
// node, the first to-be tree node at depth 3 or more with at least two live
// children and at least three live links. Children come in key order, so
// the search is deterministic.
func busyFocus(x *index.Index, s gen.Summary) (model.Node, int, error) {
	anc, err := x.Ancestors(index.ToBe, s.Deep.ID)
	if err != nil {
		return model.Node{}, 0, err
	}
	if len(anc) == 0 {
		return model.Node{}, 0, fmt.Errorf("deep node %s has no ancestors", s.Deep.Handle)
	}
	type at struct {
		n     model.Node
		depth int
	}
	queue := []at{{anc[0], 1}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		kids, err := x.Children(index.ToBe, cur.n.ID)
		if err != nil {
			return model.Node{}, 0, err
		}
		live := 0
		for _, k := range kids {
			if !k.Retired {
				live++
				queue = append(queue, at{k, cur.depth + 1})
			}
		}
		if cur.depth < 3 || live < 2 || cur.n.Retired {
			continue
		}
		ls, err := x.Links(index.ToBe, cur.n.ID, index.Both)
		if err != nil {
			return model.Node{}, 0, err
		}
		links := 0
		for _, l := range ls {
			if !l.Retired {
				links++
			}
		}
		if links >= 3 {
			return cur.n, cur.depth, nil
		}
	}
	return model.Node{}, 0, fmt.Errorf("no tree node at depth 3+ with 2+ children and 3+ links")
}

// timedSlice is one slice TestBudget times.
type timedSlice struct {
	name   string
	view   index.View
	focus  id.ID
	change id.ID // the change set sliced, zero for a node slice
}

// sliceOnce is the budget's unit of work, as the api will do it: open the
// index fresh, build the request from its change sets, build the slice and
// fit it in both formats.
func sliceOnce(g generated, ts timedSlice) (json, md Result, err error) {
	x, err := index.Open(g.root, g.dbPath)
	if err != nil {
		return
	}
	defer func() {
		if cerr := x.Close(); err == nil {
			err = cerr
		}
	}()
	infos, err := x.ChangeSets()
	if err != nil {
		return
	}
	viewName := ViewAsIs
	if ts.view == index.ToBe {
		viewName = ViewToBe
	}
	req := Request{Focus: ts.focus, View: viewName}
	for _, c := range infos {
		req.ChangeSets = append(req.ChangeSets, c.ID)
		switch c.State {
		case changeset.StateDraft, changeset.StateProposed, changeset.StateAccepted:
			es, eerr := x.Entries(c.ID)
			if eerr != nil {
				return json, md, eerr
			}
			cs, _ := changeset.Fold(c.ID, es)
			req.InFlight = append(req.InFlight, InFlightOf(cs))
		}
	}
	slices.SortFunc(req.ChangeSets, id.Compare)
	if !ts.change.IsZero() {
		es, eerr := x.Entries(ts.change)
		if eerr != nil {
			return json, md, eerr
		}
		cs, _ := changeset.Fold(ts.change, es)
		req.Change = &cs
	}
	s, err := Build(FromIndex(x, ts.view), FromIndex(x, index.AsIs), req)
	if err != nil {
		return
	}
	opts := Options{Budget: sliceTokens, JSON: true, DefaultView: DefaultView(req.Change)}
	if json, err = Fit(s, opts); err != nil {
		return
	}
	opts.JSON = false
	md, err = Fit(s, opts)
	return
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

func TestBudget(t *testing.T) {
	if testing.Short() {
		t.Skip("budget generates the 50,000-node model; skipped under -short")
	}
	budgetLock(t)
	g := defaultModel(t)
	s := g.summary
	x, err := index.Open(g.root, g.dbPath)
	if err != nil {
		t.Fatal(err)
	}
	focus, depth, err := busyFocus(x, s)
	x.Close()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("model: %d entries, %d change sets, %d nodes, %d links", s.Entries, s.ChangeSets, s.Nodes, s.Links)
	t.Logf("focus: %s %s at depth %d", focus.Kind, focus.Handle, depth)

	for _, ts := range []timedSlice{
		{name: "node slice, as-is", view: index.AsIs, focus: focus.ID},
		{name: "node slice, to-be", view: index.ToBe, focus: focus.ID},
		{name: "change slice of the accepted set, to-be", view: index.ToBe, change: s.Accepted},
	} {
		var laps []quiettest.Lap
		var json, md Result
		for range sliceRuns {
			sw := quiettest.Start()
			json, md, err = sliceOnce(g, ts)
			laps = append(laps, sw.Stop())
			if err != nil {
				t.Fatalf("%s: %v", ts.name, err)
			}
		}
		if len(json.Bytes) == 0 || len(md.Bytes) == 0 {
			t.Fatalf("%s: empty output (json %d bytes, markdown %d)", ts.name, len(json.Bytes), len(md.Bytes))
		}
		tm := quiettest.Summarize(laps)
		t.Logf("%s over %d runs: %v (json %d tokens, markdown %d)",
			ts.name, sliceRuns, tm, json.Tokens, md.Tokens)
		// Own time: the clock less waiting for a CPU, inconclusive on a
		// saturated machine (quiettest.Timing.Over).
		if over, why := tm.Over(sliceBudget); over {
			t.Errorf("%s: median own time %v, want < %v", ts.name, tm.Own, sliceBudget)
		} else if why != "" {
			t.Log(ts.name + " " + why)
		}
	}
}
