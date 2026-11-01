// Package changesapp is the changes app: the landed change sets the viewer
// has not seen, newest first, with their PR numbers. It reads the feed
// operation in-process and remembers the change set IDs it has shown in a
// file under the git common dir (.git/rhiza/feed.json: per clone, shared by
// its worktrees). It remembers IDs rather than a cursor, so work landed on a
// branch and merged later is not skipped. The app starts no daemon and binds
// no socket.
package changesapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/tui/launcher"
)

// Name is the app's name, as `rhiza tui app` takes it.
const Name = "changes"

// SeenFile is the name of the seen file under the git common dir's rhiza
// folder.
const SeenFile = "feed.json"

// pageSize is how many change sets one feed read asks for.
const pageSize = 1000

func init() {
	launcher.Register(launcher.App{Name: Name, Key: "g", Summary: "Landed change sets you have not seen yet"})
}

// Reader runs the feed operation on args and returns its JSON, as
// `rhiza feed <args> --json` prints it.
type Reader func(args []string) ([]byte, error)

// Seen is the content of the seen file.
type seenFile struct {
	Seen []string `json:"seen"`
}

// Model is the app's bubbletea program. j/k scroll; q or esc quits.
type Model struct {
	unseen   []api.FeedItem
	seenPath string
	seen     map[string]bool
	err      error
	top      int
	rows     int
}

// NewModel reads the whole feed with read and the seen IDs from seenPath
// (a missing file means none seen), and keeps the change sets whose ID is
// not seen, newest first. It writes nothing; MarkSeen does.
func NewModel(read Reader, seenPath string) *Model {
	m := &Model{seenPath: seenPath, rows: 20, seen: map[string]bool{}}
	if err := m.loadSeen(); err != nil {
		m.err = err
		return m
	}
	var all []api.FeedItem
	cursor := ""
	for {
		args := []string{"--limit", fmt.Sprint(pageSize)}
		if cursor != "" {
			args = append(args, "--cursor", cursor)
		}
		out, err := read(args)
		if err != nil {
			m.err = err
			return m
		}
		var page api.FeedResponse
		if err := json.Unmarshal(out, &page); err != nil {
			m.err = fmt.Errorf("feed: %w", err)
			return m
		}
		all = append(all, page.ChangeSets...)
		if page.NextCursor == "" {
			break
		}
		cursor = page.NextCursor
	}
	// The feed is oldest land first; the app shows newest first.
	for _, it := range slices.Backward(all) {
		if !m.seen[it.ID.String()] {
			m.unseen = append(m.unseen, it)
		}
	}
	return m
}

func (m *Model) loadSeen() error {
	b, err := os.ReadFile(m.seenPath)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	var f seenFile
	if err := json.Unmarshal(b, &f); err != nil {
		return fmt.Errorf("%s: %w", m.seenPath, err)
	}
	for _, s := range f.Seen {
		m.seen[s] = true
	}
	return nil
}

// Err is the error the read ran into, if any.
func (m *Model) Err() error { return m.err }

// Unseen returns the change sets shown, newest first.
func (m *Model) Unseen() []api.FeedItem { return m.unseen }

// MarkSeen records the shown change sets as seen in the seen file.
func (m *Model) MarkSeen() error {
	if m.err != nil || len(m.unseen) == 0 {
		return m.err
	}
	for _, it := range m.unseen {
		m.seen[it.ID.String()] = true
	}
	ids := make([]string, 0, len(m.seen))
	for s := range m.seen {
		ids = append(ids, s)
	}
	slices.Sort(ids)
	b, err := json.MarshalIndent(seenFile{Seen: ids}, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(m.seenPath), 0o755); err != nil {
		return err
	}
	tmp := m.seenPath + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, m.seenPath)
}

// Result is what --once prints: the unseen change sets, newest first, as
// canonical JSON.
func (m *Model) Result() ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	items := m.unseen
	if items == nil {
		items = []api.FeedItem{}
	}
	return canon.Marshal(api.FeedResponse{ChangeSets: items})
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.rows = max(msg.Height-4, 1)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "down", "j":
			m.scroll(1)
		case "up", "k":
			m.scroll(-1)
		case "pgdown", "space":
			m.scroll(m.rows)
		case "pgup":
			m.scroll(-m.rows)
		case "home", "g":
			m.scroll(-len(m.unseen))
		case "end", "G":
			m.scroll(len(m.unseen))
		}
	}
	return m, nil
}

func (m *Model) scroll(d int) {
	m.top = max(0, min(m.top+d, len(m.unseen)-m.rows))
}

// line is one change set's row.
func line(it api.FeedItem) string {
	pr := "     "
	if it.PR > 0 {
		pr = fmt.Sprintf("#%-4d", it.PR)
	}
	return fmt.Sprintf("%s %s  %s  %s", pr, it.LandedAt, it.ID.String()[:8], it.Title)
}

func (m *Model) View() tea.View {
	var b strings.Builder
	switch {
	case m.err != nil:
		fmt.Fprintf(&b, "rhiza: %v\n", m.err)
	case len(m.unseen) == 0:
		b.WriteString("Nothing new has landed since you last looked.\n")
	default:
		fmt.Fprintf(&b, "%d landed since you last looked\n\n", len(m.unseen))
		end := min(m.top+m.rows, len(m.unseen))
		for _, it := range m.unseen[m.top:end] {
			b.WriteString(line(it) + "\n")
		}
	}
	b.WriteString("\nj/k scroll · q quit")
	return tea.NewView(b.String())
}
