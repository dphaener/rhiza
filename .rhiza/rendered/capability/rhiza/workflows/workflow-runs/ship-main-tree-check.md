# ship-main-tree-check

capability · `01M6HWX807SBGNFEPCEAHVK0QM` · parent [workflow-runs](../workflow-runs.md)

## acceptance

Testscripts with a local bare origin and a fake gh: a ship-change run merges and passes the main step straight to close-out with no gh run list or gh workflow run call for the merge commit; when origin/main moves after the PR head's checks pass, ship merge makes no gh pr merge call and the run goes back to rebase-merge, rebases, waits on the new head's checks and then merges it; no step waits on push CI for main.

## statement

After merging, ship-change confirms main by checking locally with git that the merge commit's tree equals the tree of the PR head whose checks passed, instead of waiting for CI runs on main; a differing tree is confirmed by CI dispatched on the merge commit (ship-main-dispatches-ci). So that a base that moved after the checks never yields an untested main: `rhiza ship merge` first fetches origin and, when origin/main is not an ancestor of the tested head (another change set merged after the head's checks), merges nothing and fails with an outcome that sends ship-change back to its rebase-merge step, which rebases onto the new main and waits on the PR's checks again. A head that holds origin/main squashes to a merge commit whose tree is the tested head's tree.

## Links

### Out

- decided-by [ci-main-and-cancel](../../../../decision/ci-main-and-cancel.md)
- decided-by [ci-main-landing-order](../../../../decision/ci-main-landing-order.md)
- decided-by [ship-main-dispatch-ci](../../../../decision/ship-main-dispatch-ci.md)
- derived-from [plan-ci-runner](../../../../source/plan-ci-runner.md)

### In

- verifies [ship-main-tree-check-verified](../../../../evidence/ship-main-tree-check-verified.md)
- verifies [ship-main-tree-check-qa-1](../../../../evidence/ship-main-tree-check-qa-1.md)
- verifies [ship-main-tree-check-verified-2](../../../../evidence/ship-main-tree-check-verified-2.md)
- verifies [pr-256-update-github-runner-to](../../../../evidence/pr-256-update-github-runner-to.md)
- verifies [pr-265-fix-ship-s-new-main-tree-check](../../../../evidence/pr-265-fix-ship-s-new-main-tree-check.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
- `01M6MA6P7KKMQPFBM6ZD6EMDM3` Fix Ship's new main-tree check (ship-main-tree-check, #256) stops every run whose PR merged after another PR: GitHub's squash merge onto an advanced main gives a merge commit whose tree differs from the tested head's tree (run 01M6KXGZYGRY stopped at wait-main: 'merge commit b9229f0 has tree 6b107af, but the tested head ab3969c has tree 0e2c59e'), and with CI no longer running on push to main nothing ever turns it green. Keep the no-retest rule for the common case (trees equal: no CI on main), and when they differ dispatch CI on that main commit (add workflow_dispatch to .github/workflows/ci.yml and the ci install template; gh workflow run ci.yml --ref <sha>) and wait for it, so the owner's ask holds: "Update github runner to: Cancel stale runs, and stop re-testing on main". Runs started before #256 (old wait-main on push CI) must also finish: on resume they should take the new path
