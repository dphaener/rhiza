# delay-repro-guidance-verified

evidence · `01M6ZG631E7PGQCGSJ1KZ60G2T` · actor [claude-code-625da618](../actor/claude-code-625da618.md)

## result

PASS

## statement

At 06ceae92: go test -count=1 -trimpath ./internal/guide -run 'TestBuildText|TestReviewText' PASS (and the whole guide package incl. token caps PASS); the new delay-repro groups assert 'injecting a delay at the suspected race' and 'never by busy loops or parallel stress on the host', and the diff puts that sentence in build.md and review.md directly after the unstable-CI reproduce-locally rule. go test -count=1 -trimpath ./cmd/rhiza -run TestClaudeMDScopedTests PASS, asserting 'inject a delay at the suspected race' and 'Never use busy loops, CPU stress or parallel stress runs on the host'; the CLAUDE.md diff only adds that paragraph right after the unstable-CI rule, which is unchanged. Other build/review wording was trimmed to stay under the token caps; review.md dropped 'leaf is optional' and 'A verifier re-run follows rhiza guide build', outside the acceptance.

## Links

### Out

- verifies [delay-repro-guidance](../capability/rhiza/integrity/self-check/delay-repro-guidance.md)

## History

- `01M6ZCC27HQ7NMBTZ74D1CJFBV` Builders and fixers reproduce flaky tests by generating CPU load on the shared host. On 2026-11-09 two nebari runs did it: 01M6Z1RG ran 8 'bash' busy loops, and 01M6Z5P4 (session 42a0d9da) ran 14 'sh -c while :; do :; done' loops. Each held the 4-core VM at load 16–22 for 15–30+ min and slowed every other live run; I stopped them by PID. Load-based repro also barely works: 45/45 passed under load, while a 500 ms injected delay on the suspected request reproduced 10/10 (nebari PR #90). The build/fix guidance (rhiza guide build, the rhiza-build and rhiza-review skills) should say: reproduce timing flakes by injecting a delay at the suspected race, never by busy loops or parallel stress on the host. Optionally the watchdog flags a run's descendant processes that sit at 100% CPU for minutes.
