package api

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// renderOK runs the render operation on ctx and fails on an error.
func renderOK(t *testing.T, ctx *Ctx) RenderResponse {
	t.Helper()
	resp, err := callOp(t, ctx, "render", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	return resp.(RenderResponse)
}

func readRendered(t *testing.T, ctx *Ctx, p string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(ctx.Root, renderedDir, filepath.FromSlash(p)))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestRenderEmpty(t *testing.T) {
	t.Parallel()
	ctx := initCtx(t)
	r := renderOK(t, ctx)
	if !slices.Equal(r.Written, []string{".rhiza/rendered/index.md"}) || len(r.Removed) != 0 || r.Unchanged != 0 {
		t.Fatalf("first render = %+v", r)
	}
	if got := readRendered(t, ctx, "index.md"); got != "# Model\n\nNothing has landed yet.\n" {
		t.Errorf("index.md = %q", got)
	}
	r = renderOK(t, ctx)
	if len(r.Written) != 0 || r.Removed == nil || len(r.Removed) != 0 || r.Unchanged != 1 {
		t.Errorf("second render = %+v", r)
	}
	if out := human(t, r); !strings.HasPrefix(out, "Rendered 1 file (0 written, 0 removed, 1 unchanged).\n") {
		t.Errorf("human:\n%s", out)
	}
}

func TestRenderWritesAndRemoves(t *testing.T) {
	t.Parallel()
	ctx := initCtx(t)
	renderOK(t, ctx)
	ctx.Record = true
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
	ctx.Record = false

	// Stray files and folders: a file render does not write, one in a
	// folder of its own, an empty folder and a leftover temp file.
	dir := filepath.Join(ctx.Root, renderedDir)
	for _, p := range []string{"extra.md", "old/deep/page.md", ".render-123"} {
		p = filepath.Join(dir, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "empty"), 0o755); err != nil {
		t.Fatal(err)
	}

	r := renderOK(t, ctx)
	wantWritten := []string{".rhiza/rendered/actor/test.md", ".rhiza/rendered/capability/store.md", ".rhiza/rendered/index.md"}
	wantRemoved := []string{".rhiza/rendered/.render-123", ".rhiza/rendered/empty",
		".rhiza/rendered/extra.md", ".rhiza/rendered/old", ".rhiza/rendered/old/deep",
		".rhiza/rendered/old/deep/page.md"}
	if !slices.Equal(r.Written, wantWritten) || !slices.Equal(r.Removed, wantRemoved) || r.Unchanged != 0 {
		t.Fatalf("render = %+v", r)
	}
	if got := readRendered(t, ctx, "capability/store.md"); !strings.HasPrefix(got, "# store\n") {
		t.Errorf("store.md = %q", got)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	if !slices.Equal(names, []string{"actor", "capability", "index.md"}) {
		t.Errorf("rendered folder holds %v", names)
	}
	out := human(t, r)
	for _, s := range []string{"Rendered 3 files (3 written, 6 removed, 0 unchanged).",
		"- wrote .rhiza/rendered/index.md", "- removed .rhiza/rendered/extra.md"} {
		if !strings.Contains(out, s) {
			t.Errorf("human lacks %q:\n%s", s, out)
		}
	}

	// A changed page is rewritten; the rest is left alone.
	if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}
	r = renderOK(t, ctx)
	if !slices.Equal(r.Written, []string{".rhiza/rendered/index.md"}) || len(r.Removed) != 0 || r.Unchanged != 2 {
		t.Errorf("third render = %+v", r)
	}
}

func TestRenderErrors(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t, t.TempDir())
	_, err := callOp(t, ctx, "render", `{}`)
	wantExit(t, err, 4)
	wantNext(t, err, "rhiza init")

	ctx = initCtx(t)
	ctx.DryRun = true
	_, err = callOp(t, ctx, "render", `{}`)
	wantExit(t, err, 2)
	wantNext(t, err, "rhiza render")
	if _, serr := os.Stat(filepath.Join(ctx.Root, renderedDir)); serr == nil {
		t.Error("a dry run wrote the rendered folder")
	}
}

// TestRenderIndexedMatchesLog: with an index, render takes the change sets
// and the as-is model from the projection and writes the files it writes
// over the folded log, byte for byte: an Evidence page names the actor that
// created it, and a link to a retired node shows it retired.
func TestRenderIndexedMatchesLog(t *testing.T) {
	t.Parallel()
	ctx, _, _ := readLog(t)
	move(t, ctx, "land", `{}`)
	newChange(t, ctx, `{"intent":"Prove it"}`)
	ctx.Record = true
	mustWrite(t, ctx, "add", `{"kind":"evidence","handle":"cache-ok","fields":["statement=works","result=PASS"]}`)
	ctx.Record = false
	mustWrite(t, ctx, "link", `{"from":"cache-ok","type":"verifies","to":"cache"}`)
	move(t, ctx, "propose", `{}`)
	move(t, ctx, "accept", `{"reason":"ok"}`)
	move(t, ctx, "land", `{}`)

	indexed := renderOK(t, ctx)
	if len(indexed.Written) == 0 {
		t.Fatalf("indexed render wrote nothing: %+v", indexed)
	}
	if got := readRendered(t, ctx, "evidence/cache-ok.md"); !strings.Contains(got, " · actor [test](") {
		t.Errorf("evidence page lacks its creator:\n%s", got)
	}
	if got := readRendered(t, ctx, "decision/use-ssd.md"); !strings.Contains(got, "disk (retired)") {
		t.Errorf("decision page lacks the retired endpoint:\n%s", got)
	}

	folded := *ctx
	// The log side: no index on purpose, to compare against the indexed render.
	folded.IndexPath = ""
	r := renderOK(t, &folded)
	if len(r.Written) != 0 || len(r.Removed) != 0 || r.Unchanged != len(indexed.Written)+indexed.Unchanged {
		t.Errorf("render over the log after the indexed one = %+v, want every file unchanged", r)
	}
}

// An entry file edited in place after the index read it is not re-read by
// catch-up, so the index still holds its old bytes. Render over the index
// must notice and write what the log fold renders, not the stale index.
func TestRenderIndexedAfterInPlaceEdit(t *testing.T) {
	t.Parallel()
	ctx, _, _ := readLog(t)
	renderOK(t, ctx)
	if got := readRendered(t, ctx, "capability/store.md"); !strings.Contains(got, "Keep things") {
		t.Fatalf("store page lacks its statement:\n%s", got)
	}

	var edited string
	err := filepath.WalkDir(filepath.Join(ctx.Root, ".rhiza", "log"), func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := os.ReadFile(p)
		if err != nil || !strings.Contains(string(b), "statement=Keep things") && !strings.Contains(string(b), `"Keep things"`) {
			return err
		}
		if err := os.Chmod(p, 0o644); err != nil {
			return err
		}
		edited = p
		return os.WriteFile(p, []byte(strings.Replace(string(b), "Keep things", "Keep other things", 1)), 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
	if edited == "" {
		t.Fatal("no entry holds the store's statement")
	}

	indexed := renderOK(t, ctx)
	if got := readRendered(t, ctx, "capability/store.md"); !strings.Contains(got, "Keep other things") {
		t.Errorf("indexed render after the edit = %+v, store page:\n%s", indexed, got)
	}
	folded := *ctx
	// The log side: no index on purpose, to compare against the indexed render.
	folded.IndexPath = ""
	if r := renderOK(t, &folded); len(r.Written) != 0 || len(r.Removed) != 0 {
		t.Errorf("render over the log after the indexed one = %+v, want every file unchanged", r)
	}
}
