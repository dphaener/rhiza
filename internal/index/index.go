// Package index keeps a SQLite index of a worktree's log: its entries,
// folded change sets, problems and the as-is and to-be views, so that read
// commands answer from a few queries instead of reading and projecting the
// whole log.
//
// The vocabulary:
//
//   - The index is a SQLite database derived from the log alone. It is
//     disposable: deleting it loses nothing, and it is never committed or
//     shared. There is one index per worktree, in the user cache folder
//     (DefaultPath), keyed by the worktree's resolved absolute path.
//   - Catch-up brings the index level with the log. It is
//     input-incremental: the log is listed by file name (store.ListLog), and
//     only entry files whose path the index has not seen are read and
//     decoded. An entry file, once written, never changes, so a seen path
//     never needs reading again; a file that is unreadable or does not
//     decode is remembered with its problem, so it is not re-read either.
//     Content edits to an existing file are not noticed: finding those is
//     the job of a full check, not of the index.
//   - The index is fresh when the listing equals the one it last caught up
//     with. Open stores a fingerprint of the listing (every path and layout
//     problem), so the fresh case costs one listing of log/ and one row read,
//     with no write transaction.
//   - A rebuild empties every table and catches up from nothing. It happens
//     when a path the index has seen is no longer listed (the log is
//     append-only, so a vanished file means something rewrote it), and,
//     as a new file, when the database file is unreadable, not a database,
//     or of another schema version.
//
// What is stored: every listed entry file with its problem if it has one;
// each good entry's canonical bytes; each change set's folded state and the
// kind and intent of its open; every node and link of the as-is and to-be
// views with all their fields (fields and history as JSON) and a hash of
// their content; and the problems of the log (layout and per-file), of
// folding and of each view, each with the short name of the sentinel its
// error wraps. What-if views are not stored; they are computed on demand
// from the stored entries.
//
// When something new is read, catch-up refolds the change sets from every
// stored entry with changeset.FoldLog, so the stored change sets and fold
// problems always equal FoldLog over the whole log. A view is recomputed
// only when its inputs changed: as-is when a touched change set was or now
// is landed, to-be when one was or is landed or accepted, and both when a
// change set new to the index appears (it may be the endpoint of a link).
// A recomputed view is the whole of model.AsIs or model.ToBe over every
// folded change set, written by difference: rows whose hash changed are
// upserted, rows gone are deleted, and the view's problems are replaced.
// The index never reimplements the model's rules; it stores what
// internal/model computes.
//
// Failures: the schema version is kept in PRAGMA user_version, and a file
// with another version, or one that is not a SQLite database or is
// corrupt, is deleted (with its -wal and -shm files) and built again. The
// database runs in WAL mode, so readers never wait for a writer. Catch-up
// writes in one BEGIN IMMEDIATE transaction and re-lists the log once it
// holds the lock, since another process may have caught up meanwhile. A
// writer waits up to five seconds for the lock; a timeout is an error. The
// index never writes, creates or deletes anything under .genesis.
//
// The driver is modernc.org/sqlite, a translation of SQLite to pure Go: no
// cgo, so Genesis stays one statically built binary.
package index

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"modernc.org/sqlite"

	"github.com/dphaener/genesis/internal/store"
)

// SchemaVersion is the schema this package writes, kept in PRAGMA
// user_version. Bump it on any change to the tables or what they hold: an
// index of another version is deleted and rebuilt.
const SchemaVersion = 1

// busyTimeout is how long a connection waits for a lock another connection
// holds. Tests shorten it.
var busyTimeout = 5 * time.Second

// SQLite primary result codes the index tells apart.
const (
	codeBusy    = 5
	codeCorrupt = 11
	codeNotADB  = 26
)

// Index is an open index of one worktree's log. It is not safe for
// concurrent use by several goroutines.
type Index struct {
	db   *sql.DB
	path string
	fsys fs.FS
	// stats describes the last catch-up, for tests.
	stats stats
}

// stats records what one catch-up did.
type stats struct {
	// fresh is true when the listing matched and nothing was written.
	fresh bool
	// rebuilt is true when every table was emptied first.
	rebuilt bool
	// read is the number of entry files read.
	read int
	// asIs and toBe are true when the view was recomputed.
	asIs, toBe bool
	// writes counts node and link rows upserted or deleted, per view.
	writes map[string]int
}

