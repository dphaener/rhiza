package gen

import (
	"bytes"
	"flag"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/model"
	"github.com/dphaener/genesis/internal/op"
	"github.com/dphaener/genesis/internal/store"
)

var genOut = flag.String("gen.out", "", "write the default model into this folder (TestWriteDemo)")

func small() Options {
	return Options{Seed: 7, Nodes: 2000, Links: 800, Landed: 40, InFlight: 6}
}

// TestWriteDemo writes the default model to -gen.out for demonstrations.
func TestWriteDemo(t *testing.T) {
	if *genOut == "" {
		t.Skip("set -gen.out=DIR to write the default model there")
	}
	began := time.Now()
	s, err := Write(*genOut, Default())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("wrote %d entries in %d change sets to %s in %v; summary %+v", s.Entries, s.ChangeSets, *genOut, time.Since(began), s)
}

func TestWriteSmall(t *testing.T) {
	o := small()
	root := t.TempDir()
	s, err := Write(root, o)
	if err != nil {
		t.Fatal(err)
	}
	checkClean(t, root, o, s)
}

func TestWriteDefault(t *testing.T) {
	if testing.Short() {
		t.Skip("the default model is skipped under -short")
	}
	o := Default()
	root := t.TempDir()
	began := time.Now()
	s, err := Write(root, o)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Write(Default()) took %v: %d entries, %d change sets", time.Since(began), s.Entries, s.ChangeSets)
	checkClean(t, root, o, s)
}

func TestDeterministic(t *testing.T) {
	o := small()
	a, b, c := t.TempDir(), t.TempDir(), t.TempDir()
	for _, r := range []string{a, b} {
		if _, err := Write(r, o); err != nil {
			t.Fatal(err)
		}
	}
	o.Seed++
	if _, err := Write(c, o); err != nil {
		t.Fatal(err)
	}
	ta, tb, tc := tree(t, a), tree(t, b), tree(t, c)
	if len(ta) != len(tb) {
		t.Fatalf("same options: %d files vs %d", len(ta), len(tb))
	}
	for p, data := range ta {
		if !bytes.Equal(data, tb[p]) {
			t.Fatalf("same options: %s differs", p)
		}
	}
	same := 0
	for p, data := range ta {
		if bytes.Equal(data, tc[p]) {
			same++
		}
	}
	if same == len(ta) && len(ta) == len(tc) {
		t.Fatal("a different seed wrote the same files")
	}
	if !bytes.Equal(ta[".genesis/repo.json"], tb[".genesis/repo.json"]) || bytes.Equal(ta[".genesis/repo.json"], tc[".genesis/repo.json"]) {
		t.Fatal("repo.json should follow the seed")
	}
}

// TestRepoFile checks that repo.json is the bytes store.Init writes for
// the same ID.
func TestRepoFile(t *testing.T) {
	o := small()
	root := t.TempDir()
	if _, err := Write(root, o); err != nil {
		t.Fatal(err)
	}
	ref := t.TempDir()
	ids := id.NewGenerator(func() time.Time { return start }, rand.New(rand.NewSource(o.Seed^0x5eed)))
	if _, err := store.Init(ref, ids); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(root, store.Dir, store.RepoFile))
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join(ref, store.Dir, store.RepoFile))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("repo.json:\n%s\nstore.Init writes:\n%s", got, want)
	}
}

func TestOptionsChecked(t *testing.T) {
	if _, err := Write(t.TempDir(), Options{Seed: 1, Nodes: 10, Links: 2, Landed: 2, InFlight: 1}); err == nil {
		t.Fatal("InFlight 1 accepted")
	}
}

