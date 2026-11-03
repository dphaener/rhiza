# budgets-measure-in-quiet

capability · `01M6HR4C6TN1MN8FF6C71MZ3F2` · parent [fast-reads](../fast-reads.md)

## acceptance

make -n test shows go test ./... with -skip of TestReadBudget|TestBudget|TestGeneratedModel, then those tests alone with -p 1 over internal/api, internal/index, internal/slice and internal/tui/modelapp; the budget phase passes on the loaded VM where it failed beside the suite; CI's make check passes

## statement

The timing budgets run on their own after the rest of the test suite, one package at a time, so they measure reads rather than their neighbours' load

## Links

### Out

- decided-by [budgets-run-alone](../../../../decision/budgets-run-alone.md)

## History

- `01M6HR3VBT4DPCNS03S70QCAZF` Timing budgets measure in quiet: make check runs TestReadBudget, TestBudget and TestGeneratedModel on their own after the rest of the suite, one package at a time, so CI and the VM stop failing them on their neighbours' load
