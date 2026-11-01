package index

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/store"
	"pgregory.net/rapid"
)

// logFile is one file of a drawn log: its path under .rhiza, its bytes
// and, for an entry, its type.
type logFile struct {
	path string
	data []byte
	typ  entry.Type
}

// logGen draws a small log over three nodes and two links.
type logGen struct {
	entries, ops int
	onTarget     map[id.ID][]id.ID
	drawn        []id.ID
	files        []logFile
}

var (
	propNodes = []id.ID{nd(1), nd(2), nd(3)}
	propLinks = []id.ID{ln(1), ln(2)}
	propTypes = []op.Type{op.TypeNodeCreate, op.TypeNodeCreate, op.TypeNodeUpdate, op.TypeNodeUpdate,
		op.TypeNodeMove, op.TypeNodeRetire, op.TypeLinkCreate, op.TypeLinkRetire}
	// propStates weights landed and accepted up.
	propStates = []changeset.State{changeset.StateLanded, changeset.StateLanded, changeset.StateLanded,
		changeset.StateAccepted, changeset.StateAccepted, changeset.StateProposed, changeset.StateDraft,
		changeset.StateRejected, changeset.StateWithdrawn, changeset.StateSuperseded}
)

// base draws a base for target: usually one of its ops, sometimes any op
// or one that never existed.
func (g *logGen) base(t *rapid.T, target id.ID) id.ID {
	switch k := rapid.IntRange(0, 9).Draw(t, "base from"); {
	case k > 1 && len(g.onTarget[target]) > 0:
		return rapid.SampledFrom(g.onTarget[target]).Draw(t, "base")
	case k == 1 && len(g.drawn) > 0:
		return rapid.SampledFrom(g.drawn).Draw(t, "any base")
	}
	return opID(9999)
}

