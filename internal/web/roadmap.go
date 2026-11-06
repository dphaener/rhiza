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
		return fmt.Errorf("GET %s: %d %s", path, rec.Code, strings.TrimSpace(rec.Body.String()))
	}
	return json.Unmarshal(rec.Body.Bytes(), v)
}

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
			Intent    string   `json:"intent"`
			Ready     bool     `json:"ready"`
			DependsOn []string `json:"depends_on"`
		} `json:"change_sets"`
	} `json:"groups"`
}

// Seg is one state's share of a group's bar and legend.
type Seg struct {
	State string
	N     int
}

// RoadmapGroup is a milestone, or the backlog: its title, a bar of its change
// sets by state and a row for each.
type RoadmapGroup struct {
	Kind, Title string
	Segs        []Seg
	Rows        []Row
}

// RoadmapView is what the roadmap page and its live section render.
type RoadmapView struct {
	Milestones []RoadmapGroup
	Backlog    *RoadmapGroup
	Total      int
}

func (v *RoadmapView) Kicker() string {
	b := 0
	if v.Backlog != nil {
		b = len(v.Backlog.Rows)
	}
	return fmt.Sprintf("%s · %s · %d in backlog", plural(len(v.Milestones), "milestone"), plural(v.Total, "change set"), b)
}

func plural(n int, w string) string {
	if n == 1 {
		return "1 " + w
	}
	return fmt.Sprintf("%d %ss", n, w)
}

// loadRoadmap reads the roadmap operation, whole, and shapes it into groups.
// Milestones are named by their intent's statement, never by handle or ID.
func loadRoadmap(r *http.Request) (*RoadmapView, error) {
	var rm roadmapJSON
	if err := apiGet(r, "/api/roadmap?all=true", &rm); err != nil {
		return nil, err
	}
	v := &RoadmapView{}
	states := []string{"accepted", "proposed"}
	for _, g := range rm.Groups {
		out := RoadmapGroup{Kind: g.Kind, Title: "Not in a milestone"}
		if g.Kind != "backlog" {
			out.Title = "Milestone"
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
			sub, step := "", "—"
			if c.Ready {
				step = "ready"
			}
			if n := len(c.DependsOn); n > 0 {
				sub = "depends on " + plural(n, "change set")
			}
			out.Rows = append(out.Rows, Row{Href: "/cs/" + c.ID, Title: c.Intent, Sub: sub,
				Pill: Pill{State: c.State, Text: c.State}, Step: step})
			count[c.State]++
			v.Total++
		}
		for _, s := range states {
			if count[s] > 0 {
				out.Segs = append(out.Segs, Seg{s, count[s]})
			}
		}
		if g.Kind == "backlog" {
			v.Backlog = &out
		} else {
			v.Milestones = append(v.Milestones, out)
		}
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
	render(w, "roadmap", page{Title: "Roadmap", Nav: "roadmap", Data: v})
}
