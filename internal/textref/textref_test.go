package textref

import (
	"strings"
	"testing"
)

const (
	idA = "01M6HP0CQRJFFH7Q40JCV9H2DG"
	idB = "01M6HP0GQ76PY04WKDAPVY30Q6"
	idR = "01M6HP108NW407PT5CRA556JBA" // retired
)

var handles = map[string]string{"alpha": idA, "beta-two": idB}

// lookup resolves handles, full IDs (retired included) and 4+ character
// prefixes of IDs.
func lookup(r Ref) (string, bool) {
	if r.Kind == KindHandle {
		id, ok := handles[r.Token]
		return id, ok
	}
	var hit string
	for _, id := range []string{idA, idB, idR} {
		if strings.HasPrefix(id, r.Token) {
			if hit != "" {
				return "", false
			}
			hit = id
		}
	}
	return hit, hit != ""
}

func handleOf(id string) (string, bool) {
	switch id {
	case idA:
		return "alpha", true
	case idB:
		return "beta-two", true
	}
	return "", false
}

func TestFind(t *testing.T) {
	for text, want := range map[string][]string{
		"see @alpha.":               {"alpha"},
		"(@alpha) [@beta-two]":      {"alpha", "beta-two"},
		"@01M6HP0C and @" + idB:     {"01M6HP0C", idB},
		"@{" + idA + "} @{BOGUS}":   {idA, "BOGUS"},
		"darin@hey.com checkout@v7": nil,
		"quoted '@key'":             nil,
		"@@alpha":                   nil,
		"@Alpha @foo_bar @foo- @{}": nil,
		"@ABC @{ @{x":               nil,
	} {
		var got []string
		for _, r := range Find(text) {
			got = append(got, r.Token)
		}
		if strings.Join(got, "|") != strings.Join(want, "|") {
			t.Errorf("Find(%q) = %q, want %q", text, got, want)
		}
	}
}

func TestResolve(t *testing.T) {
	in := "a @alpha, @01M6HP0G @{" + idA + "} @@alpha x@alpha @{BOGUS} @nope"
	got, un := Resolve(in, lookup)
	want := "a @{" + idA + "}, @{" + idB + "} @{" + idA + "} @alpha x@alpha @{BOGUS} @nope"
	if got != want {
		t.Errorf("got %q\nwant %q", got, want)
	}
	if len(un) != 2 || un[0].Token != "BOGUS" || un[0].Start != strings.Index(in, "@{BOGUS}") || un[1].Token != "nope" {
		t.Errorf("unresolved = %+v", un)
	}
}

// A literal @{<full ID>} cannot be stored: it would read back as a
// reference. The @@ escape of one, and an unmapped one, are returned and
// stored without braces.
func TestResolveLiteralFullID(t *testing.T) {
	const ghost = "01ZZZZZZZZZZZZZZZZZZZZZZZZ"
	in := "a @@{" + ghost + "} b @{" + ghost + "} c @@{" + idA + "} d @@{ID}"
	got, un := Resolve(in, lookup)
	want := "a @" + ghost + " b @" + ghost + " c @" + idA + " d @{ID}"
	if got != want {
		t.Errorf("got %q\nwant %q", got, want)
	}
	if len(un) != 3 || !un[0].Escaped || un[0].Token != ghost || un[0].Start != 2 || un[0].End != 2+len("@@{")+IDLen+1 ||
		un[1].Escaped || un[1].Token != ghost || !un[2].Escaped || un[2].Token != idA {
		t.Errorf("unresolved = %+v", un)
	}
	if shown := Format(got, handleOf); strings.Contains(shown, "@{"+ghost) || strings.Contains(shown, "@alpha") {
		t.Errorf("stored %q is shown as %q, with a reference", got, shown)
	}
}

func TestFormat(t *testing.T) {
	for stored, want := range map[string]string{
		"see @{" + idA + "}.":                 "see @alpha.",
		"@{" + idA + "}x":                     "@{" + idA + "}x",
		"@{" + idR + "} gone":                 "@{" + idR + "} gone",
		"@{ID} @{BOGUS} @{" + idA[:25] + "}":  "@@{ID} @@{BOGUS} @@{" + idA[:25] + "}",
		"@{<full ID>} @{<ID>} @{...} @{ } @{": "@@{<full ID>} @@{<ID>} @@{...} @@{ } @@{",
		"(@{<26-character ID>})":              "(@@{<26-character ID>})",
		"@word":                               "@@word",
		"@@word":                              "@@@word",
		"darin@hey.com checkout@v7 '@key'":    "darin@hey.com checkout@v7 '@key'",
		"@ alone @":                           "@ alone @",
		"no refs":                             "no refs",
	} {
		if got := Format(stored, handleOf); got != want {
			t.Errorf("Format(%q) = %q, want %q", stored, got, want)
		}
	}
}

func FuzzRoundTrip(f *testing.F) {
	for _, s := range []string{"@alpha", "@@x", "(@beta-two)@{" + idA + "}", "@{" + idR + "}-", "@@@alpha @{ID}", "a @01M6 b @", "@{" + idA + "}@{" + idB + "}"} {
		f.Add(s)
	}
	f.Add("@@{01ZZZZZZZZZZZZZZZZZZZZZZZZ} @{01ZZZZZZZZZZZZZZZZZZZZZZZZ}")
	f.Fuzz(func(t *testing.T, in string) {
		gave := map[string]bool{}
		stored, un := Resolve(in, func(r Ref) (string, bool) {
			id, ok := lookup(r)
			gave[id] = gave[id] || ok
			return id, ok
		})
		// Every stored reference is one lookup gave, unresolved input or not.
		for i := 0; i < len(stored); i++ {
			if id, ok := storedRef(stored, i); ok && startOK(stored, i) && !gave[id] {
				t.Fatalf("in %q stored %q holds @{%s}, which lookup did not give", in, stored, id)
			}
		}
		if len(un) > 0 {
			return
		}
		shown := Format(stored, handleOf)
		again, un := Resolve(shown, lookup)
		if len(un) > 0 || again != stored {
			t.Fatalf("in %q stored %q shown %q again %q unresolved %v", in, stored, shown, again, un)
		}
	})
}

// FuzzStoredRoundTrip starts from arbitrary stored text, legacy look-alikes
// included: showing it and writing it back stores the same text.
func FuzzStoredRoundTrip(f *testing.F) {
	for _, s := range []string{"@word", "@@word", "@{ID}", "@{" + idA + "}", "x @{BOGUS} (@{" + idR + "})", "@"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, stored string) {
		shown := Format(stored, handleOf)
		again, un := Resolve(shown, lookup)
		if len(un) > 0 {
			return // a stored reference naming no node is refused, by design
		}
		if again != stored {
			t.Fatalf("stored %q shown %q again %q", stored, shown, again)
		}
	})
}
