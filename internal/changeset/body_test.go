package changeset

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

const (
	sSrc1 = "01JB000000000000000000SRC1"
	sSrc2 = "01JB000000000000000000SRC2"
	sOp1  = "01JB00000000000000000000P1"
	sOp2  = "01JB00000000000000000000P2"
	sOp3  = "01JB00000000000000000000P3"
	sOp4  = "01JB00000000000000000000P4"
	sNode = "01JB0000000000000000000N0D"
	sBase = "01JB000000000000000000BASE"
	sLink = "01JB0000000000000000000NK1"
	sBy   = "01JB0000000000000000000CS2"
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

// golden holds one canonical body of each entry type, spelled by hand so the
// expected bytes do not depend on EncodeBody, with the value it decodes to.
func golden(t testing.TB) []struct {
	typ  entry.Type
	doc  string
	want Body
} {
	return []struct {
		typ  entry.Type
		doc  string
		want Body
	}{
		{entry.TypeOpen, `{
  "intent": "Refresh tokens <&> before they expire",
  "kind": "intention",
  "sources": [
    "` + sSrc1 + `",
    "` + sSrc2 + `"
  ]
}`, Open{Kind: KindIntention, Intent: "Refresh tokens <&> before they expire", Sources: []id.ID{mustID(t, sSrc1), mustID(t, sSrc2)}}},
		{entry.TypeOps, `{
  "drop": [
    "` + sOp3 + `"
  ],
  "ops": [
    {
      "fields": {
        "statement": "Résumé"
      },
      "handle": "token-refresh",
      "id": "` + sOp1 + `",
      "key": "i",
      "kind": "capability",
      "node": "` + sNode + `",
      "op": "node.create",
      "parent": null
    },
    {
      "base": "` + sBase + `",
      "id": "` + sOp2 + `",
      "link": "` + sLink + `",
      "op": "link.retire"
    }
  ]
}`, Ops{
			Ops: []op.Op{
				{ID: mustID(t, sOp1), Type: op.TypeNodeCreate, Node: mustID(t, sNode), Kind: op.KindCapability,
					Handle: "token-refresh", Key: "i", Fields: map[string]json.RawMessage{"statement": json.RawMessage(`"Résumé"`)}},
				{ID: mustID(t, sOp2), Type: op.TypeLinkRetire, Link: mustID(t, sLink), Base: mustID(t, sBase)},
			},
			Drop: []id.ID{mustID(t, sOp3)},
		}},
		{entry.TypePropose, `{}`, Propose{}},
		{entry.TypeAccept, `{
  "reason": "Meets the bar",
  "rule": "owner-approves"
}`, Accept{Rule: "owner-approves", Reason: "Meets the bar"}},
		{entry.TypeReject, `{
  "reason": "",
  "rule": "owner-approves"
}`, Reject{Rule: "owner-approves"}},
		{entry.TypeWithdraw, `{}`, Withdraw{}},
		{entry.TypeSupersede, `{
  "by": "` + sBy + `"
}`, Supersede{By: mustID(t, sBy)}},
		{entry.TypeLand, `{}`, Land{}},
		{entry.TypeReopen, `{}`, Reopen{}},
		{entry.TypeNote, `{
  "text": "Line one\nline \"two\""
}`, Note{Text: "Line one\nline \"two\""}},
	}
}

func TestBodyRoundTrip(t *testing.T) {
	gs := golden(t)
	if len(gs) != len(entry.Types()) {
		t.Fatalf("%d golden bodies, want one per entry type (%d)", len(gs), len(entry.Types()))
	}
	for i, g := range gs {
		t.Run(string(g.typ), func(t *testing.T) {
			if g.typ != entry.Types()[i] {
				t.Fatalf("golden %d is %s, want %s", i, g.typ, entry.Types()[i])
			}
			if err := canon.Check([]byte(g.doc + "\n")); err != nil {
				t.Fatalf("golden is not canonical: %v", err)
			}
			b, err := DecodeBody(entry.Entry{V: EntryVersion(g.want), Type: g.typ, Body: json.RawMessage(g.doc)})
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(b, g.want) {
				t.Fatalf("DecodeBody = %#v, want %#v", b, g.want)
			}
			if b.Type() != g.typ {
				t.Fatalf("Type() = %s, want %s", b.Type(), g.typ)
			}
			out, err := EncodeBody(b)
			if err != nil {
				t.Fatal(err)
			}
			if string(out) != g.doc {
				t.Fatalf("EncodeBody(DecodeBody) =\n%s\nwant\n%s", out, g.doc)
			}
			out, err = EncodeBody(g.want)
			if err != nil {
				t.Fatal(err)
			}
			if string(out) != g.doc {
				t.Fatalf("EncodeBody(want) =\n%s\nwant\n%s", out, g.doc)
			}
		})
	}
}

func TestEmptyListsDecodeToNil(t *testing.T) {
	b, err := DecodeBody(entry.Entry{Type: entry.TypeOpen, Body: json.RawMessage(`{
  "intent": "x",
  "kind": "record",
  "sources": []
}`)})
	if err != nil {
		t.Fatal(err)
	}
	if o := b.(Open); o.Sources != nil {
		t.Fatalf("Sources = %#v, want nil", o.Sources)
	}

	doc := `{
  "drop": [
    "` + sOp1 + `"
  ],
  "ops": []
}`
	b, err = DecodeBody(entry.Entry{Type: entry.TypeOps, Body: json.RawMessage(doc)})
	if err != nil {
		t.Fatal(err)
	}
	if o := b.(Ops); o.Ops != nil {
		t.Fatalf("Ops = %#v, want nil", o.Ops)
	}
	// Nil and empty encode alike.
	for _, ops := range []Ops{{Drop: []id.ID{mustID(t, sOp1)}}, {Ops: []op.Op{}, Drop: []id.ID{mustID(t, sOp1)}}} {
		out, err := EncodeBody(ops)
		if err != nil {
			t.Fatal(err)
		}
		if string(out) != doc {
			t.Fatalf("EncodeBody = %s, want %s", out, doc)
		}
	}

	b, err = DecodeBody(entry.Entry{Type: entry.TypeOps, Body: json.RawMessage(`{
  "drop": [],
  "ops": [
    {
      "base": "` + sBase + `",
      "id": "` + sOp2 + `",
      "link": "` + sLink + `",
      "op": "link.retire"
    }
  ]
}`)})
	if err != nil {
		t.Fatal(err)
	}
	if o := b.(Ops); o.Drop != nil {
		t.Fatalf("Drop = %#v, want nil", o.Drop)
	}

	out, err := EncodeBody(Open{Kind: KindRecord, Intent: "x"})
	if err != nil {
		t.Fatal(err)
	}
	if want := "{\n  \"intent\": \"x\",\n  \"kind\": \"record\",\n  \"sources\": []\n}"; string(out) != want {
		t.Fatalf("EncodeBody(Open) = %s, want %s", out, want)
	}
}

func TestKinds(t *testing.T) {
	if got := Kinds(); !reflect.DeepEqual(got, []Kind{KindIntention, KindRecord}) {
		t.Fatalf("Kinds = %v", got)
	}
}

// base returns the golden body of typ as a generic JSON object to mutate.
func base(t testing.TB, typ entry.Type) map[string]any {
	t.Helper()
	for _, g := range golden(t) {
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
	typ   entry.Type
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

// setOp sets fields on the i-th element of an ops body's ops.
func setOp(i int, kv ...any) func(map[string]any) {
	return func(m map[string]any) {
		set(kv...)(m["ops"].([]any)[i].(map[string]any))
	}
}

// shapeCases: missing fields, unknown fields and wrong kinds, per type.
var shapeCases = []mutation{
	{"open missing kind", entry.TypeOpen, del("kind"), "kind", nil},
	{"open missing intent", entry.TypeOpen, del("intent"), "intent", nil},
	{"open missing sources", entry.TypeOpen, del("sources"), "sources", nil},
	{"open missing kind and sources", entry.TypeOpen, func(m map[string]any) { delete(m, "sources"); delete(m, "kind") }, "kind", nil},
	{"open unknown", entry.TypeOpen, set("zzz", 1, "title", "x"), "title", nil},
	{"open missing beats unknown", entry.TypeOpen, func(m map[string]any) { delete(m, "sources"); m["aaa"] = 1 }, "sources", nil},
	{"open kind a number", entry.TypeOpen, set("kind", 1), "kind", nil},
	{"open intent null", entry.TypeOpen, set("intent", nil), "intent", nil},
	{"open sources a string", entry.TypeOpen, set("sources", sSrc1), "sources", nil},
	{"open sources null", entry.TypeOpen, set("sources", nil), "sources", nil},
	{"open sources of numbers", entry.TypeOpen, set("sources", []any{1}), "sources", nil},
	{"open source lower case", entry.TypeOpen, set("sources", []any{strings.ToLower(sSrc1)}), "sources", id.ErrInvalid},

	{"ops missing ops", entry.TypeOps, del("ops"), "ops", nil},
	{"ops missing drop", entry.TypeOps, del("drop"), "drop", nil},
	{"ops unknown", entry.TypeOps, set("revision", 2), "revision", nil},
	{"ops ops an object", entry.TypeOps, set("ops", map[string]any{}), "ops", nil},
	{"ops ops null", entry.TypeOps, set("ops", nil), "ops", nil},
	{"ops op a number", entry.TypeOps, set("ops", []any{1}), "ops[0]", op.ErrInvalid},
	{"ops drop a string", entry.TypeOps, set("drop", sOp3), "drop", nil},
	{"ops drop of numbers", entry.TypeOps, set("drop", []any{3}), "drop", nil},

	{"propose unknown", entry.TypePropose, set("reason", "x"), "reason", nil},
	{"withdraw unknown", entry.TypeWithdraw, set("reason", "x"), "reason", nil},
	{"land unknown", entry.TypeLand, set("rule", "x"), "rule", nil},
	{"land unknown sorted", entry.TypeLand, set("z", 1, "a", 2), "a", nil},

	{"accept missing rule", entry.TypeAccept, del("rule"), "rule", nil},
	{"accept missing reason", entry.TypeAccept, del("reason"), "reason", nil},
	{"accept unknown", entry.TypeAccept, set("by", sBy), "by", nil},
	{"accept rule a number", entry.TypeAccept, set("rule", 1), "rule", nil},
	{"accept reason null", entry.TypeAccept, set("reason", nil), "reason", nil},

	{"reject missing rule", entry.TypeReject, del("rule"), "rule", nil},
	{"reject missing reason", entry.TypeReject, del("reason"), "reason", nil},
	{"reject unknown", entry.TypeReject, set("text", "x"), "text", nil},
	{"reject rule null", entry.TypeReject, set("rule", nil), "rule", nil},
	{"reject reason a bool", entry.TypeReject, set("reason", false), "reason", nil},

	{"supersede missing by", entry.TypeSupersede, del("by"), "by", nil},
	{"supersede unknown", entry.TypeSupersede, set("reason", "x"), "reason", nil},
	{"supersede by a number", entry.TypeSupersede, set("by", 1), "by", nil},
	{"supersede by null", entry.TypeSupersede, set("by", nil), "by", nil},
	{"supersede by not a ULID", entry.TypeSupersede, set("by", "nope"), "by", id.ErrInvalid},

	{"note missing text", entry.TypeNote, del("text"), "text", nil},
	{"note unknown target", entry.TypeNote, set("target", sNode), "target", nil},
	{"note text an array", entry.TypeNote, set("text", []any{}), "text", nil},
}

// ruleCases: well-shaped bodies that break a rule.
var ruleCases = []mutation{
	{"empty intent", entry.TypeOpen, set("intent", ""), "intent", nil},
	{"bad kind", entry.TypeOpen, set("kind", "epic"), "kind", nil},
	{"unsorted sources", entry.TypeOpen, set("sources", []any{sSrc2, sSrc1}), "sources", nil},
	{"duplicate sources", entry.TypeOpen, set("sources", []any{sSrc1, sSrc1}), "sources", nil},
	{"zero source", entry.TypeOpen, set("sources", []any{sZero}), "sources", nil},

	{"ops and drop both empty", entry.TypeOps, set("ops", []any{}, "drop", []any{}), "ops", nil},
	{"duplicate op ID in ops", entry.TypeOps, setOp(1, "id", sOp1), "ops", nil},
	{"ID in both ops and drop", entry.TypeOps, set("drop", []any{sOp2}), "drop", nil},
	{"unsorted drop", entry.TypeOps, set("drop", []any{sOp4, sOp3}), "drop", nil},
	{"duplicate drop", entry.TypeOps, set("drop", []any{sOp3, sOp3}), "drop", nil},
	{"zero drop", entry.TypeOps, set("drop", []any{sZero}), "drop", nil},
	{"invalid second op", entry.TypeOps, setOp(1, "base", nil), "ops[1]", op.ErrInvalid},
	{"second op of unknown type", entry.TypeOps, setOp(1, "op", "node.delete"), "ops[1]", op.ErrInvalid},

	{"accept empty rule", entry.TypeAccept, set("rule", ""), "rule", nil},
	{"reject empty rule", entry.TypeReject, set("rule", ""), "rule", nil},
	{"zero by", entry.TypeSupersede, set("by", sZero), "by", nil},
	{"empty note text", entry.TypeNote, set("text", ""), "text", nil},

	{"non-empty propose body", entry.TypePropose, set("note", "x"), "note", nil},
	{"non-empty withdraw body", entry.TypeWithdraw, set("why", nil), "why", nil},
	{"non-empty land body", entry.TypeLand, set("at", 1), "at", nil},
}

// canonical returns m as a canonical body, as an entry carries it.
func canonical(t testing.TB, m map[string]any) json.RawMessage {
	t.Helper()
	out, err := canon.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return json.RawMessage(bytes.TrimSuffix(out, []byte("\n")))
}

func runMutations(t *testing.T, cases []mutation) {
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := base(t, tc.typ)
			tc.edit(m)
			_, err := DecodeBody(entry.Entry{Type: tc.typ, Body: canonical(t, m)})
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

func TestBodyShapeInvalid(t *testing.T) { runMutations(t, shapeCases) }

func TestBodyRuleInvalid(t *testing.T) { runMutations(t, ruleCases) }

func TestBodyUnknownTypeInvalid(t *testing.T) {
	for _, typ := range []entry.Type{"merge", "", "Open"} {
		_, err := DecodeBody(entry.Entry{Type: typ, Body: json.RawMessage(`{}`)})
		checkErr(t, err, "type", nil)
	}
}

func TestBodySyntaxInvalid(t *testing.T) {
	for _, in := range []string{
		``,
		`[]`,
		`null`,
		`"text"`,
		`{`,
		`{"text":"x"}`,
		"{\n  \"text\": \"x\"\n}\n",
		"{\n  \"text\": \"x\",\n  \"text\": \"y\"\n}",
		`{} {}`,
	} {
		t.Run(in, func(t *testing.T) {
			if _, err := DecodeBody(entry.Entry{Type: entry.TypeNote, Body: json.RawMessage(in)}); !errors.Is(err, ErrInvalid) {
				t.Fatalf("DecodeBody(%q) = %v, want ErrInvalid", in, err)
			}
		})
	}
}

func TestEncodeBodyInvalid(t *testing.T) {
	good := op.Op{ID: mustID(t, sOp2), Type: op.TypeLinkRetire, Link: mustID(t, sLink), Base: mustID(t, sBase)}
	bad := good
	bad.Base = id.ID{}
	cases := []struct {
		name  string
		body  Body
		field string
		also  error
	}{
		{"unknown kind", Open{Kind: "epic", Intent: "x"}, "kind", nil},
		{"empty kind", Open{Intent: "x"}, "kind", nil},
		{"empty intent", Open{Kind: KindRecord}, "intent", nil},
		{"invalid UTF-8 intent", Open{Kind: KindRecord, Intent: "\xff"}, "intent", nil},
		{"unsorted sources", Open{Kind: KindRecord, Intent: "x", Sources: []id.ID{mustID(t, sSrc2), mustID(t, sSrc1)}}, "sources", nil},
		{"nothing in ops", Ops{}, "ops", nil},
		{"invalid op", Ops{Ops: []op.Op{good, bad}}, "ops[1]", op.ErrInvalid},
		{"duplicate op", Ops{Ops: []op.Op{good, good}}, "ops", nil},
		{"zero drop", Ops{Drop: []id.ID{{}}}, "drop", nil},
		{"empty accept rule", Accept{Reason: "x"}, "rule", nil},
		{"invalid UTF-8 reason", Reject{Rule: "r", Reason: "\xff"}, "reason", nil},
		{"zero by", Supersede{}, "by", nil},
		{"empty note", Note{}, "text", nil},
		{"invalid UTF-8 note", Note{Text: "\xff"}, "text", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			checkErr(t, tc.body.Validate(), tc.field, tc.also)
			_, err := EncodeBody(tc.body)
			checkErr(t, err, tc.field, tc.also)
		})
	}
	if _, err := EncodeBody(nil); !errors.Is(err, ErrInvalid) {
		t.Fatalf("EncodeBody(nil) = %v, want ErrInvalid", err)
	}
}

func TestBodyErrorNamesOpIndex(t *testing.T) {
	m := base(t, entry.TypeOps)
	setOp(1, "base", nil)(m)
	_, err := DecodeBody(entry.Entry{Type: entry.TypeOps, Body: canonical(t, m)})
	if err == nil || !strings.Contains(err.Error(), "ops[1]") || !errors.Is(err, op.ErrInvalid) || !errors.Is(err, ErrInvalid) {
		t.Fatalf("error = %v, want ops[1] wrapping op.ErrInvalid and ErrInvalid", err)
	}
}

func FuzzDecodeBody(f *testing.F) {
	for _, g := range golden(f) {
		f.Add(string(g.typ), []byte(g.doc))
	}
	for _, cases := range [][]mutation{shapeCases, ruleCases} {
		for _, tc := range cases {
			m := base(f, tc.typ)
			tc.edit(m)
			f.Add(string(tc.typ), []byte(canonical(f, m)))
		}
	}
	f.Add("merge", []byte(`{}`))
	f.Add("", []byte(``))
	f.Fuzz(func(t *testing.T, typ string, body []byte) {
		b, err := DecodeBody(entry.Entry{Type: entry.Type(typ), Body: json.RawMessage(body)})
		if err != nil {
			if !errors.Is(err, ErrInvalid) {
				t.Fatalf("DecodeBody error does not wrap ErrInvalid: %v", err)
			}
			return
		}
		if b.Type() != entry.Type(typ) {
			t.Fatalf("Type() = %s, want %s", b.Type(), typ)
		}
		out, err := EncodeBody(b)
		if err != nil {
			t.Fatalf("EncodeBody(DecodeBody(%q)) failed: %v", body, err)
		}
		if !bytes.Equal(out, body) {
			t.Fatalf("EncodeBody(DecodeBody(%q)) = %q", body, out)
		}
		again, err := DecodeBody(entry.Entry{Type: entry.Type(typ), Body: out})
		if err != nil {
			t.Fatalf("DecodeBody(EncodeBody(...)) failed: %v", err)
		}
		if !reflect.DeepEqual(again, b) {
			t.Fatalf("re-DecodeBody = %#v, want %#v", again, b)
		}
	})
}

const (
	sCarveEntry = "01JB00000000000000000000E9"
	sCarveCS    = "01JB0000000000000000000CS9"
)

// carveDocs are an ops body of each carve direction, spelled by hand, with
// the value each decodes to in a format-2 entry.
func carveDocs(t testing.TB) []struct {
	name string
	doc  string
	want Ops
} {
	create := op.Op{ID: mustID(t, sOp1), Type: op.TypeNodeCreate, Node: mustID(t, sNode), Kind: op.KindCapability,
		Handle: "token-refresh", Key: "i", Fields: map[string]json.RawMessage{"statement": json.RawMessage(`"Résumé"`)}}
	return []struct {
		name string
		doc  string
		want Ops
	}{
		{"to", `{
  "carve": {
    "entry": "` + sCarveEntry + `",
    "to": "` + sCarveCS + `"
  },
  "drop": [
    "` + sOp3 + `"
  ],
  "ops": []
}`, Ops{Drop: []id.ID{mustID(t, sOp3)}, Carve: &Carve{To: mustID(t, sCarveCS), Entry: mustID(t, sCarveEntry)}}},
		{"from", `{
  "carve": {
    "entry": "` + sCarveEntry + `",
    "from": "` + sCarveCS + `"
  },
  "drop": [],
  "ops": [
    {
      "fields": {
        "statement": "Résumé"
      },
      "handle": "token-refresh",
      "id": "` + sOp1 + `",
      "key": "i",
      "kind": "capability",
      "node": "` + sNode + `",
      "op": "node.create",
      "parent": null
    }
  ]
}`, Ops{Ops: []op.Op{create}, Carve: &Carve{From: mustID(t, sCarveCS), Entry: mustID(t, sCarveEntry)}}},
	}
}

func TestCarveRoundTrip(t *testing.T) {
	for _, c := range carveDocs(t) {
		t.Run(c.name, func(t *testing.T) {
			if err := canon.Check([]byte(c.doc + "\n")); err != nil {
				t.Fatalf("doc is not canonical: %v", err)
			}
			b, err := DecodeBody(entry.Entry{V: 2, Type: entry.TypeOps, Body: json.RawMessage(c.doc)})
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(b, c.want) {
				t.Fatalf("DecodeBody = %#v, want %#v", b, c.want)
			}
			out, err := EncodeBody(c.want)
			if err != nil {
				t.Fatal(err)
			}
			if string(out) != c.doc {
				t.Fatalf("EncodeBody =\n%s\nwant\n%s", out, c.doc)
			}
			if v := EntryVersion(c.want); v != 2 {
				t.Fatalf("EntryVersion = %d, want 2", v)
			}
		})
	}
}

func TestCarveNeedsFormat2(t *testing.T) {
	for _, c := range carveDocs(t) {
		for _, v := range []int{0, 1} {
			_, err := DecodeBody(entry.Entry{V: v, Type: entry.TypeOps, Body: json.RawMessage(c.doc)})
			checkErr(t, err, "carve", nil)
			if !strings.Contains(err.Error(), "unknown field") {
				t.Fatalf("v%d %s: error %q, want unknown field", v, c.name, err)
			}
		}
	}
	// carve belongs to ops bodies only.
	_, err := DecodeBody(entry.Entry{V: 2, Type: entry.TypeNote, Body: json.RawMessage(`{
  "carve": {},
  "text": "x"
}`)})
	checkErr(t, err, "carve", nil)
}

func TestVersionsWithoutCarve(t *testing.T) {
	for _, g := range golden(t) {
		if g.typ == entry.TypeReopen {
			continue // format 3 only; see TestReopenNeedsFormat3
		}
		for _, v := range []int{1, 2, 3} {
			b, err := DecodeBody(entry.Entry{V: v, Type: g.typ, Body: json.RawMessage(g.doc)})
			if err != nil {
				t.Fatalf("v%d %s: %v", v, g.typ, err)
			}
			if !reflect.DeepEqual(b, g.want) {
				t.Fatalf("v%d %s: DecodeBody = %#v, want %#v", v, g.typ, b, g.want)
			}
		}
		if v := EntryVersion(g.want); v != 1 {
			t.Fatalf("EntryVersion(%s) = %d, want 1", g.typ, v)
		}
	}
}

func TestReopenNeedsFormat3(t *testing.T) {
	out, err := EncodeBody(Reopen{})
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != `{}` {
		t.Fatalf("EncodeBody(Reopen{}) = %s, want {}", out)
	}
	if v := EntryVersion(Reopen{}); v != 3 {
		t.Fatalf("EntryVersion(Reopen{}) = %d, want 3", v)
	}
	for _, canonical := range []bool{false, true} {
		b, err := decodeBody(entry.Entry{V: 3, Type: entry.TypeReopen, Body: json.RawMessage(`{}`)}, canonical)
		if err != nil {
			t.Fatal(err)
		}
		if b != (Reopen{}) {
			t.Fatalf("decodeBody = %#v, want Reopen{}", b)
		}
		// Formats 1 and 2 cannot spell a reopen entry.
		for _, v := range []int{0, 1, 2} {
			_, err := decodeBody(entry.Entry{V: v, Type: entry.TypeReopen, Body: json.RawMessage(`{}`)}, canonical)
			checkErr(t, err, "type", nil)
			if !strings.Contains(err.Error(), "format 3") {
				t.Fatalf("v%d: error %q, want it to name format 3", v, err)
			}
		}
	}
	// The body is empty, as a land body is.
	_, err = DecodeBody(entry.Entry{V: 3, Type: entry.TypeReopen, Body: json.RawMessage(`{
  "x": 1
}`)})
	checkErr(t, err, "x", nil)
}

// setCarve sets fields on an ops body's carve object.
func setCarve(kv ...any) func(map[string]any) {
	return func(m map[string]any) {
		set(kv...)(m["carve"].(map[string]any))
	}
}

func delCarve(k string) func(map[string]any) {
	return func(m map[string]any) { delete(m["carve"].(map[string]any), k) }
}

func TestCarveInvalid(t *testing.T) {
	docs := carveDocs(t)
	to, from := docs[0].doc, docs[1].doc
	cases := []struct {
		name  string
		doc   string
		edit  func(map[string]any)
		field string
	}{
		{"both from and to", to, setCarve("from", sCarveCS), "carve"},
		{"neither from nor to", to, delCarve("to"), "carve"},
		{"missing entry", to, delCarve("entry"), "carve.entry"},
		{"extra key", to, setCarve("why", "x"), "carve.why"},
		{"extra keys sorted", from, setCarve("zz", 1, "aa", 2), "carve.aa"},
		{"zero entry", to, setCarve("entry", sZero), "carve.entry"},
		{"zero to", to, setCarve("to", sZero), "carve.to"},
		{"zero from", from, setCarve("from", sZero), "carve.from"},
		{"entry a number", to, setCarve("entry", 1), "carve.entry"},
		{"to not a ULID", to, setCarve("to", "nope"), "carve.to"},
		{"carve null", to, set("carve", nil), "carve"},
		{"carve an array", to, set("carve", []any{}), "carve"},
		{"carve.from with empty ops", from, func(m map[string]any) {
			m["ops"] = []any{}
			m["drop"] = []any{sOp3}
		}, "ops"},
		{"carve.to with empty drop", to, func(m map[string]any) {
			m["drop"] = []any{}
			m["ops"] = base(t, entry.TypeOps)["ops"]
		}, "drop"},
		{"missing drop beats carve", to, func(m map[string]any) {
			delete(m, "drop")
			m["carve"] = nil
		}, "drop"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var m map[string]any
			if err := json.Unmarshal([]byte(tc.doc), &m); err != nil {
				t.Fatal(err)
			}
			tc.edit(m)
			_, err := DecodeBody(entry.Entry{V: 2, Type: entry.TypeOps, Body: canonical(t, m)})
			checkErr(t, err, tc.field, nil)
		})
	}

	// The same rules hold on the way out.
	good := op.Op{ID: mustID(t, sOp2), Type: op.TypeLinkRetire, Link: mustID(t, sLink), Base: mustID(t, sBase)}
	cs, e := mustID(t, sCarveCS), mustID(t, sCarveEntry)
	drop := []id.ID{mustID(t, sOp3)}
	enc := []struct {
		name  string
		body  Ops
		field string
	}{
		{"both", Ops{Drop: drop, Carve: &Carve{From: cs, To: cs, Entry: e}}, "carve"},
		{"neither", Ops{Drop: drop, Carve: &Carve{Entry: e}}, "carve"},
		{"zero entry", Ops{Drop: drop, Carve: &Carve{To: cs}}, "carve.entry"},
		{"from without ops", Ops{Drop: drop, Carve: &Carve{From: cs, Entry: e}}, "ops"},
		{"to without drop", Ops{Ops: []op.Op{good}, Carve: &Carve{To: cs, Entry: e}}, "drop"},
	}
	for _, tc := range enc {
		t.Run("encode "+tc.name, func(t *testing.T) {
			checkErr(t, tc.body.Validate(), tc.field, nil)
			_, err := EncodeBody(tc.body)
			checkErr(t, err, tc.field, nil)
		})
	}
}

// TestDecodeBodyCanonicalMatches checks that the canonical path of
// decodeBody (FoldCanonical's) gives DecodeBody's body for every golden
// body and carve, and DecodeBody's error for every shape and rule mutation
// of them.
func TestDecodeBodyCanonicalMatches(t *testing.T) {
	for _, g := range golden(t) {
		e := entry.Entry{V: 3, Type: g.typ, Body: json.RawMessage(g.doc)}
		want, err := DecodeBody(e)
		if err != nil {
			t.Fatalf("%s: %v", g.typ, err)
		}
		got, err := decodeBody(e, true)
		if err != nil {
			t.Fatalf("%s: canonical: %v", g.typ, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("%s: canonical =\n%#v\nDecodeBody =\n%#v", g.typ, got, want)
		}
	}
	for _, d := range carveDocs(t) {
		e := entry.Entry{V: 2, Type: entry.TypeOps, Body: json.RawMessage(d.doc)}
		want, err := DecodeBody(e)
		if err != nil {
			t.Fatal(err)
		}
		got, err := decodeBody(e, true)
		if err != nil {
			t.Fatalf("canonical: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("carve: canonical =\n%#v\nDecodeBody =\n%#v", got, want)
		}
	}
	for _, cases := range [][]mutation{shapeCases, ruleCases} {
		for _, m := range cases {
			t.Run(m.name, func(t *testing.T) {
				doc := base(t, m.typ)
				m.edit(doc)
				e := entry.Entry{V: 1, Type: m.typ, Body: canonical(t, doc)}
				_, wantErr := DecodeBody(e)
				if wantErr == nil {
					t.Fatal("DecodeBody accepted the mutation")
				}
				_, gotErr := decodeBody(e, true)
				if gotErr == nil {
					t.Fatalf("canonical path accepted %s", e.Body)
				}
				if gotErr.Error() != wantErr.Error() {
					t.Fatalf("canonical error %q, DecodeBody error %q", gotErr, wantErr)
				}
			})
		}
	}
}

func TestOpenRef(t *testing.T) {
	base := Open{Kind: KindIntention, Intent: "x", Sources: []id.ID{}}
	// Without a ref the encoding and format are unchanged.
	plain, err := EncodeBody(base)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(plain), "ref") {
		t.Fatalf("encoding without a ref holds one: %s", plain)
	}
	if v := EntryVersion(base); v != 1 {
		t.Fatalf("EntryVersion without ref = %d, want 1", v)
	}
	// With a ref it round-trips and needs format 4.
	withRef := base
	withRef.Ref = "https://github.com/o/r/issues/7"
	if v := EntryVersion(withRef); v != 4 {
		t.Fatalf("EntryVersion with ref = %d, want 4", v)
	}
	raw, err := EncodeBody(withRef)
	if err != nil {
		t.Fatal(err)
	}
	for _, canonical := range []bool{false, true} {
		got, err := decodeBody(entry.Entry{V: 4, Type: entry.TypeOpen, Body: raw}, canonical)
		if err != nil {
			t.Fatal(err)
		}
		if o := got.(Open); o.Ref != withRef.Ref {
			t.Fatalf("ref = %q, want %q", o.Ref, withRef.Ref)
		}
	}
	// A reader at the previous format refuses it.
	for _, v := range []int{1, 2, 3} {
		_, err := DecodeBody(entry.Entry{V: v, Type: entry.TypeOpen, Body: raw})
		checkErr(t, err, "ref", nil)
	}
	// Empty is spelled by leaving it out.
	_, err = DecodeBody(entry.Entry{V: 4, Type: entry.TypeOpen, Body: json.RawMessage(strings.Replace(string(raw), withRef.Ref, "", 1))})
	checkErr(t, err, "ref", nil)
	// The limit is 2048 bytes.
	withRef.Ref = strings.Repeat("a", MaxRef)
	if _, err := EncodeBody(withRef); err != nil {
		t.Fatalf("ref of %d bytes: %v", MaxRef, err)
	}
	withRef.Ref += "a"
	_, err = EncodeBody(withRef)
	checkErr(t, err, "ref", nil)
}
