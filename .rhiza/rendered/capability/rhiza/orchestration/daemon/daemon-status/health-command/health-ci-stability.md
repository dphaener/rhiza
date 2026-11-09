# health-ci-stability

capability · `01M6YCRKNPNZY0F0BTDCXQAP74` · parent [health-command](../health-command.md)

## acceptance

go test -trimpath ./internal/api -run TestHealthCI shows the part ok on a clean record, warn with each flaky test and its count on a flake, warn naming main red, and the no-local-gate warning only when unstable with rhiza.ship.localGate unset; go test -trimpath ./cmd/rhiza -run TestScripts/health-ci prints the part in rhiza health and health --json

## statement

rhiza health shows a CI part: stable or unstable with the reason (flakes in the window, main red), each flaky test with its flake count, and a warning when CI is unstable and no local gate (rhiza.ship.localGate) is set

## Links

### Out

- depends-on [ci-stability-state](../../../shipping/wait-green-pr/ci-stability-state.md)
- depends-on [ship-local-gate](../../../shipping/wait-green-pr/ship-local-gate.md)
- decided-by [ci-unstable-gate-rules](../../../../../../decision/ci-unstable-gate-rules.md)

### In

- verifies [health-ci-stability-verified](../../../../../../evidence/health-ci-stability-verified.md)
- verifies [health-ci-stability-qa-1](../../../../../../evidence/health-ci-stability-qa-1.md)
- verifies [pr-403-stop-using-ci-as-a-debugger](../../../../../../evidence/pr-403-stop-using-ci-as-a-debugger.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
