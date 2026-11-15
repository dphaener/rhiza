package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Paths from the repository root, where make runs the tool.
const (
	DefaultFile    = "cmd/rhiza/testdata/script-timings.tsv"
	DefaultScripts = "cmd/rhiza/testdata/script"
	// DefaultTestsFile and DefaultTestsPkg are merge-tests' file and the
	// package whose top-level tests it lists.
	DefaultTestsFile = "internal/api/testdata/test-timings.tsv"
	DefaultTestsPkg  = "./internal/api"
)

// ScriptTimePrefix opens a script-time line, script-time<TAB><name><TAB><seconds>.
const ScriptTimePrefix = "script-time\t"

// TestTimePrefix opens a test-time line, test-time<TAB><test><TAB><seconds>:
// the time of one passing top-level test.
const TestTimePrefix = "test-time\t"

// NewScriptSeconds is the time given to a script that has neither a fresh nor
// a previous time: the same 1 second every script was seeded with.
const NewScriptSeconds = 1.0

// ErrNoTimes is returned when the input holds no passing TestScripts subtest
// and no script-time line: nothing to merge, so the file is left alone.
var ErrNoTimes = errors.New("no passing TestScripts subtest and no script-time line in the input")

// ErrNoTestTimes is ErrNoTimes for top-level tests: the input holds no passing
// top-level test and no test-time line.
var ErrNoTestTimes = errors.New("no passing top-level test and no test-time line in the input")

// FormatScriptTime renders one script-time line (no newline).
func FormatScriptTime(name string, seconds float64) string {
	return ScriptTimePrefix + name + "\t" + formatSeconds(seconds)
}

// FormatTestTime renders one test-time line (no newline).
func FormatTestTime(name string, seconds float64) string {
	return TestTimePrefix + name + "\t" + formatSeconds(seconds)
}

// ParseScriptTimeLines reads script-time lines out of a log. In a gh run view
// --log dump each line carries a job, a step and a timestamp before it, so the
// script-time field is found anywhere in the line rather than at its start.
// Other lines are ignored; a script timed twice keeps its longest time.
func ParseScriptTimeLines(r io.Reader) (map[string]float64, error) {
	return parseTimeLines(r, ScriptTimePrefix, ErrNoTimes)
}

// ParseTestTimeLines reads test-time lines out of a log, as
// ParseScriptTimeLines reads script-time lines.
func ParseTestTimeLines(r io.Reader) (map[string]float64, error) {
	return parseTimeLines(r, TestTimePrefix, ErrNoTestTimes)
}

