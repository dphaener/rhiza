package web_test

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/web"
)

func TestChangeSetPage(t *testing.T) {
	ctx := fixtureCtx(t)
	// The done point the landed change set serves.
	newChange(t, ctx, "Milestones", "record")
	callJSON(t, ctx, "add", `{"kind":"intent","handle":"reliable","fields":["statement=Deliveries are reliable","milestone=1"]}`)
	callJSON(t, ctx, "land", `{}`)

	// A dependency, landed: three verified leaves, a PR, why it exists.
	base := newChange(t, ctx, "Base work", "intention")
	callJSON(t, ctx, "link", `{"from":"`+base+`","type":"serves","to":"reliable"}`)
	for _, l := range []string{"one", "two", "three"} {
		callJSON(t, ctx, "add", `{"kind":"capability","handle":"base-`+l+`","fields":["statement=Base capability `+l+`","acceptance=Base check `+l+`"]}`)
		callJSON(t, ctx, "add", `{"kind":"evidence","handle":"base-`+l+`-ok","fields":["statement=Ran check `+l+`","result=PASS","status=verified"]}`)
		callJSON(t, ctx, "link", `{"from":"base-`+l+`-ok","type":"verifies","to":"base-`+l+`"}`)
	}
	callJSON(t, ctx, "add", `{"kind":"evidence","handle":"base-pr","fields":["statement=Base PR passed CI","result=ok","pr=https://github.com/x/y/pull/498","status=verified"]}`)
	callJSON(t, ctx, "add", `{"kind":"decision","handle":"base-dec","fields":["statement=Base work ships first","question=What ships first?","wording=Base first"]}`)
	callJSON(t, ctx, "link", `{"from":"base-one","type":"decided-by","to":"base-dec"}`)
	callJSON(t, ctx, "add", `{"kind":"question","handle":"base-q","fields":["statement=Does base work need a flag?"]}`)
	callJSON(t, ctx, "add", `{"kind":"source","handle":"base-src","fields":["title=Base planning talk","summary=Scrubbed talk","raw_sha256=`+strings.Repeat("b", 64)+`"]}`)
	callJSON(t, ctx, "propose", `{}`)
	callJSON(t, ctx, "accept", `{}`)
	callJSON(t, ctx, "land", `{}`)

	cs := newChange(t, ctx, "Retry deliveries with backoff", "intention")
	callJSON(t, ctx, "add", `{"kind":"capability","handle":"retry-a","fields":["statement=A failed delivery is retried","acceptance=Run the backoff test"]}`)
	callJSON(t, ctx, "add", `{"kind":"capability","handle":"retry-b","fields":["statement=A 429 waits for Retry-After","acceptance=Run the header test"]}`)
	callJSON(t, ctx, "add", `{"kind":"evidence","handle":"retry-a-ok","fields":["statement=Ran it","result=PASS","status=verified"]}`)
	callJSON(t, ctx, "add", `{"kind":"evidence","handle":"retry-pr","fields":["statement=PR passed CI","result=ok","pr=https://github.com/x/y/pull/525","status=verified"]}`)
	callJSON(t, ctx, "add", `{"kind":"decision","handle":"retry-dec","fields":["statement=Retries back off exponentially","question=How should retries wait?","wording=Exponential"]}`)
	callJSON(t, ctx, "link", `{"from":"retry-a","type":"decided-by","to":"retry-dec"}`)
	callJSON(t, ctx, "add", `{"kind":"source","handle":"retry-src","fields":["title=Retry planning talk","summary=Scrubbed talk","raw_sha256=`+strings.Repeat("a", 64)+`"]}`)
	callJSON(t, ctx, "link", `{"from":"retry-a-ok","type":"verifies","to":"retry-a"}`)
	callJSON(t, ctx, "add", `{"kind":"question","handle":"gone-q","fields":["statement=Should a 410 stop retries at once?"]}`)
	callJSON(t, ctx, "link", `{"from":"gone-q","type":"blocks","to":"retry-b"}`)
	callJSON(t, ctx, "propose", `{}`)

	nodeIDs := func(view string, handles ...string) []string {
		var ids []string
		for _, h := range handles {
			b, _ := json.Marshal(callJSON(t, ctx, "show", `{"ref":"`+h+`","view":"`+view+`"}`))
			var sv struct {
				Node struct {
					ID string `json:"id"`
				} `json:"node"`
			}
			json.Unmarshal(b, &sv)
			if sv.Node.ID == "" {
				t.Fatalf("no ID for %s in %s", h, view)
			}
			ids = append(ids, sv.Node.ID)
		}
		return ids
	}
	whyIDs := nodeIDs("what-if:"+cs, "retry-dec", "gone-q", "retry-src")
	baseWhyIDs := nodeIDs("as-is", "base-dec", "base-q", "base-src")

	// The dependency comes from another change set, as a landed plan's does.
	newChange(t, ctx, "Plan the order", "intention")
	callJSON(t, ctx, "link", `{"from":"`+cs+`","type":"depends-on","to":"`+base+`"}`)
	callJSON(t, ctx, "propose", `{}`)
	callJSON(t, ctx, "accept", `{}`)
	callJSON(t, ctx, "land", `{}`)
	// The page reads the change set whatever is current: here, nothing.
	callJSON(t, ctx, "change.use", `{"none":true}`)

	// A run on it.
	var common string
	if out, err := exec.Command("git", "-C", ctx.Root, "rev-parse", "--git-common-dir").Output(); err == nil {
		common = strings.TrimSpace(string(out))
		if !filepath.IsAbs(common) {
			common = filepath.Join(ctx.Root, common)
		}
	}
	rs := run.NewStore(run.Dir(common))
	at := time.Date(2026, 10, 24, 11, 0, 0, 0, time.UTC)
	if err := rs.Create("runone", run.Entry{Type: run.TypeStarted, At: at, Step: "build",
		Started: &run.Started{Workflow: "build", Change: cs}}); err != nil {
		t.Fatal(err)
	}
	// A finished run on the landed change set, four and a half minutes long.
	if err := rs.Create("runbase", run.Entry{Type: run.TypeStarted, At: at.Add(-time.Hour), Step: "build",
		Started: &run.Started{Workflow: "ship", Change: base}}); err != nil {
		t.Fatal(err)
	}
	for _, e := range []run.Entry{
		{Type: run.TypeEntered, At: at.Add(-time.Hour), Step: "build", Visit: 1},
		{Type: run.TypeEnded, At: at.Add(-time.Hour + 4*time.Minute + 30*time.Second), Step: "build", Visit: 1, Outcome: "pass"},
		{Type: run.TypeFinished, At: at.Add(-time.Hour + 4*time.Minute + 30*time.Second), Step: "build", End: run.Succeeded},
	} {
		if err := rs.Append("runbase", e); err != nil {
			t.Fatal(err)
		}
	}

	h := web.Handler(api.ReadOnlyHandler(ctx))
	// Under its project, the decision, question and source link to their model pages.
	root := web.RootWithProjects(web.Handler(nil), func(string) http.Handler { return h },
		func() []web.ProjectView { return []web.ProjectView{{Name: "proj", Initialised: true}} })
	scoped := serve(root, "GET", "/p/proj/changes/"+cs)
	if scoped.Code != 200 {
		t.Fatalf("scoped: %d %s", scoped.Code, scoped.Body)
	}
	for _, id := range whyIDs {
		if !strings.Contains(scoped.Body.String(), `href="/p/proj/nodes/`+id+`"`) {
			t.Errorf("no link to /p/proj/nodes/%s", id)
		}
	}
	rec := serve(h, "GET", "/changes/"+cs)
	if rec.Code != 200 {
		t.Fatalf("/changes/%s: %d %s", cs, rec.Code, rec.Body)
	}
	body := rec.Body.String()
	for _, s := range []string{
		"Retry deliveries with backoff", "intention", "proposed",
		"A failed delivery is retried", "Run the backoff test", "A 429 waits for Retry-After", "Run the header test",
		"Verified, has evidence", "no evidence yet", "Retries back off exponentially", "Retry planning talk", "PR #525", "2 leaves, 1 with evidence",
		"Should a 410 stop retries at once?", "Base work", `href="/changes/` + base + `"`, "landed",
		`href="/runs/runone"`, "build run", "2026-10-24 11:00", "Proposed",
	} {
		if !strings.Contains(body, s) {
			t.Errorf("missing %q in:\n%s", s, body)
		}
	}
	for _, h := range []string{"retry-a", "gone-q", "base-one"} {
		if strings.Contains(body, h) {
			t.Errorf("the page shows handle %s", h)
		}
	}
	if strings.Contains(body, ">"+cs) {
		t.Error("the page shows the change set's ID")
	}
	// A landed change set's page reads the model as it is.
	rec = serve(h, "GET", "/changes/"+base)
	if rec.Code != 200 {
		t.Fatalf("/changes/%s: %d %s", base, rec.Code, rec.Body)
	}
	body = rec.Body.String()
	for _, s := range []string{
		"Base work", "landed", "3 leaves, 3 with evidence",
		"Base capability one", "Base check one", "Base capability two", "Base check two", "Base capability three", "Base check three",
		"<dt>Shipped</dt><dd>2026-10-24 12:00</dd>",
		`<dt>Pull request</dt><dd><a href="https://github.com/x/y/pull/498">PR #498</a></dd>`,
		`<dt>Run</dt><dd><a href="/runs/runbase">ship, 4m 30s</a></dd>`,
		`<dt>Serves</dt><dd><a href="/roadmap">Deliveries are reliable</a></dd>`,
		"Base work ships first", "Does base work need a flag?", "Base planning talk",
		"<span>Proposed</span>", "<span>Accepted</span>", "<span>Landed</span>",
	} {
		if !strings.Contains(body, s) {
			t.Errorf("landed page: missing %q in:\n%s", s, body)
		}
	}
	if n := strings.Count(body, "Verified, has evidence"); n != 3 {
		t.Errorf("landed page: %d verified leaves, want 3", n)
	}
	if strings.Contains(body, base) {
		t.Error("the landed page shows the change set's ID")
	}
	for _, h := range []string{"base-one", "base-dec", "base-q", "base-src"} {
		if strings.Contains(body, h) {
			t.Errorf("the landed page shows handle %s", h)
		}
	}
	scoped = serve(root, "GET", "/p/proj/changes/"+base)
	if scoped.Code != 200 {
		t.Fatalf("scoped landed: %d %s", scoped.Code, scoped.Body)
	}
	for _, id := range baseWhyIDs {
		if !strings.Contains(scoped.Body.String(), `href="/p/proj/nodes/`+id+`"`) {
			t.Errorf("landed page: no link to /p/proj/nodes/%s", id)
		}
	}
	// A draft with no leaves says it has none yet.
	draft := newChange(t, ctx, "Empty draft", "intention")
	if body := serve(h, "GET", "/changes/"+draft).Body.String(); !strings.Contains(body, "has none yet") {
		t.Errorf("draft page:\n%s", body)
	}
	for _, target := range []string{"/changes/01ZZZZZZZZZZZZZZZZZZZZZZZZ", "/changes/nope"} {
		if rec := serve(h, "GET", target); rec.Code != 404 {
			t.Errorf("%s: %d", target, rec.Code)
		}
	}
	if rec := serve(h, "POST", "/changes/"+cs); rec.Code != 405 {
		t.Errorf("POST: %d", rec.Code)
	}
}
