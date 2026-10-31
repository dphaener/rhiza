package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestMakeUpdateDryRun checks the order of what `make update` would run: the
// in-place build into PREFIX comes first, then that binary's rex install, and
// no cp ever appears (a copied Go binary is SIGKILLed on arm64). It also
// checks that `make help` names the update target.
func TestMakeUpdateDryRun(t *testing.T) {
	if _, err := exec.LookPath("make"); err != nil {
		t.Skip("make is not on PATH")
	}
	prefix := t.TempDir()
	cmd := exec.Command("make", "-n", "update", "PREFIX="+prefix)
	cmd.Dir = "../.."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make -n update: %v\n%s", err, out)
	}
	got := string(out)
	build := strings.Index(got, "go build")
	if build < 0 || !strings.Contains(got, "-o "+prefix+"/rhiza") {
		t.Fatalf("no go build -o %s/rhiza line:\n%s", prefix, got)
	}
	install := strings.Index(got, prefix+"/rhiza rex install")
	if install < 0 {
		t.Fatalf("no %s/rhiza rex install line:\n%s", prefix, got)
	}
	if build > install {
		t.Errorf("rex install comes before the build:\n%s", got)
	}
	for _, line := range strings.Split(got, "\n") {
		for _, f := range strings.Fields(line) {
			if f == "cp" {
				t.Errorf("cp appears in make -n update:\n%s", got)
			}
		}
	}

	help := exec.Command("make", "help")
	help.Dir = "../.."
	helpOut, err := help.CombinedOutput()
	if err != nil {
		t.Fatalf("make help: %v\n%s", err, helpOut)
	}
	named := false
	for _, line := range strings.Split(string(helpOut), "\n") {
		if f := strings.Fields(line); len(f) > 0 && f[0] == "update" {
			named = true
		}
	}
	if !named {
		t.Errorf("make help does not name update:\n%s", helpOut)
	}
}

// TestClaudeMDCommandsNameMakeUpdate checks that the code block under
// CLAUDE.md's "## Commands" heading has a line starting with `make update`.
func TestClaudeMDCommandsNameMakeUpdate(t *testing.T) {
	doc, err := os.ReadFile("../../CLAUDE.md")
	if err != nil {
		t.Fatal(err)
	}
	block, ok := commandsBlock(string(doc))
	if !ok {
		t.Fatal("CLAUDE.md has no code block under ## Commands")
	}
	for _, line := range strings.Split(block, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "make update") {
			return
		}
	}
	t.Errorf("CLAUDE.md's ## Commands block has no make update line:\n%s", block)
}

// commandsBlock returns the first fenced code block after the "## Commands"
// heading and before the next "## " heading.
func commandsBlock(doc string) (string, bool) {
	lines := strings.Split(doc, "\n")
	in, fenced := false, false
	var block []string
	for _, line := range lines {
		switch {
		case !in:
			in = strings.TrimSpace(line) == "## Commands"
		case !fenced && strings.HasPrefix(line, "## "):
			return "", false
		case strings.HasPrefix(line, "```"):
			if fenced {
				return strings.Join(block, "\n"), true
			}
			fenced = true
		case fenced:
			block = append(block, line)
		}
	}
	return "", false
}
