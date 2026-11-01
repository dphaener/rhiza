# read-budget-retries

capability · `01M6B3PKC46SS5FN4YP5FKQDHD` · parent [fast-reads](../fast-reads.md)

## acceptance

A Go test in internal/api with an injected timer that is over budget once then under passes, one over budget three times fails, and the budget log records each attempt; the 50 ms budget constant is unchanged

## statement

The read-budget test times a read whose median is over its budget up to two more times and fails only if every attempt is over, logging each attempt, so a short spike of load from outside the test run does not fail make check

## Links

### Out

- decided-by [budget-retries-over-reads](../../../../decision/budget-retries-over-reads.md)

### In

- verifies [read-budget-retries-verified](../../../../evidence/read-budget-retries-verified.md)
- verifies [read-budget-retries-qa-1](../../../../evidence/read-budget-retries-qa-1.md)
- verifies [pr-168-read-budget-retries](../../../../evidence/pr-168-read-budget-retries.md)

## History

- `01M6B3NDB2E0B1A82CRJANWMTF` The read-budget test retries over-budget reads so load from outside the test run does not fail make check
