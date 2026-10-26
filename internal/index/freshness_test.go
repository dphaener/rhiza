package index

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/store"
)

// ageLog sets the modification time of log/ and of every folder in it back
// by d, past the racy window, as if the log had been written long ago.
func ageLog(t *testing.T, f *fixture, d time.Duration) {
	t.Helper()
	logDir := filepath.Join(f.gdir, store.LogDir)
	des, err := os.ReadDir(logDir)
	if err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-d)
	for _, de := range des {
		if de.IsDir() {
			if err := os.Chtimes(filepath.Join(logDir, de.Name()), old, old); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := os.Chtimes(logDir, old, old); err != nil {
		t.Fatal(err)
	}
}

// ageFolder sets one change set folder's modification time back by d.
func ageFolder(t *testing.T, f *fixture, cs id.ID, d time.Duration) {
	t.Helper()
	old := time.Now().Add(-d)
	if err := os.Chtimes(filepath.Join(f.gdir, store.LogDir, cs.String()), old, old); err != nil {
		t.Fatal(err)
	}
}

// storedStats returns the stored stat fingerprint, "" when there is none.
func storedStats(t *testing.T, x *Index) string {
	t.Helper()
	v, err := metaOf(x.db, metaStats)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// openFresh opens twice over an aged log and checks the second open is
// fresh by stats: the log is not listed.
func openFresh(t *testing.T, f *fixture, db string) {
	t.Helper()
	ageLog(t, f, time.Hour)
	x := openT(t, f, db, nil)
	if !x.stats.statStored || storedStats(t, x) == "" {
		t.Fatalf("no stat fingerprint stored: %+v", x.stats)
	}
	x.Close()
	x = openT(t, f, db, nil)
	if !x.stats.fresh || !x.stats.statFresh {
		t.Fatalf("not fresh by stats: %+v", x.stats)
	}
	x.Close()
}

// TestFreshByStats checks an unchanged, aged log is fresh without a
// listing, and reads the same model as one built from nothing.
func TestFreshByStats(t *testing.T) {
	f := newFixture(t)
	landedSet(f, 1)
	db := f.dbPath()
	openFresh(t, f, db)

	x := openT(t, f, db, nil)
	if !x.stats.statFresh || x.stats.read != 0 {
		t.Fatalf("stats %+v", x.stats)
	}
	checkModel(t, x, f.gdir)
}

// TestStatNewEntry adds an entry file to a change set folder the stored
// stats cover, then sets the folder's time back past the racy window. The
// folder's stat still differs (a new time, and a new change time where
// the os keeps one), so the open lists the log and catches the entry up.
func TestStatNewEntry(t *testing.T) {
	f := newFixture(t)
	cs := csID(1)
	f.add(cs, changeset.Open{Kind: changeset.KindIntention, Intent: "seed"})
	db := f.dbPath()
	openFresh(t, f, db)

	f.add(cs, changeset.Propose{})
	ageFolder(t, f, cs, 30*time.Minute)
	x := openT(t, f, db, nil)
	if x.stats.fresh || x.stats.statFresh || x.stats.read != 1 {
		t.Fatalf("new entry not caught up: %+v", x.stats)
	}
	if !x.stats.statStored {
		t.Fatalf("aged folder's stats not stored: %+v", x.stats)
	}
	checkModel(t, x, f.gdir)
	if s := storedChangeSets(t, x); len(s) != 1 || s[0].State != string(changeset.StateProposed) {
		t.Fatalf("change sets %+v", s)
	}
	x.Close()

	// And the open after that is fresh by stats again.
	x = openT(t, f, db, nil)
	if !x.stats.statFresh {
		t.Fatalf("not fresh by stats after catch-up: %+v", x.stats)
	}
}

// TestStatNewEntrySameTime adds an entry file and gives the folder back
// exactly the modification time it had. Only the change time (on an os
// that keeps one) tells the folders apart; the entry is caught up.
func TestStatNewEntrySameTime(t *testing.T) {
	f := newFixture(t)
	cs := csID(1)
	f.add(cs, changeset.Open{Kind: changeset.KindIntention, Intent: "seed"})
	db := f.dbPath()
	openFresh(t, f, db)

	dir := filepath.Join(f.gdir, store.LogDir, cs.String())
	fi, err := os.Stat(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := statKey(fi); !ok {
		t.Skip("no change time on this os: a restored modification time is not seen")
	}
	f.add(cs, changeset.Propose{})
	if err := os.Chtimes(dir, fi.ModTime(), fi.ModTime()); err != nil {
		t.Fatal(err)
	}
	x := openT(t, f, db, nil)
	if x.stats.statFresh || x.stats.read != 1 {
		t.Fatalf("entry hidden by a restored time: %+v", x.stats)
	}
	checkModel(t, x, f.gdir)
}

// TestStatNewChangeSet adds a change set folder; log/ gains a name, so
// the stats differ and the new change set is caught up.
func TestStatNewChangeSet(t *testing.T) {
	f := newFixture(t)
	landedSet(f, 1)
	db := f.dbPath()
	openFresh(t, f, db)

	f.add(csID(2), changeset.Open{Kind: changeset.KindIntention, Intent: "next"})
	ageLog(t, f, 30*time.Minute)
	x := openT(t, f, db, nil)
	if x.stats.fresh || x.stats.read != 1 || !x.stats.asIs || !x.stats.toBe {
		t.Fatalf("new change set not caught up: %+v", x.stats)
	}
	checkModel(t, x, f.gdir)
	if s := storedChangeSets(t, x); len(s) != 2 {
		t.Fatalf("change sets %+v", s)
	}
}

// TestStatRemovedFile removes an entry file after the stats were stored:
// the folder's stat changes, the listing no longer holds a seen path, and
// the index is rebuilt as before.
func TestStatRemovedFile(t *testing.T) {
	f := newFixture(t)
	cs := csID(1)
	f.add(cs, changeset.Open{Kind: changeset.KindIntention, Intent: "seed"})
	last := f.add(cs, changeset.Propose{})
	db := f.dbPath()
	openFresh(t, f, db)

	f.remove(last)
	ageLog(t, f, 30*time.Minute)
	x := openT(t, f, db, nil)
	if !x.stats.rebuilt || x.stats.read != 1 {
		t.Fatalf("not rebuilt: %+v", x.stats)
	}
	checkModel(t, x, f.gdir)
}

// TestRacyStats checks a folder modified inside the racy window gets no
// stored stat fingerprint, so the next open decides freshness by listing,
// and that once the folder is old enough the stats are stored.
func TestRacyStats(t *testing.T) {
	f := newFixture(t)
	f.add(csID(1), changeset.Open{Kind: changeset.KindIntention, Intent: "seed"})
	db := f.dbPath()
	x := openT(t, f, db, nil)
	if x.stats.statStored || storedStats(t, x) != "" {
		t.Fatalf("racy stats stored: %+v", x.stats)
	}
	x.Close()

	x = openT(t, f, db, nil)
	if !x.stats.fresh || x.stats.statFresh || x.stats.statStored || storedStats(t, x) != "" {
		t.Fatalf("racy open: %+v", x.stats)
	}
	x.Close()

	// One aged folder and one racy one: still nothing stored.
	f.add(csID(2), changeset.Open{Kind: changeset.KindIntention, Intent: "next"})
	ageFolder(t, f, csID(1), time.Hour)
	x = openT(t, f, db, nil)
	if x.stats.read != 1 || x.stats.statStored || storedStats(t, x) != "" {
		t.Fatalf("partly racy catch-up: %+v", x.stats)
	}
	x.Close()

	// A time in the future (a clock set back) is racy too.
	ageFolder(t, f, csID(2), -time.Hour)
	x = openT(t, f, db, nil)
	if !x.stats.fresh || x.stats.statFresh || storedStats(t, x) != "" {
		t.Fatalf("future time: %+v", x.stats)
	}
	x.Close()

	ageLog(t, f, time.Hour)
	x = openT(t, f, db, nil)
	if !x.stats.fresh || x.stats.statFresh || !x.stats.statStored {
		t.Fatalf("aged log stats not stored: %+v", x.stats)
	}
}

// TestStatTouchedFolder touches a folder without changing its names: the
// stats differ, the listing matches, so the index stays fresh by listing,
// stores the new stats, and the next open is fresh by stats.
func TestStatTouchedFolder(t *testing.T) {
	f := newFixture(t)
	landedSet(f, 1)
	db := f.dbPath()
	openFresh(t, f, db)
	x := openT(t, f, db, nil)
	before := storedStats(t, x)
	x.Close()

	ageFolder(t, f, csID(1), 10*time.Minute)
	x = openT(t, f, db, nil)
	if !x.stats.fresh || x.stats.statFresh || !x.stats.statStored || x.stats.read != 0 {
		t.Fatalf("touched open: %+v", x.stats)
	}
	if after := storedStats(t, x); after == "" || after == before {
		t.Fatalf("stat fingerprint not replaced: %q -> %q", before, after)
	}
	x.Close()

	x = openT(t, f, db, nil)
	if !x.stats.statFresh {
		t.Fatalf("not fresh by stats after touch: %+v", x.stats)
	}
	checkModel(t, x, f.gdir)
}

// TestStatOldIndex opens an index that holds a listing fingerprint but no
// stat fingerprint, as one written before stat fingerprints: it is fresh
// by listing, gains a stat fingerprint, and the next open uses it.
func TestStatOldIndex(t *testing.T) {
	f := newFixture(t)
	landedSet(f, 1)
	db := f.dbPath()
	openFresh(t, f, db)
	x := openT(t, f, db, nil)
	if _, err := x.db.Exec(`DELETE FROM meta WHERE key = ?`, metaStats); err != nil {
		t.Fatal(err)
	}
	x.Close()

	x = openT(t, f, db, nil)
	if !x.stats.fresh || x.stats.statFresh || !x.stats.statStored {
		t.Fatalf("old index open: %+v", x.stats)
	}
	x.Close()
	x = openT(t, f, db, nil)
	if !x.stats.statFresh {
		t.Fatalf("not fresh by stats: %+v", x.stats)
	}
}

// TestStatsSkipHidden checks names store.ListLog skips do not change the
// stats, and a stray file in log/ does.
func TestStatsSkipHidden(t *testing.T) {
	f := newFixture(t)
	f.add(csID(1), changeset.Open{Kind: changeset.KindIntention, Intent: "seed"})
	ageLog(t, f, time.Hour)
	fsys := os.DirFS(f.gdir)
	a := takeStats(fsys, time.Now())
	if !a.storable() {
		t.Fatalf("aged stats not storable: %+v", a)
	}
	writeFile(t, f.gdir, "log/.DS_Store", []byte("x"))
	if b := takeStats(fsys, time.Now()); b.fp != a.fp {
		t.Fatal("a hidden file changed the stats")
	}
	writeFile(t, f.gdir, "log/stray.json", []byte("{}"))
	if c := takeStats(fsys, time.Now()); c.fp == a.fp {
		t.Fatal("a stray file did not change the stats")
	}
	if p := takeStats(os.DirFS(t.TempDir()), time.Now()); p.fp != "" || p.storable() {
		t.Fatalf("missing log folder gave stats %+v", p)
	}
}
