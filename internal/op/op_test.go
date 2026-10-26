package op

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/key"
)

const (
	sOp   = "01JB0000000000000000000P01"
	sNode = "01JB0000000000000000000N0D"
	sPar  = "01JB0000000000000000000PAR"
	sBase = "01JB000000000000000000BASE"
	sLink = "01JB0000000000000000000NK1"
	sFrom = "01JB000000000000000000FR0M"
	sTo   = "01JB00000000000000000000T0"
	sDec  = "01JB0000000000000000000DEC"
	sZero = "00000000000000000000000000"
)

func mustID(t testing.TB, s string) id.ID {
	t.Helper()
	i, err := id.Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	return i
}

// golden holds one canonical op of each type, spelled by hand so the
// expected bytes do not depend on Encode.
var golden = []struct {
	typ Type
	doc string
}{
	{TypeNodeCreate, `{
  "fields": {
    "count": 3,
    "statement": "Résumé <&>",
    "tags": [
      "a",
      null
    ]
  },
  "handle": "token-refresh",
  "id": "` + sOp + `",
  "key": "i",
  "kind": "capability",
  "node": "` + sNode + `",
  "op": "node.create",
  "parent": "` + sPar + `"
}`},
	{TypeNodeUpdate, `{
  "base": "` + sBase + `",
  "handle": "renamed",
  "id": "` + sOp + `",
  "node": "` + sNode + `",
  "op": "node.update",
  "set": {
    "nested": {
      "x": true
    }
  },
  "unset": [
    "a_field",
    "b"
  ]
}`},
	{TypeNodeMove, `{
  "base": "` + sBase + `",
  "id": "` + sOp + `",
  "key": "a0i",
  "node": "` + sNode + `",
  "op": "node.move",
  "parent": null
}`},
	{TypeNodeRetire, `{
  "base": "` + sBase + `",
  "id": "` + sOp + `",
  "node": "` + sNode + `",
  "op": "node.retire",
  "reason": "merged into its sibling"
}`},
	{TypeLinkCreate, `{
  "decision": "` + sDec + `",
  "from": "` + sFrom + `",
  "id": "` + sOp + `",
  "link": "` + sLink + `",
  "note": "allowed for now",
  "op": "link.create",
  "to": "` + sTo + `",
  "type": "exempts"
}`},
	{TypeLinkRetire, `{
  "base": "` + sBase + `",
  "id": "` + sOp + `",
  "link": "` + sLink + `",
  "op": "link.retire"
}`},
}

func TestRoundTrip(t *testing.T) {
	if len(golden) != len(Types()) {
		t.Fatalf("%d golden ops, %d types", len(golden), len(Types()))
	}
	for i, g := range golden {
		if g.typ != Types()[i] {
			t.Fatalf("golden %d is %s, want %s", i, g.typ, Types()[i])
		}
		t.Run(string(g.typ), func(t *testing.T) {
			if err := canon.Check([]byte(g.doc + "\n")); err != nil {
				t.Fatalf("golden is not canonical: %v", err)
			}
			o, err := Decode(json.RawMessage(g.doc))
			if err != nil {
				t.Fatalf("Decode: %v", err)
			}
			if o.Type != g.typ {
				t.Fatalf("Type = %s, want %s", o.Type, g.typ)
			}
			out, err := Encode(o)
			if err != nil {
				t.Fatalf("Encode: %v", err)
			}
			if string(out) != g.doc {
				t.Fatalf("Encode(Decode(x)) =\n%s\nwant\n%s", out, g.doc)
			}
			again, err := Decode(out)
			if err != nil {
				t.Fatalf("re-Decode: %v", err)
			}
			if !reflect.DeepEqual(again, o) {
				t.Fatalf("re-Decode = %#v, want %#v", again, o)
			}
		})
	}
}

