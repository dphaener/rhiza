package index

import (
	"cmp"
	"database/sql"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

// Overlay is the what-if view of one change set, read as the stored as-is
// rows overridden by the nodes and links the change set's operations (and
// those of the in-flight change sets it depends on) read or wrote. It never
// holds the as-is view: building it loads only the rows those operations
// reach, and every query asks the index for the as-is rows and corrects
// them with the overrides. Each query answers what the same question of the
// full model.WhatIf would answer.
//
// An Overlay reads through the Index it came from, which must stay open,
// and is not safe for concurrent use.
type Overlay struct {
	x     *Index
	nodes map[id.ID]model.Node
	links map[id.ID]model.Link
}

// WhatIfOverlay returns the overlay of the what-if view of change set cs,
// and the problems of the change sets' own operations (model.WhatIf's,
// but for the duplicate-handle ones, which need the whole view). When cs
// is not in the log, or not in flight, the overlay is the as-is view and
// the problems are ErrMissing or ErrState, as model.WhatIf gives them.
func (x *Index) WhatIfOverlay(cs id.ID) (*Overlay, []model.Problem, error) {
	return x.WhatIfOverlayOver(cs, nil)
}

// WhatIfOverlayOver is WhatIfOverlay with change set cs taken as over
// instead of folded from its stored entries, when over is not nil: the
// view of a change set with entries not yet written (a dry run's, or a
// record change set's), which need not be stored at all.
func (x *Index) WhatIfOverlayOver(cs id.ID, over *changeset.ChangeSet) (*Overlay, []model.Problem, error) {
	tx, err := x.db.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	sets, err := whatIfSets(tx, cs, over)
	if err != nil {
		return nil, nil, err
	}
	seed, err := loadSeedBase(tx)
	if err != nil {
		return nil, nil, err
	}
	d, ps, err := model.WhatIfDelta(seed, txSource{tx}, sets, cs)
	if err != nil {
		return nil, nil, err
	}
	o := &Overlay{x: x, nodes: make(map[id.ID]model.Node, len(d.Nodes)), links: make(map[id.ID]model.Link, len(d.Links))}
	for _, n := range d.Nodes {
		o.nodes[n.ID] = n
	}
	for _, l := range d.Links {
		o.links[l.ID] = l
	}
	return o, ps, nil
}

// txSource gives a delta computation the stored as-is rows.
type txSource struct{ q querier }

func (s txSource) SeedNode(i id.ID) (model.SeedNode, bool, error) {
	rows, err := s.q.Query(nodeSelect[:len(nodeSelect)-len(" FROM nodes")]+`, aspects, placed FROM nodes WHERE view = ? AND id = ?`,
		viewAsIs, idText(i))
	if err != nil {
		return model.SeedNode{}, false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return model.SeedNode{}, false, rows.Err()
	}
	var asp, placed sql.NullString
	n, err := scanNode(rows, &asp, &placed)
	if err != nil {
		return model.SeedNode{}, false, err
	}
	sn := model.SeedNode{Node: n}
	if sn.Aspects, err = decodeAspects(asp); err != nil {
		return model.SeedNode{}, false, err
	}
	if parts := strings.Fields(placed.String); len(parts) == 3 {
		sn.Placed.Seq, _ = strconv.Atoi(parts[0])
		sn.Placed.Ref = model.Ref{ChangeSet: parseID(parts[1]), Op: parseID(parts[2])}
	} else if len(parts) == 2 {
		sn.Placed.Seq, _ = strconv.Atoi(parts[0])
	}
	return sn, true, nil
}

func (s txSource) SeedLink(i id.ID) (model.SeedLink, bool, error) {
	rows, err := s.q.Query(linkSelect[:len(linkSelect)-len(" FROM links")]+`, aspects FROM links WHERE view = ? AND id = ?`,
		viewAsIs, idText(i))
	if err != nil {
		return model.SeedLink{}, false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return model.SeedLink{}, false, rows.Err()
	}
	var asp sql.NullString
	l, err := scanLink(rows, &asp)
	if err != nil {
		return model.SeedLink{}, false, err
	}
	sl := model.SeedLink{Link: l}
	if sl.Aspects, err = decodeAspects(asp); err != nil {
		return model.SeedLink{}, false, err
	}
	return sl, true, nil
}

// HandleTaken reports which live node has handle h, the smallest ID of
// them, if any.
func (o *Overlay) HandleTaken(h string) (id.ID, bool, error) {
	ns, err := o.x.nodes(nodeSelect+` INDEXED BY nodes_handle WHERE view = ? AND handle = ? AND retired = 0`, viewAsIs, h)
	if err != nil {
		return id.ID{}, false, err
	}
	ns = append(o.dropNodes(ns), o.pickNodes(func(n model.Node) bool { return !n.Retired && n.Handle == h })...)
	if len(ns) == 0 {
		return id.ID{}, false, nil
	}
	slices.SortFunc(ns, nodeByID)
	return ns[0].ID, true, nil
}

// Overridden reports whether the overlay overrides node or link i.
func (o *Overlay) Overridden(i id.ID) bool {
	_, n := o.nodes[i]
	_, l := o.links[i]
	return n || l
}

// byID sorts nodes by ID.
func nodeByID(a, b model.Node) int { return id.Compare(a.ID, b.ID) }

func linkByID(a, b model.Link) int { return id.Compare(a.ID, b.ID) }

func nodeByKey(a, b model.Node) int {
	return cmp.Or(strings.Compare(a.Key, b.Key), id.Compare(a.ID, b.ID))
}

// dropNodes returns ns without the nodes the overlay overrides.
func (o *Overlay) dropNodes(ns []model.Node) []model.Node {
	return slices.DeleteFunc(ns, func(n model.Node) bool { _, ok := o.nodes[n.ID]; return ok })
}

// dropLinks returns ls without the links the overlay overrides.
func (o *Overlay) dropLinks(ls []model.Link) []model.Link {
	return slices.DeleteFunc(ls, func(l model.Link) bool { _, ok := o.links[l.ID]; return ok })
}

// pickNodes returns the overridden nodes that keep accepts, a copy each,
// sorted by ID.
func (o *Overlay) pickNodes(keep func(model.Node) bool) []model.Node {
	var out []model.Node
	for _, n := range o.nodes {
		if keep(n) {
			out = append(out, n)
		}
	}
	slices.SortFunc(out, nodeByID)
	return out
}

// pickLinks is pickNodes for links.
func (o *Overlay) pickLinks(keep func(model.Link) bool) []model.Link {
	var out []model.Link
	for _, l := range o.links {
		if keep(l) {
			out = append(out, l)
		}
	}
	slices.SortFunc(out, linkByID)
	return out
}

// Node returns the node with ID i, live or retired.
func (o *Overlay) Node(i id.ID) (model.Node, bool, error) {
	if n, ok := o.nodes[i]; ok {
		return n, true, nil
	}
	return o.x.Node(AsIs, i)
}

// Link returns the link with ID i, live or retired.
func (o *Overlay) Link(i id.ID) (model.Link, bool, error) {
	if l, ok := o.links[i]; ok {
		return l, true, nil
	}
	return o.x.Link(AsIs, i)
}

// Resolve finds a node as Index.Resolve does.
func (o *Overlay) Resolve(ref string) (model.Node, error) {
	if i, err := id.Parse(ref); err == nil {
		n, ok, err := o.Node(i)
		if err != nil {
			return model.Node{}, err
		}
		if ok {
			return n, nil
		}
	}
	if len(ref) >= minPrefix {
		ns, err := o.x.nodes(nodeSelect+` WHERE view = ? AND id >= ? AND id < ? ORDER BY id`, viewAsIs, ref, ref+"~")
		if err != nil {
			return model.Node{}, err
		}
		ns = append(o.dropNodes(ns), o.pickNodes(func(n model.Node) bool { return strings.HasPrefix(n.ID.String(), ref) })...)
		if len(ns) > 0 {
			slices.SortFunc(ns, nodeByID)
			return one(ref, ns)
		}
	}
	ns, err := o.x.nodes(nodeSelect+` INDEXED BY nodes_handle WHERE view = ? AND handle = ? AND retired = 0`, viewAsIs, ref)
	if err != nil {
		return model.Node{}, err
	}
	ns = append(o.dropNodes(ns), o.pickNodes(func(n model.Node) bool { return !n.Retired && n.Handle == ref })...)
	slices.SortFunc(ns, nodeByID)
	if len(ns) == 0 {
		return model.Node{}, fmt.Errorf("%w: no node matches %q", ErrNotFound, ref)
	}
	return one(ref, ns)
}

// ResolveLink finds a link as Index.ResolveLink does.
func (o *Overlay) ResolveLink(ref string) (model.Link, error) {
	if i, err := id.Parse(ref); err == nil {
		l, ok, err := o.Link(i)
		if err != nil {
			return model.Link{}, err
		}
		if ok {
			return l, nil
		}
	}
	var ls []model.Link
	if len(ref) >= minPrefix {
		var err error
		ls, err = o.x.links(linkSelect+` WHERE view = ? AND id >= ? AND id < ? ORDER BY id`, viewAsIs, ref, ref+"~")
		if err != nil {
			return model.Link{}, err
		}
		ls = append(o.dropLinks(ls), o.pickLinks(func(l model.Link) bool { return strings.HasPrefix(l.ID.String(), ref) })...)
		slices.SortFunc(ls, linkByID)
	}
	switch len(ls) {
	case 0:
		return model.Link{}, fmt.Errorf("%w: no link matches %q", ErrNotFound, ref)
	case 1:
		return ls[0], nil
	}
	names := make([]string, len(ls))
	for i, l := range ls {
		names[i] = l.ID.String()
	}
	return model.Link{}, fmt.Errorf("%w: %q matches %d links: %s", ErrAmbiguous, ref, len(ls), strings.Join(names, ", "))
}

// Children returns the nodes directly under parent, live and retired,
// sorted by key then ID; the zero parent gives the roots.
func (o *Overlay) Children(parent id.ID) ([]model.Node, error) {
	ns, err := o.x.Children(AsIs, parent)
	if err != nil {
		return nil, err
	}
	ns = append(o.dropNodes(ns), o.pickNodes(func(n model.Node) bool {
		return op.HasParent(n.Kind) && n.Parent == parent
	})...)
	slices.SortFunc(ns, nodeByKey)
	if len(ns) == 0 {
		return nil, nil
	}
	return ns, nil
}

// Ancestors returns node i's parent chain, the root first, without i; nil
// for a root or a node outside the tree. An unknown i is an error wrapping
// ErrNotFound.
func (o *Overlay) Ancestors(i id.ID) ([]model.Node, error) {
	n, ok, err := o.Node(i)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("%w: no node %s", ErrNotFound, i)
	}
	var out []model.Node
	seen := map[id.ID]bool{}
	for p := n.Parent; !p.IsZero() && !seen[p] && len(out) < maxDepth; {
		seen[p] = true
		pn, ok, err := o.Node(p)
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		out = append(out, pn)
		p = pn.Parent
	}
	slices.Reverse(out)
	return out, nil
}

