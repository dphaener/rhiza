package web_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/web"
)

// callJSON runs an operation on ctx and returns the ID it made, if any.
func callJSON(t *testing.T, ctx *api.Ctx, op, req string) any {
	t.Helper()
	o, ok := api.Lookup(op)
	if !ok {
		t.Fatalf("no op %s", op)
	}
	r, err := api.DecodeRequest(o, []byte(req))
	if err != nil {
		t.Fatalf("%s %s: %v", op, req, err)
	}
	resp, err := o.Run(ctx, r)
	if err != nil {
		t.Fatalf("%s %s: %v", op, req, err)
	}
	return resp
}

func newChange(t *testing.T, ctx *api.Ctx, intent, kind string) string {
	t.Helper()
	b, _ := json.Marshal(callJSON(t, ctx, "change.new", `{"intent":`+quote(intent)+`,"kind":"`+kind+`"}`))
	var v struct {
		CS string `json:"change_set"`
	}
	json.Unmarshal(b, &v)
	return v.CS
}

func quote(s string) string { b, _ := json.Marshal(s); return string(b) }

// plan opens an intention that serves the intent (none for the backlog),
// proposes it and, if accept, accepts it.
func plan(t *testing.T, ctx *api.Ctx, intent, serves string, accept bool) {
	t.Helper()
	cs := newChange(t, ctx, intent, "intention")
	if serves != "" {
		callJSON(t, ctx, "link", `{"from":"`+cs+`","type":"serves","to":"`+serves+`"}`)
	}
	callJSON(t, ctx, "propose", `{}`)
	if accept {
		callJSON(t, ctx, "accept", `{}`)
	}
}

func TestRoadmapPageMatchesTheOperation(t *testing.T) {
	ctx := fixtureCtx(t)
	newChange(t, ctx, "Milestones", "record")
	callJSON(t, ctx, "add", `{"kind":"intent","handle":"alpha","fields":["statement=Alpha milestone","milestone=1"]}`)
	callJSON(t, ctx, "add", `{"kind":"intent","handle":"beta","fields":["statement=Beta milestone","milestone=2"]}`)
	callJSON(t, ctx, "land", `{}`)
	plan(t, ctx, "Alpha first", "alpha", true)
	plan(t, ctx, "Alpha second", "alpha", false)
	plan(t, ctx, "Beta only", "beta", false)
	plan(t, ctx, "Loose work", "", true)

	h := web.Handler(api.ReadOnlyHandler(ctx))
	var rm struct {
		Groups []struct {
			Kind       string `json:"kind"`
			Intent     struct{ ID string }
			ChangeSets []struct {
				ID, State, Intent string
			} `json:"change_sets"`
		}
	}
	if err := json.Unmarshal([]byte(want(t, ctx, "roadmap", `{"all":true}`)), &rm); err != nil {
		t.Fatal(err)
	}
	if len(rm.Groups) != 3 || rm.Groups[0].Kind != "milestone" || rm.Groups[2].Kind != "backlog" {
		t.Fatalf("fixture roadmap: %+v", rm.Groups)
	}
	page := serve(h, "GET", "/roadmap")
	if page.Code != 200 {
		t.Fatalf("/roadmap: %d %s", page.Code, page.Body)
	}
	body := page.Body.String()
	// The titles, in the operation's order: milestone titles and change sets.
	pos := -1
	order := []string{"Alpha milestone", "Alpha first", "Alpha second", "Beta milestone", "Beta only", "Not in a milestone", "Loose work"}
	for _, s := range order {
		i := strings.Index(body[pos+1:], s)
		if i < 0 {
			t.Fatalf("%q missing or out of order in:\n%s", s, body)
		}
		pos += 1 + i
	}
	for _, g := range rm.Groups {
		for _, c := range g.ChangeSets {
			link := `href="/cs/` + c.ID + `"`
			if !strings.Contains(body, link) || !strings.Contains(body, "rz-pill--"+c.State) {
				t.Errorf("%s (%s): no %s or state pill", c.Intent, c.State, link)
			}
			if strings.Contains(body, ">"+c.ID+"<") || strings.Contains(body, ">"+c.ID[:8]) {
				t.Errorf("%s: the page shows an ID", c.Intent)
			}
		}
	}
	for _, h := range []string{"alpha", "beta"} {
		if strings.Contains(body, ">"+h+"<") {
			t.Errorf("handle %s shown", h)
		}
	}
	if !strings.Contains(body, `data-rz-fragment="/fragments/roadmap"`) || !strings.Contains(body, "aria-current=\"page\">Roadmap") {
		t.Error("no live section or nav mark")
	}

	// The live section moves when a change set does: a fragment re-fetch.
	before := serve(h, "GET", web.FragmentURL("roadmap")).Body.String()
	plan(t, ctx, "Late arrival", "beta", false)
	after := serve(h, "GET", web.FragmentURL("roadmap")).Body.String()
	if before == after || !strings.Contains(after, "Late arrival") || strings.Contains(before, "Late arrival") {
		t.Error("fragment did not follow the new change set")
	}
}

func TestRoadmapPageEmptyAndNoAPI(t *testing.T) {
	ctx := fixtureCtx(t)
	body := serve(web.Handler(api.ReadOnlyHandler(ctx)), "GET", "/roadmap")
	if body.Code != 200 || !strings.Contains(body.Body.String(), "Nothing is planned") {
		t.Errorf("empty: %d %s", body.Code, body.Body)
	}
	if rec := serve(web.Handler(nil), "GET", "/roadmap"); rec.Code != 500 {
		t.Errorf("no api: %d", rec.Code)
	}
	if rec := serve(web.Handler(api.ReadOnlyHandler(ctx)), "POST", "/roadmap"); rec.Code != 405 {
		t.Errorf("POST: %d", rec.Code)
	}
}
