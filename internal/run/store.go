package run

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dphaener/rhiza/internal/userdir"
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

// TmpDir is the run's own temp folder, outside every git repository:
// <user cache dir>/rhiza/run-tmp/<run id>, or
// <temp dir>/rhiza-run-tmp-<uid>/<run id> when there is no user cache
// dir, both placed by the os package's rules over getenv (userdir; nil
// reads nothing). A workflow's [env] values name it as
// {run_tmp}. The engine creates it as a step starts and removes it (that
// folder only, never the shared run-tmp parent) when the run ends.
func TmpDir(id string, getenv func(string) string) string {
	if c, err := userdir.Cache(getenv); err == nil && c != "" {
		return filepath.Join(c, "rhiza", "run-tmp", id)
	}
	return filepath.Join(userdir.Temp(getenv), fmt.Sprintf("rhiza-run-tmp-%d", os.Getuid()), id)
}

// Store keeps one append-only JSONL journal per run in a directory.
type Store struct {
	dir string
	// Getenv places the runs' temp folders (TmpDir) and the list cache;
	// nil reads nothing (an unset environment, never the process's).
	Getenv func(string) string
	// DiskCache makes List keep what it parsed of each journal in a file
	// under the user cache folder (ListCachePath), so a later process
	// parses only the journals that changed since. Off, List keeps it in
	// this process alone.
	DiskCache bool
}

// NewStore returns a store over dir. The directory is created on first write.
func NewStore(dir string) *Store { return &Store{dir: dir} }

// TmpDir is TmpDir for a run of this store.
func (s *Store) TmpDir(id string) string { return TmpDir(id, s.Getenv) }

// RunDir is the directory a run keeps its files in, such as a gate's failing
// output: <store>/<id>. It is not created.
func (s *Store) RunDir(id string) string { return filepath.Join(s.dir, id) }

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

// AppendSpawnEnded adds a spawn-ended entry unless the run's journal already
// records that child's end, and reports whether it added one. The spawn
// process and a restarted driver both record a child's end; whichever comes
// second finds it recorded.
func (s *Store) AppendSpawnEnded(id string, e Entry) (bool, error) {
	j, err := s.Read(id)
	if err != nil {
		return false, err
	}
	for _, sp := range j.State().Spawns {
		if sp.Session == e.Session && sp.Ended {
			return false, nil
		}
	}
	e.Type = TypeSpawnEnded
	return true, s.Append(id, e)
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
	var disk *diskCache
	if s.DiskCache {
		disk = s.openDiskCache()
	}
	var out []Journal
	for _, id := range ids {
		var j Journal
		if disk != nil {
			j, err = disk.read(s, id)
		} else {
			j, err = s.readFile(id)
		}
		if err != nil {
			return nil, err
		}
		out = append(out, j)
	}
	if disk != nil {
		disk.save()
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
	info, err := os.Stat(path)
	if err != nil {
		// Read it anyway, for the error a read gives.
		journalCache.drop(path)
		return s.parseFile(id, path)
	}
	if es, ok := journalCache.get(path, info); ok {
		return Journal{ID: id, Entries: es}, nil
	}
	j, err := s.parseFile(id, path)
	if err != nil {
		journalCache.drop(path)
		return Journal{}, err
	}
	journalCache.put(path, info, j.Entries)
	return j, nil
}

// journalCache keeps each journal's entries by path with the size and
// modification time its file had just before it was read, so the daemon's
// polls (several reads of every journal each) parse only the journals
// that changed. A journal is only appended to, so a write changes its size.
var journalCache = &entryCache{files: map[string]cachedJournal{}}

// journalReads counts the journal files parsed, process-wide; tests read
// it to see that an unchanged journal is not parsed again.
var journalReads atomic.Int64

type entryCache struct {
	mu    sync.Mutex
	files map[string]cachedJournal
}

type cachedJournal struct {
	size    int64
	mtime   time.Time
	entries []Entry
}

// get returns a copy of the cached entries of path when its file is
// unchanged: callers may append to the slice they get.
func (c *entryCache) get(path string, info os.FileInfo) ([]Entry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f, ok := c.files[path]
	if !ok || f.size != info.Size() || !f.mtime.Equal(info.ModTime()) {
		return nil, false
	}
	return slices.Clone(f.entries), true
}

func (c *entryCache) put(path string, info os.FileInfo, es []Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.files[path] = cachedJournal{size: info.Size(), mtime: info.ModTime(), entries: slices.Clone(es)}
}

// putOwned is put of a slice no one else writes into, kept without a copy.
func (c *entryCache) putOwned(path string, info os.FileInfo, es []Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.files[path] = cachedJournal{size: info.Size(), mtime: info.ModTime(), entries: es}
}

func (c *entryCache) drop(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.files, path)
}