// Links returns the links at node in direction dir, live and retired,
// sorted by ID; node may be a change set.
func (o *Overlay) Links(node id.ID, dir Direction) ([]model.Link, error) {
	return o.LinksAt([]id.ID{node}, dir)
}

// LinksAt returns the links at any of ids in direction dir, live and
// retired, sorted by ID, each once.
func (o *Overlay) LinksAt(ids []id.ID, dir Direction) ([]model.Link, error) {
	ls, err := o.x.LinksAt(AsIs, ids, dir)
	if err != nil {
		return nil, err
	}
	at := make(map[id.ID]bool, len(ids))
	for _, i := range ids {
		at[i] = true
	}
	ls = append(o.dropLinks(ls), o.pickLinks(func(l model.Link) bool {
		from, to := at[l.From], at[l.To]
		return dir == From && from || dir == To && to || dir == Both && (from || to)
	})...)
	slices.SortFunc(ls, linkByID)
	return ls, nil
}

// LinksTo returns the links of type t whose To end is one of ids, live and
// retired, sorted by ID.
func (o *Overlay) LinksTo(ids []id.ID, t op.LinkType) ([]model.Link, error) {
	ls, err := o.x.LinksTo(AsIs, ids, t)
	if err != nil {
		return nil, err
	}
	to := make(map[id.ID]bool, len(ids))
	for _, i := range ids {
		to[i] = true
	}
	ls = append(o.dropLinks(ls), o.pickLinks(func(l model.Link) bool { return l.Type == t && to[l.To] })...)
	slices.SortFunc(ls, linkByID)
	return ls, nil
}

