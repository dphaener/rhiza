// Package entry defines the log entry envelope: the fixed fields every
// entry carries around a free-form JSON object body.
//
// An entry has exactly one spelling, its canonical JSON. Encode produces it
// and Decode accepts nothing else: non-canonical bytes, missing or unknown
// fields and values of the wrong kind are all rejected rather than repaired.
// What a body means for each type is not this package's concern.
//
// The envelope is the same in every format this package knows; v says which
// format the entry is written in, so a reader can refuse an entry from a
// newer Genesis instead of misreading it. Format 2 adds no envelope field: it
// marks an entry whose body uses something format 1 cannot spell (an ops
// body's carve). A writer stamps the lowest format that can spell the entry,
// so every entry format 1 can spell is still written as format 1.
package entry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/dphaener/genesis/internal/canon"
	"github.com/dphaener/genesis/internal/id"
)

// FormatVersion is the newest entry format this package reads and writes.
// Every format from 1 through FormatVersion is read and written.
const FormatVersion = 2

// ErrInvalid is wrapped by every error for an entry that is not a valid
// entry of a format from 1 through FormatVersion.
var ErrInvalid = errors.New("entry: invalid")

// ErrNewerFormat is wrapped by Decode's error when the entry declares a
// format newer than FormatVersion. It does not wrap ErrInvalid: the entry
// may be perfectly valid to a newer Genesis.
var ErrNewerFormat = errors.New("entry: newer format")

// Type names what an entry records.
type Type string

// The nine entry types, the same in every format.
const (
	TypeOpen      Type = "open"
	TypeOps       Type = "ops"
	TypePropose   Type = "propose"
	TypeAccept    Type = "accept"
	TypeReject    Type = "reject"
	TypeWithdraw  Type = "withdraw"
	TypeSupersede Type = "supersede"
	TypeLand      Type = "land"
	TypeNote      Type = "note"
)

// Types returns the nine entry types in their documented order.
func Types() []Type {
	return []Type{TypeOpen, TypeOps, TypePropose, TypeAccept, TypeReject, TypeWithdraw, TypeSupersede, TypeLand, TypeNote}
}

// timeLayout is how at is spelled: RFC 3339 in UTC, to the second, with Z.
const timeLayout = "2006-01-02T15:04:05Z"

// Entry is one log entry.
type Entry struct {
	V         int
	ID        id.ID
	ChangeSet id.ID
	Actor     id.ID
	For       id.ID
	At        time.Time
	// After lists the entries this one follows, sorted and unique.
	After []id.ID
	Type  Type
	// Body is a JSON object in canonical form, without the final newline.
	Body json.RawMessage
}

// keys are the envelope's fields in the order Decode checks them.
var keys = []string{"v", "id", "change_set", "at", "actor", "for", "after", "type", "body"}

func invalid(field, format string, args ...any) error {
	return fmt.Errorf("%w: field %q: %s", ErrInvalid, field, fmt.Sprintf(format, args...))
}

func invalidWrap(field string, err error) error {
	return fmt.Errorf("%w: field %q: %w", ErrInvalid, field, err)
}

// Validate reports the first way e breaks the envelope rules. V must be a
// format from 1 through FormatVersion; whether the body needs that format is
// the body's concern. The time embedded in a ULID is deliberately not
// compared with At.
func (e Entry) Validate() error {
	if e.V < 1 || e.V > FormatVersion {
		return invalid("v", "%d, want 1 through %d", e.V, FormatVersion)
	}
	for _, f := range []struct {
		name string
		id   id.ID
	}{{"id", e.ID}, {"change_set", e.ChangeSet}, {"actor", e.Actor}, {"for", e.For}} {
		if f.id.IsZero() {
			return invalid(f.name, "zero id")
		}
	}
	if e.At.Location() != time.UTC {
		return invalid("at", "location %s, want UTC", e.At.Location())
	}
	if e.At.Nanosecond() != 0 {
		return invalid("at", "has a sub-second part")
	}
	// Outside these years the layout cannot spell the time and read it back.
	if y := e.At.Year(); y < 0 || y > 9999 {
		return invalid("at", "year %d out of range", y)
	}
	for i := 1; i < len(e.After); i++ {
		switch id.Compare(e.After[i-1], e.After[i]) {
		case 0:
			return invalid("after", "duplicate %s", e.After[i])
		case 1:
			return invalid("after", "not sorted at %s", e.After[i])
		}
	}
	if !slices.Contains(Types(), e.Type) {
		return invalid("type", "unknown type %q", e.Type)
	}
	if len(e.Body) == 0 || e.Body[0] != '{' {
		return invalid("body", "not an object")
	}
	out, err := canon.Canonicalize(e.Body)
	if err != nil {
		return invalidWrap("body", err)
	}
	if !bytes.Equal(bytes.TrimSuffix(out, []byte("\n")), e.Body) {
		return invalidWrap("body", canon.ErrNotCanonical)
	}
	return nil
}

