# m38-cold-at-a98fd40

evidence · `01M7GFN10445Y7HTCACYTHHDG5` · actor [claude-code-orchestrator](../actor/claude-code-orchestrator.md)

## command

go clean -testcache && /usr/bin/time -v make test

## head

a98fd4051798cc3f48643d530cb70ae53dde2e05

## result

7:03.51 wall, 967 CPU-s, 228% CPU, exit 0; target 4:00 not met

## statement

At a98fd4051, a cold make test (go clean -testcache first) on the 4-core VM with nothing else running took 423.5 s of wall time for 967 CPU-s (569.7 user + 397.2 sys, 228% CPU, exit 0): over @{01M7EQ74ZYGCCHY4QW22JZJ5SH}'s 4 minutes (2026-11-15 baseline: 854 s wall for 972 CPU-s), so milestone 38 cut waiting, not work. Measured by the orchestrator (raw: perf-data/cold-quiet/make-test.log). The target moves to milestone 39.

## status

failing

## Links

### Out

- relates-to [m38-cold](../intent/milestone-38/m38-cold.md) — Evidence that m38-cold is not met at a98fd4051 (interim result)
- derived-from [projection-milestone-interview](../source/projection-milestone-interview.md)

## History

- `01M7GFKNYD37Y9T9SAE3Z0JDZX` Close milestone 38, the test suite is fast. Owner, 2026-11-15, after the orchestrator measured a98fd40 (cold make test 423 s wall for 967 CPU-s on a quiet box, warm 41 s, 1,948 tests and 398 txtar files, too few PR CI runs since #588 to call the median): "Close it", with the cold time as an interim result and the cold target moving to the projection milestone.
