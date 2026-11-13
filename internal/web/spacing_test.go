package web

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"
)

// cssRule is one declaration block, with the @media condition it sits in
// ("" at the top level).
type cssRule struct {
	media     string
	selectors []string
	decls     map[string]string
}

// parseCSS reads the flat subset rhiza.css uses: rules and one level of @media.
func parseCSS(src string) []cssRule {
	src = regexp.MustCompile(`(?s)/\*.*?\*/`).ReplaceAllString(src, "")
	var rules []cssRule
	var walk func(s, media string)
	walk = func(s, media string) {
		for {
			open := strings.Index(s, "{")
			if open < 0 {
				return
			}
			head := strings.TrimSpace(s[:open])
			depth, end := 1, open+1
			for ; end < len(s) && depth > 0; end++ {
				switch s[end] {
				case '{':
					depth++
				case '}':
					depth--
				}
			}
			body := s[open+1 : end-1]
			s = s[end:]
			if strings.HasPrefix(head, "@media") {
				walk(body, strings.TrimSpace(strings.TrimPrefix(head, "@media")))
				continue
			}
			if strings.HasPrefix(head, "@") {
				continue
			}
			r := cssRule{media: media, decls: map[string]string{}}
			for _, sel := range strings.Split(head, ",") {
				r.selectors = append(r.selectors, strings.Join(strings.Fields(sel), " "))
			}
			for _, d := range strings.Split(body, ";") {
				k, v, ok := strings.Cut(d, ":")
				if ok {
					r.decls[strings.TrimSpace(k)] = strings.TrimSpace(v)
				}
			}
			rules = append(rules, r)
		}
	}
	walk(src, "")
	return rules
}

// declared returns the last value of prop that a rule for selector sets in
// the given media condition.
func declared(rules []cssRule, media, selector, prop string) string {
	val := ""
	for _, r := range rules {
		if r.media != media {
			continue
		}
		for _, sel := range r.selectors {
			if sel == selector {
				if v, ok := r.decls[prop]; ok {
					val = v
				}
			}
		}
	}
	return val
}

func TestSpacingRules(t *testing.T) {
	b, err := fs.ReadFile(assets, "static/css/rhiza.css")
	if err != nil {
		t.Fatal(err)
	}
	rules := parseCSS(string(b))
	const desktop = "(min-width: 721px)"

	for _, c := range []struct{ name, media, selector, prop, want string }{
		{"page head kicker to title", "", ".rz-pagehead > div", "gap", "var(--space-3)"},
		{"panel head kicker to title", "", ".rz-panel-head > div", "gap", "var(--space-1)"},
		{"panel head to content", "", ".rz-panel-head", "margin-bottom", "var(--space-4)"},
		{"between panels (stack)", "", ".rz-stack", "gap", "var(--space-5)"},
		{"between panels (grid)", "", ".rz-grid", "gap", "var(--space-5)"},
		{"between page sections", "", ".rz-main", "gap", "var(--space-6)"},
		{"panel padding on a phone", "", ".rz-panel", "padding", "var(--space-4)"},
		{"panel padding on desktop", desktop, ".rz-panel", "padding", "var(--space-5)"},
	} {
		if got := declared(rules, c.media, c.selector, c.prop); got != c.want {
			t.Errorf("%s: %s %s = %q, want %q", c.name, c.selector, c.prop, got, c.want)
		}
	}

	heading := regexp.MustCompile(`(^|[\s>+~])h[1-6]$`)
	for _, r := range rules {
		for _, sel := range r.selectors {
			if !heading.MatchString(sel) {
				continue
			}
			for prop, v := range r.decls {
				if strings.HasPrefix(prop, "margin") && v != "0" {
					t.Errorf("heading rule %q sets %s: %s; spacing belongs on the parent", sel, prop, v)
				}
			}
		}
	}

	inline := regexp.MustCompile(`style="[^"]*margin`)
	err = fs.WalkDir(templateFS, "templates", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		src, err := fs.ReadFile(templateFS, p)
		if err != nil {
			return err
		}
		if inline.Match(src) {
			t.Errorf("%s has an inline margin style", p)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
