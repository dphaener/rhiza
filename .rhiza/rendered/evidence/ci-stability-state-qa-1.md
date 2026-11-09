# ci-stability-state-qa-1

evidence · `01M6YR95XEKAJZAHBK7G8N9ZK1` · actor [claude-code-ac99a081](../actor/claude-code-ac99a081.md)

## result

PASS

## statement

QA claude-code-ac99a081, binary built from abb2ffe1 (go build -trimpath), run by hand in a hermetic scratch repo outside ~/code (env -i, own HOME, bare origin, fake gh on PATH): rhiza health ci part: empty record 'ok ci stable'; two flaky runs 'warn ci unstable: flake in the last 20 CI runs'; with 3 later clean runs, rhiza daemon set rhiza.ship.ciWindow 3 gave stable and 5 gave unstable again; ciWindow 0 and abc each exit 2 'must be a whole number of at least 1'; a red main entry 'unstable: main's latest CI result is red (m1)' (also outliving the window), cleared to stable by a later green main entry.

## Links

### Out

- verifies [ci-stability-state](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-stability-state.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
