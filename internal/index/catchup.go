package index

import (
	"bytes"
	"cmp"
	"crypto/sha256"
	"database/sql"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/store"
)

// metaListing is the meta key holding the listing fingerprint; metaStats
// the stat fingerprint taken before that listing, absent when none could
// be stored.
const (
	metaListing = "listing"
	metaStats   = "stats"
	// metaEntryFormat holds the entry format version the index was made
	// with; a rebuild keeps it.
	metaEntryFormat = "entry_format"
)

// Problem sources, as stored in problems.source; the view names double as
// nodes.view and links.view.
const (
	sourceLog  = "log"
	sourceFold = "fold"
	viewAsIs   = "asis"
	viewToBe   = "tobe"
)

// fingerprint sums a listing: every path and every layout problem. Equal
// fingerprints mean nothing new to read and no log problem changed.
func fingerprint(l store.Listing) string {
	h := sha256.New()
	for _, p := range l.Paths {
		io.WriteString(h, p)
		h.Write([]byte{'\n'})
	}
	h.Write([]byte{0})
	for _, p := range l.Problems {
		fmt.Fprintf(h, "%s\t%s\n", p.Path, p.Err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// catchUp brings the index level with the log, as the package doc
// describes, taking the stats now.
func (x *Index) catchUp() error {
	return x.catchUpFrom(takeStats(x.fsys, time.Now()))
}

// catchUpFrom is catchUp given stats taken before the call. When they
// equal the stored stat fingerprint the index is fresh and the log is not
// listed. Otherwise the log is listed, and when its fingerprint is the
// stored one the index is fresh too; the stats are then stored, if they
// may be, in a short write transaction. Only a listing that differs takes
// the catch-up transaction, which lists the log again under the lock.
func (x *Index) catchUpFrom(sp statPrint) error {
	x.stats = stats{writes: map[string]int{}}
	storedStats, err := metaOf(x.db, metaStats)
	if err != nil {
		return err
	}
	if sp.fp != "" && sp.fp == storedStats {
		x.stats.fresh, x.stats.statFresh = true, true
		return nil
	}
	fp := fingerprint(listLog(x.fsys))
	stored, err := metaOf(x.db, metaListing)
	if err != nil {
		return err
	}
	if stored == fp {
		x.stats.fresh = true
		if sp.storable() {
			x.storeStats(fp, sp.fp)
		}
		return nil
	}

	tx, err := x.db.Begin() // BEGIN IMMEDIATE, through the DSN
	if err != nil {
		return fmt.Errorf("catch up: %w", err)
	}
	defer tx.Rollback()
	if err := x.catchUpTx(tx, sp); err != nil {
		return fmt.Errorf("catch up: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("catch up: %w", err)
	}
	return nil
}

// storeStats stores statFP as the stat fingerprint of the listing whose
// fingerprint is fp, unless another process has caught up to another
// listing meanwhile. The index is fresh either way, so a failure here
// (the lock held past the busy timeout, say) fails no read: the next open
// lists the log and tries again.
func (x *Index) storeStats(fp, statFP string) {
	tx, err := x.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	if stored, err := metaOf(tx, metaListing); err != nil || stored != fp {
		return
	}
	if _, err := tx.Exec(`INSERT OR REPLACE INTO meta (key, value) VALUES (?, ?)`, metaStats, statFP); err != nil {
		return
	}
	if tx.Commit() == nil {
		x.stats.statStored = true
	}
}

// querier is what *sql.DB and *sql.Tx share.
type querier interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

// metaOf returns the stored value of a meta key, "" when there is none.
func metaOf(q querier, key string) (string, error) {
	var v string
	err := q.QueryRow(`SELECT value FROM meta WHERE key = ?`, key).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return v, err
}

// catchUpTx is catch-up once the write lock is held. sp was taken before
// the listing this lists, so a file added since changes the next stats.
func (x *Index) catchUpTx(tx *sql.Tx, sp statPrint) error {
	// Another process may have caught up while this one waited.
	listing := listLog(x.fsys)
	fp := fingerprint(listing)
	stored, err := metaOf(tx, metaListing)
	if err != nil {
		return err
	}
	if stored == fp {
		// Its stats, if it stored any, are newer than these.
		x.stats.fresh = true
		return nil
	}

	known, err := column(tx, `SELECT path FROM files`)
	if err != nil {
		return err
	}
	listed := make(map[string]bool, len(listing.Paths))
	for _, p := range listing.Paths {
		listed[p] = true
	}
	seen := make(map[string]bool, len(known))
	for _, p := range known {
		seen[p] = true
		if !listed[p] {
			x.stats.rebuilt = true
		}
	}
	if x.stats.rebuilt {
		for _, t := range []string{"meta", "files", "entries", "change_sets", "folded", "nodes", "links", "problems", "landed_links"} {
			q := `DELETE FROM ` + t
			if t == "meta" {
				q += ` WHERE key <> '` + metaEntryFormat + `'`
			}
			if _, err := tx.Exec(q); err != nil {
				return err
			}
		}
		clear(seen)
	}

	var fresh []string
	for _, p := range listing.Paths {
		if !seen[p] {
			fresh = append(fresh, p)
		}
	}
	read := store.ReadEntries(x.fsys, fresh)
	x.stats.read = len(fresh)
	touched, err := addFiles(tx, read)
	if err != nil {
		return err
	}

	if len(touched) > 0 || x.stats.rebuilt {
		if err := x.refold(tx, touched); err != nil {
			return err
		}
	}
	if err := writeLogProblems(tx, listing); err != nil {
		return err
	}
	if _, err := tx.Exec(`INSERT OR REPLACE INTO meta (key, value) VALUES (?, ?)`, metaListing, fp); err != nil {
		return err
	}
	if !sp.storable() {
		// Stats of an older listing must not outlive it.
		_, err := tx.Exec(`DELETE FROM meta WHERE key = ?`, metaStats)
		return err
	}
	if _, err := tx.Exec(`INSERT OR REPLACE INTO meta (key, value) VALUES (?, ?)`, metaStats, sp.fp); err != nil {
		return err
	}
	x.stats.statStored = true
	return nil
}

// column returns the one text column a query selects.
func column(q querier, query string, args ...any) ([]string, error) {
	rows, err := q.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// addFiles records the files read: each good entry with its canonical
// bytes, each bad file with its problem. It returns the change sets of the
// good entries.
func addFiles(tx *sql.Tx, read store.Log) (map[id.ID]bool, error) {
	file, err := tx.Prepare(`INSERT INTO files (path, change_set, entry_id, ok, problem, sentinel) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	ent, err := tx.Prepare(`INSERT INTO entries (change_set, id, raw) VALUES (?, ?, ?)`)
	if err != nil {
		return nil, err
	}
	defer ent.Close()

	touched := map[id.ID]bool{}
	for _, e := range read.Entries {
		raw, err := entry.Encode(e)
		if err != nil {
			return nil, fmt.Errorf("re-encode entry %s: %w", e.ID, err)
		}
		p := entryPath(e.ChangeSet, e.ID)
		if _, err := file.Exec(p, idText(e.ChangeSet), idText(e.ID), 1, "", ""); err != nil {
			return nil, err
		}
		if _, err := ent.Exec(idText(e.ChangeSet), idText(e.ID), raw); err != nil {
			return nil, err
		}
		touched[e.ChangeSet] = true
	}
	for _, p := range read.Problems {
		cs, eid := pathIDs(p.Path)
		if _, err := file.Exec(p.Path, cs, eid, 0, p.Err.Error(), logSentinel(p.Err)); err != nil {
			return nil, err
		}
	}
	return touched, nil
}

// entryPath is the slash path of an entry file under .rhiza.
func entryPath(cs, eid id.ID) string {
	return store.LogDir + "/" + cs.String() + "/" + eid.String() + ".json"
}

// pathIDs returns the change set and entry ID texts a log/<cs>/<id>.json
// path names, as ListLog lists them.
func pathIDs(p string) (cs, eid string) {
	parts := strings.Split(p, "/")
	if len(parts) != 3 {
		return "", ""
	}
	return parts[1], strings.TrimSuffix(parts[2], ".json")
}

// refold folds the change sets the new entries touched, from their own
// stored entries, and keeps every other change set's stored fold. It
// replaces the touched change sets' stored rows, their fold problems and
// the unknown-supersede problems the touch changed, and recomputes the
// views whose inputs changed. The result equals changeset.FoldLog over the
// whole log.
func (x *Index) refold(tx *sql.Tx, touched map[id.ID]bool) error {
	before := map[id.ID]changeset.State{}
	supersededBy := map[id.ID][]id.ID{}
	rows, err := tx.Query(`SELECT id, state, superseded_by FROM change_sets`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var s, state, by string
		if err := rows.Scan(&s, &state, &by); err != nil {
			rows.Close()
			return err
		}
		before[parseID(s)] = changeset.State(state)
		if by != "" {
			supersededBy[parseID(by)] = append(supersededBy[parseID(by)], parseID(s))
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	ids := make([]id.ID, 0, len(touched))
	for cs := range touched {
		ids = append(ids, cs)
	}
	slices.SortFunc(ids, id.Compare)

	var sets []changeset.ChangeSet
	var problems []changeset.Problem
	for _, cs := range ids {
		es, err := decodeEntriesOf(tx, cs, &x.stats)
		if err != nil {
			return err
		}
		c, ps := changeset.Fold(cs, es)
		x.stats.folded++
		sets = append(sets, c)
		problems = append(problems, ps...)
		if c.State == changeset.StateSuperseded && !touched[c.SupersededBy] {
			has, err := hasEntries(tx, c.SupersededBy)
			if err != nil {
				return err
			}
			if !has {
				problems = append(problems, unknownSupersede(c))
			}
		}
	}

	// Untouched change sets superseded by a touched one: the by now has
	// entries, so their unknown-supersede problem, if any, is gone.
	cleared := map[string]bool{}
	for cs := range touched {
		for _, s := range supersededBy[cs] {
			if !touched[s] {
				cleared[idText(s)] = true
			}
		}
	}
	if err := x.writeFoldedSets(tx, ids, sets); err != nil {
		return err
	}
	if err := replaceFoldProblems(tx, touched, cleared, problems); err != nil {
		return err
	}

	asIs, toBe := x.stats.rebuilt, x.stats.rebuilt
	for _, s := range sets {
		was, ok := before[s.ID]
		switch {
		case !ok:
			// A change set new to the index: any link may name it.
			asIs, toBe = true, true
		case was == changeset.StateLanded || s.State == changeset.StateLanded:
			asIs, toBe = true, true
		case was == changeset.StateAccepted || s.State == changeset.StateAccepted:
			toBe = true
		}
	}
	x.stats.asIs, x.stats.toBe = asIs, toBe
	if !asIs && !toBe {
		return nil
	}
	all, err := storedSets(tx)
	if err != nil {
		return err
	}
	return writeViews(tx, all, asIs, toBe, x.stats.writes)
}

// unknownSupersede is the problem FoldLog reports on a superseded change
// set whose by has no entries.
func unknownSupersede(c changeset.ChangeSet) changeset.Problem {
	return changeset.Problem{ChangeSet: c.ID, Entry: c.TerminalEntry,
		Err: fmt.Errorf("%w: superseded by %s, which has no entries in the log", changeset.ErrUnknownChangeSet, c.SupersededBy)}
}

// hasEntries reports whether cs has a good entry stored.
func hasEntries(q querier, cs id.ID) (bool, error) {
	var one int
	err := q.QueryRow(`SELECT 1 FROM entries WHERE change_set = ? LIMIT 1`, idText(cs)).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

// decodeEntriesOf decodes cs's stored entries.
func decodeEntriesOf(q querier, cs id.ID, st *stats) ([]entry.Entry, error) {
	rows, err := q.Query(`SELECT raw FROM entries WHERE change_set = ? ORDER BY id`, idText(cs))
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
		st.decoded++
		out = append(out, e)
	}
	return out, rows.Err()
}

// writeFoldedSets replaces the stored rows of the change sets ids with
// sets (the same change sets, in the same order): the summary row, the
// landed links and the folded state.
func (x *Index) writeFoldedSets(tx *sql.Tx, ids []id.ID, sets []changeset.ChangeSet) error {
	for _, t := range []string{"change_sets", "landed_links", "folded"} {
		col := "id"
		if t != "change_sets" {
			col = "change_set"
		}
		st, err := tx.Prepare(`DELETE FROM ` + t + ` WHERE ` + col + ` = ?`)
		if err != nil {
			return err
		}
		for _, cs := range ids {
			if _, err := st.Exec(idText(cs)); err != nil {
				st.Close()
				return err
			}
		}
		st.Close()
	}
	if err := writeChangeSets(tx, sets); err != nil {
		return err
	}
	st, err := tx.Prepare(`INSERT INTO folded (change_set, ops, data) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer st.Close()
	for _, s := range sets {
		var ops, data bytes.Buffer
		if err := gob.NewEncoder(&ops).Encode(opsOnly{ID: s.ID, State: s.State, Ops: s.Ops}); err != nil {
			return fmt.Errorf("encode folded %s: %w", s.ID, err)
		}
		if err := gob.NewEncoder(&data).Encode(s); err != nil {
			return fmt.Errorf("encode folded %s: %w", s.ID, err)
		}
		if _, err := st.Exec(idText(s.ID), ops.Bytes(), data.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

// decodeFolded decodes one stored folded change set.
func decodeFolded(data []byte) (changeset.ChangeSet, error) {
	var c changeset.ChangeSet
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&c); err != nil {
		return c, fmt.Errorf("decode folded change set: %w", err)
	}
	return c, nil
}

// storedSets returns every stored folded change set, sorted by ID.
func storedSets(q querier) ([]changeset.ChangeSet, error) {
	rows, err := q.Query(`SELECT data FROM folded ORDER BY change_set`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []changeset.ChangeSet
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}
		c, err := decodeFolded(data)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// FoldedOf returns the stored folded state, history included, of each of
// the change sets css; a change set with no entries has no key. The result
// equals what changeset.FoldLog gives for them over the whole log.
func (x *Index) FoldedOf(css []id.ID) (map[id.ID]changeset.ChangeSet, error) {
	return x.foldedOf(css, "data", decodeFolded)
}

// OpsOf is FoldedOf with only the ID, State and Ops of each change set
// set, read from their own stored copy: the rest of the stored state, the
// history above all, is not read, which makes it several times cheaper.
func (x *Index) OpsOf(css []id.ID) (map[id.ID]changeset.ChangeSet, error) {
	return x.foldedOf(css, "ops", decodeOps)
}

// opsOnly is the part of a folded change set stored as its ops.
type opsOnly struct {
	ID    id.ID
	State changeset.State
	Ops   []op.Op
}

// decodeOps decodes the stored ops of one folded change set.
func decodeOps(data []byte) (changeset.ChangeSet, error) {
	var c opsOnly
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&c); err != nil {
		return changeset.ChangeSet{}, fmt.Errorf("decode folded change set: %w", err)
	}
	return changeset.ChangeSet{ID: c.ID, State: c.State, Ops: c.Ops}, nil
}

// foldedOf reads column col of the stored folded change sets css, each
// through decode.
func (x *Index) foldedOf(css []id.ID, col string, decode func([]byte) (changeset.ChangeSet, error)) (map[id.ID]changeset.ChangeSet, error) {
	var blobs [][]byte
	for _, chunk := range batchChunks(css) {
		rows, err := x.db.Query(`SELECT f.`+col+` FROM json_each(?1) j CROSS JOIN folded f ON f.change_set = j.value`, idJSON(chunk))
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var data []byte
			if err := rows.Scan(&data); err != nil {
				rows.Close()
				return nil, err
			}
			blobs = append(blobs, data)
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			return nil, err
		}
	}
	// Decoding is most of the cost, so decodeWorkers do it at once.
	sets := make([]changeset.ChangeSet, len(blobs))
	errs := make([]error, len(blobs))
	next := make(chan int)
	var wg sync.WaitGroup
	for range min(decodeWorkers, len(blobs)) {
		wg.Go(func() {
			for i := range next {
				sets[i], errs[i] = decode(blobs[i])
			}
		})
	}
	for i := range blobs {
		next <- i
	}
	close(next)
	wg.Wait()
	out := make(map[id.ID]changeset.ChangeSet, len(blobs))
	for i, c := range sets {
		if errs[i] != nil {
			return nil, errs[i]
		}
		out[c.ID] = c
	}
	return out, nil
}

// replaceFoldProblems makes the stored fold problems those of the other
// change sets, minus the unknown-supersede problems of the change sets in
// cleared, plus fresh, in FoldLog's order.
func replaceFoldProblems(tx *sql.Tx, touched map[id.ID]bool, cleared map[string]bool, fresh []changeset.Problem) error {
	rows, err := tx.Query(`SELECT change_set, entry, sentinel, message FROM problems WHERE source = ? ORDER BY seq`, sourceFold)
	if err != nil {
		return err
	}
	var out []storedProblem
	for rows.Next() {
		var p storedProblem
		if err := rows.Scan(&p.ChangeSet, &p.Entry, &p.Sentinel, &p.Message); err != nil {
			rows.Close()
			return err
		}
		if touched[parseID(p.ChangeSet)] || (cleared[p.ChangeSet] && p.Sentinel == "unknown-change-set") {
			continue
		}
		out = append(out, p)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	for _, p := range fresh {
		out = append(out, storedProblem{ChangeSet: idText(p.ChangeSet), Entry: idText(p.Entry),
			Sentinel: foldSentinel(p.Err), Message: p.Err.Error()})
	}
	slices.SortStableFunc(out, func(a, b storedProblem) int {
		return cmp.Or(strings.Compare(a.ChangeSet, b.ChangeSet), strings.Compare(a.Entry, b.Entry),
			strings.Compare(a.Message, b.Message))
	})
	return replaceProblems(tx, sourceFold, out)
}

// writeChangeSets inserts the stored rows of sets.
func writeChangeSets(tx *sql.Tx, sets []changeset.ChangeSet) error {
	st, err := tx.Prepare(`INSERT INTO change_sets (id, state, kind, intent, open_entry, superseded_by, terminal_entry)
		VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer st.Close()
	for _, s := range sets {
		if _, err := st.Exec(idText(s.ID), string(s.State), string(s.Open.Kind), s.Open.Intent,
			idText(s.OpenEntry), idText(s.SupersededBy), idText(s.TerminalEntry)); err != nil {
			return err
		}
	}
	return writeLandedLinks(tx, sets)
}

// writeLandedLinks inserts the landed_links of sets.
func writeLandedLinks(tx *sql.Tx, sets []changeset.ChangeSet) error {
	st, err := tx.Prepare(`INSERT INTO landed_links (change_set, n, retire, link, from_id, to_id) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer st.Close()
	for _, s := range sets {
		if s.State != changeset.StateLanded {
			continue
		}
		for n, o := range s.Ops {
			switch {
			case o.Type == op.TypeLinkCreate && o.LinkType == op.LinkDependsOn:
				_, err = st.Exec(idText(s.ID), n, 0, idText(o.Link), idText(o.From), idText(o.To))
			case o.Type == op.TypeLinkRetire:
				_, err = st.Exec(idText(s.ID), n, 1, idText(o.Link), "", "")
			default:
				continue
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// writeLogProblems replaces the stored log problems with the listing's
// layout problems and the stored per-file problems, sorted by path then
// message as store.ReadLog sorts them.
func writeLogProblems(tx *sql.Tx, listing store.Listing) error {
	var out []storedProblem
	for _, p := range listing.Problems {
		out = append(out, storedProblem{Path: p.Path, Sentinel: logSentinel(p.Err), Message: p.Err.Error()})
	}
	rows, err := tx.Query(`SELECT path, problem, sentinel FROM files WHERE ok = 0`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var p storedProblem
		if err := rows.Scan(&p.Path, &p.Message, &p.Sentinel); err != nil {
			rows.Close()
			return err
		}
		out = append(out, p)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	slices.SortFunc(out, func(a, b storedProblem) int {
		return cmp.Or(strings.Compare(a.Path, b.Path), strings.Compare(a.Message, b.Message))
	})
	return replaceProblems(tx, sourceLog, out)
}

// replaceProblems replaces every problem of source with ps, in order.
func replaceProblems(tx *sql.Tx, source string, ps []storedProblem) error {
	if _, err := tx.Exec(`DELETE FROM problems WHERE source = ?`, source); err != nil {
		return err
	}
	if len(ps) == 0 {
		return nil
	}
	st, err := tx.Prepare(`INSERT INTO problems (source, seq, path, change_set, entry, op, sentinel, message)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer st.Close()
	for i, p := range ps {
		if _, err := st.Exec(source, i, p.Path, p.ChangeSet, p.Entry, p.Op, p.Sentinel, p.Message); err != nil {
			return err
		}
	}
	return nil
}

// storedProblem is one row of problems. IDs are canonical text, "" for
// zero. Path is set for log problems, Entry for fold problems and Op for
// view problems.
type storedProblem struct {
	Source    string
	Path      string
	ChangeSet string
	Entry     string
	Op        string
	Sentinel  string
	Message   string
}

// sentinel is a known error and its short name.
type sentinel struct {
	err  error
	name string
}

// logSentinels, foldSentinels and the model's in views.go name the errors
// problems wrap. Order matters where one error wraps another.
var (
	logSentinels = []sentinel{
		{store.ErrStray, "stray"},
		{store.ErrMismatch, "mismatch"},
		{entry.ErrNewerFormat, "newer-format"},
		{entry.ErrInvalid, "invalid"},
		{fs.ErrNotExist, "not-exist"},
		{fs.ErrPermission, "permission"},
	}
	foldSentinels = []sentinel{
		{changeset.ErrBody, "body"},
		{changeset.ErrOrder, "order"},
		{changeset.ErrTransition, "transition"},
		{changeset.ErrRevision, "revision"},
		{changeset.ErrConcurrent, "concurrent"},
		{changeset.ErrUnknownChangeSet, "unknown-change-set"},
	}
)

// nameOf returns the name of the first sentinel err wraps, "" for none.
func nameOf(err error, known []sentinel) string {
	for _, s := range known {
		if errors.Is(err, s.err) {
			return s.name
		}
	}
	return ""
}

func logSentinel(err error) string  { return nameOf(err, logSentinels) }
func foldSentinel(err error) string { return nameOf(err, foldSentinels) }

// idText is i's canonical text, "" for the zero ID.
func idText(i id.ID) string {
	if i.IsZero() {
		return ""
	}
	return i.String()
}

// parseID parses canonical text as idText writes it; "" and anything
// unparsable give the zero ID.
func parseID(s string) id.ID {
	if s == "" {
		return id.ID{}
	}
	i, err := id.Parse(s)
	if err != nil {
		return id.ID{}
	}
	return i
}

// Stamp names the log state the index is level with: the fingerprint of
// the listing it last caught up with. It changes whenever an entry file is
// added, so a value computed from the index stays good while the stamp does.
func (x *Index) Stamp() (string, error) {
	return metaOf(x.db, metaListing)
}

// Refresh catches the index up with the log again, as Open does: a log
// that has not changed costs a read of log/ and a stat per change set
// folder, with no listing, no read and no write. It stats the folders in
// the caller's goroutine: a long-lived reader that refreshes often cares
// for CPU more than time.
func (x *Index) Refresh() error {
	return x.catchUpFrom(takeStatsBy(x.fsys, time.Now(), 1))
}

// CatchUp describes what the last Open or Refresh did.
type CatchUp struct {
	// Fresh is true when the index was already level with the log.
	Fresh bool
	// Rebuilt is true when every table was emptied first.
	Rebuilt bool
	// Read is the number of entry files read; Decoded the number of stored
	// entries decoded and Folded the number of change sets folded.
	Read, Decoded, Folded int
}

// LastCatchUp reports what the last Open or Refresh did.
func (x *Index) LastCatchUp() CatchUp {
	s := x.stats
	return CatchUp{Fresh: s.fresh, Rebuilt: s.rebuilt, Read: s.read, Decoded: s.decoded, Folded: s.folded}
}

// Signature digests what a what-if view of the change sets in sets reads
// of the log: those sets' entries (their number and the newest) and the
// change sets that are no longer in flight (their number and the newest
// terminal entry), which make up as-is. It is equal for two calls when the
// view over sets is, whatever else the log gained between them, a draft's
// new entry say.
func (x *Index) Signature(sets []id.ID) (string, error) {
	var landed int
	var last string
	if err := x.db.QueryRow(`SELECT count(*), coalesce(max(terminal_entry), '') FROM change_sets
		WHERE state NOT IN ('draft', 'proposed', 'accepted')`).Scan(&landed, &last); err != nil {
		return "", err
	}
	h := sha256.New()
	fmt.Fprintf(h, "%d %s\n", landed, last)
	sorted := slices.Clone(sets)
	slices.SortFunc(sorted, id.Compare)
	for _, cs := range slices.Compact(sorted) {
		var n int
		var newest string
		if err := x.db.QueryRow(`SELECT count(*), coalesce(max(id), '') FROM entries WHERE change_set = ?`, idText(cs)).Scan(&n, &newest); err != nil {
			return "", err
		}
		fmt.Fprintf(h, "%s %d %s\n", cs, n, newest)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
