// Package changeset defines change sets: how a proposed change to the model
// is opened, revised, decided and landed, entry by entry, in the log.
//
// The vocabulary:
//
//   - A change set is one proposed change to the model, with its own ID. It
//     is never edited in place; everything that happens to it is appended to
//     the log as an entry naming that ID.
//   - An entry is one appended fact about a change set: its envelope (see
//     package entry) plus a typed body. There are nine entry types: open,
//     ops, propose, accept, reject, withdraw, supersede, land and note.
//   - An operation (see package op) is one change to the model. An ops entry
//     carries operations and the IDs of operations it drops.
//   - A revision is an ops entry after the first: it adds operations,
//     replaces earlier ones by repeating their IDs, or drops them. An
//     operation keeps its ID for the life of its change set.
//   - A carve moves work from one change set to another: the source gets an
//     ops entry dropping the moved operations, the target an ops entry
//     re-creating them under new operation IDs, and each entry's carve
//     names the other change set and entry.
//   - A base version is the ID of the last landed operation that touched an
//     operation's target when the operation was written; a stale base is a
//     conflict, not a silent overwrite.
//   - Folding reads a change set's entries in causal order (their after
//     links, ties broken by entry ID) and computes where the change set
//     stands: its operations and its lifecycle state.
//   - The lifecycle state is draft after open, then proposed, accepted, and
//     finally one of landed, rejected, withdrawn or superseded.
//   - A change set's kind is fixed when it is opened. An intention is a
//     change someone wants made: it lands only once accepted. A record
//     describes something that already happened: it may land without review.
//
// This file holds the typed bodies. Each entry type has one Go type and one
// spelling, its canonical JSON with every field always present ([] when
// empty); DecodeBody rejects missing, unknown and wrongly typed fields and
// every rule a body can break on its own. Rules that need other entries, such
// as open coming first or which lifecycle moves are legal, belong to folding.
//
// The one optional field is an ops body's carve, present only on a carve's
// two entries. Format 1 cannot spell it: DecodeBody reads it only in an entry
// of format 2 or later and treats it as an unknown field in a format-1 entry,
// and EntryVersion tells a writer which format a body needs, so every other
// entry is still written as format 1.
package changeset

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"unicode/utf8"

	"github.com/dphaener/genesis/internal/canon"
	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
)

// ErrInvalid is wrapped by every error for a body that is not valid.
var ErrInvalid = errors.New("changeset: invalid")

// Kind is whether a change set is an intention or a record.
type Kind string

// The two change set kinds.
const (
	KindIntention Kind = "intention"
	KindRecord    Kind = "record"
)

// Kinds returns the two change set kinds in their documented order.
func Kinds() []Kind {
	return []Kind{KindIntention, KindRecord}
}

// Body is the typed body of one entry. It is implemented by the nine body
// types, as values.
type Body interface {
	// Type is the entry type this body belongs to.
	Type() entry.Type
	// Validate reports the first rule the body breaks on its own.
	Validate() error
}

// Open is the body of the open entry that starts a change set:
// {kind, intent, sources}.
type Open struct {
	Kind   Kind
	Intent string
	// Sources is sorted and unique. It decodes to nil when empty.
	Sources []id.ID
}

// Ops is the body of an ops entry: {ops, drop}, plus {carve} on either side
// of a carve. At least one of ops and drop is non-empty.
type Ops struct {
	// Ops holds operations with distinct IDs. It decodes to nil when empty.
	Ops []op.Op
	// Drop is sorted and unique, and shares no ID with Ops. It decodes to nil
	// when empty.
	Drop []id.ID
	// Carve is set only on a carve's two entries, and needs format 2. With
	// Carve.From set, Ops is non-empty; with Carve.To set, Drop is non-empty.
	Carve *Carve
}

