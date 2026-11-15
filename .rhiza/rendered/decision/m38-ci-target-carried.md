# m38-ci-target-carried

decision · `01M7GFN17MMTB8MDMXX9AP1N8X`

## reasoning

Only 8 PR CI runs completed since #588 (median 181 s), fewer than the 10 the done point measures, so it can be neither checked nor called missed. The owner's answer to question 4 covered the cold target; the CI median is driven by the same suite cost, so it moves with it rather than holding milestone 38 open or closing it on an unmeasured point.

## statement

The CI median target (3 minutes over the last 10 PR runs) moves into milestone 39 with the cold target, and milestone 38 closes without calling @m38-ci

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [projection-milestone-interview](../source/projection-milestone-interview.md)

### In

- decided-by [milestone-38](../intent/milestone-38.md)
- decided-by [m38-ci](../intent/milestone-38/m38-ci.md)
- decided-by [m39-ci](../intent/milestone-39/m39-ci.md)

## History

- `01M7GFKNYD37Y9T9SAE3Z0JDZX` Close milestone 38, the test suite is fast. Owner, 2026-11-15, after the orchestrator measured a98fd40 (cold make test 423 s wall for 967 CPU-s on a quiet box, warm 41 s, 1,948 tests and 398 txtar files, too few PR CI runs since #588 to call the median): "Close it", with the cold time as an interim result and the cold target moving to the projection milestone.
