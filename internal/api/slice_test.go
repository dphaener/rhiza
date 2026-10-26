package api

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/slice"
)

// sliceCtx is initCtx with the user cache dir (where the index lives) in a
// temp folder, so tests never write to the real one.
func sliceCtx(t *testing.T) *Ctx {
	t.Helper()
	cache := t.TempDir()
	t.Setenv("HOME", cache)
	t.Setenv("XDG_CACHE_HOME", cache)
	return initCtx(t)
}

// sliceModel lands a root capability with a constraint and a child with a
// decision.
func sliceModel(t *testing.T, ctx *Ctx) {
	t.Helper()
	recordWrite(t, ctx, "add", `{"kind":"capability","handle":"store","fields":["statement=Keeps the log"]}`)
	recordWrite(t, ctx, "add", `{"kind":"capability","handle":"log","parent":"store","fields":["statement=Appends entries"]}`)
	recordWrite(t, ctx, "add", `{"kind":"constraint","handle":"no-network","fields":["statement=Never touches the network"]}`)
	recordWrite(t, ctx, "link", `{"from":"no-network","type":"constrains","to":"store"}`)
	recordWrite(t, ctx, "add", `{"kind":"decision","handle":"use-files","fields":["statement=One file per entry"]}`)
	recordWrite(t, ctx, "link", `{"from":"log","type":"decided-by","to":"use-files"}`)
}

func callSlice(t *testing.T, ctx *Ctx, req string) (SliceResponse, error) {
	t.Helper()
	resp, err := callOp(t, ctx, "slice", req)
	if err != nil {
		return SliceResponse{}, err
	}
	return resp.(SliceResponse), nil
}

// The response's fields are slice.Doc's, tags and all, so its schema is
// the document's.
func TestSliceResponseIsDoc(t *testing.T) {
	got, err := canon.Marshal(Schema(reflect.TypeFor[SliceResponse]()))
	if err != nil {
		t.Fatal(err)
	}
	want, err := canon.Marshal(Schema(reflect.TypeFor[slice.Doc]()))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("SliceResponse schema differs from slice.Doc's:\n%s\nwant\n%s", got, want)
	}
}

// --json prints exactly the bytes slice.Fit fitted, and Human the markdown.
func TestSliceBytesAreFit(t *testing.T) {
	ctx := sliceCtx(t)
	sliceModel(t, ctx)
	for _, budget := range []int{4000, 60} {
		resp, err := callSlice(t, ctx, `{"ref":"log","budget":`+itoa(budget)+`}`)
		if err != nil {
			t.Fatal(err)
		}
		gotJSON, err := canon.Marshal(resp)
		if err != nil {
			t.Fatal(err)
		}
		var md bytes.Buffer
		resp.Human(&md)

		wantJSON, wantMD := fitDirect(t, ctx, "log", budget)
		if !bytes.Equal(gotJSON, wantJSON) {
			t.Errorf("budget %d: JSON\n%s\nwant\n%s", budget, gotJSON, wantJSON)
		}
		if !bytes.Equal(md.Bytes(), wantMD) {
			t.Errorf("budget %d: markdown\n%s\nwant\n%s", budget, md.Bytes(), wantMD)
		}
		if resp.Tokens != slice.Tokens(gotJSON) {
			t.Errorf("budget %d: tokens %d, the JSON is %d", budget, resp.Tokens, slice.Tokens(gotJSON))
		}
		if !resp.OverBudget && (resp.Tokens > budget || slice.Tokens(md.Bytes()) > budget) {
			t.Errorf("budget %d: JSON %d tokens, markdown %d, not over budget", budget, resp.Tokens, slice.Tokens(md.Bytes()))
		}
	}
}

// fitDirect slices ref in as-is straight from the index, in both formats.
func fitDirect(t *testing.T, ctx *Ctx, ref string, budget int) (js, md []byte) {
	t.Helper()
	x, err := openIndex(ctx.Root)
	if err != nil {
		t.Fatal(err)
	}
	defer x.Close()
	infos, err := x.ChangeSets()
	if err != nil {
		t.Fatal(err)
	}
	var all []id.ID
	for _, c := range infos {
		all = append(all, c.ID)
	}
	n, err := x.Resolve(index.AsIs, ref)
	if err != nil {
		t.Fatal(err)
	}
	r := slice.FromIndex(x, index.AsIs)
	s, err := slice.Build(r, r, slice.Request{Focus: n.ID, View: slice.ViewAsIs, ChangeSets: all})
	if err != nil {
		t.Fatal(err)
	}
	opts := slice.Options{Budget: budget, DefaultView: slice.ViewAsIs}
	m, err := slice.Fit(s, opts)
	if err != nil {
		t.Fatal(err)
	}
	opts.JSON = true
	j, err := slice.Fit(s, opts)
	if err != nil {
		t.Fatal(err)
	}
	return j.Bytes, m.Bytes
}

func itoa(n int) string {
	b, _ := canon.Marshal(n)
	return string(b)
}

