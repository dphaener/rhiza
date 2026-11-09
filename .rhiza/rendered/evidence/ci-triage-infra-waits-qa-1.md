# ci-triage-infra-waits-qa-1

evidence · `01M6YS5THPHYYNMS22AQVFXVEW` · actor [claude-code-4adc28d2](../actor/claude-code-4adc28d2.md)

## result

PASS

## statement

QA at 91617709 with a built binary in a scratch Rhiza root (outside ~/code, RHIZA_* unset) and a fake gh using real Actions detailsUrls (/actions/runs/42/job/N) that fails every job (log: billing message, no test), reruns turning job 9 into job 10. rhiza run start --foreground on ship-change (the built-in's middle, short timeouts): pr-green failed after 1 rerun, conflict-triage fail, ci-triage 'no failing test ... infra or unparsable failure' then wait-ci-rerun. The ship state held red_head=<head> and red_checks=[{test,id 10}] with no head/checks. rhiza ship check --settled --changed in the run worktree: exit 1, 'not settled: CI on head ab3bb9b failed (test) with no failing test: an infra failure; waiting for a new head or a rerun'. Switching gh to job 11 IN_PROGRESS: 'not settled: 0 of 1 checks finished'; to job 11 SUCCESS: wait-ci-rerun passed, pr-green (visit 2) pass, signoff, run succeeded, and the state then held head/checks green with red_head cleared. A second run left red with wait-ci-rerun timeout 3s ended 'wait-ci-rerun (visit 1): timeout, then end: did not hold within 3s: PR #7 at 5cb34b7: not settled: CI on head 5cb34b7 failed (test) with no failing test: an infra failure ...'. Pushing a new head there made --settled --changed settle (exit 0) once its CI finished. --changed without --settled, or with --main, is refused (exit 2). rhiza workflow check ship-change: valid; describe shows ci-triage on.fail -> wait-ci-rerun -> on.pass pr-green.

## Links

### Out

- verifies [ci-triage-infra-waits](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-triage-infra-waits.md)

## History

- `01M6YGFGQT0A9VRZ4AZ9X6VSMJ` Fix rhiza run resume --retry on a ship-change run that failed at ci-triage can never recover once the PR is green. Run 01M6XVTXGJAMQRBZKDZKY0Z01T (PR #393): pr-green failed because GitHub refused every CI job for billing (no test ran), so ci-triage ended 'exit 1, no failing test was recorded for the head: an infra or unparsable failure' and the run failed. After billing was fixed, one rerun turned PR #393 fully green, but 'rhiza run resume --retry' restarts at ci-triage, which re-reads the stale state and fails again in 1s (journal: retried ci-triage, visit 2, outcome fail). The retry should re-evaluate the PR first: when the head (or a new head) is green, continue at pr-green/rebase; an infra-only failure (no failing test) should send the run back to wait on CI rather than end it. Unstuck by hand: merged origin/main in the run worktree, pushed, then rhiza ship check / ship land / ship merge.
