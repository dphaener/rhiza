package daemon

import (
	"slices"
	"testing"
)

func names(ps []Picked) []string {
	var out []string
	for _, p := range ps {
		out = append(out, p.Change)
	}
	return out
}

func TestRoundRobinAcrossProjects(t *testing.T) {
	projects := []ProjectPick{
		{Project: "A", Ready: []string{"A1", "A2", "A3"}},
		{Project: "B", Ready: []string{"B1", "B2"}},
		{Project: "C", Ready: []string{"C1"}},
	}
	if got := names(RoundRobin(projects, 0, 3)); !slices.Equal(got, []string{"A1", "B1", "C1"}) {
		t.Fatalf("maxRuns 3: %v", got)
	}
	if got := names(RoundRobin(projects, 0, 6)); !slices.Equal(got, []string{"A1", "B1", "C1", "A2", "B2", "A3"}) {
		t.Fatalf("maxRuns 6: %v", got)
	}
	if got := names(RoundRobin(projects, 2, 3)); !slices.Equal(got, []string{"A1"}) {
		t.Fatalf("two already running: %v", got)
	}
}

func TestRoundRobinPerProjectMax(t *testing.T) {
	projects := []ProjectPick{
		{Project: "A", Ready: []string{"A1", "A2", "A3"}, Max: 1},
		{Project: "B", Ready: []string{"B1", "B2"}},
		{Project: "C", Ready: []string{"C1"}},
	}
	if got := names(RoundRobin(projects, 0, 4)); !slices.Equal(got, []string{"A1", "B1", "C1", "B2"}) {
		t.Fatalf("A max 1: %v", got)
	}
	projects[0].Running = 1
	if got := names(RoundRobin(projects, 1, 4)); !slices.Equal(got, []string{"B1", "C1", "B2"}) {
		t.Fatalf("A already at its max: %v", got)
	}
}
