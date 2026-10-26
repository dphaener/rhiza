package slice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/index"
	"github.com/dphaener/genesis/internal/model"
	"github.com/dphaener/genesis/internal/op"
	"github.com/dphaener/genesis/internal/store"
	"pgregory.net/rapid"
)

// propNode is node n of a drawn log. Nodes come in groups of four that
// share all but their last byte, so short IDs are 6 characters between
// groups and long within one.
func propNode(n int) id.ID {
	var i id.ID
	i[0] = 5
	i[1] = byte(n / 4)
	i[15] = byte(n)
	return i
}

func entryID(n int) id.ID { return mkID(3, n) }

var (
	// propKinds are the kinds of the nodes past the first four, which are
	// capabilities and intents.
	propKinds = []op.Kind{op.KindCapability, op.KindIntent, op.KindConstraint, op.KindConstraint, op.KindDecision,
		op.KindDecision, op.KindQuestion, op.KindEvidence, op.KindAnchor, op.KindSource}
	propTree = []op.Kind{op.KindCapability, op.KindCapability, op.KindIntent}
	// propLinkTypes are the link types a slice reads, and one it does not.
	propLinkTypes = []op.LinkType{op.LinkConstrains, op.LinkExempts, op.LinkServes, op.LinkDecidedBy,
		op.LinkSupersedes, op.LinkAnswers, op.LinkBlocks, op.LinkRelatesTo, op.LinkDependsOn,
		op.LinkImplements, op.LinkVerifies, op.LinkDerivedFrom}
	// propEnds are the kinds a link type's ends usually have, from then to;
	// nil is any.
	propEnds = map[op.LinkType][2][]op.Kind{
		op.LinkConstrains: {{op.KindConstraint}, propTree},
		op.LinkExempts:    {propTree, {op.KindConstraint}},
		op.LinkServes:     {propTree, {op.KindIntent}},
		op.LinkDecidedBy:  {propTree, {op.KindDecision}},
		op.LinkSupersedes: {{op.KindDecision}, {op.KindDecision}},
		op.LinkAnswers:    {{op.KindDecision}, {op.KindQuestion}},
		op.LinkBlocks:     {{op.KindQuestion}, propTree},
		op.LinkDependsOn:  {propTree, propTree},
		op.LinkImplements: {{op.KindAnchor}, propTree},
		op.LinkVerifies:   {{op.KindEvidence}, propTree},
	}
	propStatuses = []string{`"verified"`, `"failing"`, `"pending"`, `1`}
	propFinal    = []changeset.State{changeset.StateLanded, changeset.StateLanded, changeset.StateAccepted,
		changeset.StateAccepted, changeset.StateProposed, changeset.StateDraft, changeset.StateWithdrawn}
)

// propGen draws a log: a landed seed of nodes in a tree, then change sets
// of creates, updates, moves, links and retires in several states.
type propGen struct {
	t       *rapid.T
	kinds   map[id.ID]op.Kind
	nodes   []id.ID
	links   []id.ID
	ends    []id.ID // nodes, then change sets
	ver     map[id.ID]id.ID
	ops     int
	entries int
	files   map[string][]byte
}

func (g *propGen) add(cs id.ID, head *id.ID, b changeset.Body) {
	g.entries++
	eid := entryID(g.entries)
	body, err := changeset.EncodeBody(b)
	if err != nil {
		g.t.Fatalf("encode body: %v", err)
	}
	var after []id.ID
	if !head.IsZero() {
		after = []id.ID{*head}
	}
	actor := mkID(7, 1)
	data, err := entry.Encode(entry.Entry{V: changeset.EntryVersion(b), ID: eid, ChangeSet: cs, Actor: actor, For: actor,
		At: time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC), After: after, Type: b.Type(), Body: body})
	if err != nil {
		g.t.Fatalf("encode entry: %v", err)
	}
	g.files[store.LogDir+"/"+cs.String()+"/"+eid.String()+".json"] = data
	*head = eid
}

// op fills in o's ID and base, the last op drawn on its target.
func (g *propGen) op(o op.Op) op.Op {
	g.ops++
	o.ID = opID(g.ops)
	if o.Type != op.TypeNodeCreate && o.Type != op.TypeLinkCreate {
		o.Base = g.ver[o.Target()]
		if o.Base.IsZero() {
			o.Base = opID(9999)
		}
	}
	g.ver[o.Target()] = o.ID
	return o
}

