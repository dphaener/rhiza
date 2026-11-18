package run

import (
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// isolateUserCache points the user cache and home at fresh test folders,
// so no test reads or writes the host's user cache, and returns the cache
// folder.
func isolateUserCache(t *testing.T) string {
	t.Helper()
	cache := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", cache)
	t.Setenv("HOME", t.TempDir())
	return cache
}

// clearJournalCache empties the process-wide journalCache, as a new
// process starts.
func clearJournalCache() {
	journalCache.mu.Lock()
	defer journalCache.mu.Unlock()
	clear(journalCache.files)
}

// diskStore is a fresh store over dir that keeps its list cache in the
// process environment's user cache folder.
func diskStore(dir string) *Store {
	s := NewStore(dir)
	s.Getenv = os.Getenv
	s.DiskCache = true
	return s
}

// uncachedList is List with no cache at all: every journal parsed.
func uncachedList(t *testing.T, dir string) []Journal {
	t.Helper()
	clearJournalCache()
	js, err := NewStore(dir).List()
	if err != nil {
		t.Fatal(err)
	}
	clearJournalCache()
	return js
}

// writeJournals writes three finished runs, with nested fields set, into
// a fresh store folder and returns it.
func writeJournals(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	s := NewStore(dir)
	at := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	for i, id := range []string{"r1", "r2", "r3"} {
		start := at.Add(time.Duration(i) * time.Hour)
		if err := s.Create(id, Entry{Type: TypeStarted, At: start, Step: "build", Started: &Started{
			Workflow: "wf", Definition: "def", Change: "c" + id, Includes: map[string]string{"inc": "src"},
		}}); err != nil {
			t.Fatal(err)
		}
		for _, e := range []Entry{
			{Type: TypeEntered, At: start.Add(time.Minute), Step: "build", Visit: 1},
			{Type: TypeBlocked, At: start.Add(2 * time.Minute).In(time.FixedZone("", 3600)), Step: "build", Block: &Block{}},
			{Type: "findings", At: start.Add(3 * time.Minute), Findings: []Finding{{ID: "f1", Reason: "r", Severity: "major"}}, Changes: []string{"c1"}},
			{Type: TypeFinished, At: start.Add(4 * time.Minute), End: Succeeded},
		} {
			if err := s.Append(id, e); err != nil {
				t.Fatal(err)
			}
		}
	}
	return dir
}

// A store's list cache outlives the process: with the process's journal
// cache cleared, a fresh store parses only the journal appended since the
// last List, lists what an uncached List does, and keeps the cache file
// under the user cache folder, never by the journals or in a repository.
func TestStoreListDiskCacheParsesOnlyChanged(t *testing.T) {
	cache := isolateUserCache(t)
	dir := writeJournals(t)
	clearJournalCache()
	if _, err := diskStore(dir).List(); err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 11, 4, 12, 0, 0, 0, time.UTC)
	if err := NewStore(dir).Append("r2", Entry{Type: TypeRetried, At: at, Step: "build", Detail: "again"}); err != nil {
		t.Fatal(err)
	}
	clearJournalCache()
	before := journalReads.Load()
	got, err := diskStore(dir).List()
	if err != nil {
		t.Fatal(err)
	}
	if n := journalReads.Load() - before; n != 1 {
		t.Fatalf("parsed %d journals, want only the appended one", n)
	}
	if want := uncachedList(t, dir); !reflect.DeepEqual(got, want) {
		t.Fatalf("cached List differs from uncached:\n%+v\n%+v", got, want)
	}
	if e := got[1].Entries; got[1].ID != "r2" || e[len(e)-1].Type != TypeRetried {
		t.Fatalf("appended entry missing: %+v", got[1])
	}

	file, ok := ListCachePath(dir, os.Getenv)
	if !ok {
		t.Fatal("no list cache path")
	}
	if _, err := os.Stat(file); err != nil {
		t.Fatal(err)
	}
	if rel, err := filepath.Rel(cache, file); err != nil || strings.HasPrefix(rel, "..") {
		t.Fatalf("cache file %s is not under the user cache %s", file, cache)
	}
	wd, _ := os.Getwd()
	for _, outside := range []string{dir, wd} {
		if rel, err := filepath.Rel(outside, file); err == nil && !strings.HasPrefix(rel, "..") {
			t.Fatalf("cache file %s lies under %s", file, outside)
		}
	}
	if strings.Contains(file, ".rhiza") {
		t.Fatalf("cache file %s is under .rhiza", file)
	}
	names, _ := os.ReadDir(dir)
	for _, n := range names {
		if !strings.HasSuffix(n.Name(), ".jsonl") {
			t.Fatalf("store folder holds %s", n.Name())
		}
	}
}

