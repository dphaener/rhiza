package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Marks are the release marks of one change set, read from the notes that
// discuss it: a note whose text starts "breaking:" marks the change set
// breaking, one starting "severe:" marks a bug fix severe. Each holds the
// reasons, the text after the prefix. Other notes are ignored.
type Marks struct {
	Breaking []string `json:"breaking,omitempty"`
	Severe   []string `json:"severe,omitempty"`
}

// ReadMarks reads the log under dir (a .rhiza folder's log) and returns the
// marks by change set ID. Retired notes and retired links count for nothing,
// and so do ops an entry's body.drop retracts (a retire before landing).
func ReadMarks(logDir string) (map[string]Marks, error) {
	type rawOp struct {
		ID     string `json:"id"`
		Op     string `json:"op"`
		Kind   string `json:"kind"`
		Node   string `json:"node"`
		Link   string `json:"link"`
		From   string `json:"from"`
		To     string `json:"to"`
		Type   string `json:"type"`
		Fields struct {
			Statement string `json:"statement"`
		} `json:"fields"`
	}
	type entry struct {
		Body struct {
			Ops  []rawOp  `json:"ops"`
			Drop []string `json:"drop"`
		} `json:"body"`
	}
	files, err := filepath.Glob(filepath.Join(logDir, "*", "*.json"))
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	sets := map[string]bool{}   // change sets: the log's folders
	text := map[string]string{} // note node -> statement
	links := map[string]rawOp{} // link -> discusses link
	retired := map[string]bool{}
	var ops []rawOp
	dropped := map[string]bool{} // op IDs some entry's body.drop retracts
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		sets[filepath.Base(filepath.Dir(f))] = true
		var e entry
		if err := json.Unmarshal(b, &e); err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		ops = append(ops, e.Body.Ops...)
		for _, d := range e.Body.Drop {
			dropped[d] = true
		}
	}
	for _, o := range ops {
		if o.ID != "" && dropped[o.ID] {
			continue
		}
		switch {
		case o.Op == "node.create" && o.Kind == "note":
			text[o.Node] = o.Fields.Statement
		case o.Op == "link.create" && o.Type == "discusses":
			links[o.Link] = o
		case o.Op == "node.retire":
			retired[o.Node] = true
		case o.Op == "link.retire":
			retired[o.Link] = true
		}
	}
	out := map[string]Marks{}
	var ids []string
	for l := range links {
		ids = append(ids, l)
	}
	sort.Strings(ids)
	for _, l := range ids {
		o := links[l]
		st, ok := text[o.From]
		if !ok || !sets[o.To] || retired[l] || retired[o.From] {
			continue
		}
		m := out[o.To]
		switch prefix, reason, ok := mark(st); {
		case !ok:
			continue
		case prefix == "breaking":
			m.Breaking = append(m.Breaking, reason)
		default:
			m.Severe = append(m.Severe, reason)
		}
		out[o.To] = m
	}
	return out, nil
}

// mark splits a note's text into a release mark's prefix and reason.
func mark(s string) (prefix, reason string, ok bool) {
	s = strings.TrimSpace(s)
	for _, p := range []string{"breaking", "severe"} {
		if len(s) > len(p) && strings.EqualFold(s[:len(p)], p) && s[len(p)] == ':' {
			return p, strings.TrimSpace(s[len(p)+1:]), true
		}
	}
	return "", "", false
}

// IsBreaking reports whether the change set is marked breaking.
func (m Marks) IsBreaking() bool { return len(m.Breaking) > 0 }

// IsSevere reports whether the change set is marked a severe fix.
func (m Marks) IsSevere() bool { return len(m.Severe) > 0 }
