# flaky-test-bugs

capability · `01M6YCRM4F76QFJRC071CY3DYZ` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath ./internal/api -run TestFlakyTestBugs shows one bug request queued per newly flaky test with the test name and runs, a second flake of the same test queuing none while the first request is open, two flaky tests queuing two, and fileBugs = false queuing none

## statement

With the watchdog's fileBugs on, each test the CI stability record holds a flake for queues one bug request naming the test, its flaky runs and their heads, asking to reproduce it locally and fix or quarantine it, and no second one while an open request names the same test

## Links

### Out

- depends-on [ci-stability-record](../../shipping/wait-green-pr/ci-stability-record.md)
- decided-by [ci-unstable-gate-rules](../../../../../decision/ci-unstable-gate-rules.md)

### In

- verifies [flaky-test-bugs-verified](../../../../../evidence/flaky-test-bugs-verified.md)
- verifies [flaky-test-bugs-qa-1](../../../../../evidence/flaky-test-bugs-qa-1.md)
- verifies [pr-403-stop-using-ci-as-a-debugger](../../../../../evidence/pr-403-stop-using-ci-as-a-debugger.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
