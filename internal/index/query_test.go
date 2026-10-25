package index

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/model"
	"github.com/dphaener/genesis/internal/op"
	"github.com/dphaener/genesis/internal/store"
	"pgregory.net/rapid"
)

// queryLog writes the hand-built log the query table tests read:
//
//	cs1 landed: root(1) > a(2) > b(3) > c(4) > d(5); dup(6) under root and
//	  dup(7) under a; gone(8) under root, retired; decision dec(9); root2(10);
//	  far(1000) a root whose ID differs early; links ln1 a->root (serves),
//	  ln2 root->dec (relates-to), ln3 b->root (serves).
//	cs2 accepted: sets root's statement, creates new(11) under root.
//	cs3 draft: creates draft(12), a root.
//	cs4 proposed, with an orphan entry (fold problem) and a corrupt file.
//	cs5 withdrawn.
//	cs6, cs7 landed: both update root2 from its create (a conflict).
func queryLog(t *testing.T) *fixture {
	f := newFixture(t)
	cs1 := csID(1)
	f.add(cs1, changeset.Open{Kind: changeset.KindIntention, Intent: "seed"})
	f.add(cs1, changeset.Ops{Ops: []op.Op{
		mkCap(1, 1, "root", 0, "i"),
		mkCap(2, 2, "a", 1, "m"),
		mkCap(3, 3, "b", 2, "i"),
		mkCap(4, 4, "c", 3, "i"),
		mkCap(5, 5, "d", 4, "i"),
		mkCap(6, 6, "dup", 1, "a"),
		mkCap(7, 7, "dup", 2, "b"),
		mkCap(8, 8, "gone", 1, "z"),
		{ID: opID(9), Type: op.TypeNodeCreate, Node: nd(9), Kind: op.KindDecision, Handle: "dec"},
		mkCap(10, 10, "root2", 0, "b"),
		mkCap(11, 1000, "far", 0, "c"),
		{ID: opID(12), Type: op.TypeLinkCreate, Link: ln(1), LinkType: op.LinkServes, From: nd(2), To: nd(1)},
		{ID: opID(13), Type: op.TypeLinkCreate, Link: ln(2), LinkType: op.LinkRelatesTo, From: nd(1), To: nd(9), Note: "why"},
		{ID: opID(14), Type: op.TypeLinkCreate, Link: ln(3), LinkType: op.LinkServes, From: nd(3), To: nd(1)},
	}})
	f.add(cs1, changeset.Ops{Ops: []op.Op{
		{ID: opID(15), Type: op.TypeNodeRetire, Node: nd(8), Base: opID(8), Reason: "gone"},
	}})
	f.add(cs1, changeset.Propose{})
	f.add(cs1, changeset.Accept{Rule: "owner"})
	f.add(cs1, changeset.Land{})

	cs2 := csID(2)
	f.add(cs2, changeset.Open{Kind: changeset.KindIntention, Intent: "grow"})
	f.add(cs2, changeset.Ops{Ops: []op.Op{update(20, 1, opID(1), "hello"), mkCap(21, 11, "new", 1, "q")}})
	f.add(cs2, changeset.Propose{})
	f.add(cs2, changeset.Accept{Rule: "owner"})

	cs3 := csID(3)
	f.add(cs3, changeset.Open{Kind: changeset.KindIntention, Intent: "try"})
	f.add(cs3, changeset.Ops{Ops: []op.Op{mkCap(30, 12, "draft", 0, "d")}})

	cs4 := csID(4)
	f.add(cs4, changeset.Open{Kind: changeset.KindIntention, Intent: "ask"})
	f.add(cs4, changeset.Propose{})
	p, data := encodeEntry(t, eID(900), cs4, []id.ID{eID(899)}, changeset.Note{Text: "lost"})
	writeFile(t, f.gdir, p, data)
	writeFile(t, f.gdir, entryPath(cs4, eID(901)), []byte("{not json"))

	cs5 := csID(5)
	f.add(cs5, changeset.Open{Kind: changeset.KindIntention, Intent: "drop"})
	f.add(cs5, changeset.Withdraw{})

	for k, n := range map[int]int{6: 60, 7: 70} {
		cs := csID(k)
		f.add(cs, changeset.Open{Kind: changeset.KindRecord, Intent: "fix"})
		f.add(cs, changeset.Ops{Ops: []op.Op{update(n, 10, opID(10), "x")}})
		f.add(cs, changeset.Land{})
	}
	return f
}

