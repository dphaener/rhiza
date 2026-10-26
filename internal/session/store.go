package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// MinPrefix is the shortest ID prefix Read accepts.
const MinPrefix = 4

// NotFoundError means no session matches the ID or prefix.
type NotFoundError struct{ ID string }

func (e *NotFoundError) Error() string { return fmt.Sprintf("no session matches %q", e.ID) }

// AmbiguousError means a prefix matches more than one session.
type AmbiguousError struct {
	Prefix  string
	Matches []string
}

func (e *AmbiguousError) Error() string {
	return fmt.Sprintf("%q matches %d sessions: %s", e.Prefix, len(e.Matches), strings.Join(e.Matches, ", "))
}

// Dir is where sessions live for a repository: <commondir>/rhiza/sessions.
func Dir(gitCommonDir string) string {
	return filepath.Join(gitCommonDir, "rhiza", "sessions")
}

// Store keeps one JSON file per session in a directory.
type Store struct {
	dir string
}

var updateMu sync.Mutex

// NewStore returns a store over dir. The directory is created on first write.
func NewStore(dir string) *Store { return &Store{dir: dir} }

func (s *Store) path(id string) string { return filepath.Join(s.dir, id+".json") }

// Write stores r atomically: a temp file in the same folder, then a rename.
func (s *Store) Write(r Record) error {
	if r.ID == "" || strings.ContainsAny(r.ID, `/\`) || strings.HasPrefix(r.ID, ".") {
		return fmt.Errorf("session: invalid ID %q", r.ID)
	}
	r.StartedAt = r.StartedAt.UTC()
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}
	f, err := os.CreateTemp(s.dir, "."+r.ID+".*.tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	_, err = f.Write(data)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Chmod(tmp, 0o644)
	}
	if err == nil {
		err = os.Rename(tmp, s.path(r.ID))
	}
	return err
}

// Read returns the session with the full ID, or the unique one that starts
// with a prefix of MinPrefix or more characters.
func (s *Store) Read(id string) (Record, error) {
	if id == "" || strings.ContainsAny(id, `/\`) {
		return Record{}, &NotFoundError{ID: id}
	}
	if r, err := s.readFile(s.path(id)); err == nil {
		return r, nil
	} else if !errors.Is(err, fs.ErrNotExist) {
		return Record{}, err
	}
	if len(id) < MinPrefix {
		return Record{}, &NotFoundError{ID: id}
	}
	ids, err := s.ids()
	if err != nil {
		return Record{}, err
	}
	var matches []string
	for _, c := range ids {
		if strings.HasPrefix(c, id) {
			matches = append(matches, c)
		}
	}
	switch len(matches) {
	case 0:
		return Record{}, &NotFoundError{ID: id}
	case 1:
		return s.readFile(s.path(matches[0]))
	}
	return Record{}, &AmbiguousError{Prefix: id, Matches: matches}
}

// List returns every session, oldest first, ties broken by ID.
func (s *Store) List() ([]Record, error) {
	ids, err := s.ids()
	if err != nil {
		return nil, err
	}
	var out []Record
	for _, id := range ids {
		r, err := s.readFile(s.path(id))
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if !out[i].StartedAt.Equal(out[j].StartedAt) {
			return out[i].StartedAt.Before(out[j].StartedAt)
		}
		return out[i].ID < out[j].ID
	})
	return out, nil
}

// Update re-reads the session, applies fn and writes it back, so a writer
// only changes what fn touches. It takes a full ID or a unique prefix.
func (s *Store) Update(id string, fn func(*Record)) (Record, error) {
	updateMu.Lock()
	defer updateMu.Unlock()
	r, err := s.Read(id)
	if err != nil {
		return Record{}, err
	}
	fn(&r)
	if err := s.Write(r); err != nil {
		return Record{}, err
	}
	return r, nil
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
		if e.IsDir() || strings.HasPrefix(n, ".") || !strings.HasSuffix(n, ".json") {
			continue
		}
		ids = append(ids, strings.TrimSuffix(n, ".json"))
	}
	sort.Strings(ids)
	return ids, nil
}

func (s *Store) readFile(path string) (Record, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Record{}, err
	}
	var r Record
	if err := json.Unmarshal(data, &r); err != nil {
		return Record{}, fmt.Errorf("%s: %w", path, err)
	}
	return r, nil
}
