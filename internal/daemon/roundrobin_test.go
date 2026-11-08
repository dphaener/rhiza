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
	if got := names(RoundRobin(projects, 0, 3, "")); !slices.Equal(got, []string{"A1", "B1", "C1"}) {
		t.Fatalf("maxRuns 3: %v", got)
	}
	if got := names(RoundRobin(projects, 0, 6, "")); !slices.Equal(got, []string{"A1", "B1", "C1", "A2", "B2", "A3"}) {
		t.Fatalf("maxRuns 6: %v", got)
	}
	if got := names(RoundRobin(projects, 2, 3, "")); !slices.Equal(got, []string{"A1"}) {
		t.Fatalf("two already running: %v", got)
	}
}

func TestRoundRobinPerProjectMax(t *testing.T) {
	projects := []ProjectPick{
		{Project: "A", Ready: []string{"A1", "A2", "A3"}, Max: 1},
		{Project: "B", Ready: []string{"B1", "B2"}},
		{Project: "C", Ready: []string{"C1"}},
	}
	if got := names(RoundRobin(projects, 0, 4, "")); !slices.Equal(got, []string{"A1", "B1", "C1", "B2"}) {
		t.Fatalf("A max 1: %v", got)
	}
	projects[0].Running = 1
	if got := names(RoundRobin(projects, 1, 4, "")); !slices.Equal(got, []string{"B1", "C1", "B2"}) {
		t.Fatalf("A already at its max: %v", got)
	}
}

// Turns carry on from poll to poll: with one slot and every project ready,
// each freed slot goes to the next project, not back to the first.
func TestRoundRobinRotatesAcrossPolls(t *testing.T) {
	projects := []ProjectPick{
		{Project: "A", Ready: []string{"A1", "A2", "A3"}},
		{Project: "B", Ready: []string{"B1", "B2"}},
		{Project: "C", Ready: []string{"C1"}},
	}
	after := ""
	var got []string
	for poll := 0; poll < 5; poll++ {
		picked := RoundRobin(projects, 0, 1, after)
		after = LastTurn(picked, after)
		for _, pk := range picked {
			got = append(got, pk.Project)
			for i := range projects {
				if projects[i].Project == pk.Project {
					projects[i].Ready = projects[i].Ready[1:]
				}
			}
		}
	}
	if want := []string{"A", "B", "C", "A", "B"}; !slices.Equal(got, want) {
		t.Fatalf("turns across polls: %v, want %v", got, want)
	}
	// A pick that fills more than one slot carries on after its last turn.
	projects = []ProjectPick{
		{Project: "A", Ready: []string{"A1", "A2"}},
		{Project: "B", Ready: []string{"B1"}},
		{Project: "C", Ready: []string{"C1"}},
	}
	if got := names(RoundRobin(projects, 0, 2, "A")); !slices.Equal(got, []string{"B1", "C1"}) {
		t.Fatalf("after A: %v", got)
	}
	if got := names(RoundRobin(projects, 0, 4, "C")); !slices.Equal(got, []string{"A1", "B1", "C1", "A2"}) {
		t.Fatalf("after C: %v", got)
	}
	// An after no longer registered starts at the first project.
	if got := names(RoundRobin(projects, 0, 1, "gone")); !slices.Equal(got, []string{"A1"}) {
		t.Fatalf("after gone: %v", got)
	}
	// Nothing picked keeps the rotation where it was.
	if got := LastTurn(nil, "B"); got != "B" {
		t.Fatalf("LastTurn of nothing: %q", got)
	}
}

// A project's Running counts the starts not journaled yet along with its
// live runs, so they hold slots: two live-or-unjournaled runs fill two slots
// and a third slot gives exactly one pick.
func TestRoundRobinCountsUnjournaledStarts(t *testing.T) {
	ps := []ProjectPick{{Project: "a", Ready: []string{"c1", "c2"}, Running: 2}}
	if got := RoundRobin(ps, 2, 2, ""); len(got) != 0 {
		t.Errorf("max 2: picked %v, want none", got)
	}
	if got := RoundRobin(ps, 2, 3, ""); len(got) != 1 || got[0].Change != "c1" {
		t.Errorf("max 3: picked %v, want c1 alone", got)
	}
}
