package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/store"
)

// readCtx is initCtx with an index path in a temp folder.
func readCtx(t *testing.T) *Ctx {
	t.Helper()
	ctx := initCtx(t)
	ctx.IndexPath = filepath.Join(t.TempDir(), "index.db")
	return ctx
}

// mustReader opens a reader and closes it when the test ends.
func mustReader(t *testing.T, ctx *Ctx, view string, rule viewRule) *reader {
	t.Helper()
	r, err := openReader(ctx, "show", view, rule)
	if err != nil {
		t.Fatalf("openReader(%q): %v", view, err)
	}
	t.Cleanup(func() { r.Close() })
	return r
}

// readLog builds a small log: a landed change set (a capability tree, an
// intent, a decision, links) and an accepted one on top (a new child, a
// field change, a move, two retirements, one with a live child, a retired
// link and a new link).
// It returns the ctx and both change sets; the accepted one is current.
func readLog(t *testing.T) (*Ctx, id.ID, id.ID) {
	t.Helper()
	ctx := readCtx(t)
	landed := newChange(t, ctx, `{"intent":"Base"}`).ChangeSet
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store","fields":["statement=Keep things"]}`)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"cache","parent":"store"}`)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"disk","parent":"store","json_fields":["acceptance={\"gb\":2}"]}`)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"tape","parent":"disk"}`)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"old","parent":"store"}`)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"under-old","parent":"old"}`)
	mustWrite(t, ctx, "add", `{"kind":"intent","handle":"faster"}`)
	mustWrite(t, ctx, "add", `{"kind":"decision","handle":"use-ssd"}`)
	mustWrite(t, ctx, "link", `{"from":"cache","type":"serves","to":"faster"}`)
	mustWrite(t, ctx, "link", `{"from":"disk","type":"decided-by","to":"use-ssd"}`)
	gone := mustWrite(t, ctx, "link", `{"from":"tape","type":"serves","to":"faster"}`).Links[0]
	move(t, ctx, "propose", `{}`)
	move(t, ctx, "accept", `{"reason":"ok"}`)
	move(t, ctx, "land", `{}`)

	accepted := newChange(t, ctx, `{"intent":"More"}`).ChangeSet
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"ram","parent":"cache"}`)
	mustWrite(t, ctx, "set", `{"ref":"store","fields":["statement=Keep more"]}`)
	mustWrite(t, ctx, "mv", `{"ref":"tape","parent":"cache"}`)
	mustWrite(t, ctx, "retire", `{"ref":"disk","reason":"gone"}`)
	// Undoing a write of the same change set drops it, so these come here.
	mustWrite(t, ctx, "unlink", fmt.Sprintf(`{"ref":%q}`, gone.String()))
	mustWrite(t, ctx, "retire", `{"ref":"old","reason":"replaced"}`)
	mustWrite(t, ctx, "link", `{"from":"ram","type":"serves","to":"faster"}`)
	move(t, ctx, "propose", `{}`)
	move(t, ctx, "accept", `{"reason":"ok"}`)
	return ctx, landed, accepted
}

// sameRead fails unless both results and both errors are equal.
func sameRead(t *testing.T, what string, a, b any, ea, eb error) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%s: stored %+v, what-if %+v", what, a, b)
	}
	if !reflect.DeepEqual(ea, eb) {
		t.Errorf("%s: stored error %v, what-if error %v", what, ea, eb)
	}
}

func TestReadBackendsAgree(t *testing.T) {
	ctx, _, accepted := readLog(t)
	stored := mustReader(t, ctx, "to-be", defaultView)
	whatIf := mustReader(t, ctx, "what-if:"+accepted.String(), defaultView)
	if stored.ov != nil || whatIf.ov == nil {
		t.Fatal("wrong backends")
	}
	if got, want := whatIf.View(), "what-if:"+accepted.String(); got != want {
		t.Errorf("View = %q, want %q", got, want)
	}

	nodes, err := stored.List(index.ListQuery{Retired: true, Limit: 1000})
	if err != nil || len(nodes) < 10 {
		t.Fatalf("List = %d nodes, %v", len(nodes), err)
	}
	var ls []gmodel.Link
	for _, n := range nodes {
		l, _ := stored.Links(n.ID, index.From)
		ls = append(ls, l...)
	}
	if len(ls) < 4 {
		t.Fatalf("only %d links", len(ls))
	}

	ids := []id.ID{{}, accepted}
	for _, n := range nodes {
		ids = append(ids, n.ID)
	}
	for _, i := range ids {
		a, aok, ea := stored.Node(i)
		b, bok, eb := whatIf.Node(i)
		sameRead(t, "Node "+i.String(), []any{a, aok}, []any{b, bok}, ea, eb)
		ca, ea := stored.Children(i)
		cb, eb := whatIf.Children(i)
		sameRead(t, "Children "+i.String(), ca, cb, ea, eb)
		aa, ea := stored.Ancestors(i)
		ab, eb := whatIf.Ancestors(i)
		sameRead(t, "Ancestors "+i.String(), aa, ab, ea, eb)
		for _, d := range []index.Direction{index.From, index.To, index.Both, 7} {
			la, ea := stored.Links(i, d)
			lb, eb := whatIf.Links(i, d)
			sameRead(t, fmt.Sprintf("Links %s %d", i, d), la, lb, ea, eb)
		}
		for _, depth := range []int{-1, 0, 1, 2} {
			for _, ret := range []bool{false, true} {
				ta, ea := stored.Subtree(i, depth, ret)
				tb, eb := whatIf.Subtree(i, depth, ret)
				sameRead(t, fmt.Sprintf("Subtree %s %d %v", i, depth, ret), ta, tb, ea, eb)
				ta, ea = stored.SubtreeLight(i, depth, ret)
				tb, eb = whatIf.SubtreeLight(i, depth, ret)
				sameRead(t, fmt.Sprintf("SubtreeLight %s %d %v", i, depth, ret), ta, tb, ea, eb)
			}
		}
		for _, q := range []index.ListQuery{
			{Limit: 100},
			{Limit: 2},
			{Limit: 3, After: i},
			{Limit: 100, HasParent: true, Parent: i},
			{Limit: 100, HasParent: true, Parent: i, Retired: true},
			{Limit: 100, Kinds: []op.Kind{op.KindCapability, op.KindDecision}, Retired: true},
			{Limit: 0},
		} {
			na, ea := stored.List(q)
			nb, eb := whatIf.List(q)
			sameRead(t, fmt.Sprintf("List %+v", q), na, nb, ea, eb)
		}
	}

	refs := []string{"store", "old", "disk", "ram", "nothing", "", "abc", accepted.String()}
	for _, n := range nodes {
		refs = append(refs, n.ID.String(), n.ID.String()[:4], n.ID.String()[:20], strings.ToLower(n.ID.String()))
	}
	for _, ref := range refs {
		a, ea := stored.Resolve(ref)
		b, eb := whatIf.Resolve(ref)
		sameRead(t, "Resolve "+ref, a, b, ea, eb)
	}
	lrefs := []string{"", "abc", "ZZZZ", nodes[0].ID.String()}
	for _, l := range ls {
		lrefs = append(lrefs, l.ID.String(), l.ID.String()[:4], l.ID.String()[:22])
		a, aok, ea := stored.Link(l.ID)
		b, bok, eb := whatIf.Link(l.ID)
		sameRead(t, "Link "+l.ID.String(), []any{a, aok}, []any{b, bok}, ea, eb)
	}
	for _, ref := range lrefs {
		a, ea := stored.ResolveLink(ref)
		b, eb := whatIf.ResolveLink(ref)
		sameRead(t, "ResolveLink "+ref, a, b, ea, eb)
	}

	// The batch reads, over ID sets with repeats, unknown IDs, change sets
	// and every node.
	for k, set := range [][]id.ID{nil, {{}, accepted}, ids, ids[2:4], append(slices.Clone(ids[3:]), ids[3])} {
		na, ea := stored.NodesLight(set)
		nb, eb := whatIf.NodesLight(set)
		sameRead(t, fmt.Sprintf("NodesLight %d", k), na, nb, ea, eb)
		for _, d := range []index.Direction{index.From, index.To, index.Both, 7} {
			la, ea := stored.LinksAt(set, d)
			lb, eb := whatIf.LinksAt(set, d)
			sameRead(t, fmt.Sprintf("LinksAt %d %d", k, d), la, lb, ea, eb)
		}
		for _, typ := range []op.LinkType{op.LinkServes, op.LinkDecidedBy, op.LinkBlocks} {
			la, ea := stored.LinksTo(set, typ)
			lb, eb := whatIf.LinksTo(set, typ)
			sameRead(t, fmt.Sprintf("LinksTo %d %s", k, typ), la, lb, ea, eb)
		}
	}
	if ns, _ := whatIf.NodesLight(ids); len(ns) != len(nodes) || ns[0].History != nil {
		t.Errorf("NodesLight(every node) = %d nodes, want %d light ones", len(ns), len(nodes))
	}
	ea1, ea := stored.EntriesOf([]id.ID{accepted, {}})
	eb1, eb := whatIf.EntriesOf([]id.ID{accepted, {}})
	sameRead(t, "EntriesOf", ea1, eb1, ea, eb)

	sa, ea := stored.ChangeSets()
	sb, eb := whatIf.ChangeSets()
	sameRead(t, "ChangeSets", sa, sb, ea, eb)
	if len(sa) < 2 {
		t.Errorf("ChangeSets = %d, want 2 or more", len(sa))
	}
	xa, ea := stored.Entries(accepted)
	xb, eb := whatIf.Entries(accepted)
	sameRead(t, "Entries", xa, xb, ea, eb)
	if len(xa) == 0 {
		t.Error("no entries")
	}
}

func TestReadWhatIfDraft(t *testing.T) {
	// A draft's what-if shows its writes; as-is does not.
	ctx := readCtx(t)
	cs := newChange(t, ctx, `{"intent":"X"}`).ChangeSet
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
	r := mustReader(t, ctx, "", defaultView)
	if r.View() != "what-if:"+cs.String() {
		t.Fatalf("View = %q", r.View())
	}
	if n, err := r.Resolve("store"); err != nil || n.Handle != "store" {
		t.Errorf("what-if Resolve = %+v, %v", n, err)
	}
	_, err := mustReader(t, ctx, "as-is", defaultView).Resolve("store")
	wantExit(t, err, 4)
}

func TestReadErrors(t *testing.T) {
	ctx, landed, _ := readLog(t)
	for _, v := range []string{"to-be", "what-if:" + mustCurrent(t, ctx).String()} {
		r := mustReader(t, ctx, v, defaultView)
		_, err := r.Resolve("nothing")
		wantExit(t, err, 4)
		wantNext(t, err, "rhiza ls")
		if msg := AsError(err).Message; msg != `no node matches "nothing"` {
			t.Errorf("%s: message %q", v, msg)
		}
		nodes, _ := r.List(index.ListQuery{Limit: 100})
		_, err = r.Resolve(nodes[0].ID.String()[:4])
		wantExit(t, err, 2)
		if !strings.Contains(err.Error(), "matches") || strings.Contains(err.Error(), "index:") {
			t.Errorf("%s: ambiguous message %q", v, err)
		}
		_, err = r.ResolveLink("ZZZZ")
		wantExit(t, err, 4)
		_, err = r.Ancestors(landed)
		wantExit(t, err, 4)
		_, err = r.Subtree(landed, 0, false)
		wantExit(t, err, 4)
		_, err = r.List(index.ListQuery{})
		wantExit(t, err, 1)
	}
}

// mustCurrent is the current change set, which must be set.
func mustCurrent(t *testing.T, ctx *Ctx) id.ID {
	t.Helper()
	cs := currentOf(t, ctx)
	if cs.IsZero() {
		t.Fatal("no current change set")
	}
	return cs
}

func TestReadOpenErrors(t *testing.T) {
	ctx, landed, _ := readLog(t)
	cases := []struct {
		name, view string
		exit       int
		next       string
	}{
		{"sideways", "sideways", 2, "rhiza describe show"},
		{"empty what-if", "what-if:", 2, "rhiza describe show"},
		{"unknown change set", "what-if:ZZZZ", 4, changeLsNext},
		{"unknown full ID", "what-if:" + ctxWorkspace(t, ctx).mustID(t).String(), 4, changeLsNext},
		{"landed", "what-if:" + landed.String(), 3, changeLsNext},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := openReader(ctx, "show", c.view, defaultView)
			wantExit(t, err, c.exit)
			wantNext(t, err, c.next)
			t.Logf("%s: %v", c.view, err)
		})
	}
	t.Run("no index path", func(t *testing.T) {
		c := *ctx
		c.IndexPath = ""
		_, err := openReader(&c, "show", "", defaultView)
		wantExit(t, err, 1)
		if err.Error() != "no index path" {
			t.Errorf("err = %v", err)
		}
		t.Logf("empty IndexPath: %v", err)
	})
	t.Run("no repository", func(t *testing.T) {
		c := testCtx(t.TempDir())
		c.IndexPath = filepath.Join(t.TempDir(), "x.db")
		_, err := openReader(c, "show", "", defaultView)
		wantExit(t, err, 4)
	})
	t.Run("missing change set", func(t *testing.T) {
		cs := ctxWorkspace(t, ctx).mustID(t)
		err := whatIfViewError(cs, []gmodel.Problem{
			{ChangeSet: cs, Op: cs, Err: fmt.Errorf("%w: an op", gmodel.ErrMissing)},
			{ChangeSet: cs, Err: fmt.Errorf("%w: not in the log", gmodel.ErrMissing)},
		})
		wantExit(t, err, 4)
		wantNext(t, err, changeLsNext)
		if err := whatIfViewError(cs, []gmodel.Problem{{ChangeSet: cs, Op: cs, Err: gmodel.ErrConflict}}); err != nil {
			t.Errorf("op problem: %v", err)
		}
	})
}

func TestViewParse(t *testing.T) {
	ctx, landed, accepted := readLog(t)
	w := ctxWorkspace(t, ctx)
	for _, c := range []struct {
		in, want string
		exit     int
	}{
		{"", "", 0},
		{" as-is ", "as-is", 0},
		{"to-be", "to-be", 0},
		{"what-if:" + accepted.String(), "what-if:" + accepted.String(), 0},
		{"what-if:" + strings.ToLower(accepted.String()[:20]), "what-if:" + accepted.String(), 0},
		{"what-if:" + landed.String(), "what-if:" + landed.String(), 0},
		{"what-if:" + accepted.String()[:4], "", 2},
		{"what-if:ab", "", 2},
		{"what-if:", "", 2},
		{"As-Is", "", 2},
		{"what-if", "", 2},
		{"sideways", "", 2},
		{"what-if:ZZZZ", "", 4},
	} {
		v, err := parseView(w, "ls", c.in)
		if ExitCode(err) != c.exit || v.String() != c.want {
			t.Errorf("parseView(%q) = %q, %v; want %q, exit %d", c.in, v, err, c.want, c.exit)
		}
		if c.exit == 2 && c.in != "what-if:"+accepted.String()[:4] && c.in != "what-if:ab" {
			wantNext(t, err, "rhiza describe ls")
		}
	}
}

func TestViewDefault(t *testing.T) {
	asIs, toBe := "as-is", "to-be"
	t.Run("NoChangeSet", func(t *testing.T) {
		ctx := readCtx(t)
		if got := mustReader(t, ctx, "", defaultView).View(); got != asIs {
			t.Errorf("View = %q", got)
		}
	})
	t.Run("States", func(t *testing.T) {
		ctx := readCtx(t)
		cs := newChange(t, ctx, `{"intent":"X"}`).ChangeSet
		whatIf := "what-if:" + cs.String()
		if got := mustReader(t, ctx, "", defaultView).View(); got != whatIf {
			t.Errorf("draft: View = %q", got)
		}
		move(t, ctx, "propose", `{}`)
		if got := mustReader(t, ctx, "", defaultView).View(); got != whatIf {
			t.Errorf("proposed: View = %q", got)
		}
		move(t, ctx, "accept", `{"reason":"ok"}`)
		if got := mustReader(t, ctx, "", defaultView).View(); got != toBe {
			t.Errorf("accepted: View = %q", got)
		}
		move(t, ctx, "land", `{}`)
		if got := mustReader(t, ctx, "", defaultView).View(); got != asIs {
			t.Errorf("landed, none current: View = %q", got)
		}
		// A current file still naming the landed change set.
		if err := ctxWorkspace(t, ctx).setCurrent(cs); err != nil {
			t.Fatal(err)
		}
		if got := mustReader(t, ctx, "", defaultView).View(); got != asIs {
			t.Errorf("landed current: View = %q", got)
		}
		// One not in the log is no change set.
		if err := ctxWorkspace(t, ctx).setCurrent(ctxWorkspace(t, ctx).mustID(t)); err != nil {
			t.Fatal(err)
		}
		if got := mustReader(t, ctx, "", defaultView).View(); got != asIs {
			t.Errorf("unknown current: View = %q", got)
		}
	})
	t.Run("Withdrawn", func(t *testing.T) {
		ctx := readCtx(t)
		cs := newChange(t, ctx, `{"intent":"X"}`).ChangeSet
		move(t, ctx, "withdraw", `{}`)
		ctx.ChangeRef = cs.String()
		if got := mustReader(t, ctx, "", defaultView).View(); got != asIs {
			t.Errorf("withdrawn: View = %q", got)
		}
	})
	t.Run("ChangeOverridesCurrent", func(t *testing.T) {
		ctx, landed, accepted := readLog(t)
		draft := newChange(t, ctx, `{"intent":"Y"}`).ChangeSet
		if got := mustReader(t, ctx, "", defaultView).View(); got != "what-if:"+draft.String() {
			t.Errorf("draft current: View = %q", got)
		}
		ctx.ChangeRef = accepted.String()[:20]
		if got := mustReader(t, ctx, "", defaultView).View(); got != toBe {
			t.Errorf("--change accepted: View = %q", got)
		}
		ctx.ChangeRef = ""
		ctx.ChangeSet = landed
		if got := mustReader(t, ctx, "", defaultView).View(); got != asIs {
			t.Errorf("ChangeSet landed: View = %q", got)
		}
		ctx.ChangeSet = id.ID{}
		ctx.ChangeRef = "ZZZZ"
		_, err := openReader(ctx, "show", "", defaultView)
		wantExit(t, err, 4)
		// An explicit view does not look at the effective change set.
		if got := mustReader(t, ctx, "as-is", defaultView).View(); got != asIs {
			t.Errorf("explicit as-is: View = %q", got)
		}
	})
	t.Run("OperationRule", func(t *testing.T) {
		// D17: roadmap and ready pass toBeView, to-be whatever is in effect.
		ctx := readCtx(t)
		if got := mustReader(t, ctx, "", toBeView).View(); got != toBe {
			t.Errorf("none: View = %q", got)
		}
		cs := newChange(t, ctx, `{"intent":"X"}`).ChangeSet
		if got := mustReader(t, ctx, "", toBeView).View(); got != toBe {
			t.Errorf("draft: View = %q", got)
		}
		if got := mustReader(t, ctx, "what-if:"+cs.String(), toBeView).View(); got != "what-if:"+cs.String() {
			t.Errorf("explicit: View = %q", got)
		}
	})
}

func TestReadHelpers(t *testing.T) {
	for _, lt := range op.LinkTypes() {
		if linkBackNames[lt] == "" {
			t.Errorf("no backwards name for %s", lt)
		}
	}
	if len(linkBackNames) != len(op.LinkTypes()) {
		t.Errorf("%d backwards names for %d link types", len(linkBackNames), len(op.LinkTypes()))
	}
	for lt, want := range map[op.LinkType]string{
		op.LinkServes: "served by", op.LinkDecidedBy: "decides", op.LinkDerivedFrom: "source of", "odd": "odd (from)",
	} {
		if got := backName(lt); got != want {
			t.Errorf("backName(%s) = %q, want %q", lt, got, want)
		}
	}
	i, err := id.Parse("01JB0000000000000000000000")
	if err != nil {
		t.Fatal(err)
	}
	if got := displayName(gmodel.Node{ID: i, Handle: "store"}); got != "store (01JB0000000000)" {
		t.Errorf("displayName = %q", got)
	}
	for raw, want := range map[string]string{
		`"two\nlines"`:       "two\nlines",
		`{"b":1, "a":[1,2]}`: "{\n  \"a\": [\n    1,\n    2\n  ],\n  \"b\": 1\n}",
		`7`:                  "7",
		`not json`:           "not json",
	} {
		if got := fieldText(json.RawMessage(raw)); got != want {
			t.Errorf("fieldText(%s) = %q, want %q", raw, got, want)
		}
	}
	if readError(nil) != nil {
		t.Error("readError(nil) is not nil")
	}
}

// countingCtx counts the log folds a read through ctx costs.
func countingCtx(ctx *Ctx) *int {
	n := new(int)
	ctx.readLog = func() store.Log {
		*n++
		return store.ReadLog(os.DirFS(filepath.Join(ctx.Root, store.Dir)))
	}
	return n
}

func TestDefaultViewFoldsNothing(t *testing.T) {
	ctx, landed, accepted := readLog(t)
	draft := newChange(t, ctx, `{"intent":"Later"}`).ChangeSet
	folds := countingCtx(ctx)
	cases := []struct {
		name string
		cs   id.ID
		want string
	}{
		{"landed", landed, "as-is"},
		{"accepted", accepted, "to-be"},
		{"draft", draft, "what-if:" + draft.String()},
	}
	for _, c := range cases {
		for _, byRef := range []bool{false, true} {
			rc := *ctx
			rc.readLog = ctx.readLog
			if byRef {
				rc.ChangeRef = c.cs.String()
			} else {
				rc.ChangeSet = c.cs
			}
			*folds = 0
			r, err := openReader(&rc, "show", "", defaultView)
			if err != nil {
				t.Fatalf("%s: %v", c.name, err)
			}
			got := r.View()
			r.Close()
			// What-if computes its model from the index; it folds no log.
			if got != c.want || *folds != 0 {
				t.Errorf("%s (byRef %v): view %q, %d folds; want %q, 0", c.name, byRef, got, *folds, c.want)
			}
		}
	}
	// The current change set, the landed one, costs no fold either.
	if err := ctxWorkspace(t, ctx).setCurrent(landed); err != nil {
		t.Fatal(err)
	}
	*folds = 0
	r := mustReader(t, ctx, "", defaultView)
	if r.View() != "as-is" || *folds != 0 {
		t.Errorf("current landed: view %q, %d folds", r.View(), *folds)
	}
	// Errors stay as they were.
	bad := *ctx
	bad.ChangeSet = accepted
	bad.ChangeRef = "abcd"
	if _, err := openReader(&bad, "show", "", defaultView); err == nil {
		t.Error("ChangeSet and ChangeRef together: no error")
	}
	bad = *ctx
	bad.ChangeRef = "ZZZZZZZZ"
	if _, err := openReader(&bad, "show", "", defaultView); err == nil {
		t.Error("unknown ChangeRef: no error")
	}
}

func TestDefaultViewSeesPending(t *testing.T) {
	ctx, _, _ := readLog(t)
	folds := countingCtx(ctx)
	w := ctxWorkspace(t, ctx)
	w.dryRun = true
	w.readLog = ctx.readLog
	pend, _, err := w.openChangeSet(changeset.Open{Kind: changeset.KindIntention, Intent: "Pending"}, id.ID{}, id.ID{})
	if err != nil {
		t.Fatal(err)
	}
	rc := *ctx
	rc.ChangeSet = pend
	*folds = 0
	var gotCS id.ID
	var gotState changeset.State
	r, err := openReaderIn(&rc, w, "show", "", func(cs id.ID, st changeset.State) view {
		gotCS, gotState = cs, st
		return view{stored: index.AsIs}
	})
	if err != nil {
		t.Fatalf("a pending change set: %v", err)
	}
	defer r.Close()
	if gotCS != pend || gotState != changeset.StateDraft {
		t.Errorf("rule saw %v in %q, want %v in draft", gotCS, gotState, pend)
	}
	if *folds == 0 {
		t.Error("a workspace with pending writes folded nothing")
	}
}
