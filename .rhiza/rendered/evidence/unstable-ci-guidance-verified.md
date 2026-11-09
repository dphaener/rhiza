# unstable-ci-guidance-verified

evidence · `01M6YH8QN3CX82D3R5F4GGP6VH` · actor [claude-code-cbb0802a](../actor/claude-code-cbb0802a.md)

## result

PASS

## statement

Verifier claude-code-cbb0802a at 683b17b7: go test -trimpath -count=1 ./cmd/rhiza -run TestClaudeMDScopedTests PASS (CLAUDE.md keeps the scoped-tests rule and adds the unstable-CI paragraph: while CI is unstable or a red is not explained by your change, do not push or rerun CI to see; reproduce the failing test locally; push only after it and the touched packages pass locally). go test -trimpath -count=1 ./internal/guide -run 'TestBuildText|TestReviewText' PASS, with the same rule in build.md and review.md beside the unchanged scoped-tests paragraph. Attempt 1's finding is fixed: go test -trimpath -count=1 ./internal/guide/... passes, TestCaps included. go vet and gofmt are clean. Read the diff: review.md's tightening keeps every tested phrase, and the heredoc examples are still there. TestScripts/guide and guide-install fail identically at the parent 6f44a0e5 ('a Rhiza repository already exists in /home/rhiza/code/rhiza'), an environment issue under ~/code, not this leaf.

## Links

### Out

- verifies [unstable-ci-guidance](../capability/rhiza/integrity/self-check/unstable-ci-guidance.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
