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

// ReadLog reads every entry under log/ in fsys, which is rooted at a
// .genesis folder. It never stops early: each bad file or folder becomes a
// Problem and every good entry is kept, so one corrupt or newer-format file
// cannot hide the rest of the log. Names starting with "." are skipped at
// every level, since writers stage temp files that way and tools leave
// files like .DS_Store. The result does not depend on listing order.
func ReadLog(fsys fs.FS) Log {
	r := reader{fsys: fsys, log: Log{Entries: []entry.Entry{}, Problems: []Problem{}}}
	r.walk()
	slices.SortFunc(r.log.Entries, func(a, b entry.Entry) int {
		return cmp.Or(id.Compare(a.ChangeSet, b.ChangeSet), id.Compare(a.ID, b.ID))
	})
	slices.SortFunc(r.log.Problems, func(a, b Problem) int {
		return cmp.Or(strings.Compare(a.Path, b.Path), strings.Compare(a.Err.Error(), b.Err.Error()))
	})
	return r.log
}

// ReadLog reads the repository's log; see the package-level ReadLog.
func (s *Store) ReadLog() Log { return ReadLog(os.DirFS(s.dir())) }

type reader struct {
	fsys fs.FS
	log  Log
}

func (r *reader) problem(p string, err error) {
	r.log.Problems = append(r.log.Problems, Problem{Path: p, Err: err})
}

func (r *reader) walk() {
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
		cs, err := id.Parse(name)
		if err != nil {
			// Not descended into: its contents belong to no change set.
			r.problem(p, fmt.Errorf("%w: folder name is not a change set ID: %w", ErrStray, err))
			continue
		}
		r.changeSet(p, cs)
	}
}

func (r *reader) changeSet(dir string, cs id.ID) {
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
		eid, err := id.Parse(stem)
		if err != nil {
			r.problem(p, fmt.Errorf("%w: file name is not an entry ID: %w", ErrStray, err))
			continue
		}
		data, err := fs.ReadFile(r.fsys, p)
		if err != nil {
			r.problem(p, err)
			continue
		}
		e, err := entry.Decode(data)
		if err != nil {
			r.problem(p, err)
			continue
		}
		if e.ID != eid {
			r.problem(p, fmt.Errorf("%w: id is %s, file name says %s", ErrMismatch, e.ID, eid))
			continue
		}
		if e.ChangeSet != cs {
			r.problem(p, fmt.Errorf("%w: change_set is %s, folder says %s", ErrMismatch, e.ChangeSet, cs))
			continue
		}
		r.log.Entries = append(r.log.Entries, e)
	}
}
