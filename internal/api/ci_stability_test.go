package api

import (
	"path/filepath"
	"slices"
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
// failing test of a red head whose rerun passed, a red with no flake, and
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
		if len(runs) != 1 || runs[0].Result != ciGreen || !slices.Equal(runs[0].Flakes, []string{"TestGarden", "TestPond"}) {
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