// DefaultPath returns the index file for the worktree at root:
// os.UserCacheDir()/genesis/index/<hex SHA-256 of root's absolute,
// symlink-resolved path>.db. root must exist.
func DefaultPath(root string) (string, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("index path: %w", err)
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", fmt.Errorf("index path: %w", err)
	}
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("index path: %w", err)
	}
	sum := sha256.Sum256([]byte(resolved))
	return filepath.Join(cache, "genesis", "index", hex.EncodeToString(sum[:])+".db"), nil
}

// Open opens the index at dbPath for the worktree at root, creating,
// rebuilding or catching it up as needed, and returns it fresh. Its folder
// is created when missing. It reads root's .genesis folder and never
// writes there.
func Open(root, dbPath string) (*Index, error) {
	return open(dbPath, os.DirFS(filepath.Join(root, store.Dir)))
}

// open is Open over fsys, rooted at a .genesis folder.
func open(dbPath string, fsys fs.FS) (*Index, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("open index: %w", err)
	}
	db, err := openDB(dbPath)
	if errors.Is(err, errStale) {
		if err := removeDB(dbPath); err != nil {
			return nil, fmt.Errorf("open index: %w", err)
		}
		db, err = openDB(dbPath)
	}
	if err != nil {
		return nil, fmt.Errorf("open index %s: %w", dbPath, err)
	}
	x := &Index{db: db, path: dbPath, fsys: fsys}
	if err := x.catchUp(); err != nil {
		db.Close()
		return nil, fmt.Errorf("open index %s: %w", dbPath, err)
	}
	return x, nil
}

// Close closes the database.
func (x *Index) Close() error {
	return x.db.Close()
}

// errStale marks a database file to delete and build again.
var errStale = errors.New("index: stale or unreadable database")

// dsn is the connection string: every connection the pool opens gets the
// busy timeout, WAL and the other pragmas, and its transactions start with
// BEGIN IMMEDIATE.
func dsn(path string) string {
	q := []string{
		fmt.Sprintf("_pragma=busy_timeout(%d)", busyTimeout.Milliseconds()),
		"_pragma=journal_mode(WAL)",
		"_pragma=synchronous(NORMAL)",
		"_pragma=foreign_keys(0)",
		"_txlock=immediate",
	}
	return "file:" + path + "?" + strings.Join(q, "&")
}

// openDB opens the database and makes sure it holds the current schema. An
// error wrapping errStale means the file should be deleted and recreated.
func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn(path))
	if err != nil {
		return nil, err
	}
	if err := ensureSchema(db); err != nil {
		db.Close()
		if damaged(err) {
			return nil, fmt.Errorf("%w: %w", errStale, err)
		}
		return nil, err
	}
	return db, nil
}

// damaged reports whether err says the file is not a database or is
// corrupt.
func damaged(err error) bool {
	if errors.Is(err, errStale) {
		return true
	}
	var se *sqlite.Error
	if errors.As(err, &se) {
		switch se.Code() & 0xff {
		case codeCorrupt, codeNotADB:
			return true
		}
	}
	return false
}

