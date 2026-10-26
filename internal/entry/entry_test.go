package entry

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/id"
	"pgregory.net/rapid"
)

const (
	sID     = "01JB000000000000000000000A"
	sCS     = "01JB00000000000000000000CS"
	sActor  = "01JB0000000000000000000ACT"
	sFor    = "01JB0000000000000000000F0R"
	sAfter1 = "01JB00000000000000000000A1"
	sAfter2 = "01JB00000000000000000000A2"
	sAt     = "2026-10-24T13:05:09Z"
)

func mustID(t testing.TB, s string) id.ID {
	t.Helper()
	i, err := id.Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	return i
}

// doc spells a canonical entry by hand, so the golden bytes do not depend
// on Encode. after and body are already indented for their position.
func doc(after, typ, body string) string {
	return "{\n" +
		`  "actor": "` + sActor + "\",\n" +
		`  "after": ` + after + ",\n" +
		`  "at": "` + sAt + "\",\n" +
		`  "body": ` + body + ",\n" +
		`  "change_set": "` + sCS + "\",\n" +
		`  "for": "` + sFor + "\",\n" +
		`  "id": "` + sID + "\",\n" +
		`  "type": "` + typ + "\",\n" +
		`  "v": 1` + "\n" +
		"}\n"
}

const twoAfter = "[\n    \"" + sAfter1 + "\",\n    \"" + sAfter2 + "\"\n  ]"

func baseEntry(t testing.TB) Entry {
	return Entry{
		V:         1,
		ID:        mustID(t, sID),
		ChangeSet: mustID(t, sCS),
		Actor:     mustID(t, sActor),
		For:       mustID(t, sFor),
		At:        time.Date(2026, 10, 24, 13, 5, 9, 0, time.UTC),
		After:     []id.ID{},
		Type:      TypeOpen,
		Body:      json.RawMessage(`{}`),
	}
}

var goldenCases = []struct {
	typ       Type
	after     bool
	body      string // standalone canonical body
	bodyInDoc string // the same body as it appears in the document
}{
	{TypeOpen, false, "{\n  \"title\": \"Résumé <&>  \"\n}", "{\n    \"title\": \"Résumé <&>  \"\n  }"},
	{TypeOps, true, "{\n  \"ops\": [\n    {\n      \"n\": -3\n    },\n    []\n  ]\n}",
		"{\n    \"ops\": [\n      {\n        \"n\": -3\n      },\n      []\n    ]\n  }"},
	{TypePropose, true, `{}`, `{}`},
	{TypeAccept, false, "{\n  \"a\": true,\n  \"b\": null\n}", "{\n    \"a\": true,\n    \"b\": null\n  }"},
	{TypeReject, true, "{\n  \"reason\": \"no\\n\\\"really\\\"\"\n}", "{\n    \"reason\": \"no\\n\\\"really\\\"\"\n  }"},
	{TypeWithdraw, false, `{}`, `{}`},
	{TypeSupersede, true, "{\n  \"by\": {}\n}", "{\n    \"by\": {}\n  }"},
	{TypeLand, false, "{\n  \"commit\": \"abc\"\n}", "{\n    \"commit\": \"abc\"\n  }"},
	{TypeNote, true, "{\n  \"ζ\": 9223372036854775807\n}", "{\n    \"ζ\": 9223372036854775807\n  }"},
}

func goldenEntry(t testing.TB, i int) (Entry, string) {
	tc := goldenCases[i]
	e := baseEntry(t)
	e.Type = tc.typ
	e.Body = json.RawMessage(tc.body)
	after := "[]"
	if tc.after {
		e.After = []id.ID{mustID(t, sAfter1), mustID(t, sAfter2)}
		after = twoAfter
	}
	return e, doc(after, string(tc.typ), tc.bodyInDoc)
}

func TestTypes(t *testing.T) {
	want := []Type{"open", "ops", "propose", "accept", "reject", "withdraw", "supersede", "land", "note"}
	if got := Types(); !slices.Equal(got, want) {
		t.Fatalf("Types() = %v, want %v", got, want)
	}
	if len(goldenCases) != len(want) {
		t.Fatalf("golden table covers %d types, want %d", len(goldenCases), len(want))
	}
}

func TestGolden(t *testing.T) {
	for i, tc := range goldenCases {
		t.Run(string(tc.typ), func(t *testing.T) {
			e, want := goldenEntry(t, i)
			got, err := Encode(e)
			if err != nil {
				t.Fatalf("Encode: %v", err)
			}
			if string(got) != want {
				t.Fatalf("Encode =\n%s\nwant\n%s", got, want)
			}
			back, err := Decode([]byte(want))
			if err != nil {
				t.Fatalf("Decode: %v", err)
			}
			if !reflect.DeepEqual(back, e) {
				t.Fatalf("Decode =\n%#v\nwant\n%#v", back, e)
			}
		})
	}
}

