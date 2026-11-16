package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/store"
)

// This file holds the parity cases the harness found on the index's write
// model (writeIndex.model), each on the smallest log that shows it: the
// write model now takes its view from the index's strict what-if, so each
// one agrees with the strict engine. The notes on leaf
// parity-write-index-from-strict name each case's results before the fix.

// changeSetAs is logBuilder.changeSet with the change set's ID given, for
// a log whose ops must name a change set before it is written.
func (b *logBuilder) changeSetAs(cs id.ID, kind changeset.Kind, bodies ...changeset.Body) {
	b.t.Helper()
	var prev id.ID
	for _, body := range append([]changeset.Body{changeset.Open{Kind: kind, Intent: "Test"}}, bodies...) {
		raw, err := changeset.EncodeBody(body)
		if err != nil {
			b.t.Fatal(err)
		}
		e := entry.Entry{V: changeset.EntryVersion(body), ID: b.id(), ChangeSet: cs, Actor: b.actor, For: b.actor,
			At: b.now, Type: body.Type(), Body: raw}
		if !prev.IsZero() {
			e.After = []id.ID{prev}
		}
		b.entries = append(b.entries, e)
		prev = e.ID
	}
}

// indexOf writes b's entries to a new repository and returns the index of
// it, caught up.
func indexOf(t *testing.T, b *logBuilder) *writeIndex {
	t.Helper()
	root := t.TempDir()
	gdir := filepath.Join(root, store.Dir)
	for _, e := range b.entries {
		data, err := entry.Encode(e)
		if err != nil {
			t.Fatal(err)
		}
		p := filepath.Join(gdir, store.LogDir, e.ChangeSet.String(), e.ID.String()+".json")
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	ctx := &Ctx{Root: root, IndexPath: filepath.Join(root, "index.db")}
	x, err := index.Open(ctx.Root, ctx.IndexPath)
	if err != nil {
		t.Fatalf("index.Open: %v", err)
	}
	ix := &writeIndex{ctx: ctx, x: x}
	t.Cleanup(ix.close)
	return ix
}

// indexAgrees fails t on any disagreement between the index's write model
// and the strict engine on the write view of cs, every node and link the
// log's ops name queried by ID.
func indexAgrees(t *testing.T, b *logBuilder, cs id.ID) {
	t.Helper()
	p := parityOf(t, b)
	ix := indexOf(t, b)
	m := ix.model(&workspace{}, cs).(*writeModel)
	if err := ix.failure(); err != nil {
		t.Fatalf("index model of %s: %v", viewName(cs), err)
	}
	var nodes, links []id.ID
	for _, s := range p.sets {
		for _, o := range s.Ops {
			if o.Node.IsZero() {
				links = append(links, o.Link)
			} else {
				nodes = append(nodes, o.Node)
			}
		}
	}
	got := parityView{nodes: map[id.ID]Node{}, links: map[id.ID]Link{}}
	for _, k := range slices.Compact(slices.SortedFunc(slices.Values(nodes), id.Compare)) {
		n, ok, err := m.ov.Node(k)
		if err != nil {
			t.Fatalf("index node %s: %v", k, err)
		}
		if ok {
			got.nodes[k] = m.node(n, true)
		}
	}
	for _, k := range slices.Compact(slices.SortedFunc(slices.Values(links), id.Compare)) {
		l, ok, err := m.ov.Link(k)
		if err != nil {
			t.Fatalf("index link %s: %v", k, err)
		}
		if ok {
			got.links[k] = m.link(l, true)
		}
	}
	strict, _, _ := strictView(p, cs)
	for _, d := range compareViews(viewName(cs), sideIndex, got, sideStrict, strict, false) {
		t.Errorf("%s\nlog:\n%s", d, p)
	}
}

// TestParityIndexCreateNoLongerApplies: a draft creates a node under a
// landed parent that a later landed change set retires. The strict what-if
// skips the create (model.ErrRetired), so the draft's write view has no
// such node; the index's write model kept it.
func TestParityIndexCreateNoLongerApplies(t *testing.T) {
	b := newLogBuilder(t)
	parent := b.create(op.KindCapability, "parent", nil)
	b.landed(parent)
	child := b.create(op.KindCapability, "child", nil)
	child.Parent = parent.Node
	draft := b.draft(child)
	b.landed(op.Op{ID: b.id(), Type: op.TypeNodeRetire, Node: parent.Node, Base: parent.ID, Reason: "gone"})
	indexAgrees(t, b, draft)
}

// TestParityIndexTerminalVersion: the write view of a landed change set is
// as-is with no current change set, so its own create is the node's
// Version; the index's write model left the change set's ops out.
func TestParityIndexTerminalVersion(t *testing.T) {
	b := newLogBuilder(t)
	c := b.create(op.KindCapability, "store", nil)
	landed := b.landed(c)
	indexAgrees(t, b, landed)
}

// TestParityIndexDisjointHeads: two landed updates of one node, both on
// its create, write different fields. The strict engine merges them with
// no conflict; the index's write model set Conflict on the two heads
// (Decision parity-version-rule).
func TestParityIndexDisjointHeads(t *testing.T) {
	b := newLogBuilder(t)
	c := b.create(op.KindCapability, "fork", statement("a"))
	b.landed(c)
	b.landed(b.update(c.Node, c.ID, statement("b")))
	b.landed(b.update(c.Node, c.ID, map[string]json.RawMessage{"acceptance": json.RawMessage(`"x"`)}))
	indexAgrees(t, b, id.ID{})
}

// TestParityIndexUnknownBase: a draft updates a landed node naming a base
// that is not in its history. The strict engine skips the update
// (model.ErrBase); the index's write model applied it.
func TestParityIndexUnknownBase(t *testing.T) {
	b := newLogBuilder(t)
	c := b.create(op.KindCapability, "store", statement("a"))
	b.landed(c)
	draft := b.draft(b.update(c.Node, b.id(), statement("b")))
	indexAgrees(t, b, draft)
}

// TestParityIndexCreateTaken: a draft creates a node whose ID a landed
// create already took. The strict engine skips it (model.ErrExists); the
// index's write model overwrote the node.
func TestParityIndexCreateTaken(t *testing.T) {
	b := newLogBuilder(t)
	c := b.create(op.KindCapability, "store", nil)
	b.landed(c)
	again := b.create(op.KindCapability, "other", nil)
	again.Node = c.Node
	draft := b.draft(again)
	indexAgrees(t, b, draft)
}

// TestParityIndexMissingEndpoint: a draft links a landed node to a node
// no change set creates. The strict engine skips the link
// (model.ErrMissing); the index's write model created it.
func TestParityIndexMissingEndpoint(t *testing.T) {
	b := newLogBuilder(t)
	c := b.create(op.KindCapability, "store", nil)
	b.landed(c)
	draft := b.draft(op.Op{ID: b.id(), Type: op.TypeLinkCreate, Link: b.id(), LinkType: op.LinkDerivedFrom,
		From: c.Node, To: b.id()})
	indexAgrees(t, b, draft)
}

// TestParityIndexParentWithoutChildren: a draft creates a capability under
// a landed decision, a kind with no children. The strict engine skips it
// (model.ErrParent); the index's write model created it.
func TestParityIndexParentWithoutChildren(t *testing.T) {
	b := newLogBuilder(t)
	d := b.create(op.KindDecision, "use-go", statement("Go"))
	b.landed(d)
	child := b.create(op.KindCapability, "child", nil)
	child.Parent = d.Node
	draft := b.draft(child)
	indexAgrees(t, b, draft)
}

// TestParityIndexMoveCycle: a draft moves a landed node under its own
// child. The strict engine skips the move (model.ErrCycle); the index's
// write model made the node its own ancestor.
func TestParityIndexMoveCycle(t *testing.T) {
	b := newLogBuilder(t)
	a := b.create(op.KindCapability, "a", nil)
	c := b.create(op.KindCapability, "c", nil)
	c.Parent = a.Node
	b.landed(a, c)
	draft := b.draft(op.Op{ID: b.id(), Type: op.TypeNodeMove, Node: a.Node, Base: a.ID, Parent: c.Node, Key: "i"})
	indexAgrees(t, b, draft)
}

// TestParityIndexDependencyApplyOrder: a draft depends on a later draft
// (a larger ID, so the strict engine applies the depending draft first),
// and both set one field of a landed node from its create. The strict
// what-if ends with the dependency's value, a conflict on the dependency's
// stale write; the index's write model applied every dependency's ops
// before the current change set's, so it ended with the current one's.
func TestParityIndexDependencyApplyOrder(t *testing.T) {
	b := newLogBuilder(t)
	c := b.create(op.KindCapability, "store", statement("a"))
	b.landed(c)
	cur := b.id()
	dep := b.draft(b.update(c.Node, c.ID, statement("dep")))
	b.changeSetAs(cur, changeset.KindIntention, changeset.Ops{Ops: []op.Op{
		b.update(c.Node, c.ID, statement("cur")),
		{ID: b.id(), Type: op.TypeLinkCreate, Link: b.id(), LinkType: op.LinkDependsOn, From: cur, To: dep},
	}})
	indexAgrees(t, b, cur)
}
