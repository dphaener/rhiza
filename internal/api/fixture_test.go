package api

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/dphaener/rhiza/internal/index"
)

// A snapshot is a common starting repository (the model many tests rebuild
// through the same writes), built once per test binary into a template
// folder and copied into each test's own t.TempDir(), so a test owns its
// copy and may write to it. The build runs at most once whichever test asks
// first, guarded by a sync.Once, so parallel tests may share a snapshot;
// a failed build fails every test that asks for it with the build's error.
//
// A build lays its files out under dir and returns the Ctx that its writes
// went through: the copy's Ctx is testCtx over the same folder in the copy,
// with its ID generator advanced past the IDs the build drew, so a test on
// a copy draws exactly the IDs it would have drawn after building afresh.
//
// The template carries an index too: once the build is done, an index file
// beside the template folder (never inside it, so it is not copied as part
// of the tree) is caught up with the template's log and closed. Each copy
// gets its own copy of that file, under the test's own t.TempDir(), as its
// IndexPath, so a test catches up only what it writes after the copy,
// never the whole log from an empty file, and no two tests share an index.
// The index's freshness check lets a copied file serve a copied tree: the
// stats differ (new inodes and change times), the listing does not.
// A build must write no absolute path into the files (the copy lives
// elsewhere); a git repository is copied whole, .git included.
type snapshot struct {
	name  string
	build func(dir string) (*Ctx, error)

	once   sync.Once
	builds atomic.Int32
	dir    string // the template folder
	index  string // the template's caught-up index file, beside dir
	rel    string // the build Ctx's Root, relative to dir
	draws  int    // IDs the build drew from its Ctx's generator
	err    error
}

// snapshots holds every named snapshot of the test binary.
var snapshots sync.Map // name -> *snapshot

// snapshotOf returns the snapshot called name, registering build for it on
// first use. Every caller of one name must pass the same build.
func snapshotOf(name string, build func(dir string) (*Ctx, error)) *snapshot {
	s, _ := snapshots.LoadOrStore(name, &snapshot{name: name, build: build})
	return s.(*snapshot)
}

// snapshotDir is the folder every template is built under, chosen once per
// binary by TestMain (beside the binary's temporary user cache) and removed
// when the tests end. It is never read from the environment at build time:
// a test that t.Setenvs XDG_CACHE_HOME to its own t.TempDir() may be the one
// to trigger a build, and a template under that folder would be deleted
// with it, failing every later copy.
var snapshotDir string

// snapshotBase is the folder every template is built under.
func snapshotBase() (string, error) {
	if snapshotDir == "" {
		return "", errors.New("snapshotDir is not set (TestMain sets it)")
	}
	return snapshotDir, nil
}

// template builds the snapshot once and returns its template folder.
func (s *snapshot) template() error {
	s.once.Do(func() {
		s.builds.Add(1)
		s.err = s.make()
		if s.err != nil {
			s.err = fmt.Errorf("snapshot %s: build: %w", s.name, s.err)
		}
	})
	return s.err
}

func (s *snapshot) make() error {
	base, err := snapshotBase()
	if err != nil {
		return err
	}
	dir, err := os.MkdirTemp(base, "snapshot-")
	if err != nil {
		return err
	}
	ctx, err := s.build(dir)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(dir, ctx.Root)
	if err != nil {
		return err
	}
	draws, err := drawsOf(ctx)
	if err != nil {
		return err
	}
	if err := noAbsolutePaths(dir); err != nil {
		return err
	}
	idx := dir + ".index.db"
	if err := catchUpIndex(ctx.Root, idx); err != nil {
		return err
	}
	s.dir, s.rel, s.draws, s.index = dir, rel, draws, idx
	return nil
}

