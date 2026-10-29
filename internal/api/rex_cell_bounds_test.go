package api

import (
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/dphaener/rexuikit/surface"
)

// TestRexCellBoundsScalesTheBlockToItsWindow runs the installed
// cell_bounds under luajit, with a stub rex, and checks its bounds against
// surface.BoundsForCells over the window the block's rect implies: a block
// filling the window, a half-width split and a 70/30 one. It needs luajit
// on $PATH and is skipped without it.
func TestRexCellBoundsScalesTheBlockToItsWindow(t *testing.T) {
	luajit, err := exec.LookPath("luajit")
	if err != nil {
		t.Skip("luajit not on $PATH")
	}
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
			lua := `local TERM = "term"
local LAYER_BOUNDS = { -1, -1, -1, -1 }
rex = {
  session = { view = function() return { windows = { { layers = {
    { kind = "floating", layout = { block_id = "other", rect = { 0, 0, 1, 1 } } },
    { kind = "tiled", layout = { rect = { 0, 0, 1, 1 }, split = {
      before = { block_id = "left", rect = { 0, 0, 0.5, 1 } },
      after = { block_id = "b", rect = { ` + fmt.Sprintf("%g, %g, %g, %g", c.rect[0], c.rect[1], c.rect[2], c.rect[3]) + ` } },
    } } },
  } } } } end },
  block = { call = function(_, _, a)
    return { columns = ` + strconv.Itoa(c.cols) + `, rows = ` + strconv.Itoa(c.rows) + `, cell_width_px = 9, cell_height_px = 18 }
  end },
}
` + rexCellBounds() + fmt.Sprintf(`
local b = cell_bounds("b", %d, %d)
print(string.format("%%.17g %%.17g %%.17g %%.17g", b[1], b[2], b[3], b[4]))
`, c.want[0], c.want[1])
			cmd := exec.Command(luajit, "-")
			cmd.Stdin = strings.NewReader(lua)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("luajit: %v\n%s", err, out)
			}
			var got [4]float64
			for i, f := range strings.Fields(string(out)) {
				if got[i], err = strconv.ParseFloat(f, 64); err != nil {
					t.Fatalf("bounds %q: %v", out, err)
				}
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

// TestRexBlockWindowFindsTheFocusedBlocksWindow runs the installed
// block_window under luajit, with a stub rex holding two windows, and checks
// that an app's layer goes to the window of the focused block, in a tiled or
// floating layer, and to none (the active one) without a block. It needs
// luajit on $PATH and is skipped without it.
func TestRexBlockWindowFindsTheFocusedBlocksWindow(t *testing.T) {
	luajit, err := exec.LookPath("luajit")
	if err != nil {
		t.Skip("luajit not on $PATH")
	}
	lua := `local TERM = "term"
local LAYER_BOUNDS = { -1, -1, -1, -1 }
rex = {
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
` + rexCellBounds() + `
print(block_window("a"), block_window("c"), block_window("d"), block_window("gone"), block_window(nil))
`
	cmd := exec.Command(luajit, "-")
	cmd.Stdin = strings.NewReader(lua)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("luajit: %v\n%s", err, out)
	}
	if got, want := strings.Fields(string(out)), []string{"w1", "w2", "w2", "nil", "nil"}; strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("block_window gave %v, want %v", got, want)
	}
}
