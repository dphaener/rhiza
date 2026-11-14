package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestMilestonesOperation reads a fixture with every milestone state and
// every change set state, and the roadmap of the same fixture, unchanged.
func TestMilestonesOperation(t *testing.T) {
	c := testCtx(filepath.Join(t.TempDir(), "p"))
	pinFreeBytes(c)
	if err := os.MkdirAll(c.Root, 0o755); err != nil {
		t.Fatal(err)
	}
	c.IndexPath = filepath.Join(t.TempDir(), "index.db")
	if _, err := callOp(t, c, "init", `{}`); err != nil {
		t.Fatal(err)
	}
	newChange(t, c, `{"intent":"Model"}`)
	for _, w := range []string{
		`{"op":"add","kind":"intent","handle":"a","fields":["milestone=1","statement=Milestone A"]}`,
		`{"op":"add","kind":"intent","handle":"a1","parent":"a","fields":["statement=Done point A1"]}`,
		`{"op":"add","kind":"intent","handle":"a2","parent":"a","fields":["statement=Done point A2"]}`,
		`{"op":"add","kind":"intent","handle":"b","fields":["milestone=2","statement=Milestone B"]}`,
		`{"op":"add","kind":"intent","handle":"c","fields":["milestone=3","statement=Milestone C"]}`,
		`{"op":"add","kind":"intent","handle":"d","fields":["milestone=4","statement=Milestone D"]}`,
		`{"op":"add","kind":"intent","handle":"d1","parent":"d","fields":["statement=Done point D1"]}`,
	} {
		if _, err := applyFile(t, c, `{"writes":[`+w+`]}`); err != nil {
			t.Fatal(err)
		}
	}
	for _, s := range []string{"propose", "accept", "land"} {
		move(t, c, s, `{}`)
	}
	// make creates a change set serving to, and moves it through steps.
	made := map[string]string{}
	make := func(name, to string, steps ...string) {
		t.Helper()
		cs := newChange(t, c, `{"intent":"`+name+`"}`).ChangeSet.String()
		made[name] = cs
		if _, err := callOp(t, c, "link", `{"from":"`+cs+`","type":"serves","to":"`+to+`"}`); err != nil {
			t.Fatal(err)
		}
		for _, s := range steps {
			move(t, c, s, `{}`)
		}
	}
	make("A draft", "a")
	make("A proposed", "a", "propose")
	make("A accepted", "a2", "propose", "accept")
	make("A landed", "a1", "propose", "accept", "land")
	make("A withdrawn", "a1", "propose", "withdraw")
	make("B proposed", "b", "propose")
	make("C draft", "c")
	make("D landed", "d1", "propose", "accept", "land")

	resp, err := callOp(t, c, "milestones", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	got := resp.(MilestonesResponse)
	if got.Finished != 1 || len(got.Milestones) != 4 {
		t.Fatalf("finished %d of %d milestones: %+v", got.Finished, len(got.Milestones), got)
	}
	states := map[string]string{}
	for _, m := range got.Milestones {
		states[m.Statement] = m.State
	}
	want := map[string]string{"Milestone A": "building", "Milestone B": "planned", "Milestone C": "designing", "Milestone D": "finished"}
	for k, v := range want {
		if states[k] != v {
			t.Errorf("%s is %q, want %q", k, states[k], v)
		}
	}
	a := got.Milestones[0]
	if a.Statement != "Milestone A" || len(a.DonePoints) != 2 || len(a.ChangeSets) != 5 {
		t.Fatalf("milestone A = %+v", a)
	}
	if !a.DonePoints[0].Met || a.DonePoints[1].Met {
		t.Errorf("A1 met %v, A2 met %v; want true, false", a.DonePoints[0].Met, a.DonePoints[1].Met)
	}
	if n := len(a.DonePoints[0].ChangeSets); n != 2 {
		t.Errorf("A1 has %d change sets, want landed and withdrawn", n)
	}
	byTitle := map[string]string{}
	for _, cs := range a.ChangeSets {
		byTitle[cs.Title] = string(cs.State)
	}
	for name, st := range map[string]string{"A draft": "draft", "A proposed": "proposed", "A accepted": "accepted", "A landed": "landed", "A withdrawn": "withdrawn"} {
		if byTitle[name] != st {
			t.Errorf("%s is %q, want %q (%v)", name, byTitle[name], st, byTitle)
		}
	}

	// roadmap shows proposed and accepted change sets only, as before.
	rr, err := callOp(t, c, "roadmap", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]string{}
	for _, g := range rr.(RoadmapResponse).Groups {
		for _, cs := range g.ChangeSets {
			seen[cs.Title] = string(cs.State)
		}
	}
	wantRoad := map[string]string{"A proposed": "proposed", "A accepted": "accepted", "B proposed": "proposed"}
	if len(seen) != len(wantRoad) {
		t.Errorf("roadmap = %v, want %v", seen, wantRoad)
	}
	for k, v := range wantRoad {
		if seen[k] != v {
			t.Errorf("roadmap %s = %q, want %q", k, seen[k], v)
		}
	}
	if _, err := json.Marshal(got); err != nil {
		t.Fatal(err)
	}
}
