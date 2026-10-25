// Package slice builds a slice: everything an agent needs to work on one part
// of the model, and nothing else, read from one view of it. Build is a pure
// function of its Reader and Request: no clock, no map order reaches the
// result, and the same view always gives the same Slice.
//
// The vocabulary:
//
//   - A slice is ten sections about its focus, in this order: Focus,
//     Context, Constraints, Decisions, Open questions, Anchors, Evidence,
//     Children, Neighbours, In flight. Sources, history, superseded
//     decisions and retired nodes are left out; they are lineage, for `why`.
//   - A node slice has one focus node. A change slice has a change set's
//     touched nodes as its focus: the targets of its node operations and the
//     node endpoints of the links it creates or retires. Each touched node
//     is a Focus item; every other section is the union over them,
//     de-duplicated by node ID, the first one found kept.
//   - Something applies to a node directly when it is linked to the node,
//     and inherited from an ancestor when it is linked to one, marked "from
//     <ancestor handle>". The nearest one found is kept: the focus, then
//     its ancestors from its parent up. Constraints apply through
//     `constrains`, decisions through `decided-by` and intents through
//     `serves`; a constraint linked to nothing applies nowhere.
//   - A constraint is exempt for a node when the node or an ancestor has a
//     live `exempts` link to it naming a live decision. An `exempts` naming
//     no decision, or a retired one, exempts nothing.
//   - A node's statement is its `statement` field when that is a JSON
//     string, as one line; otherwise its handle.
//   - An evidence node's status is its `status` field when that is the JSON
//     string "verified" or "failing", and "unverified" otherwise. A node's
//     evidence status is over the evidence with a live `verifies` link to
//     it: "failing" if any fails, else "verified" if any is verified, else
//     "unverified", and "none" with no evidence.
//   - A short ID is the shortest prefix of an ID, at least 6 characters,
//     unique in its namespace: one character more than it shares with the
//     nearest ID on either side. A node's namespace is the view's node IDs,
//     live and retired; a change set's is every change set ID in the log
//     (Request.ChangeSets), since change sets resolve among themselves.
//   - A one-line text has every run of white space, newlines included,
//     written as one space.
//
// What each section holds. Only live nodes and live links count; a retired
// node appears only as the Focus item of a change set that retires it.
//
//   - Focus: each focus node with every field, its evidence status and its
//     proposed change: against as-is, "new" for a node as-is lacks,
//     "retired: <reason>" for one this view retires, "moved: <old parent> →
//     <new parent>" for a new parent ("(root)" for none), and one line per
//     changed field (and the handle), "<name>: <as-is> → <proposed>", with
//     "(none)" for a missing value, every value one line. In as-is it is
//     empty.
//   - Context: each ancestor, the root first, then the intents served.
//   - Constraints, Decisions: those that apply and are not exempt. A
//     decision is current unless a live decision supersedes it.
//   - Open questions: live questions with no live `answers` from a live
//     decision. First those that block the focus or an ancestor, then those
//     linked to the focus by `relates-to`, with its note.
//   - Anchors, Evidence: through live `implements` and `verifies` links to
//     the focus; evidence with its status.
//   - Children: the live children of the focus, other than the focus
//     nodes, each with its statement and evidence status.
//   - Neighbours: `depends-on` out ("needs") and in ("needed by"), and
//     `relates-to` either way with its note, to nodes that are not
//     questions (those are under Open questions), other than the focus
//     nodes.
//   - In flight: every other in-flight change set whose operations target a
//     focus node.
//
// Every section but Context is sorted by handle, then ID, Open questions
// with the blocking ones first. Context keeps the ancestors in the order
// found, each focus node's root first, then the intents sorted the same way.
package slice

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/model"
	"github.com/dphaener/genesis/internal/op"
)

var (
	// ErrNotFound is wrapped when the focus names no node in the view.
	ErrNotFound = errors.New("slice: not found")
	// ErrRetired is wrapped when a node slice's focus is retired and the
	// change set, if any, does not retire it.
	ErrRetired = errors.New("slice: retired")
	// ErrEmpty is wrapped for a request with neither a focus nor a change
	// set.
	ErrEmpty = errors.New("slice: nothing to slice")
)

