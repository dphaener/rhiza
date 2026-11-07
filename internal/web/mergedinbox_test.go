package web_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/web"
)

func TestMergedInboxOrdersTagsAndLinks(t *testing.T) {
	src := func() []web.ProjectInbox {
		return []web.ProjectInbox{
			{Name: "alpha", Items: []web.InboxItemView{
				{Kind: "question", Title: "Alpha late question", Change: "C1", At: "2026-11-03T10:00:00Z"},
				{Kind: "failed", Title: "Alpha early failure", Run: "R1", Change: "C1", At: "2026-11-01T10:00:00Z"},
			}},
			{Name: "beta", Items: []web.InboxItemView{
				{Kind: "signoff", Title: "Beta middle sign-off", Change: "C2", At: "2026-11-02T10:00:00Z"},
			}},
		}
	}
	h := web.RootWithInbox(web.Handler(nil), func(string) http.Handler { return nil }, nil, src)
	rec := serve(h, "GET", "/inbox")
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	body := rec.Body.String()
	order := []string{"Alpha early failure", "Beta middle sign-off", "Alpha late question"}
	last := -1
	for _, w := range order {
		i := strings.Index(body, w)
		if i < 0 || i < last {
			t.Fatalf("%q missing or out of order:\n%s", w, body)
		}
		last = i
	}
	for _, want := range []string{
		`href="/p/alpha/runs/R1"`, `href="/p/beta/changes/C2"`, `href="/p/alpha/changes/C1"`,
		`<span class="rz-project-tag">alpha</span>`, `<span class="rz-project-tag">beta</span>`, "3 waiting on you",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("page lacks %q", want)
		}
	}
}

func TestMergedInboxEmpty(t *testing.T) {
	rec := serve(web.RootWithInbox(web.Handler(nil), nil, nil, nil), "GET", "/inbox")
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Nothing waits on you") {
		t.Errorf("%d %s", rec.Code, rec.Body)
	}
}
