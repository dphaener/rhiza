package api

import (
	"encoding/json"
	"math/rand"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/id"
)

// showJSON is change.show of ref as the CLI prints it.
func showJSON(t *testing.T, ctx *Ctx, ref string) string {
	t.Helper()
	resp, err := callOp(t, ctx, "change.show", `{"ref":"`+ref+`"}`)
	if err != nil {
		t.Fatalf("change.show %s: %v", ref, err)
	}
	b, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// indexedLog is readLog with a draft change set that writes into the
// accepted one's tree and carries notes, a note on a note-bearing change
// set (current), and an unrelated draft.
func indexedLog(t *testing.T) (*Ctx, []id.ID) {
	t.Helper()
	ctx, landed, accepted := readLog(t)
	draft := newChange(t, ctx, `{"intent":"Later","title":"Later"}`).ChangeSet
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"flash","parent":"store"}`)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"ssd","parent":"store"}`)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"nvme","parent":"ssd"}`)
	mustWrite(t, ctx, "link", `{"from":"ssd","type":"depends-on","to":"flash"}`)
	mustWrite(t, ctx, "set", `{"ref":"cache","fields":["statement=Cache it"]}`)
	mustWrite(t, ctx, "link", `{"from":"flash","type":"serves","to":"faster"}`)
	if _, err := callOp(t, ctx, "note.add", `{"target":"`+accepted.String()+`","text":"A note about accepted"}`); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, ctx, "retire", `{"ref":"use-ssd","reason":"gone"}`)
	other := newChange(t, ctx, `{"intent":"Elsewhere"}`).ChangeSet
	if _, err := callOp(t, ctx, "change.use", `{"ref":"`+draft.String()+`"}`); err != nil {
		t.Fatal(err)
	}
	return ctx, []id.ID{landed, accepted, draft, other}
}

// TestChangeShowIndexedMatchesFold: with an index, change show answers
// what the fold of the log answers, in every state, and reads the log not
// at all, also after an entry is appended to an unrelated change set.
func TestChangeShowIndexedMatchesFold(t *testing.T) {
	ctx, sets := indexedLog(t)
	folded := *ctx
	folded.IndexPath = ""
	want := map[id.ID]string{}
	for _, cs := range sets {
		want[cs] = showJSON(t, &folded, cs.String())
	}
	for i, sub := range []string{`"state":"landed"`, `A note about accepted`, `"action":"updated"`, `"action":"retired"`} {
		if !strings.Contains(want[sets[min(i, 2)]]+want[sets[1]], sub) {
			t.Fatalf("thin fixture: no %s in %s", sub, want[sets[1]])
		}
	}
	reads := countingCtx(ctx)
	check := func(label string) {
		t.Helper()
		for _, cs := range sets {
			*reads = 0
			if got := showJSON(t, ctx, cs.String()); got != want[cs] {
				t.Errorf("%s: change show %s\nindexed %s\nfolded  %s", label, cs, got, want[cs])
			}
			if *reads != 0 {
				t.Errorf("%s: change show %s read the log %d times", label, cs, *reads)
			}
		}
	}
	check("fresh")
	// An entry in an unrelated change set: the index catches up, and the
	// answers about the others do not change.
	unrelated := sets[3]
	ctx.ChangeSet = unrelated
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"elsewhere-cap"}`)
	ctx.ChangeSet = id.ID{}
	for _, cs := range sets {
		want[cs] = showJSON(t, &folded, cs.String())
	}
	check("after an unrelated entry")
	// The current change set and the default change set resolve from the
	// index too.
	*reads = 0
	if _, err := callOp(t, ctx, "change.show", `{}`); err != nil {
		t.Fatal(err)
	}
	if *reads != 0 {
		t.Errorf("change show of the current change set read the log %d times", *reads)
	}
}

