# ci-install-cancels-stale

capability · `01M6HWZ4ZMHV1Q7BYEZF4MPFP7` · parent [integrity](../integrity.md)

## acceptance

The ci-install testscript's expected workflow holds a concurrency group keyed on the pull request number, falling back to the commit (github.sha) for pushes to main, whose cancel-in-progress is ${{ github.event_name == 'pull_request' }}, true only for pull_request events and never for pushes to main, and keeps the push-to-main trigger and its check steps.

## statement

The workflow rhiza ci install writes cancels a pull request's in-progress run when a newer push to it starts one, through a per-pull-request concurrency group with cancel-in-progress limited to pull_request events, and still checks every push to main: each push to main gets its own group keyed on its commit, since GitHub cancels a group's pending run when a newer one queues even without cancel-in-progress, and a push to main never cancels an earlier push's check, since the push check is incremental from github.event.before and a cancelled run would leave that push's changes unchecked.

## Links

### Out

- decided-by [ci-main-and-cancel](../../../decision/ci-main-and-cancel.md)
- derived-from [plan-ci-runner](../../../source/plan-ci-runner.md)

### In

- verifies [ci-install-cancels-stale-verified](../../../evidence/ci-install-cancels-stale-verified.md)
- verifies [ci-install-cancels-stale-qa-1](../../../evidence/ci-install-cancels-stale-qa-1.md)
- verifies [ci-install-cancels-stale-verified-2](../../../evidence/ci-install-cancels-stale-verified-2.md)
- verifies [ci-install-cancels-stale-verified-3](../../../evidence/ci-install-cancels-stale-verified-3.md)
- verifies [pr-256-update-github-runner-to](../../../evidence/pr-256-update-github-runner-to.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
