package launcher

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"

	"github.com/dphaener/rexuikit/client"
	"github.com/dphaener/rexuikit/surface"
)

// fakeRex is the part of Rex a modal open uses: one window of 200x60 cells,
// recording each new layer and call. Methods it does not define panic on the
// nil embedded interface, so a stray Rex call fails the test.
type fakeRex struct {
	surface.Rex
	mu     sync.Mutex
	layers []client.NewLayerArgs
	calls  []string
}

func (f *fakeRex) View(context.Context, string) (client.View, error) {
	return client.View{
		SessionID:      "session:s",
		ActiveWindowID: "window:w",
		Windows: []client.Window{{
			WindowID: "window:w",
			Layers: []client.Layer{{
				LayerID: "layer:tiled", Kind: client.LayerTiled,
				Layout: &client.LayoutNode{Rect: client.Bounds{0, 0, 1, 1}, BlockID: "block:shell"},
			}},
		}},
	}, nil
}

func (f *fakeRex) TerminalSize(context.Context, string, string) (client.TerminalSize, error) {
	return client.TerminalSize{Columns: 200, Rows: 60, CellWidthPx: 7, CellHeightPx: 14}, nil
}

func (f *fakeRex) NewLayer(_ context.Context, a client.NewLayerArgs) (client.NewLayerResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.layers = append(f.layers, a)
	return client.NewLayerResult{LayerID: "layer:1", BlockIDs: []string{"block:1"}}, nil
}

func (f *fakeRex) Call(_ context.Context, method string, _, _ any) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, method)
	return nil
}

func (f *fakeRex) Events() <-chan client.Event { return nil }

func testApps() []App {
	return []App{
		{Name: "ready", Summary: "Leaves ready to build"},
		{Name: "roadmap", Summary: "Milestones and change sets"},
		{Name: "show", Summary: "One node"},
	}
}

func newTestModel(f *fakeRex, apps []App) *Model {
	l := Launcher{Rex: f, Session: "session:s", Bin: "/opt/rhiza/bin/rhiza", Cwd: "/work"}
	return NewModel(context.Background(), l, apps, colorprofile.TrueColor)
}

func press(m *Model, code rune) (*Model, tea.Cmd) {
	next, cmd := m.Update(tea.KeyPressMsg{Code: code})
	return next.(*Model), cmd
}

func TestAppsAreSortedAndReplaced(t *testing.T) {
	mu.Lock()
	saved := registry
	registry = map[string]App{}
	mu.Unlock()
	defer func() { mu.Lock(); registry = saved; mu.Unlock() }()

	Register(App{Name: "show", Summary: "old"})
	Register(App{Name: "ready"})
	Register(App{Name: "show", Summary: "new"})
	got := Apps()
	if len(got) != 2 || got[0].Name != "ready" || got[1].Name != "show" || got[1].Summary != "new" {
		t.Fatalf("Apps() = %+v", got)
	}
}

func TestLayerListsEachAppByName(t *testing.T) {
	m := newTestModel(&fakeRex{}, testApps())
	m.Update(tea.WindowSizeMsg{Width: 60, Height: 20})
	view := m.View().Content
	for _, a := range testApps() {
		if !strings.Contains(view, a.Name) {
			t.Errorf("layer does not list %q:\n%s", a.Name, view)
		}
	}
}

func TestSelectingAnAppOpensItsModal(t *testing.T) {
	f := &fakeRex{}
	m := newTestModel(f, testApps())
	m.Update(tea.WindowSizeMsg{Width: 60, Height: 20})
	m, _ = press(m, tea.KeyDown) // roadmap
	m, cmd := press(m, tea.KeyEnter)
	if cmd == nil {
		t.Fatal("enter returned no command")
	}
	// The list reports the choice as a message; the host acts on it.
	for cmd != nil {
		msg := cmd()
		if _, quit := msg.(tea.QuitMsg); quit {
			break
		}
		var next tea.Model
		next, cmd = m.Update(msg)
		m = next.(*Model)
	}
	if m.Err != nil {
		t.Fatal(m.Err)
	}
	if len(f.layers) != 1 {
		t.Fatalf("opened %d layers, want 1", len(f.layers))
	}
	b, _ := json.Marshal(f.layers[0].Layout)
	for _, want := range []string{`"/opt/rhiza/bin/rhiza"`, `"tui"`, `"app"`, `"roadmap"`, `"/work"`} {
		if !strings.Contains(string(b), want) {
			t.Errorf("modal layout lacks %s: %s", want, b)
		}
	}
	if !f.layers[0].Focus {
		t.Error("modal is not focused")
	}
	if len(f.calls) != 1 || f.calls[0] != "session.disown_block" {
		t.Errorf("calls = %v, want one session.disown_block", f.calls)
	}
}

func TestEscapeClosesWithoutOpening(t *testing.T) {
	f := &fakeRex{}
	m := newTestModel(f, testApps())
	m.Update(tea.WindowSizeMsg{Width: 60, Height: 20})
	_, cmd := press(m, tea.KeyEscape)
	for cmd != nil {
		msg := cmd()
		if _, quit := msg.(tea.QuitMsg); quit {
			break
		}
		next, c := m.Update(msg)
		m, cmd = next.(*Model), c
	}
	if len(f.layers) != 0 {
		t.Fatalf("esc opened %d layers", len(f.layers))
	}
}
