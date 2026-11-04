package roadmapapp

import (
	"encoding/json"
	"math/rand"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/gen"
	"github.com/dphaener/rhiza/internal/id"
)

// TestReadsARoadmapOverAPageWhole runs the real roadmap and ready
// operations over a generated roadmap of more than 100 change sets, more
// than one default page, and checks the app lists every one of them.
func TestReadsARoadmapOverAPageWhole(t *testing.T) {
	root := t.TempDir()
	// A third of the in-flight change sets are drafts, off the roadmap:
	// 156 in flight leaves 104 proposed or accepted.
	if _, err := gen.Write(root, gen.Options{Seed: 1, Nodes: 300, Links: 100, Landed: 2, InFlight: 156}); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 11, 3, 12, 0, 0, 0, time.UTC)
	clock := func() time.Time { return now }
	ctx := &api.Ctx{Root: root, IndexPath: filepath.Join(t.TempDir(), "index.db"), Now: clock,
		IDs: id.NewGenerator(clock, rand.New(rand.NewSource(1)))}
	run := func(opName, _ string, args []string) ([]byte, error) {
		o, _ := api.Lookup(opName)
		var req any = api.ReadyRequest{}
		if opName == "roadmap" {
			req = api.RoadmapRequest{All: slices.Contains(args, "--all")}
		}
		resp, err := o.Run(ctx, req)
		if err != nil {
			return nil, err
		}
		return json.Marshal(resp)
	}

	page, err := run("roadmap", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	all, err := run("roadmap", "", []string{"--all"})
	if err != nil {
		t.Fatal(err)
	}
	listed := func(out []byte) (n int, next string) {
		var rm api.RoadmapResponse
		if err := json.Unmarshal(out, &rm); err != nil {
			t.Fatal(err)
		}
		for _, g := range rm.Groups {
			if g.Kind != "umbrella" {
				n += len(g.ChangeSets)
			}
		}
		return n, rm.NextCursor
	}
	n, next := listed(all)
	if p, pnext := listed(page); n <= 100 || next != "" || p >= n || pnext == "" {
		t.Fatalf("--all lists %d (next %q), a default page %d (next %q); want more than 100 over more than one page", n, next, p, pnext)
	}

	m := NewModel(run)
	if m.err != nil {
		t.Fatal(m.err)
	}
	rows := 0
	for _, r := range m.rows {
		if !r.id.IsZero() {
			rows++
		}
	}
	if rows != n {
		t.Fatalf("the app lists %d change sets, want all %d", rows, n)
	}
}
