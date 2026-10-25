package api

// This file holds apply (D10): many node and link writes from one JSON file,
// applied in order and written as one ops entry, or nothing at all.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"

	"github.com/dphaener/genesis/internal/canon"
	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
)

func init() {
	register(Operation{
		Name:     "apply",
		Summary:  `Apply many writes as one entry: a JSON file {"writes":[{"op":"add",...},...]} whose items are add, set, mv, retire, link or unlink requests plus "op"`,
		Write:    true,
		Record:   true,
		Request:  reflect.TypeFor[ApplyRequest](),
		Response: reflect.TypeFor[WriteResponse](),
		Input: &Input{
			Summary: `{"writes":[...]}, at least one item; each item is one of add | set | mv | retire | link | unlink: that operation's request plus "op" naming it`,
			Schema:  applyFileSchema,
		},
		Run: runApply,
	})
}

// ApplyRequest names the file of writes.
type ApplyRequest struct {
	File string `json:"file" arg:"0" desc:"A JSON file of writes, or - for stdin (see the input schema); refs may name nodes added by earlier items"`
}

// applyItem is one decoded write of the file.
type applyItem struct {
	op    string
	build func(writeBase) (edit, error)
}

// applyWrites lists, in order, each op an item may name with its builder.
// It is the one list both decoding (applyBuilders) and the published file
// schema (applyFileSchema) use.
var applyWrites = []struct {
	op    string
	build func(writeBase, any) (edit, error)
}{
	{"add", func(b writeBase, r any) (edit, error) { return buildAdd(b, r.(AddRequest)) }},
	{"set", func(b writeBase, r any) (edit, error) { return buildSet(b, r.(SetRequest)) }},
	{"mv", func(b writeBase, r any) (edit, error) { return buildMv(b, r.(MvRequest)) }},
	{"retire", func(b writeBase, r any) (edit, error) { return buildRetire(b, r.(RetireRequest)) }},
	{"link", func(b writeBase, r any) (edit, error) { return buildLink(b, r.(LinkRequest)) }},
	{"unlink", func(b writeBase, r any) (edit, error) { return buildUnlink(b, r.(UnlinkRequest)) }},
}

// applyBuilders maps each op an item may name to its builder.
var applyBuilders = func() map[string]func(writeBase, any) (edit, error) {
	m := map[string]func(writeBase, any) (edit, error){}
	for _, w := range applyWrites {
		m[w.op] = w.build
	}
	return m
}()

// applyFileSchema returns the JSON Schema of an apply file: an object whose
// "writes" is a non-empty array, each item one of the applyWrites ops'
// request schemas with a required "op" fixed to that op's name.
func applyFileSchema() map[string]any {
	branches := make([]any, 0, len(applyWrites))
	for _, w := range applyWrites {
		o, ok := Lookup(w.op)
		if !ok {
			panic("api: apply: operation " + w.op + " is not registered")
		}
		s := Schema(o.Request)
		props := s["properties"].(map[string]any)
		props["op"] = map[string]any{"const": w.op, "description": "The operation this item is a request of"}
		required := append(s["required"].([]string), "op")
		slices.Sort(required)
		s["required"] = required
		branches = append(branches, s)
	}
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"writes": map[string]any{
				"type":        "array",
				"minItems":    1,
				"items":       map[string]any{"oneOf": branches},
				"description": "The writes, applied in order",
			},
		},
		"required":             []string{"writes"},
		"additionalProperties": false,
	}
}

// runApply reads and decodes the file, then writes the net edit of all its
// items through the shared write pipeline.
func runApply(ctx *Ctx, req any) (any, error) {
	r := req.(ApplyRequest)
	data, err := readApplyFile(ctx, r.File)
	if err != nil {
		return nil, err
	}
	items, err := decodeWrites(data)
	if err != nil {
		return nil, err
	}
	return runWriteOn(ctx, "apply", commandLine("apply", r.File), func(w *workspace, cs id.ID) (edit, error) {
		return buildApply(w, cs, ctx.IDs, items)
	})
}

// readApplyFile reads file, or stdin for "-".
func readApplyFile(ctx *Ctx, file string) ([]byte, error) {
	switch {
	case file == "":
		return nil, Usagef("apply needs a file, or - for stdin").WithNext("genesis describe apply")
	case file == "-":
		if ctx.Stdin == nil {
			return nil, Usagef("apply -: no stdin")
		}
		data, err := io.ReadAll(ctx.Stdin)
		if err != nil {
			return nil, Usagef("apply -: cannot read stdin: %v", err)
		}
		return data, nil
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, Usagef("apply: cannot read %s: %v", file, errors.Unwrap(err))
	}
	return data, nil
}

// decodeWrites decodes the file strictly: exactly {"writes":[...]} with at
// least one item, each an object with a known "op" whose other keys decode
// strictly into that op's request.
func decodeWrites(data []byte) ([]applyItem, error) {
	trimmed := bytes.TrimSpace(data)
	if _, err := canon.Canonicalize(trimmed); err != nil {
		return nil, Usagef("apply file: %v", err).WithNext("genesis describe apply")
	}
	var file struct {
		Writes []json.RawMessage `json:"writes"`
	}
	dec := json.NewDecoder(bytes.NewReader(trimmed))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&file); err != nil {
		return nil, Usagef("apply file: want {\"writes\":[...]}: %v", err).WithNext("genesis describe apply")
	}
	if len(file.Writes) == 0 {
		return nil, Usagef("apply file: no writes").WithNext("genesis describe apply")
	}
	items := make([]applyItem, 0, len(file.Writes))
	for i, raw := range file.Writes {
		it, err := decodeWrite(raw)
		if err != nil {
			return nil, itemError(i, err)
		}
		items = append(items, it)
	}
	return items, nil
}

