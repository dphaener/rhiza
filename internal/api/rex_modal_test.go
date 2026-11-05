package api

import (
	"math"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/colorprofile"
	"github.com/dphaener/rexuikit/surface"
	"github.com/dphaener/rexuikit/ui"
	"github.com/dphaener/rexuikit/ui/theme"
	lua "github.com/yuin/gopher-lua"

	"github.com/dphaener/rhiza/internal/tui/launcher"
	"github.com/dphaener/rhiza/internal/tui/modalshell"
)

// rexLayer is what an app's action passed to rex.session.new_layer.
type rexLayer struct {
	window string
	bounds [4]float64
}

// runRexActions loads script, the written rhiza.lua, in a Lua state with a
// fake rex: one window "w1" whose tiled layout is the single block "b" of
// cols x rows cells of 9x18 pixels. It runs every registered app's action
// as Rex would on block "b" and returns, by app name, the layer it opened.
func runRexActions(t *testing.T, script string, cols, rows int) map[string]rexLayer {
	t.Helper()
	L := lua.NewState()
	defer L.Close()
	actions := map[string]*lua.LFunction{}
	var opened *lua.LTable
	noop := func(*lua.LState) int { return 0 }
	fn := func(f lua.LGFunction) *lua.LFunction { return L.NewFunction(f) }
	tbl := func(kv map[string]lua.LValue) *lua.LTable {
		x := L.NewTable()
		for k, v := range kv {
			x.RawSetString(k, v)
		}
		return x
	}
	list := func(vs ...lua.LValue) *lua.LTable {
		x := L.NewTable()
		for _, v := range vs {
			x.Append(v)
		}
		return x
	}
	rect := list(lua.LNumber(0), lua.LNumber(0), lua.LNumber(1), lua.LNumber(1))
	view := tbl(map[string]lua.LValue{"windows": list(tbl(map[string]lua.LValue{
		"window_id": lua.LString("w1"),
		"layers": list(tbl(map[string]lua.LValue{
			"kind":   lua.LString("tiled"),
			"layout": tbl(map[string]lua.LValue{"block_id": lua.LString("b"), "rect": rect}),
		})),
	}))})
	rex := tbl(map[string]lua.LValue{
		"mode": fn(noop),
		"bind": fn(noop),
		"action": fn(func(L *lua.LState) int {
			a := L.CheckTable(1)
			if run, ok := a.RawGetString("run").(*lua.LFunction); ok {
				actions[a.RawGetString("name").String()] = run
			}
			return 0
		}),
		"layout": tbl(map[string]lua.LValue{"block": fn(func(L *lua.LState) int {
			L.Push(L.CheckTable(1))
			return 1
		})}),
		"block": tbl(map[string]lua.LValue{"call": fn(func(L *lua.LState) int {
			if L.CheckString(2) != "size" {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(tbl(map[string]lua.LValue{
				"columns": lua.LNumber(cols), "rows": lua.LNumber(rows),
				"cell_width_px": lua.LNumber(9), "cell_height_px": lua.LNumber(18),
			}))
			return 1
		})}),
		"session": tbl(map[string]lua.LValue{
			"view": fn(func(L *lua.LState) int { L.Push(view); return 1 }),
			"new_layer": fn(func(L *lua.LState) int {
				opened = L.CheckTable(1)
				L.Push(L.NewTable())
				return 1
			}),
		}),
	})
	L.SetGlobal("rex", rex)
	if err := L.DoString(script); err != nil {
		t.Fatalf("loading rhiza.lua: %v", err)
	}
	out := map[string]rexLayer{}
	for _, a := range launcher.Apps() {
		run, ok := actions[rexActionName(a.Name)]
		if !ok {
			t.Fatalf("rhiza.lua registers no action for %s", a.Name)
		}
		opened = nil
		ctx := tbl(map[string]lua.LValue{"block_id": lua.LString("b")})
		if err := L.CallByParam(lua.P{Fn: run, NRet: 1, Protect: true}, ctx); err != nil {
			t.Fatalf("%s: running its action: %v", a.Name, err)
		}
		L.Pop(1)
		if opened == nil {
			t.Fatalf("%s: its action opened no layer", a.Name)
		}
		var l rexLayer
		l.window = opened.RawGetString("window_id").String()
		b, ok := opened.RawGetString("bounds").(*lua.LTable)
		if !ok || b.Len() != 4 {
			t.Fatalf("%s: bounds %v", a.Name, opened.RawGetString("bounds"))
		}
		for i := range l.bounds {
			l.bounds[i] = float64(lua.LVAsNumber(b.RawGetInt(i + 1)))
		}
		out[a.Name] = l
	}
	return out
}

// TestRexInstallOpensEveryAppAsACentredModalSizedToItsComponent installs
// rhiza.lua, runs every registered app's action from it under a fake rex,
// and checks that each opens its layer in the focused block's window with
// exactly the bounds surface.BoundsForCells gives a Modal (zero Anchor,
// centred) of the cells of the app's component NaturalSize, and never the
// whole window. Apps hosting a Form, a List and a Dialog are registered
// beside the real ones; growing a component's content resizes its modal on
// the next install.
func TestRexInstallOpensEveryAppAsACentredModalSizedToItsComponent(t *testing.T) {
	th, prof := theme.Dark(), colorprofile.TrueColor
	rows := []ui.Row{{Label: "one"}, {Label: "two"}}
	comps := map[string]func() modalshell.Component{
		"zz-modal-form": func() modalshell.Component {
			return modalshell.Form(ui.NewForm(th, prof, ui.NewText("intent", "Intent"), ui.NewTextArea("body", "Body")))
		},
		"zz-modal-list": func() modalshell.Component { return modalshell.List(ui.NewList(th, prof, rows)) },
		"zz-modal-dialog": func() modalshell.Component {
			return modalshell.Dialog(ui.NewDialog(th, prof, "Land", []string{"Land the accepted change set?"}))
		},
	}
	for name, c := range comps {
		launcher.Register(launcher.App{Name: name, Summary: name, Component: c})
	}
	ctx, dir, _ := rexTestCtx(t)
	wc := surface.WindowCells{Cols: 240, Rows: 70, CellWPx: 9, CellHPx: 18}

	install := func() map[string]rexLayer {
		t.Helper()
		if _, err := callOp(t, ctx, "rex.install", `{}`); err != nil {
			t.Fatal(err)
		}
		return runRexActions(t, readFile(t, filepath.Join(dir, "rhiza.lua")), wc.Cols, wc.Rows)
	}
	check := func(layers map[string]rexLayer) {
		t.Helper()
		for _, a := range launcher.Apps() {
			cols, rows := a.Size()
			if c, ok := comps[a.Name]; ok {
				if nc, nr := c().NaturalSize(); cols != nc || rows != nr {
					t.Fatalf("%s: Size %dx%d, its component's NaturalSize %dx%d", a.Name, cols, rows, nc, nr)
				}
			}
			want, ok := surface.BoundsForCells(wc, cols, rows, surface.Centered)
			if !ok {
				t.Fatalf("%s: window %+v too small for %dx%d", a.Name, wc, cols, rows)
			}
			got := layers[a.Name]
			if got.window != "w1" {
				t.Errorf("%s: layer opened in window %q, want the focused block's w1", a.Name, got.window)
			}
			for i := range want {
				if math.Abs(got.bounds[i]-want[i]) > 1e-9 {
					t.Errorf("%s: bounds %v, want BoundsForCells' centred Modal %v for %dx%d", a.Name, got.bounds, want, cols, rows)
					break
				}
			}
			if got.bounds[2] >= 1 || got.bounds[3] >= 1 {
				t.Errorf("%s: bounds %v fill the window", a.Name, got.bounds)
			}
			if u := surface.Usable(wc, got.bounds); u.Columns < cols || u.Rows < rows {
				t.Errorf("%s: bounds %v give %dx%d, want %dx%d", a.Name, got.bounds, u.Columns, u.Rows, cols, rows)
			}
		}
	}
	first := install()
	check(first)

	// The list grows: its modal follows on the next install.
	before := first["zz-modal-list"].bounds
	more := append(append([]ui.Row{}, rows...), ui.Row{Label: "a far longer row than either of the first two"}, ui.Row{Label: "four"})
	comps["zz-modal-list"] = func() modalshell.Component { return modalshell.List(ui.NewList(th, prof, more)) }
	launcher.Register(launcher.App{Name: "zz-modal-list", Summary: "zz-modal-list", Component: comps["zz-modal-list"]})
	second := install()
	check(second)
	if after := second["zz-modal-list"].bounds; after[2] <= before[2] || after[3] <= before[3] {
		t.Errorf("the grown list's modal %v is not larger than %v", after, before)
	}
}
