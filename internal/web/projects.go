package web

import (
	"fmt"
	"net/http"
	"net/url"
)

// ProjectView is one registered project on the all-projects home: its live
// runs, inbox count and health. Name is the project's registered name, the
// owner's own word for it, never an ID.
type ProjectView struct {
	Name string
	// Initialised is false for a project with no .rhiza/ yet: it is listed,
	// marked not initialised, and has no runs, inbox or health.
	Initialised bool
	// Running and Waiting count the project's live runs by state.
	Running, Waiting int
	// Inbox is how many items wait on the owner in this project.
	Inbox  int
	Health []HealthPill
}

// ProjectsSource reads every registered project as it is now, by name.
type ProjectsSource func() []ProjectView

// Href is the project's overview, under the root's /p/<project>/.
func (p ProjectView) Href() string { return ProjectPrefix + url.PathEscape(p.Name) + "/" }

// Runs is the live-run count in words.
func (p ProjectView) Runs() string {
	switch n := p.Running + p.Waiting; {
	case n == 0:
		return "none live"
	case p.Waiting == 0:
		return fmt.Sprintf("%d running", p.Running)
	case p.Running == 0:
		return fmt.Sprintf("%d waiting on you", p.Waiting)
	default:
		return fmt.Sprintf("%d running, %d waiting on you", p.Running, p.Waiting)
	}
}

// HealthState is the worst state of the project's health parts: ok, warn or
// bad; empty when none is known.
func (p ProjectView) HealthState() string {
	state := ""
	for _, h := range p.Health {
		switch {
		case h.State == "bad":
			return "bad"
		case h.State == "warn":
			state = "warn"
		case state == "":
			state = "ok"
		}
	}
	return state
}

// HealthText is the health pill's word.
func (p ProjectView) HealthText() string {
	switch p.HealthState() {
	case "ok":
		return "healthy"
	case "warn":
		return "warning"
	case "bad":
		return "unhealthy"
	}
	return "unknown"
}

// ProjectsHome serves the all-projects home page at every path it is given.
func ProjectsHome(src ProjectsSource) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var projects []ProjectView
		if src != nil {
			projects = src()
		}
		render(w, "projects", page{Title: "Projects", Nav: "projects", Data: projects})
	})
}