func queryIndex(t *testing.T) (*fixture, *Index) {
	f := queryLog(t)
	return f, openT(t, f, f.dbPath(), nil)
}

func ids(ns []model.Node) []id.ID {
	var out []id.ID
	for _, n := range ns {
		out = append(out, n.ID)
	}
	return out
}

func linkIDs(ls []model.Link) []id.ID {
	var out []id.ID
	for _, l := range ls {
		out = append(out, l.ID)
	}
	return out
}

func TestQueryNode(t *testing.T) {
	_, x := queryIndex(t)
	for _, c := range []struct {
		name   string
		v      View
		i      id.ID
		ok     bool
		handle string
	}{
		{"live", AsIs, nd(2), true, "a"},
		{"retired", AsIs, nd(8), true, "gone"},
		{"to-be only, as-is", AsIs, nd(11), false, ""},
		{"to-be only, to-be", ToBe, nd(11), true, "new"},
		{"draft never stored", ToBe, nd(12), false, ""},
		{"zero", AsIs, id.ID{}, false, ""},
	} {
		t.Run(c.name, func(t *testing.T) {
			n, ok, err := x.Node(c.v, c.i)
			if err != nil || ok != c.ok || n.Handle != c.handle {
				t.Fatalf("Node = %+v, %v, %v", n, ok, err)
			}
		})
	}
	if _, _, err := x.Node(View("what"), nd(1)); err == nil {
		t.Fatal("unknown view: no error")
	}
	n, _, _ := x.Node(ToBe, nd(1))
	if string(n.Fields["statement"]) != `"hello"` {
		t.Fatalf("to-be root fields %v", n.Fields)
	}
}

func TestResolve(t *testing.T) {
	_, x := queryIndex(t)
	far := nd(1000).String()
	for _, c := range []struct {
		name string
		v    View
		ref  string
		want id.ID
		err  error
	}{
		{"full ID", AsIs, nd(3).String(), nd(3), nil},
		{"full ID of a retired node", AsIs, nd(8).String(), nd(8), nil},
		{"unique prefix", AsIs, far[:25], nd(1000), nil},
		{"prefix shorter than 4 not found", AsIs, far[:3], id.ID{}, ErrNotFound},
		{"lower-case prefix not found", AsIs, strings.ToLower(far[:25]), id.ID{}, ErrNotFound},
		{"lower-case full ID not found", AsIs, strings.ToLower(far), id.ID{}, ErrNotFound},
		{"ambiguous prefix", AsIs, far[:4], id.ID{}, ErrAmbiguous},
		{"handle", AsIs, "c", nd(4), nil},
		{"retired handle not found", AsIs, "gone", id.ID{}, ErrNotFound},
		{"ambiguous handle", AsIs, "dup", id.ID{}, ErrAmbiguous},
		{"to-be handle in as-is not found", AsIs, "new", id.ID{}, ErrNotFound},
		{"to-be handle", ToBe, "new", nd(11), nil},
		{"unknown ID not found", AsIs, nd(77).String(), id.ID{}, ErrNotFound},
		{"empty ref not found", AsIs, "", id.ID{}, ErrNotFound},
		{"LIKE wildcard not found", AsIs, "05%_", id.ID{}, ErrNotFound},
	} {
		t.Run(c.name, func(t *testing.T) {
			n, err := x.Resolve(c.v, c.ref)
			if c.err != nil {
				if !errors.Is(err, c.err) {
					t.Fatalf("Resolve(%q) = %+v, %v; want %v", c.ref, n, err, c.err)
				}
				t.Logf("%v", err)
				return
			}
			if err != nil || n.ID != c.want {
				t.Fatalf("Resolve(%q) = %v, %v; want %v", c.ref, n.ID, err, c.want)
			}
		})
	}
	_, err := x.Resolve(AsIs, "dup")
	if want := nd(6).String() + " (dup), " + nd(7).String() + " (dup)"; err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("ambiguous handle error %v does not list %s", err, want)
	}
	_, err = x.Resolve(AsIs, far[:4])
	if err == nil || !strings.Contains(err.Error(), nd(1).String()+" (root), "+nd(2).String()+" (a)") {
		t.Fatalf("ambiguous prefix error %v", err)
	}
}

