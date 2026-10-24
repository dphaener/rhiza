package api

// This file holds the workspace: an open repository plus everything a write
// needs around it, namely the current change set (.genesis/local/change),
// change set refs, the interim model, actor resolution and the one helper
// that appends entries. Operations open one per call with openWorkspace and
// pass every input explicitly; nothing here reads the environment or runs
// git.

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/dphaener/genesis/internal/canon"
	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
	"github.com/dphaener/genesis/internal/store"
)

// The current change set file, relative to the .genesis folder.
const (
	localDir    = "local"
	currentFile = "change"
)

// findRoot walks up from dir to the first folder holding a .genesis entry
// and returns it. It reports false if no ancestor of dir has one.
func findRoot(dir string) (string, bool) {
	d, err := filepath.Abs(dir)
	if err != nil {
		return "", false
	}
	for {
		if _, err := os.Stat(filepath.Join(d, store.Dir)); err == nil {
			return d, true
		}
		parent := filepath.Dir(d)
		if parent == d {
			return "", false
		}
		d = parent
	}
}

// workspace is an open repository with the clock, ID generator and dry-run
// flag of one call.
//
// In a dry run nothing is written: appendEntry builds each entry and keeps
// it in pending instead of calling store.Write, and setCurrent and
// clearCurrent do nothing. Pending entries are part of every later read
// through the workspace (log, model, change set refs), so a sequence of
// dry-run writes chains exactly as the real one would.
type workspace struct {
	root    string
	st      *store.Store
	now     func() time.Time
	ids     *id.Generator
	dryRun  bool
	pending []entry.Entry
}

// openWorkspace opens the repository in root. No repository is not_found
// with next "genesis init"; a newer-format or invalid one is an error.
func openWorkspace(root string, now func() time.Time, ids *id.Generator, dryRun bool) (*workspace, error) {
	st, err := store.Open(root)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, NotFoundf("no Genesis repository in %s", root).WithNext("genesis init")
		}
		return nil, Errorf("%s", err.Error())
	}
	return &workspace{root: root, st: st, now: now, ids: ids, dryRun: dryRun}, nil
}

// openCtxWorkspace opens the workspace of ctx: its root, clock, IDs and
// dry-run flag.
func openCtxWorkspace(ctx *Ctx) (*workspace, error) {
	return openWorkspace(ctx.Root, ctx.Now, ctx.IDs, ctx.DryRun)
}

// log reads the log, plus the entries a dry run has built so far.
func (w *workspace) log() store.Log {
	l := w.st.ReadLog()
	if len(w.pending) > 0 {
		l.Entries = append(l.Entries, w.pending...)
		slices.SortFunc(l.Entries, func(a, b entry.Entry) int {
			return cmp.Or(id.Compare(a.ChangeSet, b.ChangeSet), id.Compare(a.ID, b.ID))
		})
	}
	return l
}

// model builds the interim model over the log with current as the current
// change set (zero for none).
func (w *workspace) model(current id.ID) Model {
	return newModel(w.log(), current)
}

// entriesOf returns the log's entries of change set cs.
func entriesOf(l store.Log, cs id.ID) []entry.Entry {
	var out []entry.Entry
	for _, e := range l.Entries {
		if e.ChangeSet == cs {
			out = append(out, e)
		}
	}
	return out
}

// changeSet folds change set cs. It reports false if the log has no entry
// of cs.
func (w *workspace) changeSet(cs id.ID) (changeset.ChangeSet, bool) {
	es := entriesOf(w.log(), cs)
	if len(es) == 0 {
		return changeset.ChangeSet{}, false
	}
	c, _ := changeset.Fold(cs, es)
	return c, true
}

func (w *workspace) currentPath() string {
	return filepath.Join(w.root, store.Dir, localDir, currentFile)
}

