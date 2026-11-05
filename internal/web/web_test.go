package web

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParsePort(t *testing.T) {
	for v, want := range map[string]int{"7787": 7787, "8080": 8080, "1": 1, "65535": 65535, "off": 0} {
		if n, err := ParsePort(v); err != nil || n != want {
			t.Errorf("ParsePort(%q) = %d, %v", v, n, err)
		}
	}
	for _, v := range []string{"0", "70000", "abc", "", "-1", "080", "Off"} {
		if _, err := ParsePort(v); err == nil {
			t.Errorf("ParsePort(%q) accepted", v)
		}
	}
}

func TestAddrIsAlwaysLoopback(t *testing.T) {
	if a, _ := Addr(DefaultPort); a != "127.0.0.1:7787" {
		t.Errorf("default addr %q", a)
	}
	if a, _ := Addr("8080"); a != "127.0.0.1:8080" {
		t.Errorf("changed addr %q", a)
	}
	if a, err := Addr("off"); a != "" || err != nil {
		t.Errorf("off addr %q, %v", a, err)
	}
}

func TestHandlerServesReadOnly(t *testing.T) {
	h := Handler(nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/site.webmanifest", nil))
	if rec.Code != 200 {
		t.Errorf("manifest: %d", rec.Code)
	}
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/static/css/rhiza.css", nil))
	if rec.Code != 200 {
		t.Errorf("css: %d", rec.Code)
	}
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("POST", "/static/css/rhiza.css", nil))
	if rec.Code != 405 {
		t.Errorf("POST: %d", rec.Code)
	}
}

type fakeListener struct{ closed bool }

func (f *fakeListener) Close() error { f.closed = true; return nil }

func TestManagerFollowsSetting(t *testing.T) {
	var started []string
	var ls []*fakeListener
	fail := false
	m := &Manager{Factory: func(addr string, h http.Handler) (Listener, error) {
		if fail {
			return nil, errors.New("in use")
		}
		started = append(started, addr)
		l := &fakeListener{}
		ls = append(ls, l)
		return l, nil
	}}
	step := func(v string, wantAddr string, wantStarts int) {
		t.Helper()
		if err := m.Apply(v); err != nil {
			t.Fatalf("Apply(%q): %v", v, err)
		}
		if m.Addr() != wantAddr || len(started) != wantStarts {
			t.Fatalf("Apply(%q): addr %q, starts %v", v, m.Addr(), started)
		}
	}
	step("7787", "127.0.0.1:7787", 1)
	step("7787", "127.0.0.1:7787", 1) // unchanged: no restart
	step("8080", "127.0.0.1:8080", 2) // moved
	if !ls[0].closed || ls[1].closed {
		t.Errorf("move closed %v %v", ls[0].closed, ls[1].closed)
	}
	step("off", "", 2) // stopped
	if !ls[1].closed {
		t.Error("off left the listener open")
	}
	step("off", "", 2)
	step("7787", "127.0.0.1:7787", 3) // restarted

	if err := m.Apply("0"); err == nil {
		t.Error("bad value accepted")
	}
	m.Close()
	fail = true
	if err := m.Apply("7787"); err == nil || m.Addr() != "" {
		t.Errorf("failed start: %v, addr %q", err, m.Addr())
	}
	fail = false
	step("7787", "127.0.0.1:7787", 4) // retried at the next poll
}

func TestHandlerRendersPageAndScript(t *testing.T) {
	h := Handler(nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != 200 || !strings.HasPrefix(rec.Header().Get("Content-Type"), "text/html") {
		t.Fatalf("/: %d %q", rec.Code, rec.Header().Get("Content-Type"))
	}
	for _, want := range []string{"<title>Overview · Rhiza</title>", `href="/static/css/rhiza.css"`, `src="/static/app.js"`} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Errorf("/ lacks %q", want)
		}
	}
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/static/app.js", nil))
	if rec.Code != 200 || !strings.Contains(rec.Header().Get("Content-Type"), "javascript") || rec.Body.Len() == 0 {
		t.Errorf("app.js: %d %q", rec.Code, rec.Header().Get("Content-Type"))
	}
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/nope", nil))
	if rec.Code != 404 {
		t.Errorf("/nope: %d", rec.Code)
	}
}

func TestLiveSectionAndFragment(t *testing.T) {
	RegisterFragment("sample", &Fragment{
		Events:   []string{EventRun, EventLog},
		Template: template.Must(template.New("sample").Parse(`<p>runs: {{.}}</p>`)),
		Data:     func(*http.Request) any { return 3 },
	})
	defer delete(fragments, "sample")
	got, err := section("sample", 3)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`data-rz-fragment="/fragments/sample"`, `data-rz-events="run log"`, `<p>runs: 3</p>`} {
		if !strings.Contains(string(got), want) {
			t.Errorf("section lacks %q: %s", want, got)
		}
	}
	h := Handler(nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", FragmentURL("sample"), nil))
	if rec.Code != 200 || rec.Body.String() != `<p>runs: 3</p>` {
		t.Errorf("fragment: %d %q", rec.Code, rec.Body.String())
	}
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", FragmentURL("nope"), nil))
	if rec.Code != 404 {
		t.Errorf("unknown fragment: %d", rec.Code)
	}
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("POST", FragmentURL("sample"), nil))
	if rec.Code != 405 {
		t.Errorf("POST fragment: %d", rec.Code)
	}
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/static/app.js", nil))
	for _, want := range []string{"data-rz-fragment", "/events", "setStale(true)"} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Errorf("app.js lacks %q", want)
		}
	}
}
