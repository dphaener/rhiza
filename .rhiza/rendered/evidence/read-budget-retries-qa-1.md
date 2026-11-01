# read-budget-retries-qa-1

evidence · `01M6BF4NGEBD7QSWXB8WB5BE7B` · actor [claude-code-dfb4aaf9](../actor/claude-code-dfb4aaf9.md)

## result

PASS

## statement

QA at e8f4c34c, hands-on. (1) go test ./internal/api -run 'TestReadBudget$|TestTimeWithinBudget' -count=1 -v with RHIZA_BUDGET_LOG set: both PASS; the log has one 'attempt 1/3' line per read, medians 1.5-6.7 ms under the 50 ms budget. (2) Budget lowered to 4 ms in a throwaway edit (reverted): roadmap and ready in both views were timed 3 times each, every attempt logged as 'attempt 1/3'..'3/3' to the budget log, and the test failed with 'median over 4ms on every attempt'; reads under budget stopped after attempt 1. (3) Real 50 ms budget with 22 'yes' CPU hogs running outside the test: PASS, medians up to 16 ms. The over-then-under path was not hit with real timings at a borderline budget; TestTimeWithinBudget covers it with an injected timer and asserts readBudget == 50ms. make check with the run env unset exited 0.

## Links

### Out

- verifies [read-budget-retries](../capability/rhiza/views/fast-reads/read-budget-retries.md)

## History

- `01M6B3NDB2E0B1A82CRJANWMTF` The read-budget test retries over-budget reads so load from outside the test run does not fail make check