// pick draws a node of one of kinds, or now and then (or when there is none)
// any end.
func (g *propGen) pick(label string, kinds []op.Kind) id.ID {
	var of []id.ID
	for _, n := range g.nodes {
		if slices.Contains(kinds, g.kinds[n]) {
			of = append(of, n)
		}
	}
	if len(of) == 0 || rapid.IntRange(0, 4).Draw(g.t, label+" any?") == 0 {
		return rapid.SampledFrom(g.ends).Draw(g.t, label)
	}
	return rapid.SampledFrom(of).Draw(g.t, label)
}

func (g *propGen) create(n int, kind op.Kind) op.Op {
	i := propNode(n)
	handle := fmt.Sprintf("n%d", n)
	if rapid.IntRange(0, 5).Draw(g.t, "shared handle?") == 0 {
		handle = "same"
	}
	o := op.Op{Type: op.TypeNodeCreate, Node: i, Kind: kind, Handle: handle}
	if op.HasParent(kind) {
		var parents []id.ID
		for _, p := range g.nodes {
			if op.HasParent(g.kinds[p]) {
				parents = append(parents, p)
			}
		}
		if len(parents) > 0 && rapid.IntRange(0, 3).Draw(g.t, "child?") > 0 {
			o.Parent = rapid.SampledFrom(parents).Draw(g.t, "parent")
		}
		o.Key = rapid.SampledFrom([]string{"a", "i", "r"}).Draw(g.t, "key")
	}
	if rapid.IntRange(0, 3).Draw(g.t, "statement?") > 0 {
		o.Fields = fields("statement", q(fmt.Sprintf("%s says\nthing %d", handle, n)))
	}
	if kind == op.KindEvidence {
		if o.Fields == nil {
			o.Fields = fields("status", "null")
		}
		o.Fields["status"] = json.RawMessage(rapid.SampledFrom(propStatuses).Draw(g.t, "status"))
	}
	g.kinds[i] = kind
	g.nodes = append(g.nodes, i)
	g.ends = append(g.ends, i)
	return g.op(o)
}

// drawOp draws one op of a change set after the seed.
func (g *propGen) drawOp() op.Op {
	switch k := rapid.IntRange(0, 9).Draw(g.t, "op"); {
	case k <= 4:
		lt := rapid.SampledFrom(propLinkTypes).Draw(g.t, "link type")
		ends := propEnds[lt]
		o := op.Op{Type: op.TypeLinkCreate, Link: ln(len(g.links) + 1), LinkType: lt,
			From: g.pick("from", ends[0]), To: g.pick("to", ends[1])}
		if lt == op.LinkRelatesTo || lt == op.LinkExempts || rapid.Bool().Draw(g.t, "note?") {
			o.Note = "because"
		}
		if lt == op.LinkExempts {
			// Usually a decision, sometimes another node or none at all.
			o.Decision = g.pick("decision", []op.Kind{op.KindDecision})
			if rapid.IntRange(0, 5).Draw(g.t, "lost decision?") == 0 {
				o.Decision = propNode(250)
			}
		}
		if o.From == o.To {
			o.To = propNode(251)
		}
		g.links = append(g.links, o.Link)
		return g.op(o)
	case k == 5 && len(g.links) > 0:
		return g.op(op.Op{Type: op.TypeLinkRetire, Link: rapid.SampledFrom(g.links).Draw(g.t, "retire link")})
	case k == 6:
		return g.op(op.Op{Type: op.TypeNodeRetire, Node: rapid.SampledFrom(g.nodes).Draw(g.t, "retire node"), Reason: "done with"})
	case k == 7:
		n := rapid.SampledFrom(g.nodes).Draw(g.t, "update")
		o := op.Op{Type: op.TypeNodeUpdate, Node: n, Set: fields("statement", q(fmt.Sprintf("changed %d", g.ops)))}
		if rapid.IntRange(0, 3).Draw(g.t, "rename?") == 0 {
			o.Handle = fmt.Sprintf("renamed%d", g.ops)
		}
		return g.op(o)
	case k == 8:
		var tree []id.ID
		for _, n := range g.nodes {
			if op.HasParent(g.kinds[n]) {
				tree = append(tree, n)
			}
		}
		n := rapid.SampledFrom(tree).Draw(g.t, "move")
		o := op.Op{Type: op.TypeNodeMove, Node: n, Key: "m"}
		if rapid.Bool().Draw(g.t, "move under?") {
			o.Parent = rapid.SampledFrom(tree).Draw(g.t, "new parent")
		}
		return g.op(o)
	}
	return g.create(len(g.nodes)+1, rapid.SampledFrom(propKinds).Draw(g.t, "kind"))
}

