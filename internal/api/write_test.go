package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/check"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// writeCtx is an initialized repository with a current intention change
// set, and that change set's ID.
func writeCtx(t *testing.T) (*Ctx, id.ID) {
	t.Helper()
	ctx := initCtx(t)
	cs := newChange(t, ctx, `{"intent":"Shape"}`).ChangeSet
	return ctx, cs
}

// mustWrite runs a write operation and fails the test on an error.
func mustWrite(t *testing.T, ctx *Ctx, name, req string) WriteResponse {
	t.Helper()
	resp, err := callOp(t, ctx, name, req)
	if err != nil {
		t.Fatalf("%s %s: %v", name, req, err)
	}
	return resp.(WriteResponse)
}

// recordWrite runs a write with --record.
func recordWrite(t *testing.T, ctx *Ctx, name, req string) WriteResponse {
	t.Helper()
	c := *ctx
	c.Record = true
	return mustWrite(t, &c, name, req)
}

// nodeAt resolves ref in the to-be view of the current change set.
func nodeAt(t *testing.T, ctx *Ctx, ref string) Node {
	t.Helper()
	n, err := ctxWorkspace(t, ctx).model(currentOf(t, ctx)).Resolve(ref)
	if err != nil {
		t.Fatalf("resolve %s: %v", ref, err)
	}
	return n
}

// opsOf returns change set cs's current ops, failing on fold problems.
func opsOf(t *testing.T, ctx *Ctx, cs id.ID) []op.Op {
	t.Helper()
	mustFoldClean(t, ctxWorkspace(t, ctx))
	c, ok := ctxWorkspace(t, ctx).changeSet(cs)
	if !ok {
		t.Fatalf("no change set %s", cs)
	}
	return c.Ops
}

func onlyOp(t *testing.T, ops []op.Op) op.Op {
	t.Helper()
	if len(ops) != 1 {
		t.Fatalf("got %d ops, want 1: %+v", len(ops), ops)
	}
	return ops[0]
}

func raw(s string) json.RawMessage { return json.RawMessage(s) }

func TestWriteCommon(t *testing.T) {
	t.Run("ResponseAndHuman", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		r := mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		n := nodeAt(t, ctx, "store")
		if r.Entry.IsZero() || r.ChangeSet != cs || !slices.Equal(r.Nodes, []id.ID{n.ID}) ||
			r.Links != nil || r.DryRun || r.Record || r.Ops != nil {
			t.Fatalf("response = %+v", r)
		}
		out := human(t, r)
		for _, s := range []string{r.Entry.String(), cs.String(), n.ID.String(), "store"} {
			if !strings.Contains(out, s) {
				t.Errorf("human lacks %q:\n%s", s, out)
			}
		}
	})

	t.Run("NoChangeSet", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		if _, err := callOp(t, ctx, "change.use", `{"none":true}`); err != nil {
			t.Fatal(err)
		}
		before := len(logFiles(t, ctx.Root))
		_, err := callOp(t, ctx, "add", `{"kind":"capability","handle":"x"}`)
		wantExit(t, err, 2)
		if e := AsError(err); e.Next != newEntryNext || !strings.Contains(e.Message, "--record") {
			t.Errorf("error = %+v", e)
		}
		if n := len(logFiles(t, ctx.Root)); n != before {
			t.Errorf("wrote %d files", n-before)
		}
	})

	t.Run("TerminalChangeSet", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		landed := recordWrite(t, ctx, "add", `{"kind":"capability","handle":"a"}`).ChangeSet
		c := *ctx
		c.ChangeRef = landed.String()
		_, err := callOp(t, &c, "add", `{"kind":"capability","handle":"b"}`)
		wantExit(t, err, 3)
		if e := AsError(err); e.Next != newEntryNext || !strings.Contains(e.Message, "landed") {
			t.Errorf("error = %+v", e)
		}
	})

	t.Run("ChangeRef", func(t *testing.T) {
		ctx, first := writeCtx(t)
		second := newChange(t, ctx, `{"intent":"Other"}`).ChangeSet
		c := *ctx
		c.ChangeRef = first.String()
		r := mustWrite(t, &c, "add", `{"kind":"capability","handle":"a"}`)
		if r.ChangeSet != first || len(opsOf(t, ctx, first)) != 1 || len(opsOf(t, ctx, second)) != 0 {
			t.Errorf("wrote into %s", r.ChangeSet)
		}
	})

	t.Run("DryRun", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		before := logFiles(t, ctx.Root)
		c := *ctx
		c.DryRun = true
		resp, err := callOp(t, &c, "add", `{"kind":"capability","handle":"store","fields":["statement=Hi"]}`)
		if err != nil {
			t.Fatal(err)
		}
		r := resp.(WriteResponse)
		if !r.DryRun || !r.Entry.IsZero() || r.ChangeSet != cs || len(r.Ops) != 1 || len(r.Nodes) != 1 {
			t.Fatalf("response = %+v", r)
		}
		o, err := op.Decode(r.Ops[0])
		if err != nil || o.Type != op.TypeNodeCreate || o.Handle != "store" || o.Node != r.Nodes[0] {
			t.Errorf("op = %+v, %v", o, err)
		}
		if after := logFiles(t, ctx.Root); !slices.Equal(before, after) {
			t.Errorf("dry run wrote files: %v", after)
		}
		out := human(t, r)
		for _, s := range []string{"Dry run", "1 operation", `"handle":"store"`, "nothing was written"} {
			if !strings.Contains(out, s) {
				t.Errorf("human lacks %q:\n%s", s, out)
			}
		}
	})

	t.Run("DryRunRecordsNoActor", func(t *testing.T) {
		ctx := initCtx(t)
		c := *ctx
		c.DryRun, c.Record = true, true
		resp, err := callOp(t, &c, "add", `{"kind":"capability","handle":"a"}`)
		if err != nil {
			t.Fatal(err)
		}
		if r := resp.(WriteResponse); !r.Record || !r.DryRun || len(r.Ops) != 1 {
			t.Errorf("response = %+v", r)
		}
		if files := logFiles(t, ctx.Root); len(files) != 0 {
			t.Errorf("dry run wrote %v", files)
		}
	})

	t.Run("Record", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"draft"}`)
		r := recordWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		if !r.Record || r.ChangeSet == cs || r.Entry.IsZero() {
			t.Fatalf("response = %+v", r)
		}
		if got := currentOf(t, ctx); got != cs {
			t.Errorf("current = %s, want %s", got, cs)
		}
		if len(opsOf(t, ctx, cs)) != 1 {
			t.Errorf("current change set changed")
		}
		c, _ := ctxWorkspace(t, ctx).changeSet(r.ChangeSet)
		if c.State != changeset.StateLanded || c.Open.Kind != changeset.KindRecord || c.Open.Intent != "add capability store" {
			t.Errorf("record = %+v", c)
		}
		if n := nodeAt(t, ctx, "store"); n.Version.IsZero() {
			t.Errorf("record node not landed: %+v", n)
		}
		if !strings.Contains(human(t, r), "Recorded and landed change set "+r.ChangeSet.String()) {
			t.Errorf("human:\n%s", human(t, r))
		}
	})

	t.Run("RecordSeesOnlyLanded", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"draft"}`)
		c := *ctx
		c.Record = true
		_, err := callOp(t, &c, "set", `{"ref":"draft","fields":["a=b"]}`)
		wantExit(t, err, 4)
	})

	t.Run("Enums", func(t *testing.T) {
		kinds := make([]string, 0)
		for _, k := range op.Kinds() {
			kinds = append(kinds, string(k))
		}
		types := make([]string, 0)
		for _, k := range op.LinkTypes() {
			types = append(types, string(k))
		}
		f, _ := reflect.TypeFor[AddRequest]().FieldByName("Kind")
		if got := strings.Split(f.Tag.Get("enum"), ","); !slices.Equal(got, kinds) {
			t.Errorf("add kind enum = %v, want %v", got, kinds)
		}
		f, _ = reflect.TypeFor[LinkRequest]().FieldByName("Type")
		if got := strings.Split(f.Tag.Get("enum"), ","); !slices.Equal(got, types) {
			t.Errorf("link type enum = %v, want %v", got, types)
		}
	})
}

