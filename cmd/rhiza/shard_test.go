package main

import (
	"fmt"
	"maps"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestScriptShard(t *testing.T) {
	all, err := filepath.Glob("testdata/script/*.txtar")
	if err != nil || len(all) == 0 {
		t.Fatalf("no scripts: %v", err)
	}
	times, err := readScriptTimings(scriptTimingsFile)
	if err != nil {
		t.Fatal(err)
	}

	// fixture writes scripts named after timings' keys plus untimed into a
	// temp dir, with timings as its timings file.
	fixture := func(t *testing.T, timings map[string]float64, untimed ...string) (dir, tsv string) {
		t.Helper()
		dir = t.TempDir()
		var body strings.Builder
		for name, secs := range timings {
			fmt.Fprintf(&body, "%s\t%v\n", name, secs)
		}
		for _, name := range append(slices.Sorted(maps.Keys(timings)), untimed...) {
			if err := os.WriteFile(filepath.Join(dir, name+".txtar"), nil, 0o644); err != nil {
				t.Fatal(err)
			}
		}
		tsv = filepath.Join(t.TempDir(), "script-timings.tsv")
		if err := os.WriteFile(tsv, []byte(body.String()), 0o644); err != nil {
			t.Fatal(err)
		}
		return dir, tsv
	}
	shards := func(t *testing.T, dir, tsv string, n int) [][]string {
		t.Helper()
		out := make([][]string, n)
		for i := 1; i <= n; i++ {
			got, err := scriptFiles(dir, fmt.Sprintf("%d/%d", i, n), tsv)
			if err != nil {
				t.Fatal(err)
			}
			for _, f := range got {
				out[i-1] = append(out[i-1], scriptName(f))
			}
		}
		return out
	}
	disjointCovering := func(t *testing.T, got [][]string, want []string) {
		t.Helper()
		seen := map[string]int{}
		for _, s := range got {
			for _, name := range s {
				seen[name]++
			}
		}
		for _, name := range want {
			if seen[name] != 1 {
				t.Errorf("%s is in %d shards, want 1", name, seen[name])
			}
		}
		if len(seen) != len(want) {
			t.Errorf("shards hold %d scripts, want %d", len(seen), len(want))
		}
	}

	t.Run("unset selects all, longest first", func(t *testing.T) {
		got, err := scriptFiles("testdata/script", "", scriptTimingsFile)
		if err != nil {
			t.Fatal(err)
		}
		if want := longestFirst(all, times); !slices.Equal(got, want) {
			t.Fatalf("unset = %d scripts, want all %d longest first", len(got), len(want))
		}
	})

	t.Run("greedy, ties to the lower shard", func(t *testing.T) {
		dir, tsv := fixture(t, map[string]float64{"a": 9, "b": 7, "c": 6, "d": 5, "e": 4, "f": 3, "g": 2, "h": 1})
		// a→1 (9), b→2 (7), c→3 (6), d→3 (11), e→2 (11), f→1 (12); g finds
		// 2 and 3 tied at 11 and takes 2 (13); h→3 (12).
		got := shards(t, dir, tsv, 3)
		want := [][]string{{"a", "f"}, {"b", "e", "g"}, {"c", "d", "h"}}
		if !slices.EqualFunc(got, want, slices.Equal[[]string]) {
			t.Fatalf("shards = %v, want %v", got, want)
		}
		disjointCovering(t, got, []string{"a", "b", "c", "d", "e", "f", "g", "h"})
	})

	t.Run("missing scripts weigh the median, run last by name", func(t *testing.T) {
		dir, tsv := fixture(t, map[string]float64{"a": 9, "b": 7, "c": 6, "d": 5, "e": 4, "f": 3, "g": 2, "h": 1}, "x", "y")
		// The median of the timed eight is 4.5, so x and y place after d (5)
		// and before e (4): a→1, b→2, c→3, d→3 (11), x→2 (11.5), y→1
		// (13.5), e→3 (15), f→2 (14.5), g→1 (15.5), h→2 (15.5). In a shard
		// they still run after every timed script.
		got := shards(t, dir, tsv, 3)
		want := [][]string{{"a", "g", "y"}, {"b", "f", "h", "x"}, {"c", "d", "e"}}
		if !slices.EqualFunc(got, want, slices.Equal[[]string]) {
			t.Fatalf("shards = %v, want %v", got, want)
		}
		disjointCovering(t, got, []string{"a", "b", "c", "d", "e", "f", "g", "h", "x", "y"})
	})

	t.Run("no timings spreads by name", func(t *testing.T) {
		dir, tsv := fixture(t, nil, "a", "b", "c", "d", "e")
		// Every script weighs 1, so the ties walk the shards in turn.
		got := shards(t, dir, tsv, 2)
		want := [][]string{{"a", "c", "e"}, {"b", "d"}}
		if !slices.EqualFunc(got, want, slices.Equal[[]string]) {
			t.Fatalf("shards = %v, want %v", got, want)
		}
	})

	t.Run("checked-in timings, four shards balanced", func(t *testing.T) {
		got := shards(t, "testdata/script", scriptTimingsFile, 4)
		names := make([]string, len(all))
		for i, f := range all {
			names[i] = scriptName(f)
		}
		disjointCovering(t, got, names)
		median := medianTime(all, times)
		totals := make([]float64, len(got))
		var sum float64
		for i, s := range got {
			for idx, name := range s {
				w, ok := times[name]
				if !ok {
					w = median
				}
				totals[i] += w
				// Each shard runs in scripts-longest-first's order.
				if idx > 0 {
					prev := s[idx-1]
					tp, okP := times[prev]
					if (!okP && ok) || (okP == ok && (tp < times[name] || tp == times[name] && prev > name)) {
						t.Errorf("shard %d: %s before %s", i+1, prev, name)
					}
				}
			}
			sum += totals[i]
		}
		mean := sum / float64(len(totals))
		for i, tot := range totals {
			if math.Abs(tot-mean) > 0.15*mean {
				t.Errorf("shard %d/4 totals %.1fs, more than 15%% off the mean %.1fs (%v)", i+1, tot, mean, totals)
			}
		}
	})

	t.Run("one shard is everything", func(t *testing.T) {
		got, err := shardNames(all, times, "1/1")
		if err != nil || !slices.Equal(got, longestFirst(all, times)) {
			t.Fatalf("1/1 = %v, %v", got, err)
		}
	})

	for _, bad := range []string{"0/4", "5/4", "x", "1/0", "1/", "/4", "1/4x", "-1/4"} {
		t.Run("malformed "+bad, func(t *testing.T) {
			if _, err := shardNames(all, times, bad); err == nil {
				t.Fatalf("%q accepted", bad)
			}
		})
	}
}

func TestScriptShardEnvError(t *testing.T) {
	if _, err := os.Stat("testdata/script"); err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{"0/4", "5/4", "1/0"} {
		_, err := scriptFiles("testdata/script", bad, scriptTimingsFile)
		if err == nil {
			t.Fatalf("RHIZA_TEST_SHARD=%s accepted", bad)
		}
		if !strings.Contains(err.Error(), "RHIZA_TEST_SHARD") {
			t.Fatalf("%s: error %q does not name RHIZA_TEST_SHARD", bad, err)
		}
	}
}