// genLog draws a log and returns its files by path under .genesis.
func genLog(t *rapid.T) map[string][]byte {
	g := &propGen{t: t, kinds: map[id.ID]op.Kind{}, ver: map[id.ID]id.ID{}, files: map[string][]byte{}}
	n := rapid.IntRange(1, 4).Draw(t, "change sets")
	for k := 1; k <= n+1; k++ {
		g.ends = append(g.ends, csID(k))
	}
	// The seed: a landed record of four tree nodes and a few others.
	var head id.ID
	seed := csID(1)
	g.add(seed, &head, changeset.Open{Kind: changeset.KindRecord, Intent: "seed"})
	var ops []op.Op
	for k := 1; k <= 4; k++ {
		ops = append(ops, g.create(k, rapid.SampledFrom(propTree).Draw(t, "tree kind")))
	}
	for k := 5; k <= 4+rapid.IntRange(2, 8).Draw(t, "others"); k++ {
		ops = append(ops, g.create(k, rapid.SampledFrom(propKinds).Draw(t, "kind")))
	}
	for range rapid.IntRange(0, 6).Draw(t, "seed links") {
		ops = append(ops, g.drawOp())
	}
	g.add(seed, &head, changeset.Ops{Ops: ops})
	g.add(seed, &head, changeset.Land{})

	for k := 2; k <= n+1; k++ {
		cs := csID(k)
		var head id.ID
		g.add(cs, &head, changeset.Open{Kind: changeset.KindIntention, Intent: fmt.Sprintf("set %d\nmore", k)})
		for range rapid.IntRange(0, 2).Draw(t, "ops entries") {
			var ops []op.Op
			for range rapid.IntRange(1, 4).Draw(t, "ops") {
				ops = append(ops, g.drawOp())
			}
			g.add(cs, &head, changeset.Ops{Ops: ops})
		}
		switch rapid.SampledFrom(propFinal).Draw(t, "state") {
		case changeset.StateProposed:
			g.add(cs, &head, changeset.Propose{})
		case changeset.StateAccepted:
			g.add(cs, &head, changeset.Propose{})
			g.add(cs, &head, changeset.Accept{Rule: "owner"})
		case changeset.StateLanded:
			g.add(cs, &head, changeset.Propose{})
			g.add(cs, &head, changeset.Accept{Rule: "owner"})
			g.add(cs, &head, changeset.Land{})
		case changeset.StateWithdrawn:
			g.add(cs, &head, changeset.Withdraw{})
		}
	}
	return g.files
}

// propOutcome is a slice as printed, or the error building or fitting it.
type propOutcome struct {
	res Result
	err string
}

func slicePrinted(r, asIs Reader, req Request, opts Options) propOutcome {
	s, err := Build(r, asIs, req)
	if err != nil {
		return propOutcome{err: "build: " + err.Error()}
	}
	res, err := Fit(s, opts)
	if err != nil {
		return propOutcome{err: "fit: " + err.Error()}
	}
	return propOutcome{res: res}
}

