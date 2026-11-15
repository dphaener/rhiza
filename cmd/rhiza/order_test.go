package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// scriptTimingsFile holds each script's last CI time, name<TAB>seconds per
// line (make script-timings rewrites it).
const scriptTimingsFile = "testdata/script-timings.tsv"

// scriptFiles lists the scripts in dir that shard (RHIZA_TEST_SHARD, "i/n")
// selects, every script when shard is empty, longest first by the times in
// timingsPath; shardNames splits them by those times. It orders the list
// testscript is handed, not the start: testscript runs each script as a
// parallel subtest, and the -parallel slots go to the woken subtests in no
// fixed order.
func scriptFiles(dir, shard, timingsPath string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.txtar"))
	if err != nil {
		return nil, err
	}
	times, err := readScriptTimings(timingsPath)
	if err != nil {
		return nil, err
	}
	if shard != "" {
		return shardNames(files, times, shard)
	}
	return longestFirst(files, times), nil
}

// scriptName is a script's name as testscript and the timings file give it:
// the file's base name without .txtar.
func scriptName(path string) string {
	return strings.TrimSuffix(filepath.Base(path), ".txtar")
}

// longestFirst returns files sorted by descending seconds in times. Ties, and
// the scripts times does not list (after every timed one), go by name.
func longestFirst(files []string, times map[string]float64) []string {
	out := slices.Clone(files)
	sort.SliceStable(out, func(i, j int) bool {
		a, b := scriptName(out[i]), scriptName(out[j])
		ta, okA := times[a]
		tb, okB := times[b]
		if okA != okB {
			return okA
		}
		if ta != tb {
			return ta > tb
		}
		return a < b
	})
	return out
}

// readScriptTimings reads the timings file into script name -> seconds. A
// missing file reads as no times, so every script then goes by name.
func readScriptTimings(path string) (map[string]float64, error) {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]float64{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseScriptTimings(f, path)
}

// parseScriptTimings reads name<TAB>seconds lines; src names the input in
// errors, which carry the line number.
func parseScriptTimings(r io.Reader, src string) (map[string]float64, error) {
	times := map[string]float64{}
	sc := bufio.NewScanner(r)
	for n := 1; sc.Scan(); n++ {
		line := sc.Text()
		name, secs, ok := strings.Cut(line, "\t")
		v, err := strconv.ParseFloat(secs, 64)
		if !ok || name == "" || err != nil || v < 0 || math.IsInf(v, 0) || math.IsNaN(v) {
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

func TestScriptOrder(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a", "b", "c", "d", "e", "f"} {
		if err := os.WriteFile(filepath.Join(dir, name+".txtar"), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// Not a script: testscript's Dir globbing would skip it, and so must we.
	if err := os.WriteFile(filepath.Join(dir, "notes.md"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	write := func(t *testing.T, body string) string {
		t.Helper()
		p := filepath.Join(t.TempDir(), "script-timings.tsv")
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		return p
	}
	names := func(files []string) []string {
		out := make([]string, len(files))
		for i, f := range files {
			out[i] = scriptName(f)
		}
		return out
	}

	t.Run("longest first, ties and untimed by name", func(t *testing.T) {
		// e and c tie; a and f are missing; gone is no longer a script.
		p := write(t, "b\t2\nc\t7.5\nd\t30\ne\t7.5\ngone\t99\n")
		got, err := scriptFiles(dir, "", p)
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"d", "c", "e", "b", "a", "f"}
		if !slices.Equal(names(got), want) {
			t.Fatalf("order = %v, want %v", names(got), want)
		}
	})

	t.Run("every script exactly once", func(t *testing.T) {
		all, err := filepath.Glob("testdata/script/*.txtar")
		if err != nil || len(all) == 0 {
			t.Fatalf("no scripts: %v", err)
		}
		got, err := scriptFiles("testdata/script", "", scriptTimingsFile)
		if err != nil {
			t.Fatal(err)
		}
		sorted := slices.Clone(got)
		sort.Strings(sorted)
		sort.Strings(all)
		if !slices.Equal(sorted, all) {
			t.Fatalf("ordered list has %d entries, want each of the %d scripts once", len(got), len(all))
		}
		times, err := readScriptTimings(scriptTimingsFile)
		if err != nil {
			t.Fatal(err)
		}
		for i := 1; i < len(got); i++ {
			a, b := scriptName(got[i-1]), scriptName(got[i])
			ta, okA := times[a]
			tb, okB := times[b]
			if (!okA && okB) || (okA == okB && (ta < tb || ta == tb && a > b)) {
				t.Fatalf("%s (%v) before %s (%v)", a, ta, b, tb)
			}
		}
	})

	t.Run("a shard is ordered too", func(t *testing.T) {
		p := write(t, "a\t1\nc\t3\ne\t5\n")
		got, err := scriptFiles(dir, "1/2", p)
		if err != nil {
			t.Fatal(err)
		}
		// b, d, f weigh the median, 3. Placed e(5)→1, c(3)→2, b→2, d→1,
		// f→2, a(1)→1: shard 1 is e, d, a, run timed first, then d.
		if want := []string{"e", "a", "d"}; !slices.Equal(names(got), want) {
			t.Fatalf("shard 1/2 = %v, want %v", names(got), want)
		}
	})

	t.Run("missing timings file orders by name", func(t *testing.T) {
		got, err := scriptFiles(dir, "", filepath.Join(t.TempDir(), "none.tsv"))
		if err != nil {
			t.Fatal(err)
		}
		if want := []string{"a", "b", "c", "d", "e", "f"}; !slices.Equal(names(got), want) {
			t.Fatalf("order = %v, want %v", names(got), want)
		}
	})

	for _, tc := range []struct{ name, body, line string }{
		{"no tab", "a\t1\nb 2\n", ":2:"},
		{"seconds not a number", "a\t1\nb\t2\nc\tslow\n", ":3:"},
		{"negative seconds", "a\t-1\n", ":1:"},
		{"empty name", "\t1\n", ":1:"},
		{"listed twice", "a\t1\na\t2\n", ":2:"},
	} {
		t.Run("malformed "+tc.name, func(t *testing.T) {
			p := write(t, tc.body)
			_, err := scriptFiles(dir, "", p)
			if err == nil {
				t.Fatalf("%q accepted", tc.body)
			}
			if !strings.Contains(err.Error(), p+tc.line) {
				t.Fatalf("error %q does not name %s%s", err, p, tc.line)
			}
		})
	}
}
