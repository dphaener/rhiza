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
// results), and the machine-wide lock taken around cmd/rhiza's scripts,
// while the other packages run unlocked beside the scripts. The timing
// budgets are skipped and never run: only make test-budgets runs them.
// cmd/rhiza's scripts, and only they, run at -parallel min(2*NPROC, 8), or
// RHIZA_TEST_SCRIPT_PARALLEL when it is a positive integer.
func TestMakeTestTimeout(t *testing.T) {
	if _, err := exec.LookPath("make"); err != nil {
		t.Skip("make is not on PATH")
	}
	// NPROC and RHIZA_TEST_SCRIPT_PARALLEL on the command line override
	// the machine's CPU count and the caller's environment.
	cmd := exec.Command("make", "-n", "test", "TEST_LOCK=/locks/rhiza-test.lock", "NPROC=4", "RHIZA_TEST_SCRIPT_PARALLEL=")
	cmd.Dir = "../.."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make -n test: %v\n%s", err, out)
	}
	got := string(out)
	// The suite runs with the timing budgets skipped; every go test keeps
	// the 20m timeout and closes the lock's descriptor, so no leftover
	// process can hold it.
	const budgets = "^(TestReadBudget|TestReadBudgetApply|TestBudget|TestCatchUpBudget|TestDriverPollIsCheap|TestInboxGateQuestionsBudget|TestRootIdlePollBudget)$"
	skip := "go test -trimpath -timeout 20m -skip '" + budgets + "' "
	rest := skip + "$(go list ./... | grep -v '/cmd/rhiza$') >$rest 2>&1 9>&- &"
	// Only cmd/rhiza's scripts get -parallel (NPROC=4: 8); the other
	// packages keep go test's default.
	scripts := "go test -trimpath -timeout 20m -parallel 8 -skip '" + budgets + "' ./cmd/rhiza 9>&-"
	lockScripts := `lock "the cmd/rhiza scripts"`
	for _, want := range []string{
		rest, scripts, lockScripts,
		"lockfile='/locks/rhiza-test.lock'", "flock -n 9", "Waiting for another make test",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("make -n test has no %s line:\n%s", want, got)
		}
	}
	// make test runs no budget test: no -run of the budget regex, no
	// budget phase, no lock taken for it.
	for _, bad := range []string{
		"-run '" + budgets + "'", "Timing budgets", `lock "the timing budgets"`, "RHIZA_BUDGET_LOG",
	} {
		if strings.Contains(got, bad) {
			t.Errorf("make -n test runs the timing budgets (%s); only make test-budgets may:\n%s", bad, got)
		}
	}
	// In order: the other packages start in the background, the lock is
	// taken, the scripts run and the lock is dropped.
	at := func(s string) int { return strings.Index(got, s) }
	order := []string{rest, lockScripts, scripts, "unlock;"}
	for i := 1; i < len(order); i++ {
		if at(order[i-1]) < 0 || at(order[i]) < 0 || at(order[i-1]) > at(order[i]) {
			t.Errorf("make -n test: %q does not come before %q:\n%s", order[i-1], order[i], got)
		}
	}
	if n := strings.Count(got, "-parallel"); n != 1 {
		t.Errorf("make -n test passes -parallel %d times, want once (cmd/rhiza only):\n%s", n, got)
	}
	checkScriptParallel(t, "test")
}

