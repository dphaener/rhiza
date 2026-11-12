package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
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
// a minor one before; a finished milestone a minor; otherwise any landed bug
// fix a patch, which ships whatever is on main beside it (an unfinished
// milestone's features included); anything else, or no change, none.
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
	Intent      string // what the change set set out to do
	Kind        string // the change set's kind: "intention", or "record" for model bookkeeping
	Entry       string // the land entry's path under the log directory, slash-separated
}

// History is what the log says about landings and milestones.
type History struct {
	Landed     map[string]Landing  // change set -> landing
	Milestones map[string][]string // intent handle -> change sets serving it
	Statements map[string]string   // intent handle -> its statement
	Bugs       []time.Time         // when each bug request was filed (its open)
}

// ReadHistory reads the log under logDir. A milestone is an intent with a
// non-empty milestone field; a change set serves it when it serves the
// milestone or any intent under it (a done point, a product intent). A change
// set that ended without landing (withdrawn, rejected or superseded) serves
// nothing, and neither does a serves link that was retired or that a
// body.drop retracts.
func ReadHistory(logDir string) (History, error) {
	files, err := filepath.Glob(filepath.Join(logDir, "*", "*.json"))
	if err != nil {
		return History{}, err
	}
	sort.Strings(files)
	h := History{Landed: map[string]Landing{}, Milestones: map[string][]string{}, Statements: map[string]string{}}
	kinds := map[string]string{}   // change set -> request kind
	intents := map[string]string{} // intent node -> handle
	parents := map[string]string{} // node -> parent node
	isMilestone := map[string]bool{}
	texts := map[string]string{}   // change set -> intent text
	csKinds := map[string]string{} // change set -> change set kind
	type serve struct{ id, link, from, to string }
	var serves []serve
	sets := map[string]bool{}
	ended := map[string]bool{} // change set -> ended without landing
	type retire struct{ id, link string }
	var retires []retire
	dropped := map[string]bool{} // op IDs some entry's body.drop retracts
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
				Intent      string `json:"intent"`
				Kind        string `json:"kind"`
			}
			json.Unmarshal(e.Body, &o)
			kinds[e.ChangeSet] = o.RequestKind
			csKinds[e.ChangeSet] = o.Kind
			texts[e.ChangeSet] = o.Intent
			if o.RequestKind == "bug" {
				if at, err := time.Parse(time.RFC3339, e.At); err == nil {
					h.Bugs = append(h.Bugs, at)
				}
			}
		case "land":
			at, err := time.Parse(time.RFC3339, e.At)
			if err != nil {
				return h, fmt.Errorf("%s: %w", f, err)
			}
			rel, err := filepath.Rel(logDir, f)
			if err != nil {
				return h, err
			}
			h.Landed[e.ChangeSet] = Landing{ID: e.ChangeSet, At: at, Entry: filepath.ToSlash(rel)}
		case "withdraw", "reject", "supersede":
			ended[e.ChangeSet] = true
		case "reopen":
			ended[e.ChangeSet] = false
		case "ops":
			var o struct {
				Drop []string `json:"drop"`
				Ops  []struct {
					ID     string                     `json:"id"`
					Op     string                     `json:"op"`
					Link   string                     `json:"link"`
					Kind   string                     `json:"kind"`
					Node   string                     `json:"node"`
					Handle string                     `json:"handle"`
					Parent *string                    `json:"parent"`
					Fields map[string]json.RawMessage `json:"fields"`
					Set    map[string]json.RawMessage `json:"set"`
					Unset  []string                   `json:"unset"`
					Type   string                     `json:"type"`
					From   string                     `json:"from"`
					To     string                     `json:"to"`
				} `json:"ops"`
			}
			json.Unmarshal(e.Body, &o)
			for _, d := range o.Drop {
				dropped[d] = true
			}
			for _, op := range o.Ops {
				switch {
				case op.Op == "node.create" && op.Kind == "intent":
					intents[op.Node] = op.Handle
					if op.Parent != nil {
						parents[op.Node] = *op.Parent
					}
					isMilestone[op.Node] = fieldSet(op.Fields["milestone"])
					var st string
					json.Unmarshal(op.Fields["statement"], &st)
					h.Statements[op.Handle] = st
				case op.Op == "node.update":
					if v, ok := op.Set["milestone"]; ok {
						isMilestone[op.Node] = fieldSet(v)
					}
					if slices.Contains(op.Unset, "milestone") {
						isMilestone[op.Node] = false
					}
				case op.Op == "node.move":
					if op.Parent != nil {
						parents[op.Node] = *op.Parent
					} else {
						delete(parents, op.Node)
					}
				case op.Op == "link.create" && op.Type == "serves":
					serves = append(serves, serve{op.ID, op.Link, op.From, op.To})
				case op.Op == "link.retire":
					retires = append(retires, retire{op.ID, op.Link})
				}
			}
		}
	}
	for id, l := range h.Landed {
		l.RequestKind = kinds[id]
		l.Intent = texts[id]
		l.Kind = csKinds[id]
		h.Landed[id] = l
	}
	retired := map[string]bool{} // retired link IDs
	for _, r := range retires {
		if r.id == "" || !dropped[r.id] {
			retired[r.link] = true
		}
	}
	served := map[string]map[string]bool{} // milestone handle -> change sets
	for _, s := range serves {
		if !sets[s.from] || ended[s.from] || retired[s.link] || s.id != "" && dropped[s.id] {
			continue
		}
		m := milestoneOf(s.to, parents, isMilestone)
		if m == "" {
			continue
		}
		hd, ok := intents[m]
		if !ok {
			continue
		}
		if served[hd] == nil {
			served[hd] = map[string]bool{}
		}
		if !served[hd][s.from] {
			served[hd][s.from] = true
			h.Milestones[hd] = append(h.Milestones[hd], s.from)
		}
	}
	return h, nil
}

