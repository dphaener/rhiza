# run-steps-default-tmpdir-qa-1

evidence · `01M7R1XK9HJB18RZXE6J4RC9BS` · actor [claude-code-1fd79db3](../actor/claude-code-1fd79db3.md)

## result

PASS

## statement

QA with a scratch build: workflows with no [env] and with [env] naming no {run_tmp} gave command and cleanup steps TMPDIR=<run tmp>, the folder existed, a file written under it was gone after the run succeeded and after it failed; [env] TMPDIR=/x gave exactly one TMPDIR=/x; a stopped run kept its folder until resumed and ended

## Links

### Out

- verifies [run-steps-default-tmpdir](../capability/rhiza/workflows/workflow-definitions/workflow-file/run-steps-default-tmpdir.md)

## History

- `01M7NWFRZTGPHG4DHNN8B6NWSP` When a run ends, whatever its outcome, nothing its steps wrote to the temporary folder is left behind: the daemon sets TMPDIR to the run's {run_tmp} for every process a run's steps start unless the workflow's [env] sets TMPDIR itself, so the agents and commands of a workflow whose [env] names no {run_tmp} write under {run_tmp}, which is already removed when the run ends; the /tmp/rhiza-projection-* leak stays bug 01M7M3JZFW's
