package index

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/quiettest"
	"github.com/dphaener/rhiza/internal/store"
	"pgregory.net/rapid"
)

// checkFolded asserts the stored folded change sets, history included, and
// the fold problems equal changeset.FoldLog over the files in gdir.
func checkFolded(t tb, x *Index, gdir string) {
	t.Helper()
	sets, ps := changeset.FoldLog(store.ReadLog(os.DirFS(gdir)))
	ids := make([]id.ID, len(sets))
	for i, s := range sets {
		ids[i] = s.ID
	}
	got, err := x.FoldedOf(ids)
	if err != nil {
		t.Fatalf("FoldedOf: %v", err)
	}
	if len(got) != len(sets) {
		t.Fatalf("FoldedOf gave %d change sets, want %d", len(got), len(sets))
	}
	for _, s := range sets {
		if !reflect.DeepEqual(got[s.ID], s) {
			t.Fatalf("folded %s differs from FoldLog:\n got %+v\nwant %+v", s.ID, got[s.ID], s)
		}
	}
	var want []storedProblem
	for _, p := range ps {
		want = append(want, storedProblem{Source: sourceFold, ChangeSet: idText(p.ChangeSet), Entry: idText(p.Entry),
			Sentinel: foldSentinel(p.Err), Message: p.Err.Error()})
	}
	same(t, "fold problems", storedProblems(t, x, sourceFold), want)
}

// TestPropertyFoldedEveryStep appends a drawn log's files one at a time and
// checks, after every catch-up, that the stored folded change sets and
// fold problems equal FoldLog over the whole log.
func TestPropertyFoldedEveryStep(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		files := genFiles(t)
		order := rapid.Permutation(files).Draw(t, "order")
		root, err := os.MkdirTemp("", "index-fold-")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(root)
		gdir := filepath.Join(root, store.Dir)
		if err := os.MkdirAll(filepath.Join(gdir, store.LogDir), 0o755); err != nil {
			t.Fatal(err)
		}
		x, err := Open(root, filepath.Join(root, "index.db"))
		if err != nil {
			t.Fatal(err)
		}
		defer x.Close()
		for _, f := range order {
			writeFile(t, gdir, f.path, f.data)
			if err := x.catchUp(); err != nil {
				t.Fatalf("catchUp: %v", err)
			}
			checkFolded(t, x, gdir)
		}
	})
}

// TestCatchUpFoldsOnlyTouched appends one entry to one change set of a
// large log and counts what catch-up decodes and folds.
func TestCatchUpFoldsOnlyTouched(t *testing.T) {
	root := t.TempDir()
	gdir := filepath.Join(root, store.Dir)
	const sets, per = 120, 3
	n := 0
	for k := 1; k <= sets; k++ {
		var head []id.ID
		for j := 0; j < per; j++ {
			n++
			var b changeset.Body = changeset.Note{Text: fmt.Sprintf("n%d", n)}
			if j == 0 {
				b = changeset.Open{Kind: changeset.KindIntention, Intent: fmt.Sprintf("set %d", k)}
			}
			p, data := encodeEntry(t, eID(n), csID(k), head, b)
			writeFile(t, gdir, p, data)
			head = []id.ID{eID(n)}
		}
	}
	x, err := Open(root, filepath.Join(root, "index.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer x.Close()
	if x.stats.folded != sets || x.stats.decoded != sets*per {
		t.Fatalf("scratch build folded %d decoded %d, want %d and %d", x.stats.folded, x.stats.decoded, sets, sets*per)
	}

	n++
	p, data := encodeEntry(t, eID(n), csID(57), []id.ID{eID(57 * per)}, changeset.Note{Text: "late"})
	writeFile(t, gdir, p, data)
	if err := x.catchUp(); err != nil {
		t.Fatal(err)
	}
	if x.stats.read != 1 || x.stats.folded != 1 || x.stats.decoded != per+1 {
		t.Fatalf("catch-up read %d folded %d decoded %d, want 1, 1 and %d", x.stats.read, x.stats.folded, x.stats.decoded, per+1)
	}
	if x.stats.asIs || x.stats.toBe {
		t.Fatalf("a note on a draft recomputed a view: asIs %v toBe %v", x.stats.asIs, x.stats.toBe)
	}
	checkFolded(t, x, gdir)
}

// copyTree copies the folder src to dst.
func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, p)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		return copyFile(p, filepath.Join(dst, rel))
	})
}

// TestCatchUpBudget times catch-up after one appended entry on the
// 50,000-node model.
func TestCatchUpBudget(t *testing.T) {
	if testing.Short() {
		t.Skip("budget generates the 50,000-node model; skipped under -short")
	}
	budgetLock(t)
	g := defaultModel(t)
	root := filepath.Join(t.TempDir(), "repo")
	if err := copyTree(g.root, root); err != nil {
		t.Fatal(err)
	}
	dbPath := filepath.Join(t.TempDir(), "index.db")
	if err := copyFile(g.dbPath, dbPath); err != nil {
		t.Fatal(err)
	}
	gdir := filepath.Join(root, store.Dir)
	x, err := Open(root, dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer x.Close()

	// A draft change set of its own, indexed before the timed runs.
	cs := csID(900000)
	n := 900000
	var head []id.ID
	add := func(b changeset.Body) {
		n++
		p, data := encodeEntry(t, eID(n), cs, head, b)
		writeFile(t, gdir, p, data)
		head = []id.ID{eID(n)}
	}
	add(changeset.Open{Kind: changeset.KindIntention, Intent: "budget"})
	if err := x.catchUp(); err != nil {
		t.Fatal(err)
	}
	var laps []quiettest.Lap
	for i := range 21 {
		add(changeset.Note{Text: fmt.Sprintf("note %d", i)})
		sw := quiettest.Start()
		err := x.catchUp()
		laps = append(laps, sw.Stop())
		if err != nil {
			t.Fatal(err)
		}
		if x.stats.read != 1 || x.stats.folded != 1 || x.stats.asIs || x.stats.toBe {
			t.Fatalf("catch-up read %d folded %d asIs %v toBe %v, want one entry, one fold, no view", x.stats.read, x.stats.folded, x.stats.asIs, x.stats.toBe)
		}
	}
	tm := quiettest.Summarize(laps)
	t.Logf("catch-up after one entry: %v", tm)
	// Own time: the clock less waiting for a CPU, inconclusive on a
	// saturated machine (quiettest.Timing.Over).
	if over, why := tm.Over(50 * time.Millisecond); over {
		t.Fatalf("catch-up median own time %v, want < 50ms", tm.Own)
	} else if why != "" {
		t.Log("catch-up " + why)
	}
}
