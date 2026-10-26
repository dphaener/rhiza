package guide

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

// checkCap reports a text whose body costs more than its cap.
func checkCap(t Text) error {
	if n := Tokens(t.Body); n > t.Cap() {
		return fmt.Errorf("%s: %d bytes is %d tokens, over the cap of %d", t.Name, len(t.Body), n, t.Cap())
	}
	return nil
}

// embedded returns the names of the embedded markdown files.
func embedded() []string {
	ents, err := files.ReadDir("text")
	if err != nil {
		panic(err)
	}
	var out []string
	for _, e := range ents {
		if strings.HasSuffix(e.Name(), ".md") {
			out = append(out, e.Name())
		}
	}
	return out
}

func TestCaps(t *testing.T) {
	for _, tx := range List() {
		n := Tokens(tx.Body)
		t.Logf("%-12s %5d bytes %5d tokens (cap %d)", tx.Name, len(tx.Body), n, tx.Cap())
		if err := checkCap(tx); err != nil {
			t.Error(err)
		}
	}
}

func TestCapCatchesOversize(t *testing.T) {
	g, _ := Lookup("genesis")
	g.Body += strings.Repeat("x", 9000)
	if err := checkCap(g); err == nil {
		t.Fatalf("checkCap(genesis + 9000 bytes) = nil, want an error")
	} else if !strings.Contains(err.Error(), "genesis") {
		t.Fatalf("checkCap error %q does not name the text", err)
	}
	topic, _ := Lookup("apply")
	topic.Body = strings.Repeat("x", 4*TopicCap+1)
	if checkCap(topic) == nil {
		t.Fatalf("checkCap(topic of %d tokens) = nil, want an error", TopicCap+1)
	}
	topic.Body = strings.Repeat("x", 4*TopicCap)
	if err := checkCap(topic); err != nil {
		t.Fatalf("checkCap(topic of %d tokens) = %v, want nil", TopicCap, err)
	}
}

func TestListOrder(t *testing.T) {
	var names []string
	for _, tx := range List() {
		names = append(names, tx.Name)
	}
	want := []string{"genesis", "scrub", "apply", "carving", "change-sets", "sources"}
	if !slices.Equal(names, want) {
		t.Fatalf("List names = %v, want %v", names, want)
	}
	for i, tx := range List() {
		if tx.Skill != (i < 2) {
			t.Errorf("%s: Skill = %v", tx.Name, tx.Skill)
		}
	}
}

func TestEveryFileListed(t *testing.T) {
	listed := map[string]bool{}
	for _, e := range slices.Concat(skills, topics) {
		listed[e.file] = true
	}
	files := embedded()
	if len(files) != len(listed) {
		t.Errorf("embedded %d files %v, listed %d", len(files), files, len(listed))
	}
	for _, f := range files {
		if !listed[f] {
			t.Errorf("text/%s is embedded but not listed", f)
		}
	}
}

func TestLookup(t *testing.T) {
	for _, tx := range List() {
		got, ok := Lookup(tx.Name)
		if !ok || got != tx {
			t.Errorf("Lookup(%q) = %v, %v; want the listed text", tx.Name, got.Name, ok)
		}
		if strings.TrimSpace(tx.Body) == "" {
			t.Errorf("%s: empty body", tx.Name)
		}
	}
	g, ok := Lookup("")
	if !ok || g.Name != "genesis" {
		t.Errorf(`Lookup("") = %q, %v; want genesis`, g.Name, ok)
	}
	if _, ok := Lookup("nope"); ok {
		t.Errorf(`Lookup("nope") found a text`)
	}
	if _, ok := Lookup("genesis.md"); ok {
		t.Errorf(`Lookup("genesis.md") found a text; names have no extension`)
	}
}

func TestTokens(t *testing.T) {
	for _, tc := range []struct {
		s    string
		want int
	}{
		{"", 0},
		{"a", 1},
		{"abcd", 1},
		{"abcde", 2},
		{strings.Repeat("x", 8), 2},
		{strings.Repeat("x", 8000), 2000},
		{strings.Repeat("x", 8001), 2001},
	} {
		if got := Tokens(tc.s); got != tc.want {
			t.Errorf("Tokens(%d bytes) = %d, want %d", len(tc.s), got, tc.want)
		}
	}
}

func TestDescriptions(t *testing.T) {
	for _, tx := range List() {
		if strings.TrimSpace(tx.Description) == "" {
			t.Errorf("%s: empty description", tx.Name)
		}
		if len(tx.Description) > MaxDescription {
			t.Errorf("%s: description is %d bytes, want at most %d", tx.Name, len(tx.Description), MaxDescription)
		}
	}
}

func TestList(t *testing.T) {
	a := List()
	a[0].Body = "changed"
	if b, _ := Lookup("genesis"); b.Body == "changed" {
		t.Fatalf("List returned the package's own slice")
	}
}
