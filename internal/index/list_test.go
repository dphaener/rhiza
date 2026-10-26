package index

import (
	"errors"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/store"
)

// listIndex is queryLog with one more landed change set, cs8: hid(13)
// under b(3) with kid(14) under it, then hid retired, so a live node sits
// under a retired one; and link ln4 kid->root, retired. It returns the
// index and the as-is and to-be views computed directly, the oracles.
func listIndex(t *testing.T) (*Index, map[View]model.Model) {
	f := queryLog(t)
	cs := csID(8)
	f.add(cs, changeset.Open{Kind: changeset.KindIntention, Intent: "hide"})
	f.add(cs, changeset.Ops{Ops: []op.Op{
		mkCap(80, 13, "hid", 3, "k"),
		mkCap(81, 14, "kid", 13, "k"),
		{ID: opID(82), Type: op.TypeLinkCreate, Link: ln(4), LinkType: op.LinkServes, From: nd(14), To: nd(1)},
		{ID: opID(85), Type: op.TypeLinkCreate, Link: ln(1000), LinkType: op.LinkRelatesTo, From: nd(5), To: nd(9), Note: "far"},
	}})
	f.add(cs, changeset.Ops{Ops: []op.Op{
		{ID: opID(83), Type: op.TypeNodeRetire, Node: nd(13), Base: opID(80), Reason: "hid"},
		{ID: opID(84), Type: op.TypeLinkRetire, Link: ln(4), Base: opID(82)},
	}})
	f.add(cs, changeset.Propose{})
	f.add(cs, changeset.Accept{Rule: "owner"})
	f.add(cs, changeset.Land{})

	x := openT(t, f, f.dbPath(), nil)
	sets, _ := changeset.FoldLog(store.ReadLog(os.DirFS(f.gdir)))
	asIs, _ := model.AsIs(sets)
	toBe, _ := model.ToBe(sets)
	if n, ok := asIs.Node(nd(14)); !ok || n.Retired {
		t.Fatalf("fixture: kid %+v, %v", n, ok)
	}
	return x, map[View]model.Model{AsIs: asIs, ToBe: toBe}
}

// wantList filters m.Nodes() as q selects, ignoring After and Limit.
func wantList(m model.Model, q ListQuery) []model.Node {
	var out []model.Node
	for _, n := range m.Nodes() {
		if len(q.Kinds) > 0 && !slices.Contains(q.Kinds, n.Kind) ||
			q.HasParent && n.Parent != q.Parent || !q.Retired && n.Retired {
			continue
		}
		out = append(out, n)
	}
	return out
}

func TestList(t *testing.T) {
	x, ms := listIndex(t)
	for _, c := range []struct {
		name string
		q    ListQuery
	}{
		{"all live", ListQuery{}},
		{"all with retired", ListQuery{Retired: true}},
		{"capabilities", ListQuery{Kinds: []op.Kind{op.KindCapability}}},
		{"decisions", ListQuery{Kinds: []op.Kind{op.KindDecision}}},
		{"two kinds", ListQuery{Kinds: []op.Kind{op.KindDecision, op.KindCapability}, Retired: true}},
		{"no such kind", ListQuery{Kinds: []op.Kind{op.KindEvidence}}},
		{"under root, live", ListQuery{Parent: nd(1), HasParent: true}},
		{"under root, retired", ListQuery{Parent: nd(1), HasParent: true, Retired: true}},
		{"under retired hid, live", ListQuery{Parent: nd(13), HasParent: true}},
		{"no parent", ListQuery{HasParent: true}},
		{"no parent, capabilities", ListQuery{HasParent: true, Kinds: []op.Kind{op.KindCapability}}},
	} {
		for _, v := range []View{AsIs, ToBe} {
			t.Run(c.name+"/"+string(v), func(t *testing.T) {
				want := wantList(ms[v], c.q)
				// Pages of one, each After the last, concatenate to the set.
				var got []model.Node
				q := c.q
				q.Limit = 1
				for range len(want) + 2 {
					page, err := x.List(v, q)
					if err != nil {
						t.Fatal(err)
					}
					if len(page) > 1 {
						t.Fatalf("page of %d with limit 1", len(page))
					}
					if len(page) == 0 {
						break
					}
					got = append(got, page...)
					q.After = page[0].ID
				}
				same(t, "paged list", got, want)
				// One page big enough holds them all.
				q = c.q
				q.Limit = 1000
				all, err := x.List(v, q)
				if err != nil {
					t.Fatal(err)
				}
				same(t, "whole list", all, want)
				// A page of three from the second on.
				if len(want) > 1 {
					q.After, q.Limit = want[0].ID, 3
					page, err := x.List(v, q)
					if err != nil {
						t.Fatal(err)
					}
					same(t, "middle page", ids(page), ids(want[1:min(4, len(want))]))
				}
			})
		}
	}
}