// catchUpIndex catches the index file at path up with the log of the
// repository at root and closes it, leaving one self-contained file (no
// write-ahead log beside it) to copy.
func catchUpIndex(root, path string) error {
	x, err := index.Open(root, path)
	if err != nil {
		return err
	}
	if err := x.Close(); err != nil {
		return err
	}
	for _, side := range []string{"-wal", "-journal"} {
		if _, err := os.Stat(path + side); !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("index %s left %s behind after closing (%v)", path, side, err)
		}
	}
	return nil
}

// drawsOf counts the IDs ctx's generator has drawn since testCtxNoIndex made
// it: it draws the next ID from ctx (whose build is over) and counts the
// draws a fresh testCtxNoIndex generator takes to reach one before it.
func drawsOf(ctx *Ctx) (int, error) {
	next, err := ctx.IDs.New()
	if err != nil {
		return 0, err
	}
	fresh := testCtxNoIndex(ctx.Root).IDs
	for n := 0; n < 1_000_000; n++ {
		v, err := fresh.New()
		if err != nil {
			return 0, err
		}
		if v == next {
			return n, nil
		}
	}
	return 0, fmt.Errorf("the build's next ID %s is not one testCtxNoIndex's generator draws", next)
}

// noAbsolutePaths fails when a file under dir names dir, which a copy
// elsewhere would still point back at.
func noAbsolutePaths(dir string) error {
	return filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || !d.Type().IsRegular() {
			return err
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if bytes.Contains(b, []byte(dir)) {
			return fmt.Errorf("%s names the template folder %s", p, dir)
		}
		return nil
	})
}

// copyTo copies the template into dst (an empty or missing folder) and
// returns the Ctx over the copy. When indexPath is set (a file outside dst
// that no other Ctx uses), the template's index is copied there and is the
// Ctx's IndexPath; when it is empty the Ctx has no index.
func (s *snapshot) copyTo(dst, indexPath string) (*Ctx, error) {
	if err := s.template(); err != nil {
		return nil, err
	}
	if err := copyTree(s.dir, dst); err != nil {
		return nil, fmt.Errorf("snapshot %s: copy: %w", s.name, err)
	}
	ctx := testCtxNoIndex(filepath.Join(dst, s.rel))
	if indexPath != "" {
		if err := copyFile(s.index, indexPath); err != nil {
			return nil, fmt.Errorf("snapshot %s: copy the index: %w", s.name, err)
		}
		ctx.IndexPath = indexPath
	}
	for range s.draws {
		if _, err := ctx.IDs.New(); err != nil {
			return nil, err
		}
	}
	return ctx, nil
}

// copyFile copies the regular file src to dst, which must not exist.
func copyFile(src, dst string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// copy gives t its own copy of the snapshot in a new t.TempDir(), with its
// own copy of the template's index in another, and returns the Ctx over
// it, guarded as testCtx's is, failing t with the build's error when the
// build failed.
func (s *snapshot) copy(t testing.TB) (string, *Ctx) {
	t.Helper()
	dst := t.TempDir()
	ctx, err := s.copyTo(dst, filepath.Join(t.TempDir(), "index.db"))
	if err != nil {
		t.Fatal(err)
	}
	guardProjections(t, ctx)
	return dst, ctx
}

// copyTree copies the tree at src into dst, keeping every file's and
// folder's mode (but dst's own) and every symbolic link as a link. Folders get their mode
// after their contents are written, so a read-only folder copies too.
func copyTree(src, dst string) error {
	type dirMode struct {
		path string
		mode fs.FileMode
	}
	var dirs []dirMode
	err := filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		to := filepath.Join(dst, rel)
		info, err := d.Info()
		if err != nil {
			return err
		}
		switch {
		case d.IsDir():
			// dst itself keeps its own mode: it is the test's folder.
			if p != src {
				dirs = append(dirs, dirMode{to, info.Mode().Perm()})
			}
			return os.MkdirAll(to, 0o755)
		case info.Mode()&fs.ModeSymlink != 0:
			link, err := os.Readlink(p)
			if err != nil {
				return err
			}
			return os.Symlink(link, to)
		case info.Mode().IsRegular():
			b, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			if err := os.WriteFile(to, b, info.Mode().Perm()); err != nil {
				return err
			}
			// WriteFile's mode is masked by the umask; set it exactly.
			return os.Chmod(to, info.Mode().Perm())
		}
		return fmt.Errorf("%s: cannot copy a %s", p, info.Mode().Type())
	})
	if err != nil {
		return err
	}
	for i := len(dirs) - 1; i >= 0; i-- {
		if err := os.Chmod(dirs[i].path, dirs[i].mode); err != nil {
			return err
		}
	}
	return nil
}