// current returns the current change set, zero if none is set. Content that
// is not one ULID (surrounding space aside) is an error.
func (w *workspace) current() (id.ID, error) {
	data, err := os.ReadFile(w.currentPath())
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return id.ID{}, nil
		}
		return id.ID{}, Errorf("read current change set: %v", err)
	}
	cs, err := id.Parse(string(bytes.TrimSpace(data)))
	if err != nil {
		return id.ID{}, Errorf("%s: not a change set ID: %v", w.currentPath(), err).
			WithNext("genesis change use <change set>")
	}
	return cs, nil
}

// setCurrent makes cs the current change set. The file is replaced
// atomically (a temp file in the same folder, then a rename) and local/ is
// created on demand. A dry run writes nothing.
func (w *workspace) setCurrent(cs id.ID) error {
	if w.dryRun {
		return nil
	}
	dir := filepath.Dir(w.currentPath())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Errorf("write current change set: %v", err)
	}
	f, err := os.CreateTemp(dir, "."+currentFile+".*.tmp")
	if err != nil {
		return Errorf("write current change set: %v", err)
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	_, err = f.WriteString(cs.String() + "\n")
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Chmod(tmp, 0o644)
	}
	if err == nil {
		err = os.Rename(tmp, w.currentPath())
	}
	if err != nil {
		return Errorf("write current change set: %v", err)
	}
	return nil
}

// clearCurrent unsets the current change set; none set is fine. A dry run
// writes nothing.
func (w *workspace) clearCurrent() error {
	if w.dryRun {
		return nil
	}
	if err := os.Remove(w.currentPath()); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return Errorf("clear current change set: %v", err)
	}
	return nil
}

// minPrefix is the shortest ID prefix a ref may be.
const minPrefix = 4

// resolveChangeSet resolves a change set ref, a full ID or a unique prefix
// of 4+ characters (any case), over the change sets in the log. No match is
// not_found; several are a usage error listing them.
func (w *workspace) resolveChangeSet(ref string) (id.ID, error) {
	ref = strings.TrimSpace(ref)
	var sets []id.ID
	for _, e := range w.log().Entries {
		if len(sets) == 0 || sets[len(sets)-1] != e.ChangeSet {
			sets = append(sets, e.ChangeSet)
		}
	}
	if full, err := id.Parse(ref); err == nil {
		if slices.Contains(sets, full) {
			return full, nil
		}
		return id.ID{}, NotFoundf("no change set %s", full).WithNext("genesis change ls")
	}
	if len(ref) < minPrefix {
		return id.ID{}, Usagef("change set ref %q is too short: give a full ID or %d+ characters", ref, minPrefix)
	}
	prefix := strings.ToUpper(ref)
	var matches []string
	var found id.ID
	for _, cs := range sets {
		if strings.HasPrefix(cs.String(), prefix) {
			matches = append(matches, cs.String())
			found = cs
		}
	}
	switch len(matches) {
	case 0:
		return id.ID{}, NotFoundf("no change set matches %q", ref).WithNext("genesis change ls")
	case 1:
		return found, nil
	}
	return id.ID{}, Usagef("change set ref %q matches several change sets: %s", ref, strings.Join(matches, ", "))
}

// effectiveChangeSet is the change set a call acts on: ref resolved if
// given, else the current change set, else zero.
func (w *workspace) effectiveChangeSet(ref string) (id.ID, error) {
	if strings.TrimSpace(ref) != "" {
		return w.resolveChangeSet(ref)
	}
	return w.current()
}

// newEntryNext is the next command for a write that needs a change set.
const newEntryNext = `genesis change new "<intent>"`

