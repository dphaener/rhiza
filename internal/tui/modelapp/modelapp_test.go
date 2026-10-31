package modelapp

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/gen"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/quiettest"
)

func key(s string) tea.KeyPressMsg {
	if s == "tab" {
		return tea.KeyPressMsg{Code: tea.KeyTab}
	}
	r := []rune(s)[0]
	return tea.KeyPressMsg{Code: r, Text: s}
}

func TestTabSwitchesAndScrolls(t *testing.T) {
	read := func(op string, args []string) (string, error) {
		var b strings.Builder
		for i := range 100 {
			fmt.Fprintf(&b, "%s line %d\n", op, i)
		}
		return b.String(), nil
	}
	m := NewModel(read, "")
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 13})
	if v := m.View().Content; !strings.Contains(v, "[tree]") || !strings.Contains(v, "tree line 0") {
		t.Fatalf("want the tree first, got %q", v)
	}
	m.Update(key("j"))
	if v := m.View().Content; strings.Contains(v, "tree line 0\n") || !strings.Contains(v, "tree line 1\n") {
		t.Fatalf("j did not scroll: %q", v)
	}
	m.Update(key("tab"))
	if v := m.View().Content; !strings.Contains(v, "[roadmap]") || !strings.Contains(v, "roadmap line 0") {
		t.Fatalf("tab did not show the roadmap: %q", v)
	}
	if _, cmd := m.Update(key("q")); cmd == nil {
		t.Fatal("q did not quit")
	}
}

func TestViewFitsTheWindow(t *testing.T) {
	read := func(op string, args []string) (string, error) {
		return strings.Repeat(op+" line\n", 100), nil
	}
	m := NewModel(read, "")
	for _, h := range []int{5, 13, 24} {
		m.Update(tea.WindowSizeMsg{Width: 80, Height: h})
		v := m.View().Content
		if n := len(strings.Split(v, "\n")); n > h {
			t.Fatalf("height %d: view has %d lines", h, n)
		}
		if !strings.HasPrefix(v, "[tree]") {
			t.Fatalf("height %d: header is not the first line: %q", h, v)
		}
	}
}

func TestReadErrorIsShown(t *testing.T) {
	read := func(op string, args []string) (string, error) {
		if op == "roadmap" {
			return "", fmt.Errorf("no index")
		}
		return "- a (1)\n", nil
	}
	m := NewModel(read, "")
	if m.Err() == nil || m.Text("roadmap") != "" || m.Text("tree") != "- a (1)\n" {
		t.Fatalf("Err %v, roadmap %q, tree %q", m.Err(), m.Text("roadmap"), m.Text("tree"))
	}
	m.Update(key("tab"))
	if !strings.Contains(m.View().Content, "no index") {
		t.Fatalf("error not shown: %q", m.View().Content)
	}
}

// budgetLock orders this test with the other packages' budget tests, which
// each build the 50,000-node model: a read timed beside another build
// measures its contention, not the read.
func budgetLock(t *testing.T) {
	t.Helper()
	f, err := os.OpenFile(filepath.Join(os.TempDir(), "rhiza-budget.lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		f.Close()
	})
}

// TestGeneratedModel builds the app over the 50,000-node model, checks its
// two panes against the text of the tree and roadmap operations, and
// times the read, as the CLI does it: a fresh Ctx per operation. It first
// waits until go test runs nothing else, as TestReadBudget does.
func TestGeneratedModel(t *testing.T) {
	if testing.Short() {
		t.Skip("generates the 50,000-node model; skipped under -short")
	}
	q := quiettest.Wait(t)
	t.Logf("other packages: %s", q)
	budgetLock(t)
	dir := t.TempDir()
	root := filepath.Join(dir, "repo")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	sum, err := gen.Write(root, gen.Default())
	if err != nil {
		t.Fatal(err)
	}
	indexPath := filepath.Join(dir, "index.db")
	x, err := index.Open(root, indexPath)
	if err != nil {
		t.Fatal(err)
	}
	// The subtree is the one the api read budget times: the deep sample's
	// grandparent, two levels above the deepest.
	as, err := x.Ancestors(index.AsIs, sum.Deep.ID)
	if err != nil || len(as) < 2 {
		t.Fatalf("deep node's ancestors: %v, %v", as, err)
	}
	subtree := as[len(as)-2].ID.String()
	if err := x.Close(); err != nil {
		t.Fatal(err)
	}

	read := func(opName string, args []string) (string, error) {
		o, _ := api.Lookup(opName)
		var req any
		switch opName {
		case "tree":
			req = api.TreeRequest{Ref: subtree, View: "as-is"}
		default:
			req = api.RoadmapRequest{}
		}
		resp, err := o.Run(&api.Ctx{Root: root, IndexPath: indexPath, Now: time.Now, IDs: id.Default()}, req)
		if err != nil {
			return "", err
		}
		var b strings.Builder
		resp.(api.Humaner).Human(&b)
		return b.String(), nil
	}

	const runs = 20
	var times []time.Duration
	var m *Model
	for range runs {
		start := time.Now()
		m = NewModel(read, subtree)
		times = append(times, time.Since(start))
		if err := m.Err(); err != nil {
			t.Fatal(err)
		}
	}
	slices.Sort(times)
	median := (times[runs/2-1] + times[runs/2]) / 2
	t.Logf("tree and roadmap read: median %v (min %v, max %v)", median, times[0], times[runs-1])
	if median >= 50*time.Millisecond {
		t.Errorf("read median %v, want < 50ms", median)
	}
	if late := quiettest.Late("tree and roadmap", "as-is"); late != "" {
		t.Log(late)
	}

	for _, c := range []struct{ pane, op string }{{"tree", "tree"}, {"roadmap", "roadmap"}} {
		want, err := read(c.op, []string{"x"})
		if err != nil {
			t.Fatal(err)
		}
		if got := m.Text(c.pane); got != want || len(got) < 1000 {
			t.Errorf("%s pane differs from the operation's text (%d vs %d bytes)", c.pane, len(got), len(want))
		}
	}
}
