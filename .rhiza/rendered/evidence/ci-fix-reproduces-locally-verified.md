# ci-fix-reproduces-locally-verified

evidence · `01M6YHGW1KNZAC6ZSDHMWBD4NK` · actor [claude-code-584bd6c6](../actor/claude-code-584bd6c6.md)

## result

PASS

## statement

Verifier at d93f1b12: go test -trimpath -count=1 ./internal/workflow -run TestShipChangeCIFixPrompt passes; it asserts the ship-change ci-fix prompt carries reproduce-locally (repeat until it fails again or plainly passes, never push to see), fix or quarantine with a bug request filed and named at the skip, and commit only once the tests and touched packages pass locally, and that it names no go test/go vet/make/testscript/TestScripts/-trimpath. go test -trimpath -count=1 ./internal/api -run 'TestShipCheckRerunBudget|TestWatchdogCIRerun' passes (3 reruns ok; one earlier run tripped rootguard on a concurrent write to the enclosing root's state folder with every test PASS, not reproduced). TestWatchdogCIRerun pins one rerun at ciReruns 1; new TestWatchdogCIRerunBudgetN shows reruns exactly while attempt <= N for N=1..3 and none above; TestShipCheckRerunBudget shows --reruns 2 reruns twice then settles red, resetting on a new head (TestShipCheckRerunDefaultOnce, also passing, pins the default of 1). Diff read: prompt change in internal/workflow/builtin/ship-change.toml plus the two tests and the builder's actor record.

## Links

### Out

- verifies [ci-fix-reproduces-locally](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-fix-reproduces-locally.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
