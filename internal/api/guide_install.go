package api

// This file holds the guide.install operation: it writes the agent guide
// into Claude Code as a skills-directory plugin, a folder in
// ~/.claude/skills/ (or a repository's .claude/skills/) that Claude Code
// loads with no install step. The plugin holds the genesis skill, the
// genesis-scrub skill and a SessionStart hook that runs genesis status
// --hook.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dphaener/genesis/internal/canon"
	"github.com/dphaener/genesis/internal/guide"
)

// The files guide.install writes, relative to the plugin folder, with
// slashes, sorted.
const (
	pluginManifestPath = ".claude-plugin/plugin.json"
	pluginSkillPath    = "SKILL.md"
	pluginHooksPath    = "hooks/hooks.json"
	pluginScrubPath    = "skills/genesis-scrub/SKILL.md"
)

// pluginName is the plugin's name, and the folder's under .claude/skills.
const pluginName = "genesis"

// pluginDescription is the plugin manifest's description.
const pluginDescription = "The Genesis agent guide: the genesis and genesis-scrub skills, and a " +
	"SessionStart hook that runs genesis status --hook."

// pluginMarker is the manifest's ownership marker: a plugin.json holding
// it, under the name genesis, was written by guide.install.
const pluginMarker = "genesis"

// GuideInstallRequest asks for the guide installed for an agent, in the
// user's skills folder or, with Project, the repository's.
type GuideInstallRequest struct {
	Agent   string `json:"agent" enum:"claude" desc:"The agent to install the guide for: claude (Claude Code)"`
	Project bool   `json:"project,omitempty" desc:"Install into this repository's .claude/skills/ instead of ~/.claude/skills/"`
	Force   bool   `json:"force,omitempty" desc:"Overwrite a folder genesis guide install did not write; other files in it are kept"`
}

// GuideInstallResponse says where guide.install wrote the plugin, which
// files it wrote and which were already current. Paths are relative to
// Dir, with slashes, sorted.
type GuideInstallResponse struct {
	Dir       string   `json:"dir" desc:"The plugin folder, absolute"`
	Written   []string `json:"written" desc:"The files written: new, or different from what was on disk"`
	Unchanged []string `json:"unchanged" desc:"The files that were already current"`
}

func init() {
	register(Operation{
		Name:     "guide.install",
		Summary:  "Install the agent guide into Claude Code as a skills-directory plugin",
		Request:  reflect.TypeFor[GuideInstallRequest](),
		Response: reflect.TypeFor[GuideInstallResponse](),
		Run:      runGuideInstall,
	})
}

// runGuideInstall writes the plugin's four files. A folder that exists and
// is not ours (see pluginOwned) is a conflict unless req.Force; then
// nothing is written. Each file is written atomically, and only when its
// bytes change; no other file in the folder is touched. guide.install
// writes no log entry, so --dry-run is a usage error, as for ci.install.
func runGuideInstall(ctx *Ctx, req any) (any, error) {
	r := req.(GuideInstallRequest)
	if r.Agent != "claude" {
		return nil, Usagef("unknown agent %q: want claude", r.Agent).WithNext("genesis guide install --agent claude")
	}
	cmd := "genesis guide install --agent " + r.Agent
	if ctx.DryRun {
		return nil, Usagef("guide install does not take --dry-run: it writes no log entry").WithNext(cmd + projectFlag(r.Project))
	}
	var dir string
	if r.Project {
		if _, err := openStore(ctx.Root); err != nil {
			return nil, err
		}
		dir = filepath.Join(ctx.Root, ".claude", "skills", pluginName)
	} else {
		if ctx.Home == "" {
			return nil, Usagef("guide install: no home directory to install into").
				WithNext(cmd + " --project")
		}
		dir = filepath.Join(ctx.Home, ".claude", "skills", pluginName)
	}

	want, err := pluginFiles(ctx.Version)
	if err != nil {
		return nil, Errorf("guide install: %v", err)
	}
	if _, err := os.Lstat(dir); err == nil {
		if !r.Force && !pluginOwned(dir) {
			return nil, Conflictf("%s exists and was not written by genesis guide install", dir).
				WithNext(cmd + " --force" + projectFlag(r.Project))
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, Errorf("guide install: %v", err)
	}

	resp := GuideInstallResponse{Dir: dir, Written: []string{}, Unchanged: []string{}}
	for _, p := range []string{pluginManifestPath, pluginSkillPath, pluginHooksPath, pluginScrubPath} {
		dst := filepath.Join(dir, filepath.FromSlash(p))
		have, err := os.ReadFile(dst)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, Errorf("guide install: %v", err)
		}
		if err == nil && bytes.Equal(have, want[p]) {
			resp.Unchanged = append(resp.Unchanged, p)
			continue
		}
		if err := writeAtomic(dst, want[p]); err != nil {
			return nil, Errorf("guide install: %v", err)
		}
		resp.Written = append(resp.Written, p)
	}
	return resp, nil
}

