package api

import (
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/id"
)

// lsCall calls ls with req.
func lsCall(t *testing.T, ctx *Ctx, req string) LsResponse {
	t.Helper()
	resp, err := callOp(t, ctx, "ls", req)
	if err != nil {
		t.Fatalf("ls %s: %v", req, err)
	}
	return resp.(LsResponse)
}

// lsHandles is the handles of a page, a retired one marked with "!".
func lsHandles(r LsResponse) string {
	hs := make([]string, len(r.Nodes))
	for i, n := range r.Nodes {
		hs[i] = n.Handle
		if n.Retired {
			hs[i] += "!"
		}
	}
	return strings.Join(hs, ",")
}

func TestLsFilters(t *testing.T) {
	ctx, _, _ := readLog(t)
	all := lsCall(t, ctx, `{"view":"to-be"}`)
	// The test actor, then the nodes in creation (ID) order, live only (under-old lives on under retired old).
	if got, want := lsHandles(all), "test,store,cache,tape,under-old,faster,use-ssd,ram"; got != want {
		t.Errorf("ls = %s, want %s", got, want)
	}
	if all.NextCursor != "" || all.View != "to-be" {
		t.Errorf("next %q view %q", all.NextCursor, all.View)
	}
	for i := 1; i < len(all.Nodes); i++ {
		if id.Compare(all.Nodes[i-1].ID, all.Nodes[i].ID) >= 0 {
			t.Errorf("not sorted by ID at %d", i)
		}
	}
	if got := lsHandles(lsCall(t, ctx, `{"view":"to-be","retired":true}`)); got != "test,store,cache,disk!,tape,old!,under-old,faster,use-ssd,ram" {
		t.Errorf("--retired = %s", got)
	}
	if got := lsHandles(lsCall(t, ctx, `{"view":"to-be","kinds":["intent","decision"]}`)); got != "faster,use-ssd" {
		t.Errorf("--kind = %s", got)
	}
	if got := lsHandles(lsCall(t, ctx, `{"view":"to-be","parent":"cache"}`)); got != "tape,ram" {
		t.Errorf("--parent cache = %s", got)
	}
	if got := lsHandles(lsCall(t, ctx, `{"view":"as-is","parent":"store","retired":true}`)); got != "cache,disk,old" {
		t.Errorf("as-is --parent store = %s", got)
	}
	r := lsCall(t, ctx, `{"view":"to-be","parent":"store"}`)
	if len(r.Nodes) != 1 || r.Nodes[0].Parent.IsZero() {
		t.Errorf("parent field = %+v", r.Nodes)
	}
	if got := lsHandles(lsCall(t, ctx, `{"view":"to-be","kinds":["question"]}`)); got != "" {
		t.Errorf("no questions, got %s", got)
	}
}

func TestLsPaging(t *testing.T) {
	ctx, _, _ := readLog(t)
	want := lsHandles(lsCall(t, ctx, `{"view":"to-be","retired":true}`))
	var got []string
	cursor := ""
	for pages := 0; ; pages++ {
		if pages > 20 {
			t.Fatal("paging does not end")
		}
		r := lsCall(t, ctx, `{"view":"to-be","retired":true,"limit":3,"cursor":"`+cursor+`"}`)
		if len(r.Nodes) > 3 {
			t.Fatalf("page of %d", len(r.Nodes))
		}
		got = append(got, lsHandles(r))
		if r.NextCursor == "" {
			break
		}
		if len(r.Nodes) != 3 {
			t.Errorf("short page %d with a next cursor", len(r.Nodes))
		}
		if c := lsEncodeCursor(r.Nodes[2].ID); c != r.NextCursor {
			t.Errorf("cursor %q, want %q", r.NextCursor, c)
		}
		h := human(t, r)
		if !strings.Contains(h, "\nnext: rhiza ls --retired --limit 3 --view to-be --cursor "+r.NextCursor+"\n") {
			t.Errorf("human next:\n%s", h)
		}
		cursor = r.NextCursor
	}
	if strings.Join(got, ",") != want {
		t.Errorf("pages %v, want %s", got, want)
	}
	// Exactly a page's worth: no next cursor.
	if r := lsCall(t, ctx, `{"view":"to-be","limit":8}`); len(r.Nodes) != 8 || r.NextCursor != "" {
		t.Errorf("limit 8 = %d nodes, next %q", len(r.Nodes), r.NextCursor)
	}
}

func TestLsHuman(t *testing.T) {
	ctx, _, _ := readLog(t)
	r := lsCall(t, ctx, `{"view":"to-be","retired":true,"kinds":["capability"],"parent":"store"}`)
	h := human(t, r)
	lines := strings.Split(strings.TrimSuffix(h, "\n"), "\n")
	if len(lines) != 3 || !strings.HasPrefix(lines[0], "- capability cache ") ||
		lines[1] != "- capability disk "+r.Nodes[1].ID.String()+" [retired]" || strings.Contains(h, "next:") {
		t.Errorf("human:\n%s", h)
	}
	if h := human(t, lsCall(t, ctx, `{"view":"to-be","kinds":["question"]}`)); h != "No nodes match.\n\nnext: rhiza add\n" {
		t.Errorf("empty human %q", h)
	}
	r = lsCall(t, ctx, `{"view":"to-be","kinds":["capability","intent"],"parent":"cache","limit":1}`)
	if h := human(t, r); !strings.HasSuffix(h, "\nnext: rhiza ls --kind capability --kind intent --parent cache --limit 1 --view to-be --cursor "+r.NextCursor+"\n") {
		t.Errorf("next line:\n%s", h)
	}
}

func TestLsErrors(t *testing.T) {
	ctx, _, _ := readLog(t)
	for _, req := range []string{
		`{"limit":5000}`,
		`{"limit":-1}`,
		`{"cursor":"!!"}`,
		`{"cursor":"` + lsEncodeCursor(id.ID{})[:4] + `x"}`,
		`{"cursor":"bm90LWFuLWlk"}`,
		`{"kinds":["widget"]}`,
	} {
		_, err := callOp(t, ctx, "ls", req)
		wantExit(t, err, 2)
		wantNext(t, err, "rhiza describe ls")
	}
	_, err := callOp(t, ctx, "ls", `{"parent":"nosuch"}`)
	wantExit(t, err, 4)
	wantNext(t, err, "rhiza ls")
	_, err = callOp(t, ctx, "ls", `{"view":"what-if:ZZZZ"}`)
	wantExit(t, err, 4)
	if r := lsCall(t, ctx, `{"limit":1000}`); r.NextCursor != "" {
		t.Errorf("limit 1000 next %q", r.NextCursor)
	}
}

func TestLsWhatIfDraft(t *testing.T) {
	ctx := readCtx(t)
	cs := newChange(t, ctx, `{"intent":"X"}`).ChangeSet
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
	r := lsCall(t, ctx, `{"kinds":["capability"]}`)
	if r.View != "what-if:"+cs.String() || lsHandles(r) != "store" {
		t.Errorf("draft ls = %+v", r)
	}
	if got := lsHandles(lsCall(t, ctx, `{"kinds":["capability"],"view":"as-is"}`)); got != "" {
		t.Errorf("as-is ls = %s", got)
	}
}
