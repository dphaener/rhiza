package api

import (
	"errors"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
	"github.com/dphaener/genesis/internal/store"
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

// logFiles lists the files under .genesis/log, relative to it, skipping
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
			AsError(err).Next != "genesis init" {
			t.Fatalf("open without repo: %v", err)
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
		if ExitCode(err) != 3 || !strings.Contains(AsError(err).Next, "genesis change new") {
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
