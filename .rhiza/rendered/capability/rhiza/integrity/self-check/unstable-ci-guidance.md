# unstable-ci-guidance

capability · `01M6YCRMA48VMM2KTZ2Q92EY10` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath ./cmd/rhiza -run TestClaudeMDScopedTests finds CLAUDE.md's scoped-tests rule unchanged plus the unstable-CI reproduce-locally rule; go test -trimpath ./internal/guide -run 'TestBuildText|TestReviewText' finds the same two rules in the rhiza-build and rhiza-review guide texts

## statement

This repo's CLAUDE.md and the rhiza-build and rhiza-review guide topics keep scoped tests as the default and add that, while CI is unstable or a red is not explained by the change, an agent reproduces the failing test locally and pushes only after it and the touched packages pass, never pushing or rerunning CI to see

## Links

### Out

- decided-by [no-ci-as-debugger](../../../../decision/no-ci-as-debugger.md)

### In

- verifies [unstable-ci-guidance-verified](../../../../evidence/unstable-ci-guidance-verified.md)
- verifies [unstable-ci-guidance-qa-1](../../../../evidence/unstable-ci-guidance-qa-1.md)
- verifies [pr-403-stop-using-ci-as-a-debugger](../../../../evidence/pr-403-stop-using-ci-as-a-debugger.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