func TestListErrors(t *testing.T) {
	x, _ := listIndex(t)
	for _, limit := range []int{0, -1} {
		_, err := x.List(AsIs, ListQuery{Limit: limit})
		if err == nil || errors.Is(err, ErrNotFound) {
			t.Fatalf("limit %d: %v", limit, err)
		}
		t.Logf("limit %d: %v", limit, err)
	}
	_, err := x.List(View("what"), ListQuery{Limit: 1})
	if err == nil {
		t.Fatal("unknown view: no error")
	}
	t.Logf("unknown view: %v", err)
}

// wantSubtree walks m.Children from root in pre-order, as Subtree is
// documented to.
func wantSubtree(m model.Model, root id.ID, depth int, retired bool) []TreeNode {
	var out []TreeNode
	var visit func(n model.Node, d int)
	visit = func(n model.Node, d int) {
		if n.Retired && !retired {
			return
		}
		out = append(out, TreeNode{Node: n, Depth: d})
		if depth > 0 && d >= depth {
			return
		}
		for _, c := range m.Children(n.ID) {
			visit(c, d+1)
		}
	}
	if root.IsZero() {
		for _, n := range m.Children(id.ID{}) {
			visit(n, 0)
		}
	} else if n, ok := m.Node(root); ok {
		visit(n, 0)
	}
	return out
}

func treeIDs(ts []TreeNode) []string {
	var out []string
	for _, n := range ts {
		out = append(out, strings.Repeat(".", n.Depth)+n.Node.Handle)
	}
	return out
}

func TestSubtree(t *testing.T) {
	x, ms := listIndex(t)
	for _, root := range []id.ID{{}, nd(1), nd(2), nd(3), nd(13), nd(5), nd(9), nd(8), nd(11)} {
		for _, depth := range []int{0, -1, 1, 2, 3, 1000} {
			for _, retired := range []bool{false, true} {
				for _, v := range []View{AsIs, ToBe} {
					want := wantSubtree(ms[v], root, depth, retired)
					got, err := x.Subtree(v, root, depth, retired)
					if v == AsIs && root == nd(11) {
						if !errors.Is(err, ErrNotFound) {
							t.Fatalf("to-be node in as-is: %v", err)
						}
						continue
					}
					if err != nil {
						t.Fatalf("Subtree(%s, %s, %d, %v): %v", v, root, depth, retired, err)
					}
					same(t, "subtree "+string(v)+" "+root.String(), got, want)
				}
			}
		}
	}

	// Spot checks of the oracle itself, by handle with a dot per level.
	got, _ := x.Subtree(AsIs, id.ID{}, 0, false)
	same(t, "whole as-is tree", treeIDs(got), []string{
		"root2", "far", "root", ".dup", ".a", "..dup", "..b", "...c", "....d"})
	got, _ = x.Subtree(AsIs, nd(3), 0, true)
	same(t, "b with retired", treeIDs(got), []string{"b", ".c", "..d", ".hid", "..kid"})
	got, _ = x.Subtree(ToBe, nd(1), 1, false)
	same(t, "root depth 1", treeIDs(got), []string{"root", ".dup", ".a", ".new"})
	got, _ = x.Subtree(AsIs, nd(9), 0, false)
	same(t, "outside the tree", treeIDs(got), []string{"dec"})
	got, _ = x.Subtree(AsIs, nd(8), 0, false)
	same(t, "retired root left out", treeIDs(got), nil)
}

