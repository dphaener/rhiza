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
	infos := p.changeSetInfos(w.pending)
	if i := slices.IndexFunc(infos, func(c index.ChangeSetInfo) bool { return c.ID == cs }); i < 0 ||
		infos[i].State != changeset.StateDraft || infos[i].Intent != "Pending" || infos[i].Kind != changeset.KindIntention {
		t.Errorf("changeSetInfos lacks the pending change set %s as a draft: %+v", cs, infos)
	}
	if got, want := len(infos), len(p.changeSetIDs(w.pending)); got != want ||
		!slices.IsSortedFunc(infos, func(a, b index.ChangeSetInfo) int { return id.Compare(a.ID, b.ID) }) {
		t.Errorf("changeSetInfos = %d rows (sorted: %v), want %d sorted", got,
			slices.IsSortedFunc(infos, func(a, b index.ChangeSetInfo) int { return id.Compare(a.ID, b.ID) }), want)
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

// withTempDir makes ctx's temporary folder ($TMPDIR) dir.
func withTempDir(ctx *Ctx, dir string) {
	prev := ctx.Getenv
	ctx.Getenv = func(k string) string {
		if k == "TMPDIR" {
			return dir
		}
		if prev != nil {
			return prev(k)
		}
		return ""
	}
}

// projectionAnswers is a status, a show and a dry-run add on ctx, as JSON.
func projectionAnswers(t *testing.T, ctx *Ctx) []string {
	t.Helper()
	var out []string
	for _, c := range []struct{ op, req string }{{"status", `{}`}, {"show", `{"ref":"ssd"}`}} {
		resp, err := callOp(t, ctx, c.op, c.req)
		if err != nil {
			t.Fatalf("%s: %v", c.op, err)
		}
		b, err := json.Marshal(resp)
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, string(b))
	}
	add, err := dryAdd(t, ctx)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	return append(out, add)
}

// emptyDir fails the test unless dir holds nothing.
func emptyDir(t *testing.T, dir string) {
	t.Helper()
	left, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range left {
		t.Errorf("left behind in the temporary folder: %s", e.Name())
	}
}

// throwawayMatches checks that a status, a show and an add on ctx (whose
// index path the caller has made empty or unusable) open a throwaway
// index, answer what they answer over want's index and leave nothing in
// the temporary folder.
func throwawayMatches(t *testing.T, ctx *Ctx, want []string) {
	t.Helper()
	tmp := t.TempDir()
	withTempDir(ctx, tmp)
	opens := 0
	ctx.openedProjection = func() { opens++ }
	w := dryRunWorkspace(t, ctx)
	if p := w.projection(); p == nil {
		t.Fatal("no projection")
	} else if err := p.failure(); err != nil {
		t.Fatalf("throwaway projection failed: %v", err)
	}
	w.closeProjection()
	emptyDir(t, tmp)
	got := projectionAnswers(t, ctx)
	if opens < 2 {
		t.Errorf("opened a projection %d times, want the add's too", opens)
	}
	for i, op := range []string{"status", "show", "add"} {
		if got[i] != want[i] {
			t.Errorf("%s over a throwaway index = %s, want %s", op, got[i], want[i])
		}
	}
	emptyDir(t, tmp)
}

// TestProjectionThrowawayWithoutIndex: with ctx.IndexPath empty the
// workspace's projection is a throwaway index in the temporary folder: a
// status, a show and an add answer what they answer with an index, and the
// folder is removed when the projection closes.
func TestProjectionThrowawayWithoutIndex(t *testing.T) {
	t.Parallel()
	ctx, _ := indexedLog(t)
	want := projectionAnswers(t, ctx)
	ctx.IndexPath = ""
	throwawayMatches(t, ctx, want)
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

// TestProjectionIndexCannotOpen: an index that cannot be created (its
// folder is read-only) gives the call a throwaway index that answers what
// the index answers.
func TestProjectionIndexCannotOpen(t *testing.T) {
	t.Parallel()
	if os.Getuid() == 0 {
		t.Skip("root writes into a read-only folder")
	}
	ctx, _ := indexedLog(t)
	want := projectionAnswers(t, ctx)
	ro := t.TempDir()
	if err := os.Chmod(ro, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(ro, 0o755) })
	ctx.IndexPath = filepath.Join(ro, "index.db")
	if _, err := index.Open(ctx.Root, ctx.IndexPath); err == nil {
		t.Fatal("index.Open in a read-only folder succeeded")
	}
	throwawayMatches(t, ctx, want)
}

// TestProjectionNoTempDir: when neither the index nor a temporary folder
// can be created, the call fails with an error naming the index path and
// the temporary folder.
func TestProjectionNoTempDir(t *testing.T) {
	t.Parallel()
	if os.Getuid() == 0 {
		t.Skip("root writes into a read-only folder")
	}
	ctx, _ := indexedLog(t)
	ro, tmp := t.TempDir(), t.TempDir()
	for _, d := range []string{ro, tmp} {
		if err := os.Chmod(d, 0o555); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { os.Chmod(d, 0o755) })
	}
	ctx.IndexPath = filepath.Join(ro, "index.db")
	withTempDir(ctx, tmp)
	_, err := dryAdd(t, ctx)
	if err == nil {
		t.Fatal("add with no index and no temporary folder succeeded")
	}
	for _, s := range []string{ctx.IndexPath, tmp} {
		if !strings.Contains(err.Error(), s) {
			t.Errorf("error %q does not name %s", err, s)
		}
	}
}

// TestProjectionNoneGiven: a workspace opened with openWorkspace and given
// no projection opens none, not a throwaway index over a bare Ctx: its
// reads answer nothing, projectionFailure names the missing projection,
// and no rhiza-projection-* folder is made in the Ctx's TMPDIR.
func TestProjectionNoneGiven(t *testing.T) {
	t.Parallel()
	ctx, _ := indexedLog(t)
	ctx.IndexPath = ""
	tmp := t.TempDir()
	withTempDir(ctx, tmp)
	w, err := openWorkspace(ctx.Root, ctx.Now, ctx.IDs, false)
	if err != nil {
		t.Fatal(err)
	}
	if got := w.changeSets(); len(got) != 0 {
		t.Errorf("changeSets with no projection = %d change sets, want none", len(got))
	}
	err = w.projectionFailure()
	if err == nil || !strings.Contains(err.Error(), "no projection") {
		t.Errorf("projectionFailure = %v, want an error naming the missing projection", err)
	}
	if _, ok := w.proj.p.(failedProjection); !ok {
		t.Errorf("projection = %T, want a failed projection opening no index", w.proj.p)
	}
	w.closeProjection()
	left, err := filepath.Glob(filepath.Join(tmp, "rhiza-projection-*"))
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("throwaway indexes made with no projection: %v", left)
	}
}
