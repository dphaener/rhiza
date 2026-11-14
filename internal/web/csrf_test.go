package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func csrfGate() (http.Handler, *int) {
	hits := 0
	return WithToken(func() string { return testToken }, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			hits++
		}
		w.Write([]byte(CSRF(r)))
	})), &hits
}

func post(h http.Handler, hdr map[string]string, form url.Values) *httptest.ResponseRecorder {
	r := httptest.NewRequest("POST", "/act", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range hdr {
		r.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	return rec
}

func TestWebCSRF(t *testing.T) {
	h, hits := csrfGate()
	cookie := map[string]string{"Cookie": TokenCookie + "=" + testToken}

	// The page's value: handed only to a request that carries the token.
	if v := strings.TrimSpace(do(h, "GET", "/", nil).Body.String()); v != "" {
		t.Errorf("CSRF without the token = %q, want empty", v)
	}
	csrf := do(h, "GET", "/", cookie).Body.String()
	if csrf == "" || strings.Contains(csrf, testToken) {
		t.Fatalf("CSRF = %q, want a value that is not the token", csrf)
	}

	if rec := post(h, cookie, url.Values{CSRFField: {csrf}}); rec.Code != 200 || *hits != 1 {
		t.Fatalf("cookie POST with the page's CSRF = %d, hits %d", rec.Code, *hits)
	}
	for name, form := range map[string]url.Values{
		"missing": {"x": {"1"}},
		"wrong":   {CSRFField: {"nope"}},
		"empty":   {CSRFField: {""}},
	} {
		if rec := post(h, cookie, form); rec.Code != http.StatusForbidden {
			t.Errorf("%s: cookie POST = %d, want 403", name, rec.Code)
		}
	}
	if *hits != 1 {
		t.Errorf("handler reached %d times, want 1", *hits)
	}

	// A bearer header is not sent by a browser unasked, so it needs no field.
	if rec := post(h, map[string]string{"Authorization": "Bearer " + testToken}, nil); rec.Code != 200 {
		t.Errorf("bearer POST = %d", rec.Code)
	}
	// A wrong token with the right-looking field is still refused.
	if rec := post(h, map[string]string{"Cookie": TokenCookie + "=nope"}, url.Values{CSRFField: {csrf}}); rec.Code != http.StatusUnauthorized {
		t.Errorf("wrong cookie = %d, want 401", rec.Code)
	}
}
