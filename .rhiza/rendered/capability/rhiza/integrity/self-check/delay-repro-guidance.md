# delay-repro-guidance

capability · `01M6ZCM35BZXJ8X5SQ96639XZ4` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath ./internal/guide -run 'TestBuildText|TestReviewText' finds the delay-injection rule (inject a delay at the suspected race; never busy loops or parallel stress on the host) in both the build and review topics beside the unstable-CI rule; go test -trimpath ./cmd/rhiza -run TestClaudeMDScopedTests finds the same rule in CLAUDE.md beside the unstable-CI reproduce-locally rule, with that rule unchanged

## statement

The rhiza-build and rhiza-review guide topics (so the rhiza-build and rhiza-review skills) and this repo's CLAUDE.md tell builders, verifiers, checks and fixers to reproduce a timing flake by injecting a delay at the suspected race, and never by busy loops, CPU stress or parallel stress runs on the shared host

## Links

### Out

- decided-by [decision-delay-repro-scope](../../../../decision/decision-delay-repro-scope.md)

### In

- verifies [delay-repro-guidance-verified](../../../../evidence/delay-repro-guidance-verified.md)
- verifies [delay-repro-guidance-qa-1](../../../../evidence/delay-repro-guidance-qa-1.md)
- verifies [pr-411-builders-and-fixers-reproduce](../../../../evidence/pr-411-builders-and-fixers-reproduce.md)

## History

- `01M6ZCC27HQ7NMBTZ74D1CJFBV` Builders and fixers reproduce flaky tests by generating CPU load on the shared host. On 2026-11-09 two nebari runs did it: 01M6Z1RG ran 8 'bash' busy loops, and 01M6Z5P4 (session 42a0d9da) ran 14 'sh -c while :; do :; done' loops. Each held the 4-core VM at load 16–22 for 15–30+ min and slowed every other live run; I stopped them by PID. Load-based repro also barely works: 45/45 passed under load, while a 500 ms injected delay on the suspected request reproduced 10/10 (nebari PR #90). The build/fix guidance (rhiza guide build, the rhiza-build and rhiza-review skills) should say: reproduce timing flakes by injecting a delay at the suspected race, never by busy loops or parallel stress on the host. Optionally the watchdog flags a run's descendant processes that sit at 100% CPU for minutes.
