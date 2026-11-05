package modalopen

import (
	"context"
	"math"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/dphaener/rexuikit/client"
	"github.com/dphaener/rexuikit/surface"

	"github.com/dphaener/rhiza/internal/tui/modalshell"
)

// fakeRex is a session with two windows: w1 holds the opener block "b" and
// w2, the active one, block "other". Each tiled block is the whole window
// of 240x70 cells at 9x18 px. It records new_layer and resize_layer calls;
// close sends the block_closed event of the modal's block.
type fakeRex struct {
	mu      sync.Mutex
	opened  []client.NewLayerArgs
	resized []client.Bounds
	events  chan client.Event
}

func newFake() *fakeRex { return &fakeRex{events: make(chan client.Event, 1)} }

// close sends block_closed for the modal's block: its program exited.
func (f *fakeRex) close() {
	f.events <- client.Event{Type: "session_event",
		Payload: []byte(`{"session_id":"s","event":{"type":"block_closed","payload":{"block_id":"block:m"}}}`)}
}

var window = surface.WindowCells{Cols: 240, Rows: 70, CellWPx: 9, CellHPx: 18}

func (f *fakeRex) View(context.Context, string) (client.View, error) {
	tiled := func(block string) []client.Layer {
		return []client.Layer{{LayerID: "tiled:" + block, Kind: client.LayerTiled,
			Layout: &client.LayoutNode{Rect: client.Bounds{0, 0, 1, 1}, BlockID: block}}}
	}
	return client.View{SessionID: "s", ActiveWindowID: "w2", Windows: []client.Window{
		{WindowID: "w1", Layers: tiled("b")},
		{WindowID: "w2", Layers: tiled("other")},
	}}, nil
}

func (f *fakeRex) TerminalSize(context.Context, string, string) (client.TerminalSize, error) {
	return client.TerminalSize{Columns: window.Cols, Rows: window.Rows, CellWidthPx: window.CellWPx, CellHeightPx: window.CellHPx}, nil
}

func (f *fakeRex) NewLayer(_ context.Context, a client.NewLayerArgs) (client.NewLayerResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.opened = append(f.opened, a)
	return client.NewLayerResult{LayerID: "layer:m", BlockIDs: []string{"block:m"}, Revision: 1}, nil
}

func (f *fakeRex) ResizeLayer(_ context.Context, _, _ string, b client.Bounds) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.resized = append(f.resized, b)
	return 2, nil
}

func (f *fakeRex) CloseLayer(context.Context, string, string) (int64, error) { return 0, nil }
func (f *fakeRex) RaiseLayer(context.Context, string, string, string) (int64, error) {
	return 0, nil
}
func (f *fakeRex) NewSplit(context.Context, client.NewSplitArgs) (client.NewSplitResult, error) {
	return client.NewSplitResult{}, nil
}
func (f *fakeRex) NewBlock(context.Context, client.NewBlockArgs) (string, error) { return "", nil }
func (f *fakeRex) MoveBlock(context.Context, client.MoveBlockArgs) (client.MoveBlockResult, error) {
	return client.MoveBlockResult{}, nil
}
func (f *fakeRex) FocusBlock(context.Context, string, string) (int64, error) { return 0, nil }
func (f *fakeRex) SetBlockLabel(context.Context, string, string, string) (int64, error) {
	return 0, nil
}
func (f *fakeRex) ListBlocks(context.Context, string) ([]client.BlockInfo, error) { return nil, nil }
func (f *fakeRex) CloseBlock(context.Context, string, string) (int64, error)      { return 0, nil }
func (f *fakeRex) Call(context.Context, string, any, any) error                   { return nil }
func (f *fakeRex) Events() <-chan client.Event                                    { return f.events }

func (f *fakeRex) resizes() []client.Bounds {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]client.Bounds(nil), f.resized...)
}

func near(a, b client.Bounds) bool {
	for i := range a {
		if math.Abs(a[i]-b[i]) > 1e-9 {
			return false
		}
	}
	return true
}

