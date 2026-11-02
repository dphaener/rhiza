// Package op defines the six operations that change the model: creating,
// updating, moving and retiring nodes, and creating and retiring links.
//
// An operation is one change to the model. It has its own ID, which stays the
// same for the life of its change set: a later revision replaces an
// operation by repeating its ID. Nodes and links are never deleted, only
// retired.
//
// Every operation except a create names a base version: the ID of the last
// landed operation that touched its target when the change was written. A
// change made against a stale base is a conflict, not a silent overwrite.
//
// An operation has exactly one spelling, its canonical JSON, with every
// field always present (null, [] or {} when empty). Decode rejects missing,
// unknown and wrongly typed fields and every rule below. This package checks
// shape only: whether a target, parent, base or endpoint exists, is current
// or has the right kind needs the model, and is projection's job.
package op

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"unicode/utf8"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/key"
)

// ErrInvalid is wrapped by every error for an operation that is not valid.
var ErrInvalid = errors.New("op: invalid")

// Type names what an operation does. Its JSON field is "op".
type Type string

// The six operation types.
const (
	TypeNodeCreate Type = "node.create"
	TypeNodeUpdate Type = "node.update"
	TypeNodeMove   Type = "node.move"
	TypeNodeRetire Type = "node.retire"
	TypeLinkCreate Type = "link.create"
	TypeLinkRetire Type = "link.retire"
)

// Types returns the six operation types in their documented order.
func Types() []Type {
	return []Type{TypeNodeCreate, TypeNodeUpdate, TypeNodeMove, TypeNodeRetire, TypeLinkCreate, TypeLinkRetire}
}

// Kind is what a node is.
type Kind string

// The ten node kinds.
const (
	KindIntent     Kind = "intent"
	KindCapability Kind = "capability"
	KindConstraint Kind = "constraint"
	KindDecision   Kind = "decision"
	KindQuestion   Kind = "question"
	KindEvidence   Kind = "evidence"
	KindSource     Kind = "source"
	KindActor      Kind = "actor"
	KindAnchor     Kind = "anchor"
	KindNote       Kind = "note"
)

// Kinds returns the ten node kinds in their documented order.
func Kinds() []Kind {
	return []Kind{KindIntent, KindCapability, KindConstraint, KindDecision, KindQuestion, KindEvidence, KindSource, KindActor, KindAnchor, KindNote}
}

// HasParent reports whether nodes of kind k sit in a tree, with a parent
// (or none, for a root) and a sibling key. Only capabilities and intents do.
func HasParent(k Kind) bool {
	return k == KindCapability || k == KindIntent
}

// LinkType is what a link says about its endpoints.
type LinkType string

// The fifteen link types.
const (
	LinkServes         LinkType = "serves"
	LinkConstrains     LinkType = "constrains"
	LinkExempts        LinkType = "exempts"
	LinkDependsOn      LinkType = "depends-on"
	LinkRelatesTo      LinkType = "relates-to"
	LinkDecidedBy      LinkType = "decided-by"
	LinkAnswers        LinkType = "answers"
	LinkBlocks         LinkType = "blocks"
	LinkSupersedes     LinkType = "supersedes"
	LinkDerivedFrom    LinkType = "derived-from"
	LinkDiscoveredFrom LinkType = "discovered-from"
	LinkImplements     LinkType = "implements"
	LinkVerifies       LinkType = "verifies"
	LinkObserves       LinkType = "observes"
	LinkDiscusses      LinkType = "discusses"
)

// LinkTypes returns the fifteen link types in their documented order.
func LinkTypes() []LinkType {
	return []LinkType{LinkServes, LinkConstrains, LinkExempts, LinkDependsOn, LinkRelatesTo, LinkDecidedBy, LinkAnswers,
		LinkBlocks, LinkSupersedes, LinkDerivedFrom, LinkDiscoveredFrom, LinkImplements, LinkVerifies, LinkObserves, LinkDiscusses}
}

