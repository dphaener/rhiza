package api

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func stabilityRuns(t *testing.T, f *sessionFixture) []CIRun {
	t.Helper()
	rec, err := readCIStability(filepath.Join(f.common, filepath.FromSlash(ciStabilityFile)))
	if err != nil {
		t.Fatal(err)
	}
	return rec.Runs
}

// ship check records each settled CI run: a clean green, a flake per
// failing job and step of a red head whose rerun passed, a red with no flake, and
// main's result, from the tested head's or from the dispatched run.
func TestShipCheckStabilityRecord(t *testing.T) {
	const logCall = "run view 42 --log-failed"
	const testLog = "test\tRun\t--- FAIL: TestGarden (0.1s)\n--- FAIL: TestPond (0.1s)\n"

	t.Run("green first attempt", func(t *testing.T) {
		f, g, _ := shipFixture(t)
		g.answers[shipPRView] = openPR(greenRollup)
		shipCheck(t, f, `{}`)
		shipCheck(t, f, `{}`)
		runs := stabilityRuns(t, f)
		if len(runs) != 1 || runs[0].Head != "abc1234def" || runs[0].PR != 7 || runs[0].Main || runs[0].Result != ciGreen || len(runs[0].Flakes) != 0 || runs[0].Time == "" {
			t.Errorf("runs = %+v", runs)
		}
	})

	t.Run("red then green on rerun", func(t *testing.T) {
		f, g, _ := shipFixture(t)
		g.answers[shipPRView] = openPR(rerunRollup)
		g.answers["run rerun 42 --failed"] = ""
		g.answers[logCall] = testLog
		if r := shipCheck(t, f, `{"settled":true}`); r.Settled {
			t.Fatalf("rerun: %+v", r)
		}
		if runs := stabilityRuns(t, f); len(runs) != 0 {
			t.Fatalf("recorded before the rerun settled: %+v", runs)
		}
		g.answers[shipPRView] = openPR(greenRollup)
		shipCheck(t, f, `{}`)
		runs := stabilityRuns(t, f)
		if len(runs) != 1 || runs[0].Result != ciGreen || !slices.Equal(runs[0].Flakes, []string{"TestGarden", "TestPond"}) || len(runs[0].Infra) != 0 {
			t.Errorf("runs = %+v", runs)
		}
	})

	// Each case: a failed log, rerun green: the flakes and infra recorded.
	for _, c := range []struct {
		name, log     string
		flakes, infra []string
	}{
		{"subtests by full name", "test\tRun\t--- FAIL: TestFoo (0.1s)\ntest\tRun\t--- FAIL: TestScripts (1s)\ntest\tRun\t    --- FAIL: TestScripts/bar (1s)\n", []string{"TestFoo", "TestScripts/bar"}, nil},
		{"unknown step", "scripts (4/4)\tUNKNOWN STEP\t--- FAIL\n", nil, []string{"scripts (4/4) / UNKNOWN STEP"}},
		{"no job or step", "FAIL\n", nil, []string{"CI run 42"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			f, g, _ := shipFixture(t)
			g.answers[shipPRView] = openPR(rerunRollup)
			g.answers["run rerun 42 --failed"] = ""
			g.answers[logCall] = c.log
			shipCheck(t, f, `{"settled":true}`)
			g.answers[shipPRView] = openPR(greenRollup)
			shipCheck(t, f, `{}`)
			runs := stabilityRuns(t, f)
			if len(runs) != 1 || runs[0].Result != ciGreen || !slices.Equal(runs[0].Flakes, c.flakes) || !slices.Equal(runs[0].Infra, c.infra) {
				t.Errorf("runs = %+v", runs)
			}
		})
	}

	t.Run("old record reads job / step flakes as infra", func(t *testing.T) {
		f, _, _ := shipFixture(t)
		if err := writeAtomic(ciStabilityPath(f.common), []byte(`{"runs":[{"head":"h","result":"green","time":"x","flakes":["test / Run","CI run 9","TestA"]}]}`)); err != nil {
			t.Fatal(err)
		}
		runs := stabilityRuns(t, f)
		if len(runs) != 1 || !slices.Equal(runs[0].Flakes, []string{"TestA"}) || !slices.Equal(runs[0].Infra, []string{"test / Run", "CI run 9"}) {
			t.Errorf("runs = %+v", runs)
		}
	})

	t.Run("red after rerun", func(t *testing.T) {
		f, g, _ := shipFixture(t)
		g.answers[shipPRView] = openPR(rerunRollup)
		g.answers["run rerun 42 --failed"] = ""
		g.answers[logCall] = testLog
		shipCheck(t, f, `{"settled":true}`)
		shipCheck(t, f, `{"settled":true}`)
		shipCheck(t, f, `{"settled":true}`)
		runs := stabilityRuns(t, f)
		if len(runs) != 1 || runs[0].Result != ciRed || len(runs[0].Flakes) != 0 {
			t.Errorf("runs = %+v", runs)
		}
	})

	t.Run("main tree equal", func(t *testing.T) {
		f, g, _ := shipFixture(t)
		g.answers[shipPRView] = openPR(greenRollup)
		shipCheck(t, f, `{}`)
		mainHeld(g)
		g.answers[shipMergeArg] = ""
		g.answers[shipMergedPR] = `{"state":"MERGED","mergeCommit":{"oid":"fff0000aaa"}}`
		if _, err := callOp(t, f.ctx, "ship.merge", `{}`); err != nil {
			t.Fatal(err)
		}
		g.answers[gitMergeTree] = "tree1111"
		g.answers[gitHeadTree] = "tree1111"
		shipCheck(t, f, `{"main":true}`)
		runs := stabilityRuns(t, f)
		if len(runs) != 2 || !runs[1].Main || runs[1].Head != "fff0000aaa" || runs[1].Result != ciGreen {
			t.Errorf("runs = %+v", runs)
		}
		for _, c := range g.calls {
			if len(c) > 6 && c[:6] == "gh run" {
				t.Errorf("gh call %q for a tree-equal main", c)
			}
		}
	})

	t.Run("main tree differs", func(t *testing.T) {
		f, g, state := shipFixture(t)
		st := shipState{PR: 7, Head: "abc1234def", MergeCommit: "fff0000aaa"}
		if err := writeShipState(state, st); err != nil {
			t.Fatal(err)
		}
		g.answers[gitFetch] = ""
		g.answers[gitMergeTree] = "tree1111"
		g.answers[gitHeadTree] = "tree2222"
		g.answers["git ls-tree --name-only refs/remotes/origin/main .github/workflows/"] = ".github/workflows/ci.yml\n"
		g.answers["git show refs/remotes/origin/main:.github/workflows/ci.yml"] = "run-name: CI ${{ inputs.sha }}\non:\n  pull_request:\n  workflow_dispatch:\n    inputs:\n      sha:\n        required: true\n"
		g.answers["workflow run ci.yml --ref main -f sha=fff0000aaa"] = ""
		g.answers["run list --workflow ci.yml --event workflow_dispatch --json databaseId,name,displayTitle,status,conclusion --limit 100"] =
			`[{"databaseId":9,"name":"ci","displayTitle":"CI fff0000aaa","status":"completed","conclusion":"success"}]`
		r := shipCheck(t, f, `{"main":true}`)
		if !r.Passed {
			t.Fatalf("main: %+v", r)
		}
		runs := stabilityRuns(t, f)
		if len(runs) != 1 || !runs[0].Main || runs[0].Head != "fff0000aaa" || runs[0].Result != ciGreen {
			t.Errorf("runs = %+v", runs)
		}
	})
}

