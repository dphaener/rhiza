package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Kind is the size of the next release.
type Kind string

const (
	KindNone  Kind = "none"
	KindPatch Kind = "patch"
	KindMinor Kind = "minor"
	KindMajor Kind = "major"
)

// Inputs is what decides the next version.
type Inputs struct {
	Breaking  []string // reasons: marked or detected breaking changes
	Milestone []string // milestones finished since the tag
	Fixes     []string // bug-fix change sets landed since the tag
	Landed    int      // change sets landed since the tag
}

// Next is the computed next release.
type Next struct {
	Tag     string   `json:"tag"`
	Version string   `json:"version,omitempty"`
	Kind    Kind     `json:"kind"`
	Reasons []string `json:"reasons,omitempty"`
}

// NextVersion computes the next version from the last v* tag ("" for none,
// read as v0.0.0). A breaking change makes a major release from v1.0.0 on and
// a minor one before; a finished milestone a minor; bug fixes alone a patch;
// anything else, or no change, none.
func NextVersion(tag string, in Inputs) (Next, error) {
	base := tag
	if base == "" {
		base = "v0.0.0"
	}
	m := tagRE.FindStringSubmatch(base)
	if m == nil {
		return Next{}, fmt.Errorf("tag %q is not vX.Y.Z", tag)
	}
	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	patch, _ := strconv.Atoi(m[3])
	n := Next{Tag: tag, Kind: KindNone}
	switch {
	case len(in.Breaking) > 0 && major >= 1:
		n.Kind = KindMajor
		n.Reasons = prefixed("breaking", in.Breaking)
	case len(in.Breaking) > 0:
		n.Kind = KindMinor
		n.Reasons = prefixed("breaking before v1.0.0 (minor)", in.Breaking)
	case len(in.Milestone) > 0:
		n.Kind = KindMinor
		n.Reasons = prefixed("milestone finished", in.Milestone)
	case len(in.Fixes) > 0:
		n.Kind = KindPatch
		n.Reasons = prefixed("bug fix", in.Fixes)
	}
	switch n.Kind {
	case KindMajor:
		n.Version = fmt.Sprintf("v%d.0.0", major+1)
	case KindMinor:
		n.Version = fmt.Sprintf("v%d.%d.0", major, minor+1)
	case KindPatch:
		n.Version = fmt.Sprintf("v%d.%d.%d", major, minor, patch+1)
	}
	return n, nil
}

func prefixed(p string, rs []string) []string {
	var out []string
	for _, r := range rs {
		out = append(out, p+": "+r)
	}
	return out
}

// Landing is a change set that landed.
type Landing struct {
	ID          string
	At          time.Time
	RequestKind string // "bug", "milestone" or "" for a feature
}

// History is what the log says about landings and milestones.
type History struct {
	Landed     map[string]Landing  // change set -> landing
	Milestones map[string][]string // intent handle -> change sets serving it
}

// ReadHistory reads the log under logDir. A milestone is an intent that
// change sets serve.
func ReadHistory(logDir string) (History, error) {
	files, err := filepath.Glob(filepath.Join(logDir, "*", "*.json"))
	if err != nil {
		return History{}, err
	}
	sort.Strings(files)
	h := History{Landed: map[string]Landing{}, Milestones: map[string][]string{}}
	kinds := map[string]string{}   // change set -> request kind
	intents := map[string]string{} // intent node -> handle
	type serve struct{ from, to string }
	var serves []serve
	sets := map[string]bool{}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return h, err
		}
		var e struct {
			At        string          `json:"at"`
			ChangeSet string          `json:"change_set"`
			Type      string          `json:"type"`
			Body      json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(b, &e); err != nil {
			return h, fmt.Errorf("%s: %w", f, err)
		}
		sets[e.ChangeSet] = true
		switch e.Type {
		case "open":
			var o struct {
				RequestKind string `json:"request_kind"`
			}
			json.Unmarshal(e.Body, &o)
			kinds[e.ChangeSet] = o.RequestKind
		case "land":
			at, err := time.Parse(time.RFC3339, e.At)
			if err != nil {
				return h, fmt.Errorf("%s: %w", f, err)
			}
			h.Landed[e.ChangeSet] = Landing{ID: e.ChangeSet, At: at}
		case "ops":
			var o struct {
				Ops []struct {
					Op     string `json:"op"`
					Kind   string `json:"kind"`
					Node   string `json:"node"`
					Handle string `json:"handle"`
					Type   string `json:"type"`
					From   string `json:"from"`
					To     string `json:"to"`
				} `json:"ops"`
			}
			json.Unmarshal(e.Body, &o)
			for _, op := range o.Ops {
				switch {
				case op.Op == "node.create" && op.Kind == "intent":
					intents[op.Node] = op.Handle
				case op.Op == "link.create" && op.Type == "serves":
					serves = append(serves, serve{op.From, op.To})
				}
			}
		}
	}
	for id, l := range h.Landed {
		l.RequestKind = kinds[id]
		h.Landed[id] = l
	}
	for _, s := range serves {
		if handle, ok := intents[s.to]; ok && sets[s.from] {
			h.Milestones[handle] = append(h.Milestones[handle], s.from)
		}
	}
	return h, nil
}

