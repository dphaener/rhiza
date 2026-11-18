# run-sessions-default-tmpdir

capability · `01M7QZVZ213HF4N41XDNGP7P26` · parent [workflow-file](../workflow-file.md)

## acceptance

go test -trimpath ./internal/api/... -run 'SessionRun|Abandon|ShipRunEnv|Gotmp|RunTmp' passes with new tests: session run --dry-run --json of a session in a run whose workflow has no [env] (and one whose [env] names no {run_tmp}) shows TMPDIR=<run tmp> in env and the folder exists after a real start; abandoning a run whose step wrote a file under that TMPDIR leaves nothing of it (beside run_abandon_gotmp_test.go); (failure case) a workflow with TMPDIR = "/x" in [env] gives the session TMPDIR=/x only; the env key description in internal/workflow/spec.go and docs/custom-workflows.md say TMPDIR defaults to {run_tmp} unless [env] sets it, docs/reference regenerated with rhiza ship generate, and make self-check passes

## files

internal/api/session_run.go, internal/api/ship_build.go, internal/api/session_run_tmpdir_test.go, internal/api/run_abandon_gotmp_test.go, internal/workflow/spec.go, docs/custom-workflows.md, docs/reference/workflow-file.md, cmd/rhiza/testdata/script/describe-golden.txtar

## must_not

Change the engine (internal/run) or the helper's behaviour from run-steps-default-tmpdir; change where {run_tmp} lives; remove the /tmp/rhiza-projection-* folders (bug 01M7M3JZFW's)

## statement

A session a run spawns (rhiza session run for a run's session, such as a conductor's children) and the watchdog's view of a run (shipRunEnv) get the same run env as the run's steps, TMPDIR={run_tmp} included unless [env] sets TMPDIR, with the folder made, and an abandoned run leaves nothing under the tmp root; the [env] docs say so

## Links

### Out

- depends-on [run-steps-default-tmpdir](run-steps-default-tmpdir.md)
- decided-by [run-tmpdir-one-helper](../../../../../decision/run-tmpdir-one-helper.md)

### In

- verifies [run-sessions-default-tmpdir-verified-2](../../../../../evidence/run-sessions-default-tmpdir-verified-2.md)
- verifies [run-sessions-default-tmpdir-qa-1](../../../../../evidence/run-sessions-default-tmpdir-qa-1.md)

## History

- `01M7NWFRZTGPHG4DHNN8B6NWSP` When a run ends, whatever its outcome, nothing its steps wrote to the temporary folder is left behind: the daemon sets TMPDIR to the run's {run_tmp} for every process a run's steps start unless the workflow's [env] sets TMPDIR itself, so the agents and commands of a workflow whose [env] names no {run_tmp} write under {run_tmp}, which is already removed when the run ends; the /tmp/rhiza-projection-* leak stays bug 01M7M3JZFW's
