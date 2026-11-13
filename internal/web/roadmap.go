package web

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

// apiKey carries the read-only API handler to the page and fragment code, so
// the pages read the same operations the CLI does without importing api.
type apiKey struct{}

func withAPI(r *http.Request, api http.Handler) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), apiKey{}, api))
}

// apiGet runs GET path on the API handler of r and decodes its JSON into v.
func apiGet(r *http.Request, path string, v any) error {
	api, _ := r.Context().Value(apiKey{}).(http.Handler)
	if api == nil {
		return fmt.Errorf("no API to read %s", path)
	}
	rec := httptest.NewRecorder()
	api.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	if rec.Code != http.StatusOK {
		if rec.Code == http.StatusNotFound || rec.Code == http.StatusBadRequest {
			return &notFoundError{path}
		}
		return fmt.Errorf("GET %s: %d %s", path, rec.Code, strings.TrimSpace(rec.Body.String()))
	}
	return json.Unmarshal(rec.Body.Bytes(), v)
}

// notFoundError is an API read that found nothing, or took no such ref.
type notFoundError struct{ path string }

func (e *notFoundError) Error() string { return "not found: " + e.path }

// The roadmap operation's response, as far as the page reads it.
type roadmapJSON struct {
	Groups []struct {
		Kind   string `json:"kind"`
		Intent struct {
			ID string `json:"id"`
		} `json:"intent"`
		ChangeSets []struct {
			ID        string   `json:"id"`
			State     string   `json:"state"`
			Title     string   `json:"title"`
			Intent    string   `json:"intent"`
			Ready     bool     `json:"ready"`
			DependsOn []string `json:"depends_on"`
		} `json:"change_sets"`
	} `json:"groups"`
}

// The milestones operation's response, as far as the page reads it.
type milestonesJSON struct {
	Milestones []struct {
		Statement  string `json:"statement"`
		State      string `json:"state"`
		DonePoints []struct {
			Statement string `json:"statement"`
			Met       bool   `json:"met"`
		} `json:"done_points"`
		ChangeSets []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			State string `json:"state"`
		} `json:"change_sets"`
	} `json:"milestones"`
}

// Seg is one state's share of a group's bar and legend.
type Seg struct {
	State string
	N     int
}

// shownRows is how many of a group's change sets show before its fold.
const shownRows = 5

// RowGroup is a labelled run of change set rows: the first few shown and
// the rest in a fold that opens in place.
type RowGroup struct {
	Label string
	Shown []Row
	More  []Row
	// MoreText is the fold's summary, e.g. "7 more landed, 1 withdrawn".
	MoreText string
}

// DonePointView is a milestone's done point and whether it is met.
type DonePointView struct {
	Text string
	Met  bool
}

// RoadmapGroup is a milestone, the backlog or an umbrella: its title, a bar
// of its change sets by state and its rows. A milestone has its state word,
// its done points and its rows in RowGroups; the backlog and an umbrella
// have Rows.
type RoadmapGroup struct {
	Kind, Title string
	State       string
	Segs        []Seg
	DonePoints  []DonePointView
	Met         int
	RowGroups   []RowGroup
	Rows        []Row
}

// FinishedMilestone is a finished milestone's name and done points.
type FinishedMilestone struct {
	Name       string
	DonePoints []DonePointView
}

// RoadmapView is what the roadmap page and its live section render: the
// open milestones, then the backlog and the umbrellas, and the finished
// milestones.
type RoadmapView struct {
	Groups   []RoadmapGroup
	Finished []FinishedMilestone
	Total    int
}

func (v *RoadmapView) Kicker() string {
	m, b := 0, 0
	for _, g := range v.Groups {
		switch g.Kind {
		case "milestone":
			m++
		case "backlog":
			b += len(g.Rows)
		}
	}
	return fmt.Sprintf("%s · %s · %d in backlog", plural(m, "milestone"), plural(v.Total, "change set"), b)
}

func plural(n int, w string) string {
	if n == 1 {
		return "1 " + w
	}
	return fmt.Sprintf("%d %ss", n, w)
}

