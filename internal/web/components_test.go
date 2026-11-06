package web

import (
	"bytes"
	"flag"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "rewrite the golden files")

// renderPartial renders one named template of the overview page's set.
func renderPartial(t *testing.T, name string, data any) string {
	t.Helper()
	var buf bytes.Buffer
	if err := pages["overview"].ExecuteTemplate(&buf, name, data); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return buf.String()
}

func golden(t *testing.T, name, got string) {
	t.Helper()
	path := filepath.Join("testdata", name+".golden.html")
	if *update {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%v (regenerate with -update)", err)
	}
	if string(want) != got {
		t.Errorf("%s differs from %s (regenerate with -update):\n%s", name, path, got)
	}
}

func TestLayoutGolden(t *testing.T) {
	for name, p := range map[string]page{
		"layout-overview": {Title: "Overview", Nav: "overview", Daemon: Pill{"running", "daemon · 2/4"}},
		"layout-runs":     {Title: "Runs", Nav: "runs"},
	} {
		var buf bytes.Buffer
		if err := pages["overview"].ExecuteTemplate(&buf, "layout", p); err != nil {
			t.Fatal(err)
		}
		golden(t, name, buf.String())
	}
}

func TestPartialsGolden(t *testing.T) {
	var pills strings.Builder
	for _, s := range []string{"running", "waiting", "done", "failed", "interrupted", "draft", "proposed", "accepted", "landed", "withdrawn"} {
		pills.WriteString(renderPartial(t, "pill", Pill{s, s}) + "\n")
	}
	golden(t, "pills", pills.String())

	var glyphs strings.Builder
	for _, k := range []string{"intent", "capability", "decision", "constraint", "question"} {
		glyphs.WriteString(renderPartial(t, "glyph", Glyph{k, false}) + renderPartial(t, "glyph", Glyph{k, true}) + "\n")
	}
	golden(t, "glyphs", glyphs.String())

	rows := []Row{
		{Href: "/runs/01M6PF31", Title: "Retry <webhook> delivery", Sub: "Webhook retries · M3", Pill: Pill{"running", "running"}, Step: "build · 6m", Live: true, Time: "4s ago", Datetime: "2026-11-05T10:42:44Z"},
		{Href: "/runs/01M6PE7Q", Title: "Idempotency keys", Pill: Pill{"waiting", "waiting on you"}, Step: "gate", Time: "12m ago"},
	}
	golden(t, "rows", renderPartial(t, "rows", rows))

	golden(t, "steps", renderPartial(t, "steps", []Step{
		{Name: "Plan", State: "done", Pill: Pill{"done", "planned"}, Duration: "2m 41s", Note: "3 leaves.", Sessions: []Session{{"/sessions/01M6PF3A", "planner", "opus", "41k tok · 2m 38s"}}},
		{Name: "Build", State: "running", Pill: Pill{"running", "running"}, Duration: "6m 12s"},
		{Name: "Review", State: "pending"},
	}))

	golden(t, "term", renderPartial(t, "term", Term{Label: "Live log", Head: "log · tail", Path: "session 01M6PG2C", Live: true, Lines: []TermLine{
		{"10:41:52", "", "build   reading"},
		{"10:42:31", "g", "✓ 14 passed <ok>"},
		{"10:42:33", "w", "warn"},
	}}))

	golden(t, "empty", renderPartial(t, "empty", Empty{"Nothing running", "Start a run from the CLI."}))
	golden(t, "stale", renderPartial(t, "stale", nil))

	body := template.HTML(renderPartial(t, "empty", Empty{Title: "None"}))
	golden(t, "panel", renderPartial(t, "panel", Panel{Kicker: "Inbox", Title: "3 waiting on you", Aside: "oldest 38m", Flush: true, Waiting: true, Body: body}))
}
