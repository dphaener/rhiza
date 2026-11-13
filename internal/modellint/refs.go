// Package modellint holds deterministic lints over a change set's to-be
// model. It is pure: it reads a model and an fs.FS, runs no git and writes
// nothing.
package modellint

import (
	"encoding/json"
	"io/fs"
	"path"
	"regexp"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

// Kinds of Hit.
const (
	// KindPath is a repo path the head tree lacks.
	KindPath = "path"
	// KindTest is a Go test name no *_test.go declares.
	KindTest = "test"
	// KindScript is a testscript TestScripts/<name> with no <name>.txtar.
	KindScript = "script"
	// KindUnretired is a capability whose path references are all gone.
	KindUnretired = "unretired-capability"
	// KindUnapplied is a Decision no live node applies with decided-by.
	KindUnapplied = "unapplied-decision"
)

// Hit is one lint finding on a node.
type Hit struct {
	Node   id.ID  `json:"node"`
	Handle string `json:"handle"`
	// Field is "statement" or "acceptance".
	Field string `json:"field"`
	Kind  string `json:"kind"`
	// Ref is the reference as written, empty for KindUnretired.
	Ref string `json:"ref,omitempty"`
}

var (
	testRe   = regexp.MustCompile(`\bTest[A-Z0-9_][A-Za-z0-9_]*`)
	scriptRe = regexp.MustCompile(`\bTestScripts/([A-Za-z0-9_][A-Za-z0-9_.-]*)`)
	funcRe   = regexp.MustCompile(`(?m)^func (Test[A-Za-z0-9_]*)\(`)
	lineRe   = regexp.MustCompile(`:\d+(:\d+)?$`)
)

var fields = []string{"statement", "acceptance"}

// StaleRefs reports, for each node cs adds or changes that is live in m (the
// to-be model), every reference in its statement or acceptance that head
// lacks: a repo path, a Go test name, a testscript name. A live capability
// whose path references are all gone is also reported as KindUnretired. Hits
// are sorted by node, field, then position of the reference. Nodes cs did
// not write are never looked at.
func StaleRefs(m model.Model, cs changeset.ChangeSet, head fs.FS) []Hit {
	ids := writtenNodes(cs)

	c := &checker{head: head}
	var out []Hit
	for _, i := range ids {
		n, ok := m.Node(i)
		if !ok || n.Retired {
			continue
		}
		paths, gone := 0, 0
		var first Hit
		for _, f := range fields {
			text := fieldText(n, f)
			seen := map[string]bool{}
			for _, ref := range pathRefs(text) {
				if seen[ref] {
					continue
				}
				seen[ref] = true
				if !c.isPath(ref) {
					continue
				}
				paths++
				if c.exists(ref) {
					continue
				}
				gone++
				h := Hit{Node: n.ID, Handle: n.Handle, Field: f, Kind: KindPath, Ref: ref}
				if first.Kind == "" {
					first = h
				}
				out = append(out, h)
			}
			for _, name := range testNames(text) {
				if !c.hasTest(name) {
					out = append(out, Hit{Node: n.ID, Handle: n.Handle, Field: f, Kind: KindTest, Ref: name})
				}
			}
			for _, name := range scriptNames(text) {
				if !c.hasScript(name) {
					out = append(out, Hit{Node: n.ID, Handle: n.Handle, Field: f, Kind: KindScript, Ref: "TestScripts/" + name})
				}
			}
		}
		if n.Kind == op.KindCapability && paths > 0 && gone == paths {
			out = append(out, Hit{Node: n.ID, Handle: n.Handle, Field: first.Field, Kind: KindUnretired})
		}
	}
	return out
}

// writtenNodes returns the nodes cs creates, updates or moves, sorted by ID.
func writtenNodes(cs changeset.ChangeSet) []id.ID {
	written := map[id.ID]bool{}
	for _, o := range cs.Ops {
		switch o.Type {
		case op.TypeNodeCreate, op.TypeNodeUpdate, op.TypeNodeMove:
			written[o.Node] = true
		}
	}
	ids := make([]id.ID, 0, len(written))
	for i := range written {
		ids = append(ids, i)
	}
	slices.SortFunc(ids, id.Compare)
	return ids
}

func fieldText(n model.Node, field string) string {
	raw, ok := n.Fields[field]
	if !ok {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) != nil {
		return ""
	}
	return s
}

// pathRefs returns the repo path candidates in text, in order: tokens whose
// first segment, after stripping a leading ./ and a trailing /..., is a
// top-level entry is decided by the caller; here only the token shape is.
func pathRefs(text string) []string {
	var out []string
	for _, tok := range strings.FieldsFunc(text, isSep) {
		tok = strings.Trim(tok, ".,;:!?'\"`*")
		tok = lineRe.ReplaceAllString(tok, "")
		tok = strings.TrimPrefix(tok, "./")
		tok = strings.TrimSuffix(tok, "/...")
		tok = strings.TrimSuffix(tok, "/")
		if tok == "" || strings.ContainsAny(tok, "*<>{}$|\\@") || strings.Contains(tok, "://") {
			continue
		}
		out = append(out, tok)
	}
	return out
}

func isSep(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r', '(', ')', '[', ']', '=', ',':
		return true
	}
	return false
}

func testNames(text string) []string {
	var out []string
	for _, m := range testRe.FindAllString(text, -1) {
		if !slices.Contains(out, m) {
			out = append(out, m)
		}
	}
	return out
}

func scriptNames(text string) []string {
	var out []string
	for _, m := range scriptRe.FindAllStringSubmatch(text, -1) {
		name := strings.TrimRight(m[1], ".-")
		if name != "" && !slices.Contains(out, name) {
			out = append(out, name)
		}
	}
	return out
}

// checker answers the head-tree questions, each scan done once and lazily.
type checker struct {
	head    fs.FS
	top     map[string]bool
	topDone bool
	scanned bool
	tests   map[string]bool
	scripts map[string]bool
}

// isPath reports whether ref is a path reference: a valid path whose first
// segment is a top-level entry of head.
func (c *checker) isPath(ref string) bool {
	if !c.topDone {
		c.topDone = true
		c.top = map[string]bool{}
		entries, _ := fs.ReadDir(c.head, ".")
		for _, e := range entries {
			c.top[e.Name()] = true
		}
	}
	first, _, _ := strings.Cut(ref, "/")
	return c.top[first] && fs.ValidPath(ref)
}

// exists reports whether the path reference ref is in head.
func (c *checker) exists(ref string) bool {
	_, err := fs.Stat(c.head, path.Clean(ref))
	return err == nil
}

func (c *checker) scan() {
	if c.scanned {
		return
	}
	c.scanned = true
	c.tests = map[string]bool{}
	c.scripts = map[string]bool{}
	_ = fs.WalkDir(c.head, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if p != "." && (d.Name() == ".git" || d.Name() == "node_modules") {
				return fs.SkipDir
			}
			return nil
		}
		switch {
		case strings.HasSuffix(p, ".txtar"):
			c.scripts[strings.TrimSuffix(path.Base(p), ".txtar")] = true
		case strings.HasSuffix(p, "_test.go"):
			b, err := fs.ReadFile(c.head, p)
			if err != nil {
				return nil
			}
			for _, m := range funcRe.FindAllSubmatch(b, -1) {
				c.tests[string(m[1])] = true
			}
		}
		return nil
	})
}

func (c *checker) hasTest(name string) bool   { c.scan(); return c.tests[name] }
func (c *checker) hasScript(name string) bool { c.scan(); return c.scripts[name] }
