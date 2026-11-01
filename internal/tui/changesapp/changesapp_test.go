package changesapp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/api"
)

// fakeFeed serves items (oldest first) two to a page, as the feed does.
func fakeFeed(items []map[string]any) Reader {
	return func(args []string) ([]byte, error) {
		start := 0
		for i, a := range args {
			if a == "--cursor" {
				start = int(args[i+1][0] - '0')
			}
		}
		end := min(start+2, len(items))
		resp := map[string]any{"change_sets": items[start:end]}
		if end < len(items) {
			resp["next_cursor"] = string(rune('0' + end))
		}
		return json.Marshal(resp)
	}
}

func item(ulid, title string, pr int) map[string]any {
	m := map[string]any{"id": ulid, "title": title, "kind": "record", "leaves": 0,
		"landed_at": "2026-10-25T12:00:00Z", "actor": "a", "land_entry": ulid}
	if pr > 0 {
		m["pr"] = pr
	}
	return m
}

const (
	idA = "01M5PWGEF85Y185JZ2ZGCNH40N"
	idB = "01M5PWGEF85Y185JZ2ZGCNH40P"
	idC = "01M5PWGEF85Y185JZ2ZGCNH40Q"
)

func titles(m *Model) []string {
	var out []string
	for _, it := range m.Unseen() {
		out = append(out, it.Title)
	}
	return out
}

func TestUnseenNewestFirstAcrossPagesThenSeen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rhiza", SeenFile)
	items := []map[string]any{item(idA, "a", 1), item(idB, "b", 0)}
	m := NewModel(fakeFeed(items), path)
	if got := strings.Join(titles(m), ","); got != "b,a" {
		t.Fatalf("unseen = %s, want b,a", got)
	}
	if err := m.MarkSeen(); err != nil {
		t.Fatal(err)
	}
	m = NewModel(fakeFeed(items), path)
	if len(m.Unseen()) != 0 {
		t.Fatalf("seen change sets listed again: %v", titles(m))
	}
	// Three pages' worth; a change set landed later (even with an older ID,
	// as a merged branch's would be) shows up alone.
	items = append(items, item(idC, "c", 3))
	m = NewModel(fakeFeed(items), path)
	if got := strings.Join(titles(m), ","); got != "c" {
		t.Fatalf("unseen = %s, want c", got)
	}
	out, err := m.Result()
	if err != nil || !strings.Contains(string(out), `"pr": 3`) {
		t.Fatalf("result %s, %v", out, err)
	}
}

func TestSeenFileHoldsIDsSorted(t *testing.T) {
	path := filepath.Join(t.TempDir(), SeenFile)
	m := NewModel(fakeFeed([]map[string]any{item(idB, "b", 0), item(idA, "a", 0)}), path)
	if err := m.MarkSeen(); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(path)
	var f seenFile
	if err := json.Unmarshal(b, &f); err != nil || len(f.Seen) != 2 || f.Seen[0] != idA || f.Seen[1] != idB {
		t.Fatalf("seen file %s (%v)", b, err)
	}
}

func TestEmptyResultIsAnEmptyList(t *testing.T) {
	m := NewModel(fakeFeed(nil), filepath.Join(t.TempDir(), SeenFile))
	out, err := m.Result()
	if err != nil || !strings.Contains(string(out), `"change_sets": []`) {
		t.Fatalf("result %s, %v", out, err)
	}
	if !strings.Contains(m.View().Content, "Nothing new") {
		t.Fatalf("empty view: %q", m.View().Content)
	}
	var _ = api.FeedItem{}
}