// The evidence statuses.
const (
	StatusVerified   = "verified"
	StatusFailing    = "failing"
	StatusUnverified = "unverified"
	StatusNone       = "none"
)

// The relations an item can have to the focus.
const (
	RelAncestor  = "ancestor"
	RelServes    = "serves"
	RelBlocks    = "blocks"
	RelRelatesTo = "relates-to"
	RelNeeds     = "needs"
	RelNeededBy  = "needed by"
)

// InFlight is an in-flight change set: its ID, state, intent and the nodes
// its operations target.
type InFlight struct {
	ID      id.ID
	State   changeset.State
	Intent  string
	Targets []id.ID
}

// InFlightOf returns cs as an InFlight, its targets the nodes of its node
// operations, sorted and unique.
func InFlightOf(cs changeset.ChangeSet) InFlight {
	var ts []id.ID
	for _, o := range cs.Ops {
		if o.Type != op.TypeLinkCreate && o.Type != op.TypeLinkRetire {
			ts = append(ts, o.Node)
		}
	}
	slices.SortFunc(ts, id.Compare)
	return InFlight{ID: cs.ID, State: cs.State, Intent: cs.Open.Intent, Targets: slices.Compact(ts)}
}

// Request says what to slice. A zero Focus with a Change is a change slice;
// a Focus is a node slice, viewed through Change when it is set.
type Request struct {
	Focus  id.ID
	Change *changeset.ChangeSet
	// View is the view's name ("as-is", "to-be", "what-if"), carried
	// through to the Slice.
	View string
	// InFlight is the in-flight change sets; Change, when among them, is
	// left out of the slice.
	InFlight []InFlight
	// ChangeSets is every change set ID in the log, sorted, for change set
	// short IDs. Empty, they are 6 characters.
	ChangeSets []id.ID
}

// Change is a change set as a slice shows it: the header of a change slice,
// or an In flight item.
type Change struct {
	ID     id.ID
	Short  string
	State  changeset.State
	Intent string
}

// Item is one line of a section. From is the ancestor's handle for an
// inherited item; Status an evidence status; Relation how the item relates
// to the focus; Note the link's note. Each is "" where it does not apply.
type Item struct {
	ID        id.ID
	Short     string
	Handle    string
	Kind      op.Kind
	Statement string
	Note      string
	From      string
	Status    string
	Relation  string
}

// FocusItem is a focus node: its item (Status is its evidence status), every
// field, whether it is retired and why, and its proposed change against
// as-is, one line each.
type FocusItem struct {
	Item
	Fields   map[string]json.RawMessage
	Retired  bool
	Reason   string
	Proposed []string
}

// Slice is the ten sections. Title is the focus handle for a node slice and
// `change <short> "<intent>"` for a change slice; Change is the change set
// sliced or viewed through, nil for none. ChangeSlice is true for a change
// slice and false for a node slice, viewed through Change or not. Empty
// sections are nil.
type Slice struct {
	Title       string
	View        string
	Change      *Change
	ChangeSlice bool
	Focus       []FocusItem
	Context     []Item
	Constraints []Item
	Decisions   []Item
	Questions   []Item
	Anchors     []Item
	Evidence    []Item
	Children    []Item
	Neighbours  []Item
	InFlight    []Change
}

// Build builds the slice req asks for from r. asIs is the as-is view, to
// show each focus node's proposed change against; in as-is it is r itself.
// A Focus naming no node in r is an error wrapping ErrNotFound.
func Build(r, asIs Reader, req Request) (Slice, error) {
	b := &builder{r: r, asIs: asIs, req: req, shorts: map[id.ID]string{}}
	return b.build()
}

type builder struct {
	r, asIs Reader
	req     Request
	shorts  map[id.ID]string
}

// section collects items de-duplicated by key, the first kept.
type section struct {
	seen  map[string]bool
	items []Item
}

