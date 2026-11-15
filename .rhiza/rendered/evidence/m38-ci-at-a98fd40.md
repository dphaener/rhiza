# m38-ci-at-a98fd40

evidence · `01M7GFN12PW139S5GSACMMS9SM` · actor [claude-code-orchestrator](../actor/claude-code-orchestrator.md)

## command

gh run list --workflow CI --event pull_request --status completed --limit 10 --json createdAt,updatedAt

## head

a98fd4051798cc3f48643d530cb70ae53dde2e05

## result

8 runs since #588, median 181 s; too few to call

## statement

Since #588 only 8 PR CI runs completed, fewer than the 10 @m38-ci measures: 165, 73, 325, 197, 228, 83, 212 and 80 s, a median of 181 s (just over 3 minutes). Too few runs to call, so the done point is neither met nor missed; the target moves to milestone 39. Measured by the orchestrator with the done point's gh run list command.

## status

unverified

## Links

### Out

- relates-to [m38-ci](../intent/milestone-38/m38-ci.md) — Evidence that m38-ci could not be called at a98fd4051 (too few runs)
- derived-from [projection-milestone-interview](../source/projection-milestone-interview.md)

## History

- `01M7GFKNYD37Y9T9SAE3Z0JDZX` Close milestone 38, the test suite is fast. Owner, 2026-11-15, after the orchestrator measured a98fd40 (cold make test 423 s wall for 967 CPU-s on a quiet box, warm 41 s, 1,948 tests and 398 txtar files, too few PR CI runs since #588 to call the median): "Close it", with the cold time as an interim result and the cold target moving to the projection milestone.
