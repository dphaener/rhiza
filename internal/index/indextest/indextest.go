// Package indextest draws random logs for property tests: internal/index's
// own (TestPropertyIncremental and the tests that share its generator) and
// internal/api's projection parity test. Only tests import it.
//
// It writes and removes the drawn files but never opens an index: each test
// catches its own index up after each step.
package indextest

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/store"
	"pgregory.net/rapid"
)

// TB is what *testing.T and *rapid.T share.
type TB interface {
	Helper()
	Fatalf(format string, args ...any)
}

// mkID is fixture ID n of family f: IDs of one family grow with n. The
// families match internal/index's own test fixtures.
func mkID(f byte, n int) id.ID {
	var i id.ID
	i[0] = f
	binary.BigEndian.PutUint32(i[12:], uint32(n))
	return i
}

// ChangeSet, Entry, Op, Node and Link are the fixture IDs of each family;
// Actor is the one actor every drawn entry is stamped with.
func ChangeSet(n int) id.ID { return mkID(2, n) }
func Entry(n int) id.ID     { return mkID(3, n) }
func Op(n int) id.ID        { return mkID(4, n) }
func Node(n int) id.ID      { return mkID(5, n) }
func Link(n int) id.ID      { return mkID(6, n) }
func Actor() id.ID          { return mkID(7, 1) }

// At is the time every drawn entry is stamped with.
var At = time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)

// Nodes and Links are every node and link ID a drawn log's ops target.
var (
	Nodes = []id.ID{Node(1), Node(2), Node(3)}
	Links = []id.ID{Link(1), Link(2)}
)

// File is one file of a drawn log: its path under .rhiza (slash form), its
// bytes and, for an entry, its type.
type File struct {
	Path string
	Data []byte
	Type entry.Type
}

var (
	opTypes = []op.Type{op.TypeNodeCreate, op.TypeNodeCreate, op.TypeNodeUpdate, op.TypeNodeUpdate,
		op.TypeNodeMove, op.TypeNodeRetire, op.TypeLinkCreate, op.TypeLinkRetire}
	// states weights landed and accepted up.
	states = []changeset.State{changeset.StateLanded, changeset.StateLanded, changeset.StateLanded,
		changeset.StateAccepted, changeset.StateAccepted, changeset.StateProposed, changeset.StateDraft,
		changeset.StateRejected, changeset.StateWithdrawn, changeset.StateSuperseded}
)

// Log is a drawn log: its files, and what drawing more entries on it needs.
type Log struct {
	Files []File
	// ChangeSets is the IDs of the drawn change sets, in the order drawn.
	ChangeSets []id.ID
	// States is the state each change set was moved to.
	States map[id.ID]changeset.State

	entries, ops int
	onTarget     map[id.ID][]id.ID
	drawn        []id.ID
	ends         []id.ID
	heads        map[id.ID]id.ID
}

// base draws a base for target: usually one of its ops, sometimes any op
// or one that never existed.
func (g *Log) base(t *rapid.T, target id.ID) id.ID {
	switch k := rapid.IntRange(0, 9).Draw(t, "base from"); {
	case k > 1 && len(g.onTarget[target]) > 0:
		return rapid.SampledFrom(g.onTarget[target]).Draw(t, "base")
	case k == 1 && len(g.drawn) > 0:
		return rapid.SampledFrom(g.drawn).Draw(t, "any base")
	}
	return Op(9999)
}