func TestAdd(t *testing.T) {
	t.Run("KeysAndParent", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"api"}`)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"log","parent":"store"}`)
		store, api, log := nodeAt(t, ctx, "store"), nodeAt(t, ctx, "api"), nodeAt(t, ctx, "log")
		if store.Key != "i" || !(api.Key > store.Key) || !api.Parent.IsZero() {
			t.Errorf("keys: store %q api %q", store.Key, api.Key)
		}
		if log.Parent != store.ID || log.Key != "i" {
			t.Errorf("log = %+v", log)
		}
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"mid","after":"store"}`)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"first","before":"store"}`)
		mid, first := nodeAt(t, ctx, "mid"), nodeAt(t, ctx, "first")
		if !(first.Key < store.Key && store.Key < mid.Key && mid.Key < api.Key) {
			t.Errorf("keys: first %q store %q mid %q api %q", first.Key, store.Key, mid.Key, api.Key)
		}
		if got := len(opsOf(t, ctx, cs)); got != 5 {
			t.Errorf("got %d ops", got)
		}
	})

	t.Run("NonTreeKinds", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		r := mustWrite(t, ctx, "add", `{"kind":"decision","handle":"d1","fields":["statement=Use files"]}`)
		if n := nodeAt(t, ctx, "d1"); n.Key != "" || !n.Parent.IsZero() || n.ID != r.Nodes[0] {
			t.Errorf("decision = %+v", n)
		}
		for _, req := range []string{
			`{"kind":"decision","handle":"d2","parent":"store"}`,
			`{"kind":"decision","handle":"d2","before":"store"}`,
			`{"kind":"decision","handle":"d2","after":"store"}`,
		} {
			_, err := callOp(t, ctx, "add", req)
			wantExit(t, err, 2)
		}
	})

	t.Run("BadParentAndPosition", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"log","parent":"store"}`)
		mustWrite(t, ctx, "add", `{"kind":"intent","handle":"goal"}`)
		for req, code := range map[string]int{
			`{"kind":"capability","handle":"x","parent":"goal"}`:                  2, // other kind
			`{"kind":"capability","handle":"x","parent":"nosuch"}`:                4,
			`{"kind":"capability","handle":"x","before":"log"}`:                   2, // not a sibling
			`{"kind":"capability","handle":"x","before":"store","after":"store"}`: 2,
		} {
			_, err := callOp(t, ctx, "add", req)
			if ExitCode(err) != code {
				t.Errorf("%s: %v, exit %d; want %d", req, err, ExitCode(err), code)
			}
		}
	})

	t.Run("Fields", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store","fields":["statement=A <b> & c"],"json_fields":["n=[1, 2]","o={\"b\":1,\"a\":true}"]}`)
		f := nodeAt(t, ctx, "store").Fields
		if string(f["statement"]) != `"A <b> & c"` {
			t.Errorf("statement = %s", f["statement"])
		}
		if string(f["n"]) != "[\n  1,\n  2\n]" || !strings.HasPrefix(string(f["o"]), "{\n  \"a\": true") {
			t.Errorf("json fields = %s, %s", f["n"], f["o"])
		}
		for _, req := range []string{
			`{"kind":"capability","handle":"x","fields":["a=1","a=2"]}`,
			`{"kind":"capability","handle":"x","fields":["a=1"],"json_fields":["a=1"]}`,
			`{"kind":"capability","handle":"x","json_fields":["a={"]}`,
			`{"kind":"capability","handle":"x","json_fields":["a=1.5"]}`,
			`{"kind":"capability","handle":"x","fields":["noequals"]}`,
			`{"kind":"capability","handle":"x","fields":["Bad=1"]}`,
		} {
			_, err := callOp(t, ctx, "add", req)
			wantExit(t, err, 2)
		}
	})

	t.Run("BadKindAndHandle", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		_, err := callOp(t, ctx, "add", `{"kind":"widget","handle":"x"}`)
		wantExit(t, err, 2)
		_, err = callOp(t, ctx, "add", `{"kind":"capability","handle":"Not OK"}`)
		wantExit(t, err, 2)
		if !strings.Contains(err.Error(), "handle") {
			t.Errorf("message = %q", err)
		}
	})

	t.Run("HandleTaken", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		r := mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		_, err := callOp(t, ctx, "add", `{"kind":"decision","handle":"store"}`)
		wantExit(t, err, 3)
		if !strings.Contains(err.Error(), r.Nodes[0].String()) {
			t.Errorf("message lacks the holder: %q", err)
		}
		// The actor's handle is taken too.
		_, err = callOp(t, ctx, "add", `{"kind":"capability","handle":"test"}`)
		wantExit(t, err, 3)
	})

	t.Run("HandleFreeAfterRetire", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		recordWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		mustWrite(t, ctx, "retire", `{"ref":"store","reason":"Gone"}`)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
	})
}

