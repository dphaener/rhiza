package index

import (
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/store"
)

// tb is what *testing.T and *rapid.T share.
type tb interface {
	Helper()
	Fatalf(format string, args ...any)
}

// mkID is fixture ID n of family f: IDs of one family grow with n.
func mkID(f byte, n int) id.ID {
	var i id.ID
	i[0] = f
	binary.BigEndian.PutUint32(i[12:], uint32(n))
	return i
}

func csID(n int) id.ID             { return mkID(2, n) }
func eID(n int) id.ID              { return mkID(3, n) }
func opID(n int) id.ID             { return mkID(4, n) }
func nd(n int) id.ID               { return mkID(5, n) }
func ln(n int) id.ID               { return mkID(6, n) }
func actorID() id.ID               { return mkID(7, 1) }
func raw(s string) json.RawMessage { return json.RawMessage(s) }

var fixtureTime = time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)

// encodeEntry builds and encodes an entry of cs with body b after after,
// returning its path under .rhiza and its bytes.
func encodeEntry(t tb, eid, cs id.ID, after []id.ID, b changeset.Body) (string, []byte) {
	t.Helper()
	body, err := changeset.EncodeBody(b)
	if err != nil {
		t.Fatalf("encode body: %v", err)
	}
	e := entry.Entry{V: changeset.EntryVersion(b), ID: eid, ChangeSet: cs, Actor: actorID(), For: actorID(),
		At: fixtureTime, After: after, Type: b.Type(), Body: body}
	data, err := entry.Encode(e)
	if err != nil {
		t.Fatalf("encode entry: %v", err)
	}
	return entryPath(cs, eid), data
}

// writeFile writes data at path p (slash form) under the .rhiza folder
// gdir.
func writeFile(t tb, gdir, p string, data []byte) {
	t.Helper()
	full := filepath.Join(gdir, filepath.FromSlash(p))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(full, data, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

// fixture writes a log one entry at a time, each change set's entries in
// a chain.
type fixture struct {
	t     *testing.T
	root  string
	gdir  string
	n     int
	heads map[id.ID]id.ID
	paths map[id.ID]string
}

func newFixture(t *testing.T) *fixture {
	root := t.TempDir()
	gdir := filepath.Join(root, store.Dir)
	if err := os.MkdirAll(filepath.Join(gdir, store.LogDir), 0o755); err != nil {
		t.Fatal(err)
	}
	return &fixture{t: t, root: root, gdir: gdir, heads: map[id.ID]id.ID{}, paths: map[id.ID]string{}}
}

// add writes an entry with body b to cs after its last one and returns the
// entry's ID.
func (f *fixture) add(cs id.ID, b changeset.Body) id.ID {
	f.t.Helper()
	f.n++
	eid := eID(f.n)
	var after []id.ID
	if h, ok := f.heads[cs]; ok {
		after = []id.ID{h}
	}
	p, data := encodeEntry(f.t, eid, cs, after, b)
	writeFile(f.t, f.gdir, p, data)
	f.heads[cs] = eid
	f.paths[eid] = p
	return eid
}

// remove deletes entry eid's file.
func (f *fixture) remove(eid id.ID) {
	f.t.Helper()
	if err := os.Remove(filepath.Join(f.gdir, filepath.FromSlash(f.paths[eid]))); err != nil {
		f.t.Fatal(err)
	}
}

func (f *fixture) dbPath() string { return filepath.Join(f.t.TempDir(), "index.db") }

// countingFS counts the entry files opened through it.
type countingFS struct {
	fsys  fs.FS
	mu    sync.Mutex
	reads map[string]int
}

func newCountingFS(gdir string) *countingFS {
	return &countingFS{fsys: os.DirFS(gdir), reads: map[string]int{}}
}

// Open implements fs.FS. countingFS implements nothing else, so fs.ReadFile
// and fs.ReadDir both come through here.
func (c *countingFS) Open(name string) (fs.File, error) {
	if strings.HasSuffix(name, ".json") {
		c.mu.Lock()
		c.reads[name]++
		c.mu.Unlock()
	}
	return c.fsys.Open(name)
}

func (c *countingFS) total() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := 0
	for _, v := range c.reads {
		n += v
	}
	return n
}

// storedNodes reads a whole stored view's nodes, sorted by ID, nil when
// there are none, as model.Model.Nodes returns them.
func storedNodes(t tb, x *Index, view string) []model.Node {
	t.Helper()
	rows, err := x.db.Query(nodeSelect+` WHERE view = ? ORDER BY id`, view)
	if err != nil {
		t.Fatalf("query nodes: %v", err)
	}
	defer rows.Close()
	var out []model.Node
	for rows.Next() {
		n, err := scanNode(rows)
		if err != nil {
			t.Fatalf("scan node: %v", err)
		}
		out = append(out, n)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("nodes: %v", err)
	}
	return out
}

// storedLinks reads a whole stored view's links, as storedNodes.
func storedLinks(t tb, x *Index, view string) []model.Link {
	t.Helper()
	rows, err := x.db.Query(linkSelect+` WHERE view = ? ORDER BY id`, view)
	if err != nil {
		t.Fatalf("query links: %v", err)
	}
	defer rows.Close()
	var out []model.Link
	for rows.Next() {
		l, err := scanLink(rows)
		if err != nil {
			t.Fatalf("scan link: %v", err)
		}
		out = append(out, l)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("links: %v", err)
	}
	return out
}

// storedProblems reads every problem of source in order, nil when none.
func storedProblems(t tb, x *Index, source string) []storedProblem {
	t.Helper()
	rows, err := x.db.Query(`SELECT source, path, change_set, entry, op, sentinel, message
		FROM problems WHERE source = ? ORDER BY seq`, source)
	if err != nil {
		t.Fatalf("query problems: %v", err)
	}
	defer rows.Close()
	var out []storedProblem
	for rows.Next() {
		var p storedProblem
		if err := rows.Scan(&p.Source, &p.Path, &p.ChangeSet, &p.Entry, &p.Op, &p.Sentinel, &p.Message); err != nil {
			t.Fatalf("scan problem: %v", err)
		}
		out = append(out, p)
	}
	return out
}

// csRow is one stored change set.
type csRow struct {
	ID, State, Kind, Intent, OpenEntry, SupersededBy, TerminalEntry string
}

func storedChangeSets(t tb, x *Index) []csRow {
	t.Helper()
	rows, err := x.db.Query(`SELECT id, state, kind, intent, open_entry, superseded_by, terminal_entry
		FROM change_sets ORDER BY id`)
	if err != nil {
		t.Fatalf("query change sets: %v", err)
	}
	defer rows.Close()
	var out []csRow
	for rows.Next() {
		var r csRow
		if err := rows.Scan(&r.ID, &r.State, &r.Kind, &r.Intent, &r.OpenEntry, &r.SupersededBy, &r.TerminalEntry); err != nil {
			t.Fatalf("scan change set: %v", err)
		}
		out = append(out, r)
	}
	return out
}

func csRows(sets []changeset.ChangeSet) []csRow {
	var out []csRow
	for _, s := range sets {
		out = append(out, csRow{idText(s.ID), string(s.State), string(s.Open.Kind), s.Open.Intent,
			idText(s.OpenEntry), idText(s.SupersededBy), idText(s.TerminalEntry)})
	}
	return out
}

func viewProblems(view string, ps []model.Problem) []storedProblem {
	var out []storedProblem
	for _, p := range ps {
		out = append(out, storedProblem{Source: view, ChangeSet: idText(p.ChangeSet), Op: idText(p.Op),
			Sentinel: modelSentinel(p.Err), Message: p.Err.Error()})
	}
	return out
}

// checkModel fails unless everything stored equals what store.ReadLog,
// changeset.FoldLog and model.AsIs/ToBe compute from the .rhiza folder
// gdir.
func checkModel(t tb, x *Index, gdir string) {
	t.Helper()
	l := store.ReadLog(os.DirFS(gdir))
	sets, fps := changeset.FoldLog(l)

	var logWant []storedProblem
	for _, p := range l.Problems {
		logWant = append(logWant, storedProblem{Source: sourceLog, Path: p.Path, Sentinel: logSentinel(p.Err), Message: p.Err.Error()})
	}
	same(t, "log problems", storedProblems(t, x, sourceLog), logWant)

	var foldWant []storedProblem
	for _, p := range fps {
		foldWant = append(foldWant, storedProblem{Source: sourceFold, ChangeSet: idText(p.ChangeSet), Entry: idText(p.Entry),
			Sentinel: foldSentinel(p.Err), Message: p.Err.Error()})
	}
	same(t, "fold problems", storedProblems(t, x, sourceFold), foldWant)
	same(t, "change sets", storedChangeSets(t, x), csRows(sets))

	for _, v := range []struct {
		name string
		f    func([]changeset.ChangeSet) (model.Model, []model.Problem)
	}{{viewAsIs, model.AsIs}, {viewToBe, model.ToBe}} {
		m, ps := v.f(sets)
		same(t, v.name+" nodes", storedNodes(t, x, v.name), m.Nodes())
		same(t, v.name+" links", storedLinks(t, x, v.name), m.Links())
		same(t, v.name+" problems", storedProblems(t, x, v.name), viewProblems(v.name, ps))
	}
}

func same[T any](t tb, what string, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s differ:\ngot  %+v\nwant %+v", what, got, want)
	}
}

