# run-steps-default-tmpdir-verified

evidence · `01M7R03FDAZ0B6W70WFDW2FTWX` · actor [claude-code-83056e5a](../actor/claude-code-83056e5a.md)

## result

PASS

## statement

At builder commit 0c5d5f2d9: re-ran Checks (capped go test -trimpath -p 2 ./internal/workflow/... ./internal/run/... exit 0, also -count=1 for TestRunEnvPairs, TestStepsGetDefaultTMPDIR (no [env] and [env] without {run_tmp}, each succeeded and failed: build/ci/agent/tidy cleanup steps get TMPDIR=Store.TmpDir(run), folder exists, a file written under it, nothing left after the run), TestOwnTMPDIRWins (TMPDIR=/x, one pair, FOO {run_tmp} and PLAIN unchanged), TestNoEnvNoRunTmp: all PASS; rhiza ship generate exit 0; rhiza render exit 0). Diff: stepEnv uses NeedsRunTmp/RunEnvPairs and is shared by runStep (all step kinds), arbiter, advisor and cleanup; EnvPairs unchanged (describe keeps the file's [env]); run.TmpDir and removal untouched; internal/api untouched; no TODOs. engine_test.go also changed (outside files): existing env expectations gained TMPDIR and the fake runner got an onEnv hook, both needed by the change.

## Links

### Out

- verifies [run-steps-default-tmpdir](../capability/rhiza/workflows/workflow-definitions/workflow-file/run-steps-default-tmpdir.md)

## History

- `01M7NWFRZTGPHG4DHNN8B6NWSP` When a run ends, whatever its outcome, nothing its steps wrote to the temporary folder is left behind: the daemon sets TMPDIR to the run's {run_tmp} for every process a run's steps start unless the workflow's [env] sets TMPDIR itself, so the agents and commands of a workflow whose [env] names no {run_tmp} write under {run_tmp}, which is already removed when the run ends; the /tmp/rhiza-projection-* leak stays bug 01M7M3JZFW's
