package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var update = flag.Bool("update", false, "rewrite golden files")

func TestNotes(t *testing.T) {
	dir := t.TempDir()
	before, mid, after := "2026-11-01T00:00:00Z", "2026-11-10T00:00:00Z", "2026-11-11T00:00:00Z"
	writeAt(t, dir, "CS1", "a", "open", before, `{"intent":"Old feature"}`)
	writeAt(t, dir, "CS1", "b", "land", before, `{}`)
	writeAt(t, dir, "CS2", "a", "open", mid, `{"intent":"Fix the crash on empty logs","request_kind":"bug"}`)
	writeAt(t, dir, "CS2", "b", "land", mid, `{}`)
	writeAt(t, dir, "CS3", "a", "open", mid, `{"intent":"Release marks"}`)
	writeAt(t, dir, "CS3", "b", "ops", mid, `{"ops":[{"op":"node.create","kind":"intent","node":"I2","handle":"m33","fields":{"milestone":"33","statement":"Releases are cut by rule"}},
		{"op":"link.create","type":"serves","from":"CS3","to":"I2"},{"op":"link.create","type":"serves","from":"CS4","to":"I2"}]}`)
	writeAt(t, dir, "CS3", "c", "land", mid, `{}`)
	writeAt(t, dir, "CS4", "a", "open", after, `{"intent":"Evaluate the gates"}`)
	writeAt(t, dir, "CS4", "b", "land", after, `{}`)
	writeAt(t, dir, "CS5", "a", "open", after, `{"intent":"Drop the --old flag"}`)
	writeAt(t, dir, "CS5", "b", "land", after, `{}`)
	writeAt(t, dir, "CS6", "a", "open", after, `{"intent":"Repair a corrupted log","request_kind":"bug"}`)
	writeAt(t, dir, "CS6", "b", "land", after, `{}`)
	writeAt(t, dir, "CS7", "a", "open", after, `{"intent":"Record agent actor claude-code-1 for session s","kind":"record"}`)
	writeAt(t, dir, "CS7", "b", "land", after, `{}`)
	writeAt(t, dir, "CS8", "a", "open", after, `{"intent":"apply /tmp/ev.json","kind":"record"}`)
	writeAt(t, dir, "CS8", "b", "land", after, `{}`)
	h, err := ReadHistory(dir)
	if err != nil {
		t.Fatal(err)
	}
	marks := map[string]Marks{
		"CS5": {Breaking: []string{"removes --old"}},
		"CS6": {Severe: []string{"log corruption"}},
	}
	if h.Landed["CS7"].Kind != "record" || h.Landed["CS3"].Kind != "" {
		t.Fatalf("change set kinds not read: %+v", h.Landed)
	}
	got := Notes(h, time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC), marks, []string{"cli: command gone: rhiza old"})
	golden := filepath.Join("testdata", "notes.golden")
	if *update {
		os.MkdirAll("testdata", 0o755)
		os.WriteFile(golden, []byte(got), 0o644)
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if got != string(want) {
		t.Fatalf("notes differ from %s:\n%s", golden, got)
	}
	if strings.Contains(got, "Record agent actor") || strings.Contains(got, "ev.json") {
		t.Fatalf("record change sets listed:\n%s", got)
	}
	if again := Notes(h, time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC), marks, []string{"cli: command gone: rhiza old"}); again != got {
		t.Fatal("not deterministic")
	}
}