// drawOp draws one operation, or ok false when it is not valid.
func (g *Log) drawOp(t *rapid.T) (o op.Op, ok bool) {
	g.ops++
	o = op.Op{ID: Op(g.ops)}
	pick := func(label string, from []id.ID) id.ID { return rapid.SampledFrom(from).Draw(t, label) }
	switch o.Type = rapid.SampledFrom(opTypes).Draw(t, "type"); o.Type {
	case op.TypeNodeCreate:
		o.Node = pick("node", Nodes)
		o.Kind = rapid.SampledFrom([]op.Kind{op.KindCapability, op.KindCapability, op.KindDecision}).Draw(t, "kind")
		o.Handle = rapid.SampledFrom([]string{"a", "b", "c"}).Draw(t, "handle")
		if op.HasParent(o.Kind) {
			if rapid.Bool().Draw(t, "child?") {
				o.Parent = pick("parent", Nodes)
			}
			o.Key = rapid.SampledFrom([]string{"8", "i", "r"}).Draw(t, "key")
		}
		if rapid.Bool().Draw(t, "fields?") {
			o.Fields = fieldsOf("statement", fmt.Sprintf(`"s%d"`, g.ops))
		}
	case op.TypeNodeUpdate:
		o.Node = pick("node", Nodes)
		o.Base = g.base(t, o.Node)
		if rapid.Bool().Draw(t, "handle?") {
			o.Handle = rapid.SampledFrom([]string{"a", "b", "c"}).Draw(t, "handle")
		}
		name := rapid.SampledFrom([]string{"statement", "statement", "title"}).Draw(t, "set")
		o.Set = fieldsOf(name, fmt.Sprintf("{\n  \"n\": %d,\n  \"t\": \"<&>\"\n}", g.ops))
	case op.TypeNodeMove:
		o.Node = pick("node", Nodes)
		o.Base = g.base(t, o.Node)
		if rapid.Bool().Draw(t, "child?") {
			o.Parent = pick("parent", Nodes)
		}
		o.Key = rapid.SampledFrom([]string{"8", "i", "r"}).Draw(t, "key")
	case op.TypeNodeRetire:
		o.Node = pick("node", Nodes)
		o.Base = g.base(t, o.Node)
		o.Reason = "gone"
	case op.TypeLinkCreate:
		o.Link = pick("link", Links)
		o.LinkType = rapid.SampledFrom([]op.LinkType{op.LinkServes, op.LinkRelatesTo, op.LinkDependsOn, op.LinkDependsOn}).Draw(t, "link type")
		o.From = pick("from", g.ends)
		o.To = pick("to", g.ends)
		if o.LinkType == op.LinkRelatesTo || rapid.Bool().Draw(t, "note?") {
			o.Note = "why" // required on relates-to
		}
	case op.TypeLinkRetire:
		o.Link = pick("link", Links)
		o.Base = g.base(t, o.Link)
	}
	if o.Validate() != nil {
		return op.Op{}, false
	}
	g.onTarget[o.Target()] = append(g.onTarget[o.Target()], o.ID)
	g.drawn = append(g.drawn, o.ID)
	return o, true
}

func fieldsOf(name, value string) map[string]json.RawMessage {
	return map[string]json.RawMessage{name: json.RawMessage(value)}
}

// EntryPath is the slash path of an entry file under .rhiza.
func EntryPath(cs, eid id.ID) string {
	return store.LogDir + "/" + cs.String() + "/" + eid.String() + ".json"
}

// NewEntry builds the entry eid of cs with body b after after (zero for
// none), stamped with Actor and At.
func NewEntry(t TB, eid, cs, after id.ID, b changeset.Body) entry.Entry {
	t.Helper()
	body, err := changeset.EncodeBody(b)
	if err != nil {
		t.Fatalf("encode body: %v", err)
	}
	var a []id.ID
	if !after.IsZero() {
		a = []id.ID{after}
	}
	return entry.Entry{V: changeset.EntryVersion(b), ID: eid, ChangeSet: cs, Actor: Actor(), For: Actor(),
		At: At, After: a, Type: b.Type(), Body: body}
}

// next builds the next entry of cs, after its last one, and moves its head.
func (g *Log) next(t *rapid.T, cs id.ID, b changeset.Body) entry.Entry {
	g.entries++
	e := NewEntry(t, Entry(g.entries), cs, g.heads[cs], b)
	g.heads[cs] = e.ID
	return e
}

// add encodes the next entry of cs and records its file.
func (g *Log) add(t *rapid.T, cs id.ID, b changeset.Body) {
	g.addEntry(t, g.next(t, cs, b))
}

