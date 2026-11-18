package api

import (
	"errors"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/store"
)

// testWorkspace makes a repository in a temp folder and opens it with a
// fixed clock and seeded IDs.
func testWorkspace(t *testing.T, dryRun bool) *workspace {
	t.Helper()
	now := time.Date(2026, 10, 24, 12, 0, 0, 500, time.FixedZone("x", 3600))
	gen := id.NewGenerator(func() time.Time { return now }, rand.New(rand.NewSource(7)))
	root := t.TempDir()
	if _, err := store.Init(root, gen); err != nil {
		t.Fatal(err)
	}
	w, err := openWorkspace(root, func() time.Time { return now }, gen, dryRun)
	if err != nil {
		t.Fatal(err)
	}
	// No Ctx here: one over the root, clock and IDs whose TMPDIR is t's
	// own, so the throwaway index lands under t and is closed with it.
	tmp := t.TempDir()
	w.useProjection(&Ctx{Root: root, Now: w.now, IDs: gen, Getenv: func(k string) string {
		if k == "TMPDIR" {
			return tmp
		}
		return testCache(k)
	}})
	t.Cleanup(w.closeProjection)
	return w
}

// gitConfig fakes `git config` with a fixed map; a missing key is an error.
func gitConfig(m map[string]string) func(string) (string, error) {
	return func(key string) (string, error) {
		if v, ok := m[key]; ok {
			return v, nil
		}
		return "", errors.New("unset")
	}
}