// runOp is callOp for a snapshot build, which has no *testing.T: it
// returns a decode failure as an error.
func runOp(ctx *Ctx, name, reqJSON string) (any, error) {
	o, ok := Lookup(name)
	if !ok {
		return nil, fmt.Errorf("%s is not registered", name)
	}
	req, err := DecodeRequest(o, []byte(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("DecodeRequest(%s, %s): %w", name, reqJSON, err)
	}
	resp, err := o.Run(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w", name, reqJSON, err)
	}
	return resp, nil
}

// buildInit is the snapshot of a repository just after init: the repo is
// dir itself.
func buildInit(dir string) (*Ctx, error) {
	// No index: the build runs once per binary, outside any test's TempDir.
	ctx := testCtxNoIndex(dir)
	_, err := runOp(ctx, "init", `{}`)
	return ctx, err
}

// initSnapshot is initCtx's repository.
var initSnapshot = snapshotOf("init", buildInit)

// filesOf maps every file and link under dir, by slash path relative to
// dir, to its mode and contents (a link's target).
func filesOf(t *testing.T, dir string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		var body string
		switch {
		case info.Mode()&fs.ModeSymlink != 0:
			l, err := os.Readlink(p)
			if err != nil {
				return err
			}
			body = "-> " + l
		case info.Mode().IsRegular():
			b, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			body = string(b)
		}
		out[filepath.ToSlash(rel)] = info.Mode().String() + " " + body
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// TestSnapshotFixtureFoldsLikeFresh shows a copy of each fixture is the
// repository its build makes afresh: the same files, folding to the same
// model, with a Ctx that draws the same next ID.
func TestSnapshotFixtureFoldsLikeFresh(t *testing.T) {
	t.Parallel()
	for _, s := range []*snapshot{initSnapshot, linkEndpoints} {
		t.Run(s.name, func(t *testing.T) {
			t.Parallel()
			fresh := t.TempDir()
			freshCtx, err := s.build(fresh)
			if err != nil {
				t.Fatal(err)
			}
			dst, ctx := s.copy(t)
			if got, want := filesOf(t, dst), filesOf(t, fresh); !reflect.DeepEqual(got, want) {
				t.Fatalf("copy's files differ from a fresh build's:\ncopy  %v\nfresh %v", got, want)
			}
			if ctx.Root != filepath.Join(dst, s.rel) {
				t.Errorf("copy's Root = %s, want under %s", ctx.Root, dst)
			}
			got, want := mustFoldClean(t, ctxWorkspace(t, ctx)), mustFoldClean(t, ctxWorkspace(t, freshCtx))
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("copy folds to %d change sets %+v, fresh to %d %+v", len(got), got, len(want), want)
			}
			if gc, wc := currentOf(t, ctx), currentOf(t, freshCtx); gc != wc {
				t.Errorf("copy's current change set %s, fresh %s", gc, wc)
			}
			gotID, err1 := ctx.IDs.New()
			wantID, err2 := freshCtx.IDs.New()
			if err1 != nil || err2 != nil || gotID != wantID {
				t.Errorf("copy's next ID %s (%v), fresh %s (%v)", gotID, err1, wantID, err2)
			}
		})
	}
}

// TestSnapshotFixtureCopiesAreIsolated shows a write in one test's copy is
// seen neither by another copy nor by the template later copies come from.
func TestSnapshotFixtureCopiesAreIsolated(t *testing.T) {
	t.Parallel()
	_, a := initSnapshot.copy(t)
	bDir, b := initSnapshot.copy(t)
	before := filesOf(t, bDir)
	cs := newChange(t, a, `{"intent":"Only in a"}`).ChangeSet
	if got := currentOf(t, a); got != cs {
		t.Fatalf("a's current = %s, want %s", got, cs)
	}
	if got := filesOf(t, bDir); !reflect.DeepEqual(got, before) {
		t.Fatalf("a write in copy a changed copy b:\nbefore %v\nafter  %v", before, got)
	}
	if _, ok := ctxWorkspace(t, b).changeSet(cs); ok {
		t.Errorf("copy b sees change set %s written in copy a", cs)
	}
	cDir, _ := initSnapshot.copy(t)
	if got := filesOf(t, cDir); !reflect.DeepEqual(got, before) {
		t.Errorf("a later copy differs from b, so the write reached the template:\nb %v\nc %v", before, got)
	}
	// b writes on its own, drawing the IDs a drew: the copies share no
	// state, not even the generator.
	if got := newChange(t, b, `{"intent":"Only in b"}`).ChangeSet; got != cs {
		t.Errorf("b's first change set %s, want the same ID as a's %s", got, cs)
	}
}

// TestSnapshotFixtureBuildsOnce shows many goroutines asking for one
// fixture at once build it once, and each gets a whole copy of its own.
func TestSnapshotFixtureBuildsOnce(t *testing.T) {
	t.Parallel()
	s := &snapshot{name: "builds-once", build: buildInit}
	base := t.TempDir()
	const n = 16
	var wg sync.WaitGroup
	ctxs := make([]*Ctx, n)
	errs := make([]error, n)
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctxs[i], errs[i] = s.copyTo(filepath.Join(base, fmt.Sprint(i)), "")
		}()
	}
	wg.Wait()
	if got := s.builds.Load(); got != 1 {
		t.Errorf("%d goroutines built the fixture %d times, want once", n, got)
	}
	roots := map[string]bool{}
	for i, ctx := range ctxs {
		if errs[i] != nil {
			t.Fatalf("copy %d: %v", i, errs[i])
		}
		if roots[ctx.Root] {
			t.Fatalf("two copies share Root %s", ctx.Root)
		}
		roots[ctx.Root] = true
		mustFoldClean(t, ctxWorkspace(t, ctx))
	}
}

