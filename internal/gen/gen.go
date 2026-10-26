// Package gen writes a synthetic, deterministic Genesis repository: a
// repo.json and a log of change sets that projects with no problems. It
// exists for tests, benchmarks and demonstrations, never for real data.
//
// The shape, for Options o:
//
//   - o.Landed change sets go through the whole lifecycle: open (an
//     intention or a record), one to three ops entries, propose, accept
//     and land, each entry's after naming the one before it.
//   - The first three fifths of them create o.Nodes nodes of every kind
//     except anchor (milestone 2): intents and capabilities in a tree up to
//     six levels deep under a few roots, sibling keys from package key, and
//     constraints, decisions, questions, evidence, sources and actors
//     outside it. The first set creates one node of each of those kinds.
//   - The last four fifths create o.Links links across all fourteen
//     link types, between live nodes of plausible kinds.
//   - The last three quarters also update fields and handles, move leaf
//     nodes to new parents, and retire leaf nodes, other nodes and links,
//     each against the target's current version. A retired node is never
//     touched again.
//   - Then o.InFlight change sets: the first accepted, the second a draft,
//     the third proposed, and so on round, each with a few updates against
//     current as-is versions and a create or two. No two of them touch the
//     same target, so to-be has no problems either.
//
// Every handle is unique across the whole repository. IDs come from an
// id.Generator over a clock that starts at 2026-01-01T00:00:00Z and moves
// one second per entry, with entropy seeded from o.Seed, so the same
// Options always give the same bytes. Each file is the canonical encoding
// store.Write would produce, but written without fsync: this is test data.
package gen

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/dphaener/genesis/internal/canon"
	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/key"
	"github.com/dphaener/genesis/internal/op"
	"github.com/dphaener/genesis/internal/store"
)

// Options sizes a generated repository.
type Options struct {
	// Seed decides every ID and every random choice.
	Seed int64
	// Nodes and Links are the numbers of nodes and links in as-is, live and
	// retired. Creates in in-flight change sets come on top.
	Nodes, Links int
	// Landed is the number of landed change sets, InFlight the number of
	// draft, proposed and accepted ones after them. InFlight must be at
	// least 2, so that there is an accepted and a draft change set.
	Landed, InFlight int
}

// Default returns the model the budget test and benchmarks use: 35,000
// nodes and 15,000 links over 500 landed and 50 in-flight change sets.
func Default() Options {
	return Options{Seed: 1, Nodes: 35000, Links: 15000, Landed: 500, InFlight: 50}
}

// Sample is one node of the generated as-is model.
type Sample struct {
	ID     id.ID
	Handle string
}

// Summary describes what Write wrote. The samples are as the as-is view
// holds them; none of them is retired.
//
// To append a valid ops entry to the draft change set, write an entry with
// ChangeSet = Draft, After = [DraftHead], Actor and For = Actor, At after
// Last (Last plus a second, say), and an ops body whose operations are new
// creates or an update of Spare with Base = Spare's version (SpareVersion).
// No in-flight change set touches Spare. To land the accepted change set,
// write a land entry with ChangeSet = Accepted and After = [AcceptedHead].
// IDs for new entries may come from any generator: they only need to be
// unused.
type Summary struct {
	// Entries is the number of entry files, ChangeSets the number of change
	// sets. Nodes and Links are as-is's counts, live and retired.
	Entries, ChangeSets int
	Nodes, Links        int

	// Deep is a live node at Depth, the deepest level of the tree, where a
	// root is at depth 1. Depth is at least 4 for any reasonably sized
	// Options.
	Deep  Sample
	Depth int
	// Parent is a live tree node with LiveChildren live children (the most
	// of any node).
	Parent       Sample
	LiveChildren int
	// Linked is a live node with live links in both directions: LinksOut
	// from it and LinksIn to it.
	Linked            Sample
	LinksOut, LinksIn int
	Accepted          id.ID
	AcceptedHead      id.ID
	Draft             id.ID
	DraftHead         id.ID
	Actor             id.ID
	Spare             Sample
	SpareVersion      id.ID
	Last              time.Time
}

// start is the clock of the first entry.
var start = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// maxDepth is the deepest level a tree node is created or moved to.
const maxDepth = 6

