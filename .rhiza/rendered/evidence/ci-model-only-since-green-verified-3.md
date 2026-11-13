# ci-model-only-since-green-verified-3

evidence · `01M7915VHRKEX5B480ZY92QW0Q` · actor [claude-code-136ef0c3](../actor/claude-code-136ef0c3.md)

## result

PASS

## statement

At ec8a3ff45, go test -trimpath -count=1 -v ./internal/ciworkflow/... passed: TestChangesScript runs .github/scripts/ci-changes.sh (BASE, GREEN inputs) on scratch git repos and gets code=false for a land-like .rhiza/-only commit after green and for a merge of main bringing only .rhiza/, code=true for a code change since green, a merge of main bringing code, and no green with the PR leaving .rhiza/ (plus code=false for no green staying in .rhiza/, the merge-base fallback); TestCIChangesDecidesFromGitDiff asserts ci.yml's changes job runs the script with steps.green.outputs.sha (found via gh api with github.token), and TestCICheckAggregatesEveryJob passes. Read the script and changes job: the not-ancestor/empty GREEN fallback to git merge-base matches the updated statement. make self-check passed in a clean detached worktree with the run env cleared.

## Links

### Out

- verifies [ci-model-only-since-green](../capability/rhiza/integrity/self-check/ci-model-only-since-green.md)

## History

- `01M76XSFD0JZ9424EWBSQA3K8X` CI on a pull request runs only the model checks when nothing outside .rhiza/ changed since the pull request's last commit whose CI passed, so ship land's push and merges of main that bring only .rhiza/ changes cost a model check instead of a full CI round. Evidence: wait-landed is a full CI round in every run, 5.9m median (perf report rec 1).
