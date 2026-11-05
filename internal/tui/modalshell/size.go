package modalshell

// The shell reports its content's size so the opener (internal/tui/modalopen)
// can resize the modal as the content changes: each time the hosted model's
// NaturalSize changes, the shell writes it to the size file the opener named
// with SizeFileFlag, and the opener resizes the layer to fit it.

import (
	"fmt"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
)

// SizeFileFlag is the flag of `rhiza tui app` naming the file the shell
// writes its content's size to.
const SizeFileFlag = "--size-file"

var sizeFile string

// ReportSizeTo makes Run and Shell.Run write the hosted model's NaturalSize
// to path each time it changes; an empty path reports nothing.
func ReportSizeTo(path string) { sizeFile = path }

// sized wraps m to report its size when sizeFile is set.
func sized(m tea.Model) tea.Model {
	if sizeFile == "" {
		return m
	}
	return &sizeReporter{Model: m, path: sizeFile}
}

// sizeReporter writes the NaturalSize of its model, when it has one, to
// path after Init and after every Update that changes it.
type sizeReporter struct {
	tea.Model
	path string
	last [2]int
}

func (r *sizeReporter) Init() tea.Cmd {
	cmd := r.Model.Init()
	r.report()
	return cmd
}

func (r *sizeReporter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m, cmd := r.Model.Update(msg)
	r.Model = m
	r.report()
	return r, cmd
}

func (r *sizeReporter) report() {
	n, ok := r.Model.(interface{ NaturalSize() (int, int) })
	if !ok {
		return
	}
	cols, rows := n.NaturalSize()
	if cols <= 0 || rows <= 0 || [2]int{cols, rows} == r.last {
		return
	}
	if WriteSize(r.path, cols, rows) == nil {
		r.last = [2]int{cols, rows}
	}
}

// WriteSize writes cols x rows to path atomically (a temp file renamed over
// it), so a reader never sees half of it.
func WriteSize(path string, cols, rows int) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".size-*")
	if err != nil {
		return err
	}
	_, werr := fmt.Fprintf(f, "%d %d\n", cols, rows)
	if cerr := f.Close(); werr == nil {
		werr = cerr
	}
	if werr == nil {
		werr = os.Rename(f.Name(), path)
	}
	if werr != nil {
		os.Remove(f.Name())
	}
	return werr
}

// ReadSize reads a size WriteSize wrote.
func ReadSize(path string) (cols, rows int, err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0, 0, err
	}
	if _, err := fmt.Sscanf(string(b), "%d %d", &cols, &rows); err != nil {
		return 0, 0, fmt.Errorf("size file %s: %w", path, err)
	}
	return cols, rows, nil
}