// Write writes the repository o describes into root/.genesis, which must
// not exist yet, and returns its Summary.
func Write(root string, o Options) (Summary, error) {
	if o.Nodes < 1 || o.Links < 0 || o.Landed < 1 || o.InFlight < 2 {
		return Summary{}, fmt.Errorf("gen: options %+v: want Nodes >= 1, Links >= 0, Landed >= 1, InFlight >= 2", o)
	}
	g := newGenerator(root, o)
	if err := g.run(); err != nil {
		return Summary{}, err
	}
	return g.summary(), nil
}

// node is one node as the generator tracks it in as-is.
type node struct {
	id      id.ID
	kind    op.Kind
	handle  string
	parent  *node
	depth   int // 0 outside the tree
	version id.ID
	retired bool
	// children counts nodes directly under this one, live or retired.
	children int
	// liveChildren counts the live ones.
	liveChildren int
	lastKey      string
	out, in      int // live links from and to it
	hasNotes     bool
}

// link is one link as the generator tracks it in as-is.
type link struct {
	id       id.ID
	from, to *node
	version  id.ID
}

// pool is a set of items with uniform random choice and O(1) removal.
type pool[T comparable] struct {
	items []T
	at    map[T]int
}

func newPool[T comparable]() *pool[T] { return &pool[T]{at: map[T]int{}} }

func (p *pool[T]) add(x T) {
	if _, ok := p.at[x]; ok {
		return
	}
	p.at[x] = len(p.items)
	p.items = append(p.items, x)
}

func (p *pool[T]) remove(x T) {
	i, ok := p.at[x]
	if !ok {
		return
	}
	last := p.items[len(p.items)-1]
	p.items[i] = last
	p.at[last] = i
	p.items = p.items[:len(p.items)-1]
	delete(p.at, x)
}

func (p *pool[T]) len() int { return len(p.items) }

func (p *pool[T]) pick(r *rand.Rand) T { return p.items[r.Intn(len(p.items))] }

type generator struct {
	root string
	o    Options
	r    *rand.Rand
	ids  *id.Generator
	now  time.Time

	actor   id.ID
	handles int

	nodes  []*node // in creation order
	links  []*link
	roots  int
	live   *pool[*node]
	byKind map[op.Kind]*pool[*node]
	// parents holds live tree nodes above maxDepth.
	parents   *pool[*node]
	liveLinks *pool[*link]
	rootKey   string

	// touched holds the nodes and links the current change set has written;
	// frozen those an in-flight change set has.
	touched map[id.ID]bool
	frozen  map[id.ID]bool

	entries, sets int
	accepted      id.ID
	acceptedHead  id.ID
	draft         id.ID
	draftHead     id.ID
	madeDirs      map[id.ID]bool
}

func newGenerator(root string, o Options) *generator {
	clock := start
	g := &generator{
		root:      root,
		o:         o,
		r:         rand.New(rand.NewSource(o.Seed)),
		now:       clock,
		live:      newPool[*node](),
		byKind:    map[op.Kind]*pool[*node]{},
		parents:   newPool[*node](),
		liveLinks: newPool[*link](),
		frozen:    map[id.ID]bool{},
		madeDirs:  map[id.ID]bool{},
	}
	for _, k := range op.Kinds() {
		g.byKind[k] = newPool[*node]()
	}
	entropy := rand.New(rand.NewSource(o.Seed ^ 0x5eed))
	g.ids = id.NewGenerator(func() time.Time { return g.now }, entropy)
	return g
}

func (g *generator) newID() id.ID {
	x, err := g.ids.New()
	if err != nil {
		// Monotonic entropy overflows only after ~2^48 IDs in one
		// millisecond; the generator mints a few hundred.
		panic(err)
	}
	return x
}

// str spells s as a canonical JSON field value.
func str(s string) json.RawMessage {
	out, err := canon.Marshal(s)
	if err != nil {
		panic(err)
	}
	return json.RawMessage(out[:len(out)-1])
}

func (g *generator) handle(word string) string {
	g.handles++
	return fmt.Sprintf("%s-%d", word, g.handles)
}

