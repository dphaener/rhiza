# ci-model-only-since-green

capability · `01M76XZFT4TP194AK2Z1NCJH5A` · parent [self-check](../self-check.md)

## acceptance

The changes job's decision is a script file that takes the earlier green commit as an input (the workflow finds it with the job's GITHUB_TOKEN; the test never touches the network). A Go test in internal/ciworkflow (go test -trimpath ./internal/ciworkflow) runs the script against scratch git repos and asserts code=false for a land-like commit changing only .rhiza/ after a green commit and for a merge of main bringing only .rhiza/ changes, and code=true for a code change since green, for a merge of main bringing code, and when no green commit is given and the PR diff from base leaves .rhiza/; and asserts ci.yml's changes job runs that script and the check job still aggregates the verdict.

## statement

Rhiza's own CI (.github/workflows/ci.yml) runs only the model job on a pull request push when no path outside .rhiza/ differs between the tree this run tests and the tree of the pull request's latest run whose check verdict succeeded. With no usable earlier green commit (none, or one that is not an ancestor of the tested tree, as after a force push or rebase) the decision diffs from the merge base with the pull request's base instead, so a pull request that only touches .rhiza/ still runs only the model job. Any difference outside .rhiza/ (including code main brought in by a merge) runs the Go jobs as today.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [ci-model-only-since-green-verified](../../../../evidence/ci-model-only-since-green-verified.md)
- verifies [ci-model-only-since-green-qa-1](../../../../evidence/ci-model-only-since-green-qa-1.md)
- verifies [ci-model-only-since-green-verified-2](../../../../evidence/ci-model-only-since-green-verified-2.md)
- verifies [ci-model-only-since-green-verified-3](../../../../evidence/ci-model-only-since-green-verified-3.md)
- verifies [pr-512-ci-since-last-green](../../../../evidence/pr-512-ci-since-last-green.md)

## History

- `01M76XSFD0JZ9424EWBSQA3K8X` CI on a pull request runs only the model checks when nothing outside .rhiza/ changed since the pull request's last commit whose CI passed, so ship land's push and merges of main that bring only .rhiza/ changes cost a model check instead of a full CI round. Evidence: wait-landed is a full CI round in every run, 5.9m median (perf report rec 1).