// ensureSchema checks user_version and creates the schema in a new file.
// A version other than SchemaVersion, or tables in a file with no version,
// is errStale.
func ensureSchema(db *sql.DB) error {
	var v int
	if err := db.QueryRow("PRAGMA user_version").Scan(&v); err != nil {
		return err
	}
	switch v {
	case SchemaVersion:
		return nil
	case 0:
	default:
		return fmt.Errorf("%w: schema version %d, want %d", errStale, v, SchemaVersion)
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Another process may have created it while this one waited.
	if err := tx.QueryRow("PRAGMA user_version").Scan(&v); err != nil {
		return err
	}
	if v == SchemaVersion {
		return nil
	}
	var n int
	if err := tx.QueryRow("SELECT count(*) FROM sqlite_master").Scan(&n); err != nil {
		return err
	}
	if v != 0 || n != 0 {
		return fmt.Errorf("%w: schema version %d with %d objects", errStale, v, n)
	}
	for _, s := range schema {
		if _, err := tx.Exec(s); err != nil {
			return fmt.Errorf("create schema: %w", err)
		}
	}
	if _, err := tx.Exec(fmt.Sprintf("PRAGMA user_version = %d", SchemaVersion)); err != nil {
		return err
	}
	return tx.Commit()
}

// removeDB deletes the database file and its WAL and shared-memory files.
func removeDB(path string) error {
	for _, p := range []string{path, path + "-wal", path + "-shm"} {
		if err := os.Remove(p); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	return nil
}

// schema creates the tables. IDs are stored as their canonical text, the
// zero ID as "", so text order is ID order.
var schema = []string{
	// meta holds the listing fingerprint the index last caught up with.
	`CREATE TABLE meta (
		key   TEXT PRIMARY KEY,
		value TEXT NOT NULL
	) WITHOUT ROWID`,
	// files is every listed entry file the index has read. ok is 1 for a
	// good entry; otherwise problem holds why it is bad, and the file is
	// not read again.
	`CREATE TABLE files (
		path       TEXT PRIMARY KEY,
		change_set TEXT NOT NULL,
		entry_id   TEXT NOT NULL,
		ok         INTEGER NOT NULL,
		problem    TEXT NOT NULL,
		sentinel   TEXT NOT NULL
	) WITHOUT ROWID`,
	// entries holds each good entry's canonical bytes.
	`CREATE TABLE entries (
		change_set TEXT NOT NULL,
		id         TEXT NOT NULL,
		raw        BLOB NOT NULL,
		PRIMARY KEY (change_set, id)
	) WITHOUT ROWID`,
	// change_sets is each change set as folded. state is "" when it has no
	// valid open.
	`CREATE TABLE change_sets (
		id             TEXT PRIMARY KEY,
		state          TEXT NOT NULL,
		kind           TEXT NOT NULL,
		intent         TEXT NOT NULL,
		open_entry     TEXT NOT NULL,
		superseded_by  TEXT NOT NULL,
		terminal_entry TEXT NOT NULL
	) WITHOUT ROWID`,
	// nodes and links hold the as-is ("asis") and to-be ("tobe") views.
	// fields is a JSON object of raw values or NULL when empty; history a
	// JSON array of [change set, op] pairs or NULL when empty.
	`CREATE TABLE nodes (
		view    TEXT NOT NULL,
		id      TEXT NOT NULL,
		kind    TEXT NOT NULL,
		handle  TEXT NOT NULL,
		parent  TEXT NOT NULL,
		key     TEXT NOT NULL,
		fields  TEXT,
		retired INTEGER NOT NULL,
		reason  TEXT NOT NULL,
		version TEXT NOT NULL,
		created TEXT NOT NULL,
		history TEXT,
		hash    BLOB NOT NULL,
		PRIMARY KEY (view, id)
	) WITHOUT ROWID`,
	`CREATE INDEX nodes_parent ON nodes (view, parent, key)`,
	`CREATE INDEX nodes_handle ON nodes (view, handle)`,
	`CREATE TABLE links (
		view     TEXT NOT NULL,
		id       TEXT NOT NULL,
		type     TEXT NOT NULL,
		from_id  TEXT NOT NULL,
		to_id    TEXT NOT NULL,
		note     TEXT NOT NULL,
		decision TEXT NOT NULL,
		retired  INTEGER NOT NULL,
		version  TEXT NOT NULL,
		created  TEXT NOT NULL,
		history  TEXT,
		hash     BLOB NOT NULL,
		PRIMARY KEY (view, id)
	) WITHOUT ROWID`,
	`CREATE INDEX links_from ON links (view, from_id)`,
	`CREATE INDEX links_to ON links (view, to_id)`,
	// problems holds every problem, by source: "log" (layout and per-file,
	// path set), "fold" (change_set and entry set), "asis" and "tobe"
	// (change_set and op set). seq is the problem's place in its source's
	// sorted list. sentinel is the short name of the error it wraps.
	`CREATE TABLE problems (
		source     TEXT NOT NULL,
		seq        INTEGER NOT NULL,
		path       TEXT NOT NULL,
		change_set TEXT NOT NULL,
		entry      TEXT NOT NULL,
		op         TEXT NOT NULL,
		sentinel   TEXT NOT NULL,
		message    TEXT NOT NULL,
		PRIMARY KEY (source, seq)
	) WITHOUT ROWID`,
}
