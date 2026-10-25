package render

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/model"
	"github.com/dphaener/genesis/internal/op"
	"github.com/dphaener/genesis/internal/store"
)

// -update rewrites the golden files under testdata from the current output:
// go test ./internal/render -run Golden -update
var update = flag.Bool("update", false, "rewrite the golden files under testdata")

// corpusDir holds sample logs from every released format; see its README.
const corpusDir = "../../testdata/corpus"

// fid returns a fixed, valid ID ending in the decimal digits of n.
func fid(n int) id.ID {
	v, err := id.Parse(fmt.Sprintf("01K8Z%021d", n))
	if err != nil {
		panic(err)
	}
	return v
}

func str(s string) json.RawMessage {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return b
}

func landed(cs id.ID, intent string, ops ...op.Op) changeset.ChangeSet {
	return changeset.ChangeSet{ID: cs, State: changeset.StateLanded, Open: changeset.Open{Intent: intent}, Ops: ops}
}

func create(o, node id.ID, k op.Kind, handle string, parent id.ID, key string, fields map[string]json.RawMessage) op.Op {
	return op.Op{ID: o, Type: op.TypeNodeCreate, Node: node, Kind: k, Handle: handle, Parent: parent, Key: key, Fields: fields}
}

func link(o, l id.ID, t op.LinkType, from, to id.ID, note string) op.Op {
	return op.Op{ID: o, Type: op.TypeLinkCreate, Link: l, LinkType: t, From: from, To: to, Note: note}
}

// The hand-built model's IDs: change sets 1xx, nodes 2xx, links 3xx, ops
// the rest.
var (
	csSeed, csWire, csDraft = fid(101), fid(102), fid(103)

	nCheckout, nCart, nPayment, nCard = fid(201), fid(202), fid(203), fid(204)
	nSell, nPCI, nPCIDup, nStripe     = fid(205), fid(206), fid(207), fid(208)
	nWhich, nLegacy, nLegacyChild     = fid(209), fid(210), fid(211)
	nAuditor, nGone                   = fid(212), fid(213)

	lServes, lConstrains, lDecided, lAnswers = fid(301), fid(302), fid(303), fid(304)
	lRelates, lRetired, lSelf                = fid(305), fid(306), fid(307)
)

// handBuilt returns change sets whose as-is covers every branch of a page
// and the index: a nested tree three deep, a duplicate handle, retired nodes
// (one with a live child), a change set endpoint, notes, a non-string field,
// links in and out, a retired link and a node touched by two change sets.
func handBuilt() []changeset.ChangeSet {
	seed := landed(csSeed, "Seed the coffee shop",
		create(fid(1), nCheckout, op.KindCapability, "checkout", id.ID{}, "i", map[string]json.RawMessage{
			"statement": str("Customers pay **online**.\r\n\nNo cash.\n"),
			"limits":    json.RawMessage(`{"max":3,"note":"three ` + "```" + ` ticks"}`),
		}),
		create(fid(2), nCart, op.KindCapability, "cart", nCheckout, "i", nil),
		create(fid(3), nPayment, op.KindCapability, "payment", nCheckout, "r", nil),
		create(fid(4), nCard, op.KindCapability, "card", nPayment, "i", map[string]json.RawMessage{
			"empty": str(""),
		}),
		create(fid(5), nSell, op.KindIntent, "sell-coffee", id.ID{}, "i", map[string]json.RawMessage{
			"statement": str("Sell coffee."),
		}),
		create(fid(6), nPCI, op.KindConstraint, "pci", id.ID{}, "", nil),
		create(fid(7), nPCIDup, op.KindConstraint, "pci", id.ID{}, "", nil),
		create(fid(8), nStripe, op.KindDecision, "use-stripe", id.ID{}, "", nil),
		create(fid(9), nWhich, op.KindQuestion, "which-psp", id.ID{}, "", nil),
		create(fid(10), nLegacy, op.KindCapability, "legacy", id.ID{}, "r", nil),
		create(fid(11), nLegacyChild, op.KindCapability, "kiosk", nLegacy, "i", nil),
		create(fid(12), nAuditor, op.KindActor, "auditor", id.ID{}, "", nil),
		create(fid(13), nGone, op.KindSource, "old-notes", id.ID{}, "", nil),
	)
	wire := landed(csWire, "Wire it\nup",
		link(fid(20), lServes, op.LinkServes, nCheckout, nSell, "the main path"),
		link(fid(21), lConstrains, op.LinkConstrains, nPCI, nCheckout, ""),
		link(fid(22), lDecided, op.LinkDecidedBy, nCheckout, nStripe, ""),
		link(fid(23), lAnswers, op.LinkAnswers, nStripe, nWhich, "settled"),
		link(fid(24), lRelates, op.LinkRelatesTo, nCheckout, csDraft, ""),
		link(fid(25), lRetired, op.LinkRelatesTo, nCheckout, nCart, ""),
		link(fid(26), lSelf, op.LinkRelatesTo, nAuditor, nAuditor, "me, myself"),
		op.Op{ID: fid(30), Type: op.TypeNodeUpdate, Node: nCheckout, Base: fid(1),
			Set: map[string]json.RawMessage{"owner": str("payments team")}},
		op.Op{ID: fid(31), Type: op.TypeNodeRetire, Node: nWhich, Base: fid(9), Reason: "answered by\nuse-stripe"},
		op.Op{ID: fid(32), Type: op.TypeNodeRetire, Node: nLegacy, Base: fid(10), Reason: "folded into checkout"},
		op.Op{ID: fid(33), Type: op.TypeNodeRetire, Node: nGone, Base: fid(13)},
		op.Op{ID: fid(34), Type: op.TypeLinkRetire, Link: lRetired, Base: fid(25)},
	)
	draft := changeset.ChangeSet{ID: csDraft, State: changeset.StateDraft,
		Open: changeset.Open{Intent: "Try a gift card flow"},
		Ops:  []op.Op{create(fid(40), fid(220), op.KindCapability, "gift-card", nCheckout, "z", nil)}}
	return []changeset.ChangeSet{seed, wire, draft}
}