func TestSet(t *testing.T) {
	t.Run("RevisesCreate", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store","fields":["a=1","b=2"]}`)
		created := onlyOp(t, opsOf(t, ctx, cs))
		r := mustWrite(t, ctx, "set", `{"ref":"store","handle":"files","fields":["c=3"],"unset":["a"]}`)
		o := onlyOp(t, opsOf(t, ctx, cs))
		if o.ID != created.ID || o.Type != op.TypeNodeCreate || o.Handle != "files" {
			t.Fatalf("op = %+v", o)
		}
		want := map[string]json.RawMessage{"b": raw(`"2"`), "c": raw(`"3"`)}
		if !reflect.DeepEqual(o.Fields, want) {
			t.Errorf("fields = %v", o.Fields)
		}
		if !slices.Equal(r.Nodes, []id.ID{o.Node}) {
			t.Errorf("nodes = %v", r.Nodes)
		}
	})

	t.Run("NewUpdateThenMerge", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		recordWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		landed := nodeAt(t, ctx, "store")
		mustWrite(t, ctx, "set", `{"ref":"store","fields":["a=1","b=2"]}`)
		first := onlyOp(t, opsOf(t, ctx, cs))
		if first.Type != op.TypeNodeUpdate || first.Base != landed.Version || first.Node != landed.ID || len(first.Set) != 2 {
			t.Fatalf("update = %+v", first)
		}
		mustWrite(t, ctx, "set", `{"ref":"store","fields":["c=3"],"unset":["a","z"]}`)
		o := onlyOp(t, opsOf(t, ctx, cs))
		if o.ID != first.ID || o.Base != landed.Version {
			t.Fatalf("merged = %+v", o)
		}
		want := map[string]json.RawMessage{"b": raw(`"2"`), "c": raw(`"3"`)}
		if !reflect.DeepEqual(o.Set, want) || !slices.Equal(o.Unset, []string{"a", "z"}) {
			t.Errorf("set %v unset %v", o.Set, o.Unset)
		}
		// Setting an unset field moves it back to set.
		mustWrite(t, ctx, "set", `{"ref":"store","fields":["a=4"],"handle":"files"}`)
		o = onlyOp(t, opsOf(t, ctx, cs))
		if string(o.Set["a"]) != `"4"` || !slices.Equal(o.Unset, []string{"z"}) || o.Handle != "files" {
			t.Errorf("set %v unset %v handle %q", o.Set, o.Unset, o.Handle)
		}
	})

	t.Run("Rename", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"api"}`)
		_, err := callOp(t, ctx, "set", `{"ref":"api","handle":"store"}`)
		wantExit(t, err, 3)
		mustWrite(t, ctx, "set", `{"ref":"api","handle":"api"}`)
	})

	t.Run("Errors", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		for req, code := range map[string]int{
			`{"ref":"store"}`:                                2, // nothing to change
			`{"ref":"nosuch","fields":["a=b"]}`:              4,
			`{"ref":"store","fields":["a=1"],"unset":["a"]}`: 2,
			`{"ref":"store","unset":["Bad"]}`:                2,
			`{"ref":"store","handle":"Bad"}`:                 2,
		} {
			_, err := callOp(t, ctx, "set", req)
			if ExitCode(err) != code {
				t.Errorf("%s: %v, exit %d; want %d", req, err, ExitCode(err), code)
			}
		}
	})

	t.Run("AmbiguousPrefix", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		a := mustWrite(t, ctx, "add", `{"kind":"capability","handle":"a"}`).Nodes[0]
		b := mustWrite(t, ctx, "add", `{"kind":"capability","handle":"b"}`).Nodes[0]
		prefix := commonPrefix(a.String(), b.String())
		if len(prefix) < 4 {
			t.Skipf("prefix %q too short", prefix)
		}
		_, err := callOp(t, ctx, "set", `{"ref":"`+prefix+`","fields":["a=b"]}`)
		wantExit(t, err, 2)
	})

	t.Run("Retired", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		r := recordWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		recordWrite(t, ctx, "retire", `{"ref":"store","reason":"Gone"}`)
		for _, name := range []string{"set", "mv", "retire"} {
			req := map[string]string{
				"set":    `{"ref":"` + r.Nodes[0].String() + `","fields":["a=b"]}`,
				"mv":     `{"ref":"` + r.Nodes[0].String() + `","root":true}`,
				"retire": `{"ref":"` + r.Nodes[0].String() + `","reason":"x"}`,
			}[name]
			_, err := callOp(t, ctx, name, req)
			wantExit(t, err, 3)
		}
	})

	t.Run("LandedConflict", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		recordWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		n := nodeAt(t, ctx, "store")
		w := ctxWorkspace(t, ctx)
		for _, v := range []string{`"1"`, `"2"`} {
			u := op.Op{ID: w.mustID(t), Type: op.TypeNodeUpdate, Node: n.ID, Base: n.Version,
				Set: map[string]json.RawMessage{"a": raw(v)}}
			if _, _, err := w.writeRecord("Fork", []op.Op{u}, w.mustID(t), w.mustID(t)); err != nil {
				t.Fatal(err)
			}
		}
		if !nodeAt(t, ctx, "store").Conflict {
			t.Fatal("no landed conflict")
		}
		for name, req := range map[string]string{
			"set":    `{"ref":"store","fields":["a=3"]}`,
			"mv":     `{"ref":"store","root":true}`,
			"retire": `{"ref":"store","reason":"x"}`,
		} {
			_, err := callOp(t, ctx, name, req)
			wantExit(t, err, 3)
			wantNext(t, err, "rhiza check")
			if !strings.Contains(err.Error(), "conflicting changes") {
				t.Errorf("%s message = %q", name, err)
			}
		}
	})
}

func TestMv(t *testing.T) {
	setup := func(t *testing.T) (*Ctx, id.ID) {
		ctx, cs := writeCtx(t)
		recordWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		recordWrite(t, ctx, "add", `{"kind":"capability","handle":"api"}`)
		recordWrite(t, ctx, "add", `{"kind":"capability","handle":"log","parent":"store"}`)
		return ctx, cs
	}

	t.Run("LandedThenReplace", func(t *testing.T) {
		ctx, cs := setup(t)
		log, api := nodeAt(t, ctx, "log"), nodeAt(t, ctx, "api")
		mustWrite(t, ctx, "mv", `{"ref":"log","root":true}`)
		first := onlyOp(t, opsOf(t, ctx, cs))
		if first.Type != op.TypeNodeMove || first.Base != log.Version || !first.Parent.IsZero() || !(first.Key > api.Key) {
			t.Fatalf("move = %+v (api key %q)", first, api.Key)
		}
		r := mustWrite(t, ctx, "mv", `{"ref":"log","parent":"api"}`)
		o := onlyOp(t, opsOf(t, ctx, cs))
		if o.ID != first.ID || o.Parent != api.ID || o.Key != "i" || o.Base != log.Version {
			t.Errorf("replaced = %+v", o)
		}
		if !slices.Equal(r.Nodes, []id.ID{log.ID}) {
			t.Errorf("nodes = %v", r.Nodes)
		}
	})

	t.Run("BeforeAfter", func(t *testing.T) {
		ctx, _ := setup(t)
		mustWrite(t, ctx, "mv", `{"ref":"api","before":"store"}`)
		if a, s := nodeAt(t, ctx, "api"), nodeAt(t, ctx, "store"); !(a.Key < s.Key) {
			t.Errorf("api %q store %q", a.Key, s.Key)
		}
		mustWrite(t, ctx, "mv", `{"ref":"api","after":"store"}`)
		if a, s := nodeAt(t, ctx, "api"), nodeAt(t, ctx, "store"); !(a.Key > s.Key) {
			t.Errorf("api %q store %q", a.Key, s.Key)
		}
		// Placing next to itself is not placing among its siblings.
		_, err := callOp(t, ctx, "mv", `{"ref":"api","after":"api"}`)
		wantExit(t, err, 2)
	})

	t.Run("RevisesCreate", func(t *testing.T) {
		ctx, cs := setup(t)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"new"}`)
		created := onlyOp(t, opsOf(t, ctx, cs))
		mustWrite(t, ctx, "mv", `{"ref":"new","parent":"store"}`)
		o := onlyOp(t, opsOf(t, ctx, cs))
		store, log := nodeAt(t, ctx, "store"), nodeAt(t, ctx, "log")
		if o.ID != created.ID || o.Type != op.TypeNodeCreate || o.Parent != store.ID || !(o.Key > log.Key) {
			t.Errorf("create = %+v", o)
		}
	})

	t.Run("Errors", func(t *testing.T) {
		ctx, _ := setup(t)
		mustWrite(t, ctx, "add", `{"kind":"decision","handle":"d"}`)
		mustWrite(t, ctx, "add", `{"kind":"intent","handle":"goal"}`)
		for req, code := range map[string]int{
			`{"ref":"d","root":true}`:                  2, // not a tree kind
			`{"ref":"log"}`:                            2, // nothing to do
			`{"ref":"log","root":true,"parent":"api"}`: 2,
			`{"ref":"store","parent":"store"}`:         2, // itself
			`{"ref":"store","parent":"log"}`:           2, // its descendant
			`{"ref":"store","parent":"goal"}`:          2, // other kind
			`{"ref":"store","parent":"nosuch"}`:        4,
			`{"ref":"api","before":"log"}`:             2, // not a sibling
		} {
			_, err := callOp(t, ctx, "mv", req)
			if ExitCode(err) != code {
				t.Errorf("%s: %v, exit %d; want %d", req, err, ExitCode(err), code)
			}
		}
	})
}

