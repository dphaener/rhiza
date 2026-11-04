package api

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/gen"
)

// TestRoadmapOverAPageIsReadWhole writes a roadmap of more than 100 change
// sets, more than one default page, and checks that a plain read pages it
// while --all and health's roadmap check see every change set.
func TestRoadmapOverAPageIsReadWhole(t *testing.T) {
	root := t.TempDir()
	// A third of the in-flight change sets are drafts, off the roadmap:
	// 156 in flight leaves 104 proposed or accepted.
	if _, err := gen.Write(root, gen.Options{Seed: 1, Nodes: 300, Links: 100, Landed: 2, InFlight: 156}); err != nil {
		t.Fatal(err)
	}
	ctx := testCtx(root)
	ctx.IndexPath = filepath.Join(t.TempDir(), "index.db")
	count := func(r RoadmapResponse) int {
		c := 0
		for _, g := range r.Groups {
			c += len(g.ChangeSets)
		}
		return c
	}

	out, err := callOp(t, ctx, "roadmap", `{"all":true}`)
	if err != nil {
		t.Fatal(err)
	}
	all := out.(RoadmapResponse)
	n := count(all)
	if n <= 100 || all.NextCursor != "" {
		t.Fatalf("--all holds %d change sets, next cursor %q; want more than 100 and none", n, all.NextCursor)
	}

	out, err = callOp(t, ctx, "roadmap", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	if page := out.(RoadmapResponse); count(page) >= n || page.NextCursor == "" {
		t.Fatalf("default page holds %d of %d change sets, next cursor %q; want a short page and a cursor", count(page), n, page.NextCursor)
	}

	p := healthRoadmap(ctx)
	if want := fmt.Sprintf(" %d change sets,", n); !strings.Contains(p.Summary, want) {
		t.Fatalf("health roadmap = %q, want it to count %d change sets", p.Summary, n)
	}
}