// logFiles lists the files under .rhiza/log, relative to it, skipping
// dot-files (init's .gitkeep) as the log reader does.
func logFiles(t *testing.T, root string) []string {
	t.Helper()
	var out []string
	dir := filepath.Join(root, store.Dir, store.LogDir)
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && !strings.HasPrefix(d.Name(), ".") {
			rel, _ := filepath.Rel(dir, p)
			out = append(out, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func mustFoldClean(t *testing.T, w *workspace) []changeset.ChangeSet {
	t.Helper()
	l := w.st.ReadLog()
	if len(l.Problems) > 0 {
		t.Fatalf("log problems: %v", l.Problems)
	}
	sets, problems := changeset.FoldLog(l)
	if len(problems) > 0 {
		t.Fatalf("fold problems: %v", problems)
	}
	return sets
}

func TestWorkspace(t *testing.T) {
	t.Parallel()
	t.Run("CurrentRoundTrip", func(t *testing.T) {
		w := testWorkspace(t, false)
		got, err := w.current()
		if err != nil || !got.IsZero() {
			t.Fatalf("initial current %s, %v; want zero, nil", got, err)
		}
		cs, _ := w.ids.New()
		if err := w.setCurrent(cs); err != nil {
			t.Fatal(err)
		}
		if got, err := w.current(); err != nil || got != cs {
			t.Fatalf("current %s, %v; want %s", got, err, cs)
		}
		data, err := os.ReadFile(filepath.Join(w.root, store.Dir, "local", "change"))
		if err != nil || string(data) != cs.String()+"\n" {
			t.Fatalf("file %q, %v", data, err)
		}
		if err := w.clearCurrent(); err != nil {
			t.Fatal(err)
		}
		if got, err := w.current(); err != nil || !got.IsZero() {
			t.Fatalf("after clear %s, %v", got, err)
		}
		if err := w.clearCurrent(); err != nil {
			t.Fatalf("second clear: %v", err)
		}
		if err := os.WriteFile(filepath.Join(w.root, store.Dir, "local", "change"), []byte("junk\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := w.current(); ExitCode(err) != 1 {
			t.Fatalf("junk current: %v, want exit 1", err)
		}
	})

	t.Run("FindRootFromSubfolder", func(t *testing.T) {
		w := testWorkspace(t, false)
		sub := filepath.Join(w.root, "a", "b")
		if err := os.MkdirAll(sub, 0o755); err != nil {
			t.Fatal(err)
		}
		got, ok := FindRoot(sub)
		if !ok || got != w.root {
			t.Fatalf("FindRoot = %q, %v; want %q", got, ok, w.root)
		}
		if got, ok := FindRoot(t.TempDir()); ok {
			t.Fatalf("FindRoot outside a repo = %q, true", got)
		}
		if _, err := openWorkspace(t.TempDir(), time.Now, id.Default(), false); ExitCode(err) != 4 ||
			AsError(err).Next != "rhiza init" {
			t.Fatalf("open without repo: %v", err)
		}
	})

	t.Run("FindRootOldFolder", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, ".genesis"), 0o755); err != nil {
			t.Fatal(err)
		}
		sub := filepath.Join(root, "src")
		if err := os.Mkdir(sub, 0o755); err != nil {
			t.Fatal(err)
		}
		for _, d := range []string{root, sub} {
			if got, ok := FindRoot(d); !ok || got != root {
				t.Fatalf("FindRoot(%s) = %q, %v; want the old-style root %q", d, got, ok, root)
			}
		}
		// A nearer .rhiza wins over an older .genesis above it.
		inner := filepath.Join(sub, "pkg")
		if err := os.MkdirAll(filepath.Join(inner, ".rhiza"), 0o755); err != nil {
			t.Fatal(err)
		}
		if got, ok := FindRoot(inner); !ok || got != inner {
			t.Fatalf("FindRoot(inner) = %q, %v; want %q", got, ok, inner)
		}
		// Opening the old-style root is a conflict naming the rename.
		_, err := openWorkspace(root, time.Now, id.Default(), false)
		if e := AsError(err); ExitCode(err) != 3 || e.Code != CodeConflict || e.Next != "git mv .genesis .rhiza" ||
			!strings.Contains(e.Message, ".genesis/") || !isOldFolderError(err) {
			t.Fatalf("open old root: %v", err)
		}
	})

	t.Run("ResolveChangeSetPrefix", func(t *testing.T) {
		w := testWorkspace(t, false)
		var sets []id.ID
		for range 2 {
			cs, _, err := w.openChangeSet(changeset.Open{Kind: changeset.KindIntention, Intent: "x"}, w.mustID(t), w.mustID(t))
			if err != nil {
				t.Fatal(err)
			}
			sets = append(sets, cs)
		}
		// Same clock, same millisecond: the IDs share a long prefix.
		full := sets[1].String()
		common := 0
		for common < len(full) && sets[0].String()[common] == full[common] {
			common++
		}
		if common < 4 {
			t.Fatalf("test needs a shared prefix, got %d", common)
		}
		got, err := w.resolveChangeSet(full[:common+1])
		if err != nil || got != sets[1] {
			t.Fatalf("unique prefix: %s, %v", got, err)
		}
		if got, err := w.resolveChangeSet(strings.ToLower(full)); err != nil || got != sets[1] {
			t.Fatalf("full lower-case id: %s, %v", got, err)
		}
		_, err = w.resolveChangeSet(full[:common])
		if ExitCode(err) != 2 || !strings.Contains(err.Error(), sets[0].String()) {
			t.Fatalf("ambiguous: %v", err)
		}
		if _, err := w.resolveChangeSet("ZZZZZZ"); ExitCode(err) != 4 {
			t.Fatalf("unknown: %v", err)
		}
		if _, err := w.resolveChangeSet("01"); ExitCode(err) != 2 {
			t.Fatalf("too short: %v", err)
		}

		// Effective: the ref wins, then current, then none.
		if got, err := w.effectiveChangeSet(&Ctx{}); err != nil || !got.IsZero() {
			t.Fatalf("effective with nothing: %s, %v", got, err)
		}
		if err := w.setCurrent(sets[0]); err != nil {
			t.Fatal(err)
		}
		if got, err := w.effectiveChangeSet(&Ctx{}); err != nil || got != sets[0] {
			t.Fatalf("effective current: %s, %v", got, err)
		}
		if got, err := w.effectiveChangeSet(&Ctx{ChangeRef: full}); err != nil || got != sets[1] {
			t.Fatalf("effective ref: %s, %v", got, err)
		}
	})

	t.Run("ActorRecordedOnce", func(t *testing.T) {
		w := testWorkspace(t, false)
		git := gitConfig(map[string]string{"user.email": "Ada.Lovelace@example.com", "user.name": "Ada"})
		actor, forID, err := resolveRefs(w, "", "", git)
		if err != nil {
			t.Fatal(err)
		}
		if actor.IsZero() || forID != actor {
			t.Fatalf("actor %s for %s", actor, forID)
		}
		files := logFiles(t, w.root)
		if len(files) != 3 {
			t.Fatalf("log files %v, want 3", files)
		}
		sets := mustFoldClean(t, w)
		if len(sets) != 1 || sets[0].State != changeset.StateLanded || sets[0].Open.Kind != changeset.KindRecord ||
			sets[0].Open.Intent != "Record actor Ada.Lovelace@example.com" {
			t.Fatalf("change sets %+v", sets)
		}
		for _, e := range w.st.ReadLog().Entries {
			if e.Actor != actor || e.For != actor {
				t.Errorf("entry %s actor %s for %s, want %s", e.ID, e.Actor, e.For, actor)
			}
		}
		n, ok := w.model(id.ID{}).Node(actor)
		if !ok || n.Kind != op.KindActor || n.Handle != "ada-lovelace" ||
			string(n.Fields["email"]) != `"Ada.Lovelace@example.com"` || string(n.Fields["name"]) != `"Ada"` {
			t.Fatalf("actor node %+v", n)
		}

		again, forAgain, err := resolveRefs(w, "", "", git)
		if err != nil || again != actor || forAgain != actor {
			t.Fatalf("second call %s %s %v", again, forAgain, err)
		}
		if got := logFiles(t, w.root); !slices.Equal(got, files) {
			t.Fatalf("second call wrote: %v", got)
		}

		// An explicit ref and a for ref resolve through the model.
		other, _, err := resolveRefs(w, "", "", gitConfig(map[string]string{"user.email": "bob@example.com"}))
		if err != nil {
			t.Fatal(err)
		}
		a, f, err := resolveRefs(w, "ada-lovelace", "bob", git)
		if err != nil || a != actor || f != other {
			t.Fatalf("explicit refs: %s %s %v; want %s %s", a, f, err, actor, other)
		}
		if _, ok := w.model(id.ID{}).Node(other); !ok {
			t.Fatal("bob not recorded")
		}
	})

	t.Run("ActorWrongKind", func(t *testing.T) {
		w := testWorkspace(t, false)
		git := gitConfig(map[string]string{"user.email": "ada@example.com"})
		actor, _, err := resolveRefs(w, "", "", git)
		if err != nil {
			t.Fatal(err)
		}
		c := op.Op{ID: w.mustID(t), Type: op.TypeNodeCreate, Node: w.mustID(t), Kind: op.KindDecision, Handle: "use-go"}
		if _, _, err := w.writeRecord("Decide", []op.Op{c}, actor, actor); err != nil {
			t.Fatal(err)
		}
		_, _, err = resolveRefs(w, "use-go", "", git)
		if ExitCode(err) != 2 || !strings.Contains(err.Error(), "is a decision, not an actor") {
			t.Fatalf("wrong kind actor: %v", err)
		}
		_, _, err = resolveRefs(w, "", "use-go", git)
		if ExitCode(err) != 2 {
			t.Fatalf("wrong kind for: %v", err)
		}
		if _, _, err := resolveRefs(w, "nobody", "", git); ExitCode(err) != 4 {
			t.Fatalf("unknown actor: %v", err)
		}
	})

	t.Run("NoIdentity", func(t *testing.T) {
		w := testWorkspace(t, false)
		_, _, err := resolveRefs(w, "", "", gitConfig(nil))
		e := AsError(err)
		if ExitCode(err) != 2 || !strings.Contains(e.Next, "git config user.email") {
			t.Fatalf("no identity: %+v", e)
		}
		if files := logFiles(t, w.root); len(files) != 0 {
			t.Fatalf("wrote %v", files)
		}
	})

	t.Run("AfterIsHeads", func(t *testing.T) {
		w := testWorkspace(t, false)
		a := w.mustID(t)
		cs, open, err := w.openChangeSet(changeset.Open{Kind: changeset.KindIntention, Intent: "Plan"}, a, a)
		if err != nil {
			t.Fatal(err)
		}
		if len(open.After) != 0 || open.Type != entry.TypeOpen || open.ChangeSet != cs || open.ID == cs {
			t.Fatalf("open %+v", open)
		}
		if !open.At.Equal(time.Date(2026, 10, 24, 11, 0, 0, 0, time.UTC)) || open.At.Location() != time.UTC {
			t.Fatalf("at %v", open.At)
		}
		c := op.Op{ID: w.mustID(t), Type: op.TypeNodeCreate, Node: w.mustID(t), Kind: op.KindQuestion, Handle: "why"}
		ops1, err := w.appendEntry(cs, changeset.Ops{Ops: []op.Op{c}}, a, a)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(ops1.After, []id.ID{open.ID}) {
			t.Fatalf("ops after %v, want [%s]", ops1.After, open.ID)
		}
		c2 := op.Op{ID: w.mustID(t), Type: op.TypeNodeCreate, Node: w.mustID(t), Kind: op.KindQuestion, Handle: "how"}
		ops2, err := w.appendEntry(cs, changeset.Ops{Ops: []op.Op{c2}}, a, a)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(ops2.After, []id.ID{ops1.ID}) || ops2.V != 1 {
			t.Fatalf("second ops %+v", ops2)
		}
		mustFoldClean(t, w)

		// A terminal change set takes no more ops.
		if _, err := w.appendEntry(cs, changeset.Withdraw{}, a, a); err != nil {
			t.Fatal(err)
		}
		_, err = w.appendEntry(cs, changeset.Ops{Ops: []op.Op{c2}}, a, a)
		if ExitCode(err) != 3 || !strings.Contains(AsError(err).Next, "rhiza change new") {
			t.Fatalf("ops into withdrawn: %v", err)
		}
		missing := w.mustID(t)
		if _, err := w.appendEntry(missing, changeset.Land{}, a, a); ExitCode(err) != 4 {
			t.Fatalf("unknown change set: %v", err)
		}
	})

	t.Run("DryRunWritesNothing", func(t *testing.T) {
		w := testWorkspace(t, true)
		// No actor node: a dry run records nothing and returns zero.
		actor, forID, err := resolveRefs(w, "", "", gitConfig(map[string]string{"user.email": "ada@example.com"}))
		if err != nil || !actor.IsZero() || !forID.IsZero() {
			t.Fatalf("dry-run actor %s %s %v", actor, forID, err)
		}
		a := w.mustID(t)
		c := op.Op{ID: w.mustID(t), Type: op.TypeNodeCreate, Node: w.mustID(t), Kind: op.KindQuestion, Handle: "why"}
		cs, entries, err := w.writeRecord("Ask", []op.Op{c}, a, a)
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 3 || entries[0].ChangeSet != cs ||
			!slices.Equal(entries[1].After, []id.ID{entries[0].ID}) ||
			!slices.Equal(entries[2].After, []id.ID{entries[1].ID}) {
			t.Fatalf("dry-run entries %+v", entries)
		}
		if err := w.setCurrent(cs); err != nil {
			t.Fatal(err)
		}
		if files := logFiles(t, w.root); len(files) != 0 {
			t.Fatalf("dry run wrote %v", files)
		}
		if _, err := os.Stat(filepath.Join(w.root, store.Dir, "local")); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("dry run made local/: %v", err)
		}
	})

	t.Run("HandleSlugUnique", func(t *testing.T) {
		w := testWorkspace(t, false)
		var handles []string
		for _, email := range []string{"x@a.example", "X@b.example", "x@c.example"} {
			a, _, err := resolveRefs(w, "", "", gitConfig(map[string]string{"user.email": email}))
			if err != nil {
				t.Fatal(err)
			}
			n, _ := w.model(id.ID{}).Node(a)
			handles = append(handles, n.Handle)
		}
		if want := []string{"x", "x-2", "x-3"}; !slices.Equal(handles, want) {
			t.Fatalf("handles %v, want %v", handles, want)
		}
		mustFoldClean(t, w)

		for in, want := range map[string]string{
			"Ada.Lovelace":                 "ada-lovelace",
			"--a__b--":                     "a-b",
			"...":                          "actor",
			"":                             "actor",
			strings.Repeat("ab", 40):       strings.Repeat("ab", 32),
			strings.Repeat("a", 63) + ".b": strings.Repeat("a", 63),
		} {
			if got := handleSlug(in); got != want {
				t.Errorf("handleSlug(%q) = %q, want %q", in, got, want)
			}
		}
		taken := func(h string) bool { return h == strings.Repeat("a", 64) }
		if got := uniqueHandle(strings.Repeat("a", 64), taken); got != strings.Repeat("a", 62)+"-2" {
			t.Errorf("uniqueHandle at 64: %q", got)
		}
	})
}