func TestRetire(t *testing.T) {
	t.Run("Landed", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		recordWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		n := nodeAt(t, ctx, "store")
		mustWrite(t, ctx, "set", `{"ref":"store","fields":["a=1"]}`)
		r := mustWrite(t, ctx, "retire", `{"ref":"store","reason":"Gone"}`)
		ops := opsOf(t, ctx, cs)
		if len(ops) != 2 || ops[1].Type != op.TypeNodeRetire || ops[1].Base != n.Version || ops[1].Reason != "Gone" {
			t.Fatalf("ops = %+v", ops)
		}
		if !slices.Equal(r.Nodes, []id.ID{n.ID}) {
			t.Errorf("nodes = %v", r.Nodes)
		}
		_, err := callOp(t, ctx, "retire", `{"ref":"`+n.ID.String()+`","reason":"Again"}`)
		wantExit(t, err, 3)
	})

	t.Run("CreatedHereIsDropped", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		recordWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"tmp"}`)
		mustWrite(t, ctx, "set", `{"ref":"store","fields":["a=1"]}`)
		mustWrite(t, ctx, "link", `{"from":"tmp","type":"depends-on","to":"store"}`)
		created := opsOf(t, ctx, cs)[0]
		r := mustWrite(t, ctx, "retire", `{"ref":"tmp","reason":"Never mind"}`)
		ops := opsOf(t, ctx, cs)
		if len(ops) != 1 || ops[0].Type != op.TypeNodeUpdate {
			t.Fatalf("ops = %+v", ops)
		}
		if len(r.Links) != 1 || !slices.Equal(r.Nodes, []id.ID{created.Node}) {
			t.Errorf("response = %+v", r)
		}
		// Its handle is free again.
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"tmp"}`)
	})

	t.Run("DryRunShowsDrop", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"tmp"}`)
		created := onlyOp(t, opsOf(t, ctx, cs))
		c := *ctx
		c.DryRun = true
		resp, err := callOp(t, &c, "retire", `{"ref":"tmp","reason":"x"}`)
		if err != nil {
			t.Fatal(err)
		}
		r := resp.(WriteResponse)
		if len(r.Ops) != 0 || !slices.Equal(r.Drop, []id.ID{created.ID}) {
			t.Errorf("response = %+v", r)
		}
		if !strings.Contains(human(t, r), "drop "+created.ID.String()) {
			t.Errorf("human:\n%s", human(t, r))
		}
	})

	t.Run("NeedsReason", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"tmp"}`)
		_, err := callOp(t, ctx, "retire", `{"ref":"tmp","reason":"  "}`)
		wantExit(t, err, 2)
	})
}

