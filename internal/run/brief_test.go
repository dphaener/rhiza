package run

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func briefInput() BriefInput {
	return BriefInput{
		ChangeTitle: "Every agent session a run", ChangeID: "01ABC", Intent: "files",
		Leaves:      []BriefLeaf{{"a", "does a", "go test a"}, {"b", "does b", "go test b"}},
		Decisions:   []string{"d1"},
		Constraints: []string{"c1", "c2"},
		LocalGate:   "make gate",
	}
}

func TestBriefFiles(t *testing.T) {
	in := briefInput()
	got := RenderBrief(in)
	if got != RenderBrief(in) {
		t.Fatal("same input rendered different bytes")
	}
	last := -1
	for _, s := range []string{"## Change set", "Every agent session a run", "01ABC", "files", "### a", "go test a", "### b", "go test b",
		"## Decisions", "- d1", "## Constraints", "- c2", "## Check commands", "rhiza.ship.localGate: make gate", "## Contract"} {
		i := strings.Index(got, s)
		if i < 0 || i < last {
			t.Fatalf("%q missing or out of order in:\n%s", s, got)
		}
		last = i
	}
	if strings.Contains(got, "prePushCheck") {
		t.Error("unset prePushCheck listed")
	}
	if !strings.HasSuffix(got, "Full contract: rhiza guide contract\n") {
		t.Errorf("last line wrong:\n%s", got)
	}
	in.LocalGate = ""
	if none := RenderBrief(in); !strings.Contains(none, "## Check commands\n\nnone\n") {
		t.Errorf("no check commands should say none:\n%s", none)
	}
	in.PrePushCheck = "make self-check"
	if both := RenderBrief(in); !strings.Contains(both, "rhiza.ship.prePushCheck: make self-check") {
		t.Errorf("prePushCheck missing:\n%s", both)
	}
}

func TestSessionDir(t *testing.T) {
	cd := "/repo/.git"
	dirs := map[string]string{
		"agent-1":         SessionDir(cd, "r1", "conduct", 1, AgentSessionName()),
		"agent-2":         SessionDir(cd, "r1", "conduct", 2, AgentSessionName()),
		"builder-a-1":     SessionDir(cd, "r1", "build", 1, LeafSessionName("builder", "a", 1)),
		"builder-a-2":     SessionDir(cd, "r1", "build", 1, LeafSessionName("builder", "a", 2)),
		"check":           SessionDir(cd, "r1", "review", 1, CheckSessionName("qa", "check", 1)),
		"critic-1":        SessionDir(cd, "r1", "review", 1, CheckSessionName("qa", "critic", 1)),
		"critic-2":        SessionDir(cd, "r1", "review", 1, CheckSessionName("quality", "critic", 1)),
		"check-cycle2":    SessionDir(cd, "r1", "review", 1, CheckSessionName("qa", "check", 2)),
		"leaffixer-1":     SessionDir(cd, "r1", "review", 1, LeafFixSessionName("fixer", "a", 1, 1)),
		"leafverifier-1":  SessionDir(cd, "r1", "review", 1, LeafFixSessionName("verifier", "a", 1, 1)),
		"leaffixer-2":     SessionDir(cd, "r1", "review", 1, LeafFixSessionName("fixer", "a", 1, 2)),
		"csfixer-1":       SessionDir(cd, "r1", "review", 1, ChangeSetFixSessionName(1, 1)),
		"csfixer-2":       SessionDir(cd, "r1", "review", 1, ChangeSetFixSessionName(1, 2)),
		"leaffixer-other": SessionDir(cd, "r1", "review", 1, LeafFixSessionName("fixer", "b", 1, 1)),
	}
	seen := map[string]string{}
	for name, d := range dirs {
		if other, ok := seen[d]; ok {
			t.Errorf("%s and %s share %s", name, other, d)
		}
		seen[d] = name
		if !strings.HasPrefix(d, filepath.Join(cd, "rhiza", "runs", "r1")) {
			t.Errorf("%s: %s outside the run's directory", name, d)
		}
	}
	if got, want := dirs["builder-a-1"], filepath.Join(cd, "rhiza/runs/r1/build-1/builder-a-1"); got != want {
		t.Errorf("builder dir = %s, want %s", got, want)
	}
}

func TestSessionDirWrite(t *testing.T) {
	root := t.TempDir()
	dir := SessionDir(root, "r1", "build", 1, LeafSessionName("builder", "a", 1))
	if err := WriteSessionFiles(dir, "brief", []byte(`{"x":1}`)); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "slice.json")); string(b) != `{"x":1}` {
		t.Errorf("slice.json = %q", b)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "brief.md")); string(b) != "brief" {
		t.Errorf("brief.md = %q", b)
	}
	file := filepath.Join(root, "file")
	os.WriteFile(file, nil, 0o644)
	bad := filepath.Join(file, "sub")
	err := WriteSessionFiles(bad, "b", nil)
	if err == nil || !strings.Contains(err.Error(), bad) {
		t.Errorf("error = %v, want one naming %s", err, bad)
	}
}

func TestHandbackFile(t *testing.T) {
	got, err := RenderHandback(HandbackCause{Step: "gate", Outcome: "failed",
		Commands: []FailedCommand{{"make gate", "/o/gate.log"}}})
	if err != nil || !strings.Contains(got, "gate") || !strings.Contains(got, "failed") ||
		!strings.Contains(got, "`make gate`") || !strings.Contains(got, "/o/gate.log") {
		t.Fatalf("gate: %v\n%s", err, got)
	}
	got, _ = RenderHandback(HandbackCause{Step: "ci", Outcome: "red", CIRuns: []FailedCI{{"run 42", "/l/42.log"}}})
	if !strings.Contains(got, "run 42") || !strings.Contains(got, "/l/42.log") {
		t.Errorf("ci:\n%s", got)
	}
	got, _ = RenderHandback(HandbackCause{Step: "merge", Outcome: "conflict", Conflicts: []string{"a.go", "b/c.go"}})
	if !strings.Contains(got, "- a.go") || !strings.Contains(got, "- b/c.go") {
		t.Errorf("conflict:\n%s", got)
	}
	if _, err := RenderHandback(HandbackCause{Outcome: "failed"}); err == nil {
		t.Error("cause with no failing step accepted")
	}
	dir := filepath.Join(t.TempDir(), "s")
	if err := WriteHandback(dir, HandbackCause{Step: "gate", Outcome: "failed"}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "handback.md")); err != nil {
		t.Error(err)
	}
	if err := WriteHandback(dir, HandbackCause{}); err == nil {
		t.Error("empty cause written")
	}
}
