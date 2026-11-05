package gosweep

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func write(t *testing.T, path string, size int, age time.Duration, now time.Time) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, make([]byte, size), 0o644); err != nil {
		t.Fatal(err)
	}
	mt := now.Add(-age)
	if err := os.Chtimes(path, mt, mt); err != nil {
		t.Fatal(err)
	}
}

func exists(p string) bool { _, err := os.Stat(p); return err == nil }

func setup(t *testing.T, now time.Time) (cache, tmp string) {
	cache, tmp = t.TempDir(), t.TempDir()
	write(t, filepath.Join(cache, "ab", "old1-a"), 10, 4*time.Hour, now)
	write(t, filepath.Join(cache, "ab", "old2-d"), 20, 5*time.Hour, now)
	write(t, filepath.Join(cache, "cd", "fresh-a"), 40, time.Hour, now)
	write(t, filepath.Join(cache, "README"), 5, 9*time.Hour, now)

	write(t, filepath.Join(tmp, "go-build111", "x", "f"), 100, time.Hour, now)
	write(t, filepath.Join(tmp, "go-build222", "f"), 200, time.Hour, now)
	write(t, filepath.Join(tmp, "other-old", "f"), 300, time.Hour, now)
	for _, d := range []string{"go-build111", "other-old"} {
		old := now.Add(-3 * time.Hour)
		if err := os.Chtimes(filepath.Join(tmp, d), old, old); err != nil {
			t.Fatal(err)
		}
	}
	fresh := now.Add(-30 * time.Minute)
	if err := os.Chtimes(filepath.Join(tmp, "go-build222"), fresh, fresh); err != nil {
		t.Fatal(err)
	}
	return cache, tmp
}

func TestSweepDirsRemovesOnlyStale(t *testing.T) {
	now := time.Now()
	cache, tmp := setup(t, now)
	got := SweepDirs(cache, tmp, now)
	want := Freed{CacheEntries: 2, CacheBytes: 30, TempFolders: 1, TempBytes: 100}
	if got != want {
		t.Fatalf("freed %+v, want %+v", got, want)
	}
	for p, keep := range map[string]bool{
		filepath.Join(cache, "ab", "old1-a"):  false,
		filepath.Join(cache, "ab", "old2-d"):  false,
		filepath.Join(cache, "cd", "fresh-a"): true,
		filepath.Join(cache, "README"):        true,
		filepath.Join(tmp, "go-build111"):     false,
		filepath.Join(tmp, "go-build222"):     true,
		filepath.Join(tmp, "other-old"):       true,
	} {
		if exists(p) != keep {
			t.Errorf("%s exists=%v, want %v", p, exists(p), keep)
		}
	}
	if !got.Any() || got.Bytes() != 130 {
		t.Errorf("Any/Bytes wrong: %v %d", got.Any(), got.Bytes())
	}
}

func TestCacheOffTouchesNoEntry(t *testing.T) {
	now := time.Now()
	cache, tmp := setup(t, now)
	t.Setenv("GOCACHE", "off")
	if d := CacheDir(); d != "" {
		t.Fatalf("CacheDir = %q with GOCACHE=off", d)
	}
	got := SweepDirs(CacheDir(), tmp, now)
	if got.CacheEntries != 0 || got.TempFolders != 1 {
		t.Fatalf("freed %+v", got)
	}
	if !exists(filepath.Join(cache, "ab", "old1-a")) {
		t.Error("cache entry removed with GOCACHE=off")
	}
}

func TestCacheDirDefaults(t *testing.T) {
	t.Setenv("GOCACHE", "/x/y")
	if CacheDir() != "/x/y" {
		t.Error("GOCACHE not honoured")
	}
	t.Setenv("GOCACHE", "")
	t.Setenv("XDG_CACHE_HOME", "/c")
	t.Setenv("HOME", "/h")
	if d := CacheDir(); d != "/c/go-build" && d != "/h/Library/Caches/go-build" {
		t.Errorf("CacheDir = %q", d)
	}
}

func TestMissingDirsAreFine(t *testing.T) {
	if f := SweepDirs("/nonexistent/a", "/nonexistent/b", time.Now()); f.Any() {
		t.Errorf("freed %+v", f)
	}
}