func (g *generator) run() error {
	dir := filepath.Join(g.root, store.Dir)
	if err := os.Mkdir(dir, 0o755); err != nil {
		return err
	}
	if err := os.Mkdir(filepath.Join(dir, store.LogDir), 0o755); err != nil {
		return err
	}
	repoID := g.newID()
	g.actor = g.newID()
	data, err := canon.Marshal(map[string]any{"format": store.FormatVersion, "id": repoID, "remotes": []string{}})
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, store.RepoFile), data, 0o644); err != nil {
		return err
	}

	L := g.o.Landed
	createEnd := max(1, L*3/5)
	linkStart := L / 5
	mutStart := L / 4
	for i := range L {
		var nodes, links int
		if i < createEnd {
			nodes = share(g.o.Nodes, createEnd, i)
		}
		if i >= linkStart {
			links = share(g.o.Links, L-linkStart, i-linkStart)
		}
		if err := g.landed(i, nodes, links, i >= mutStart); err != nil {
			return err
		}
	}
	for j := range g.o.InFlight {
		if err := g.inFlight(j); err != nil {
			return err
		}
	}
	return nil
}

// share is part i of total split into n near-equal parts.
func share(total, n, i int) int {
	s := total / n
	if i < total%n {
		s++
	}
	return s
}

// changeSet writes entries for one change set in a linear chain.
type changeSet struct {
	g    *generator
	id   id.ID
	head id.ID
}

func (g *generator) open(kind changeset.Kind, intent string) (*changeSet, error) {
	g.tick()
	c := &changeSet{g: g, id: g.newID()}
	g.sets++
	g.touched = map[id.ID]bool{}
	return c, c.write(changeset.Open{Kind: kind, Intent: intent}, false)
}

// tick moves the clock to the next entry's second.
func (g *generator) tick() {
	g.now = start.Add(time.Duration(g.entries) * time.Second)
}

func (c *changeSet) write(b changeset.Body, tick bool) error {
	g := c.g
	if tick {
		g.tick()
	}
	body, err := changeset.EncodeBody(b)
	if err != nil {
		return fmt.Errorf("gen: %s body: %w", b.Type(), err)
	}
	e := entry.Entry{V: changeset.EntryVersion(b), ID: g.newID(), ChangeSet: c.id, Actor: g.actor, For: g.actor,
		At: g.now.UTC().Truncate(time.Second), Type: b.Type(), Body: body}
	if !c.head.IsZero() {
		e.After = []id.ID{c.head}
	}
	data, err := entry.Encode(e)
	if err != nil {
		return fmt.Errorf("gen: entry: %w", err)
	}
	dir := filepath.Join(g.root, store.Dir, store.LogDir, c.id.String())
	if !g.madeDirs[c.id] {
		if err := os.Mkdir(dir, 0o755); err != nil {
			return err
		}
		g.madeDirs[c.id] = true
	}
	if err := os.WriteFile(filepath.Join(dir, e.ID.String()+".json"), data, 0o644); err != nil {
		return err
	}
	c.head = e.ID
	g.entries++
	return nil
}

// ops writes ops in one to three ops entries.
func (c *changeSet) ops(ops []op.Op) error {
	if len(ops) == 0 {
		return nil
	}
	n := 1 + c.g.r.Intn(3)
	n = min(n, len(ops))
	for i := range n {
		part := ops[len(ops)*i/n : len(ops)*(i+1)/n]
		if err := c.write(changeset.Ops{Ops: part}, true); err != nil {
			return err
		}
	}
	return nil
}

// kindWeights are the odds, in 100, of each kind a create picks.
var kindWeights = []struct {
	kind   op.Kind
	weight int
}{
	{op.KindIntent, 25}, {op.KindCapability, 45},
	{op.KindConstraint, 5}, {op.KindDecision, 5}, {op.KindQuestion, 5},
	{op.KindEvidence, 5}, {op.KindSource, 5}, {op.KindActor, 5},
}

func (g *generator) pickKind() op.Kind {
	n := g.r.Intn(100)
	for _, w := range kindWeights {
		if n < w.weight {
			return w.kind
		}
		n -= w.weight
	}
	return op.KindCapability
}

