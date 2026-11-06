package web

import (
	"net/http/httptest"
	"strings"
	"testing"
)

// Each Overview section is a live section: it declares the events it
// re-renders on, and its fragment URL serves that section alone.
func TestOverviewSectionsAreLive(t *testing.T) {
	sections := []struct {
		name   string
		events []string // the events the section must declare
		marker string   // what only this section's fragment renders
	}{
		{"runs", []string{EventRun}, "No live runs"},
		{"inbox", []string{EventInbox}, "Nothing waits on you"},
		{"milestones", []string{EventChange}, "No open milestones"},
		{"daemon", []string{EventDaemon}, "daemon"},
	}
	h := HandlerWithRuns(func() DaemonView { return DaemonView{} }, func() []RunView { return nil }, nil)
	page := httptest.NewRecorder()
	h.ServeHTTP(page, httptest.NewRequest("GET", "/", nil))
	if page.Code != 200 {
		t.Fatalf("/: %d", page.Code)
	}
	for _, s := range sections {
		url := FragmentURL(s.name)
		i := strings.Index(page.Body.String(), `data-rz-fragment="`+url+`"`)
		if i < 0 {
			t.Errorf("/ has no live section for %s", url)
			continue
		}
		tag := page.Body.String()[i:]
		tag = tag[:strings.Index(tag, ">")]
		for _, ev := range s.events {
			if !strings.Contains(tag, ev) {
				t.Errorf("%s declares %q, lacks event %q", s.name, tag, ev)
			}
		}

		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", url, nil))
		body := rec.Body.String()
		if rec.Code != 200 || !strings.Contains(body, s.marker) {
			t.Errorf("%s: %d, lacks %q:\n%s", url, rec.Code, s.marker, body)
		}
		if strings.Contains(body, "<html") || strings.Contains(body, "rz-section") {
			t.Errorf("%s is not the section alone:\n%s", url, body)
		}
		for _, o := range sections {
			if o.name != s.name && o.name != "daemon" && strings.Contains(body, o.marker) {
				t.Errorf("%s also renders the %s section", url, o.name)
			}
		}
	}
}
