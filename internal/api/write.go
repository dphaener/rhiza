package api

// This file holds the pipeline every node and link write shares (add, set,
// mv, retire, link, unlink): pick the change set, build the operations
// against the model, resolve the actor, and write one ops entry, or, with
// --record, one landed record change set. Building is separate from
// writing: each verb's builder is a function of a writeBase and its request
// that returns an edit, so apply can run the same builders one after
// another.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// writeBase is what a builder builds against: the to-be model, the current
// ops of the change set being written (nil for a new record change set),
// and the ID generator.
type writeBase struct {
	// root is the project root anchors resolve against.
	root string
	m    Model
	ops  []op.Op
	ids  *id.Generator
	// changeSet resolves a change set ref over the log, for links whose
	// end is a change set; nil where there is none.
	changeSet func(ref string) (id.ID, error)
	// changeSets folds every change set in the log, the change set being
	// written included; nil where there is none.
	changeSets func() []changeset.ChangeSet
}

// newID draws a fresh ID.
func (b writeBase) newID() (id.ID, error) {
	v, err := b.ids.New()
	if err != nil {
		return id.ID{}, Errorf("%s", err.Error())
	}
	return v, nil
}

// pendingOp returns the change set's op of type t among the op IDs pending
// (a Node's or Link's Pending).
func (b writeBase) pendingOp(pending []id.ID, t op.Type) (op.Op, bool) {
	for _, o := range b.ops {
		if o.Type == t && slices.Contains(pending, o.ID) {
			return o, true
		}
	}
	return op.Op{}, false
}

// edit is what one write does to its change set: new ops and revisions of
// existing ones (a revision repeats the op's ID), op IDs to drop, and the
// nodes and links it touches, with their handles for human output.
type edit struct {
	ops     []op.Op
	drop    []id.ID
	nodes   []id.ID
	links   []id.ID
	handles map[id.ID]string
}

// touchNode notes that the write touches node n, whose handle is handle.
func (e *edit) touchNode(n id.ID, handle string) {
	if !slices.Contains(e.nodes, n) {
		e.nodes = append(e.nodes, n)
	}
	if e.handles == nil {
		e.handles = map[id.ID]string{}
	}
	e.handles[n] = handle
}

// touchLink notes that the write touches link l.
func (e *edit) touchLink(l id.ID) {
	if !slices.Contains(e.links, l) {
		e.links = append(e.links, l)
	}
}

// check validates every op; a broken one is a usage error, since only the
// request can have made it so.
func (e edit) check(opName string) error {
	for _, o := range e.ops {
		if err := o.Validate(); err != nil {
			return Usagef("%s", err.Error()).WithNext("rhiza describe " + opName)
		}
	}
	return nil
}

// WriteResponse is what every node and link write returns.
type WriteResponse struct {
	Entry     id.ID             `json:"entry,omitzero" desc:"The ops entry written; absent in a dry run"`
	ChangeSet id.ID             `json:"change_set" desc:"The change set written to (with --record, the new record change set)"`
	Nodes     []id.ID           `json:"nodes" desc:"The nodes the write created or changed, sorted; link endpoints are not included"`
	Links     []id.ID           `json:"links,omitempty" desc:"The links the write created or retired, sorted"`
	Record    bool              `json:"record,omitempty" desc:"True if the write opened and landed a record change set"`
	DryRun    bool              `json:"dry_run,omitempty" desc:"True if nothing was written"`
	Ops       []json.RawMessage `json:"ops,omitempty" desc:"Dry run only: the ops the entry would carry, as canonical op JSON"`
	Drop      []id.ID           `json:"drop,omitempty" desc:"Dry run only: the op IDs the entry would drop"`

	handles map[id.ID]string
}