// dump returns every row of every table, for comparing two indexes.
func dump(t tb, x *Index) map[string][]string {
	t.Helper()
	out := map[string][]string{}
	for _, table := range []string{"meta", "files", "entries", "change_sets", "nodes", "links", "problems"} {
		rows, err := x.db.Query(`SELECT * FROM ` + table)
		if err != nil {
			t.Fatalf("dump %s: %v", table, err)
		}
		cols, _ := rows.Columns()
		for rows.Next() {
			vals := make([]any, len(cols))
			ptrs := make([]any, len(cols))
			for i := range vals {
				ptrs[i] = &vals[i]
			}
			if err := rows.Scan(ptrs...); err != nil {
				t.Fatalf("dump %s: %v", table, err)
			}
			out[table] = append(out[table], fmt.Sprintf("%q", vals))
		}
		rows.Close()
		// Row order follows the primary key in every table.
	}
	return out
}

// userVersion reads PRAGMA user_version from the file at path.
func userVersion(t *testing.T, path string) int {
	t.Helper()
	db, err := sql.Open("sqlite", dsn(path))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var v int
	if err := db.QueryRow("PRAGMA user_version").Scan(&v); err != nil {
		t.Fatal(err)
	}
	return v
}

// snapshot hashes every file under dir by relative path.
func snapshot(t *testing.T, dir string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(dir, p)
		if d.IsDir() {
			out[rel+"/"] = "dir"
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(data)
		out[rel] = hex.EncodeToString(sum[:])
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// Shorthands for op fixtures.
func mkCap(n, node int, handle string, parent int, key string) op.Op {
	o := op.Op{ID: opID(n), Type: op.TypeNodeCreate, Node: nd(node), Kind: op.KindCapability, Handle: handle, Key: key}
	if parent != 0 {
		o.Parent = nd(parent)
	}
	return o
}
