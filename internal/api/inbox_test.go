package api

import (
	"slices"
	"testing"

	"github.com/dphaener/rhiza/internal/run"
)

func TestUncoveredAcceptItems(t *testing.T) {
	const a, b = "01M7A1PXPS84ZD222S9ADHF0SV", "01M7A26CAT1FTRENTDNPR4XD18"
	tests := []struct {
		name  string
		items []InboxItem
		runs  []run.Journal
		want  []string
	}{
		{"no runs lists both", nil, nil, []string{a, b}},
		{"a live plan-review run takes the item away", nil, []run.Journal{proposalRun("r1", "plan-review", a, nil, "", "")}, []string{b}},
		{"a failed plan-review run takes the item away", nil, []run.Journal{proposalRun("r1", "plan-review", a, nil, run.Failed, "boom")}, []string{b}},
		{"an abandoned run does not", nil, []run.Journal{proposalRun("r1", "plan-review", a, nil, run.Failed, run.AbandonedPrefix+"the owner")}, []string{a, b}},
		{
			"a gated accept item from a non-plan run is not listed twice",
			[]InboxItem{{Kind: InboxAccept, Change: a, Detail: "waiting for acceptance"}},
			[]run.Journal{proposalRun("r1", "hold", "REC", []string{a}, "", "")},
			[]string{b},
		},
		{"an item of another kind does not hide it", []InboxItem{{Kind: InboxFailed, Change: a}}, nil, []string{a, b}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := uncoveredAcceptItems(tc.items, []string{a, b}, tc.runs)
			var changes []string
			for _, it := range got {
				if it.Kind != InboxAccept || it.Detail != "waiting for acceptance: no plan run covers it" || it.Next != ownerAcceptNext(it.Change) {
					t.Errorf("item %+v: want an accept item with the uncovered detail and next", it)
				}
				changes = append(changes, it.Change)
			}
			if !slices.Equal(changes, tc.want) {
				t.Errorf("uncoveredAcceptItems changes = %v, want %v", changes, tc.want)
			}
		})
	}
}
