# workflow-env-verified

evidence · `01M7791P430NNMQY49YTRFF22T` · actor [claude-code-172d4a41](../actor/claude-code-172d4a41.md)

## result

PASS

## statement

Verified at aae248c36 (diff 9ea73a5ea..HEAD read). go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/workflow-env$' passes: a workflow with [env] FOO={run_tmp}, PLAIN=x runs a command step that sees FOO as an existing folder and PLAIN=x, with no GOTMPDIR, and the FOO folder is gone after the run; a workflow without [env] sees neither GOTMPDIR nor FOO. engine.stepEnv and session_run no longer add GOTMPDIR (EnvGoTmp removed); they set the [env] pairs and create the run tmp folder only when {run_tmp} is used. go test ./internal/workflow/... ./internal/run/... pass (TestEnvTable, TestStepsGetEnvRunTmp, TestRunTmpGoneAfterFailedRun, TestNoEnvNoRunTmp); TestSessionRun in internal/api passes (the package-level FAIL is the root guard seeing the live daemon rewrite ~/code/.rhiza-root/usage.live.json, which is unrelated). Attempt 1's finding is fixed: with this commit's build, rhiza workflow check passes on all five .rhiza/workflows/*.toml, and [env] sits below autonomous in bugfix and plan-bug.

## Links

### Out

- verifies [workflow-env](../capability/rhiza/workflows/workflow-definitions/workflow-file/workflow-env.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
