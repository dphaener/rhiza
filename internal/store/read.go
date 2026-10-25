package store

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
)

var (
	// ErrStray is wrapped by a Problem's error for a file or folder under
	// log/ that has no place in the layout: anything but a folder directly
	// in log/, a folder not named by an ID, a folder inside a change set,
	// or a file not named <ID>.json.
	ErrStray = errors.New("store: unexpected file or folder")
	// ErrMismatch is wrapped by a Problem's error for an entry file whose
	// decoded id or change_set differs from its file or folder name.
	ErrMismatch = errors.New("store: entry id or change set does not match its path")
)

// Problem is one thing wrong under log/. Path is relative to the .genesis
// folder, in slash form (log/<change_set>/<id>.json), and Err is never nil.
// Decode failures keep the entry error, so errors.Is matches
// entry.ErrInvalid or entry.ErrNewerFormat.
type Problem struct {
	Path string
	Err  error
}

// Log is everything ReadLog found. Entries are sorted by change set, then
// ID; Problems by path, then error text. Neither slice is ever nil.
type Log struct {
	Entries  []entry.Entry
	Problems []Problem
}

// Listing is everything ListLog found. Paths are the entry files whose
// folder and file names parse as IDs, in slash form relative to the
// .genesis folder (log/<change_set>/<id>.json), sorted. Problems are the
// layout problems, sorted by path, then error text. Neither slice is ever
// nil.
type Listing struct {
	Paths    []string
	Problems []Problem
}

// ReadLog reads every entry under log/ in fsys, which is rooted at a
// .genesis folder. It never stops early: each bad file or folder becomes a
// Problem and every good entry is kept, so one corrupt or newer-format file
// cannot hide the rest of the log. Names starting with "." are skipped at
// every level, since writers stage temp files that way and tools leave
// files like .DS_Store. The result does not depend on listing order.
//
// ReadLog is ListLog followed by ReadEntries over the listed paths, with
// the problems of both merged.
func ReadLog(fsys fs.FS) Log {
	l := ListLog(fsys)
	log := ReadEntries(fsys, l.Paths)
	log.Problems = append(log.Problems, l.Problems...)
	sortProblems(log.Problems)
	return log
}

// ReadLog reads the repository's log; see the package-level ReadLog.
func (s *Store) ReadLog() Log { return ReadLog(os.DirFS(s.dir())) }

// ListLog lists the entry files under log/ in fsys, which is rooted at a
// .genesis folder, by name alone: it opens no entry file, so a corrupt
// entry is listed, not reported. It walks log/ as ReadLog does, skipping
// names that start with ".", and reports the same layout problems: a
// missing or unreadable log/ or change set folder, and every stray file or
// folder (ErrStray). A folder whose name is not an ID is not descended
// into. The result does not depend on listing order.
func ListLog(fsys fs.FS) Listing {
	r := lister{fsys: fsys, l: Listing{Paths: []string{}, Problems: []Problem{}}}
	r.walk()
	slices.Sort(r.l.Paths)
	sortProblems(r.l.Problems)
	return r.l
}

// ReadEntries reads and decodes exactly the given entry files in fsys,
// which is rooted at a .genesis folder. Each path is in slash form,
// log/<change_set>/<id>.json, as ListLog returns it; the change set and ID
// an entry must carry are parsed from its path, and a path not of that
// form is an ErrStray problem. Like ReadLog it never stops early: a file
// that cannot be read, does not decode, or whose id or change_set differs
// from its path (ErrMismatch) becomes a Problem, and every good entry is
// kept. Entries and Problems are sorted as in Log and never nil.
func ReadEntries(fsys fs.FS, paths []string) Log {
	log := Log{Entries: []entry.Entry{}, Problems: []Problem{}}
	for _, p := range paths {
		e, err := readEntry(fsys, p)
		if err != nil {
			log.Problems = append(log.Problems, Problem{Path: p, Err: err})
			continue
		}
		log.Entries = append(log.Entries, e)
	}
	slices.SortFunc(log.Entries, func(a, b entry.Entry) int {
		return cmp.Or(id.Compare(a.ChangeSet, b.ChangeSet), id.Compare(a.ID, b.ID))
	})
	sortProblems(log.Problems)
	return log
}

