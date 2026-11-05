package api

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"

	"github.com/dphaener/rhiza/internal/tui/launcher"
)

// rexExec is what an app's action passed to rex.exec.
type rexExec struct {
	argv   []string
	detach bool
	cwd    string
	layers int // rex.session.new_layer calls: an action must open none
}

// runRexActions loads script, the written rhiza.lua, in a Lua state with a
// fake rex whose block "b" in session "s" runs a shell in /repo, with SHELL
// set to shell ("" for none). It runs every registered app's action as Rex
// would on block "b" and returns, by app name, what it ran.
func runRexActions(t *testing.T, script, shell string) map[string]rexExec {
	t.Helper()
	L := lua.NewState()
	defer L.Close()
	actions := map[string]*lua.LFunction{}
	var ran rexExec
	noop := func(*lua.LState) int { return 0 }
	fn := func(f lua.LGFunction) *lua.LFunction { return L.NewFunction(f) }
	tbl := func(kv map[string]lua.LValue) *lua.LTable {
		x := L.NewTable()
		for k, v := range kv {
			x.RawSetString(k, v)
		}
		return x
	}
	L.SetField(L.GetGlobal("os"), "getenv", fn(func(L *lua.LState) int {
		if L.CheckString(1) == "SHELL" && shell != "" {
			L.Push(lua.LString(shell))
			return 1
		}
		L.Push(lua.LNil)
		return 1
	}))
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
			if L.CheckString(2) != "process" {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(tbl(map[string]lua.LValue{"child": tbl(map[string]lua.LValue{"cwd": lua.LString("/repo")})}))
			return 1
		})}),
		"session": tbl(map[string]lua.LValue{
			"new_layer": fn(func(L *lua.LState) int {
				ran.layers++
				L.Push(L.NewTable())
				return 1
			}),
		}),
		"exec": fn(func(L *lua.LState) int {
			argv := L.CheckTable(1)
			for i := 1; i <= argv.Len(); i++ {
				ran.argv = append(ran.argv, argv.RawGetInt(i).String())
			}
			if o, ok := L.Get(2).(*lua.LTable); ok {
				ran.detach = lua.LVAsBool(o.RawGetString("detach"))
				if c := o.RawGetString("cwd"); c != lua.LNil {
					ran.cwd = c.String()
				}
			}
			return 0
		}),
	})
	L.SetGlobal("rex", rex)
	if err := L.DoString(script); err != nil {
		t.Fatalf("loading rhiza.lua: %v", err)
	}
	out := map[string]rexExec{}
	for _, a := range launcher.Apps() {
		run, ok := actions[rexActionName(a.Name)]
		if !ok {
			t.Fatalf("rhiza.lua registers no action for %s", a.Name)
		}
		ran = rexExec{}
		ctx := tbl(map[string]lua.LValue{"block_id": lua.LString("b"), "session_id": lua.LString("s")})
		if err := L.CallByParam(lua.P{Fn: run, NRet: 1, Protect: true}, ctx); err != nil {
			t.Fatalf("%s: running its action: %v", a.Name, err)
		}
		L.Pop(1)
		out[a.Name] = ran
	}
	return out
}

// TestRexInstallActionsStartTheModalOpener installs rhiza.lua, runs every
// registered app's action under a fake rex, and checks each opens no layer
// itself but starts `rhiza tui open <app> --session s --block b`, detached,
// in the invoking block's directory, through the owner's login shell; or,
// with settings.bin set, runs that binary as given. The opener sizes and
// resizes the modal (internal/tui/modalopen).
func TestRexInstallActionsStartTheModalOpener(t *testing.T) {
	launcher.Register(launcher.App{Name: "zz-opened", Summary: "Opened"})
	ctx, dir, _ := rexTestCtx(t)
	if _, err := callOp(t, ctx, "rex.install", `{}`); err != nil {
		t.Fatal(err)
	}
	script := readFile(t, filepath.Join(dir, "rhiza.lua"))
	check := func(ran map[string]rexExec, prefix []string) {
		t.Helper()
		for _, a := range launcher.Apps() {
			got := ran[a.Name]
			want := append(slices.Clone(prefix), "tui", "open", a.Name, "--session", "s", "--block", "b")
			if !slices.Equal(got.argv, want) {
				t.Errorf("%s: ran %q, want %q", a.Name, got.argv, want)
			}
			if !got.detach || got.cwd != "/repo" {
				t.Errorf("%s: detach %v cwd %q, want detached in /repo", a.Name, got.detach, got.cwd)
			}
			if got.layers != 0 {
				t.Errorf("%s: the action opened %d layers itself", a.Name, got.layers)
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
}
