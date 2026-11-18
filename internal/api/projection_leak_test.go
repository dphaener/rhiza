package api

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// testTmp is every test Ctx's TMPDIR (testCache answers it): a
// rhiza-api-tmp- folder TestMain makes once, checks with tmpLeakCode after
// m.Run and removes.
var testTmp string

// leftProjections names the rhiza-projection-* folders directly under dir:
// throwaway indexes no one removed.
func leftProjections(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var left []string
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "rhiza-projection-") {
			left = append(left, e.Name())
		}
	}
	sort.Strings(left)
	return left, nil
}

// tmpLeakCode is the test binary's exit code once m.Run answered code: code,
// or 1 when a throwaway index's folder is left in dir (each named on w) or
// dir cannot be read.
func tmpLeakCode(code int, dir string, w io.Writer) int {
	left, err := leftProjections(dir)
	if err != nil {
		fmt.Fprintf(w, "checking for throwaway indexes left: %v\n", err)
		return 1
	}
	if len(left) == 0 {
		return code
	}
	fmt.Fprintf(w, "throwaway indexes left in %s (a projection opened and never closed):\n", dir)
	for _, name := range left {
		fmt.Fprintf(w, "\t%s\n", name)
	}
	return 1
}

func TestTmpLeakReportsLeftProjection(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "rhiza-projection-x"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "rhiza-projection-file"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "rhiza-run-1"), 0o755); err != nil {
		t.Fatal(err)
	}
	left, err := leftProjections(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) != 1 || left[0] != "rhiza-projection-x" {
		t.Fatalf("leftProjections = %v, want [rhiza-projection-x]", left)
	}
	var out bytes.Buffer
	if code := tmpLeakCode(0, dir, &out); code == 0 {
		t.Fatalf("tmpLeakCode = 0 with a projection folder left; output %q", out.String())
	}
	if !strings.Contains(out.String(), "rhiza-projection-x") {
		t.Fatalf("output %q does not name rhiza-projection-x", out.String())
	}
}

func TestTmpLeakEmptyFolderPasses(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	var out bytes.Buffer
	if code := tmpLeakCode(0, dir, &out); code != 0 || out.Len() != 0 {
		t.Fatalf("tmpLeakCode(empty) = %d, output %q; want 0 and nothing", code, out.String())
	}
	// A failed run's code is kept.
	if code := tmpLeakCode(3, dir, &out); code != 3 {
		t.Fatalf("tmpLeakCode(3, empty) = %d, want 3", code)
	}
}

// testCache answers the run's own folder for TMPDIR, and the user cache for
// XDG_CACHE_HOME as before; a snapshot build's Ctx (testCtxNoIndex) puts its
// throwaway index there, never in the host's /tmp.
func TestTmpLeakTestCacheAnswersPackageTmp(t *testing.T) {
	t.Parallel()
	if testTmp == "" || !strings.HasPrefix(filepath.Base(testTmp), "rhiza-api-tmp-") {
		t.Fatalf("testTmp = %q, want a rhiza-api-tmp- folder", testTmp)
	}
	if got := testCache("TMPDIR"); got != testTmp {
		t.Fatalf("testCache(TMPDIR) = %q, want %q", got, testTmp)
	}
	if got, want := testCache("XDG_CACHE_HOME"), os.Getenv("XDG_CACHE_HOME"); got != want || got == "" {
		t.Fatalf("testCache(XDG_CACHE_HOME) = %q, want %q", got, want)
	}
	if got := testCache("HOME"); got != "" {
		t.Fatalf("testCache(HOME) = %q, want empty", got)
	}
	ctx := testCtxNoIndex(t.TempDir())
	if got := ctx.Getenv("TMPDIR"); got != testTmp {
		t.Fatalf("testCtxNoIndex's TMPDIR = %q, want %q", got, testTmp)
	}
}

// throwawayProjections makes ctx open throwaway indexes (no index path,
// no user cache folder for another root's) under a fresh TMPDIR, and
// returns a check that none is left there.
func throwawayProjections(t *testing.T, ctx *Ctx) func() {
	t.Helper()
	tmp := t.TempDir()
	ctx.IndexPath = ""
	inner := ctx.Getenv
	ctx.Getenv = func(k string) string {
		switch k {
		case "TMPDIR":
			return tmp
		case "XDG_CACHE_HOME", "HOME", "LocalAppData":
			return ""
		}
		if inner == nil {
			return ""
		}
		return inner(k)
	}
	return func() {
		t.Helper()
		left, err := filepath.Glob(filepath.Join(tmp, "rhiza-projection-*"))
		if err != nil {
			t.Fatal(err)
		}
		if len(left) > 0 {
			t.Errorf("projections left open: %v", left)
		}
	}
}

func TestShipCheckModelConflictClosesProjection(t *testing.T) {
	t.Parallel()
	f, _, _ := modelConflictFixture(t)
	if _, err := callOp(t, f.ctx, "ship.land", `{}`); err == nil {
		t.Fatal("land succeeded")
	}
	check := throwawayProjections(t, f.ctx)
	if r := shipCheck(t, f, `{"model_conflict":true}`); r.ModelConflict == nil || r.ModelConflict.Intent != "Ship it" {
		t.Fatalf("check = %+v", r)
	}
	check()
}

func TestRunStatusGatedClosesProjection(t *testing.T) {
	t.Parallel()
	f := newMessageFixture(t)
	cs := currentOf(t, f.ctx).String()
	rs := run.NewStore(run.Dir(f.common))
	t0 := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	if err := rs.Create("R1", run.Entry{Type: run.TypeStarted, At: t0, Started: &run.Started{Workflow: "w", Change: cs}}); err != nil {
		t.Fatal(err)
	}
	if err := rs.Append("R1", run.Entry{Type: run.TypeGated, At: t0.Add(time.Minute), Changes: []string{cs}}); err != nil {
		t.Fatal(err)
	}
	check := throwawayProjections(t, f.ctx)
	if _, err := callOp(t, f.ctx, "run.status", `{"run":"R1"}`); err != nil {
		t.Fatal(err)
	}
	check()
}

func TestWorktreeWorkspaceClosesReplacedProjection(t *testing.T) {
	t.Parallel()
	main := newMessageFixture(t)
	other := newSessionFixture(t, "other")
	cs := currentOf(t, main.ctx)
	check := throwawayProjections(t, main.ctx)
	w, err := openCtxWorkspace(main.ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer w.closeProjection()
	// The other root lacks cs, so its first workspace is replaced after
	// the copy; that first one's projection must be closed.
	ww, err := worktreeWorkspace(main.ctx, w, other.ctx.Root, cs)
	if err != nil {
		t.Fatal(err)
	}
	ww.closeProjection()
	w.closeProjection()
	check()

	// A change set the main log lacks fails the copy; the open workspace
	// is closed on that error too.
	if _, err := worktreeWorkspace(main.ctx, w, other.ctx.Root, newID(t, main.ctx)); err == nil {
		t.Fatal("unknown change set opened")
	}
	w.closeProjection()
	check()
}

func TestSessionWorkspaceMissClosesProjection(t *testing.T) {
	t.Parallel()
	f := newSessionFixture(t, "miss")
	check := throwawayProjections(t, f.ctx)
	// The worktree's log lacks the change set the record names, so the
	// workspace opened for it is closed on the not_found return.
	rec := session.Record{ID: "S1", Worktree: f.worktree, Change: newID(t, f.ctx).String()}
	_, _, err := sessionWorkspace(f.ctx, rec)
	if e := AsError(err); e == nil || e.Code != CodeNotFound {
		t.Fatalf("err = %v, want not_found", err)
	}
	check()
}
