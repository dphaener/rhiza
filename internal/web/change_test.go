package web_test

import (
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
	// A dependency, landed.
	base := newChange(t, ctx, "Base work", "intention")
	callJSON(t, ctx, "add", `{"kind":"capability","handle":"base-cap","fields":["statement=Base capability","acceptance=Base check"]}`)
	callJSON(t, ctx, "propose", `{}`)
	callJSON(t, ctx, "accept", `{}`)
	callJSON(t, ctx, "land", `{}`)

	cs := newChange(t, ctx, "Retry deliveries with backoff", "intention")
	callJSON(t, ctx, "add", `{"kind":"capability","handle":"retry-a","fields":["statement=A failed delivery is retried","acceptance=Run the backoff test"]}`)
	callJSON(t, ctx, "add", `{"kind":"capability","handle":"retry-b","fields":["statement=A 429 waits for Retry-After","acceptance=Run the header test"]}`)
	callJSON(t, ctx, "add", `{"kind":"evidence","handle":"retry-a-ok","fields":["statement=Ran it","result=PASS"]}`)
	callJSON(t, ctx, "link", `{"from":"retry-a-ok","type":"verifies","to":"retry-a"}`)
	callJSON(t, ctx, "add", `{"kind":"question","handle":"gone-q","fields":["statement=Should a 410 stop retries at once?"]}`)
	callJSON(t, ctx, "link", `{"from":"gone-q","type":"blocks","to":"retry-b"}`)
	callJSON(t, ctx, "propose", `{}`)

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

	h := web.Handler(api.ReadOnlyHandler(ctx))
	rec := serve(h, "GET", "/changes/"+cs)
	if rec.Code != 200 {
		t.Fatalf("/changes/%s: %d %s", cs, rec.Code, rec.Body)
	}
	body := rec.Body.String()
	for _, s := range []string{
		"Retry deliveries with backoff", "intention", "proposed",
		"A failed delivery is retried", "Run the backoff test", "A 429 waits for Retry-After", "Run the header test",
		"has evidence", "no evidence yet", "2 leaves, 1 with evidence",
		"Should a 410 stop retries at once?", "Base work", `href="/changes/` + base + `"`, "landed",
		`href="/runs/runone"`, "build run", "2026-10-24 11:00", "Proposed",
	} {
		if !strings.Contains(body, s) {
			t.Errorf("missing %q in:\n%s", s, body)
		}
	}
	for _, h := range []string{"retry-a", "gone-q", "base-cap"} {
		if strings.Contains(body, h) {
			t.Errorf("the page shows handle %s", h)
		}
	}
	if strings.Contains(body, ">"+cs) {
		t.Error("the page shows the change set's ID")
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