func TestVersions(t *testing.T) {
	if FormatVersion != 2 {
		t.Fatalf("FormatVersion = %d, want 2", FormatVersion)
	}
	// Both formats read and write; the envelope is the same in each.
	for v := 1; v <= FormatVersion; v++ {
		t.Run(strconv.Itoa(v), func(t *testing.T) {
			in := swap(`"v": 1`, `"v": `+strconv.Itoa(v))
			e, err := Decode([]byte(in))
			if err != nil {
				t.Fatalf("Decode: %v", err)
			}
			if e.V != v {
				t.Fatalf("V = %d, want %d", e.V, v)
			}
			out, err := Encode(e)
			if err != nil {
				t.Fatalf("Encode: %v", err)
			}
			if string(out) != in {
				t.Fatalf("Encode =\n%s\nwant\n%s", out, in)
			}
		})
	}
}

func TestDecodeEmptyAfterIsNonNil(t *testing.T) {
	e, err := Decode([]byte(doc("[]", "note", "{}")))
	if err != nil {
		t.Fatal(err)
	}
	if e.After == nil || len(e.After) != 0 {
		t.Fatalf("After = %#v, want non-nil empty", e.After)
	}
}

func TestEncodeNilAfter(t *testing.T) {
	e := baseEntry(t)
	e.After = nil
	got, err := Encode(e)
	if err != nil {
		t.Fatal(err)
	}
	if want := doc("[]", "open", "{}"); string(got) != want {
		t.Fatalf("Encode =\n%s\nwant\n%s", got, want)
	}
}

func TestEncodeAllowsAnyULIDTime(t *testing.T) {
	// D16: the ULID's embedded time is never compared with at.
	e := baseEntry(t)
	e.At = time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)
	data, err := Encode(e)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Decode(data); err != nil {
		t.Fatal(err)
	}
}

var valid = doc("[]", "note", "{}")

func swap(old, new string) string {
	if !strings.Contains(valid, old) {
		panic("swap: " + old + " not in valid doc")
	}
	return strings.Replace(valid, old, new, 1)
}

var decodeFailures = []struct {
	name  string
	in    string
	is    error
	field string // quoted field name the message must contain; "" for none
	msg   string // other text the message must contain
}{
	{"v 0", swap(`"v": 1`, `"v": 0`), ErrInvalid, "v", ""},
	{"v negative", swap(`"v": 1`, `"v": -1`), ErrInvalid, "v", ""},
	{"v string", swap(`"v": 1`, `"v": "1"`), ErrInvalid, "v", ""},
	{"v missing", swap(`,`+"\n"+`  "v": 1`, ``), ErrInvalid, "v", ""},
	{"v 3", swap(`"v": 1`, `"v": 3`), ErrNewerFormat, "v", ""},
	{"v 3 with unknown keys", "{\n  \"future\": [],\n  \"v\": 3\n}\n", ErrNewerFormat, "v", ""},
	{"missing for", swap(`  "for": "`+sFor+`",`+"\n", ``), ErrInvalid, "for", "missing"},
	{"extra key", swap(`  "for"`, `  "extra": 1,`+"\n"+`  "for"`), ErrInvalid, "extra", "unknown"},
	{"lower-case id", swap(sID, strings.ToLower(sID)), ErrInvalid, "id", ""},
	{"id number", swap(`"`+sID+`"`, `7`), ErrInvalid, "id", ""},
	{"zero id", swap(sID, "00000000000000000000000000"), ErrInvalid, "id", ""},
	{"lower-case change_set", swap(sCS, strings.ToLower(sCS)), ErrInvalid, "change_set", ""},
	{"bad actor", swap(sActor, "nope"), ErrInvalid, "actor", ""},
	{"at milliseconds", swap(sAt, "2026-10-24T13:05:09.000Z"), ErrInvalid, "at", ""},
	{"at offset", swap(sAt, "2026-10-24T13:05:09+00:00"), ErrInvalid, "at", ""},
	{"at lower z", swap(sAt, "2026-10-24T13:05:09z"), ErrInvalid, "at", ""},
	{"at number", swap(`"`+sAt+`"`, `0`), ErrInvalid, "at", ""},
	{"after unsorted", swap(`"after": []`, "\"after\": [\n    \""+sAfter2+"\",\n    \""+sAfter1+"\"\n  ]"),
		ErrInvalid, "after", "not sorted"},
	{"after duplicate", swap(`"after": []`, "\"after\": [\n    \""+sAfter1+"\",\n    \""+sAfter1+"\"\n  ]"),
		ErrInvalid, "after", "duplicate"},
	{"after null", swap(`"after": []`, `"after": null`), ErrInvalid, "after", ""},
	{"after object", swap(`"after": []`, `"after": {}`), ErrInvalid, "after", ""},
	{"after numbers", swap(`"after": []`, "\"after\": [\n    1\n  ]"), ErrInvalid, "after", ""},
	{"after lower-case", swap(`"after": []`, "\"after\": [\n    \""+strings.ToLower(sAfter1)+"\"\n  ]"),
		ErrInvalid, "after", ""},
	{"unknown type", swap(`"type": "note"`, `"type": "bogus"`), ErrInvalid, "type", ""},
	{"type number", swap(`"type": "note"`, `"type": 1`), ErrInvalid, "type", ""},
	{"body array", swap(`"body": {}`, `"body": []`), ErrInvalid, "body", ""},
	{"body string", swap(`"body": {}`, `"body": "x"`), ErrInvalid, "body", ""},
	{"body null", swap(`"body": {}`, `"body": null`), ErrInvalid, "body", ""},
}

