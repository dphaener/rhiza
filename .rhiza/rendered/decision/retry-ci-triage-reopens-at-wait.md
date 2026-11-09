# retry-ci-triage-reopens-at-wait

decision · `01M6YGPS3159C8B71XAJ6XYMRC`

## reasoning

In internal/workflow/builtin/ship-change.toml, ci-triage has two entry points, pr-green and rebased-green, both through on.fail, and no wait step feeds it directly. So 01M6XX4R8's rule finds nothing. Restarting at ci-triage itself is the current bug: it re-reads the head's stale ship-check state and fails in 1s (run 01M6XVTXGJAMQRBZKDZKY0Z01T, PR #393). Restarting at pr-green or rebased-green directly would check CI before it settles after a manual rerun or push. The preceding wait (`rhiza ship check --settled`) followed by its check re-evaluates the PR exactly as the change set's intent asks: a green head (or new head) continues, a still-red head gets a fresh triage. The journal records which branch led to ci-triage, so the choice is deterministic and keeps the run on the same path (pre-signoff vs post-rebase), so a post-signoff run is not sent back through signoff. One rule, built on 01M6XX4R8 instead of a special case per step, keeps run resume general.

## statement

A `rhiza run resume --retry` of a ship-change run that failed at ci-triage restarts at the wait step that led to the check which routed into ci-triage: wait-pr when ci-triage was reached from pr-green's on.fail, wait-rebased when reached from rebased-green's on.fail. Which one is read from the run's journal (the step visited just before ci-triage), not guessed. The run then waits for CI to settle on the current head (or a new one) and re-runs that check: green continues to signoff or land; red goes to ci-triage again with fresh state. This extends 01M6XX4R8's rule (reopen at the wait whose on.pass is the failed step) one hop: when the failed step has no such wait, follow the on.fail edge back to the step that routed into it and reopen at that step's wait.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-bc50c2e7-1](../question/question-bc50c2e7-1.md)
- derived-from [resume-retry-ci-triage-plan](../source/resume-retry-ci-triage-plan.md)

### In

- decided-by [retry-ci-triage-at-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/retry-ci-triage-at-wait.md)

## History

- `01M6YGFGQT0A9VRZ4AZ9X6VSMJ` Fix rhiza run resume --retry on a ship-change run that failed at ci-triage can never recover once the PR is green. Run 01M6XVTXGJAMQRBZKDZKY0Z01T (PR #393): pr-green failed because GitHub refused every CI job for billing (no test ran), so ci-triage ended 'exit 1, no failing test was recorded for the head: an infra or unparsable failure' and the run failed. After billing was fixed, one rerun turned PR #393 fully green, but 'rhiza run resume --retry' restarts at ci-triage, which re-reads the stale state and fails again in 1s (journal: retried ci-triage, visit 2, outcome fail). The retry should re-evaluate the PR first: when the head (or a new head) is green, continue at pr-green/rebase; an infra-only failure (no failing test) should send the run back to wait on CI rather than end it. Unstuck by hand: merged origin/main in the run worktree, pushed, then rhiza ship check / ship land / ship merge.
