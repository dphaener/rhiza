package index

import (
	"errors"
	"slices"
	"sync"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

// light is n without Fields and History, as the light queries return it.
func light(n model.Node) model.Node {
	n.Fields, n.History = nil, nil
	return n
}

// batchIDSets are the ID sets the batch tests ask about: empty, repeats,
// retired and unknown nodes, a change set, and one over batchSize, so it
// takes several chunks, with known IDs in the first and the last.
func batchIDSets() map[string][]id.ID {
	big := []id.ID{nd(1)}
	for n := range batchSize + 10 {
		big = append(big, nd(100+n))
	}
	big = append(big, nd(9), nd(1000))
	return map[string][]id.ID{
		"none":     nil,
		"one":      {nd(1)},
		"repeats":  {nd(3), nd(1), nd(3), nd(1)},
		"retired":  {nd(8), nd(13), nd(14)},
		"unknown":  {nd(77), {}},
		"to-be":    {nd(11), nd(12), nd(5)},
		"links":    {nd(1), nd(9), nd(2), nd(5)},
		"cs":       {csID(1), csID(8)},
		"big":      big,
		"everyone": {nd(1), nd(2), nd(3), nd(4), nd(5), nd(6), nd(7), nd(8), nd(9), nd(10), nd(1000), nd(11), nd(13), nd(14)},
	}
}

func TestNodesLight(t *testing.T) {
	x, ms := listIndex(t)
	for name, ids := range batchIDSets() {
		for _, v := range []View{AsIs, ToBe} {
			var want []model.Node
			for _, n := range ms[v].Nodes() {
				if slices.Contains(ids, n.ID) {
					want = append(want, light(n))
				}
			}
			got, err := x.NodesLight(v, ids)
			if err != nil {
				t.Fatalf("%s %s: %v", name, v, err)
			}
			same(t, "nodes light "+name+" "+string(v), got, want)
		}
	}
	// The light nodes drop what the full ones have.
	n, _, _ := x.Node(ToBe, nd(1))
	if n.Fields == nil || n.History == nil {
		t.Fatalf("fixture: root has no fields or history: %+v", n)
	}
	if _, err := x.NodesLight(View("what"), nil); err == nil {
		t.Fatal("unknown view: no error")
	}
}

func TestNodesFields(t *testing.T) {
	x, ms := listIndex(t)
	for name, ids := range batchIDSets() {
		for _, v := range []View{AsIs, ToBe} {
			var want []model.Node
			for _, n := range ms[v].Nodes() {
				if slices.Contains(ids, n.ID) {
					n.History = nil
					want = append(want, n)
				}
			}
			got, err := x.NodesFields(v, ids)
			if err != nil {
				t.Fatalf("%s %s: %v", name, v, err)
			}
			same(t, "nodes fields "+name+" "+string(v), got, want)
		}
	}
	if _, err := x.NodesFields(View("what"), nil); err == nil {
		t.Fatal("unknown view: no error")
	}
}

func TestSubtreeLight(t *testing.T) {
	x, ms := listIndex(t)
	for _, root := range []id.ID{{}, nd(1), nd(3), nd(13), nd(9), nd(8)} {
		for _, depth := range []int{0, 1, 2} {
			for _, retired := range []bool{false, true} {
				for _, v := range []View{AsIs, ToBe} {
					want := wantSubtree(ms[v], root, depth, retired)
					for k := range want {
						want[k].Node = light(want[k].Node)
					}
					got, err := x.SubtreeLight(v, root, depth, retired)
					if err != nil {
						t.Fatalf("SubtreeLight(%s, %s, %d, %v): %v", v, root, depth, retired, err)
					}
					same(t, "subtree light "+string(v)+" "+root.String(), got, want)
				}
			}
		}
	}
	if _, err := x.SubtreeLight(AsIs, nd(77), 0, false); !errors.Is(err, ErrNotFound) {
		t.Fatalf("unknown root: %v", err)
	}
}

func TestLinksAt(t *testing.T) {
	x, ms := listIndex(t)
	for name, ids := range batchIDSets() {
		for _, dir := range []Direction{From, To, Both} {
			for _, v := range []View{AsIs, ToBe} {
				var want []model.Link
				for _, l := range ms[v].Links() {
					from, to := slices.Contains(ids, l.From), slices.Contains(ids, l.To)
					if dir != To && from || dir != From && to {
						want = append(want, l)
					}
				}
				got, err := x.LinksAt(v, ids, dir)
				if err != nil {
					t.Fatalf("%s %d %s: %v", name, dir, v, err)
				}
				same(t, "links at "+name+" "+string(v), got, want)
			}
		}
	}
	if _, err := x.LinksAt(AsIs, []id.ID{nd(1)}, Direction(9)); err == nil {
		t.Fatal("unknown direction: no error")
	}
}

func TestLinksTo(t *testing.T) {
	x, ms := listIndex(t)
	for name, ids := range batchIDSets() {
		for _, typ := range []op.LinkType{op.LinkServes, op.LinkRelatesTo, op.LinkBlocks} {
			for _, v := range []View{AsIs, ToBe} {
				var want []model.Link
				for _, l := range ms[v].Links() {
					if l.Type == typ && slices.Contains(ids, l.To) {
						want = append(want, l)
					}
				}
				got, err := x.LinksTo(v, ids, typ)
				if err != nil {
					t.Fatalf("%s %s %s: %v", name, typ, v, err)
				}
				same(t, "links to "+name+" "+string(typ)+" "+string(v), got, want)
			}
		}
	}
	// The fixture has both: a live serves link to root and a retired one.
	got, _ := x.LinksTo(AsIs, []id.ID{nd(1)}, op.LinkServes)
	same(t, "serves root", linkIDs(got), []id.ID{ln(1), ln(3), ln(4)})
}

func TestEntriesOf(t *testing.T) {
	x, _ := listIndex(t)
	css := []id.ID{csID(1), csID(4), csID(8), csID(77), csID(1)}
	got, err := x.EntriesOf(css)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Fatalf("EntriesOf has %d change sets, want 3", len(got))
	}
	for _, cs := range []id.ID{csID(1), csID(4), csID(8)} {
		want, err := x.Entries(cs)
		if err != nil {
			t.Fatal(err)
		}
		same(t, "entries of "+cs.String(), got[cs], want)
	}
	if _, ok := got[csID(77)]; ok {
		t.Fatal("unknown change set has a key")
	}
	if got, err := x.EntriesOf(nil); len(got) != 0 || err != nil {
		t.Fatalf("no change sets: %v, %v", got, err)
	}
}