// A list cache file that is corrupt, or of another format version, is
// ignored: a fresh store parses every journal and lists what an uncached
// List does.
func TestStoreListDiskCacheIgnoresBadFile(t *testing.T) {
	for name, spoil := range map[string]func(t *testing.T, file string){
		"corrupt": func(t *testing.T, file string) {
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(file, data[:len(data)/2], 0o600); err != nil {
				t.Fatal(err)
			}
		},
		"garbage": func(t *testing.T, file string) {
			if err := os.WriteFile(file, []byte("not a cache"), 0o600); err != nil {
				t.Fatal(err)
			}
		},
		"other version": func(t *testing.T, file string) {
			f, err := os.Open(file)
			if err != nil {
				t.Fatal(err)
			}
			dec := gob.NewDecoder(f)
			var h listCacheHeader
			var recs []listCacheRecord
			if err := dec.Decode(&h); err != nil {
				t.Fatal(err)
			}
			if err := dec.Decode(&recs); err != nil {
				t.Fatal(err)
			}
			f.Close()
			if len(recs) != 3 {
				t.Fatalf("cache holds %d records, want 3", len(recs))
			}
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			h.Version++
			if err := enc.Encode(h); err != nil {
				t.Fatal(err)
			}
			if err := enc.Encode(recs); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(file, buf.Bytes(), 0o600); err != nil {
				t.Fatal(err)
			}
		},
	} {
		t.Run(name, func(t *testing.T) {
			isolateUserCache(t)
			dir := writeJournals(t)
			clearJournalCache()
			if _, err := diskStore(dir).List(); err != nil {
				t.Fatal(err)
			}
			file, ok := ListCachePath(dir, os.Getenv)
			if !ok {
				t.Fatal("no list cache path")
			}
			spoil(t, file)
			clearJournalCache()
			before := journalReads.Load()
			got, err := diskStore(dir).List()
			if err != nil {
				t.Fatal(err)
			}
			if n := journalReads.Load() - before; n != 3 {
				t.Fatalf("parsed %d journals, want all 3", n)
			}
			if want := uncachedList(t, dir); !reflect.DeepEqual(got, want) {
				t.Fatalf("List differs from uncached:\n%+v\n%+v", got, want)
			}
		})
	}
}

// A list cache file written by a binary whose Entry has another shape (a
// field added, removed or retyped, here or in a type it holds) is ignored:
// gob would decode its records leaving a new field zero, so every journal
// is parsed again.
func TestStoreListDiskCacheIgnoresOtherEntryShape(t *testing.T) {
	isolateUserCache(t)
	dir := writeJournals(t)
	clearJournalCache()
	shape := listCacheShape
	listCacheShape = "an older Entry"
	_, err := diskStore(dir).List()
	listCacheShape = shape
	if err != nil {
		t.Fatal(err)
	}
	clearJournalCache()
	before := journalReads.Load()
	got, err := diskStore(dir).List()
	if err != nil {
		t.Fatal(err)
	}
	if n := journalReads.Load() - before; n != 3 {
		t.Fatalf("parsed %d journals, want all 3", n)
	}
	if want := uncachedList(t, dir); !reflect.DeepEqual(got, want) {
		t.Fatalf("List differs from uncached:\n%+v\n%+v", got, want)
	}
}

// The Entry shape key changes with any field of Entry or of a type it
// holds, however deep.
func TestTypeShapeRecursive(t *testing.T) {
	type inner struct{ A int }
	type inner2 struct{ A string }
	type outer struct {
		X  string
		In *inner
	}
	type outer2 struct {
		X  string
		In *inner2
	}
	type self struct {
		Next *self
		V    int
	}
	if typeShape(reflect.TypeOf(outer{})) == typeShape(reflect.TypeOf(outer2{})) {
		t.Fatal("a nested field's type change kept the shape")
	}
	if typeShape(reflect.TypeOf(self{})) == "" {
		t.Fatal("a self-referencing type has no shape")
	}
	if listCacheShape != typeShape(reflect.TypeOf(Entry{})) {
		t.Fatal("listCacheShape is not Entry's shape")
	}
}

// A store without DiskCache, or with no user cache folder, writes no list
// cache.
func TestStoreListNoDiskCacheWritesNothing(t *testing.T) {
	cache := isolateUserCache(t)
	dir := writeJournals(t)
	s := NewStore(dir)
	s.Getenv = os.Getenv
	if _, err := s.List(); err != nil {
		t.Fatal(err)
	}
	nilEnv := NewStore(dir)
	nilEnv.DiskCache = true
	if _, err := nilEnv.List(); err != nil {
		t.Fatal(err)
	}
	if names, _ := os.ReadDir(cache); len(names) != 0 {
		t.Fatalf("user cache holds %v", names)
	}
}