// tree reads every file under root, by slash path.
func tree(t *testing.T, root string) map[string][]byte {
	t.Helper()
	out := map[string][]byte{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, p)
		out[filepath.ToSlash(rel)] = data
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// checkClean reads, folds and projects root and checks the properties every
// generated repository has.
func checkClean(t *testing.T, root string, o Options, s Summary) {
	t.Helper()
	if _, err := store.Open(root); err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	l := store.ReadLog(os.DirFS(filepath.Join(root, store.Dir)))
	for _, p := range l.Problems {
		t.Errorf("log problem: %s: %v", p.Path, p.Err)
	}
	if len(l.Entries) != s.Entries {
		t.Errorf("entries: read %d, summary says %d", len(l.Entries), s.Entries)
	}
	sets, fps := changeset.FoldLog(l)
	for _, p := range fps {
		t.Errorf("fold problem: %v", p)
	}
	if len(sets) != s.ChangeSets || len(sets) != o.Landed+o.InFlight {
		t.Errorf("change sets: %d, summary %d, options %d", len(sets), s.ChangeSets, o.Landed+o.InFlight)
	}
	asIs, ps := model.AsIs(sets)
	for _, p := range ps {
		t.Errorf("as-is problem: %v", p)
	}
	toBe, ps := model.ToBe(sets)
	for _, p := range ps {
		t.Errorf("to-be problem: %v", p)
	}
	for _, cs := range []id.ID{s.Accepted, s.Draft} {
		if _, ps := model.WhatIf(sets, cs); len(ps) > 0 {
			t.Errorf("what-if %s: %v", cs, ps[0])
		}
	}
	if t.Failed() {
		t.FailNow()
	}

	nodes, links := asIs.Nodes(), asIs.Links()
	if len(nodes) != o.Nodes || len(nodes) != s.Nodes {
		t.Errorf("as-is nodes: %d, options %d, summary %d", len(nodes), o.Nodes, s.Nodes)
	}
	if len(links) != o.Links || len(links) != s.Links {
		t.Errorf("as-is links: %d, options %d, summary %d", len(links), o.Links, s.Links)
	}
	if len(toBe.Nodes()) <= len(nodes) {
		t.Errorf("to-be has %d nodes, as-is %d: accepted creates missing", len(toBe.Nodes()), len(nodes))
	}

	kinds := map[op.Kind]int{}
	for _, n := range nodes {
		kinds[n.Kind]++
	}
	for _, k := range op.Kinds() {
		if k == op.KindAnchor {
			if kinds[k] > 0 {
				t.Errorf("%d anchors", kinds[k])
			}
			continue
		}
		if kinds[k] == 0 {
			t.Errorf("no %s nodes", k)
		}
	}
	types := map[op.LinkType]bool{}
	for _, lk := range links {
		types[lk.Type] = true
	}
	if len(types) < 8 {
		t.Errorf("%d link types, want at least 8", len(types))
	}

	ops := map[op.Type]int{}
	for _, cs := range sets {
		if cs.State != changeset.StateLanded {
			continue
		}
		for _, x := range cs.Ops {
			ops[x.Type]++
		}
	}
	for _, typ := range op.Types() {
		if ops[typ] == 0 {
			t.Errorf("no landed %s", typ)
		}
	}
	t.Logf("landed ops: %v; kinds: %v; link types: %d", ops, kinds, len(types))

	depth := func(i id.ID) int {
		d := 0
		for !i.IsZero() {
			n, ok := asIs.Node(i)
			if !ok {
				t.Fatalf("ancestor %s missing", i)
			}
			d++
			i = n.Parent
		}
		return d
	}
	maxDepth := 0
	for _, n := range nodes {
		if op.HasParent(n.Kind) {
			maxDepth = max(maxDepth, depth(n.ID))
		}
	}
	if maxDepth < 4 {
		t.Errorf("tree depth %d, want at least 4", maxDepth)
	}

	live := func(sm Sample) model.Node {
		t.Helper()
		n, ok := asIs.Node(sm.ID)
		if !ok || n.Retired || n.Handle != sm.Handle {
			t.Fatalf("sample %+v: node %+v, found %v", sm, n, ok)
		}
		return n
	}
	live(s.Deep)
	if d := depth(s.Deep.ID); d != s.Depth || d != maxDepth || d < 4 {
		t.Errorf("deep node at depth %d, summary %d, tree %d", d, s.Depth, maxDepth)
	}
	live(s.Parent)
	liveKids := 0
	for _, c := range asIs.Children(s.Parent.ID) {
		if !c.Retired {
			liveKids++
		}
	}
	if liveKids < 2 || liveKids != s.LiveChildren {
		t.Errorf("parent has %d live children, summary %d", liveKids, s.LiveChildren)
	}
	live(s.Linked)
	in, out := 0, 0
	for _, lk := range links {
		if lk.Retired {
			continue
		}
		if lk.From == s.Linked.ID {
			out++
		}
		if lk.To == s.Linked.ID {
			in++
		}
	}
	if in == 0 || out == 0 || in != s.LinksIn || out != s.LinksOut {
		t.Errorf("linked node: %d in, %d out; summary %d, %d", in, out, s.LinksIn, s.LinksOut)
	}
	if n := live(s.Spare); n.Version != s.SpareVersion {
		t.Errorf("spare version %s, summary %s", n.Version, s.SpareVersion)
	}

	byID := map[id.ID]changeset.ChangeSet{}
	for _, cs := range sets {
		byID[cs.ID] = cs
	}
	for _, c := range []struct {
		cs, head id.ID
		state    changeset.State
	}{{s.Accepted, s.AcceptedHead, changeset.StateAccepted}, {s.Draft, s.DraftHead, changeset.StateDraft}} {
		cs, ok := byID[c.cs]
		if !ok || cs.State != c.state || len(cs.Heads) != 1 || cs.Heads[0] != c.head {
			t.Errorf("change set %s: want %s with head %s, got %+v", c.cs, c.state, c.head, cs)
		}
		for _, x := range cs.Ops {
			if x.Target() == s.Spare.ID {
				t.Errorf("in-flight %s touches the spare node", c.cs)
			}
		}
	}
	if s.Actor.IsZero() {
		t.Error("zero actor")
	}
}