// Op is one operation. It is a single flat struct: each type uses the fields
// its shape lists and leaves every other field at its zero value, which
// Validate enforces. Null is the zero value: a zero ID for parent and
// decision, "" for key and for the update's handle.
//
// Decode always returns nil, never empty, for Fields, Set and Unset when they
// hold nothing, so equal operations are reflect.DeepEqual. Encode accepts
// nil and empty alike.
type Op struct {
	ID   id.ID
	Type Type

	// Node is the target of the node operations, Link of the link ones.
	Node id.ID
	Link id.ID
	// Base is the version the change was written against: update, move,
	// retire and link.retire.
	Base id.ID

	// node.create: Kind, Handle, Parent, Key, Fields.
	// node.update: Handle ("" leaves it unchanged), Set, Unset.
	// node.move: Parent, Key.
	Kind   Kind
	Handle string
	Parent id.ID
	Key    string
	// Field values are canonical JSON without a final newline.
	Fields map[string]json.RawMessage
	Set    map[string]json.RawMessage
	// Unset is sorted and unique.
	Unset []string

	// node.retire.
	Reason string

	// link.create: LinkType is its "type" field.
	LinkType LinkType
	From     id.ID
	To       id.ID
	Note     string
	Decision id.ID
}

// shapes lists each type's JSON fields in the order Decode checks them.
var shapes = map[Type][]string{
	TypeNodeCreate: {"id", "op", "node", "kind", "handle", "parent", "key", "fields"},
	TypeNodeUpdate: {"id", "op", "node", "base", "handle", "set", "unset"},
	TypeNodeMove:   {"id", "op", "node", "base", "parent", "key"},
	TypeNodeRetire: {"id", "op", "node", "base", "reason"},
	TypeLinkCreate: {"id", "op", "link", "type", "from", "to", "note", "decision"},
	TypeLinkRetire: {"id", "op", "link", "base"},
}

// Target returns the ID of the node or link o acts on, or the zero ID for
// an unknown type.
func (o Op) Target() id.ID {
	switch o.Type {
	case TypeNodeCreate, TypeNodeUpdate, TypeNodeMove, TypeNodeRetire:
		return o.Node
	case TypeLinkCreate, TypeLinkRetire:
		return o.Link
	}
	return id.ID{}
}

func invalid(field, format string, args ...any) error {
	return fmt.Errorf("%w: field %q: %s", ErrInvalid, field, fmt.Sprintf(format, args...))
}

func invalidWrap(field string, err error) error {
	return fmt.Errorf("%w: field %q: %w", ErrInvalid, field, err)
}

const maxNameLen = 64

func isLower(c byte) bool { return 'a' <= c && c <= 'z' }
func isDigit(c byte) bool { return '0' <= c && c <= '9' }

// validHandle reports whether s is one slug segment: ^[a-z0-9]+(-[a-z0-9]+)*$,
// 1 to 64 bytes.
func validHandle(s string) error {
	if s == "" || len(s) > maxNameLen {
		return fmt.Errorf("handle %.16q: length %d, want 1 to %d", s, len(s), maxNameLen)
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case isLower(c) || isDigit(c):
		case c == '-' && i > 0 && i < len(s)-1 && s[i-1] != '-':
		default:
			return fmt.Errorf("handle %q is not lower-case words joined by single dashes", s)
		}
	}
	return nil
}

// validName reports whether s is a field name: ^[a-z][a-z0-9_]*$, 1 to 64
// bytes.
func validName(s string) error {
	if s == "" || len(s) > maxNameLen {
		return fmt.Errorf("name %.16q: length %d, want 1 to %d", s, len(s), maxNameLen)
	}
	if !isLower(s[0]) {
		return fmt.Errorf("name %q does not start with a lower-case letter", s)
	}
	for i := 1; i < len(s); i++ {
		if c := s[i]; !isLower(c) && !isDigit(c) && c != '_' {
			return fmt.Errorf("name %q has %q, want only a-z, 0-9 and _", s, c)
		}
	}
	return nil
}

// validValues checks a field map's names and that each value is canonical.
// With canonical set the values are known to be canonical (object cut them
// out of canonical text) and only their names and presence are checked.
func validValues(field string, m map[string]json.RawMessage, canonical bool) error {
	names := make([]string, 0, len(m))
	for k := range m {
		names = append(names, k)
	}
	slices.Sort(names)
	for _, k := range names {
		if err := validName(k); err != nil {
			return invalidWrap(field, err)
		}
		v := m[k]
		if len(v) == 0 {
			return invalid(field, "value of %q is empty", k)
		}
		if canonical {
			continue
		}
		out, err := canon.Canonicalize(v)
		if err != nil {
			return invalid(field, "value of %q: %v", k, err)
		}
		if !bytes.Equal(bytes.TrimSuffix(out, []byte("\n")), v) {
			return invalid(field, "value of %q: %v", k, canon.ErrNotCanonical)
		}
	}
	return nil
}

