package guide

import (
	"regexp"
	"strings"
	"testing"
)

// TestBuildIsAnIndex: the build topic is a short index that names the
// three role topics, each by its skill and its rhiza guide route, and says
// which role loads which; the role guidance is not repeated in it.
func TestBuildIsAnIndex(t *testing.T) {
	build, ok := Lookup("build")
	if !ok || build.Skill {
		t.Fatalf("build: found %v, skill %v; want a topic", ok, build.Skill)
	}
	checkPhrases(t, "build", map[string][]string{
		"builder":  {"Builder: `/rhiza-builder` (`rhiza guide builder`)"},
		"retry":    {"RHIZA_ATTEMPT` > 1", "`/rhiza-build-retry` (`rhiza guide build-retry`)"},
		"verifier": {"Verifier: `/rhiza-verifier` (`rhiza guide verifier`)"},
		"done":     {"rhiza session done"},
	})
	if n := Tokens(build.Body); n > 250 {
		t.Errorf("build index is %d tokens, want a short index", n)
	}
	for _, banned := range []string{"Record agent actor", "rhiza run findings", "Stream CSV exports", "Claude"} {
		if strings.Contains(build.Body, banned) {
			t.Errorf("build index holds %q: that belongs in a role topic", banned)
		}
	}
	for _, name := range []string{"builder", "verifier", "build-retry"} {
		if _, ok := Lookup(name); !ok {
			t.Errorf("build names the %s topic, which does not exist", name)
		}
	}
}

// finalParagraph is the last blank-line-separated paragraph of s.
func finalParagraph(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.LastIndex(s, "\n\n"); i >= 0 {
		return s[i+2:]
	}
	return s
}

func TestGuidesFeedFindingsFromStdin(t *testing.T) {
	re := regexp.MustCompile(`rhiza run (findings|verdict)[^\n]*`)
	for _, name := range []string{"verifier", "review", "plan-review"} {
		g, ok := Lookup(name)
		if !ok {
			t.Fatalf("no guide %q", name)
		}
		body := g.Body
		hits := re.FindAllString(body, -1)
		for _, h := range hits {
			if strings.Contains(h, "<file>") {
				t.Errorf("%s: file form remains: %q", name, h)
			}
		}
		if strings.Count(body, "<<'EOF'") == 0 {
			t.Errorf("%s: no stdin heredoc form", name)
		}
	}
}
