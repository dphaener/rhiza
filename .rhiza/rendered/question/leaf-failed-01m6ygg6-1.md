# leaf-failed-01m6ygg6-1

question · `01M6YY6J9N4W6WEN6M5J07TN8B`

## statement

Review check quality still fails after 3 cycles: finding f1 (major at internal/workflow/builtin/ship-change.toml:193) on leaf ship-local-gate: A gate failing after land still cannot reach a working fix. rebase-merge's new on.fail (rebase-gate-failed -> gate-fix), like gate-rebase-merge's on.fail = gat…. Its last cycle: cycle 3: confirmed f1 (major). How should it be fixed, or should leaf ship-local-gate change?

## Links

### Out

- blocks [ship-local-gate](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-local-gate.md)

### In

- answers [ship-local-gate-post-land-fix-loop](../decision/ship-local-gate-post-land-fix-loop.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
