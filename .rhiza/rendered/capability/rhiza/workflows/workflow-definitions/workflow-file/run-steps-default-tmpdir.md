# run-steps-default-tmpdir

capability · `01M7QZVYZYEW2YQJQR0HXV9K9G` · parent [workflow-file](../workflow-file.md)

## acceptance

go test -trimpath ./internal/workflow/... ./internal/run/... passes with new tests (beside gotmp_test.go): a workflow whose [env] names no {run_tmp} (and one with no [env] at all) gives its command, wait, agent and cleanup steps TMPDIR = Store.TmpDir(run), the folder exists as each step runs, a step writes a file under $TMPDIR, and after the run ends succeeded and again after it ends failed nothing of it is left under the tmp root; (failure case) a workflow whose [env] sets TMPDIR = "/x" gets TMPDIR=/x and no second TMPDIR pair in the step env; a workflow's other [env] keys and their {run_tmp} filling are unchanged

## files

internal/workflow/workflow.go, internal/workflow/workflow_test.go, internal/run/engine.go, internal/run/gotmp_test.go

## must_not

Change where {run_tmp} lives (run.TmpDir) or when the engine removes it; add TMPDIR to a workflow's [env] table or to EnvPairs' output as the file's own (workflow describe keeps showing only the file's [env]); touch internal/api

## statement

Every process a run's engine starts (command, wait, gate and agent steps, the arbiter and advisor, cleanup steps) gets TMPDIR set to the run's {run_tmp}, made before the step starts, unless the workflow's [env] sets TMPDIR itself, so what it writes under $TMPDIR is removed with {run_tmp} when the run ends

## Links

### Out

- decided-by [run-tmpdir-one-helper](../../../../../decision/run-tmpdir-one-helper.md)

### In

- depends-on [run-sessions-default-tmpdir](run-sessions-default-tmpdir.md)
- verifies [run-steps-default-tmpdir-verified](../../../../../evidence/run-steps-default-tmpdir-verified.md)
- verifies [run-steps-default-tmpdir-qa-1](../../../../../evidence/run-steps-default-tmpdir-qa-1.md)

## History

- `01M7NWFRZTGPHG4DHNN8B6NWSP` When a run ends, whatever its outcome, nothing its steps wrote to the temporary folder is left behind: the daemon sets TMPDIR to the run's {run_tmp} for every process a run's steps start unless the workflow's [env] sets TMPDIR itself, so the agents and commands of a workflow whose [env] names no {run_tmp} write under {run_tmp}, which is already removed when the run ends; the /tmp/rhiza-projection-* leak stays bug 01M7M3JZFW's
