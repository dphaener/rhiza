package index

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/model"
	"github.com/dphaener/genesis/internal/op"
)

// ListQuery selects the nodes List returns.
type ListQuery struct {
	// Kinds keeps the nodes of these kinds; all kinds when empty.
	Kinds []op.Kind
	// Parent, when HasParent, keeps the nodes whose parent is Parent. The
	// zero Parent keeps every node without one: the roots and the nodes of
	// kinds outside the tree.
	Parent    id.ID
	HasParent bool
	// Retired keeps retired nodes too; otherwise only live ones.
	Retired bool
	// After keeps the nodes whose ID is above it; the zero After keeps all.
	// A page's last ID is the next page's After.
	After id.ID
	// Limit is the most nodes returned; it must be positive.
	Limit int
}

// List returns the nodes of view v that q selects, sorted by ID, at most
// q.Limit of them; nil when there are none. A Limit of zero or less is an
// error. A page reads the primary key from After on, or with HasParent the
// parent index, so it costs about its own size, not the view's.
func (x *Index) List(v View, q ListQuery) ([]model.Node, error) {
	if err := v.check(); err != nil {
		return nil, err
	}
	if q.Limit <= 0 {
		return nil, fmt.Errorf("index: list limit %d is not positive", q.Limit)
	}
	query := nodeSelect + ` WHERE view = ? AND id > ?`
	args := []any{string(v), idText(q.After)}
	if q.HasParent {
		query = nodeSelect + ` INDEXED BY nodes_parent WHERE view = ? AND parent = ? AND id > ?`
		args = []any{string(v), idText(q.Parent), idText(q.After)}
	}
	if len(q.Kinds) > 0 {
		marks := make([]string, len(q.Kinds))
		for i, k := range q.Kinds {
			marks[i] = "?"
			args = append(args, string(k))
		}
		query += ` AND kind IN (` + strings.Join(marks, ", ") + `)`
	}
	if !q.Retired {
		query += ` AND retired = 0`
	}
	args = append(args, q.Limit)
	return x.nodes(query+` ORDER BY id LIMIT ?`, args...)
}

// TreeNode is one node of a Subtree with its depth: 0 for the subtree's
// root, or for each root of the tree, and one more per level down.
type TreeNode struct {
	Node  model.Node
	Depth int
}

// maxDepth bounds every walk of the tree, as Ancestors is bounded.
const maxDepth = 1000

// Subtree returns the nodes of view v under root, depth first in pre-order,
// siblings sorted by key then ID as Children sorts them; nil when there are
// none. The zero root walks every root of the tree, each at depth 0;
// otherwise root itself is at depth 0, even when its kind is outside the
// tree (it then has nothing under it). depth is the deepest Depth returned,
// so depth 1 is the roots and their children; zero or less means no limit.
// Retired nodes, and everything under them however live, are left out
// unless retired, a retired root included. An unknown root is an error
// wrapping ErrNotFound. It is one recursive query down the parent index.
func (x *Index) Subtree(v View, root id.ID, depth int, retired bool) ([]TreeNode, error) {
	return x.subtree(v, root, depth, retired, false)
}

// SubtreeLight is Subtree with light nodes (see NodesLight): the same nodes
// in the same order, without the fields and history a reader that only
// names them would pay to read and decode.
func (x *Index) SubtreeLight(v View, root id.ID, depth int, retired bool) ([]TreeNode, error) {
	return x.subtree(v, root, depth, retired, true)
}