// TestOpenOpensACentredModalAndFollowsItsContent opens an app's modal over
// the opener block's window and checks that new_layer gets the bounds
// surface.BoundsForCells gives a centred Modal of the content's cells, never
// the whole window; that each new content size resizes the layer to
// BoundsForCells of it; and that Open returns once the modal is gone.
func TestOpenOpensACentredModalAndFollowsItsContent(t *testing.T) {
	f := newFake()
	sizes := make(chan client.Size)
	done := make(chan error, 1)
	go func() {
		done <- Open(context.Background(), f, Request{
			SessionID: "s", BlockID: "b", Label: "inbox",
			Argv: []string{"/bin/rhiza", "tui", "app", "inbox"}, Cwd: "/repo",
			Cols: 60, Rows: 12, Sizes: sizes,
		})
	}()
	sizes <- client.Size{Columns: 60, Rows: 12} // the first size again: no resize
	sizes <- client.Size{Columns: 90, Rows: 30} // the list grows
	sizes <- client.Size{Columns: 40, Rows: 8}  // and shrinks
	sizes <- client.Size{Columns: 40, Rows: 8}

	f.mu.Lock()
	if len(f.opened) != 1 {
		f.mu.Unlock()
		t.Fatalf("new_layer called %d times, want 1", len(f.opened))
	}
	a := f.opened[0]
	f.mu.Unlock()
	if a.WindowID != "w1" {
		t.Errorf("modal opened in window %q, want the opener block's w1", a.WindowID)
	}
	want, _ := surface.BoundsForCells(window, 60, 12, surface.Centered)
	if !near(a.Bounds, want) {
		t.Errorf("modal bounds %v, want BoundsForCells' centred %v", a.Bounds, want)
	}
	if a.Bounds[2] >= 1 || a.Bounds[3] >= 1 {
		t.Errorf("modal bounds %v fill the window", a.Bounds)
	}
	if !a.Focus {
		t.Error("modal does not take focus")
	}
	if u := surface.Usable(window, a.Bounds); u.Columns < 60 || u.Rows < 12 {
		t.Errorf("modal gives %dx%d, want at least 60x12", u.Columns, u.Rows)
	}

	f.close() // the app exited: the modal is gone
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Open did not return once the modal was gone")
	}
	grown, _ := surface.BoundsForCells(window, 90, 30, surface.Centered)
	shrunk, _ := surface.BoundsForCells(window, 40, 8, surface.Centered)
	got := f.resizes()
	if len(got) != 2 || !near(got[0], grown) || !near(got[1], shrunk) {
		t.Errorf("resizes %v, want [%v %v]", got, grown, shrunk)
	}
}

// TestOpenWithoutAKnownBlockUsesTheActiveWindow checks a block no window
// holds opens the modal in the active window.
func TestOpenWithoutAKnownBlockUsesTheActiveWindow(t *testing.T) {
	f := newFake()
	f.close()
	if err := Open(context.Background(), f, Request{SessionID: "s", BlockID: "gone", Label: "x",
		Argv: []string{"x"}, Cols: 10, Rows: 5}); err != nil {
		t.Fatal(err)
	}
	if w := f.opened[0].WindowID; w != "w2" {
		t.Errorf("window %q, want the active w2", w)
	}
}

// TestWatchSizeSendsEachNewSize checks WatchSize sends what modalshell
// writes, once per change.
func TestWatchSizeSendsEachNewSize(t *testing.T) {
	path := filepath.Join(t.TempDir(), "size")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := WatchSize(ctx, path, time.Millisecond)
	next := func() client.Size {
		t.Helper()
		select {
		case sz := <-ch:
			return sz
		case <-time.After(5 * time.Second):
			t.Fatal("no size")
		}
		return client.Size{}
	}
	if err := modalshell.WriteSize(path, 30, 4); err != nil {
		t.Fatal(err)
	}
	if sz := next(); sz != (client.Size{Columns: 30, Rows: 4}) {
		t.Fatalf("got %+v", sz)
	}
	if err := modalshell.WriteSize(path, 31, 9); err != nil {
		t.Fatal(err)
	}
	if sz := next(); sz != (client.Size{Columns: 31, Rows: 9}) {
		t.Fatalf("got %+v", sz)
	}
	cancel()
	for range ch {
	}
}

func TestServerAddr(t *testing.T) {
	env := func(kv map[string]string) func(string) string { return func(k string) string { return kv[k] } }
	if a, _ := ServerAddr(env(map[string]string{"REX_SERVER": "unix:///x.sock"})); a != "unix:///x.sock" {
		t.Errorf("REX_SERVER: %q", a)
	}
	dir := t.TempDir()
	sock := filepath.Join(dir, "rex", "server.sock")
	if a, _ := ServerAddr(env(map[string]string{"XDG_RUNTIME_DIR": dir, "HOME": dir})); a == "unix://"+sock {
		t.Errorf("no socket under XDG_RUNTIME_DIR, but got %q", a)
	}
	if err := os.MkdirAll(filepath.Dir(sock), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sock, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if a, _ := ServerAddr(env(map[string]string{"XDG_RUNTIME_DIR": dir})); a != "unix://"+sock {
		t.Errorf("XDG_RUNTIME_DIR: %q, want unix://%s", a, sock)
	}
}
