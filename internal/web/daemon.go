package web

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
)

// DaemonView is what the Overview page shows of the daemon: whether it runs
// or is paused, its run slots and each part of the system's health.
type DaemonView struct {
	Running bool
	Paused  bool
	// ProjectPaused is whether the project the page is of is paused.
	ProjectPaused bool
	// CSRF is the field value of the pause forms, empty when the request
	// carried no web token: the panel then shows no form.
	CSRF string
	// Slots is how many runs hold a slot, of Limit.
	Slots, Limit int
	// Live and Starting split Slots: the live runs and the starts not
	// journaled yet, as `rhiza daemon status` counts them.
	Live, Starting int
	Health         []HealthPill
	// HealthPending is set until the first health report has landed.
	HealthPending bool
}

// HealthPill is one part of the health report with its state: ok, warn or
// bad. Summary is the line the pill's title carries.
type HealthPill struct{ Part, State, Summary string }

// DaemonSource reads the daemon as it is now.
type DaemonSource func() DaemonView

// State is the word the daemon panel leads with.
func (d DaemonView) State() string {
	switch {
	case !d.Running:
		return "Not running"
	case d.Paused:
		return "Paused"
	}
	return "Running"
}

// Heading is the panel's title.
func (d DaemonView) Heading() string {
	if !d.Running {
		return "Not running"
	}
	return fmt.Sprintf("%s, %d of %d slots in use", d.State(), d.Slots, d.Limit)
}

// RunsLine is the count in the form of daemon status and health.
func (d DaemonView) RunsLine() string {
	return fmt.Sprintf("runs: %d/%d live, %d starting", d.Live, d.Limit, d.Starting)
}

// SlotCells is one entry per slot, true where a run holds it.
func (d DaemonView) SlotCells() []bool {
	cells := make([]bool, max(d.Limit, 0))
	for i := range cells {
		cells[i] = i < d.Slots
	}
	return cells
}

// Bad is the number of health parts that are not ok.
func (d DaemonView) Bad() int {
	n := 0
	for _, h := range d.Health {
		if h.State != "ok" {
			n++
		}
	}
	return n
}

type daemonKey struct{}

// daemonOf is the DaemonView a request's handler serves.
func daemonOf(r *http.Request) DaemonView {
	if src, _ := r.Context().Value(daemonKey{}).(DaemonSource); src != nil {
		v := src()
		v.CSRF = CSRF(r)
		return v
	}
	return DaemonView{CSRF: CSRF(r)}
}

func init() {
	RegisterFragment("daemon", &Fragment{
		Events:   []string{EventDaemon, EventInbox, EventChange},
		Template: template.Must(template.New("daemon").ParseFS(templateFS, "templates/partials.html", "templates/daemon.html")).Lookup("daemon"),
		Data:     func(r *http.Request) any { return daemonOf(r) },
	})
}

// PausePath is where the panel's forms post, under the project's prefix.
func (d DaemonView) PausePath() string { return PausePath }

// Pill is the header's daemon pill.
func (d DaemonView) Pill() Pill {
	switch {
	case !d.Running:
		return Pill{State: "failed", Text: "daemon · off"}
	case d.Paused:
		return Pill{State: "waiting", Text: fmt.Sprintf("daemon · paused · %d/%d", d.Slots, d.Limit)}
	}
	return Pill{State: "running", Text: fmt.Sprintf("daemon · %d/%d", d.Slots, d.Limit)}
}

// WithDaemon serves h with daemon as the daemon every page's header pill
// reads, for pages served outside a project's handler (the root's projects
// home and merged inbox). A project's handler sets its own daemon over it.
func WithDaemon(daemon DaemonSource, h http.Handler) http.Handler {
	if daemon == nil {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), daemonKey{}, daemon)))
	})
}
