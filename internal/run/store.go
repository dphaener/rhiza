package run

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// MinPrefix is the shortest ID prefix Read accepts.
const MinPrefix = 4

// NotFoundError means no run matches the ID or prefix.
type NotFoundError struct{ ID string }

func (e *NotFoundError) Error() string { return fmt.Sprintf("no run matches %q", e.ID) }

// AmbiguousError means a prefix matches more than one run.
type AmbiguousError struct {
	Prefix  string
	Matches []string
}

func (e *AmbiguousError) Error() string {
	return fmt.Sprintf("%q matches %d runs: %s", e.Prefix, len(e.Matches), strings.Join(e.Matches, ", "))
}

// Dir is where runs live for a repository: <commondir>/rhiza/runs.
func Dir(gitCommonDir string) string {
	return filepath.Join(gitCommonDir, "rhiza", "runs")
}

// Store keeps one append-only JSONL journal per run in a directory.
type Store struct {
	dir string
}

// NewStore returns a store over dir. The directory is created on first write.
func NewStore(dir string) *Store { return &Store{dir: dir} }

func (s *Store) path(id string) string { return filepath.Join(s.dir, id+".jsonl") }

func validID(id string) error {
	if id == "" || strings.ContainsAny(id, `/\`) || strings.HasPrefix(id, ".") {
		return fmt.Errorf("run: invalid ID %q", id)
	}
	return nil
}

func line(e Entry) ([]byte, error) {
	if !e.At.IsZero() {
		e.At = e.At.UTC()
	}
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

// Create starts the journal of a new run with its first entry. It fails if
// the run exists.
func (s *Store) Create(id string, first Entry) error {
	if err := validID(id); err != nil {
		return err
	}
	data, err := line(first)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(s.path(id), os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_APPEND, 0o644)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return fmt.Errorf("run: %s already exists", id)
		}
		return err
	}
	_, err = f.Write(data)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	return err
}

// Append adds one entry to an existing run's journal.
func (s *Store) Append(id string, e Entry) error {
	if err := validID(id); err != nil {
		return err
	}
	data, err := line(e)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(s.path(id), os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return &NotFoundError{ID: id}
		}
		return err
	}
	_, err = f.Write(data)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	return err
}

// Read returns the run with the full ID, or the unique one that starts with
// a prefix of MinPrefix or more characters.
func (s *Store) Read(id string) (Journal, error) {
	if id == "" || strings.ContainsAny(id, `/\`) {
		return Journal{}, &NotFoundError{ID: id}
	}
	if j, err := s.readFile(id); err == nil {
		return j, nil
	} else if !errors.Is(err, fs.ErrNotExist) {
		return Journal{}, err
	}
	if len(id) < MinPrefix || strings.HasPrefix(id, ".") {
		return Journal{}, &NotFoundError{ID: id}
	}
	ids, err := s.ids()
	if err != nil {
		return Journal{}, err
	}
	var matches []string
	for _, c := range ids {
		if strings.HasPrefix(c, id) {
			matches = append(matches, c)
		}
	}
	switch len(matches) {
	case 0:
		return Journal{}, &NotFoundError{ID: id}
	case 1:
		return s.readFile(matches[0])
	}
	return Journal{}, &AmbiguousError{Prefix: id, Matches: matches}
}

// List returns every run, oldest started first, ties broken by ID.
func (s *Store) List() ([]Journal, error) {
	ids, err := s.ids()
	if err != nil {
		return nil, err
	}
	var out []Journal
	for _, id := range ids {
		j, err := s.readFile(id)
		if err != nil {
			return nil, err
		}
		out = append(out, j)
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i].started(), out[j].started()
		if !a.Equal(b) {
			return a.Before(b)
		}
		return out[i].ID < out[j].ID
	})
	return out, nil
}

func (s *Store) ids() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var ids []string
	for _, e := range entries {
		n := e.Name()
		if e.IsDir() || strings.HasPrefix(n, ".") || !strings.HasSuffix(n, ".jsonl") {
			continue
		}
		ids = append(ids, strings.TrimSuffix(n, ".jsonl"))
	}
	sort.Strings(ids)
	return ids, nil
}

// readFile reads a journal. A final line with no newline is a torn append
// and is ignored; any other unreadable line is an error.
func (s *Store) readFile(id string) (Journal, error) {
	path := s.path(id)
	data, err := os.ReadFile(path)
	if err != nil {
		return Journal{}, err
	}
	j := Journal{ID: id}
	for n := 1; len(data) > 0; n++ {
		i := bytes.IndexByte(data, '\n')
		if i < 0 {
			break
		}
		raw := bytes.TrimSpace(data[:i])
		data = data[i+1:]
		if len(raw) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(raw, &e); err != nil {
			return Journal{}, fmt.Errorf("%s:%d: %w", path, n, err)
		}
		j.Entries = append(j.Entries, e)
	}
	return j, nil
}

func (j Journal) started() (t time.Time) {
	for _, e := range j.Entries {
		if e.Type == TypeStarted {
			return e.At
		}
	}
	return
}

// LiveOf returns the first run of the change set, other than except, whose
// journal has not finished: a run that is going, stopped at a wait or
// interrupted all count.
func (s *Store) LiveOf(change, except string) (Journal, bool, error) {
	runs, err := s.List()
	if err != nil {
		return Journal{}, false, err
	}
	for _, j := range runs {
		if j.ID == except {
			continue
		}
		if st := j.State(); st.Change == change && !st.Ended {
			return j, true, nil
		}
	}
	return Journal{}, false, nil
}
