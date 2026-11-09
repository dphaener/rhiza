// Package textref finds the references to other nodes inside free text and
// turns them both ways.
//
// A reference is written @handle, @<ID or 4+ character ID prefix> or
// @{<ID>}. It starts with @ at the start of the text or after whitespace or an
// opening ( or [, so ada@example.com, checkout@v7 and a quoted '@key' hold none.
// @@ at such a position is an escape for a literal @ and is never a reference.
//
// Stored text holds each reference as @{<full 26-character ID>}. Only that
// form counts as a stored reference; any other @{...} in stored text is
// literal. Format turns stored text back into what a reader sees: @<handle>
// for a live node, @{<full ID>} otherwise, and a literal @ that would start a
// reference is shown as @@, so shown text written back stores the same text.
//
// The package is pure: it knows nothing of the model. Callers supply lookups.
package textref

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Kind says how a reference was spelled.
type Kind int

const (
	// KindHandle is a bare lower-case handle: words of [a-z0-9] joined by single
	// dashes. A token of digits only is a handle too; a lookup may also try it
	// as an ID prefix.
	KindHandle Kind = iota + 1
	// KindID is an ID or an ID prefix of 4 or more characters (upper-case
	// Crockford base32), bare or inside @{ }. Anything inside @{ } is KindID.
	KindID
)

// Ref is one reference found in a text.
type Ref struct {
	Start, End int    // byte offsets of the whole reference, @ included
	Token      string // the handle, or the ID or prefix, without @ or braces
	Kind       Kind
	Braced     bool
	// Escaped marks an @@{<full ID>} escape Resolve returns: its literal @
	// would be stored as @{<full ID>}, which reads back as a reference.
	Escaped bool
}

// IDLen is the length of a full ID.
const IDLen = 26

const crockford = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

type item struct {
	Ref
	escape bool
}

// startOK reports whether a reference may start at byte i of text.
func startOK(text string, i int) bool {
	if i == 0 {
		return true
	}
	r, _ := utf8.DecodeLastRuneInString(text[:i])
	return r == '(' || r == '[' || unicode.IsSpace(r)
}

func isWordByte(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_' || c == '-'
}

func isLowerAlnum(c byte) bool { return c >= '0' && c <= '9' || c >= 'a' && c <= 'z' }

func isCrockford(c byte) bool { return strings.IndexByte(crockford, c) >= 0 }

// parseAt reads what text[i:] (text[i] == '@', a start position) holds.
func parseAt(text string, i int) (item, bool) {
	rest := text[i+1:]
	if rest == "" {
		return item{}, false
	}
	switch rest[0] {
	case '@':
		return item{Ref: Ref{Start: i, End: i + 2}, escape: true}, true
	case '{':
		j := strings.IndexByte(rest, '}')
		if j < 2 {
			return item{}, false
		}
		tok := rest[1:j]
		for k := 0; k < len(tok); k++ {
			if !isWordByte(tok[k]) || tok[k] == '_' {
				return item{}, false
			}
		}
		return item{Ref: Ref{Start: i, End: i + 1 + j + 1, Token: tok, Kind: KindID, Braced: true}}, true
	}
	// An ID or prefix: upper-case Crockford with at least one letter.
	n := 0
	for n < len(rest) && (rest[n] >= '0' && rest[n] <= '9' || rest[n] >= 'A' && rest[n] <= 'Z') {
		n++
	}
	if n >= 4 && n <= IDLen && (n == len(rest) || !isWordByte(rest[n])) {
		tok := rest[:n]
		letters, ok := false, true
		for k := 0; k < n; k++ {
			if !isCrockford(tok[k]) {
				ok = false
			}
			if tok[k] > '9' {
				letters = true
			}
		}
		if ok && letters {
			return item{Ref: Ref{Start: i, End: i + 1 + n, Token: tok, Kind: KindID}}, true
		}
	}
	// A handle: [a-z0-9]+ words joined by single dashes.
	n = 0
	for {
		m := n
		for m < len(rest) && isLowerAlnum(rest[m]) {
			m++
		}
		if m == n {
			break
		}
		n = m
		if n+1 < len(rest) && rest[n] == '-' && isLowerAlnum(rest[n+1]) {
			n++
			continue
		}
		break
	}
	if n == 0 || (n < len(rest) && isWordByte(rest[n])) {
		return item{}, false
	}
	return item{Ref: Ref{Start: i, End: i + 1 + n, Token: rest[:n], Kind: KindHandle}}, true
}

