package slice

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/model"
	"github.com/dphaener/genesis/internal/op"
)

// touched returns the change set's touched nodes that the view has, sorted
// by handle then ID. A retired one is kept only when the change retires it.
func (b *builder) touched() ([]model.Node, error) {
	var ids []id.ID
	for _, o := range b.req.Change.Ops {
		switch o.Type {
		case op.TypeLinkCreate:
			ids = append(ids, o.From, o.To)
		case op.TypeLinkRetire:
			l, ok, err := b.link(o.Link)
			if err != nil {
				return nil, err
			}
			if ok {
				ids = append(ids, l.From, l.To)
			}
		default:
			ids = append(ids, o.Node)
		}
	}
	slices.SortFunc(ids, id.Compare)
	var out []model.Node
	for _, i := range slices.Compact(ids) {
		n, ok, err := b.r.Node(i)
		if err != nil {
			return nil, err
		}
		// Not a node here: a change set endpoint, or a node this view
		// lacks.
		if !ok || n.Retired && !b.retires(i) {
			continue
		}
		out = append(out, n)
	}
	slices.SortFunc(out, func(a, b model.Node) int {
		return cmp.Or(strings.Compare(a.Handle, b.Handle), id.Compare(a.ID, b.ID))
	})
	return out, nil
}

// link looks link i up in the view, then in as-is, when they can.
func (b *builder) link(i id.ID) (model.Link, bool, error) {
	for _, r := range []Reader{b.r, b.asIs} {
		if lr, ok := r.(LinkReader); ok {
			if l, ok, err := lr.Link(i); err != nil || ok {
				return l, ok, err
			}
		}
	}
	return model.Link{}, false, nil
}

// retires reports whether the request's change set retires node i.
func (b *builder) retires(i id.ID) bool {
	if b.req.Change == nil {
		return false
	}
	return slices.ContainsFunc(b.req.Change.Ops, func(o op.Op) bool {
		return o.Type == op.TypeNodeRetire && o.Node == i
	})
}

func (b *builder) short(i id.ID) (string, error) {
	if s, ok := b.shorts[i]; ok {
		return s, nil
	}
	s, err := b.r.ShortID(i)
	if err != nil {
		return "", err
	}
	b.shorts[i] = s
	return s, nil
}

// changeShort returns change set i's short ID among req.ChangeSets, which
// are sorted. i need not be among them.
func (b *builder) changeShort(i id.ID) string {
	cs := b.req.ChangeSets
	k, found := slices.BinarySearchFunc(cs, i, id.Compare)
	var prev, next id.ID
	if k > 0 {
		prev = cs[k-1]
	}
	if found {
		k++
	}
	if k < len(cs) {
		next = cs[k]
	}
	return Short(i, prev, next)
}

// item returns n as an item, inherited from the handle from ("" for none).
func (b *builder) item(n model.Node, from string) (Item, error) {
	short, err := b.short(n.ID)
	if err != nil {
		return Item{}, err
	}
	return Item{ID: n.ID, Short: short, Handle: n.Handle, Kind: n.Kind, Statement: statement(n), From: from}, nil
}

// live returns node i when the view has it live, of kind k ("" for any).
func (b *builder) live(i id.ID, k op.Kind) (model.Node, bool, error) {
	n, ok, err := b.r.Node(i)
	if err != nil || !ok || n.Retired || k != "" && n.Kind != k {
		return model.Node{}, false, err
	}
	return n, true, nil
}

// liveLinks returns node i's live links of type t in direction dir.
func (b *builder) liveLinks(i id.ID, dir Direction, t op.LinkType) ([]model.Link, error) {
	ls, err := b.r.Links(i, dir)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(ls, func(l model.Link) bool { return l.Retired || l.Type != t }), nil
}

// hasLive reports whether node i has a live link of type t in direction dir
// whose other end is a live node of kind k.
func (b *builder) hasLive(i id.ID, dir Direction, t op.LinkType, k op.Kind) (bool, error) {
	ls, err := b.liveLinks(i, dir, t)
	if err != nil {
		return false, err
	}
	for _, l := range ls {
		other := l.From
		if dir == From {
			other = l.To
		}
		if _, ok, err := b.live(other, k); err != nil || ok {
			return ok, err
		}
	}
	return false, nil
}

