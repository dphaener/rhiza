package api

// This file holds the CI stability record: the repo's account of the CI runs
// ship.check settled, on run PR heads and on main. It lives under the git
// common dir, so every worktree of the repo reads and extends the same one.
// A run that went red and then passed on a rerun of the same head leaves one
// flake per failing test; the rolling window and the unstable/stable verdict
// are read from it elsewhere.

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"time"
)

// ciStabilityFile is the record's path under the git common dir.
const ciStabilityFile = "rhiza/ci-stability.json"

// CI result values of a record entry.
const (
	ciGreen = "green"
	ciRed   = "red"
)

// ciStabilityKeep bounds the record: older entries are dropped.
const ciStabilityKeep = 500

// CIRun is one settled CI run in the record.
type CIRun struct {
	Head   string `json:"head" desc:"The commit CI ran on: a run's PR head, or main's merge commit"`
	PR     int    `json:"pr,omitempty" desc:"The pull request's number; 0 if unknown"`
	Main   bool   `json:"main,omitempty" desc:"True for main's result (the merge commit), false for a PR head"`
	Result string `json:"result" desc:"green or red"`
	Time   string `json:"time" desc:"When the run was recorded, RFC 3339"`
	// Flakes are the tests that failed on the head's first attempt and
	// whose job then passed on a rerun of the same head: set only on a
	// green result that followed a rerun.
	Flakes []string `json:"flakes,omitempty" desc:"The tests that failed before a rerun of the same head passed"`
}

// CIStability is the record, oldest entry first.
type CIStability struct {
	Runs []CIRun `json:"runs"`
}

func ciStabilityPath(common string) string {
	return filepath.Join(common, filepath.FromSlash(ciStabilityFile))
}

// readCIStability reads the record at path; a missing file is empty.
func readCIStability(path string) (CIStability, error) {
	var rec CIStability
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return rec, nil
	}
	if err != nil {
		return rec, Errorf("read CI stability record: %v", err)
	}
	if err := json.Unmarshal(b, &rec); err != nil {
		return rec, Errorf("CI stability record %s: %v", path, err)
	}
	return rec, nil
}

// recordCIRun appends run to the record, once per (head, main): a repeated
// check of a settled head adds nothing. A dry run records nothing.
func recordCIRun(ctx *Ctx, t shipTarget, run CIRun) error {
	if ctx.DryRun {
		return nil
	}
	path := ciStabilityPath(t.common)
	rec, err := readCIStability(path)
	if err != nil {
		return err
	}
	if slices.ContainsFunc(rec.Runs, func(r CIRun) bool { return r.Head == run.Head && r.Main == run.Main }) {
		return nil
	}
	now := time.Now
	if ctx.Now != nil {
		now = ctx.Now
	}
	run.Time = now().UTC().Format(time.RFC3339)
	rec.Runs = append(rec.Runs, run)
	if n := len(rec.Runs); n > ciStabilityKeep {
		rec.Runs = rec.Runs[n-ciStabilityKeep:]
	}
	b, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return Errorf("CI stability record: %v", err)
	}
	if err := writeAtomic(path, append(b, '\n')); err != nil {
		return Errorf("write CI stability record: %v", err)
	}
	return nil
}

// recordSettledCI records a settled check's result for commit: red when it
// is not green; green with the flakes read before the head's rerun when it
// was rerun. Callers pass a settled response.
func recordSettledCI(ctx *Ctx, t shipTarget, st shipState, resp ShipCheckResponse, main bool) error {
	run := CIRun{Head: resp.Commit, PR: resp.PR, Main: main, Result: ciRed}
	if resp.Green {
		run.Result = ciGreen
		if st.RerunHead == resp.Commit {
			run.Flakes = slices.Clone(st.RerunTests)
		}
	}
	return recordCIRun(ctx, t, run)
}
