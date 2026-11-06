package web

import (
	"html/template"
	"net/http"
)

// RunView is one live run on the Overview: its change set's title, its
// workflow and the step it is at. ID only builds the link; it is never shown
// (constraint owner-surfaces-no-ids).
type RunView struct {
	ID, Title, Workflow, Step string
	// State is running, waiting or stopped.
	State string
	// Started is when the run started, RFC 3339; empty if unknown.
	Started string
}

// RunsSource reads the live runs as they are now: those not finished.
type RunsSource func() []RunView

// runRows is the runs as dense rows, each linking to its run page.
func runRows(runs []RunView) []Row {
	rows := make([]Row, 0, len(runs))
	for _, r := range runs {
		text := r.State
		if r.State == "waiting" {
			text = "waiting on you"
		}
		rows = append(rows, Row{
			Href: "/runs/" + r.ID, Title: r.Title, Sub: r.Workflow,
			Pill: Pill{State: r.State, Text: text}, Step: r.Step, Live: r.State == "running",
			Datetime: r.Started,
		})
	}
	return rows
}

// RunsPanel is what the live runs section renders.
type RunsPanel struct {
	Rows []Row
}

type runsKey struct{}

func runsOf(r *http.Request) RunsPanel {
	if src, _ := r.Context().Value(runsKey{}).(RunsSource); src != nil {
		return RunsPanel{Rows: runRows(src())}
	}
	return RunsPanel{}
}

func init() {
	RegisterFragment("runs", &Fragment{
		Events:   []string{EventRun, EventDaemon},
		Template: template.Must(template.New("runs").ParseFS(templateFS, "templates/partials.html", "templates/runs.html")).Lookup("runs"),
		Data:     func(r *http.Request) any { return runsOf(r) },
	})
}
