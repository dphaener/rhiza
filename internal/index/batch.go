package index

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/model"
	"github.com/dphaener/genesis/internal/op"
)

// This file holds the batch queries: one query for many IDs where a reader
// would otherwise ask once per ID, for reads such as the plan's that reach
// hundreds of nodes. Each splits its IDs into chunks of batchSize, well
// under SQLite's limit on bound parameters.

// batchSize is the most IDs one batch query binds.
const batchSize = 500

// lightNodeSelect is nodeSelect without fields and history, which scanNode
// reads back as empty.
const lightNodeSelect = `SELECT id, kind, handle, parent, key, NULL, retired, reason, version, created, NULL FROM nodes`

// NodesLight returns the nodes of view v with the given IDs, live or
// retired, sorted by ID, each once; IDs the view lacks are left out, and
// nil comes back when it has none of them. The nodes are light: every
// field of model.Node but Fields and History, which are nil. It is one
// primary-key query per chunk of IDs.
func (x *Index) NodesLight(v View, ids []id.ID) ([]model.Node, error) {
	if err := v.check(); err != nil {
		return nil, err
	}
	var out []model.Node
	for _, chunk := range batchChunks(ids) {
		ns, err := x.nodes(lightNodeSelect+` WHERE view = ? AND id IN `+sqlMarks(len(chunk)), batchArgs(v, chunk)...)
		if err != nil {
			return nil, err
		}
		out = append(out, ns...)
	}
	slices.SortFunc(out, func(a, b model.Node) int { return id.Compare(a.ID, b.ID) })
	return out, nil
}

// LinksTo returns the links of view v of type t whose To end is one of
// ids, live and retired, sorted by ID, each once; nil when there are none.
// It reads the to index once per chunk of IDs.
func (x *Index) LinksTo(v View, ids []id.ID, t op.LinkType) ([]model.Link, error) {
	if err := v.check(); err != nil {
		return nil, err
	}
	var out []model.Link
	for _, chunk := range batchChunks(ids) {
		ls, err := x.links(linkSelect+` INDEXED BY links_to WHERE view = ? AND to_id IN `+sqlMarks(len(chunk))+` AND type = ?`,
			append(batchArgs(v, chunk), string(t))...)
		if err != nil {
			return nil, err
		}
		out = append(out, ls...)
	}
	slices.SortFunc(out, func(a, b model.Link) int { return id.Compare(a.ID, b.ID) })
	return out, nil
}

// LinksAt returns the links of view v at any of ids in direction dir, as
// Links gives them for one node: live and retired, sorted by ID, each
// once; nil when there are none. ids may hold change sets.
func (x *Index) LinksAt(v View, ids []id.ID, dir Direction) ([]model.Link, error) {
	if err := v.check(); err != nil {
		return nil, err
	}
	if dir != From && dir != To && dir != Both {
		return nil, fmt.Errorf("index: unknown direction %d", dir)
	}
	var out []model.Link
	for _, chunk := range batchChunks(ids) {
		in := sqlMarks(len(chunk))
		var query string
		args := batchArgs(v, chunk)
		switch dir {
		case From:
			query = linkSelect + ` INDEXED BY links_from WHERE view = ? AND from_id IN ` + in
		case To:
			query = linkSelect + ` INDEXED BY links_to WHERE view = ? AND to_id IN ` + in
		case Both:
			query = linkSelect + ` INDEXED BY links_from WHERE view = ? AND from_id IN ` + in +
				` UNION ` + linkSelect + ` INDEXED BY links_to WHERE view = ? AND to_id IN ` + in
			args = append(args, args...)
		}
		ls, err := x.links(query, args...)
		if err != nil {
			return nil, err
		}
		out = append(out, ls...)
	}
	slices.SortFunc(out, func(a, b model.Link) int { return id.Compare(a.ID, b.ID) })
	return slices.CompactFunc(out, func(a, b model.Link) bool { return a.ID == b.ID }), nil
}

// EntriesOf returns the good entries of each of the change sets css, as
// Entries gives them for one: decoded from their stored bytes, sorted by
// ID. A change set with none has no key. The bytes are read with one query
// per chunk of change sets and decoded by decodeWorkers at once, since
// decoding is most of the cost.
func (x *Index) EntriesOf(css []id.ID) (map[id.ID][]entry.Entry, error) {
	type stored struct {
		cs  id.ID
		raw []byte
	}
	var all []stored
	for _, chunk := range batchChunks(css) {
		rows, err := x.db.Query(`SELECT change_set, raw FROM entries WHERE change_set IN `+sqlMarks(len(chunk))+
			` ORDER BY change_set, id`, idArgs(chunk)...)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var cs string
			var raw []byte
			if err := rows.Scan(&cs, &raw); err != nil {
				rows.Close()
				return nil, err
			}
			all = append(all, stored{parseID(cs), raw})
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			return nil, err
		}
	}
	es := make([]entry.Entry, len(all))
	errs := make([]error, len(all))
	next := make(chan int)
	var wg sync.WaitGroup
	for range min(decodeWorkers, len(all)) {
		wg.Go(func() {
			for i := range next {
				es[i], errs[i] = entry.Decode(all[i].raw)
			}
		})
	}
	for i := range all {
		next <- i
	}
	close(next)
	wg.Wait()
	out := map[id.ID][]entry.Entry{}
	for i, s := range all {
		if errs[i] != nil {
			return nil, fmt.Errorf("decode stored entry: %w", errs[i])
		}
		out[s.cs] = append(out[s.cs], es[i])
	}
	return out, nil
}

// decodeWorkers is how many entries EntriesOf decodes at once.
const decodeWorkers = 8

// batchChunks is ids sorted, without repeats, in chunks of at most
// batchSize; nil for no IDs.
func batchChunks(ids []id.ID) [][]id.ID {
	ids = slices.Clone(ids)
	slices.SortFunc(ids, id.Compare)
	return slices.Collect(slices.Chunk(slices.Compact(ids), batchSize))
}

// batchArgs is a batch query's arguments: the view, then each ID's text.
func batchArgs(v View, ids []id.ID) []any {
	return append([]any{string(v)}, idArgs(ids)...)
}

// idArgs is each ID's text, as query arguments.
func idArgs(ids []id.ID) []any {
	args := make([]any, len(ids))
	for k, i := range ids {
		args[k] = idText(i)
	}
	return args
}

// sqlMarks is an SQL list of n placeholders: "(?, ?, ?)".
func sqlMarks(n int) string {
	return "(" + strings.Repeat("?, ", n-1) + "?)"
}
