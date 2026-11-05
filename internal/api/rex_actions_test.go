package api

import (
	"math"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"

	"github.com/dphaener/rexuikit/surface"
	"github.com/dphaener/rhiza/internal/tui/launcher"
)

// rexLayer is one rex.session.new_layer call an app's action made.
type rexLayer struct {
	window  string
	bounds  []float64
	flavor  string
	label   string
	command []string
	cwd     string
	focus   bool
}

// rexRan is what one app's action did under the fake rex.
type rexRan struct {
	layers []rexLayer
	execs  int // rex.exec calls: an action must make none
}

// The fake session: block "b" fills window "w2" of a 200x50 terminal at
// 9x18 px cells; window "w1" holds another block.
const (
	fakeCols, fakeRows   = 200, 50
	fakeCellW, fakeCellH = 9, 18
)

// luaTable turns kv into a Lua table of L.
func luaTable(L *lua.LState, kv map[string]lua.LValue) *lua.LTable {
	x := L.NewTable()
	for k, v := range kv {
		x.RawSetString(k, v)
	}
	return x
}

// luaStrings returns the array part of t as strings.
func luaStrings(t lua.LValue) []string {
	tb, ok := t.(*lua.LTable)
	if !ok {
		return nil
	}
	var out []string
	for i := 1; i <= tb.Len(); i++ {
		out = append(out, tb.RawGetInt(i).String())
	}
	return out
}

// runRexActions loads script, the written rhiza.lua, in a Lua state with a
// fake rex whose block "b" in session "s" runs a shell in /repo, with SHELL
// set to shell ("" for none). It runs every registered app's action as Rex
// would on block "b" and returns, by app name, what it did.
func runRexActions(t *testing.T, script, shell string) map[string]rexRan {
	t.Helper()
	L := lua.NewState()
	defer L.Close()
	actions := map[string]*lua.LFunction{}
	var ran rexRan
	noop := func(*lua.LState) int { return 0 }
	fn := func(f lua.LGFunction) *lua.LFunction { return L.NewFunction(f) }
	tbl := func(kv map[string]lua.LValue) *lua.LTable { return luaTable(L, kv) }
	rect := func(x, y, w, h float64) *lua.LTable {
		r := L.NewTable()
		for _, v := range []float64{x, y, w, h} {
			r.Append(lua.LNumber(v))
		}
		return r
	}
	L.SetField(L.GetGlobal("os"), "getenv", fn(func(L *lua.LState) int {
		if L.CheckString(1) == "SHELL" && shell != "" {
			L.Push(lua.LString(shell))
			return 1
		}
		L.Push(lua.LNil)
		return 1
	}))
	windows := L.NewTable()
	windows.Append(tbl(map[string]lua.LValue{"window_id": lua.LString("w1"), "layers": func() *lua.LTable {
		ls := L.NewTable()
		ls.Append(tbl(map[string]lua.LValue{"kind": lua.LString("tiled"), "layout": tbl(map[string]lua.LValue{
			"block_id": lua.LString("other"), "rect": rect(0, 0, 1, 1)})}))
		return ls
	}()}))
	windows.Append(tbl(map[string]lua.LValue{"window_id": lua.LString("w2"), "layers": func() *lua.LTable {
		ls := L.NewTable()
		ls.Append(tbl(map[string]lua.LValue{"kind": lua.LString("tiled"), "layout": tbl(map[string]lua.LValue{
			"block_id": lua.LString("b"), "rect": rect(0, 0, 1, 1)})}))
		return ls
	}()}))
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
		"block": tbl(map[string]lua.LValue{"call": fn(func(L *lua.LState) int {
			switch L.CheckString(2) {
			case "process":
				L.Push(tbl(map[string]lua.LValue{"child": tbl(map[string]lua.LValue{"cwd": lua.LString("/repo")})}))
			case "size":
				L.Push(tbl(map[string]lua.LValue{
					"columns": lua.LNumber(fakeCols), "rows": lua.LNumber(fakeRows),
					"cell_width_px": lua.LNumber(fakeCellW), "cell_height_px": lua.LNumber(fakeCellH),
				}))
			default:
				L.Push(lua.LNil)
			}
			return 1
		})}),
		"layout": tbl(map[string]lua.LValue{"block": fn(func(L *lua.LState) int {
			L.Push(L.CheckTable(1))
			return 1
		})}),
		"session": tbl(map[string]lua.LValue{
			"view": fn(func(L *lua.LState) int {
				L.Push(tbl(map[string]lua.LValue{"windows": windows}))
				return 1
			}),
			"new_layer": fn(func(L *lua.LState) int {
				o := L.CheckTable(1)
				var l rexLayer
				l.window = lua.LVAsString(o.RawGetString("window_id"))
				if b, ok := o.RawGetString("bounds").(*lua.LTable); ok {
					for i := 1; i <= b.Len(); i++ {
						l.bounds = append(l.bounds, float64(lua.LVAsNumber(b.RawGetInt(i))))
					}
				}
				l.focus = lua.LVAsBool(o.RawGetString("focus"))
				if lay, ok := o.RawGetString("layout").(*lua.LTable); ok {
					l.flavor = lua.LVAsString(lay.RawGetString("flavor"))
					l.label = lua.LVAsString(lay.RawGetString("label"))
					if opts, ok := lay.RawGetString("options").(*lua.LTable); ok {
						l.command = luaStrings(opts.RawGetString("command"))
						l.cwd = lua.LVAsString(opts.RawGetString("cwd"))
					}
				}
				ran.layers = append(ran.layers, l)
				L.Push(L.NewTable())
				return 1
			}),
		}),
		"exec": fn(func(L *lua.LState) int {
			ran.execs++
			return 0
		}),
	})
	L.SetGlobal("rex", rex)
	if err := L.DoString(script); err != nil {
		t.Fatalf("loading rhiza.lua: %v", err)
	}
	out := map[string]rexRan{}
	for _, a := range launcher.Apps() {
		run, ok := actions[rexActionName(a.Name)]
		if !ok {
			t.Fatalf("rhiza.lua registers no action for %s", a.Name)
		}
		ran = rexRan{}
		ctx := tbl(map[string]lua.LValue{"block_id": lua.LString("b"), "session_id": lua.LString("s")})
		if err := L.CallByParam(lua.P{Fn: run, NRet: 1, Protect: true}, ctx); err != nil {
			t.Fatalf("%s: running its action: %v", a.Name, err)
		}
		L.Pop(1)
		out[a.Name] = ran
	}
	return out
}