var shapeFailures = []struct {
	name string
	in   string
	is   []error
	msg  string
}{
	{"non-canonical whitespace", `{"actor":"` + sActor + `","after":[],"at":"` + sAt + `","body":{},"change_set":"` + sCS +
		`","for":"` + sFor + `","id":"` + sID + `","type":"note","v":1}` + "\n",
		[]error{ErrInvalid, canon.ErrNotCanonical}, "not canonical"},
	{"missing final newline", strings.TrimSuffix(valid, "\n"), []error{ErrInvalid, canon.ErrNotCanonical}, "not canonical"},
	{"duplicate key", strings.Replace(valid, `  "v": 1`, `  "v": 1,`+"\n"+`  "v": 1`, 1),
		[]error{ErrInvalid, canon.ErrInvalid}, "duplicate"},
	{"not json", "nope", []error{ErrInvalid, canon.ErrInvalid}, ""},
	{"empty", "", []error{ErrInvalid, canon.ErrInvalid}, ""},
	{"top-level array", "[]\n", []error{ErrInvalid}, "object"},
	{"top-level null", "null\n", []error{ErrInvalid}, "object"},
	{"top-level string", "\"x\"\n", []error{ErrInvalid}, "object"},
}

func TestDecodeFailures(t *testing.T) {
	for _, tc := range decodeFailures {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Decode([]byte(tc.in))
			if err == nil {
				t.Fatalf("Decode(%q) succeeded", tc.in)
			}
			if !errors.Is(err, tc.is) {
				t.Fatalf("error %q does not wrap %v", err, tc.is)
			}
			if tc.is == ErrNewerFormat && errors.Is(err, ErrInvalid) {
				t.Fatalf("newer-format error %q also wraps ErrInvalid", err)
			}
			if want := `field "` + tc.field + `"`; !strings.Contains(err.Error(), want) {
				t.Fatalf("error %q does not name %s", err, want)
			}
			if !strings.Contains(err.Error(), tc.msg) {
				t.Fatalf("error %q does not contain %q", err, tc.msg)
			}
		})
	}
	for _, tc := range shapeFailures {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Decode([]byte(tc.in))
			if err == nil {
				t.Fatalf("Decode(%q) succeeded", tc.in)
			}
			for _, is := range tc.is {
				if !errors.Is(err, is) {
					t.Fatalf("error %q does not wrap %v", err, is)
				}
			}
			if !strings.Contains(err.Error(), tc.msg) {
				t.Fatalf("error %q does not contain %q", err, tc.msg)
			}
		})
	}
}