// changeIntent reads a change set's intent, for the sub-line of its row.
func changeIntent(r *http.Request, id string) string {
	var cs changeJSON
	if err := apiGet(r, "/api/change.show?ref="+url.QueryEscape(id), &cs); err != nil {
		return ""
	}
	return cs.Intent
}

// activeRuns is the change sets with a run not finished: going, waiting on
// the owner or stopped.
func activeRuns(r *http.Request) map[string]string {
	var runs struct {
		Runs []struct {
			Change string `json:"change"`
			State  string `json:"state"`
		} `json:"runs"`
	}
	out := map[string]string{}
	if err := apiGet(r, "/api/run.ls?limit=1000", &runs); err != nil {
		return out
	}
	for _, ru := range runs.Runs {
		if ru.Change == "" || ru.State == "succeeded" || ru.State == "failed" {
			continue
		}
		if cur, ok := out[ru.Change]; !ok || cur != "running" {
			out[ru.Change] = ru.State
		}
	}
	return out
}

// fold splits rows into the first shownRows and the rest.
func fold(label string, rows []Row, more string) RowGroup {
	g := RowGroup{Label: label, Shown: rows}
	if len(rows) > shownRows {
		g.Shown, g.More = rows[:shownRows], rows[shownRows:]
		g.MoreText = more
	}
	return g
}

// milestoneGroup shapes an open milestone: its bar, its done points and its
// change sets grouped by what is happening to them.
func milestoneGroup(r *http.Request, m milestonesJSON, i int, runs map[string]string) RoadmapGroup {
	ms := m.Milestones[i]
	g := RoadmapGroup{Kind: "milestone", Title: ms.Statement, State: ms.State}
	for _, d := range ms.DonePoints {
		g.DonePoints = append(g.DonePoints, DonePointView{d.Statement, d.Met})
		if d.Met {
			g.Met++
		}
	}
	row := func(c struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		State string `json:"state"`
	}, p Pill) Row {
		return Row{Href: "/changes/" + c.ID, Title: c.Title, Sub: changeIntent(r, c.ID), Pill: p, Step: "—"}
	}
	var building, accepted, proposed, landed, withdrawn, drafts []Row
	count := map[string]int{}
	for _, c := range ms.ChangeSets {
		switch c.State {
		case "accepted":
			count["accepted"]++
			if st, ok := runs[c.ID]; ok {
				p := Pill{"running", "building"}
				if st != "running" {
					p = Pill{"waiting", "stuck"}
				}
				rw := row(c, p)
				rw.Live = st == "running"
				building = append(building, rw)
			} else {
				accepted = append(accepted, row(c, Pill{"accepted", "accepted"}))
			}
		case "proposed":
			count["proposed"]++
			proposed = append(proposed, row(c, Pill{"proposed", "proposed"}))
		case "landed":
			count["landed"]++
			landed = append(landed, row(c, Pill{"landed", "landed"}))
		case "withdrawn":
			count["withdrawn"]++
			withdrawn = append(withdrawn, row(c, Pill{"withdrawn", "withdrawn"}))
		case "draft":
			drafts = append(drafts, row(c, Pill{"draft", "draft"}))
		}
	}
	for _, s := range []string{"landed", "accepted", "proposed", "withdrawn"} {
		if count[s] > 0 {
			g.Segs = append(g.Segs, Seg{s, count[s]})
		}
	}
	moreOf := func(rows []Row, n int, word string) string {
		return fmt.Sprintf("%d more %s", len(rows)-n, word)
	}
	add := func(label string, rows []Row, word string) {
		if len(rows) > 0 {
			g.RowGroups = append(g.RowGroups, fold(label, rows, moreOf(rows, shownRows, word)))
		}
	}
	add("Building now", building, "building")
	add("Accepted, waiting", accepted, "accepted")
	add("Proposed", proposed, "proposed")
	landed = append(landed, withdrawn...)
	if len(landed) > 0 {
		lg := fold("Landed", landed, "")
		if len(lg.More) > 0 {
			nl, nw := 0, 0
			for _, rw := range lg.More {
				if rw.Pill.State == "withdrawn" {
					nw++
				} else {
					nl++
				}
			}
			var parts []string
			if nl > 0 {
				parts = append(parts, fmt.Sprintf("%d more landed", nl))
			}
			if nw > 0 {
				parts = append(parts, fmt.Sprintf("%d withdrawn", nw))
			}
			lg.MoreText = strings.Join(parts, ", ")
		}
		g.RowGroups = append(g.RowGroups, lg)
	}
	add("Drafts", drafts, "drafts")
	return g
}

