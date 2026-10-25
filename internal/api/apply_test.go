package api

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/genesis/internal/op"
)

// writeFile writes an apply file into a temp folder and returns its path.
func writeFile(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "writes.json")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

// applyFile runs apply on a file holding content.
func applyFile(t *testing.T, ctx *Ctx, content string) (WriteResponse, error) {
	t.Helper()
	resp, err := callOp(t, ctx, "apply", fmt.Sprintf(`{"file":%q}`, writeFile(t, content)))
	if err != nil {
		return WriteResponse{}, err
	}
	return resp.(WriteResponse), nil
}

const treeWrites = `{"writes":[
  {"op":"add","kind":"capability","handle":"store"},
  {"op":"add","kind":"capability","handle":"log","parent":"store"},
  {"op":"add","kind":"decision","handle":"d1"},
  {"op":"link","from":"log","type":"decided-by","to":"d1"},
  {"op":"set","ref":"store","fields":["statement=Keeps the log"]}
]}`

func TestApply(t *testing.T) {
	t.Run("OneEntry", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		before := len(logFiles(t, ctx.Root))
		r, err := applyFile(t, ctx, treeWrites)
		if err != nil {
			t.Fatal(err)
		}
		if n := len(logFiles(t, ctx.Root)); n != before+1 {
			t.Fatalf("wrote %d files, want 1", n-before)
		}
		ops := opsOf(t, ctx, cs)
		if len(ops) != 4 {
			t.Fatalf("got %d ops, want 4: %+v", len(ops), ops)
		}
		store, log, d1 := nodeAt(t, ctx, "store"), nodeAt(t, ctx, "log"), nodeAt(t, ctx, "d1")
		if log.Parent != store.ID || string(store.Fields["statement"]) != `"Keeps the log"` {
			t.Errorf("store = %+v, log = %+v", store, log)
		}
		// D9 within the file: set revised store's create.
		creates := 0
		for _, o := range ops {
			if o.Type == op.TypeNodeCreate && o.Node == store.ID {
				creates++
			}
			if o.Type == op.TypeNodeUpdate {
				t.Errorf("unexpected update %+v", o)
			}
		}
		if creates != 1 {
			t.Errorf("store creates = %d", creates)
		}
		if r.ChangeSet != cs || r.Entry.IsZero() || len(r.Links) != 1 {
			t.Errorf("response = %+v", r)
		}
		want := []string{store.ID.String(), log.ID.String(), d1.ID.String()}
		got := []string{}
		for _, n := range r.Nodes {
			got = append(got, n.String())
		}
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("nodes = %v, want %v", got, want)
		}
		if !strings.Contains(human(t, r), "Wrote entry") {
			t.Errorf("human = %s", human(t, r))
		}
	})

	t.Run("RevisesExistingChangeSet", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		if _, err := applyFile(t, ctx, `{"writes":[
		  {"op":"set","ref":"store","fields":["statement=A"]},
		  {"op":"add","kind":"capability","handle":"x"},
		  {"op":"retire","ref":"x","reason":"oops"}]}`); err != nil {
			t.Fatal(err)
		}
		o := onlyOp(t, opsOf(t, ctx, cs))
		if o.Type != op.TypeNodeCreate || string(o.Fields["statement"]) != `"A"` {
			t.Errorf("op = %+v", o)
		}
	})

	t.Run("DropsStartOps", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		r, err := applyFile(t, ctx, `{"writes":[{"op":"retire","ref":"store","reason":"no"}]}`)
		if err != nil {
			t.Fatal(err)
		}
		if ops := opsOf(t, ctx, cs); len(ops) != 0 {
			t.Errorf("ops = %+v", ops)
		}
		if len(r.Nodes) != 1 {
			t.Errorf("nodes = %v", r.Nodes)
		}
	})

	t.Run("NothingToWrite", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		before := len(logFiles(t, ctx.Root))
		_, err := applyFile(t, ctx, `{"writes":[
		  {"op":"add","kind":"capability","handle":"x"},
		  {"op":"retire","ref":"x","reason":"oops"}]}`)
		wantExit(t, err, 2)
		if !strings.Contains(err.Error(), "nothing to write") {
			t.Errorf("err = %v", err)
		}
		if n := len(logFiles(t, ctx.Root)); n != before {
			t.Errorf("wrote %d files", n-before)
		}
	})

	t.Run("BadItems", func(t *testing.T) {
		for name, tc := range map[string]struct {
			content, want string
			code          int
		}{
			"UnknownField":  {`{"writes":[{"op":"add","kind":"decision","handle":"a"},{"op":"add","kind":"decision","handle":"b","nope":1}]}`, "writes[1]: ", 2},
			"UnknownOp":     {`{"writes":[{"op":"frob"}]}`, "writes[0]: ", 2},
			"MissingOp":     {`{"writes":[{"kind":"decision","handle":"a"}]}`, "writes[0]: ", 2},
			"OpNotString":   {`{"writes":[{"op":1}]}`, "writes[0]: ", 2},
			"ItemNotObject": {`{"writes":[[]]}`, "writes[0]: ", 2},
			"NotFound":      {`{"writes":[{"op":"add","kind":"decision","handle":"a"},{"op":"set","ref":"nosuch","handle":"b"}]}`, "writes[1]: ", 4},
			"Conflict":      {`{"writes":[{"op":"add","kind":"decision","handle":"a"},{"op":"add","kind":"decision","handle":"a"}]}`, "writes[1]: ", 3},
			"InvalidOp":     {`{"writes":[{"op":"add","kind":"decision","handle":"a","fields":["Bad Name=x"]}]}`, "writes[0]: ", 2},
			"EmptyWrites":   {`{"writes":[]}`, "no writes", 2},
			"NoWrites":      {`{}`, "no writes", 2},
			"UnknownTop":    {`{"writes":[],"x":1}`, "", 2},
			"NotJSON":       {`nope`, "", 2},
		} {
			t.Run(name, func(t *testing.T) {
				ctx, _ := writeCtx(t)
				before := len(logFiles(t, ctx.Root))
				_, err := applyFile(t, ctx, tc.content)
				wantExit(t, err, tc.code)
				if !strings.Contains(err.Error(), tc.want) {
					t.Errorf("err = %v, want %q", err, tc.want)
				}
				if n := len(logFiles(t, ctx.Root)); n != before {
					t.Errorf("wrote %d files", n-before)
				}
			})
		}
	})

	t.Run("KeepsNext", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		_, err := applyFile(t, ctx, `{"writes":[{"op":"add","kind":"widget","handle":"a"}]}`)
		wantExit(t, err, 2)
		wantNext(t, err, "genesis describe add")
	})

	t.Run("MissingFile", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		_, err := callOp(t, ctx, "apply", `{"file":"/no/such/file.json"}`)
		wantExit(t, err, 2)
		_, err = callOp(t, ctx, "apply", `{}`)
		wantExit(t, err, 2)
	})

	t.Run("Stdin", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		ctx.Stdin = strings.NewReader(treeWrites)
		if _, err := callOp(t, ctx, "apply", `{"file":"-"}`); err != nil {
			t.Fatal(err)
		}
		if ops := opsOf(t, ctx, cs); len(ops) != 4 {
			t.Errorf("got %d ops", len(ops))
		}
	})

	t.Run("NoChangeSet", func(t *testing.T) {
		ctx, _ := writeCtx(t)
		if _, err := callOp(t, ctx, "change.use", `{"none":true}`); err != nil {
			t.Fatal(err)
		}
		_, err := applyFile(t, ctx, treeWrites)
		wantExit(t, err, 2)
		wantNext(t, err, newEntryNext)
	})

	t.Run("Record", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		c := *ctx
		c.Record = true
		before := len(logFiles(t, ctx.Root))
		r, err := applyFile(t, &c, treeWrites)
		if err != nil {
			t.Fatal(err)
		}
		if !r.Record || r.ChangeSet == cs {
			t.Fatalf("response = %+v", r)
		}
		if n := len(logFiles(t, ctx.Root)); n != before+3 {
			t.Errorf("wrote %d files, want 3", n-before)
		}
		if ops := opsOf(t, ctx, cs); len(ops) != 0 {
			t.Errorf("current got ops %+v", ops)
		}
		rec, _ := ctxWorkspace(t, ctx).changeSet(r.ChangeSet)
		if len(rec.Ops) != 4 || !strings.HasPrefix(rec.Open.Intent, "apply ") {
			t.Errorf("record = %+v", rec)
		}
		if nodeAt(t, ctx, "log").Version.IsZero() {
			t.Error("log is not landed")
		}
	})

	t.Run("DryRun", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		c := *ctx
		c.DryRun = true
		before := len(logFiles(t, ctx.Root))
		r, err := applyFile(t, &c, treeWrites)
		if err != nil {
			t.Fatal(err)
		}
		if n := len(logFiles(t, ctx.Root)); n != before {
			t.Errorf("wrote %d files", n-before)
		}
		if !r.DryRun || !r.Entry.IsZero() || len(r.Ops) != 4 || r.ChangeSet != cs {
			t.Errorf("response = %+v", r)
		}
	})

	t.Run("Many", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		var items []string
		for i := range 30 {
			items = append(items, fmt.Sprintf(`{"op":"add","kind":"capability","handle":"c%d"}`, i))
		}
		before := len(logFiles(t, ctx.Root))
		if _, err := applyFile(t, ctx, `{"writes":[`+strings.Join(items, ",")+`]}`); err != nil {
			t.Fatal(err)
		}
		if n := len(logFiles(t, ctx.Root)); n != before+1 {
			t.Errorf("wrote %d files, want 1", n-before)
		}
		if ops := opsOf(t, ctx, cs); len(ops) != 30 {
			t.Errorf("got %d ops", len(ops))
		}
	})
}

