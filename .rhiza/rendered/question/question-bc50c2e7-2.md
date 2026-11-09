# question-bc50c2e7-2

question · `01M6YGREJVNGXG1XFHEZPWXBFD`

## session

bc50c2e7-c330-4961-babf-01a1bd082164

## statement

In ship-change, ci-triage (`rhiza ship check --test-failures`) has no on.fail. So a red PR head whose CI logs name no failing test (an infra or unparsable failure, e.g. GitHub refusing every job for billing, after `rhiza ship check` has used its rerun budget) ends the run. The change set's intent says such a failure should send the run back to wait on CI instead. Sending it straight to wait-pr would loop at once, because the head is already settled red, until max_visits ends the run. What should ci-triage do on an infra-only failure?

## Links

### In

- answers [ci-triage-infra-waits-for-new-ci](../decision/ci-triage-infra-waits-for-new-ci.md)

## History

- `01M6YGFGQT0A9VRZ4AZ9X6VSMJ` Fix rhiza run resume --retry on a ship-change run that failed at ci-triage can never recover once the PR is green. Run 01M6XVTXGJAMQRBZKDZKY0Z01T (PR #393): pr-green failed because GitHub refused every CI job for billing (no test ran), so ci-triage ended 'exit 1, no failing test was recorded for the head: an infra or unparsable failure' and the run failed. After billing was fixed, one rerun turned PR #393 fully green, but 'rhiza run resume --retry' restarts at ci-triage, which re-reads the stale state and fails again in 1s (journal: retried ci-triage, visit 2, outcome fail). The retry should re-evaluate the PR first: when the head (or a new head) is green, continue at pr-green/rebase; an infra-only failure (no failing test) should send the run back to wait on CI rather than end it. Unstuck by hand: merged origin/main in the run worktree, pushed, then rhiza ship check / ship land / ship merge.
