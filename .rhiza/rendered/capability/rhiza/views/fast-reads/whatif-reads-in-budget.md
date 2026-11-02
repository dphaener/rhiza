# whatif-reads-in-budget

capability · `01M6C8ES8FHJJXHJND4RZJ8DPW` · parent [fast-reads](../fast-reads.md)

## acceptance

TestReadBudget times every budgetReads read in what-if:<the fixture's draft> under 50 ms median, and the not-budgeted note on BenchmarkShowWhatIf is removed; read_test.go's stored-versus-what-if comparison passes for every read; a property test over generated logs finds the overlay reader's answers equal to the full model's

## statement

Every budgeted read in a what-if view answers from the stored as-is rows overlaid with the change set's records, under the 50 ms median on the 50,000-node model, with the same answers as the full model

## Links

### Out

- decided-by [whatif-overlay-reader](../../../../decision/whatif-overlay-reader.md)

### In

- verifies [whatif-reads-in-budget-verified](../../../../evidence/whatif-reads-in-budget-verified.md)
- verifies [whatif-reads-in-budget-qa-1](../../../../evidence/whatif-reads-in-budget-qa-1.md)
- verifies [whatif-reads-in-budget-verified-2](../../../../evidence/whatif-reads-in-budget-verified-2.md)
- verifies [pr-183-an-incremental-what-if-folds](../../../../evidence/pr-183-an-incremental-what-if-folds.md)

## History

- `01M5SWTFB3ZAX4YKX4KA2ZVSZR` An incremental what-if folds one change set over the indexed as-is, so what-if reads meet the 50 ms target too.
