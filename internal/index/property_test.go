package index

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/index/indextest"
	"github.com/dphaener/rhiza/internal/store"
	"pgregory.net/rapid"
)

// logFile is one file of a drawn log: its path under .rhiza, its bytes
// and, for an entry, its type.
type logFile struct {
	path string
	data []byte
	typ  entry.Type
}

// genFiles draws a log with indextest.Draw (see there) and returns its
// files.
func genFiles(t *rapid.T) []logFile {
	return logFiles(indextest.Draw(t).Files)
}

func logFiles(fs []indextest.File) []logFile {
	out := make([]logFile, len(fs))
	for i, f := range fs {
		out[i] = logFile{path: f.Path, data: f.Data, typ: f.Type}
	}
	return out
}

// TestPropertyIncremental: catching up through any split of a log into
// steps, including steps that remove files already indexed (rebuilds)
// and add them back, gives the same index as one scratch build, and both
// equal what store, changeset and model compute from the files.
func TestPropertyIncremental(t *testing.T) {
	var lands, accepts, rebuilds, reopens, steps int
	rapid.Check(t, func(t *rapid.T) {
		files := indextest.Draw(t).Files
		root, err := os.MkdirTemp("", "index-prop-")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(root)
		gdir := filepath.Join(root, store.Dir)
		if err := os.MkdirAll(filepath.Join(gdir, store.LogDir), 0o755); err != nil {
			t.Fatal(err)
		}
		db := filepath.Join(root, "index.db")

		x, err := Open(root, db)
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		defer func() { x.Close() }()

		st := indextest.NewSteps(gdir, files)
		for !st.Done() {
			steps++
			step := st.Next(t)
			if step.Landed {
				lands++
			}
			if step.Accepted {
				accepts++
			}
			if rapid.Bool().Draw(t, "reopen?") {
				reopens++
				x.Close()
				if x, err = Open(root, db); err != nil {
					t.Fatalf("Open: %v", err)
				}
			} else if err := x.catchUp(); err != nil {
				t.Fatalf("catchUp: %v", err)
			}
			if x.stats.rebuilt {
				rebuilds++
			}
		}

		scratch, err := Open(root, filepath.Join(root, "scratch.db"))
		if err != nil {
			t.Fatalf("scratch Open: %v", err)
		}
		defer scratch.Close()
		same(t, "incremental vs scratch", dump(t, x), dump(t, scratch))
		checkModel(t, x, gdir)
	})
	t.Logf("steps %d, with a land %d, with an accept %d, rebuilds %d, reopens %d", steps, lands, accepts, rebuilds, reopens)
	if lands == 0 || accepts == 0 || rebuilds == 0 {
		t.Errorf("vacuous run: lands %d, accepts %d, rebuilds %d", lands, accepts, rebuilds)
	}
}
