package main

import (
	"bytes"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestShard(t *testing.T) {
	t.Run("greedy longest first, a missing test at the median", func(t *testing.T) {
		// Median of the timed 10, 7, 6, 3 is 6.5, so TestE (untimed) goes
		// third: A to 1 (10/0), B to 2 (10/7), E to 2 (10/13.5), C to 1
		// (16/13.5), D to 2 (16/16.5).
		times := map[string]float64{"TestA": 10, "TestB": 7, "TestC": 6, "TestD": 3}
		names := []string{"TestA", "TestB", "TestC", "TestD", "TestE"}
		if got, want := Shard(names, times, 1, 2), []string{"TestA", "TestC"}; !reflect.DeepEqual(got, want) {
			t.Errorf("shard 1/2 = %v, want %v", got, want)
		}
		if got, want := Shard(names, times, 2, 2), []string{"TestB", "TestE", "TestD"}; !reflect.DeepEqual(got, want) {
			t.Errorf("shard 2/2 = %v, want %v", got, want)
		}
	})

	t.Run("ties to the lower shard", func(t *testing.T) {
		// P to 1, Q to 2 (2/2, a tie), R to 1, S to 2.
		times := map[string]float64{"TestP": 2, "TestQ": 2, "TestR": 1, "TestS": 1}
		names := []string{"TestS", "TestR", "TestQ", "TestP"}
		if got, want := Shard(names, times, 1, 2), []string{"TestP", "TestR"}; !reflect.DeepEqual(got, want) {
			t.Errorf("shard 1/2 = %v, want %v", got, want)
		}
		// Nothing timed: every test weighs 1, by name, alternating.
		if got, want := Shard([]string{"TestZ", "TestY", "TestX"}, nil, 1, 2), []string{"TestX", "TestZ"}; !reflect.DeepEqual(got, want) {
			t.Errorf("untimed shard 1/2 = %v, want %v", got, want)
		}
	})

	t.Run("disjoint and covering", func(t *testing.T) {
		times := map[string]float64{"TestA": 10, "TestB": 7, "TestC": 6, "TestD": 3, "TestF": 0}
		names := []string{"TestA", "TestB", "TestC", "TestD", "TestE", "TestF", "TestG"}
		for n := 1; n <= 4; n++ {
			seen := map[string]int{}
			for i := 1; i <= n; i++ {
				for _, name := range Shard(names, times, i, n) {
					seen[name]++
				}
			}
			for _, name := range names {
				if seen[name] != 1 {
					t.Errorf("n=%d: %s is in %d shards, want 1", n, name, seen[name])
				}
			}
			if len(seen) != len(names) {
				t.Errorf("n=%d: shards hold %d names, want %d", n, len(seen), len(names))
			}
		}
	})
}

// With the checked-in file, the two shards' totals are within 15% of their
// mean.
func TestShardCheckedInFileBalanced(t *testing.T) {
	times, err := ReadFile(filepath.Join("..", "..", DefaultTestsFile))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for name := range times {
		names = append(names, name)
	}
	sort.Strings(names)
	var totals [2]float64
	for i := range totals {
		for _, name := range Shard(names, times, i+1, 2) {
			totals[i] += times[name]
		}
	}
	mean := (totals[0] + totals[1]) / 2
	for i, total := range totals {
		if math.Abs(total-mean) > 0.15*mean {
			t.Errorf("shard %d/2 takes %.1fs, more than 15%% off the mean %.1fs (%v)", i+1, total, mean, totals)
		}
	}
}

func TestShardMode(t *testing.T) {
	file := filepath.Join(t.TempDir(), "times.tsv")
	if err := os.WriteFile(file, []byte("TestA\t10\nTestB\t7\nTestC\t6\nTestD\t3\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	in := "TestA\nTestB\n\nTestC\nTestD\nTestE\nTestA\n"
	for shard, want := range map[string]string{"1/2": "TestA\nTestC\n", "2/2": "TestB\nTestE\nTestD\n"} {
		var out, errOut bytes.Buffer
		if err := run([]string{"shard", "-shard", shard, "-file", file}, strings.NewReader(in), &out, &errOut, nil, nil); err != nil {
			t.Fatalf("shard %s: %v\n%s", shard, err, errOut.String())
		}
		if out.String() != want {
			t.Errorf("shard %s printed %q, want %q", shard, out.String(), want)
		}
	}

	// A missing file weighs every test 1.
	var out bytes.Buffer
	if err := run([]string{"shard", "-shard", "2/2", "-file", file + ".none"}, strings.NewReader("TestB\nTestA\n"), &out, &bytes.Buffer{}, nil, nil); err != nil || out.String() != "TestB\n" {
		t.Errorf("no file: got %q, %v; want TestB", out.String(), err)
	}

	for _, bad := range []string{"", "0/2", "3/2", "1/0", "1/2x", "a/b", "1"} {
		var out bytes.Buffer
		err := run([]string{"shard", "-shard", bad, "-file", file}, strings.NewReader("TestA\n"), &out, &bytes.Buffer{}, nil, nil)
		if err == nil || !strings.Contains(err.Error(), "want i/n with 1 <= i <= n") {
			t.Errorf("shard %q: err = %v, want the i/n error", bad, err)
		}
		if out.Len() != 0 {
			t.Errorf("shard %q printed %q", bad, out.String())
		}
	}
	if err := run([]string{"shard", "-shard", "1/2", "-file", file, "extra"}, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{}, nil, nil); err == nil {
		t.Error("shard with a positional argument accepted")
	}
}