func TestSliceViews(t *testing.T) {
	ctx := sliceCtx(t)
	sliceModel(t, ctx)
	cs := newChange(t, ctx, `{"intent":"Add an index"}`).ChangeSet
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"index","parent":"store"}`)

	// A draft change set: the default is what-if, where the new node is.
	resp, err := callSlice(t, ctx, `{"ref":"index"}`)
	if err != nil {
		t.Fatal(err)
	}
	if resp.View != slice.ViewWhatIf || resp.Change == nil || resp.Change.ID != cs {
		t.Fatalf("view %s, change %+v; want what-if of %s", resp.View, resp.Change, cs)
	}
	if len(resp.Focus) != 1 || resp.Focus[0].Handle != "index" || !reflect.DeepEqual(resp.Focus[0].Proposed, []string{"new"}) {
		t.Errorf("focus %+v; want index, new", resp.Focus)
	}
	if len(resp.Constraints) != 1 || resp.Constraints[0].From != "store" {
		t.Errorf("constraints %+v; want no-network from store", resp.Constraints)
	}

	// The change slice: no ref.
	resp, err = callSlice(t, ctx, `{}`)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Ref != resp.Change.Short || len(resp.Focus) != 1 {
		t.Errorf("change slice ref %q, focus %d; want %q, 1", resp.Ref, len(resp.Focus), resp.Change.Short)
	}

	// as-is lacks the new node.
	_, err = callSlice(t, ctx, `{"ref":"index","view":"as-is"}`)
	wantExit(t, err, 4)

	// --change by prefix gives the same change set.
	c := *ctx
	c.ChangeRef = strings.ToLower(cs.String()[:24])
	resp, err = callSlice(t, &c, `{"ref":"log"}`)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Change == nil || resp.Change.ID != cs || resp.View != slice.ViewWhatIf {
		t.Errorf("change %+v view %s; want %s what-if", resp.Change, resp.View, cs)
	}
}

func TestSliceErrors(t *testing.T) {
	ctx := sliceCtx(t)
	sliceModel(t, ctx)

	o, _ := Lookup("slice")
	for _, req := range []string{`{"budget":0}`, `{"budget":-3}`} {
		_, err := DecodeRequest(o, []byte(req))
		wantExit(t, err, 2)
		if !strings.Contains(err.Error(), "budget must be at least 1") {
			t.Errorf("%s: %v", req, err)
		}
	}
	if _, err := o.Run(ctx, SliceRequest{Ref: "log", Budget: -1}); ExitCode(err) != 2 {
		t.Errorf("budget -1: %v", err)
	}

	_, err := callSlice(t, ctx, `{"ref":"nope"}`)
	wantExit(t, err, 4)
	if AsError(err).Message != `no node matches "nope" in the as-is view` {
		t.Errorf("unknown ref: %v", err)
	}

	_, err = callSlice(t, ctx, `{}`)
	wantExit(t, err, 2)
	if AsError(err).Next != "rhiza slice <ref>" {
		t.Errorf("no change set: next %q", AsError(err).Next)
	}

	_, err = callSlice(t, ctx, `{"ref":"log","view":"what-if"}`)
	wantExit(t, err, 2)
	if AsError(err).Next == "" {
		t.Error("what-if without a change set has no next")
	}

	_, err = callSlice(t, ctx, `{"ref":"log","view":"sideways"}`)
	wantExit(t, err, 2)

	recordWrite(t, ctx, "retire", `{"ref":"use-files","reason":"Replaced"}`)
	_, err = callSlice(t, ctx, `{"ref":"use-files"}`)
	wantExit(t, err, 4)

	c := *ctx
	c.ChangeRef = "ZZZZZZ"
	_, err = callSlice(t, &c, `{"ref":"log"}`)
	wantExit(t, err, 4)

	out := testCtx(t.TempDir())
	_, err = callSlice(t, out, `{"ref":"log"}`)
	wantExit(t, err, 4)
	if AsError(err).Next != "rhiza init" {
		t.Errorf("outside a repo: next %q", AsError(err).Next)
	}
}

// resolveInModel follows index.Resolve's rules.
func TestSliceResolveInModel(t *testing.T) {
	ctx := sliceCtx(t)
	sliceModel(t, ctx)
	x, err := openIndex(ctx.Root)
	if err != nil {
		t.Fatal(err)
	}
	defer x.Close()
	m, _, err := x.WhatIf(id.ID{})
	if err != nil {
		t.Fatal(err)
	}
	n, err := x.Resolve(index.AsIs, "log")
	if err != nil {
		t.Fatal(err)
	}
	for _, ref := range []string{"log", n.ID.String(), n.ID.String()[:20]} {
		got, err := resolveInModel(m, ref)
		if err != nil || got.ID != n.ID {
			t.Errorf("%s: %v, %v", ref, got.ID, err)
		}
	}
	if _, err := resolveInModel(m, "nope"); err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("nope: %v", err)
	}
	// Every node shares the first characters of its ID here.
	if _, err := resolveInModel(m, n.ID.String()[:4]); err == nil || !strings.Contains(err.Error(), "ambiguous") {
		t.Errorf("short prefix: %v", err)
	}
}