// splitIDs returns the ids the overlay overrides as nodes, and the rest.
func (o *Overlay) splitIDs(ids []id.ID) (over []model.Node, rest []id.ID) {
	for _, i := range ids {
		if n, ok := o.nodes[i]; ok {
			over = append(over, n)
		} else {
			rest = append(rest, i)
		}
	}
	return over, rest
}

// dedupNodes sorts ns by ID and drops repeats.
func dedupNodes(ns []model.Node) []model.Node {
	slices.SortFunc(ns, nodeByID)
	return slices.CompactFunc(ns, func(a, b model.Node) bool { return a.ID == b.ID })
}

// NodesLight returns the nodes with the given IDs, live or retired, sorted
// by ID, each once, without Fields and History.
func (o *Overlay) NodesLight(ids []id.ID) ([]model.Node, error) {
	over, rest := o.splitIDs(ids)
	ns, err := o.x.NodesLight(AsIs, rest)
	if err != nil {
		return nil, err
	}
	for _, n := range over {
		n.Fields, n.History = nil, nil
		ns = append(ns, n)
	}
	return dedupNodes(ns), nil
}

// NodesFields is NodesLight with each node's Fields.
func (o *Overlay) NodesFields(ids []id.ID) ([]model.Node, error) {
	over, rest := o.splitIDs(ids)
	ns, err := o.x.NodesFields(AsIs, rest)
	if err != nil {
		return nil, err
	}
	for _, n := range over {
		n.History = nil
		ns = append(ns, n)
	}
	return dedupNodes(ns), nil
}

