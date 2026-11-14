package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// The log operation's response, as far as a node's page reads it.
type nodeLogJSON struct {
	ChangeSets []struct {
		ID     string `json:"id"`
		State  string `json:"state"`
		Intent string `json:"intent"`
		At     string `json:"at"`
		Ops    []struct {
			Type string `json:"type"`
		} `json:"ops"`
	} `json:"change_sets"`
}

// NodeLink is a node, source or change set another node links to or from.
type NodeLink struct {
	Href, Text, Sub string
	Glyph           string
}

// NodeGroup is the links of one kind of relation, under a heading.
type NodeGroup struct {
	Title string
	Links []NodeLink
}

// NodeEvent is one change set of a node's history.
type NodeEvent struct {
	Href, Text, Time, Datetime string
}

// NodeView is what a model node's page renders. It carries no handle or ID:
// nodes appear by their statements and change sets by their titles.
type NodeView struct {
	Kind, KindLabel, Statement string
	Retired                    bool
	// Quote is the owner's own words behind a decision, with where they were
	// said; Wording is the agent's wording the owner picked, when there is no
	// quote. Question is the question the decision answered.
	Quote, QuoteWhere, Wording, Question string
	Summary                              string
	Groups                               []NodeGroup
	Origins                              []NodeLink
	History                              []NodeEvent
}

// groupTitles word a link as read from the page's node, by the link's name.
var groupTitles = map[string]string{
	"serves": "Serves", "served-by": "Served by",
	"constrains": "Constrains", "constrained-by": "Constrained by",
	"decides": "What it decided", "decided-by": "Decided by",
	"answers": "Answers", "answered-by": "Answered by",
	"blocks": "Blocks", "blocked-by": "Blocked by",
	"depends-on": "Depends on", "depended-on-by": "Needed by",
	"supersedes": "Replaces", "superseded-by": "Replaced by",
	"verifies": "Verifies", "verified-by": "Evidence",
	"implements": "Implements", "implemented-by": "Implemented by",
	"relates-to": "Related", "related-to": "Related",
	"observes": "Observes", "observed-by": "Observed by",
	"discusses": "Discusses", "discussed-by": "Discussed by",
	"exempts": "Exempts", "exempted-by": "Exempted by",
	"discovered-from": "Found while working on", "discovered": "Found along the way",
	"derived-from": "Sources", "source-of": "Drawn on by",
}

var kindLabels = map[string]string{
	"capability": "Capability", "intent": "Intent", "constraint": "Constraint",
	"decision": "Decision", "question": "Question", "source": "Source",
	"evidence": "Evidence", "note": "Note",
}

// groupOrder keeps the groups in a steady order: what the node decided or
// links to first, what points at it after.
var groupOrder = []string{"decides", "serves", "constrains", "answers", "blocks", "depends-on", "verified-by", "implements", "supersedes", "superseded-by", "decided-by", "relates-to"}

// words is a node's own words: its statement, else its title or summary.
func (s *showJSON) words() string {
	for _, f := range []string{"statement", "title", "summary"} {
		if t := s.field(f); t != "" {
			return t
		}
	}
	return ""
}

func nodeGlyph(kind string) string {
	switch kind {
	case "intent", "capability", "constraint", "decision", "question":
		return kind
	}
	return "other"
}