func TestChildren(t *testing.T) {
	_, x := queryIndex(t)
	for _, c := range []struct {
		name   string
		v      View
		parent id.ID
		want   []id.ID
	}{
		{"roots by key, not the decision", AsIs, id.ID{}, []id.ID{nd(10), nd(1000), nd(1)}},
		{"under root by key, retired kept", AsIs, nd(1), []id.ID{nd(6), nd(2), nd(8)}},
		{"to-be adds", ToBe, nd(1), []id.ID{nd(6), nd(2), nd(11), nd(8)}},
		{"leaf", AsIs, nd(5), nil},
		{"unknown", AsIs, nd(77), nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, err := x.Children(c.v, c.parent)
			if err != nil {
				t.Fatal(err)
			}
			same(t, "children", ids(got), c.want)
		})
	}
}

func TestAncestors(t *testing.T) {
	_, x := queryIndex(t)
	for _, c := range []struct {
		name string
		i    id.ID
		want []id.ID
	}{
		{"deep chain", nd(5), []id.ID{nd(1), nd(2), nd(3), nd(4)}},
		{"one up", nd(2), []id.ID{nd(1)}},
		{"root", nd(1), nil},
		{"outside the tree", nd(9), nil},
		{"retired", nd(8), []id.ID{nd(1)}},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, err := x.Ancestors(AsIs, c.i)
			if err != nil {
				t.Fatal(err)
			}
			same(t, "ancestors", ids(got), c.want)
		})
	}
	if _, err := x.Ancestors(AsIs, nd(77)); !errors.Is(err, ErrNotFound) {
		t.Fatalf("unknown node: %v", err)
	}
	if _, err := x.Ancestors(AsIs, nd(11)); !errors.Is(err, ErrNotFound) {
		t.Fatalf("to-be node in as-is: %v", err)
	}
	got, err := x.Ancestors(ToBe, nd(11))
	if err != nil || !reflect.DeepEqual(ids(got), []id.ID{nd(1)}) {
		t.Fatalf("to-be ancestors %v, %v", ids(got), err)
	}
}

func TestLinks(t *testing.T) {
	_, x := queryIndex(t)
	for _, c := range []struct {
		name string
		node id.ID
		dir  Direction
		want []id.ID
	}{
		{"from", nd(1), From, []id.ID{ln(2)}},
		{"to", nd(1), To, []id.ID{ln(1), ln(3)}},
		{"both", nd(1), Both, []id.ID{ln(1), ln(2), ln(3)}},
		{"to the decision", nd(9), To, []id.ID{ln(2)}},
		{"none", nd(5), Both, nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, err := x.Links(AsIs, c.node, c.dir)
			if err != nil {
				t.Fatal(err)
			}
			same(t, "links", linkIDs(got), c.want)
		})
	}
	if _, err := x.Links(AsIs, nd(1), Direction(9)); err == nil {
		t.Fatal("unknown direction: no error")
	}
}

func TestNeighbours(t *testing.T) {
	_, x := queryIndex(t)
	for _, c := range []struct {
		name       string
		v          View
		i          id.ID
		prev, next id.ID
	}{
		{"first", AsIs, nd(1), id.ID{}, nd(2)},
		{"middle", AsIs, nd(5), nd(4), nd(6)},
		{"after a retired node", AsIs, nd(9), nd(8), nd(10)},
		{"last", AsIs, nd(1000), nd(10), id.ID{}},
		{"not in the view", AsIs, nd(11), nd(10), nd(1000)},
		{"to-be has it", ToBe, nd(10), nd(9), nd(11)},
		{"to-be, itself", ToBe, nd(11), nd(10), nd(1000)},
		{"draft never stored", ToBe, nd(12), nd(11), nd(1000)},
		{"zero", AsIs, id.ID{}, id.ID{}, nd(1)},
		{"past the end", AsIs, nd(5000), nd(1000), id.ID{}},
	} {
		t.Run(c.name, func(t *testing.T) {
			prev, next, err := x.Neighbours(c.v, c.i)
			if err != nil || prev != c.prev || next != c.next {
				t.Fatalf("Neighbours = %s, %s, %v; want %s, %s", prev, next, err, c.prev, c.next)
			}
		})
	}
	if _, _, err := x.Neighbours(View("what"), nd(1)); err == nil {
		t.Fatal("unknown view: no error")
	}

	f := newFixture(t)
	empty := openT(t, f, f.dbPath(), nil)
	for _, v := range []View{AsIs, ToBe} {
		prev, next, err := empty.Neighbours(v, nd(1))
		if err != nil || !prev.IsZero() || !next.IsZero() {
			t.Fatalf("empty %s: Neighbours = %s, %s, %v", v, prev, next, err)
		}
	}
}

