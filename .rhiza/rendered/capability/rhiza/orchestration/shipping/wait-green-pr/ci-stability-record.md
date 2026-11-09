# ci-stability-record

capability · `01M6YCRKDX9JJQ25FRHFJND2ZK` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/api -run TestShipCheckStabilityRecord shows a green first attempt recorded as a clean run, a red head that passes on its rerun recorded with one flake per failing test, a head still red after its rerun recorded red with no flake, a tree-equal ship check --main recording a green main result for the merge commit (the tested head's result, no CI run), and a tree-differing ship check --main recording the dispatched CI run's result on main

## statement

rhiza ship check records each CI run it settles, on a run's PR head or on main, in the repo's CI stability record (head, PR or main, result, time), and each failing test whose job then passed on a rerun of the same head as a flake against that test

## Links

### Out

- decided-by [ci-unstable-gate-rules](../../../../../decision/ci-unstable-gate-rules.md)

### In

- depends-on [ci-stability-state](ci-stability-state.md)
- depends-on [flaky-test-bugs](../../daemon/daemon-watchdog/flaky-test-bugs.md)
- verifies [ci-stability-record-verified](../../../../../evidence/ci-stability-record-verified.md)
- verifies [ci-stability-record-qa-1](../../../../../evidence/ci-stability-record-qa-1.md)
- verifies [pr-403-stop-using-ci-as-a-debugger](../../../../../evidence/pr-403-stop-using-ci-as-a-debugger.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
