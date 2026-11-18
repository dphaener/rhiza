# run-sessions-default-tmpdir-qa-1

evidence · `01M7R1XKBNJKX5SQS3GE8GK73J` · actor [claude-code-1fd79db3](../actor/claude-code-1fd79db3.md)

## result

PASS

## statement

QA with a scratch build: session run --dry-run --json of a run's session showed TMPDIR=<run tmp> without making it; a real start (fake agent) got that TMPDIR, the folder existed, and abandoning the run removed what it wrote; a workflow setting TMPDIR=/x gave the session /x only; scoped internal/api tests incl. shipRunEnv and watchdogMergeEnded pass

## Links

### Out

- verifies [run-sessions-default-tmpdir](../capability/rhiza/workflows/workflow-definitions/workflow-file/run-sessions-default-tmpdir.md)

## History

- `01M7NWFRZTGPHG4DHNN8B6NWSP` When a run ends, whatever its outcome, nothing its steps wrote to the temporary folder is left behind: the daemon sets TMPDIR to the run's {run_tmp} for every process a run's steps start unless the workflow's [env] sets TMPDIR itself, so the agents and commands of a workflow whose [env] names no {run_tmp} write under {run_tmp}, which is already removed when the run ends; the /tmp/rhiza-projection-* leak stays bug 01M7M3JZFW's
