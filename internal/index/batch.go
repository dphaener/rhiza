package index

import (
	"fmt"
	"runtime"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

// This file holds the batch queries: one query for many IDs where a reader
// would otherwise ask once per ID, for reads such as the plan's that reach
// hundreds of nodes. Each splits its IDs into chunks of batchSize and binds
// a chunk as one JSON array, which json_each walks as the outer loop of a
// join, a primary-key or index lookup per ID: cheaper than an IN list of
// one bound parameter per ID, which SQLite copies and sorts first.

// batchSize is the most IDs one batch query asks for.
const batchSize = 500

// batchIDs is the join that the batch queries read their IDs through: the
// chunk's JSON array, bound as ?2, then the table, aliased t; the view, if
// any, is ?1.
const batchIDs = ` FROM json_each(?2) j CROSS JOIN `

// lightNodeSelect is nodeSelect without fields and history, which scanNode
// reads back as empty, over the IDs of a batch.
const lightNodeSelect = `SELECT t.id, t.kind, t.handle, t.parent, t.key, NULL, t.retired, t.reason, t.version, t.created, NULL` +
	batchIDs + `nodes t ON t.view = ?1 AND t.id = j.value`

// NodesLight returns the nodes of view v with the given IDs, live or
// retired, sorted by ID, each once; IDs the view lacks are left out, and
// nil comes back when it has none of them. The nodes are light: every
// field of model.Node but Fields and History, which are nil. It is one
// primary-key query per chunk of IDs.
func (x *Index) NodesLight(v View, ids []id.ID) ([]model.Node, error) {
	if err := v.check(); err != nil {
		return nil, err
	}
	out, err := batchRun(ids, func(chunk []id.ID) ([]model.Node, error) {
		return x.nodes(lightNodeSelect, string(v), idJSON(chunk))
	})
	if err != nil {
		return nil, err
	}
	slices.SortFunc(out, func(a, b model.Node) int { return id.Compare(a.ID, b.ID) })
	return out, nil
}

// Parent is a node's ID and parent as Parents gives them; Parent is zero
// for a root or a node outside the tree.
type Parent struct {
	ID, Parent id.ID
}

// Parents returns the ID and parent of each node of view v among ids,
// live or retired, sorted by ID, each once; IDs the view lacks are left
// out, and nil comes back when it has none of them. It is NodesLight for a
// walk up the tree that needs only each node's parent: two columns read
// per row instead of eleven, one primary-key query per chunk of IDs.
func (x *Index) Parents(v View, ids []id.ID) ([]Parent, error) {
	if err := v.check(); err != nil {
		return nil, err
	}
	out, err := batchRun(ids, func(chunk []id.ID) ([]Parent, error) {
		rows, err := x.db.Query(`SELECT t.id, t.parent`+batchIDs+`nodes t ON t.view = ?1 AND t.id = j.value`,
			string(v), idJSON(chunk))
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var ps []Parent
		for rows.Next() {
			var nid, parent string
			if err := rows.Scan(&nid, &parent); err != nil {
				return nil, err
			}
			ps = append(ps, Parent{ID: parseID(nid), Parent: parseID(parent)})
		}
		return ps, rows.Err()
	})
	if err != nil {
		return nil, err
	}
	slices.SortFunc(out, func(a, b Parent) int { return id.Compare(a.ID, b.ID) })
	return out, nil
}

// fieldsNodeSelect is lightNodeSelect with the fields.
const fieldsNodeSelect = `SELECT t.id, t.kind, t.handle, t.parent, t.key, t.fields, t.retired, t.reason, t.version, t.created, NULL` +
	batchIDs + `nodes t ON t.view = ?1 AND t.id = j.value`

// NodesFields is NodesLight with each node's Fields; History is nil. It
// is one primary-key query per chunk of IDs.
func (x *Index) NodesFields(v View, ids []id.ID) ([]model.Node, error) {
	if err := v.check(); err != nil {
		return nil, err
	}
	out, err := batchRun(ids, func(chunk []id.ID) ([]model.Node, error) {
		return x.nodes(fieldsNodeSelect, string(v), idJSON(chunk))
	})
	if err != nil {
		return nil, err
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
	out, err := batchRun(ids, func(chunk []id.ID) ([]model.Link, error) {
		return x.links(batchLinkSelect+`links t INDEXED BY links_to ON t.view = ?1 AND t.to_id = j.value WHERE t.type = ?3`,
			string(v), idJSON(chunk), string(t))
	})
	if err != nil {
		return nil, err
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
	const (
		from = batchLinkSelect + `links t INDEXED BY links_from ON t.view = ?1 AND t.from_id = j.value`
		to   = batchLinkSelect + `links t INDEXED BY links_to ON t.view = ?1 AND t.to_id = j.value`
	)
	var query string
	switch dir {
	case From:
		query = from
	case To:
		query = to
	case Both:
		query = from + ` UNION ` + to
	}
	out, err := batchRun(ids, func(chunk []id.ID) ([]model.Link, error) {
		return x.links(query, string(v), idJSON(chunk))
	})
	if err != nil {
		return nil, err
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
	var all []storedEntry
	for _, chunk := range batchChunks(css) {
		var err error
		if all, err = x.storedEntries(all, `SELECT t.change_set, t.raw FROM json_each(?1) j CROSS JOIN entries t ON t.change_set = j.value
			ORDER BY t.change_set, t.id`, idJSON(chunk)); err != nil {
			return nil, err
		}
	}
	return decodeStored(all)
}

// EachInFlight calls each with the good entries of every change set whose
// state is draft, proposed or accepted, as ChangeSets gives it, decoded
// from their stored bytes and sorted by ID, as EntriesOf gives them; a
// change set with none is left out. One query finds the change sets and
// reads their entries, so a reader need not list them first. The rows come
// a change set at a time, and decodeWorkers decode each change set's and
// call each while the rest are still read, so decoding, and whatever each
// does, overlaps the query: each is called from several goroutines at once.
// EachInFlight returns once the last call has returned. Its error is the
// query's, or else the first entry's, in the order read, that does not
// decode; each may have been called for other change sets by then.
func (x *Index) EachInFlight(each func(cs id.ID, es []entry.Entry)) error {
	rows, err := x.db.Query(`SELECT t.change_set, t.raw FROM change_sets c CROSS JOIN entries t ON t.change_set = c.id
		WHERE c.state IN (?, ?, ?) ORDER BY c.id, t.id`,
		string(changeset.StateDraft), string(changeset.StateProposed), string(changeset.StateAccepted))
	if err != nil {
		return err
	}
	defer rows.Close()
	// group is one change set's rows: the place of its first in the order
	// read, and their bytes.
	type group struct {
		first int
		cs    id.ID
		raws  [][]byte
	}
	var (
		// bad is the place of the first entry that failed to decode, -1
		// for none, and badErr its error.
		mu     sync.Mutex
		bad    = -1
		badErr error
		wg     sync.WaitGroup
		work   = make(chan group, decodeWorkers)
		// cur is the change set being read, curText its ID's text, and
		// read the rows read so far.
		cur     group
		curText string
		read    int
		cs      string
		raw     []byte
		scanErr error
	)
	for range decodeWorkers {
		wg.Go(func() {
			for g := range work {
				es := make([]entry.Entry, len(g.raws))
				for k, raw := range g.raws {
					var err error
					if es[k], err = entry.DecodeCanonical(raw); err != nil {
						mu.Lock()
						if bad < 0 || g.first+k < bad {
							bad, badErr = g.first+k, err
						}
						mu.Unlock()
						es = nil
						break
					}
				}
				if es != nil {
					each(g.cs, es)
				}
			}
		})
	}
	for rows.Next() {
		if scanErr = rows.Scan(&cs, &raw); scanErr != nil {
			break
		}
		if len(cur.raws) == 0 || cs != curText {
			if len(cur.raws) > 0 {
				work <- cur
			}
			cur, curText = group{first: read, cs: parseID(cs)}, cs
		}
		cur.raws = append(cur.raws, raw)
		raw = nil
		read++
	}
	if len(cur.raws) > 0 && scanErr == nil {
		work <- cur
	}
	close(work)
	wg.Wait()
	if scanErr != nil {
		return scanErr
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if badErr != nil {
		return fmt.Errorf("decode stored entry: %w", badErr)
	}
	return nil
}

// storedEntry is one stored entry's change set and bytes.
type storedEntry struct {
	cs  id.ID
	raw []byte
}

// storedEntries appends to all the rows of query, which selects an entry's
// change set and raw bytes.
func (x *Index) storedEntries(all []storedEntry, query string, args ...any) ([]storedEntry, error) {
	rows, err := x.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var cs string
		var raw []byte
		if err := rows.Scan(&cs, &raw); err != nil {
			return nil, err
		}
		all = append(all, storedEntry{parseID(cs), raw})
	}
	return all, rows.Err()
}

// decodeStored decodes all by decodeWorkers at once, grouped by change
// set in the order of all.
func decodeStored(all []storedEntry) (map[id.ID][]entry.Entry, error) {
	es := make([]entry.Entry, len(all))
	errs := make([]error, len(all))
	next := make(chan int)
	var wg sync.WaitGroup
	for range min(decodeWorkers, len(all)) {
		wg.Go(func() {
			for i := range next {
				es[i], errs[i] = entry.DecodeCanonical(all[i].raw)
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

// batchParts is the most connections a batch query's chunks are read on
// at once, and batchMinPart the fewest IDs worth a part of their own. Each
// lookup is SQLite CPU work, so two parts on two connections take about
// half as long. database/sql keeps two idle connections, so a second one
// is opened once per Index and then reused, where a third would be opened
// again for every batch and cost more than it saves on two CPUs.
const (
	batchParts   = 2
	batchMinPart = 50
)

// BatchOne is the most IDs a batch query reads on one connection: fewer
// than two parts' worth. A reader that holds the other connection already,
// as the plan does while its entries are read, asks for at most this many
// at once, so that the index opens no third.
const BatchOne = 2*batchMinPart - 1

// batchRun runs query over ids sorted, without repeats, in chunks of at
// most batchSize, as many as batchParts at once when there are enough IDs
// and CPUs, and returns what the chunks return, in chunk order. Its error
// is the first chunk's that failed.
func batchRun[T any](ids []id.ID, query func([]id.ID) ([]T, error)) ([]T, error) {
	ids = slices.Clone(ids)
	slices.SortFunc(ids, id.Compare)
	ids = slices.Compact(ids)
	parts := max(1, min(batchParts, runtime.GOMAXPROCS(0), len(ids)/batchMinPart))
	size := max(1, min(batchSize, (len(ids)+parts-1)/parts))
	chunks := slices.Collect(slices.Chunk(ids, size))
	res := make([][]T, len(chunks))
	errs := make([]error, len(chunks))
	run := func(next *atomic.Int64) {
		for k := int(next.Add(1)) - 1; k < len(chunks); k = int(next.Add(1)) - 1 {
			res[k], errs[k] = query(chunks[k])
		}
	}
	var next atomic.Int64
	var wg sync.WaitGroup
	for range min(parts, len(chunks)) - 1 {
		wg.Go(func() { run(&next) })
	}
	run(&next)
	wg.Wait()
	var out []T
	for k := range chunks {
		if errs[k] != nil {
			return nil, errs[k]
		}
		out = append(out, res[k]...)
	}
	return out, nil
}

// batchChunks is ids sorted, without repeats, in chunks of at most
// batchSize; nil for no IDs.
func batchChunks(ids []id.ID) [][]id.ID {
	ids = slices.Clone(ids)
	slices.SortFunc(ids, id.Compare)
	return slices.Collect(slices.Chunk(slices.Compact(ids), batchSize))
}

// batchLinkSelect is linkSelect over the IDs of a batch; the join's table
// and condition follow.
const batchLinkSelect = `SELECT t.id, t.type, t.from_id, t.to_id, t.note, t.decision, t.retired, t.version, t.created, t.history` + batchIDs

// idJSON is ids' text as a JSON array, for json_each. ID text is digits
// and upper-case letters, which need no escaping.
func idJSON(ids []id.ID) string {
	var b strings.Builder
	b.Grow(len(ids)*(len(id.ID{}.String())+3) + 2)
	b.WriteByte('[')
	for k, i := range ids {
		if k > 0 {
			b.WriteByte(',')
		}
		b.WriteByte('"')
		b.WriteString(idText(i))
		b.WriteByte('"')
	}
	b.WriteByte(']')
	return b.String()
}
