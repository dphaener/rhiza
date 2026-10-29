package api

// This file holds the rex.install and rex.uninstall operations: they put
// the quick-launch script into the owner's Rex config as its own file,
// rhiza.lua, beside Rex's init.lua (as forge keeps forge.lua), and a single
// require("rhiza") line in init.lua that loads it.

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/dphaener/rexuikit/surface"
	"github.com/dphaener/rhiza/internal/tui/launcher"
)

const (
	rexScriptName = "rhiza.lua"
	rexInitName   = "init.lua"
	// rexRequireLine is the line install adds to init.lua. The comment
	// marks it as ours: uninstall removes only a line that matches it, so
	// a require("rhiza") the owner wrote is theirs and stays.
	rexRequireLine = `require("rhiza") -- added by rhiza rex install`
	// rexScriptMarker is the first line of rhiza.lua; a file without it
	// was not written by install.
	rexScriptMarker = "-- Written by rhiza rex install."
	// The settings block is the part of rhiza.lua a reinstall keeps.
	rexSettingsBegin = "-- settings begin\n"
	rexSettingsEnd   = "-- settings end\n"
	// rexDefaultSettings is the settings block of a first install. The
	// launch key is ctrl+a>a (see decision tui-launch-key); it enters the
	// rhiza key mode. keys remaps an app's key in that mode.
	rexDefaultSettings = rexSettingsBegin +
		"local settings = {\n  launch_key = \"ctrl+a>a\",\n" +
		"  -- Remap an app's key in the rhiza key mode, e.g. keys = { roadmap = \"m\" }.\n" +
		"  keys = {},\n}\n" +
		rexSettingsEnd
)

// rexLaunchLua is the part of rhiza.lua after the settings and BIN: the
// action that opens the quick-launch layer (in the palette only; the launch
// key enters the key mode, see rexKeyModeLua). The layer runs
// `rhiza tui launch`, which lists the apps and opens the chosen one in a
// modal. The layer starts in the invoking block's directory, since the
// launcher hands its own to each app; without one Rex starts it in $HOME.
const rexLaunchLua = `local TERM = "com.superlogical.terminal"
local LAYER_BOUNDS = { 0.3, 0.25, 0.4, 0.5 }

-- block_cwd returns the working directory of the terminal block's
-- foreground program, else of its shell, else nil.
local function block_cwd(block_id)
  if block_id == nil then return nil end
  local p = rex.block.call(TERM, "process", { block_id = block_id })
  if p == nil then return nil end
  return (p.foreground and p.foreground.cwd) or (p.child and p.child.cwd)
end

rex.action{
  name = "rhiza_launch",
  title = "Rhiza apps",
  description = "Open the Rhiza quick-launch layer",
  category = "rhiza",
  keywords = "rhiza apps launch",
  run = function(ctx)
    local res, err = rex.session.new_layer{
      bounds = LAYER_BOUNDS,
      layout = rex.layout.block{
        flavor = TERM .. ".shell",
        label = "Rhiza apps",
        options = { command = { BIN, "tui", "launch" }, cwd = block_cwd(ctx.block_id) },
      },
      focus = true,
    }
    if err ~= nil then error(err, 0) end
    return res
  end,
}

`

// rexKeyModeLua is the rhiza key mode: the launch key enters it, and each
// app's key in it runs that app's action. DEFAULT_KEYS (written by
// rexKeyModeTable) gives each app a default key; settings.keys overrides
// any of them, and a reinstall keeps settings.keys. An exclusive mode drops
// every key it does not bind, so a mistyped key never reaches the terminal.
const rexKeyModeLua = `
local MODE = "rhiza"
rex.mode(MODE, { exclusive = true })
rex.bind(MODE .. "/escape", "client.mode.exit")
local overrides = type(settings.keys) == "table" and settings.keys or {}
for app, key in pairs(DEFAULT_KEYS) do
  local k = overrides[app]
  if type(k) ~= "string" or k == "" then k = key end
  rex.bind(MODE .. "/" .. k, "rhiza_app_" .. app)
end
if type(settings.launch_key) == "string" and settings.launch_key ~= "" then
  rex.bind((settings.launch_key:gsub("%s+", ">")), "client.mode.enter", { name = MODE, once = true })
end
`

