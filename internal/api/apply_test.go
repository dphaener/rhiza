package api

import (
	"fmt"
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