// Carve links the two entries of one carve: {entry, to} on the source's
// entry, which drops the moved operations, or {entry, from} on the target's,
// which re-creates them. Exactly one of From and To is set.
type Carve struct {
	// From is the change set the work was carved from, on the target's entry.
	From id.ID
	// To is the change set the work was carved into, on the source's entry.
	To id.ID
	// Entry is the other side's entry.
	Entry id.ID
}

// Propose is the body of a propose entry: {}.
type Propose struct{}

// Accept is the body of an accept entry: {rule, reason}. Rule names the
// rule the decision was made under; Reason may be empty.
type Accept struct {
	Rule   string
	Reason string
}

// Reject is the body of a reject entry: {rule, reason}, as for Accept.
type Reject struct {
	Rule   string
	Reason string
}

// Withdraw is the body of a withdraw entry: {}.
type Withdraw struct{}

// Supersede is the body of a supersede entry: {by}, the change set that
// replaces this one.
type Supersede struct {
	By id.ID
}

// Land is the body of a land entry: {}.
type Land struct{}

// Note is the body of a note entry: {text}.
type Note struct {
	Text string
}

func (Open) Type() entry.Type      { return entry.TypeOpen }
func (Ops) Type() entry.Type       { return entry.TypeOps }
func (Propose) Type() entry.Type   { return entry.TypePropose }
func (Accept) Type() entry.Type    { return entry.TypeAccept }
func (Reject) Type() entry.Type    { return entry.TypeReject }
func (Withdraw) Type() entry.Type  { return entry.TypeWithdraw }
func (Supersede) Type() entry.Type { return entry.TypeSupersede }
func (Land) Type() entry.Type      { return entry.TypeLand }
func (Note) Type() entry.Type      { return entry.TypeNote }

// shapes lists each entry type's body fields in the order DecodeBody checks
// them.
var shapes = map[entry.Type][]string{
	entry.TypeOpen:      {"kind", "intent", "sources"},
	entry.TypeOps:       {"ops", "drop"},
	entry.TypePropose:   {},
	entry.TypeAccept:    {"rule", "reason"},
	entry.TypeReject:    {"rule", "reason"},
	entry.TypeWithdraw:  {},
	entry.TypeSupersede: {"by"},
	entry.TypeLand:      {},
	entry.TypeNote:      {"text"},
}

// carveShape lists a carve's fields in the order DecodeBody checks them;
// exactly one of the last two is present.
var carveShape = []string{"entry", "from", "to"}

// EntryVersion returns the lowest entry format that can spell b: 2 for an
// ops body with a carve, 1 for every other body. A writer stamps it as the
// entry's v.
func EntryVersion(b Body) int {
	if o, ok := b.(Ops); ok && o.Carve != nil {
		return 2
	}
	return 1
}

func invalid(field, format string, args ...any) error {
	return fmt.Errorf("%w: field %q: %s", ErrInvalid, field, fmt.Sprintf(format, args...))
}

func invalidWrap(field string, err error) error {
	return fmt.Errorf("%w: field %q: %w", ErrInvalid, field, err)
}

func opField(i int) string { return fmt.Sprintf("ops[%d]", i) }

// text checks a free-text field: valid UTF-8 and, if required, non-empty.
func text(field, s string, required bool) error {
	if required && s == "" {
		return invalid(field, "empty")
	}
	if !utf8.ValidString(s) {
		return invalid(field, "invalid UTF-8")
	}
	return nil
}

// sortedIDs checks that ids are non-zero, sorted and unique.
func sortedIDs(field string, ids []id.ID) error {
	for i, x := range ids {
		if x.IsZero() {
			return invalid(field, "zero id")
		}
		if i > 0 {
			switch id.Compare(ids[i-1], x) {
			case 0:
				return invalid(field, "duplicate %s", x)
			case 1:
				return invalid(field, "not sorted at %s", x)
			}
		}
	}
	return nil
}

