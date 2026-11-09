// Package doclinks checks that relative links in the Markdown docs resolve:
// the target file exists and a #fragment names a heading in it.
package doclinks

import (
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	linkRe    = regexp.MustCompile(`!?\[[^\]]*\]\(([^)\s]+)(?:\s+"[^"]*")?\)`)
	headingRe = regexp.MustCompile(`^ {0,3}#{1,6}\s+(.*?)\s*#*\s*$`)
	slugDrop  = regexp.MustCompile(`[^\p{L}\p{N}\s_-]`)
	inlineRe  = regexp.MustCompile("`([^`]*)`")
	mdLinkRe  = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
)

// Check walks README.md and docs/ under root and returns one problem per
// broken relative link, sorted. External links (with a scheme) are skipped.
func Check(root string) ([]string, error) {
	var files []string
	if _, err := os.Stat(filepath.Join(root, "README.md")); err == nil {
		files = append(files, "README.md")
	}
	docs := filepath.Join(root, "docs")
	if _, err := os.Stat(docs); err == nil {
		err := filepath.WalkDir(docs, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(p, ".md") {
				rel, _ := filepath.Rel(root, p)
				files = append(files, rel)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	var problems []string
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(root, f))
		if err != nil {
			return nil, err
		}
		inFence := false
		for i, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "```") {
				inFence = !inFence
				continue
			}
			if inFence {
				continue
			}
			for _, m := range linkRe.FindAllStringSubmatch(inlineRe.ReplaceAllString(line, ""), -1) {
				if msg := checkLink(root, f, m[1]); msg != "" {
					problems = append(problems, fmt.Sprintf("%s:%d: %s", f, i+1, msg))
				}
			}
		}
	}
	sort.Strings(problems)
	return problems, nil
}

func checkLink(root, from, target string) string {
	if u, err := url.Parse(target); err == nil && (u.Scheme != "" || u.Host != "") {
		return ""
	}
	path, frag, _ := strings.Cut(target, "#")
	file := from
	if path != "" {
		if un, err := url.PathUnescape(path); err == nil {
			path = un
		}
		file = filepath.Join(filepath.Dir(from), path)
		if filepath.IsAbs(path) {
			file = strings.TrimPrefix(path, "/")
		}
		file = filepath.Clean(file)
		if file == ".." || strings.HasPrefix(file, "../") {
			return fmt.Sprintf("link %q leaves the repo", target)
		}
		if _, err := os.Stat(filepath.Join(root, file)); err != nil {
			return fmt.Sprintf("link %q: %s does not exist", target, file)
		}
	}
	if frag == "" || !strings.HasSuffix(file, ".md") {
		return ""
	}
	if un, err := url.PathUnescape(frag); err == nil {
		frag = un
	}
	anchors, err := anchorsOf(filepath.Join(root, file))
	if err != nil {
		return fmt.Sprintf("link %q: %v", target, err)
	}
	if !anchors[strings.ToLower(frag)] {
		return fmt.Sprintf("link %q: no heading #%s in %s", target, frag, file)
	}
	return ""
}

// anchorsOf returns the GitHub-style heading slugs of a Markdown file.
func anchorsOf(path string) (map[string]bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	out := map[string]bool{}
	seen := map[string]int{}
	inFence := false
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		m := headingRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		text := mdLinkRe.ReplaceAllString(m[1], "$1")
		text = strings.ReplaceAll(text, "`", "")
		slug := strings.ReplaceAll(slugDrop.ReplaceAllString(strings.ToLower(text), ""), " ", "-")
		n := seen[slug]
		seen[slug]++
		if n > 0 {
			slug = fmt.Sprintf("%s-%d", slug, n)
		}
		out[slug] = true
	}
	return out, nil
}