// Parents returns the ID and parent of each node among ids, live or
// retired, sorted by ID, each once.
func (o *Overlay) Parents(ids []id.ID) ([]Parent, error) {
	over, rest := o.splitIDs(ids)
	ps, err := o.x.Parents(AsIs, rest)
	if err != nil {
		return nil, err
	}
	for _, n := range over {
		ps = append(ps, Parent{ID: n.ID, Parent: n.Parent})
	}
	slices.SortFunc(ps, func(a, b Parent) int { return id.Compare(a.ID, b.ID) })
	return slices.CompactFunc(ps, func(a, b Parent) bool { return a.ID == b.ID }), nil
}

// List returns the nodes q selects, sorted by ID, at most q.Limit, as
// Index.List does.
func (o *Overlay) List(q ListQuery) ([]model.Node, error) {
	if q.Limit <= 0 {
		return nil, fmt.Errorf("index: list limit %d is not positive", q.Limit)
	}
	over := o.pickNodes(func(n model.Node) bool {
		return (q.After.IsZero() || id.Compare(n.ID, q.After) > 0) &&
			(!q.HasParent || n.Parent == q.Parent) &&
			(len(q.Kinds) == 0 || slices.Contains(q.Kinds, n.Kind)) &&
			(q.Retired || !n.Retired)
	})
	// Every overridden row may sit among the stored ones read, so ask for
	// as many more as there are.
	stored := q
	stored.Limit = q.Limit + len(o.nodes)
	ns, err := o.x.List(AsIs, stored)
	if err != nil {
		return nil, err
	}
	ns = append(o.dropNodes(ns), over...)
	slices.SortFunc(ns, nodeByID)
	if len(ns) > q.Limit {
		ns = ns[:q.Limit]
	}
	if len(ns) == 0 {
		return nil, nil
	}
	return ns, nil
}

// Search returns the live nodes that may match terms, of kinds when given,
// sorted by ID: the stored ones that match, and every overridden live node
// of the kinds, which the caller filters by its text.
func (o *Overlay) Search(terms []string, kinds []op.Kind) ([]model.Node, error) {
	ns, err := o.x.Search(AsIs, terms, kinds)
	if err != nil {
		return nil, err
	}
	ns = append(o.dropNodes(ns), o.pickNodes(func(n model.Node) bool {
		return !n.Retired && (len(kinds) == 0 || slices.Contains(kinds, n.Kind))
	})...)
	slices.SortFunc(ns, nodeByID)
	return ns, nil
}

// Subtree returns the nodes under root, depth first in pre-order, as
// Index.Subtree does.
func (o *Overlay) Subtree(root id.ID, depth int, retired bool) ([]TreeNode, error) {
	return o.subtree(root, depth, retired, false)
}

// SubtreeLight is Subtree with light nodes.
func (o *Overlay) SubtreeLight(root id.ID, depth int, retired bool) ([]TreeNode, error) {
	return o.subtree(root, depth, retired, true)
}

func (o *Overlay) subtree(root id.ID, depth int, retired, light bool) ([]TreeNode, error) {
	if depth <= 0 || depth > maxDepth {
		depth = maxDepth
	}
	var tops []model.Node
	if root.IsZero() {
		var err error
		if tops, err = o.Children(id.ID{}); err != nil {
			return nil, err
		}
	} else {
		n, ok, err := o.Node(root)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("%w: no node %s", ErrNotFound, root)
		}
		tops = []model.Node{n}
	}
	var out []TreeNode
	done := map[id.ID]bool{}
	var walkErr error
	var visit func(n model.Node, d int)
	visit = func(n model.Node, d int) {
		if walkErr != nil || done[n.ID] || n.Retired && !retired {
			return
		}
		done[n.ID] = true
		if light {
			n.Fields, n.History = nil, nil
		}
		out = append(out, TreeNode{Node: n, Depth: d})
		if d >= depth {
			return
		}
		kids, err := o.Children(n.ID)
		if err != nil {
			walkErr = err
			return
		}
		for _, c := range kids {
			visit(c, d+1)
		}
	}
	for _, n := range tops {
		visit(n, 0)
	}
	if walkErr != nil {
		return nil, walkErr
	}
	return out, nil
}
