package changeset

// This file holds a change set's title: the short name an open body may
// carry, and the display title every change set has, its title or else the
// first clause of its intent. One display title names the change set's
// branch, worktree and terminal tab.

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// MaxTitle is the most runes a title may have.
const MaxTitle = 30

// clauseBreaks end an intent's first clause wherever they appear. A '.' or
// ',' ends it too, but only before a space or the end (see firstClause).
var clauseBreaks = []string{":", ";", " — ", " - ", " (", "\n", "\r"}

// validTitle checks an open body's title: empty (none), or valid UTF-8 on
// one line, with no control characters, no space at either end and at most
// MaxTitle runes.
func validTitle(s string) error {
	if s == "" {
		return nil
	}
	if !utf8.ValidString(s) {
		return invalid("title", "invalid UTF-8")
	}
	if strings.IndexFunc(s, unicode.IsControl) >= 0 {
		return invalid("title", "not a single line: it holds a line break or control character")
	}
	first, _ := utf8.DecodeRuneInString(s)
	last, _ := utf8.DecodeLastRuneInString(s)
	if unicode.IsSpace(first) || unicode.IsSpace(last) {
		return invalid("title", "leading or trailing space")
	}
	if n := utf8.RuneCountInString(s); n > MaxTitle {
		return invalid("title", "%d characters, at most %d", n, MaxTitle)
	}
	return nil
}

// DisplayTitle returns the change set's title when it has one, else the
// first clause of its intent: the text before the first ':', ';', " — ",
// " - ", " (" or line break, or a '.' or ',' followed by a space or the end,
// trimmed and cut to MaxTitle runes at the last word boundary that fits
// (hard-cut when there is none). When the clause comes out empty it is the
// intent's first MaxTitle runes instead, so it is empty only for an empty
// intent.
func (o Open) DisplayTitle() string {
	if o.Title != "" {
		return o.Title
	}
	if t := shorten(firstClause(o.Intent)); t != "" {
		return t
	}
	if t := strings.TrimSpace(headRunes(strings.TrimSpace(o.Intent), MaxTitle)); t != "" {
		return t
	}
	return headRunes(o.Intent, MaxTitle)
}

// firstClause returns s up to its first clause break, trimmed.
func firstClause(s string) string {
	end := len(s)
	for _, b := range clauseBreaks {
		if i := strings.Index(s, b); i >= 0 && i < end {
			end = i
		}
	}
	for i := 0; i < end; i++ {
		if c := s[i]; c == '.' || c == ',' {
			if i+1 == len(s) || unicode.IsSpace(rune(s[i+1])) {
				end = i
				break
			}
		}
	}
	return strings.TrimSpace(s[:end])
}

// shorten cuts s to at most MaxTitle runes: at the last space that leaves
// at most MaxTitle runes before it, else hard at MaxTitle runes.
func shorten(s string) string {
	if utf8.RuneCountInString(s) <= MaxTitle {
		return s
	}
	cut, runes := -1, 0
	for i, r := range s {
		if runes > MaxTitle {
			break
		}
		if unicode.IsSpace(r) {
			cut = i
		}
		runes++
	}
	if cut > 0 {
		return strings.TrimRightFunc(s[:cut], unicode.IsSpace)
	}
	return headRunes(s, MaxTitle)
}

// headRunes returns the first n runes of s.
func headRunes(s string, n int) string {
	for i := range s {
		if n == 0 {
			return s[:i]
		}
		n--
	}
	return s
}
