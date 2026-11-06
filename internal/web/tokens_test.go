package web

import (
	"regexp"
	"strings"
	"testing"
)

// Colour tokens of the delivered Claude Design system (commit 1f4b22c,
// static/css/rhiza-tokens.css), pinned per theme.
var darkTokens = map[string]string{
	"ground": "#0d0c09", "panel": "#17140f", "panel-hover": "#1e1a14", "line": "#2b261d",
	"ink": "#ede6d6", "ink-2": "#b5ac97", "ink-3": "#918975",
	"gold": "#e8d49a", "teal": "#86c9b7", "rose": "#d9a48f",
	"gold-soft": "#e8d49a1a", "teal-soft": "#86c9b71a", "rose-soft": "#d9a48f1f",
	"gold-glow": "#e8d49a59", "on-solid": "#0d0c09",
	"term-bg": "#0d0c09", "term-line": "#2b261d", "term-ink": "#ede6d6", "term-ink-3": "#918975",
	"term-gold": "#e8d49a", "term-teal": "#86c9b7", "term-rose": "#d9a48f",
}

var lightTokens = map[string]string{
	"ground": "#f2eee4", "panel": "#faf8f2", "panel-hover": "#f2eee4", "line": "#dad3c2",
	"ink": "#1d1a14", "ink-2": "#4c4535", "ink-3": "#665e4c",
	"gold": "#82600f", "teal": "#2b7364", "rose": "#9a4f37",
	"gold-soft": "#82600f14", "teal-soft": "#2b736414", "rose-soft": "#9a4f3714",
	"gold-glow": "#82600f40", "on-solid": "#f2eee4",
	"term-bg": "#0d0c09", "term-line": "#2b261d", "term-ink": "#ede6d6", "term-ink-3": "#918975",
	"term-gold": "#e8d49a", "term-teal": "#86c9b7", "term-rose": "#d9a48f",
}

// block returns the declarations of the first rule whose selector starts with sel.
func block(t *testing.T, css, sel string) string {
	t.Helper()
	i := strings.Index(css, sel+" {")
	if i < 0 {
		t.Fatalf("no rule %q", sel)
	}
	rest := css[i+len(sel)+2:]
	return rest[:strings.Index(rest, "}")]
}

func decls(body string) map[string]string {
	m := map[string]string{}
	for _, d := range regexp.MustCompile(`--([\w-]+):\s*([^;]+);`).FindAllStringSubmatch(body, -1) {
		m[d[1]] = strings.TrimSpace(d[2])
	}
	return m
}

func TestTokensEqualDeliveredDesignSystem(t *testing.T) {
	b, err := assets.ReadFile("static/css/rhiza-tokens.css")
	if err != nil {
		t.Fatal(err)
	}
	css := string(b)

	check := func(name string, got, want map[string]string) {
		for k, w := range want {
			if got[k] != w {
				t.Errorf("%s --%s = %q, want %q", name, k, got[k], w)
			}
		}
	}
	check("dark default", decls(block(t, css, `:root, [data-theme="dark"]`)), darkTokens)
	media := css[strings.Index(css, "@media (prefers-color-scheme: light)"):]
	check("light under prefers-color-scheme", decls(block(t, media, `:root:not([data-theme="dark"])`)), lightTokens)
	check("light forced", decls(block(t, css, `[data-theme="light"]`)), lightTokens)

	base := decls(block(t, css, ":root"))
	for k, w := range map[string]string{
		"radius-md": "12px", "radius-lg": "16px", "rule": "1px",
	} {
		if base[k] != w {
			t.Errorf("--%s = %q, want %q", k, base[k], w)
		}
	}
	for k, face := range map[string]string{"font-display": "Petrona", "font-body": "Instrument Sans", "font-mono": "JetBrains Mono"} {
		if !strings.Contains(base[k], face) {
			t.Errorf("--%s = %q lacks %s", k, base[k], face)
		}
	}
}

func TestStylesheetsHaveNoGradients(t *testing.T) {
	for _, f := range []string{"static/css/rhiza-tokens.css", "static/css/rhiza.css"} {
		b, err := assets.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(strings.ToLower(string(b)), "gradient") {
			t.Errorf("%s has a gradient", f)
		}
	}
}