// subtree is Subtree, with light nodes if light.
func (x *Index) subtree(v View, root id.ID, depth int, retired, light bool) ([]TreeNode, error) {
	if err := v.check(); err != nil {
		return nil, err
	}
	if depth <= 0 || depth > maxDepth {
		depth = maxDepth
	}
	live := ""
	if !retired {
		live = ` AND n.retired = 0`
	}
	// The named root is read live or retired, so that a missing row means
	// an unknown node; the walk goes on only from rows that pass.
	cols := `n.id, n.kind, n.handle, n.parent, n.key, n.fields, n.retired, n.reason, n.version, n.created, n.history`
	if light {
		cols = `n.id, n.kind, n.handle, n.parent, n.key, NULL, n.retired, n.reason, n.version, n.created, NULL`
	}
	start := `SELECT ` + cols + `, 0 FROM nodes n WHERE n.view = ?1 AND n.id = ?2`
	if root.IsZero() {
		start = `SELECT ` + cols + `, 0
			FROM nodes n INDEXED BY nodes_parent WHERE n.view = ?1 AND n.parent = '' AND n.kind IN ` + treeKinds + live
	}
	walk := ` AND t.retired = 0`
	if retired {
		walk = ""
	}
	// CROSS JOIN keeps t as the outer loop, so each step is a lookup in the
	// parent index.
	rows, err := x.db.Query(`WITH RECURSIVE t(id, kind, handle, parent, key, fields, retired, reason, version, created, history, depth) AS (
			`+start+`
			UNION ALL
			SELECT `+cols+`, t.depth + 1
			FROM t CROSS JOIN nodes n INDEXED BY nodes_parent ON n.view = ?1 AND n.parent = t.id
			WHERE n.kind IN `+treeKinds+live+walk+` AND t.depth < ?3
		)
		SELECT id, kind, handle, parent, key, fields, retired, reason, version, created, history FROM t`,
		string(v), idText(root), depth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	nodes := map[id.ID]model.Node{}
	under := map[id.ID][]id.ID{}
	var tops []model.Node
	for rows.Next() {
		n, err := scanNode(rows)
		if err != nil {
			return nil, err
		}
		if _, seen := nodes[n.ID]; seen {
			continue
		}
		nodes[n.ID] = n
		if root.IsZero() && n.Parent.IsZero() || n.ID == root {
			tops = append(tops, n)
		} else {
			under[n.Parent] = append(under[n.Parent], n.ID)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if !root.IsZero() && len(tops) == 0 {
		return nil, fmt.Errorf("%w: no node %s", ErrNotFound, root)
	}
	if !retired && !root.IsZero() && tops[0].Retired {
		return nil, nil
	}
	byKey := func(a, b model.Node) int {
		return cmp.Or(strings.Compare(a.Key, b.Key), id.Compare(a.ID, b.ID))
	}
	slices.SortFunc(tops, byKey)
	// done guards the walk against a cycle the model should never produce.
	var out []TreeNode
	done := map[id.ID]bool{}
	var visit func(n model.Node, d int)
	visit = func(n model.Node, d int) {
		if done[n.ID] {
			return
		}
		done[n.ID] = true
		out = append(out, TreeNode{Node: n, Depth: d})
		kids := make([]model.Node, 0, len(under[n.ID]))
		for _, c := range under[n.ID] {
			kids = append(kids, nodes[c])
		}
		slices.SortFunc(kids, byKey)
		for _, c := range kids {
			visit(c, d+1)
		}
	}
	for _, n := range tops {
		visit(n, 0)
	}
	return out, nil
}

// ResolveLink finds a link of view v by ref, which is the full canonical ID
// of a link or a unique prefix of four characters or more of the link IDs,
// matched as given as Resolve matches node IDs; live or retired either way.
// No match is an error wrapping ErrNotFound, several an error wrapping
// ErrAmbiguous that lists them.
func (x *Index) ResolveLink(v View, ref string) (model.Link, error) {
	if err := v.check(); err != nil {
		return model.Link{}, err
	}
	if i, err := id.Parse(ref); err == nil {
		l, ok, err := x.Link(v, i)
		if err != nil {
			return model.Link{}, err
		}
		if ok {
			return l, nil
		}
	}
	var ls []model.Link
	if len(ref) >= minPrefix {
		// The range trick of Resolve: [ref, ref+"~") is the IDs starting
		// with ref, read from the primary key.
		var err error
		ls, err = x.links(linkSelect+` WHERE view = ? AND id >= ? AND id < ? ORDER BY id`,
			string(v), ref, ref+"~")
		if err != nil {
			return model.Link{}, err
		}
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

// links runs a query selecting linkSelect's columns and scans every row;
// nil when there are none.
func (x *Index) links(query string, args ...any) ([]model.Link, error) {
	rows, err := x.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Link
	for rows.Next() {
		l, err := scanLink(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}
