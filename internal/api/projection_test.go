package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
)

// indexStamp is the stamp of the index at ctx.IndexPath, opened apart from
// any call.
func indexStamp(t *testing.T, ctx *Ctx) string {
	t.Helper()
	x, err := index.Open(ctx.Root, ctx.IndexPath)
	if err != nil {
		t.Fatalf("index.Open: %v", err)
	}
	defer x.Close()
	s, err := x.Stamp()
	if err != nil {
		t.Fatalf("Stamp: %v", err)
	}
	return s
}

// dryRunWorkspace opens ctx's workspace as a dry run with its projection.
func dryRunWorkspace(t *testing.T, ctx *Ctx) *workspace {
	t.Helper()
	c := *ctx
	c.DryRun = true
	w, err := openCtxWorkspace(&c)
	if err != nil {
		t.Fatal(err)
	}
	w.useProjection(&c)
	t.Cleanup(w.closeProjection)
	return w
}

// TestProjectionDryRunSeesPending: a dry run's pending entries (a new
// change set and an ops entry in it) are part of the projection's change
// sets, entries and model, and the index file is left as it was.
func TestProjectionDryRunSeesPending(t *testing.T) {
	t.Parallel()
	ctx, _ := indexedLog(t)
	before := indexStamp(t, ctx)
	w := dryRunWorkspace(t, ctx)
	cs, _, err := w.openChangeSet(changeset.Open{Kind: changeset.KindIntention, Intent: "Pending"}, id.ID{}, id.ID{})
	if err != nil {
		t.Fatal(err)
	}
	items, err := decodeWrites([]byte(`{"writes":[{"op":"add","kind":"capability","handle":"zip","parent":"store"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	e, err := buildApply(w, cs, w.ids, items)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.appendEntry(cs, changeset.Ops{Ops: e.ops}, id.ID{}, id.ID{}); err != nil {
		t.Fatal(err)
	}
	if len(w.pending) != 2 {
		t.Fatalf("pending = %d entries, want 2", len(w.pending))
	}

	p := w.projection()
	if p == nil {
		t.Fatal("no projection with ctx.IndexPath set")
	}
	if !slices.Contains(p.changeSetIDs(w.pending), cs) {
		t.Errorf("changeSetIDs lacks the pending change set %s", cs)
	}
	var found *changeset.ChangeSet
	for _, c := range p.changeSets(w.pending) {
		if c.ID == cs {
			found = &c
		}
	}
	if found == nil || found.State != changeset.StateDraft || len(found.Ops) != 1 {
		t.Errorf("changeSets: pending change set = %+v, want a draft with one op", found)
	}
	if es := p.entries(w.pending, cs); len(es) != 2 {
		t.Errorf("entries of the pending change set = %d, want 2", len(es))
	}
	if n, err := p.model(w.pending, cs).Resolve("zip"); err != nil || len(n.Pending) != 1 {
		t.Errorf("model of the pending change set: zip = %+v, %v; want it with one pending op", n, err)
	}
	if _, err := p.model(w.pending, id.ID{}).Resolve("zip"); err == nil {
		t.Error("as-is resolves zip, which only the pending change set adds")
	}
	x, err := index.Open(ctx.Root, ctx.IndexPath)
	if err != nil {
		t.Fatal(err)
	}
	stored, err := x.Problems()
	x.Close()
	if err != nil {
		t.Fatal(err)
	}
	if got := p.problems(); len(got) != len(stored) {
		t.Errorf("problems = %v, want the stored %v", got, stored)
	}
	if err := w.projectionFailure(); err != nil {
		t.Fatalf("projection failed: %v", err)
	}
	w.closeProjection()
	if after := indexStamp(t, ctx); after != before {
		t.Errorf("index stamp moved in a dry run: %s, then %s", before, after)
	}
}

// TestProjectionOpensOnce: one call that asks the projection for the
// model, the change sets and a change set's entries three times opens the
// index once, and a write opens it once.
func TestProjectionOpensOnce(t *testing.T) {
	t.Parallel()
	ctx, sets := indexedLog(t)
	opens := 0
	ctx.openedProjection = func() { opens++ }
	w := dryRunWorkspace(t, ctx)
	if opens != 0 {
		t.Fatalf("workspace opened the index %d times before a read", opens)
	}
	for range 3 {
		w.model(sets[2])
		w.changeSets()
		w.entriesOfSet(sets[2])
	}
	if opens != 1 {
		t.Errorf("three rounds of reads opened the index %d times, want 1", opens)
	}
	if err := w.projectionFailure(); err != nil {
		t.Fatal(err)
	}

	opens = 0
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"zap","parent":"store"}`)
	if opens != 1 {
		t.Errorf("a write opened the index %d times, want 1", opens)
	}
}

// TestProjectionNilWithoutIndex: with ctx.IndexPath empty the workspace's
// projection is nil, nothing opens an index and an op answers from the log.
func TestProjectionNilWithoutIndex(t *testing.T) {
	t.Parallel()
	ctx, _ := indexedLog(t)
	ctx.IndexPath = ""
	opens := 0
	ctx.openedProjection = func() { opens++ }
	w := dryRunWorkspace(t, ctx)
	if p := w.projection(); p != nil {
		t.Fatalf("projection = %T, want nil", p)
	}
	reads := countingCtx(ctx)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"zap","parent":"store"}`)
	if opens != 0 {
		t.Errorf("opened an index %d times", opens)
	}
	if *reads == 0 {
		t.Error("the write did not read the log")
	}
}

// dryAdd is a dry-run add on ctx with seeded IDs, as JSON.
func dryAdd(t *testing.T, ctx *Ctx) (string, error) {
	t.Helper()
	c := *ctx
	c.DryRun = true
	seededIDs(&c, 7)
	resp, err := callOp(t, &c, "add", `{"kind":"capability","handle":"zap","parent":"store"}`)
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	return string(b), nil
}

// TestProjectionIndexNotADatabase: an index file that is not a database
// is rebuilt and the call answers what it answers from the log.
func TestProjectionIndexNotADatabase(t *testing.T) {
	t.Parallel()
	ctx, _ := indexedLog(t)
	folded := *ctx
	folded.IndexPath = ""
	want, err := dryAdd(t, &folded)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ctx.IndexPath, []byte(strings.Repeat("not a database\n", 200)), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := dryAdd(t, ctx)
	if err != nil {
		t.Fatalf("add over a garbage index: %v", err)
	}
	if got != want {
		t.Errorf("add over a rebuilt index = %s, want %s", got, want)
	}
	indexStamp(t, ctx)
}

// TestProjectionIndexCannotOpen: an index that cannot be opened fails the
// call with the open's error; it does not fold the log instead.
func TestProjectionIndexCannotOpen(t *testing.T) {
	t.Parallel()
	ctx, _ := indexedLog(t)
	blocker := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(blocker, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	ctx.IndexPath = filepath.Join(blocker, "index.db")
	_, want := index.Open(ctx.Root, ctx.IndexPath)
	if want == nil {
		t.Fatal("index.Open under a file succeeded")
	}
	reads := countingCtx(ctx)
	_, err := dryAdd(t, ctx)
	if err == nil {
		t.Fatal("add over an index that cannot open succeeded")
	}
	if err.Error() != want.Error() {
		t.Errorf("error = %q, want the open's %q", err, want)
	}
	if *reads != 0 {
		t.Errorf("read the log %d times instead of failing", *reads)
	}
}
