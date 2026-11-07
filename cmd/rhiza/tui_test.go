package main

import (
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
