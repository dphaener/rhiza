package api

import (
	"github.com/dphaener/rhiza/internal/web"
)

// WebMilestonesSource reads the open milestones for the web UI's Overview
// from the repository base works in, as rhiza roadmap does: each milestone's
// change sets counted by state. A milestone is named by its statement.
func WebMilestonesSource(base *Ctx) web.MilestonesSource {
	return func() []web.MilestoneView {
		c := *base
		c.Record, c.DryRun, c.Out, c.JSON = false, false, nil, true
		c.Stdin, c.Rex, c.Start, c.Detach, c.Interrupts = nil, nil, nil, nil, nil
		out, err := runRoadmap(&c, RoadmapRequest{All: true})
		if err != nil {
			return nil
		}
		var ms []web.MilestoneView
		for _, g := range out.(RoadmapResponse).Groups {
			if g.Kind != "milestone" {
				continue
			}
			name := oneLine(g.Intent.Statement)
			if name == "" {
				name = "Milestone"
			}
			m := web.MilestoneView{Name: name, Counts: map[string]int{}}
			for _, cs := range g.ChangeSets {
				m.Counts[string(cs.State)]++
				m.Total++
			}
			ms = append(ms, m)
		}
		return ms
	}
}