// TestRexInstallActionsOpenALayerPerApp installs rhiza.lua, runs every
// registered app's action under a fake rex, and checks each opens exactly
// one focused layer in the invoking block's window, sized to the app's
// cells (launcher.App.Size) as surface.BoundsForCells centres them, with a
// terminal block running `rhiza tui app <name>` in the invoking block's
// directory: through the owner's login shell, or, with settings.bin set,
// that binary as given. No action calls rex.exec, whose detach starts
// nothing on Rex r10.
func TestRexInstallActionsOpenALayerPerApp(t *testing.T) {
	launcher.Register(launcher.App{Name: "zz-opened", Summary: "Opened"})
	ctx, dir, _ := rexTestCtx(t)
	if _, err := callOp(t, ctx, "rex.install", `{}`); err != nil {
		t.Fatal(err)
	}
	script := readFile(t, filepath.Join(dir, "rhiza.lua"))
	wc := surface.WindowCells{Cols: fakeCols, Rows: fakeRows, CellWPx: fakeCellW, CellHPx: fakeCellH}
	check := func(ran map[string]rexRan, prefix []string) {
		t.Helper()
		for _, a := range launcher.Apps() {
			got := ran[a.Name]
			if got.execs != 0 {
				t.Errorf("%s: the action called rex.exec %d times", a.Name, got.execs)
			}
			if len(got.layers) != 1 {
				t.Errorf("%s: the action opened %d layers, want 1", a.Name, len(got.layers))
				continue
			}
			l := got.layers[0]
			want := append(slices.Clone(prefix), "tui", "app", a.Name)
			if !slices.Equal(l.command, want) {
				t.Errorf("%s: layer runs %q, want %q", a.Name, l.command, want)
			}
			if l.cwd != "/repo" || l.window != "w2" || !l.focus {
				t.Errorf("%s: layer cwd %q window %q focus %v, want /repo, w2, focused", a.Name, l.cwd, l.window, l.focus)
			}
			if l.flavor != "com.superlogical.terminal.shell" || l.label != a.Name {
				t.Errorf("%s: layer flavor %q label %q", a.Name, l.flavor, l.label)
			}
			cols, rows := a.Size()
			wb, ok := surface.BoundsForCells(wc, cols, rows, surface.Centered)
			if !ok {
				t.Fatalf("%s: %dx%d does not fit the fake window", a.Name, cols, rows)
			}
			if len(l.bounds) != 4 {
				t.Errorf("%s: bounds %v", a.Name, l.bounds)
				continue
			}
			for i := range wb {
				if math.Abs(l.bounds[i]-wb[i]) > 1e-9 {
					t.Errorf("%s: bounds %v, want %v for %dx%d", a.Name, l.bounds, wb, cols, rows)
					break
				}
			}
		}
	}
	check(runRexActions(t, script, "/bin/zsh"), []string{"/bin/zsh", "-lc", `exec "$0" "$@"`, "rhiza"})
	check(runRexActions(t, script, ""), []string{"/bin/sh", "-lc", `exec "$0" "$@"`, "rhiza"})

	// settings.bin names the binary: it runs as given, with no shell.
	pinned := strings.Replace(script, `  -- bin = "rhiza",`, `  bin = "/opt/rhiza/bin/rhiza",`, 1)
	if pinned == script {
		t.Fatal("no bin setting to set in rhiza.lua")
	}
	check(runRexActions(t, pinned, "/bin/zsh"), []string{"/opt/rhiza/bin/rhiza"})
	if strings.Contains(script, "rex.exec(") {
		t.Errorf("rhiza.lua calls rex.exec")
	}
}
