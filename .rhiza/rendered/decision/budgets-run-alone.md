# budgets-run-alone

decision · `01M6HR4FE2FYXS6YE0CAGNGVR8`

## statement

Run the timing budgets in a second go test phase after the suite instead of beside it; quiettest's wait gave up after 3 minutes while cmd/rhiza ran 5-11, so every run measured under load (PR #210 failed twice then passed on one SHA; #221 and #224 failed on TestReadBudget alone)

## wording

the orchestrator's, as the owner's proxy (standing instruction 'act on my behalf', 2026-11-03); assumed

## Links

### In

- decided-by [budgets-measure-in-quiet](../capability/rhiza/views/fast-reads/budgets-measure-in-quiet.md)

## History

- `01M6HR3VBT4DPCNS03S70QCAZF` Timing budgets measure in quiet: make check runs TestReadBudget, TestBudget and TestGeneratedModel on their own after the rest of the suite, one package at a time, so CI and the VM stop failing them on their neighbours' load