// unused reports the first field set on o that o's shape does not have.
func (o Op) unused() error {
	has := shapes[o.Type]
	for _, f := range []struct {
		name string
		set  bool
	}{
		{"node", !o.Node.IsZero()},
		{"link", !o.Link.IsZero()},
		{"base", !o.Base.IsZero()},
		{"kind", o.Kind != ""},
		{"handle", o.Handle != ""},
		{"parent", !o.Parent.IsZero()},
		{"key", o.Key != ""},
		{"fields", o.Fields != nil},
		{"set", o.Set != nil},
		{"unset", o.Unset != nil},
		{"reason", o.Reason != ""},
		{"type", o.LinkType != ""},
		{"from", !o.From.IsZero()},
		{"to", !o.To.IsZero()},
		{"note", o.Note != ""},
		{"decision", !o.Decision.IsZero()},
	} {
		if f.set && !slices.Contains(has, f.name) {
			return invalid(f.name, "not part of %s", o.Type)
		}
	}
	return nil
}

// Validate reports the first way o breaks its type's rules. It checks shape
// only, never whether anything it names exists.
func (o Op) Validate() error {
	return o.validate(false)
}

// validate is Validate; with canonical set, o's field values are known to
// be canonical (decode cut them out of canonical text) and are not checked
// to be, which is most of Validate's cost on a create or an update.
func (o Op) validate(canonical bool) error {
	if o.ID.IsZero() {
		return invalid("id", "zero id")
	}
	if _, ok := shapes[o.Type]; !ok {
		return invalid("op", "unknown type %q", o.Type)
	}
	if err := o.unused(); err != nil {
		return err
	}
	if t := o.Target(); t.IsZero() {
		return invalid(shapes[o.Type][2], "zero id")
	}
	if o.Type != TypeNodeCreate && o.Type != TypeLinkCreate && o.Base.IsZero() {
		return invalid("base", "zero id")
	}

	switch o.Type {
	case TypeNodeCreate:
		if !slices.Contains(Kinds(), o.Kind) {
			return invalid("kind", "unknown kind %q", o.Kind)
		}
		if err := validHandle(o.Handle); err != nil {
			return invalidWrap("handle", err)
		}
		if HasParent(o.Kind) {
			if err := key.Valid(o.Key); err != nil {
				return invalidWrap("key", err)
			}
		} else {
			if !o.Parent.IsZero() {
				return invalid("parent", "a %s has no parent", o.Kind)
			}
			if o.Key != "" {
				return invalid("key", "a %s has no key", o.Kind)
			}
		}
		return validValues("fields", o.Fields, canonical)

	case TypeNodeUpdate:
		if o.Handle != "" {
			if err := validHandle(o.Handle); err != nil {
				return invalidWrap("handle", err)
			}
		}
		if err := validValues("set", o.Set, canonical); err != nil {
			return err
		}
		for i, name := range o.Unset {
			if err := validName(name); err != nil {
				return invalidWrap("unset", err)
			}
			if i > 0 {
				switch prev := o.Unset[i-1]; {
				case prev == name:
					return invalid("unset", "duplicate %q", name)
				case prev > name:
					return invalid("unset", "not sorted at %q", name)
				}
			}
			if _, ok := o.Set[name]; ok {
				return invalid("unset", "%q is also in set", name)
			}
		}
		if o.Handle == "" && len(o.Set) == 0 && len(o.Unset) == 0 {
			return invalid("set", "nothing to change: handle is null, set and unset are empty")
		}

	case TypeNodeMove:
		if err := key.Valid(o.Key); err != nil {
			return invalidWrap("key", err)
		}

	case TypeNodeRetire:
		if o.Reason == "" {
			return invalid("reason", "empty")
		}
		if !utf8.ValidString(o.Reason) {
			return invalid("reason", "invalid UTF-8")
		}

	case TypeLinkCreate:
		if !slices.Contains(LinkTypes(), o.LinkType) {
			return invalid("type", "unknown link type %q", o.LinkType)
		}
		if o.From.IsZero() {
			return invalid("from", "zero id")
		}
		if o.To.IsZero() {
			return invalid("to", "zero id")
		}
		if o.From == o.To {
			return invalid("to", "same as from %s", o.From)
		}
		if !utf8.ValidString(o.Note) {
			return invalid("note", "invalid UTF-8")
		}
		if o.Note == "" && (o.LinkType == LinkRelatesTo || o.LinkType == LinkExempts) {
			return invalid("note", "empty, required on %s", o.LinkType)
		}
		if o.LinkType == LinkExempts && o.Decision.IsZero() {
			return invalid("decision", "null, required on %s", o.LinkType)
		}
		if o.LinkType != LinkExempts && !o.Decision.IsZero() {
			return invalid("decision", "only %s names a decision, not %s", LinkExempts, o.LinkType)
		}
	}
	return nil
}

