# ci-yml-cancels-stale

capability · `01M6HX48DQR5E1PQ4JM4QN5H0C` · parent [integrity](../integrity.md)

## acceptance

A Go test in the build parses .github/workflows/ci.yml and asserts that pull_request is its only trigger (no push), and that it has a concurrency group keyed on the pull request number (github.event.pull_request.number) with cancel-in-progress: true; make check runs it.

## statement

Rhiza's own CI runs only on pull requests, not on pushes to main, and a newer push to a pull request cancels that pull request's in-progress run.

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

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