// drawOp draws one operation, or ok false when it is not valid.
func (g *logGen) drawOp(t *rapid.T, ends []id.ID) (o op.Op, ok bool) {
	g.ops++
	o = op.Op{ID: opID(g.ops)}
	pick := func(label string, from []id.ID) id.ID { return rapid.SampledFrom(from).Draw(t, label) }
	switch o.Type = rapid.SampledFrom(propTypes).Draw(t, "type"); o.Type {
	case op.TypeNodeCreate:
		o.Node = pick("node", propNodes)
		o.Kind = rapid.SampledFrom([]op.Kind{op.KindCapability, op.KindCapability, op.KindDecision}).Draw(t, "kind")
		o.Handle = rapid.SampledFrom([]string{"a", "b", "c"}).Draw(t, "handle")
		if op.HasParent(o.Kind) {
			if rapid.Bool().Draw(t, "child?") {
				o.Parent = pick("parent", propNodes)
			}
			o.Key = rapid.SampledFrom([]string{"8", "i", "r"}).Draw(t, "key")
		}
		if rapid.Bool().Draw(t, "fields?") {
			o.Fields = fieldsOf("statement", fmt.Sprintf(`"s%d"`, g.ops))
		}
	case op.TypeNodeUpdate:
		o.Node = pick("node", propNodes)
		o.Base = g.base(t, o.Node)
		if rapid.Bool().Draw(t, "handle?") {
			o.Handle = rapid.SampledFrom([]string{"a", "b", "c"}).Draw(t, "handle")
		}
		name := rapid.SampledFrom([]string{"statement", "statement", "title"}).Draw(t, "set")
		o.Set = fieldsOf(name, fmt.Sprintf("{\n  \"n\": %d,\n  \"t\": \"<&>\"\n}", g.ops))
	case op.TypeNodeMove:
		o.Node = pick("node", propNodes)
		o.Base = g.base(t, o.Node)
		if rapid.Bool().Draw(t, "child?") {
			o.Parent = pick("parent", propNodes)
		}
		o.Key = rapid.SampledFrom([]string{"8", "i", "r"}).Draw(t, "key")
	case op.TypeNodeRetire:
		o.Node = pick("node", propNodes)
		o.Base = g.base(t, o.Node)
		o.Reason = "gone"
	case op.TypeLinkCreate:
		o.Link = pick("link", propLinks)
		o.LinkType = rapid.SampledFrom([]op.LinkType{op.LinkServes, op.LinkRelatesTo, op.LinkDependsOn, op.LinkDependsOn}).Draw(t, "link type")
		o.From = pick("from", ends)
		o.To = pick("to", ends)
		if o.LinkType == op.LinkRelatesTo || rapid.Bool().Draw(t, "note?") {
			o.Note = "why" // required on relates-to
		}
	case op.TypeLinkRetire:
		o.Link = pick("link", propLinks)
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
	return map[string]json.RawMessage{name: raw(value)}
}

// add encodes an entry of cs after the previous one and records its file.
func (g *logGen) add(t *rapid.T, cs id.ID, after *id.ID, b changeset.Body) {
	g.entries++
	eid := eID(g.entries)
	var a []id.ID
	if !after.IsZero() {
		a = []id.ID{*after}
	}
	p, data := encodeEntry(t, eid, cs, a, b)
	g.files = append(g.files, logFile{path: p, data: data, typ: b.Type()})
	*after = eid
}

// genFiles draws 1 to 5 change sets, each opened, given ops entries and
// moved through the lifecycle to a drawn state, plus now and then a corrupt
// entry file, a stray file or an entry out of order.
func genFiles(t *rapid.T) []logFile {
	n := rapid.IntRange(1, 5).Draw(t, "change sets")
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i + 1
	}
	nums = rapid.Permutation(nums).Draw(t, "change set IDs")
	ends := slices.Clone(propNodes)
	for _, k := range nums {
		ends = append(ends, csID(k))
	}
	g := &logGen{onTarget: map[id.ID][]id.ID{}}

	for i, k := range nums {
		cs := csID(k)
		var head id.ID
		kind := changeset.KindIntention
		if rapid.IntRange(0, 4).Draw(t, "record?") == 0 {
			kind = changeset.KindRecord
		}
		g.add(t, cs, &head, changeset.Open{Kind: kind, Intent: fmt.Sprintf("set %d", k)})
		if i == 0 && rapid.IntRange(0, 3).Draw(t, "seed?") > 0 {
			var seed []op.Op
			for _, node := range propNodes {
				g.ops++
				o := op.Op{ID: opID(g.ops), Type: op.TypeNodeCreate, Node: node, Kind: op.KindCapability,
					Handle: fmt.Sprintf("n%d", g.ops), Key: "i"}
				g.onTarget[node] = append(g.onTarget[node], o.ID)
				g.drawn = append(g.drawn, o.ID)
				seed = append(seed, o)
			}
			g.add(t, cs, &head, changeset.Ops{Ops: seed})
		}
		for range rapid.IntRange(0, 3).Draw(t, "ops entries") {
			var ops []op.Op
			for range rapid.IntRange(1, 3).Draw(t, "ops") {
				if o, ok := g.drawOp(t, ends); ok {
					ops = append(ops, o)
				}
			}
			if len(ops) > 0 {
				g.add(t, cs, &head, changeset.Ops{Ops: ops})
			}
		}

		state := rapid.SampledFrom(propStates).Draw(t, "state")
		reach := func(s changeset.State) {
			switch s {
			case changeset.StateProposed:
				g.add(t, cs, &head, changeset.Propose{})
			case changeset.StateAccepted:
				g.add(t, cs, &head, changeset.Propose{})
				g.add(t, cs, &head, changeset.Accept{Rule: "owner"})
			}
		}
		switch state {
		case changeset.StateProposed, changeset.StateAccepted:
			reach(state)
		case changeset.StateLanded:
			if kind == changeset.KindIntention {
				reach(changeset.StateAccepted)
			}
			g.add(t, cs, &head, changeset.Land{})
		case changeset.StateRejected:
			reach(changeset.StateProposed)
			g.add(t, cs, &head, changeset.Reject{Rule: "owner", Reason: "no"})
		case changeset.StateWithdrawn:
			g.add(t, cs, &head, changeset.Withdraw{})
		case changeset.StateSuperseded:
			// Sometimes by a change set not in the log.
			by := rapid.SampledFrom(append([]id.ID{csID(99)}, ends[len(propNodes):]...)).Draw(t, "by")
			if by == cs {
				by = csID(99)
			}
			g.add(t, cs, &head, changeset.Supersede{By: by})
		}
		if rapid.IntRange(0, 9).Draw(t, "late ops?") == 0 {
			if o, ok := g.drawOp(t, ends); ok {
				g.add(t, cs, &head, changeset.Ops{Ops: []op.Op{o}}) // skipped after a terminal state
			}
		}
		if rapid.IntRange(0, 9).Draw(t, "orphan?") == 0 {
			orphan := opID(8888) // names no entry: an order problem
			g.add(t, cs, &orphan, changeset.Note{Text: "lost"})
		}
		if rapid.IntRange(0, 9).Draw(t, "corrupt?") == 0 {
			g.entries++
			g.files = append(g.files, logFile{path: entryPath(cs, eID(g.entries)), data: []byte("{not json")})
		}
	}
	if rapid.IntRange(0, 9).Draw(t, "stray?") == 0 {
		g.files = append(g.files, logFile{path: "log/stray.txt", data: []byte("x")})
	}
	return g.files
}

