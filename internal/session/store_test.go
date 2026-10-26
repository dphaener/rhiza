package session

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func rec(id string, at time.Time) Record {
	return Record{ID: id, Agent: "a", Change: "c", StartedAt: at, PermissionMode: "manual"}
}

func newStore(t *testing.T) *Store {
	t.Helper()
	return NewStore(Dir(t.TempDir()))
}

func TestStoreDir(t *testing.T) {
	if got := Dir("/r/.git"); got != "/r/.git/rhiza/sessions" {
		t.Errorf("Dir = %s", got)
	}
}

func TestStoreWriteReadAtomic(t *testing.T) {
	s := newStore(t)
	at := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := s.Write(rec("aaaa1111", at)); err != nil {
		t.Fatal(err)
	}
	if err := s.Write(rec("aaaa1111", at)); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "aaaa1111.json" {
		t.Fatalf("leftover files: %v", entries)
	}
	got, err := s.Read("aaaa1111")
	if err != nil || !got.StartedAt.Equal(at) || got.PermissionMode != "manual" {
		t.Fatalf("read %+v, %v", got, err)
	}
}

func TestStoreReadPrefix(t *testing.T) {
	s := newStore(t)
	at := time.Now()
	for _, id := range []string{"abcd1111", "abce2222", "ffff3333"} {
		if err := s.Write(rec(id, at)); err != nil {
			t.Fatal(err)
		}
	}
	got, err := s.Read("ffff")
	if err != nil || got.ID != "ffff3333" {
		t.Fatalf("prefix: %+v, %v", got, err)
	}
	got, err = s.Read("abcd1")
	if err != nil || got.ID != "abcd1111" {
		t.Fatalf("prefix: %+v, %v", got, err)
	}
}

func TestStoreReadUnknownID(t *testing.T) {
	s := newStore(t)
	if _, err := s.Read("deadbeef"); !isNotFound(err) {
		t.Errorf("empty store: %v", err)
	}
	if err := s.Write(rec("abcd1111", time.Now())); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"deadbeef", "abc", "", "../x", "abcd11112"} {
		if _, err := s.Read(id); !isNotFound(err) {
			t.Errorf("Read(%q) = %v, want not found", id, err)
		}
	}
}

func TestStoreReadAmbiguousPrefix(t *testing.T) {
	s := newStore(t)
	for _, id := range []string{"abcd1111", "abcd2222"} {
		if err := s.Write(rec(id, time.Now())); err != nil {
			t.Fatal(err)
		}
	}
	_, err := s.Read("abcd")
	var amb *AmbiguousError
	if !errors.As(err, &amb) || len(amb.Matches) != 2 {
		t.Fatalf("want ambiguous, got %v", err)
	}
	if !strings.Contains(err.Error(), "abcd2222") {
		t.Errorf("message should name matches: %v", err)
	}
	if _, err := s.Read("abcd1"); err != nil {
		t.Errorf("longer prefix should resolve: %v", err)
	}
}

func TestStoreListOrder(t *testing.T) {
	s := newStore(t)
	if got, err := s.List(); err != nil || len(got) != 0 {
		t.Fatalf("empty list: %v %v", got, err)
	}
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, r := range []Record{rec("cccc", t0), rec("bbbb", t0), rec("aaaa", t0.Add(time.Hour)), rec("dddd", t0.Add(-time.Hour))} {
		if err := s.Write(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(s.dir, ".stray.tmp"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, r := range got {
		ids = append(ids, r.ID)
	}
	if strings.Join(ids, ",") != "dddd,bbbb,cccc,aaaa" {
		t.Errorf("order %v", ids)
	}
}

func TestStoreUpdateTouchesOnlyItsSource(t *testing.T) {
	s := newStore(t)
	if err := s.Write(rec("abcd1111", time.Now())); err != nil {
		t.Fatal(err)
	}
	at := time.Now().UTC()
	if _, err := s.Update("abcd", func(r *Record) { r.Status.Hook = &Report{State: StateWorking, At: at} }); err != nil {
		t.Fatal(err)
	}
	code := 0
	if _, err := s.Update("abcd1111", func(r *Record) { r.Status.Exit = &Report{State: StateExited, Code: &code, At: at} }); err != nil {
		t.Fatal(err)
	}
	got, _ := s.Read("abcd1111")
	if got.Status.Hook == nil || got.Status.Exit == nil || got.Status.Rex != nil {
		t.Errorf("status %+v", got.Status)
	}
	if _, err := s.Update("nope1234", func(*Record) {}); !isNotFound(err) {
		t.Errorf("update unknown: %v", err)
	}
}

func isNotFound(err error) bool {
	var nf *NotFoundError
	return errors.As(err, &nf)
}
