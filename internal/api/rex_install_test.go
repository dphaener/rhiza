package api

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/tui/launcher"
)

// rexTestCtx is a ctx with a fake Rex config folder and a fake rex that
// records its argv.
func rexTestCtx(t *testing.T) (ctx *Ctx, dir string, calls *[][]string) {
	t.Helper()
	ctx = testCtx(t.TempDir())
	dir = t.TempDir()
	ctx.Getenv = func(k string) string {
		if k == "REX_CONFIG" {
			return filepath.Join(dir, "init.lua")
		}
		return ""
	}
	calls = &[][]string{}
	ctx.Exec = func(name string, args []string, _ string, _ []string) ([]byte, error) {
		*calls = append(*calls, append([]string{name}, args...))
		return nil, nil
	}
	return ctx, dir, calls
}

func readFile(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestRexInstallRoundTrip(t *testing.T) {
	for _, owner := range []string{
		"-- mine\nrex.bind(\"ctrl+a r\", \"client.config.reload\")\n",
		"-- no trailing newline\nrex.bind(\"x\", \"y\")",
		"",
		"require(\"forge\")\n\n\n  -- odd   spacing\t\n",
	} {
		ctx, dir, calls := rexTestCtx(t)
		initPath := filepath.Join(dir, "init.lua")
		writePath(t, initPath, owner)

		if _, err := callOp(t, ctx, "rex.install", `{}`); err != nil {
			t.Fatal(err)
		}
		afterFirst := readFile(t, initPath)
		script := readFile(t, filepath.Join(dir, "rhiza.lua"))
		if !strings.HasPrefix(afterFirst, owner) || strings.Count(afterFirst, `require("rhiza")`) != 1 {
			t.Fatalf("owner %q: init.lua = %q", owner, afterFirst)
		}

		resp, err := callOp(t, ctx, "rex.install", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		if got := resp.(RexInstallResponse); len(got.Written) != 0 || len(got.Unchanged) != 2 {
			t.Errorf("second install = %+v", got)
		}
		if readFile(t, initPath) != afterFirst || readFile(t, filepath.Join(dir, "rhiza.lua")) != script {
			t.Errorf("owner %q: second install changed a file", owner)
		}
		if len(*calls) != 2 || strings.Join((*calls)[0], " ") != "rex config reload" {
			t.Errorf("calls = %v", *calls)
		}

		if _, err := callOp(t, ctx, "rex.uninstall", `{}`); err != nil {
			t.Fatal(err)
		}
		if owner == "" {
			if _, err := os.Stat(initPath); !errors.Is(err, os.ErrNotExist) {
				t.Errorf("empty owner: init.lua should be gone, err %v", err)
			}
		} else if got := readFile(t, initPath); got != owner {
			t.Errorf("uninstall left %q, want %q", got, owner)
		}
		if _, err := os.Stat(filepath.Join(dir, "rhiza.lua")); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("rhiza.lua remains: %v", err)
		}
		if len(*calls) != 3 {
			t.Errorf("uninstall did not reload: %v", *calls)
		}
	}
}

func TestRexInstallKeepsSettingsAndOwnersRequire(t *testing.T) {
	ctx, dir, _ := rexTestCtx(t)
	owner := "require('rhiza')\n"
	writePath(t, filepath.Join(dir, "init.lua"), owner)
	if _, err := callOp(t, ctx, "rex.install", `{}`); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, filepath.Join(dir, "init.lua")); got != owner {
		t.Errorf("owner's own require was changed: %q", got)
	}
	script := filepath.Join(dir, "rhiza.lua")
	edited := strings.Replace(readFile(t, script), `"ctrl+a>a"`, `"ctrl+a l"`, 1)
	writePath(t, script, edited)
	if _, err := callOp(t, ctx, "rex.install", `{}`); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, script); got != edited {
		t.Errorf("reinstall lost the setting:\n%s", got)
	}
	// Uninstall leaves the owner's require line alone.
	if _, err := callOp(t, ctx, "rex.uninstall", `{}`); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, filepath.Join(dir, "init.lua")); got != owner {
		t.Errorf("uninstall touched the owner's require: %q", got)
	}
}

