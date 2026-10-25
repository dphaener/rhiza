package index

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/model"
	"github.com/dphaener/genesis/internal/op"
)

// The queries below read the index as Open left it: they never catch up
// again, so an index answers for the log as it was listed when opened.

// ErrNotFound is wrapped by query errors when nothing matches: an unknown
// node or a ref that names none.
var ErrNotFound = errors.New("index: not found")

// ErrAmbiguous is wrapped by Resolve's error when a ref matches several
// nodes; the message lists their IDs, sorted.
var ErrAmbiguous = errors.New("index: ambiguous")

// View names a stored view of the model.
type View string

// The stored views.
const (
	// AsIs is the as-is view: the landed change sets (model.AsIs).
	AsIs View = viewAsIs
	// ToBe is the to-be view: as-is with the accepted change sets on top
	// (model.ToBe).
	ToBe View = viewToBe
)

// check returns an error unless v is a stored view.
func (v View) check() error {
	if v != AsIs && v != ToBe {
		return fmt.Errorf("index: unknown view %q", string(v))
	}
	return nil
}

// Direction picks which links of a node Links returns.
type Direction int

// The directions.
const (
	// From is the links whose From is the node.
	From Direction = iota
	// To is the links whose To is the node.
	To
	// Both is the links with the node at either end, each once.
	Both
)

// Index use: without ANALYZE statistics SQLite prefers walking a
// WITHOUT ROWID table's primary key over a secondary index that needs a
// lookup per row, so the handle, parent and link queries name their index
// with INDEXED BY (and sort in Go where ORDER BY id would compete). The
// schema is fixed by SchemaVersion, so the indexes always exist.

// minPrefix is the shortest ID prefix Resolve matches.
const minPrefix = 4

// treeKinds is an SQL list of the node kinds that sit in the tree.
var treeKinds = func() string {
	var ks []string
	for _, k := range op.Kinds() {
		if op.HasParent(k) {
			ks = append(ks, "'"+string(k)+"'")
		}
	}
	return "(" + strings.Join(ks, ", ") + ")"
}()

// Node returns the node with ID i in view v, live or retired.
func (x *Index) Node(v View, i id.ID) (model.Node, bool, error) {
	if err := v.check(); err != nil {
		return model.Node{}, false, err
	}
	ns, err := x.nodes(nodeSelect+` WHERE view = ? AND id = ?`, string(v), idText(i))
	if err != nil || len(ns) == 0 {
		return model.Node{}, false, err
	}
	return ns[0], true, nil
}

// Resolve finds a node of view v by ref, which is, in this order: the full
// canonical ID of a node, live or retired; a unique prefix of four
// characters or more of the node IDs, matched as given (IDs are upper case,
// so a lower-case ref never matches one), live or retired; the handle of a
// live node. No match is an error wrapping ErrNotFound; several matches at
// the step that matched one an error wrapping ErrAmbiguous that lists them.
func (x *Index) Resolve(v View, ref string) (model.Node, error) {
	if err := v.check(); err != nil {
		return model.Node{}, err
	}
	if i, err := id.Parse(ref); err == nil {
		n, ok, err := x.Node(v, i)
		if err != nil {
			return model.Node{}, err
		}
		if ok {
			return n, nil
		}
	}
	if len(ref) >= minPrefix {
		// IDs use only digits and upper-case letters, all below '~', so
		// [ref, ref+"~") holds exactly the IDs starting with ref, and the
		// range reads the primary key.
		ns, err := x.nodes(nodeSelect+` WHERE view = ? AND id >= ? AND id < ? ORDER BY id`,
			string(v), ref, ref+"~")
		if err != nil {
			return model.Node{}, err
		}
		if len(ns) > 0 {
			return one(ref, ns)
		}
	}
	ns, err := x.nodes(nodeSelect+` INDEXED BY nodes_handle WHERE view = ? AND handle = ? AND retired = 0`, string(v), ref)
	if err != nil {
		return model.Node{}, err
	}
	slices.SortFunc(ns, func(a, b model.Node) int { return id.Compare(a.ID, b.ID) })
	if len(ns) == 0 {
		return model.Node{}, fmt.Errorf("%w: no node matches %q", ErrNotFound, ref)
	}
	return one(ref, ns)
}

