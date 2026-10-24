// Package store keeps a Genesis repository on disk.
//
// A repository lives in a .genesis folder at the root of a project:
//
//	.genesis/
//	  repo.json                     {"format": 1, "id": "<ULID>", "remotes": []}
//	  log/
//	    <change_set>/
//	      <id>.json                 one canonical entry per file
//
// The log is append-only. An entry file, once written, is never modified,
// replaced or deleted: Write creates it atomically and refuses an ID that is
// already there, and no method of this package changes an existing file.
// A writer stages each file as a dot-prefixed .tmp file beside its final
// name; readers skip those.
//
// ReadLog reads the whole log at once and never stops at a bad file: it
// returns every valid entry, in (change set, ID) order, beside a Problem
// for each file or folder that is stray, corrupt or from a newer format.
package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/dphaener/genesis/internal/canon"
	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
)

// The on-disk layout, relative to the project root.
const (
	// Dir is the repository folder at the project root.
	Dir = ".genesis"
	// RepoFile is the repository description inside Dir.
	RepoFile = "repo.json"
	// LogDir holds one folder per change set inside Dir.
	LogDir = "log"
)

// FormatVersion is the repo.json format this package reads and writes.
const FormatVersion = 1

var (
	// ErrExists is wrapped by Init's error when the root already has a Dir.
	ErrExists = errors.New("store: repository already exists")
	// ErrNotFound is wrapped by Open's error when there is no Dir or RepoFile.
	ErrNotFound = errors.New("store: no repository")
	// ErrInvalid is wrapped by Open's error when repo.json is not a valid
	// format-1 repository description.
	ErrInvalid = errors.New("store: invalid repo.json")
	// ErrEntryExists is wrapped by Write's error when an entry with the same
	// ID is already in the change set folder.
	ErrEntryExists = errors.New("store: entry already exists")
	// ErrNewerFormat is wrapped by Open's error when repo.json declares a
	// format newer than FormatVersion. It is entry.ErrNewerFormat, so one
	// errors.Is check covers both files, and it does not wrap ErrInvalid.
	ErrNewerFormat = entry.ErrNewerFormat
)

// Repo is the content of repo.json.
type Repo struct {
	Format int
	ID     id.ID
	// Remotes is never nil once decoded.
	Remotes []string
}

// repoKeys are repo.json's fields, in the order Open checks them.
var repoKeys = []string{"format", "id", "remotes"}

func (r Repo) encode() ([]byte, error) {
	remotes := r.Remotes
	if remotes == nil {
		remotes = []string{}
	}
	return canon.Marshal(map[string]any{
		"format":  r.Format,
		"id":      r.ID,
		"remotes": remotes,
	})
}

func invalid(field, format string, args ...any) error {
	return fmt.Errorf("%w: field %q: %s", ErrInvalid, field, fmt.Sprintf(format, args...))
}

func invalidWrap(field string, err error) error {
	return fmt.Errorf("%w: field %q: %w", ErrInvalid, field, err)
}

// decodeRepo parses canonical repo.json bytes. The format is read before
// anything else, since a newer format may have a different shape.
func decodeRepo(data []byte) (Repo, error) {
	if err := canon.Check(data); err != nil {
		return Repo{}, fmt.Errorf("%w: %w", ErrInvalid, err)
	}
	if data[0] != '{' {
		return Repo{}, fmt.Errorf("%w: not an object", ErrInvalid)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return Repo{}, fmt.Errorf("%w: %w", ErrInvalid, err)
	}

	var r Repo
	raw, ok := fields["format"]
	if !ok {
		return Repo{}, invalid("format", "missing")
	}
	f, err := strconv.ParseInt(string(raw), 10, 64)
	if err != nil {
		return Repo{}, invalid("format", "not an integer")
	}
	if f > FormatVersion {
		return Repo{}, fmt.Errorf("%w: field %q: %d, this Genesis reads %d", ErrNewerFormat, "format", f, FormatVersion)
	}
	if f != FormatVersion {
		return Repo{}, invalid("format", "%d, want %d", f, FormatVersion)
	}
	r.Format = int(f)

	for _, k := range repoKeys {
		if _, ok := fields[k]; !ok {
			return Repo{}, invalid(k, "missing")
		}
	}
	if len(fields) != len(repoKeys) {
		var names []string
		for k := range fields {
			if !slices.Contains(repoKeys, k) {
				names = append(names, k)
			}
		}
		slices.Sort(names)
		return Repo{}, invalid(names[0], "unknown field")
	}

	s, err := str(fields["id"])
	if err != nil {
		return Repo{}, invalidWrap("id", err)
	}
	if r.ID, err = id.Parse(s); err != nil {
		return Repo{}, invalidWrap("id", err)
	}

	raw = fields["remotes"]
	if raw[0] != '[' {
		return Repo{}, invalid("remotes", "not an array")
	}
	var remotes []json.RawMessage
	if err := json.Unmarshal(raw, &remotes); err != nil {
		return Repo{}, invalidWrap("remotes", err)
	}
	r.Remotes = make([]string, 0, len(remotes))
	for _, rr := range remotes {
		s, err := str(rr)
		if err != nil {
			return Repo{}, invalidWrap("remotes", err)
		}
		r.Remotes = append(r.Remotes, s)
	}
	return r, nil
}

