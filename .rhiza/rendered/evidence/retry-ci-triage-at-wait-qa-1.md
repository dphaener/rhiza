# retry-ci-triage-at-wait-qa-1

evidence · `01M6YS5TKST2Z1Z2XM3ADN2AMS` · actor [claude-code-4adc28d2](../actor/claude-code-4adc28d2.md)

## result

PASS

## statement

QA at 91617709 with a built binary in a scratch Rhiza root and a fake gh: a ship-change workflow (wait-pr, pr-green, conflict-triage, ci-triage with no on.fail) run with every job failing ended 'ci-triage (visit 1): fail, then end'. With gh switched to green, rhiza run resume --foreground <run> --retry walked back through conflict-triage and pr-green and reported 'retried at wait-pr, after: step "ci-triage" ended with outcome "fail"'; status showed wait-pr (visit 2, retried): pass, then pr-green; pr-green (visit 2): pass, then signoff; signoff pass, then end; the run succeeded.

## Links

### Out

- verifies [retry-ci-triage-at-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/retry-ci-triage-at-wait.md)

## History

- `01M6YGFGQT0A9VRZ4AZ9X6VSMJ` Fix rhiza run resume --retry on a ship-change run that failed at ci-triage can never recover once the PR is green. Run 01M6XVTXGJAMQRBZKDZKY0Z01T (PR #393): pr-green failed because GitHub refused every CI job for billing (no test ran), so ci-triage ended 'exit 1, no failing test was recorded for the head: an infra or unparsable failure' and the run failed. After billing was fixed, one rerun turned PR #393 fully green, but 'rhiza run resume --retry' restarts at ci-triage, which re-reads the stale state and fails again in 1s (journal: retried ci-triage, visit 2, outcome fail). The retry should re-evaluate the PR first: when the head (or a new head) is green, continue at pr-green/rebase; an infra-only failure (no failing test) should send the run back to wait on CI rather than end it. Unstuck by hand: merged origin/main in the run worktree, pushed, then rhiza ship check / ship land / ship merge.