// corpusCases returns the corpus case folders, as "format-N/<case>" names
// mapped to their paths.
func corpusCases(t *testing.T) map[string]string {
	t.Helper()
	roots, err := filepath.Glob(filepath.Join(corpusDir, "format-*", "*"))
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, r := range roots {
		if info, err := os.Stat(r); err != nil || !info.IsDir() {
			continue
		}
		out[filepath.Base(filepath.Dir(r))+"/"+filepath.Base(r)] = r
	}
	if len(out) == 0 {
		t.Fatalf("no corpus cases under %s", corpusDir)
	}
	return out
}

// corpusSets reads and folds one corpus case, failing on any problem.
func corpusSets(t *testing.T, root string) []changeset.ChangeSet {
	t.Helper()
	l := store.ReadLog(os.DirFS(filepath.Join(root, store.Dir)))
	for _, p := range l.Problems {
		t.Fatalf("store problem at %s: %v", p.Path, p.Err)
	}
	sets, ps := changeset.FoldLog(l)
	for _, p := range ps {
		t.Fatalf("fold problem: %v", p)
	}
	return sets
}

func renderSets(t *testing.T, sets []changeset.ChangeSet) map[string][]byte {
	t.Helper()
	m, _ := model.AsIs(sets)
	return Render(m, sets)
}

// TestGolden renders every corpus case and the hand-built model and
// compares the files with testdata/<name>/, file by file.
func TestGolden(t *testing.T) {
	cases := map[string][]changeset.ChangeSet{}
	for name, root := range corpusCases(t) {
		cases[name] = corpusSets(t, root)
	}
	cases["hand-built"] = handBuilt()
	for name, sets := range cases {
		t.Run(name, func(t *testing.T) {
			got := renderSets(t, sets)
			checkFiles(t, got)
			dir := filepath.Join("testdata", filepath.FromSlash(name))
			if *update {
				writeGolden(t, dir, got)
			}
			compareGolden(t, dir, got)
		})
	}
}

// checkFiles checks the rules every rendered file keeps: LF only, exactly
// one final newline, no absolute paths in links.
func checkFiles(t *testing.T, files map[string][]byte) {
	t.Helper()
	if _, ok := files["index.md"]; !ok {
		t.Error("no index.md")
	}
	for p, b := range files {
		if strings.HasPrefix(p, "/") || strings.Contains(p, "\\") || strings.Contains(p, "..") {
			t.Errorf("%s: not a clean relative slash path", p)
		}
		if bytes.ContainsRune(b, '\r') {
			t.Errorf("%s: contains CR", p)
		}
		if !bytes.HasSuffix(b, []byte("\n")) || bytes.HasSuffix(b, []byte("\n\n")) {
			t.Errorf("%s: does not end with exactly one newline", p)
		}
		if bytes.Contains(b, []byte("](/")) {
			t.Errorf("%s: has an absolute link", p)
		}
	}
}