var errNotString = errors.New("not a string")

// str decodes a JSON string value, rejecting every other kind.
func str(raw json.RawMessage) (string, error) {
	if len(raw) == 0 || raw[0] != '"' {
		return "", errNotString
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", err
	}
	return s, nil
}

// Store is an open repository.
type Store struct {
	root string
	repo Repo
}

// Root returns the project root the repository lives in.
func (s *Store) Root() string { return s.root }

// Repo returns the repository description read from repo.json.
func (s *Store) Repo() Repo { return s.repo }

// dir returns the path of the .genesis folder.
func (s *Store) dir() string { return filepath.Join(s.root, Dir) }

// logDir returns the path of the .genesis/log folder.
func (s *Store) logDir() string { return filepath.Join(s.root, Dir, LogDir) }

// Init creates a new repository in root, with an ID from gen and no
// remotes. It fails with ErrExists, changing nothing, if root already has
// a Dir of any kind.
func Init(root string, gen *id.Generator) (*Store, error) {
	s := &Store{root: root}
	// Mkdir, not MkdirAll: it fails atomically if anything is already
	// there, so two Inits cannot both succeed.
	if err := os.Mkdir(s.dir(), 0o755); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return nil, fmt.Errorf("%w: %s", ErrExists, s.dir())
		}
		return nil, err
	}
	if err := os.Mkdir(s.logDir(), 0o755); err != nil {
		return nil, err
	}
	rid, err := gen.New()
	if err != nil {
		return nil, err
	}
	s.repo = Repo{Format: FormatVersion, ID: rid, Remotes: []string{}}
	data, err := s.repo.encode()
	if err != nil {
		return nil, err
	}
	if err := writeNew(s.dir(), RepoFile, data); err != nil {
		return nil, err
	}
	if err := syncDir(root); err != nil {
		return nil, err
	}
	return s, nil
}

// Open opens the repository in root. It fails with ErrNotFound if there is
// none, with ErrNewerFormat if repo.json is from a newer Genesis, and with
// ErrInvalid if repo.json is not exactly a canonical format-1 description.
func Open(root string) (*Store, error) {
	s := &Store{root: root}
	p := filepath.Join(s.dir(), RepoFile)
	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("%w: %s", ErrNotFound, p)
		}
		return nil, err
	}
	if s.repo, err = decodeRepo(data); err != nil {
		return nil, fmt.Errorf("%s: %w", p, err)
	}
	return s, nil
}

// Write adds e to the log at log/<change_set>/<id>.json. It fails with
// ErrEntryExists, leaving the existing file untouched, if that ID is
// already written, and writes nothing at all if e is not a valid entry.
func (s *Store) Write(e entry.Entry) error {
	data, err := entry.Encode(e)
	if err != nil {
		return err
	}
	dir := filepath.Join(s.logDir(), e.ChangeSet.String())
	created := true
	if err := os.Mkdir(dir, 0o755); err != nil {
		if !errors.Is(err, fs.ErrExist) {
			return err
		}
		created = false
	}
	name := e.ID.String() + ".json"
	if err := writeNew(dir, name, data); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return fmt.Errorf("%w: %s", ErrEntryExists, filepath.Join(dir, name))
		}
		return err
	}
	// A new change set folder is itself a new name in log/.
	if created {
		return syncDir(s.logDir())
	}
	return nil
}

// writeNew creates dir/name holding data, durably and atomically, and
// fails with an error matching fs.ErrExist if the name is taken. The bytes
// go to a temp file first, so the final name never shows a partial file;
// the temp file is then hard-linked into place, because link (unlike
// rename) refuses to replace an existing file.
func writeNew(dir, name string, data []byte) error {
	f, err := os.CreateTemp(dir, "."+strings.TrimSuffix(name, ".json")+".*.tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp, 0o644); err != nil {
		return err
	}
	if err := os.Link(tmp, filepath.Join(dir, name)); err != nil {
		return err
	}
	return syncDir(dir)
}

// syncDir fsyncs a folder so the names created in it survive a crash.
func syncDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	if err := d.Sync(); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}
