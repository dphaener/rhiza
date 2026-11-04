// Package ciworkflow holds the test that pins the shape of Rhiza's own CI
// workflow file.
package ciworkflow

import (
	"os"
	"strings"
	"testing"
)

type line struct {
	indent int
	text   string
}

// block returns the lines nested under the first key line equal to key at
// the given indent within lines, or nil. Only the small YAML subset the
// workflow uses is understood: block mappings, one key per line.
func block(lines []line, indent int, key string) []line {
	for i, l := range lines {
		if l.indent != indent || !strings.HasPrefix(l.text, key+":") {
			continue
		}
		var out []line
		for _, n := range lines[i+1:] {
			if n.indent <= indent {
				break
			}
			out = append(out, n)
		}
		return out
	}
	return nil
}

// value returns the inline value of key at indent, or "" and false.
func value(lines []line, indent int, key string) (string, bool) {
	for _, l := range lines {
		if l.indent == indent && strings.HasPrefix(l.text, key+":") {
			return strings.Trim(strings.TrimSpace(strings.TrimPrefix(l.text, key+":")), `"'`), true
		}
	}
	return "", false
}

func parse(t *testing.T) []line {
	t.Helper()
	data, err := os.ReadFile("../../.github/workflows/ci.yml")
	if err != nil {
		t.Fatal(err)
	}
	var out []line
	for _, raw := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		out = append(out, line{len(raw) - len(strings.TrimLeft(raw, " ")), trimmed})
	}
	return out
}

func TestCIRunsOnPullRequestsOnly(t *testing.T) {
	on := block(parse(t), 0, "on")
	if on == nil {
		t.Fatal("ci.yml has no on: block")
	}
	var triggers []string
	for _, l := range on {
		if l.indent == on[0].indent {
			triggers = append(triggers, strings.TrimSuffix(strings.SplitN(l.text, ":", 2)[0], ":"))
		}
	}
	if len(triggers) != 1 || triggers[0] != "pull_request" {
		t.Errorf("ci.yml triggers = %v; want only pull_request", triggers)
	}
}

func TestCICancelsStaleRuns(t *testing.T) {
	c := block(parse(t), 0, "concurrency")
	if c == nil {
		t.Fatal("ci.yml has no concurrency: block")
	}
	group, _ := value(c, c[0].indent, "group")
	if !strings.Contains(group, "github.event.pull_request.number") {
		t.Errorf("concurrency group = %q; want it keyed on github.event.pull_request.number", group)
	}
	if v, _ := value(c, c[0].indent, "cancel-in-progress"); v != "true" {
		t.Errorf("cancel-in-progress = %q; want true", v)
	}
}