// decodeWrite decodes one item.
func decodeWrite(raw json.RawMessage) (applyItem, error) {
	var fields map[string]json.RawMessage
	if !bytes.HasPrefix(bytes.TrimSpace(raw), []byte("{")) || json.Unmarshal(raw, &fields) != nil {
		return applyItem{}, Usagef("not a JSON object").WithNext("genesis describe apply")
	}
	var name string
	if err := json.Unmarshal(fields["op"], &name); err != nil || fields["op"] == nil {
		return applyItem{}, Usagef(`want a string "op": one of add, set, mv, retire, link, unlink`).WithNext("genesis describe apply")
	}
	build, ok := applyBuilders[name]
	if !ok {
		return applyItem{}, Usagef(`unknown op %q: want one of add, set, mv, retire, link, unlink`, name).WithNext("genesis describe apply")
	}
	o, ok := Lookup(name)
	if !ok {
		return applyItem{}, Errorf("operation %s is not registered", name)
	}
	delete(fields, "op")
	rest, err := json.Marshal(fields)
	if err != nil {
		return applyItem{}, Usagef("%v", err)
	}
	req, err := DecodeRequest(o, rest)
	if err != nil {
		return applyItem{}, AsError(err).WithNext("genesis describe " + name)
	}
	return applyItem{op: name, build: func(b writeBase) (edit, error) { return build(b, req) }}, nil
}

// itemError prefixes err's message with the item's index, keeping its code
// and next.
func itemError(i int, err error) error {
	e := AsError(err)
	return &Error{Code: e.Code, Message: fmt.Sprintf("writes[%d]: %s", i, e.Message), Next: e.Next}
}

// buildApply applies the items in order over a dry-run copy of w, each one
// as an ops entry of its own in the target change set (with --record, a
// draft record change set made up for the purpose), so every item builds
// against the state the earlier ones left. The result is the net edit
// between the target's ops before and after: new or changed ops, and drops
// of ops it had before and no longer has. Ops added and dropped within the
// file vanish.
func buildApply(w *workspace, cs id.ID, ids *id.Generator, items []applyItem) (edit, error) {
	s := *w
	s.dryRun = true
	s.pending = slices.Clone(w.pending)
	target := cs
	var start []op.Op
	if cs.IsZero() {
		var err error
		if target, _, err = s.openChangeSet(changeset.Open{Kind: changeset.KindRecord, Intent: "apply"}, id.ID{}, id.ID{}); err != nil {
			return edit{}, err
		}
	} else {
		var err error
		if start, err = scratchOps(&s, target); err != nil {
			return edit{}, err
		}
	}

	var touched edit
	ops := start
	for i, it := range items {
		e, err := it.build(writeBase{m: s.model(target), ops: ops, ids: ids})
		if err == nil {
			err = e.check(it.op)
		}
		if err != nil {
			return edit{}, itemError(i, err)
		}
		if len(e.ops)+len(e.drop) > 0 {
			slices.SortFunc(e.drop, id.Compare)
			body := changeset.Ops{Ops: e.ops, Drop: slices.Compact(e.drop)}
			if _, err := s.appendEntry(target, body, id.ID{}, id.ID{}); err != nil {
				return edit{}, itemError(i, err)
			}
			if ops, err = scratchOps(&s, target); err != nil {
				return edit{}, itemError(i, err)
			}
		}
		for _, n := range e.nodes {
			touched.touchNode(n, e.handles[n])
		}
		for _, l := range e.links {
			touched.touchLink(l)
		}
	}

	var net edit
	before := map[id.ID]json.RawMessage{}
	for _, o := range start {
		raw, err := op.Encode(o)
		if err != nil {
			return edit{}, Usagef("%s", err.Error())
		}
		before[o.ID] = raw
	}
	after := map[id.ID]bool{}
	targets := map[id.ID]bool{}
	for _, o := range ops {
		after[o.ID] = true
		raw, err := op.Encode(o)
		if err != nil {
			return edit{}, Usagef("%s", err.Error())
		}
		if old, ok := before[o.ID]; !ok || !bytes.Equal(old, raw) {
			net.ops = append(net.ops, o)
			targets[o.Target()] = true
		}
	}
	for _, o := range start {
		if !after[o.ID] {
			net.drop = append(net.drop, o.ID)
			targets[o.Target()] = true
		}
	}
	if len(net.ops)+len(net.drop) == 0 {
		return edit{}, Usagef("nothing to write: the file's writes cancel out").WithNext("genesis change show")
	}
	for _, n := range touched.nodes {
		if targets[n] {
			net.touchNode(n, touched.handles[n])
		}
	}
	for _, l := range touched.links {
		if targets[l] {
			net.touchLink(l)
		}
	}
	return net, nil
}

// scratchOps folds change set cs of the scratch workspace; a fold problem
// means an item could not be applied as built.
func scratchOps(s *workspace, cs id.ID) ([]op.Op, error) {
	c, problems := changeset.Fold(cs, entriesOf(s.log(), cs))
	if len(problems) > 0 {
		return nil, Errorf("%s", problems[0].Err.Error())
	}
	return slices.Clone(c.Ops), nil
}
