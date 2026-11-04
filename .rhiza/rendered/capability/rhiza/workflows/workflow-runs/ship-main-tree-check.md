# ship-main-tree-check

capability · `01M6HWX807SBGNFEPCEAHVK0QM` · parent [workflow-runs](../workflow-runs.md)

## acceptance

Testscripts with a local bare origin and a fake gh: a ship-change run merges and passes the main step straight to close-out with no gh run list call for the merge commit; when origin/main moves after the PR head's checks pass, ship merge makes no gh pr merge call and the run goes back to rebase-merge, rebases, waits on the new head's checks and then merges it; a run whose recorded merge commit's tree differs from the recorded tested head fails the main step with a stop: line naming both commits; no step waits on CI for main.

## statement

After merging, ship-change confirms main by checking locally with git that the merge commit's tree equals the tree of the PR head whose checks passed, instead of waiting for CI runs on main; a differing tree fails the step with a stop: line. So that a base that moved after the checks never yields an untested main: `rhiza ship merge` first fetches origin and, when origin/main is not an ancestor of the tested head (another change set merged after the head's checks), merges nothing and fails with an outcome that sends ship-change back to its rebase-merge step, which rebases onto the new main and waits on the PR's checks again. A head that holds origin/main squashes to a merge commit whose tree is the tested head's tree.

## Links

### Out

- decided-by [ci-main-and-cancel](../../../../decision/ci-main-and-cancel.md)
- decided-by [ci-main-landing-order](../../../../decision/ci-main-landing-order.md)
- derived-from [plan-ci-runner](../../../../source/plan-ci-runner.md)

### In

- verifies [ship-main-tree-check-verified](../../../../evidence/ship-main-tree-check-verified.md)
- verifies [ship-main-tree-check-qa-1](../../../../evidence/ship-main-tree-check-qa-1.md)
- verifies [ship-main-tree-check-verified-2](../../../../evidence/ship-main-tree-check-verified-2.md)
- verifies [pr-256-update-github-runner-to](../../../../evidence/pr-256-update-github-runner-to.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
