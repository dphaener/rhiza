package web_test

import (
	"encoding/json"
	"fmt"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
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
	callJSON(t, ctx, "add", `{"kind":"intent","handle":"areas","fields":["statement=Product areas","umbrella=yes"]}`)
	callJSON(t, ctx, "land", `{}`)
	plan(t, ctx, "Alpha first", "alpha", true)
	plan(t, ctx, "Alpha second", "alpha", false)
	plan(t, ctx, "Beta only", "beta", false)
	plan(t, ctx, "Loose work", "", true)
	plan(t, ctx, "Area upkeep", "areas", false)

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
	if len(rm.Groups) != 4 || rm.Groups[0].Kind != "milestone" || rm.Groups[2].Kind != "backlog" || rm.Groups[3].Kind != "umbrella" {
		t.Fatalf("fixture roadmap: %+v", rm.Groups)
	}
	page := serve(h, "GET", "/roadmap")
	if page.Code != 200 {
		t.Fatalf("/roadmap: %d %s", page.Code, page.Body)
	}
	body := page.Body.String()
	// The titles, in the operation's order: milestone titles and change sets.
	pos := -1
	order := []string{"Alpha milestone", "Alpha first", "Alpha second", "Beta milestone", "Beta only", "Not in a milestone", "Loose work", "Umbrella", "Product areas", "Area upkeep"}
	for _, s := range order {
		i := strings.Index(body[pos+1:], s)
		if i < 0 {
			t.Fatalf("%q missing or out of order in:\n%s", s, body)
		}
		pos += 1 + i
	}
	for _, g := range rm.Groups {
		for _, c := range g.ChangeSets {
			link := `href="/changes/` + c.ID + `"`
			if !strings.Contains(body, link) || !strings.Contains(body, "rz-pill--"+c.State) {
				t.Errorf("%s (%s): no %s or state pill", c.Intent, c.State, link)
			}
			if strings.Contains(body, ">"+c.ID+"<") || strings.Contains(body, ">"+c.ID[:8]) {
				t.Errorf("%s: the page shows an ID", c.Intent)
			}
		}
	}
	// The umbrella is not counted as a milestone.
	if !strings.Contains(body, "2 milestones · 5 change sets · 1 in backlog") {
		t.Errorf("kicker miscounts the groups:\n%s", body)
	}
	for _, h := range []string{"alpha", "beta", "areas"} {
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

func TestRoadmapRowShowsTitleThenIntent(t *testing.T) {
	ctx := fixtureCtx(t)
	intent := "Rework how rows are listed, which is long enough that it would run right off the end of a phone"
	cs := newChange(t, ctx, intent, "intention")
	_ = cs
	callJSON(t, ctx, "propose", `{}`)
	callJSON(t, ctx, "accept", `{}`)
	body := serve(web.Handler(api.ReadOnlyHandler(ctx)), "GET", "/roadmap").Body.String()
	link := `>Rework how rows are listed<small>` + intent + `</small></a>`
	if !strings.Contains(body, link) {
		t.Errorf("row is not title then intent in the secondary line:\n%s", body)
	}
	css, err := os.ReadFile("static/css/rhiza.css")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(css), ".rz-row-title small {") || !strings.Contains(string(css), "-webkit-line-clamp: 2") {
		t.Error("stylesheet does not clamp the secondary line to two lines")
	}
}

// milestoneProject is a project with a milestone in every state of work: A
// has a draft, a proposed, two accepted (one with a run going), seven landed
// and a withdrawn change set; D is finished.
func milestoneProject(t *testing.T, seed int64, prefix string) (http.Handler, *api.Ctx) {
	t.Helper()
	ctx := fixtureCtx(t)
	ctx.IDs = id.NewGenerator(ctx.Now, rand.New(rand.NewSource(seed)))
	newChange(t, ctx, "Milestones", "record")
	for _, w := range []string{
		`{"kind":"intent","handle":"a","fields":["statement=` + prefix + `Alpha milestone","milestone=1"]}`,
		`{"kind":"intent","handle":"a1","parent":"a","fields":["statement=` + prefix + `Alpha done one"]}`,
		`{"kind":"intent","handle":"a2","parent":"a","fields":["statement=` + prefix + `Alpha done two"]}`,
		`{"kind":"intent","handle":"d","fields":["statement=` + prefix + `Delta milestone","milestone=4"]}`,
		`{"kind":"intent","handle":"d1","parent":"d","fields":["statement=` + prefix + `Delta done one"]}`,
	} {
		callJSON(t, ctx, "add", w)
	}
	callJSON(t, ctx, "land", `{}`)
	mk := func(title, serves string, steps ...string) string {
		cs := newChange(t, ctx, prefix+title, "intention")
		callJSON(t, ctx, "link", `{"from":"`+cs+`","type":"serves","to":"`+serves+`"}`)
		for _, s := range steps {
			callJSON(t, ctx, s, `{}`)
		}
		return cs
	}
	mk("A draft", "a")
	mk("A proposed", "a", "propose")
	building := mk("A building", "a2", "propose", "accept")
	mk("A waiting", "a2", "propose", "accept")
	mk("A withdrawn", "a1", "propose", "withdraw")
	for i := 1; i <= 7; i++ {
		mk(fmt.Sprintf("A landed %d", i), "a1", "propose", "accept", "land")
	}
	mk("D landed", "d1", "propose", "accept", "land")
	rs := run.NewStore(run.Dir(filepath.Join(ctx.Root, ".git")))
	runID := fmt.Sprintf("01RUN%d", seed)
	if err := rs.Create(runID, run.Entry{Type: run.TypeStarted, Started: &run.Started{Workflow: "build", Change: building, Title: prefix + "A building"}}); err != nil {
		t.Fatal(err)
	}
	if err := rs.Append(runID, run.Entry{Type: run.TypeEntered, Step: "leaves", Visit: 1}); err != nil {
		t.Fatal(err)
	}
	return web.HandlerWithRuns(nil, api.WebRunsSource(ctx), api.ReadOnlyHandler(ctx)), ctx
}

func TestRoadmapMilestones(t *testing.T) {
	ha, _ := milestoneProject(t, 21, "")
	hb, _ := milestoneProject(t, 22, "Other ")
	list := func() []web.ProjectView {
		return []web.ProjectView{{Name: "a", Initialised: true}, {Name: "b", Initialised: true}}
	}
	h := web.RootWithProjects(web.Handler(nil), func(n string) http.Handler {
		if n == "a" {
			return ha
		}
		return hb
	}, list)
	rec := serve(h, "GET", "/p/a/roadmap")
	if rec.Code != 200 {
		t.Fatalf("roadmap: %d %s", rec.Code, rec.Body)
	}
	body := rec.Body.String()
	if strings.Contains(body, "Other ") {
		t.Error("a row of the other project")
	}
	if !strings.Contains(body, "<summary>1 finished milestone</summary>") || !strings.Contains(body, "Delta milestone") {
		t.Error("the finished milestone is not in the page's fold")
	}
	open := body[:strings.Index(body, "finished milestone")]
	if strings.Contains(open, "Delta") {
		t.Error("a finished milestone among the open ones")
	}
	fin := body[strings.Index(body, "finished milestone"):]
	if strings.Contains(fin, "href=") && strings.Index(fin, "href=") < strings.Index(fin, "</details>") {
		t.Error("the finished fold links to another page")
	}
	pos := -1
	for _, s := range []string{"Milestone · building", "Alpha milestone", "Building now", "A building", "Accepted, waiting", "A waiting",
		"Proposed", "A proposed", "Landed", "A landed 1", "Drafts", "A draft"} {
		i := strings.Index(body[pos+1:], s)
		if i < 0 {
			t.Fatalf("%q missing or out of order in:\n%s", s, body)
		}
		pos += 1 + i
	}
	// 7 landed, 2 accepted, 1 proposed, 1 withdrawn in the bar and legend.
	for _, s := range []string{`aria-label="7 landed, 2 accepted, 1 proposed, 1 withdrawn"`, `flex:7`, `<b>7</b> landed`, `<b>1</b> withdrawn`} {
		if !strings.Contains(body, s) {
			t.Errorf("bar lacks %q", s)
		}
	}
	// A change set whose run is going is under building now, not accepted.
	now := body[strings.Index(body, "Building now"):strings.Index(body, "Accepted, waiting")]
	if !strings.Contains(now, "A building") || strings.Contains(now, "A waiting") {
		t.Errorf("building now is wrong: %s", now)
	}
	// Done points in a closed fold: one of two met.
	if !strings.Contains(body, "<details class=\"rz-fold\"><summary>Done points · 1 of 2 met") || strings.Contains(body, "<details open") {
		t.Error("done points are not in a closed fold")
	}
	// Past five landed, the rest fold in place with the withdrawn one.
	if !strings.Contains(body, "<summary>2 more landed, 1 withdrawn</summary>") {
		t.Errorf("landed fold is wrong:\n%s", body)
	}
	// Title from the title field, the intent as the clamped sub-line.
	if !strings.Contains(body, `>A building<small>A building</small></a>`) {
		t.Error("row is not title then intent")
	}
	if strings.Contains(body, ">01") || strings.Contains(body, "Alpha milestone</h3></div>") == false {
		t.Error("an ID shows, or the milestone is not named by its statement")
	}
	for _, w := range []string{"01RUN21", "01RUN22"} {
		if strings.Contains(body, w) {
			t.Errorf("run ID %s shown", w)
		}
	}
}