// Validate reports the first rule o breaks.
func (o Open) Validate() error {
	if !slices.Contains(Kinds(), o.Kind) {
		return invalid("kind", "unknown kind %q", o.Kind)
	}
	if err := text("intent", o.Intent, true); err != nil {
		return err
	}
	return sortedIDs("sources", o.Sources)
}

// Validate reports the first rule o breaks.
func (o Ops) Validate() error {
	if len(o.Ops) == 0 && len(o.Drop) == 0 {
		return invalid("ops", "nothing to change: ops and drop are empty")
	}
	seen := make(map[id.ID]int, len(o.Ops))
	for i, x := range o.Ops {
		if err := x.Validate(); err != nil {
			return invalidWrap(opField(i), err)
		}
		if j, ok := seen[x.ID]; ok {
			return invalid("ops", "op id %s at %d repeats %s", x.ID, i, opField(j))
		}
		seen[x.ID] = i
	}
	if err := sortedIDs("drop", o.Drop); err != nil {
		return err
	}
	for _, x := range o.Drop {
		if j, ok := seen[x]; ok {
			return invalid("drop", "%s is also in %s", x, opField(j))
		}
	}
	if o.Carve == nil {
		return nil
	}
	if err := o.Carve.Validate(); err != nil {
		return err
	}
	if !o.Carve.From.IsZero() && len(o.Ops) == 0 {
		return invalid("ops", "empty, but a carve from %s re-creates operations", o.Carve.From)
	}
	if !o.Carve.To.IsZero() && len(o.Drop) == 0 {
		return invalid("drop", "empty, but a carve into %s drops operations", o.Carve.To)
	}
	return nil
}

// Validate reports the first rule c breaks.
func (c Carve) Validate() error {
	if c.Entry.IsZero() {
		return invalid("carve.entry", "zero id")
	}
	switch from, to := !c.From.IsZero(), !c.To.IsZero(); {
	case from && to:
		return invalid("carve", "both from and to")
	case !from && !to:
		return invalid("carve", "neither from nor to")
	}
	return nil
}

// Validate always succeeds: a propose body has no fields.
func (Propose) Validate() error { return nil }

// Validate reports the first rule a breaks.
func (a Accept) Validate() error { return decision(a.Rule, a.Reason) }

// Validate reports the first rule r breaks.
func (r Reject) Validate() error { return decision(r.Rule, r.Reason) }

func decision(rule, reason string) error {
	if err := text("rule", rule, true); err != nil {
		return err
	}
	return text("reason", reason, false)
}

// Validate always succeeds: a withdraw body has no fields.
func (Withdraw) Validate() error { return nil }

// Validate reports the first rule s breaks.
func (s Supersede) Validate() error {
	if s.By.IsZero() {
		return invalid("by", "zero id")
	}
	return nil
}

// Validate always succeeds: a land body has no fields.
func (Land) Validate() error { return nil }

// Validate reports the first rule n breaks.
func (n Note) Validate() error { return text("text", n.Text, true) }

// list spells a nil slice as [].
func list[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

// EncodeBody validates b and returns its canonical JSON without the final
// newline, ready for entry.Entry.Body. Nil and empty lists both encode as [].
func EncodeBody(b Body) (json.RawMessage, error) {
	if b == nil {
		return nil, fmt.Errorf("%w: nil body", ErrInvalid)
	}
	if err := b.Validate(); err != nil {
		return nil, err
	}
	var m map[string]any
	switch b := b.(type) {
	case Open:
		m = map[string]any{"kind": b.Kind, "intent": b.Intent, "sources": list(b.Sources)}
	case Ops:
		ops := make([]json.RawMessage, 0, len(b.Ops))
		for i, x := range b.Ops {
			raw, err := op.Encode(x)
			if err != nil {
				return nil, invalidWrap(opField(i), err)
			}
			ops = append(ops, raw)
		}
		m = map[string]any{"ops": ops, "drop": list(b.Drop)}
		if c := b.Carve; c != nil {
			cm := map[string]any{"entry": c.Entry}
			if !c.From.IsZero() {
				cm["from"] = c.From
			} else {
				cm["to"] = c.To
			}
			m["carve"] = cm
		}
	case Accept:
		m = map[string]any{"rule": b.Rule, "reason": b.Reason}
	case Reject:
		m = map[string]any{"rule": b.Rule, "reason": b.Reason}
	case Supersede:
		m = map[string]any{"by": b.By}
	case Note:
		m = map[string]any{"text": b.Text}
	case Propose, Withdraw, Land:
		m = map[string]any{}
	default:
		return nil, fmt.Errorf("%w: unknown body type %T", ErrInvalid, b)
	}
	out, err := canon.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalid, err)
	}
	return json.RawMessage(bytes.TrimSuffix(out, []byte("\n"))), nil
}