func (g *generator) landed(i, nodes, links int, mutate bool) error {
	kind := changeset.KindIntention
	if i%4 == 3 {
		kind = changeset.KindRecord
	}
	c, err := g.open(kind, fmt.Sprintf("Landed change %d", i+1))
	if err != nil {
		return err
	}
	var ops []op.Op
	for k := range nodes {
		kd := g.pickKind()
		if i == 0 && k < len(kindWeights) {
			kd = kindWeights[k].kind
		}
		ops = append(ops, g.create(kd, false))
	}
	for range links {
		ops = append(ops, g.createLink(false))
	}
	if mutate {
		per := max(1, g.o.Nodes/g.o.Landed/4)
		for range per {
			if o, ok := g.update(false); ok {
				ops = append(ops, o)
			}
		}
		for range 3 {
			if o, ok := g.move(); ok {
				ops = append(ops, o)
			}
		}
		for range 2 {
			if o, ok := g.retireNode(); ok {
				ops = append(ops, o)
			}
		}
		if g.liveLinks.len() > 0 {
			ops = append(ops, g.retireLink())
		}
	}
	if len(ops) == 0 {
		// A change set needs something to land; a field on any live node.
		if o, ok := g.update(false); ok {
			ops = append(ops, o)
		}
	}
	if err := c.ops(ops); err != nil {
		return err
	}
	if err := c.write(changeset.Propose{}, true); err != nil {
		return err
	}
	if err := c.write(changeset.Accept{Rule: "owner", Reason: "generated"}, true); err != nil {
		return err
	}
	return c.write(changeset.Land{}, true)
}

// inFlight writes in-flight change set j: accepted, draft, proposed, round.
func (g *generator) inFlight(j int) error {
	c, err := g.open(changeset.KindIntention, fmt.Sprintf("In-flight change %d", j+1))
	if err != nil {
		return err
	}
	var ops []op.Op
	for range 2 + g.r.Intn(3) {
		if o, ok := g.update(true); ok {
			ops = append(ops, o)
		}
	}
	for range 1 + g.r.Intn(2) {
		kd := op.KindCapability
		if g.r.Intn(2) == 0 {
			kd = op.KindDecision
		}
		ops = append(ops, g.create(kd, true))
	}
	if err := c.ops(ops); err != nil {
		return err
	}
	state := []changeset.State{changeset.StateAccepted, changeset.StateDraft, changeset.StateProposed}[j%3]
	if state != changeset.StateDraft {
		if err := c.write(changeset.Propose{}, true); err != nil {
			return err
		}
	}
	if state == changeset.StateAccepted {
		if err := c.write(changeset.Accept{Rule: "owner", Reason: "generated"}, true); err != nil {
			return err
		}
	}
	switch {
	case state == changeset.StateAccepted && g.accepted.IsZero():
		g.accepted, g.acceptedHead = c.id, c.head
	case state == changeset.StateDraft && g.draft.IsZero():
		g.draft, g.draftHead = c.id, c.head
	}
	return nil
}

// eligible reports whether n may be written by the current change set.
func (g *generator) eligible(n *node) bool {
	if n.retired || g.touched[n.id] || g.frozen[n.id] {
		return false
	}
	return true
}

// pickParent returns a live tree node above maxDepth with room for a key,
// or nil for a root.
func (g *generator) pickParent(not *node) *node {
	for range 20 {
		if g.parents.len() == 0 {
			return nil
		}
		p := g.parents.pick(g.r)
		if p != not && (not == nil || p != not.parent) && len(p.lastKey) < 64 {
			return p
		}
	}
	return nil
}

func (g *generator) nextKey(p *node) string {
	last := g.rootKey
	if p != nil {
		last = p.lastKey
	}
	k, err := key.Between(last, "")
	if err != nil {
		panic(err)
	}
	if p != nil {
		p.lastKey = k
	} else {
		g.rootKey = k
	}
	return k
}