// milestoneOf returns the milestone node at or above node, or "" when no
// ancestor is a milestone.
func milestoneOf(node string, parents map[string]string, isMilestone map[string]bool) string {
	seen := map[string]bool{}
	for n := node; n != "" && !seen[n]; n = parents[n] {
		if isMilestone[n] {
			return n
		}
		seen[n] = true
	}
	return ""
}

// fieldSet reports whether a field's JSON value is present and neither ""
// nor null.
func fieldSet(v json.RawMessage) bool {
	v = bytes.TrimSpace(v)
	return len(v) > 0 && string(v) != `""` && string(v) != "null"
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

// Report is the evaluated next release: its version and kind, the gates for
// the candidate commit and the drafted notes.
type Report struct {
	Tag    string // last v* tag; "" for none
	Next   Next
	Commit string
	Gates  []Gate
	Notes  string
}

// evaluate computes the next version and its kind since the last v* tag in
// the repo at dir, each release gate for the candidate commit and the notes.
// Only an annotated tag is a release: the lightweight tag that publishing a
// major's draft GitHub release makes (its sign-off) is not, so describe
// runs without --tags. It reads only: nothing is tagged or written.
func evaluate(dir string, env CheckEnv) (Report, error) {
	tag, err := gitOut(dir, "describe", "--match", "v[0-9]*.[0-9]*.[0-9]*", "--abbrev=0")
	if err != nil {
		tag = ""
	}
	since := time.Time{}
	if tag != "" {
		d, err := gitOut(dir, "log", "-1", "--format=%cI", tag)
		if err != nil {
			return Report{}, err
		}
		if since, err = time.Parse(time.RFC3339, d); err != nil {
			return Report{}, err
		}
	}
	logDir := filepath.Join(dir, ".rhiza", "log")
	h, err := ReadHistory(logDir)
	if err != nil {
		return Report{}, err
	}
	marks, err := ReadMarks(logDir)
	if err != nil {
		return Report{}, err
	}
	in := h.Inputs(since, marks)
	var bs []Break
	if tag != "" {
		detect := env.Breaking
		if detect == nil {
			detect = breakingAgainstLastTag
		}
		if _, bs, err = detect(dir); err != nil {
			return Report{}, err
		}
	}
	for _, b := range bs {
		in.Breaking = append(in.Breaking, b.String())
	}
	n, err := NextVersion(tag, in)
	if err != nil {
		return Report{}, err
	}
	detected := make([]string, 0, len(bs))
	for _, b := range bs {
		detected = append(detected, b.String())
	}
	notes := Notes(h, since, marks, detected)
	gin, err := gateInputs(dir, logDir, env, h, in, n.Kind, marks)
	if err != nil {
		return Report{}, err
	}
	return Report{Tag: tag, Next: n, Commit: gin.Commit, Gates: Gates(gin), Notes: notes}, nil
}

// runCheck prints the next version and its kind since the last v* tag in the
// repo at dir, then each release gate for the candidate commit with why it
// passes or fails. It reads only: nothing is tagged or written.
func runCheck(dir string, w io.Writer, env CheckEnv) error {
	r, err := evaluate(dir, env)
	if err != nil {
		return err
	}
	if r.Next.Kind == KindNone {
		fmt.Fprintf(w, "none: no release since %s\n", orNone(r.Tag))
		printGates(w, r.Commit, r.Gates)
		return nil
	}
	fmt.Fprintf(w, "%s %s (since %s)\n", r.Next.Version, r.Next.Kind, orNone(r.Tag))
	for _, reason := range r.Next.Reasons {
		fmt.Fprintln(w, "  "+reason)
	}
	printGates(w, r.Commit, r.Gates)
	fmt.Fprintf(w, "\n# %s\n\n%s", r.Next.Version, r.Notes)
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
