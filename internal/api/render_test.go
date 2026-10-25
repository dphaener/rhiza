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
	ctx := initCtx(t)
	r := renderOK(t, ctx)
	if !slices.Equal(r.Written, []string{".genesis/rendered/index.md"}) || len(r.Removed) != 0 || r.Unchanged != 0 {
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
	wantWritten := []string{".genesis/rendered/actor/test.md", ".genesis/rendered/capability/store.md", ".genesis/rendered/index.md"}
	wantRemoved := []string{".genesis/rendered/.render-123", ".genesis/rendered/empty",
		".genesis/rendered/extra.md", ".genesis/rendered/old", ".genesis/rendered/old/deep",
		".genesis/rendered/old/deep/page.md"}
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
		"- wrote .genesis/rendered/index.md", "- removed .genesis/rendered/extra.md"} {
		if !strings.Contains(out, s) {
			t.Errorf("human lacks %q:\n%s", s, out)
		}
	}

	// A changed page is rewritten; the rest is left alone.
	if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}
	r = renderOK(t, ctx)
	if !slices.Equal(r.Written, []string{".genesis/rendered/index.md"}) || len(r.Removed) != 0 || r.Unchanged != 2 {
		t.Errorf("third render = %+v", r)
	}
}

func TestRenderErrors(t *testing.T) {
	ctx := testCtx(t.TempDir())
	_, err := callOp(t, ctx, "render", `{}`)
	wantExit(t, err, 4)
	wantNext(t, err, "genesis init")

	ctx = initCtx(t)
	ctx.DryRun = true
	_, err = callOp(t, ctx, "render", `{}`)
	wantExit(t, err, 2)
	wantNext(t, err, "genesis render")
	if _, serr := os.Stat(filepath.Join(ctx.Root, renderedDir)); serr == nil {
		t.Error("a dry run wrote the rendered folder")
	}
}
