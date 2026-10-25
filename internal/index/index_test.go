package index

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
	"github.com/dphaener/genesis/internal/store"
)

// openT opens the index for f through fsys (os.DirFS when nil).
func openT(t *testing.T, f *fixture, dbPath string, c *countingFS) *Index {
	t.Helper()
	var x *Index
	var err error
	if c != nil {
		x, err = open(dbPath, c)
	} else {
		x, err = Open(f.root, dbPath)
	}
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { x.Close() })
	return x
}

func catchUpT(t *testing.T, x *Index) {
	t.Helper()
	if err := x.catchUp(); err != nil {
		t.Fatalf("catchUp: %v", err)
	}
}

// landedSet writes a landed change set creating nodes 1 and 2 (2 under 1)
// and a link between them, and returns its ID.
func landedSet(f *fixture, n int) id.ID {
	cs := csID(n)
	f.add(cs, changeset.Open{Kind: changeset.KindIntention, Intent: "seed"})
	f.add(cs, changeset.Ops{Ops: []op.Op{
		mkCap(1, 1, "root", 0, "i"),
		mkCap(2, 2, "child", 1, "i"),
		{ID: opID(3), Type: op.TypeLinkCreate, Link: ln(1), LinkType: op.LinkServes, From: nd(2), To: nd(1)},
	}})
	f.add(cs, changeset.Propose{})
	f.add(cs, changeset.Accept{Rule: "owner"})
	f.add(cs, changeset.Land{})
	return cs
}

func update(n, node int, base id.ID, statement string) op.Op {
	return op.Op{ID: opID(n), Type: op.TypeNodeUpdate, Node: nd(node), Base: base,
		Set: map[string]json.RawMessage{"statement": raw(`"` + statement + `"`)}}
}