// inherited adds what applies to f through its chain (f, then its
// ancestors from the parent up): intents, constraints not exempt for f,
// current decisions and blocking open questions.
func (b *builder) inherited(f model.Node, chain []model.Node, intents, cons, decs, blocking *section) error {
	type candidate struct {
		n    model.Node
		from string
	}
	var constraints []candidate
	exempt := map[id.ID]bool{}
	for _, c := range chain {
		from := ""
		if c.ID != f.ID {
			from = c.Handle
		}
		ls, err := b.r.Links(c.ID, Both)
		if err != nil {
			return err
		}
		for _, l := range ls {
			if l.Retired {
				continue
			}
			switch {
			case l.Type == op.LinkConstrains && l.To == c.ID:
				if n, ok, err := b.live(l.From, ""); err != nil {
					return err
				} else if ok {
					constraints = append(constraints, candidate{n, from})
				}
			case l.Type == op.LinkExempts && l.From == c.ID && !l.Decision.IsZero():
				if _, ok, err := b.live(l.Decision, op.KindDecision); err != nil {
					return err
				} else if ok {
					exempt[l.To] = true
				}
			case l.Type == op.LinkServes && l.From == c.ID:
				if err := b.addLive(intents, l.To, "", from, RelServes); err != nil {
					return err
				}
			case l.Type == op.LinkDecidedBy && l.From == c.ID:
				n, ok, err := b.live(l.To, "")
				if err != nil {
					return err
				}
				if !ok {
					continue
				}
				if superseded, err := b.hasLive(n.ID, To, op.LinkSupersedes, op.KindDecision); err != nil {
					return err
				} else if !superseded {
					if err := b.add(decs, n, from, ""); err != nil {
						return err
					}
				}
			case l.Type == op.LinkBlocks && l.To == c.ID:
				if err := b.addQuestion(blocking, l.From, from, RelBlocks, ""); err != nil {
					return err
				}
			}
		}
	}
	for _, c := range constraints {
		if !exempt[c.n.ID] {
			if err := b.add(cons, c.n, c.from, ""); err != nil {
				return err
			}
		}
	}
	return nil
}

