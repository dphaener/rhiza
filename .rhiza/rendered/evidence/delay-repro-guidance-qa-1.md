# delay-repro-guidance-qa-1

evidence · `01M6ZGASJXTJCJ3SF8KA7TXSRR` · actor [claude-code-2e1281bf](../actor/claude-code-2e1281bf.md)

## result

PASS

## statement

QA at 741b9b74: built the branch (go build -trimpath) into a scratch binary; 'rhiza guide build' and 'rhiza guide review' each print 'Reproduce a timing flake by injecting a delay at the suspected race, never by busy loops or parallel stress on the host.' directly after the unstable-CI reproduce-locally rule; 'rhiza guide install --agent claude' with a scratch HOME wrote skills/rhiza-build/SKILL.md and skills/rhiza-review/SKILL.md, both containing that sentence. CLAUDE.md adds the inject-a-delay / never-busy-loops paragraph right after the unchanged unstable-CI rule. go test -count=1 -trimpath ./internal/guide/... PASS; go test -count=1 -trimpath ./cmd/rhiza -run TestClaudeMDScopedTests PASS.

## Links

### Out

- verifies [delay-repro-guidance](../capability/rhiza/integrity/self-check/delay-repro-guidance.md)

## History

- `01M6ZCC27HQ7NMBTZ74D1CJFBV` Builders and fixers reproduce flaky tests by generating CPU load on the shared host. On 2026-11-09 two nebari runs did it: 01M6Z1RG ran 8 'bash' busy loops, and 01M6Z5P4 (session 42a0d9da) ran 14 'sh -c while :; do :; done' loops. Each held the 4-core VM at load 16–22 for 15–30+ min and slowed every other live run; I stopped them by PID. Load-based repro also barely works: 45/45 passed under load, while a 500 ms injected delay on the suspected request reproduced 10/10 (nebari PR #90). The build/fix guidance (rhiza guide build, the rhiza-build and rhiza-review skills) should say: reproduce timing flakes by injecting a delay at the suspected race, never by busy loops or parallel stress on the host. Optionally the watchdog flags a run's descendant processes that sit at 100% CPU for minutes.
