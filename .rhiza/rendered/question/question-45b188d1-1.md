# question-45b188d1-1

question · `01M6M1PBF4M7FV5WG305ASD8D4`

## session

45b188d1-8ca7-4482-a8d3-01a1a81a2270

## statement

Change set 'Fix Test temp folders leak' serves no intent yet. Which intent should it serve? (The earlier test-suite fix, test-suite-shared, served milestone-25, the backlog intent.) / How should the leaked test temp folders (rhiza-api/index/slice-budget-*, rhiza-quiet-<ppid>) be cleaned up?

## Links

### In

- answers [test-temp-cleanup-by-prefix](../decision/test-temp-cleanup-by-prefix.md)

## History

- `01M6M1M10N7AP6F0DHGANKXEA1` Fix Test temp folders leak and fill the disk: TestMain in internal/api/budget_test.go removes the generated rhiza-api-budget-* model (~146MB) only when the test binary exits normally, so every go test timeout (common under load) leaves one behind; quiettest leaves rhiza-quiet-* the same way, and slice/index budgets their own. On the VM 124 budget folders (~18GB) and 369 quiet folders piled up in TMPDIR by 2026-11-04 and the disk fell to 19GB free. Have each package remove its own stale folders (older than an hour, by prefix) when it starts, or put them under one per-run directory make test deletes, and test it