// DecodeBody parses e.Body as the body of an entry of type e.Type. The body
// must be canonical on its own, as entry.Decode leaves it. Missing fields are
// reported in schema order, then the first unknown field in sorted order,
// then each field's kind and the body's rules. An invalid operation is
// reported on its index, as field "ops[i]", and its error still wraps
// op.ErrInvalid. An ops body may also hold carve when e.V is 2 or later; in
// a format-1 entry carve is an unknown field. A carve's own fields are
// checked the same way and reported as "carve.<field>".
func DecodeBody(e entry.Entry) (Body, error) {
	shape, ok := shapes[e.Type]
	if !ok {
		return nil, invalid("type", "unknown entry type %q", e.Type)
	}
	data, err := canon.Canonicalize(e.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalid, err)
	}
	if !bytes.Equal(bytes.TrimSuffix(data, []byte("\n")), e.Body) {
		return nil, fmt.Errorf("%w: %w", ErrInvalid, canon.ErrNotCanonical)
	}
	if data[0] != '{' {
		return nil, fmt.Errorf("%w: not an object", ErrInvalid)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalid, err)
	}
	for _, k := range shape {
		if _, ok := fields[k]; !ok {
			return nil, invalid(k, "missing")
		}
	}
	known := shape
	if e.Type == entry.TypeOps && e.V >= 2 {
		known = append(slices.Clip(shape), "carve")
	}
	if name, ok := unknownField(fields, known); ok {
		return nil, invalid(name, "unknown field")
	}

	var b Body
	switch e.Type {
	case entry.TypeOpen:
		b, err = decodeOpen(fields)
	case entry.TypeOps:
		b, err = decodeOps(fields)
	case entry.TypePropose:
		b = Propose{}
	case entry.TypeAccept, entry.TypeReject:
		var rule, reason string
		if rule, err = str(fields["rule"]); err != nil {
			return nil, invalidWrap("rule", err)
		}
		if reason, err = str(fields["reason"]); err != nil {
			return nil, invalidWrap("reason", err)
		}
		if e.Type == entry.TypeAccept {
			b = Accept{Rule: rule, Reason: reason}
		} else {
			b = Reject{Rule: rule, Reason: reason}
		}
	case entry.TypeWithdraw:
		b = Withdraw{}
	case entry.TypeSupersede:
		var by id.ID
		if by, err = ulid(fields["by"]); err != nil {
			return nil, invalidWrap("by", err)
		}
		b = Supersede{By: by}
	case entry.TypeLand:
		b = Land{}
	case entry.TypeNote:
		var s string
		if s, err = str(fields["text"]); err != nil {
			return nil, invalidWrap("text", err)
		}
		b = Note{Text: s}
	}
	if err != nil {
		return nil, err
	}
	if err := b.Validate(); err != nil {
		return nil, err
	}
	return b, nil
}

func decodeOpen(fields map[string]json.RawMessage) (Body, error) {
	kind, err := str(fields["kind"])
	if err != nil {
		return nil, invalidWrap("kind", err)
	}
	intent, err := str(fields["intent"])
	if err != nil {
		return nil, invalidWrap("intent", err)
	}
	sources, err := ids(fields["sources"])
	if err != nil {
		return nil, invalidWrap("sources", err)
	}
	return Open{Kind: Kind(kind), Intent: intent, Sources: sources}, nil
}