// checkScriptParallel checks the -parallel that make -n <target> passes
// cmd/rhiza's scripts: twice NPROC capped at 8, or RHIZA_TEST_SCRIPT_PARALLEL
// when it is a positive integer, else the NPROC default. Each case sets both
// on make's command line, so neither the machine nor the caller leaks in.
func checkScriptParallel(t *testing.T, target ...string) {
	t.Helper()
	for _, c := range []struct {
		nproc, override, want string
	}{
		{"4", "", "8"},
		{"2", "", "4"},
		{"16", "", "8"},
		{"1", "", "2"},
		{"4", "3", "3"},
		{"4", "12", "12"},
		// Leading zeros are read in decimal: go's flag would read 08 as
		// octal and reject it.
		{"4", "08", "8"},
		{"4", "007", "7"},
		// Not a positive integer: the NPROC default, never -parallel 0.
		{"4", "0", "8"},
		{"4", "00", "8"},
		{"2", "-2", "4"},
		{"16", "abc", "8"},
		{"2", "2x", "4"},
		// No CPU count to read: one CPU.
		{"abc", "", "2"},
	} {
		args := append([]string{"--no-print-directory", "-n"}, target...)
		args = append(args, "TEST_LOCK=/locks/rhiza-test.lock", "NPROC="+c.nproc, "RHIZA_TEST_SCRIPT_PARALLEL="+c.override)
		cmd := exec.Command("make", args...)
		cmd.Dir = "../.."
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("make %s: %v\n%s", strings.Join(args, " "), err, out)
			continue
		}
		var got []string
		fields := strings.Fields(string(out))
		for i, f := range fields {
			if f == "-parallel" && i+1 < len(fields) {
				got = append(got, fields[i+1])
			}
		}
		if len(got) != 1 || got[0] != c.want {
			t.Errorf("make %s: -parallel %v, want [%s]:\n%s", strings.Join(args, " "), got, c.want, out)
		}
	}
}