// TestApplyInputSchema checks describe apply publishes the file schema: a
// oneOf over the item ops, each that op's request schema plus a required
// "op" fixed to its name.
func TestApplyInputSchema(t *testing.T) {
	resp, err := callDescribe(t, `{"op":"apply"}`)
	if err != nil {
		t.Fatal(err)
	}
	var file struct {
		Required   []string `json:"required"`
		Properties struct {
			Writes struct {
				Type     string `json:"type"`
				MinItems int    `json:"minItems"`
				Items    struct {
					OneOf []map[string]json.RawMessage `json:"oneOf"`
				} `json:"items"`
			} `json:"writes"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(resp.Operations[0].Input, &file); err != nil {
		t.Fatalf("input schema %s: %v", resp.Operations[0].Input, err)
	}
	w := file.Properties.Writes
	if !slices.Equal(file.Required, []string{"writes"}) || w.Type != "array" || w.MinItems != 1 {
		t.Errorf("file schema = %s", resp.Operations[0].Input)
	}
	verbs := []string{"add", "set", "mv", "retire", "link", "unlink"}
	if len(w.Items.OneOf) != len(verbs) {
		t.Fatalf("oneOf has %d branches, want %d", len(w.Items.OneOf), len(verbs))
	}
	for i, verb := range verbs {
		b := w.Items.OneOf[i]
		var props map[string]json.RawMessage
		var required []string
		if err := json.Unmarshal(b["properties"], &props); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(b["required"], &required); err != nil {
			t.Fatal(err)
		}
		var opProp struct {
			Const string `json:"const"`
		}
		if err := json.Unmarshal(props["op"], &opProp); err != nil || opProp.Const != verb {
			t.Errorf("branch %d op = %s, want const %q", i, props["op"], verb)
		}
		if !slices.Contains(required, "op") {
			t.Errorf("branch %s does not require op: %v", verb, required)
		}
		o, _ := Lookup(verb)
		reqSchema := Schema(o.Request)
		want := map[string]bool{"op": true}
		for k := range reqSchema["properties"].(map[string]any) {
			want[k] = true
		}
		got := map[string]bool{}
		for k := range props {
			got[k] = true
		}
		if !maps.Equal(got, want) {
			t.Errorf("branch %s properties = %v, want %v", verb, got, want)
		}
		wantReq := append(slices.Clone(reqSchema["required"].([]string)), "op")
		slices.Sort(wantReq)
		if !slices.Equal(required, wantReq) {
			t.Errorf("branch %s required = %v, want %v", verb, required, wantReq)
		}
	}
}

// TestApplySchemaMatchesDecoder checks the ops the schema publishes are
// exactly the ops the decoder accepts.
func TestApplySchemaMatchesDecoder(t *testing.T) {
	var schemaOps []string
	for _, w := range applyWrites {
		schemaOps = append(schemaOps, w.op)
	}
	decoderOps := slices.Sorted(maps.Keys(applyBuilders))
	if !slices.Equal(slices.Sorted(slices.Values(schemaOps)), decoderOps) {
		t.Errorf("schema ops %v, decoder ops %v", schemaOps, decoderOps)
	}
	for _, o := range Operations() {
		_, err := decodeWrite(json.RawMessage(fmt.Sprintf(`{"op":%q}`, o.Name)))
		unknown := err != nil && strings.Contains(err.Error(), "unknown op")
		if unknown == slices.Contains(schemaOps, o.Name) {
			t.Errorf("decodeWrite op %s: err %v; in schema: %v", o.Name, err, !unknown)
		}
	}
}