func TestRoundTripIndented(t *testing.T) {
	// As an op appears inside an entry body: indented and compacted forms
	// decode to the same op and encode to the standalone canonical form.
	for _, g := range golden {
		t.Run(string(g.typ), func(t *testing.T) {
			indented := strings.ReplaceAll(g.doc, "\n", "\n      ")
			var compact bytes.Buffer
			if err := json.Compact(&compact, []byte(g.doc)); err != nil {
				t.Fatal(err)
			}
			for _, in := range []string{indented, compact.String()} {
				o, err := Decode(json.RawMessage(in))
				if err != nil {
					t.Fatalf("Decode(%q): %v", in, err)
				}
				out, err := Encode(o)
				if err != nil {
					t.Fatal(err)
				}
				if string(out) != g.doc {
					t.Fatalf("Encode = %s, want %s", out, g.doc)
				}
			}
		})
	}
}

func TestDecodeValues(t *testing.T) {
	o, err := Decode(json.RawMessage(golden[0].doc))
	if err != nil {
		t.Fatal(err)
	}
	want := Op{
		ID:     mustID(t, sOp),
		Type:   TypeNodeCreate,
		Node:   mustID(t, sNode),
		Kind:   KindCapability,
		Handle: "token-refresh",
		Parent: mustID(t, sPar),
		Key:    "i",
		Fields: map[string]json.RawMessage{
			"count":     json.RawMessage(`3`),
			"statement": json.RawMessage(`"Résumé <&>"`),
			"tags":      json.RawMessage("[\n  \"a\",\n  null\n]"),
		},
	}
	if !reflect.DeepEqual(o, want) {
		t.Fatalf("Decode = %#v\nwant %#v", o, want)
	}
	if o.Target() != mustID(t, sNode) {
		t.Fatalf("Target = %s", o.Target())
	}
	l, err := Decode(json.RawMessage(golden[5].doc))
	if err != nil {
		t.Fatal(err)
	}
	if l.Target() != mustID(t, sLink) {
		t.Fatalf("link Target = %s", l.Target())
	}
}