// TestNeighboursPlan: both neighbour queries are a range search of the
// nodes primary key.
func TestNeighboursPlan(t *testing.T) {
	_, x := queryIndex(t)
	for _, q := range []string{prevQuery, nextQuery} {
		rows, err := x.db.Query(`EXPLAIN QUERY PLAN `+q, string(AsIs), idText(nd(5)))
		if err != nil {
			t.Fatal(err)
		}
		var plan []string
		for rows.Next() {
			var a, b, c int
			var detail string
			if err := rows.Scan(&a, &b, &c, &detail); err != nil {
				t.Fatal(err)
			}
			plan = append(plan, detail)
		}
		rows.Close()
		if len(plan) != 1 || !strings.Contains(plan[0], "USING PRIMARY KEY (view=? AND id") {
			t.Fatalf("plan of %s:\n%s", q, strings.Join(plan, "\n"))
		}
	}
}

func TestQueryLink(t *testing.T) {
	f := newFixture(t)
	cs := csID(1)
	f.add(cs, changeset.Open{Kind: changeset.KindRecord, Intent: "seed"})
	f.add(cs, changeset.Ops{Ops: []op.Op{
		mkCap(1, 1, "a", 0, "i"),
		mkCap(2, 2, "b", 0, "j"),
		{ID: opID(3), Type: op.TypeLinkCreate, Link: ln(1), LinkType: op.LinkServes, From: nd(1), To: nd(2)},
		{ID: opID(4), Type: op.TypeLinkCreate, Link: ln(2), LinkType: op.LinkRelatesTo, From: nd(2), To: nd(1), Note: "why"},
		{ID: opID(5), Type: op.TypeLinkRetire, Link: ln(2), Base: opID(4)},
	}})
	f.add(cs, changeset.Land{})
	x := openT(t, f, f.dbPath(), nil)
	for _, c := range []struct {
		name    string
		i       id.ID
		ok      bool
		retired bool
		from    id.ID
	}{
		{"live", ln(1), true, false, nd(1)},
		{"retired", ln(2), true, true, nd(2)},
		{"missing", ln(3), false, false, id.ID{}},
		{"a node is no link", nd(1), false, false, id.ID{}},
	} {
		t.Run(c.name, func(t *testing.T) {
			l, ok, err := x.Link(AsIs, c.i)
			if err != nil || ok != c.ok || l.Retired != c.retired || l.From != c.from {
				t.Fatalf("Link = %+v, %v, %v", l, ok, err)
			}
		})
	}
	want := storedLinks(t, x, viewAsIs)
	got, _, _ := x.Link(AsIs, ln(2))
	same(t, "retired link", got, want[1])
	if _, _, err := x.Link(View("what"), ln(1)); err == nil {
		t.Fatal("unknown view: no error")
	}
}

func TestQueryChangeSets(t *testing.T) {
	f, x := queryIndex(t)
	got, err := x.ChangeSets()
	if err != nil {
		t.Fatal(err)
	}
	want := map[id.ID]changeset.State{
		csID(1): changeset.StateLanded, csID(2): changeset.StateAccepted, csID(3): changeset.StateDraft,
		csID(4): changeset.StateProposed, csID(5): changeset.StateWithdrawn,
		csID(6): changeset.StateLanded, csID(7): changeset.StateLanded,
	}
	if len(got) != len(want) {
		t.Fatalf("change sets %+v", got)
	}
	for i, c := range got {
		if i > 0 && id.Compare(got[i-1].ID, c.ID) >= 0 {
			t.Fatalf("not sorted: %+v", got)
		}
		if want[c.ID] != c.State {
			t.Fatalf("%s state %s, want %s", c.ID, c.State, want[c.ID])
		}
	}
	if got[0].Kind != changeset.KindIntention || got[0].Intent != "seed" || got[0].OpenEntry.IsZero() || got[0].TerminalEntry.IsZero() {
		t.Fatalf("cs1 %+v", got[0])
	}
	if got[5].Kind != changeset.KindRecord {
		t.Fatalf("cs6 %+v", got[5])
	}

	sets, _ := changeset.FoldLog(store.ReadLog(os.DirFS(f.gdir)))
	var fromModel []ChangeSetInfo
	for _, s := range sets {
		fromModel = append(fromModel, ChangeSetInfo{s.ID, s.State, s.Open.Kind, s.Open.Intent, s.OpenEntry, s.SupersededBy, s.TerminalEntry})
	}
	same(t, "change sets", got, fromModel)
}

