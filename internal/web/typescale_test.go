package web

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var (
	wantSizes = []string{"11px", "12px", "13px", "14px", "15px", "16px", "17px", "20px", "22px", "26px", "28px", "32px", "36px"}
	wantLines = []string{"1", "1.6", "16px", "18px", "20px", "22px", "26px", "30px", "34px", "36px", "40px"}

	sizeTok = regexp.MustCompile(`--fs-[\w-]+:\s*([^;]+);`)
	lineTok = regexp.MustCompile(`--lh-[\w-]+:\s*([^;]+);`)
	sizeVar = regexp.MustCompile(`^var\(--fs-[\w-]+\)$`)
	lineVar = regexp.MustCompile(`^var\(--lh-[\w-]+\)$`)

	// a font shorthand: [style] [weight] size[/line-height] family
	fontDecl = regexp.MustCompile(`(?:^|[;{\s"'])(font|font-size|line-height)\s*:\s*([^;}"']+)`)
	sizeAtom = regexp.MustCompile(`^\d*\.?\d+(px|rem|em|%)$`)
)

func values(re *regexp.Regexp, src string) []string {
	var out []string
	for _, m := range re.FindAllStringSubmatch(src, -1) {
		out = append(out, strings.TrimSpace(m[1]))
	}
	sort.Strings(out)
	return out
}

func sorted(s []string) []string {
	c := append([]string(nil), s...)
	sort.Strings(c)
	return c
}

func same(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// sizeAndLine splits a declaration into its size and line height ("" when it
// sets none).
func sizeAndLine(prop, val string) (size, line string) {
	val = strings.TrimSpace(val)
	switch prop {
	case "font-size":
		return val, ""
	case "line-height":
		return "", val
	}
	// shorthand: the size token is the one that carries the optional /line-height
	for _, f := range strings.Fields(val) {
		if s, l, ok := strings.Cut(f, "/"); ok && s != "" {
			return s, l
		}
	}
	for _, f := range strings.Fields(val) {
		if sizeAtom.MatchString(f) || sizeVar.MatchString(f) {
			return f, ""
		}
	}
	return "", ""
}

func TestTypeScale(t *testing.T) {
	tb, err := fs.ReadFile(assets, "static/css/rhiza-tokens.css")
	if err != nil {
		t.Fatal(err)
	}
	if got := values(sizeTok, string(tb)); !same(got, sorted(wantSizes)) {
		t.Errorf("size tokens = %v, want %v", got, sorted(wantSizes))
	}
	if got := values(lineTok, string(tb)); !same(got, sorted(wantLines)) {
		t.Errorf("line-height tokens = %v, want %v", got, sorted(wantLines))
	}
	sizes, lines := map[string]bool{}, map[string]bool{}
	for _, s := range wantSizes {
		sizes[s] = true
	}
	for _, l := range wantLines {
		lines[l] = true
	}

	// in the stylesheet and templates every size and line height is a token
	check := func(file, src string) {
		for _, m := range fontDecl.FindAllStringSubmatch(src, -1) {
			size, line := sizeAndLine(m[1], m[2])
			if size != "" && !sizeVar.MatchString(size) {
				t.Errorf("%s: %s: %s sets size %q outside the tokens", file, m[1], strings.TrimSpace(m[2]), size)
			}
			if line != "" && !lineVar.MatchString(line) {
				t.Errorf("%s: %s: %s sets line height %q outside the tokens", file, m[1], strings.TrimSpace(m[2]), line)
			}
		}
	}
	css, err := fs.ReadFile(assets, "static/css/rhiza.css")
	if err != nil {
		t.Fatal(err)
	}
	check("static/css/rhiza.css", string(css))
	err = fs.WalkDir(templateFS, "templates", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		src, err := fs.ReadFile(templateFS, p)
		if err != nil {
			return err
		}
		for _, m := range regexp.MustCompile(`style="[^"]*"`).FindAllString(string(src), -1) {
			check(p, m)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// the design pages use only sizes the tokens have
	pages, err := filepath.Glob("pages/*.html")
	if err != nil || len(pages) == 0 {
		t.Fatalf("no design pages: %v", err)
	}
	for _, p := range pages {
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range fontDecl.FindAllStringSubmatch(string(b), -1) {
			size, line := sizeAndLine(m[1], m[2])
			if size != "" && !strings.HasPrefix(size, "var(") && !sizes[size] {
				t.Errorf("%s: %s: %s sets size %s, which the tokens lack", p, m[1], strings.TrimSpace(m[2]), size)
			}
			if line != "" && !strings.HasPrefix(line, "var(") && !lines[line] {
				t.Errorf("%s: %s: %s sets line height %s, which the tokens lack", p, m[1], strings.TrimSpace(m[2]), line)
			}
		}
	}
}
