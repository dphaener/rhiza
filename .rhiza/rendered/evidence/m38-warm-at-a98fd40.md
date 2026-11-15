# m38-warm-at-a98fd40

evidence · `01M7GFN0VG0ZSCYVTK14CSGJV8` · actor [claude-code-orchestrator](../actor/claude-code-orchestrator.md)

## command

make test, then /usr/bin/time -v make test

## head

a98fd4051798cc3f48643d530cb70ae53dde2e05

## result

0:40.65 wall, 126% CPU, exit 0

## statement

At a98fd4051 (milestone 38 finished, #595), a warm make test on the 4-core VM took 40.65 s of wall time (25.9 user + 25.6 sys CPU-s, 126% CPU, exit 0), within @{01M7EQ750SRBDSEFNF18FJ9HQT}'s 60 s (2026-11-15 baseline: 254 s). Measured by the orchestrator: make test, then /usr/bin/time -v make test (perf-data/after2/time-warm.txt).

## status

verified

## Links

### Out

- relates-to [m38-warm](../intent/milestone-38/m38-warm.md) — Evidence that m38-warm is met at a98fd4051
- derived-from [projection-milestone-interview](../source/projection-milestone-interview.md)

## History

- `01M7GFKNYD37Y9T9SAE3Z0JDZX` Close milestone 38, the test suite is fast. Owner, 2026-11-15, after the orchestrator measured a98fd40 (cold make test 423 s wall for 967 CPU-s on a quiet box, warm 41 s, 1,948 tests and 398 txtar files, too few PR CI runs since #588 to call the median): "Close it", with the cold time as an interim result and the cold target moving to the projection milestone.
