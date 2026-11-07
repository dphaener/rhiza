package modalshell

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rexuikit/ui"
	"github.com/dphaener/rexuikit/ui/theme"
)

// dataFor is the list each project's app shows.
var dataFor = map[string][]string{
	"alpha": {"alpha-item-1", "alpha-item-2"},
	"beta":  {"beta-item-1"},
}

func switching(cwd *string) Switching {
	return Switching{
		Projects: []string{"alpha", "beta"},
		Folder:   func(p string) string { return "/root/" + p },
		Chdir:    func(d string) error { *cwd = strings.TrimPrefix(d, "/root/"); return nil },
		Header:   func() string { return "Project: " + *cwd },
		Build: func() Built {
			var rows []ui.Row
			for _, l := range dataFor[*cwd] {
				rows = append(rows, ui.Row{Label: l})
			}
			sh := New(List(ui.NewList(theme.Dark(), prof, rows)))
			return Built{Model: sh, Result: sh.Result}
		},
	}
}

// run feeds msg and then every message its command yields, as the program does.
func run(t *testing.T, m tea.Model, msg tea.Msg) tea.Model {
	t.Helper()
	for i := 0; msg != nil && i < 10; i++ {
		var cmd tea.Cmd
		m, cmd = m.Update(msg)
		msg = nil
		if cmd != nil {
			msg = cmd()
			if _, ok := msg.(tea.BatchMsg); ok {
				msg = nil
			}
		}
	}
	return m
}

func key(s string) tea.KeyPressMsg {
	switch s {
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	case "down":
		return tea.KeyPressMsg{Code: tea.KeyDown}
	case "esc":
		return tea.KeyPressMsg{Code: tea.KeyEscape}
	}
	return tea.KeyPressMsg{Code: 'o', Mod: tea.ModCtrl}
}

func TestSwitchProjectReloadsHeaderAndData(t *testing.T) {
	cwd := "alpha"
	s := NewSwitcher(switching(&cwd))
	view := func() string { return s.View().Content }
	if v := view(); !strings.Contains(v, "Project: alpha") || !strings.Contains(v, "alpha-item-1") {
		t.Fatalf("start view:\n%s", v)
	}
	var m tea.Model = s
	m = run(t, m, key("ctrl+o"))
	if !s.Picking() {
		t.Fatal("the switch key did not open the picker")
	}
	if v := view(); !strings.Contains(v, "alpha") || !strings.Contains(v, "beta") {
		t.Fatalf("picker lists the projects:\n%s", v)
	}
	m = run(t, m, key("down"))
	m = run(t, m, key("enter"))
	if s.Picking() {
		t.Fatal("picker still open after picking")
	}
	v := view()
	if !strings.Contains(v, "Project: beta") || !strings.Contains(v, "beta-item-1") || strings.Contains(v, "alpha-item") {
		t.Fatalf("view after the switch is not beta's:\n%s", v)
	}
	if cwd != "beta" {
		t.Errorf("cwd = %q, want beta", cwd)
	}
}

func TestSwitchPickerEscKeepsProject(t *testing.T) {
	cwd := "alpha"
	s := NewSwitcher(switching(&cwd))
	var m tea.Model = s
	m = run(t, m, key("ctrl+o"))
	m = run(t, m, key("esc"))
	if s.Picking() || cwd != "alpha" || !strings.Contains(s.View().Content, "alpha-item-1") {
		t.Fatalf("esc should close only the picker; cwd=%q picking=%v", cwd, s.Picking())
	}
	if s.Model().(*Shell).Closed() {
		t.Error("esc in the picker closed the app")
	}
}

func TestSwitchNoProjectsIgnoresKey(t *testing.T) {
	cwd := "alpha"
	sw := switching(&cwd)
	sw.Projects = nil
	s := NewSwitcher(sw)
	run(t, s, key("ctrl+o"))
	if s.Picking() {
		t.Error("picker opened with no projects")
	}
}
