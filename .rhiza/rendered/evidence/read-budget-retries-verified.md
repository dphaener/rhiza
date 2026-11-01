# read-budget-retries-verified

evidence · `01M6BENPKD55TH5TNXRJ7K1H3V` · actor [claude-code-62874986](../actor/claude-code-62874986.md)

## result

PASS

## statement

At 042e3a91 ran go test ./internal/api -run 'TestTimeWithinBudget|TestReadBudget' -count=1 -v: both PASS. TestTimeWithinBudget injects a measure func: over-then-under returns ok after 2 calls with 2 log lines; over three times returns not-ok after 3 calls with 3 log lines; under at once stops after 1 call; it also asserts readBudget == 50ms. TestReadBudget logs each read as 'attempt n/3'. The diff leaves const readBudget = 50 * time.Millisecond unchanged and adds readBudgetRetries = 2 with timeWithinBudget. make check exited 0.

## Links

### Out

- verifies [read-budget-retries](../capability/rhiza/views/fast-reads/read-budget-retries.md)

## History

- `01M6B3NDB2E0B1A82CRJANWMTF` The read-budget test retries over-budget reads so load from outside the test run does not fail make check
