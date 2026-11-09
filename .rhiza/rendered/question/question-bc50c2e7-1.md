# question-bc50c2e7-1

question · `01M6YGMPABACEVZCB6Y4SD5MG8`

## session

bc50c2e7-c330-4961-babf-01a1bd082164

## statement

A ship-change run that failed at ci-triage (reached from pr-green's or rebased-green's on.fail) is retried with `rhiza run resume --retry`. Today the retry restarts at ci-triage, which reads the stale state and fails again. The accepted change set 01M6XX4R8 (retry-reopens-at-wait, not built yet) will reopen a retry at the wait step whose on.pass is the failed step, but ci-triage has no such wait. Where should a retry of a run that failed at ci-triage restart?

## Links

### In

- answers [retry-ci-triage-reopens-at-wait](../decision/retry-ci-triage-reopens-at-wait.md)

## History

- `01M6YGFGQT0A9VRZ4AZ9X6VSMJ` Fix rhiza run resume --retry on a ship-change run that failed at ci-triage can never recover once the PR is green. Run 01M6XVTXGJAMQRBZKDZKY0Z01T (PR #393): pr-green failed because GitHub refused every CI job for billing (no test ran), so ci-triage ended 'exit 1, no failing test was recorded for the head: an infra or unparsable failure' and the run failed. After billing was fixed, one rerun turned PR #393 fully green, but 'rhiza run resume --retry' restarts at ci-triage, which re-reads the stale state and fails again in 1s (journal: retried ci-triage, visit 2, outcome fail). The retry should re-evaluate the PR first: when the head (or a new head) is green, continue at pr-green/rebase; an infra-only failure (no failing test) should send the run back to wait on CI rather than end it. Unstuck by hand: merged origin/main in the run worktree, pushed, then rhiza ship check / ship land / ship merge.
