# ci-yml-cancels-stale-verified

evidence · `01M6KZSATNESPK1ZQHHM658XKB` · actor [claude-code-981fcdb2](../actor/claude-code-981fcdb2.md)

## result

PASS

## statement

At 6a7e233b: .github/workflows/ci.yml has only the pull_request trigger (push removed) and concurrency group ci-${{ github.event.pull_request.number }} with cancel-in-progress: true. New internal/ciworkflow/ciworkflow_test.go reads ci.yml and asserts pull_request is the only trigger, the group contains github.event.pull_request.number, and cancel-in-progress is true. go test -count=1 -v ./internal/ciworkflow: both tests PASS. Negative check in a scratch copy: adding a push trigger, setting cancel-in-progress false, or keying the group on github.ref each fails the matching test. The package is in go list ./..., and make test ran it (ok internal/ciworkflow). fmt-check and vet pass. The rest of make check was red for reasons outside this leaf: timing budgets (modelapp TestGeneratedModel, then api TestReadBudget) flaked under machine load ~7-8 on 4 CPUs, and self-check flagged stale .rhiza/rendered for actor records the run committed.

## Links

### Out

- verifies [ci-yml-cancels-stale](../capability/rhiza/integrity/ci-yml-cancels-stale.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