func scan(text string) []item {
	var out []item
	for i := 0; i < len(text); {
		if text[i] == '@' && startOK(text, i) {
			if it, ok := parseAt(text, i); ok {
				out = append(out, it)
				i = it.End
				continue
			}
		}
		i++
	}
	return out
}

// Find returns the references in text, in order. Escapes (@@) are not
// references.
func Find(text string) []Ref {
	var out []Ref
	for _, it := range scan(text) {
		if !it.escape {
			out = append(out, it.Ref)
		}
	}
	return out
}

// Resolve rewrites every reference in text to @{<full ID>} and every @@ to @.
// lookup maps a reference to the full ID of the node it names. A reference
// lookup cannot map is left as written and returned, with its offset in
// Ref.Start, so the caller can refuse the write.
//
// Stored text cannot hold a literal @{<full ID>}: it would read back as a
// reference. So an @@{<full ID>} escape is returned too, with Ref.Escaped
// set, and both it and an unmapped @{<full ID>} are written as @<full ID>,
// without braces, so the text Resolve returns never holds a stored reference
// that lookup did not give.
func Resolve(text string, lookup func(Ref) (string, bool)) (string, []Ref) {
	var (
		b          strings.Builder
		unresolved []Ref
		last       int
	)
	for _, it := range scan(text) {
		b.WriteString(text[last:it.Start])
		last = it.End
		switch {
		case it.escape:
			if id, ok := storedRef(text, it.Start+1); ok {
				end := it.Start + 1 + 2 + IDLen + 1
				b.WriteString("@" + id)
				unresolved = append(unresolved, Ref{Start: it.Start, End: end, Token: id, Kind: KindID, Braced: true, Escaped: true})
				last = end
				continue
			}
			b.WriteByte('@')
		default:
			if full, ok := lookup(it.Ref); ok {
				b.WriteString("@{" + full + "}")
				continue
			}
			if id, ok := storedRef(text, it.Start); ok {
				b.WriteString("@" + id)
			} else {
				b.WriteString(text[it.Start:it.End])
			}
			unresolved = append(unresolved, it.Ref)
		}
	}
	b.WriteString(text[last:])
	return b.String(), unresolved
}

// storedRef reads a stored reference, @{<26-character ID>}, at text[i:].
func storedRef(text string, i int) (string, bool) {
	const n = 2 + IDLen + 1
	if len(text)-i < n || text[i:i+2] != "@{" || text[i+n-1] != '}' {
		return "", false
	}
	id := text[i+2 : i+2+IDLen]
	for k := 0; k < len(id); k++ {
		if !isCrockford(id[k]) || id[k] >= 'a' {
			return "", false
		}
	}
	return id, true
}

func validHandle(h string) bool {
	if h == "" || h[0] == '-' || h[len(h)-1] == '-' {
		return false
	}
	prev := byte('-')
	for k := 0; k < len(h); k++ {
		c := h[k]
		if c == '-' && prev == '-' || c != '-' && !isLowerAlnum(c) {
			return false
		}
		prev = c
	}
	return true
}

// Format turns stored text into shown text. handleOf gives the current handle
// of a live node by full ID; it reports false for a retired or missing node,
// which is shown as @{<full ID>}. Any other literal @ that would start a
// reference, and any other @{ such as '@{<full ID>}' or '@{...}', is shown as
// @@, so Resolve of the shown text stores the same text.
func Format(text string, handleOf func(id string) (string, bool)) string {
	var b strings.Builder
	for i := 0; i < len(text); {
		if text[i] != '@' || !startOK(text, i) {
			b.WriteByte(text[i])
			i++
			continue
		}
		if id, ok := storedRef(text, i); ok {
			end := i + 2 + IDLen + 1
			h, live := handleOf(id)
			bare := live && validHandle(h) && !allDigits(h)
			if bare && end < len(text) {
				r, _ := utf8.DecodeRuneInString(text[end:])
				bare = unicode.IsSpace(r) || strings.ContainsRune(".,;:!?)]", r)
			}
			if bare {
				b.WriteString("@" + h)
			} else {
				b.WriteString(text[i:end])
			}
			i = end
			continue
		}
		// Any other @{ is literal text: escape it whether or not parseAt would
		// take it, so stored '@{<full ID>}' is shown as '@@{<full ID>}'.
		if _, ok := parseAt(text, i); ok || strings.HasPrefix(text[i:], "@{") {
			b.WriteString("@@")
		} else {
			b.WriteByte('@')
		}
		i++
	}
	return b.String()
}

func allDigits(s string) bool {
	for k := 0; k < len(s); k++ {
		if s[k] < '0' || s[k] > '9' {
			return false
		}
	}
	return true
}
