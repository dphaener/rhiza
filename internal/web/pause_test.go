package web

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type pauseCalls struct{ calls []string }

func pausePage(token string) (http.Handler, *pauseCalls) {
	calls := &pauseCalls{}
	daemon := func() DaemonView { return DaemonView{Running: true, Limit: 2, ProjectPaused: true} }
	a := Actions{
		Pause: func(scope string, paused bool) (string, error) {
			calls.calls = append(calls.calls, scope+"|"+map[bool]string{true: "pause", false: "resume"}[paused])
			return "done", nil
		},
		Answer:  func(string, string, string) (string, error) { return "x", nil },
		Approve: func(string) (string, error) { return "x", nil },
	}
	return WithToken(func() string { return token }, WithActions(a, HandlerWith(daemon, nil))), calls
}

func TestWebPauseFormsCarryCSRFAndAreAbsentWithoutToken(t *testing.T) {
	h, _ := pausePage(testToken)
	cookie := map[string]string{"Cookie": TokenCookie + "=" + testToken}
	body := do(h, "GET", "/", cookie).Body.String()
	for _, want := range []string{
		`action="` + PausePath + `"`, `name="csrf" value="` + csrfFor(testToken) + `"`,
		`name="scope" value="daemon"`, `name="scope" value="project"`,
		`>Pause daemon</button>`, `>Resume project</button>`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("page lacks %q", want)
		}
	}
	if strings.Contains(body, testToken) {
		t.Error("page shows the token")
	}
	if b := do(h, "GET", "/", nil).Body.String(); strings.Contains(b, "<form") {
		t.Error("page carries a form without the token")
	}
	none, _ := pausePage("")
	if b := do(none, "GET", "/", nil).Body.String(); strings.Contains(b, "<form") {
		t.Error("page carries a form with no token set")
	}
}

func TestWebPausePostsAndRefusesMissingCSRF(t *testing.T) {
	h, calls := pausePage(testToken)
	cookie := map[string]string{"Cookie": TokenCookie + "=" + testToken}
	csrf := csrfFor(testToken)
	for _, f := range []url.Values{
		{CSRFField: {csrf}, "scope": {"daemon"}, "paused": {"true"}},
		{CSRFField: {csrf}, "scope": {"daemon"}, "paused": {"false"}},
		{CSRFField: {csrf}, "scope": {"project"}, "paused": {"true"}},
		{CSRFField: {csrf}, "scope": {"project"}, "paused": {"false"}},
	} {
		if rec := postTo(h, PausePath, cookie, f); rec.Code != 200 || !strings.Contains(rec.Body.String(), "done") {
			t.Fatalf("pause %v = %d %s", f, rec.Code, rec.Body)
		}
	}
	if got := strings.Join(calls.calls, ","); got != "daemon|pause,daemon|resume,project|pause,project|resume" {
		t.Errorf("calls = %s", got)
	}
	n := len(calls.calls)
	for name, f := range map[string]url.Values{
		"missing": {"scope": {"daemon"}, "paused": {"true"}},
		"wrong":   {CSRFField: {"nope"}, "scope": {"daemon"}, "paused": {"true"}},
	} {
		if rec := postTo(h, PausePath, cookie, f); rec.Code != http.StatusForbidden {
			t.Errorf("%s CSRF = %d, want 403", name, rec.Code)
		}
	}
	if rec := postTo(h, PausePath, cookie, url.Values{CSRFField: {csrf}, "scope": {"machine"}, "paused": {"true"}}); rec.Code != http.StatusBadRequest {
		t.Errorf("bad scope = %d", rec.Code)
	}
	if len(calls.calls) != n {
		t.Errorf("a refused POST ran: %v", calls.calls[n:])
	}
}

func TestWebActionsOnly(t *testing.T) {
	h, _ := pausePage(testToken)
	cookie := map[string]string{"Cookie": TokenCookie + "=" + testToken}
	form := url.Values{CSRFField: {csrfFor(testToken)}, "run": {"R1"}}
	for _, path := range []string{"/", "/retry", "/abandon", "/merge", "/actions/retry", "/actions/abandon", "/actions/merge", "/runs/R1/retry", "/runs/R1/abandon", "/changes/C1/merge", "/api/daemon.set", "/api/run.retry", "/actions/pause/x"} {
		if rec := postTo(h, path, cookie, form); rec.Code != http.StatusNotFound && rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("POST %s = %d, want 404 or 405", path, rec.Code)
		}
	}
}
