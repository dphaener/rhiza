package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// The operations' responses, as far as a change set's page reads them.
type changeJSON struct {
	ID     string   `json:"id"`
	Kind   string   `json:"kind"`
	Intent string   `json:"intent"`
	Title  string   `json:"title"`
	State  string   `json:"state"`
	Serves []string `json:"serves"`
	Leaves []struct {
		ID string `json:"id"`
	} `json:"leaves"`
	History []struct {
		Type string `json:"type"`
		At   string `json:"at"`
	} `json:"history"`
	Changes struct {
		Nodes []struct {
			Node   string `json:"node"`
			Kind   string `json:"kind"`
			Action string `json:"action"`
		} `json:"nodes"`
	} `json:"changes"`
}

type showLinkJSON struct {
	Type    string `json:"type"`
	Retired bool   `json:"retired"`
	Other   struct {
		ID string `json:"id"`
	} `json:"other"`
}

type showJSON struct {
	Node struct {
		Kind     string                     `json:"kind"`
		Retired  bool                       `json:"retired"`
		Fields   map[string]json.RawMessage `json:"fields"`
		LinksIn  []showLinkJSON             `json:"links_in"`
		LinksOut []showLinkJSON             `json:"links_out"`
	} `json:"node"`
}

func (s *showJSON) field(name string) string {
	var v string
	json.Unmarshal(s.Node.Fields[name], &v)
	return v
}

func (s *showJSON) liveIn(typ string) []string {
	var out []string
	for _, l := range s.Node.LinksIn {
		if l.Type == typ && !l.Retired {
			out = append(out, l.Other.ID)
		}
	}
	return out
}

// Rel is a related change set or intent: its words, its state and where it
// links to.
type Rel struct {
	Href, Text string
	Pill       Pill
}

// LeafView is a leaf with its acceptance and whether it has Evidence.
type LeafView struct {
	Statement, Acceptance string
	Evidence              bool
}

// HistoryItem is one lifecycle event.
type HistoryItem struct{ Time, Datetime, Text string }

// ChangeView is what a change set's page renders.
type ChangeView struct {
	Title, Intent, Kind string
	State               Pill
	Serves, DependsOn   []Rel
	Leaves              []LeafView
	WithEvidence        int
	Questions           []string
	Runs                []Row
	History             []HistoryItem
}

var lifecycleWords = map[string]string{
	"propose": "Proposed", "accept": "Accepted", "reject": "Rejected", "withdraw": "Withdrawn",
	"supersede": "Superseded", "land": "Landed", "reopen": "Reopened",
}

var runStatePill = map[string]Pill{
	"running":   {"running", "running"},
	"blocked":   {"waiting", "waiting on you"},
	"stopped":   {"interrupted", "stopped"},
	"succeeded": {"done", "done"},
	"failed":    {"failed", "failed"},
}

// inFlight is the change set states that have a what-if view.
var inFlight = map[string]bool{"draft": true, "proposed": true, "accepted": true}

// showRef reads a node as the change set sees it, whatever change set the
// daemon's repo has current: its what-if view while it is in flight, and
// the as-is model once it has landed (or will never land).
func showRef(r *http.Request, cs *changeJSON, ref string) (*showJSON, error) {
	view := "as-is"
	if inFlight[cs.State] {
		view = "what-if:" + cs.ID
	}
	var s showJSON
	return &s, apiGet(r, "/api/show?ref="+url.QueryEscape(ref)+"&view="+url.QueryEscape(view), &s)
}

// dependsOn is the change sets the model says cs depends on, whichever
// change set added the links: the roadmap operation's depends_on.
func dependsOn(r *http.Request, cs string) []string {
	var rm roadmapJSON
	if err := apiGet(r, "/api/roadmap?all=true", &rm); err != nil {
		return nil
	}
	for _, g := range rm.Groups {
		for _, c := range g.ChangeSets {
			if c.ID == cs {
				return c.DependsOn
			}
		}
	}
	return nil
}

