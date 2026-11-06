package web

import (
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

var fonts = []struct{ file, licence string }{
	{"Petrona-Variable.woff2", "OFL-Petrona.txt"},
	{"InstrumentSans-Variable.woff2", "OFL-InstrumentSans.txt"},
	{"JetBrainsMono-Variable.woff2", "OFL-JetBrainsMono.txt"},
}

func TestFontsServedFromHandler(t *testing.T) {
	h := Handler(nil)
	for _, f := range fonts {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", "/static/fonts/"+f.file, nil))
		if rec.Code != 200 {
			t.Errorf("%s: status %d", f.file, rec.Code)
			continue
		}
		if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "font/woff2") {
			t.Errorf("%s: content type %q", f.file, ct)
		}
		if !strings.HasPrefix(rec.Body.String(), "wOF2") {
			t.Errorf("%s: not a woff2 file", f.file)
		}
		lic, err := assets.ReadFile("static/fonts/" + f.licence)
		if err != nil || !strings.Contains(string(lic), "SIL OPEN FONT LICENSE") {
			t.Errorf("%s: licence not embedded: %v", f.licence, err)
		}
	}
}

func TestPagesLoadNothingFromTheNetwork(t *testing.T) {
	rec := httptest.NewRecorder()
	Handler(nil).ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	for _, u := range regexp.MustCompile(`(?:src|href)="(https?:[^"]*)"`).FindAllStringSubmatch(rec.Body.String(), -1) {
		t.Errorf("layout references external %s", u[1])
	}
	for _, css := range []string{"rhiza-tokens.css", "rhiza.css"} {
		b, err := assets.ReadFile("static/css/" + css)
		if err != nil {
			t.Fatal(err)
		}
		if m := regexp.MustCompile(`(?:url\(|@import)\s*["']?(https?:)?//`).FindString(string(b)); m != "" {
			t.Errorf("%s loads from the network: %q", css, m)
		}
	}
}
