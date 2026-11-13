package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const testToken = "s3cret-token-value-1234"

func tokenHandler(token string) (http.Handler, *int) {
	hits := 0
	return WithToken(func() string { return token }, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			hits++
		}
		w.Write([]byte("ok"))
	})), &hits
}

func do(h http.Handler, method, target string, hdr map[string]string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	for k, v := range hdr {
		r.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	return rec
}

func TestWebTokenNoneSetRefusesEveryPost(t *testing.T) {
	h, hits := tokenHandler("")
	for _, hdr := range []map[string]string{nil, {"Authorization": "Bearer "}, {"Authorization": "Bearer anything"}, {"Cookie": TokenCookie + "="}} {
		if rec := do(h, "POST", "/api/x", hdr); rec.Code != http.StatusForbidden {
			t.Errorf("POST with %v and no token set = %d, want 403", hdr, rec.Code)
		}
	}
	if *hits != 0 {
		t.Errorf("handler reached %d times", *hits)
	}
	if rec := do(h, "GET", "/runs?token=x", nil); rec.Code != 200 || rec.Header().Get("Set-Cookie") != "" {
		t.Errorf("GET with no token set = %d, cookie %q", rec.Code, rec.Header().Get("Set-Cookie"))
	}
}

func TestWebTokenWrongOrMissingRefused(t *testing.T) {
	h, hits := tokenHandler(testToken)
	for name, hdr := range map[string]map[string]string{
		"missing":      nil,
		"wrong bearer": {"Authorization": "Bearer nope"},
		"wrong cookie": {"Cookie": TokenCookie + "=nope"},
		"basic":        {"Authorization": "Basic " + testToken},
	} {
		if rec := do(h, "POST", "/api/x", hdr); rec.Code != http.StatusUnauthorized {
			t.Errorf("%s: POST = %d, want 401", name, rec.Code)
		}
	}
	if *hits != 0 {
		t.Errorf("handler reached %d times", *hits)
	}
	// A wrong ?token= sets nothing; reads stay open.
	rec := do(h, "GET", "/runs?token=nope", nil)
	if rec.Code != 200 || rec.Header().Get("Set-Cookie") != "" {
		t.Errorf("GET wrong token = %d, cookie %q", rec.Code, rec.Header().Get("Set-Cookie"))
	}
}

func TestWebTokenCookieAfterQueryVisit(t *testing.T) {
	h, hits := tokenHandler(testToken)
	rec := do(h, "GET", "/runs?token="+testToken+"&x=1", nil)
	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/runs?x=1" {
		t.Fatalf("visit = %d to %q, want a redirect without the token", rec.Code, rec.Header().Get("Location"))
	}
	cs := rec.Result().Cookies()
	if len(cs) != 1 || cs[0].Name != TokenCookie || cs[0].Value != testToken || !cs[0].HttpOnly {
		t.Fatalf("cookies = %+v, want one HttpOnly %s", cs, TokenCookie)
	}
	if rec := post(h, map[string]string{"Cookie": cs[0].String()}, url.Values{CSRFField: {csrfFor(testToken)}}); rec.Code != 200 {
		t.Errorf("POST with cookie = %d", rec.Code)
	}
	if *hits != 1 {
		t.Errorf("hits = %d, want 1", *hits)
	}
}

func TestWebTokenBearerAccepted(t *testing.T) {
	h, hits := tokenHandler(testToken)
	if rec := do(h, "POST", "/api/x", map[string]string{"Authorization": "Bearer " + testToken}); rec.Code != 200 || strings.TrimSpace(rec.Body.String()) != "ok" {
		t.Fatalf("POST with bearer = %d %q", rec.Code, rec.Body)
	}
	if *hits != 1 {
		t.Errorf("hits = %d", *hits)
	}
}

func TestWebTokenActionsOn(t *testing.T) {
	var on []bool
	h := WithToken(func() string { return testToken }, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { on = append(on, ActionsOn(r)) }))
	do(h, "GET", "/", nil)
	do(h, "GET", "/", map[string]string{"Authorization": "Bearer " + testToken})
	if len(on) != 2 || on[0] || !on[1] {
		t.Errorf("ActionsOn = %v, want [false true]", on)
	}
}

func TestValidateToken(t *testing.T) {
	if err := ValidateToken(testToken); err != nil {
		t.Error(err)
	}
	for _, v := range []string{"", "short", "has space in the token value", "tab\tinside-the-token-xx"} {
		if ValidateToken(v) == nil {
			t.Errorf("ValidateToken(%q) accepted", v)
		}
	}
}
