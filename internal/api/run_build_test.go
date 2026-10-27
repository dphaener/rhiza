package api

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/check"
	"github.com/dphaener/rhiza/internal/id"
)

// buildCheck finds a leaf OK when live Evidence verifying it was created
// by an actor named in no Actor: trailer of a commit whose Leaf: names the
// leaf; otherwise it says why: no Evidence, or the Evidence's actor wrote
// the leaf's commit. A commit naming another leaf, or the actor by ID,
// counts as described.
func TestBuildCheck(t *testing.T) {
	ctx := testCtx(t.TempDir())
	if _, err := callOp(t, ctx, "init", `{}`); err != nil {
		t.Fatal(err)
	}
	cs := newChange(t, ctx, `{"intent":"Build"}`).ChangeSet
	ada := recordAda(t, ctx, "ada@example.com")
	bob := recordAda(t, ctx, "bob@example.com")
	for _, h := range []string{"alpha", "beta", "gamma", "delta"} {
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"`+h+`"}`)
	}
	as := func(actor id.ID, ev, leaf string) {
		t.Helper()
		v := *ctx
		v.Actor = actor
		mustWrite(t, &v, "add", `{"kind":"evidence","handle":"`+ev+`","fields":["statement=works","result=PASS"]}`)
		mustWrite(t, &v, "link", `{"from":"`+ev+`","type":"verifies","to":"`+leaf+`"}`)
	}
	// alpha: ada verifies, bob built. beta: nothing. gamma: bob verifies
	// and built it. delta: bob verifies, and his commit naming delta
	// carries his ID rather than his handle.
	as(ada, "alpha-ok", "alpha")
	as(bob, "gamma-ok", "gamma")
	as(bob, "delta-ok", "delta")

	w := ctxWorkspace(t, ctx)
	m := w.model(cs)
	c, _ := w.changeSet(cs)
	adaH, bobH := nodeName(m, ada), nodeName(m, bob)
	gammaCommit := strings.Repeat("c", 40)
	commits := []check.Commit{
		{Hash: strings.Repeat("d", 40), Actors: []string{bob.String()}, Leaves: []string{"delta"}},
		{Hash: gammaCommit, Actors: []string{bobH}, Leaves: []string{"gamma"}},
		{Hash: strings.Repeat("a", 40), Actors: []string{bobH}, Leaves: []string{"alpha", "beta"}},
		// ada wrote a commit for another leaf: it does not count against alpha.
		{Hash: strings.Repeat("e", 40), Actors: []string{adaH}, Leaves: []string{"epsilon"}},
	}
	got := buildCheck(m, w.changeSets(), changeLeaves(m, c.Ops), commits)
	byHandle := map[string]RunBuildLeaf{}
	for _, l := range got {
		byHandle[l.Handle] = l
	}
	want := map[string]RunBuildLeaf{
		"alpha": {Handle: "alpha", OK: true, Evidence: []string{"alpha-ok"}},
		"beta":  {Handle: "beta", Evidence: []string{}, Why: "no Evidence verifies it"},
		"gamma": {Handle: "gamma", Evidence: []string{"gamma-ok"}, Why: "its Evidence gamma-ok was written by " + bobH + ", who wrote its commit cccccccccccc"},
		"delta": {Handle: "delta", Evidence: []string{"delta-ok"}, Why: "its Evidence delta-ok was written by " + bobH + ", who wrote its commit dddddddddddd"},
	}
	if len(got) != len(want) {
		t.Fatalf("buildCheck = %+v, want 4 leaves", got)
	}
	for h, wl := range want {
		gl := byHandle[h]
		if gl.Leaf == "" {
			t.Errorf("%s: no leaf ID", h)
		}
		gl.Leaf = ""
		if !reflect.DeepEqual(gl, wl) {
			t.Errorf("%s = %+v, want %+v", h, gl, wl)
		}
	}

	// With no commits at all, every leaf with Evidence is OK.
	for _, l := range buildCheck(m, w.changeSets(), changeLeaves(m, c.Ops), nil) {
		if l.OK != (l.Handle != "beta") {
			t.Errorf("no commits: %s OK = %v", l.Handle, l.OK)
		}
	}
}
