# resume-retry-ci-triage-plan

source · `01M6YGX0K9MV6C076HCG8CEANC`

## quotes

```json
[]
```

## raw_sha256

ae3779ee4a735ab3b2e2b904eb1a3c017737297039230a5ae108fcbafbcc8eb3

## session

bc50c2e7-c330-4961-babf-01a1bd082164

## summary

Autonomous plan run for the bug where `rhiza run resume --retry` of ship-change run 01M6XVTXGJAM (PR #393) restarted at ci-triage and failed again on stale state, after a billing outage had failed every CI job with no failing test. The planner found the cause in the code: Engine.Retry reopens at State.FailedStep, and ci-triage (`rhiza ship check --test-failures`) has no on.fail, so an infra-only red head ends the run. The accepted change set 01M6XX4R8 (retry-reopens-at-wait) covers retries of green steps entered from a wait, but not ci-triage, which pr-green's or rebased-green's on.fail reaches. The change set was placed under milestone-25, like the other ship-change fixes. Two questions went to the arbiter. It decided retry-ci-triage-reopens-at-wait: a retry follows the journal back from ci-triage to the check that routed into it and reopens at that check's wait. It also decided ci-triage-infra-waits-for-new-ci: ci-triage's on.fail goes to a new wait step that settles only once CI moved past the judged state, with a 6h timeout and 3 visits, then re-checks at pr-green. The change set depends on 01M6XX4R8, whose engine rule it extends. No owner words: the run was autonomous.

## title

Resume-retry at ci-triage plan, autonomous

## Links

### In

- derived-from [retry-ci-triage-reopens-at-wait](../decision/retry-ci-triage-reopens-at-wait.md)
- derived-from [ci-triage-infra-waits-for-new-ci](../decision/ci-triage-infra-waits-for-new-ci.md)
- derived-from [retry-ci-triage-at-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/retry-ci-triage-at-wait.md)
- derived-from [ci-triage-infra-waits](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-triage-infra-waits.md)

## History

- `01M6YGFGQT0A9VRZ4AZ9X6VSMJ` Fix rhiza run resume --retry on a ship-change run that failed at ci-triage can never recover once the PR is green. Run 01M6XVTXGJAMQRBZKDZKY0Z01T (PR #393): pr-green failed because GitHub refused every CI job for billing (no test ran), so ci-triage ended 'exit 1, no failing test was recorded for the head: an infra or unparsable failure' and the run failed. After billing was fixed, one rerun turned PR #393 fully green, but 'rhiza run resume --retry' restarts at ci-triage, which re-reads the stale state and fails again in 1s (journal: retried ci-triage, visit 2, outcome fail). The retry should re-evaluate the PR first: when the head (or a new head) is green, continue at pr-green/rebase; an infra-only failure (no failing test) should send the run back to wait on CI rather than end it. Unstuck by hand: merged origin/main in the run worktree, pushed, then rhiza ship check / ship land / ship merge.