// parseFile reads and parses the journal file of run id at path.
func (s *Store) parseFile(id, path string) (Journal, error) {
	journalReads.Add(1)
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

// The list cache's file starts with this header; a file with another
// format, version or Entry shape, or one that does not decode, is ignored
// and replaced on the next save.
const (
	listCacheFormat  = "rhiza run list cache"
	listCacheVersion = 1
)

type listCacheHeader struct {
	Format  string
	Version int
	Shape   string
}

// listCacheShape keys the list cache to the Entry type it holds: a file
// written by a binary whose Entry, or a type it holds, has other fields is
// ignored, since gob would decode its records into this Entry leaving a
// new field zero.
var listCacheShape = typeShape(reflect.TypeOf(Entry{}))

// typeShape is a hash of t's shape: its kind, and recursively each
// exported field's name and shape, or its element's, key's and value's.
func typeShape(t reflect.Type) string {
	var b strings.Builder
	writeTypeShape(&b, t, map[reflect.Type]bool{})
	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}

func writeTypeShape(b *strings.Builder, t reflect.Type, seen map[reflect.Type]bool) {
	b.WriteString(t.String())
	switch t.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Array:
		b.WriteString("[")
		writeTypeShape(b, t.Elem(), seen)
		b.WriteString("]")
	case reflect.Map:
		b.WriteString("[")
		writeTypeShape(b, t.Key(), seen)
		b.WriteString("]")
		writeTypeShape(b, t.Elem(), seen)
	case reflect.Struct:
		if seen[t] {
			return
		}
		seen[t] = true
		b.WriteString("{")
		for i := range t.NumField() {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}
			b.WriteString(f.Name)
			b.WriteString(" ")
			writeTypeShape(b, f.Type, seen)
			b.WriteString(";")
		}
		b.WriteString("}")
	}
}

// listCacheRecord is what List parsed of one journal: its entries, as the
// file at Path was when it had Size bytes and ModTime (Unix nanoseconds).
type listCacheRecord struct {
	Path    string
	Size    int64
	ModTime int64
	Entries []Entry
}

// ListCachePath is the file a DiskCache store keeps its list cache in:
// <user cache dir>/rhiza/run-list/<hash of the store's folder>.gob, placed
// by the os package's rules over getenv (userdir). ok is false when there
// is no user cache dir.
func ListCachePath(dir string, getenv func(string) string) (path string, ok bool) {
	c, err := userdir.Cache(getenv)
	if err != nil || c == "" {
		return "", false
	}
	if abs, err := filepath.Abs(dir); err == nil {
		dir = abs
	}
	sum := sha256.Sum256([]byte(dir))
	return filepath.Join(c, "rhiza", "run-list", hex.EncodeToString(sum[:16])+".gob"), true
}

// diskCache is one List call's view of the list cache: the records the
// file held, and the ones the call would save, one per journal it listed.
//
// The file is rewritten only when a finished journal is new or changed,
// or a recorded one is gone: a live run's journal changes all the time,
// and rewriting every record for it would cost what parsing them did. So
// a live journal's record is refreshed with the next such save, and until
// then the journal is parsed whenever it changed since its record.
type diskCache struct {
	file  string
	old   map[string]listCacheRecord
	next  []listCacheRecord
	dirty bool
}

