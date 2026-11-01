package api

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/index"
)

// showCall calls show with req.
func showCall(t *testing.T, ctx *Ctx, req string) ShowResponse {
	t.Helper()
	resp, err := callOp(t, ctx, "show", req)
	if err != nil {
		t.Fatalf("show %s: %v", req, err)
	}
	return resp.(ShowResponse)
}

func TestShowNode(t *testing.T) {
	ctx, landed, _ := readLog(t)
	r := showCall(t, ctx, `{"ref":"store","view":"to-be"}`)
	if r.Kind != "node" || r.Link != nil || r.Node == nil || r.View != "to-be" {
		t.Fatalf("show store = %+v", r)
	}
	n := r.Node
	if n.Kind != "capability" || n.Handle != "store" || len(n.Path) != 0 || n.Retired || n.Created != landed {
		t.Errorf("node = %+v", n)
	}
	if got := string(n.Fields["statement"]); got != `"Keep more"` {
		t.Errorf("statement = %s", got)
	}
	var kids []string
	for _, c := range n.Children {
		kids = append(kids, c.Handle+map[bool]string{true: "!", false: ""}[c.Retired])
	}
	if got := strings.Join(kids, ","); got != "cache,disk!,old!" {
		t.Errorf("children = %s", got)
	}

	tape := showCall(t, ctx, `{"ref":"tape","view":"to-be"}`).Node
	if len(tape.Path) != 2 || tape.Path[0].Handle != "store" || tape.Path[1].Handle != "cache" {
		t.Errorf("tape path = %+v", tape.Path)
	}

	faster := showCall(t, ctx, `{"ref":"faster","view":"to-be"}`).Node
	if len(faster.LinksOut) != 0 || len(faster.LinksIn) != 3 {
		t.Fatalf("faster links = %+v / %+v", faster.LinksOut, faster.LinksIn)
	}
	retired := 0
	for _, l := range faster.LinksIn {
		if l.Name != "served by" || l.Type != "serves" || l.Other.Handle == "" {
			t.Errorf("link in = %+v", l)
		}
		if l.Retired {
			retired++
			if l.Other.Handle != "tape" {
				t.Errorf("retired link from %s", l.Other.Handle)
			}
		}
	}
	if retired != 1 {
		t.Errorf("%d retired links in", retired)
	}

	disk := showCall(t, ctx, `{"ref":"`+n.Children[1].ID.String()+`","view":"to-be"}`).Node
	if !disk.Retired || disk.Reason != "gone" || len(disk.LinksOut) != 1 || disk.LinksOut[0].Name != "decided-by" {
		t.Errorf("disk = %+v", disk)
	}
	if got := string(disk.Fields["acceptance"]); got != "{\n  \"gb\": 2\n}" {
		t.Errorf("size = %s", got)
	}
	use := showCall(t, ctx, `{"ref":"use-ssd","view":"to-be"}`).Node
	if len(use.LinksIn) != 1 || use.LinksIn[0].Name != "decides" || use.LinksIn[0].Other.Handle != "disk" {
		t.Errorf("use-ssd links in = %+v", use.LinksIn)
	}

	// The default view in an accepted change set is to-be; as-is has no ram.
	if r := showCall(t, ctx, `{"ref":"ram"}`); r.View != "to-be" {
		t.Errorf("default view %q", r.View)
	}
	_, err := callOp(t, ctx, "show", `{"ref":"ram","view":"as-is"}`)
	wantExit(t, err, 4)
}

func TestShowLink(t *testing.T) {
	ctx, landed, _ := readLog(t)
	cache := showCall(t, ctx, `{"ref":"cache","view":"as-is"}`).Node
	if len(cache.LinksOut) != 1 {
		t.Fatalf("cache links out = %+v", cache.LinksOut)
	}
	lid := cache.LinksOut[0].ID
	for _, ref := range []string{lid.String(), lid.String()[:25]} {
		r := showCall(t, ctx, `{"ref":"`+ref+`","view":"as-is"}`)
		if r.Kind != "link" || r.Node != nil || r.Link == nil {
			t.Fatalf("show %s = %+v", ref, r)
		}
		l := r.Link
		if l.ID != lid || l.Type != "serves" || l.From.Handle != "cache" || l.To.Handle != "faster" || l.Created != landed || l.Retired {
			t.Errorf("link = %+v", l)
		}
	}
	h := human(t, showCall(t, ctx, `{"ref":"`+lid.String()+`","view":"as-is"}`))
	for _, want := range []string{"# link " + shortID(lid) + "\n", "- Type: serves\n", "- From: cache (", "- To: faster ("} {
		if !strings.Contains(h, want) {
			t.Errorf("human missing %q:\n%s", want, h)
		}
	}
}

