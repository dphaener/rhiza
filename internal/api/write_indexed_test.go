package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/dphaener/rhiza/internal/id"
)

// seededIDs gives ctx a generator that draws the same IDs on every call
// with the same n.
func seededIDs(ctx *Ctx, n int64) {
	ctx.IDs = id.NewGenerator(ctx.Now, rand.New(rand.NewSource(n)))
}

// TestApplyIndexedReadsNoLog: with a fresh index an apply of several writes
// (and the single-write commands) reads the log zero times, and what it
// writes is what the log-folding write writes.
func TestApplyIndexedReadsNoLog(t *testing.T) {
	ctx, _ := indexedLog(t)
	// Warm the index.
	if _, err := callOp(t, ctx, "change.show", `{}`); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "writes.json")
	body := `{"writes":[
		{"op":"add","kind":"capability","handle":"zip","parent":"store"},
		{"op":"set","ref":"zip","fields":["statement=Zip it"]},
		{"op":"add","kind":"capability","handle":"zap","parent":"zip"},
		{"op":"link","from":"zap","type":"depends-on","to":"flash"},
		{"op":"mv","ref":"zap","parent":"store"},
		{"op":"retire","ref":"zip","reason":"not needed"}
	]}`
	if err := os.WriteFile(file, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	reads := countingCtx(ctx)
	if _, err := callOp(t, ctx, "apply", fmt.Sprintf(`{"file":%q}`, file)); err != nil {
		t.Fatal(err)
	}
	if *reads != 0 {
		t.Errorf("apply read the log %d times", *reads)
	}
	for _, c := range []struct{ name, req string }{
		{"add", `{"kind":"capability","handle":"zig","parent":"store"}`},
		{"set", `{"ref":"zig","fields":["statement=Zig"]}`},
		{"link", `{"from":"zig","type":"depends-on","to":"zap"}`},
		{"mv", `{"ref":"zig","parent":"flash"}`},
		{"unlink", `{"ref":"zig","type":"depends-on","to":"zap"}`},
		{"retire", `{"ref":"zig","reason":"gone"}`},
	} {
		*reads = 0
		if _, err := callOp(t, ctx, c.name, c.req); err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if *reads != 0 {
			t.Errorf("%s read the log %d times", c.name, *reads)
		}
	}
}

// TestWritesThroughIndexMatchFold: over generated writes, each accepted or
// refused through the index, with the same error or the same ops, as
// through the fold of the log.
func TestWritesThroughIndexMatchFold(t *testing.T) {
	for seed := int64(1); seed <= 3; seed++ {
		ctx, _ := indexedLog(t)
		folded := *ctx
		folded.IndexPath = ""
		rng := rand.New(rand.NewSource(seed))
		handles := []string{"store", "cache", "disk", "tape", "old", "under-old", "flash", "ssd", "nvme", "faster", "use-ssd", "ghost"}
		pick := func() string { return handles[rng.Intn(len(handles))] }
		stable := []string{"store", "cache", "disk", "tape", "old"}
		types := []string{"serves", "depends-on", "relates-to", "decided-by"}
		dir := t.TempDir()
		accepted, refused := 0, 0
		for step := 0; step < 60; step++ {
			var name, req string
			switch rng.Intn(7) {
			case 0:
				name, req = "add", fmt.Sprintf(`{"kind":"capability","handle":"n%d-%d","parent":%q}`, seed, rng.Intn(20), pick())
			case 1:
				name, req = "set", fmt.Sprintf(`{"ref":%q,"fields":["statement=s%d"]}`, pick(), step)
			case 2:
				// The parent is one that exists before the change set: the
				// what-if view drops a node moved under one the change set
				// creates after it, which the interim model keeps.
				name, req = "mv", fmt.Sprintf(`{"ref":%q,"parent":%q}`, pick(), stable[rng.Intn(len(stable))])
			case 3:
				// Not a node the change set creates and a child hangs under: the
				// what-if view drops the child, which the interim model keeps.
				retirable := append([]string{"nvme", "use-ssd", "faster", "ghost"}, stable...)
				name, req = "retire", fmt.Sprintf(`{"ref":%q,"reason":"r"}`, retirable[rng.Intn(len(retirable))])
			case 4:
				name, req = "link", fmt.Sprintf(`{"from":%q,"type":%q,"to":%q,"note":"n"}`, pick(), types[rng.Intn(len(types))], pick())
			case 5:
				name, req = "unlink", fmt.Sprintf(`{"ref":%q,"type":%q,"to":%q}`, pick(), types[rng.Intn(len(types))], pick())
			default:
				file := filepath.Join(dir, fmt.Sprintf("w%d.json", step))
				body := fmt.Sprintf(`{"writes":[{"op":"add","kind":"capability","handle":"a%d-%d","parent":%q},{"op":"set","ref":%q,"fields":["statement=x%d"]},{"op":"link","from":%q,"type":"depends-on","to":%q}]}`,
					seed, step, pick(), pick(), step, pick(), pick())
				if err := os.WriteFile(file, []byte(body), 0o644); err != nil {
					t.Fatal(err)
				}
				name, req = "apply", fmt.Sprintf(`{"file":%q}`, file)
			}
			answer := func(c *Ctx, dry bool) string {
				run := *c
				run.DryRun = dry
				seededIDs(&run, int64(step)+100)
				resp, err := callOp(t, &run, name, req)
				if err != nil {
					return fmt.Sprintf("error %d: %v", ExitCode(err), err)
				}
				b, _ := json.Marshal(resp)
				return string(b)
			}
			want, got := answer(&folded, true), answer(ctx, true)
			if got != want {
				t.Fatalf("seed %d step %d: %s %s\nindexed %s\nfolded  %s", seed, step, name, req, got, want)
			}
			if want[:5] == "error" {
				refused++
				continue
			}
			accepted++
			real := answer(&folded, false)
			if real[:5] == "error" {
				t.Fatalf("seed %d step %d: dry run accepted, real write failed: %s", seed, step, real)
			}
		}
		t.Logf("seed %d: %d accepted, %d refused", seed, accepted, refused)
		if accepted < 10 || refused < 5 {
			t.Errorf("seed %d: thin run: %d accepted, %d refused", seed, accepted, refused)
		}
	}
}
