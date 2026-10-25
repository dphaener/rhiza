package index

import (
	"cmp"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strings"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/store"
)

// metaListing is the meta key holding the listing fingerprint.
const metaListing = "listing"

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
		fmt.Fprintf(h, "%s\n", p)
	}
	h.Write([]byte{0})
	for _, p := range l.Problems {
		fmt.Fprintf(h, "%s\t%s\n", p.Path, p.Err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// catchUp brings the index level with the log, as the package doc
// describes. When the listing's fingerprint is the stored one it returns
// without a write transaction.
func (x *Index) catchUp() error {
	x.stats = stats{writes: map[string]int{}}
	fp := fingerprint(store.ListLog(x.fsys))
	stored, err := listingOf(x.db)
	if err != nil {
		return err
	}
	if stored == fp {
		x.stats.fresh = true
		return nil
	}

	tx, err := x.db.Begin() // BEGIN IMMEDIATE, through the DSN
	if err != nil {
		return fmt.Errorf("catch up: %w", err)
	}
	defer tx.Rollback()
	if err := x.catchUpTx(tx); err != nil {
		return fmt.Errorf("catch up: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("catch up: %w", err)
	}
	return nil
}

// querier is what *sql.DB and *sql.Tx share.
type querier interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

// listingOf returns the stored listing fingerprint, "" when there is none.
func listingOf(q querier) (string, error) {
	var v string
	err := q.QueryRow(`SELECT value FROM meta WHERE key = ?`, metaListing).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return v, err
}

// catchUpTx is catch-up once the write lock is held.
func (x *Index) catchUpTx(tx *sql.Tx) error {
	// Another process may have caught up while this one waited.
	listing := store.ListLog(x.fsys)
	fp := fingerprint(listing)
	stored, err := listingOf(tx)
	if err != nil {
		return err
	}
	if stored == fp {
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
		for _, t := range []string{"meta", "files", "entries", "change_sets", "nodes", "links", "problems"} {
			if _, err := tx.Exec(`DELETE FROM ` + t); err != nil {
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
	_, err = tx.Exec(`INSERT OR REPLACE INTO meta (key, value) VALUES (?, ?)`, metaListing, fp)
	return err
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

// entryPath is the slash path of an entry file under .genesis.
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

// refold folds every stored entry with changeset.FoldLog, replaces the
// stored change sets and fold problems, and recomputes the views whose
// inputs changed.
func (x *Index) refold(tx *sql.Tx, touched map[id.ID]bool) error {
	before := map[id.ID]changeset.State{}
	rows, err := tx.Query(`SELECT id, state FROM change_sets`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var s, state string
		if err := rows.Scan(&s, &state); err != nil {
			rows.Close()
			return err
		}
		before[parseID(s)] = changeset.State(state)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	sets, problems, err := foldStored(tx)
	if err != nil {
		return err
	}
	if err := writeChangeSets(tx, sets); err != nil {
		return err
	}
	if err := writeFoldProblems(tx, problems); err != nil {
		return err
	}

	asIs, toBe := x.stats.rebuilt, x.stats.rebuilt
	for _, s := range sets {
		if !touched[s.ID] {
			continue
		}
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
	return writeViews(tx, sets, asIs, toBe, x.stats.writes)
}

// foldStored decodes every stored entry and folds the log.
func foldStored(q querier) ([]changeset.ChangeSet, []changeset.Problem, error) {
	rows, err := q.Query(`SELECT raw FROM entries ORDER BY change_set, id`)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var log store.Log
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return nil, nil, err
		}
		e, err := entry.Decode(raw)
		if err != nil {
			return nil, nil, fmt.Errorf("decode stored entry: %w", err)
		}
		log.Entries = append(log.Entries, e)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	sets, problems := changeset.FoldLog(log)
	return sets, problems, nil
}

// writeChangeSets replaces the stored change sets.
func writeChangeSets(tx *sql.Tx, sets []changeset.ChangeSet) error {
	if _, err := tx.Exec(`DELETE FROM change_sets`); err != nil {
		return err
	}
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
	return nil
}

// writeFoldProblems replaces the stored fold problems, which FoldLog
// returns sorted.
func writeFoldProblems(tx *sql.Tx, ps []changeset.Problem) error {
	out := make([]storedProblem, 0, len(ps))
	for _, p := range ps {
		out = append(out, storedProblem{ChangeSet: idText(p.ChangeSet), Entry: idText(p.Entry),
			Sentinel: foldSentinel(p.Err), Message: p.Err.Error()})
	}
	return replaceProblems(tx, sourceFold, out)
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
