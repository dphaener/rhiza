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

func TestAddr(t *testing.T) {
	for _, c := range []struct{ host, port, want string }{
		{DefaultHost, DefaultPort, "127.0.0.1:7787"},
		{"", DefaultPort, "127.0.0.1:7787"},
		{"0.0.0.0", DefaultPort, "0.0.0.0:7787"},
		{"::1", DefaultPort, "[::1]:7787"},
		{"::", "8080", "[::]:8080"},
		{DefaultHost, "off", ""},
	} {
		if a, err := Addr(c.host, c.port); a != c.want || err != nil {
			t.Errorf("Addr(%q, %q) = %q, %v; want %q", c.host, c.port, a, err, c.want)
		}
	}
	for _, h := range []string{"localhost", "abc", "127.0.0.1:80", "[::1]", " x"} {
		if a, err := Addr(h, DefaultPort); a != "" || err == nil {
			t.Errorf("Addr(%q) = %q, %v; want an error", h, a, err)
		}
		if _, err := ParseHost(h); err == nil {
			t.Errorf("ParseHost(%q) accepted", h)
		}
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
		if err := m.Apply(DefaultHost, v); err != nil {
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

	if err := m.Apply(DefaultHost, "0"); err == nil {
		t.Error("bad value accepted")
	}
	m.Close()
	fail = true
	if err := m.Apply(DefaultHost, "7787"); err == nil || m.Addr() != "" {
		t.Errorf("failed start: %v, addr %q", err, m.Addr())
	}
	fail = false
	step("7787", "127.0.0.1:7787", 4) // retried at the next poll
}

func TestManagerFollowsHost(t *testing.T) {
	var started []string
	var ls []*fakeListener
	m := &Manager{Factory: func(addr string, h http.Handler) (Listener, error) {
		started = append(started, addr)
		l := &fakeListener{}
		ls = append(ls, l)
		return l, nil
	}}
	for _, host := range []string{DefaultHost, "0.0.0.0"} {
		if err := m.Apply(host, "7787"); err != nil {
			t.Fatal(err)
		}
	}
	if m.Addr() != "0.0.0.0:7787" || len(started) != 2 || !ls[0].closed || ls[1].closed {
		t.Fatalf("addr %q, starts %v", m.Addr(), started)
	}
	if err := m.Apply("localhost", "7787"); err == nil {
		t.Error("bad host accepted")
	}
	if len(started) != 2 {
		t.Errorf("bad host bound a socket: %v", started)
	}
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