// Inputs derives the release inputs for the landings after since: bug fixes,
// milestones whose serving change sets have all landed with the last one
// after since, and the marks of the landed change sets.
func (h History) Inputs(since time.Time, marks map[string]Marks) Inputs {
	var in Inputs
	ids := make([]string, 0, len(h.Landed))
	for id, l := range h.Landed {
		if l.At.After(since) {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	for _, id := range ids {
		l := h.Landed[id]
		in.Landed++
		if l.RequestKind == "bug" {
			in.Fixes = append(in.Fixes, id)
		}
		for _, r := range marks[id].Breaking {
			in.Breaking = append(in.Breaking, id+": "+r)
		}
	}
	var handles []string
	for hd := range h.Milestones {
		handles = append(handles, hd)
	}
	sort.Strings(handles)
	for _, hd := range handles {
		done, last := true, time.Time{}
		for _, cs := range h.Milestones[hd] {
			l, ok := h.Landed[cs]
			if !ok {
				done = false
				break
			}
			if l.At.After(last) {
				last = l.At
			}
		}
		if done && last.After(since) {
			in.Milestone = append(in.Milestone, hd)
		}
	}
	return in
}

// runCheck prints the next version and its kind since the last v* tag in the
// repo at dir. It reads only: nothing is tagged or written.
func runCheck(dir string, w io.Writer) error {
	tag, err := gitOut(dir, "describe", "--tags", "--match", "v[0-9]*.[0-9]*.[0-9]*", "--abbrev=0")
	if err != nil {
		tag = ""
	}
	since := time.Time{}
	if tag != "" {
		d, err := gitOut(dir, "log", "-1", "--format=%cI", tag)
		if err != nil {
			return err
		}
		if since, err = time.Parse(time.RFC3339, d); err != nil {
			return err
		}
	}
	logDir := filepath.Join(dir, ".rhiza", "log")
	h, err := ReadHistory(logDir)
	if err != nil {
		return err
	}
	marks, err := ReadMarks(logDir)
	if err != nil {
		return err
	}
	in := h.Inputs(since, marks)
	var bs []Break
	if tag != "" {
		if _, bs, err = breakingAgainstLastTag(dir); err != nil {
			return err
		}
	}
	for _, b := range bs {
		in.Breaking = append(in.Breaking, b.String())
	}
	n, err := NextVersion(tag, in)
	if err != nil {
		return err
	}
	if n.Kind == KindNone {
		fmt.Fprintf(w, "none: no release since %s\n", orNone(tag))
		return nil
	}
	fmt.Fprintf(w, "%s %s (since %s)\n", n.Version, n.Kind, orNone(tag))
	for _, r := range n.Reasons {
		fmt.Fprintln(w, "  "+r)
	}
	return nil
}

func orNone(tag string) string {
	if tag == "" {
		return "no tag"
	}
	return tag
}

func gitOut(dir string, args ...string) (string, error) {
	c := exec.Command("git", args...)
	c.Dir = dir
	b, err := c.Output()
	return strings.TrimSpace(string(b)), err
}
