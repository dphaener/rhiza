package web

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"
)

// The Runs page, /runs: the live runs first, then the recent ones, newest
// first. It reads run.ls through the API handler the page is served with,
// and its list is a live section that re-renders when a run's journal moves.

// The most runs read from run.ls, and the most recent (not live) ones shown.
const (
	runsRead   = 200
	runsRecent = 50
)

// runItem is the part of a run.ls item the page uses.
type runItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	State     string `json:"state"`
	Current   string `json:"current"`
	StartedAt string `json:"started_at"`
	Live      bool   `json:"live"`
	Workflow  string `json:"workflow"`
}

// runsView is what the runs list section is rendered with.
type runsView struct {
	Live, Recent []Row
	Empty        *Empty
}

type apiKey struct{}

// withAPI makes the API handler reachable from a request's context, for the
// pages that read it.
func withAPI(api http.Handler, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), apiKey{}, api)))
	})
}

// apiGet reads one API operation's JSON into v.
func apiGet(r *http.Request, target string, v any) error {
	api, _ := r.Context().Value(apiKey{}).(http.Handler)
	if api == nil {
		return errNoAPI
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, target, nil).WithContext(r.Context())
	api.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		return &apiError{target, rec.Code}
	}
	return json.Unmarshal(rec.Body.Bytes(), v)
}

type apiError struct {
	target string
	code   int
}

func (e *apiError) Error() string { return "GET " + e.target + ": " + http.StatusText(e.code) }

var errNoAPI = &apiError{"/api/", http.StatusNotFound}

// runPill maps a run's state to its pill.
func runPill(state string) Pill {
	switch state {
	case "running":
		return Pill{"running", "running"}
	case "blocked":
		return Pill{"waiting", "waiting on you"}
	case "stopped":
		return Pill{"interrupted", "interrupted"}
	case "succeeded":
		return Pill{"done", "done"}
	case "failed":
		return Pill{"failed", "failed"}
	}
	return Pill{"interrupted", state}
}

// runRow is a run's dense row. It links to the run's page and shows no ID.
func runRow(it runItem) Row {
	row := Row{
		Href: "/runs/" + it.ID, Title: it.Title, Sub: it.Workflow,
		Pill: runPill(it.State), Step: it.Current, Live: it.State == "running",
	}
	if row.Title == "" {
		row.Title = "Untitled change set"
	}
	if t, err := time.Parse(time.RFC3339, it.StartedAt); err == nil {
		row.Datetime = it.StartedAt
		row.Time = t.UTC().Format("Jan 2 15:04")
	}
	return row
}

// buildRuns splits runs (newest first) into the live ones and the recent
// ones, keeping their order. A run is live while a process holds it or it
// waits on the owner.
func buildRuns(items []runItem) runsView {
	var v runsView
	for _, it := range items {
		if it.Live || it.State == "blocked" {
			v.Live = append(v.Live, runRow(it))
		} else if len(v.Recent) < runsRecent {
			v.Recent = append(v.Recent, runRow(it))
		}
	}
	return v
}

func runsData(r *http.Request) any {
	var resp struct {
		Runs []runItem `json:"runs"`
	}
	if err := apiGet(r, "/api/run.ls?limit="+strconv.Itoa(runsRead), &resp); err != nil {
		return runsView{Empty: &Empty{Title: "Runs are unavailable", Text: err.Error()}}
	}
	v := buildRuns(resp.Runs)
	if len(v.Live)+len(v.Recent) == 0 {
		v.Empty = &Empty{Title: "No runs yet", Text: "Runs appear here when one starts."}
	}
	return v
}

func init() {
	RegisterFragment("runs", &Fragment{
		Events:   []string{EventRun},
		Template: template.Must(template.New("runs").ParseFS(templateFS, "templates/partials.html", "templates/runs-list.html")).Lookup("runs-list"),
		Data:     runsData,
	})
}