func TestLink(t *testing.T) {
	setup := func(t *testing.T) (*Ctx, id.ID) {
		ctx, cs := writeCtx(t)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"log"}`)
		mustWrite(t, ctx, "add", `{"kind":"decision","handle":"d"}`)
		mustWrite(t, ctx, "add", `{"kind":"constraint","handle":"c"}`)
		return ctx, cs
	}

	t.Run("Create", func(t *testing.T) {
		ctx, cs := setup(t)
		r := mustWrite(t, ctx, "link", `{"from":"log","type":"depends-on","to":"store"}`)
		ops := opsOf(t, ctx, cs)
		o := ops[len(ops)-1]
		log, store := nodeAt(t, ctx, "log"), nodeAt(t, ctx, "store")
		if o.Type != op.TypeLinkCreate || o.From != log.ID || o.To != store.ID || o.LinkType != op.LinkDependsOn {
			t.Fatalf("op = %+v", o)
		}
		if !slices.Equal(r.Links, []id.ID{o.Link}) || len(r.Nodes) != 0 || r.Nodes == nil {
			t.Errorf("response = %+v", r)
		}
		if !strings.Contains(human(t, r), o.Link.String()) {
			t.Errorf("human:\n%s", human(t, r))
		}
		_, err := callOp(t, ctx, "link", `{"from":"log","type":"depends-on","to":"store"}`)
		wantExit(t, err, 3)
	})

	t.Run("NoteAndDecision", func(t *testing.T) {
		ctx, _ := setup(t)
		mustWrite(t, ctx, "link", `{"from":"log","type":"relates-to","to":"store","note":"Same data"}`)
		mustWrite(t, ctx, "link", `{"from":"log","type":"exempts","to":"c","note":"Prototype","decision":"d"}`)
		for _, req := range []string{
			`{"from":"log","type":"bogus","to":"store"}`,
			`{"from":"store","type":"relates-to","to":"log"}`,
			`{"from":"store","type":"exempts","to":"c","note":"x"}`,
			`{"from":"store","type":"exempts","to":"c","note":"x","decision":"log"}`,
			`{"from":"log","type":"serves","to":"store","decision":"d"}`,
			`{"from":"log","type":"serves","to":"log"}`,
		} {
			_, err := callOp(t, ctx, "link", req)
			wantExit(t, err, 2)
		}
		_, err := callOp(t, ctx, "link", `{"from":"nosuch","type":"serves","to":"log"}`)
		wantExit(t, err, 4)
	})

	t.Run("RetiredEndpoint", func(t *testing.T) {
		ctx, _ := setup(t)
		r := recordWrite(t, ctx, "add", `{"kind":"capability","handle":"old"}`)
		recordWrite(t, ctx, "retire", `{"ref":"old","reason":"x"}`)
		_, err := callOp(t, ctx, "link", `{"from":"log","type":"depends-on","to":"`+r.Nodes[0].String()+`"}`)
		wantExit(t, err, 3)
	})
}

// TestLinkEndpointKinds shows buildLink agrees with check.Allowed for every
// link type and every pair of end kinds, change sets included: a pair the
// table allows is never refused as a kind mismatch, and any other is a
// usage error naming both kinds that writes nothing.
func TestLinkEndpointKinds(t *testing.T) {
	ctx, cs := writeCtx(t)
	handle := func(prefix string, e check.End) string { return prefix + "-" + strings.ReplaceAll(string(e), " ", "-") }
	for _, k := range op.Kinds() {
		for _, p := range []string{"x", "y"} {
			if h := handle(p, check.NodeEnd(k)); op.HasParent(k) {
				recordWrite(t, ctx, "add", fmt.Sprintf(`{"kind":%q,"handle":%q}`, k, h))
			} else {
				landNode(t, ctx, k, h)
			}
		}
	}
	other := newChange(t, ctx, `{"intent":"Other"}`).ChangeSet
	ref := func(prefix string, e check.End) string {
		switch {
		case e != check.ChangeSet:
			return handle(prefix, e)
		case prefix == "x":
			return cs.String()
		}
		return other.String()
	}
	for _, lt := range op.LinkTypes() {
		for _, from := range check.Ends() {
			for _, to := range check.Ends() {
				r := LinkRequest{From: ref("x", from), Type: string(lt), To: ref("y", to), Note: "n"}
				if lt == op.LinkExempts {
					r.Decision = "x-decision"
				}
				req, err := json.Marshal(r)
				if err != nil {
					t.Fatal(err)
				}
				_, err = callOp(t, ctx, "link", string(req))
				mismatch := err != nil && strings.Contains(err.Error(), "does not go from")
				switch {
				case check.Allowed(lt, from, to) && mismatch:
					t.Errorf("%s from %s to %s is allowed, but refused: %v", lt, from, to, err)
				case !check.Allowed(lt, from, to) && !mismatch:
					t.Errorf("%s from %s to %s is not allowed, but got %v", lt, from, to, err)
				case mismatch:
					want := fmt.Sprintf("link type %q does not go from %s to %s", lt, from, to)
					if ExitCode(err) != 2 || !strings.Contains(err.Error(), want) {
						t.Errorf("%s from %s to %s: %v (exit %d), want usage %q", lt, from, to, err, ExitCode(err), want)
					}
				}
			}
		}
	}
}

// TestLinkCycle refuses a depends-on link between change sets that would
// close a loop, counting the links the same write batch added.
func TestLinkCycle(t *testing.T) {
	// link writes "from depends-on to" into change set in.
	link := func(ctx *Ctx, in, from, to id.ID) error {
		c := *ctx
		c.ChangeRef = in.String()
		_, err := callOp(t, &c, "link", fmt.Sprintf(`{"from":%q,"type":"depends-on","to":%q}`, from, to))
		return err
	}
	// wantLoop checks err refuses a loop and names to as the one in the way.
	wantLoop := func(t *testing.T, err error, to id.ID) {
		t.Helper()
		wantExit(t, err, 3)
		e := AsError(err)
		if !strings.Contains(e.Message, "would close a loop") || e.Next != "rhiza change show "+shortID(to) {
			t.Errorf("err = %q, next %q", e.Message, e.Next)
		}
	}
	sets := func(t *testing.T, n int) (*Ctx, []id.ID) {
		ctx := initCtx(t)
		var out []id.ID
		for i := range n {
			out = append(out, newChange(t, ctx, fmt.Sprintf(`{"intent":"S%d"}`, i)).ChangeSet)
		}
		return ctx, out
	}

	t.Run("Direct", func(t *testing.T) {
		ctx, s := sets(t, 2)
		a, b := s[0], s[1]
		if err := link(ctx, b, b, a); err != nil {
			t.Fatal(err)
		}
		before := len(logFiles(t, ctx.Root))
		wantLoop(t, link(ctx, a, a, b), b)
		if n := len(logFiles(t, ctx.Root)); n != before {
			t.Errorf("wrote %d files", n-before)
		}
	})

	t.Run("ThreeSets", func(t *testing.T) {
		ctx, s := sets(t, 3)
		a, b, c := s[0], s[1], s[2]
		if err := link(ctx, b, b, a); err != nil {
			t.Fatal(err)
		}
		if err := link(ctx, c, c, b); err != nil {
			t.Fatal(err)
		}
		wantLoop(t, link(ctx, a, a, c), c)
	})

	t.Run("NoLoop", func(t *testing.T) {
		ctx, s := sets(t, 3)
		a, b, c := s[0], s[1], s[2]
		if err := link(ctx, b, b, a); err != nil {
			t.Fatal(err)
		}
		if err := link(ctx, c, c, a); err != nil {
			t.Fatal(err)
		}
		if err := link(ctx, c, c, b); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("SelfLink", func(t *testing.T) {
		ctx, s := sets(t, 1)
		wantExit(t, link(ctx, s[0], s[0], s[0]), 2)
	})

	t.Run("ApplyBatch", func(t *testing.T) {
		ctx, s := sets(t, 3)
		a, b, c := s[0], s[1], s[2]
		before := len(logFiles(t, ctx.Root))
		// Written in c, neither link is from c, yet both would count once
		// c lands.
		_, err := applyFile(t, ctx, fmt.Sprintf(`{"writes":[
		  {"op":"link","from":%q,"type":"depends-on","to":%q},
		  {"op":"link","from":%q,"type":"depends-on","to":%q}]}`, a, b, b, a))
		wantLoop(t, err, a)
		if e := AsError(err); !strings.HasPrefix(e.Message, "writes[1]: ") {
			t.Errorf("err = %q", e.Message)
		}
		if n := len(logFiles(t, ctx.Root)); n != before {
			t.Errorf("wrote %d files", n-before)
		}
		if ops := opsOf(t, ctx, c); len(ops) != 0 {
			t.Errorf("ops = %+v", ops)
		}
	})

}

// TestLinkCapabilityLoop refuses a depends-on link between capabilities
// that would close a loop, naming it, whether the links in the way are
// earlier writes or earlier items of the same apply.
func TestLinkCapabilityLoop(t *testing.T) {
	// caps adds capabilities a, b and c to a new change set.
	caps := func(t *testing.T) *Ctx {
		ctx, _ := writeCtx(t)
		for _, h := range []string{"a", "b", "c"} {
			mustWrite(t, ctx, "add", `{"kind":"capability","handle":"`+h+`"}`)
		}
		return ctx
	}
	const loop = "depends-on from c to a would close a loop: a depends on b, b depends on c"

	t.Run("Direct", func(t *testing.T) {
		ctx := caps(t)
		mustWrite(t, ctx, "link", `{"from":"a","type":"depends-on","to":"b"}`)
		mustWrite(t, ctx, "link", `{"from":"b","type":"depends-on","to":"c"}`)
		before := len(logFiles(t, ctx.Root))
		_, err := callOp(t, ctx, "link", `{"from":"c","type":"depends-on","to":"a"}`)
		wantExit(t, err, 3)
		if e := AsError(err); e.Message != loop || e.Next != "rhiza show a" {
			t.Errorf("err = %q, next %q", e.Message, e.Next)
		}
		if n := len(logFiles(t, ctx.Root)); n != before {
			t.Errorf("wrote %d files", n-before)
		}
	})

	t.Run("Two", func(t *testing.T) {
		ctx := caps(t)
		mustWrite(t, ctx, "link", `{"from":"a","type":"depends-on","to":"b"}`)
		_, err := callOp(t, ctx, "link", `{"from":"b","type":"depends-on","to":"a"}`)
		wantExit(t, err, 3)
		if e := AsError(err); e.Message != "depends-on from b to a would close a loop: a depends on b" {
			t.Errorf("err = %q", e.Message)
		}
	})

	t.Run("ApplyBatch", func(t *testing.T) {
		ctx := caps(t)
		before := len(logFiles(t, ctx.Root))
		_, err := applyFile(t, ctx, `{"writes":[
		  {"op":"link","from":"a","type":"depends-on","to":"b"},
		  {"op":"link","from":"b","type":"depends-on","to":"c"},
		  {"op":"link","from":"c","type":"depends-on","to":"a"}]}`)
		wantExit(t, err, 3)
		if e := AsError(err); e.Message != "writes[2]: "+loop {
			t.Errorf("err = %q", e.Message)
		}
		if n := len(logFiles(t, ctx.Root)); n != before {
			t.Errorf("wrote %d files", n-before)
		}
	})

	t.Run("Unlinked", func(t *testing.T) {
		// A retired link no longer counts.
		ctx := caps(t)
		mustWrite(t, ctx, "link", `{"from":"a","type":"depends-on","to":"b"}`)
		mustWrite(t, ctx, "unlink", `{"ref":"a","type":"depends-on","to":"b"}`)
		mustWrite(t, ctx, "link", `{"from":"b","type":"depends-on","to":"a"}`)
	})

	t.Run("NoLoop", func(t *testing.T) {
		ctx := caps(t)
		mustWrite(t, ctx, "link", `{"from":"a","type":"depends-on","to":"b"}`)
		mustWrite(t, ctx, "link", `{"from":"b","type":"depends-on","to":"c"}`)
		mustWrite(t, ctx, "link", `{"from":"a","type":"depends-on","to":"c"}`)
		mustWrite(t, ctx, "link", `{"from":"c","type":"relates-to","to":"a","note":"x"}`)
	})
}

func TestUnlink(t *testing.T) {
	t.Run("CreatedHereIsDropped", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"log"}`)
		l := mustWrite(t, ctx, "link", `{"from":"log","type":"depends-on","to":"store"}`).Links[0]
		r := mustWrite(t, ctx, "unlink", `{"ref":"`+l.String()+`"}`)
		if len(opsOf(t, ctx, cs)) != 2 || !slices.Equal(r.Links, []id.ID{l}) || len(r.Nodes) != 0 {
			t.Errorf("response = %+v", r)
		}
	})

	t.Run("Landed", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		recordWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		recordWrite(t, ctx, "add", `{"kind":"capability","handle":"log"}`)
		l := recordWrite(t, ctx, "link", `{"from":"log","type":"depends-on","to":"store"}`).Links[0]
		k, err := ctxWorkspace(t, ctx).model(cs).ResolveLink(l.String())
		if err != nil {
			t.Fatal(err)
		}
		mustWrite(t, ctx, "unlink", `{"ref":"log","type":"depends-on","to":"store"}`)
		o := onlyOp(t, opsOf(t, ctx, cs))
		if o.Type != op.TypeLinkRetire || o.Link != l || o.Base != k.Version {
			t.Errorf("op = %+v", o)
		}
		_, err = callOp(t, ctx, "unlink", `{"ref":"`+l.String()+`"}`)
		wantExit(t, err, 3)
		_, err = callOp(t, ctx, "unlink", `{"ref":"log","type":"depends-on","to":"store"}`)
		wantExit(t, err, 4)
	})

	t.Run("Usage", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		_, err := callOp(t, ctx, "unlink", `{"ref":"a","type":"depends-on"}`)
		wantExit(t, err, 2)
		_, err = callOp(t, ctx, "unlink", `{"ref":"ZZZZZZ"}`)
		wantExit(t, err, 4)
	})
}