// appendEntry appends body to change set cs as actor for forID and returns
// the entry. It reads the log, folds cs and sets after to its heads; an open
// starts a new change set, so cs must have no entries yet and after is
// empty. At is the clock in UTC, truncated to the second.
//
// A non-open body needs cs in the log (else not_found), and anything but a
// note into a terminal change set is a conflict. Whether a lifecycle move is
// legal in a non-terminal state is the caller's check.
func (w *workspace) appendEntry(cs id.ID, body changeset.Body, actor, forID id.ID) (entry.Entry, error) {
	es := entriesOf(w.log(), cs)
	var after []id.ID
	if body.Type() == entry.TypeOpen {
		if len(es) > 0 {
			return entry.Entry{}, Errorf("change set %s already has entries", cs)
		}
	} else {
		if len(es) == 0 {
			return entry.Entry{}, NotFoundf("no change set %s", cs).WithNext("genesis change ls")
		}
		c, _ := changeset.Fold(cs, es)
		if c.State.Terminal() && body.Type() != entry.TypeNote {
			return entry.Entry{}, Conflictf("change set %s is %s", cs, c.State).WithNext(newEntryNext)
		}
		after = c.Heads
	}
	raw, err := changeset.EncodeBody(body)
	if err != nil {
		return entry.Entry{}, Usagef("%s", err.Error())
	}
	eid, err := w.ids.New()
	if err != nil {
		return entry.Entry{}, Errorf("%s", err.Error())
	}
	e := entry.Entry{
		V:         entry.FormatVersion,
		ID:        eid,
		ChangeSet: cs,
		Actor:     actor,
		For:       forID,
		At:        w.now().UTC().Truncate(time.Second),
		After:     after,
		Type:      body.Type(),
		Body:      raw,
	}
	if w.dryRun {
		w.pending = append(w.pending, e)
		return e, nil
	}
	if err := w.st.Write(e); err != nil {
		return entry.Entry{}, Errorf("%s", err.Error())
	}
	return e, nil
}

// openChangeSet starts a new change set with a fresh ID and returns that ID
// and its open entry.
func (w *workspace) openChangeSet(open changeset.Open, actor, forID id.ID) (id.ID, entry.Entry, error) {
	cs, err := w.ids.New()
	if err != nil {
		return id.ID{}, entry.Entry{}, Errorf("%s", err.Error())
	}
	e, err := w.appendEntry(cs, open, actor, forID)
	if err != nil {
		return id.ID{}, entry.Entry{}, err
	}
	return cs, e, nil
}

// writeRecord writes ops as one landed record change set: open (kind
// record, the given intent), one ops entry and a land, chained in that
// order. It returns the change set ID and the three entries. It never
// touches the current change set.
func (w *workspace) writeRecord(intent string, ops []op.Op, actor, forID id.ID) (id.ID, []entry.Entry, error) {
	cs, open, err := w.openChangeSet(changeset.Open{Kind: changeset.KindRecord, Intent: intent}, actor, forID)
	if err != nil {
		return id.ID{}, nil, err
	}
	out := []entry.Entry{open}
	for _, b := range []changeset.Body{changeset.Ops{Ops: ops}, changeset.Land{}} {
		e, err := w.appendEntry(cs, b, actor, forID)
		if err != nil {
			return id.ID{}, nil, err
		}
		out = append(out, e)
	}
	return cs, out, nil
}

// resolveActor returns who acts and for whom (D5). current is the change
// set the model sees as current (zero for none); actorRef and forRef are
// the raw --actor/--for refs, empty if not given; gitConfig looks up git
// config keys.
//
// The actor is actorRef resolved (it must be an actor node), else the actor
// node whose email is git's user.email, else a new actor node recorded in
// its own landed record change set, written by itself. In a dry run with no
// such node nothing is recorded and the actor is zero: dry runs only print
// operations. With neither actorRef nor a git email it is a usage error.
// For is forRef resolved (also an actor node), else the actor.
func (w *workspace) resolveActor(current id.ID, actorRef, forRef string, gitConfig func(key string) (string, error)) (actor, forID id.ID, err error) {
	m := w.model(current)
	if actorRef != "" {
		if actor, err = resolveActorRef(m, actorRef); err != nil {
			return id.ID{}, id.ID{}, err
		}
	} else {
		email := gitValue(gitConfig, "user.email")
		if email == "" {
			return id.ID{}, id.ID{}, Usagef("no actor: set git user.email or pass --actor").
				WithNext("git config user.email you@example.com")
		}
		if n, ok := m.ActorByEmail(email); ok {
			actor = n.ID
		} else if w.dryRun {
			return id.ID{}, id.ID{}, nil
		} else {
			if actor, err = w.recordActor(m, email, gitValue(gitConfig, "user.name")); err != nil {
				return id.ID{}, id.ID{}, err
			}
			m = w.model(current)
		}
	}
	if forRef == "" {
		return actor, actor, nil
	}
	if forID, err = resolveActorRef(m, forRef); err != nil {
		return id.ID{}, id.ID{}, err
	}
	return actor, forID, nil
}