// Encode validates e and returns its canonical JSON, ending in one newline.
func Encode(e Entry) ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}
	after := e.After
	if after == nil {
		after = []id.ID{}
	}
	return canon.Marshal(map[string]any{
		"v":          e.V,
		"id":         e.ID,
		"change_set": e.ChangeSet,
		"at":         e.At.Format(timeLayout),
		"actor":      e.Actor,
		"for":        e.For,
		"after":      after,
		"type":       e.Type,
		"body":       e.Body,
	})
}

// Decode parses one canonical entry. A newer format fails with
// ErrNewerFormat whatever else the entry holds, since a newer format may
// add fields this one does not know; every other failure wraps ErrInvalid.
func Decode(data []byte) (Entry, error) {
	if err := canon.Check(data); err != nil {
		return Entry{}, fmt.Errorf("%w: %w", ErrInvalid, err)
	}
	if data[0] != '{' {
		return Entry{}, fmt.Errorf("%w: not an object", ErrInvalid)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return Entry{}, fmt.Errorf("%w: %w", ErrInvalid, err)
	}

	var e Entry
	raw, ok := fields["v"]
	if !ok {
		return Entry{}, invalid("v", "missing")
	}
	v, err := strconv.ParseInt(string(raw), 10, 64)
	if err != nil {
		return Entry{}, invalid("v", "not an integer")
	}
	if v > FormatVersion {
		return Entry{}, fmt.Errorf("%w: field %q: %d, this Genesis reads %d", ErrNewerFormat, "v", v, FormatVersion)
	}
	if v < 1 {
		return Entry{}, invalid("v", "%d, want 1 through %d", v, FormatVersion)
	}
	e.V = int(v)

	for _, k := range keys {
		if _, ok := fields[k]; !ok {
			return Entry{}, invalid(k, "missing")
		}
	}
	if len(fields) != len(keys) {
		names := make([]string, 0, len(fields))
		for k := range fields {
			if !slices.Contains(keys, k) {
				names = append(names, k)
			}
		}
		slices.Sort(names)
		return Entry{}, invalid(names[0], "unknown field")
	}

	for _, f := range []struct {
		name string
		dst  *id.ID
	}{{"id", &e.ID}, {"change_set", &e.ChangeSet}, {"actor", &e.Actor}, {"for", &e.For}} {
		s, err := str(fields[f.name])
		if err != nil {
			return Entry{}, invalidWrap(f.name, err)
		}
		if *f.dst, err = id.Parse(s); err != nil {
			return Entry{}, invalidWrap(f.name, err)
		}
	}

	at, err := str(fields["at"])
	if err != nil {
		return Entry{}, invalidWrap("at", err)
	}
	// time.Parse accepts fractional seconds the layout does not mention, so
	// only a value that reformats to itself is the one spelling.
	if e.At, err = time.Parse(timeLayout, at); err != nil || e.At.Format(timeLayout) != at {
		return Entry{}, invalid("at", "%q is not %s", at, timeLayout)
	}

	raw = fields["after"]
	if raw[0] != '[' {
		return Entry{}, invalid("after", "not an array")
	}
	var after []json.RawMessage
	if err := json.Unmarshal(raw, &after); err != nil {
		return Entry{}, invalidWrap("after", err)
	}
	e.After = make([]id.ID, 0, len(after))
	for _, r := range after {
		s, err := str(r)
		if err != nil {
			return Entry{}, invalidWrap("after", err)
		}
		i, err := id.Parse(s)
		if err != nil {
			return Entry{}, invalidWrap("after", err)
		}
		e.After = append(e.After, i)
	}

	typ, err := str(fields["type"])
	if err != nil {
		return Entry{}, invalidWrap("type", err)
	}
	e.Type = Type(typ)

	raw = fields["body"]
	if raw[0] != '{' {
		return Entry{}, invalid("body", "not an object")
	}
	// The body sits indented inside the entry; standalone it has its own
	// canonical form.
	body, err := canon.Canonicalize(raw)
	if err != nil {
		return Entry{}, invalidWrap("body", err)
	}
	e.Body = json.RawMessage(bytes.TrimSuffix(body, []byte("\n")))

	if err := e.Validate(); err != nil {
		return Entry{}, err
	}
	return e, nil
}

var errNotString = errors.New("not a string")

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
