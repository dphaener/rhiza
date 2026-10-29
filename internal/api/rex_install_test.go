package api

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
	edited := strings.Replace(readFile(t, script), `"ctrl+a a"`, `"ctrl+a l"`, 1)
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
		`launch_key = "ctrl+a a"`,
		`local BIN = "/opt/my tools/rhiza"`,
		`name = "rhiza_launch"`,
		`options = { command = { BIN, "tui", "launch" }, cwd = block_cwd(ctx.block_id) }`,
		`rex.block.call(TERM, "process", { block_id = block_id })`,
		`rex.bind((settings.launch_key:gsub("%s+", ">")), "rhiza_launch")`,
	} {
		if !strings.Contains(script, want) {
			t.Errorf("rhiza.lua lacks %s:\n%s", want, script)
		}
	}

	// A changed launch key survives a reinstall, and the binary path is
	// refreshed.
	writePath(t, filepath.Join(dir, "rhiza.lua"), strings.Replace(script, `"ctrl+a a"`, `"ctrl+a g"`, 1))
	ctx.Self = "/elsewhere/rhiza"
	if _, err := callOp(t, ctx, "rex.install", `{}`); err != nil {
		t.Fatal(err)
	}
	script = readFile(t, filepath.Join(dir, "rhiza.lua"))
	if !strings.Contains(script, `launch_key = "ctrl+a g"`) || !strings.Contains(script, `local BIN = "/elsewhere/rhiza"`) {
		t.Errorf("reinstall lost the key or kept the old binary:\n%s", script)
	}
}
