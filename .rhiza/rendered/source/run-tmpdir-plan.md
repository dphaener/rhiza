# run-tmpdir-plan

source · `01M7QZVBWQNA65EN997T5MH4PJ`

## quotes

```json
[]
```

## raw_sha256

ca591c14adc41d76a435576a5149dc63bb6ce2dd24a5da9160389e3c46227c61

## session

0aa52902-753a-42c5-8fc5-01a1effc83e3

## summary

Autonomous plan of change set 01M7NWFR in run 01M7QZRN; no owner present, nothing asked. Research: a workflow's [env] reaches processes in three places: internal/run/engine.go stepEnv (command, wait, gate, agent, arbiter/advisor and cleanup steps), internal/api/session_run.go (sessions a run spawns, e.g. a conductor's children; today only when len(def.Env) > 0) and internal/api/ship_build.go shipRunEnv (the watchdog acting on a run from outside). All three build from workflow.Definition.EnvPairs and make {run_tmp} only when UsesRunTmp. {run_tmp} is removed by engine.Run's defer and by run abandon (internal/api/run.go) after cleanup steps run. Claude Code puts its scratch folder under $TMPDIR, so agent scratch lands in {run_tmp} too. Plan: one workflow-package helper giving the run's env with TMPDIR={run_tmp} unless [env] sets TMPDIR, used at all three sites, the folder made whenever the helper names it; then the api side (spawned sessions, abandon) and the docs.

## title

Run TMPDIR plan: research pass (autonomous conduct)

## Links

### In

- derived-from [run-tmpdir-one-helper](../decision/run-tmpdir-one-helper.md)

## History

- `01M7NWFRZTGPHG4DHNN8B6NWSP` When a run ends, whatever its outcome, nothing its steps wrote to the temporary folder is left behind: the daemon sets TMPDIR to the run's {run_tmp} for every process a run's steps start unless the workflow's [env] sets TMPDIR itself, so the agents and commands of a workflow whose [env] names no {run_tmp} write under {run_tmp}, which is already removed when the run ends; the /tmp/rhiza-projection-* leak stays bug 01M7M3JZFW's
