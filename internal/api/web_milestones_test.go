package api

import (
	"path/filepath"
	"testing"
)

// The Overview names each open milestone by its statement, read from the
// repository, so two milestones read apart.
func TestWebMilestonesSourceNamesByStatement(t *testing.T) {
	ctx, _ := daemonRepo(t)
	ctx.IndexPath = filepath.Join(t.TempDir(), "index.db")
	ctx.Record = true
	if _, err := applyFile(t, ctx, `{"writes":[
  {"op":"add","kind":"intent","handle":"m1","fields":["milestone=1","statement=Ship the web overview"]},
  {"op":"add","kind":"intent","handle":"m2","fields":["milestone=2","statement=Run on the VM"]}
]}`); err != nil {
		t.Fatal(err)
	}
	ctx.Record = false
	for _, m := range []string{"m1", "m2"} {
		cs := newChange(t, ctx, `{"intent":"Serve `+m+`"}`).ChangeSet
		if _, err := applyFile(t, ctx, `{"writes":[{"op":"link","from":"`+cs.String()+`","type":"serves","to":"`+m+`"}]}`); err != nil {
			t.Fatal(err)
		}
		if _, err := callOp(t, ctx, "propose", `{}`); err != nil {
			t.Fatal(err)
		}
	}
	got := WebMilestonesSource(ctx)()
	if len(got) != 2 || got[0].Name != "Ship the web overview" || got[1].Name != "Run on the VM" {
		t.Fatalf("milestones = %+v", got)
	}
	if got[0].Total != 1 || got[1].Total != 1 {
		t.Errorf("milestones = %+v", got)
	}
}
