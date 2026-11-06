package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"testing"
)

func TestScriptShard(t *testing.T) {
	all, err := filepath.Glob("testdata/script/*.txtar")
	if err != nil || len(all) == 0 {
		t.Fatalf("no scripts: %v", err)
	}
	sort.Strings(all)

	t.Run("unset selects all", func(t *testing.T) {
		got, err := shardFiles("testdata/script", "")
		if err != nil || got != nil {
			t.Fatalf("got %v, %v; want nil (testscript runs every script)", got, err)
		}
	})

	t.Run("four shards", func(t *testing.T) {
		seen := map[string]int{}
		for i := 1; i <= 4; i++ {
			got, err := shardFiles("testdata/script", fmt.Sprintf("%d/4", i))
			if err != nil {
				t.Fatal(err)
			}
			for idx, name := range got {
				seen[name]++
				want := all[(i-1)+4*idx]
				if name != want {
					t.Fatalf("shard %d/4 entry %d = %s, want %s", i, idx, name, want)
				}
			}
		}
		if len(seen) != len(all) {
			t.Fatalf("shards cover %d scripts, want %d", len(seen), len(all))
		}
		for name, c := range seen {
			if c != 1 {
				t.Errorf("%s is in %d shards", name, c)
			}
		}
	})

	t.Run("one shard is everything", func(t *testing.T) {
		got, err := shardNames(all, "1/1")
		if err != nil || !slices.Equal(got, all) {
			t.Fatalf("1/1 = %v, %v", got, err)
		}
	})

	for _, bad := range []string{"0/4", "5/4", "x", "1/0", "1/", "/4", "1/4x", "-1/4"} {
		t.Run("malformed "+bad, func(t *testing.T) {
			if _, err := shardNames(all, bad); err == nil {
				t.Fatalf("%q accepted", bad)
			}
		})
	}
}

func TestScriptShardEnvError(t *testing.T) {
	if _, err := os.Stat("testdata/script"); err != nil {
		t.Fatal(err)
	}
	if _, err := shardFiles("testdata/script", "0/4"); err == nil {
		t.Fatal("0/4 accepted")
	}
}