func TestQueryEntries(t *testing.T) {
	f, x := queryIndex(t)
	l := store.ReadLog(os.DirFS(f.gdir))
	for _, cs := range []id.ID{csID(1), csID(4)} {
		got, err := x.Entries(cs)
		if err != nil {
			t.Fatal(err)
		}
		var want []entry.Entry
		for _, e := range l.Entries {
			if e.ChangeSet == cs {
				want = append(want, e)
			}
		}
		same(t, "entries", got, want)
	}
	if got, _ := x.Entries(csID(1)); len(got) != 6 {
		t.Fatalf("cs1 has %d entries", len(got))
	}
	// cs4: open, propose and the orphan; the corrupt file is not an entry.
	if got, _ := x.Entries(csID(4)); len(got) != 3 {
		t.Fatalf("cs4 has %d entries", len(got))
	}
	if got, err := x.Entries(csID(77)); got != nil || err != nil {
		t.Fatalf("unknown change set: %v, %v", got, err)
	}
}

func TestQueryProblems(t *testing.T) {
	f, x := queryIndex(t)
	got, err := x.Problems()
	if err != nil {
		t.Fatal(err)
	}
	order := map[string]int{sourceLog: 0, sourceFold: 1, viewAsIs: 2, viewToBe: 3}
	seen := map[string]Problem{}
	for i, p := range got {
		if i > 0 && order[got[i-1].Source] > order[p.Source] {
			t.Fatalf("not sorted by source: %+v", got)
		}
		seen[p.Source+" "+p.Sentinel] = p
	}
	if p, ok := seen["log invalid"]; !ok || p.Path != entryPath(csID(4), eID(901)) {
		t.Fatalf("no corrupt-file problem: %+v", got)
	}
	if p, ok := seen["fold order"]; !ok || p.ChangeSet != csID(4) || p.Entry != eID(900) {
		t.Fatalf("no fold order problem: %+v", got)
	}
	for _, v := range []string{viewAsIs, viewToBe} {
		if p, ok := seen[v+" conflict"]; !ok || p.ChangeSet.IsZero() || p.Op.IsZero() {
			t.Fatalf("no %s conflict: %+v", v, got)
		}
	}

	// Equal, in order, to the stored rows of each source.
	var want []Problem
	for _, src := range []string{sourceLog, sourceFold, viewAsIs, viewToBe} {
		for _, p := range storedProblems(t, x, src) {
			want = append(want, Problem{p.Source, p.Path, parseID(p.ChangeSet), parseID(p.Entry), parseID(p.Op), p.Sentinel, p.Message})
		}
	}
	same(t, "problems", got, want)
	checkModel(t, x, f.gdir)
}

func TestWhatIf(t *testing.T) {
	f, x := queryIndex(t)
	sets, _ := changeset.FoldLog(store.ReadLog(os.DirFS(f.gdir)))

	m, ps, err := x.WhatIf(csID(3))
	if err != nil {
		t.Fatal(err)
	}
	if n, ok := m.Node(nd(12)); !ok || n.Handle != "draft" {
		t.Fatalf("what-if of the draft lacks its node: %+v %v", n, ok)
	}
	wm, wps := model.WhatIf(sets, csID(3))
	same(t, "what-if draft", m, wm)
	same(t, "what-if draft problems", ps, wps)

	m, ps, err = x.WhatIf(csID(1))
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, p := range ps {
		found = found || (errors.Is(p.Err, model.ErrState) && p.ChangeSet == csID(1))
	}
	if !found {
		t.Fatalf("what-if of a landed set: no ErrState problem in %v", ps)
	}
	wm, wps = model.WhatIf(sets, csID(1))
	same(t, "what-if landed", m, wm)
	same(t, "what-if landed problems", ps, wps)

	if _, ps, _ = x.WhatIf(csID(77)); !hasErr(ps, model.ErrMissing) {
		t.Fatalf("what-if of an unknown set: %v", ps)
	}
}

