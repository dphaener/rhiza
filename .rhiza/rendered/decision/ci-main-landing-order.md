# ci-main-landing-order

decision · `01M6HY0KT6ZKQZRRTG5PHA0ZZN`

## statement

This change set's own ship lands safely with the push-to-main trigger dropped in the same PR, because the wait-main step keeps its command `rhiza ship check --main --settled` unchanged and needs only ship state the old `rhiza ship merge` already records (the PR, the tested head, the merge commit). Its own ship runs the old installed binary, which sees no CI runs on the merge commit and keeps waiting (no checks is not settled, never a failure). The wait step re-runs `rhiza` from PATH every 30s, and the daemon's self-update (internal/daemon/selfupdate.go) fetches origin/main at most every 5 minutes and runs make update when main moves, so within minutes the new binary answers wait-main by the tree check and the step passes well inside its 1h timeout. If wait-main is still waiting 15 minutes after the merge, the orchestrator runs make update on the VM's main checkout by hand; it never edits the run's state to pass the step.

## wording

(reviser, derived from the code: internal/workflow/builtin/ship-change.toml wait-main polls `rhiza ship check --main --settled` every 30s with a 1h timeout; internal/api/ship.go checksVerdict treats no checks as not settled; internal/daemon/selfupdate.go rebuilds the installed binary when origin/main moves)

## Links

### Out

- derived-from [plan-ci-runner](../source/plan-ci-runner.md)

### In

- decided-by [ship-main-tree-check](../capability/rhiza/workflows/workflow-runs/ship-main-tree-check.md)
- decided-by [ci-yml-cancels-stale](../capability/rhiza/integrity/ci-yml-cancels-stale.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
