package modalshell

import (
	"os"
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// growing is a model whose NaturalSize grows by a row on each "j".
type growing struct{ rows int }

func (g *growing) Init() tea.Cmd { return nil }
func (g *growing) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok && k.String() == "j" {
		g.rows++
	}
	return g, nil
}
func (g *growing) View() tea.View          { return tea.NewView("") }
func (g *growing) NaturalSize() (int, int) { return 20, g.rows }

// TestReportsSizeAsContentChanges checks a shell with a size file writes
// its model's NaturalSize after Init and again whenever it changes, and
// that a model with no NaturalSize writes nothing.
func TestReportsSizeAsContentChanges(t *testing.T) {
	path := filepath.Join(t.TempDir(), "size")
	ReportSizeTo(path)
	defer ReportSizeTo("")
	m := sized(&growing{rows: 3})
	m.Init()
	if c, r, err := ReadSize(path); err != nil || c != 20 || r != 3 {
		t.Fatalf("after Init: %dx%d %v, want 20x3", c, r, err)
	}
	m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	if c, r, err := ReadSize(path); err != nil || c != 20 || r != 5 {
		t.Fatalf("after growing: %dx%d %v, want 20x5", c, r, err)
	}

	other := filepath.Join(t.TempDir(), "size")
	ReportSizeTo(other)
	plain := sized(struct{ tea.Model }{&growing{}})
	plain.Init()
	if _, err := os.Stat(other); !os.IsNotExist(err) {
		t.Errorf("a model without NaturalSize wrote a size: %v", err)
	}

	ReportSizeTo("")
	if g := (&growing{}); sized(g) != tea.Model(g) {
		t.Error("without a size file the model is wrapped")
	}
}