func TestParents(t *testing.T) {
	x, ms := listIndex(t)
	for name, ids := range batchIDSets() {
		for _, v := range []View{AsIs, ToBe} {
			var want []Parent
			for _, n := range ms[v].Nodes() {
				if slices.Contains(ids, n.ID) {
					want = append(want, Parent{ID: n.ID, Parent: n.Parent})
				}
			}
			got, err := x.Parents(v, ids)
			if err != nil {
				t.Fatalf("%s %s: %v", name, v, err)
			}
			same(t, "parents "+name+" "+string(v), got, want)
			// Parents is NodesLight's ID and parent, in NodesLight's order.
			ns, err := x.NodesLight(v, ids)
			if err != nil {
				t.Fatal(err)
			}
			var fromLight []Parent
			for _, n := range ns {
				fromLight = append(fromLight, Parent{ID: n.ID, Parent: n.Parent})
			}
			same(t, "parents against nodes light "+name+" "+string(v), got, fromLight)
		}
	}
	if _, err := x.Parents(View("what"), nil); err == nil {
		t.Fatal("unknown view: no error")
	}
}

func TestEachInFlight(t *testing.T) {
	x, _ := listIndex(t)
	sets, err := x.ChangeSets()
	if err != nil {
		t.Fatal(err)
	}
	var flying []id.ID
	for _, c := range sets {
		switch c.State {
		case changeset.StateDraft, changeset.StateProposed, changeset.StateAccepted:
			flying = append(flying, c.ID)
		}
	}
	if len(flying) < 2 {
		t.Fatalf("fixture: %d change sets in flight, want 2 or more", len(flying))
	}
	want, err := x.EntriesOf(flying)
	if err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	got := map[id.ID][]entry.Entry{}
	calls := 0
	err = x.EachInFlight(func(cs id.ID, es []entry.Entry) {
		mu.Lock()
		defer mu.Unlock()
		calls++
		got[cs] = es
	})
	if err != nil {
		t.Fatal(err)
	}
	if calls != len(want) {
		t.Fatalf("each called %d times, want once per change set with entries (%d)", calls, len(want))
	}
	same(t, "entries in flight", got, want)
}
