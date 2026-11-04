package modelapp

import (
	"bytes"
	"math/rand"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/gen"
	"github.com/dphaener/rhiza/internal/id"
)

// TestRoadmapTabReadsARoadmapOverAPageWhole runs the real roadmap
// operation over a generated roadmap of more than 100 change sets, more
// than one default page, and checks the roadmap tab holds every one of them
// and no next: line.
func TestRoadmapTabReadsARoadmapOverAPageWhole(t *testing.T) {
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
	read := func(opName string, args []string) (string, error) {
		if opName != "roadmap" {
			return "- tree\n", nil
		}
		o, _ := api.Lookup(opName)
		resp, err := o.Run(ctx, api.RoadmapRequest{All: slices.Contains(args, "--all")})
		if err != nil {
			return "", err
		}
		var b bytes.Buffer
		resp.(api.Humaner).Human(&b)
		return b.String(), nil
	}
	entries := func(s string) int {
		n := 0
		for _, l := range strings.Split(s, "\n") {
			if strings.HasPrefix(l, "- [") {
				n++
			}
		}
		return n
	}

	all, err := read("roadmap", []string{"--all"})
	if err != nil {
		t.Fatal(err)
	}
	page, err := read("roadmap", nil)
	if err != nil {
		t.Fatal(err)
	}
	n := entries(all)
	if n <= 100 || entries(page) >= n || !strings.Contains(page, "\nnext: ") {
		t.Fatalf("--all lists %d, a default page %d; want more than 100 over more than one page", n, entries(page))
	}

	m := NewModel(read, "")
	if m.Err() != nil {
		t.Fatal(m.Err())
	}
	got := m.Text("roadmap")
	if entries(got) != n || strings.Contains(got, "\nnext: ") {
		t.Fatalf("the roadmap tab lists %d change sets, want all %d and no next: line", entries(got), n)
	}
}