// TestPropertyIndexAgrees: on drawn logs, a slice built from the index's
// as-is and to-be views prints the same bytes as one built from
// model.AsIs/ToBe over FoldLog(ReadLog), for node and change slices, in
// both formats at any budget, with the same errors; and building it again
// prints the same bytes.
func TestPropertyIndexAgrees(t *testing.T) {
	var slicesN, errs, changes, cuts, retiredLinks int
	rapid.Check(t, func(t *rapid.T) {
		files := genLog(t)
		root, err := os.MkdirTemp("", "slice-prop-")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(root)
		gdir := filepath.Join(root, store.Dir)
		for p, data := range files {
			full := filepath.Join(gdir, filepath.FromSlash(p))
			if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(full, data, 0o644); err != nil {
				t.Fatal(err)
			}
		}
		x, err := index.Open(root, filepath.Join(root, "index.db"))
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		defer x.Close()

		sets, _ := changeset.FoldLog(store.ReadLog(os.DirFS(gdir)))
		var inFlight []InFlight
		var all []id.ID
		for _, s := range sets {
			all = append(all, s.ID)
			switch s.State {
			case changeset.StateDraft, changeset.StateProposed, changeset.StateAccepted:
				inFlight = append(inFlight, InFlightOf(s))
			}
		}
		slices.SortFunc(all, id.Compare)
		asIsModel, _ := model.AsIs(sets)
		modelAsIs, indexAsIs := FromModel(asIsModel), FromIndex(x, index.AsIs)

		for _, v := range []struct {
			name  string
			view  index.View
			model func([]changeset.ChangeSet) (model.Model, []model.Problem)
		}{{ViewAsIs, index.AsIs, model.AsIs}, {ViewToBe, index.ToBe, model.ToBe}} {
			m, _ := v.model(sets)
			fromModel, fromIndex := FromModel(m), FromIndex(x, v.view)
			var focuses []id.ID
			for _, n := range m.Nodes() {
				focuses = append(focuses, n.ID)
			}
			focuses = append(focuses, propNode(252)) // in no view
			var reqs []Request
			for range rapid.IntRange(1, 4).Draw(t, "node slices") {
				req := Request{Focus: rapid.SampledFrom(focuses).Draw(t, "focus"), View: v.name, InFlight: inFlight, ChangeSets: all}
				if rapid.IntRange(0, 2).Draw(t, "through a change set?") == 0 {
					req.Change = &sets[rapid.IntRange(0, len(sets)-1).Draw(t, "through")]
				}
				reqs = append(reqs, req)
			}
			for range rapid.IntRange(0, 2).Draw(t, "change slices") {
				cs := &sets[rapid.IntRange(0, len(sets)-1).Draw(t, "change")]
				reqs = append(reqs, Request{Change: cs, View: v.name, InFlight: inFlight, ChangeSets: all})
			}
			for _, req := range reqs {
				opts := Options{Budget: rapid.IntRange(1, 1500).Draw(t, "budget"), JSON: rapid.Bool().Draw(t, "json?"),
					DefaultView: DefaultView(req.Change)}
				asIsIndex, asIsM := indexAsIs, modelAsIs
				if v.view == index.AsIs {
					asIsIndex, asIsM = fromIndex, fromModel
				}
				got := slicePrinted(fromIndex, asIsIndex, req, opts)
				want := slicePrinted(fromModel, asIsM, req, opts)
				if got.err != want.err || !bytes.Equal(got.res.Bytes, want.res.Bytes) {
					t.Fatalf("%s %+v: index and model differ:\nindex %s\n%s\nmodel %s\n%s",
						v.name, opts, got.err, got.res.Bytes, want.err, want.res.Bytes)
				}
				same := slicePrinted(fromIndex, asIsIndex, req, opts)
				if same.err != got.err || !bytes.Equal(same.res.Bytes, got.res.Bytes) {
					t.Fatalf("%s %+v: a second run differs:\n%s\n%s", v.name, opts, got.res.Bytes, same.res.Bytes)
				}
				slicesN++
				if got.err != "" {
					errs++
				}
				if len(got.res.Omitted) > 0 {
					cuts++
				}
				if req.Focus.IsZero() {
					changes++
					if slices.ContainsFunc(req.Change.Ops, func(o op.Op) bool { return o.Type == op.TypeLinkRetire }) {
						retiredLinks++
					}
				}
			}
		}
	})
	t.Logf("slices %d (errors %d, cut %d), change slices %d (retiring links %d)", slicesN, errs, cuts, changes, retiredLinks)
	if slicesN == errs || changes == 0 || cuts == 0 || retiredLinks == 0 {
		t.Errorf("vacuous run: slices %d, errors %d, change slices %d, cut %d, retiring links %d",
			slicesN, errs, changes, cuts, retiredLinks)
	}
}