// nullable spells a zero ID as null.
func nullable(i id.ID) any {
	if i.IsZero() {
		return nil
	}
	return i
}

func values(m map[string]json.RawMessage) map[string]json.RawMessage {
	if m == nil {
		return map[string]json.RawMessage{}
	}
	return m
}

// Encode validates o and returns its canonical JSON without the final
// newline.
func Encode(o Op) (json.RawMessage, error) {
	if err := o.Validate(); err != nil {
		return nil, err
	}
	m := map[string]any{"id": o.ID, "op": o.Type}
	switch o.Type {
	case TypeNodeCreate:
		var k any
		if o.Key != "" {
			k = o.Key
		}
		m["node"], m["kind"], m["handle"], m["parent"], m["key"], m["fields"] =
			o.Node, o.Kind, o.Handle, nullable(o.Parent), k, values(o.Fields)
	case TypeNodeUpdate:
		var h any
		if o.Handle != "" {
			h = o.Handle
		}
		unset := o.Unset
		if unset == nil {
			unset = []string{}
		}
		m["node"], m["base"], m["handle"], m["set"], m["unset"] = o.Node, o.Base, h, values(o.Set), unset
	case TypeNodeMove:
		m["node"], m["base"], m["parent"], m["key"] = o.Node, o.Base, nullable(o.Parent), o.Key
	case TypeNodeRetire:
		m["node"], m["base"], m["reason"] = o.Node, o.Base, o.Reason
	case TypeLinkCreate:
		m["link"], m["type"], m["from"], m["to"], m["note"], m["decision"] =
			o.Link, o.LinkType, o.From, o.To, o.Note, nullable(o.Decision)
	case TypeLinkRetire:
		m["link"], m["base"] = o.Link, o.Base
	}
	out, err := canon.Marshal(m)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(bytes.TrimSuffix(out, []byte("\n"))), nil
}

// Decode parses one operation as it appears inside an entry body. The bytes
// need not be canonical on their own, since the op sits indented inside its
// entry; Decode canonicalizes them first, so invalid JSON and duplicate keys
// are still rejected. The op type is read first, as it decides the shape;
// then missing fields are reported in schema order, then the first unknown
// field in sorted order, then each field's kind and the type's rules.
func Decode(raw json.RawMessage) (Op, error) {
	data, err := canon.Canonicalize(raw)
	if err != nil {
		return Op{}, fmt.Errorf("%w: %w", ErrInvalid, err)
	}
	return decode(bytes.TrimSuffix(data, []byte("\n")), false)
}

// DecodeCanonical is Decode for an op cut out of a canonical entry body,
// such as one the index stored after entry.Decode accepted it: canonical
// text, indented by its place in the body. It runs every check Decode does
// but the canonical-form ones, of the op and of its field values, which
// are most of Decode's cost; on bytes that are not canonical its result is
// unspecified, though it still fails on bytes the checks it runs reject.
func DecodeCanonical(raw json.RawMessage) (Op, error) {
	if len(raw) == 0 {
		return Op{}, fmt.Errorf("%w: %w", ErrInvalid, canon.Check(raw))
	}
	return decode(deindent(raw), true)
}

// deindent returns the canonical object or array data, written at some
// indentation inside a larger canonical document, as it is written on its
// own: the indentation of its closing line, which every line but the first
// starts with, is removed from each. Canonical text holds no raw newline
// inside a string, so every newline in data is structural. Data without a
// newline (a scalar, or an empty object) is returned as it is.
func deindent(data []byte) []byte {
	last := bytes.LastIndexByte(data, '\n')
	if last < 0 {
		return data
	}
	indent := data[last+1 : len(data)-1]
	if len(indent) == 0 || bytes.IndexFunc(indent, func(r rune) bool { return r != ' ' }) >= 0 {
		return data
	}
	return bytes.ReplaceAll(data, append([]byte{'\n'}, indent...), []byte{'\n'})
}

