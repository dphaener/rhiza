# ci-main-and-cancel

decision · `01M6HWFGAPK603YSDHDPB4BHAK`

## reasoning

Milestone 25 is the open milestone and its intent is building and shipping the queued owner requests; this change set's intent is one of those requests, so it serves milestone-25. On main: the question states ship-change rebases a PR that main has moved past and merges only a head whose checks passed, so with squash merges the merge commit's tree equals the tested head's tree; comparing trees keeps the guarantee wait-main gave (main is what was tested) without the CI run the change set removes, and failing loudly on a mismatch keeps a real gap visible instead of silently passing. On the template: cancelling superseded PR runs is safe for any repo, but dropping the main check is only safe where a merge tool guarantees the merged tree was tested, which ci install cannot assume about other repos. On landing: memory and the question both say the agents' gh token lacks the `workflow` scope, so any agent push touching .github/workflows/ is rejected and would stall the run; keeping the workflow edit out of agent leaves lets the Go work ship through ship-change, and ordering the hand edit after it means main never stops testing while ship still waits on main runs.

## statement

(1) The change set serves milestone-25 (empty the backlog): it is an owner request from the queue, and cheaper, non-redundant CI helps the daemon ship the backlog on its own. (2) Main stops running CI on push, and ship-change stops waiting for CI on main: the wait-main step (`rhiza ship check --main --settled`) instead passes when the merge commit's tree equals the tree of the PR head whose checks passed, checked locally with git and no wait on GitHub; if the trees differ it fails with a stop: line for a look, and never waits an hour for runs that will not come. (3) Both change: this repo's .github/workflows/ci.yml drops the push-to-main trigger and gains a concurrency group per pull request with cancel-in-progress; the `rhiza ci install` template in internal/api/ci.go gains the same cancel-stale concurrency group, but keeps its push-to-main check, since other repos may merge by hand without ship-change's tested-tree guarantee. (4) Superseded by Decision ci-yml-agent-built, which holds instead: this part had kept the ci.yml edit out of agent leaves for the owner to push by hand, because the agents' gh token lacks the `workflow` scope; origin now pushes over SSH, so the ci.yml edit is an agent-built leaf.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-226cf21d-1](../question/question-226cf21d-1.md)
- derived-from [plan-ci-runner](../source/plan-ci-runner.md)

### In

- decided-by [ship-main-tree-check](../capability/rhiza/workflows/workflow-runs/ship-main-tree-check.md)
- decided-by [ci-install-cancels-stale](../capability/rhiza/integrity/ci-install-cancels-stale.md)
- decided-by [ci-yml-cancels-stale](../capability/rhiza/integrity/ci-yml-cancels-stale.md)
- supersedes [ci-yml-agent-built](ci-yml-agent-built.md) — Supersedes part (4) only: the ci.yml edit is an agent-built leaf, not a hand push by the owner. Parts (1) to (3) stand.

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
