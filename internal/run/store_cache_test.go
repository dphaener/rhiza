package run

import (
	"testing"
	"time"
)

// An unchanged journal is parsed once however often it is read; an
// appended one is parsed again, and appending to what a read returned
// never reaches the cache.
func TestStoreCachesUnchangedJournals(t *testing.T) {
	s := NewStore(t.TempDir())
	at := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	if err := s.Create("r1", Entry{Type: TypeStarted, At: at, Change: "c"}); err != nil {
		t.Fatal(err)
	}
	before := journalReads.Load()
	for range 3 {
		js, err := s.List()
		if err != nil {
			t.Fatal(err)
		}
		if len(js) != 1 || len(js[0].Entries) != 1 {
			t.Fatalf("got %+v", js)
		}
		js[0].Entries = append(js[0].Entries, Entry{Type: TypeFinished})
	}
	if n := journalReads.Load() - before; n != 1 {
		t.Fatalf("parsed %d times, want 1", n)
	}
	if err := s.Append("r1", Entry{Type: TypeFinished, At: at.Add(time.Minute), End: Failed}); err != nil {
		t.Fatal(err)
	}
	j, err := s.Read("r1")
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Entries) != 2 || j.Entries[1].Type != TypeFinished {
		t.Fatalf("after append got %+v", j.Entries)
	}
	if n := journalReads.Load() - before; n != 2 {
		t.Fatalf("parsed %d times, want 2", n)
	}
}
