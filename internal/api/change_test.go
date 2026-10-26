package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/store"
)

// initCtx runs init in a temp folder and returns a Ctx over it.
func initCtx(t *testing.T) *Ctx {
	t.Helper()
	ctx := testCtx(t.TempDir())
	if _, err := callOp(t, ctx, "init", `{}`); err != nil {
		t.Fatal(err)
	}
	return ctx
}

// ctxWorkspace opens the workspace of ctx, for building logs in tests.
func ctxWorkspace(t *testing.T, ctx *Ctx) *workspace {
	t.Helper()
	w, err := openWorkspace(ctx.Root, ctx.Now, ctx.IDs, false)
	if err != nil {
		t.Fatal(err)
	}
	return w
}

func newChange(t *testing.T, ctx *Ctx, req string) ChangeNewResponse {
	t.Helper()
	resp, err := callOp(t, ctx, "change.new", req)
	if err != nil {
		t.Fatalf("change.new %s: %v", req, err)
	}
	return resp.(ChangeNewResponse)
}

func currentOf(t *testing.T, ctx *Ctx) id.ID {
	t.Helper()
	cs, err := ctxWorkspace(t, ctx).current()
	if err != nil {
		t.Fatal(err)
	}
	return cs
}

func wantExit(t *testing.T, err error, code int) {
	t.Helper()
	if ExitCode(err) != code {
		t.Fatalf("err = %v, exit %d; want %d", err, ExitCode(err), code)
	}
}

// landSource records a landed source node with the given handle and
// returns its ID.
func landNode(t *testing.T, ctx *Ctx, kind op.Kind, handle string) id.ID {
	t.Helper()
	w := ctxWorkspace(t, ctx)
	node := w.mustID(t)
	c := op.Op{ID: w.mustID(t), Type: op.TypeNodeCreate, Node: node, Kind: kind, Handle: handle}
	if kind == op.KindSource {
		c.Fields = map[string]json.RawMessage{"url": json.RawMessage(`"https://example.com"`)}
	}
	if _, _, err := w.writeRecord("Add "+handle, []op.Op{c}, w.mustID(t), w.mustID(t)); err != nil {
		t.Fatal(err)
	}
	return node
}

func TestChangeOutsideRepo(t *testing.T) {
	for _, c := range []struct{ op, req string }{
		{"change.new", `{"intent":"x"}`},
		{"change.use", `{"none":true}`},
		{"change.show", `{}`},
		{"change.ls", `{}`},
	} {
		_, err := callOp(t, testCtx(t.TempDir()), c.op, c.req)
		if ExitCode(err) != 4 || AsError(err).Next != "rhiza init" {
			t.Errorf("%s outside a repo: %v (exit %d, next %q)", c.op, err, ExitCode(err), AsError(err).Next)
		}
	}
}