// decode parses data, which Decode has canonicalized or DecodeCanonical
// takes to be canonical (canonical). Both read it with the scanners, which
// split canonical text without parsing it, and fall back to encoding/json
// and its errors for text not shaped as they expect.
func decode(data []byte, canonical bool) (Op, error) {
	if len(data) == 0 || data[0] != '{' {
		return Op{}, fmt.Errorf("%w: not an object", ErrInvalid)
	}
	fields := scanMembers(data)
	if fields == nil {
		if err := json.Unmarshal(data, &fields); err != nil {
			return Op{}, fmt.Errorf("%w: %w", ErrInvalid, err)
		}
	}

	r, ok := fields["op"]
	if !ok {
		return Op{}, invalid("op", "missing")
	}
	typ, err := str(r)
	if err != nil {
		return Op{}, invalidWrap("op", err)
	}
	shape, ok := shapes[Type(typ)]
	if !ok {
		return Op{}, invalid("op", "unknown type %q", typ)
	}
	for _, k := range shape {
		if _, ok := fields[k]; !ok {
			return Op{}, invalid(k, "missing")
		}
	}
	if len(fields) != len(shape) {
		names := make([]string, 0, len(fields))
		for k := range fields {
			if !slices.Contains(shape, k) {
				names = append(names, k)
			}
		}
		slices.Sort(names)
		return Op{}, invalid(names[0], "unknown field")
	}

	o := Op{Type: Type(typ)}
	for _, name := range shape {
		r := fields[name]
		var err error
		switch name {
		case "op":
		case "id":
			o.ID, err = ulid(r)
		case "node":
			o.Node, err = ulid(r)
		case "link":
			o.Link, err = ulid(r)
		case "base":
			o.Base, err = ulid(r)
		case "from":
			o.From, err = ulid(r)
		case "to":
			o.To, err = ulid(r)
		case "parent":
			o.Parent, err = nullULID(r)
		case "decision":
			o.Decision, err = nullULID(r)
		case "kind":
			var s string
			s, err = str(r)
			o.Kind = Kind(s)
		case "type":
			var s string
			s, err = str(r)
			o.LinkType = LinkType(s)
		case "reason":
			o.Reason, err = str(r)
		case "note":
			o.Note, err = str(r)
		case "handle":
			// Null only means "unchanged" on an update; a create needs one.
			if o.Type == TypeNodeUpdate && isNull(r) {
				break
			}
			if o.Handle, err = str(r); err == nil {
				err = validHandle(o.Handle)
			}
		case "key":
			// Null only means "no key" on a create; a move needs one.
			if o.Type == TypeNodeCreate && isNull(r) {
				break
			}
			if o.Key, err = str(r); err == nil {
				err = key.Valid(o.Key)
			}
		case "fields":
			o.Fields, err = object(r)
		case "set":
			o.Set, err = object(r)
		case "unset":
			o.Unset, err = stringList(r)
		}
		if err != nil {
			return Op{}, invalidWrap(name, err)
		}
	}

	if err := o.validate(canonical); err != nil {
		return Op{}, err
	}
	return o, nil
}

var (
	errNotString = errors.New("not a string")
	errNotObject = errors.New("not an object")
	errNotArray  = errors.New("not an array")
	errZeroID    = errors.New("zero id")
)

func isNull(raw json.RawMessage) bool {
	return string(raw) == "null"
}

