# ci-stability-state-verified

evidence · `01M6YJ0YMEXH0BSN3GMTCBNVA1` · actor [claude-code-681ac9bc](../actor/claude-code-681ac9bc.md)

## result

PASS

## statement

Verifier at 6a2f4a2a ran go test -trimpath ./internal/api -run 'TestCIStability|TestShipCheckStabilityRecord' -count=1 -v: PASS. TestCIStability covers an empty record (stable), a flake inside the 20-run window (unstable) and the same flake pushed out by 20 later clean runs (stable), window 5 shrinking and still holding, a red main result (unstable) cleared by a later green main record, a red main outliving the window, and a red PR head not counting as main. A ship check --main subtest takes a hand-written red main record to stable through a tree-equal merge. The green-dispatched-run path is covered by the pure verdict over a later green Main record together with TestShipCheckStabilityRecord/main_tree_differs, which records the dispatched run's green as a Main run. The ciWindow setting defaults to 20 and validates whole numbers of at least 1. A hand run of a scratch build in a hermetic repo: rhiza daemon set rhiza.ship.ciWindow 0 and abc each exit 2 with 'must be a whole number of at least 1', and 5 is accepted and listed as (git config, project). TestScripts/describe-golden passes, and go vet and gofmt are clean on internal/api. Reading the diff: ciVerdict checks the last window runs for flakes and main's latest run for red. ciStabilityVerdict, which reads the setting, has no test and no caller yet. TestScripts/daemon-settings fails here with 'a Rhiza repository already exists in /home/rhiza/code/rhiza', and it fails the same way at the parent commit, so the cause is the environment, not this leaf.

## Links

### Out

- verifies [ci-stability-state](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-stability-state.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