func TestEncodeFailures(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*Entry)
		field  string
	}{
		{"v 0", func(e *Entry) { e.V = 0 }, "v"},
		{"v 3", func(e *Entry) { e.V = 3 }, "v"},
		{"v negative", func(e *Entry) { e.V = -1 }, "v"},
		{"zero id", func(e *Entry) { e.ID = id.ID{} }, "id"},
		{"zero change_set", func(e *Entry) { e.ChangeSet = id.ID{} }, "change_set"},
		{"zero actor", func(e *Entry) { e.Actor = id.ID{} }, "actor"},
		{"zero for", func(e *Entry) { e.For = id.ID{} }, "for"},
		{"non-UTC at", func(e *Entry) { e.At = e.At.In(time.FixedZone("X", 3600)) }, "at"},
		{"UTC-offset zone at", func(e *Entry) { e.At = e.At.In(time.FixedZone("UTC", 0)) }, "at"},
		{"sub-second at", func(e *Entry) { e.At = e.At.Add(time.Millisecond) }, "at"},
		{"year 10000 at", func(e *Entry) { e.At = time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC) }, "at"},
		{"negative year at", func(e *Entry) { e.At = time.Date(-1, 1, 1, 0, 0, 0, 0, time.UTC) }, "at"},
		{"unsorted after", func(e *Entry) {
			e.After = []id.ID{mustID(t, sAfter2), mustID(t, sAfter1)}
		}, "after"},
		{"duplicate after", func(e *Entry) {
			e.After = []id.ID{mustID(t, sAfter1), mustID(t, sAfter1)}
		}, "after"},
		{"unknown type", func(e *Entry) { e.Type = "bogus" }, "type"},
		{"empty type", func(e *Entry) { e.Type = "" }, "type"},
		{"non-canonical body", func(e *Entry) { e.Body = json.RawMessage(`{"a":1}`) }, "body"},
		{"body with newline", func(e *Entry) { e.Body = json.RawMessage("{}\n") }, "body"},
		{"body array", func(e *Entry) { e.Body = json.RawMessage(`[]`) }, "body"},
		{"body nil", func(e *Entry) { e.Body = nil }, "body"},
		{"body invalid", func(e *Entry) { e.Body = json.RawMessage(`{`) }, "body"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := baseEntry(t)
			tc.mutate(&e)
			if err := e.Validate(); !errors.Is(err, ErrInvalid) {
				t.Fatalf("Validate = %v, want ErrInvalid", err)
			}
			_, err := Encode(e)
			if !errors.Is(err, ErrInvalid) {
				t.Fatalf("Encode error = %v, want ErrInvalid", err)
			}
			if want := `field "` + tc.field + `"`; !strings.Contains(err.Error(), want) {
				t.Fatalf("error %q does not name %s", err, want)
			}
		})
	}
}

func FuzzDecode(f *testing.F) {
	for i := range goldenCases {
		_, golden := goldenEntry(f, i)
		f.Add([]byte(golden))
	}
	for _, tc := range decodeFailures {
		f.Add([]byte(tc.in))
	}
	for _, tc := range shapeFailures {
		f.Add([]byte(tc.in))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		e, err := Decode(data)
		if err != nil {
			if !errors.Is(err, ErrInvalid) && !errors.Is(err, ErrNewerFormat) {
				t.Fatalf("Decode error wraps neither sentinel: %v", err)
			}
			return
		}
		out, err := Encode(e)
		if err != nil {
			t.Fatalf("Encode(Decode(%q)) failed: %v", data, err)
		}
		if !bytes.Equal(out, data) {
			t.Fatalf("Encode(Decode(%q)) = %q", data, out)
		}
	})
}

func genID() *rapid.Generator[id.ID] {
	return rapid.Custom(func(t *rapid.T) id.ID {
		var i id.ID
		copy(i[:], rapid.SliceOfN(rapid.Byte(), 16, 16).Draw(t, "id bytes"))
		return i
	}).Filter(func(i id.ID) bool { return !i.IsZero() })
}

var bodyKeys = []string{"a", "b", "z", "Z", "é", "ζ", "名", "😀", "a b", ""}

func genValue(t *rapid.T, depth int) any {
	kinds := 4
	if depth < 3 {
		kinds = 6
	}
	switch rapid.IntRange(0, kinds-1).Draw(t, "kind") {
	case 0:
		return nil
	case 1:
		return rapid.Bool().Draw(t, "bool")
	case 2:
		return rapid.Int64().Draw(t, "int")
	case 3:
		return rapid.String().Draw(t, "string")
	case 4:
		n := rapid.IntRange(0, 3).Draw(t, "len")
		arr := make([]any, n)
		for i := range arr {
			arr[i] = genValue(t, depth+1)
		}
		return arr
	default:
		return genObject(t, depth+1)
	}
}

func genObject(t *rapid.T, depth int) map[string]any {
	n := rapid.IntRange(0, 4).Draw(t, "members")
	obj := make(map[string]any, n)
	for range n {
		obj[rapid.SampledFrom(bodyKeys).Draw(t, "key")] = genValue(t, depth)
	}
	return obj
}

