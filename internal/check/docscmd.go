package check

// This file holds the docs-command rule: README.md and docs/ may only name
// rhiza commands the registry has.

import (
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strings"
)

// docsCommandRe finds `rhiza ...` in a line; the text before it must not be
// part of a word or path (.rhiza/, rhiza.toml, github.com/x/rhiza).
var docsCommandRe = regexp.MustCompile("(?:^|[^A-Za-z0-9_./-])rhiza ([^`]*)")

// docsCommands fails on each rhiza command that README.md or a markdown
// file under docs/ names, in backticks or a code block, and the registry
// does not have. Input.Commands nil skips the rule.
func (c *checker) docsCommands() {
	if c.in.Commands == nil || c.in.Tree == nil {
		return
	}
	cmds := map[string]bool{}
	prefixes := map[string]bool{}
	for _, n := range c.in.Commands {
		cmds[n] = true
		parts := strings.Fields(n)
		for i := range parts {
			prefixes[strings.Join(parts[:i+1], " ")] = true
		}
	}
	var files []string
	if _, err := fs.Stat(c.in.Tree, "README.md"); err == nil {
		files = append(files, "README.md")
	}
	_ = fs.WalkDir(c.in.Tree, "docs", func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && path.Ext(p) == ".md" {
			files = append(files, p)
		}
		return nil
	})
	sort.Strings(files)
	for _, f := range files {
		b, err := fs.ReadFile(c.in.Tree, f)
		if err != nil {
			continue
		}
		seen := map[string]bool{}
		fenced := false
		for _, line := range strings.Split(string(b), "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "```") {
				fenced = !fenced
				continue
			}
			for _, text := range docsCodeSpans(line, fenced) {
				for _, m := range docsCommandRe.FindAllStringSubmatch(text, -1) {
					bad := badCommand(m[1], cmds, prefixes)
					if bad == "" || seen[bad] {
						continue
					}
					seen[bad] = true
					c.add(Finding{Severity: SeverityFail, Rule: RuleDocsCommand, Path: f,
						Message: "names `rhiza " + bad + "`, which is not a command", Next: "rhiza describe"})
				}
			}
		}
	}
}

// docsCodeSpans is the code in a line: the whole line inside a fence,
// else each backtick span.
func docsCodeSpans(line string, fenced bool) []string {
	if fenced {
		return []string{line}
	}
	var out []string
	parts := strings.Split(line, "`")
	for i := 1; i < len(parts)-1; i += 2 {
		out = append(out, parts[i])
	}
	return out
}

var wordRe = regexp.MustCompile(`^[a-z][a-z-]*$`)

// badCommand walks the words after "rhiza" and returns the words that name
// no command, "" when they do (arguments after a command are ignored).
func badCommand(rest string, cmds, prefixes map[string]bool) string {
	cur := ""
	for _, w := range strings.Fields(rest) {
		if !wordRe.MatchString(w) {
			return ""
		}
		next := w
		if cur != "" {
			next = cur + " " + w
		}
		if !prefixes[next] {
			return next
		}
		if cmds[next] {
			return ""
		}
		cur = next
	}
	return ""
}