// TestSnapshotFixtureBuildError shows a fixture whose build fails fails
// every caller with the build's error, never a missing folder, and builds
// once.
func TestSnapshotFixtureBuildError(t *testing.T) {
	t.Parallel()
	boom := errors.New("boom: the build broke")
	s := &snapshot{name: "broken", build: func(dir string) (*Ctx, error) {
		if _, err := buildInit(dir); err != nil {
			return nil, err
		}
		return nil, boom
	}}
	base := t.TempDir()
	const n = 8
	var wg sync.WaitGroup
	errs := make([]error, n)
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, errs[i] = s.copyTo(filepath.Join(base, fmt.Sprint(i)), "")
		}()
	}
	wg.Wait()
	for i, err := range errs {
		if !errors.Is(err, boom) || errors.Is(err, fs.ErrNotExist) {
			t.Errorf("caller %d: err = %v, want the build's %q", i, err, boom)
		}
		if _, statErr := os.Stat(filepath.Join(base, fmt.Sprint(i))); !errors.Is(statErr, fs.ErrNotExist) {
			t.Errorf("caller %d got a folder from a failed build: %v", i, statErr)
		}
	}
	if got := s.builds.Load(); got != 1 {
		t.Errorf("the failing build ran %d times, want once", got)
	}
	// A later caller still gets the build's error.
	if _, err := s.copyTo(filepath.Join(base, "late"), ""); !errors.Is(err, boom) {
		t.Errorf("late caller: err = %v, want %q", err, boom)
	}
}

