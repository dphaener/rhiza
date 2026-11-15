# m38-nothing-lost-at-a98fd40

evidence · `01M7GFN0XKEXST5V9ED7PS729A` · actor [claude-code-orchestrator](../actor/claude-code-orchestrator.md)

## command

go test -list '.*' ./... | grep -c '^Test'; ls cmd/rhiza/testdata/script/*.txtar | wc -l

## head

a98fd4051798cc3f48643d530cb70ae53dde2e05

## result

tests 1,892 -> 1,948; scripts 395 -> 398

## statement

At a98fd4051, go test -list '.*' ./... lists 1,948 tests and cmd/rhiza/testdata/script holds 398 txtar files, against 1,892 and 395 at the base 88516d956: no drop, so @m38-nothing-lost holds. Measured by the orchestrator with the done point's commands on both commits.

## status

verified

## Links

### Out

- relates-to [m38-nothing-lost](../intent/milestone-38/m38-nothing-lost.md) — Evidence that m38-nothing-lost holds at a98fd4051
- derived-from [projection-milestone-interview](../source/projection-milestone-interview.md)

## History

- `01M7GFKNYD37Y9T9SAE3Z0JDZX` Close milestone 38, the test suite is fast. Owner, 2026-11-15, after the orchestrator measured a98fd40 (cold make test 423 s wall for 967 CPU-s on a quiet box, warm 41 s, 1,948 tests and 398 txtar files, too few PR CI runs since #588 to call the median): "Close it", with the cold time as an interim result and the cold target moving to the projection milestone.