// Human reports the entry and what it touched, or, in a dry run, the ops.
func (r WriteResponse) Human(w io.Writer) {
	switch {
	case r.DryRun:
		fmt.Fprintf(w, "Dry run: would write %d operation(s)", len(r.Ops))
		if len(r.Drop) > 0 {
			fmt.Fprintf(w, " and drop %d", len(r.Drop))
		}
		fmt.Fprintf(w, " to change set %s; nothing was written.\n\n", r.ChangeSet)
		for _, o := range r.Ops {
			var buf bytes.Buffer
			if json.Compact(&buf, o) != nil {
				buf.Reset()
				buf.Write(o)
			}
			fmt.Fprintf(w, "- %s\n", buf.String())
		}
		for _, x := range r.Drop {
			fmt.Fprintf(w, "- drop %s\n", x)
		}
		return
	case r.Record:
		fmt.Fprintf(w, "Recorded and landed change set %s (entry %s).\n", r.ChangeSet, r.Entry)
	default:
		fmt.Fprintf(w, "Wrote entry %s to change set %s.\n", r.Entry, r.ChangeSet)
	}
	if len(r.Nodes)+len(r.Links) > 0 {
		fmt.Fprintln(w)
	}
	for _, n := range r.Nodes {
		if h := r.handles[n]; h != "" {
			fmt.Fprintf(w, "- node %s (%s)\n", n, h)
		} else {
			fmt.Fprintf(w, "- node %s\n", n)
		}
	}
	for _, l := range r.Links {
		fmt.Fprintf(w, "- link %s\n", l)
	}
}

// runWrite is the shared pipeline. intent is the command line, used as a
// record change set's intent; build makes the edit against a writeBase.
//
// Without --record the change set is ctx.ChangeSet, else ctx.ChangeRef
// resolved, else the current one; there must be one, and it must not be
// terminal. With --record the write sees only landed state and goes into a
// new record change set; ctx's change set choice and the current change
// set are left alone.
//
// The edit is built and checked before the actor is resolved, so a bad
// request never records an actor; if resolving records one, the edit is
// built again over the new log.
func runWrite(ctx *Ctx, opName, intent string, build func(writeBase) (edit, error)) (any, error) {
	return runWriteOn(ctx, opName, intent, func(w *workspace, cs id.ID) (edit, error) {
		b := writeBase{root: w.root, m: w.model(cs), ids: ctx.IDs, changeSet: w.resolveChangeSet, changeSets: w.changeSets}
		if !cs.IsZero() {
			c, _ := w.changeSet(cs)
			b.ops = c.Ops
		}
		return build(b)
	})
}

// runWriteOn is runWrite with a builder that gets the workspace and the
// target change set (zero with --record) instead of a writeBase, for writes
// such as apply that build in several steps.
func runWriteOn(ctx *Ctx, opName, intent string, build func(w *workspace, cs id.ID) (edit, error)) (any, error) {
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	var cs id.ID
	if !ctx.Record {
		if cs, err = writeTarget(w, ctx); err != nil {
			return nil, err
		}
	}
	buildEdit := func() (edit, error) {
		e, err := build(w, cs)
		if err != nil {
			return edit{}, err
		}
		return e, e.check(opName)
	}
	e, err := buildEdit()
	if err != nil {
		return nil, err
	}
	before := len(w.log().Entries)
	actor, forID, err := w.resolveActor(cs, ctx)
	if err != nil {
		return nil, err
	}
	if len(w.log().Entries) != before {
		if e, err = buildEdit(); err != nil {
			return nil, err
		}
	}

	slices.SortFunc(e.drop, id.Compare)
	e.drop = slices.Compact(e.drop)
	resp := WriteResponse{Record: ctx.Record, DryRun: ctx.DryRun, handles: e.handles}
	if ctx.Record {
		rec, entries, err := w.writeRecord(intent, e.ops, actor, forID)
		if err != nil {
			return nil, err
		}
		resp.ChangeSet, resp.Entry = rec, entries[1].ID
	} else {
		en, err := w.appendEntry(cs, changeset.Ops{Ops: e.ops, Drop: e.drop}, actor, forID)
		if err != nil {
			return nil, err
		}
		resp.ChangeSet, resp.Entry = cs, en.ID
	}
	resp.Nodes = orEmpty(slices.SortedFunc(slices.Values(e.nodes), id.Compare))
	if len(e.links) > 0 {
		resp.Links = slices.SortedFunc(slices.Values(e.links), id.Compare)
	}
	if ctx.DryRun {
		resp.Entry = id.ID{}
		for _, o := range e.ops {
			raw, err := op.Encode(o)
			if err != nil {
				return nil, Usagef("%s", err.Error())
			}
			resp.Ops = append(resp.Ops, raw)
		}
		resp.Drop = e.drop
	}
	return resp, nil
}