func TestCIStability(t *testing.T) {
	clean := func(n int) []CIRun {
		var rs []CIRun
		for range n {
			rs = append(rs, CIRun{Result: ciGreen})
		}
		return rs
	}
	flake := CIRun{Head: "h1", Result: ciGreen, Flakes: []string{"TestPond"}}
	flake2 := CIRun{Head: "h2", Result: ciGreen, Flakes: []string{"TestPond"}}
	cases := []struct {
		name     string
		runs     []CIRun
		window   int
		unstable bool
	}{
		{"empty", nil, 20, false},
		{"one flake in window is ignored", append([]CIRun{flake}, clean(19)...), 20, false},
		{"same head twice is ignored", []CIRun{flake, flake}, 20, false},
		{"two tests once each", []CIRun{{Head: "h1", Result: ciGreen, Flakes: []string{"TestX"}}, {Head: "h2", Result: ciGreen, Flakes: []string{"TestY"}}}, 20, false},
		{"two heads in window", append([]CIRun{flake, flake2}, clean(18)...), 20, true},
		{"older flake pushed out by 19 clean runs", append([]CIRun{flake, flake2}, clean(19)...), 20, false},
		{"window 5 shrinks", append([]CIRun{flake, flake2}, clean(4)...), 5, false},
		{"window 5 still holds", append([]CIRun{flake, flake2}, clean(3)...), 5, true},
		{"red main", []CIRun{{Main: true, Result: ciRed}}, 20, true},
		{"red main cleared by later green main", []CIRun{{Main: true, Result: ciRed}, {Main: true, Result: ciGreen}}, 20, false},
		{"red main outlives the window", append([]CIRun{{Main: true, Result: ciRed}}, clean(30)...), 20, true},
		{"red PR head is not main", []CIRun{{Result: ciRed}}, 20, false},
		{"job / step flake reads as infra", []CIRun{{Result: ciGreen, Flakes: []string{"scripts (4/4) / Run scripts"}}}, 20, false},
		{"unidentified flakes only", []CIRun{{Result: ciGreen, Flakes: []string{"ci / UNKNOWN STEP"}}, {Result: ciGreen, Flakes: []string{"CI run 123"}}}, 20, false},
		{"identified beside unidentified", []CIRun{{Head: "h1", Result: ciGreen, Flakes: []string{"ci / UNKNOWN STEP", "TestPond"}}, flake2}, 20, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := ciVerdict(CIStability{Runs: c.runs}, c.window)
			if v.Unstable != c.unstable || (v.Unstable && v.Reason == "") {
				t.Errorf("verdict = %+v, want unstable=%v", v, c.unstable)
			}
		})
	}

	t.Run("ignored unidentified flake named in the reason", func(t *testing.T) {
		v := ciVerdict(CIStability{Runs: []CIRun{{Result: ciGreen, Flakes: []string{"ci / UNKNOWN STEP"}}, {Result: ciGreen, Flakes: []string{"CI run 123", "ci / UNKNOWN STEP"}}}}, 20)
		if v.Unstable || v.Reason != "ignored infra failures: CI run 123, ci / UNKNOWN STEP" {
			t.Errorf("verdict = %+v", v)
		}
		if v := ciVerdict(CIStability{Runs: []CIRun{{Result: ciGreen, Flakes: []string{"CI run 123"}}, {Main: true, Result: ciRed, Head: "m1"}}}, 20); !v.Unstable || !strings.Contains(v.Reason, "main's latest CI result is red") {
			t.Errorf("red main beside an unidentified flake = %+v", v)
		}
	})

	t.Run("unstable reason names the test, its count and when it stops", func(t *testing.T) {
		for _, c := range []struct {
			window, clean int
			want          string
		}{
			{20, 3, "TestX flaked 2 times on 2 heads in the last 20 CI runs (stops applying after 16 more completed CI runs)"},
			{5, 2, "TestX flaked 2 times on 2 heads in the last 5 CI runs (stops applying after 2 more completed CI runs)"},
		} {
			rec := CIStability{Runs: append([]CIRun{{Head: "a", Result: ciGreen, Flakes: []string{"TestX"}}, {Head: "b", Result: ciGreen, Flakes: []string{"TestX"}}}, clean(c.clean)...)}
			if v := ciVerdict(rec, c.window); !v.Unstable || v.Reason != c.want {
				t.Fatalf("verdict = %+v, want %q", v, c.want)
			}
			// Stable again after exactly N more clean runs, not before.
			n := c.window - c.clean - 1
			rec.Runs = append(rec.Runs, clean(n-1)...)
			if !ciVerdict(rec, c.window).Unstable {
				t.Errorf("window %d: stable after %d more runs, one too soon", c.window, n-1)
			}
			rec.Runs = append(rec.Runs, clean(1)...)
			if ciVerdict(rec, c.window).Unstable {
				t.Errorf("window %d: still unstable after %d more runs", c.window, n)
			}
		}
	})

	t.Run("single flake reason says seen once and ignored", func(t *testing.T) {
		v := ciVerdict(CIStability{Runs: []CIRun{{Head: "a", Result: ciGreen, Flakes: []string{"TestX"}}}}, 20)
		if v.Unstable || v.Reason != "TestX was seen once and ignored" {
			t.Errorf("verdict = %+v", v)
		}
	})

	t.Run("ship check main recording clears red", func(t *testing.T) {
		f, g, state := shipFixture(t)
		path := ciStabilityPath(f.common)
		if err := writeAtomic(path, []byte(`{"runs":[{"head":"old","main":true,"result":"red","time":"x"}]}`)); err != nil {
			t.Fatal(err)
		}
		rec, _ := readCIStability(path)
		if !ciVerdict(rec, 20).Unstable {
			t.Fatal("red main not unstable")
		}
		if err := writeShipState(state, shipState{PR: 7, Head: "abc1234def", MergeCommit: "fff0000aaa"}); err != nil {
			t.Fatal(err)
		}
		g.answers[gitFetch] = ""
		g.answers[gitMergeTree] = "tree1111"
		g.answers[gitHeadTree] = "tree1111"
		shipCheck(t, f, `{"main":true}`)
		rec, _ = readCIStability(path)
		if v := ciVerdict(rec, 20); v.Unstable {
			t.Errorf("still unstable after green main: %+v", v)
		}
	})

	t.Run("ciWindow setting", func(t *testing.T) {
		d, ok := findSetting(CIWindowKey)
		if !ok || d.def != "20" || d.validate("abc") == nil || d.validate("0") == nil || d.validate("5") != nil {
			t.Errorf("setting = %+v, %v", d, ok)
		}
	})
}
