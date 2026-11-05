package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
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

// writeGen draws random writes over a growing pool of refs.
type writeGen struct {
	t       *testing.T
	rng     *rand.Rand
	dir     string
	handles []string
	sets    []id.ID
	n       int
	// base counts the handles at the start of handles every view has.
	base int
}

func (g *writeGen) pick() string {
	// Now and then a change set, which a link may name at either end.
	if len(g.sets) > 0 && g.rng.Intn(12) == 0 {
		return g.sets[g.rng.Intn(len(g.sets))].String()
	}
	// Mostly one of the landed handles, which every view resolves; else
	// any handle drawn so far, which may live in another change set only,
	// or nowhere.
	if g.rng.Intn(2) == 0 {
		return g.handles[g.rng.Intn(min(len(g.handles), g.base))]
	}
	return g.handles[g.rng.Intn(len(g.handles))]
}

// write returns a random write: its operation and request.
func (g *writeGen) write() (string, string) {
	g.n++
	types := []string{"serves", "depends-on", "relates-to", "decided-by"}
	kinds := []string{"capability", "capability", "capability", "decision", "intent"}
	switch g.rng.Intn(8) {
	case 0, 1:
		h := fmt.Sprintf("h%d", g.rng.Intn(40))
		g.handles = append(g.handles, h)
		kind := kinds[g.rng.Intn(len(kinds))]
		if kind != "capability" {
			return "add", fmt.Sprintf(`{"kind":%q,"handle":%q}`, kind, h)
		}
		return "add", fmt.Sprintf(`{"kind":"capability","handle":%q,"parent":%q}`, h, g.pick())
	case 2:
		return "set", fmt.Sprintf(`{"ref":%q,"fields":["statement=s%d"]}`, g.pick(), g.n)
	case 3:
		return "mv", fmt.Sprintf(`{"ref":%q,"parent":%q}`, g.pick(), g.pick())
	case 4:
		return "retire", fmt.Sprintf(`{"ref":%q,"reason":"r"}`, g.pick())
	case 5:
		return "link", fmt.Sprintf(`{"from":%q,"type":%q,"to":%q,"note":"n"}`, g.pick(), types[g.rng.Intn(len(types))], g.pick())
	case 6:
		return "unlink", fmt.Sprintf(`{"ref":%q,"type":%q,"to":%q}`, g.pick(), types[g.rng.Intn(len(types))], g.pick())
	}
	h := fmt.Sprintf("a%d", g.n)
	g.handles = append(g.handles, h)
	file := filepath.Join(g.dir, fmt.Sprintf("w%d.json", g.n))
	body := fmt.Sprintf(`{"writes":[{"op":"add","kind":"capability","handle":%q,"parent":%q},{"op":"set","ref":%q,"fields":["statement=x%d"]},{"op":"mv","ref":%q,"parent":%q},{"op":"link","from":%q,"type":"depends-on","to":%q}]}`,
		h, g.pick(), g.pick(), g.n, g.pick(), h, h, g.pick())
	if err := os.WriteFile(file, []byte(body), 0o644); err != nil {
		g.t.Fatal(err)
	}
	return "apply", fmt.Sprintf(`{"file":%q}`, file)
}

// genLog grows the log of ctx with random change sets: each takes random
// writes, now and then depends on an earlier one, and ends in a random
// state (draft, proposed, accepted, landed or withdrawn). Writes and moves
// the log refuses are skipped.
func genLog(t *testing.T, ctx *Ctx, g *writeGen) {
	t.Helper()
	for i, n := 0, 3+g.rng.Intn(3); i < n; i++ {
		cs := newChange(t, ctx, fmt.Sprintf(`{"intent":"Generated %d"}`, i)).ChangeSet
		for k, m := 0, 4+g.rng.Intn(8); k < m; k++ {
			name, req := g.write()
			callOp(t, ctx, name, req)
		}
		if len(g.sets) > 0 && g.rng.Intn(2) == 0 {
			callOp(t, ctx, "link", fmt.Sprintf(`{"from":%q,"type":"depends-on","to":%q}`, cs, g.sets[g.rng.Intn(len(g.sets))]))
		}
		g.sets = append(g.sets, cs)
	}
	for _, cs := range g.sets {
		ctx.ChangeSet = cs
		steps := [][2]string{{"propose", `{}`}, {"accept", `{"reason":"ok"}`}, {"land", `{}`}}
		switch r := g.rng.Intn(5); r {
		case 4:
			callOp(t, ctx, "withdraw", `{}`)
		default:
			for _, s := range steps[:r] {
				if _, err := callOp(t, ctx, s[0], s[1]); err != nil {
					break
				}
			}
		}
	}
	ctx.ChangeSet = id.ID{}
}

// TestWritesThroughIndexMatchFold: over generated logs and generated
// writes, each write is accepted or refused through the index, with the
// same error or the same ops, as through the fold of the log.
func TestWritesThroughIndexMatchFold(t *testing.T) {
	seeds := int64(4)
	if testing.Short() {
		seeds = 1
	}
	for seed := int64(1); seed <= seeds; seed++ {
		ctx, landed, accepted := readLog(t)
		folded := *ctx
		folded.IndexPath = ""
		rng := rand.New(rand.NewSource(seed))
		g := &writeGen{t: t, rng: rng, dir: t.TempDir(),
			handles: []string{"store", "cache", "disk", "tape", "old", "under-old", "faster", "use-ssd", "ram", "ghost"},
			sets:    []id.ID{landed, accepted}, base: 8}
		genLog(t, ctx, g)
		var live []id.ID
		states := map[changeset.State]int{}
		for _, cs := range g.sets {
			c := foldedState(t, ctx, cs)
			states[c.State]++
			if inFlightState(c.State) {
				live = append(live, cs)
			}
		}
		t.Logf("seed %d: change sets by state %v", seed, states)
		accepts, refusals := 0, 0
		for step := 0; step < 80; step++ {
			if step%20 == 0 {
				// Writes go mostly to an in-flight change set, now and
				// then to one in any state.
				ctx.ChangeSet = g.sets[rng.Intn(len(g.sets))]
				if len(live) > 0 && rng.Intn(4) > 0 {
					ctx.ChangeSet = live[rng.Intn(len(live))]
				}
				folded.ChangeSet = ctx.ChangeSet
			}
			name, req := g.write()
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
			if strings.HasPrefix(want, "error") {
				if testing.Verbose() {
					t.Logf("refused %s %s: %s", name, req, want)
				}
				refusals++
				continue
			}
			accepts++
			// The real write goes through the index or the fold by turns.
			real := &folded
			if step%2 == 1 {
				real = ctx
			}
			if r := answer(real, false); strings.HasPrefix(r, "error") {
				t.Fatalf("seed %d step %d: dry run accepted, real write failed: %s", seed, step, r)
			}
		}
		t.Logf("seed %d: %d accepted, %d refused", seed, accepts, refusals)
		if accepts < 10 || refusals < 5 {
			t.Errorf("seed %d: thin run: %d accepted, %d refused", seed, accepts, refusals)
		}
	}
}