func genEntry() *rapid.Generator[Entry] {
	return rapid.Custom(func(t *rapid.T) Entry {
		after := rapid.SliceOfNDistinct(genID(), 0, 4, func(i id.ID) id.ID { return i }).Draw(t, "after")
		slices.SortFunc(after, id.Compare)
		if after == nil {
			after = []id.ID{}
		}
		body, err := canon.Marshal(genObject(t, 0))
		if err != nil {
			t.Fatalf("canon.Marshal body: %v", err)
		}
		return Entry{
			V:         rapid.IntRange(1, FormatVersion).Draw(t, "v"),
			ID:        genID().Draw(t, "id"),
			ChangeSet: genID().Draw(t, "change_set"),
			Actor:     genID().Draw(t, "actor"),
			For:       genID().Draw(t, "for"),
			// Years 0000 through 9999, the range RFC 3339 can spell.
			At:    time.Unix(rapid.Int64Range(-62167219200, 253402300799).Draw(t, "at"), 0).UTC(),
			After: after,
			Type:  rapid.SampledFrom(Types()).Draw(t, "type"),
			Body:  json.RawMessage(bytes.TrimSuffix(body, []byte("\n"))),
		}
	})
}

func TestRoundTripProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		e := genEntry().Draw(t, "entry")
		data, err := Encode(e)
		if err != nil {
			t.Fatalf("Encode: %v", err)
		}
		if err := canon.Check(data); err != nil {
			t.Fatalf("Encode output not canonical: %v", err)
		}
		back, err := Decode(data)
		if err != nil {
			t.Fatalf("Decode(%q): %v", data, err)
		}
		if !reflect.DeepEqual(back, e) {
			t.Fatalf("Decode(Encode(e)) =\n%#v\nwant\n%#v", back, e)
		}
	})
}

func TestEncodeDeterministicProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		e := genEntry().Draw(t, "entry")
		a, err := Encode(e)
		if err != nil {
			t.Fatalf("Encode: %v", err)
		}
		b, err := Encode(e)
		if err != nil {
			t.Fatalf("Encode: %v", err)
		}
		if !bytes.Equal(a, b) {
			t.Fatalf("Encode not deterministic:\n%s\n%s", a, b)
		}
	})
}

func TestDecodeCanonicalMatchesDecode(t *testing.T) {
	for i := range goldenCases {
		e, doc := goldenEntry(t, i)
		want, err := Decode([]byte(doc))
		if err != nil {
			t.Fatal(err)
		}
		got, err := DecodeCanonical([]byte(doc))
		if err != nil {
			t.Fatalf("%s: %v", e.Type, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("%s: DecodeCanonical = %+v, Decode = %+v", e.Type, got, want)
		}
	}
}

func TestDecodeCanonicalRejects(t *testing.T) {
	good := doc("[]", "note", "{}")
	for _, tc := range []struct {
		name, in string
		is       error
		field    string
	}{
		{"missing field", strings.Replace(good, "  \"for\": \""+sFor+"\",\n", "", 1), ErrInvalid, "for"},
		{"unknown field", strings.Replace(good, "  \"id\"", "  \"extra\": 1,\n  \"id\"", 1), ErrInvalid, "extra"},
		{"bad version", strings.Replace(good, `"v": 1`, `"v": 0`, 1), ErrInvalid, "v"},
		{"newer version", strings.Replace(good, `"v": 1`, `"v": 99`, 1), ErrNewerFormat, "v"},
		{"empty", "", ErrInvalid, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecodeCanonical([]byte(tc.in))
			if err == nil {
				t.Fatalf("DecodeCanonical(%q) succeeded", tc.in)
			}
			if !errors.Is(err, tc.is) {
				t.Fatalf("error %q does not wrap %v", err, tc.is)
			}
			if tc.is == ErrNewerFormat && errors.Is(err, ErrInvalid) {
				t.Fatalf("newer-format error %q also wraps ErrInvalid", err)
			}
			if tc.field != "" && !strings.Contains(err.Error(), `field "`+tc.field+`"`) {
				t.Fatalf("error %q does not name field %q", err, tc.field)
			}
		})
	}
}

func TestDecodeCanonicalProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		data, err := Encode(genEntry().Draw(t, "entry"))
		if err != nil {
			t.Fatalf("Encode: %v", err)
		}
		want, err := Decode(data)
		if err != nil {
			t.Fatalf("Decode(%q): %v", data, err)
		}
		got, err := DecodeCanonical(data)
		if err != nil {
			t.Fatalf("DecodeCanonical(%q): %v", data, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("DecodeCanonical =\n%#v\nDecode =\n%#v", got, want)
		}
	})
}
