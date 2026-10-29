// Package launcher is the registry of rhiza apps. rhiza rex install writes
// one Rex action and one key of the rhiza key mode per registered app, each
// opening a layer that runs `rhiza tui app <name>`.
package launcher

import (
	"sort"
	"sync"
)

// App is one rhiza app.
type App struct {
	// Name is the app's name, the title of its action and passed to
	// `rhiza tui app`.
	Name string
	// Summary is the description of the app's action.
	Summary string
	// Cols and Rows are the app layer's size in cells; zero means the
	// default.
	Cols, Rows int
}

var (
	mu       sync.Mutex
	registry = map[string]App{}
)

// Register adds a to the apps, replacing one of the same name.
func Register(a App) {
	mu.Lock()
	defer mu.Unlock()
	registry[a.Name] = a
}

// Apps returns the registered apps, sorted by name.
func Apps() []App {
	mu.Lock()
	defer mu.Unlock()
	out := make([]App, 0, len(registry))
	for _, a := range registry {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

const (
	defaultCols = 100
	defaultRows = 30
)

// Size returns the app layer's size in cells, the defaults where the app
// sets none.
func (a App) Size() (cols, rows int) {
	cols, rows = a.Cols, a.Rows
	if cols == 0 {
		cols = defaultCols
	}
	if rows == 0 {
		rows = defaultRows
	}
	return cols, rows
}
