package index

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
	"github.com/dphaener/genesis/internal/store"
)

// BenchmarkBuild builds the default model's index from scratch.
func BenchmarkBuild(b *testing.B) {
	g := defaultModel(b)
	dir := b.TempDir()
	b.ResetTimer()
	for i := range b.N {
		x, err := Open(g.root, filepath.Join(dir, "build-"+strconv.Itoa(i)+".db"))
		if err != nil {
			b.Fatal(err)
		}
		if err := x.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkFreshRead is TestBudget's unit of work.
func BenchmarkFreshRead(b *testing.B) {
	g := defaultModel(b)
	for b.Loop() {
		if _, _, _, _, err := freshRead(g); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCatchUpDraft catches up one new ops entry in the draft change
// set: the log is refolded, the views are not recomputed.
func BenchmarkCatchUpDraft(b *testing.B) {
	g := defaultModel(b)
	s := g.summary
	o := op.Op{ID: newID(b), Type: op.TypeNodeUpdate, Node: s.Spare.ID, Base: s.SpareVersion,
		Set: map[string]json.RawMessage{"notes": json.RawMessage(`"Benchmark note"`)}}
	benchCatchUp(b, g, s.Draft, s.DraftHead, changeset.Ops{Ops: []op.Op{o}}, func(st stats) bool {
		return st.read == 1 && !st.rebuilt && !st.asIs && !st.toBe
	})
}

// BenchmarkCatchUpLand catches up a land entry for the accepted change set,
// recomputing both views.
func BenchmarkCatchUpLand(b *testing.B) {
	g := defaultModel(b)
	s := g.summary
	benchCatchUp(b, g, s.Accepted, s.AcceptedHead, changeset.Land{}, func(st stats) bool {
		return st.read == 1 && !st.rebuilt && st.asIs && st.toBe
	})
}

// BenchmarkWhatIf computes the draft change set's what-if on an open index.
func BenchmarkWhatIf(b *testing.B) {
	g := defaultModel(b)
	x, err := Open(g.root, g.dbPath)
	if err != nil {
		b.Fatal(err)
	}
	defer x.Close()
	for b.Loop() {
		if _, _, err := x.WhatIf(g.summary.Draft); err != nil {
			b.Fatal(err)
		}
	}
}

// benchCatchUp times Open catching up one new entry with body body in cs
// after head. Outside the timer each iteration restores the starting
// state: the entry file is removed and the built database copied back (a
// vanished file would otherwise mean a rebuild). want checks the first
// catch-up did what the benchmark claims to measure.
func benchCatchUp(b *testing.B, g generated, cs, head id.ID, body changeset.Body, want func(stats) bool) {
	s := g.summary
	data := encodeBenchEntry(b, cs, head, s.Actor, s.Last.Add(time.Second), body)
	e, err := entry.Decode(data)
	if err != nil {
		b.Fatal(err)
	}
	file := filepath.Join(g.root, store.Dir, filepath.FromSlash(entryPath(cs, e.ID)))
	b.Cleanup(func() { os.Remove(file) })
	dbPath := filepath.Join(b.TempDir(), "index.db")
	restore := func() {
		for _, p := range []string{dbPath + "-wal", dbPath + "-shm"} {
			if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
				b.Fatal(err)
			}
		}
		if err := copyFile(g.dbPath, dbPath); err != nil {
			b.Fatal(err)
		}
		if err := os.WriteFile(file, data, 0o644); err != nil {
			b.Fatal(err)
		}
	}
	checked := false
	b.StopTimer()
	b.ResetTimer()
	for range b.N {
		restore()
		b.StartTimer()
		x, err := Open(g.root, dbPath)
		b.StopTimer()
		if err != nil {
			b.Fatal(err)
		}
		if !checked {
			if !want(x.stats) {
				b.Fatalf("catch-up stats %+v", x.stats)
			}
			ps, err := x.Problems()
			if err != nil {
				b.Fatal(err)
			}
			if len(ps) != 0 {
				b.Fatalf("catch-up left problems: %+v", ps)
			}
			checked = true
		}
		if err := x.Close(); err != nil {
			b.Fatal(err)
		}
		if err := os.Remove(file); err != nil {
			b.Fatal(err)
		}
	}
}

// encodeBenchEntry encodes a new entry of cs after head, as the gen.Summary
// recipe describes.
func encodeBenchEntry(b *testing.B, cs, head, actor id.ID, at time.Time, body changeset.Body) []byte {
	raw, err := changeset.EncodeBody(body)
	if err != nil {
		b.Fatal(err)
	}
	e := entry.Entry{V: changeset.EntryVersion(body), ID: newID(b), ChangeSet: cs, Actor: actor, For: actor,
		At: at.UTC().Truncate(time.Second), After: []id.ID{head}, Type: body.Type(), Body: raw}
	data, err := entry.Encode(e)
	if err != nil {
		b.Fatal(err)
	}
	return data
}

func newID(b *testing.B) id.ID {
	i, err := id.Default().New()
	if err != nil {
		b.Fatal(err)
	}
	return i
}