// create returns a node.create of kind k. In-flight creates are not added
// to the as-is state.
func (g *generator) create(k op.Kind, inFlight bool) op.Op {
	n := &node{id: g.newID(), kind: k, handle: g.handle(string(k))}
	o := op.Op{ID: g.newID(), Type: op.TypeNodeCreate, Node: n.id, Kind: k, Handle: n.handle,
		Fields: map[string]json.RawMessage{
			"title":  str(fmt.Sprintf("Generated %s %d", k, g.handles)),
			"status": str("open"),
		}}
	if op.HasParent(k) {
		var p *node
		// A few roots, then everything under an existing node.
		rootTarget := max(1, g.o.Nodes/2500)
		if g.roots >= rootTarget || inFlight {
			p = g.pickParent(nil)
		}
		if p == nil && inFlight {
			// No parent available: an in-flight create stays out of the tree.
			k = op.KindDecision
			n.kind, o.Kind = k, k
			return o
		}
		o.Key = g.nextKey(p)
		if p != nil {
			o.Parent = p.id
		}
		if inFlight {
			return o
		}
		n.parent = p
		n.depth = 1
		if p == nil {
			g.roots++
		} else {
			n.depth = p.depth + 1
			p.children++
			p.liveChildren++
		}
	}
	if inFlight {
		return o
	}
	n.version = o.ID
	g.nodes = append(g.nodes, n)
	g.live.add(n)
	g.byKind[k].add(n)
	if op.HasParent(k) && n.depth < maxDepth {
		g.parents.add(n)
	}
	g.touched[n.id] = true
	return o
}

// linkKinds suggests endpoint kinds for each link type.
var linkKinds = map[op.LinkType][2]op.Kind{
	op.LinkServes:         {op.KindCapability, op.KindIntent},
	op.LinkConstrains:     {op.KindConstraint, op.KindCapability},
	op.LinkExempts:        {op.KindCapability, op.KindConstraint},
	op.LinkDependsOn:      {op.KindCapability, op.KindCapability},
	op.LinkRelatesTo:      {op.KindIntent, op.KindIntent},
	op.LinkDecidedBy:      {op.KindCapability, op.KindDecision},
	op.LinkAnswers:        {op.KindDecision, op.KindQuestion},
	op.LinkBlocks:         {op.KindQuestion, op.KindCapability},
	op.LinkSupersedes:     {op.KindDecision, op.KindDecision},
	op.LinkDerivedFrom:    {op.KindIntent, op.KindSource},
	op.LinkDiscoveredFrom: {op.KindQuestion, op.KindEvidence},
	op.LinkImplements:     {op.KindCapability, op.KindIntent},
	op.LinkVerifies:       {op.KindEvidence, op.KindCapability},
	op.LinkObserves:       {op.KindActor, op.KindCapability},
}

// endpoint picks a live node, of kind k when there is one.
func (g *generator) endpoint(k op.Kind, not *node) *node {
	for range 10 {
		p := g.byKind[k]
		if p.len() == 0 {
			p = g.live
		}
		n := p.pick(g.r)
		if n != not {
			return n
		}
	}
	for _, n := range g.live.items {
		if n != not {
			return n
		}
	}
	return nil
}

func (g *generator) createLink(inFlight bool) op.Op {
	t := op.LinkTypes()[g.r.Intn(len(op.LinkTypes()))]
	if g.live.len() < 2 {
		panic("gen: links need at least two nodes")
	}
	kinds := linkKinds[t]
	from := g.endpoint(kinds[0], nil)
	to := g.endpoint(kinds[1], from)
	l := &link{id: g.newID(), from: from, to: to}
	o := op.Op{ID: g.newID(), Type: op.TypeLinkCreate, Link: l.id, LinkType: t, From: from.id, To: to.id}
	switch t {
	case op.LinkRelatesTo:
		o.Note = "related by the generator"
	case op.LinkExempts:
		o.Note = "exempt by the generator"
		o.Decision = g.endpoint(op.KindDecision, nil).id
	}
	if inFlight {
		return o
	}
	l.version = o.ID
	g.links = append(g.links, l)
	g.liveLinks.add(l)
	from.out++
	to.in++
	g.touched[l.id] = true
	return o
}

