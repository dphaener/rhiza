# pr-403-stop-using-ci-as-a-debugger

evidence · `01M6Z2W564GGD7H4NNX688PXRF` · actor [dphaener](../actor/dphaener.md)

## head

5faa99c27567a081a91b739eb00b317f0289b96c

## pr

https://github.com/dphaener/rhiza/pull/403

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #403 (https://github.com/dphaener/rhiza/pull/403) passed CI at 5faa99c

## status

verified

## Links

### Out

- verifies [ci-stability-record](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-stability-record.md)
- verifies [ci-stability-state](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-stability-state.md)
- verifies [health-ci-stability](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-ci-stability.md)
- verifies [ship-local-gate](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-local-gate.md)
- verifies [ci-fix-reproduces-locally](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-fix-reproduces-locally.md)
- verifies [flaky-test-bugs](../capability/rhiza/orchestration/daemon/daemon-watchdog/flaky-test-bugs.md)
- verifies [unstable-ci-guidance](../capability/rhiza/integrity/self-check/unstable-ci-guidance.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
