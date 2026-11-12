# ci-model-only-since-green-verified-2

evidence · `01M790RD7CDBYQWP81ABYKMMA6` · actor [claude-code-85f60a52](../actor/claude-code-85f60a52.md)

## result

PASS

## statement

Verifier attempt 2 at fedd6cd77 (detached checkout). go test -trimpath -count=1 ./internal/ciworkflow/... exit 0: TestChangesScript passes its subtests land-like .rhiza/-only commit after green (code=false), merge of main bringing only .rhiza/ (code=false), code change since green (code=true), merge of main bringing code (code=true), no green commit with PR leaving .rhiza/ (code=true); TestCIChangesDecidesFromGitDiff, TestCIGoJobsRunOnlyOnCodeChanges, TestCIModelJobRunsOnlyOnModelChanges and TestCICheckAggregatesEveryJob pass. Read ci.yml: the changes job finds the PR's last green commit with gh and the job token and runs .github/scripts/ci-changes.sh "$BASE" "$GREEN"; check needs every job with if: always(). Re-ran the commit's Checks with a scratch build: rhiza ship generate exit 0 (0 files changed), rhiza render exit 0 (only this verifier's own untracked actor rendered), make self-check exit 0. The fix commit sets ci-model-only-pr's statement and acceptance to defer the PR decision rule to ci-model-only-since-green; it no longer mentions a plain git diff against the merge base, resolving the model check's finding f1.

## Links

### Out

- verifies [ci-model-only-since-green](../capability/rhiza/integrity/self-check/ci-model-only-since-green.md)

## History

- `01M76XSFD0JZ9424EWBSQA3K8X` CI on a pull request runs only the model checks when nothing outside .rhiza/ changed since the pull request's last commit whose CI passed, so ship land's push and merges of main that bring only .rhiza/ changes cost a model check instead of a full CI round. Evidence: wait-landed is a full CI round in every run, 5.9m median (perf report rec 1).