// TestClaudeMDCommandsNameMakeUpdate checks that the code block under
// AGENTS.md's "## Commands" heading has a line starting with `make update`.
func TestClaudeMDCommandsNameMakeUpdate(t *testing.T) {
	doc, err := os.ReadFile("../../AGENTS.md")
	if err != nil {
		t.Fatal(err)
	}
	block, ok := commandsBlock(string(doc))
	if !ok {
		t.Fatal("AGENTS.md has no code block under ## Commands")
	}
	for _, line := range strings.Split(block, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "make update") {
			return
		}
	}
	t.Errorf("AGENTS.md's ## Commands block has no make update line:\n%s", block)
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
// listed tests, budgets excluded; test-budgets prints the read-budget timings;
// test-scripts passes the scripts' -parallel, test-rest and test-budgets none.
func TestMakeTestShardTargets(t *testing.T) {
	if _, err := exec.LookPath("make"); err != nil {
		t.Skip("make is not on PATH")
	}
	makeRun := func(env []string, args ...string) (string, error) {
		t.Helper()
		// --no-print-directory: under a parent make (CI runs make
		// test-rest) the child would add Entering/Leaving directory lines.
		cmd := exec.Command("make", append([]string{"--no-print-directory"}, args...)...)
		cmd.Dir = "../.."
		cmd.Env = append(os.Environ(), env...)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	make := func(args ...string) string {
		t.Helper()
		out, err := makeRun(nil, args...)
		if err != nil {
			t.Fatalf("make %s: %v\n%s", strings.Join(args, " "), err, out)
		}
		return out
	}
	const budgets = "^(TestReadBudget|TestReadBudgetApply|TestBudget|TestCatchUpBudget|TestDriverPollIsCheap|TestInboxGateQuestionsBudget|TestRootIdlePollBudget)$"
	lock := "TEST_LOCK=/locks/rhiza-test.lock"

	// test-scripts pipes go test -json through tools/scripttimings replay
	// under the lock, and keeps go test's status (sh has no pipefail, so it
	// goes through a file), falling back to the replay's.
	scripts := make("-n", "test-scripts", "SHARD=2/3", lock, "NPROC=4", "RHIZA_TEST_SCRIPT_PARALLEL=")
	for _, want := range []string{
		`go build -trimpath -o "$replay" ./tools/scripttimings`,
		`lock "the cmd/rhiza scripts"`,
		"{ RHIZA_TEST_SHARD='2/3' go test -trimpath -timeout 20m -parallel 8 -json -run '^TestScripts$' -skip '" + budgets + `' ./cmd/rhiza 9>&-; echo $? >"$st"; }`,
		`| "$replay" replay 9>&-`,
		`rstatus=$?; unlock`,
		`status=$(cat "$st")`,
		`[ "$status" != 0 ] || status=$rstatus`,
		`exit $status`,
	} {
		if !strings.Contains(scripts, want) {
			t.Errorf("make -n test-scripts has no %s:\n%s", want, scripts)
		}
	}

	checkScriptParallel(t, "test-scripts", "SHARD=1/4")

	rest := make("-n", "test-rest")
	if want := "go test -trimpath -timeout 20m -skip '" + budgets[:len(budgets)-1] + "$|^TestScripts$'"; !strings.Contains(rest, want) ||
		!strings.Contains(rest, "grep -v '/internal/api$'") {
		t.Errorf("make -n test-rest: want %s and no internal/api:\n%s", want, rest)
	}

	budget := make("-n", "test-budgets", lock)
	if want := "go test -trimpath -timeout 20m -count=1 -p 1 -run '" + budgets + "' ./internal/api ./internal/index ./internal/slice 9>&-"; !strings.Contains(budget, want) {
		t.Errorf("make -n test-budgets has no %s:\n%s", want, budget)
	}

	// -parallel is cmd/rhiza's scripts' alone.
	for name, out := range map[string]string{"test-rest": rest, "test-budgets": budget} {
		if strings.Contains(out, "-parallel") {
			t.Errorf("make -n %s passes -parallel, which only cmd/rhiza's scripts get:\n%s", name, out)
		}
	}

	lint := make("-n", "lint")
	for _, want := range []string{"gofmt -l", "go vet ./...", "check"} {
		if !strings.Contains(lint, want) {
			t.Errorf("make -n lint has no %s:\n%s", want, lint)
		}
	}

	// test-api: a fake go test records the arguments each shard passes in
	// FAKE_ARGS and prints the go test -json stream in FAKE_JSON, then exits
	// FAKE_EXIT. make test-api pipes it through tools/scripttimings replay.
	dir := t.TempDir()
	fakeGo := dir + "/fakego"
	if err := os.WriteFile(fakeGo, []byte("#!/bin/sh\nprintf '%s\\n' \"$*\" >\"$FAKE_ARGS\"\ncat \"$FAKE_JSON\"\nexit \"$FAKE_EXIT\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	const (
		apiPkg  = `"Package":"github.com/x/rhiza/internal/api"`
		apiPass = `{"Action":"start",` + apiPkg + "}\n" +
			`{"Action":"run",` + apiPkg + `,"Test":"TestA"}` + "\n" +
			`{"Action":"output",` + apiPkg + `,"Test":"TestA","Output":"=== RUN   TestA\n"}` + "\n" +
			`{"Action":"pass",` + apiPkg + `,"Test":"TestA","Elapsed":1.5}` + "\n" +
			`{"Action":"output",` + apiPkg + `,"Output":"ok  \tgithub.com/x/rhiza/internal/api\t1.6s\n"}` + "\n" +
			`{"Action":"pass",` + apiPkg + `,"Elapsed":1.6}` + "\n"
		apiFail = `{"Action":"start",` + apiPkg + "}\n" +
			`{"Action":"run",` + apiPkg + `,"Test":"TestB"}` + "\n" +
			`{"Action":"output",` + apiPkg + `,"Test":"TestB","Output":"--- FAIL: TestB (0.10s)\n"}` + "\n" +
			`{"Action":"fail",` + apiPkg + `,"Test":"TestB","Elapsed":0.1}` + "\n" +
			`{"Action":"output",` + apiPkg + `,"Output":"FAIL\tgithub.com/x/rhiza/internal/api\t0.2s\n"}` + "\n" +
			`{"Action":"fail",` + apiPkg + `,"Elapsed":0.2}` + "\n"
		apiBuild = `{"Action":"build-output","ImportPath":"github.com/x/rhiza/internal/api","Output":"internal/api/x.go:1:1: undefined: nope\n"}` + "\n" +
			`{"Action":"build-fail","ImportPath":"github.com/x/rhiza/internal/api"}` + "\n" +
			`{"Action":"start",` + apiPkg + "}\n" +
			`{"Action":"output",` + apiPkg + `,"Output":"FAIL\tgithub.com/x/rhiza/internal/api [build failed]\n"}` + "\n" +
			`{"Action":"fail",` + apiPkg + `,"Elapsed":0}` + "\n"
		apiKilled = `{"Action":"start",` + apiPkg + "}\n" +
			`{"Action":"run",` + apiPkg + `,"Test":"TestC"}` + "\n" +
			`{"Action":"output",` + apiPkg + `,"Test":"TestC","Output":"=== RUN   TestC\n"}` + "\n"
	)
	fakeEnv := func(name, stream, exit string) []string {
		t.Helper()
		js := dir + "/" + name + ".json"
		if err := os.WriteFile(js, []byte(stream), 0o644); err != nil {
			t.Fatal(err)
		}
		return []string{"FAKE_ARGS=" + dir + "/" + name + ".args", "FAKE_JSON=" + js, "FAKE_EXIT=" + exit}
	}

	// make test-api keeps go test's exit status through the pipe: it fails
	// on a failed test, a failed build, a killed go test (no package end) and
	// a non-zero go test exit whatever its stream says; it passes, printing
	// the test-time lines, only when both pass.
	for _, c := range []struct {
		name, stream, exit string
		ok                 bool
		want               string
	}{
		{"pass", apiPass, "0", true, "ok  \tgithub.com/x/rhiza/internal/api\t1.6s\ntest-time\tTestA\t1.5\n"},
		{"failed test", apiFail, "1", false, "--- FAIL: TestB (0.10s)\n"},
		{"failed test, go test exit 0", apiFail, "0", false, "--- FAIL: TestB (0.10s)\n"},
		{"failed build", apiBuild, "1", false, "undefined: nope"},
		{"killed", apiKilled, "137", false, "no package end event"},
		{"passing stream, go test exit 1", apiPass, "1", false, "test-time\tTestA\t1.5\n"},
	} {
		out, err := makeRun(fakeEnv("status", c.stream, c.exit), "test-api", "SHARD=1/2", "GO_TEST="+fakeGo)
		if (err == nil) != c.ok {
			t.Errorf("test-api with %s: err = %v, want ok %v:\n%s", c.name, err, c.ok, out)
		}
		if !strings.Contains(out, c.want) {
			t.Errorf("test-api with %s: output has no %q:\n%s", c.name, c.want, out)
		}
	}

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
		env := fakeEnv("shard", apiPass, "0")
		if out, err := makeRun(env, "test-api", "SHARD="+shard, "GO_TEST="+fakeGo); err != nil {
			t.Fatalf("make test-api SHARD=%s: %v\n%s", shard, err, out)
		}
		args, err := os.ReadFile(strings.TrimPrefix(env[0], "FAKE_ARGS="))
		if err != nil {
			t.Fatal(err)
		}
		out := string(args)
		if !strings.Contains(out, "-trimpath") || !strings.Contains(out, " -json ") || !strings.HasSuffix(strings.TrimSpace(out), "./internal/api") {
			t.Errorf("test-api %s line lacks -trimpath, -json or the package:\n%s", shard, out)
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

	// Every name in the budget regex is a test in the budget packages, so
	// no -run or -skip carries a dead entry.
	budgetList := exec.Command("go", "test", "-trimpath", "-list", budgets, "./internal/api", "./internal/index", "./internal/slice")
	budgetList.Dir = "../.."
	budgetListed, err := budgetList.Output()
	if err != nil {
		t.Fatalf("go test -list %s: %v", budgets, err)
	}
	found := map[string]bool{}
	for _, name := range strings.Fields(string(budgetListed)) {
		found[name] = true
	}
	for _, b := range strings.Split(strings.Trim(budgets, "^()$"), "|") {
		if !found[b] {
			t.Errorf("budget test %s is in no budget package (go test -list):\n%s", b, budgetListed)
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
