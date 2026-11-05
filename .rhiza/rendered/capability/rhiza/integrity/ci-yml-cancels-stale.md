# ci-yml-cancels-stale

capability · `01M6HX48DQR5E1PQ4JM4QN5H0C` · parent [integrity](../integrity.md)

## acceptance

A Go test in the build parses .github/workflows/ci.yml and asserts that push is not among its triggers and pull_request is, and that it has a concurrency group keyed on the pull request number (github.event.pull_request.number) with cancel-in-progress: true; make check runs it.

## statement

Rhiza's own CI never runs on pushes to main: it runs on pull requests (and on dispatch, see ci-yml-dispatchable), and a newer push to a pull request cancels that pull request's in-progress run.

## Links

### Out

- decided-by [ci-main-and-cancel](../../../decision/ci-main-and-cancel.md)
- decided-by [ci-yml-agent-built](../../../decision/ci-yml-agent-built.md)
- decided-by [ci-main-landing-order](../../../decision/ci-main-landing-order.md)
- derived-from [plan-ci-runner](../../../source/plan-ci-runner.md)

### In

- verifies [ci-yml-cancels-stale-verified](../../../evidence/ci-yml-cancels-stale-verified.md)
- verifies [ci-yml-cancels-stale-qa-1](../../../evidence/ci-yml-cancels-stale-qa-1.md)
- verifies [pr-256-update-github-runner-to](../../../evidence/pr-256-update-github-runner-to.md)
- verifies [pr-265-fix-ship-s-new-main-tree-check](../../../evidence/pr-265-fix-ship-s-new-main-tree-check.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
- `01M6MA6P7KKMQPFBM6ZD6EMDM3` Fix Ship's new main-tree check (ship-main-tree-check, #256) stops every run whose PR merged after another PR: GitHub's squash merge onto an advanced main gives a merge commit whose tree differs from the tested head's tree (run 01M6KXGZYGRY stopped at wait-main: 'merge commit b9229f0 has tree 6b107af, but the tested head ab3969c has tree 0e2c59e'), and with CI no longer running on push to main nothing ever turns it green. Keep the no-retest rule for the common case (trees equal: no CI on main), and when they differ dispatch CI on that main commit (add workflow_dispatch to .github/workflows/ci.yml and the ci install template; gh workflow run ci.yml --ref <sha>) and wait for it, so the owner's ask holds: "Update github runner to: Cancel stale runs, and stop re-testing on main". Runs started before #256 (old wait-main on push CI) must also finish: on resume they should take the new path
