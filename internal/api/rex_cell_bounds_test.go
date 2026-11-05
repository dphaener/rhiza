package api

import (
	"fmt"
	"math"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"

	"github.com/dphaener/rexuikit/surface"
)

// runCellBoundsLua runs the installed cell-bounds Lua after stub, a Lua
// chunk defining the fake rex, then body, in gopher-lua, and returns the
// values body leaves in the global "out".
func runCellBoundsLua(t *testing.T, stub, body string) lua.LValue {
	t.Helper()
	L := lua.NewState()
	defer L.Close()
	src := "local TERM = \"term\"\nlocal LAYER_BOUNDS = { -1, -1, -1, -1 }\n" + stub + rexCellBounds() + body
	if err := L.DoString(src); err != nil {
		t.Fatalf("lua: %v\n%s", err, src)
	}
	return L.GetGlobal("out")
}

// TestRexCellBoundsScalesTheBlockToItsWindow runs the installed
// cell_bounds with a stub rex and checks its bounds against
// surface.BoundsForCells over the window the block's rect implies: a block
// filling the window, a half-width split and a 70/30 one.
func TestRexCellBoundsScalesTheBlockToItsWindow(t *testing.T) {
	for _, c := range []struct {
		name       string
		cols, rows int        // the block's terminal size
		rect       [4]float64 // the block's rect in the window
		want       [2]int     // the app's cells
	}{
		{"whole window", 200, 50, [4]float64{0, 0, 1, 1}, [2]int{100, 30}},
		{"half width", 100, 50, [4]float64{0.5, 0, 0.5, 1}, [2]int{100, 30}},
		{"70/30 below", 200, 15, [4]float64{0, 0.7, 1, 0.3}, [2]int{80, 20}},
	} {
		t.Run(c.name, func(t *testing.T) {
			stub := fmt.Sprintf(`rex = {
  session = { view = function() return { windows = { { layers = {
    { kind = "floating", layout = { block_id = "other", rect = { 0, 0, 1, 1 } } },
    { kind = "tiled", layout = { rect = { 0, 0, 1, 1 }, split = {
      before = { block_id = "left", rect = { 0, 0, 0.5, 1 } },
      after = { block_id = "b", rect = { %g, %g, %g, %g } },
    } } },
  } } } } end },
  block = { call = function(_, _, a)
    return { columns = %d, rows = %d, cell_width_px = 9, cell_height_px = 18 }
  end },
}
`, c.rect[0], c.rect[1], c.rect[2], c.rect[3], c.cols, c.rows)
			out := runCellBoundsLua(t, stub, fmt.Sprintf("\nout = cell_bounds(\"b\", %d, %d)\n", c.want[0], c.want[1]))
			tb, ok := out.(*lua.LTable)
			if !ok || tb.Len() != 4 {
				t.Fatalf("cell_bounds gave %v", out)
			}
			var got [4]float64
			for i := range got {
				got[i] = float64(lua.LVAsNumber(tb.RawGetInt(i + 1)))
			}
			wc := surface.WindowCells{
				Cols:    int(math.Round(float64(c.cols) / c.rect[2])),
				Rows:    int(math.Round(float64(c.rows) / c.rect[3])),
				CellWPx: 9, CellHPx: 18,
			}
			want, ok := surface.BoundsForCells(wc, c.want[0], c.want[1], surface.Centered)
			if !ok {
				t.Fatalf("window %+v too small for %v", wc, c.want)
			}
			for i := range got {
				if math.Abs(got[i]-want[i]) > 1e-9 {
					t.Fatalf("bounds %v, want %v (window %+v)", got, want, wc)
				}
			}
			if u := surface.Usable(wc, got); u.Columns < c.want[0] || u.Rows < c.want[1] {
				t.Errorf("bounds %v give %dx%d, want %v", got, u.Columns, u.Rows, c.want)
			}
		})
	}
}

// TestRexCellBoundsFallsBackWithoutASize checks cell_bounds gives
// LAYER_BOUNDS for no block, a block in a floating layer, and a block the
// view does not hold.
func TestRexCellBoundsFallsBackWithoutASize(t *testing.T) {
	stub := `rex = {
  session = { view = function() return { windows = { { layers = {
    { kind = "floating", layout = { block_id = "f", rect = { 0, 0, 1, 1 } } },
  } } } } end },
  block = { call = function() return { columns = 80, rows = 24, cell_width_px = 9, cell_height_px = 18 } end },
}
`
	out := runCellBoundsLua(t, stub, "\nout = { cell_bounds(nil, 10, 10)[1], cell_bounds(\"f\", 10, 10)[1], cell_bounds(\"gone\", 10, 10)[1] }\n")
	tb := out.(*lua.LTable)
	for i := 1; i <= 3; i++ {
		if v := lua.LVAsNumber(tb.RawGetInt(i)); v != -1 {
			t.Errorf("case %d: bounds start %v, want LAYER_BOUNDS", i, v)
		}
	}
}

// TestRexBlockWindowFindsTheFocusedBlocksWindow runs the installed
// block_window with a stub rex holding two windows, and checks that an
// app's layer goes to the window of the focused block, in a tiled or
// floating layer, and to none (the active one) without a block.
func TestRexBlockWindowFindsTheFocusedBlocksWindow(t *testing.T) {
	stub := `rex = {
  session = { view = function() return { windows = {
    { window_id = "w1", layers = {
      { kind = "tiled", layout = { block_id = "a", rect = { 0, 0, 1, 1 } } },
    } },
    { window_id = "w2", layers = {
      { kind = "tiled", layout = { rect = { 0, 0, 1, 1 }, split = {
        before = { block_id = "b", rect = { 0, 0, 0.5, 1 } },
        after = { block_id = "c", rect = { 0.5, 0, 0.5, 1 } },
      } } },
      { kind = "floating", layout = { block_id = "d", rect = { 0, 0, 1, 1 } } },
    } },
  } } end },
}
`
	out := runCellBoundsLua(t, stub, `
out = { tostring(block_window("a")), tostring(block_window("c")), tostring(block_window("d")), tostring(block_window("gone")), tostring(block_window(nil)) }
`)
	got := luaStrings(out)
	if want := []string{"w1", "w2", "w2", "nil", "nil"}; strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("block_window gave %v, want %v", got, want)
	}
}
