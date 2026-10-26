package api

import (
	"encoding/json"
	"math/rand"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/store"
)

// logBuilder writes change sets as log entries with deterministic IDs.
type logBuilder struct {
	t       *testing.T
	now     time.Time
	gen     *id.Generator
	actor   id.ID
	entries []entry.Entry
}

func newLogBuilder(t *testing.T) *logBuilder {
	b := &logBuilder{t: t, now: time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)}
	b.gen = id.NewGenerator(func() time.Time { return b.now }, rand.New(rand.NewSource(1)))
	b.actor = b.id()
	return b
}

func (b *logBuilder) id() id.ID {
	b.t.Helper()
	b.now = b.now.Add(time.Second)
	v, err := b.gen.New()
	if err != nil {
		b.t.Fatal(err)
	}
	return v
}

// changeSet opens a change set of the given kind and chains bodies after the
// open. It returns the change set's ID.
func (b *logBuilder) changeSet(kind changeset.Kind, bodies ...changeset.Body) id.ID {
	b.t.Helper()
	cs := b.id()
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
	return cs
}

func (b *logBuilder) landed(ops ...op.Op) id.ID {
	return b.changeSet(changeset.KindRecord, changeset.Ops{Ops: ops}, changeset.Land{})
}

func (b *logBuilder) draft(ops ...op.Op) id.ID {
	return b.changeSet(changeset.KindIntention, changeset.Ops{Ops: ops})
}

func (b *logBuilder) accepted(ops ...op.Op) id.ID {
	return b.changeSet(changeset.KindIntention, changeset.Ops{Ops: ops}, changeset.Propose{},
		changeset.Accept{Rule: "owner"})
}

func (b *logBuilder) model(current id.ID) Model {
	b.t.Helper()
	_, problems := changeset.FoldLog(store.Log{Entries: b.entries})
	if len(problems) > 0 {
		b.t.Fatalf("fold problems: %v", problems)
	}
	return newModel(store.Log{Entries: b.entries}, current)
}

func (b *logBuilder) create(kind op.Kind, handle string, fields map[string]json.RawMessage) op.Op {
	o := op.Op{ID: b.id(), Type: op.TypeNodeCreate, Node: b.id(), Kind: kind, Handle: handle, Fields: fields}
	if op.HasParent(kind) {
		o.Key = "i"
	}
	return o
}

func (b *logBuilder) update(node, base id.ID, set map[string]json.RawMessage) op.Op {
	return op.Op{ID: b.id(), Type: op.TypeNodeUpdate, Node: node, Base: base, Set: set}
}

func statement(s string) map[string]json.RawMessage {
	raw, _ := json.Marshal(s)
	return map[string]json.RawMessage{"statement": raw}
}

