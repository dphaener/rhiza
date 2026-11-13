# unstable-ci-guidance

capability · `01M6YCRMA48VMM2KTZ2Q92EY10` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath ./cmd/rhiza -run TestClaudeMDScopedTests finds AGENTS.md's scoped-tests rule unchanged plus the unstable-CI reproduce-locally rule; go test -trimpath ./internal/guide -run 'TestBuilderText|TestReviewText|TestSharedPhrasesInExactlyOneFile' finds the rule in the shared unstable-ci block, which the builder topic carries, and in the review text

## statement

This repo's AGENTS.md, the builder guide topic (through the shared unstable-ci block) and the review guide topic, so the rhiza-builder and rhiza-review skills, keep scoped tests as the default and add that, while CI is unstable or a red is not explained by the change, an agent reproduces the failing test locally and pushes only after it and the touched packages pass, never pushing or rerunning CI to see

## Links

### Out

- decided-by [no-ci-as-debugger](../../../../decision/no-ci-as-debugger.md)

### In

- verifies [unstable-ci-guidance-verified](../../../../evidence/unstable-ci-guidance-verified.md)
- verifies [unstable-ci-guidance-qa-1](../../../../evidence/unstable-ci-guidance-qa-1.md)
- verifies [pr-403-stop-using-ci-as-a-debugger](../../../../evidence/pr-403-stop-using-ci-as-a-debugger.md)
- verifies [unstable-ci-guidance-agents-md-qa-1](../../../../evidence/unstable-ci-guidance-agents-md-qa-1.md)
- verifies [pr-494-agents-md-entry-point](../../../../evidence/pr-494-agents-md-entry-point.md)
- verifies [pr-521-build-role-topics](../../../../evidence/pr-521-build-role-topics.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
- `01M76NWXV16NF94MBV6WYESRM0` AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md
- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
