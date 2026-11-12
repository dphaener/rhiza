# workflow-env-verified-2

evidence · `01M77BXNQ5B4XJVP0BZVGNH2QY` · actor [claude-code-6141ef24](../actor/claude-code-6141ef24.md)

## result

PASS

## statement

Verifier at ca5773613 (attempt 2). go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(workflow-env|daemon-settings-scope)$' passes: workflow-env.txtar runs a workflow with [env] FOO={run_tmp} and PLAIN=x, its command step sees FOO as an existing folder and PLAIN=x and no GOTMPDIR, the folder is gone after the run, and a workflow without [env] sees neither GOTMPDIR nor FOO. No non-test Go code under internal/ or cmd/ sets GOTMPDIR. go test -trimpath ./internal/workflow/... ./internal/run/... pass. The fix commit adds TestSessionRunWorkflowEnv (session.run gives an agent a run's [env]: {run_tmp}/go with the folder made, and a plain value) and an env key to FileKeys and docs/reference/workflow-file.md; go test -trimpath -count=1 ./internal/api -run TestSession passes in a scratch worktree outside ~/code (under ~/code the rootguard reports usage.live.json, which the live Claude Code statusline rewrites every second, not this change). A built rhiza's workflow check accepts every .rhiza/workflows/*.toml, bugfix.toml and plan-bug.toml included.

## Links

### Out

- verifies [workflow-env](../capability/rhiza/workflows/workflow-definitions/workflow-file/workflow-env.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