// str decodes a JSON string value, rejecting every other kind.
func str(raw json.RawMessage) (string, error) {
	if len(raw) == 0 || raw[0] != '"' {
		return "", errNotString
	}
	// A string without an escape is its own text between the quotes; the
	// bytes are valid JSON by here, so this is what Unmarshal would give.
	if n := len(raw); n >= 2 && raw[n-1] == '"' && bytes.IndexByte(raw, '\\') < 0 {
		return string(raw[1 : n-1]), nil
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

// nullULID decodes null or a ULID. A zero ULID is rejected, since the zero
// ID is how null is held.
func nullULID(raw json.RawMessage) (id.ID, error) {
	if isNull(raw) {
		return id.ID{}, nil
	}
	i, err := ulid(raw)
	if err == nil && i.IsZero() {
		err = errZeroID
	}
	return i, err
}

// object decodes a JSON object of field values, each in its own canonical
// form. An empty object decodes to nil. raw is a member of a canonical op,
// so each value is canonical once the two levels of indentation it sits
// under are removed (deindent); a value the scanner cannot split falls
// back to encoding/json and canon, as the op's own text does.
func object(raw json.RawMessage) (map[string]json.RawMessage, error) {
	if len(raw) == 0 || raw[0] != '{' {
		return nil, errNotObject
	}
	m := scanMembers(raw)
	if m == nil {
		if err := json.Unmarshal(raw, &m); err != nil {
			return nil, err
		}
		for k, v := range m {
			out, err := canon.Canonicalize(v)
			if err != nil {
				return nil, err
			}
			m[k] = json.RawMessage(bytes.TrimSuffix(out, []byte("\n")))
		}
	} else {
		for k, v := range m {
			m[k] = deindent(v)
		}
	}
	if len(m) == 0 {
		return nil, nil
	}
	return m, nil
}

// stringList decodes a JSON array of strings. An empty array decodes to nil.
func stringList(raw json.RawMessage) ([]string, error) {
	if len(raw) == 0 || raw[0] != '[' {
		return nil, errNotArray
	}
	items := scanElements(raw)
	if items == nil {
		if err := json.Unmarshal(raw, &items); err != nil {
			return nil, err
		}
	}
	if len(items) == 0 {
		return nil, nil
	}
	out := make([]string, 0, len(items))
	for _, r := range items {
		s, err := str(r)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

// scanMembers splits a canonical JSON object into its members without
// parsing their values, which come back as slices of data. It is nil when
// data is not shaped as expected (a key with an escape, say), so that the
// caller falls back to json.Unmarshal and its errors. Canonical text has
// no duplicate keys, so the map holds what Unmarshal would give.
func scanMembers(data []byte) map[string]json.RawMessage {
	fields := make(map[string]json.RawMessage, 8)
	i := skipSpace(data, 1)
	for i < len(data) && data[i] != '}' {
		if data[i] != '"' {
			return nil
		}
		end := scanValue(data, i)
		if end < 0 || bytes.IndexByte(data[i:end], '\\') >= 0 {
			return nil
		}
		key := string(data[i+1 : end-1])
		i = skipSpace(data, end)
		if i >= len(data) || data[i] != ':' {
			return nil
		}
		i = skipSpace(data, i+1)
		if end = scanValue(data, i); end < 0 {
			return nil
		}
		fields[key] = data[i:end]
		if i = skipSpace(data, end); i < len(data) && data[i] == ',' {
			i = skipSpace(data, i+1)
		}
	}
	if i >= len(data) {
		return nil
	}
	return fields
}

// scanElements is scanMembers for an array: its elements, never nil for an
// array that scans, so an empty one is [].
func scanElements(data []byte) []json.RawMessage {
	out := []json.RawMessage{}
	i := skipSpace(data, 1)
	for i < len(data) && data[i] != ']' {
		end := scanValue(data, i)
		if end < 0 {
			return nil
		}
		out = append(out, data[i:end])
		if i = skipSpace(data, end); i < len(data) && data[i] == ',' {
			i = skipSpace(data, i+1)
		}
	}
	if i >= len(data) {
		return nil
	}
	return out
}

// skipSpace skips the spaces and newlines of canonical text from data[i].
func skipSpace(data []byte, i int) int {
	for i < len(data) && (data[i] == ' ' || data[i] == '\n') {
		i++
	}
	return i
}

// scanValue returns the end of the JSON value starting at data[i], or -1
// when it does not end.
func scanValue(data []byte, i int) int {
	if i >= len(data) {
		return -1
	}
	depth := 0
	for j := i; j < len(data); j++ {
		switch data[j] {
		case '"':
			for j++; j < len(data) && data[j] != '"'; j++ {
				if data[j] == '\\' {
					j++
				}
			}
			if j >= len(data) {
				return -1
			}
			if depth == 0 {
				return j + 1
			}
		case '{', '[':
			depth++
		case '}', ']':
			depth--
			if depth <= 0 {
				if depth < 0 {
					return j // a scalar ended at its container's close
				}
				return j + 1
			}
		case ',', ' ', '\n':
			if depth == 0 {
				return j
			}
		}
	}
	return -1
}
