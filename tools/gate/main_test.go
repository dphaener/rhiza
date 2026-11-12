package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func sh(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%v: %v\n%s", args, err, b)
	}
	return strings.TrimSpace(string(b))
}

func write(t *testing.T, root, name, body string) {
	t.Helper()
	p := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// repo makes a module with a <- b <- cmd/rhiza and an unrelated c, committed
// once; it returns the root and that commit.
func repo(t *testing.T) (string, string) {
	root := t.TempDir()
	sh(t, root, "git", "init", "-q")
	write(t, root, "go.mod", "module example.com/x\n\ngo 1.21\n")
	write(t, root, "internal/a/a.go", "package a\n")
	write(t, root, "internal/b/b.go", "package b\nimport _ \"example.com/x/internal/a\"\n")
	write(t, root, "internal/c/c.go", "package c\n")
	write(t, root, "cmd/rhiza/main.go", "package main\nimport _ \"example.com/x/internal/b\"\nfunc main() {}\n")
	write(t, root, "cmd/rhiza/testdata/script/foo.txtar", "x\n")
	write(t, root, ".rhiza/log/x.json", "{}\n")
	sh(t, root, "git", "add", ".")
	sh(t, root, "git", "commit", "-qm", "base")
	return root, sh(t, root, "git", "rev-parse", "HEAD")
}

func plan(t *testing.T, root, base, change string) string {
	t.Helper()
	sh(t, root, "git", "add", ".")
	sh(t, root, "git", "commit", "-qm", "change", "--allow-empty")
	env := map[string]string{"RHIZA_GATE_BASE": base, "RHIZA_CHANGE": change}
	var out bytes.Buffer
	accept := func(string) []string { return []string{"go test -trimpath ./cmd/rhiza -run TestScripts/bar"} }
	if err := run(root, true, func(k string) string { return env[k] }, &out, accept); err != nil {
		t.Fatal(err)
	}
	return out.String()
}

func TestPlanOnePackageAndDependents(t *testing.T) {
	root, base := repo(t)
	write(t, root, "internal/a/a.go", "package a\n\nvar X = 1\n")
	got := plan(t, root, base, "")
	for _, want := range []string{"go test -trimpath ./internal/a ./internal/b\n", "go test -trimpath ./cmd/rhiza -skip", "make self-check"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in\n%s", want, got)
		}
	}
	if strings.Contains(got, "internal/c") || strings.Contains(got, "make test") || strings.Contains(got, "-run") {
		t.Errorf("too much in\n%s", got)
	}
}

func TestPlanRhizaOnly(t *testing.T) {
	root, base := repo(t)
	write(t, root, ".rhiza/log/y.json", "{}\n")
	got := plan(t, root, base, "")
	if strings.Contains(got, "go test") || !strings.Contains(got, "make self-check") {
		t.Errorf("want self-check only:\n%s", got)
	}
}

func TestPlanFallbacks(t *testing.T) {
	root, base := repo(t)
	if got := plan(t, root, "", ""); !strings.Contains(got, "make test") || !strings.Contains(got, "make self-check") {
		t.Errorf("empty base:\n%s", got)
	}
	for _, f := range []string{"docs/quickstart.md", "AGENTS.md", ".github/workflows/ci.yml"} {
		root, base := repo(t)
		write(t, root, f, "changed\n")
		if got := plan(t, root, base, ""); !strings.Contains(got, "make test") || !strings.Contains(got, "make self-check") {
			t.Errorf("%s:\n%s", f, got)
		}
	}
	write(t, root, "go.mod", "module example.com/x\n\ngo 1.22\n")
	if got := plan(t, root, base, ""); !strings.Contains(got, "make test") || !strings.Contains(got, "make self-check") {
		t.Errorf("go.mod:\n%s", got)
	}
}

func TestPlanScripts(t *testing.T) {
	root, base := repo(t)
	write(t, root, "cmd/rhiza/testdata/script/foo.txtar", "y\n")
	got := plan(t, root, base, "01CHANGE")
	if !strings.Contains(got, "-run 'TestScripts/(bar|foo)$'") {
		t.Errorf("scripts:\n%s", got)
	}
	if strings.Contains(got, "internal/") {
		t.Errorf("only cmd/rhiza changed:\n%s", got)
	}
}

func TestScriptNames(t *testing.T) {
	got := ScriptNames([]string{"-run 'TestScripts/(a|b)$' and TestScripts/c-d"})
	if strings.Join(got, ",") != "a,b,c-d" {
		t.Errorf("got %v", got)
	}
}
