package web

import (
	"net/http/httptest"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// liveSection matches one live section's opening tag on the Overview page:
// its fragment URL and its space-separated event names.
var liveSection = regexp.MustCompile(`<div class="rz-section" data-rz-fragment="([^"]*)" data-rz-events="([^"]*)">`)

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
		{"daemon", []string{EventDaemon}, "Paused, 1 of 3 slots in use"},
	}
	daemon := func() DaemonView { return DaemonView{Running: true, Paused: true, Slots: 1, Limit: 3} }
	h := HandlerWithRuns(daemon, func() []RunView { return nil }, nil)
	page := httptest.NewRecorder()
	h.ServeHTTP(page, httptest.NewRequest("GET", "/", nil))
	if page.Code != 200 {
		t.Fatalf("/: %d", page.Code)
	}
	declared := map[string][]string{} // fragment URL -> its event tokens
	for _, m := range liveSection.FindAllStringSubmatch(page.Body.String(), -1) {
		declared[m[1]] = strings.Fields(m[2])
	}
	for _, s := range sections {
		url := FragmentURL(s.name)
		events, ok := declared[url]
		if !ok {
			t.Errorf("/ has no live section for %s", url)
		}
		for _, ev := range s.events {
			if !slices.Contains(events, ev) {
				t.Errorf("%s declares events %q, lacks %q", s.name, events, ev)
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
			if o.name != s.name && strings.Contains(body, o.marker) {
				t.Errorf("%s also renders the %s section", url, o.name)
			}
		}
	}
}
