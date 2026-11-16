package index

import (
	"bytes"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/store"
)

// HoldsLog reports whether the index stores exactly the good entries of
// l: no entry missing, none extra, none different. Catch-up never re-reads
// a path it has seen, so an entry file edited in place after it was
// indexed leaves the index answering from the old bytes; a caller that has
// read the whole log anyway (render) uses this to find out and fold the
// log instead. It writes nothing.
//
// Both sides are canonical (Decode accepts only canonical bytes, and the
// index stores what Decode accepted), so equal fields mean equal bytes:
// each stored entry is decoded on DecodeCanonical's fast path and its
// fields compared with l's entry, on decodeWorkers at once, instead of
// re-encoding every entry of l. The stored rows are read in key order,
// which is l's order (change set, then ID).
func (x *Index) HoldsLog(l store.Log) (bool, error) {
	var n int
	if err := x.db.QueryRow(`SELECT count(*) FROM entries`).Scan(&n); err != nil {
		return false, err
	}
	if n != len(l.Entries) {
		return false, nil
	}
	rows, err := x.db.Query(`SELECT raw FROM entries ORDER BY change_set, id`)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	raws := make([][]byte, 0, n)
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return false, err
		}
		raws = append(raws, raw)
	}
	if err := rows.Err(); err != nil {
		return false, err
	}
	if len(raws) != len(l.Entries) {
		return false, nil
	}

	var differs atomic.Bool
	var next atomic.Int64
	var wg sync.WaitGroup
	for range min(decodeWorkers, len(raws)) {
		wg.Go(func() {
			for k := int(next.Add(1)) - 1; k < len(raws) && !differs.Load(); k = int(next.Add(1)) - 1 {
				stored, err := entry.DecodeCanonical(raws[k])
				if err != nil || !sameEntry(stored, l.Entries[k]) {
					differs.Store(true)
				}
			}
		})
	}
	wg.Wait()
	return !differs.Load(), nil
}

// sameEntry reports whether a and b hold the same fields, and so encode
// to the same canonical bytes.
func sameEntry(a, b entry.Entry) bool {
	return a.V == b.V && a.ID == b.ID && a.ChangeSet == b.ChangeSet &&
		a.Actor == b.Actor && a.For == b.For && a.At.Equal(b.At) &&
		slices.Equal(a.After, b.After) && a.Type == b.Type && bytes.Equal(a.Body, b.Body)
}
