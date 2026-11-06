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

// TestMakeTestTimeout checks that `make test` gives go test 20 minutes, so the
// script suite does not hit go test's 10-minute default on CI's slow runner.
// It also checks what keeps many worktrees' suites cheap on one machine:
// -trimpath (a second worktree reuses the first one's builds and cached
// results), and the machine-wide lock taken around cmd/rhiza's scripts and
// around the timing budgets, while the other packages run unlocked beside
// the scripts.
func TestMakeTestTimeout(t *testing.T) {
	if _, err := exec.LookPath("make"); err != nil {
		t.Skip("make is not on PATH")
	}
	cmd := exec.Command("make", "-n", "test", "TEST_LOCK=/locks/rhiza-test.lock")
	cmd.Dir = "../.."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make -n test: %v\n%s", err, out)
	}
	got := string(out)
	// The suite runs with the timing budgets skipped, then the budgets run
	// alone, one package at a time; every phase keeps the 20m timeout and
	// closes the lock's descriptor, so no leftover process can hold it.
	skip := "go test -trimpath -timeout 20m -skip '^(TestReadBudget|TestReadBudgetApply|TestBudget|TestCatchUpBudget|TestDriverPollIsCheap|TestGeneratedModel)$' "
	rest := skip + "$(go list ./... | grep -v '/cmd/rhiza$') >$rest 2>&1 9>&- &"
	scripts := skip + "./cmd/rhiza 9>&-"
	budgets := "go test -trimpath -timeout 20m -count=1 -p 1 -run '^(TestReadBudget|TestReadBudgetApply|TestBudget|TestCatchUpBudget|TestDriverPollIsCheap|TestGeneratedModel)$' ./internal/api ./internal/index ./internal/slice 9>&-"
	lockScripts := `lock "the cmd/rhiza scripts"`
	lockBudgets := `lock "the timing budgets"`
	for _, want := range []string{
		rest, scripts, budgets, lockScripts, lockBudgets,
		"lockfile='/locks/rhiza-test.lock'", "flock -n 9", "Waiting for another make test",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("make -n test has no %s line:\n%s", want, got)
		}
	}
	// In order: the other packages start in the background, the lock is
	// taken, the scripts run and the lock is dropped; then the budgets run
	// under the lock again.
	at := func(s string) int { return strings.Index(got, s) }
	order := []string{rest, lockScripts, scripts, "unlock;", lockBudgets, budgets}
	for i := 1; i < len(order); i++ {
		if at(order[i-1]) < 0 || at(order[i]) < 0 || at(order[i-1]) > at(order[i]) {
			t.Errorf("make -n test: %q does not come before %q:\n%s", order[i-1], order[i], got)
		}
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

// TestMakeTestShardTargets checks the per-job test targets CI runs in
// parallel: each go test line keeps -trimpath and either skips the timing
// budgets or runs only them; test-api's shards partition internal/api's
// listed tests, budgets excluded; test-budgets prints the read-budget timings.
func TestMakeTestShardTargets(t *testing.T) {
	if _, err := exec.LookPath("make"); err != nil {
		t.Skip("make is not on PATH")
	}
	make := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("make", args...)
		cmd.Dir = "../.."
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("make %s: %v\n%s", strings.Join(args, " "), err, out)
		}
		return string(out)
	}
	const budgets = "^(TestReadBudget|TestReadBudgetApply|TestBudget|TestCatchUpBudget|TestDriverPollIsCheap|TestGeneratedModel)$"
	lock := "TEST_LOCK=/locks/rhiza-test.lock"

	scripts := make("-n", "test-scripts", "SHARD=2/3", lock)
	for _, want := range []string{
		"RHIZA_TEST_SHARD='2/3' go test -trimpath -timeout 20m -run '^TestScripts$' -skip '" + budgets + "' ./cmd/rhiza 9>&-",
		`lock "the cmd/rhiza scripts"`,
	} {
		if !strings.Contains(scripts, want) {
			t.Errorf("make -n test-scripts has no %s:\n%s", want, scripts)
		}
	}

	rest := make("-n", "test-rest")
	if want := "go test -trimpath -timeout 20m -skip '" + budgets[:len(budgets)-1] + "$|^TestScripts$'"; !strings.Contains(rest, want) ||
		!strings.Contains(rest, "grep -v '/internal/api$'") {
		t.Errorf("make -n test-rest: want %s and no internal/api:\n%s", want, rest)
	}

	budget := make("-n", "test-budgets", lock)
	if want := "go test -trimpath -timeout 20m -count=1 -p 1 -run '" + budgets + "' ./internal/api ./internal/index ./internal/slice 9>&-"; !strings.Contains(budget, want) {
		t.Errorf("make -n test-budgets has no %s:\n%s", want, budget)
	}

	lint := make("-n", "lint")
	for _, want := range []string{"gofmt -l", "go vet ./...", "check"} {
		if !strings.Contains(lint, want) {
			t.Errorf("make -n lint has no %s:\n%s", want, lint)
		}
	}

	// test-api: GO_TEST=echo prints the arguments each shard would pass.
	list := exec.Command("go", "test", "-trimpath", "-list", "^Test", "./internal/api")
	list.Dir = "../.."
	listed, err := list.Output()
	if err != nil {
		t.Fatalf("go test -list: %v", err)
	}
	want := map[string]bool{}
	for _, name := range strings.Fields(string(listed)) {
		if strings.HasPrefix(name, "Test") {
			want[name] = true
		}
	}
	for _, b := range strings.Split(strings.Trim(budgets, "^()$"), "|") {
		delete(want, b)
	}
	got := map[string]int{}
	for _, shard := range []string{"1/2", "2/2"} {
		out := make("test-api", "SHARD="+shard, "GO_TEST=echo")
		if !strings.Contains(out, "-trimpath") || !strings.HasSuffix(strings.TrimSpace(out), "./internal/api") {
			t.Errorf("test-api %s line lacks -trimpath or the package:\n%s", shard, out)
		}
		_, after, ok := strings.Cut(out, "-run ")
		if !ok {
			t.Fatalf("test-api %s has no -run:\n%s", shard, out)
		}
		re := strings.Fields(after)[0]
		if !strings.HasPrefix(re, "^(") || !strings.HasSuffix(re, ")$") {
			t.Fatalf("test-api %s regex %q is not anchored", shard, re)
		}
		for _, name := range strings.Split(strings.TrimSuffix(strings.TrimPrefix(re, "^("), ")$"), "|") {
			got[name]++
		}
	}
	for name := range want {
		if got[name] != 1 {
			t.Errorf("test %s is in %d shards, want 1", name, got[name])
		}
	}
	for name := range got {
		if !want[name] {
			t.Errorf("shards run %s, which go test -list does not show (or is a budget test)", name)
		}
	}

	// test-budgets prints what the budget tests append to RHIZA_BUDGET_LOG.
	fake := t.TempDir() + "/fakego"
	if err := os.WriteFile(fake, []byte("#!/bin/sh\necho 'read-budget: 12ms' >> \"$RHIZA_BUDGET_LOG\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	out := make("test-budgets", "GO_TEST="+fake, "TEST_LOCK=")
	if !strings.Contains(out, "Read-budget timings\nread-budget: 12ms") {
		t.Errorf("test-budgets did not print the timings:\n%s", out)
	}
}