// attached adds what is linked to f itself: related questions, anchors,
// evidence and neighbours.
func (b *builder) attached(f model.Node, related, anchors, evidence, neighbours *section) error {
	ls, err := b.r.Links(f.ID, Both)
	if err != nil {
		return err
	}
	for _, l := range ls {
		if l.Retired {
			continue
		}
		var err error
		switch {
		case l.Type == op.LinkImplements && l.To == f.ID:
			err = b.addLive(anchors, l.From, "", "", "")
		case l.Type == op.LinkVerifies && l.To == f.ID:
			n, ok, e := b.live(l.From, "")
			if e != nil || !ok {
				err = e
				break
			}
			it, e := b.item(n, "")
			if e != nil {
				err = e
				break
			}
			it.Status = evidenceOf(n)
			evidence.add(n.ID.String(), it)
		case l.Type == op.LinkDependsOn && l.From == f.ID && l.To != f.ID:
			err = b.addNeighbour(neighbours, l.To, RelNeeds, "")
		case l.Type == op.LinkDependsOn && l.To == f.ID && l.From != f.ID:
			err = b.addNeighbour(neighbours, l.From, RelNeededBy, "")
		case l.Type == op.LinkRelatesTo && l.From != l.To:
			other := l.From
			if other == f.ID {
				other = l.To
			}
			n, ok, e := b.live(other, "")
			if e != nil || !ok {
				err = e
				break
			}
			if n.Kind == op.KindQuestion {
				err = b.addQuestion(related, n.ID, "", RelRelatesTo, l.Note)
			} else {
				err = b.addNeighbour(neighbours, n.ID, RelRelatesTo, l.Note)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// add adds n to s, keyed by its ID.
func (b *builder) add(s *section, n model.Node, from, rel string) error {
	it, err := b.item(n, from)
	if err != nil {
		return err
	}
	it.Relation = rel
	s.add(n.ID.String(), it)
	return nil
}

// addLive adds node i to s when it is live.
func (b *builder) addLive(s *section, i id.ID, k op.Kind, from, rel string) error {
	n, ok, err := b.live(i, k)
	if err != nil || !ok {
		return err
	}
	return b.add(s, n, from, rel)
}

// addQuestion adds question i to s when it is live and open.
func (b *builder) addQuestion(s *section, i id.ID, from, rel, note string) error {
	n, ok, err := b.live(i, op.KindQuestion)
	if err != nil || !ok {
		return err
	}
	if answered, err := b.hasLive(i, To, op.LinkAnswers, op.KindDecision); err != nil || answered {
		return err
	}
	it, err := b.item(n, from)
	if err != nil {
		return err
	}
	it.Relation, it.Note = rel, oneLine(note)
	s.add(n.ID.String(), it)
	return nil
}

// addNeighbour adds live node i to s, keyed by its ID and relation.
func (b *builder) addNeighbour(s *section, i id.ID, rel, note string) error {
	n, ok, err := b.live(i, "")
	if err != nil || !ok {
		return err
	}
	it, err := b.item(n, "")
	if err != nil {
		return err
	}
	it.Relation, it.Note = rel, oneLine(note)
	s.add(n.ID.String()+" "+rel, it)
	return nil
}

// evidenceStatus returns node i's evidence status.
func (b *builder) evidenceStatus(i id.ID) (string, error) {
	ls, err := b.liveLinks(i, To, op.LinkVerifies)
	if err != nil {
		return "", err
	}
	status := StatusNone
	for _, l := range ls {
		n, ok, err := b.live(l.From, "")
		if err != nil {
			return "", err
		}
		if !ok {
			continue
		}
		switch e := evidenceOf(n); {
		case e == StatusFailing:
			return StatusFailing, nil
		case e == StatusVerified, status == StatusNone:
			status = e
		}
	}
	return status, nil
}

// focusItem returns f with every field, its evidence status and its
// proposed change.
func (b *builder) focusItem(f model.Node) (FocusItem, error) {
	it, err := b.item(f, "")
	if err != nil {
		return FocusItem{}, err
	}
	if it.Status, err = b.evidenceStatus(f.ID); err != nil {
		return FocusItem{}, err
	}
	proposed, err := b.proposed(f)
	if err != nil {
		return FocusItem{}, err
	}
	return FocusItem{Item: it, Fields: f.Fields, Retired: f.Retired, Reason: f.Reason, Proposed: proposed}, nil
}

// proposed returns f's proposed change against as-is, nil for none.
func (b *builder) proposed(f model.Node) ([]string, error) {
	if b.asIs == nil {
		return nil, nil
	}
	old, ok, err := b.asIs.Node(f.ID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return []string{"new"}, nil
	}
	var out []string
	if f.Retired && !old.Retired {
		if r := oneLine(f.Reason); r != "" {
			out = append(out, "retired: "+r)
		} else {
			out = append(out, "retired")
		}
	}
	if f.Parent != old.Parent {
		from, err := parentName(b.asIs, old.Parent)
		if err != nil {
			return nil, err
		}
		to, err := parentName(b.r, f.Parent)
		if err != nil {
			return nil, err
		}
		out = append(out, fmt.Sprintf("moved: %s → %s", from, to))
	}
	if f.Handle != old.Handle {
		out = append(out, fmt.Sprintf("handle: %s → %s", old.Handle, f.Handle))
	}
	names := slices.Sorted(maps.Keys(f.Fields))
	for k := range old.Fields {
		if _, ok := f.Fields[k]; !ok {
			names = append(names, k)
		}
	}
	slices.Sort(names)
	for _, k := range names {
		if !bytes.Equal(old.Fields[k], f.Fields[k]) {
			out = append(out, fmt.Sprintf("%s: %s → %s", oneLine(k), value(old.Fields[k]), value(f.Fields[k])))
		}
	}
	return out, nil
}

// parentName is the handle of parent p in r, "(root)" for none.
func parentName(r Reader, p id.ID) (string, error) {
	if p.IsZero() {
		return "(root)", nil
	}
	n, ok, err := r.Node(p)
	if err != nil {
		return "", err
	}
	if !ok {
		return p.String(), nil
	}
	return n.Handle, nil
}

// value is a field value as one line: a JSON string's text, other JSON as
// it is stored, "(none)" when missing.
func value(v json.RawMessage) string {
	if v == nil {
		return "(none)"
	}
	var s string
	if json.Unmarshal(v, &s) == nil {
		return oneLine(s)
	}
	return oneLine(string(v))
}

// inFlight returns the in-flight change sets, other than the request's,
// that target a focus node, sorted by ID.
func (b *builder) inFlight(isFocus map[id.ID]bool) ([]Change, error) {
	var out []Change
	for _, c := range b.req.InFlight {
		if b.req.Change != nil && c.ID == b.req.Change.ID || !slices.ContainsFunc(c.Targets, func(t id.ID) bool { return isFocus[t] }) {
			continue
		}
		out = append(out, Change{ID: c.ID, Short: b.changeShort(c.ID), State: c.State, Intent: oneLine(c.Intent)})
	}
	slices.SortFunc(out, func(a, b Change) int { return id.Compare(a.ID, b.ID) })
	return slices.CompactFunc(out, func(a, b Change) bool { return a.ID == b.ID }), nil
}

// sorted returns items sorted by handle, ID, then relation.
func sorted(items []Item) []Item {
	out := slices.Clone(items)
	slices.SortFunc(out, func(a, b Item) int {
		return cmp.Or(strings.Compare(a.Handle, b.Handle), id.Compare(a.ID, b.ID), strings.Compare(a.Relation, b.Relation))
	})
	return out
}

// statement is n's statement field as one line, or its handle when that
// is not a non-blank JSON string.
func statement(n model.Node) string {
	var s string
	if json.Unmarshal(n.Fields["statement"], &s) == nil && oneLine(s) != "" {
		return oneLine(s)
	}
	return n.Handle
}

// evidenceOf is evidence node n's status.
func evidenceOf(n model.Node) string {
	var s string
	if json.Unmarshal(n.Fields["status"], &s) == nil && (s == StatusVerified || s == StatusFailing) {
		return s
	}
	return StatusUnverified
}

func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