func TestRexInstallForeignScript(t *testing.T) {
	ctx, dir, _ := rexTestCtx(t)
	writePath(t, filepath.Join(dir, "rhiza.lua"), "-- not ours\n")
	if _, err := callOp(t, ctx, "rex.install", `{}`); err == nil {
		t.Fatal("want a conflict")
	}
	if got := readFile(t, filepath.Join(dir, "rhiza.lua")); got != "-- not ours\n" {
		t.Errorf("foreign file overwritten: %q", got)
	}
	if _, err := callOp(t, ctx, "rex.install", `{"force":true}`); err != nil {
		t.Fatal(err)
	}
}

func TestRexInstallThroughSymlink(t *testing.T) {
	ctx, dir, _ := rexTestCtx(t)
	real := filepath.Join(t.TempDir(), "dotfiles-init.lua")
	writePath(t, real, "-- real\n")
	if err := os.Symlink(real, filepath.Join(dir, "init.lua")); err != nil {
		t.Fatal(err)
	}
	if _, err := callOp(t, ctx, "rex.install", `{}`); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Lstat(filepath.Join(dir, "init.lua")); err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Errorf("symlink was replaced: %v %v", info, err)
	}
	if !strings.Contains(readFile(t, real), `require("rhiza")`) {
		t.Error("require not written through the link")
	}
}

func TestRexInstallScriptOpensTheLauncher(t *testing.T) {
	ctx, dir, _ := rexTestCtx(t)
	ctx.Self = "/opt/my tools/rhiza"
	if _, err := callOp(t, ctx, "rex.install", `{}`); err != nil {
		t.Fatal(err)
	}
	script := readFile(t, filepath.Join(dir, "rhiza.lua"))
	for _, want := range []string{
		`launch_key = "ctrl+a>a"`,
		`local BIN = "/opt/my tools/rhiza"`,
		`name = "rhiza_launch"`,
		`options = { command = { BIN, "tui", "launch" }, cwd = block_cwd(ctx.block_id) }`,
		`rex.block.call(TERM, "process", { block_id = block_id })`,
		`rex.bind((settings.launch_key:gsub("%s+", ">")), "client.mode.enter", { name = MODE, once = true })`,
	} {
		if !strings.Contains(script, want) {
			t.Errorf("rhiza.lua lacks %s:\n%s", want, script)
		}
	}

	// A changed launch key survives a reinstall, and the binary path is
	// refreshed.
	writePath(t, filepath.Join(dir, "rhiza.lua"), strings.Replace(script, `"ctrl+a>a"`, `"ctrl+a g"`, 1))
	ctx.Self = "/elsewhere/rhiza"
	if _, err := callOp(t, ctx, "rex.install", `{}`); err != nil {
		t.Fatal(err)
	}
	script = readFile(t, filepath.Join(dir, "rhiza.lua"))
	if !strings.Contains(script, `launch_key = "ctrl+a g"`) || !strings.Contains(script, `local BIN = "/elsewhere/rhiza"`) {
		t.Errorf("reinstall lost the key or kept the old binary:\n%s", script)
	}
}

func TestRexInstallRegistersAnActionPerApp(t *testing.T) {
	launcher.Register(launcher.App{Name: "zz-extra", Summary: "One more app"})
	ctx, dir, _ := rexTestCtx(t)
	if _, err := callOp(t, ctx, "rex.install", `{}`); err != nil {
		t.Fatal(err)
	}
	script := readFile(t, filepath.Join(dir, "rhiza.lua"))
	for _, want := range []string{
		`name = "rhiza_app_zz_extra"`,
		`title = "zz-extra"`,
		`command = { BIN, "tui", "app", "zz-extra" }`,
	} {
		if !strings.Contains(script, want) {
			t.Errorf("rhiza.lua lacks %s:\n%s", want, script)
		}
	}
	if got, want := strings.Count(script, "rex.action{"), 1+len(launcher.Apps()); got != want {
		t.Errorf("%d rex.action calls, want %d", got, want)
	}

	// An app registered later gets its action on a reinstall.
	launcher.Register(launcher.App{Name: "zz-later", Summary: "Added after the first install"})
	if _, err := callOp(t, ctx, "rex.install", `{}`); err != nil {
		t.Fatal(err)
	}
	if script = readFile(t, filepath.Join(dir, "rhiza.lua")); !strings.Contains(script, `command = { BIN, "tui", "app", "zz-later" }`) {
		t.Errorf("reinstall lacks the new app's action:\n%s", script)
	}
}

