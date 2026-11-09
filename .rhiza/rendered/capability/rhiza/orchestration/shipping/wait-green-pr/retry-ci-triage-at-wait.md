# retry-ci-triage-at-wait

capability · `01M6YGXM7SRCDR03S5J1K60QVR` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

A Go test in internal/run retries a run failed at a step reached from a command step's on.fail, itself entered from a wait's on.pass, and shows it reopening at that wait, choosing between two such predecessors from the journal; the same for a run whose journal already holds two retried entries and three visits of the failed step (retried twice there), which still reopens at the wait; a run failed at a step with neither still retries at that step; testscript run-resume-retry (or a new run-ship-retry-ci-triage with a fake gh) retries a ship-change run failed at ci-triage after the PR head turned green and it passes wait-pr and pr-green on to signoff; run with go test -trimpath ./internal/run/... ./internal/api/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-resume-retry|run-ship-retry)' only

## statement

rhiza run resume --retry on a run that failed at a step with no wait whose on.pass it is (ship-change's ci-triage) walks the journal back from the failed visit, skipping retried entries and the failed step's own earlier visits (each earlier retry adds a retried entry and a fresh visit of that step), to the last other step visited: the check whose on.fail routed into it (pr-green or rebased-green). It reopens the run at that check's wait step (wait-pr or wait-rebased), so the retry waits for CI on the current or a new head and re-runs the check: green goes on to signoff or land, red goes to ci-triage again with fresh state. A run retried at ci-triage once or more before still reopens at that wait. A failed step whose last other step is not a command step routing into it by on.fail still retries at itself

## Links

### Out

- decided-by [retry-ci-triage-reopens-at-wait](../../../../../decision/retry-ci-triage-reopens-at-wait.md)
- derived-from [resume-retry-ci-triage-plan](../../../../../source/resume-retry-ci-triage-plan.md)

### In

- verifies [retry-ci-triage-at-wait-verified](../../../../../evidence/retry-ci-triage-at-wait-verified.md)
- verifies [retry-ci-triage-at-wait-qa-1](../../../../../evidence/retry-ci-triage-at-wait-qa-1.md)
- verifies [pr-402-fix-rhiza-run-resume-retry-01m6ygfgqt0a](../../../../../evidence/pr-402-fix-rhiza-run-resume-retry-01m6ygfgqt0a.md)

## History

- `01M6YGFGQT0A9VRZ4AZ9X6VSMJ` Fix rhiza run resume --retry on a ship-change run that failed at ci-triage can never recover once the PR is green. Run 01M6XVTXGJAMQRBZKDZKY0Z01T (PR #393): pr-green failed because GitHub refused every CI job for billing (no test ran), so ci-triage ended 'exit 1, no failing test was recorded for the head: an infra or unparsable failure' and the run failed. After billing was fixed, one rerun turned PR #393 fully green, but 'rhiza run resume --retry' restarts at ci-triage, which re-reads the stale state and fails again in 1s (journal: retried ci-triage, visit 2, outcome fail). The retry should re-evaluate the PR first: when the head (or a new head) is green, continue at pr-green/rebase; an infra-only failure (no failing test) should send the run back to wait on CI rather than end it. Unstuck by hand: merged origin/main in the run worktree, pushed, then rhiza ship check / ship land / ship merge.
