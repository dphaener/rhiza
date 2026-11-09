# pr-402-fix-rhiza-run-resume-retry-01m6ygfgqt0a

evidence · `01M6YYHPCNFFVDPZ42X291AR2J` · actor [dphaener](../actor/dphaener.md)

## head

eba2e25cae48bbbf41cfe38293cb454932b008b1

## pr

https://github.com/dphaener/rhiza/pull/402

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #402 (https://github.com/dphaener/rhiza/pull/402) passed CI at eba2e25

## status

verified

## Links

### Out

- verifies [retry-ci-triage-at-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/retry-ci-triage-at-wait.md)
- verifies [ci-triage-infra-waits](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-triage-infra-waits.md)
- verifies [ship-change-ci-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-ci-fix-step.md)
- verifies [retry-reopens-at-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/retry-reopens-at-wait.md)
- verifies [ship-change-conflict-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-conflict-fix-step.md)

## History

- `01M6YGFGQT0A9VRZ4AZ9X6VSMJ` Fix rhiza run resume --retry on a ship-change run that failed at ci-triage can never recover once the PR is green. Run 01M6XVTXGJAMQRBZKDZKY0Z01T (PR #393): pr-green failed because GitHub refused every CI job for billing (no test ran), so ci-triage ended 'exit 1, no failing test was recorded for the head: an infra or unparsable failure' and the run failed. After billing was fixed, one rerun turned PR #393 fully green, but 'rhiza run resume --retry' restarts at ci-triage, which re-reads the stale state and fails again in 1s (journal: retried ci-triage, visit 2, outcome fail). The retry should re-evaluate the PR first: when the head (or a new head) is green, continue at pr-green/rebase; an infra-only failure (no failing test) should send the run back to wait on CI rather than end it. Unstuck by hand: merged origin/main in the run worktree, pushed, then rhiza ship check / ship land / ship merge.
