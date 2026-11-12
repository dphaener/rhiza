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
	"strconv"
	"strings"
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
	// Flakes name the tests the head's first attempt failed (full names,
	// subtests as 'TestScripts/bar') before a rerun of the same head
	// passed: set only on a green result that followed a rerun.
	Flakes []string `json:"flakes,omitempty" desc:"The tests that failed before a rerun of the same head passed"`
	// Infra names the failures that name no test ('job / step', 'UNKNOWN
	// STEP', 'CI run <id>') before such a rerun passed. They are kept
	// apart from flakes: they make CI neither unstable nor a bug.
	Infra []string `json:"infra,omitempty" desc:"Failures that named no test ('job / step', 'CI run <id>') before a rerun of the same head passed"`
}

// infraEntry reports whether a recorded flake names no test: a 'job /
// step' (UNKNOWN STEP included) or 'CI run <id>'.
func infraEntry(name string) bool {
	return strings.HasPrefix(name, "CI run ") || strings.Contains(name, " / ")
}

// normalized reads an entry recorded before infra failures were kept
// apart: flakes holding only 'job / step' or 'CI run' names are infra.
func (r CIRun) normalized() CIRun {
	if !slices.ContainsFunc(r.Flakes, infraEntry) {
		return r
	}
	var flakes, infra []string
	infra = slices.Clone(r.Infra)
	for _, f := range r.Flakes {
		if !infraEntry(f) {
			flakes = append(flakes, f)
		} else if !slices.Contains(infra, f) {
			infra = append(infra, f)
		}
	}
	r.Flakes, r.Infra = flakes, infra
	return r
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
	for i, r := range rec.Runs {
		rec.Runs[i] = r.normalized()
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
			run.Infra = slices.Clone(st.RerunInfra)
			run = run.normalized()
		}
	}
	return recordCIRun(ctx, t, run)
}

// CIVerdict is Rhiza's judgment of a repo's CI, read from the record.
type CIVerdict struct {
	Unstable bool   `json:"unstable" desc:"True while the window holds an identified flake or main's latest result is red"`
	Reason   string `json:"reason,omitempty" desc:"Why CI is unstable; when stable, the unidentified flakes it ignored, if any"`
}

// ciVerdict judges rec: CI is unstable while the last window runs hold an
// identified flake, or while main's latest result is red. Infra
// failures are ignored, and a stable verdict's reason names them. window
// below 1 counts as 1.
func ciVerdict(rec CIStability, window int) CIVerdict {
	window = max(window, 1)
	runs := rec.Runs
	if len(runs) > window {
		runs = runs[len(runs)-window:]
	}
	var ignored []string
	for _, r := range runs {
		r = r.normalized()
		if len(r.Flakes) > 0 {
			return CIVerdict{Unstable: true, Reason: "flake in the last " + strconv.Itoa(window) + " CI runs: " + strings.Join(r.Flakes, ", ")}
		}
		for _, f := range r.Infra {
			if !slices.Contains(ignored, f) {
				ignored = append(ignored, f)
			}
		}
	}
	for i := len(rec.Runs) - 1; i >= 0; i-- {
		if r := rec.Runs[i]; r.Main {
			if r.Result == ciRed {
				return CIVerdict{Unstable: true, Reason: "main's latest CI result is red (" + r.Head + ")"}
			}
			break
		}
	}
	if len(ignored) > 0 {
		return CIVerdict{Reason: "ignored unidentified flake in the last " + strconv.Itoa(window) + " CI runs: " + strings.Join(ignored, ", ")}
	}
	return CIVerdict{}
}

// ciStabilityVerdict reads the repo's record and the ciWindow setting.
func ciStabilityVerdict(ctx *Ctx) (CIVerdict, error) {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return CIVerdict{}, err
	}
	rec, err := readCIStability(ciStabilityPath(common))
	if err != nil {
		return CIVerdict{}, err
	}
	rows, err := settingsOf(ctx, nil)
	if err != nil {
		return CIVerdict{}, err
	}
	window, err := strconv.Atoi(settingValue(rows, CIWindowKey))
	if err != nil {
		return CIVerdict{}, Usagef("%s: %v", CIWindowKey, err).WithNext("rhiza daemon settings")
	}
	return ciVerdict(rec, window), nil
}