// update returns a node.update of a random live node the current change
// set has not written.
func (g *generator) update(inFlight bool) (op.Op, bool) {
	for range 20 {
		n := g.live.pick(g.r)
		if !g.eligible(n) {
			continue
		}
		o := op.Op{ID: g.newID(), Type: op.TypeNodeUpdate, Node: n.id, Base: n.version}
		switch g.r.Intn(4) {
		case 0:
			o.Set = map[string]json.RawMessage{"status": str([]string{"open", "doing", "done"}[g.r.Intn(3)])}
		case 1:
			o.Set = map[string]json.RawMessage{"notes": str(fmt.Sprintf("Note %d", g.r.Intn(1000000)))}
		case 2:
			if n.hasNotes && !inFlight {
				o.Unset = []string{"notes"}
			} else {
				o.Set = map[string]json.RawMessage{"title": str(fmt.Sprintf("Renamed %s", n.handle))}
			}
		case 3:
			o.Handle = g.handle(string(n.kind))
		}
		g.touched[n.id] = true
		if inFlight {
			g.frozen[n.id] = true
			return o, true
		}
		if o.Set["notes"] != nil {
			n.hasNotes = true
		}
		if len(o.Unset) > 0 {
			n.hasNotes = false
		}
		if o.Handle != "" {
			n.handle = o.Handle
		}
		n.version = o.ID
		return o, true
	}
	return op.Op{}, false
}

// move returns a node.move of a live tree leaf to a new parent.
func (g *generator) move() (op.Op, bool) {
	for range 20 {
		n := g.live.pick(g.r)
		if !op.HasParent(n.kind) || n.children > 0 || !g.eligible(n) {
			continue
		}
		p := g.pickParent(n)
		if p == nil {
			return op.Op{}, false
		}
		o := op.Op{ID: g.newID(), Type: op.TypeNodeMove, Node: n.id, Base: n.version, Parent: p.id, Key: g.nextKey(p)}
		if old := n.parent; old != nil {
			old.children--
			old.liveChildren--
		} else {
			g.roots--
		}
		n.parent = p
		p.children++
		p.liveChildren++
		n.depth = p.depth + 1
		if n.depth < maxDepth {
			g.parents.add(n)
		} else {
			g.parents.remove(n)
		}
		n.version = o.ID
		g.touched[n.id] = true
		return o, true
	}
	return op.Op{}, false
}

// retireNode returns a node.retire of a live node with no children.
func (g *generator) retireNode() (op.Op, bool) {
	for range 20 {
		n := g.live.pick(g.r)
		if n.children > 0 || !g.eligible(n) {
			continue
		}
		o := op.Op{ID: g.newID(), Type: op.TypeNodeRetire, Node: n.id, Base: n.version, Reason: "retired by the generator"}
		n.retired = true
		n.version = o.ID
		if n.parent != nil {
			n.parent.liveChildren--
		}
		g.live.remove(n)
		g.byKind[n.kind].remove(n)
		g.parents.remove(n)
		g.touched[n.id] = true
		return o, true
	}
	return op.Op{}, false
}

// retireLink returns a link.retire of a random live link.
func (g *generator) retireLink() op.Op {
	l := g.liveLinks.pick(g.r)
	o := op.Op{ID: g.newID(), Type: op.TypeLinkRetire, Link: l.id, Base: l.version}
	l.version = o.ID
	g.liveLinks.remove(l)
	l.from.out--
	l.to.in--
	return o
}

func (g *generator) summary() Summary {
	s := Summary{Entries: g.entries, ChangeSets: g.sets, Nodes: len(g.nodes), Links: len(g.links),
		Accepted: g.accepted, AcceptedHead: g.acceptedHead, Draft: g.draft, DraftHead: g.draftHead,
		Actor: g.actor, Last: g.now}
	var deep, parent, linked, spare *node
	for _, n := range g.nodes {
		if n.retired {
			continue
		}
		if deep == nil || n.depth > deep.depth {
			deep = n
		}
		if op.HasParent(n.kind) && (parent == nil || n.liveChildren > parent.liveChildren) {
			parent = n
		}
		if min(n.in, n.out) > 0 && (linked == nil || min(n.in, n.out) > min(linked.in, linked.out)) {
			linked = n
		}
		if spare == nil && !g.frozen[n.id] {
			spare = n
		}
	}
	if deep != nil {
		s.Deep, s.Depth = Sample{deep.id, deep.handle}, deep.depth
	}
	if parent != nil {
		s.Parent, s.LiveChildren = Sample{parent.id, parent.handle}, parent.liveChildren
	}
	if linked != nil {
		s.Linked, s.LinksOut, s.LinksIn = Sample{linked.id, linked.handle}, linked.out, linked.in
	}
	if spare != nil {
		s.Spare, s.SpareVersion = Sample{spare.id, spare.handle}, spare.version
	}
	return s
}
