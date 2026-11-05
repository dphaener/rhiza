# ci-yml-cancels-stale-qa-1

evidence · `01M6M04Q83K5ACAA3688Q6W2CR` · actor [claude-code-ca676fb4](../actor/claude-code-ca676fb4.md)

## result

PASS

## statement

At 83445164 I parsed .github/workflows/ci.yml with PyYAML: on holds only pull_request (branches [main]), no push; concurrency group 'ci-${{ github.event.pull_request.number }}' with cancel-in-progress true; the job still runs make check. go test -trimpath ./internal/ciworkflow -count=1 -v ran TestCIRunsOnPullRequestsOnly and TestCICancelsStaleRuns, both PASS; the package is under ./... so make test runs it.

## Links

### Out

- verifies [ci-yml-cancels-stale](../capability/rhiza/integrity/ci-yml-cancels-stale.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