func writeGolden(t *testing.T, dir string, files map[string][]byte) {
	t.Helper()
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	for p, b := range files {
		full := filepath.Join(dir, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, b, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// compareGolden names every file that differs, is missing or is extra.
func compareGolden(t *testing.T, dir string, got map[string][]byte) {
	t.Helper()
	want := map[string][]byte{}
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		b, err := os.ReadFile(p)
		want[filepath.ToSlash(rel)] = b
		return err
	})
	if err != nil {
		t.Fatalf("reading golden files (run with -update to create them): %v", err)
	}
	for p, b := range got {
		w, ok := want[p]
		switch {
		case !ok:
			t.Errorf("%s: rendered but not in the golden files (run with -update)", filepath.Join(dir, p))
		case !bytes.Equal(w, b):
			t.Errorf("%s: differs from the golden file (run with -update)\n--- got\n%s\n--- want\n%s",
				filepath.Join(dir, p), b, w)
		}
	}
	for p := range want {
		if _, ok := got[p]; !ok {
			t.Errorf("%s: golden file not rendered", filepath.Join(dir, p))
		}
	}
}

func TestEmpty(t *testing.T) {
	for _, sets := range [][]changeset.ChangeSet{nil, {{ID: fid(1), State: changeset.StateDraft,
		Ops: []op.Op{create(fid(2), fid(3), op.KindActor, "a", id.ID{}, "", nil)}}}} {
		got := renderSets(t, sets)
		want := map[string][]byte{"index.md": []byte("# Model\n\nNothing has landed yet.\n")}
		if len(got) != 1 || !bytes.Equal(got["index.md"], want["index.md"]) {
			t.Errorf("Render = %q, want %q", got, want)
		}
	}
}

// TestPaths checks the page paths of the hand-built model: the tree, the
// duplicate handle and no pages for retired nodes.
func TestPaths(t *testing.T) {
	got := renderSets(t, handBuilt())
	var paths []string
	for p := range got {
		paths = append(paths, p)
	}
	slices.Sort(paths)
	want := []string{
		"actor/auditor.md",
		"capability/checkout.md",
		"capability/checkout/cart.md",
		"capability/checkout/payment.md",
		"capability/checkout/payment/card.md",
		"capability/legacy/kiosk.md",
		"constraint/pci." + nPCIDup.String() + ".md",
		"constraint/pci.md",
		"decision/use-stripe.md",
		"index.md",
		"intent/sell-coffee.md",
	}
	if !slices.Equal(paths, want) {
		t.Errorf("paths = %q\nwant %q", paths, want)
	}
}

// TestHandBuiltProblems checks the hand-built model applies as meant: its
// only problem is the duplicate handle.
func TestHandBuiltProblems(t *testing.T) {
	_, ps := model.AsIs(handBuilt())
	if len(ps) != 1 || !errors.Is(ps[0].Err, model.ErrHandle) {
		t.Errorf("problems = %v, want one duplicate handle", ps)
	}
}

func TestRel(t *testing.T) {
	tests := []struct{ from, to, want string }{
		{"index.md", "capability/a.md", "capability/a.md"},
		{"capability/a.md", "capability/a/b.md", "a/b.md"},
		{"capability/a/b.md", "capability/a.md", "../a.md"},
		{"capability/a/b/c.md", "intent/x.md", "../../../intent/x.md"},
		{"constraint/p.md", "constraint/q.md", "q.md"},
		{"capability/a/b.md", "capability/a/c/d.md", "c/d.md"},
	}
	for _, tc := range tests {
		if got := rel(tc.from, tc.to); got != tc.want {
			t.Errorf("rel(%q, %q) = %q, want %q", tc.from, tc.to, got, tc.want)
		}
	}
}

// TestNoAliasing checks Render keeps no reference into its output: two
// calls give equal, separate bytes.
func TestNoAliasing(t *testing.T) {
	sets := handBuilt()
	a, b := renderSets(t, sets), renderSets(t, sets)
	for p := range a {
		if len(a[p]) > 0 {
			a[p][0] = 'X'
		}
		if bytes.HasPrefix(b[p], []byte("X")) {
			t.Errorf("%s: outputs share bytes", p)
		}
	}
}

// BenchmarkRender renders a synthetic model of n nodes, half capabilities in
// a shallow tree, each with two outgoing links. ns/op should grow roughly
// linearly with n.
func BenchmarkRender(b *testing.B) {
	for _, n := range []int{1000, 4000} {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			sets := synthetic(n)
			m, _ := model.AsIs(sets)
			b.ResetTimer()
			for b.Loop() {
				Render(m, sets)
			}
		})
	}
}

// synthetic returns one landed change set creating n nodes and 2n links.
func synthetic(n int) []changeset.ChangeSet {
	var ops []op.Op
	next := 1000
	nid := func() id.ID { next++; return fid(next) }
	nodes := make([]id.ID, n)
	for i := range n {
		nodes[i] = nid()
		if i%2 == 0 {
			var parent id.ID
			if i >= 20 {
				parent = nodes[i%20]
			}
			ops = append(ops, create(nid(), nodes[i], op.KindCapability, fmt.Sprintf("c%d", i), parent, "i", nil))
		} else {
			ops = append(ops, create(nid(), nodes[i], op.KindDecision, fmt.Sprintf("d%d", i), id.ID{}, "", nil))
		}
	}
	for i := range n {
		ops = append(ops,
			link(nid(), nid(), op.LinkRelatesTo, nodes[i], nodes[(i+1)%n], ""),
			link(nid(), nid(), op.LinkDecidedBy, nodes[i], nodes[(i*7+1)%n], "note"))
	}
	return []changeset.ChangeSet{landed(fid(1), "synthetic", ops...)}
}