func (g *Log) addEntry(t *rapid.T, e entry.Entry) {
	data, err := entry.Encode(e)
	if err != nil {
		t.Fatalf("encode entry: %v", err)
	}
	g.Files = append(g.Files, File{Path: EntryPath(e.ChangeSet, e.ID), Data: data, Type: e.Type})
}

// Draw draws 1 to 5 change sets, each opened, given ops entries and moved
// through the lifecycle to a drawn state, plus now and then a corrupt entry
// file, a stray file or an entry out of order.
func Draw(t *rapid.T) *Log {
	n := rapid.IntRange(1, 5).Draw(t, "change sets")
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i + 1
	}
	nums = rapid.Permutation(nums).Draw(t, "change set IDs")
	g := &Log{onTarget: map[id.ID][]id.ID{}, heads: map[id.ID]id.ID{}, States: map[id.ID]changeset.State{},
		ends: slices.Clone(Nodes)}
	for _, k := range nums {
		g.ends = append(g.ends, ChangeSet(k))
	}

	for i, k := range nums {
		cs := ChangeSet(k)
		g.ChangeSets = append(g.ChangeSets, cs)
		kind := changeset.KindIntention
		if rapid.IntRange(0, 4).Draw(t, "record?") == 0 {
			kind = changeset.KindRecord
		}
		g.add(t, cs, changeset.Open{Kind: kind, Intent: fmt.Sprintf("set %d", k)})
		if i == 0 && rapid.IntRange(0, 3).Draw(t, "seed?") > 0 {
			var seed []op.Op
			for _, node := range Nodes {
				g.ops++
				o := op.Op{ID: Op(g.ops), Type: op.TypeNodeCreate, Node: node, Kind: op.KindCapability,
					Handle: fmt.Sprintf("n%d", g.ops), Key: "i"}
				g.onTarget[node] = append(g.onTarget[node], o.ID)
				g.drawn = append(g.drawn, o.ID)
				seed = append(seed, o)
			}
			g.add(t, cs, changeset.Ops{Ops: seed})
		}
		for range rapid.IntRange(0, 3).Draw(t, "ops entries") {
			if ops := g.drawOps(t); len(ops) > 0 {
				g.add(t, cs, changeset.Ops{Ops: ops})
			}
		}

		state := rapid.SampledFrom(states).Draw(t, "state")
		g.States[cs] = state
		reach := func(s changeset.State) {
			switch s {
			case changeset.StateProposed:
				g.add(t, cs, changeset.Propose{})
			case changeset.StateAccepted:
				g.add(t, cs, changeset.Propose{})
				g.add(t, cs, changeset.Accept{Rule: "owner"})
			}
		}
		switch state {
		case changeset.StateProposed, changeset.StateAccepted:
			reach(state)
		case changeset.StateLanded:
			if kind == changeset.KindIntention {
				reach(changeset.StateAccepted)
			}
			g.add(t, cs, changeset.Land{})
		case changeset.StateRejected:
			reach(changeset.StateProposed)
			g.add(t, cs, changeset.Reject{Rule: "owner", Reason: "no"})
		case changeset.StateWithdrawn:
			g.add(t, cs, changeset.Withdraw{})
		case changeset.StateSuperseded:
			// Sometimes by a change set not in the log.
			by := rapid.SampledFrom(append([]id.ID{ChangeSet(99)}, g.ends[len(Nodes):]...)).Draw(t, "by")
			if by == cs {
				by = ChangeSet(99)
			}
			g.add(t, cs, changeset.Supersede{By: by})
		}
		if rapid.IntRange(0, 9).Draw(t, "late ops?") == 0 {
			if o, ok := g.drawOp(t); ok {
				g.add(t, cs, changeset.Ops{Ops: []op.Op{o}}) // skipped after a terminal state
			}
		}
		if rapid.IntRange(0, 9).Draw(t, "orphan?") == 0 {
			// After names no entry: an order problem. The change set's head
			// does not move.
			g.entries++
			g.addEntry(t, NewEntry(t, Entry(g.entries), cs, Op(8888), changeset.Note{Text: "lost"}))
		}
		if rapid.IntRange(0, 9).Draw(t, "corrupt?") == 0 {
			g.entries++
			g.Files = append(g.Files, File{Path: EntryPath(cs, Entry(g.entries)), Data: []byte("{not json")})
		}
	}
	if rapid.IntRange(0, 9).Draw(t, "stray?") == 0 {
		g.Files = append(g.Files, File{Path: "log/stray.txt", Data: []byte("x")})
	}
	return g
}

