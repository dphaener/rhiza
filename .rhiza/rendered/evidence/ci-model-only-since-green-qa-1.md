# ci-model-only-since-green-qa-1

evidence · `01M790GBBC2JZC15DTP4KACRSS` · actor [claude-code-c8e6c442](../actor/claude-code-c8e6c442.md)

## result

PASS

## statement

QA at 4ed003163: ran .github/scripts/ci-changes.sh by hand in scratch git repos: code=true for code since base with no green; code=false for a .rhiza/-only commit after green, for an Actions-style merge-commit HEAD whose PR side added only .rhiza/ since green, and for a merge of main bringing only .rhiza/; code=true for a merge of main bringing code, a bogus green sha, a green that is not an ancestor (rebase), and a path with a space. Ran ci.yml's green-commit lookup (gh api pulls/N/commits + check-runs named check) read-only against real PRs 509 and 510: it picked each PR's newest earlier commit whose check succeeded. go test -trimpath ./internal/ciworkflow/... passes; bash -n clean; a built rhiza check passes from a clean detached worktree of HEAD.

## Links

### Out

- verifies [ci-model-only-since-green](../capability/rhiza/integrity/self-check/ci-model-only-since-green.md)

## History

- `01M76XSFD0JZ9424EWBSQA3K8X` CI on a pull request runs only the model checks when nothing outside .rhiza/ changed since the pull request's last commit whose CI passed, so ship land's push and merges of main that bring only .rhiza/ changes cost a model check instead of a full CI round. Evidence: wait-landed is a full CI round in every run, 5.9m median (perf report rec 1).
