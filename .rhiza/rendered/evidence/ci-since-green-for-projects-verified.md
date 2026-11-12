# ci-since-green-for-projects-verified

evidence · `01M7905W557GWKEZCMG8X0X0C3` · actor [claude-code-5506a83e](../actor/claude-code-5506a83e.md)

## result

PASS

## statement

At 479d4d7f8, go test -trimpath ./internal/api -run 'Docs|CIInstall' passes (TestDocsCISinceGreen, TestCIInstall, TestCIInstallLatest, TestCIInstallErrors). Read docs/ci-since-last-green.md and the diffs of 289b40082 and 479d4d7f8: it is the owner's chosen form (decision ci-since-green-recipe), a documented recipe in the user docs, linked from README.md, with an example changes job that uses only git and gh, plus the rule (no green commit, a non-ancestor green commit, or any path outside .rhiza/ means everything runs) and an aggregate job that fails unless changes succeeded and another job ran. Attempt 1's finding f1 is fixed: the aggregate job and CHECK are now ci-passed, which no job in the workflow ci.install writes uses; a commit counts green only when every ci-passed check run succeeded, so a clash fails safe; and TestDocsCISinceGreen checks the name against workflowTemplate's jobs and checks the all-runs jq. rhiza ci install still writes only rhiza.yml (TestCIInstall). make self-check in a clean detached worktree at this commit reports one stale render, .rhiza/rendered/actor/claude-code-50966a9b.md and index.md. The cause is the earlier run evidence commit 9e331b53b for the leaf ci-model-only-since-green, which also fails at that commit, not this leaf's change. The land step's rhiza render must cover it.

## Links

### Out

- verifies [ci-since-green-for-projects](../capability/rhiza/integrity/ci-install/ci-since-green-for-projects.md)

## History

- `01M76XSFD0JZ9424EWBSQA3K8X` CI on a pull request runs only the model checks when nothing outside .rhiza/ changed since the pull request's last commit whose CI passed, so ship land's push and merges of main that bring only .rhiza/ changes cost a model check instead of a full CI round. Evidence: wait-landed is a full CI round in every run, 5.9m median (perf report rec 1).