// projectFlag is " --project" when project, else "".
func projectFlag(project bool) string {
	if project {
		return " --project"
	}
	return ""
}

// pluginFiles returns the plugin's files by path, for a genesis of the
// given version.
func pluginFiles(version string) (map[string][]byte, error) {
	manifest, err := canon.Marshal(map[string]any{
		"name":        pluginName,
		"version":     version,
		"description": pluginDescription,
		pluginMarker:  map[string]any{"installed_by": "genesis guide install"},
	})
	if err != nil {
		return nil, err
	}
	hooks, err := canon.Marshal(map[string]any{
		"hooks": map[string]any{
			"SessionStart": []any{map[string]any{
				"hooks": []any{map[string]any{
					"type":    "command",
					"command": "genesis",
					"args":    []string{"status", "--hook"},
					"timeout": 10,
				}},
			}},
		},
	})
	if err != nil {
		return nil, err
	}
	skill, err := skillFile("genesis", "genesis")
	if err != nil {
		return nil, err
	}
	scrub, err := skillFile("genesis-scrub", "scrub")
	if err != nil {
		return nil, err
	}
	return map[string][]byte{
		pluginManifestPath: manifest,
		pluginSkillPath:    skill,
		pluginHooksPath:    hooks,
		pluginScrubPath:    scrub,
	}, nil
}

// skillFile is a SKILL.md: frontmatter naming the skill and describing it
// with the guide text's description, then the text.
func skillFile(name, text string) ([]byte, error) {
	t, ok := guide.Lookup(text)
	if !ok {
		return nil, fmt.Errorf("no guide text %q", text)
	}
	return fmt.Appendf(nil, "---\nname: %s\ndescription: %s\n---\n\n%s", name, yamlQuote(t.Description), t.Body), nil
}

// yamlQuote returns s as a YAML double-quoted scalar, which is valid
// whatever s holds: backslash, double quote and control characters are
// escaped.
func yamlQuote(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, c := range s {
		switch {
		case c == '"' || c == '\\':
			b.WriteByte('\\')
			b.WriteRune(c)
		case c == '\n':
			b.WriteString(`\n`)
		case c == '\t':
			b.WriteString(`\t`)
		case c < 0x20 || c == 0x7f:
			fmt.Fprintf(&b, `\x%02x`, c)
		default:
			b.WriteRune(c)
		}
	}
	b.WriteByte('"')
	return b.String()
}

// pluginManifest is the part of a plugin.json that says whose it is.
type pluginManifest struct {
	name    string
	version string
	marked  bool
}

// readPluginManifest reads dir's plugin.json; ok is false when it is
// missing or does not parse as a JSON object.
func readPluginManifest(dir string) (m pluginManifest, ok bool) {
	b, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(pluginManifestPath)))
	if err != nil {
		return m, false
	}
	var raw map[string]json.RawMessage
	if json.Unmarshal(b, &raw) != nil {
		return m, false
	}
	json.Unmarshal(raw["name"], &m.name)
	json.Unmarshal(raw["version"], &m.version)
	var marker map[string]json.RawMessage
	m.marked = json.Unmarshal(raw[pluginMarker], &marker) == nil && marker != nil
	return m, true
}

// pluginOwned reports whether dir holds a plugin guide.install wrote: a
// plugin.json named genesis that carries the marker.
func pluginOwned(dir string) bool {
	m, ok := readPluginManifest(dir)
	return ok && m.name == pluginName && m.marked
}

// installedVersion is the genesis version that installed the plugin in
// dir; ok is false unless dir holds a plugin guide.install wrote.
func installedVersion(dir string) (version string, ok bool) {
	if !pluginOwned(dir) {
		return "", false
	}
	m, _ := readPluginManifest(dir)
	return m.version, true
}

// Human reports the folder and the counts, then each path written.
func (r GuideInstallResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "Installed the Genesis guide in %s (%d written, %d unchanged).\n", r.Dir, len(r.Written), len(r.Unchanged))
	if len(r.Written) == 0 {
		return
	}
	fmt.Fprintln(w)
	for _, p := range r.Written {
		fmt.Fprintf(w, "- wrote %s\n", p)
	}
}