func (w *workspace) mustID(t *testing.T) id.ID {
	t.Helper()
	v, err := w.ids.New()
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// resolveRefs resolves the actor with no current change set, choosing by
// actorRef and forRef over git.
func resolveRefs(w *workspace, actorRef, forRef string, git func(string) (string, error)) (id.ID, id.ID, error) {
	return w.resolveActor(id.ID{}, &Ctx{ActorRef: actorRef, ForRef: forRef, GitConfig: git})
}

// logOnceWorkspace opens a workspace over ctx's root with no index path (a
// throwaway index) whose disk reads of the log are counted through
// Ctx.logRead.
func logOnceWorkspace(t *testing.T, ctx *Ctx, dryRun bool) (*workspace, *int) {
	t.Helper()
	ctx.IndexPath = ""
	reads := countingCtx(ctx)
	rc := *ctx
	rc.DryRun = dryRun
	w, err := openCtxWorkspace(&rc)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(w.closeProjection)
	return w, reads
}

// childHandles is the handles of every capability under parent in m, depth
// first: what two models are compared by.
func childHandles(m Model, parent id.ID) []string {
	var out []string
	for _, n := range m.Children(op.KindCapability, parent) {
		out = append(out, n.Handle)
		out = append(out, childHandles(m, n.ID)...)
	}
	return out
}

func TestWorkspaceLogReadsDiskOnce(t *testing.T) {
	t.Parallel()
	ctx, landed, accepted := readLog(t)
	w, reads := logOnceWorkspace(t, ctx, false)
	fresh := ctxWorkspace(t, ctx)
	defer fresh.closeProjection()
	for i := range 3 {
		got, want := childHandles(w.model(accepted), id.ID{}), childHandles(fresh.model(accepted), id.ID{})
		if !slices.Equal(got, want) || len(got) == 0 {
			t.Fatalf("call %d: model %v differs from a fresh read's %v", i, got, want)
		}
		if got, want := w.changeSets(), fresh.changeSets(); !reflect.DeepEqual(got, want) || len(got) == 0 {
			t.Fatalf("call %d: %d change sets, differ from a fresh read: %v", i, len(got), !reflect.DeepEqual(got, want))
		}
		if got, want := w.entriesOfSet(landed), fresh.entriesOfSet(landed); !reflect.DeepEqual(got, want) || len(got) == 0 {
			t.Fatalf("call %d: entries of the landed set differ from a fresh read", i)
		}
	}
	// The model, change sets and entries come from the projection: the
	// log is never read from disk.
	if *reads != 0 {
		t.Errorf("disk reads = %d, want 0", *reads)
	}
	// A caller changing its copy of the log leaves the kept read as it was.
	l := w.log()
	l.Entries[0] = entry.Entry{}
	l.Entries = l.Entries[:0]
	if got := w.log(); len(got.Entries) == 0 || got.Entries[0].ID.IsZero() {
		t.Error("changing a returned log changed the kept read")
	}
	if *reads != 1 {
		t.Errorf("disk reads after two log calls = %d, want 1", *reads)
	}
}

func TestWorkspaceLogDryRunChains(t *testing.T) {
	t.Parallel()
	ctx, landed, _ := readLog(t)
	before := logFiles(t, ctx.Root)
	w, reads := logOnceWorkspace(t, ctx, true)
	n := len(w.changeSets())
	var opened []id.ID
	for i := range 3 {
		cs, _, err := w.openChangeSet(changeset.Open{Kind: changeset.KindIntention, Intent: "Pending"}, id.ID{}, id.ID{})
		if err != nil {
			t.Fatal(err)
		}
		opened = append(opened, cs)
		// Each later read sees every earlier pending entry, laid over the
		// projection.
		if got := len(w.changeSets()); got != n+i+1 {
			t.Fatalf("after %d dry writes: %d change sets, want %d", i+1, got, n+i+1)
		}
		for _, o := range opened {
			if _, ok := w.changeSet(o); !ok {
				t.Fatalf("after %d dry writes: pending change set %s not read", i+1, o)
			}
			if got, err := w.resolveChangeSet(o.String()); err != nil || got != o {
				t.Fatalf("resolve pending %s: %s, %v", o, got, err)
			}
		}
		if len(w.entriesOfSet(landed)) == 0 {
			t.Fatal("the landed set's entries are gone")
		}
	}
	if *reads != 0 {
		t.Errorf("disk reads = %d, want 0", *reads)
	}
	if after := logFiles(t, ctx.Root); !slices.Equal(before, after) {
		t.Errorf("a dry run wrote to the log: %d files, then %d", len(before), len(after))
	}
}

func TestWorkspaceLogWriteDropsRead(t *testing.T) {
	t.Parallel()
	ctx, _, _ := readLog(t)
	w, reads := logOnceWorkspace(t, ctx, false)
	n := len(w.changeSets())
	first := w.entriesOfSet(w.changeSets()[0].ID)[0]
	cs, e, err := w.openChangeSet(changeset.Open{Kind: changeset.KindIntention, Intent: "Written"}, first.Actor, first.For)
	if err != nil {
		t.Fatal(err)
	}
	// The projection catches up with the entry written.
	if got := len(w.changeSets()); got != n+1 {
		t.Fatalf("after a write: %d change sets, want %d", got, n+1)
	}
	es := w.entriesOfSet(cs)
	if len(es) != 1 || es[0].ID != e.ID {
		t.Fatalf("entries of the written set = %v, want its open entry %s", es, e.ID)
	}
	w.model(cs)
	if *reads != 0 {
		t.Errorf("disk reads = %d, want 0", *reads)
	}
}

func TestWorkspaceLogDamagedEntry(t *testing.T) {
	t.Parallel()
	ctx, landed, _ := readLog(t)
	bad, _ := ctx.IDs.New()
	path := filepath.Join(ctx.Root, store.Dir, store.LogDir, landed.String(), bad.String()+".json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	w, reads := logOnceWorkspace(t, ctx, false)
	first := w.log()
	if len(first.Problems) == 0 {
		t.Fatal("a damaged entry gave no problems")
	}
	firstSets := w.changeSets()
	second := w.log()
	if !reflect.DeepEqual(first.Problems, second.Problems) {
		t.Errorf("problems differ: %v, then %v", first.Problems, second.Problems)
	}
	if !reflect.DeepEqual(firstSets, w.changeSets()) {
		t.Error("change sets differ between calls")
	}
	if fresh := ctxWorkspace(t, ctx).log(); !reflect.DeepEqual(fresh.Problems, second.Problems) {
		t.Errorf("problems %v, a fresh read's %v", second.Problems, fresh.Problems)
	}
	if *reads != 1 {
		t.Errorf("disk reads = %d, want 1", *reads)
	}
}
