package op

import "slices"

// knownFields lists, per node kind, every field a node of that kind may
// carry, sorted. It covers each field in use in the log today. Anchors carry
// no fields yet.
var knownFields = map[Kind][]string{
	KindIntent:     {"milestone", "statement", "umbrella"},
	KindCapability: {"acceptance", "built_in", "statement"},
	KindConstraint: {"acceptance", "statement"},
	KindDecision: {"accepted_by", "authorized_by", "notes_quote", "options", "picked", "question", "quote", "quote_where",
		"reason", "reasoning", "refines", "revised", "ruling", "session", "settled_by", "statement", "wording"},
	KindQuestion: {"session", "statement"},
	KindEvidence: {"command", "head", "observed", "pr", "result", "statement", "status", "test", "tests"},
	KindSource:   {"quotes", "raw_sha256", "session", "summary", "title"},
	KindActor:    {"agent", "email", "for", "model", "name", "role", "session"},
	KindAnchor:   {},
	KindNote:     {"op", "statement"},
}

// KnownFields returns the sorted field names a node of kind k may carry, or
// nil for an unknown kind. The caller may keep the result.
func KnownFields(k Kind) []string {
	return slices.Clone(knownFields[k])
}

// FieldKnown reports whether a node of kind k may carry a field called name.
func FieldKnown(k Kind, name string) bool {
	return slices.Contains(knownFields[k], name)
}
