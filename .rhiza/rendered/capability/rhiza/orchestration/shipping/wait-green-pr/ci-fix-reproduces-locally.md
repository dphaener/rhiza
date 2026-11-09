# ci-fix-reproduces-locally

capability · `01M6YCRM0GKZ3W3J2QDD87GWHQ` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/workflow -run TestShipChangeCIFixPrompt asserts the ci-fix prompt carries the reproduce-locally, fix-or-quarantine-with-a-bug and pass-locally-before-push rules and names no Go, make or testscript command; go test -trimpath ./internal/api -run 'TestShipCheckRerunBudget|TestWatchdogCIRerun' shows neither ship check nor the watchdog rerunning a head's failed jobs beyond its rerun budget: one rerun at the default budget of 1, and exactly N at --reruns N / ciReruns N

## statement

A head still red once its rerun budget is spent (ship check --reruns and the watchdog's ciReruns, both default 1) is never rerun again: ship-change's ci-fix fixer runs each failing test locally, repeatedly until it reproduces or plainly passes, fixes it or quarantines it with a bug request filed and named at the skip, and commits only once those tests and the packages it touched pass locally

## Links

### Out

- decided-by [ci-unstable-gate-rules](../../../../../decision/ci-unstable-gate-rules.md)
- decided-by [no-ci-as-debugger](../../../../../decision/no-ci-as-debugger.md)

### In

- verifies [ci-fix-reproduces-locally-verified](../../../../../evidence/ci-fix-reproduces-locally-verified.md)
- verifies [ci-fix-reproduces-locally-qa-1](../../../../../evidence/ci-fix-reproduces-locally-qa-1.md)
- verifies [pr-403-stop-using-ci-as-a-debugger](../../../../../evidence/pr-403-stop-using-ci-as-a-debugger.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