func decodeOps(fields map[string]json.RawMessage) (Body, error) {
	items, err := array(fields["ops"])
	if err != nil {
		return nil, invalidWrap("ops", err)
	}
	var ops []op.Op
	for i, r := range items {
		// The op sits indented inside the body; op.Decode canonicalizes it on
		// its own.
		o, err := op.Decode(r)
		if err != nil {
			return nil, invalidWrap(opField(i), err)
		}
		ops = append(ops, o)
	}
	drop, err := ids(fields["drop"])
	if err != nil {
		return nil, invalidWrap("drop", err)
	}
	b := Ops{Ops: ops, Drop: drop}
	if raw, ok := fields["carve"]; ok {
		if b.Carve, err = decodeCarve(raw); err != nil {
			return nil, err
		}
	}
	return b, nil
}

// decodeCarve parses an ops body's carve: missing entry, then the first
// unknown field in sorted order, then both or neither of from and to, then
// each field's kind. A zero ID is reported on its own field.
func decodeCarve(raw json.RawMessage) (*Carve, error) {
	if len(raw) == 0 || raw[0] != '{' {
		return nil, invalid("carve", "not an object")
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return nil, invalidWrap("carve", err)
	}
	if _, ok := fields["entry"]; !ok {
		return nil, invalid("carve.entry", "missing")
	}
	if name, ok := unknownField(fields, carveShape); ok {
		return nil, invalid("carve."+name, "unknown field")
	}
	_, from := fields["from"]
	_, to := fields["to"]
	switch {
	case from && to:
		return nil, invalid("carve", "both from and to")
	case !from && !to:
		return nil, invalid("carve", "neither from nor to")
	}
	var c Carve
	for _, f := range []struct {
		name string
		dst  *id.ID
	}{{"entry", &c.Entry}, {"from", &c.From}, {"to", &c.To}} {
		r, ok := fields[f.name]
		if !ok {
			continue
		}
		x, err := ulid(r)
		if err != nil {
			return nil, invalidWrap("carve."+f.name, err)
		}
		if x.IsZero() {
			return nil, invalid("carve."+f.name, "zero id")
		}
		*f.dst = x
	}
	return &c, nil
}

// unknownField returns the first key of fields, in sorted order, that is not
// in known.
func unknownField(fields map[string]json.RawMessage, known []string) (string, bool) {
	var names []string
	for k := range fields {
		if !slices.Contains(known, k) {
			names = append(names, k)
		}
	}
	if len(names) == 0 {
		return "", false
	}
	slices.Sort(names)
	return names[0], true
}

var (
	errNotString = errors.New("not a string")
	errNotArray  = errors.New("not an array")
)

// str decodes a JSON string value, rejecting every other kind.
func str(raw json.RawMessage) (string, error) {
	if len(raw) == 0 || raw[0] != '"' {
		return "", errNotString
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", err
	}
	return s, nil
}

func ulid(raw json.RawMessage) (id.ID, error) {
	s, err := str(raw)
	if err != nil {
		return id.ID{}, err
	}
	return id.Parse(s)
}

// array decodes a JSON array into its raw elements.
func array(raw json.RawMessage) ([]json.RawMessage, error) {
	if len(raw) == 0 || raw[0] != '[' {
		return nil, errNotArray
	}
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// ids decodes a JSON array of ULIDs. An empty array decodes to nil.
func ids(raw json.RawMessage) ([]id.ID, error) {
	items, err := array(raw)
	if err != nil {
		return nil, err
	}
	var out []id.ID
	for i, r := range items {
		x, err := ulid(r)
		if err != nil {
			return nil, fmt.Errorf("index %d: %w", i, err)
		}
		out = append(out, x)
	}
	return out, nil
}