// one returns the only node of ns, or an ErrAmbiguous error listing them
// all; ns is sorted by ID.
func one(ref string, ns []model.Node) (model.Node, error) {
	if len(ns) == 1 {
		return ns[0], nil
	}
	names := make([]string, len(ns))
	for i, n := range ns {
		names[i] = n.ID.String() + " (" + n.Handle + ")"
	}
	return model.Node{}, fmt.Errorf("%w: %q matches %d nodes: %s", ErrAmbiguous, ref, len(ns), strings.Join(names, ", "))
}

// Children returns the nodes of view v directly under parent, live and
// retired, sorted by key then ID, as model.Model.Children does. The zero
// parent gives the roots: nodes of a kind that sits in the tree, with no
// parent. It is nil when there are none.
func (x *Index) Children(v View, parent id.ID) ([]model.Node, error) {
	if err := v.check(); err != nil {
		return nil, err
	}
	return x.nodes(nodeSelect+` INDEXED BY nodes_parent WHERE view = ? AND parent = ? AND kind IN `+treeKinds+` ORDER BY key, id`,
		string(v), idText(parent))
}

// Ancestors returns the parent chain of node i in view v, the root first,
// without i itself; nil for a root or a node outside the tree. An unknown i
// is an error wrapping ErrNotFound.
func (x *Index) Ancestors(v View, i id.ID) ([]model.Node, error) {
	if _, ok, err := x.Node(v, i); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("%w: no node %s", ErrNotFound, i)
	}
	// The depth bound ends even a cycle the model should never produce.
	// CROSS JOIN keeps up as the outer loop, so each step is a primary key
	// lookup.
	rows, err := x.nodes(`WITH RECURSIVE up(id, depth) AS (
			SELECT parent, 1 FROM nodes WHERE view = ?1 AND id = ?2 AND parent != ''
			UNION
			SELECT n.parent, up.depth + 1 FROM nodes n JOIN up ON n.view = ?1 AND n.id = up.id
			WHERE n.parent != '' AND up.depth < 1000
		)
		SELECT n.id, n.kind, n.handle, n.parent, n.key, n.fields, n.retired, n.reason, n.version, n.created, n.history
		FROM up CROSS JOIN nodes n ON n.view = ?1 AND n.id = up.id
		ORDER BY up.depth DESC`, string(v), idText(i))
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Links returns the links of view v at node in direction dir, live and
// retired, sorted by ID; nil when there are none. node may also be a
// change set, which links can name too.
func (x *Index) Links(v View, node id.ID, dir Direction) ([]model.Link, error) {
	if err := v.check(); err != nil {
		return nil, err
	}
	var query string
	args := []any{string(v), idText(node)}
	switch dir {
	case From:
		query = linkSelect + ` INDEXED BY links_from WHERE view = ? AND from_id = ?`
	case To:
		query = linkSelect + ` INDEXED BY links_to WHERE view = ? AND to_id = ?`
	case Both:
		query = linkSelect + ` INDEXED BY links_from WHERE view = ?1 AND from_id = ?2
			UNION ` + linkSelect + ` INDEXED BY links_to WHERE view = ?1 AND to_id = ?2`
	default:
		return nil, fmt.Errorf("index: unknown direction %d", dir)
	}
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
	slices.SortFunc(out, func(a, b model.Link) int { return id.Compare(a.ID, b.ID) })
	return out, rows.Err()
}

