package index

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"pgregory.net/rapid"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/store"
)

// modelSubtree walks m from root as Overlay.Subtree promises.
func modelSubtree(m model.Model, root id.ID, depth int, retired bool) []TreeNode {
	if depth <= 0 || depth > maxDepth {
		depth = maxDepth
	}
	tops := m.Children(id.ID{})
	if !root.IsZero() {
		n, _ := m.Node(root)
		tops = []model.Node{n}
	}
	var out []TreeNode
	done := map[id.ID]bool{}
	var visit func(n model.Node, d int)
	visit = func(n model.Node, d int) {
		if done[n.ID] || n.Retired && !retired {
			return
		}
		done[n.ID] = true
		out = append(out, TreeNode{Node: n, Depth: d})
		if d >= depth {
			return
		}
		for _, c := range m.Children(n.ID) {
			visit(c, d+1)
		}
	}
	for _, n := range tops {
		visit(n, 0)
	}
	return out
}

// TestPropertyOverlayAgrees: on drawn logs, every query of a what-if
// overlay answers what the same question of the full model.WhatIf does.
func TestPropertyOverlayAgrees(t *testing.T) {
	var nodes, overlays int
	rapid.Check(t, func(t *rapid.T) {
		files := genFiles(t)
		root, err := os.MkdirTemp("", "index-overlay-")
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
		for _, s := range sets {
			if !inFlightState(s.State) {
				continue
			}
			overlays++
			o, _, err := x.WhatIfOverlay(s.ID)
			if err != nil {
				t.Fatal(err)
			}
			m, _ := model.WhatIf(sets, s.ID)
			roots, err := o.Children(id.ID{})
			if err != nil {
				t.Fatal(err)
			}
			same(t, "roots", roots, m.Children(id.ID{}))
			all := m.Nodes()
			var ids []id.ID
			for _, want := range all {
				nodes++
				ids = append(ids, want.ID)
				got, ok, err := o.Node(want.ID)
				if err != nil || !ok {
					t.Fatalf("Node(%s) = %v, %v", want.ID, ok, err)
				}
				same(t, "node", got, want)
				if r, err := o.Resolve(want.ID.String()); err != nil {
					t.Fatal(err)
				} else {
					same(t, "resolve", r, want)
				}
				kids, err := o.Children(want.ID)
				if err != nil {
					t.Fatal(err)
				}
				same(t, "children", kids, m.Children(want.ID))
				anc, err := o.Ancestors(want.ID)
				if err != nil {
					t.Fatal(err)
				}
				same(t, "ancestors", anc, modelAncestors(m, want.ID))
				for _, dir := range []Direction{From, To, Both} {
					ls, err := o.Links(want.ID, dir)
					if err != nil {
						t.Fatal(err)
					}
					same(t, "links", ls, modelLinks(m, want.ID, dir))
				}
				for _, retired := range []bool{false, true} {
					for _, depth := range []int{0, 1} {
						ts, err := o.Subtree(want.ID, depth, retired)
						if err != nil {
							t.Fatal(err)
						}
						want := modelSubtree(m, want.ID, depth, retired)
						if len(want) == 0 {
							want = nil
						}
						same(t, "subtree", ts, want)
					}
				}
			}
			for _, retired := range []bool{false, true} {
				ts, err := o.Subtree(id.ID{}, 0, retired)
				if err != nil {
					t.Fatal(err)
				}
				want := modelSubtree(m, id.ID{}, 0, retired)
				if len(want) == 0 {
					want = nil
				}
				same(t, "whole tree", ts, want)
			}
			for _, l := range m.Links() {
				got, ok, err := o.Link(l.ID)
				if err != nil || !ok {
					t.Fatalf("Link(%s) = %v, %v", l.ID, ok, err)
				}
				same(t, "link", got, l)
				if r, err := o.ResolveLink(l.ID.String()); err != nil {
					t.Fatal(err)
				} else {
					same(t, "resolve link", r, l)
				}
			}
			// Batches over every node, and a few more IDs no node has.
			ids = append(slices.Clone(ids), csID(77))
			light, err := o.NodesLight(ids)
			if err != nil {
				t.Fatal(err)
			}
			var wantLight, wantFields []model.Node
			for _, n := range all {
				n.Fields, n.History = nil, nil
				wantLight = append(wantLight, n)
			}
			same(t, "light", light, wantLight)
			fields, err := o.NodesFields(ids)
			if err != nil {
				t.Fatal(err)
			}
			for _, n := range all {
				n.History = nil
				wantFields = append(wantFields, n)
			}
			same(t, "fields", fields, wantFields)
			for _, dir := range []Direction{From, To, Both} {
				ls, err := o.LinksAt(ids, dir)
				if err != nil {
					t.Fatal(err)
				}
				var want []model.Link
				for _, l := range m.Links() {
					from, to := slices.Contains(ids, l.From), slices.Contains(ids, l.To)
					if dir == From && from || dir == To && to || dir == Both && (from || to) {
						want = append(want, l)
					}
				}
				same(t, "links at", ls, want)
			}
			// Pages of the listing, three at a time.
			var after id.ID
			var paged []model.Node
			for {
				page, err := o.List(ListQuery{Retired: true, After: after, Limit: 3})
				if err != nil {
					t.Fatal(err)
				}
				paged = append(paged, page...)
				if len(page) < 3 {
					break
				}
				after = page[len(page)-1].ID
			}
			same(t, "list", paged, all)
		}
	})
	t.Logf("nodes checked %d, overlays %d", nodes, overlays)
	if nodes == 0 || overlays == 0 {
		t.Errorf("vacuous run: nodes %d, overlays %d", nodes, overlays)
	}
}
