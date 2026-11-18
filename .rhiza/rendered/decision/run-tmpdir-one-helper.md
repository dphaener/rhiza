# run-tmpdir-one-helper

decision · `01M7QZVYY04XMGRFM45R6QSX1D`

## statement

The default TMPDIR={run_tmp} comes from one helper in internal/workflow beside EnvPairs, used by every place a run's [env] reaches a process (engine stepEnv, session_run for spawned sessions, shipRunEnv for the watchdog), so steps, arbiter and advisor, cleanup steps and spawned sessions all get it; a TMPDIR key in the workflow's own [env] wins, and the folder is made whenever the env names it

## wording

(assumed, would have escalated) chosen by the conductor from the research pass; autonomous run, no owner present

## Links

### Out

- derived-from [run-tmpdir-plan](../source/run-tmpdir-plan.md)

### In

- decided-by [run-steps-default-tmpdir](../capability/rhiza/workflows/workflow-definitions/workflow-file/run-steps-default-tmpdir.md)
- decided-by [run-sessions-default-tmpdir](../capability/rhiza/workflows/workflow-definitions/workflow-file/run-sessions-default-tmpdir.md)

## History

- `01M7NWFRZTGPHG4DHNN8B6NWSP` When a run ends, whatever its outcome, nothing its steps wrote to the temporary folder is left behind: the daemon sets TMPDIR to the run's {run_tmp} for every process a run's steps start unless the workflow's [env] sets TMPDIR itself, so the agents and commands of a workflow whose [env] names no {run_tmp} write under {run_tmp}, which is already removed when the run ends; the /tmp/rhiza-projection-* leak stays bug 01M7M3JZFW's