// TestSnapshotFixtureBuildUnderSetenvServesLaterCopies shows a build
// triggered by a test that sets XDG_CACHE_HOME to its own t.TempDir() (as
// sliceCtx and the inbox gate tests do before initCtx) still serves copies
// after that test's folder is gone: the template's folder is chosen once by
// TestMain, not read from the environment at build time.
func TestSnapshotFixtureBuildUnderSetenvServesLaterCopies(t *testing.T) {
	s := &snapshot{name: "built-under-setenv", build: buildInit}
	var cache string
	t.Run("build", func(t *testing.T) {
		cache = t.TempDir()
		t.Setenv("XDG_CACHE_HOME", cache)
		s.copy(t)
	})
	if _, err := os.Stat(cache); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("the building test's cache %s outlived it: %v", cache, err)
	}
	if rel, err := filepath.Rel(snapshotDir, s.dir); err != nil || !filepath.IsLocal(rel) {
		t.Errorf("template %s is not under snapshotDir %s", s.dir, snapshotDir)
	}
	_, ctx := s.copy(t)
	mustFoldClean(t, ctxWorkspace(t, ctx))
	if got := s.builds.Load(); got != 1 {
		t.Errorf("built %d times, want once", got)
	}
}

// TestSnapshotFixtureCopyCatchesUpFromCopiedIndex shows a copy's index is
// the template's, caught up: opened before any write it reads no entry
// file, and after an apply on the copy it reads only the entry the apply
// wrote (no rebuild from empty), which it then serves back. The template's
// own index file is untouched by it, so no copy writes the template's.
func TestSnapshotFixtureCopyCatchesUpFromCopiedIndex(t *testing.T) {
	t.Parallel()
	if err := initSnapshot.template(); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(initSnapshot.index)
	if err != nil {
		t.Fatal(err)
	}
	_, ctx := initSnapshot.copy(t)
	if ctx.IndexPath == "" || ctx.IndexPath == initSnapshot.index {
		t.Fatalf("copy's IndexPath = %q, want its own copy of %s", ctx.IndexPath, initSnapshot.index)
	}
	open := func() *index.Index {
		t.Helper()
		x, err := index.Open(ctx.Root, ctx.IndexPath)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { x.Close() })
		return x
	}
	if c := open().LastCatchUp(); !c.Fresh || c.Rebuilt || c.Read != 0 {
		t.Fatalf("copy's index before any write: %+v, want fresh with no entry read", c)
	}

	cs := newChange(t, ctx, `{"intent":"Shape"}`).ChangeSet
	x := open()
	if c := x.LastCatchUp(); c.Rebuilt || c.Read != 1 {
		t.Fatalf("copy's index after change.new: %+v, want the one new entry read, no rebuild", c)
	}
	es, err := x.Entries(cs)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := applyFile(t, ctx, treeWrites); err != nil {
		t.Fatal(err)
	}
	x = open()
	if c := x.LastCatchUp(); c.Rebuilt || c.Read != 1 {
		t.Errorf("copy's index after the apply: %+v, want the apply's one entry read, no rebuild", c)
	}
	after, err := x.Entries(cs)
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(es)+1 {
		t.Fatalf("copy's index serves %d entries of %s after the apply, want %d", len(after), cs, len(es)+1)
	}
	if got := nodeAt(t, ctx, "store"); got.Handle != "store" {
		t.Errorf("the applied node reads back as %+v", got)
	}

	now, err := os.Stat(initSnapshot.index)
	if err != nil {
		t.Fatal(err)
	}
	if now.Size() != before.Size() || !now.ModTime().Equal(before.ModTime()) {
		t.Errorf("the template's index changed: size %d -> %d, modified %s -> %s", before.Size(), now.Size(), before.ModTime(), now.ModTime())
	}
	for _, side := range []string{"-wal", "-shm"} {
		if _, err := os.Stat(initSnapshot.index + side); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("the template's index has a %s file: %v", side, err)
		}
	}
}
