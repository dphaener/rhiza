package session

import (
	"fmt"
	"path/filepath"
	"strings"
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
	t := repo + ": " + title
	if step != "" && step != "build" {
		t += " · " + step
	}
	return t
}

// NewID formats 16 bytes as a version-4, RFC 4122 variant UUID.
func NewID(entropy [16]byte) string {
	b := entropy
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