// nodes runs a query selecting nodeSelect's columns and scans every row;
// nil when there are none.
func (x *Index) nodes(query string, args ...any) ([]model.Node, error) {
	rows, err := x.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Node
	for rows.Next() {
		n, err := scanNode(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

// ChangeSetInfo is one folded change set as the index stores it.
type ChangeSetInfo struct {
	ID    id.ID
	State changeset.State
	// Kind and Intent come from the change set's open entry.
	Kind   changeset.Kind
	Intent string
	// OpenEntry is the open entry, SupersededBy the change set named by a
	// supersede and TerminalEntry the entry that ended it; each zero when
	// there is none.
	OpenEntry     id.ID
	SupersededBy  id.ID
	TerminalEntry id.ID
}

// ChangeSets returns every folded change set, sorted by ID; nil when there
// are none.
func (x *Index) ChangeSets() ([]ChangeSetInfo, error) {
	rows, err := x.db.Query(`SELECT id, state, kind, intent, open_entry, superseded_by, terminal_entry
		FROM change_sets ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ChangeSetInfo
	for rows.Next() {
		var c ChangeSetInfo
		var cid, state, kind, open, by, term string
		if err := rows.Scan(&cid, &state, &kind, &c.Intent, &open, &by, &term); err != nil {
			return nil, err
		}
		c.ID, c.State, c.Kind = parseID(cid), changeset.State(state), changeset.Kind(kind)
		c.OpenEntry, c.SupersededBy, c.TerminalEntry = parseID(open), parseID(by), parseID(term)
		out = append(out, c)
	}
	return out, rows.Err()
}

// Entries returns the good entries of change set cs, decoded from their
// stored bytes and sorted by ID; nil when there are none.
func (x *Index) Entries(cs id.ID) ([]entry.Entry, error) {
	rows, err := x.db.Query(`SELECT raw FROM entries WHERE change_set = ? ORDER BY id`, idText(cs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entry.Entry
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		e, err := entry.Decode(raw)
		if err != nil {
			return nil, fmt.Errorf("decode stored entry: %w", err)
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// Problem is one stored problem of the log, of folding or of a view.
type Problem struct {
	// Source is where it came from: "log" (layout and entry files), "fold"
	// (changeset.FoldLog), "asis" or "tobe" (model.AsIs, model.ToBe).
	Source string
	// Path is the file or folder under .genesis, for log problems.
	Path string
	// ChangeSet is set for fold and view problems, Entry for fold problems
	// and Op for view problems; each zero otherwise.
	ChangeSet id.ID
	Entry     id.ID
	Op        id.ID
	// Sentinel is the short name of the sentinel error the problem wraps
	// (for example "invalid", "order" or "conflict"), "" for none.
	Sentinel string
	// Message is the error's text.
	Message string
}

// Problems returns every stored problem: log, then fold, then as-is, then
// to-be, each source in the order its package sorts them; nil when there
// are none.
func (x *Index) Problems() ([]Problem, error) {
	rows, err := x.db.Query(`SELECT source, path, change_set, entry, op, sentinel, message FROM problems
		ORDER BY CASE source WHEN ? THEN 0 WHEN ? THEN 1 WHEN ? THEN 2 WHEN ? THEN 3 ELSE 4 END, source, seq`,
		sourceLog, sourceFold, viewAsIs, viewToBe)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Problem
	for rows.Next() {
		var p Problem
		var cs, e, o string
		if err := rows.Scan(&p.Source, &p.Path, &cs, &e, &o, &p.Sentinel, &p.Message); err != nil {
			return nil, err
		}
		p.ChangeSet, p.Entry, p.Op = parseID(cs), parseID(e), parseID(o)
		out = append(out, p)
	}
	return out, rows.Err()
}

// WhatIf returns the what-if view of change set cs, as model.WhatIf
// computes it over the change sets folded from every stored entry: as-is
// with cs on top. Problems are model.WhatIf's (ErrMissing when cs is not
// in the log, ErrState when it is not in flight).
func (x *Index) WhatIf(cs id.ID) (model.Model, []model.Problem, error) {
	sets, _, err := foldStored(x.db)
	if err != nil {
		return model.Model{}, nil, err
	}
	m, ps := model.WhatIf(sets, cs)
	return m, ps, nil
}