// gitValue returns git config key trimmed, or "" if it is unset or fails.
func gitValue(gitConfig func(string) (string, error), key string) string {
	if gitConfig == nil {
		return ""
	}
	v, err := gitConfig(key)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(v)
}

// resolveActorRef resolves ref and requires an actor node.
func resolveActorRef(m Model, ref string) (id.ID, error) {
	n, err := m.Resolve(ref)
	if err != nil {
		return id.ID{}, err
	}
	if n.Kind != op.KindActor {
		return id.ID{}, Usagef("%s is a %s, not an actor", ref, n.Kind)
	}
	return n.ID, nil
}

// recordActor records a new actor node for email (and name, if set) in a
// landed record change set whose entries it authors itself.
func (w *workspace) recordActor(m Model, email, name string) (id.ID, error) {
	node, err := w.ids.New()
	if err != nil {
		return id.ID{}, Errorf("%s", err.Error())
	}
	opID, err := w.ids.New()
	if err != nil {
		return id.ID{}, Errorf("%s", err.Error())
	}
	fields := map[string]string{"email": email}
	if name != "" {
		fields["name"] = name
	}
	raw, err := textFields(fields)
	if err != nil {
		return id.ID{}, err
	}
	local, _, _ := strings.Cut(email, "@")
	handle := uniqueHandle(handleSlug(local), func(h string) bool {
		_, taken := m.HandleTaken(h)
		return taken
	})
	c := op.Op{ID: opID, Type: op.TypeNodeCreate, Node: node, Kind: op.KindActor, Handle: handle, Fields: raw}
	if _, _, err := w.writeRecord("Record actor "+email, []op.Op{c}, node, node); err != nil {
		return id.ID{}, err
	}
	return node, nil
}

// textFields encodes string field values as canonical JSON strings.
func textFields(m map[string]string) (map[string]json.RawMessage, error) {
	out := make(map[string]json.RawMessage, len(m))
	for k, v := range m {
		b, err := canon.Marshal(v)
		if err != nil {
			return nil, Errorf("field %s: %v", k, err)
		}
		out[k] = bytes.TrimSuffix(b, []byte("\n"))
	}
	return out, nil
}

// maxHandle is the longest handle op accepts.
const maxHandle = 64

// handleSlug turns s into a handle: lower case, each run of characters
// outside [a-z0-9] one '-', no '-' at either end, at most 64 bytes. An empty
// result is "actor".
func handleSlug(s string) string {
	var sb strings.Builder
	dash := false
	for _, c := range []byte(strings.ToLower(s)) {
		if ('a' <= c && c <= 'z') || ('0' <= c && c <= '9') {
			if dash && sb.Len() > 0 {
				sb.WriteByte('-')
			}
			dash = false
			sb.WriteByte(c)
		} else {
			dash = true
		}
	}
	h := sb.String()
	if len(h) > maxHandle {
		h = strings.TrimRight(h[:maxHandle], "-")
	}
	if h == "" {
		return "actor"
	}
	return h
}

// uniqueHandle returns base if it is free, else the first free of base-2,
// base-3 and so on, shortening base so the result stays within 64 bytes.
func uniqueHandle(base string, taken func(string) bool) string {
	if !taken(base) {
		return base
	}
	for n := 2; ; n++ {
		suffix := "-" + strconv.Itoa(n)
		b := base
		if len(b)+len(suffix) > maxHandle {
			b = strings.TrimRight(b[:maxHandle-len(suffix)], "-")
		}
		if h := b + suffix; !taken(h) {
			return h
		}
	}
}