// rexDefaultKeys gives each app a default key of the rhiza key mode, in
// the order of apps: the first letter of its name that no earlier app took
// (word initials first), else the first free letter or digit. The same apps always give the same keys.
func rexDefaultKeys(apps []launcher.App) map[string]string {
	used := map[string]bool{}
	keys := map[string]string{}
	for _, a := range apps {
		pick := ""
		name := strings.ToLower(a.Name)
		for _, r := range rexInitials(name) + name + "abcdefghijklmnopqrstuvwxyz0123456789" {
			k := string(r)
			if (r >= 'a' && r <= 'z' || r >= '0' && r <= '9') && !used[k] {
				pick = k
				break
			}
		}
		if pick == "" {
			continue
		}
		used[pick] = true
		keys[a.Name] = pick
	}
	return keys
}

// rexInitials returns the first letter of each hyphen-separated word of name.
func rexInitials(name string) string {
	var b strings.Builder
	for _, w := range strings.Split(name, "-") {
		if w != "" {
			b.WriteByte(w[0])
		}
	}
	return b.String()
}

// rexKeyMode returns the Lua of the key mode: the table of default keys,
// then rexKeyModeLua. Apps without a free key get no default, only the
// action in the palette.
func rexKeyMode() string {
	apps := launcher.Apps()
	keys := rexDefaultKeys(apps)
	var b strings.Builder
	b.WriteString("\nlocal DEFAULT_KEYS = {\n")
	for _, a := range apps {
		if k, ok := keys[a.Name]; ok {
			fmt.Fprintf(&b, "  [%s] = %s,\n", strconv.Quote(a.Name), strconv.Quote(k))
		}
	}
	b.WriteString("}\n")
	return b.String() + rexKeyModeLua
}

// rexCellBoundsLua is the Lua that sizes a layer to a number of cells: the
// chrome model of rexuikit's surface.BoundsForCells, centred, with its
// constants (frame, border and margin in pixels, cols then rows) filled in
// by rexCellBounds. The window is taken to be the invoking block's
// terminal; without a size to work from it gives LAYER_BOUNDS.
const rexCellBoundsLua = `
-- cell_bounds returns bounds that give a program cols x rows cells.
local function cell_bounds(block_id, cols, rows)
  if block_id == nil then return LAYER_BOUNDS end
  local ok, sz = pcall(rex.block.call, TERM, "size", { block_id = block_id })
  if not ok or type(sz) ~= "table" or (sz.cell_width_px or 0) <= 0 or (sz.cell_height_px or 0) <= 0 then
    return LAYER_BOUNDS
  end
  local function edge(v) return math.floor(v + 0.5) end
  local function fit_axis(win, frame, cell, border, margin, want)
    local px = win * cell + frame
    local span = math.min(1, (want * cell + border + margin) / px)
    while true do
      local pos = math.max(0, (1 - span) / 2)
      local got = math.floor((edge(px * (pos + span)) - edge(px * pos) - border - margin) / cell)
      if got >= want or span >= 1 then return pos, span end
      span = math.min(1, span + 1 / px)
    end
  end
  local x, w = fit_axis(sz.columns, %[1]d, sz.cell_width_px, %[2]d, %[3]d, cols)
  local y, h = fit_axis(sz.rows, %[4]d, sz.cell_height_px, %[5]d, %[6]d, rows)
  return { x, y, w, h }
end
`

func rexCellBounds() string {
	return fmt.Sprintf(rexCellBoundsLua,
		surface.FramePxCols, surface.BorderPxCols, surface.InverseMarginPxCols,
		surface.FramePxRows, surface.BorderPxRows, surface.InverseMarginPxRows)
}

