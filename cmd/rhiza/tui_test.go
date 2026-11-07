package main

import (
	"os/exec"

	"github.com/charmbracelet/colorprofile"
	"github.com/dphaener/rexuikit/ui"
	"github.com/dphaener/rexuikit/ui/theme"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/rootcfg"
	"github.com/dphaener/rhiza/internal/tui/autonomyapp"
	"github.com/dphaener/rhiza/internal/tui/daemonapp"
	"github.com/dphaener/rhiza/internal/tui/inboxapp"
	"github.com/dphaener/rhiza/internal/tui/requestapp"
	"github.com/dphaener/rhiza/internal/tui/runsapp"
)

// TestTUIAppHeaders opens each app, as `rhiza tui app <name>` does, in a
// project's folder, in its sibling worktree and at the root, and checks
// the header names the project, the project and no project, and fits the
// window with the app under it.
func TestTUIAppHeaders(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := t.TempDir()
	if _, err := rootcfg.Init(root); err != nil {
		t.Fatal(err)
	}
	if err := rootcfg.Add(root, "alpha", rootcfg.Project{Path: "alpha"}); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"alpha", "alpha.feature"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	places := []struct{ name, dir, want string }{
		{"project folder", "alpha", "Project: alpha"},
		{"sibling worktree", "alpha.feature", "Project: alpha"},
		{"root", ".", "No project"},
	}
	apps := []string{inboxapp.Name, runsapp.Name, autonomyapp.Name, daemonapp.Name, requestapp.Name}
	getenv := func(string) string { return "" }
	size := tea.WindowSizeMsg{Width: 80, Height: 30}
	for _, p := range places {
		for _, app := range apps {
			t.Run(p.name+"/"+app, func(t *testing.T) {
				t.Chdir(filepath.Join(root, p.dir))
				s, err := appSwitcher(app, newTUIModels(globals{}, getenv))
				if err != nil {
					t.Fatalf("open: %v", err)
				}
				if got := s.Header(); got != p.want {
					t.Errorf("header %q, want %q", got, p.want)
				}
				m, _ := s.Update(size)
				view := m.View().Content
				if first, _, _ := strings.Cut(view, "\n"); first != p.want {
					t.Errorf("first line %q, want %q", first, p.want)
				}
				if n := strings.Count(view, "\n") + 1; n > size.Height {
					t.Errorf("view is %d lines in a %d-line window", n, size.Height)
				}
			})
		}
	}
}

// TestTUISettingsScopes opens the Autonomy and Daemon apps in one of two
// projects: a toggled autonomy stage lands in that project's git config
// only, and a new run limit lands in the root's rhiza.toml.
func TestTUISettingsScopes(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(t.TempDir(), "gitconfig"))
	root := t.TempDir()
	if _, err := rootcfg.Init(root); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"alpha", "beta"} {
		if err := rootcfg.Add(root, p, rootcfg.Project{Path: p}); err != nil {
			t.Fatal(err)
		}
		dir := filepath.Join(root, p)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if out, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
			t.Fatalf("git init: %v\n%s", err, out)
		}
	}
	t.Chdir(filepath.Join(root, "beta"))
	apps := newTUIModels(globals{}, func(string) string { return "" })

	// Ship defaults to auto: toggle it to owner in beta.
	am := apps.autonomy()
	am.Update(ui.FormSubmittedMsg{Values: ui.Values{
		"rhiza.autonomy.plan": "owner", "rhiza.autonomy.accept": "auto",
		"rhiza.autonomy.ship": "owner", "rhiza.autonomy.signoff": "auto"}})
	gitGet := func(project string) string {
		out, _ := exec.Command("git", "-C", filepath.Join(root, project), "config", "--get", "rhiza.autonomy.ship").Output()
		return strings.TrimSpace(string(out))
	}
	if got := gitGet("beta"); got != "owner" {
		t.Errorf("beta git config ship = %q, want owner", got)
	}
	if got := gitGet("alpha"); got != "" {
		t.Errorf("alpha git config ship = %q, want unset", got)
	}

	da, err := daemonapp.NewIn(func(op string, args []string) ([]byte, error) {
		return apps.call(op, false, false, args)
	}, theme.Dark(), colorprofile.TrueColor, "beta")
	if err != nil {
		t.Fatal(err)
	}
	da.Update(ui.FormSubmittedMsg{Values: ui.Values{
		"paused": false, "maxRuns": "7", "projectPaused": false, "projectMaxRuns": ""}})
	toml, err := os.ReadFile(filepath.Join(root, "rhiza.toml"))
	if err != nil || !strings.Contains(string(toml), "7") {
		t.Errorf("rhiza.toml = %q, %v: want the run limit 7", toml, err)
	}
	if out, _ := exec.Command("git", "-C", filepath.Join(root, "beta"), "config", "--get", "rhiza.daemon.maxRuns").Output(); len(out) > 0 {
		t.Errorf("run limit leaked into git config: %q", out)
	}
}
