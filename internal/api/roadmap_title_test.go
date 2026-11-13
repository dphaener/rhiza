package api

import (
	"path/filepath"
	"testing"
)

// TestRoadmapChangeSetsCarryTheirDisplayTitle checks that roadmap gives each
// change set the title change ls shows: the set title, else the intent's
// first clause.
func TestRoadmapChangeSetsCarryTheirDisplayTitle(t *testing.T) {
	ctx := testCtx(t.TempDir())
	ctx.IndexPath = filepath.Join(t.TempDir(), "index.db")
	if _, err := callOp(t, ctx, "init", `{}`); err != nil {
		t.Fatal(err)
	}
	for _, req := range []string{
		`{"intent":"Stream the export, then compress it for the slow link","kind":"intention","title":"Stream exports"}`,
		`{"intent":"Cache the roadmap, then drop the old reader entirely","kind":"intention"}`,
	} {
		if _, err := callOp(t, ctx, "change.new", req); err != nil {
			t.Fatal(err)
		}
		for _, op := range []string{"propose", "accept"} {
			if _, err := callOp(t, ctx, op, `{}`); err != nil {
				t.Fatal(err)
			}
		}
	}
	out, err := callOp(t, ctx, "roadmap", `{"all":true}`)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, g := range out.(RoadmapResponse).Groups {
		for _, c := range g.ChangeSets {
			got[c.Intent] = c.Title
		}
	}
	want := map[string]string{
		"Stream the export, then compress it for the slow link": "Stream exports",
		"Cache the roadmap, then drop the old reader entirely":  "Cache the roadmap",
	}
	for intent, title := range want {
		if got[intent] != title {
			t.Errorf("title of %q = %q, want %q (all: %v)", intent, got[intent], title, got)
		}
	}
}
