package index

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index/indextest"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/store"
	"pgregory.net/rapid"
)

// TestAsIsProperty: after every step of catching an index up through a
// drawn log (rebuilds included), sometimes with a landed duplicate handle
// added, AsIs gives the model and problems
// model.AsIs computes over changeset.FoldLog of the same log.
func TestAsIsProperty(t *testing.T) {
	var heads, dups, rebuilds, steps int
	rapid.Check(t, func(t *rapid.T) {
		files := indextest.Draw(t).Files
		if rapid.Bool().Draw(t, "duplicate handle?") {
			files = append(files, duplicateHandle(t)...)
		}
		root, err := os.MkdirTemp("", "index-asis-")
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
			t.Fatalf("Open: %v", err)
		}
		defer x.Close()

		var head, dup, rebuilt bool
		st := indextest.NewSteps(gdir, files)
		for !st.Done() {
			steps++
			st.Next(t)
			if err := x.catchUp(); err != nil {
				t.Fatalf("catchUp: %v", err)
			}
			rebuilt = rebuilt || x.stats.rebuilt
			sets, _ := changeset.FoldLog(store.ReadLog(os.DirFS(gdir)))
			head = head || twoLandedHeadsOf(sets)
			want, wps := model.AsIs(sets)
			for _, p := range wps {
				dup = dup || model.IsDuplicateHandle(p.Err)
			}
			got, gps, err := x.AsIs()
			if err != nil {
				t.Fatalf("AsIs: %v", err)
			}
			sameModel(t, "as-is", got, want)
			same(t, "as-is problems", viewProblems(viewAsIs, gps), viewProblems(viewAsIs, wps))
		}
		if head {
			heads++
		}
		if dup {
			dups++
		}
		if rebuilt {
			rebuilds++
		}
	})
	t.Logf("steps %d, logs with two landed heads %d, with a duplicate handle %d, with a rebuild %d", steps, heads, dups, rebuilds)
	if heads == 0 || dups == 0 || rebuilds == 0 {
		t.Errorf("vacuous run: two landed heads %d, duplicate handles %d, rebuilds %d", heads, dups, rebuilds)
	}
}

// duplicateHandle returns the files of a landed record change set that
// creates two root capabilities with handle "a", the handle drawn logs
// give their nodes most often: the drawn log alone holds a duplicate
// handle too seldom for every run to see one.
func duplicateHandle(t *rapid.T) []indextest.File {
	cs := indextest.ChangeSet(50)
	var ops []op.Op
	for i := range 2 {
		ops = append(ops, op.Op{ID: indextest.Op(7000 + i), Type: op.TypeNodeCreate, Node: indextest.Node(90 + i),
			Kind: op.KindCapability, Handle: "a", Key: "i"})
	}
	var out []indextest.File
	var after id.ID
	for i, b := range []changeset.Body{changeset.Open{Kind: changeset.KindRecord, Intent: "dup"},
		changeset.Ops{Ops: ops}, changeset.Land{}} {
		e := indextest.NewEntry(t, indextest.Entry(7000+i), cs, after, b)
		data, err := entry.Encode(e)
		if err != nil {
			t.Fatalf("encode entry: %v", err)
		}
		out = append(out, indextest.File{Path: indextest.EntryPath(cs, e.ID), Data: data, Type: e.Type})
		after = e.ID
	}
	return out
}

// sameModel fails unless got and want hold the same nodes, links and tree.
func sameModel(t tb, what string, got, want model.Model) {
	t.Helper()
	same(t, what+" nodes", got.Nodes(), want.Nodes())
	same(t, what+" links", got.Links(), want.Links())
	same(t, what+" roots", got.Children(id.ID{}), want.Children(id.ID{}))
	for _, n := range want.Nodes() {
		same(t, what+" children of "+n.ID.String(), got.Children(n.ID), want.Children(n.ID))
	}
}

// twoLandedHeadsOf reports whether some target has two ops among the
// landed change sets' that no other landed op on it names as base.
func twoLandedHeadsOf(sets []changeset.ChangeSet) bool {
	byTarget := map[id.ID][]op.Op{}
	for _, cs := range sets {
		if cs.State == changeset.StateLanded {
			for _, o := range cs.Ops {
				byTarget[o.Target()] = append(byTarget[o.Target()], o)
			}
		}
	}
	for _, ops := range byTarget {
		based := map[id.ID]bool{}
		for _, o := range ops {
			based[o.Base] = true
		}
		n := 0
		for _, o := range ops {
			if !based[o.ID] {
				n++
			}
		}
		if n > 1 {
			return true
		}
	}
	return false
}

// TestAsIsEmpty: an index of an empty log gives an empty model, no
// problems and no error.
func TestAsIsEmpty(t *testing.T) {
	f := newFixture(t)
	x, err := Open(f.root, f.dbPath())
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer x.Close()
	m, ps, err := x.AsIs()
	if err != nil {
		t.Fatalf("AsIs: %v", err)
	}
	if m.Nodes() != nil || m.Links() != nil || m.Children(id.ID{}) != nil || ps != nil {
		t.Fatalf("AsIs of an empty log: nodes %v, links %v, problems %v", m.Nodes(), m.Links(), ps)
	}
}

// TestAsIsClosed: a closed index returns an error.
func TestAsIsClosed(t *testing.T) {
	f := newFixture(t)
	f.add(csID(1), changeset.Open{Kind: changeset.KindIntention, Intent: "x"})
	x, err := Open(f.root, f.dbPath())
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if err := x.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if _, _, err := x.AsIs(); err == nil {
		t.Fatal("AsIs of a closed index: no error")
	}
}