// rexAppActionLua is the Lua of one app's action: its own entry in Rex's
// command palette, titled with the app's name, that opens the app straight
// away, in a layer sized to the app's own cells, running
// `rhiza tui app <name>`; no list step comes first. The arguments are Lua
// string literals (action name, title, description, keywords, app name),
// quoted by rexAppActions, then the app's columns and rows.
const rexAppActionLua = `
rex.action{
  name = %[1]s,
  title = %[2]s,
  description = %[3]s,
  category = "rhiza",
  keywords = %[4]s,
  run = function(ctx)
    local res, err = rex.session.new_layer{
      bounds = cell_bounds(ctx.block_id, %[6]d, %[7]d),
      layout = rex.layout.block{
        flavor = TERM .. ".shell",
        label = %[2]s,
        options = { command = { BIN, "tui", "app", %[5]s }, cwd = block_cwd(ctx.block_id) },
      },
      focus = true,
    }
    if err ~= nil then error(err, 0) end
    return res
  end,
}
`

// rexAppActions returns the Lua registering one rex.action per registered
// app, in the order launcher.Apps gives (by name), so the same apps always
// give the same bytes.
func rexAppActions() string {
	var b strings.Builder
	for _, a := range launcher.Apps() {
		desc := a.Summary
		if desc == "" {
			desc = "Open the " + a.Name + " app"
		}
		cols, rows := a.Size()
		fmt.Fprintf(&b, rexAppActionLua,
			strconv.Quote("rhiza_app_"+a.Name), strconv.Quote(a.Name), strconv.Quote(desc),
			strconv.Quote("rhiza app "+a.Name), strconv.Quote(a.Name), cols, rows)
	}
	return b.String()
}

const rexScriptTail = `
return settings
`

// RexInstallRequest asks for the quick-launch script installed in Rex.
type RexInstallRequest struct {
	Force bool `json:"force,omitempty" desc:"Overwrite a rhiza.lua rhiza rex install did not write"`
}

// RexInstallResponse says where the files are and what changed.
type RexInstallResponse struct {
	Script    string   `json:"script" desc:"The rhiza.lua file, absolute"`
	Init      string   `json:"init" desc:"Rex's init.lua, absolute"`
	Written   []string `json:"written" desc:"The files written: rhiza.lua and init.lua, those that changed"`
	Unchanged []string `json:"unchanged" desc:"The files that were already current"`
}

// RexUninstallRequest asks for the quick-launch script removed from Rex.
type RexUninstallRequest struct{}

// RexUninstallResponse says what uninstall removed.
type RexUninstallResponse struct {
	Removed []string `json:"removed" desc:"The files changed or removed: rhiza.lua and init.lua, those that held our part"`
}

func init() {
	register(Operation{
		Name:     "rex.install",
		Summary:  "Install the quick-launch script into the Rex config and reload Rex",
		Request:  reflect.TypeFor[RexInstallRequest](),
		Response: reflect.TypeFor[RexInstallResponse](),
		Run:      runRexInstall,
	})
	register(Operation{
		Name:     "rex.uninstall",
		Summary:  "Remove the quick-launch script from the Rex config and reload Rex",
		Request:  reflect.TypeFor[RexUninstallRequest](),
		Response: reflect.TypeFor[RexUninstallResponse](),
		Run:      runRexUninstall,
	})
}

// rexConfigPaths returns the paths of Rex's init.lua and of rhiza.lua
// beside it. Rex reads REX_CONFIG when set, else rex/init.lua in
// $XDG_CONFIG_HOME or ~/.config. init.lua may be a symlink (dotfiles); it is
// resolved so a write goes through it, but rhiza.lua stays in the folder
// Rex reads, where require("rhiza") looks for it.
func rexConfigPaths(ctx *Ctx, cmd string) (initPath, scriptPath string, err error) {
	get := ctx.Getenv
	if get == nil {
		get = func(string) string { return "" }
	}
	switch {
	case get("REX_CONFIG") != "":
		initPath = get("REX_CONFIG")
	case get("XDG_CONFIG_HOME") != "":
		initPath = filepath.Join(get("XDG_CONFIG_HOME"), "rex", rexInitName)
	case ctx.Home != "":
		initPath = filepath.Join(ctx.Home, ".config", "rex", rexInitName)
	default:
		return "", "", Usagef("%s: no Rex config folder: set REX_CONFIG", cmd)
	}
	scriptPath = filepath.Join(filepath.Dir(initPath), rexScriptName)
	if r, rerr := filepath.EvalSymlinks(initPath); rerr == nil {
		initPath = r
	}
	return initPath, scriptPath, nil
}