func sortProblems(ps []Problem) {
	slices.SortFunc(ps, func(a, b Problem) int {
		return cmp.Or(strings.Compare(a.Path, b.Path), strings.Compare(a.Err.Error(), b.Err.Error()))
	})
}

// entryIDs parses the change set and entry ID out of an entry file path,
// log/<change_set>/<id>.json.
func entryIDs(p string) (cs, eid id.ID, err error) {
	parts := strings.Split(p, "/")
	if len(parts) != 3 || parts[0] != LogDir {
		return cs, eid, fmt.Errorf("%w: not a log/<change_set>/<id>.json path", ErrStray)
	}
	if cs, err = id.Parse(parts[1]); err != nil {
		return cs, eid, fmt.Errorf("%w: folder name is not a change set ID: %w", ErrStray, err)
	}
	stem, ok := strings.CutSuffix(parts[2], ".json")
	if !ok {
		return cs, eid, fmt.Errorf("%w: not a .json file", ErrStray)
	}
	if eid, err = id.Parse(stem); err != nil {
		return cs, eid, fmt.Errorf("%w: file name is not an entry ID: %w", ErrStray, err)
	}
	return cs, eid, nil
}

func readEntry(fsys fs.FS, p string) (entry.Entry, error) {
	cs, eid, err := entryIDs(p)
	if err != nil {
		return entry.Entry{}, err
	}
	data, err := fs.ReadFile(fsys, p)
	if err != nil {
		return entry.Entry{}, err
	}
	e, err := entry.Decode(data)
	if err != nil {
		return entry.Entry{}, err
	}
	if e.ID != eid {
		return entry.Entry{}, fmt.Errorf("%w: id is %s, file name says %s", ErrMismatch, e.ID, eid)
	}
	if e.ChangeSet != cs {
		return entry.Entry{}, fmt.Errorf("%w: change_set is %s, folder says %s", ErrMismatch, e.ChangeSet, cs)
	}
	return e, nil
}

type lister struct {
	fsys fs.FS
	l    Listing
}

func (r *lister) problem(p string, err error) {
	r.l.Problems = append(r.l.Problems, Problem{Path: p, Err: err})
}

func (r *lister) walk() {
	des, err := fs.ReadDir(r.fsys, LogDir)
	if err != nil {
		r.problem(LogDir, err)
		return
	}
	for _, de := range des {
		name := de.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		p := path.Join(LogDir, name)
		if !de.IsDir() {
			r.problem(p, fmt.Errorf("%w: not a change set folder", ErrStray))
			continue
		}
		if _, err := id.Parse(name); err != nil {
			// Not descended into: its contents belong to no change set.
			r.problem(p, fmt.Errorf("%w: folder name is not a change set ID: %w", ErrStray, err))
			continue
		}
		r.changeSet(p)
	}
}

func (r *lister) changeSet(dir string) {
	des, err := fs.ReadDir(r.fsys, dir)
	if err != nil {
		r.problem(dir, err)
		return
	}
	for _, de := range des {
		name := de.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		p := path.Join(dir, name)
		if de.IsDir() {
			r.problem(p, fmt.Errorf("%w: a folder inside a change set folder", ErrStray))
			continue
		}
		stem, ok := strings.CutSuffix(name, ".json")
		if !ok {
			r.problem(p, fmt.Errorf("%w: not a .json file", ErrStray))
			continue
		}
		if _, err := id.Parse(stem); err != nil {
			r.problem(p, fmt.Errorf("%w: file name is not an entry ID: %w", ErrStray, err))
			continue
		}
		r.l.Paths = append(r.l.Paths, p)
	}
}
