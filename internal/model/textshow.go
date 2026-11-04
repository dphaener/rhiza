package model

import (
	"encoding/json"
	"maps"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/textref"
)

// verbatimFields are stored, and so shown, exactly as given: an owner's words.
var verbatimFields = map[string]bool{"quote": true, "quote_where": true, "wording": true}

// TextVerbatim reports whether field name of a node of kind is stored and
// shown as given, its @ never a reference: every field of a source, and the
// owner's words of a decision.
func TextVerbatim(kind op.Kind, name string) bool {
	return kind == op.KindSource || verbatimFields[name]
}

// HandleOf gives the current handle of the live node with ID i, and false
// for a retired or unknown node.
type HandleOf func(i id.ID) (string, bool)

// ShowText turns one stored text into the text a reader sees, as
// textref.Format does, with handleOf over IDs.
func ShowText(text string, handleOf HandleOf) string {
	if !strings.Contains(text, "@") {
		return text
	}
	return textref.Format(text, func(s string) (string, bool) {
		i, err := id.Parse(s)
		if err != nil {
			return "", false
		}
		return handleOf(i)
	})
}

// ShowFieldValue is ShowText of a JSON string value; any other value is
// returned as given.
func ShowFieldValue(kind op.Kind, name string, raw json.RawMessage, handleOf HandleOf) json.RawMessage {
	if TextVerbatim(kind, name) || len(raw) == 0 || raw[0] != '"' {
		return raw
	}
	var s string
	if json.Unmarshal(raw, &s) != nil || !strings.Contains(s, "@") {
		return raw
	}
	shown := ShowText(s, handleOf)
	if shown == s {
		return raw
	}
	b, err := json.Marshal(shown)
	if err != nil {
		return raw
	}
	return b
}

// ShowFields returns fields, those of a node of kind, with each text shown by
// current handle. The map given is not changed.
func ShowFields(kind op.Kind, fields map[string]json.RawMessage, handleOf HandleOf) map[string]json.RawMessage {
	if fields == nil || kind == op.KindSource {
		return fields
	}
	var out map[string]json.RawMessage
	for name, raw := range fields {
		shown := ShowFieldValue(kind, name, raw, handleOf)
		if len(shown) != len(raw) || string(shown) != string(raw) {
			if out == nil {
				out = maps.Clone(fields)
			}
			out[name] = shown
		}
	}
	if out == nil {
		return fields
	}
	return out
}