// TestPropertyIncremental: catching up through any split of a log into
// steps, including steps that remove files already indexed (rebuilds)
// and add them back, gives the same index as one scratch build, and both
// equal what store, changeset and model compute from the files.
func TestPropertyIncremental(t *testing.T) {
	var lands, accepts, rebuilds, reopens, steps int
	rapid.Check(t, func(t *rapid.T) {
		files := genFiles(t)
		root, err := os.MkdirTemp("", "index-prop-")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(root)
		gdir := filepath.Join(root, store.Dir)
		if err := os.MkdirAll(filepath.Join(gdir, store.LogDir), 0o755); err != nil {
			t.Fatal(err)
		}
		db := filepath.Join(root, "index.db")

		x, err := Open(root, db)
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		defer func() { x.Close() }()

		remaining := make([]int, len(files))
		for i := range remaining {
			remaining[i] = i
		}
		var written []int
		removals := 0
		for len(remaining) > 0 {
			steps++
			if len(written) > 0 && removals < 2 && rapid.IntRange(0, 3).Draw(t, "remove?") == 0 {
				removals++
				out := rapid.SliceOfNDistinct(rapid.SampledFrom(written), 1, len(written), rapid.ID[int]).Draw(t, "remove")
				for _, i := range out {
					if err := os.Remove(filepath.Join(gdir, filepath.FromSlash(files[i].path))); err != nil {
						t.Fatal(err)
					}
					written = slices.DeleteFunc(written, func(w int) bool { return w == i })
					remaining = append(remaining, i)
				}
			} else {
				// Small steps half the time, so that orders that matter (a
				// link's endpoint change set arriving after the link lands)
				// come up often.
				most := len(remaining)
				if rapid.Bool().Draw(t, "small step?") {
					most = min(most, 2)
				}
				in := rapid.SliceOfNDistinct(rapid.SampledFrom(remaining), 1, most, rapid.ID[int]).Draw(t, "add")
				landed, accepted := false, false
				for _, i := range in {
					writeFile(t, gdir, files[i].path, files[i].data)
					remaining = slices.DeleteFunc(remaining, func(r int) bool { return r == i })
					written = append(written, i)
					landed = landed || files[i].typ == entry.TypeLand
					accepted = accepted || files[i].typ == entry.TypeAccept
				}
				if landed {
					lands++
				}
				if accepted {
					accepts++
				}
			}
			if rapid.Bool().Draw(t, "reopen?") {
				reopens++
				x.Close()
				if x, err = Open(root, db); err != nil {
					t.Fatalf("Open: %v", err)
				}
			} else if err := x.catchUp(); err != nil {
				t.Fatalf("catchUp: %v", err)
			}
			if x.stats.rebuilt {
				rebuilds++
			}
		}

		scratch, err := Open(root, filepath.Join(root, "scratch.db"))
		if err != nil {
			t.Fatalf("scratch Open: %v", err)
		}
		defer scratch.Close()
		same(t, "incremental vs scratch", dump(t, x), dump(t, scratch))
		checkModel(t, x, gdir)
	})
	t.Logf("steps %d, with a land %d, with an accept %d, rebuilds %d, reopens %d", steps, lands, accepts, rebuilds, reopens)
	if lands == 0 || accepts == 0 || rebuilds == 0 {
		t.Errorf("vacuous run: lands %d, accepts %d, rebuilds %d", lands, accepts, rebuilds)
	}
}