func TestChangeNew(t *testing.T) {
	t.Run("OpensAndMakesCurrent", func(t *testing.T) {
		ctx := initCtx(t)
		r := newChange(t, ctx, `{"intent":"Try it"}`)
		if r.ChangeSet.IsZero() || r.Entry.IsZero() || r.Actor.IsZero() {
			t.Fatalf("response = %+v", r)
		}
		if r.Kind != changeset.KindIntention || r.Intent != "Try it" || r.DryRun {
			t.Errorf("response = %+v", r)
		}
		if got := currentOf(t, ctx); got != r.ChangeSet {
			t.Errorf("current = %s, want %s", got, r.ChangeSet)
		}
		w := ctxWorkspace(t, ctx)
		sets := mustFoldClean(t, w)
		// The actor was recorded on first use, then the change set opened.
		if len(sets) != 2 {
			t.Fatalf("got %d change sets, want 2", len(sets))
		}
		var opened changeset.ChangeSet
		for _, s := range sets {
			if s.ID == r.ChangeSet {
				opened = s
			} else if s.Open.Kind != changeset.KindRecord || s.State != changeset.StateLanded ||
				s.Open.Intent != "Record actor test@example.com" {
				t.Errorf("actor change set = %+v", s)
			}
		}
		if opened.State != changeset.StateDraft || opened.OpenEntry != r.Entry || opened.Open.Intent != "Try it" {
			t.Errorf("opened = %+v", opened)
		}
		out := human(t, r)
		for _, s := range []string{r.ChangeSet.String(), "intention", "Try it", "current"} {
			if !strings.Contains(out, s) {
				t.Errorf("human output lacks %q:\n%s", s, out)
			}
		}

		// A second one reuses the actor and becomes current.
		r2 := newChange(t, ctx, `{"intent":"Other","kind":"record"}`)
		if r2.Actor != r.Actor || r2.Kind != changeset.KindRecord {
			t.Errorf("second = %+v", r2)
		}
		if got := currentOf(t, ctx); got != r2.ChangeSet {
			t.Errorf("current = %s, want %s", got, r2.ChangeSet)
		}
		if n := len(mustFoldClean(t, ctxWorkspace(t, ctx))); n != 3 {
			t.Errorf("got %d change sets, want 3", n)
		}
	})

	t.Run("BadInputWritesNothing", func(t *testing.T) {
		for _, tc := range []struct {
			req  string
			code int
		}{
			{`{"intent":"x","kind":"bogus"}`, 2},
			{`{"intent":""}`, 2},
			{`{"intent":"x","sources":["ZZZ"]}`, 4},
			{`{"intent":"x","title":"This title is far too long to be accepted"}`, 2},
			{`{"intent":"x","title":"two\nlines"}`, 2},
		} {
			req := tc.req
			ctx := initCtx(t)
			_, err := callOp(t, ctx, "change.new", req)
			if c := ExitCode(err); c != tc.code {
				t.Errorf("%s: %v, exit %d; want %d", req, err, c, tc.code)
			}
			if files := logFiles(t, ctx.Root); len(files) != 0 {
				t.Errorf("%s wrote %v", req, files)
			}
		}
	})

	t.Run("Title", func(t *testing.T) {
		ctx := initCtx(t)
		titled := newChange(t, ctx, `{"intent":"Make tabs readable: name them","title":"Readable tabs"}`)
		plain := newChange(t, ctx, `{"intent":"Carving keeps IDs, always"}`)
		if titled.Title != "Readable tabs" || plain.Title != "Carving keeps IDs" {
			t.Errorf("titles = %q, %q", titled.Title, plain.Title)
		}
		// Only the titled open needs format 2.
		versions := map[id.ID]int{}
		for _, e := range ctxWorkspace(t, ctx).log().Entries {
			versions[e.ID] = e.V
		}
		if versions[titled.Entry] != 2 || versions[plain.Entry] != 1 {
			t.Errorf("entry formats = %d, %d; want 2, 1", versions[titled.Entry], versions[plain.Entry])
		}
		for cs, want := range map[id.ID]string{titled.ChangeSet: "Readable tabs", plain.ChangeSet: "Carving keeps IDs"} {
			resp, err := callOp(t, ctx, "change.show", `{"ref":"`+cs.String()+`"}`)
			if err != nil {
				t.Fatal(err)
			}
			if got := resp.(ChangeShowResponse).Title; got != want {
				t.Errorf("show %s title = %q, want %q", cs, got, want)
			}
		}
		resp, err := callOp(t, ctx, "change.ls", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		ls := resp.(ChangeLsResponse)
		got := map[id.ID]string{}
		for _, c := range ls.ChangeSets {
			got[c.ID] = c.Title
		}
		if got[titled.ChangeSet] != "Readable tabs" || got[plain.ChangeSet] != "Carving keeps IDs" {
			t.Errorf("ls titles = %v", got)
		}
		// The human list names a set title beside the intent; the fallback,
		// a prefix of the intent, adds nothing.
		var b strings.Builder
		ls.Human(&b)
		if out := b.String(); !strings.Contains(out, ": Readable tabs (Make tabs readable: name them)") ||
			!strings.Contains(out, ": Carving keeps IDs, always (current)") {
			t.Errorf("ls human =\n%s", out)
		}
	})

	t.Run("RecordIsUsage", func(t *testing.T) {
		ctx := initCtx(t)
		ctx.Record = true
		_, err := callOp(t, ctx, "change.new", `{"intent":"x"}`)
		wantExit(t, err, 2)
	})

	t.Run("ChangeRefIgnored", func(t *testing.T) {
		ctx := initCtx(t)
		ctx.ChangeRef = "NOPE-NOT-A-SET"
		newChange(t, ctx, `{"intent":"x"}`)
	})

	t.Run("Sources", func(t *testing.T) {
		ctx := initCtx(t)
		a := landNode(t, ctx, op.KindSource, "paper-a")
		b := landNode(t, ctx, op.KindSource, "paper-b")
		landNode(t, ctx, op.KindQuestion, "why")
		r := newChange(t, ctx, `{"intent":"x","sources":["paper-b","paper-a","`+b.String()+`"]}`)
		want := []id.ID{a, b}
		slices.SortFunc(want, id.Compare)
		if !slices.Equal(r.Sources, want) {
			t.Errorf("sources = %v, want %v", r.Sources, want)
		}
		cs, _ := ctxWorkspace(t, ctx).changeSet(r.ChangeSet)
		if !slices.Equal(cs.Open.Sources, want) {
			t.Errorf("open sources = %v, want %v", cs.Open.Sources, want)
		}

		_, err := callOp(t, ctx, "change.new", `{"intent":"x","sources":["why"]}`)
		wantExit(t, err, 2)
		_, err = callOp(t, ctx, "change.new", `{"intent":"x","sources":["nothing-here"]}`)
		wantExit(t, err, 4)
	})

	t.Run("DryRunWritesNothing", func(t *testing.T) {
		ctx := initCtx(t)
		ctx.DryRun = true
		r := newChange(t, ctx, `{"intent":"Nope"}`)
		if !r.DryRun || r.ChangeSet.IsZero() || !r.Entry.IsZero() || !r.Actor.IsZero() {
			t.Errorf("response = %+v", r)
		}
		if files := logFiles(t, ctx.Root); len(files) != 0 {
			t.Errorf("dry run wrote %v", files)
		}
		if got := currentOf(t, ctx); !got.IsZero() {
			t.Errorf("dry run set current %s", got)
		}
		if out := human(t, r); !strings.Contains(out, "dry run") {
			t.Errorf("human output does not say dry run:\n%s", out)
		}
	})
}

func TestChangeUse(t *testing.T) {
	ctx := initCtx(t)
	first := newChange(t, ctx, `{"intent":"First"}`).ChangeSet
	second := newChange(t, ctx, `{"intent":"Second"}`).ChangeSet

	resp, err := callOp(t, ctx, "change.use", `{"ref":"`+strings.ToLower(first.String())+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(ChangeUseResponse); r.ChangeSet != first || r.Intent != "First" {
		t.Errorf("response = %+v", r)
	}
	if got := currentOf(t, ctx); got != first {
		t.Errorf("current = %s, want %s", got, first)
	}
	if out := human(t, resp); !strings.Contains(out, first.String()) {
		t.Errorf("human output:\n%s", out)
	}

	resp, err = callOp(t, ctx, "change.use", `{"none":true}`)
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(ChangeUseResponse); !r.ChangeSet.IsZero() {
		t.Errorf("cleared response = %+v", r)
	}
	if got := currentOf(t, ctx); !got.IsZero() {
		t.Errorf("current after --none = %s", got)
	}
	if out := human(t, resp); !strings.Contains(out, "No current change set") {
		t.Errorf("human output:\n%s", out)
	}

	t.Run("Usage", func(t *testing.T) {
		for _, req := range []string{`{}`, `{"ref":"` + second.String() + `","none":true}`, `{"ref":"ab"}`} {
			_, err := callOp(t, ctx, "change.use", req)
			if ExitCode(err) != 2 {
				t.Errorf("%s: %v, exit %d; want 2", req, err, ExitCode(err))
			}
		}
	})

	t.Run("Unknown", func(t *testing.T) {
		_, err := callOp(t, ctx, "change.use", `{"ref":"ZZZZZZ"}`)
		wantExit(t, err, 4)
	})

	t.Run("TerminalIsConflict", func(t *testing.T) {
		sets := mustFoldClean(t, ctxWorkspace(t, ctx))
		var landed id.ID
		for _, s := range sets {
			if s.State == changeset.StateLanded {
				landed = s.ID
			}
		}
		_, err := callOp(t, ctx, "change.use", `{"ref":"`+landed.String()+`"}`)
		wantExit(t, err, 3)
		if AsError(err).Next != newEntryNext || !strings.Contains(err.Error(), "landed") {
			t.Errorf("err = %v, next %q", err, AsError(err).Next)
		}
	})

	t.Run("DryRunKeepsCurrent", func(t *testing.T) {
		dry := *ctx
		dry.DryRun = true
		resp, err := callOp(t, &dry, "change.use", `{"ref":"`+second.String()+`"}`)
		if err != nil {
			t.Fatal(err)
		}
		if !resp.(ChangeUseResponse).DryRun {
			t.Errorf("response = %+v", resp)
		}
		if got := currentOf(t, ctx); !got.IsZero() {
			t.Errorf("dry run set current %s", got)
		}
	})
}

func TestChangeShow(t *testing.T) {
	ctx := initCtx(t)
	_, err := callOp(t, ctx, "change.show", `{}`)
	wantExit(t, err, 2)
	if AsError(err).Next != "rhiza change ls" {
		t.Errorf("next = %q", AsError(err).Next)
	}

	// An intention with an op, proposed and accepted.
	w := ctxWorkspace(t, ctx)
	actor := w.mustID(t)
	cs, _, err := w.openChangeSet(changeset.Open{Kind: changeset.KindIntention, Intent: "Decide"}, actor, actor)
	if err != nil {
		t.Fatal(err)
	}
	c := op.Op{ID: w.mustID(t), Type: op.TypeNodeCreate, Node: w.mustID(t), Kind: op.KindDecision, Handle: "use-go"}
	for _, b := range []changeset.Body{
		changeset.Ops{Ops: []op.Op{c}},
		changeset.Propose{},
		changeset.Accept{Rule: "owner", Reason: "Fine"},
	} {
		if _, err := w.appendEntry(cs, b, actor, actor); err != nil {
			t.Fatal(err)
		}
	}
	other := newChange(t, ctx, `{"intent":"Other"}`).ChangeSet

	show := func(ctx *Ctx, req string) ChangeShowResponse {
		t.Helper()
		resp, err := callOp(t, ctx, "change.show", req)
		if err != nil {
			t.Fatal(err)
		}
		return resp.(ChangeShowResponse)
	}

	// Default: the current change set.
	if r := show(ctx, `{}`); r.ID != other || !r.Current || r.State != changeset.StateDraft || len(r.Ops) != 0 {
		t.Errorf("current show = %+v", r)
	}

	// --change, then an explicit ref, which wins.
	refCtx := *ctx
	refCtx.ChangeRef = other.String()
	r := show(&refCtx, `{"ref":"`+cs.String()+`"}`)
	if r.ID != cs || r.Current || r.Kind != changeset.KindIntention || r.Intent != "Decide" ||
		r.State != changeset.StateAccepted {
		t.Errorf("show = %+v", r)
	}
	wantOp, _ := op.Encode(c)
	if len(r.Ops) != 1 || string(r.Ops[0]) != string(wantOp) {
		t.Errorf("ops = %s, want [%s]", r.Ops, wantOp)
	}
	if r.Decided == nil || r.Decided.Type != entry.TypeAccept || r.Decided.Rule != "owner" ||
		r.Decided.Reason != "Fine" || r.Decided.Actor != actor {
		t.Errorf("decided = %+v", r.Decided)
	}
	if len(r.Heads) != 1 || len(r.Applied) != 4 || r.Sources == nil || r.Problems == nil {
		t.Errorf("heads %v, applied %v, sources %v, problems %v", r.Heads, r.Applied, r.Sources, r.Problems)
	}
	out := human(t, r)
	for _, s := range []string{"# Change set " + cs.String(), "Decide", "accepted", "## Operations",
		"node.create", c.Node.String(), "use-go", "owner"} {
		if !strings.Contains(out, s) {
			t.Errorf("human output lacks %q:\n%s", s, out)
		}
	}
	if r := show(&refCtx, `{}`); r.ID != other {
		t.Errorf("--change show = %s, want %s", r.ID, other)
	}

	_, err = callOp(t, ctx, "change.show", `{"ref":"ZZZZZZ"}`)
	wantExit(t, err, 4)
}

// relanded lands an intention change set on its branch, reopens it, adds
// one more node and lands it again, leaving none current.
func relanded(t *testing.T, ctx *Ctx) id.ID {
	t.Helper()
	cs := landed(t, ctx)
	ctx.Git = branchGit("").run
	ctx.ChangeRef = cs.String()
	move(t, ctx, "reopen", `{}`)
	ctx.ChangeRef = ""
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"more"}`)
	move(t, ctx, "land", `{}`)
	return cs
}

func TestChangeShowHistory(t *testing.T) {
	ctx := initCtx(t)
	cs := relanded(t, ctx)
	resp, err := callOp(t, ctx, "change.show", `{"ref":"`+cs.String()+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(ChangeShowResponse)
	var types []entry.Type
	for _, e := range r.History {
		types = append(types, e.Type)
		if e.Entry.IsZero() || e.Actor.IsZero() || !lineageAtRe.MatchString(e.At) {
			t.Errorf("history event %+v", e)
		}
	}
	want := []entry.Type{entry.TypePropose, entry.TypeAccept, entry.TypeLand, entry.TypeReopen, entry.TypeLand}
	if !slices.Equal(types, want) {
		t.Errorf("history = %v, want %v", types, want)
	}
	out := human(t, r)
	block := "\n## History\n\n"
	for _, e := range r.History {
		block += "- " + string(e.Type) + " " + e.At + " by " + e.Actor.String() + "\n"
	}
	if !strings.Contains(out, block) {
		t.Errorf("human output lacks %q:\n%s", block, out)
	}

	// A fresh draft has an empty history, never null.
	other := newChange(t, ctx, `{"intent":"Other"}`).ChangeSet
	resp, err = callOp(t, ctx, "change.show", `{"ref":"`+other.String()+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	if h := resp.(ChangeShowResponse).History; h == nil || len(h) != 0 {
		t.Errorf("draft history = %#v", h)
	}
	if out := human(t, resp); !strings.Contains(out, "## History\n\nNo lifecycle moves yet.\n") {
		t.Errorf("draft human:\n%s", out)
	}
}

func TestChangeLs(t *testing.T) {
	ctx := initCtx(t)
	resp, err := callOp(t, ctx, "change.ls", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(ChangeLsResponse)
	if r.ChangeSets == nil || len(r.ChangeSets) != 0 || r.Problems == nil {
		t.Errorf("empty ls = %+v", r)
	}
	if out := human(t, r); !strings.Contains(out, "No change sets yet.") ||
		!strings.Contains(out, `next: rhiza change new "<intent>"`) {
		t.Errorf("empty human:\n%s", out)
	}

	a := newChange(t, ctx, `{"intent":"First"}`).ChangeSet
	b := newChange(t, ctx, `{"intent":"Second"}`).ChangeSet
	resp, err = callOp(t, ctx, "change.ls", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r = resp.(ChangeLsResponse)
	if len(r.ChangeSets) != 3 || len(r.Problems) != 0 {
		t.Fatalf("ls = %+v", r)
	}
	if !slices.IsSortedFunc(r.ChangeSets, func(x, y ChangeSetSummary) int { return id.Compare(x.ID, y.ID) }) {
		t.Errorf("not sorted: %+v", r.ChangeSets)
	}
	for _, s := range r.ChangeSets {
		switch s.ID {
		case a:
			if s.Current || s.State != changeset.StateDraft || s.Intent != "First" || s.Kind != changeset.KindIntention {
				t.Errorf("a = %+v", s)
			}
		case b:
			if !s.Current || s.Intent != "Second" {
				t.Errorf("b = %+v", s)
			}
		default:
			if s.State != changeset.StateLanded || s.Kind != changeset.KindRecord {
				t.Errorf("actor set = %+v", s)
			}
		}
	}
	out := human(t, r)
	for _, s := range []string{a.String(), b.String(), "First", "Second", "(current)", "Record actor test@example.com"} {
		if !strings.Contains(out, s) {
			t.Errorf("human lacks %q:\n%s", s, out)
		}
	}

	// Problems are listed, never hidden.
	dir := filepath.Join(ctx.Root, store.Dir, store.LogDir, a.String())
	if err := os.WriteFile(filepath.Join(dir, "garbage.json"), []byte(`{"x":1}`), 0o644); err != nil {
		t.Fatal(err)
	}
	resp, err = callOp(t, ctx, "change.ls", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r = resp.(ChangeLsResponse)
	if len(r.Problems) == 0 || !strings.HasPrefix(r.Problems[0], "log: ") {
		t.Errorf("problems = %q", r.Problems)
	}
	if out := human(t, r); !strings.Contains(out, "## Problems") {
		t.Errorf("human lacks problems:\n%s", out)
	}
}

func TestChangeShowLeaves(t *testing.T) {
	ctx := initCtx(t)
	landed(t, ctx) // store, a capability outside the change set
	cs := newChange(t, ctx, `{"intent":"Leaves"}`).ChangeSet
	for _, req := range []string{
		`{"kind":"capability","handle":"p"}`,
		`{"kind":"capability","handle":"c","parent":"p"}`,
		`{"kind":"capability","handle":"b","parent":"p"}`,
		`{"kind":"capability","handle":"a","parent":"p"}`,
		`{"kind":"capability","handle":"q"}`,
		`{"kind":"capability","handle":"e","parent":"q"}`,
		`{"kind":"capability","handle":"d","parent":"q","before":"e"}`,
		`{"kind":"capability","handle":"r"}`,
		`{"kind":"decision","handle":"not-a-capability"}`,
	} {
		mustWrite(t, ctx, "add", req)
	}
	for _, req := range []string{
		`{"from":"c","type":"depends-on","to":"b"}`,
		`{"from":"b","type":"depends-on","to":"a"}`,
		`{"from":"a","type":"depends-on","to":"store"}`,
		`{"from":"e","type":"depends-on","to":"store"}`,
		`{"from":"e","type":"depends-on","to":"d"}`,
	} {
		mustWrite(t, ctx, "link", req)
	}
	mustWrite(t, ctx, "retire", `{"ref":"r","reason":"Gone"}`)

	show := func(cs id.ID) ChangeShowResponse {
		t.Helper()
		resp, err := callOp(t, ctx, "change.show", `{"ref":"`+cs.String()+`"}`)
		if err != nil {
			t.Fatal(err)
		}
		return resp.(ChangeShowResponse)
	}
	r := show(cs)
	var got []string
	for _, l := range r.Leaves {
		got = append(got, l.Handle+"<"+strings.Join(l.Needs, ",")+">")
		if l.ID.IsZero() || l.Needs == nil {
			t.Errorf("leaf %+v", l)
		}
	}
	// a, b, c by depends-on; d before e by tree order (d was put before e)
	// and by e's need of d; store is outside the change set, so it is a need
	// but not ordered.
	want := []string{"a<store>", "b<a>", "c<b>", "d<>", "e<d,store>"}
	if !slices.Equal(got, want) {
		t.Errorf("leaves = %v, want %v", got, want)
	}
	if raw, _ := json.Marshal(r); !strings.Contains(string(raw), `"handle":"d","needs":[]`) {
		t.Errorf("json = %s", raw)
	}
	out := human(t, r)
	block := "\n## Leaves\n\n- a (needs store)\n- b (needs a)\n- c (needs b)\n- d\n- e (needs d, store)\n"
	if !strings.Contains(out, block) {
		t.Errorf("human output lacks %q:\n%s", block, out)
	}

	// Unlinked leaves come out in tree order, not ID order.
	ctx2 := initCtx(t)
	cs2 := newChange(t, ctx2, `{"intent":"Tree"}`).ChangeSet
	for _, req := range []string{
		`{"kind":"capability","handle":"z"}`,
		`{"kind":"capability","handle":"y","before":"z"}`,
		`{"kind":"capability","handle":"x","after":"y"}`,
	} {
		mustWrite(t, ctx2, "add", req)
	}
	resp, err := callOp(t, ctx2, "change.show", `{"ref":"`+cs2.String()+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	got = nil
	for _, l := range resp.(ChangeShowResponse).Leaves {
		got = append(got, l.Handle)
	}
	if want := []string{"y", "x", "z"}; !slices.Equal(got, want) {
		t.Errorf("tree-order leaves = %v, want %v", got, want)
	}

	// No leaves: an empty list, never null, and no Leaves section.
	empty := newChange(t, ctx2, `{"intent":"Empty"}`).ChangeSet
	resp, err = callOp(t, ctx2, "change.show", `{"ref":"`+empty.String()+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	if l := resp.(ChangeShowResponse).Leaves; l == nil || len(l) != 0 {
		t.Errorf("empty leaves = %#v", l)
	}
	if out := human(t, resp); strings.Contains(out, "## Leaves") {
		t.Errorf("empty human has a Leaves section:\n%s", out)
	}
}

func TestChangeShowLeavesLanded(t *testing.T) {
	ctx := initCtx(t)
	cs := landed(t, ctx)
	leaves := func() []string {
		t.Helper()
		resp, err := callOp(t, ctx, "change.show", `{"ref":"`+cs.String()+`"}`)
		if err != nil {
			t.Fatal(err)
		}
		var out []string
		for _, l := range resp.(ChangeShowResponse).Leaves {
			out = append(out, l.Handle)
		}
		return out
	}
	if got := leaves(); !slices.Equal(got, []string{"store"}) {
		t.Errorf("landed leaves = %v", got)
	}
	// A later change set retires store and lands: it is no longer a leaf.
	newChange(t, ctx, `{"intent":"Drop store"}`)
	mustWrite(t, ctx, "retire", `{"ref":"store","reason":"Gone"}`)
	for _, s := range []string{"propose", "accept", "land"} {
		move(t, ctx, s, `{}`)
	}
	if got := leaves(); len(got) != 0 {
		t.Errorf("leaves after retire = %v", got)
	}
}