// parseTimeLines reads the lines carrying prefix, <name><TAB><seconds> after
// it, returning none when there are none.
func parseTimeLines(r io.Reader, prefix string, none error) (map[string]float64, error) {
	kind := strings.TrimSuffix(prefix, "\t")
	times := map[string]float64{}
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for sc.Scan() {
		line := strings.TrimRight(sc.Text(), "\r")
		i := strings.Index(line, prefix)
		if i < 0 {
			continue
		}
		fields := strings.Split(line[i+len(prefix):], "\t")
		if len(fields) != 2 {
			return nil, fmt.Errorf("malformed %s line %q", kind, line)
		}
		name := fields[0]
		secs, err := parseSeconds(strings.TrimSpace(fields[1]))
		if err != nil || !validName(name) {
			return nil, fmt.Errorf("malformed %s line %q", kind, line)
		}
		times[name] = math.Max(times[name], secs)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(times) == 0 {
		return nil, none
	}
	return times, nil
}

// testEvent is the part of a go test -json event the tool reads.
type testEvent struct {
	Action  string
	Test    string
	Elapsed float64
}

// ParseTestJSON reads the passing TestScripts/<name> subtests of a go test
// -json stream with their elapsed seconds. Failed, skipped and unfinished
// scripts are left out, so they keep their previous time. Lines that are not
// JSON (a build error, say) are ignored.
func ParseTestJSON(r io.Reader) (map[string]float64, error) {
	times := map[string]float64{}
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for sc.Scan() {
		line := bytes.TrimSpace(sc.Bytes())
		if len(line) == 0 || line[0] != '{' {
			continue
		}
		var ev testEvent
		if json.Unmarshal(line, &ev) != nil || ev.Action != "pass" {
			continue
		}
		name, ok := strings.CutPrefix(ev.Test, "TestScripts/")
		if !ok || !validName(name) {
			continue
		}
		times[name] = math.Max(times[name], ev.Elapsed)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(times) == 0 {
		return nil, ErrNoTimes
	}
	return times, nil
}

// ParseTopTestJSON reads the passing top-level tests of a go test -json
// stream with their elapsed seconds, as ParseTestJSON reads the scripts:
// subtests, failed, skipped and unfinished tests are left out.
func ParseTopTestJSON(r io.Reader) (map[string]float64, error) {
	times := map[string]float64{}
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for sc.Scan() {
		line := bytes.TrimSpace(sc.Bytes())
		if len(line) == 0 || line[0] != '{' {
			continue
		}
		var ev testEvent
		if json.Unmarshal(line, &ev) != nil || ev.Action != "pass" || !validName(ev.Test) {
			continue
		}
		times[ev.Test] = math.Max(times[ev.Test], ev.Elapsed)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(times) == 0 {
		return nil, ErrNoTestTimes
	}
	return times, nil
}

// ReadFile parses a timings file. A missing file is an empty one; a line that
// is not name<TAB>seconds, a bad time or a name twice is an error.
func ReadFile(path string) (map[string]float64, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]float64{}, nil
	}
	if err != nil {
		return nil, err
	}
	return Parse(bytes.NewReader(b), path)
}

// Parse reads name<TAB>seconds lines; src names the input in errors.
func Parse(r io.Reader, src string) (map[string]float64, error) {
	times := map[string]float64{}
	sc := bufio.NewScanner(r)
	for n := 1; sc.Scan(); n++ {
		line := sc.Text()
		name, secs, ok := strings.Cut(line, "\t")
		v, err := parseSeconds(secs)
		if !ok || err != nil || !validName(name) {
			return nil, fmt.Errorf("%s:%d: malformed line %q (want name<TAB>seconds)", src, n, line)
		}
		if _, dup := times[name]; dup {
			return nil, fmt.Errorf("%s:%d: %s listed twice", src, n, name)
		}
		times[name] = v
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return times, nil
}

// Merge gives each script its fresh time, else its previous one, else
// NewScriptSeconds. Names not in scripts are dropped.
func Merge(prev, fresh map[string]float64, scripts []string) map[string]float64 {
	out := make(map[string]float64, len(scripts))
	for _, s := range scripts {
		if v, ok := fresh[s]; ok {
			out[s] = v
		} else if v, ok := prev[s]; ok {
			out[s] = v
		} else {
			out[s] = NewScriptSeconds
		}
	}
	return out
}

// MergeTests gives each test its fresh time, else its previous one. A test
// with neither is left out (the shards weigh a missing test at the median),
// and names not in tests are dropped.
func MergeTests(prev, fresh map[string]float64, tests []string) map[string]float64 {
	out := make(map[string]float64, len(tests))
	for _, s := range tests {
		if v, ok := fresh[s]; ok {
			out[s] = v
		} else if v, ok := prev[s]; ok {
			out[s] = v
		}
	}
	return out
}

// Format renders a timings file: name<TAB>seconds lines sorted by name.
func Format(times map[string]float64) []byte {
	names := make([]string, 0, len(times))
	for n := range times {
		names = append(names, n)
	}
	sort.Strings(names)
	var b bytes.Buffer
	for _, n := range names {
		fmt.Fprintf(&b, "%s\t%s\n", n, formatSeconds(times[n]))
	}
	return b.Bytes()
}

// Scripts lists the testscripts in dir by name (the file name without .txtar).
func Scripts(dir string) ([]string, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*.txtar"))
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no *.txtar scripts in %s", dir)
	}
	names := make([]string, len(paths))
	for i, p := range paths {
		names[i] = strings.TrimSuffix(filepath.Base(p), ".txtar")
	}
	return names, nil
}

// MergeFile merges fresh times into the file at path for the scripts in dir
// and writes it back atomically, returning how many scripts it lists. On any
// error the file is left unchanged.
func MergeFile(path, dir string, fresh map[string]float64) (int, error) {
	if len(fresh) == 0 {
		return 0, ErrNoTimes
	}
	prev, err := ReadFile(path)
	if err != nil {
		return 0, err
	}
	scripts, err := Scripts(dir)
	if err != nil {
		return 0, err
	}
	merged := Merge(prev, fresh, scripts)
	return len(merged), writeFile(path, Format(merged))
}

// MergeTestsFile merges fresh times into the file at path for tests (the
// package's top-level tests) and writes it back atomically, returning how
// many tests it lists. On any error the file is left unchanged.
func MergeTestsFile(path string, tests []string, fresh map[string]float64) (int, error) {
	if len(fresh) == 0 {
		return 0, ErrNoTestTimes
	}
	if len(tests) == 0 {
		return 0, errors.New("no tests to time")
	}
	prev, err := ReadFile(path)
	if err != nil {
		return 0, err
	}
	merged := MergeTests(prev, fresh, tests)
	return len(merged), writeFile(path, Format(merged))
}

// writeFile replaces path with data through a temporary file beside it.
func writeFile(path string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".script-timings-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

func parseSeconds(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v < 0 || math.IsInf(v, 0) || math.IsNaN(v) {
		return 0, fmt.Errorf("bad seconds %q", s)
	}
	return v, nil
}

// formatSeconds keeps two decimals, without trailing zeros: 1, 0.5, 12.34.
func formatSeconds(v float64) string {
	return strconv.FormatFloat(math.Round(v*100)/100, 'f', -1, 64)
}

// validName rejects names that would break a line: empty, or with a tab,
// newline or slash (a nested subtest).
func validName(name string) bool {
	return name != "" && !strings.ContainsAny(name, "\t\n\r/")
}