func TestSubtreeErrors(t *testing.T) {
	x, _ := listIndex(t)
	_, err := x.Subtree(AsIs, nd(77), 0, false)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("unknown root: %v", err)
	}
	t.Logf("unknown root: %v", err)
	_, err = x.Subtree(View("what"), id.ID{}, 0, false)
	if err == nil {
		t.Fatal("unknown view: no error")
	}
	t.Logf("unknown view: %v", err)
}

func TestLink(t *testing.T) {
	x, ms := listIndex(t)
	for _, v := range []View{AsIs, ToBe} {
		for _, i := range []id.ID{ln(1), ln(2), ln(3), ln(4), ln(1000), ln(77), {}} {
			got, ok, err := x.Link(v, i)
			want, wok := ms[v].Link(i)
			if err != nil || ok != wok {
				t.Fatalf("Link(%s, %s) = %v, %v; want %v", v, i, ok, err, wok)
			}
			same(t, "link "+i.String(), got, want)
		}
	}
	if _, _, err := x.Link(View("what"), ln(1)); err == nil {
		t.Fatal("unknown view: no error")
	}
}

func TestResolveLink(t *testing.T) {
	x, ms := listIndex(t)
	full := ln(3).String()
	// ln1 to ln4 share all but the last character of their IDs.
	prefix := full[:len(full)-1]
	for _, c := range []struct {
		name string
		v    View
		ref  string
		want id.ID
		err  error
	}{
		{"full ID", AsIs, full, ln(3), nil},
		{"full ID of a retired link", AsIs, ln(4).String(), ln(4), nil},
		{"unique prefix", AsIs, ln(1000).String()[:len(full)-1], ln(1000), nil},
		{"ambiguous prefix", AsIs, prefix, id.ID{}, ErrAmbiguous},
		{"ambiguous short prefix", ToBe, full[:4], id.ID{}, ErrAmbiguous},
		{"prefix shorter than 4 not found", AsIs, full[:3], id.ID{}, ErrNotFound},
		{"lower-case not found", AsIs, strings.ToLower(ln(1000).String()), id.ID{}, ErrNotFound},
		{"unknown ID not found", AsIs, ln(77).String(), id.ID{}, ErrNotFound},
		{"node ID not found", AsIs, nd(1).String(), id.ID{}, ErrNotFound},
		{"empty not found", AsIs, "", id.ID{}, ErrNotFound},
	} {
		t.Run(c.name, func(t *testing.T) {
			l, err := x.ResolveLink(c.v, c.ref)
			if c.err != nil {
				if !errors.Is(err, c.err) {
					t.Fatalf("ResolveLink(%q) = %+v, %v; want %v", c.ref, l, err, c.err)
				}
				t.Logf("%v", err)
				return
			}
			want, _ := ms[c.v].Link(c.want)
			if err != nil {
				t.Fatalf("ResolveLink(%q): %v", c.ref, err)
			}
			same(t, "link", l, want)
		})
	}
	_, err := x.ResolveLink(AsIs, prefix)
	if want := strings.Join([]string{ln(1).String(), ln(2).String(), ln(3).String(), ln(4).String()}, ", "); err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("ambiguous error %v does not list %s", err, want)
	}
	if _, err := x.ResolveLink(View("what"), full); err == nil || errors.Is(err, ErrNotFound) {
		t.Fatalf("unknown view: %v", err)
	}
}