func TestEmptyLog(t *testing.T) {
	f := newFixture(t)
	db := f.dbPath()
	x := openT(t, f, db, nil)
	if x.stats.fresh || x.stats.read != 0 {
		t.Fatalf("first open stats %+v", x.stats)
	}
	checkModel(t, x, f.gdir)
	x.Close()

	x = openT(t, f, db, nil)
	if !x.stats.fresh {
		t.Fatalf("second open not fresh: %+v", x.stats)
	}

	// No .genesis at all: a log problem, no error.
	root := t.TempDir()
	y, err := Open(root, filepath.Join(t.TempDir(), "i.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer y.Close()
	ps := storedProblems(t, y, sourceLog)
	if len(ps) != 1 || ps[0].Path != "log" || ps[0].Sentinel != "not-exist" {
		t.Fatalf("log problems %+v", ps)
	}
}

func TestFirstBuild(t *testing.T) {
	f := newFixture(t)
	landedSet(f, 1)
	acc := csID(2)
	f.add(acc, changeset.Open{Kind: changeset.KindIntention, Intent: "edit"})
	f.add(acc, changeset.Ops{Ops: []op.Op{update(10, 1, opID(1), "hello")}})
	f.add(acc, changeset.Propose{})
	f.add(acc, changeset.Accept{Rule: "owner"})

	x := openT(t, f, f.dbPath(), nil)
	if !x.stats.asIs || !x.stats.toBe || x.stats.read != 9 {
		t.Fatalf("stats %+v", x.stats)
	}
	checkModel(t, x, f.gdir)
	if n := storedNodes(t, x, viewToBe); len(n) != 2 || string(n[0].Fields["statement"]) != `"hello"` {
		t.Fatalf("to-be nodes %+v", n)
	}
	if n := storedNodes(t, x, viewAsIs); n[0].Fields != nil {
		t.Fatalf("as-is node 1 has fields %v", n[0].Fields)
	}
}

// TestDirty walks one change set through its life and checks which views
// each new entry recomputes.
func TestDirty(t *testing.T) {
	f := newFixture(t)
	landedSet(f, 1)
	x := openT(t, f, f.dbPath(), nil)
	cs := csID(2)

	steps := []struct {
		name       string
		body       changeset.Body
		asIs, toBe bool
		writes     map[string]int
	}{
		// A change set new to the index may be a link endpoint: both
		// views are recomputed, but no row changes.
		{"open", changeset.Open{Kind: changeset.KindIntention, Intent: "x"}, true, true, map[string]int{}},
		{"draft ops", changeset.Ops{Ops: []op.Op{update(10, 1, opID(1), "a")}}, false, false, map[string]int{}},
		{"propose", changeset.Propose{}, false, false, map[string]int{}},
		{"accept", changeset.Accept{Rule: "owner"}, false, true, map[string]int{viewToBe: 1}},
		{"accepted ops", changeset.Ops{Ops: []op.Op{update(11, 2, opID(2), "b")}}, false, true, map[string]int{viewToBe: 1}},
		{"land", changeset.Land{}, true, true, map[string]int{viewAsIs: 2}},
		{"note after land", changeset.Note{Text: "done"}, true, true, map[string]int{}},
	}
	for _, s := range steps {
		before := dump(t, x)
		f.add(cs, s.body)
		catchUpT(t, x)
		st := x.stats
		if st.fresh || st.rebuilt || st.read != 1 || st.asIs != s.asIs || st.toBe != s.toBe {
			t.Fatalf("%s: stats %+v, want asIs %v toBe %v", s.name, st, s.asIs, s.toBe)
		}
		if !maps.Equal(st.writes, s.writes) {
			t.Fatalf("%s: writes %v, want %v", s.name, st.writes, s.writes)
		}
		after := dump(t, x)
		if s.writes[viewAsIs]+s.writes[viewToBe] == 0 {
			same(t, s.name+" nodes", after["nodes"], before["nodes"])
			same(t, s.name+" links", after["links"], before["links"])
		}
		checkModel(t, x, f.gdir)
	}
}

func TestVanishRebuilds(t *testing.T) {
	f := newFixture(t)
	landedSet(f, 1)
	cs := csID(2)
	f.add(cs, changeset.Open{Kind: changeset.KindIntention, Intent: "x"})
	ops := f.add(cs, changeset.Ops{Ops: []op.Op{update(10, 1, opID(1), "a")}})
	x := openT(t, f, f.dbPath(), nil)
	checkModel(t, x, f.gdir)

	f.remove(ops)
	catchUpT(t, x)
	if !x.stats.rebuilt || x.stats.read != 6 || !x.stats.asIs || !x.stats.toBe {
		t.Fatalf("after vanish: %+v", x.stats)
	}
	checkModel(t, x, f.gdir)

	// Back again: an ordinary new path, read alone.
	f.add(cs, changeset.Note{Text: "n"}) // after the vanished entry: an order problem
	catchUpT(t, x)
	if x.stats.rebuilt || x.stats.read != 1 {
		t.Fatalf("after re-add: %+v", x.stats)
	}
	checkModel(t, x, f.gdir)
	if len(storedProblems(t, x, sourceFold)) == 0 {
		t.Fatal("want an order problem for the note after the vanished entry")
	}
}

func TestCorruptEntryNotReread(t *testing.T) {
	f := newFixture(t)
	landedSet(f, 1)
	bad := entryPath(csID(3), eID(900))
	writeFile(t, f.gdir, bad, []byte("not json"))
	db := f.dbPath()

	c := newCountingFS(f.gdir)
	x := openT(t, f, db, c)
	if c.reads[bad] != 1 || c.total() != 6 {
		t.Fatalf("first open reads %v", c.reads)
	}
	ps := storedProblems(t, x, sourceLog)
	if len(ps) != 1 || ps[0].Path != bad || ps[0].Sentinel != "invalid" {
		t.Fatalf("log problems %+v", ps)
	}
	checkModel(t, x, f.gdir)
	x.Close()

	c = newCountingFS(f.gdir)
	x = openT(t, f, db, c)
	if c.total() != 0 || !x.stats.fresh {
		t.Fatalf("second open reads %v, stats %+v", c.reads, x.stats)
	}
	x.Close()

	f.add(csID(4), changeset.Open{Kind: changeset.KindIntention, Intent: "y"})
	c = newCountingFS(f.gdir)
	x = openT(t, f, db, c)
	if c.total() != 1 || c.reads[bad] != 0 {
		t.Fatalf("third open reads %v", c.reads)
	}
	checkModel(t, x, f.gdir)

	// A stray file is a layout problem, re-listed but never read.
	writeFile(t, f.gdir, "log/README", []byte("hi"))
	catchUpT(t, x)
	if x.stats.fresh || x.stats.read != 0 {
		t.Fatalf("stray stats %+v", x.stats)
	}
	checkModel(t, x, f.gdir)
}

func TestSchemaStaleVersion(t *testing.T) {
	for _, tc := range []struct {
		name   string
		break_ func(t *testing.T, db *sql.DB)
	}{
		{"newer version", func(t *testing.T, db *sql.DB) {
			if _, err := db.Exec("PRAGMA user_version = 99"); err != nil {
				t.Fatal(err)
			}
		}},
		{"no version", func(t *testing.T, db *sql.DB) {
			if _, err := db.Exec("PRAGMA user_version = 0"); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			landedSet(f, 1)
			db := f.dbPath()
			x := openT(t, f, db, nil)
			x.Close()

			raw, err := sql.Open("sqlite", dsn(db))
			if err != nil {
				t.Fatal(err)
			}
			tc.break_(t, raw)
			if _, err := raw.Exec("DELETE FROM nodes"); err != nil {
				t.Fatal(err)
			}
			raw.Close()

			x = openT(t, f, db, nil)
			if x.stats.read != 5 || !x.stats.asIs {
				t.Fatalf("not rebuilt: %+v", x.stats)
			}
			checkModel(t, x, f.gdir)
			if v := userVersion(t, db); v != SchemaVersion {
				t.Fatalf("user_version %d", v)
			}
		})
	}
}

func TestCorruptDBFile(t *testing.T) {
	for _, tc := range []struct {
		name string
		data []byte
	}{
		{"garbage", []byte(strings.Repeat("this is not a database ", 300))},
		{"truncated header", []byte("SQLite format 3\x00\x10\x00")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			landedSet(f, 1)
			db := f.dbPath()
			if err := os.WriteFile(db, tc.data, 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(db+"-wal", []byte("junk"), 0o644); err != nil {
				t.Fatal(err)
			}
			x := openT(t, f, db, nil)
			checkModel(t, x, f.gdir)
			if v := userVersion(t, db); v != SchemaVersion {
				t.Fatalf("user_version %d", v)
			}
		})
	}
}

// TestLockWaits: catch-up waits for another connection's write lock, and
// fails when the busy timeout runs out first.
func TestLockWaits(t *testing.T) {
	f := newFixture(t)
	landedSet(f, 1)
	db := f.dbPath()
	x := openT(t, f, db, nil)
	x.Close()

	hold := func() (release func()) {
		other, err := sql.Open("sqlite", dsn(db))
		if err != nil {
			t.Fatal(err)
		}
		ctx := context.Background()
		conn, err := other.Conn(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := conn.ExecContext(ctx, "BEGIN IMMEDIATE"); err != nil {
			t.Fatal(err)
		}
		return func() {
			conn.ExecContext(ctx, "ROLLBACK")
			conn.Close()
			other.Close()
		}
	}
	old := busyTimeout
	t.Cleanup(func() { busyTimeout = old })

	// Nothing new: no write lock needed, so a held lock does not matter.
	release := hold()
	busyTimeout = 50 * time.Millisecond
	x = openT(t, f, db, nil)
	if !x.stats.fresh {
		t.Fatalf("stats %+v", x.stats)
	}
	x.Close()

	// Something new and a short timeout: an error.
	f.add(csID(2), changeset.Open{Kind: changeset.KindIntention, Intent: "x"})
	if y, err := Open(f.root, db); err == nil {
		y.Close()
		t.Fatal("Open succeeded while the lock was held")
	} else {
		t.Logf("held lock: %v", err)
	}

	// A longer timeout: Open waits and succeeds once the lock is released.
	busyTimeout = 5 * time.Second
	start := time.Now()
	go func() {
		time.Sleep(300 * time.Millisecond)
		release()
	}()
	x = openT(t, f, db, nil)
	waited := time.Since(start)
	if waited < 200*time.Millisecond {
		t.Fatalf("Open returned after %v, before the lock was released", waited)
	}
	t.Logf("waited %v", waited)
	if x.stats.read != 1 {
		t.Fatalf("stats %+v", x.stats)
	}
	checkModel(t, x, f.gdir)
}

func TestDefaultPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CACHE_HOME", filepath.Join(home, "cache"))
	cache, err := os.UserCacheDir()
	if err != nil {
		t.Fatal(err)
	}

	base := t.TempDir()
	a := filepath.Join(base, "a")
	b := filepath.Join(base, "b")
	link := filepath.Join(base, "link")
	for _, d := range []string{a, b} {
		if err := os.Mkdir(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(a, link); err != nil {
		t.Fatal(err)
	}
	pa, err := DefaultPath(a)
	if err != nil {
		t.Fatal(err)
	}
	pl, err := DefaultPath(link)
	if err != nil {
		t.Fatal(err)
	}
	pb, err := DefaultPath(b)
	if err != nil {
		t.Fatal(err)
	}
	if pa != pl {
		t.Fatalf("symlinked root: %s vs %s", pa, pl)
	}
	if pa == pb {
		t.Fatalf("different roots share %s", pa)
	}
	if dir := filepath.Join(cache, "genesis", "index"); filepath.Dir(pa) != dir || filepath.Ext(pa) != ".db" {
		t.Fatalf("path %s not a .db in %s", pa, dir)
	}
	if _, err := DefaultPath(filepath.Join(base, "missing")); err == nil {
		t.Fatal("missing root: no error")
	}
}

// TestReadOnly: building, catching up and rebuilding never change anything
// under .genesis.
func TestReadOnly(t *testing.T) {
	f := newFixture(t)
	landedSet(f, 1)
	writeFile(t, f.gdir, entryPath(csID(3), eID(900)), []byte("{}"))
	writeFile(t, f.gdir, "log/stray.txt", []byte("x"))
	db := f.dbPath()

	check := func(step string, do func()) {
		t.Helper()
		before := snapshot(t, f.gdir)
		do()
		same(t, step+": .genesis", snapshot(t, f.gdir), before)
		t.Logf("%s: %d paths under .genesis unchanged", step, len(before))
	}
	check("build", func() { openT(t, f, db, nil).Close() })
	check("fresh", func() { openT(t, f, db, nil).Close() })
	f.add(csID(2), changeset.Open{Kind: changeset.KindIntention, Intent: "x"})
	check("catch up", func() { openT(t, f, db, nil).Close() })
	f.remove(f.heads[csID(2)])
	check("rebuild", func() {
		x := openT(t, f, db, nil)
		if !x.stats.rebuilt {
			t.Fatalf("not rebuilt: %+v", x.stats)
		}
		x.Close()
	})
	if err := os.WriteFile(db, []byte("garbage garbage garbage"), 0o644); err != nil {
		t.Fatal(err)
	}
	check("garbage db", func() { openT(t, f, db, nil).Close() })
	if _, err := os.Stat(filepath.Join(f.gdir, store.LogDir)); errors.Is(err, os.ErrNotExist) {
		t.Fatal("log folder gone")
	}
}

// TestNewChangeSetEndpoint: a landed link may name a change set the log
// does not hold yet; when that change set's first entry arrives, even as a
// draft, the link applies and both views change.
func TestNewChangeSetEndpoint(t *testing.T) {
	f := newFixture(t)
	cs := csID(1)
	f.add(cs, changeset.Open{Kind: changeset.KindRecord, Intent: "seed"})
	f.add(cs, changeset.Ops{Ops: []op.Op{
		mkCap(1, 1, "root", 0, "i"),
		{ID: opID(2), Type: op.TypeLinkCreate, Link: ln(1), LinkType: op.LinkRelatesTo, From: nd(1), To: csID(5), Note: "see"},
	}})
	f.add(cs, changeset.Land{})
	x := openT(t, f, f.dbPath(), nil)
	checkModel(t, x, f.gdir)
	if l := storedLinks(t, x, viewAsIs); l != nil {
		t.Fatalf("link to a missing change set applied: %+v", l)
	}

	f.add(csID(5), changeset.Open{Kind: changeset.KindIntention, Intent: "later"})
	catchUpT(t, x)
	if !x.stats.asIs || !x.stats.toBe || x.stats.writes[viewAsIs] != 1 {
		t.Fatalf("stats %+v", x.stats)
	}
	checkModel(t, x, f.gdir)
}