// rexReload makes Rex read its config again.
func rexReload(ctx *Ctx, cmd string) error {
	if ctx.Exec == nil {
		return Errorf("%s: cannot run rex to reload it", cmd)
	}
	if _, err := ctx.Exec("rex", []string{"config", "reload"}, "", nil); err != nil {
		return Errorf("%s: files are in place, but rex config reload failed: %v", cmd, err).WithNext("rex config reload")
	}
	return nil
}

func readIfExists(p string) ([]byte, bool, error) {
	b, err := os.ReadFile(p)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, false, nil
	}
	return b, err == nil, err
}

// rexScript returns rhiza.lua's bytes, keeping the settings block of have
// (the file now on disk, nil if none). bin is the absolute path of rhiza,
// run by the layer the launch key opens.
func rexScript(have []byte, bin string) []byte {
	settings := rexDefaultSettings
	s := string(have)
	if i := strings.Index(s, rexSettingsBegin); i >= 0 {
		if j := strings.Index(s[i:], rexSettingsEnd); j >= 0 {
			settings = s[i : i+j+len(rexSettingsEnd)]
		}
	}
	return []byte(rexScriptMarker + "\n" +
		"-- Edit the settings block; a reinstall keeps it. Everything else is rewritten.\n" +
		settings +
		"\nlocal BIN = " + strconv.Quote(bin) + "\n" +
		rexLaunchLua + rexCellBounds() + rexAppActions() + rexKeyMode() + rexScriptTail)
}

// rexHasRequire reports whether init.lua already loads rhiza: any line that
// is a require of "rhiza", ours or the owner's.
func rexHasRequire(init []byte) bool {
	for _, l := range strings.Split(string(init), "\n") {
		l = strings.TrimSpace(l)
		for _, q := range []string{`"`, `'`} {
			if strings.HasPrefix(l, "require("+q+"rhiza"+q+")") || strings.HasPrefix(l, "require "+q+"rhiza"+q) {
				return true
			}
		}
	}
	return false
}

// rexAddRequire appends our require line, touching no existing byte. A file
// without a final newline gets one before the line, and none after, so
// rexRemoveRequire can tell the two shapes apart and undo either exactly.
func rexAddRequire(init []byte) []byte {
	switch {
	case len(init) == 0 || init[len(init)-1] == '\n':
		return append(init[:len(init):len(init)], rexRequireLine+"\n"...)
	default:
		return append(init[:len(init):len(init)], "\n"+rexRequireLine...)
	}
}

// rexRemoveRequire undoes rexAddRequire; changed is false if our line is
// not there.
func rexRemoveRequire(init []byte) (out []byte, changed bool) {
	if bytes.HasSuffix(init, []byte("\n"+rexRequireLine)) {
		return init[:len(init)-len(rexRequireLine)-1], true
	}
	line := []byte(rexRequireLine + "\n")
	if bytes.HasPrefix(init, line) {
		return init[len(line):], true
	}
	if i := bytes.LastIndex(init, append([]byte("\n"), line...)); i >= 0 {
		return append(init[:i+1:i+1], init[i+1+len(line):]...), true
	}
	return init, false
}