// drawOps draws one to three ops, leaving out those that are not valid.
func (g *Log) drawOps(t *rapid.T) []op.Op {
	var ops []op.Op
	for range rapid.IntRange(1, 3).Draw(t, "ops") {
		if o, ok := g.drawOp(t); ok {
			ops = append(ops, o)
		}
	}
	return ops
}

// Pending draws an ops entry of cs that is not among Files, after cs's
// last entry, as a dry run would hold it in pending; ok is false when no
// valid op was drawn. Each call draws the next one after the last.
func (g *Log) Pending(t *rapid.T, cs id.ID) (e entry.Entry, ok bool) {
	ops := g.drawOps(t)
	if len(ops) == 0 {
		return entry.Entry{}, false
	}
	return g.next(t, cs, changeset.Ops{Ops: ops}), true
}

// WriteFile writes data at path p (slash form) under the .rhiza folder
// gdir.
func WriteFile(t TB, gdir, p string, data []byte) {
	t.Helper()
	full := filepath.Join(gdir, filepath.FromSlash(p))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(full, data, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

// Steps splits writing a drawn log under a .rhiza folder into steps: each
// adds some of the files not yet written or, at most twice per log, removes
// some already written (which makes an index rebuild) to add them back
// later.
type Steps struct {
	gdir      string
	files     []File
	remaining []int
	written   []int
	removals  int
}

// Step is what one step did.
type Step struct {
	// Removed is true when the step removed files rather than added them.
	Removed bool
	// Landed and Accepted are true when it added a land or an accept entry.
	Landed, Accepted bool
}

// NewSteps writes files under gdir step by step.
func NewSteps(gdir string, files []File) *Steps {
	s := &Steps{gdir: gdir, files: files, remaining: make([]int, len(files))}
	for i := range s.remaining {
		s.remaining[i] = i
	}
	return s
}

// Done is true once every file is written.
func (s *Steps) Done() bool { return len(s.remaining) == 0 }

// Next draws and takes one step.
func (s *Steps) Next(t *rapid.T) Step {
	if len(s.written) > 0 && s.removals < 2 && rapid.IntRange(0, 3).Draw(t, "remove?") == 0 {
		s.removals++
		out := rapid.SliceOfNDistinct(rapid.SampledFrom(s.written), 1, len(s.written), rapid.ID[int]).Draw(t, "remove")
		for _, i := range out {
			if err := os.Remove(filepath.Join(s.gdir, filepath.FromSlash(s.files[i].Path))); err != nil {
				t.Fatal(err)
			}
			s.written = slices.DeleteFunc(s.written, func(w int) bool { return w == i })
			s.remaining = append(s.remaining, i)
		}
		return Step{Removed: true}
	}
	// Small steps half the time, so that orders that matter (a link's
	// endpoint change set arriving after the link lands) come up often.
	most := len(s.remaining)
	if rapid.Bool().Draw(t, "small step?") {
		most = min(most, 2)
	}
	var st Step
	in := rapid.SliceOfNDistinct(rapid.SampledFrom(s.remaining), 1, most, rapid.ID[int]).Draw(t, "add")
	for _, i := range in {
		WriteFile(t, s.gdir, s.files[i].Path, s.files[i].Data)
		s.remaining = slices.DeleteFunc(s.remaining, func(r int) bool { return r == i })
		s.written = append(s.written, i)
		st.Landed = st.Landed || s.files[i].Type == entry.TypeLand
		st.Accepted = st.Accepted || s.files[i].Type == entry.TypeAccept
	}
	return st
}
