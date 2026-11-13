# delay-repro-guidance

capability · `01M6ZCM35BZXJ8X5SQ96639XZ4` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath ./internal/guide/... passes, with TestSharedPhrasesInExactlyOneFile finding the delay-injection rule (inject a delay at the suspected race; never busy loops or parallel stress on the host) in the shared unstable-ci block beside the unstable-CI rule, TestBuilderText, TestFixerText, TestCIFixText and TestGateFixText asserting the builder, fixer, ci-fix and gate-fix topics carry that block, and TestReviewRoleText finding the rule in the expanded review role topics; go test -trimpath ./cmd/rhiza -run TestClaudeMDScopedTests finds the same rule in AGENTS.md beside the unstable-CI reproduce-locally rule, with that rule unchanged

## statement

The builder guide topic and the fixer, ci-fix and gate-fix review role topics (each through the shared unstable-ci block, beside the unstable-CI rule), so the rhiza-builder, rhiza-fixer, rhiza-ci-fix and rhiza-gate-fix skills, and this repo's AGENTS.md tell builders and those three fixer roles to reproduce a timing flake by injecting a delay at the suspected race, and never by busy loops, CPU stress or parallel stress runs on the shared host; the review index does not hold the rule

## Links

### Out

- decided-by [decision-delay-repro-scope](../../../../decision/decision-delay-repro-scope.md)

### In

- blocks [review-01m79a41352xbz4fc2fjt5903j-r3-f8](../../../../question/review-01m79a41352xbz4fc2fjt5903j-r3-f8.md)
- verifies [delay-repro-guidance-verified](../../../../evidence/delay-repro-guidance-verified.md)
- verifies [delay-repro-guidance-qa-1](../../../../evidence/delay-repro-guidance-qa-1.md)
- verifies [pr-411-builders-and-fixers-reproduce](../../../../evidence/pr-411-builders-and-fixers-reproduce.md)
- verifies [delay-repro-guidance-agents-md-qa-1](../../../../evidence/delay-repro-guidance-agents-md-qa-1.md)
- verifies [pr-494-agents-md-entry-point](../../../../evidence/pr-494-agents-md-entry-point.md)
- verifies [pr-521-build-role-topics](../../../../evidence/pr-521-build-role-topics.md)
- verifies [pr-537-review-role-topics](../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M6ZCC27HQ7NMBTZ74D1CJFBV` Builders and fixers reproduce flaky tests by generating CPU load on the shared host. On 2026-11-09 two nebari runs did it: 01M6Z1RG ran 8 'bash' busy loops, and 01M6Z5P4 (session 42a0d9da) ran 14 'sh -c while :; do :; done' loops. Each held the 4-core VM at load 16–22 for 15–30+ min and slowed every other live run; I stopped them by PID. Load-based repro also barely works: 45/45 passed under load, while a 500 ms injected delay on the suspected request reproduced 10/10 (nebari PR #90). The build/fix guidance (rhiza guide build, the rhiza-build and rhiza-review skills) should say: reproduce timing flakes by injecting a delay at the suspected race, never by busy loops or parallel stress on the host. Optionally the watchdog flags a run's descendant processes that sit at 100% CPU for minutes.
- `01M76NWXV16NF94MBV6WYESRM0` AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md
- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
