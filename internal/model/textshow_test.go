package model

import (
	"bytes"
	"encoding/json"
	"github.com/dphaener/rhiza/internal/canon"
	"testing"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

func TestShowFields(t *testing.T) {
	live := "01M6HNNWFR1J2AHQRNKFTPME46"
	gone := "01M6HNNWFR1J2AHQRNKFTPME47"
	liveID, _ := id.Parse(live)
	handleOf := func(i id.ID) (string, bool) { return "renamed", i == liveID }
	fields := map[string]json.RawMessage{
		"statement": json.RawMessage(`"See @{` + live + `} and @{` + gone + `}, @word"`),
		"quote":     json.RawMessage(`"@word"`),
		"n":         json.RawMessage(`3`),
	}
	got := ShowFields(op.KindDecision, fields, handleOf)
	if want := `"See @renamed and @{` + gone + `}, @@word"`; string(got["statement"]) != want {
		t.Errorf("statement = %s, want %s", got["statement"], want)
	}
	if string(got["quote"]) != `"@word"` || string(got["n"]) != "3" {
		t.Errorf("verbatim or non-text field changed: %v", got)
	}
	if string(fields["statement"]) == string(got["statement"]) {
		t.Error("input map was changed")
	}
	if src := ShowFields(op.KindSource, fields, handleOf); string(src["statement"]) != string(fields["statement"]) {
		t.Error("source fields must be shown as given")
	}
}

func TestShowFieldValueCanonical(t *testing.T) {
	live := "01M6HNNWFR1J2AHQRNKFTPME46"
	liveID, _ := id.Parse(live)
	handleOf := func(i id.ID) (string, bool) { return "renamed", i == liveID }
	text := "See @{" + live + "} <branch>, a & b, x > y   end"
	raw, err := canon.Marshal(text)
	if err != nil {
		t.Fatal(err)
	}
	raw = bytes.TrimSuffix(raw, []byte("\n"))
	got := ShowFieldValue(op.KindEvidence, "statement", raw, handleOf)
	want, _ := canon.Marshal("See @renamed <branch>, a & b, x > y   end")
	want = bytes.TrimSuffix(want, []byte("\n"))
	if string(got) != string(want) {
		t.Errorf("got %s, want %s", got, want)
	}
	if err := canon.Check(append(got[:len(got):len(got)], '\n')); err != nil {
		t.Errorf("not canonical: %v", err)
	}
}
