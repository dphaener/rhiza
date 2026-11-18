package api

import (
	"bytes"
	"slices"
	"strings"
	"testing"
	"time"

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

// A live run's open spawn-held entries are one spawn-held item each, a
// spawned child's naming its role and the step session's its step, both
// with the entry's detail; a hold released by its spawn or the step's
// session entry, or one in an ended or unheld run, gives none, and a hold
// keeps its one ID.
func TestInboxSpawnHeldItems(t *testing.T) {
	t.Parallel()
	at := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	const load = "machine hold: load-high: 1-minute load 4.00 per core, over the limit of 3.00"
	const mem = "machine hold: memory-low: 1.0 GiB available, under the minimum of 2.0 GiB"
	journal := func(extra ...run.Entry) run.Journal {
		es := []run.Entry{
			{Type: run.TypeStarted, At: at, Step: "build", Started: &run.Started{Change: "CS1", Workflow: "conduct"}},
			{Type: run.TypeEntered, At: at, Step: "build", Visit: 1},
			{Type: run.TypeSpawnHeld, At: at.Add(time.Second), Step: "build", Detail: mem},
			{Type: run.TypeSession, At: at.Add(2 * time.Second), Step: "build", Session: "P1"},
			{Type: run.TypeSpawnHeld, At: at.Add(3 * time.Second), Step: "build", Parent: "P1", Role: "verifier", Detail: load},
		}
		return run.Journal{ID: "R1", Entries: append(es, extra...)}
	}
	items := func(j run.Journal, live bool) []InboxItem {
		js := []run.Journal{j}
		return spawnHeldItems(js, runStates(js), func(string) bool { return live })
	}

	got := items(journal(), true)
	if len(got) != 1 {
		t.Fatalf("items %+v, want the verifier's spawn only", got)
	}
	it := got[0]
	if it.Kind != InboxSpawnHeld || it.Run != "R1" || it.Change != "CS1" || it.Step != "build" || it.Session != "P1" || it.At != "2030-01-01T00:00:03Z" {
		t.Errorf("item %+v", it)
	}
	if want := "a verifier session spawned by P1 waits: " + load; it.Detail != want {
		t.Errorf("detail %q, want %q", it.Detail, want)
	}
	if again := items(journal(), true); again[0].ID != it.ID {
		t.Errorf("ID %s then %s, want one ID per hold", it.ID, again[0].ID)
	}

	step := run.Journal{ID: "R1", Entries: journal().Entries[:3]}
	if got := items(step, true); len(got) != 1 || got[0].Session != "" || got[0].Detail != "the step's session waits: "+mem {
		t.Errorf("step session hold: %+v", got)
	}

	spawned := journal(run.Entry{Type: run.TypeSpawn, At: at.Add(time.Minute), Step: "build", Parent: "P1", Session: "C1", Role: "verifier"})
	if got := items(spawned, true); len(got) != 0 {
		t.Errorf("released by its spawn: %+v, want none", got)
	}
	ended := journal(run.Entry{Type: run.TypeFinished, At: at.Add(time.Minute), End: run.Failed})
	if got := items(ended, true); len(got) != 0 {
		t.Errorf("ended run: %+v, want none", got)
	}
	if got := items(journal(), false); len(got) != 0 {
		t.Errorf("run no process holds: %+v, want none", got)
	}

	var out bytes.Buffer
	InboxResponse{Items: []InboxItem{it}}.Human(&out)
	for _, want := range []string{"spawn-held run R1, change set CS1, step build", "a verifier session spawned by P1 waits: " + load} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("human output %q lacks %q", out.String(), want)
		}
	}
}
