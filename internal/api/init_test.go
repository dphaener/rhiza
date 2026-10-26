package api

import (
	"bytes"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/store"
)

// testCtx is a Ctx over root with a fixed clock, seeded IDs and a git
// config whose user.email is test@example.com.
func testCtx(root string) *Ctx {
	now := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	clock := func() time.Time { return now }
	return &Ctx{
		Root:      root,
		Now:       clock,
		IDs:       id.NewGenerator(clock, rand.New(rand.NewSource(3))),
		GitConfig: gitConfig(map[string]string{"user.email": "test@example.com", "user.name": "Test User"}),
	}
}

// callOp decodes reqJSON for the named operation and runs it on ctx.
func callOp(t *testing.T, ctx *Ctx, name, reqJSON string) (any, error) {
	t.Helper()
	o, ok := Lookup(name)
	if !ok {
		t.Fatalf("%s is not registered", name)
	}
	req, err := DecodeRequest(o, []byte(reqJSON))
	if err != nil {
		t.Fatalf("DecodeRequest(%s, %s): %v", name, reqJSON, err)
	}
	return o.Run(ctx, req)
}

// human renders a response's markdown.
func human(t *testing.T, resp any) string {
	t.Helper()
	h, ok := resp.(Humaner)
	if !ok {
		t.Fatalf("%T is not a Humaner", resp)
	}
	var buf bytes.Buffer
	h.Human(&buf)
	return buf.String()
}

func TestInit(t *testing.T) {
	t.Run("Fresh", func(t *testing.T) {
		root := t.TempDir()
		resp, err := callOp(t, testCtx(root), "init", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		r := resp.(InitResponse)
		st, err := store.Open(root)
		if err != nil {
			t.Fatal(err)
		}
		if r.Repo != st.Repo().ID || r.Repo.IsZero() {
			t.Errorf("repo = %s, want %s", r.Repo, st.Repo().ID)
		}
		if r.Root != root {
			t.Errorf("root = %q, want %q", r.Root, root)
		}
		data, err := os.ReadFile(filepath.Join(root, store.Dir, ".gitignore"))
		if err != nil || string(data) != "local/\n" {
			t.Errorf(".gitignore = %q, %v; want \"local/\\n\"", data, err)
		}
		if data, err := os.ReadFile(filepath.Join(root, store.Dir, store.LogDir, ".gitkeep")); err != nil || len(data) != 0 {
			t.Errorf("log/.gitkeep = %q, %v; want an empty file", data, err)
		}
		out := human(t, resp)
		for _, s := range []string{root, r.Repo.String(), `rhiza change new "<intent>"`} {
			if !strings.Contains(out, s) {
				t.Errorf("human output lacks %q:\n%s", s, out)
			}
		}
	})

	t.Run("Exists", func(t *testing.T) {
		root := t.TempDir()
		if _, err := callOp(t, testCtx(root), "init", `{}`); err != nil {
			t.Fatal(err)
		}
		_, err := callOp(t, testCtx(root), "init", `{}`)
		if ExitCode(err) != 3 {
			t.Fatalf("second init: %v, exit %d; want 3", err, ExitCode(err))
		}
		if !strings.Contains(err.Error(), "already") {
			t.Errorf("message = %q", err)
		}
	})

	t.Run("RelativeRootIsMadeAbsolute", func(t *testing.T) {
		root := t.TempDir()
		t.Chdir(root)
		resp, err := callOp(t, testCtx("."), "init", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		if got := resp.(InitResponse).Root; !filepath.IsAbs(got) {
			t.Errorf("root = %q, want absolute", got)
		}
	})

	t.Run("DryRunIsUsage", func(t *testing.T) {
		root := t.TempDir()
		ctx := testCtx(root)
		ctx.DryRun = true
		_, err := callOp(t, ctx, "init", `{}`)
		if ExitCode(err) != 2 {
			t.Fatalf("dry-run init: %v, exit %d; want 2", err, ExitCode(err))
		}
		if _, err := os.Stat(filepath.Join(root, store.Dir)); !os.IsNotExist(err) {
			t.Errorf("dry-run init created .rhiza: %v", err)
		}
	})
}
