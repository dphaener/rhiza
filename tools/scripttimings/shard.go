package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"
)

// ParseShard reads a shard as "i/n" (1-based, 1 <= i <= n).
func ParseShard(shard string) (i, n int, err error) {
	var rest string
	k, _ := fmt.Sscanf(shard, "%d/%d%s", &i, &n, &rest)
	if k != 2 || n < 1 || i < 1 || i > n {
		return 0, 0, fmt.Errorf("SHARD=%s: want i/n with 1 <= i <= n", shard)
	}
	return i, n, nil
}

// Shard returns the names shard i of n (1-based) takes, split by the seconds
// in times so the n shards take about as long: each name, longest first, goes
// to the shard with the least total so far, ties to the lower shard. A name
// times does not list weighs the median of the timed names' seconds (1 when
// none is timed), for placement only. On equal weight a timed name goes first,
// then by name, so the placement is the same on every machine. The shards are
// disjoint and together cover every name; shard i comes back in placement
// order. This is cmd/rhiza's shardNames for the scripts shards, over test
// names.
func Shard(names []string, times map[string]float64, i, n int) []string {
	median := medianOf(names, times)
	weight := func(name string) (float64, bool) {
		if s, ok := times[name]; ok {
			return s, true
		}
		return median, false
	}
	order := append([]string(nil), names...)
	sort.SliceStable(order, func(a, b int) bool {
		wa, okA := weight(order[a])
		wb, okB := weight(order[b])
		if wa != wb {
			return wa > wb
		}
		if okA != okB {
			return okA
		}
		return order[a] < order[b]
	})
	totals := make([]float64, n)
	out := []string{}
	for _, name := range order {
		least := 0
		for s := 1; s < n; s++ {
			if totals[s] < totals[least] {
				least = s
			}
		}
		w, _ := weight(name)
		totals[least] += w
		if least == i-1 {
			out = append(out, name)
		}
	}
	return out
}

// medianOf is the median of the seconds times gives names, or 1 when none of
// them is timed.
func medianOf(names []string, times map[string]float64) float64 {
	vs := []float64{}
	for _, name := range names {
		if v, ok := times[name]; ok {
			vs = append(vs, v)
		}
	}
	if len(vs) == 0 {
		return 1
	}
	sort.Float64s(vs)
	m := len(vs) / 2
	if len(vs)%2 == 1 {
		return vs[m]
	}
	return (vs[m-1] + vs[m]) / 2
}

// shard reads test names on stdin, one per line (blank lines ignored, a name
// twice counted once), and prints the names shard -shard takes, one per line.
func shard(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("shard", flag.ContinueOnError)
	fs.SetOutput(stderr)
	spec := fs.String("shard", "", "the shard to print, i/n (1-based)")
	file := fs.String("file", DefaultTestsFile, "the timings file the tests are weighed by")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		return errors.New("shard: takes no positional arguments; pipe the test names into it")
	}
	i, n, err := ParseShard(*spec)
	if err != nil {
		return err
	}
	times, err := ReadFile(*file)
	if err != nil {
		return err
	}
	var names []string
	seen := map[string]bool{}
	sc := bufio.NewScanner(stdin)
	for sc.Scan() {
		name := strings.TrimSpace(sc.Text())
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		names = append(names, name)
	}
	if err := sc.Err(); err != nil {
		return err
	}
	for _, name := range Shard(names, times, i, n) {
		fmt.Fprintln(stdout, name)
	}
	return nil
}