func TestRexInstallSizesEachAppLayerToItsCells(t *testing.T) {
	launcher.Register(launcher.App{Name: "zz-sized", Summary: "Sized", Cols: 77, Rows: 21})
	launcher.Register(launcher.App{Name: "zz-default", Summary: "Default size"})
	ctx, dir, _ := rexTestCtx(t)
	if _, err := callOp(t, ctx, "rex.install", `{}`); err != nil {
		t.Fatal(err)
	}
	script := readFile(t, filepath.Join(dir, "rhiza.lua"))
	action := func(name string) string {
		i := strings.Index(script, `name = "`+rexActionName(name)+`"`)
		if i < 0 {
			t.Fatalf("no action for %s:\n%s", name, script)
		}
		end := strings.Index(script[i:], "\n}\n")
		return script[i : i+end]
	}
	sized := action("zz-sized")
	if !strings.Contains(sized, "bounds = cell_bounds(ctx.block_id, 77, 21)") {
		t.Errorf("zz-sized is not sized from its own cells:\n%s", sized)
	}
	if !strings.Contains(sized, `command = { BIN, "tui", "app", "zz-sized" }`) || strings.Contains(sized, `"launch"`) {
		t.Errorf("zz-sized does not run its app directly:\n%s", sized)
	}
	dc, dr := launcher.App{}.Size()
	if want := fmt.Sprintf("cell_bounds(ctx.block_id, %d, %d)", dc, dr); !strings.Contains(action("zz-default"), want) {
		t.Errorf("zz-default lacks %s", want)
	}
	if strings.Contains(strings.Replace(script, `{ BIN, "tui", "launch" }`, "", 1), `"launch"`) {
		t.Errorf("an app action goes through tui launch:\n%s", script)
	}
}

func TestRexInstallBindsAKeyPerAppInTheKeyMode(t *testing.T) {
	ctx, dir, _ := rexTestCtx(t)
	if _, err := callOp(t, ctx, "rex.install", `{}`); err != nil {
		t.Fatal(err)
	}
	script := readFile(t, filepath.Join(dir, "rhiza.lua"))
	keys := rexDefaultKeys(launcher.Apps())
	if len(keys) != len(launcher.Apps()) {
		t.Fatalf("%d keys for %d apps", len(keys), len(launcher.Apps()))
	}
	seen := map[string]string{}
	for app, key := range keys {
		if other, dup := seen[key]; dup {
			t.Errorf("apps %s and %s share key %s", app, other, key)
		}
		seen[key] = app
		if want := fmt.Sprintf("[%q] = %q,", app, key); !strings.Contains(script, want) {
			t.Errorf("rhiza.lua lacks the default %s:\n%s", want, script)
		}
		action := rexActionName(app)
		if want := fmt.Sprintf("[%q] = %q,", app, action); !strings.Contains(script, want) {
			t.Errorf("rhiza.lua lacks the action %s:\n%s", want, script)
		}
		if !strings.Contains(script, fmt.Sprintf("name = %q", action)) {
			t.Errorf("key of %s binds %s, which no rex.action defines", app, action)
		}
	}
	for _, want := range []string{
		`rex.mode(MODE, { exclusive = true })`,
		`rex.bind(MODE .. "/" .. k, APP_ACTIONS[app])`,
		`keys = {},`,
	} {
		if !strings.Contains(script, want) {
			t.Errorf("rhiza.lua lacks %s", want)
		}
	}
}

func TestRexDefaultKeysAreMnemonicAndUnique(t *testing.T) {
	got := rexDefaultKeys([]launcher.App{{Name: "model"}, {Name: "mine"}, {Name: "roadmap"}, {Name: "x"}, {Name: "x2"}})
	want := map[string]string{"model": "m", "mine": "i", "roadmap": "r", "x": "x", "x2": "2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("keys = %v, want %v", got, want)
	}
}

func TestRexInstallKeepsKeyOverridesOnReinstall(t *testing.T) {
	ctx, dir, _ := rexTestCtx(t)
	if _, err := callOp(t, ctx, "rex.install", `{}`); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(dir, "rhiza.lua")
	writePath(t, p, strings.Replace(readFile(t, p), "keys = {},", `keys = { roadmap = "z" },`, 1))
	if _, err := callOp(t, ctx, "rex.install", `{}`); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(readFile(t, p), `keys = { roadmap = "z" },`) {
		t.Error("reinstall lost the key override")
	}
}