// TestStackedWrites writes in a change set on top of the in-flight change
// sets it depends on.
func TestStackedWrites(t *testing.T) {
	// in returns ctx writing into change set cs.
	in := func(ctx *Ctx, cs id.ID) *Ctx {
		c := *ctx
		c.ChangeRef = cs.String()
		return &c
	}
	at := func(t *testing.T, ctx *Ctx, cs id.ID, ref string) Node {
		t.Helper()
		n, err := ctxWorkspace(t, ctx).model(cs).Resolve(ref)
		if err != nil {
			t.Fatalf("resolve %s in %s: %v", ref, cs, err)
		}
		return n
	}
	// setup gives A, with capability base, and B, depending on A when dep.
	setup := func(t *testing.T, dep bool) (ctx *Ctx, a, b id.ID) {
		ctx = initCtx(t)
		a = newChange(t, ctx, `{"intent":"A"}`).ChangeSet
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"base","fields":["statement=Base"]}`)
		b = newChange(t, ctx, `{"intent":"B"}`).ChangeSet
		if dep {
			planLink(t, ctx, op.LinkDependsOn, b, a)
		}
		return ctx, a, b
	}

	t.Run("SeesDependency", func(t *testing.T) {
		ctx, a, b := setup(t, true)
		created := onlyOp(t, opsOf(t, ctx, a))
		base := at(t, ctx, b, "base")
		if base.Version != created.ID || base.Conflict || base.Pending != nil {
			t.Fatalf("base in B = %+v", base)
		}
		bc := in(ctx, b)
		leaf := mustWrite(t, bc, "add", `{"kind":"capability","handle":"leaf","parent":"base"}`).Nodes[0]
		if n := at(t, ctx, b, "leaf"); n.ID != leaf || n.Parent != base.ID {
			t.Errorf("leaf = %+v", n)
		}
		mustWrite(t, bc, "add", `{"kind":"capability","handle":"other"}`)
		mustWrite(t, bc, "link", `{"from":"leaf","type":"depends-on","to":"base"}`)
		mustWrite(t, bc, "set", `{"ref":"base","fields":["statement=Changed"]}`)
		mustWrite(t, bc, "mv", `{"ref":"base","parent":"other"}`)
		mustWrite(t, bc, "retire", `{"ref":"base","reason":"Gone"}`)
		var kinds []op.Type
		for _, o := range opsOf(t, ctx, b) {
			if o.Node == base.ID {
				kinds = append(kinds, o.Type)
				if o.Base != created.ID {
					t.Errorf("%s base = %s, want A's create %s", o.Type, o.Base, created.ID)
				}
			}
		}
		if want := []op.Type{op.TypeNodeUpdate, op.TypeNodeMove, op.TypeNodeRetire}; !slices.Equal(kinds, want) {
			t.Errorf("ops on base = %v, want %v", kinds, want)
		}
		// A's own view is unchanged.
		if n := at(t, ctx, a, "base"); n.Retired || !n.Version.IsZero() || !n.Parent.IsZero() {
			t.Errorf("base in A = %+v", n)
		}
	})

	t.Run("NoDependencyNoView", func(t *testing.T) {
		ctx, _, b := setup(t, false)
		_, err := callOp(t, in(ctx, b), "add", `{"kind":"capability","handle":"leaf","parent":"base"}`)
		wantExit(t, err, 4)
		wantNext(t, err, lsNext)
	})

	t.Run("WithdrawnDependency", func(t *testing.T) {
		ctx, a, b := setup(t, true)
		move(t, in(ctx, a), "withdraw", `{}`)
		_, err := callOp(t, in(ctx, b), "add", `{"kind":"capability","handle":"leaf","parent":"base"}`)
		wantExit(t, err, 4)
	})

	t.Run("LandedDependency", func(t *testing.T) {
		ctx, a, b := setup(t, true)
		for _, s := range []string{"propose", "accept", "land"} {
			move(t, in(ctx, a), s, `{}`)
		}
		base := at(t, ctx, b, "base")
		mustWrite(t, in(ctx, b), "set", `{"ref":"base","fields":["statement=Changed"]}`)
		var o op.Op
		for _, x := range opsOf(t, ctx, b) {
			if x.Node == base.ID {
				o = x
			}
		}
		if o.Type != op.TypeNodeUpdate || o.Base != base.Version || base.Version != onlyOp(t, opsOf(t, ctx, a)).ID {
			t.Errorf("update = %+v, base %+v", o, base)
		}
	})

	t.Run("Transitive", func(t *testing.T) {
		ctx, a, b := setup(t, true)
		mustWrite(t, in(ctx, b), "add", `{"kind":"capability","handle":"leaf","parent":"base"}`)
		c := newChange(t, ctx, `{"intent":"C"}`).ChangeSet
		planLink(t, ctx, op.LinkDependsOn, c, b)
		cc := in(ctx, c)
		mustWrite(t, cc, "add", `{"kind":"capability","handle":"x","parent":"base"}`)
		mustWrite(t, cc, "add", `{"kind":"capability","handle":"y","parent":"leaf"}`)
		r := mustWrite(t, cc, "set", `{"ref":"base","fields":["statement=C"]}`)
		if o := opsOf(t, ctx, c); o[len(o)-1].Base != onlyOp(t, opsOf(t, ctx, a)).ID || r.ChangeSet != c {
			t.Errorf("update = %+v", o[len(o)-1])
		}
	})
}