func TestEmptyCollectionsDecodeToNil(t *testing.T) {
	in := `{"id":"` + sOp + `","op":"node.create","node":"` + sNode + `","kind":"decision","handle":"d","parent":null,"key":null,"fields":{}}`
	o, err := Decode(json.RawMessage(in))
	if err != nil {
		t.Fatal(err)
	}
	if o.Fields != nil {
		t.Fatalf("Fields = %#v, want nil", o.Fields)
	}
	upd := `{"id":"` + sOp + `","op":"node.update","node":"` + sNode + `","base":"` + sBase + `","handle":"h","set":{},"unset":[]}`
	u, err := Decode(json.RawMessage(upd))
	if err != nil {
		t.Fatal(err)
	}
	if u.Set != nil || u.Unset != nil {
		t.Fatalf("Set, Unset = %#v, %#v, want nil", u.Set, u.Unset)
	}
	// Encode spells nil collections as {} and [], and empty non-nil ones the same.
	u.Set, u.Unset = map[string]json.RawMessage{}, []string{}
	out, err := Encode(u)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"set": {}`) || !strings.Contains(string(out), `"unset": []`) {
		t.Fatalf("Encode = %s", out)
	}
}

func TestEnumerations(t *testing.T) {
	if got := Types(); !reflect.DeepEqual(got, []Type{"node.create", "node.update", "node.move", "node.retire", "link.create", "link.retire"}) {
		t.Fatalf("Types = %v", got)
	}
	if got := Kinds(); !reflect.DeepEqual(got, []Kind{"intent", "capability", "constraint", "decision", "question", "evidence", "source", "actor", "anchor"}) {
		t.Fatalf("Kinds = %v", got)
	}
	if got := LinkTypes(); !reflect.DeepEqual(got, []LinkType{"serves", "constrains", "exempts", "depends-on", "relates-to", "decided-by", "answers", "blocks", "supersedes", "derived-from", "discovered-from", "implements", "verifies", "observes"}) {
		t.Fatalf("LinkTypes = %v", got)
	}
	for _, k := range Kinds() {
		if want := k == KindCapability || k == KindIntent; HasParent(k) != want {
			t.Errorf("HasParent(%s) = %v", k, !want)
		}
	}
	if HasParent("nope") {
		t.Error("HasParent of an unknown kind")
	}
}

// base returns a valid op of type typ as a generic JSON object to mutate.
func base(t testing.TB, typ Type) map[string]any {
	t.Helper()
	for _, g := range golden {
		if g.typ == typ {
			var m map[string]any
			if err := json.Unmarshal([]byte(g.doc), &m); err != nil {
				t.Fatal(err)
			}
			return m
		}
	}
	t.Fatalf("no golden %s", typ)
	return nil
}

type mutation struct {
	name  string
	typ   Type
	edit  func(m map[string]any)
	field string
	also  error // another sentinel the error must wrap
}

func del(k string) func(map[string]any) {
	return func(m map[string]any) { delete(m, k) }
}

func set(kv ...any) func(map[string]any) {
	return func(m map[string]any) {
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
}

// shapeCases: missing fields, unknown fields, wrong kinds and zero IDs.
var shapeCases = []mutation{
	{"missing op", TypeNodeCreate, del("op"), "op", nil},
	{"op not a string", TypeNodeCreate, set("op", 1), "op", nil},
	{"op null", TypeNodeCreate, set("op", nil), "op", nil},

	// Missing fields, reported in schema order.
	{"create missing id", TypeNodeCreate, del("id"), "id", nil},
	{"create missing kind and fields", TypeNodeCreate, func(m map[string]any) { delete(m, "fields"); delete(m, "kind") }, "kind", nil},
	{"create missing fields", TypeNodeCreate, del("fields"), "fields", nil},
	{"update missing unset", TypeNodeUpdate, del("unset"), "unset", nil},
	{"update missing base", TypeNodeUpdate, del("base"), "base", nil},
	{"move missing parent", TypeNodeMove, del("parent"), "parent", nil},
	{"retire missing reason", TypeNodeRetire, del("reason"), "reason", nil},
	{"link create missing decision", TypeLinkCreate, del("decision"), "decision", nil},
	{"link create missing type", TypeLinkCreate, del("type"), "type", nil},
	{"link retire missing link", TypeLinkRetire, del("link"), "link", nil},
	{"missing beats unknown", TypeLinkRetire, func(m map[string]any) { delete(m, "base"); m["aaa"] = 1 }, "base", nil},

	// Unknown fields: the first in sorted order.
	{"create unknown", TypeNodeCreate, set("zzz", 1, "extra", 2), "extra", nil},
	{"update unknown", TypeNodeUpdate, set("parent", nil), "parent", nil},
	{"move unknown", TypeNodeMove, set("handle", "x"), "handle", nil},
	{"retire unknown", TypeNodeRetire, set("why", "x"), "why", nil},
	{"link create unknown", TypeLinkCreate, set("base", sBase), "base", nil},
	{"link retire unknown", TypeLinkRetire, set("node", sNode), "node", nil},

	// Wrong kinds.
	{"id a number", TypeNodeCreate, set("id", 7), "id", nil},
	{"id lower case", TypeNodeCreate, set("id", strings.ToLower(sOp)), "id", id.ErrInvalid},
	{"node an array", TypeNodeCreate, set("node", []any{}), "node", nil},
	{"kind a number", TypeNodeCreate, set("kind", 1), "kind", nil},
	{"handle null on create", TypeNodeCreate, set("handle", nil), "handle", nil},
	{"parent a bool", TypeNodeCreate, set("parent", true), "parent", nil},
	{"key a number", TypeNodeCreate, set("key", 5), "key", nil},
	{"fields an array", TypeNodeCreate, set("fields", []any{}), "fields", nil},
	{"fields null", TypeNodeCreate, set("fields", nil), "fields", nil},
	{"base null", TypeNodeUpdate, set("base", nil), "base", nil},
	{"update handle a number", TypeNodeUpdate, set("handle", 1), "handle", nil},
	{"set a string", TypeNodeUpdate, set("set", "x"), "set", nil},
	{"unset an object", TypeNodeUpdate, set("unset", map[string]any{}), "unset", nil},
	{"unset of numbers", TypeNodeUpdate, set("unset", []any{1}), "unset", nil},
	{"move key null", TypeNodeMove, set("key", nil), "key", nil},
	{"move parent a number", TypeNodeMove, set("parent", 1), "parent", nil},
	{"reason null", TypeNodeRetire, set("reason", nil), "reason", nil},
	{"link type a number", TypeLinkCreate, set("type", 1), "type", nil},
	{"from null", TypeLinkCreate, set("from", nil), "from", nil},
	{"note null", TypeLinkCreate, set("note", nil), "note", nil},
	{"decision a number", TypeLinkCreate, set("decision", 9), "decision", nil},
	{"link retire base a number", TypeLinkRetire, set("base", 9), "base", nil},

	// Zero IDs.
	{"zero id", TypeNodeCreate, set("id", sZero), "id", nil},
	{"zero node", TypeNodeCreate, set("node", sZero), "node", nil},
	{"zero parent", TypeNodeCreate, set("parent", sZero), "parent", nil},
	{"zero update base", TypeNodeUpdate, set("base", sZero), "base", nil},
	{"zero move node", TypeNodeMove, set("node", sZero), "node", nil},
	{"zero move parent", TypeNodeMove, set("parent", sZero), "parent", nil},
	{"zero retire base", TypeNodeRetire, set("base", sZero), "base", nil},
	{"zero link", TypeLinkCreate, set("link", sZero), "link", nil},
	{"zero from", TypeLinkCreate, set("from", sZero), "from", nil},
	{"zero to", TypeLinkCreate, set("to", sZero), "to", nil},
	{"zero decision", TypeLinkCreate, set("decision", sZero), "decision", nil},
	{"zero link retire base", TypeLinkRetire, set("base", sZero), "base", nil},
}

// ruleCases: well-shaped ops that break a rule.
var ruleCases = []mutation{
	{"unknown op type", TypeNodeCreate, set("op", "node.delete"), "op", nil},
	{"unknown kind", TypeNodeCreate, set("kind", "epic"), "kind", nil},
	{"unknown link type", TypeLinkCreate, set("type", "loves"), "type", nil},

	{"non-tree kind with a parent", TypeNodeCreate, set("kind", "decision", "key", nil), "parent", nil},
	{"non-tree kind with a key", TypeNodeCreate, set("kind", "evidence", "parent", nil), "key", nil},
	{"tree kind with a null key", TypeNodeCreate, set("key", nil), "key", key.ErrInvalid},
	{"tree kind with a bad key", TypeNodeCreate, set("key", "a0"), "key", key.ErrInvalid},
	{"tree kind with an empty key", TypeNodeCreate, set("kind", "intent", "key", ""), "key", key.ErrInvalid},
	{"tree kind with an upper-case key", TypeNodeCreate, set("key", "A"), "key", key.ErrInvalid},
	{"move with a bad key", TypeNodeMove, set("key", "x0"), "key", key.ErrInvalid},

	{"empty handle", TypeNodeCreate, set("handle", ""), "handle", nil},
	{"handle with upper case", TypeNodeCreate, set("handle", "Token"), "handle", nil},
	{"handle with a slash", TypeNodeCreate, set("handle", "a/b"), "handle", nil},
	{"handle with a leading dash", TypeNodeCreate, set("handle", "-a"), "handle", nil},
	{"handle with a double dash", TypeNodeCreate, set("handle", "a--b"), "handle", nil},
	{"handle too long", TypeNodeCreate, set("handle", strings.Repeat("a", 65)), "handle", nil},
	{"update with an empty handle", TypeNodeUpdate, set("handle", ""), "handle", nil},
	{"update with a bad handle", TypeNodeUpdate, set("handle", "a_b"), "handle", nil},

	{"field name with upper case", TypeNodeCreate, set("fields", map[string]any{"Name": 1}), "fields", nil},
	{"field name starting with a digit", TypeNodeCreate, set("fields", map[string]any{"1a": 1}), "fields", nil},
	{"field name empty", TypeNodeCreate, set("fields", map[string]any{"": 1}), "fields", nil},
	{"field name too long", TypeNodeCreate, set("fields", map[string]any{"a" + strings.Repeat("b", 64): 1}), "fields", nil},
	{"set name with a dash", TypeNodeUpdate, set("set", map[string]any{"a-b": 1}), "set", nil},
	{"unset name with a dash", TypeNodeUpdate, set("unset", []any{"a-b"}), "unset", nil},

	{"update with nothing to change", TypeNodeUpdate, set("handle", nil, "set", map[string]any{}, "unset", []any{}), "set", nil},
	{"name in both set and unset", TypeNodeUpdate, set("set", map[string]any{"b": 1}), "unset", nil},
	{"unsorted unset", TypeNodeUpdate, set("unset", []any{"b", "a"}), "unset", nil},
	{"duplicate in unset", TypeNodeUpdate, set("unset", []any{"a", "a"}), "unset", nil},

	{"empty retire reason", TypeNodeRetire, set("reason", ""), "reason", nil},

	{"relates-to without a note", TypeLinkCreate, set("type", "relates-to", "decision", nil, "note", ""), "note", nil},
	{"exempts without a note", TypeLinkCreate, set("note", ""), "note", nil},
	{"exempts without a decision", TypeLinkCreate, set("decision", nil), "decision", nil},
	{"decision on a non-exempts link", TypeLinkCreate, set("type", "serves"), "decision", nil},
	{"from equals to", TypeLinkCreate, set("to", sFrom), "to", nil},
}

func runMutations(t *testing.T, cases []mutation) {
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := base(t, tc.typ)
			tc.edit(m)
			in, err := json.Marshal(m)
			if err != nil {
				t.Fatal(err)
			}
			_, err = Decode(in)
			checkErr(t, err, tc.field, tc.also)
		})
	}
}

func checkErr(t *testing.T, err error, field string, also error) {
	t.Helper()
	if err == nil {
		t.Fatalf("no error, want one on field %q", field)
	}
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("error does not wrap ErrInvalid: %v", err)
	}
	if also != nil && !errors.Is(err, also) {
		t.Fatalf("error does not wrap %v: %v", also, err)
	}
	if want := `field "` + field + `"`; !strings.Contains(err.Error(), want) {
		t.Fatalf("error %q does not name %s", err, want)
	}
}

func TestDecodeShapeInvalid(t *testing.T) { runMutations(t, shapeCases) }

func TestDecodeRuleInvalid(t *testing.T) { runMutations(t, ruleCases) }

func TestDecodeSyntaxInvalid(t *testing.T) {
	for _, in := range []string{
		``,
		`[]`,
		`"op"`,
		`null`,
		`{`,
		`{"op":"node.retire","op":"node.retire"}`,
		`{"op":1.5}`,
		`{} {}`,
	} {
		t.Run(in, func(t *testing.T) {
			if _, err := Decode(json.RawMessage(in)); !errors.Is(err, ErrInvalid) {
				t.Fatalf("Decode(%q) = %v, want ErrInvalid", in, err)
			}
		})
	}
}

func validOp(t testing.TB, typ Type) Op {
	t.Helper()
	for _, g := range golden {
		if g.typ == typ {
			o, err := Decode(json.RawMessage(g.doc))
			if err != nil {
				t.Fatal(err)
			}
			return o
		}
	}
	t.Fatalf("no golden %s", typ)
	return Op{}
}

func TestValidateInvalid(t *testing.T) {
	cases := []struct {
		name  string
		typ   Type
		edit  func(o *Op)
		field string
	}{
		{"zero id", TypeNodeRetire, func(o *Op) { o.ID = id.ID{} }, "id"},
		{"unknown type", TypeNodeRetire, func(o *Op) { o.Type = "node.delete" }, "op"},
		{"zero node", TypeNodeRetire, func(o *Op) { o.Node = id.ID{} }, "node"},
		{"zero link", TypeLinkRetire, func(o *Op) { o.Link = id.ID{} }, "link"},
		{"zero base", TypeNodeMove, func(o *Op) { o.Base = id.ID{} }, "base"},
		{"link on a node op", TypeNodeRetire, func(o *Op) { o.Link = o.Node }, "link"},
		{"base on a create", TypeNodeCreate, func(o *Op) { o.Base = o.Node }, "base"},
		{"fields on an update", TypeNodeUpdate, func(o *Op) { o.Fields = map[string]json.RawMessage{"a": json.RawMessage("1")} }, "fields"},
		{"empty fields map on an update", TypeNodeUpdate, func(o *Op) { o.Fields = map[string]json.RawMessage{} }, "fields"},
		{"reason on a move", TypeNodeMove, func(o *Op) { o.Reason = "x" }, "reason"},
		{"link type on a retire", TypeLinkRetire, func(o *Op) { o.LinkType = LinkServes }, "type"},
		{"note on a node create", TypeNodeCreate, func(o *Op) { o.Note = "x" }, "note"},
		{"kind on a move", TypeNodeMove, func(o *Op) { o.Kind = KindIntent }, "kind"},
		{"non-canonical field value", TypeNodeCreate, func(o *Op) { o.Fields["count"] = json.RawMessage(" 3") }, "fields"},
		{"empty field value", TypeNodeCreate, func(o *Op) { o.Fields["count"] = nil }, "fields"},
		{"invalid field value", TypeNodeCreate, func(o *Op) { o.Fields["count"] = json.RawMessage("1.5") }, "fields"},
		{"non-canonical set value", TypeNodeUpdate, func(o *Op) { o.Set["nested"] = json.RawMessage(`{"x":true}`) }, "set"},
		{"invalid UTF-8 reason", TypeNodeRetire, func(o *Op) { o.Reason = "\xff" }, "reason"},
		{"invalid UTF-8 note", TypeLinkCreate, func(o *Op) { o.Note = "\xff" }, "note"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := validOp(t, tc.typ)
			tc.edit(&o)
			checkErr(t, o.Validate(), tc.field, nil)
			_, err := Encode(o)
			checkErr(t, err, tc.field, nil)
		})
	}
}

func TestTargetOfUnknownType(t *testing.T) {
	if got := (Op{Type: "x", Node: mustID(t, sNode)}).Target(); !got.IsZero() {
		t.Fatalf("Target = %s", got)
	}
}

func FuzzDecode(f *testing.F) {
	for _, g := range golden {
		f.Add([]byte(g.doc))
	}
	for _, cases := range [][]mutation{shapeCases, ruleCases} {
		for _, tc := range cases {
			m := base(f, tc.typ)
			tc.edit(m)
			in, err := json.Marshal(m)
			if err != nil {
				f.Fatal(err)
			}
			f.Add(in)
		}
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		o, err := Decode(data)
		if err != nil {
			if !errors.Is(err, ErrInvalid) {
				t.Fatalf("Decode error does not wrap ErrInvalid: %v", err)
			}
			return
		}
		out, err := Encode(o)
		if err != nil {
			t.Fatalf("Encode(Decode(%q)) failed: %v", data, err)
		}
		want, err := canon.Canonicalize(data)
		if err != nil {
			t.Fatalf("Decode accepted non-JSON %q", data)
		}
		if !bytes.Equal(out, bytes.TrimSuffix(want, []byte("\n"))) {
			t.Fatalf("Encode(Decode(%q)) = %q, want %q", data, out, want)
		}
		again, err := Decode(out)
		if err != nil {
			t.Fatalf("Decode(Encode(...)) failed: %v", err)
		}
		if !reflect.DeepEqual(again, o) {
			t.Fatalf("re-Decode = %#v, want %#v", again, o)
		}
		trusted, err := DecodeCanonical(indented(string(out), 2))
		if err != nil {
			t.Fatalf("DecodeCanonical(Encode(...)) failed: %v", err)
		}
		if !reflect.DeepEqual(trusted, o) {
			t.Fatalf("DecodeCanonical = %#v, want %#v", trusted, o)
		}
		if o.Target().IsZero() {
			t.Fatalf("decoded op has a zero target: %#v", o)
		}
	})
}

// indented returns doc as it appears inside an entry body: every line but
// the first starts with depth levels of two-space indentation.
func indented(doc string, depth int) json.RawMessage {
	return json.RawMessage(strings.ReplaceAll(doc, "\n", "\n"+strings.Repeat("  ", depth)))
}

// canonicalDoc returns m's canonical JSON without the final newline.
func canonicalDoc(t testing.TB, m map[string]any) string {
	t.Helper()
	out, err := canon.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSuffix(string(out), "\n")
}

// TestDecodeCanonicalMatchesDecode checks that DecodeCanonical gives
// Decode's op for every golden op at every indentation an entry body puts
// it at, and Decode's error, word for word, for every shape and rule
// mutation.
func TestDecodeCanonicalMatchesDecode(t *testing.T) {
	for _, g := range golden {
		want, err := Decode(json.RawMessage(g.doc))
		if err != nil {
			t.Fatalf("%s: %v", g.typ, err)
		}
		for depth := 0; depth <= 3; depth++ {
			in := indented(g.doc, depth)
			got, err := DecodeCanonical(in)
			if err != nil {
				t.Fatalf("%s at depth %d: %v", g.typ, depth, err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("%s at depth %d: DecodeCanonical =\n%#v\nDecode =\n%#v", g.typ, depth, got, want)
			}
		}
	}
	for _, cases := range [][]mutation{shapeCases, ruleCases} {
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				m := base(t, tc.typ)
				tc.edit(m)
				in := indented(canonicalDoc(t, m), 2)
				_, wantErr := Decode(in)
				if wantErr == nil {
					t.Fatal("Decode accepted the mutation")
				}
				_, gotErr := DecodeCanonical(in)
				if gotErr == nil {
					t.Fatalf("DecodeCanonical accepted %s", in)
				}
				if gotErr.Error() != wantErr.Error() {
					t.Fatalf("DecodeCanonical error %q, Decode error %q", gotErr, wantErr)
				}
			})
		}
	}
}

func TestDecodeCanonicalSyntaxInvalid(t *testing.T) {
	for _, in := range []string{``, `[]`, `"op"`, `null`, `{`, `{} {}`} {
		t.Run(in, func(t *testing.T) {
			if _, err := DecodeCanonical(json.RawMessage(in)); !errors.Is(err, ErrInvalid) {
				t.Fatalf("DecodeCanonical(%q) = %v, want ErrInvalid", in, err)
			}
		})
	}
}

func TestDeindent(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"{\n    \"a\": 1,\n    \"b\": [\n      2\n    ]\n  }", "{\n  \"a\": 1,\n  \"b\": [\n    2\n  ]\n}"},
		{"[\n  1\n]", "[\n  1\n]"},
		{"{}", "{}"},
		{"3", "3"},
		{`"a\nb"`, `"a\nb"`},
		{"{\n\"a\": 1\n}", "{\n\"a\": 1\n}"},
	} {
		if got := string(deindent([]byte(tc.in))); got != tc.want {
			t.Errorf("deindent(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// benchDecode decodes the golden ops as they sit in an entry body.
func benchDecode(b *testing.B, decode func(json.RawMessage) (Op, error)) {
	docs := make([]json.RawMessage, len(golden))
	for i, g := range golden {
		docs[i] = indented(g.doc, 2)
	}
	b.ReportAllocs()
	for b.Loop() {
		for _, d := range docs {
			if _, err := decode(d); err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkDecode(b *testing.B)          { benchDecode(b, Decode) }
func BenchmarkDecodeCanonical(b *testing.B) { benchDecode(b, DecodeCanonical) }
