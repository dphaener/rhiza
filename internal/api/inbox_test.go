package api

import (
	"slices"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/run"
)

func TestUncoveredAcceptItems(t *testing.T) {
	t.Parallel()
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

// TestInboxBlockersListProposedNotDraft builds a log whose proposed change
// set depends on a landed one and is blocked on a node it touches and on
// itself, beside a draft blocked the same way: inboxBlockers, which reads
// no dependency sets, lists the proposed set's two blockers with their
// statements, as planGather and planCore's blockers say, and never the
// draft's.
func TestInboxBlockersListProposedNotDraft(t *testing.T) {
	t.Parallel()
	ctx := readCtx(t)
	recordWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
	recordWrite(t, ctx, "add", `{"kind":"capability","handle":"cache"}`)
	l := newChange(t, ctx, `{"intent":"Landed"}`).ChangeSet
	mustWrite(t, ctx, "set", `{"ref":"store","fields":["statement=Stored"]}`)
	move(t, ctx, "propose", `{}`)
	move(t, ctx, "accept", `{"reason":"ok"}`)
	move(t, ctx, "land", `{}`)

	p := newChange(t, ctx, `{"intent":"Proposed"}`).ChangeSet
	planLink(t, ctx, op.LinkDependsOn, p, l)
	mustWrite(t, ctx, "set", `{"ref":"store","fields":["statement=Kept"]}`)
	move(t, ctx, "propose", `{}`)
	recordWrite(t, ctx, "add", `{"kind":"question","handle":"store-q","fields":["statement=Which store?"]}`)
	recordWrite(t, ctx, "link", `{"from":"store-q","type":"blocks","to":"store"}`)
	recordWrite(t, ctx, "add", `{"kind":"question","handle":"own-q","fields":["statement=Ship it?"]}`)
	planLink(t, ctx, op.LinkBlocks, nodeAt(t, ctx, "own-q").ID, p)

	d := newChange(t, ctx, `{"intent":"Drafty"}`).ChangeSet
	mustWrite(t, ctx, "set", `{"ref":"cache","fields":["statement=Cached"]}`)
	recordWrite(t, ctx, "add", `{"kind":"question","handle":"cache-q","fields":["statement=Which cache?"]}`)
	recordWrite(t, ctx, "link", `{"from":"cache-q","type":"blocks","to":"cache"}`)
	planLink(t, ctx, op.LinkBlocks, nodeAt(t, ctx, "cache-q").ID, d)

	// What the full plan path says blocks each non-draft change set.
	r, err := openReader(ctx, "inbox", "", toBeView)
	if err != nil {
		t.Fatal(err)
	}
	full, err := planGather(r)
	if err != nil {
		t.Fatal(err)
	}
	var want []string
	for _, c := range planCore(full).order {
		for _, b := range c.blockers {
			if c.info.State != changeset.StateDraft {
				want = append(want, c.info.ID.String()+" "+planHandle(b.question, full.nodes)+" "+planHandle(b.target, full.nodes))
			}
		}
	}
	lean, err := planGatherBlockers(r)
	r.Close()
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range lean.sets {
		if s.State == changeset.StateLanded || s.ID.IsZero() {
			t.Errorf("planGatherBlockers read set %s (%s): want the in-flight sets alone", s.ID, s.State)
		}
	}

	items, err := inboxBlockers(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	texts := map[string]InboxItem{}
	for _, it := range items {
		blocks := it.Blocks
		if blocks == "this change set" {
			blocks = planHandle(p, full.nodes)
		}
		got = append(got, it.Change+" "+it.Handle+" "+blocks)
		texts[it.Handle] = it
	}
	slices.Sort(got)
	slices.Sort(want)
	if len(want) != 2 || !slices.Equal(got, want) {
		t.Errorf("blockers = %q, want %q (two, from the full plan path)", got, want)
	}
	if it := texts["store-q"]; it.Change != p.String() || it.Question != "Which store?" || it.Blocks != "store" || it.BlocksTitle != "Stored" || it.Next != "rhiza show store-q" {
		t.Errorf("store-q item = %+v", it)
	}
	if it := texts["own-q"]; it.Change != p.String() || it.Question != "Ship it?" || it.Blocks != "this change set" || it.BlocksTitle != "" {
		t.Errorf("own-q item = %+v", it)
	}
	if _, ok := texts["cache-q"]; ok {
		t.Errorf("the draft's blocker cache-q is listed: %+v", texts["cache-q"])
	}
}
