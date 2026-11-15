package main

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
)

func replayFixture(t *testing.T, name string) (string, error) {
	t.Helper()
	f, err := os.Open("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var out bytes.Buffer
	err = run([]string{"replay"}, f, &out, nil, nil, nil)
	return out.String(), err
}

// A passing package: passing and skipped tests' text is dropped, the
// package's ok line is kept, and each passing script gets a script-time line.
func TestReplayPass(t *testing.T) {
	got, err := replayFixture(t, "replay-pass.json")
	if err != nil {
		t.Fatalf("replay: %v\n%s", err, got)
	}
	want := "ok  \tgithub.com/x/rhiza/cmd/rhiza\t16.9s\n" +
		"script-time\talpha\t4.21\n" +
		"script-time\tbravo\t12.5\n" +
		"test-time\tTestScripts\t16.8\n"
	if got != want {
		t.Errorf("output =\n%q\nwant\n%q", got, want)
	}
}

// A package of top-level tests (make test-api): the failed test's text, the
// package's FAIL lines, then a test-time line per passing top-level test
// only (no subtest, no failed or skipped test), and a non-zero exit.
func TestReplayTestTimes(t *testing.T) {
	got, err := replayFixture(t, "api.json")
	if !errors.Is(err, ErrReplayFailed) {
		t.Fatalf("err = %v, want ErrReplayFailed\n%s", err, got)
	}
	want := "=== RUN   TestBravo\n" +
		"    api_test.go:12: want 1, got 2\n" +
		"--- FAIL: TestBravo (1.20s)\n" +
		"FAIL\n" +
		"FAIL\tgithub.com/x/rhiza/internal/api\t16.1s\n" +
		"test-time\tTestAlpha\t3.5\n" +
		"test-time\tTestDelta\t0.25\n" +
		"test-time\tTestGone\t2\n" +
		"test-time\tTestReadBudget\t9\n"
	if got != want {
		t.Errorf("output =\n%s\nwant\n%s", got, want)
	}

	// The lines parse back to what the go test -json gives directly.
	fromLog, err := ParseTestTimeLines(strings.NewReader(got))
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open("testdata/api.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	fromJSON, err := ParseTopTestJSON(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(fromLog) != 4 || len(fromJSON) != len(fromLog) {
		t.Fatalf("from the log %v, from the json %v", fromLog, fromJSON)
	}
	for n, v := range fromJSON {
		if fromLog[n] != v {
			t.Errorf("%s: %v from the log, %v from the json", n, fromLog[n], v)
		}
	}
}

// A failed script: its text with the --- FAIL: line, the parent's --- FAIL:,
// the package's FAIL lines, then the passing scripts' times; a non-zero exit.
func TestReplayFail(t *testing.T) {
	got, err := replayFixture(t, "replay-fail.json")
	if !errors.Is(err, ErrReplayFailed) {
		t.Fatalf("err = %v, want ErrReplayFailed\n%s", err, got)
	}
	want := "=== RUN   TestScripts/charlie\n" +
		"    testscript.go:584: > exec rhiza show x\n" +
		"        FAIL: testdata/script/charlie.txtar:3: unexpected command failure\n" +
		"--- FAIL: TestScripts/charlie (2.00s)\n" +
		"=== RUN   TestScripts\n" +
		"--- FAIL: TestScripts (6.60s)\n" +
		"FAIL\n" +
		"FAIL\tgithub.com/x/rhiza/cmd/rhiza\t6.7s\n" +
		"script-time\talpha\t4.21\n" +
		"script-time\tdelta\t0.33\n"
	if got != want {
		t.Errorf("output =\n%s\nwant\n%s", got, want)
	}
	if strings.Contains(got, "skipped") {
		t.Errorf("a skipped script's text was printed:\n%s", got)
	}
}

// A build failure prints the compiler's output and the package's FAIL line,
// and exits non-zero with no script-time line.
func TestReplayBuildFailure(t *testing.T) {
	got, err := replayFixture(t, "replay-build.json")
	if !errors.Is(err, ErrReplayFailed) {
		t.Fatalf("err = %v, want ErrReplayFailed\n%s", err, got)
	}
	for _, want := range []string{
		"cmd/rhiza/main.go:12:2: undefined: nope\n",
		"FAIL\tgithub.com/x/rhiza/cmd/rhiza [build failed]\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output has no %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, ScriptTimePrefix) {
		t.Errorf("a build failure printed script times:\n%s", got)
	}
}

// A killed go test leaves no package end event: the unfinished tests' text
// is printed, the passing scripts still get their times, and it exits
// non-zero.
func TestReplayKilled(t *testing.T) {
	got, err := replayFixture(t, "replay-killed.json")
	if !errors.Is(err, ErrReplayFailed) {
		t.Fatalf("err = %v, want ErrReplayFailed\n%s", err, got)
	}
	for _, want := range []string{
		"=== RUN   TestScripts/bravo\n    testscript.go:584: bravo log line\n",
		"no package end event for github.com/x/rhiza/cmd/rhiza",
		"script-time\talpha\t4.21\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output has no %q:\n%s", want, got)
		}
	}
}

// Empty input, or input that is not go test -json (go test died before it
// started), has no package end event either. Non-JSON lines pass through.
func TestReplayNoPackageEnd(t *testing.T) {
	for _, in := range []string{"", "signal: killed\n"} {
		var out bytes.Buffer
		err := run([]string{"replay"}, strings.NewReader(in), &out, nil, nil, nil)
		if !errors.Is(err, ErrReplayFailed) {
			t.Errorf("input %q: err = %v, want ErrReplayFailed", in, err)
		}
		if !strings.HasPrefix(out.String(), in) || !strings.Contains(out.String(), "no package end event") {
			t.Errorf("input %q: output = %q", in, out.String())
		}
	}
	if err := run([]string{"replay", "x"}, strings.NewReader(""), &bytes.Buffer{}, nil, nil, nil); err == nil {
		t.Error("replay with an argument ran")
	}
}
