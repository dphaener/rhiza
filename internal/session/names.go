package session

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// Slug lowercases title, turns each run of characters outside [a-z0-9] into
// one dash and trims dashes. It is never empty.
func Slug(title string) string {
	var b strings.Builder
	dash := false
	for _, c := range strings.ToLower(title) {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			if dash && b.Len() > 0 {
				b.WriteByte('-')
			}
			dash = false
			b.WriteRune(c)
		} else {
			dash = true
		}
	}
	if b.Len() == 0 {
		return "session"
	}
	return b.String()
}

// Branch is the branch for a slug.
func Branch(slug string) string { return "cs/" + slug }

// Worktree is the sibling folder of repoRoot named <repo>.<slug>.
func Worktree(repoRoot, repo, slug string) string {
	return filepath.Join(filepath.Dir(repoRoot), repo+"."+slug)
}

// TabTitle is "<repo>: <title>", plus " · <step>" when step is set and is
// not build. Pass the change set's display title.
func TabTitle(repo, title, step string) string {
	tail := ""
	if step != "" && step != "build" {
		tail = " · " + step
	}
	return fitTitle(repo+": ", title, tail)
}

// RunTabTitle is the tab title of a run's session: "<repo>: <title> · run
// <run> · <step>", plus " · <label>" when label is set, such as "builder
// export-csv #2". Unlike TabTitle it keeps a build step's name, so every
// tab of a run names both the run and the step.
func RunTabTitle(repo, title, run, step, label string) string {
	tail := " · run " + run
	if step != "" {
		tail += " · " + step
	}
	if label != "" {
		tail += " · " + label
	}
	return fitTitle(repo+": ", title, tail)
}

// MaxTabTitle is the longest tab title Rex accepts, in characters.
const MaxTabTitle = 128

// fitTitle joins head, title and tail, shortening title with an ellipsis
// so the whole fits in MaxTabTitle characters. The head and tail are kept
// whole, as they tell a run's tabs apart; if they alone are too long, the
// end is cut.
func fitTitle(head, title, tail string) string {
	t := head + title + tail
	if utf8.RuneCountInString(t) <= MaxTabTitle {
		return t
	}
	room := MaxTabTitle - utf8.RuneCountInString(head+tail) - 1
	if room > 0 {
		r := []rune(title)
		if room < len(r) {
			return head + strings.TrimRight(string(r[:room]), " ") + "…" + tail
		}
	}
	r := []rune(t)
	return string(r[:MaxTabTitle-1]) + "…"
}

// NewID formats 16 bytes as a version-4, RFC 4122 variant UUID.
func NewID(entropy [16]byte) string {
	b := entropy
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
