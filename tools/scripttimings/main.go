// Command scripttimings keeps cmd/rhiza/testdata/script-timings.tsv, the
// recorded wall time of each cmd/rhiza testscript (name, tab, seconds, sorted
// by name), which the scripts' scheduler and the CI shards weigh by, and
// internal/api/testdata/test-timings.tsv, the same for internal/api's
// top-level tests.
//
// It has three modes, as subcommands. merge (make script-timings) reads fresh
// times from a CI run's scripts jobs (-run <id>, through gh run view --log,
// which carries script-time lines) or from a go test -json file of
// TestScripts (-json <file>), and merges them into the previous file. replay
// (make test-scripts) reads go test -json on stdin and prints go test's text
// for what did not pass, then a script-time line per passing script, so each
// CI scripts job's log carries its scripts' times; it exits non-zero when go
// test failed; it also prints a test-time line per passing top-level test
// (make test-api). merge-tests (make api-timings) is merge for those: fresh
// times from test-time lines (-run) or a go test -json file (-json), merged
// into the file for the package's top-level tests (go test -list) that -skip
// does not match.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, ghLog, goTestList); err != nil {
		fmt.Fprintln(os.Stderr, "scripttimings:", err)
		os.Exit(1)
	}
}

// ghLog returns gh run view <id> --log: every job's log of one CI run.
func ghLog(id string) (io.Reader, error) {
	cmd := exec.Command("gh", "run", "view", id, "--log")
	cmd.Stderr = os.Stderr
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gh run view %s --log: %w", id, err)
	}
	return bytes.NewReader(b), nil
}

// goTestList returns the top-level tests go test -list '^Test' prints for pkg.
func goTestList(pkg string) ([]string, error) {
	cmd := exec.Command("go", "test", "-trimpath", "-list", "^Test", pkg)
	cmd.Stderr = os.Stderr
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("go test -list ^Test %s: %w", pkg, err)
	}
	var tests []string
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "Test") {
			tests = append(tests, line)
		}
	}
	return tests, nil
}

// run runs one mode: merge reports to stderr, replay reads stdin and writes
// stdout.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer, runLog func(id string) (io.Reader, error), listTests func(pkg string) ([]string, error)) error {
	if len(args) == 0 {
		return errors.New("usage: scripttimings merge (-run <id> | -json <file>) [-file <tsv>] [-scripts <dir>]\n       scripttimings merge-tests (-run <id> | -json <file>) [-file <tsv>] [-pkg <pkg>] [-skip <regexp>]\n       scripttimings replay < go-test.json")
	}
	switch args[0] {
	case "merge":
		return merge(args[1:], stderr, runLog)
	case "merge-tests":
		return mergeTests(args[1:], stderr, runLog, listTests)
	case "replay":
		if len(args) > 1 {
			return errors.New("replay: takes no arguments; pipe go test -json into it")
		}
		return Replay(stdin, stdout)
	default:
		return fmt.Errorf("unknown mode %q (want merge, merge-tests or replay)", args[0])
	}
}

func merge(args []string, out io.Writer, runLog func(id string) (io.Reader, error)) error {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	fs.SetOutput(out)
	runID := fs.String("run", "", "CI run id whose scripts jobs print script-time lines (read with gh run view --log)")
	jsonFile := fs.String("json", "", "go test -json output of TestScripts")
	file := fs.String("file", DefaultFile, "the timings file to rewrite")
	scripts := fs.String("scripts", DefaultScripts, "the testscript directory; scripts no longer in it are dropped")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if (*runID == "") == (*jsonFile == "") {
		return errors.New("merge: pass exactly one of -run <id> or -json <file>")
	}
	var fresh map[string]float64
	var err error
	if *runID != "" {
		var r io.Reader
		if r, err = runLog(*runID); err != nil {
			return err
		}
		fresh, err = ParseScriptTimeLines(r)
	} else {
		var f *os.File
		if f, err = os.Open(*jsonFile); err != nil {
			return err
		}
		fresh, err = ParseTestJSON(f)
		f.Close()
	}
	if err != nil {
		return err
	}
	n, err := MergeFile(*file, *scripts, fresh)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "scripttimings: %d fresh times, %d scripts in %s\n", len(fresh), n, *file)
	return nil
}

func mergeTests(args []string, out io.Writer, runLog func(id string) (io.Reader, error), listTests func(pkg string) ([]string, error)) error {
	fs := flag.NewFlagSet("merge-tests", flag.ContinueOnError)
	fs.SetOutput(out)
	runID := fs.String("run", "", "CI run id whose jobs print test-time lines (read with gh run view --log)")
	jsonFile := fs.String("json", "", "go test -json output of the package")
	file := fs.String("file", DefaultTestsFile, "the timings file to rewrite")
	pkg := fs.String("pkg", DefaultTestsPkg, "the package whose top-level tests the file lists; tests no longer in it are dropped")
	skip := fs.String("skip", "", "regexp of tests left out of the file (the timing budgets)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if (*runID == "") == (*jsonFile == "") {
		return errors.New("merge-tests: pass exactly one of -run <id> or -json <file>")
	}
	var skipRE *regexp.Regexp
	if *skip != "" {
		var err error
		if skipRE, err = regexp.Compile(*skip); err != nil {
			return fmt.Errorf("merge-tests: -skip: %w", err)
		}
	}
	var fresh map[string]float64
	var err error
	if *runID != "" {
		var r io.Reader
		if r, err = runLog(*runID); err != nil {
			return err
		}
		fresh, err = ParseTestTimeLines(r)
	} else {
		var f *os.File
		if f, err = os.Open(*jsonFile); err != nil {
			return err
		}
		fresh, err = ParseTopTestJSON(f)
		f.Close()
	}
	if err != nil {
		return err
	}
	listed, err := listTests(*pkg)
	if err != nil {
		return err
	}
	var tests []string
	for _, t := range listed {
		if skipRE == nil || !skipRE.MatchString(t) {
			tests = append(tests, t)
		}
	}
	n, err := MergeTestsFile(*file, tests, fresh)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "scripttimings: %d fresh times, %d tests in %s\n", len(fresh), n, *file)
	return nil
}