func hasErr(ps []model.Problem, target error) bool {
	for _, p := range ps {
		if errors.Is(p.Err, target) {
			return true
		}
	}
	return false
}

// modelAncestors walks the parent chain of i in m, root first.
func modelAncestors(m model.Model, i id.ID) []model.Node {
	var out []model.Node
	n, _ := m.Node(i)
	for !n.Parent.IsZero() {
		p, ok := m.Node(n.Parent)
		if !ok {
			break
		}
		out = append([]model.Node{p}, out...)
		n = p
	}
	return out
}

// modelLinks filters m's links at node in direction dir.
func modelLinks(m model.Model, node id.ID, dir Direction) []model.Link {
	var out []model.Link
	for _, l := range m.Links() {
		if (dir != To && l.From == node) || (dir != From && l.To == node) {
			out = append(out, l)
		}
	}
	return out
}

// TestPropertyQueryAgrees: on drawn logs, every query of each stored view
// equals the same question asked of model.AsIs/ToBe over
// FoldLog(ReadLog), and WhatIf equals model.WhatIf.
func TestPropertyQueryAgrees(t *testing.T) {
	var nodes, whatIfs int
	rapid.Check(t, func(t *rapid.T) {
		files := genFiles(t)
		root, err := os.MkdirTemp("", "index-query-")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(root)
		gdir := filepath.Join(root, store.Dir)
		if err := os.MkdirAll(filepath.Join(gdir, store.LogDir), 0o755); err != nil {
			t.Fatal(err)
		}
		for _, f := range files {
			writeFile(t, gdir, f.path, f.data)
		}
		x, err := Open(root, filepath.Join(root, "index.db"))
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		defer x.Close()

		sets, _ := changeset.FoldLog(store.ReadLog(os.DirFS(gdir)))
		for _, v := range []struct {
			view View
			f    func([]changeset.ChangeSet) (model.Model, []model.Problem)
		}{{AsIs, model.AsIs}, {ToBe, model.ToBe}} {
			m, _ := v.f(sets)
			roots, err := x.Children(v.view, id.ID{})
			if err != nil {
				t.Fatal(err)
			}
			same(t, "roots", roots, m.Children(id.ID{}))
			for _, want := range m.Nodes() {
				nodes++
				got, ok, err := x.Node(v.view, want.ID)
				if err != nil || !ok {
					t.Fatalf("Node(%s) = %v, %v", want.ID, ok, err)
				}
				same(t, "node", got, want)
				kids, err := x.Children(v.view, want.ID)
				if err != nil {
					t.Fatal(err)
				}
				same(t, "children", kids, m.Children(want.ID))
				anc, err := x.Ancestors(v.view, want.ID)
				if err != nil {
					t.Fatal(err)
				}
				same(t, "ancestors", anc, modelAncestors(m, want.ID))
				for _, dir := range []Direction{From, To, Both} {
					ls, err := x.Links(v.view, want.ID, dir)
					if err != nil {
						t.Fatal(err)
					}
					same(t, "links", ls, modelLinks(m, want.ID, dir))
				}
			}
			for _, s := range sets {
				ls, err := x.Links(v.view, s.ID, Both)
				if err != nil {
					t.Fatal(err)
				}
				same(t, "change set links", ls, modelLinks(m, s.ID, Both))
			}
		}
		for _, s := range sets {
			switch s.State {
			case changeset.StateDraft, changeset.StateProposed, changeset.StateAccepted:
			default:
				continue
			}
			whatIfs++
			m, ps, err := x.WhatIf(s.ID)
			if err != nil {
				t.Fatal(err)
			}
			wm, wps := model.WhatIf(sets, s.ID)
			same(t, "what-if", m, wm)
			same(t, "what-if problems", ps, wps)
		}
	})
	t.Logf("nodes checked %d, what-ifs %d", nodes, whatIfs)
	if nodes == 0 || whatIfs == 0 {
		t.Errorf("vacuous run: nodes %d, what-ifs %d", nodes, whatIfs)
	}
}