func TestShowHuman(t *testing.T) {
	ctx, _, _ := readLog(t)
	h := human(t, showCall(t, ctx, `{"ref":"store","view":"to-be"}`))
	for _, want := range []string{
		"# capability store\n\n- ID: ",
		"- View: to-be\n\n## Fields\n\n### statement\n\nKeep more\n",
		"## Children\n\n- capability cache (",
		") [retired]\n",
	} {
		if !strings.Contains(h, want) {
			t.Errorf("human missing %q:\n%s", want, h)
		}
	}
	if strings.Contains(h, "## Links") || strings.Contains(h, "- Path:") {
		t.Errorf("empty sections shown:\n%s", h)
	}
	if h2 := human(t, showCall(t, ctx, `{"ref":"store","view":"to-be"}`)); h2 != h {
		t.Error("human output is not deterministic")
	}
	h = human(t, showCall(t, ctx, `{"ref":"faster","view":"to-be"}`))
	if !strings.Contains(h, "## Links\n\n- served by cache (") || !strings.Contains(h, "- served by tape (") ||
		!strings.Contains(h, "[retired]") || strings.Contains(h, "## Children") || strings.Contains(h, "## Fields") {
		t.Errorf("faster human:\n%s", h)
	}
	h = human(t, showCall(t, ctx, `{"ref":"tape","view":"to-be"}`))
	if !strings.Contains(h, "- Path: store / cache\n") {
		t.Errorf("tape human:\n%s", h)
	}
	h = human(t, showCall(t, ctx, `{"ref":"disk","view":"as-is"}`))
	if !strings.Contains(h, "### acceptance\n\n{\n  \"gb\": 2\n}\n") || !strings.Contains(h, "- decided-by use-ssd (") {
		t.Errorf("disk human:\n%s", h)
	}
	oldID := showCall(t, ctx, `{"ref":"store","view":"to-be"}`).Node.Children[2].ID
	h = human(t, showCall(t, ctx, `{"ref":"`+oldID.String()+`","view":"to-be"}`))
	if !strings.Contains(h, "# capability old\n") || !strings.Contains(h, "- Retired: replaced\n") ||
		!strings.Contains(h, "## Children\n\n- capability under-old (") {
		t.Errorf("old human:\n%s", h)
	}
}

func TestShowJSON(t *testing.T) {
	ctx, _, _ := readLog(t)
	b, err := json.Marshal(showCall(t, ctx, `{"ref":"use-ssd","view":"as-is"}`))
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, want := range []string{`"kind":"node"`, `"view":"as-is"`, `"path":[]`, `"fields":{}`, `"links_out":[]`, `"children":[]`, `"name":"decides"`} {
		if !strings.Contains(s, want) {
			t.Errorf("json missing %s: %s", want, s)
		}
	}
	if strings.Contains(s, `"link":`) {
		t.Errorf("node response has a link: %s", s)
	}
}

func TestShowErrors(t *testing.T) {
	ctx, _, _ := readLog(t)
	_, err := callOp(t, ctx, "show", `{"ref":"nosuch"}`)
	wantExit(t, err, 4)
	wantNext(t, err, "rhiza ls")
	if msg := AsError(err).Message; msg != `no node or link matches "nosuch"` {
		t.Errorf("message %q", msg)
	}
	// A retired node's handle no longer resolves.
	_, err = callOp(t, ctx, "show", `{"ref":"old","view":"to-be"}`)
	wantExit(t, err, 4)
	_, err = callOp(t, ctx, "show", `{"ref":" "}`)
	wantExit(t, err, 2)
	_, err = callOp(t, ctx, "show", `{"ref":"store","view":"sideways"}`)
	wantExit(t, err, 2)
	wantNext(t, err, "rhiza describe show")

	// An ambiguous node prefix stays a usage error.
	r := mustReader(t, ctx, "to-be", defaultView)
	nodes, _ := r.List(index.ListQuery{Limit: 100})
	_, err = callOp(t, ctx, "show", `{"ref":"`+nodes[0].ID.String()[:4]+`","view":"to-be"}`)
	wantExit(t, err, 2)
}

func TestShowWhatIfDraft(t *testing.T) {
	ctx := readCtx(t)
	cs := newChange(t, ctx, `{"intent":"X"}`).ChangeSet
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
	r := showCall(t, ctx, `{"ref":"store"}`)
	if r.View != "what-if:"+cs.String() || r.Node.Handle != "store" || r.Node.Created != cs {
		t.Errorf("draft show = %+v", r)
	}
	_, err := callOp(t, ctx, "show", `{"ref":"store","view":"as-is"}`)
	wantExit(t, err, 4)
}

func TestShowEnds(t *testing.T) {
	ctx, landed, _ := readLog(t)
	r := mustReader(t, ctx, "as-is", defaultView)
	e := &showEnds{r: r}
	got, err := e.end(landed)
	if err != nil || !got.ChangeSet || got.Intent != "Base" || got.Handle != "" {
		t.Errorf("end(change set) = %+v, %v", got, err)
	}
	if s := showEndName(got); s != "change set "+shortID(landed)+": Base" {
		t.Errorf("name %q", s)
	}
	var zero ShowEnd
	zero.ID = landed
	if s := showEndName(zero); s != landed.String() {
		t.Errorf("bare name %q", s)
	}
}
