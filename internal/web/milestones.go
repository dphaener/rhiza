package web

import (
	"html/template"
	"net/http"
)

// MilestoneView is one open milestone: its statement and its change sets
// counted by state. It carries no ID (constraint owner-surfaces-no-ids).
type MilestoneView struct {
	Name   string
	Counts map[string]int
	Total  int
}

// MilestonesSource reads the open milestones as they are now.
type MilestonesSource func() []MilestoneView

// milestoneStates are the states counted, in order.
var milestoneStates = []string{"accepted", "proposed", "draft", "landed"}

// MilestoneRow is a milestone ready to render.
type MilestoneRow struct {
	Name   string
	Total  int
	Counts []MilestoneCount
}

// MilestoneCount is one state's count.
type MilestoneCount struct {
	Pill  Pill
	Count int
}

// MilestonesPanel is what the milestones section renders.
type MilestonesPanel struct {
	Milestones []MilestoneRow
}

type milestonesKey struct{}

func milestonesOf(r *http.Request) MilestonesPanel {
	src, _ := r.Context().Value(milestonesKey{}).(MilestonesSource)
	if src == nil {
		return MilestonesPanel{}
	}
	var p MilestonesPanel
	for _, m := range src() {
		row := MilestoneRow{Name: m.Name, Total: m.Total}
		for _, s := range milestoneStates {
			if n := m.Counts[s]; n > 0 {
				row.Counts = append(row.Counts, MilestoneCount{Pill: Pill{State: s, Text: s}, Count: n})
			}
		}
		p.Milestones = append(p.Milestones, row)
	}
	return p
}

func init() {
	RegisterFragment("milestones", &Fragment{
		Events:   []string{EventChange},
		Template: template.Must(template.New("milestones").ParseFS(templateFS, "templates/partials.html", "templates/milestones.html")).Lookup("milestones"),
		Data:     func(r *http.Request) any { return milestonesOf(r) },
	})
}