// runRexInstall writes rhiza.lua (keeping its settings), adds the require
// line to init.lua if no line loads rhiza yet, and reloads Rex. A run that
// finds both current writes nothing but still reloads. A rhiza.lua install
// did not write is a conflict unless Force. Like guide.install it writes no
// log entry, so --dry-run is a usage error.
func runRexInstall(ctx *Ctx, req any) (any, error) {
	r := req.(RexInstallRequest)
	const cmd = "rhiza rex install"
	if ctx.DryRun {
		return nil, Usagef("rex install does not take --dry-run: it writes no log entry").WithNext(cmd)
	}
	initPath, scriptPath, err := rexConfigPaths(ctx, cmd)
	if err != nil {
		return nil, err
	}
	haveScript, scriptExists, err := readIfExists(scriptPath)
	if err != nil {
		return nil, Errorf("rex install: %v", err)
	}
	if scriptExists && !bytes.HasPrefix(haveScript, []byte(rexScriptMarker)) && !r.Force {
		return nil, Conflictf("%s exists and was not written by rhiza rex install", scriptPath).WithNext(cmd + " --force")
	}
	initBytes, _, err := readIfExists(initPath)
	if err != nil {
		return nil, Errorf("rex install: %v", err)
	}

	resp := RexInstallResponse{Script: scriptPath, Init: initPath, Written: []string{}, Unchanged: []string{}}
	bin := ctx.Self
	if bin == "" {
		bin = "rhiza"
	}
	wantScript := rexScript(haveScript, bin)
	if scriptExists && bytes.Equal(haveScript, wantScript) {
		resp.Unchanged = append(resp.Unchanged, rexScriptName)
	} else {
		if err := writeAtomic(scriptPath, wantScript); err != nil {
			return nil, Errorf("rex install: %v", err)
		}
		resp.Written = append(resp.Written, rexScriptName)
	}
	if rexHasRequire(initBytes) {
		resp.Unchanged = append(resp.Unchanged, rexInitName)
	} else {
		if err := writeKeepingMode(initPath, rexAddRequire(initBytes)); err != nil {
			return nil, Errorf("rex install: %v", err)
		}
		resp.Written = append(resp.Written, rexInitName)
	}
	if err := rexReload(ctx, cmd); err != nil {
		return nil, err
	}
	return resp, nil
}

// runRexUninstall removes rhiza.lua (if install wrote it) and the require
// line install added, leaving every other byte of init.lua as it was, and
// reloads Rex. An init.lua install created and that is now empty is removed.
func runRexUninstall(ctx *Ctx, req any) (any, error) {
	const cmd = "rhiza rex uninstall"
	if ctx.DryRun {
		return nil, Usagef("rex uninstall does not take --dry-run: it writes no log entry").WithNext(cmd)
	}
	initPath, scriptPath, err := rexConfigPaths(ctx, cmd)
	if err != nil {
		return nil, err
	}
	resp := RexUninstallResponse{Removed: []string{}}
	haveScript, scriptExists, err := readIfExists(scriptPath)
	if err != nil {
		return nil, Errorf("rex uninstall: %v", err)
	}
	if scriptExists {
		if !bytes.HasPrefix(haveScript, []byte(rexScriptMarker)) {
			return nil, Conflictf("%s was not written by rhiza rex install: remove it yourself", scriptPath)
		}
		if err := os.Remove(scriptPath); err != nil {
			return nil, Errorf("rex uninstall: %v", err)
		}
		resp.Removed = append(resp.Removed, rexScriptName)
	}
	initBytes, initExists, err := readIfExists(initPath)
	if err != nil {
		return nil, Errorf("rex uninstall: %v", err)
	}
	if initExists {
		if out, changed := rexRemoveRequire(initBytes); changed {
			if len(out) == 0 {
				err = os.Remove(initPath)
			} else {
				err = writeKeepingMode(initPath, out)
			}
			if err != nil {
				return nil, Errorf("rex uninstall: %v", err)
			}
			resp.Removed = append(resp.Removed, rexInitName)
		}
	}
	if err := rexReload(ctx, cmd); err != nil {
		return nil, err
	}
	return resp, nil
}

// writeKeepingMode writes data to p atomically, keeping p's permissions if
// it exists.
func writeKeepingMode(p string, data []byte) error {
	info, statErr := os.Stat(p)
	if err := writeAtomic(p, data); err != nil {
		return err
	}
	if statErr == nil {
		return os.Chmod(p, info.Mode().Perm())
	}
	return nil
}

// Human reports the files, then each path written.
func (r RexInstallResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "Installed the Rhiza quick-launch script in %s (%d written, %d unchanged); Rex reloaded.\n", r.Script, len(r.Written), len(r.Unchanged))
	for _, p := range r.Written {
		fmt.Fprintf(w, "- wrote %s\n", p)
	}
}

// Human reports what was removed.
func (r RexUninstallResponse) Human(w io.Writer) {
	if len(r.Removed) == 0 {
		fmt.Fprintln(w, "Nothing to remove; Rex reloaded.")
		return
	}
	fmt.Fprintf(w, "Removed the Rhiza quick-launch script from Rex; Rex reloaded.\n")
	for _, p := range r.Removed {
		fmt.Fprintf(w, "- changed %s\n", p)
	}
}
