package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const prevFile = "alpha\t1\nbravo\t1\ncharlie\t7\ndelta\t1\nold\t5\n"

// The merge of either fixture: fresh times for the passing scripts, charlie
// (failed) keeps its 7, echo (new, unmeasured) gets 1, and old (no longer a
// script) and gone (measured but no longer a script) are dropped.
const wantFile = "alpha\t4.21\nbravo\t12.5\ncharlie\t7\ndelta\t0.33\necho\t1\n"

// setup makes a script directory of alpha..echo and a previous timings file.
func setup(t *testing.T, prev string) (file, dir string) {
	t.Helper()
	root := t.TempDir()
	dir = filepath.Join(root, "script")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, n := range []string{"alpha", "bravo", "charlie", "delta", "echo"} {
		if err := os.WriteFile(filepath.Join(dir, n+".txtar"), []byte("exec true\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("not a script\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	file = filepath.Join(root, "script-timings.tsv")
	if err := os.WriteFile(file, []byte(prev), 0o644); err != nil {
		t.Fatal(err)
	}
	return file, dir
}

func fixtureLog(t *testing.T) func(string) (io.Reader, error) {
	return func(id string) (io.Reader, error) {
		if id != "123" {
			t.Errorf("run id = %q, want 123", id)
		}
		f, err := os.Open("testdata/run-log.txt")
		if err != nil {
			return nil, err
		}
		t.Cleanup(func() { f.Close() })
		return f, nil
	}
}

func readString(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestMergeFromJSON(t *testing.T) {
	file, dir := setup(t, prevFile)
	var out bytes.Buffer
	if err := run([]string{"merge", "-json", "testdata/scripts.json", "-file", file, "-scripts", dir}, nil, nil, &out, nil, nil); err != nil {
		t.Fatalf("merge: %v\n%s", err, out.String())
	}
	if got := readString(t, file); got != wantFile {
		t.Errorf("file =\n%s\nwant\n%s", got, wantFile)
	}
}

func TestMergeFromRunLog(t *testing.T) {
	file, dir := setup(t, prevFile)
	var out bytes.Buffer
	if err := run([]string{"merge", "-run", "123", "-file", file, "-scripts", dir}, nil, nil, &out, fixtureLog(t), nil); err != nil {
		t.Fatalf("merge: %v\n%s", err, out.String())
	}
	if got := readString(t, file); got != wantFile {
		t.Errorf("file =\n%s\nwant\n%s", got, wantFile)
	}
}

func TestMergeNoPreviousFile(t *testing.T) {
	file, dir := setup(t, "")
	os.Remove(file)
	if err := run([]string{"merge", "-json", "testdata/scripts.json", "-file", file, "-scripts", dir}, nil, nil, io.Discard, nil, nil); err != nil {
		t.Fatal(err)
	}
	want := "alpha\t4.21\nbravo\t12.5\ncharlie\t1\ndelta\t0.33\necho\t1\n"
	if got := readString(t, file); got != want {
		t.Errorf("file =\n%s\nwant\n%s", got, want)
	}
}

// Failure cases leave the previous file byte for byte as it was.
func TestMergeFailuresLeaveFileUnchanged(t *testing.T) {
	empty := func(string) (io.Reader, error) {
		return strings.NewReader("scripts (1/4)\tRun scripts\t2026-11-15T10:01:00Z ok\n"), nil
	}
	cases := []struct {
		name string
		prev string
		args []string
		log  func(string) (io.Reader, error)
		want string
	}{
		{"json with no passing subtest", prevFile, []string{"-json", "testdata/none.json"}, nil, "no passing TestScripts subtest"},
		{"run log with no script-time line", prevFile, []string{"-run", "9"}, empty, "no passing TestScripts subtest"},
		{"malformed previous file, no tab", "alpha 1\n", []string{"-json", "testdata/scripts.json"}, nil, "malformed line"},
		{"malformed previous file, bad seconds", "alpha\tfast\n", []string{"-run", "123"}, fixtureLog(t), "malformed line"},
		{"malformed previous file, name twice", "alpha\t1\nalpha\t2\n", []string{"-json", "testdata/scripts.json"}, nil, "listed twice"},
		{"neither source", prevFile, nil, nil, "exactly one of"},
		{"both sources", prevFile, []string{"-run", "1", "-json", "x"}, nil, "exactly one of"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			file, dir := setup(t, c.prev)
			args := append([]string{"merge"}, c.args...)
			args = append(args, "-file", file, "-scripts", dir)
			err := run(args, nil, nil, io.Discard, c.log, nil)
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("err = %v, want one containing %q", err, c.want)
			}
			if got := readString(t, file); got != c.prev {
				t.Errorf("file changed to\n%s", got)
			}
			if m, _ := filepath.Glob(filepath.Join(filepath.Dir(file), ".script-timings-*")); len(m) > 0 {
				t.Errorf("temp files left: %v", m)
			}
		})
	}
}

func TestParseScriptTimeLines(t *testing.T) {
	in := "job\tstep\t2026-11-15T10:00:00Z script-time\ta\t1.5\r\n" +
		"script-time\ta\t0.5\n" + // a script timed twice keeps its longest time
		"script-time\tb\t2\n" +
		"unrelated line\n"
	got, err := ParseScriptTimeLines(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got["a"] != 1.5 || got["b"] != 2 {
		t.Errorf("got %v", got)
	}
	if _, err := ParseScriptTimeLines(strings.NewReader("x script-time\ta\n")); err == nil {
		t.Error("a script-time line without seconds parsed")
	}
	if _, err := ParseScriptTimeLines(strings.NewReader("")); !errors.Is(err, ErrNoTimes) {
		t.Errorf("empty input: err = %v, want ErrNoTimes", err)
	}
}

func TestFormatScriptTimeRoundTrips(t *testing.T) {
	line := FormatScriptTime("alpha", 4.214)
	if line != "script-time\talpha\t4.21" {
		t.Errorf("line = %q", line)
	}
	got, err := ParseScriptTimeLines(strings.NewReader(line + "\n"))
	if err != nil || got["alpha"] != 4.21 {
		t.Errorf("got %v, %v", got, err)
	}
}

// The checked-in file parses and is sorted name<TAB>seconds lines. (Scripts
// added or removed since the last refresh are fine: merge fills them in.)
func TestCheckedInFileParses(t *testing.T) {
	file := filepath.Join("..", "..", DefaultFile)
	times, err := ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(times) == 0 {
		t.Fatalf("%s is missing or empty", file)
	}
	if got := readString(t, file); got != string(Format(times)) {
		t.Errorf("%s is not sorted name<TAB>seconds lines", file)
	}
}

// The api fixtures' previous file: TestBravo failed and TestCharlie was
// skipped in the fresh run, TestOld is no longer a test.
const prevTestsFile = "TestAlpha\t1\nTestBravo\t7\nTestCharlie\t4\nTestOld\t5\n"

// The merge of either api fixture: fresh times for the passing tests (the
// subtest TestAlpha/sub not among them), TestBravo (failed) and TestCharlie
// (skipped, unrun) keep theirs, TestEcho (new, unmeasured) is left out,
// TestOld and TestGone (no longer listed) are dropped, and TestReadBudget
// (a budget, -skip) is left out though it was timed.
const wantTestsFile = "TestAlpha\t3.5\nTestBravo\t7\nTestCharlie\t4\nTestDelta\t0.25\n"

const budgetSkip = "^(TestReadBudget|TestBudget)$"

func setupTests(t *testing.T, prev string) string {
	t.Helper()
	file := filepath.Join(t.TempDir(), "test-timings.tsv")
	if err := os.WriteFile(file, []byte(prev), 0o644); err != nil {
		t.Fatal(err)
	}
	return file
}

// fakeList stands in for go test -list on ./internal/api.
func fakeList(t *testing.T) func(string) ([]string, error) {
	return func(pkg string) ([]string, error) {
		if pkg != DefaultTestsPkg {
			t.Errorf("pkg = %q, want %s", pkg, DefaultTestsPkg)
		}
		return []string{"TestAlpha", "TestBravo", "TestCharlie", "TestDelta", "TestEcho", "TestReadBudget"}, nil
	}
}

func apiLog(t *testing.T) func(string) (io.Reader, error) {
	return func(id string) (io.Reader, error) {
		if id != "456" {
			t.Errorf("run id = %q, want 456", id)
		}
		f, err := os.Open("testdata/api-run-log.txt")
		if err != nil {
			return nil, err
		}
		t.Cleanup(func() { f.Close() })
		return f, nil
	}
}

func TestMergeTests(t *testing.T) {
	for name, src := range map[string][]string{
		"json":    {"-json", "testdata/api.json"},
		"run log": {"-run", "456"},
	} {
		t.Run(name, func(t *testing.T) {
			file := setupTests(t, prevTestsFile)
			args := append([]string{"merge-tests"}, src...)
			args = append(args, "-file", file, "-skip", budgetSkip)
			var out bytes.Buffer
			if err := run(args, nil, nil, &out, apiLog(t), fakeList(t)); err != nil {
				t.Fatalf("merge-tests: %v\n%s", err, out.String())
			}
			if got := readString(t, file); got != wantTestsFile {
				t.Errorf("file =\n%s\nwant\n%s", got, wantTestsFile)
			}
			if !strings.Contains(out.String(), "4 tests in "+file) {
				t.Errorf("report = %q", out.String())
			}
		})
	}
}

// Failure cases leave the previous file byte for byte as it was.
func TestMergeTestsFailuresLeaveFileUnchanged(t *testing.T) {
	scriptsOnly := func(string) (io.Reader, error) {
		return strings.NewReader("scripts (1/4)\tRun scripts\t2026-11-15T10:01:00Z script-time\talpha\t1\n"), nil
	}
	listFails := func(string) ([]string, error) { return nil, errors.New("go test -list failed") }
	cases := []struct {
		name string
		prev string
		args []string
		log  func(string) (io.Reader, error)
		list func(string) ([]string, error)
		want string
	}{
		{"json with no passing top-level test", prevTestsFile, []string{"-json", "testdata/api-none.json"}, nil, fakeList(t), "no passing top-level test"},
		{"run log with no test-time line", prevTestsFile, []string{"-run", "9"}, scriptsOnly, fakeList(t), "no passing top-level test"},
		{"go test -list fails", prevTestsFile, []string{"-json", "testdata/api.json"}, nil, listFails, "go test -list failed"},
		{"bad skip", prevTestsFile, []string{"-json", "testdata/api.json", "-skip", "("}, nil, fakeList(t), "-skip"},
		{"malformed previous file", "TestAlpha 1\n", []string{"-json", "testdata/api.json"}, nil, fakeList(t), "malformed line"},
		{"neither source", prevTestsFile, nil, nil, fakeList(t), "exactly one of"},
		{"both sources", prevTestsFile, []string{"-run", "1", "-json", "x"}, nil, fakeList(t), "exactly one of"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			file := setupTests(t, c.prev)
			args := append([]string{"merge-tests"}, c.args...)
			args = append(args, "-file", file)
			err := run(args, nil, nil, io.Discard, c.log, c.list)
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("err = %v, want one containing %q", err, c.want)
			}
			if got := readString(t, file); got != c.prev {
				t.Errorf("file changed to\n%s", got)
			}
		})
	}
}

func TestFormatTestTimeRoundTrips(t *testing.T) {
	line := FormatTestTime("TestAlpha", 3.504)
	if line != "test-time\tTestAlpha\t3.5" {
		t.Errorf("line = %q", line)
	}
	got, err := ParseTestTimeLines(strings.NewReader("job\tstep\t2026-11-15T10:00:00Z " + line + "\r\n"))
	if err != nil || len(got) != 1 || got["TestAlpha"] != 3.5 {
		t.Errorf("got %v, %v", got, err)
	}
	if _, err := ParseTestTimeLines(strings.NewReader("x test-time\tTestA\n")); err == nil {
		t.Error("a test-time line without seconds parsed")
	}
	if _, err := ParseTestTimeLines(strings.NewReader("x script-time\ta\t1\n")); !errors.Is(err, ErrNoTestTimes) {
		t.Errorf("script-time only: err = %v, want ErrNoTestTimes", err)
	}
}

// The checked-in api file parses and is sorted name<TAB>seconds lines.
func TestCheckedInTestsFileParses(t *testing.T) {
	file := filepath.Join("..", "..", DefaultTestsFile)
	times, err := ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(times) == 0 {
		t.Fatalf("%s is missing or empty", file)
	}
	if got := readString(t, file); got != string(Format(times)) {
		t.Errorf("%s is not sorted name<TAB>seconds lines", file)
	}
}
