package index

import (
	"slices"
	"testing"
)

func TestTokens(t *testing.T) {
	got := Tokens("Search-FTS: asis, 42 Dîsk")
	want := []string{"search", "fts", "asis", "42", "dîsk"}
	if !slices.Equal(got, want) {
		t.Errorf("Tokens = %q, want %q", got, want)
	}
	if len(Tokens(" -- ")) != 0 {
		t.Errorf("punctuation alone gave tokens")
	}
}

func TestSearchNeedsTerms(t *testing.T) {
	f := newFixture(t)
	x, err := Open(f.root, f.dbPath())
	if err != nil {
		t.Fatal(err)
	}
	defer x.Close()
	if _, err := x.Search(AsIs, nil, nil); err == nil {
		t.Error("Search with no terms succeeded")
	}
	if ns, err := x.Search(AsIs, []string{"disk"}, nil); err != nil || len(ns) != 0 {
		t.Errorf("Search on an empty log = %v, %v", ns, err)
	}
}