func TestModel(t *testing.T) {
	t.Run("Version", func(t *testing.T) {
		b := newLogBuilder(t)
		c := b.create(op.KindDecision, "use-go", statement("Go"))
		b.landed(c)
		u := b.update(c.Node, c.ID, statement("Go 1.25"))
		b.landed(u)

		n, err := b.model(id.ID{}).Resolve("use-go")
		if err != nil {
			t.Fatal(err)
		}
		if n.Version != u.ID || n.Conflict || n.Pending != nil {
			t.Errorf("version %s conflict %v pending %v, want %s false nil", n.Version, n.Conflict, n.Pending, u.ID)
		}
		if got := string(n.Fields["statement"]); got != `"Go 1.25"` {
			t.Errorf("statement %s", got)
		}
		if n.Kind != op.KindDecision || n.ID != c.Node {
			t.Errorf("node %+v", n)
		}
	})

	t.Run("CurrentVisible", func(t *testing.T) {
		b := newLogBuilder(t)
		c := b.create(op.KindDecision, "draft-node", nil)
		cs := b.draft(c)

		if _, err := b.model(id.ID{}).Resolve("draft-node"); ExitCode(err) != 4 {
			t.Fatalf("without current: err %v, want not found", err)
		}
		n, err := b.model(cs).Resolve("draft-node")
		if err != nil {
			t.Fatal(err)
		}
		if !n.Version.IsZero() {
			t.Errorf("version %s, want zero", n.Version)
		}
		if !slices.Equal(n.Pending, []id.ID{c.ID}) {
			t.Errorf("pending %v, want [%s]", n.Pending, c.ID)
		}
	})

	t.Run("CurrentOnTopOfLanded", func(t *testing.T) {
		b := newLogBuilder(t)
		c := b.create(op.KindDecision, "base", statement("old"))
		b.landed(c)
		u := b.update(c.Node, c.ID, statement("new"))
		cs := b.draft(u)

		n, ok := b.model(cs).Node(c.Node)
		if !ok {
			t.Fatal("node missing")
		}
		if n.Version != c.ID || string(n.Fields["statement"]) != `"new"` || !slices.Equal(n.Pending, []id.ID{u.ID}) {
			t.Errorf("node %+v", n)
		}
		n, _ = b.model(id.ID{}).Node(c.Node)
		if string(n.Fields["statement"]) != `"old"` {
			t.Errorf("as-is statement %s", n.Fields["statement"])
		}
	})

	t.Run("RetireFreesHandle", func(t *testing.T) {
		b := newLogBuilder(t)
		c := b.create(op.KindDecision, "gone", nil)
		b.landed(c)
		r := op.Op{ID: b.id(), Type: op.TypeNodeRetire, Node: c.Node, Base: c.ID, Reason: "no longer"}
		b.landed(r)
		m := b.model(id.ID{})

		if got, taken := m.HandleTaken("gone"); taken {
			t.Errorf("handle still taken by %s", got)
		}
		if _, err := m.Resolve("gone"); ExitCode(err) != 4 {
			t.Errorf("resolve by handle: %v, want not found", err)
		}
		n, err := m.Resolve(c.Node.String())
		if err != nil {
			t.Fatal(err)
		}
		if !n.Retired || n.Version != r.ID {
			t.Errorf("node %+v", n)
		}
	})

	t.Run("Prefix", func(t *testing.T) {
		b := newLogBuilder(t)
		c := b.create(op.KindDecision, "only", nil)
		b.landed(c)
		n, err := b.model(id.ID{}).Resolve(c.Node.String()[:8])
		if err != nil || n.ID != c.Node {
			t.Errorf("resolve prefix: %v %v", n.ID, err)
		}
		if _, err := b.model(id.ID{}).Resolve(strings.ToLower(c.Node.String()[:8])); ExitCode(err) != 4 {
			t.Errorf("lower-case prefix: %v, want not found", err)
		}
	})

	t.Run("Ambiguous", func(t *testing.T) {
		b := newLogBuilder(t)
		c1 := b.create(op.KindDecision, "one", nil)
		c2 := b.create(op.KindDecision, "two", nil)
		b.landed(c1, c2)
		prefix := commonPrefix(c1.Node.String(), c2.Node.String())
		if len(prefix) < 4 {
			t.Fatalf("test IDs share only %q", prefix)
		}
		_, err := b.model(id.ID{}).Resolve(prefix)
		if ExitCode(err) != 2 {
			t.Fatalf("err %v, exit %d, want 2", err, ExitCode(err))
		}
		for _, want := range []string{c1.Node.String() + " (one)", c2.Node.String() + " (two)"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("message %q lacks %q", err, want)
			}
		}
		if _, err := b.model(id.ID{}).Resolve(prefix[:3]); ExitCode(err) != 4 {
			t.Errorf("3-char prefix: %v, want not found", err)
		}
	})

	t.Run("Unknown", func(t *testing.T) {
		b := newLogBuilder(t)
		b.landed(b.create(op.KindDecision, "known", nil))
		m := b.model(id.ID{})
		for _, ref := range []string{"nosuch", "ZZZZZZ", b.id().String(), ""} {
			if _, err := m.Resolve(ref); ExitCode(err) != 4 {
				t.Errorf("Resolve(%q): %v, exit %d, want 4", ref, err, ExitCode(err))
			}
		}
		if _, err := m.ResolveLink("nosuch"); ExitCode(err) != 4 {
			t.Errorf("ResolveLink: %v, want not found", err)
		}
	})

	t.Run("Conflict", func(t *testing.T) {
		b := newLogBuilder(t)
		c := b.create(op.KindDecision, "fork", statement("a"))
		b.landed(c)
		u1 := b.update(c.Node, c.ID, statement("b"))
		b.landed(u1)
		u2 := b.update(c.Node, c.ID, statement("c"))
		b.landed(u2)

		n, err := b.model(id.ID{}).Resolve("fork")
		if err != nil {
			t.Fatal(err)
		}
		if !n.Conflict {
			t.Error("Conflict false, want true")
		}
		// The larger head wins, both for Version and for displayed fields.
		if n.Version != u2.ID || string(n.Fields["statement"]) != `"c"` {
			t.Errorf("version %s statement %s, want %s \"c\"", n.Version, n.Fields["statement"], u2.ID)
		}
	})

	t.Run("UnlandedNotAsIs", func(t *testing.T) {
		b := newLogBuilder(t)
		c := b.create(op.KindDecision, "pending", nil)
		cs := b.accepted(c)
		if _, err := b.model(id.ID{}).Resolve("pending"); ExitCode(err) != 4 {
			t.Errorf("accepted, not current: %v, want not found", err)
		}
		if _, taken := b.model(id.ID{}).HandleTaken("pending"); taken {
			t.Error("handle taken by an unlanded change set")
		}
		if _, err := b.model(cs).Resolve("pending"); err != nil {
			t.Errorf("accepted and current: %v", err)
		}
	})

	t.Run("LinkBetween", func(t *testing.T) {
		b := newLogBuilder(t)
		a := b.create(op.KindDecision, "a", nil)
		z := b.create(op.KindDecision, "z", nil)
		l := op.Op{ID: b.id(), Type: op.TypeLinkCreate, Link: b.id(), LinkType: op.LinkDependsOn, From: a.Node, To: z.Node}
		b.landed(a, z, l)
		m := b.model(id.ID{})

		got, err := m.LinkBetween(a.Node, z.Node, op.LinkDependsOn)
		if err != nil {
			t.Fatal(err)
		}
		if got.ID != l.Link || got.Version != l.ID || got.Type != op.LinkDependsOn || got.Retired {
			t.Errorf("link %+v", got)
		}
		if _, err := m.LinkBetween(z.Node, a.Node, op.LinkDependsOn); ExitCode(err) != 4 {
			t.Errorf("reversed: %v, want not found", err)
		}
		if got, err := m.ResolveLink(l.Link.String()[:10]); err != nil || got.ID != l.Link {
			t.Errorf("ResolveLink prefix: %v %v", got.ID, err)
		}

		r := op.Op{ID: b.id(), Type: op.TypeLinkRetire, Link: l.Link, Base: l.ID}
		cs := b.draft(r)
		m = b.model(cs)
		if _, err := m.LinkBetween(a.Node, z.Node, op.LinkDependsOn); ExitCode(err) != 4 {
			t.Errorf("retired in current: %v, want not found", err)
		}
		got, err = m.ResolveLink(l.Link.String())
		if err != nil || !got.Retired || got.Version != l.ID || !slices.Equal(got.Pending, []id.ID{r.ID}) {
			t.Errorf("retired link %+v, %v", got, err)
		}
	})

	t.Run("Children", func(t *testing.T) {
		b := newLogBuilder(t)
		root := b.create(op.KindCapability, "root", nil)
		mk := func(handle, key string) op.Op {
			o := b.create(op.KindCapability, handle, nil)
			o.Parent, o.Key = root.Node, key
			return o
		}
		x, y, w := mk("x", "r"), mk("y", "c"), mk("w", "r")
		other := b.create(op.KindIntent, "intent", nil)
		b.landed(root, x, y, w, other)
		m := b.model(id.ID{})

		var got []string
		for _, n := range m.Children(op.KindCapability, root.Node) {
			got = append(got, n.Handle)
		}
		// Key c first, then the two r keys in ID order (x was created before w).
		if want := []string{"y", "x", "w"}; !slices.Equal(got, want) {
			t.Errorf("children %v, want %v", got, want)
		}
		roots := m.Children(op.KindCapability, id.ID{})
		if len(roots) != 1 || roots[0].ID != root.Node {
			t.Errorf("roots %+v", roots)
		}
		if n := m.Children(op.KindIntent, id.ID{}); len(n) != 1 || n[0].Handle != "intent" {
			t.Errorf("intent roots %+v", n)
		}
	})

	t.Run("ActorByEmail", func(t *testing.T) {
		b := newLogBuilder(t)
		email := func(s string) map[string]json.RawMessage {
			raw, _ := json.Marshal(s)
			return map[string]json.RawMessage{"email": raw}
		}
		ann := b.create(op.KindActor, "ann", email("ann@example.com"))
		bob := b.create(op.KindDecision, "bob", email("bob@example.com"))
		b.landed(ann, bob)
		m := b.model(id.ID{})

		if n, ok := m.ActorByEmail("ann@example.com"); !ok || n.ID != ann.Node {
			t.Errorf("ann: %v %v", n.ID, ok)
		}
		if _, ok := m.ActorByEmail("bob@example.com"); ok {
			t.Error("found a non-actor node")
		}
		if _, ok := m.ActorByEmail("nobody@example.com"); ok {
			t.Error("found an unknown email")
		}
	})

	t.Run("HandleTaken", func(t *testing.T) {
		b := newLogBuilder(t)
		c := b.create(op.KindDecision, "taken", nil)
		b.landed(c)
		if got, ok := b.model(id.ID{}).HandleTaken("taken"); !ok || got != c.Node {
			t.Errorf("HandleTaken: %s %v", got, ok)
		}
	})

	t.Run("BadLogNeverPanics", func(t *testing.T) {
		b := newLogBuilder(t)
		c := b.create(op.KindDecision, "fine", nil)
		b.landed(c)
		// An update whose node was never created, and an entry with a body
		// that does not decode.
		orphan := b.update(b.id(), b.id(), statement("x"))
		b.landed(orphan)
		bad := b.entries[0]
		bad.ID, bad.Body = b.id(), json.RawMessage(`{"nope":1}`)
		entries := append(slices.Clone(b.entries), bad)

		m := newModel(store.Log{Entries: entries}, id.ID{})
		if _, err := m.Resolve("fine"); err != nil {
			t.Error(err)
		}
		if _, ok := m.Node(orphan.Node); ok {
			t.Error("a node with no create exists")
		}
	})

	// A change set sees the in-flight change sets it depends on, a
	// dependency before its dependents whatever their IDs.
	t.Run("Dependencies", func(t *testing.T) {
		b := newLogBuilder(t)
		c := b.create(op.KindCapability, "base", statement("A"))
		u := b.update(c.Node, c.ID, statement("B"))
		depB := b.draft(u) // B's ID sorts before A's
		depA := b.draft(c)
		top := b.changeSet(changeset.KindIntention)
		other := b.changeSet(changeset.KindIntention)
		b.landed(
			op.Op{ID: b.id(), Type: op.TypeLinkCreate, Link: b.id(), LinkType: op.LinkDependsOn, From: depB, To: depA},
			op.Op{ID: b.id(), Type: op.TypeLinkCreate, Link: b.id(), LinkType: op.LinkDependsOn, From: top, To: depB},
		)

		n, err := b.model(top).Resolve("base")
		if err != nil {
			t.Fatal(err)
		}
		if n.Version != u.ID || n.Conflict || n.Pending != nil || string(n.Fields["statement"]) != `"B"` {
			t.Errorf("transitive: %+v", n)
		}
		n, err = b.model(depB).Resolve("base")
		if err != nil {
			t.Fatal(err)
		}
		if n.Version != c.ID || n.Conflict || !slices.Equal(n.Pending, []id.ID{u.ID}) || string(n.Fields["statement"]) != `"B"` {
			t.Errorf("direct: %+v", n)
		}
		if n, err := b.model(depA).Resolve("base"); err != nil || !n.Version.IsZero() || !slices.Equal(n.Pending, []id.ID{c.ID}) {
			t.Errorf("own create: %+v, %v", n, err)
		}
		for _, cs := range []id.ID{other, {}} {
			if _, err := b.model(cs).Resolve("base"); ExitCode(err) != 4 {
				t.Errorf("no dependency (%s): %v, want not found", cs, err)
			}
		}
	})

	t.Run("DependencyCycle", func(t *testing.T) {
		b := newLogBuilder(t)
		c := b.create(op.KindCapability, "base", nil)
		u := b.update(c.Node, c.ID, statement("B"))
		depA := b.draft(c)
		depB := b.draft(u)
		top := b.changeSet(changeset.KindIntention)
		b.landed(
			op.Op{ID: b.id(), Type: op.TypeLinkCreate, Link: b.id(), LinkType: op.LinkDependsOn, From: depA, To: depB},
			op.Op{ID: b.id(), Type: op.TypeLinkCreate, Link: b.id(), LinkType: op.LinkDependsOn, From: depB, To: depA},
			op.Op{ID: b.id(), Type: op.TypeLinkCreate, Link: b.id(), LinkType: op.LinkDependsOn, From: top, To: depB},
		)
		n, err := b.model(top).Resolve("base")
		if err != nil || n.Version != u.ID || string(n.Fields["statement"]) != `"B"` {
			t.Errorf("cycle: %+v, %v", n, err)
		}
	})
}

func commonPrefix(a, b string) string {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return a[:i]
}
