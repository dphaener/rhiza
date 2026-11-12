# ci-model-only-since-green-verified

evidence · `01M78ZW3KFJ3Y1SWGENPPC91X7` · actor [claude-code-50966a9b](../actor/claude-code-50966a9b.md)

## result

PASS

## statement

Verifier at d5c74fb19: go test -trimpath -count=1 ./internal/ciworkflow/... -v passes (go vet clean, bash -n on the script clean). TestChangesScript runs .github/scripts/ci-changes.sh BASE [GREEN] against scratch git repos (no network) and asserts code=false for a land-like .rhiza/-only commit after green and for a merge of main bringing only .rhiza/, code=true for a code change since green, for a merge of main bringing code, and with no green commit when the PR diff from base leaves .rhiza/ (plus code=false when it stays in .rhiza/). TestCIChangesDecidesFromGitDiff asserts the changes job runs .github/scripts/ci-changes.sh fed steps.green.outputs.sha, and TestCICheckAggregatesEveryJob still passes on the check job. Read the diff: the green step lists the PR's commits newest first via gh api with the job's github.token, head excluded, and takes the first whose check run named check concluded success; a green that is not an ancestor of HEAD falls back to the merge base. Note: ci.yml declares no permissions block, so the gh api reads rely on the repo's default token permissions (fine for a public repo); and since HEAD is the PR merge commit, an advanced main makes the diff conservative (code=true), never unsafe.

## Links

### Out

- verifies [ci-model-only-since-green](../capability/rhiza/integrity/self-check/ci-model-only-since-green.md)

## History

- `01M76XSFD0JZ9424EWBSQA3K8X` CI on a pull request runs only the model checks when nothing outside .rhiza/ changed since the pull request's last commit whose CI passed, so ship land's push and merges of main that bring only .rhiza/ changes cost a model check instead of a full CI round. Evidence: wait-landed is a full CI round in every run, 5.9m median (perf report rec 1).