// loadNode reads a node through the operations and shapes it.
func loadNode(r *http.Request, ref string) (*NodeView, error) {
	var s struct {
		Node struct {
			Kind     string                     `json:"kind"`
			Retired  bool                       `json:"retired"`
			Fields   map[string]json.RawMessage `json:"fields"`
			LinksIn  []nodeLinkJSON             `json:"links_in"`
			LinksOut []nodeLinkJSON             `json:"links_out"`
		} `json:"node"`
	}
	q := url.QueryEscape(ref)
	if err := apiGet(r, "/api/show?ref="+q+"&view=to-be", &s); err != nil {
		return nil, err
	}
	self := showJSON{}
	self.Node.Fields = s.Node.Fields
	v := &NodeView{Kind: s.Node.Kind, Retired: s.Node.Retired, Statement: self.field("statement")}
	v.KindLabel = kindLabels[v.Kind]
	if v.KindLabel == "" && v.Kind != "" {
		v.KindLabel = strings.ToUpper(v.Kind[:1]) + v.Kind[1:]
	}
	if v.Statement == "" {
		v.Statement = self.words()
	}
	if v.Kind == "source" {
		v.Summary = self.field("summary")
		if self.field("title") != "" && v.Summary == v.Statement {
			v.Summary = ""
		}
	}
	v.Quote, v.QuoteWhere, v.Wording = self.field("quote"), self.field("quote_where"), self.field("wording")
	v.Question = self.field("question")
	if v.Quote != "" {
		// The decision quotes the owner: the wording is only the assistant's
		// summary of it, not shown as a second claim.
		v.Wording = ""
	}

	// The other ends, read once each.
	type end struct {
		words, kind, href string
		ok                bool
	}
	ends := map[string]end{}
	resolve := func(l nodeLinkJSON) end {
		if e, ok := ends[l.Other.ID]; ok {
			return e
		}
		var e end
		if l.Other.ChangeSet {
			var cs changeJSON
			if apiGet(r, "/api/change.show?ref="+url.QueryEscape(l.Other.ID), &cs) == nil {
				e = end{words: firstOf(cs.Title, cs.Intent), kind: "change set", href: "/changes/" + cs.ID, ok: true}
			}
		} else {
			var o showJSON
			if apiGet(r, "/api/show?ref="+url.QueryEscape(l.Other.ID)+"&view=to-be", &o) == nil && !o.Node.Retired {
				e = end{words: o.words(), kind: o.Node.Kind, href: "/nodes/" + l.Other.ID, ok: o.words() != ""}
			}
		}
		ends[l.Other.ID] = e
		return e
	}

	byName := map[string][]NodeLink{}
	add := func(links []nodeLinkJSON) {
		for _, l := range links {
			if l.Retired {
				continue
			}
			e := resolve(l)
			if !e.ok {
				continue
			}
			nl := NodeLink{Href: e.href, Text: e.words, Sub: kindLabels[e.kind], Glyph: nodeGlyph(e.kind)}
			if e.kind == "change set" {
				nl.Sub, nl.Glyph = "Change set", "intent"
			}
			switch {
			case l.Name == "derived-from" && e.kind == "source":
				v.Origins = append(v.Origins, nl)
			case l.Name == "answers" && v.Question == "" && v.Kind == "decision":
				v.Question = e.words
				byName[l.Name] = append(byName[l.Name], nl)
			default:
				byName[l.Name] = append(byName[l.Name], nl)
			}
		}
	}
	add(s.Node.LinksOut)
	add(s.Node.LinksIn)

	seen := map[string]bool{}
	for _, n := range groupOrder {
		if ls := byName[n]; len(ls) > 0 {
			v.Groups = append(v.Groups, NodeGroup{Title: groupTitle(n), Links: ls})
			seen[n] = true
		}
	}
	var rest []string
	for n := range byName {
		if !seen[n] {
			rest = append(rest, n)
		}
	}
	sort.Strings(rest)
	for _, n := range rest {
		v.Groups = append(v.Groups, NodeGroup{Title: groupTitle(n), Links: byName[n]})
	}

	// Its history: the change sets that touched it, in order.
	var lg nodeLogJSON
	if apiGet(r, "/api/log?ref="+q+"&view=to-be", &lg) == nil {
		for i, c := range lg.ChangeSets {
			var cs changeJSON
			title := c.Intent
			if apiGet(r, "/api/change.show?ref="+url.QueryEscape(c.ID), &cs) == nil {
				title = firstOf(cs.Title, cs.Intent)
			}
			verb := "Changed in"
			if i == 0 {
				verb = "Recorded in"
			}
			ev := NodeEvent{Href: "/changes/" + c.ID, Text: verb + " " + title}
			if t, err := time.Parse(time.RFC3339, c.At); err == nil {
				ev.Datetime, ev.Time = c.At, t.UTC().Format("2006-01-02 15:04")
			} else {
				ev.Time = c.State
			}
			v.History = append(v.History, ev)
			v.Origins = append(v.Origins, NodeLink{Href: ev.Href, Text: title, Sub: "Change set · " + c.State, Glyph: "intent"})
		}
	}
	return v, nil
}

type nodeLinkJSON struct {
	Name    string `json:"name"`
	Retired bool   `json:"retired"`
	Other   struct {
		ID        string `json:"id"`
		ChangeSet bool   `json:"change_set"`
	} `json:"other"`
}

func firstOf(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func groupTitle(name string) string {
	if t, ok := groupTitles[name]; ok {
		return t
	}
	w := strings.ReplaceAll(name, "-", " ")
	return strings.ToUpper(w[:1]) + w[1:]
}

// NotFound is what the 404 page names: the ref asked for.
type NodeMissing struct{ Ref string }

func serveNode(w http.ResponseWriter, r *http.Request) {
	ref := r.PathValue("ref")
	v, err := loadNode(r, ref)
	var nf *notFoundError
	if errors.As(err, &nf) {
		w.WriteHeader(http.StatusNotFound)
		render(w, r, "nodemissing", page{Title: "Not found", Nav: "roadmap", Data: NodeMissing{Ref: ref}})
		return
	}
	if err != nil {
		http.Error(w, "node: "+err.Error(), http.StatusInternalServerError)
		return
	}
	render(w, r, "node", page{Title: v.Statement, Nav: "roadmap", Data: v})
}
