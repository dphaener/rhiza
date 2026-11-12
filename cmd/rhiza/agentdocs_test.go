package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	mdLink   = regexp.MustCompile(`\]\(([^)\s]+)\)`)
	codeSpan = regexp.MustCompile("`[^`]*`")
)

// checkAgentDocs returns the problems in the agent docs under root: an
// agents/*.md missing from agents/README.md's index, an AGENTS.md without a
// link to agents/README.md, and relative links that point at nothing.
func checkAgentDocs(root string) []string {
	var problems []string
	read := func(rel string) string {
		b, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			problems = append(problems, "cannot read "+rel+": "+err.Error())
		}
		return string(b)
	}
	agents := read("AGENTS.md")
	index := read("agents/README.md")

	files := []string{"AGENTS.md"}
	docs, _ := filepath.Glob(filepath.Join(root, "agents", "*.md"))
	for _, d := range docs {
		name := filepath.Base(d)
		files = append(files, "agents/"+name)
		if name != "README.md" && !strings.Contains(index, "]("+name+")") {
			problems = append(problems, "agents/README.md does not index "+name)
		}
	}

	linksIndex := false
	for _, l := range mdLink.FindAllStringSubmatch(agents, -1) {
		if p, _, _ := strings.Cut(l[1], "#"); p == "agents/README.md" {
			linksIndex = true
		}
	}
	if !linksIndex {
		problems = append(problems, "AGENTS.md does not link agents/README.md")
	}

	for _, f := range files {
		b, err := os.ReadFile(filepath.Join(root, f))
		if err != nil {
			continue
		}
		for _, l := range mdLink.FindAllStringSubmatch(codeSpan.ReplaceAllString(string(b), ""), -1) {
			target, _, _ := strings.Cut(l[1], "#")
			if target == "" || strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") {
				continue
			}
			if _, err := os.Stat(filepath.Join(root, filepath.Dir(f), target)); err != nil {
				problems = append(problems, f+" links missing "+l[1])
			}
		}
	}
	return problems
}

func TestAgentDocsConsistent(t *testing.T) {
	if p := checkAgentDocs("../.."); len(p) > 0 {
		t.Errorf("agent docs inconsistent:\n%s", strings.Join(p, "\n"))
	}
}

func TestAgentDocsConsistentFixtures(t *testing.T) {
	good := map[string]string{
		"AGENTS.md":        "See [index](agents/README.md).\n",
		"agents/README.md": "- [a.md](a.md): x\n[up](../AGENTS.md)\n",
		"agents/a.md":      "# a\n",
	}
	cases := []struct {
		name   string
		mutate func(map[string]string)
		want   string
	}{
		{"ok", func(map[string]string) {}, ""},
		{"unindexed doc", func(m map[string]string) { m["agents/b.md"] = "# b\n" }, "does not index b.md"},
		{"no AGENTS link", func(m map[string]string) { m["AGENTS.md"] = "nothing\n" }, "does not link agents/README.md"},
		{"broken link", func(m map[string]string) { m["agents/a.md"] = "[x](gone.md)\n" }, "links missing gone.md"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := map[string]string{}
			for k, v := range good {
				m[k] = v
			}
			c.mutate(m)
			root := t.TempDir()
			for k, v := range m {
				p := filepath.Join(root, k)
				if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(p, []byte(v), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			got := strings.Join(checkAgentDocs(root), "\n")
			if c.want == "" && got != "" || !strings.Contains(got, c.want) {
				t.Errorf("problems = %q, want %q", got, c.want)
			}
		})
	}
}
