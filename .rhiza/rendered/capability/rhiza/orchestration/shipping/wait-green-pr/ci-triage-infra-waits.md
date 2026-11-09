# ci-triage-infra-waits

capability · `01M6YGXMGVZQBYW9NX74P1TQ4W` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

Go tests in internal/api show a failing ship check recording red_head and red_checks (IDs and attempts) and no Head/Checks, ship check --settled --changed not settling on the unchanged settled red head and naming the infra failure in its output, settling on a new head or on a newer attempt of a check run once it finishes, and acting as --settled with no baseline; a Go test in internal/run shows a wait whose until never holds timing out with a detail that holds the until command's last output; testscript run-ship-rerun-fail (or a new run-ship-infra-wait with a fake gh) drives a ship-change run whose CI fails with no failing test into wait-ci-rerun, then a rerun turns the head green and the run passes pr-green to signoff; run with go test -trimpath ./internal/api/... ./internal/run/... ./internal/workflow/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship-rerun-fail|run-ship-infra-wait|run-ship-ci-fix)' only

## statement

When ship-change's ci-triage finds no failing test for a red head (an infra or unparsable failure), the run goes to a new wait step wait-ci-rerun instead of ending. The baseline it waits past is recorded by the failing ship check (pr-green or rebased-green, rhiza ship check without --settled): on a red head it now records in the state file the red head SHA and each of its check runs' ID and attempt (red_head, red_checks), while still recording no green Head/Checks; ci-triage (--test-failures) stays read-only. wait-ci-rerun's until, rhiza ship check --settled --changed, settles only once CI moved past that baseline (a head SHA other than red_head, or a check run ID or attempt not in red_checks) and that new CI finished; with no baseline recorded it behaves as --settled. While not settled its output names the infra failure (the red head and its failed checks with no failing test). The wait runs up to 6h with max_visits 3 and its on.pass goes to pr-green, which re-checks the head (a run back from after the rebase passes signoff again). Its timeout ends the run with a line naming the infra failure: the run engine's wait keeps the until command's last output and adds it to a wait's timeout detail (did not hold within <timeout>: <last output>)

## Links

### Out

- decided-by [ci-triage-infra-waits-for-new-ci](../../../../../decision/ci-triage-infra-waits-for-new-ci.md)
- derived-from [resume-retry-ci-triage-plan](../../../../../source/resume-retry-ci-triage-plan.md)

### In

- verifies [ci-triage-infra-waits-verified](../../../../../evidence/ci-triage-infra-waits-verified.md)
- verifies [ci-triage-infra-waits-qa-1](../../../../../evidence/ci-triage-infra-waits-qa-1.md)
- verifies [ci-triage-infra-waits-verified-2](../../../../../evidence/ci-triage-infra-waits-verified-2.md)
- verifies [pr-402-fix-rhiza-run-resume-retry-01m6ygfgqt0a](../../../../../evidence/pr-402-fix-rhiza-run-resume-retry-01m6ygfgqt0a.md)

## History

- `01M6YGFGQT0A9VRZ4AZ9X6VSMJ` Fix rhiza run resume --retry on a ship-change run that failed at ci-triage can never recover once the PR is green. Run 01M6XVTXGJAMQRBZKDZKY0Z01T (PR #393): pr-green failed because GitHub refused every CI job for billing (no test ran), so ci-triage ended 'exit 1, no failing test was recorded for the head: an infra or unparsable failure' and the run failed. After billing was fixed, one rerun turned PR #393 fully green, but 'rhiza run resume --retry' restarts at ci-triage, which re-reads the stale state and fails again in 1s (journal: retried ci-triage, visit 2, outcome fail). The retry should re-evaluate the PR first: when the head (or a new head) is green, continue at pr-green/rebase; an infra-only failure (no failing test) should send the run back to wait on CI rather than end it. Unstuck by hand: merged origin/main in the run worktree, pushed, then rhiza ship check / ship land / ship merge.