// openDiskCache loads the store's list cache, or returns nil when there is
// no user cache dir. A missing, corrupt or other-version file loads empty.
func (s *Store) openDiskCache() *diskCache {
	file, ok := ListCachePath(s.dir, s.Getenv)
	if !ok {
		return nil
	}
	d := &diskCache{file: file, old: map[string]listCacheRecord{}}
	f, err := os.Open(file)
	if err != nil {
		return d
	}
	defer f.Close()
	dec := gob.NewDecoder(bufio.NewReader(f))
	var h listCacheHeader
	if dec.Decode(&h) != nil || h.Format != listCacheFormat || h.Version != listCacheVersion ||
		h.Shape != listCacheShape {
		return d
	}
	var recs []listCacheRecord
	if dec.Decode(&recs) != nil {
		return d
	}
	for _, r := range recs {
		d.old[r.Path] = r
	}
	return d
}

// read is readFile through the list cache: an unchanged journal comes from
// this process's cache or the file's record, a new or changed one is
// parsed and recorded.
func (d *diskCache) read(s *Store, id string) (Journal, error) {
	path := s.path(id)
	info, err := os.Stat(path)
	if err != nil {
		// Read it anyway, for the error a read gives.
		return s.readFile(id)
	}
	size, mtime := info.Size(), info.ModTime().UnixNano()
	rec, have := d.old[path]
	have = have && rec.Size == size && rec.ModTime == mtime
	if es, ok := journalCache.get(path, info); ok {
		if have {
			d.next = append(d.next, rec)
		} else {
			d.record(id, path, size, mtime, es)
		}
		return Journal{ID: id, Entries: es}, nil
	}
	if have {
		// The record is this call's own: the process cache shares its
		// entries, clipped so an append by either never reaches the other.
		rec.Entries = slices.Clip(rec.Entries)
		journalCache.putOwned(path, info, rec.Entries)
		d.next = append(d.next, rec)
		return Journal{ID: id, Entries: rec.Entries}, nil
	}
	j, err := s.parseFile(id, path)
	if err != nil {
		journalCache.drop(path)
		return Journal{}, err
	}
	journalCache.put(path, info, j.Entries)
	d.record(id, path, size, mtime, j.Entries)
	return j, nil
}

// record adds a journal's entries to the records to save, when they
// decode back from gob exactly as they are; otherwise the journal goes
// unrecorded and is parsed whenever it is listed. A finished journal's
// record makes the file worth saving.
func (d *diskCache) record(id, path string, size, mtime int64, es []Entry) {
	if !gobExact(es) {
		return
	}
	if (Journal{ID: id, Entries: es}).State().Ended {
		d.dirty = true
	}
	d.next = append(d.next, listCacheRecord{Path: path, Size: size, ModTime: mtime, Entries: es})
}

// gobExact reports whether es comes back from gob as it went in: gob
// drops what JSON keeps, such as an empty slice that is not nil.
func gobExact(es []Entry) bool {
	var buf bytes.Buffer
	if gob.NewEncoder(&buf).Encode(es) != nil {
		return false
	}
	var back []Entry
	if gob.NewDecoder(&buf).Decode(&back) != nil {
		return false
	}
	return reflect.DeepEqual(back, es)
}

// save writes the records, through a temporary file renamed over the
// cache file, when a finished journal was recorded or a recorded journal
// was not listed. The cache is a speed-up only, so a failed write is
// dropped.
func (d *diskCache) save() {
	if !d.dirty && !d.lost() {
		return
	}
	dir := filepath.Dir(d.file)
	if os.MkdirAll(dir, 0o700) != nil {
		return
	}
	f, err := os.CreateTemp(dir, ".run-list-*")
	if err != nil {
		return
	}
	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w)
	err = enc.Encode(listCacheHeader{Format: listCacheFormat, Version: listCacheVersion, Shape: listCacheShape})
	if err == nil {
		err = enc.Encode(d.next)
	}
	if err == nil {
		err = w.Flush()
	}
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Rename(f.Name(), d.file)
	}
	if err != nil {
		os.Remove(f.Name())
	}
}

// lost reports whether a record of the file names no journal this call
// listed: its run was removed.
func (d *diskCache) lost() bool {
	if len(d.old) == 0 {
		return false
	}
	listed := make(map[string]bool, len(d.next))
	for _, r := range d.next {
		listed[r.Path] = true
	}
	for p := range d.old {
		if !listed[p] {
			return true
		}
	}
	return false
}