// loadRoadmap reads the milestones operation for the open milestones and the
// roadmap operation for the backlog and the umbrellas, and shapes them into
// groups. Milestones and umbrellas are named by their intent's statement,
// never by handle or ID; each row shows the change set's title, its intent
// beneath and its state.
func loadRoadmap(r *http.Request) (*RoadmapView, error) {
	var rm roadmapJSON
	if err := apiGet(r, "/api/roadmap?all=true", &rm); err != nil {
		return nil, err
	}
	var ms milestonesJSON
	if err := apiGet(r, "/api/milestones", &ms); err != nil {
		return nil, err
	}
	v := &RoadmapView{}
	runs := activeRuns(r)
	for i, m := range ms.Milestones {
		if m.State == "finished" {
			f := FinishedMilestone{Name: m.Statement}
			for _, d := range m.DonePoints {
				f.DonePoints = append(f.DonePoints, DonePointView{d.Statement, d.Met})
			}
			v.Finished = append(v.Finished, f)
			continue
		}
		g := milestoneGroup(r, ms, i, runs)
		for _, c := range m.ChangeSets {
			if c.State == "draft" || c.State == "proposed" || c.State == "accepted" || c.State == "landed" || c.State == "withdrawn" {
				v.Total++
			}
		}
		v.Groups = append(v.Groups, g)
	}
	for _, g := range rm.Groups {
		if g.Kind == "milestone" {
			continue
		}
		out := RoadmapGroup{Kind: g.Kind, Title: "Not in a milestone"}
		if g.Kind == "umbrella" {
			out.Title = "Umbrella"
			var s struct {
				Node struct {
					Fields map[string]string `json:"fields"`
				} `json:"node"`
			}
			if err := apiGet(r, "/api/show?ref="+url.QueryEscape(g.Intent.ID), &s); err == nil {
				if st := s.Node.Fields["statement"]; st != "" {
					out.Title = st
				}
			}
		}
		count := map[string]int{}
		for _, c := range g.ChangeSets {
			note, step := "", "—"
			if c.Ready {
				step = "ready"
			}
			if n := len(c.DependsOn); n > 0 {
				note = "depends on " + plural(n, "change set")
			}
			out.Rows = append(out.Rows, Row{Href: "/changes/" + c.ID, Title: c.Title, Sub: c.Intent, Note: note,
				Pill: Pill{State: c.State, Text: c.State}, Step: step})
			count[c.State]++
			v.Total++
		}
		for _, s := range []string{"accepted", "proposed"} {
			if count[s] > 0 {
				out.Segs = append(out.Segs, Seg{s, count[s]})
			}
		}
		v.Groups = append(v.Groups, out)
	}
	return v, nil
}

func init() {
	RegisterFragment("roadmap", &Fragment{
		Events: []string{EventChange},
		Template: template.Must(template.New("roadmap").Funcs(template.FuncMap{"section": section}).
			ParseFS(templateFS, "templates/partials.html", "templates/roadmap.html")).Lookup("roadmap-body"),
		Data: func(r *http.Request) any {
			v, err := loadRoadmap(r)
			if err != nil {
				return &RoadmapView{}
			}
			return v
		},
	})
}

func serveRoadmap(w http.ResponseWriter, r *http.Request) {
	v, err := loadRoadmap(r)
	if err != nil {
		http.Error(w, "roadmap: "+err.Error(), http.StatusInternalServerError)
		return
	}
	render(w, r, "roadmap", page{Title: "Roadmap", Nav: "roadmap", Data: v})
}