func (s *section) add(key string, it Item) {
	if s.seen == nil {
		s.seen = map[string]bool{}
	}
	if !s.seen[key] {
		s.seen[key] = true
		s.items = append(s.items, it)
	}
}

func (b *builder) build() (Slice, error) {
	req := b.req
	s := Slice{View: req.View}
	if req.Change != nil {
		s.Change = &Change{ID: req.Change.ID, Short: b.changeShort(req.Change.ID), State: req.Change.State, Intent: oneLine(req.Change.Open.Intent)}
	}
	var focus []model.Node
	switch {
	case !req.Focus.IsZero():
		n, ok, err := b.r.Node(req.Focus)
		if err != nil {
			return Slice{}, err
		}
		if !ok {
			return Slice{}, fmt.Errorf("%w: no node %s", ErrNotFound, req.Focus)
		}
		if n.Retired && !b.retires(n.ID) {
			return Slice{}, fmt.Errorf("%w: node %s (%s) is retired", ErrRetired, n.Handle, n.ID)
		}
		focus = []model.Node{n}
		s.Title = n.Handle
	case req.Change != nil:
		var err error
		if focus, err = b.touched(); err != nil {
			return Slice{}, err
		}
		s.Title = fmt.Sprintf("change %s %q", s.Change.Short, s.Change.Intent)
		s.ChangeSlice = true
	default:
		return Slice{}, ErrEmpty
	}

	var context, intents, cons, decs, blocking, related, anchors, evidence, children, neighbours section
	isFocus := map[id.ID]bool{}
	for _, f := range focus {
		isFocus[f.ID] = true
	}
	for _, f := range focus {
		fi, err := b.focusItem(f)
		if err != nil {
			return Slice{}, err
		}
		s.Focus = append(s.Focus, fi)

		anc, err := b.r.Ancestors(f.ID)
		if err != nil {
			return Slice{}, err
		}
		for _, a := range anc {
			it, err := b.item(a, "")
			if err != nil {
				return Slice{}, err
			}
			it.Relation = RelAncestor
			context.add(a.ID.String(), it)
		}
		// chain is the focus, then its ancestors from the parent up.
		chain := []model.Node{f}
		for _, a := range slices.Backward(anc) {
			chain = append(chain, a)
		}
		if err := b.inherited(f, chain, &intents, &cons, &decs, &blocking); err != nil {
			return Slice{}, err
		}
		if err := b.attached(f, &related, &anchors, &evidence, &neighbours); err != nil {
			return Slice{}, err
		}
		kids, err := b.r.Children(f.ID)
		if err != nil {
			return Slice{}, err
		}
		for _, k := range kids {
			if k.Retired {
				continue
			}
			it, err := b.item(k, "")
			if err != nil {
				return Slice{}, err
			}
			if it.Status, err = b.evidenceStatus(k.ID); err != nil {
				return Slice{}, err
			}
			children.add(k.ID.String(), it)
		}
	}

	s.Context = slices.Concat(context.items, sorted(intents.items))
	s.Constraints = sorted(cons.items)
	s.Decisions = sorted(decs.items)
	// A question that both blocks and relates is listed once, as blocking.
	s.Questions = sorted(blocking.items)
	for _, q := range sorted(related.items) {
		if !blocking.seen[q.ID.String()] {
			s.Questions = append(s.Questions, q)
		}
	}
	s.Anchors = sorted(anchors.items)
	s.Evidence = sorted(evidence.items)
	// A focus node is not listed again as another focus node's child or
	// neighbour; in a change slice it is a Focus item already.
	inFocus := func(it Item) bool { return isFocus[it.ID] }
	s.Children = sorted(slices.DeleteFunc(children.items, inFocus))
	s.Neighbours = sorted(slices.DeleteFunc(neighbours.items, inFocus))
	inFlight, err := b.inFlight(isFocus)
	if err != nil {
		return Slice{}, err
	}
	s.InFlight = inFlight
	return s, nil
}