// writeTarget is the change set a write without --record goes into, chosen
// by effectiveChangeSet. None is a usage error, a terminal one a conflict.
func writeTarget(w *workspace, ctx *Ctx) (id.ID, error) {
	cs, err := w.effectiveChangeSet(ctx)
	if err != nil {
		return id.ID{}, err
	}
	if cs.IsZero() {
		return id.ID{}, Usagef("no change set: open one with rhiza change new, or add --record to write and land this in one step").
			WithNext(newEntryNext)
	}
	c, ok := w.changeSet(cs)
	if !ok {
		return id.ID{}, NotFoundf("no change set %s", cs).WithNext(changeLsNext)
	}
	if c.State.Terminal() {
		next := newEntryNext
		if cur, err := w.current(); err == nil && cur == cs {
			next = clearCurrentNext
		}
		return id.ID{}, Conflictf("change set %s is %s; nothing can be written to it", cs, c.State).
			WithNext(next)
	}
	return cs, nil
}

// commandLine joins a command and its arguments into a record change set's
// intent, quoting arguments with spaces.
func commandLine(parts ...string) string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p == "" {
			continue
		}
		if strings.ContainsAny(p, " \t\n\"") {
			p = fmt.Sprintf("%q", p)
		}
		out = append(out, p)
	}
	return strings.Join(out, " ")
}

// parseFields turns --field name=text and --json-field name=<json> into
// canonical field values. A name given twice, in either or both, a value
// without "=", and invalid JSON are usage errors; names are checked by op
// validation later. Nil if there are none.
func parseFields(text, jsonFields []string) (map[string]json.RawMessage, error) {
	if len(text)+len(jsonFields) == 0 {
		return nil, nil
	}
	out := map[string]json.RawMessage{}
	add := func(flag, kv string, value func(string) (json.RawMessage, error)) error {
		name, v, ok := strings.Cut(kv, "=")
		if !ok {
			return Usagef("--%s %q: want name=value", flag, kv)
		}
		if _, dup := out[name]; dup {
			return Usagef("field %q is given more than once", name)
		}
		raw, err := value(v)
		if err != nil {
			return err
		}
		out[name] = raw
		return nil
	}
	for _, kv := range text {
		err := add("field", kv, func(v string) (json.RawMessage, error) {
			m, err := textFields(map[string]string{"x": v})
			return m["x"], err
		})
		if err != nil {
			return nil, err
		}
	}
	for _, kv := range jsonFields {
		err := add("json-field", kv, func(v string) (json.RawMessage, error) {
			b, err := canon.Canonicalize([]byte(v))
			if err != nil {
				return nil, Usagef("--json-field %s: %v", kv, err)
			}
			return json.RawMessage(bytes.TrimSuffix(b, []byte("\n"))), nil
		})
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// mergeFields returns base with set applied and the unset names removed,
// as a new map; nil if empty.
func mergeFields(base, set map[string]json.RawMessage, unset []string) map[string]json.RawMessage {
	out := maps.Clone(base)
	if out == nil {
		out = map[string]json.RawMessage{}
	}
	maps.Copy(out, set)
	for _, k := range unset {
		delete(out, k)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
