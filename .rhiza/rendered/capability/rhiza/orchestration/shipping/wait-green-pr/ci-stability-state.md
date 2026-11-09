# ci-stability-state

capability · `01M6YCRKH9VEHGSH0A29RKY1ZD` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/api -run TestCIStability shows: an empty record stable; a flake inside the window unstable, and stable once 20 later clean runs push it out; a red dispatched main result unstable, and stable once a later ship check --main records green, whether by a tree-equal merge or a green dispatched run; rhiza.ship.ciWindow=5 shrinking the window; an invalid window refused by rhiza daemon set

## statement

Rhiza judges a repo's CI unstable while its last completed CI runs in the stability record (20 by default, the project setting rhiza.ship.ciWindow) hold a flake, or while main's latest result in the record is red, and stable again once neither holds. Main's latest result is the one the latest rhiza ship check --main recorded: a merge whose tree equals its tested PR head's records that head's green result for the merge commit, and a tree-differing merge records the dispatched CI run's result, so a red main clears at the next merge that settles green

## Links

### Out

- depends-on [ci-stability-record](ci-stability-record.md)
- decided-by [ci-unstable-gate-rules](../../../../../decision/ci-unstable-gate-rules.md)

### In

- depends-on [health-ci-stability](../../daemon/daemon-status/health-command/health-ci-stability.md)
- depends-on [ship-local-gate](ship-local-gate.md)
- verifies [ci-stability-state-verified](../../../../../evidence/ci-stability-state-verified.md)
- verifies [ci-stability-state-qa-1](../../../../../evidence/ci-stability-state-qa-1.md)
- verifies [pr-403-stop-using-ci-as-a-debugger](../../../../../evidence/pr-403-stop-using-ci-as-a-debugger.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
