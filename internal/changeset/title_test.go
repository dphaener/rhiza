package changeset

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/entry"
)

func TestDisplayTitle(t *testing.T) {
	cases := []struct {
		intent, title, want string
	}{
		{"The Claude Code adapter: Genesis launches agents itself", "", "The Claude Code adapter"},
		{"Make tabs readable: name them", "Readable tabs", "Readable tabs"},
		{"Carving keeps IDs, always", "", "Carving keeps IDs"},
		{"Try it", "", "Try it"},
		{"Fix the flaky test.", "", "Fix the flaky test"},
		{"Fix the flaky test. Then ship", "", "Fix the flaky test"},
		{"Bump Go to 1.25 everywhere", "", "Bump Go to 1.25 everywhere"},
		{"Keep a,b intact", "", "Keep a,b intact"},
		{"Plan it; then do it", "", "Plan it"},
		{"Ship reports — CSV first", "", "Ship reports"},
		{"Ship reports - CSV first", "", "Ship reports"},
		{"Ship reports (CSV first)", "", "Ship reports"},
		{"Ship-ready reports", "", "Ship-ready reports"},
		{"  Padded intent  ", "", "Padded intent"},
		{"First line\nsecond line", "", "First line"},
		// Over MaxTitle: cut at the last space that fits.
		{"Record actor test@rhiza.invalid", "", "Record actor"},
		{"One two three four five six seven", "", "One two three four five six"},
		// A space right after rune 30 still cuts there.
		{"abcdefghij abcdefghij abcdefgh xyz", "", "abcdefghij abcdefghij abcdefgh"},
		// No space to cut at: a hard cut.
		{"Supercalifragilisticexpialidocious-ish", "", "Supercalifragilisticexpialidoc"},
		{"Éléphants éléphants éléphants éléphants", "", "Éléphants éléphants éléphants"},
		// The clause is empty: the intent's first runes.
		{": a leading colon", "", ": a leading colon"},
		{" - leading dash and then a long tail of words", "", "- leading dash and then a long"},
		{"(aside) and the rest of a long intent here", "", "(aside) and the rest of a long"},
		{"", "", ""},
	}
	for _, tc := range cases {
		got := Open{Kind: KindIntention, Intent: tc.intent, Title: tc.title}.DisplayTitle()
		if got != tc.want {
			t.Errorf("DisplayTitle(%q, title %q) = %q, want %q", tc.intent, tc.title, got, tc.want)
		}
		if n := utf8.RuneCountInString(got); n > MaxTitle {
			t.Errorf("DisplayTitle(%q) has %d runes, over %d", tc.intent, n, MaxTitle)
		}
	}
}

func TestDisplayTitleIsValid(t *testing.T) {
	for _, intent := range []string{
		"The Claude Code adapter: Genesis launches agents itself",
		"Supercalifragilisticexpialidocious-ish words",
		"Record actor test@rhiza.invalid",
		"  spaces  everywhere   in   this   very   long   intent  ",
	} {
		got := Open{Intent: intent}.DisplayTitle()
		if got == "" {
			t.Errorf("DisplayTitle(%q) is empty", intent)
		}
		if err := validTitle(got); err != nil {
			t.Errorf("DisplayTitle(%q) = %q is not a valid title: %v", intent, got, err)
		}
	}
}

func TestTitleValidate(t *testing.T) {
	ok := []string{"", "Readable tabs", strings.Repeat("é", MaxTitle), "a"}
	for _, title := range ok {
		if err := (Open{Kind: KindIntention, Intent: "x", Title: title}).Validate(); err != nil {
			t.Errorf("title %q: %v", title, err)
		}
	}
	bad := []string{
		strings.Repeat("a", MaxTitle+1),
		"This title is far too long to be accepted",
		" leading",
		"trailing ",
		"two\nlines",
		"carriage\rreturn",
		"tab\there",
		"\xff",
	}
	for _, title := range bad {
		body := Open{Kind: KindIntention, Intent: "x", Title: title}
		checkErr(t, body.Validate(), "title", nil)
		_, err := EncodeBody(body)
		checkErr(t, err, "title", nil)
	}
}

// titleDoc is an open body with a title, spelled by hand.
const titleDoc = `{
  "intent": "Make tabs readable: name them",
  "kind": "intention",
  "sources": [],
  "title": "Readable tabs"
}`

func TestTitleRoundTrip(t *testing.T) {
	if err := canon.Check([]byte(titleDoc + "\n")); err != nil {
		t.Fatalf("doc is not canonical: %v", err)
	}
	want := Open{Kind: KindIntention, Intent: "Make tabs readable: name them", Title: "Readable tabs"}
	for _, canonical := range []bool{false, true} {
		b, err := decodeBody(entry.Entry{V: 2, Type: entry.TypeOpen, Body: json.RawMessage(titleDoc)}, canonical)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(b, want) {
			t.Fatalf("DecodeBody = %#v, want %#v", b, want)
		}
	}
	out, err := EncodeBody(want)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != titleDoc {
		t.Fatalf("EncodeBody =\n%s\nwant\n%s", out, titleDoc)
	}
	if v := EntryVersion(want); v != 2 {
		t.Fatalf("EntryVersion = %d, want 2", v)
	}
	// Without a title the body is format 1 and has no title key.
	want.Title = ""
	if v := EntryVersion(want); v != 1 {
		t.Fatalf("EntryVersion without a title = %d, want 1", v)
	}
	out, err = EncodeBody(want)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "title") {
		t.Fatalf("EncodeBody without a title = %s, want no title key", out)
	}
}

func TestTitleNeedsFormat2(t *testing.T) {
	for _, v := range []int{0, 1} {
		_, err := DecodeBody(entry.Entry{V: v, Type: entry.TypeOpen, Body: json.RawMessage(titleDoc)})
		checkErr(t, err, "title", nil)
		if !strings.Contains(err.Error(), "unknown field") {
			t.Fatalf("v%d: error %q, want unknown field", v, err)
		}
	}
	// title belongs to open bodies only.
	_, err := DecodeBody(entry.Entry{V: 2, Type: entry.TypeNote, Body: json.RawMessage(`{
  "text": "x",
  "title": "x"
}`)})
	checkErr(t, err, "title", nil)
}

func TestTitleInvalid(t *testing.T) {
	cases := []struct {
		name  string
		edit  func(map[string]any)
		field string
	}{
		{"empty title", set("title", ""), "title"},
		{"title a number", set("title", 1), "title"},
		{"title null", set("title", nil), "title"},
		{"title too long", set("title", strings.Repeat("a", MaxTitle+1)), "title"},
		{"title on two lines", set("title", "a\nb"), "title"},
		{"title with a leading space", set("title", " a"), "title"},
		{"missing kind beats title", func(m map[string]any) { delete(m, "kind"); m["title"] = 1 }, "kind"},
		{"unknown beside title", set("zzz", 1), "zzz"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var m map[string]any
			if err := json.Unmarshal([]byte(titleDoc), &m); err != nil {
				t.Fatal(err)
			}
			tc.edit(m)
			e := entry.Entry{V: 2, Type: entry.TypeOpen, Body: canonical(t, m)}
			_, err := DecodeBody(e)
			checkErr(t, err, tc.field, nil)
			_, cerr := decodeBody(e, true)
			if cerr == nil || cerr.Error() != err.Error() {
				t.Fatalf("canonical error %v, DecodeBody error %v", cerr, err)
			}
		})
	}
}
