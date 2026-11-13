package web_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/web"
)

func TestNodePages(t *testing.T) {
	ctx := fixtureCtx(t)
	newChange(t, ctx, "Plan the retry work", "intention")
	callJSON(t, ctx, "add", `{"kind":"source","handle":"retry-talk","fields":["title=Retry planning talk","summary=The planning conversation, scrubbed","raw_sha256=`+strings.Repeat("a", 64)+`"]}`)
	callJSON(t, ctx, "add", `{"kind":"intent","handle":"reliable-intent","fields":["statement=Deliveries are reliable"]}`)
	callJSON(t, ctx, "add", `{"kind":"capability","handle":"retry-cap","fields":["statement=A failed delivery is retried","acceptance=Run the backoff test"]}`)
	callJSON(t, ctx, "link", `{"from":"retry-cap","type":"serves","to":"reliable-intent"}`)
	callJSON(t, ctx, "add", `{"kind":"question","handle":"retry-q","fields":["statement=Should a 410 stop retries at once?"]}`)
	callJSON(t, ctx, "add", `{"kind":"decision","handle":"quoted-dec","fields":["statement=A 410 stops retries","quote=Yes, stop at once, no more waiting.","quote_where=2026-11-13, to the orchestrator","wording=the assistant's summary"]}`)
	callJSON(t, ctx, "link", `{"from":"quoted-dec","type":"answers","to":"retry-q"}`)
	callJSON(t, ctx, "link", `{"from":"quoted-dec","type":"derived-from","to":"retry-talk"}`)
	callJSON(t, ctx, "link", `{"from":"retry-cap","type":"decided-by","to":"quoted-dec"}`)
	callJSON(t, ctx, "add", `{"kind":"decision","handle":"worded-dec","fields":["statement=Retries back off exponentially","question=How should retries wait?","wording=Exponential backoff (Recommended), picked by the owner"]}`)
	callJSON(t, ctx, "link", `{"from":"retry-cap","type":"decided-by","to":"worded-dec"}`)
	callJSON(t, ctx, "propose", `{}`)
	callJSON(t, ctx, "accept", `{}`)
	callJSON(t, ctx, "land", `{}`)

	h := web.Handler(api.ReadOnlyHandler(ctx))
	cases := []struct {
		ref   string
		wants []string
	}{
		{"retry-cap", []string{"Capability", "A failed delivery is retried", "Decided by", "A 410 stops retries", "Retries back off exponentially", "Serves", "Deliveries are reliable", "Plan the retry work", "Recorded in", `href="/nodes/`}},
		{"quoted-dec", []string{"Decision", "In your words", "Yes, stop at once, no more waiting.", "2026-11-13, to the orchestrator", "Should a 410 stop retries at once?", "Retry planning talk", "What it decided", "A failed delivery is retried"}},
		{"worded-dec", []string{"Picked by you", "Exponential backoff (Recommended), picked by the owner", "How should retries wait?"}},
		{"retry-q", []string{"Question", "Should a 410 stop retries at once?", "A 410 stops retries"}},
		{"retry-talk", []string{"Source", "Retry planning talk", "The planning conversation, scrubbed", "Plan the retry work"}},
		{"reliable-intent", []string{"Intent", "Deliveries are reliable", "A failed delivery is retried"}},
	}
	for _, c := range cases {
		rec := serve(h, "GET", "/nodes/"+c.ref)
		if rec.Code != 200 {
			t.Fatalf("/nodes/%s: %d %s", c.ref, rec.Code, rec.Body)
		}
		body := rec.Body.String()
		for _, s := range c.wants {
			if !strings.Contains(body, s) {
				t.Errorf("%s: missing %q in:\n%s", c.ref, s, body)
			}
		}
		// The page never shows its node's handle or ID.
		if strings.Contains(body, c.ref) {
			t.Errorf("%s: the page shows the handle", c.ref)
		}
		id := nodeID(t, ctx, c.ref)
		if strings.Contains(body, id) {
			t.Errorf("%s: the page shows the node's ID %s", c.ref, id)
		}
	}
	// No handle of any node shows on any page.
	for _, c := range cases {
		body := serve(h, "GET", "/nodes/quoted-dec").Body.String()
		if strings.Contains(body, c.ref) {
			t.Errorf("a page shows handle %s", c.ref)
		}
	}
	// The ID in a link opens the same page.
	if rec := serve(h, "GET", "/nodes/"+nodeID(t, ctx, "retry-cap")); rec.Code != 200 {
		t.Errorf("by ID: %d", rec.Code)
	}

	rec := serve(h, "GET", "/nodes/no-such-node")
	if rec.Code != 404 || !strings.Contains(rec.Body.String(), "no-such-node") {
		t.Errorf("unknown ref: %d %s", rec.Code, rec.Body)
	}
	if rec := serve(h, "POST", "/nodes/retry-cap"); rec.Code != 405 {
		t.Errorf("POST: %d", rec.Code)
	}
}

// nodeID is a node's ID, read through the show operation.
func nodeID(t *testing.T, ctx *api.Ctx, ref string) string {
	t.Helper()
	b, err := json.Marshal(callJSON(t, ctx, "show", `{"ref":"`+ref+`","view":"to-be"}`))
	if err != nil {
		t.Fatal(err)
	}
	var v struct {
		Node struct {
			ID string `json:"id"`
		} `json:"node"`
	}
	json.Unmarshal(b, &v)
	return v.Node.ID
}