// loadChange reads a change set through the operations and shapes it. No IDs
// or handles reach the view: nodes appear by their statements.
func loadChange(r *http.Request, ref string) (*ChangeView, error) {
	var cs changeJSON
	if err := apiGet(r, "/api/change.show?ref="+url.QueryEscape(ref), &cs); err != nil {
		return nil, err
	}
	v := &ChangeView{Title: cs.Title, Intent: cs.Intent, Kind: cs.Kind, State: Pill{cs.State, cs.State}}
	if v.Title == "" {
		v.Title = cs.Intent
	}

	for _, id := range cs.Serves {
		if s, err := showRef(r, &cs, id); err == nil {
			if t := s.field("statement"); t != "" {
				v.Serves = append(v.Serves, Rel{Href: "/roadmap", Text: t})
			}
		}
	}

	seen := map[string]bool{}
	for _, to := range dependsOn(r, cs.ID) {
		if seen[to] {
			continue
		}
		seen[to] = true
		var d changeJSON
		if err := apiGet(r, "/api/change.show?ref="+url.QueryEscape(to), &d); err != nil {
			continue
		}
		t := d.Title
		if t == "" {
			t = d.Intent
		}
		v.DependsOn = append(v.DependsOn, Rel{Href: "/changes/" + d.ID, Text: t, Pill: Pill{d.State, d.State}})
	}

	qs := map[string]bool{}
	var qorder []string
	addQ := func(id string) {
		if !qs[id] {
			qs[id] = true
			qorder = append(qorder, id)
		}
	}
	for _, lf := range cs.Leaves {
		s, err := showRef(r, &cs, lf.ID)
		if err != nil {
			continue
		}
		lv := LeafView{Statement: s.field("statement"), Acceptance: s.field("acceptance")}
		for _, e := range s.liveIn("verifies") {
			if en, err := showRef(r, &cs, e); err == nil && !en.Node.Retired {
				lv.Evidence = true
			}
		}
		if lv.Evidence {
			v.WithEvidence++
		}
		for _, q := range s.liveIn("blocks") {
			addQ(q)
		}
		v.Leaves = append(v.Leaves, lv)
	}
	for _, n := range cs.Changes.Nodes {
		if n.Kind == "question" && n.Action == "created" {
			addQ(n.Node)
		}
	}
	for _, id := range qorder {
		q, err := showRef(r, &cs, id)
		if err != nil || q.Node.Retired || len(q.liveIn("answers")) > 0 {
			continue
		}
		if t := q.field("statement"); t != "" {
			v.Questions = append(v.Questions, t)
		}
	}

	var runs struct {
		Runs []struct {
			ID        string `json:"id"`
			Change    string `json:"change"`
			State     string `json:"state"`
			Current   string `json:"current"`
			Workflow  string `json:"workflow"`
			StartedAt string `json:"started_at"`
			Live      bool   `json:"live"`
		} `json:"runs"`
	}
	if err := apiGet(r, "/api/run.ls?limit=1000", &runs); err == nil {
		for _, ru := range runs.Runs {
			if ru.Change != cs.ID {
				continue
			}
			p, ok := runStatePill[ru.State]
			if !ok {
				p = Pill{"interrupted", ru.State}
			}
			step := ru.Current
			if step == "" {
				step = "—"
			}
			row := Row{Href: "/runs/" + ru.ID, Title: ru.Workflow + " run", Pill: p, Step: step, Live: ru.Live}
			if t, err := time.Parse(time.RFC3339, ru.StartedAt); err == nil {
				row.Datetime, row.Time = ru.StartedAt, t.UTC().Format("2006-01-02 15:04")
			}
			v.Runs = append(v.Runs, row)
		}
	}

	for _, h := range cs.History {
		item := HistoryItem{Text: lifecycleWords[h.Type]}
		if item.Text == "" {
			item.Text = strings.ToUpper(h.Type[:1]) + h.Type[1:]
		}
		if t, err := time.Parse(time.RFC3339, h.At); err == nil {
			item.Datetime, item.Time = h.At, t.UTC().Format("2006-01-02 15:04")
		}
		v.History = append(v.History, item)
	}
	return v, nil
}

func (v *ChangeView) Facts() string {
	n := "1 leaf"
	if len(v.Leaves) != 1 {
		n = strconv.Itoa(len(v.Leaves)) + " leaves"
	}
	return n + ", " + strconv.Itoa(v.WithEvidence) + " with evidence"
}

func serveChange(w http.ResponseWriter, r *http.Request) {
	v, err := loadChange(r, r.PathValue("id"))
	var nf *notFoundError
	if errors.As(err, &nf) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "change set: "+err.Error(), http.StatusInternalServerError)
		return
	}
	render(w, "change", page{Title: v.Title, Nav: "roadmap", Data: v})
}