func TestChangeShowIndexedErrors(t *testing.T) {
	ctx, _ := indexedLog(t)
	folded := *ctx
	folded.IndexPath = ""
	for _, req := range []string{`{"ref":"ZZZZZZZZ"}`, `{"ref":"abc"}`, `{"ref":"01M6NE0XK2H70R52C5JAPWK71M"}`} {
		_, a := callOp(t, ctx, "change.show", req)
		_, b := callOp(t, &folded, "change.show", req)
		if (a == nil) != (b == nil) || (a != nil && a.Error() != b.Error()) {
			t.Errorf("%s: indexed error %v, folded error %v", req, a, b)
		}
	}
	ctx.IndexPath = filepath.Join(t.TempDir(), "other.db")
	if _, err := callOp(t, ctx, "change.show", `{}`); err != nil {
		t.Fatal(err)
	}
}

// TestShipCheckSettledReadsNoLog: ship check --settled finds its change
// set from the index and reads the log not at all.
func TestShipCheckSettledReadsNoLog(t *testing.T) {
	f, g, _ := shipFixture(t)
	f.ctx.IndexPath = filepath.Join(t.TempDir(), "index.db")
	g.answers[shipPRView] = openPR(greenRollup)
	reads := countingCtx(f.ctx)
	for _, label := range []string{"fresh", "after an unrelated entry"} {
		*reads = 0
		r := shipCheck(t, f, `{"settled":true}`)
		if !r.Passed || r.Change == "" {
			t.Errorf("%s: %+v", label, r)
		}
		if *reads != 0 {
			t.Errorf("%s: ship check --settled read the log %d times", label, *reads)
		}
		cur := currentOf(t, f.ctx)
		other := newChange(t, f.ctx, `{"intent":"Elsewhere"}`).ChangeSet
		mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"elsewhere-cap"}`)
		if _, err := callOp(t, f.ctx, "change.use", `{"ref":"`+cur.String()+`"}`); err != nil {
			t.Fatal(err)
		}
		_ = other
	}
}

// TestChangeShowIndexedKeepsMovedUnderLater: a node moved under one the
// change set creates after it keeps its place in change show over the
// index, as in the fold: its create no longer applies strictly, but the
// interim model keeps it, so the leaf is the moved node, not its new parent.
func TestChangeShowIndexedKeepsMovedUnderLater(t *testing.T) {
	ctx, _, _ := readLog(t)
	folded := *ctx
	folded.IndexPath = ""
	cs := newChange(t, ctx, `{"intent":"Edge"}`).ChangeSet
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"early"}`)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"late"}`)
	mustWrite(t, ctx, "mv", `{"ref":"early","parent":"late"}`)
	want := showJSON(t, &folded, cs.String())
	if !strings.Contains(want, `"handle":"early"`) {
		t.Fatalf("thin fixture: no early in %s", want)
	}
	if got := showJSON(t, ctx, cs.String()); got != want {
		t.Errorf("change show %s\nindexed %s\nfolded  %s", cs, got, want)
	}
}

// TestChangeShowIndexedMatchesFoldGenerated: over generated logs, change
// show of every change set answers through the index what it answers
// through the fold of the log.
func TestChangeShowIndexedMatchesFoldGenerated(t *testing.T) {
	seeds := int64(6)
	if testing.Short() {
		seeds = 2
	}
	for seed := int64(1); seed <= seeds; seed++ {
		ctx, landed, accepted := readLog(t)
		folded := *ctx
		folded.IndexPath = ""
		g := &writeGen{t: t, rng: rand.New(rand.NewSource(seed)), dir: t.TempDir(),
			handles: []string{"store", "cache", "disk", "tape", "old", "under-old", "faster", "use-ssd", "ram", "ghost"},
			sets:    []id.ID{landed, accepted}, base: 8}
		genLog(t, ctx, g)
		for _, cs := range g.sets {
			want := showJSON(t, &folded, cs.String())
			if got := showJSON(t, ctx, cs.String()); got != want {
				t.Errorf("seed %d: change show %s\nindexed %s\nfolded  %s", seed, cs, got, want)
			}
		}
	}
}
