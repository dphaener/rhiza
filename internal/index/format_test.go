package index

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
)

// TestEntryFormatChangeRebuilds stores an entry as newer-format (as an
// older binary would), then opens the unchanged log with a binary that
// reads that format: the index is rebuilt, the change set folds landed
// with no problems, and a further open does not rebuild.
func TestEntryFormatChangeRebuilds(t *testing.T) {
	f := newFixture(t)
	cs := csID(1)
	f.add(cs, changeset.Open{Kind: changeset.KindIntention, Intent: "seed"})
	pin := f.add(cs, changeset.Propose{})
	f.add(cs, changeset.Accept{Rule: "owner"})
	f.add(cs, changeset.Land{})

	// Stand for the newer-format entry: bump its v on disk.
	pinPath := filepath.Join(f.gdir, filepath.FromSlash(f.paths[pin]))
	good, err := os.ReadFile(pinPath)
	if err != nil {
		t.Fatal(err)
	}
	old := `"v": 1`
	if !bytes.Contains(good, []byte(old)) {
		t.Fatalf("no %s in entry:\n%s", old, good)
	}
	newer := bytes.Replace(good, []byte(old), []byte(fmt.Sprintf(`"v": %d`, entry.FormatVersion+1)), 1)
	if err := os.WriteFile(pinPath, newer, 0o644); err != nil {
		t.Fatal(err)
	}

	db := f.dbPath()
	openFresh(t, f, db)
	x := openT(t, f, db, nil)
	if len(storedProblems(t, x, sourceLog)) == 0 {
		t.Fatalf("no newer-format problem stored")
	}
	x.Close()

	// The binary is upgraded; the log's files are as the new binary reads
	// them, and the stat fingerprint still matches.
	if err := os.WriteFile(pinPath, good, 0o644); err != nil {
		t.Fatal(err)
	}
	ageLog(t, f, time.Hour)
	prev := entryFormat
	entryFormat = prev + 1
	t.Cleanup(func() { entryFormat = prev })

	x = openT(t, f, db, nil)
	if x.stats.read != 4 {
		t.Fatalf("not rebuilt, stats %+v", x.stats)
	}
	if p := storedProblems(t, x, sourceLog); len(p) != 0 {
		t.Fatalf("problems after rebuild: %+v", p)
	}
	checkModel(t, x, f.gdir)
	sets := storedChangeSets(t, x)
	if len(sets) != 1 || sets[0].State != "landed" {
		t.Fatalf("change sets %+v", sets)
	}
	x.Close()

	x = openT(t, f, db, nil)
	if x.stats.read != 0 || !x.stats.fresh {
		t.Fatalf("rebuilt again: %+v", x.stats)
	}
}
