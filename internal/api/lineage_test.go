package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// lineageLog builds a log for why and log. Landed change set one: store
// with child cache; decisions d0-d4 and dx; sources s1 (summary, quotes)
// and s2; question q1; cache decided-by d1 and dx, derived-from s2,
// relates-to d9; store decided-by d4; d1 decided-by d2 and d2 decided-by
// d1 (a cycle); d2 derived-from s1; d3 supersedes d1, d1 supersedes d0, d0
// supersedes d3 (a cycle); d1 answers q1. Landed change set two (sources
// s1): retires cache's decided-by dx link, updates cache. Accepted change
// set three: updates cache again. It returns the ctx, the three change
// sets and the dx link.
func lineageLog(t *testing.T) (*Ctx, []id.ID, id.ID) {
	t.Helper()
	ctx := readCtx(t)
	one := newChange(t, ctx, `{"intent":"Base"}`).ChangeSet
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"cache","parent":"store"}`)
	for _, d := range []string{"d0", "d1", "d2", "d3", "d4", "dx", "d9"} {
		mustWrite(t, ctx, "add", fmt.Sprintf(`{"kind":"decision","handle":%q,"fields":["statement=Because %s"]}`, d, d))
	}
	mustWrite(t, ctx, "add", `{"kind":"source","handle":"s1","fields":["summary=The talk"],"json_fields":["quotes=[\"we must\",\"we shall\"]"]}`)
	mustWrite(t, ctx, "add", `{"kind":"source","handle":"s2"}`)
	mustWrite(t, ctx, "add", `{"kind":"question","handle":"q1"}`)
	mustWrite(t, ctx, "link", `{"from":"cache","type":"decided-by","to":"d1"}`)
	dx := mustWrite(t, ctx, "link", `{"from":"cache","type":"decided-by","to":"dx"}`).Links[0]
	mustWrite(t, ctx, "link", `{"from":"cache","type":"derived-from","to":"s2"}`)
	mustWrite(t, ctx, "link", `{"from":"cache","type":"relates-to","to":"d9","note":"near"}`)
	mustWrite(t, ctx, "link", `{"from":"store","type":"decided-by","to":"d4"}`)
	mustWrite(t, ctx, "link", `{"from":"d1","type":"decided-by","to":"d2"}`)
	mustWrite(t, ctx, "link", `{"from":"d2","type":"decided-by","to":"d1"}`)
	mustWrite(t, ctx, "link", `{"from":"d2","type":"derived-from","to":"s1"}`)
	mustWrite(t, ctx, "link", `{"from":"d3","type":"supersedes","to":"d1"}`)
	mustWrite(t, ctx, "link", `{"from":"d1","type":"supersedes","to":"d0"}`)
	mustWrite(t, ctx, "link", `{"from":"d0","type":"supersedes","to":"d3"}`)
	mustWrite(t, ctx, "link", `{"from":"d1","type":"answers","to":"q1"}`)
	move(t, ctx, "propose", `{}`)
	move(t, ctx, "accept", `{"reason":"ok"}`)
	move(t, ctx, "land", `{}`)

	two := newChange(t, ctx, `{"intent":"Prune","sources":["s1"]}`).ChangeSet
	mustWrite(t, ctx, "unlink", fmt.Sprintf(`{"ref":%q}`, dx.String()))
	mustWrite(t, ctx, "set", `{"ref":"cache","fields":["statement=Fast"]}`)
	move(t, ctx, "propose", `{}`)
	move(t, ctx, "accept", `{"reason":"ok"}`)
	move(t, ctx, "land", `{}`)

	three := newChange(t, ctx, `{"intent":"Faster"}`).ChangeSet
	mustWrite(t, ctx, "set", `{"ref":"cache","fields":["statement=Faster"]}`)
	move(t, ctx, "propose", `{}`)
	move(t, ctx, "accept", `{"reason":"ok"}`)
	return ctx, []id.ID{one, two, three}, dx
}

func whyOf(t *testing.T, ctx *Ctx, req string) WhyResponse {
	t.Helper()
	resp, err := callOp(t, ctx, "why", req)
	if err != nil {
		t.Fatalf("why %s: %v", req, err)
	}
	return resp.(WhyResponse)
}

func logOf(t *testing.T, ctx *Ctx, req string) LogResponse {
	t.Helper()
	resp, err := callOp(t, ctx, "log", req)
	if err != nil {
		t.Fatalf("log %s: %v", req, err)
	}
	return resp.(LogResponse)
}

// whyLine is a decision as handle/relation/via/inherited.
func whyLine(d WhyDecision) string {
	return fmt.Sprintf("%s %s %s %v", d.Via.Handle, d.Relation, d.Handle, d.Inherited)
}

func TestWhyWalk(t *testing.T) {
	ctx, css, _ := lineageLog(t)
	r := whyOf(t, ctx, `{"ref":"cache","view":"as-is"}`)
	if r.Node.Handle != "cache" || r.Node.Kind != op.KindCapability || r.View != "as-is" {
		t.Errorf("node %+v view %q", r.Node, r.View)
	}
	var got []string
	for _, d := range r.Decisions {
		got = append(got, whyLine(d))
	}
	want := []string{
		"cache decided by d1 false",
		"store decided by d4 true",
		"d1 decided by d2 false",
		"d1 superseded by d3 false",
		"d1 supersedes d0 false",
	}
	if !slices.Equal(got, want) {
		t.Errorf("decisions:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if a := r.Decisions[0].Answers; len(a) != 1 || a[0].Handle != "q1" {
		t.Errorf("d1 answers %+v, want q1", a)
	}
	if a := r.Decisions[1].Answers; a == nil || len(a) != 0 {
		t.Errorf("d4 answers %#v, want []", a)
	}
	if got := fieldText(r.Decisions[0].Fields["statement"]); got != "Because d1" {
		t.Errorf("d1 statement %q", got)
	}
	var srcs []string
	for _, s := range r.Sources {
		srcs = append(srcs, s.Via.Handle+" "+s.Handle)
	}
	if want := []string{"cache s2", "d2 s1"}; !slices.Equal(srcs, want) {
		t.Errorf("sources %v, want %v", srcs, want)
	}
	if got := fieldText(r.Sources[1].Fields["quotes"]); got != "[\n  \"we must\",\n  \"we shall\"\n]" {
		t.Errorf("s1 quotes %q", got)
	}
	if r.Sources[0].Fields == nil {
		t.Error("s2 fields nil, want {}")
	}
	if len(r.ChangeSets) != 2 || r.ChangeSets[0].ID != css[0] || r.ChangeSets[1].ID != css[1] {
		t.Fatalf("change sets %+v, want %v", r.ChangeSets, css[:2])
	}
	c0, c1 := r.ChangeSets[0], r.ChangeSets[1]
	if c0.State != changeset.StateLanded || c0.Kind != changeset.KindIntention || c0.Intent != "Base" || len(c0.Sources) != 0 || c0.Sources == nil {
		t.Errorf("change set one %+v", c0)
	}
	s1 := r.Sources[1].ID
	if c1.Intent != "Prune" || !slices.Equal(c1.Sources, []id.ID{s1}) {
		t.Errorf("change set two %+v, want sources [%s]", c1, s1)
	}

	// to-be adds the accepted change set.
	tb := whyOf(t, ctx, `{"ref":"cache","view":"to-be"}`)
	if len(tb.ChangeSets) != 3 || tb.ChangeSets[2].ID != css[2] || tb.ChangeSets[2].State != changeset.StateAccepted {
		t.Errorf("to-be change sets %+v", tb.ChangeSets)
	}

	var b bytes.Buffer
	r.Human(&b)
	out := b.String()
	for _, s := range []string{"# Why cache\n", "## Decisions\n", "### d1\n", "- cache decided by d1\n",
		"- store decided by d4 (inherited)\n", "- statement: Because d2\n", "- Answers: q1 (",
		"## Sources\n", "- summary: The talk\n", "- quotes: [\n  \"we must\",\n  \"we shall\"\n]\n", "## Change sets\n",
		"intention: Prune (sources: " + s1.String() + ")"} {
		if !strings.Contains(out, s) {
			t.Errorf("Human lacks %q:\n%s", s, out)
		}
	}
	if strings.Contains(out, "d9") || strings.Contains(out, "dx") {
		t.Errorf("Human walked relates-to or a retired link:\n%s", out)
	}
}

func TestWhyFromDecision(t *testing.T) {
	ctx, _, _ := lineageLog(t)
	r := whyOf(t, ctx, `{"ref":"d1","view":"as-is"}`)
	var got []string
	for _, d := range r.Decisions {
		got = append(got, whyLine(d))
	}
	want := []string{"d1 decided by d2 false", "d1 superseded by d3 false", "d1 supersedes d0 false"}
	if !slices.Equal(got, want) {
		t.Errorf("decisions %q, want %q", got, want)
	}
	if len(r.Sources) != 1 || r.Sources[0].Handle != "s1" {
		t.Errorf("sources %+v", r.Sources)
	}
	// The decision asked about leads, with its own fields: its quote is the
	// first thing why gives back.
	if got := fieldText(r.Node.Fields["statement"]); got != "Because d1" {
		t.Errorf("node statement %q, want %q", got, "Because d1")
	}
	var b bytes.Buffer
	r.Human(&b)
	out := b.String()
	if k, d := strings.Index(out, "- statement: Because d1\n"), strings.Index(out, "## Decisions\n"); k < 0 || k > d {
		t.Errorf("Human does not lead with d1's fields:\n%s", out)
	}

	// From d0, the supersedes cycle d0 -> d3 -> d1 -> d0 ends.
	r = whyOf(t, ctx, `{"ref":"d0","view":"as-is"}`)
	got = nil
	for _, d := range r.Decisions {
		got = append(got, whyLine(d))
	}
	want = []string{"d0 superseded by d1 false", "d0 supersedes d3 false", "d1 decided by d2 false"}
	if !slices.Equal(got, want) {
		t.Errorf("from d0 %q, want %q", got, want)
	}
}

func TestWhyEmptyAndErrors(t *testing.T) {
	ctx, _, _ := lineageLog(t)
	r := whyOf(t, ctx, `{"ref":"q1","view":"as-is"}`)
	if r.Decisions == nil || r.Sources == nil || len(r.Decisions)+len(r.Sources) != 0 {
		t.Errorf("q1 %+v", r)
	}
	var b bytes.Buffer
	r.Human(&b)
	if strings.Count(b.String(), "None.") != 2 {
		t.Errorf("Human:\n%s", b.String())
	}
	_, err := callOp(t, ctx, "why", `{"ref":"nosuch"}`)
	wantExit(t, err, 4)
	_, err = callOp(t, ctx, "why", `{"ref":"cache","view":"bogus"}`)
	wantExit(t, err, 2)
}

var lineageAtRe = regexp.MustCompile(`^\d{4}-\d\d-\d\dT\d\d:\d\d:\d\dZ$`)

func TestLogNode(t *testing.T) {
	ctx, css, _ := lineageLog(t)
	r := logOf(t, ctx, `{"ref":"cache","view":"to-be"}`)
	if r.Target.Kind != "node" || r.Target.Handle != "cache" || r.View != "to-be" {
		t.Errorf("target %+v view %q", r.Target, r.View)
	}
	if len(r.ChangeSets) != 3 {
		t.Fatalf("change sets %+v", r.ChangeSets)
	}
	wantTypes := []op.Type{op.TypeNodeCreate, op.TypeNodeUpdate, op.TypeNodeUpdate}
	for k, c := range r.ChangeSets {
		if c.ID != css[k] || len(c.Ops) != 1 || c.Ops[0].Type != wantTypes[k] || c.Ops[0].ID.IsZero() {
			t.Errorf("change set %d: %+v", k, c)
		}
		landed := k < 2
		if landed != (c.State == changeset.StateLanded) || landed != lineageAtRe.MatchString(c.At) || landed == c.Actor.IsZero() {
			t.Errorf("change set %d: state %s at %q actor %s", k, c.State, c.At, c.Actor)
		}
		if c.Events != nil {
			t.Errorf("change set %d never reopened has events %+v", k, c.Events)
		}
	}
	if r.ChangeSets[0].At > r.ChangeSets[1].At {
		t.Errorf("land times out of order: %s, %s", r.ChangeSets[0].At, r.ChangeSets[1].At)
	}
	asis := logOf(t, ctx, `{"ref":"cache","view":"as-is"}`)
	if len(asis.ChangeSets) != 2 {
		t.Errorf("as-is change sets %+v", asis.ChangeSets)
	}

	var b bytes.Buffer
	r.Human(&b)
	out := b.String()
	for _, s := range []string{"# Log cache\n", "## " + css[0].String() + " landed intention: Base\n",
		"- Landed " + r.ChangeSets[0].At + " by " + r.ChangeSets[0].Actor.String() + "\n",
		"- node.create " + r.ChangeSets[0].Ops[0].ID.String() + "\n",
		"## " + css[2].String() + " accepted intention: Faster\n"} {
		if !strings.Contains(out, s) {
			t.Errorf("Human lacks %q:\n%s", s, out)
		}
	}
}

func TestLogReopened(t *testing.T) {
	ctx := readCtx(t)
	cs := relanded(t, ctx)
	r := logOf(t, ctx, `{"ref":"store"}`)
	if len(r.ChangeSets) != 1 {
		t.Fatalf("change sets %+v", r.ChangeSets)
	}
	c := r.ChangeSets[0]
	if c.ID != cs || c.State != changeset.StateLanded || len(c.Events) != 3 {
		t.Fatalf("change set %+v", c)
	}
	var lines string
	for k, e := range c.Events {
		want := []entry.Type{entry.TypeLand, entry.TypeReopen, entry.TypeLand}[k]
		if e.Type != want || e.Entry.IsZero() || e.Actor.IsZero() || !lineageAtRe.MatchString(e.At) {
			t.Errorf("event %d = %+v, want %s", k, e, want)
		}
		verb := map[entry.Type]string{entry.TypeLand: "Landed", entry.TypeReopen: "Reopened"}[e.Type]
		lines += "- " + verb + " " + e.At + " by " + e.Actor.String() + "\n"
	}
	// The existing fields still name the last land.
	if last := c.Events[2]; c.At != last.At || c.Actor != last.Actor {
		t.Errorf("at %q actor %s, want the last land %+v", c.At, c.Actor, last)
	}
	var b bytes.Buffer
	r.Human(&b)
	if out := b.String(); !strings.Contains(out, " landed intention: X\n\n"+lines+"- ") {
		t.Errorf("Human lacks\n%s in:\n%s", lines, out)
	}
}

func TestLogReopenedAccepted(t *testing.T) {
	// Reopened and not yet relanded: one Landed and one Reopened line, no
	// stale at or actor.
	ctx := readCtx(t)
	cs := landed(t, ctx)
	ctx.Git = branchGit("").run
	ctx.ChangeRef = cs.String()
	move(t, ctx, "reopen", `{}`)
	r := logOf(t, ctx, `{"ref":"store","view":"to-be"}`)
	if len(r.ChangeSets) != 1 {
		t.Fatalf("change sets %+v", r.ChangeSets)
	}
	c := r.ChangeSets[0]
	if c.State != changeset.StateAccepted || c.At != "" || !c.Actor.IsZero() || len(c.Events) != 2 {
		t.Fatalf("change set %+v", c)
	}
	var b bytes.Buffer
	r.Human(&b)
	out := b.String()
	want := "- Landed " + c.Events[0].At + " by " + c.Events[0].Actor.String() + "\n- Reopened " +
		c.Events[1].At + " by " + c.Events[1].Actor.String() + "\n"
	if !strings.Contains(out, want) || strings.Count(out, "- Landed") != 1 {
		t.Errorf("Human lacks %q:\n%s", want, out)
	}
}

func TestLogLink(t *testing.T) {
	ctx, css, dx := lineageLog(t)
	r := logOf(t, ctx, fmt.Sprintf(`{"ref":%q,"view":"as-is"}`, dx.String()[:20]))
	if r.Target.Kind != "link" || r.Target.ID != dx || r.Target.Handle != "" {
		t.Errorf("target %+v", r.Target)
	}
	if len(r.ChangeSets) != 2 || r.ChangeSets[0].ID != css[0] || r.ChangeSets[1].ID != css[1] ||
		r.ChangeSets[0].Ops[0].Type != op.TypeLinkCreate || r.ChangeSets[1].Ops[0].Type != op.TypeLinkRetire {
		t.Errorf("change sets %+v", r.ChangeSets)
	}
	var b bytes.Buffer
	r.Human(&b)
	if !strings.HasPrefix(b.String(), "# Log link "+dx.String()+"\n") {
		t.Errorf("Human:\n%s", b.String())
	}
	_, err := callOp(t, ctx, "log", `{"ref":"nosuch"}`)
	wantExit(t, err, 4)
	if e := AsError(err); e.Next != lsNext {
		t.Errorf("next %q", e.Next)
	}
}

// carveByHand moves node handle's create from change set from to change set
// to the way a carve does: one entry in from dropping the create, one in to
// re-creating it under a new op ID, each naming the other. It returns the
// source's and the target's entry IDs.
func carveByHand(t *testing.T, ctx *Ctx, handle string, from, to id.ID) (id.ID, id.ID) {
	t.Helper()
	w := ctxWorkspace(t, ctx)
	src, _ := w.changeSet(from)
	tgt, _ := w.changeSet(to)
	var create op.Op
	for _, o := range src.Ops {
		if o.Type == op.TypeNodeCreate && o.Handle == handle {
			create = o
		}
	}
	moved := create
	moved.ID = newID(t, ctx)
	srcEntry, tgtEntry := newID(t, ctx), newID(t, ctx)
	actor := entriesOf(w.log(), from)[0].Actor
	write := func(cs, eid id.ID, after []id.ID, body changeset.Ops) {
		raw, err := changeset.EncodeBody(body)
		if err != nil {
			t.Fatal(err)
		}
		e := entry.Entry{V: changeset.EntryVersion(body), ID: eid, ChangeSet: cs, Actor: actor, For: actor,
			At: ctx.Now().UTC().Truncate(time.Second), After: after, Type: body.Type(), Body: raw}
		if err := w.st.Write(e); err != nil {
			t.Fatal(err)
		}
	}
	write(from, srcEntry, src.Heads, changeset.Ops{Ops: []op.Op{}, Drop: []id.ID{create.ID},
		Carve: &changeset.Carve{To: to, Entry: tgtEntry}})
	write(to, tgtEntry, tgt.Heads, changeset.Ops{Ops: []op.Op{moved}, Drop: []id.ID{},
		Carve: &changeset.Carve{From: from, Entry: srcEntry}})
	return srcEntry, tgtEntry
}

func newID(t *testing.T, ctx *Ctx) id.ID {
	t.Helper()
	v, err := ctx.IDs.New()
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func TestLogCarve(t *testing.T) {
	ctx := readCtx(t)
	src := newChange(t, ctx, `{"intent":"Bundle"}`).ChangeSet
	mustWrite(t, ctx, "add", `{"kind":"decision","handle":"moved"}`)
	mustWrite(t, ctx, "add", `{"kind":"decision","handle":"stays"}`)
	tgt := newChange(t, ctx, `{"intent":"Milestone"}`).ChangeSet
	mustWrite(t, ctx, "add", `{"kind":"decision","handle":"other"}`)
	srcEntry, tgtEntry := carveByHand(t, ctx, "moved", src, tgt)

	r := logOf(t, ctx, fmt.Sprintf(`{"ref":"moved","view":"what-if:%s"}`, tgt))
	if len(r.ChangeSets) != 2 || r.ChangeSets[0].ID != src || r.ChangeSets[1].ID != tgt {
		t.Fatalf("change sets %+v, want %s then %s", r.ChangeSets, src, tgt)
	}
	s, g := r.ChangeSets[0], r.ChangeSets[1]
	if s.Intent != "Bundle" || s.State != changeset.StateDraft || s.Ops == nil || len(s.Ops) != 0 {
		t.Errorf("source %+v, want draft Bundle with no ops", s)
	}
	if want := []LogCarve{{Entry: srcEntry, CarvedInto: tgt, CarveEntry: tgtEntry}}; !slices.Equal(s.Carves, want) {
		t.Errorf("source carves %+v, want %+v", s.Carves, want)
	}
	if len(g.Ops) != 1 || g.Ops[0].Type != op.TypeNodeCreate {
		t.Errorf("target ops %+v", g.Ops)
	}
	if want := []LogCarve{{Entry: tgtEntry, CarvedFrom: src, CarveEntry: srcEntry}}; !slices.Equal(g.Carves, want) {
		t.Errorf("target carves %+v, want %+v", g.Carves, want)
	}
	var b bytes.Buffer
	r.Human(&b)
	line := "- Carved from " + shortID(src) + " into " + shortID(tgt) + "\n"
	if strings.Count(b.String(), line) != 2 {
		t.Errorf("Human lacks %q twice:\n%s", line, b.String())
	}

	// A node the carve did not move logs as before: no carves.
	r = logOf(t, ctx, fmt.Sprintf(`{"ref":"stays","view":"what-if:%s"}`, src))
	if len(r.ChangeSets) != 1 || r.ChangeSets[0].ID != src || r.ChangeSets[0].Carves != nil {
		t.Errorf("stays %+v", r.ChangeSets)
	}
	b.Reset()
	r.Human(&b)
	if strings.Contains(b.String(), "Carved") {
		t.Errorf("Human:\n%s", b.String())
	}
}

func TestWhyCitedQuotes(t *testing.T) {
	ctx := readCtx(t)
	newChange(t, ctx, `{"intent":"Base"}`)
	mustWrite(t, ctx, "add", `{"kind":"source","handle":"talk","json_fields":["quotes=[\"we must  keep\\nit small\",\"we shall grow\"]"]}`)
	mustWrite(t, ctx, "add", `{"kind":"decision","handle":"small","fields":["quote=must keep it"]}`)
	mustWrite(t, ctx, "link", `{"from":"small","type":"derived-from","to":"talk"}`)
	move(t, ctx, "propose", `{}`)
	move(t, ctx, "accept", `{"reason":"ok"}`)
	move(t, ctx, "land", `{}`)

	r := whyOf(t, ctx, `{"ref":"small","view":"as-is"}`)
	if len(r.Sources) != 1 {
		t.Fatalf("sources %+v", r.Sources)
	}
	var got []string
	if err := json.Unmarshal(r.Sources[0].Fields["quotes"], &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "we must  keep\nit small" {
		t.Errorf("quotes %q, want only the cited one", got)
	}
	if r.Sources[0].QuotesLeftOut != 1 {
		t.Errorf("left out %d, want 1", r.Sources[0].QuotesLeftOut)
	}
}

func TestWhyCitedQuotesKeepsAll(t *testing.T) {
	src := map[string]json.RawMessage{"quotes": json.RawMessage(`["a b","c d"]`)}
	for name, dec := range map[string]map[string]json.RawMessage{
		"no match": {"quote": json.RawMessage(`"zzz"`)},
		"no quote": {},
		"empty":    {"quote": json.RawMessage(`""`)},
	} {
		got, n := whyCitedQuotes(src, dec)
		if n != 0 || string(got["quotes"]) != `["a b","c d"]` {
			t.Errorf("%s: got %s, left out %d", name, got["quotes"], n)
		}
	}
}
